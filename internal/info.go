package internal

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
)

func InfoSend(spAddress string, nodeIP, fileKeys []string, fileSize int) (*NodesResponse, error) {
	const logLoc = "internal.InfoSend->"

	addressPath := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spAddress)

	_, err := os.Stat(addressPath)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	shared.MU.Lock()
	spFsFile, fileBytes, err := shared.ReadFile(filepath.Join(addressPath, paths.SpFsFilename))
	if err != nil {
		shared.MU.Unlock()
		return nil, logger.CreateDetails(logLoc, err)
	}

	defer spFsFile.Close()

	spFs := &shared.StorageProviderData{}

	err = json.Unmarshal(fileBytes, spFs)
	if err != nil {
		shared.MU.Unlock()
		return nil, logger.CreateDetails(logLoc, err)
	}
	spFs.Address = spAddress

	shared.MU.Unlock()
	spFsFile.Close()

	succesfulNodes := make([]string, 0)
	for _, node := range nodeIP {
		backUpNodes, err := BackUpFileKeys(node, addressPath, spFs, fileSize, fileKeys)
		if err != nil {
			logger.Log(logger.CreateDetails(logLoc, err))
		}

		succesfulNodes = append(succesfulNodes, backUpNodes...)
	}

	if len(succesfulNodes) == 0 {
		return nil, logger.CreateDetails(logLoc, errors.New("no nodes"))
	}

	return &NodesResponse{
		Nodes: succesfulNodes,
	}, nil
}
