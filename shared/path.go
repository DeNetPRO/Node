package shared

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	WorkDir string
	AccDir  string
)

// GetHomeDirectory return path to the home directory of dfile
func CreateIfNotExistAccDirs() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Fatal error")
	}

	workDir := filepath.Join(homeDir, "dfile")

	_, err = os.Stat(workDir)
	if err != nil {
		errPart := strings.Split(err.Error(), ":")

		if strings.Trim(errPart[1], " ") != "no such file or directory" {
			log.Fatal("Fatal error")
		}

		err = os.MkdirAll(workDir, os.ModePerm|os.ModeDir)
		if err != nil {
			log.Fatal("Fatal error")
		}
	}

	WorkDir = workDir

	accDir := filepath.Join(WorkDir, "accounts")

	_, err = os.Stat(accDir)
	if err != nil {
		errPart := strings.Split(err.Error(), ":")

		if strings.Trim(errPart[1], " ") != "no such file or directory" {
			log.Fatal("Fatal error")
		}

		err = os.MkdirAll(accDir, os.ModePerm|os.ModeDir)
		if err != nil {
			log.Fatal("Fatal error")
		}
	}

	AccDir = accDir

}

// GetHomeDirectory return path to the app data of dfile secondary node
func GetDirectoryDFileSecondaryNode() (string, error) {

	nodeDir := filepath.Join(WorkDir, "dfile-secondary-node")

	_, err := os.Stat(nodeDir)
	if err != nil {
		err = os.MkdirAll(nodeDir, os.ModePerm|os.ModeDir)
		if err != nil {
			return "", err
		}
	}

	return nodeDir, nil
}

// GetConfigsDirectory return path to the app data of dfile secondary node
func GetConfigsDirectory() (string, error) {

	homeDir, err := GetDirectoryDFileSecondaryNode()
	if err != nil {
		return "", err
	}

	homeDir = filepath.Join(homeDir, "configs")

	_, err = os.Stat(homeDir)
	if err != nil {
		err = os.MkdirAll(homeDir, os.ModePerm|os.ModeDir)
		if err != nil {
			return "", err
		}
	}

	return homeDir, nil
}
