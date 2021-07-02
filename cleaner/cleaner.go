package cleaner

import (
	"dfile-secondary-node/paths"
	"dfile-secondary-node/shared"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func Start() {

	const logInfo = "cleaner.Start->"

	regAddr := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	regFileName := regexp.MustCompile("[0-9A-Za-z_]")

	nodeAddr, err := shared.DecryptNodeAddr()
	if err != nil {
		shared.LogError(logInfo, shared.GetDetailedError(err))
	}

	for {
		time.Sleep(5 * time.Second)

		pathToAccStorage := filepath.Join(paths.AccsDirPath, nodeAddr.String(), paths.StorageDirName)

		storageProviderAddresses := []string{}

		err = filepath.WalkDir(pathToAccStorage,
			func(path string, info fs.DirEntry, err error) error {
				if err != nil {
					shared.LogError(logInfo, shared.GetDetailedError(err))
				}

				if regAddr.MatchString(info.Name()) {
					storageProviderAddresses = append(storageProviderAddresses, info.Name())
				}

				return nil
			})

		if err != nil {
			shared.LogError(logInfo, shared.GetDetailedError(err))
			continue
		}

		if len(storageProviderAddresses) == 0 {
			continue
		}

		fileNames := []string{}

		for _, spAddress := range storageProviderAddresses {

			pathToFsTree := filepath.Join(paths.AccsDirPath, nodeAddr.String(), paths.StorageDirName, spAddress, "tree.json")

			pathToStorProviderFiles := filepath.Join(pathToAccStorage, spAddress)

			err = filepath.WalkDir(pathToStorProviderFiles,
				func(path string, info fs.DirEntry, err error) error {
					if err != nil {
						shared.LogError(logInfo, shared.GetDetailedError(err))
					}

					if regFileName.MatchString(info.Name()) && len(info.Name()) == 64 {
						fileNames = append(fileNames, info.Name())
					}

					return nil
				})
			if err != nil {
				shared.LogError(logInfo, shared.GetDetailedError(err))
				continue
			}

			shared.MU.Lock()
			fileFsTree, err := os.Open(pathToFsTree)
			if err != nil {
				shared.MU.Unlock()
				shared.LogError(logInfo, shared.GetDetailedError(err))
			}

			treeBytes, err := io.ReadAll(fileFsTree)
			if err != nil {
				fileFsTree.Close()
				shared.MU.Unlock()
				shared.LogError(logInfo, shared.GetDetailedError(err))
			}
			fileFsTree.Close()
			shared.MU.Unlock()

			var storageFsStruct shared.StorageInfo

			err = json.Unmarshal(treeBytes, &storageFsStruct)
			if err != nil {
				shared.LogError(logInfo, shared.GetDetailedError(err))
			}

			fsFiles := map[string]bool{}

			for _, hashes := range storageFsStruct.Tree {
				for _, hash := range hashes {
					fsFiles[hex.EncodeToString(hash)] = true
				}
			}

			for _, fileName := range fileNames {

				fmt.Println(fileName, fsFiles[fileName])

				if !fsFiles[fileName] {
					shared.MU.Lock()
					fmt.Println("removing file", fileName)
					err = os.Remove(filepath.Join(pathToStorProviderFiles, fileName))
					if err != nil {
						shared.LogError(logInfo, shared.GetDetailedError(err))
					}

					shared.MU.Lock()

					fsFiles = nil
					continue
				}
			}

		}

	}

}
