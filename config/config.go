package config

import (
	"bufio"
	"dfile-secondary-node/account"
	"dfile-secondary-node/crypto"
	"dfile-secondary-node/shared"
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
	Name                   string `json:"config_name"`
	HTTPPort               int    `json:"HTTP_port"`
	HTTPSPort              int    `json:"HTTPS_port"`
	Address                string `json:"public_address"`
	PathToStorageDirectory string `json:"path_to_storage"`
	LimitGB                int    `json:"storage_limit"`
	AccountName            string `json:"account_name"`
}

func (config *SecondaryNodeConfig) Create() error {

	normalInput := false

	// Config name
	for !normalInput {
		fmt.Println("Enter config name (use only letters, numbers and symbol '_')")
		name, err := readFromConsole()
		if err != nil {
			return err
		}

		match, err := regexp.MatchString("[A-Za-z0-9_]+", name)
		if err != nil {
			return err
		}

		if match {
			configNames, err := GetConfigsList()
			if err != nil {
				return err
			}
			if _, ok := configNames[name]; !ok {
				config.Name = name
				normalInput = true
			} else {
				fmt.Println("This config name is exist, try again")
			}

		} else {
			fmt.Println("Name is incorrect, try again")
		}
	}

	// Account name

	fmt.Println("Please select on of account addresses")
	accounts, err := account.GetAllAccounts()
	if err != nil {
		return err
	}
	for i, a := range accounts {
		fmt.Println(i, a)
	}

	normalInput = false

	for !normalInput {
		accountAddress, err := readFromConsole()
		if err != nil {
			return err
		}

		inAccounts := containsAccount(accounts, accountAddress)

		if inAccounts {
			normalInput = true
			config.Address = accountAddress
		} else {
			fmt.Println("Account is incorrect, try again")
		}
	}

	//storage path

	normalInput = false

	for !normalInput {
		fmt.Println("Enter storage path (path to store DFile decentralized files)")
		storagePath, err := readFromConsole()
		if err != nil {
			return err
		}

		storagePath = filepath.FromSlash(storagePath)
		err = os.MkdirAll(storagePath, os.ModePerm|os.ModeDir)
		if err != nil {
			fmt.Println("Bad path, try again")
		} else {
			normalInput = true
			config.PathToStorageDirectory = storagePath
		}
	}

	// storage space
	normalInput = false

	for !normalInput {
		fmt.Println("Enter disk space for usage in GB (should be positive integer number)")

		availableSpace := shared.GetAvailableSpace(config.PathToStorageDirectory)
		fmt.Println("Available space:", availableSpace)
		space, err := readFromConsole()
		if err != nil {
			return err
		}

		match, err := regexp.MatchString("[0-9]+", space)
		if match {
			intSpace, err := strconv.Atoi(space)
			if err != nil {
				fmt.Println("Bad number, try again")
				continue
			}
			if (intSpace >= 0) && (intSpace <= availableSpace) {
				normalInput = true
				config.LimitGB = intSpace
			} else {
				fmt.Println("Bad number, try again")
				continue
			}
		} else {
			fmt.Println("Bad number, try again")
			continue
		}

	}

	// http port
	normalInput = false

	for !normalInput {
		fmt.Println("Enter http port number (press enter to use default 48654):")

		httpPort, err := readFromConsole()
		if err != nil {
			return err
		}

		match, err := regexp.MatchString("[0-9]+|", httpPort)
		if match {

			if httpPort == "" {
				normalInput = true
				config.HTTPPort = 48654
				continue
			}
			intHttpPort, err := strconv.Atoi(httpPort)
			if err != nil {
				fmt.Println("Bad number, try again")
				continue
			}
			if (intHttpPort >= 1) && (intHttpPort <= 65535) {
				normalInput = true
				config.HTTPPort = intHttpPort
				continue
			} else {
				fmt.Println("Bad number, try again")
				continue
			}
		} else {
			fmt.Println("Bad number, try again")
			continue
		}
	}

	// https port
	normalInput = false

	for !normalInput {
		fmt.Println("Enter https port number (press enter to use default 48654):")

		httpsPort, err := readFromConsole()
		if err != nil {
			return err
		}

		match, err := regexp.MatchString("[0-9]+|", httpsPort)
		if match {

			if httpsPort == "" {
				normalInput = true
				config.HTTPSPort = 48655
				continue
			}
			intHttpsPort, err := strconv.Atoi(httpsPort)
			if err != nil {
				fmt.Println("Bad number, try again")
				continue
			}
			if (intHttpsPort >= 1) && (intHttpsPort <= 65535) {
				normalInput = true
				config.HTTPSPort = intHttpsPort
				continue
			} else {
				fmt.Println("Bad number, try again")
				continue
			}
		} else {
			fmt.Println("Bad number, try again")
			continue
		}
	}
	accountName := crypto.Sha256String([]byte(config.Address + config.Name))
	config.AccountName = accountName

	return nil
}

func readFromConsole() (string, error) {
	fmt.Print("Enter text: ")
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// remove the delimiter from the string
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")

	return input, err
}

func containsAccount(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
