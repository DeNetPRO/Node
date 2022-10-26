package fsysinfo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/pb"

	nodeFile "git.denetwork.xyz/DeNet/dfile-secondary-node/node_file"
	nodeTypes "git.denetwork.xyz/DeNet/dfile-secondary-node/node_types"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
)

var mutex sync.Mutex

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func Save(fsInfo *pb.FsInfo, fsTree [][][]byte) error {
	const location = "fsys_info.Save->"

	pathToSpFiles := filepath.Join(paths.List().Storages[0], fsInfo.Network, fsInfo.SpAddress)

	pathToSpFs := filepath.Join(pathToSpFiles, paths.List().SpFsFilename)

	stat, err := os.Stat(pathToSpFs)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return logger.MarkLocation(location, errors.New("no files of "+fsInfo.SpAddress))
	}

	mutex.Lock()
	defer mutex.Unlock()

	newFsInfo := nodeTypes.StorageProviderData{
		Nonce:        fsInfo.Nonce,
		Storage:      fsInfo.Storage,
		SignedFsInfo: fsInfo.Signature,
		Tree:         fsTree,
	}

	if stat == nil {

		file, err := os.Create(pathToSpFs)
		if err != nil {
			return logger.MarkLocation(location, err)
		}
		defer file.Close()

		err = nodeFile.Write(file, newFsInfo)
		if err != nil {
			return logger.MarkLocation(location, err)
		}

		file.Sync()

		return nil

	}

	file, bytes, err := nodeFile.Read(pathToSpFs)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	defer file.Close()

	var previousFsInfo nodeTypes.StorageProviderData

	err = json.Unmarshal(bytes, &previousFsInfo)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	if newFsInfo.Nonce < previousFsInfo.Nonce {
		return logger.MarkLocation(location, fmt.Errorf("%v fs info is up to date", fsInfo.SpAddress))
	}

	err = nodeFile.Write(file, newFsInfo)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	file.Sync()

	return nil

}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

//Returns storage provider filesystem path
func BackUpSPFsys(spAddress string, fileSystemHeader *multipart.FileHeader) error {
	const location = "spFiles.UpdateStorageFilesystem"

	mutex.Lock()
	defer mutex.Unlock()

	fileSystem, err := fileSystemHeader.Open()
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	defer fileSystem.Close()

	path := filepath.Join(paths.List().SysDir, spAddress)

	file, err := os.Create(path)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	defer file.Close()

	_, err = io.Copy(file, fileSystem)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
