package internal

import (
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"

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

func CopyFile(req *http.Request, spData *shared.StorageProviderData, config *config.SecondaryNodeConfig, pathToConfig string, fileSize int, enoughSpace bool) (*NodeAddressResponse, error) {
	const logLoc = "internal.CopyFile->"

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
