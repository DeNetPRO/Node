package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	blckChain "git.denetwork.xyz/DeNet/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/shared"
	termEmul "git.denetwork.xyz/DeNet/dfile-secondary-node/term_emul"
	tstpkg "git.denetwork.xyz/DeNet/dfile-secondary-node/tst_pkg"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/upnp"
)

type NodeConfig struct {
	Address              string          `json:"nodeAddress"`
	IpAddress            string          `json:"ipAddress"`
	HTTPPort             string          `json:"portHTTP"`
	Network              string          `json:"network"`
	StorageLimit         int             `json:"storageLimit"`
	StoragePaths         []string        `json:"storagePaths"`
	UsedStorageSpace     int64           `json:"usedStorageSpace"`
	SendBugReports       bool            `json:"sendBugReports"`
	RegisteredInNetworks map[string]bool `json:"registeredInNetworks"`
}

const (
	CreateStatus = "Create"
	UpdateStatus = "Update"
)

var fullyReservedIPs = map[string]bool{
	"0":   true,
	"10":  true,
	"127": true,
}

var partiallyReservedIPs = map[string]int{
	"172": 31,
	"192": 168,
}

var TestConfig = NodeConfig{
	Address:              tstpkg.TestAccAddr,
	StoragePaths:         []string{filepath.Join(tstpkg.TestWorkDirName, paths.StorageDirName, tstpkg.TestAccAddr)},
	SendBugReports:       true,
	RegisteredInNetworks: map[string]bool{},
	IpAddress:            tstpkg.TestIP,
	HTTPPort:             tstpkg.TestPort,
	Network:              tstpkg.TestNetwork,
	StorageLimit:         tstpkg.TestStorageLimit,
	UsedStorageSpace:     int64(tstpkg.TestUsedStorageSpace),
}

//Create is used for creating a config file.
func Create(address string) (NodeConfig, error) {
	const location = "config.Create->"

	var nodeConfig NodeConfig

	paths.ConfigDirPath = filepath.Join(paths.AccsDirPath, address, paths.ConfDirName)

	if tstpkg.TestMode {
		nodeConfig = TestConfig
		paths.StoragePaths = nodeConfig.StoragePaths
	} else {
		nodeConfig = NodeConfig{
			Address:              address,
			StoragePaths:         []string{filepath.Join(paths.WorkDirPath, paths.StorageDirName, address)},
			SendBugReports:       true,
			RegisteredInNetworks: map[string]bool{},
		}

		paths.StoragePaths = nodeConfig.StoragePaths

		network, err := SelectNetwork()
		if err != nil {
			return nodeConfig, logger.CreateDetails(location, err)
		}

		nodeConfig.Network = network
		blckChain.CurrentNetwork = network

		fmt.Println("Please enter disk space for usage in GB (should be positive number)")

		err = SetStorageLimit(&nodeConfig, CreateStatus)
		if err != nil {
			return nodeConfig, logger.CreateDetails(location, err)
		}

		if upnp.InternetDevice != nil {
			ip, err := upnp.InternetDevice.PublicIP()
			if err != nil {
				return nodeConfig, logger.CreateDetails(location, err)
			}
			nodeConfig.IpAddress = ip
			fmt.Println("Your public IP address", ip, "is added to config")
		} else {
			fmt.Println("Please enter your public ip address")
			err = SetIpAddr(&nodeConfig, CreateStatus)
			if err != nil {
				return nodeConfig, logger.CreateDetails(location, err)
			}
		}

		err = SetPort(&nodeConfig, CreateStatus)
		if err != nil {
			return nodeConfig, logger.CreateDetails(location, err)
		}

		fmt.Println("Due to testing stage bug reports from your device are going to be received by developers")
		fmt.Println("You can stop sending reports by updating config")

	}

	err := os.MkdirAll(paths.ConfigDirPath, 0700)
	if err != nil {
		fmt.Println(err)
		return nodeConfig, logger.CreateDetails(location, err)
	}

	confFile, err := os.Create(filepath.Join(paths.ConfigDirPath, paths.ConfFileName))
	if err != nil {
		return nodeConfig, logger.CreateDetails(location, err)
	}
	defer confFile.Close()

	confJSON, err := json.Marshal(nodeConfig)
	if err != nil {
		return nodeConfig, logger.CreateDetails(location, err)
	}

	_, err = confFile.Write(confJSON)
	if err != nil {
		return nodeConfig, logger.CreateDetails(location, err)
	}

	fmt.Println("Saving config...")

	confFile.Sync()

	return nodeConfig, nil
}

// ====================================================================================

// returns selected network
func SelectNetwork() (string, error) {
	const location = "config.SelectNetwork"

	fmt.Println("Choose a network")

	currentNets := make([]string, 0, len(blckChain.Networks))

	indx := 1

	for network := range blckChain.Networks {
		fmt.Println(indx, network)
		currentNets = append(currentNets, network)
		indx++
	}

	for {
		number, err := termEmul.ReadInput()
		if err != nil {
			return "", logger.CreateDetails(location, err)
		}

		netIndx, err := strconv.Atoi(number)
		if err != nil {
			fmt.Println("Incorrect value, try again")
			continue
		}

		if netIndx < 1 || netIndx > len(currentNets) {
			fmt.Println("Incorrect value, try again")
			continue
		}

		netwrok := currentNets[netIndx-1]

		return netwrok, nil
	}

}

// ====================================================================================

//Set storage limit in config file
func SetStorageLimit(nodeConfig *NodeConfig, state string) error {
	const location = "config.SetStorageLimit->"
	regNum := regexp.MustCompile(("[0-9]+"))

	for {
		availableSpace, err := shared.GetAvailableSpace()
		if err != nil {
			return logger.CreateDetails(location, err)
		}

		fmt.Println("Available space:", availableSpace, "GB")
		space, err := termEmul.ReadInput()
		if err != nil {
			return logger.CreateDetails(location, err)
		}

		if state == UpdateStatus && space == "" {
			break
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

		if intSpace <= 0 || intSpace >= availableSpace {

			if tstpkg.TestMode {
				return logger.CreateDetails(location, errors.New("out of range"))
			}

			fmt.Println("Passed value is out of avaliable space range, please try again")
			continue
		}

		nodeConfig.StorageLimit = intSpace
		break
	}

	return nil
}

// ====================================================================================

//Set ip address in config file
func SetIpAddr(nodeConfig *NodeConfig, state string) error {
	const location = "config.SetIpAddr->"

	regIp := regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`) // check regex

	for {
		ipAddr, err := termEmul.ReadInput()
		if err != nil {
			return logger.CreateDetails(location, err)
		}

		if state == UpdateStatus && ipAddr == "" {
			break
		}

		match := regIp.MatchString(ipAddr)

		if !match {
			fmt.Println("Value is incorrect, please try again")
			continue
		}

		splitIPAddr := strings.Split(ipAddr, ".")

		if fullyReservedIPs[splitIPAddr[0]] {
			fmt.Println("Address", ipAddr, "can't be used as a public ip address")
			continue
		}

		reservedSecAddrPart, partiallyReserved := partiallyReservedIPs[splitIPAddr[0]]

		if partiallyReserved {
			secondAddrPart, err := strconv.Atoi(splitIPAddr[1])
			if err != nil {
				return logger.CreateDetails(location, err)
			}

			if secondAddrPart <= reservedSecAddrPart {
				fmt.Println("Address", ipAddr, "can't be used as a public ip address")
				continue
			}
		}

		nodeConfig.IpAddress = ipAddr
		break
	}

	return nil
}

// ====================================================================================

//Set port in config file
func SetPort(nodeConfig *NodeConfig, state string) error {
	const location = "config.SetPort->"

	regPort := regexp.MustCompile("[0-9]+|")

	for {
		fmt.Println("Enter http port number (value from 49152 to 65535) or press enter to use default port number 55050")

		httpPort, err := termEmul.ReadInput()
		if err != nil {
			return logger.CreateDetails(location, err)
		}

		if state == CreateStatus && httpPort == "" {
			nodeConfig.HTTPPort = fmt.Sprint(55050)
			break
		}

		if state == UpdateStatus && httpPort == "" {
			break
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

		nodeConfig.HTTPPort = fmt.Sprint(intHttpPort)
		break
	}

	return nil
}

// ====================================================================================

//Changing the sending logs agreement
func ChangeAgreeSendLogs(nodeConfig *NodeConfig, state string) error {
	const location = "config.ChangeAgreeSendLogs->"
	regPort := regexp.MustCompile("^(?:y|n)$")

	for {
		agree, err := termEmul.ReadInput()
		if err != nil {
			return logger.CreateDetails(location, err)
		}

		if state == UpdateStatus && agree == "" {
			break
		}

		if !regPort.MatchString(agree) {
			fmt.Println("Value is incorrect, please try again. [y/n]")
			continue
		}

		if agree == "y" {
			nodeConfig.SendBugReports = true
			logger.SendLogs = true
		} else {
			nodeConfig.SendBugReports = false
			logger.SendLogs = false
		}

		break
	}

	return nil
}

// ====================================================================================

//Saving config file
//TODO add os.WriteFile
func Save(confFile *os.File, nodeConfig NodeConfig) error {
	confJSON, err := json.Marshal(nodeConfig)
	if err != nil {
		return err
	}

	err = confFile.Truncate(0)
	if err != nil {
		return err
	}

	_, err = confFile.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = confFile.Write(confJSON)
	if err != nil {
		return err
	}

	err = confFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
