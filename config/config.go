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
	HTTPPort     string `json:"portHTTP"`
	HTTPSPort    string `json:"portHTTPS"`
	Address      string `json:"publicAddress"`
	PathToConfig string `json:"pathToConfig"`
	StorageLimit int    `json:"storageLimit"`
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
		fmt.Println("Enter https port number (value from 49152 to 65535)  or press enter to use default port number 55051")

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

	confFile, err := os.Create(filepath.Join(config.PathToConfig, "config.json"))
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
