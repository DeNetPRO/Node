package spfiles

import (
	"os"
	"path/filepath"

	"github.com/DeNetPRO/src/logger"
	"github.com/DeNetPRO/src/paths"
)

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func SaveChunk(pathToSpFiles, fileName string, spFileChunk []byte) error {
	const location = "files.SaveChunk->"

	pathToFile := filepath.Join(pathToSpFiles, fileName)

	f, err := os.Create(pathToFile)
	if err != nil {
		return logger.MarkLocation(location, err)
	}
	defer f.Close()

	_, err = f.Write(spFileChunk)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	f.Sync()

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

// DeleteParts deletes parts of the file that wasn't fully uploaded to the node for some reason.
func deleteParts(addressPath string, fileHashes []string) {
	logger.Log("deleting file parts after error...")

	for _, hash := range fileHashes {
		pathToFile := filepath.Join(addressPath, hash)

		os.Remove(pathToFile)
	}
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

// Return storage provider filesystem path
func SearchStorageFilesystem(spAddress string) (string, bool) {
	path := filepath.Join(paths.List().SysDir, spAddress)
	stat, _ := os.Stat(path)
	if stat == nil {
		return "", false
	}

	return path, true
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
