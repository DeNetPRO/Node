package cleaner

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	blckChain "git.denetwork.xyz/DeNet/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/errs"
	nodeFile "git.denetwork.xyz/DeNet/dfile-secondary-node/node_file"

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
		time.Sleep(time.Minute * 20)

		pathToAccStorage := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, blckChain.CurrentNetwork)

		stat, err := os.Stat(pathToAccStorage)
		if err != nil {
			err = errs.CheckStatErr(err)
			if err != nil {
				logger.Log(logger.CreateDetails(location, err))
				log.Fatal(err)
			}
		}

		if stat == nil {
			fmt.Println("no files from", blckChain.CurrentNetwork, "to delete")
			continue
		}

		storageProviderAddresses := []string{}

		err = filepath.WalkDir(pathToAccStorage,
			func(path string, info fs.DirEntry, err error) error {
				if err != nil {
					logger.Log(logger.CreateDetails(location, err))
				}

				if regAddr.MatchString(info.Name()) {
					storageProviderAddresses = append(storageProviderAddresses, info.Name())
				}

				return nil
			})

		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			continue
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

			fileNames := []string{}

			pathToStorProviderFiles := filepath.Join(pathToAccStorage, spAddress)

			err = filepath.WalkDir(pathToStorProviderFiles,
				func(path string, info fs.DirEntry, err error) error {
					if err != nil {
						logger.Log(logger.CreateDetails(location, err))
					}

					if regFileName.MatchString(info.Name()) && len(info.Name()) == 64 {
						fileNames = append(fileNames, info.Name())
					}

					return nil
				})
			if err != nil {
				logger.Log(logger.CreateDetails(location, err))
				continue
			}

			pathToFsTree := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, blckChain.CurrentNetwork, spAddress, paths.SpFsFilename)

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

			testMode := os.Getenv("DENET_TEST")

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

					if testMode != "1" {
						logger.SendStatistic(spAddress, "", logger.Delete, stat.Size())
					}

					removedTotal++

					shared.MU.Unlock()
				}
			}

		}

		if removedTotal > 0 {
			pathToConfig := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.ConfDirName, paths.ConfFileName)

			shared.MU.Lock()
			confFile, fileBytes, err := nodeFile.Read(pathToConfig)
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
