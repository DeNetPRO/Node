package files

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/minio/sha256-simd"

	blockchainprovider "git.denetwork.xyz/dfile/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/dfile/dfile-secondary-node/config"
	dnetsignature "git.denetwork.xyz/dfile/dfile-secondary-node/dnet_signature"
	fsysinfo "git.denetwork.xyz/dfile/dfile-secondary-node/fsys_info"

	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

type NodesResponse struct {
	Nodes []string `json:"nodes"`
}

type NodeAddressResponse struct {
	NodeAddress string `json:"node_address"`
}

const eightKB = 8192

//Copy makes a copy of file parts that are stored on the node, or sends them on oher node for replication.
func Copy(req *http.Request, spData *shared.StorageProviderData, config *config.SecondaryNodeConfig, pathToConfig string, fileSize int, enoughSpace bool) (*NodeAddressResponse, error) {
	const location = "files.Copy->"

	pathToSpFiles := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spData.Address)

	hashes := req.MultipartForm.File["hashes"]

	if len(hashes) == 0 {
		return nil, logger.CreateDetails(location, errors.New("empty hashes"))
	}

	hashesFileHeader, err := hashes[0].Open()
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	hashesInfo, err := io.ReadAll(hashesFileHeader)
	if err != nil {
		hashesFileHeader.Close()
		return nil, logger.CreateDetails(location, err)
	}

	hashesMap := make(map[string]string)
	err = json.Unmarshal(hashesInfo, &hashesMap)
	if err != nil {
		hashesFileHeader.Close()
		return nil, logger.CreateDetails(location, err)
	}

	hashesFileHeader.Close()

	if !enoughSpace {
		nftNode, err := blockchainprovider.GetNodeNFT()
		if err != nil {
			return nil, logger.CreateDetails(location, err)
		}

		totalNodes, err := nftNode.TotalSupply(&bind.CallOpts{})
		if err != nil {
			return nil, logger.CreateDetails(location, err)
		}

		intTotal := totalNodes.Int64()
		fastReq := fasthttp.AcquireRequest()
		fastResp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(fastReq)
		defer fasthttp.ReleaseResponse(fastResp)

		rand.Seed(time.Now().UnixNano())

		for i := int64(0); i < intTotal; i++ {
			randID := rand.Int63n(intTotal)
			nodeInfo, err := nftNode.GetNodeById(&bind.CallOpts{}, big.NewInt(randID))
			if err != nil {
				continue
			}

			ipBuilder := strings.Builder{}
			for i, v := range nodeInfo.IpAddress {
				stringPart := strconv.Itoa(int(v))
				ipBuilder.WriteString(stringPart)

				if i < 3 {
					ipBuilder.WriteString(".")
				}
			}

			stringPort := strconv.Itoa(int(nodeInfo.Port))
			ipBuilder.WriteString(":")
			ipBuilder.WriteString(stringPort)

			nodeIP := ipBuilder.String()

			if nodeIP == config.IpAddress+":"+config.HTTPPort {
				continue
			}

			url := "http://" + nodeIP
			fastReq.Reset()
			fastResp.Reset()

			fastReq.Header.SetRequestURI(url)
			fastReq.Header.SetMethod("GET")
			fastReq.Header.Set("Connection", "close")

			err = fasthttp.Do(fastReq, fastResp)
			if err != nil {
				continue
			}

			err = copyOnOtherNode(nodeIP, pathToSpFiles, req.MultipartForm, hashesMap, fileSize)
			if err != nil {
				continue
			}

			return &NodeAddressResponse{
				NodeAddress: nodeIP,
			}, nil
		}

		return nil, logger.CreateDetails(location, errors.New("no available nodes"))
	}

	err = fsysinfo.Save(pathToSpFiles, spData)
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	savedParts := make([]string, 0, len(hashesMap))
	for oldHash, newHash := range hashesMap {
		savedParts = append(savedParts, newHash)

		path := filepath.Join(pathToSpFiles, oldHash)
		file, err := os.Open(path)
		if err != nil {
			file.Close()
			deleteParts(pathToSpFiles, savedParts)
			return nil, logger.CreateDetails(location, err)
		}

		err = savePart(file, pathToSpFiles, newHash)
		if err != nil {
			file.Close()
			deleteParts(pathToSpFiles, savedParts)
			return nil, logger.CreateDetails(location, err)
		}

		file.Close()
	}

	return &NodeAddressResponse{
		NodeAddress: config.IpAddress + ":" + config.HTTPPort,
	}, nil
}

// ====================================================================================
//Save is used for chaecking and saving file parts from the inoming request to the node's storage.
func Save(req *http.Request, spData *shared.StorageProviderData, pathToConfig string, fileSize int) error {
	const location = "files.Save->"

	pathToSpFiles := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spData.Address)

	err := fsysinfo.Save(pathToSpFiles, spData)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	reqFileParts := req.MultipartForm.File["files"]

	oneMBHashes, err := shared.GetOneMbHashes(reqFileParts)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	var wholeFileHash string

	if len(oneMBHashes) == 1 {
		wholeFileHash = oneMBHashes[0]
	} else {
		sort.Strings(oneMBHashes)
		wholeFileHash, _, err = shared.CalcRootHash(oneMBHashes)
		if err != nil {
			return logger.CreateDetails(location, err)
		}
	}

	fsContainsFile := false

	for _, fileHash := range spData.Fs {
		if fileHash == wholeFileHash {
			fsContainsFile = true
			break
		}
	}

	if !fsContainsFile {
		return logger.CreateDetails(location, shared.ErrWrongFile)
	}

	count := 0
	savedParts := make([]string, 0, len(oneMBHashes))
	for _, reqFilePart := range reqFileParts {
		savedParts = append(savedParts, reqFilePart.Filename)

		rqFile, err := reqFilePart.Open()
		if err != nil {
			return logger.CreateDetails(location, err)
		}

		err = savePart(rqFile, pathToSpFiles, reqFilePart.Filename)
		if err != nil {
			rqFile.Close()
			deleteParts(pathToSpFiles, savedParts)
			return logger.CreateDetails(location, err)
		}

		rqFile.Close()

		count++
		logger.Log("Saved file " + reqFilePart.Filename + " (" + strconv.Itoa(count) + "/" + strconv.Itoa(len(oneMBHashes)) + ")" + " from " + spData.Address) //TODO remove
	}

	return nil
}

// ====================================================================================
func savePart(file io.Reader, pathToSpFiles, fileName string) error {
	const location = "files.saveFilePart->"

	pathToFile := filepath.Join(pathToSpFiles, fileName)

	newFile, err := os.Create(pathToFile)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	_, err = io.Copy(newFile, file)
	if err != nil {
		newFile.Close()
		return logger.CreateDetails(location, err)
	}

	newFile.Sync()
	newFile.Close()

	return nil
}

// ====================================================================================

func Serve(spAddress, fileKey, signatureFromReq string) (string, error) {
	const location = "files.Serve->"

	signature, err := hex.DecodeString(signatureFromReq)
	if err != nil {
		return "", logger.CreateDetails(location, err)
	}

	hash := sha256.Sum256([]byte(fileKey + spAddress))

	err = dnetsignature.Check(spAddress, signature, hash)
	if err != nil {
		return "", logger.CreateDetails(location, err)
	}

	pathToFile := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spAddress, fileKey)

	_, err = os.Stat(pathToFile)
	if err != nil {
		return "", logger.CreateDetails(location, err)
	}

	return pathToFile, nil
}

// ====================================================================================

//Sends files to other node if they can't be copied locally and returns that nodes address
func copyOnOtherNode(nodeAddress, pathToSpFiles string, multiForm *multipart.Form, hashesMap map[string]string, fileSize int) error {
	const location = "files_helpers.backUp->"

	pipeConns := fasthttputil.NewPipeConns()
	pr := pipeConns.Conn1()
	pw := pipeConns.Conn2()

	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		defer pw.Close()

		address := multiForm.Value["address"]
		nonce := multiForm.Value["nonce"]
		fsRootHash := multiForm.Value["fsRootHash"]

		err := writer.WriteField("address", address[0])
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			return
		}

		err = writer.WriteField("nonce", nonce[0])
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			return
		}

		err = writer.WriteField("fsRootHash", fsRootHash[0])
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			return
		}

		wholeFileHashes := multiForm.Value["fs"]
		for _, wholeHash := range wholeFileHashes {
			err = writer.WriteField("fs", wholeHash)
			if err != nil {
				logger.Log(logger.CreateDetails(location, err))
				return
			}
		}

		for oldHash, newHash := range hashesMap {
			path := filepath.Join(pathToSpFiles, oldHash)
			file, err := os.Open(path)
			if err != nil {
				logger.Log(logger.CreateDetails(location, err))
				return
			}

			filePart, err := writer.CreateFormFile("files", newHash)
			if err != nil {
				file.Close()
				logger.Log(logger.CreateDetails(location, err))
				return
			}

			_, err = io.Copy(filePart, file)
			if err != nil {
				file.Close()
				logger.Log(logger.CreateDetails(location, err))
				return
			}

			file.Close()
		}
	}()

	req, err := http.NewRequest("POST", "http://"+nodeAddress+"/upload/"+strconv.Itoa(fileSize), pr)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	defer resp.Body.Close()

	if string(body) != "OK" {
		return logger.CreateDetails(location, shared.ErrFileSaving)
	}

	return nil
}

// ====================================================================================

// DeleteParts deletes parts of the file that wasn't fully uploaded to the node for some reason.
func deleteParts(addressPath string, fileHashes []string) {
	logger.Log("deleting file parts after error...")

	for _, hash := range fileHashes {
		pathToFile := filepath.Join(addressPath, hash)

		os.Remove(pathToFile)
	}
}
