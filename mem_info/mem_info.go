package meminfo

import (
	"encoding/json"

	"git.denetwork.xyz/dfile/dfile-secondary-node/config"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	nodeFile "git.denetwork.xyz/dfile/dfile-secondary-node/node_file"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
)

//RestoreMemoryInfo sets previous value of used storage space info.
func Restore(pathToConfig string, intFileSize int) {
	location := "files.restoreMemoryInfo->"

	shared.MU.Lock()
	defer shared.MU.Unlock()

	confFile, fileBytes, err := nodeFile.Read(pathToConfig)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		return
	}
	defer confFile.Close()

	var nodeConfig config.SecondaryNodeConfig

	err = json.Unmarshal(fileBytes, &nodeConfig)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		return
	}

	nodeConfig.UsedStorageSpace -= int64(intFileSize)

	err = config.Save(confFile, nodeConfig)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		return
	}
}
