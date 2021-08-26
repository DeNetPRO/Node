package files

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	blockchainprovider "git.denetwork.xyz/dfile/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/dfile/dfile-secondary-node/config"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/valyala/fasthttp"
)

type NodeAddressResponse struct {
	NodeAddress string `json:"node_address"`
}

const eightKB = 8192

func Copy(req *http.Request, spData *shared.StorageProviderData, config *config.SecondaryNodeConfig, pathToConfig string, fileSize int, enoughSpace bool) (*NodeAddressResponse, error) {
	const logLoc = "files.Copy->"

	addressPath := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spData.Address)

	if !enoughSpace {
		nftNode, err := blockchainprovider.GetNodeNFT()
		if err != nil {
			RestoreMemoryInfo(pathToConfig, fileSize)
			return nil, logger.CreateDetails(logLoc, err)
		}

		total, err := nftNode.TotalSupply(&bind.CallOpts{})
		if err != nil {
			if err != nil {
				RestoreMemoryInfo(pathToConfig, fileSize)
				return nil, logger.CreateDetails(logLoc, err)
			}

			intTotal := total.Int64()

			fastReq := fasthttp.AcquireRequest()
			fastResp := fasthttp.AcquireResponse()
			defer fasthttp.ReleaseRequest(fastReq)
			defer fasthttp.ReleaseResponse(fastResp)

			for i := int64(0); i < intTotal; i++ {
				node, err := nftNode.GetNodeById(&bind.CallOpts{}, big.NewInt(i))
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

				nodeAddress, err := backUpCopy(nodeIP, addressPath, req.MultipartForm, fileSize)
				if err != nil {
					continue
				}

				return &NodeAddressResponse{
					NodeAddress: nodeAddress,
				}, nil
			}

			return nil, logger.CreateDetails(logLoc, errors.New("no available nodes"))
		}
	}

	err := initSPFile(addressPath, spData)
	if err != nil {
		RestoreMemoryInfo(pathToConfig, fileSize)
		return nil, logger.CreateDetails(logLoc, err)
	}

	hashes := req.MultipartForm.File["hashes"]
	hashesFile, err := hashes[0].Open()
	if err != nil {
		RestoreMemoryInfo(pathToConfig, fileSize)
		return nil, logger.CreateDetails(logLoc, err)
	}

	hashesBody, err := io.ReadAll(hashesFile)
	if err != nil {
		hashesFile.Close()
		RestoreMemoryInfo(pathToConfig, fileSize)
		return nil, logger.CreateDetails(logLoc, err)
	}

	hashDif := make(map[string]string)
	err = json.Unmarshal(hashesBody, &hashDif)
	if err != nil {
		hashesFile.Close()
		RestoreMemoryInfo(pathToConfig, fileSize)
		return nil, logger.CreateDetails(logLoc, err)
	}

	hashesFile.Close()

	for old, new := range hashDif {
		path := filepath.Join(addressPath, old)
		file, err := os.Open(path)
		if err != nil {
			logger.Log(logger.CreateDetails(logLoc, err))
			RestoreMemoryInfo(pathToConfig, fileSize)
			return nil, logger.CreateDetails(logLoc, err)
		}

		defer file.Close()

		newPath := filepath.Join(addressPath, new)
		newFile, err := os.Create(newPath)
		if err != nil {
			logger.Log(logger.CreateDetails(logLoc, err))
			RestoreMemoryInfo(pathToConfig, fileSize)
			return nil, logger.CreateDetails(logLoc, err)
		}

		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			logger.Log(logger.CreateDetails(logLoc, err))
			RestoreMemoryInfo(pathToConfig, fileSize)
			return nil, logger.CreateDetails(logLoc, err)
		}

		newFile.Sync()
		newFile.Close()
	}

	return &NodeAddressResponse{
		NodeAddress: config.IpAddress + ":" + config.HTTPPort,
	}, nil
}

// ====================================================================================

func Save(req *http.Request, spData *shared.StorageProviderData, pathToConfig string, fileSize int) error {
	const logLoc = "files.Save->"

	pathToSpFiles := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spData.Address)

	err := initSPFile(pathToSpFiles, spData)
	if err != nil {
		RestoreMemoryInfo(pathToConfig, fileSize)
		return logger.CreateDetails(logLoc, err)
	}

	reqFileParts := req.MultipartForm.File["files"]

	oneMBHashes, err := GetOneMbHashes(reqFileParts)
	if err != nil {
		RestoreMemoryInfo(pathToConfig, fileSize)
		return logger.CreateDetails(logLoc, err)
	}

	var wholeFileHash string

	if len(oneMBHashes) == 1 {
		wholeFileHash = oneMBHashes[0]
	} else {
		sort.Strings(oneMBHashes)
		wholeFileHash, _, err = shared.CalcRootHash(oneMBHashes)
		if err != nil {
			RestoreMemoryInfo(pathToConfig, fileSize)
			return logger.CreateDetails(logLoc, err)
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
		RestoreMemoryInfo(pathToConfig, fileSize)
		return logger.CreateDetails(logLoc, shared.ErrWrongFile)
	}

	err = SaveFileParts(reqFileParts, spData.Address, pathToSpFiles, len(oneMBHashes))
	if err != nil {
		RestoreMemoryInfo(pathToConfig, fileSize)
		DeleteParts(pathToSpFiles, oneMBHashes)
		return logger.CreateDetails(logLoc, err)
	}

	return nil
}

func SaveFileParts(reqFileParts []*multipart.FileHeader, spAddress, pathToSpFiles string, total int) error {
	count := 1

	for _, reqFilePart := range reqFileParts {
		rqFile, err := reqFilePart.Open()
		if err != nil {
			return err
		}
		defer rqFile.Close()

		pathToFile := filepath.Join(pathToSpFiles, reqFilePart.Filename)

		newFile, err := os.Create(pathToFile)
		if err != nil {
			return err
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, rqFile)
		if err != nil {
			return err
		}

		logger.Log("Saved file " + reqFilePart.Filename + " (" + fmt.Sprint(count) + "/" + fmt.Sprint(total) + ")" + " from " + spAddress) //TODO remove

		newFile.Sync()
		rqFile.Close()
		newFile.Close()

		count++
	}

	return nil
}

func GetOneMbHashes(reqFileParts []*multipart.FileHeader) ([]string, error) {

	eightKBHashes := make([]string, 0, 128)
	oneMBHashes := make([]string, 0, len(reqFileParts))

	for _, reqFilePart := range reqFileParts {

		var buf bytes.Buffer

		rqFile, err := reqFilePart.Open()
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(&buf, rqFile)
		if err != nil {
			rqFile.Close()
			return nil, err
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
			return nil, err
		}

		if reqFilePart.Filename != oneMBHash {
			return nil, err
		}

		oneMBHashes = append(oneMBHashes, oneMBHash)
	}

	return oneMBHashes, nil
}

// ====================================================================================

func Serve(spAddress, fileKey, signatureFromReq string) (string, error) {
	const logLoc = "files.Serve->"

	signature, err := hex.DecodeString(signatureFromReq)
	if err != nil {
		return "", logger.CreateDetails(logLoc, err)
	}

	hash := sha256.Sum256([]byte(fileKey + spAddress))

	sigPublicKey, err := crypto.SigToPub(hash[:], signature)
	if err != nil {
		return "", logger.CreateDetails(logLoc, err)
	}

	signatureAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if spAddress != signatureAddress.String() {
		return "", logger.CreateDetails(logLoc, shared.ErrWrongSignature)
	}

	pathToFile := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spAddress, fileKey)

	_, err = os.Stat(pathToFile)
	if err != nil {
		return "", logger.CreateDetails(logLoc, err)
	}

	return pathToFile, nil
}

// ====================================================================================

func BackUp(req *http.Request, spData *shared.StorageProviderData, pathToConfig string, fileSize int) error {
	const logLoc = "files.BackUp->"

	pathToSpFiles := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spData.Address)

	err := initSPFile(pathToSpFiles, spData)
	if err != nil {
		RestoreMemoryInfo(pathToConfig, fileSize)
		return logger.CreateDetails(logLoc, err)
	}

	reqFileParts := req.MultipartForm.File["files"]

	oneMBHashes, err := GetOneMbHashes(reqFileParts)
	if err != nil {
		RestoreMemoryInfo(pathToConfig, fileSize)
		return logger.CreateDetails(logLoc, err)
	}

	fsContainsFile := false

	var wholeFileHash string

	if len(oneMBHashes) == 1 {
		wholeFileHash = oneMBHashes[0]
	} else {
		sort.Strings(oneMBHashes)
		wholeFileHash, _, err = shared.CalcRootHash(oneMBHashes)
		if err != nil {
			RestoreMemoryInfo(pathToConfig, fileSize)
			return logger.CreateDetails(logLoc, err)
		}
	}

	for _, fileHash := range spData.Fs {
		if fileHash == wholeFileHash {
			fsContainsFile = true
		}
	}

	if !fsContainsFile {
		RestoreMemoryInfo(pathToConfig, fileSize)
		return logger.CreateDetails(logLoc, shared.ErrWrongFile)
	}

	hashes := req.MultipartForm.File["hashes"]
	hashesFile, err := hashes[0].Open()
	if err != nil {
		RestoreMemoryInfo(pathToConfig, fileSize)
		return logger.CreateDetails(logLoc, err)
	}

	hashesBody, err := io.ReadAll(hashesFile)
	if err != nil {
		hashesFile.Close()
		RestoreMemoryInfo(pathToConfig, fileSize)
		return logger.CreateDetails(logLoc, err)
	}

	hashDif := make(map[string]string)
	err = json.Unmarshal(hashesBody, &hashDif)
	if err != nil {
		hashesFile.Close()
		RestoreMemoryInfo(pathToConfig, fileSize)
		return logger.CreateDetails(logLoc, err)
	}

	hashesFile.Close()

	err = SaveFileParts(reqFileParts, spData.Address, pathToSpFiles, len(oneMBHashes))
	if err != nil {
		RestoreMemoryInfo(pathToConfig, fileSize)
		DeleteParts(pathToSpFiles, oneMBHashes)
		return logger.CreateDetails(logLoc, err)
	}

	return nil
}
