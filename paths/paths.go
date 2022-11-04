package paths

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	nodeTypes "git.denetwork.xyz/DeNet/dfile-secondary-node/node_types"
	tstpkg "git.denetwork.xyz/DeNet/dfile-secondary-node/tst_pkg"
)

var (
	workDirName  = "denet-node"
	confDirName  = "config"
	confFileName = "config.json"
)

var paths = nodeTypes.Paths{SpFsFilename: "sp_fs.json"}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

//Initializes default node paths
func Init() error {
	const location = "shared.InitPaths->"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	if tstpkg.Data().TestMode {
		workDirName = tstpkg.Data().WorkDirName
	}

	paths.WorkDir = filepath.Join(homeDir, workDirName)
	paths.AccsDir = filepath.Join(paths.WorkDir, "accounts")
	paths.UpdateDir = filepath.Join(paths.WorkDir, "update")
	paths.SysDir = filepath.Join(paths.WorkDir, "systems")

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func SetConfigPath(addr string) {
	paths.ConfigDir = filepath.Join(paths.AccsDir, addr, confDirName)
	paths.ConfigFile = filepath.Join(paths.ConfigDir, confFileName)
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func SetStoragePaths(ps []string) {
	paths.Storages = ps
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func List() nodeTypes.Paths {
	return paths
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

//Creates account dir if it doesn't already exist
func CreateAccDirs() error {
	const location = "shared.CreateIfNotExistAccDirs->"
	statWDP, err := os.Stat(paths.WorkDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return logger.MarkLocation(location, err)
	}

	if statWDP == nil {
		err = os.MkdirAll(paths.WorkDir, os.ModePerm|os.ModeDir)
		if err != nil {
			return logger.MarkLocation(location, err)
		}
	}

	statADP, err := os.Stat(paths.AccsDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return logger.MarkLocation(location, err)
	}

	if statADP == nil {
		err = os.MkdirAll(paths.AccsDir, os.ModePerm|os.ModeDir)
		if err != nil {
			return logger.MarkLocation(location, err)
		}
	}

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func GetMountPoints() ([]string, error) {
	const location = "config.getMountPoints =>"

	mountPoints := []string{}

	if runtime.GOOS == "windows" {
		return mountPoints, nil
	} else {
		var stdout bytes.Buffer

		cmd := exec.Command("df")
		cmd.Stdout = &stdout
		err := cmd.Run()
		if err != nil {
			return mountPoints, logger.MarkLocation(location, err)
		}

		splitRes := strings.Split(stdout.String(), "\n")

		for _, strng := range splitRes {
			splitStrng := strings.Split(strng, " ")

			if strings.HasPrefix(splitStrng[0], "/dev") && !strings.Contains(splitStrng[0], "loop") {
				mountPoint := splitStrng[len(splitStrng)-1]

				if mountPoint == "/" {
					continue
				}

				mountPoints = append(mountPoints, mountPoint)
			}
		}
	}

	return mountPoints, nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func CreateStorage(path string) error {

	err := os.MkdirAll(path, 0700)
	if err != nil {
		return errors.New("couldn't create storage with path " + path)
	}

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
