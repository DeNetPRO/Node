package config

import (
	"dfile-secondary-node/account"
	"dfile-secondary-node/crypto"
	"dfile-secondary-node/shared"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func GetConfigsList() (map[string]string, error) {
	configs := make(map[string]string)
	path, err := shared.GetConfigsDirectory()
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		configs[strings.TrimSuffix(f.Name(), ".json")] = filepath.Join(path, f.Name())
	}

	return configs, err
}

type SecondaryNodeConfig struct {
	Name         string `json:"configName"`
	HTTPPort     string `json:"portHTTP"`
	HTTPSPort    string `json:"portHTTPS"`
	Address      string `json:"publicAddress"`
	PathToConfig string `json:"pathToConfig"`
	StorageLimit int    `json:"storageLimit"`
	AccountName  string `json:"accountName"`
}

func Create(address string) (SecondaryNodeConfig, error) {

	config := SecondaryNodeConfig{}

	addressIsCorrect := false

	if address == "" {
		fmt.Println("Please select one of account addresses")
		accounts := account.GetAllAccounts()

		for i, a := range accounts {
			fmt.Println(i+1, a)
		}

		for !addressIsCorrect {
			accountAddress, err := shared.ReadFromConsole()
			if err != nil {
				return config, err
			}

			addressMatches := shared.ContainsAccount(accounts, accountAddress)

			if !addressMatches {
				fmt.Println("Account is incorrect, try again")
				continue
			}

			addressIsCorrect = true
			config.Address = accountAddress
			address = accountAddress

		}
	} else {
		config.Address = address
		fmt.Println("Now, a config file creation is needed.")
	}

	config.PathToConfig = filepath.Join(shared.AccDir, address, "config")

	regName := regexp.MustCompile("[A-Za-z0-9_]+")
	configNameIsCorrect := false

	for !configNameIsCorrect {
		fmt.Println("Please, enter config file name (use only letters, numbers and symbol '_')")

		name, err := shared.ReadFromConsole()
		if err != nil {
			return config, err
		}

		match := regName.MatchString(name)

		if !match {
			fmt.Println("Name is incorrect, please try again.")
			continue
		}

		configNames, err := GetConfigsList()
		if err != nil {
			return config, err
		}
		_, nameExists := configNames[name]

		if nameExists {
			fmt.Println("This config name exists, please try a new name")
			continue
		}
		config.Name = name
		configNameIsCorrect = true
	}

	spaceValueIsCorrect := false

	regNum := regexp.MustCompile(("[0-9]+"))

	for !spaceValueIsCorrect {
		fmt.Println("Enter disk space for usage in GB (should be positive number)")

		availableSpace := shared.GetAvailableSpace(config.PathToConfig)
		fmt.Println("Available space:", availableSpace, "GB")
		space, err := shared.ReadFromConsole()
		if err != nil {
			return config, err
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
		config.StorageLimit = intSpace

	}

	portHTTPValueIsCorrect := false
	regPort := regexp.MustCompile("[0-9]+|")

	for !portHTTPValueIsCorrect {
		fmt.Println("Enter http port number (value from 49152 to 65535) or press enter to use default port number 55050")

		httpPort, err := shared.ReadFromConsole()
		if err != nil {
			return config, err
		}

		match := regPort.MatchString(httpPort)
		if !match {
			fmt.Println("Value is incorrect, please try again")
			continue

		}

		if httpPort == "" {
			portHTTPValueIsCorrect = true
			config.HTTPPort = fmt.Sprint(55050)
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
		config.HTTPPort = fmt.Sprint(intHttpPort)
	}

	portHTTPSValueIsCorrect := false

	for !portHTTPSValueIsCorrect {
		fmt.Println("Enter http port number (value from 49152 to 65535)  or press enter to use default port number 55051")

		httpsPort, err := shared.ReadFromConsole()
		if err != nil {
			return config, err
		}

		match := regPort.MatchString(httpsPort)
		if !match {
			fmt.Println("Value is incorrect, please try again")
			continue
		}

		if httpsPort == "" {
			portHTTPSValueIsCorrect = true
			config.HTTPSPort = fmt.Sprint(55051)
			continue
		}
		intHttpsPort, err := strconv.Atoi(httpsPort)
		if err != nil {
			fmt.Println("Value is incorrect, please try again")
			continue
		}
		if intHttpsPort < 49152 || intHttpsPort > 65535 {
			fmt.Println("Value is incorrect, please try again")
			continue

		}
		portHTTPSValueIsCorrect = true
		config.HTTPSPort = fmt.Sprint(intHttpsPort)
		continue
	}

	config.AccountName = crypto.Sha256String([]byte(config.Address + config.Name))

	confFile, err := os.Create(filepath.Join(config.PathToConfig, config.Name+".json"))
	if err != nil {
		return config, err
	}
	defer confFile.Close()

	confJSON, err := json.Marshal(config)
	if err != nil {
		return config, err
	}

	_, err = confFile.Write(confJSON)
	if err != nil {
		return config, err
	}

	confFile.Sync()

	return config, nil
}
