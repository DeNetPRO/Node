package config

import (
	"dfile-secondary-node/account"
	"dfile-secondary-node/shared"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type SecondaryNodeConfig struct {
	HTTPPort         string `json:"portHTTP"`
	Address          string `json:"publicAddress"`
	StorageLimit     int    `json:"storageLimit"`
	UsedStorageSpace int64  `json:"usedStorageSpace"`
}

func Create(address string) (SecondaryNodeConfig, error) {

	dFileConf := SecondaryNodeConfig{}

	addressIsCorrect := false

	if address == "" {
		fmt.Println("Please select one of account addresses")
		accounts := account.List()

		for i, a := range accounts {
			fmt.Println(i+1, a)
		}

		for !addressIsCorrect {
			accountAddress, err := shared.ReadFromConsole()
			if err != nil {
				return dFileConf, err
			}

			addressMatches := shared.ContainsAccount(accounts, accountAddress)

			if !addressMatches {
				fmt.Println("Account is incorrect, try again")
				continue
			}

			addressIsCorrect = true
			dFileConf.Address = accountAddress
			address = accountAddress

		}
	} else {
		dFileConf.Address = address
		fmt.Println("Now, a config file creation is needed.")
	}

	pathToConfig := filepath.Join(shared.AccsDirPath, address, shared.ConfDirName)

	spaceValueIsCorrect := false

	regNum := regexp.MustCompile(("[0-9]+"))

	for !spaceValueIsCorrect {
		fmt.Println("Please enter disk space for usage in GB (should be positive number)")

		availableSpace := shared.GetAvailableSpace(pathToConfig)
		fmt.Println("Available space:", availableSpace, "GB")
		space, err := shared.ReadFromConsole()
		if err != nil {
			return dFileConf, err
		}

		match := regNum.MatchString(space)

		if !match {
			fmt.Println("Value is incorrect, please try again")
			continue
		}

		intSpace, err := strconv.Atoi(space)
		if err != nil {
			fmt.Println("Value is incorrect, please try again")
			continue
		}

		if intSpace < 0 || intSpace >= availableSpace {
			fmt.Println("Value is incorrect, please try again")
			continue
		}

		spaceValueIsCorrect = true
		dFileConf.StorageLimit = intSpace

	}

	portHTTPValueIsCorrect := false
	regPort := regexp.MustCompile("[0-9]+|")

	for !portHTTPValueIsCorrect {
		fmt.Println("Enter http port number (value from 49152 to 65535) or press enter to use default port number 55050")

		httpPort, err := shared.ReadFromConsole()
		if err != nil {
			return dFileConf, err
		}

		if httpPort == "" {
			portHTTPValueIsCorrect = true
			dFileConf.HTTPPort = fmt.Sprint(55050)
			continue
		}

		match := regPort.MatchString(httpPort)
		if !match {
			fmt.Println("Value is incorrect, please try again")
			continue

		}

		intHttpPort, err := strconv.Atoi(httpPort)
		if err != nil {
			fmt.Println("Value is incorrect, please try again")
			continue
		}
		if intHttpPort < 49152 || intHttpPort > 65535 {
			fmt.Println("Value is incorrect, please try again")
			continue

		}

		portHTTPValueIsCorrect = true
		dFileConf.HTTPPort = fmt.Sprint(intHttpPort)
	}

	confFile, err := os.Create(filepath.Join(pathToConfig, "config.json"))
	if err != nil {
		return dFileConf, err
	}
	defer confFile.Close()

	confJSON, err := json.Marshal(dFileConf)
	if err != nil {
		return dFileConf, err
	}

	_, err = confFile.Write(confJSON)
	if err != nil {
		return dFileConf, err
	}

	confFile.Sync()

	return dFileConf, nil
}
