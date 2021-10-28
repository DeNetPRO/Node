package config

import (
	"context"

	"time"

	"encoding/json"
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
	AgreeSendLogs        bool            `json:"agreeSendLogs"`
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

//Create is used for creating a config file.
func Create(address, password string) (NodeConfig, error) {
	const location = "config.Create->"

	nodeConfig := NodeConfig{
		Address:       address,
		StoragePaths:  []string{filepath.Join(paths.WorkDirPath, paths.StorageDirName, address)},
		AgreeSendLogs: true,
	}

	pathToConfig := filepath.Join(paths.AccsDirPath, address, paths.ConfDirName)

	if shared.TestMode {
		nodeConfig.IpAddress = "127.0.01"
		nodeConfig.HTTPPort = "55050"
		nodeConfig.Network = "kovan"
		nodeConfig.StorageLimit = 1
		nodeConfig.UsedStorageSpace = 0
		nodeConfig.StoragePaths = []string{filepath.Join(paths.WorkDirPath, paths.StorageDirName, address)}

	} else {
		network, err := SelectNetwork()
		if err != nil {
			return nodeConfig, logger.CreateDetails(location, err)
		}

		nodeConfig.Network = network
		blckChain.CurrentNetwork = network

		fmt.Println("Please enter disk space for usage in GB (should be positive number)")

		err = SetStorageLimit(pathToConfig, CreateStatus, &nodeConfig)
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		fmt.Println("Registering node...")

		err = blckChain.RegisterNode(ctx, address, password, nodeConfig.IpAddress, nodeConfig.HTTPPort)
		if err != nil {
			return nodeConfig, logger.CreateDetails(location, err)
		}

		nodeConfig.RegisteredInNetworks[blckChain.CurrentNetwork] = true
	}

	paths.StoragePaths = nodeConfig.StoragePaths

	confFile, err := os.Create(filepath.Join(pathToConfig, paths.ConfFileName))
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

	currentNets := make(map[int]string)

	indx := 1

	for network := range blckChain.Networks {
		fmt.Println(indx, network)
		currentNets[indx] = network
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

		if netIndx < 1 || netIndx > indx-1 {
			fmt.Println("Incorrect value, try again")
			continue
		}

		netwrok := currentNets[netIndx]

		return netwrok, nil
	}

}

// ====================================================================================

//Set storage limit in config file
func SetStorageLimit(pathToConfig, state string, nodeConfig *NodeConfig) error {
	const location = "config.SetStorageLimit->"
	regNum := regexp.MustCompile(("[0-9]+"))

	for {
		availableSpace := shared.GetAvailableSpace(pathToConfig)
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

		if intSpace < int(nodeConfig.UsedStorageSpace) || intSpace >= availableSpace {
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

	regIp := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)

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
			nodeConfig.AgreeSendLogs = true
			logger.SendLogs = true
		} else {
			nodeConfig.AgreeSendLogs = false
			logger.SendLogs = false
		}

		break
	}

	return nil
}

// ====================================================================================

//Saving config file
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
