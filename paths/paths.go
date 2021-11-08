package paths

import (
	"errors"
	"os"
	"path/filepath"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/shared"
)

var (
	WorkDirPath    string
	AccsDirPath    string
	WorkDirName    = "denet-node"
	ConfDirName    = "config"
	ConfigDirPath  string
	ConfFileName   = "config.json"
	StorageDirName = "storage"
	StoragePaths   []string
	SpFsFilename   = "tree.json"
	UpdateDirPath  string
	SystemsDirName = "systems"
	SystemsDirPath string
)

// ====================================================================================

//Initializes default node paths
func Init() error {
	const location = "shared.InitPaths->"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	if shared.TestMode {
		WorkDirName = shared.TestWorkDirName
	}

	WorkDirPath = filepath.Join(homeDir, WorkDirName)
	AccsDirPath = filepath.Join(WorkDirPath, "accounts")
	UpdateDirPath = filepath.Join(WorkDirPath, "update")
	SystemsDirPath = filepath.Join(WorkDirPath, SystemsDirName)

	return nil
}

// ====================================================================================

//Creates account dir if it doesn't already exist
func CreateAccDirs() error {
	const location = "shared.CreateIfNotExistAccDirs->"
	statWDP, err := os.Stat(WorkDirPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return logger.CreateDetails(location, err)
	}

	if statWDP == nil {
		err = os.MkdirAll(WorkDirPath, os.ModePerm|os.ModeDir)
		if err != nil {
			return logger.CreateDetails(location, err)
		}
	}

	statADP, err := os.Stat(AccsDirPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return logger.CreateDetails(location, err)
	}

	if statADP == nil {
		err = os.MkdirAll(AccsDirPath, os.ModePerm|os.ModeDir)
		if err != nil {
			return logger.CreateDetails(location, err)
		}
	}

	return nil
}
