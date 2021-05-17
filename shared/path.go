package shared

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

var WorkingDir string
var AccDir string

//GetAccountDirectory return account directory of dfile products
func GetAccountDirectory() {

	accDir := filepath.Join(WorkingDir, "accounts")

	_, err := os.Stat(accDir)
	if err != nil {
		err = os.MkdirAll(accDir, os.ModePerm|os.ModeDir)
		if err != nil {
			log.Fatal("Fatal error")
		}
	}

	AccDir = accDir
}

// GetHomeDirectory return path to the home directory of dfile
func GetOrCreateWorkDir() {

	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatal("Fatal error")
	}

	homeDir = filepath.Join(homeDir, "dfile")

	_, err = os.Stat(homeDir)
	if err != nil {
		err = os.MkdirAll(homeDir, os.ModePerm|os.ModeDir)
		if err != nil {
			log.Fatal("Fatal error")
		}
	}

	WorkingDir = homeDir
}

// GetHomeDirectory return path to the app data of dfile secondary node
func GetDirectoryDFileSecondaryNode() (string, error) {

	nodeDir := filepath.Join(WorkingDir, "dfile-secondary-node")

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
