package common

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

//GetAccountDirectory return account directory of dfile products
func GetAccountDirectory() (string, error) {

	homeDir, err := GetOrCreateHomeDirectory()
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

// GetHomeDirectory return path to the home directory of dfile
func GetOrCreateHomeDirectory() (string, error) {

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

// GetHomeDirectory return path to the app data of dfile secondary node
func GetDirectoryDFileSecondaryNode() (string, error) {

	homeDir, err := GetOrCreateHomeDirectory()
	if err != nil {
		return "", err
	}

	homeDir = filepath.Join(homeDir, "dfile-secondary-node")

	_, err = os.Stat(homeDir)
	if err != nil {
		err = os.MkdirAll(homeDir, os.ModePerm|os.ModeDir)
		if err != nil {
			return "", err
		}
	}

	return homeDir, nil
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
