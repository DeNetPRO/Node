package cleaner

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	blckChain "git.denetwork.xyz/DeNet/dfile-secondary-node/blockchain_provider"
	nodeFile "git.denetwork.xyz/DeNet/dfile-secondary-node/node_file"
	tstpkg "git.denetwork.xyz/DeNet/dfile-secondary-node/tst_pkg"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/config"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/shared"
)

const oneMB = 1048576

//Starts cleaner, that checks if stored file part is in Storage Provider's file system and deletes it if it was not found.
func Start() {
	const location = "cleaner.Start->"

	regAddr := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	regFileName := regexp.MustCompile("[0-9A-Za-z_]")

	for {
		time.Sleep(time.Minute * 1)

		for network := range blckChain.Networks {
			pathToAccStorage := filepath.Join(paths.StoragePaths[0], network)

			stat, err := os.Stat(pathToAccStorage)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				logger.Log(logger.CreateDetails(location, err))
				log.Fatal(err)
			}

			if stat == nil {
				fmt.Println("no files from", network, "users to delete")
				continue
			}

			dirFiles, err := nodeFile.ReadDirFiles(pathToAccStorage)
			if err != nil {
				logger.Log(logger.CreateDetails(location, err))
				continue
			}

			storageProviderAddresses := []string{}

			for _, f := range dirFiles {
				if regAddr.MatchString(f.Name()) {
					storageProviderAddresses = append(storageProviderAddresses, f.Name())
				}
			}

			if len(storageProviderAddresses) == 0 {
				err := os.Remove(pathToAccStorage)
				if err != nil {
					logger.Log(logger.CreateDetails(location, err))
				}
				continue
			}

			removedTotal := 0

			for _, spAddress := range storageProviderAddresses {

				pathToStorProviderFiles := filepath.Join(pathToAccStorage, spAddress)

				dirFiles, err := nodeFile.ReadDirFiles(pathToStorProviderFiles)
				if err != nil {
					logger.Log(logger.CreateDetails(location, err))
					continue
				}

				fileNames := []string{}

				for _, f := range dirFiles {
					if len(f.Name()) == 64 && regFileName.MatchString(f.Name()) {
						fileNames = append(fileNames, f.Name())
					}
				}

				pathToFsTree := filepath.Join(paths.StoragePaths[0], network, spAddress, paths.SpFsFilename)

				if len(fileNames) == 0 {
					err := os.Remove(pathToFsTree)
					if err != nil {
						logger.Log(logger.CreateDetails(location, err))
					}

					err = os.Remove(pathToStorProviderFiles)
					if err != nil {
						logger.Log(logger.CreateDetails(location, err))
					}
					continue
				}

				shared.MU.Lock()
				fileFsTree, treeBytes, err := nodeFile.Read(pathToFsTree)
				if err != nil {
					shared.MU.Unlock()
					logger.Log(logger.CreateDetails(location, err))
					continue
				}

				fileFsTree.Close()
				shared.MU.Unlock()

				var spFs shared.StorageProviderData

				err = json.Unmarshal(treeBytes, &spFs)
				if err != nil {
					logger.Log(logger.CreateDetails(location, err))
				}

				fsFiles := map[string]bool{}

				for _, hashes := range spFs.Tree {
					for _, hash := range hashes {
						fsFiles[hex.EncodeToString(hash)] = true
					}
				}

				for _, fileName := range fileNames {

					if !fsFiles[fileName] {
						shared.MU.Lock()
						logger.Log("removing file: " + fileName + " of " + spAddress)
						stat, err := os.Stat(filepath.Join(pathToStorProviderFiles, fileName))
						if err != nil {
							logger.Log(logger.CreateDetails(location, err))
							continue
						}

						err = os.Remove(filepath.Join(pathToStorProviderFiles, fileName))
						if err != nil {
							logger.Log(logger.CreateDetails(location, err))
							continue
						}

						if !tstpkg.TestMode {
							logger.SendStatistic(spAddress, network, "", logger.Delete, stat.Size())
						}

						removedTotal++

						shared.MU.Unlock()
					}
				}

			}

			if removedTotal > 0 {
				pathToConfigFile := filepath.Join(paths.ConfigDirPath, paths.ConfFileName)

				shared.MU.Lock()
				confFile, fileBytes, err := nodeFile.Read(pathToConfigFile)
				if err != nil {
					shared.MU.Unlock()
					logger.Log(logger.CreateDetails(location, err))
					continue
				}

				var nodeConfig config.NodeConfig

				err = json.Unmarshal(fileBytes, &nodeConfig)
				if err != nil {
					shared.MU.Unlock()
					confFile.Close()
					logger.Log(logger.CreateDetails(location, err))
					continue
				}

				nodeConfig.UsedStorageSpace -= int64(removedTotal * oneMB)

				err = config.Save(confFile, nodeConfig)
				if err != nil {
					shared.MU.Unlock()
					confFile.Close()
					logger.Log(logger.CreateDetails(location, err))
					continue
				}
				confFile.Close()
				shared.MU.Unlock()
			}
		}

	}
}
