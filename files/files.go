package files

import (
	"bytes"
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
	"time"

	"github.com/minio/sha256-simd"

	blockchainprovider "git.denetwork.xyz/dfile/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/dfile/dfile-secondary-node/config"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/valyala/fasthttp"
)

type NodeAddressResponse struct {
	NodeAddress string `json:"node_address"`
}

const eightKB = 8192

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
			node, err := nftNode.GetNodeById(&bind.CallOpts{}, big.NewInt(randID))
			if err != nil {
				continue
			}

			nodeIP := getNodeIP(node)

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

			nodeAddress, err := backUp(nodeIP, pathToSpFiles, req.MultipartForm, hashesMap, fileSize)
			if err != nil {
				continue
			}

			return &NodeAddressResponse{
				NodeAddress: nodeAddress,
			}, nil
		}

		return nil, logger.CreateDetails(location, errors.New("no available nodes"))
	}

	err = saveSpFsInfo(pathToSpFiles, spData)
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
			DeleteParts(pathToSpFiles, savedParts)
			return nil, logger.CreateDetails(location, err)
		}

		err = saveFilePart(file, pathToSpFiles, newHash)
		if err != nil {
			file.Close()
			DeleteParts(pathToSpFiles, savedParts)
			return nil, logger.CreateDetails(location, err)
		}

		file.Close()
	}

	return &NodeAddressResponse{
		NodeAddress: config.IpAddress + ":" + config.HTTPPort,
	}, nil
}

// ====================================================================================

func Save(req *http.Request, spData *shared.StorageProviderData, pathToConfig string, fileSize int) error {
	const location = "files.Save->"

	pathToSpFiles := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spData.Address)

	err := saveSpFsInfo(pathToSpFiles, spData)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	reqFileParts := req.MultipartForm.File["files"]

	oneMBHashes, err := GetOneMbHashes(reqFileParts)
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
		defer rqFile.Close()

		err = saveFilePart(rqFile, pathToSpFiles, reqFilePart.Filename)
		if err != nil {
			DeleteParts(pathToSpFiles, savedParts)
			return logger.CreateDetails(location, err)
		}

		count++
		logger.Log("Saved file " + reqFilePart.Filename + " (" + strconv.Itoa(count) + "/" + strconv.Itoa(len(oneMBHashes)) + ")" + " from " + spData.Address) //TODO remove
	}

	return nil
}

// ====================================================================================

func saveFilePart(file io.Reader, pathToSpFiles, fileName string) error {
	const location = "files.saveFilePart->"

	pathToFile := filepath.Join(pathToSpFiles, fileName)

	newFile, err := os.Create(pathToFile)
	if err != nil {
		return logger.CreateDetails(location, err)
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	newFile.Sync()
	newFile.Close()

	return nil
}

// ====================================================================================

func GetOneMbHashes(reqFileParts []*multipart.FileHeader) ([]string, error) {
	const location = "files.GetOneMbHashes->"
	eightKBHashes := make([]string, 0, 128)
	oneMBHashes := make([]string, 0, len(reqFileParts))

	for _, reqFilePart := range reqFileParts {

		var buf bytes.Buffer

		rqFile, err := reqFilePart.Open()
		if err != nil {
			return nil, logger.CreateDetails(location, err)
		}

		_, err = io.Copy(&buf, rqFile)
		if err != nil {
			rqFile.Close()
			return nil, logger.CreateDetails(location, err)
		}

		rqFile.Close()

		bufBytes := buf.Bytes()
		eightKBHashes = eightKBHashes[:0]

		for i := 0; i < len(bufBytes); i += eightKB {
			hSum := sha256.Sum256(bufBytes[i : i+eightKB])
			eightKBHashes = append(eightKBHashes, hex.EncodeToString(hSum[:]))
		}

		oneMBHash, _, err := shared.CalcRootHash(eightKBHashes)
		if err != nil {
			return nil, logger.CreateDetails(location, err)
		}

		if reqFilePart.Filename != oneMBHash {
			return nil, logger.CreateDetails(location, err)
		}

		oneMBHashes = append(oneMBHashes, oneMBHash)
	}

	return oneMBHashes, nil
}

// ====================================================================================

func Serve(spAddress, fileKey, signatureFromReq string) (string, error) {
	const location = "files.Serve->"

	signature, err := hex.DecodeString(signatureFromReq)
	if err != nil {
		return "", logger.CreateDetails(location, err)
	}

	hash := sha256.Sum256([]byte(fileKey + spAddress))

	err = CheckDataSign(spAddress, signature, hash)
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
