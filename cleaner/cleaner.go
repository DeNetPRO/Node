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
	"sync"
	"time"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/config"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/networks"
	nodeFile "git.denetwork.xyz/DeNet/dfile-secondary-node/node_file"
	nodeTypes "git.denetwork.xyz/DeNet/dfile-secondary-node/node_types"
	tstpkg "git.denetwork.xyz/DeNet/dfile-secondary-node/tst_pkg"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
)

const oneMB = 1048576

var mutex sync.Mutex

var unusedFiles = map[string]bool{}

//Starts cleaner, that checks if stored file part is in Storage Provider's file system and deletes it if it was not found.
func Start() {
	const location = "cleaner.Start->"

	var markedToDelete = map[string]int64{} // in case if file was uploaded but fs tree wasn't updated yet

	regAddr := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	regFileName := regexp.MustCompile("[0-9A-Za-z_]")

	for {
		time.Sleep(time.Minute * 10)

		currentNets := networks.List()

		for _, network := range currentNets {
			pathToAccStorage := filepath.Join(paths.List().Storages[0], network)

			stat, err := os.Stat(pathToAccStorage)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				logger.Log(logger.MarkLocation(location, err))
				log.Fatal(err)
			}

			if stat == nil {
				continue
			}

			dirFiles, err := nodeFile.ReadDirFiles(pathToAccStorage)
			if err != nil {
				logger.Log(logger.MarkLocation(location, err))
				continue
			}

			storageProviderAddresses := []string{}

			for _, f := range dirFiles {
				if regAddr.MatchString(f.Name()) {
					storageProviderAddresses = append(storageProviderAddresses, f.Name())
				}
			}

			if len(storageProviderAddresses) == 0 {
				continue
			}

			removedTotal := 0

			for _, spAddress := range storageProviderAddresses {

				pathToStorProviderFiles := filepath.Join(pathToAccStorage, spAddress)

				// _, unusedFiles := unusedFiles[spAddress]

				// if unusedFiles {

				// 	dirFiles, err := nodeFile.ReadDirFiles(pathToStorProviderFiles)
				// 	if err != nil {
				// 		logger.Log(logger.MarkLocation(location, err))
				// 		continue
				// 	}

				// 	err = os.RemoveAll(pathToStorProviderFiles)
				// 	if err != nil {
				// 		continue
				// 	}

				// 	fmt.Println("removed", len(dirFiles), "unused files of", spAddress)

				// 	err = restoreSpaceInConfig(len(dirFiles))
				// 	if err != nil {
				// 		logger.Log(logger.MarkLocation(location, err))
				// 		continue
				// 	}

				// 	unmarkUnused(spAddress)
				// 	continue
				// }

				dirFiles, err := nodeFile.ReadDirFiles(pathToStorProviderFiles)
				if err != nil {
					logger.Log(logger.MarkLocation(location, err))
					continue
				}

				fileNames := []string{}

				for _, f := range dirFiles {
					if len(f.Name()) == 64 && regFileName.MatchString(f.Name()) {
						fileNames = append(fileNames, f.Name())
					}
				}

				pathToFsTree := filepath.Join(paths.List().Storages[0], network, spAddress, paths.List().SpFsFilename)

				if len(fileNames) == 0 {
					err := os.Remove(pathToFsTree)
					if err != nil {
						logger.Log(logger.MarkLocation(location, err))
					}

					err = os.Remove(pathToStorProviderFiles)
					if err != nil {
						logger.Log(logger.MarkLocation(location, err))
					}
					continue
				}

				treeBytes, err := os.ReadFile(pathToFsTree)
				if err != nil {
					logger.Log(logger.MarkLocation(location, err))
					continue
				}

				var spFs nodeTypes.StorageProviderData

				err = json.Unmarshal(treeBytes, &spFs)
				if err != nil {
					logger.Log(logger.MarkLocation(location, err))
					continue
				}

				fsInfo := map[string]bool{}

				for _, hashes := range spFs.Tree {
					for _, hash := range hashes {
						fsInfo[hex.EncodeToString(hash)] = true
					}
				}

				if len(fsInfo) == 0 {
					continue
				}

				for _, fileName := range fileNames {

					_, marked := markedToDelete[fileName]

					if !fsInfo[fileName] {

						if !marked {
							markedToDelete[fileName] = time.Now().Unix()
							continue
						}

						if time.Now().Unix()-markedToDelete[fileName] > 60*60*2 {
							mutex.Lock()
							fmt.Println("removing file: " + fileName + " of " + spAddress)
							stat, err := os.Stat(filepath.Join(pathToStorProviderFiles, fileName))
							if err != nil {
								mutex.Unlock()
								logger.Log(logger.MarkLocation(location, err))
								continue
							}

							err = os.Remove(filepath.Join(pathToStorProviderFiles, fileName))
							if err != nil {
								mutex.Unlock()
								logger.Log(logger.MarkLocation(location, err))
								continue
							}

							if !tstpkg.Data().TestMode {
								logger.SendStatistic(spAddress, network, "", logger.Delete, stat.Size())
							}

							removedTotal++

							delete(markedToDelete, fileName)

							mutex.Unlock()
						}

					} else {
						if marked {
							delete(markedToDelete, fileName)
						}
					}

				}

			}

			if removedTotal > 0 {
				err := restoreSpaceInConfig(removedTotal)
				if err != nil {
					logger.Log(logger.MarkLocation(location, err))
					continue
				}
			}
		}

	}
}

func restoreSpaceInConfig(space int) error {

	const location = "cleaner.restoreSpaceInConfig ->"

	mutex.Lock()
	confFile, fileBytes, err := nodeFile.Read(paths.List().ConfigFile)
	if err != nil {
		mutex.Unlock()
		return logger.MarkLocation(location, err)
	}

	var nodeConfig nodeTypes.Config

	err = json.Unmarshal(fileBytes, &nodeConfig)
	if err != nil {
		mutex.Unlock()
		confFile.Close()
		return logger.MarkLocation(location, err)
	}

	nodeConfig.UsedStorageSpace -= int64(space * oneMB)

	err = config.Save(confFile, nodeConfig)
	if err != nil {
		mutex.Unlock()
		confFile.Close()
		return logger.MarkLocation(location, err)
	}
	mutex.Unlock()

	fmt.Println("cleaned", space, "Mbytes")

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func MarkUnused(spAddr string) {
	mutex.Lock()

	_, alreadyMarked := unusedFiles[spAddr]

	if alreadyMarked {
		mutex.Unlock()
		return
	}

	unusedFiles[spAddr] = true
	fmt.Println("marked", spAddr, "files as unused")

	mutex.Unlock()
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func unmarkUnused(spAddr string) {
	mutex.Lock()

	_, alreadyUnarked := unusedFiles[spAddr]

	if alreadyUnarked {
		mutex.Unlock()
		return
	}

	delete(unusedFiles, spAddr)
	mutex.Unlock()
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
