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

	blockchainprovider "git.denetwork.xyz/dfile/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	termEmul "git.denetwork.xyz/dfile/dfile-secondary-node/term_emul"
	"git.denetwork.xyz/dfile/dfile-secondary-node/upnp"
)

type NodeConfig struct {
	Address          string `json:"nodeAddress"`
	ChnClntAddr      string `json:"chainAddress"`
	IpAddress        string `json:"ipAddress"`
	HTTPPort         string `json:"portHTTP"`
	NFT              string `json:"nft"`
	StorageLimit     int    `json:"storageLimit"`
	UsedStorageSpace int64  `json:"usedStorageSpace"`
	AgreeSendLogs    bool   `json:"agreeSendLogs"`
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
		AgreeSendLogs: true,
		ChnClntAddr:   "https://kovan.infura.io/v3/6433ee0efa38494a85541b00cd377c5f",
		NFT:           "0xBfAfdaE6B77a02A4684D39D1528c873961528342",
	}

	blockchainprovider.NFT = nodeConfig.NFT
	blockchainprovider.ChainClientAddr = nodeConfig.ChnClntAddr

	fmt.Println("Now, a config file creation is needed.")

	pathToConfig := filepath.Join(paths.AccsDirPath, address, paths.ConfDirName)

	fmt.Println("Please enter disk space for usage in GB (should be positive number)")

	err := SetStorageLimit(pathToConfig, CreateStatus, &nodeConfig)
	if err != nil {
		return nodeConfig, logger.CreateDetails(location, err)
	}

	var splitIPAddr []string

	if upnp.InternetDevice != nil {
		ip, err := upnp.InternetDevice.PublicIP()
		if err != nil {
			return nodeConfig, logger.CreateDetails(location, err)
		}

		nodeConfig.IpAddress = ip
		splitIPAddr = strings.Split(ip, ".")
		fmt.Println("Your public IP address", ip, "is added to config")
	} else {
		fmt.Println("Please enter your public ip address")
		splitIPAddr, err = SetIpAddr(&nodeConfig, CreateStatus)
		if err != nil {
			return nodeConfig, logger.CreateDetails(location, err)
		}
	}

	err = SetPort(&nodeConfig, CreateStatus)
	if err != nil {
		return nodeConfig, logger.CreateDetails(location, err)
	}

	if !shared.TestMode {
		fmt.Println("Due to testing stage bug reports from your device are going to be received by developers")
		fmt.Println("You can stop sending reports by updating config")
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		fmt.Println("Registering node...")

		err = blockchainprovider.RegisterNode(ctx, address, password, splitIPAddr, nodeConfig.HTTPPort)
		if err != nil {
			return nodeConfig, logger.CreateDetails(location, err)
		}
	}

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

//Set storage limit in config file
func SetStorageLimit(pathToConfig, state string, nodeConfig *NodeConfig) error {
	const location = "config.SetStorageLimit->"
	regNum := regexp.MustCompile(("[0-9]+"))

	if shared.TestMode {
		nodeConfig.StorageLimit = shared.TestLimit
		return nil
	}

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
func SetIpAddr(nodeConfig *NodeConfig, state string) ([]string, error) {
	const location = "config.SetIpAddr->"

	var splitIPAddr []string

	if shared.TestMode {
		ipAddr := shared.TestAddress
		splitIPAddr := strings.Split(ipAddr, ".")
		nodeConfig.IpAddress = ipAddr
		return splitIPAddr, nil
	}

	regIp := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)

	for {
		ipAddr, err := termEmul.ReadInput()
		if err != nil {
			return nil, logger.CreateDetails(location, err)
		}

		if state == UpdateStatus && ipAddr == "" {
			break
		}

		match := regIp.MatchString(ipAddr)

		if !match {
			fmt.Println("Value is incorrect, please try again")
			continue
		}

		splitIPAddr = strings.Split(ipAddr, ".")

		if fullyReservedIPs[splitIPAddr[0]] {
			fmt.Println("Address", ipAddr, "can't be used as a public ip address")
			continue
		}

		reservedSecAddrPart, partiallyReserved := partiallyReservedIPs[splitIPAddr[0]]

		if partiallyReserved {
			secondAddrPart, err := strconv.Atoi(splitIPAddr[1])
			if err != nil {
				return nil, logger.CreateDetails(location, err)
			}

			if secondAddrPart <= reservedSecAddrPart {
				fmt.Println("Address", ipAddr, "can't be used as a public ip address")
				continue
			}
		}

		nodeConfig.IpAddress = ipAddr
		break
	}

	return splitIPAddr, nil
}

// ====================================================================================

//Set port in config file
func SetPort(nodeConfig *NodeConfig, state string) error {
	const location = "config.SetPort->"

	if shared.TestMode {
		nodeConfig.HTTPPort = shared.TestPort
		return nil
	}

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
