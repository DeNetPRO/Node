package common

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
)

// GetHomeDirectory return path to the home directory of application files for the dfile-secondary-node application
// if not exist, create this directory
func GetAccountDirectory() (string, error) {

	homeDir, err := GetHomeDirectory()
	if err != nil {
		return "", err
	}

	homeDir = filepath.Join(homeDir, "accounts")

	_, err = os.Stat(homeDir)
	if err != nil {
		err = os.MkdirAll(homeDir, os.ModePerm|os.ModeDir)
		if err != nil {
			return "", err
		}
	}

	return homeDir, nil
}


func GetHomeDirectory() (string, error) {

	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	homeDir = filepath.Join(homeDir, "dfile")

	_, err = os.Stat(homeDir)
	if err != nil {
		err = os.MkdirAll(homeDir, os.ModePerm|os.ModeDir)
		if err != nil {
			return "", err
		}
	}

	return homeDir, nil
}