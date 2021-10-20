package meminfo

import (
	"encoding/json"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/config"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	nodeFile "git.denetwork.xyz/DeNet/dfile-secondary-node/node_file"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/shared"
)

//RestoreMemoryInfo sets previous value of used storage space info.
func Restore(pathToConfig string, fileSize int) {
	location := "files.restoreMemoryInfo->"

	shared.MU.Lock()
	defer shared.MU.Unlock()

	confFile, fileBytes, err := nodeFile.Read(pathToConfig)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		return
	}
	defer confFile.Close()

	var nodeConfig config.NodeConfig

	err = json.Unmarshal(fileBytes, &nodeConfig)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		return
	}

	nodeConfig.UsedStorageSpace -= int64(fileSize)

	err = config.Save(confFile, nodeConfig)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		return
	}
}
