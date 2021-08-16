package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"git.denetwork.xyz/dfile/dfile-secondary-node/config"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"github.com/valyala/fasthttp"
)

func CheckFile(spAddress, fileKey string) (*NodeAddressResponse, error) {
	const logLoc = "internal.CheckFile"

	filePath := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spAddress, fileKey)
	file, err := os.Stat(filePath)
	err = shared.CheckStatErr(err)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	if file == nil {
		rating, err := shared.GetRating()
		if err != nil {
			return nil, logger.CreateDetails(logLoc, err)
		}

		if rating < 40 {
			return nil, logger.CreateDetails(logLoc, errors.New("no file"))
		}

		connectedNode, err := shared.GetConnectionNodes()
		if err != nil {
			return nil, logger.CreateDetails(logLoc, err)
		}

		fastReq := fasthttp.AcquireRequest()
		fastResp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(fastReq)
		defer fasthttp.ReleaseResponse(fastResp)

		var nodeAddress *NodeAddressResponse

		for _, node := range connectedNode {
			fastReq.Reset()
			fastResp.Reset()

			nodeAddress, err = checkFileOnNode(node, spAddress, fileKey, fastReq, fastResp)
			if err != nil {
				err = shared.RemoveConnectionNode(node)
				if err != nil {
					logger.Log(logger.CreateDetails(logLoc, err))
				}
				continue
			}

			return nodeAddress, nil
		}
	}

	pathToConfig := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.ConfDirName, paths.ConfFileName)

	var nodeConfig config.SecondaryNodeConfig

	shared.MU.Lock()
	confFile, fileBytes, err := shared.ReadFile(pathToConfig)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	confFile.Close()

	err = json.Unmarshal(fileBytes, &nodeConfig)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}
	shared.MU.Unlock()

	return &NodeAddressResponse{
		NodeAddress: nodeConfig.IpAddress + ":" + nodeConfig.HTTPPort,
	}, nil
}

func checkFileOnNode(node, spAddress, fileKey string, req *fasthttp.Request, resp *fasthttp.Response) (*NodeAddressResponse, error) {
	const logLoc = "internal.checkFileOnNode"

	req.SetRequestURI("http://" + node + "/check/" + spAddress + "/" + fileKey)
	req.Header.SetMethod("GET")

	err := fasthttp.Do(req, resp)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	if string(resp.Body()) == "no file" {
		fmt.Println(string(resp.Body()))
		return nil, logger.CreateDetails(logLoc, err)
	}

	nodeAddress := &NodeAddressResponse{}

	err = json.Unmarshal(resp.Body(), nodeAddress)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	return nodeAddress, nil
}
