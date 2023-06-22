package meminfo

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/DeNetPRO/src/config"
	"github.com/DeNetPRO/src/logger"
	nodeFile "github.com/DeNetPRO/src/node_file"
	nodeTypes "github.com/DeNetPRO/src/node_types"
)

var mutex sync.Mutex

// RestoreMemoryInfo sets previous value of used storage space info.
func Restore(pathToConfig string, fileSize int) {
	location := "files.restoreMemoryInfo->"

	mutex.Lock()
	defer mutex.Unlock()

	confFile, fileBytes, err := nodeFile.Read(pathToConfig)
	if err != nil {
		logger.Log(logger.MarkLocation(location, err))
		return
	}
	defer confFile.Close()

	var nodeConfig nodeTypes.Config

	err = json.Unmarshal(fileBytes, &nodeConfig)
	if err != nil {
		logger.Log(logger.MarkLocation(location, err))
		return
	}

	nodeConfig.UsedStorageSpace -= int64(fileSize)

	if nodeConfig.UsedStorageSpace < 0 {
		logger.Log(logger.MarkLocation(location, errors.New("used storage space is less than 0")))
		return
	}

	err = config.Save(confFile, nodeConfig)
	if err != nil {
		logger.Log(logger.MarkLocation(location, err))
		return
	}
}
