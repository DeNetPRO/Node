package internal

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
)

func BackUpCopy(req *http.Request, spData *shared.StorageProviderData, pathToConfig string, fileSize int) error {
	const logLoc = "internal.BackUpCopy->"

	addressPath := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spData.Address)

	err := initSPFile(addressPath, spData)
	if err != nil {
		RestoreMemoryInfo(pathToConfig, fileSize)
		return logger.CreateDetails(logLoc, err)
	}

	reqFileParts := req.MultipartForm.File["files"]

	const eightKB = 8192

	oneMBHashes := make([]string, 0, len(reqFileParts))

	for _, reqFilePart := range reqFileParts {

		eightKBHashes := make([]string, 0, 128)

		var buf bytes.Buffer

		rqFile, err := reqFilePart.Open()
		if err != nil {
			RestoreMemoryInfo(pathToConfig, fileSize)
			return logger.CreateDetails(logLoc, err)
		}

		_, err = io.Copy(&buf, rqFile)
		if err != nil {
			rqFile.Close()
			RestoreMemoryInfo(pathToConfig, fileSize)
			return logger.CreateDetails(logLoc, err)
		}

		rqFile.Close()

		bufBytes := buf.Bytes()

		for i := 0; i < len(bufBytes); i += eightKB {
			hSum := sha256.Sum256(bufBytes[i : i+eightKB])
			eightKBHashes = append(eightKBHashes, hex.EncodeToString(hSum[:]))
		}

		oneMBHash, _, err := shared.CalcRootHash(eightKBHashes)
		if err != nil {
			RestoreMemoryInfo(pathToConfig, fileSize)
			return logger.CreateDetails(logLoc, err)
		}

		if reqFilePart.Filename != oneMBHash {
			RestoreMemoryInfo(pathToConfig, fileSize)
			return logger.CreateDetails(logLoc, shared.ErrWrongFile)
		}

		oneMBHashes = append(oneMBHashes, oneMBHash)
	}

	fsContainsFile := false

	var wholeFileHash string

	if len(oneMBHashes) == 1 {
		wholeFileHash = oneMBHashes[0]
	} else {
		sort.Strings(oneMBHashes)
		wholeFileHash, _, err = shared.CalcRootHash(oneMBHashes)
		if err != nil {
			RestoreMemoryInfo(pathToConfig, fileSize)
			return logger.CreateDetails(logLoc, err)
		}
	}

	for _, fileHash := range spData.Fs {
		if fileHash == wholeFileHash {
			fsContainsFile = true
		}
	}

	if !fsContainsFile {
		RestoreMemoryInfo(pathToConfig, fileSize)
		return logger.CreateDetails(logLoc, shared.ErrWrongFile)
	}

	count := 1
	total := len(oneMBHashes)

	hashes := req.MultipartForm.File["hashes"]
	hashesFile, err := hashes[0].Open()
	if err != nil {
		RestoreMemoryInfo(pathToConfig, fileSize)
		return logger.CreateDetails(logLoc, err)
	}

	hashesBody, err := io.ReadAll(hashesFile)
	if err != nil {
		hashesFile.Close()
		RestoreMemoryInfo(pathToConfig, fileSize)
		return logger.CreateDetails(logLoc, err)
	}

	hashDif := make(map[string]string)
	err = json.Unmarshal(hashesBody, &hashDif)
	if err != nil {
		hashesFile.Close()
		RestoreMemoryInfo(pathToConfig, fileSize)
		return logger.CreateDetails(logLoc, err)
	}

	hashesFile.Close()

	for _, reqFilePart := range reqFileParts {
		rqFile, err := reqFilePart.Open()
		if err != nil {
			deleteFileParts(addressPath, oneMBHashes)
			RestoreMemoryInfo(pathToConfig, fileSize)
			return logger.CreateDetails(logLoc, err)
		}
		defer rqFile.Close()

		pathToFile := filepath.Join(addressPath, hashDif[reqFilePart.Filename])

		newFile, err := os.Create(pathToFile)
		if err != nil {
			deleteFileParts(addressPath, oneMBHashes)
			RestoreMemoryInfo(pathToConfig, fileSize)
			return logger.CreateDetails(logLoc, err)
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, rqFile)
		if err != nil {
			deleteFileParts(addressPath, oneMBHashes)
			RestoreMemoryInfo(pathToConfig, fileSize)
			return logger.CreateDetails(logLoc, err)
		}

		logger.Log("Saved file " + hashDif[reqFilePart.Filename] + " (" + fmt.Sprint(count) + "/" + fmt.Sprint(total) + ")" + " from " + spData.Address) //TODO remove

		newFile.Sync()
		rqFile.Close()
		newFile.Close()

		count++
	}

	return nil
}
