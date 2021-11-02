package spfiles

import (
	"encoding/hex"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/minio/sha256-simd"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/errs"

	fsysinfo "git.denetwork.xyz/DeNet/dfile-secondary-node/fsys_info"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/hash"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/sign"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/shared"
)

type NodesResponse struct {
	Nodes []string `json:"nodes"`
}

type NodeAddressResponse struct {
	NodeAddress string `json:"node_address"`
}

// ====================================================================================
//Save is used for checking and saving file parts from the inoming request to the node's storage.
func Save(req *http.Request, spData *shared.StorageProviderData, pathToSpFiles string) error {
	const location = "files.Save->"

	stat, err := os.Stat(pathToSpFiles)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return logger.CreateDetails(location, err)
	}

	if stat == nil {
		err = os.MkdirAll(pathToSpFiles, 0700)
		if err != nil {
			return logger.CreateDetails(location, err)
		}
	}

	err = fsysinfo.Save(pathToSpFiles, spData)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	reqFileParts := req.MultipartForm.File["files"]

	oneMBHashes, err := hash.OneMbParts(reqFileParts)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	var wholeFileHash string

	if len(oneMBHashes) == 1 {
		wholeFileHash = oneMBHashes[0]
	} else {
		sort.Strings(oneMBHashes)
		wholeFileHash, _, err = hash.CalcRoot(oneMBHashes)
		if err != nil {
			return logger.CreateDetails(location, err)
		}
	}

	fsContainsFile := false

	for _, fileHash := range spData.Fs {
		if fileHash == wholeFileHash {
			fsContainsFile = true
			break
		}
	}

	if !fsContainsFile {
		return logger.CreateDetails(location, errs.WrongFile)
	}

	count := 0
	savedParts := make([]string, 0, len(oneMBHashes))
	for _, reqFilePart := range reqFileParts {
		savedParts = append(savedParts, reqFilePart.Filename)

		rqFile, err := reqFilePart.Open()
		if err != nil {
			return logger.CreateDetails(location, err)
		}

		err = savePart(rqFile, pathToSpFiles, reqFilePart.Filename)
		if err != nil {
			rqFile.Close()
			deleteParts(pathToSpFiles, savedParts)
			return logger.CreateDetails(location, err)
		}

		rqFile.Close()

		count++

		logger.Log("Saved file " + reqFilePart.Filename + " (" + strconv.Itoa(count) + "/" + strconv.Itoa(len(oneMBHashes)) + ")" + " from " + spData.Address) //TODO remove

	}

	return nil
}

// ====================================================================================
func savePart(file io.Reader, pathToSpFiles, fileName string) error {
	const location = "files.saveFilePart->"

	pathToFile := filepath.Join(pathToSpFiles, fileName)

	newFile, err := os.Create(pathToFile)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	_, err = io.Copy(newFile, file)
	if err != nil {
		newFile.Close()
		return logger.CreateDetails(location, err)
	}

	newFile.Sync()
	newFile.Close()

	return nil
}

// ====================================================================================

func Serve(spAddress, fileKey, signatureFromReq, network string) (string, error) {
	const location = "files.Serve->"

	signature, err := hex.DecodeString(signatureFromReq)
	if err != nil {
		return "", logger.CreateDetails(location, err)
	}

	hash := sha256.Sum256([]byte(fileKey + spAddress))

	err = sign.Check(spAddress, signature, hash)
	if err != nil {
		return "", logger.CreateDetails(location, err)
	}

	pathToFile := filepath.Join(paths.StoragePaths[0], network, spAddress, fileKey)

	_, err = os.Stat(pathToFile)
	if err != nil {
		return "", logger.CreateDetails(location, err)
	}

	return pathToFile, nil
}

// ====================================================================================

// DeleteParts deletes parts of the file that wasn't fully uploaded to the node for some reason.
func deleteParts(addressPath string, fileHashes []string) {
	logger.Log("deleting file parts after error...")

	for _, hash := range fileHashes {
		pathToFile := filepath.Join(addressPath, hash)

		os.Remove(pathToFile)
	}
}

//Return storage provider filesystem path
func SearchStorageFilesystem(spAddress string) (string, bool) {
	path := filepath.Join(paths.SystemsDirPath, spAddress)
	stat, _ := os.Stat(path)
	if stat == nil {
		return "", false
	}

	return path, true
}

//Return storage provider filesystem path
func UpdateStorageFilesystem(spAddress string, fileSystemHeader *multipart.FileHeader) error {
	const location = "spFiles.UpdateStorageFilesystem"

	shared.MU.Lock()
	defer shared.MU.Unlock()

	fileSystem, err := fileSystemHeader.Open()
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	defer fileSystem.Close()

	path := filepath.Join(paths.SystemsDirPath, spAddress)

	file, err := os.Create(path)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	defer file.Close()

	_, err = io.Copy(file, fileSystem)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	return nil
}
