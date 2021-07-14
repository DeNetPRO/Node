package cleaner

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"git.denetwork.xyz/dfile/dfile-secondary-node/config"
	"git.denetwork.xyz/dfile/dfile-secondary-node/encryption"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
)

const oneMB = 1048576

func Start() {

	const logInfo = "cleaner.Start->"

	regAddr := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	regFileName := regexp.MustCompile("[0-9A-Za-z_]")

	nodeAddr, err := encryption.DecryptNodeAddr()
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
	}

	for {
		time.Sleep(time.Minute) // add period

		pathToAccStorage := filepath.Join(paths.AccsDirPath, nodeAddr.String(), paths.StorageDirName)

		storageProviderAddresses := []string{}

		err = filepath.WalkDir(pathToAccStorage,
			func(path string, info fs.DirEntry, err error) error {
				if err != nil {
					logger.Log(logger.CreateDetails(logInfo, err))
				}

				if regAddr.MatchString(info.Name()) {
					storageProviderAddresses = append(storageProviderAddresses, info.Name())
				}

				return nil
			})

		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			continue
		}

		if len(storageProviderAddresses) == 0 {
			continue
		}

		removedTotal := 0

		for _, spAddress := range storageProviderAddresses {

			fileNames := []string{}

			pathToStorProviderFiles := filepath.Join(pathToAccStorage, spAddress)

			err = filepath.WalkDir(pathToStorProviderFiles,
				func(path string, info fs.DirEntry, err error) error {
					if err != nil {
						logger.Log(logger.CreateDetails(logInfo, err))
					}

					if regFileName.MatchString(info.Name()) && len(info.Name()) == 64 {
						fileNames = append(fileNames, info.Name())
					}

					return nil
				})
			if err != nil {
				logger.Log(logger.CreateDetails(logInfo, err))
				continue
			}

			pathToFsTree := filepath.Join(paths.AccsDirPath, nodeAddr.String(), paths.StorageDirName, spAddress, "tree.json")

			shared.MU.Lock()
			fileFsTree, err := os.Open(pathToFsTree)
			if err != nil {
				shared.MU.Unlock()
				logger.Log(logger.CreateDetails(logInfo, err))
			}

			treeBytes, err := io.ReadAll(fileFsTree)
			if err != nil {
				fileFsTree.Close()
				shared.MU.Unlock()
				logger.Log(logger.CreateDetails(logInfo, err))
			}
			fileFsTree.Close()
			shared.MU.Unlock()

			var storageFsStruct shared.StorageInfo

			err = json.Unmarshal(treeBytes, &storageFsStruct)
			if err != nil {
				logger.Log(logger.CreateDetails(logInfo, err))
			}

			fsFiles := map[string]bool{}

			for _, hashes := range storageFsStruct.Tree {
				for _, hash := range hashes {
					fsFiles[hex.EncodeToString(hash)] = true
				}
			}

			for _, fileName := range fileNames {

				if !fsFiles[fileName] {
					shared.MU.Lock()
					logger.Log("removing file: " + fileName + " of " + spAddress)
					err = os.Remove(filepath.Join(pathToStorProviderFiles, fileName))
					if err != nil {
						logger.Log(logger.CreateDetails(logInfo, err))
					}

					removedTotal++

					shared.MU.Unlock()
				}
			}

		}

		if removedTotal > 0 {
			pathToConfig := filepath.Join(paths.AccsDirPath, nodeAddr.String(), paths.ConfDirName, "config.json")

			shared.MU.Lock()
			confFile, err := os.OpenFile(pathToConfig, os.O_RDWR, 0755)
			if err != nil {
				shared.MU.Unlock()
				logger.Log(logger.CreateDetails(logInfo, err))
			}
			defer confFile.Close()

			fileBytes, err := io.ReadAll(confFile)
			if err != nil {
				shared.MU.Unlock()
				logger.Log(logger.CreateDetails(logInfo, err))
			}

			var dFileConf config.SecondaryNodeConfig

			err = json.Unmarshal(fileBytes, &dFileConf)
			if err != nil {
				shared.MU.Unlock()
				logger.Log(logger.CreateDetails(logInfo, err))
			}

			dFileConf.UsedStorageSpace -= int64(removedTotal * oneMB)

			err = config.Save(confFile, dFileConf)
			if err != nil {
				shared.MU.Unlock()
				logger.Log(logger.CreateDetails(logInfo, err))
			}
			confFile.Close()
			shared.MU.Unlock()
		}

	}

}
