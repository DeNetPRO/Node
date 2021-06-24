package config

import (
	"context"
	blockchainprovider "dfile-secondary-node/blockchain_provider"
	"dfile-secondary-node/shared"
	"time"

	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type SecondaryNodeConfig struct {
	Address          string `json:"nodeAddress"`
	IpAddress        string `json:"ipAddress"`
	HTTPPort         string `json:"portHTTP"`
	StorageLimit     int    `json:"storageLimit"`
	UsedStorageSpace int64  `json:"usedStorageSpace"`
	AgreeSendLogs    bool   `json:"agreeSendLogs"`
}

type configState struct {
	Create string
	Update string
}

var State = configState{
	Create: "Create",
	Update: "Update",
}

var fullyReservedIPs = map[string]bool{
	"0":   true,
	"10":  true,
	"127": true,
}

var partiallyReservedIPs = map[string]int{
	"172": 31,
	"192": 168,
}

func Create(address, password string) (SecondaryNodeConfig, error) {
	const logInfo = "config.Create->"
	dFileConf := SecondaryNodeConfig{
		Address:       address,
		AgreeSendLogs: true,
	}

	fmt.Println("Now, a config file creation is needed.")

	pathToConfig := filepath.Join(shared.AccsDirPath, address, shared.ConfDirName)

	fmt.Println("Please enter disk space for usage in GB (should be positive number)")

	err := SetStorageLimit(pathToConfig, State.Create, &dFileConf)
	if err != nil {
		return dFileConf, fmt.Errorf("%s %w", logInfo, err)
	}

	ip, err := shared.InternetDevice.ExternalIP()
	if err != nil {
		return dFileConf, fmt.Errorf("%s %w", logInfo, err)
	}

	dFileConf.IpAddress = ip
	splitIPAddr := strings.Split(ip, ".")

	fmt.Println("Your IP address:", ip)
	fmt.Println("Remember if you don't have a static ip address it may change")
	fmt.Println("After router reset ip address info update may be needed")

	err = SetPort(&dFileConf, State.Create)
	if err != nil {
		return dFileConf, fmt.Errorf("%s %w", logInfo, err)
	}

	fmt.Println("Bug reports will be sent to developers")
	fmt.Println("You can switch off reports by updating config")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	fmt.Println("Registering node...")

	err = blockchainprovider.RegisterNode(ctx, address, password, splitIPAddr, dFileConf.HTTPPort)
	if err != nil {
		return dFileConf, fmt.Errorf("%s %w", logInfo, err)
	}

	confFile, err := os.Create(filepath.Join(pathToConfig, "config.json"))
	if err != nil {
		return dFileConf, fmt.Errorf("%s %w", logInfo, shared.GetDetailedError(err))
	}
	defer confFile.Close()

	confJSON, err := json.Marshal(dFileConf)
	if err != nil {
		return dFileConf, fmt.Errorf("%s %w", logInfo, shared.GetDetailedError(err))
	}

	_, err = confFile.Write(confJSON)
	if err != nil {
		return dFileConf, fmt.Errorf("%s %w", logInfo, shared.GetDetailedError(err))
	}

	fmt.Println("Saving config...")

	confFile.Sync()

	return dFileConf, nil
}

func SetStorageLimit(pathToConfig, state string, dFileConf *SecondaryNodeConfig) error {
	const logInfo = "config.SetStorageLimit->"
	regNum := regexp.MustCompile(("[0-9]+"))

	for {
		availableSpace := shared.GetAvailableSpace(pathToConfig)
		fmt.Println("Available space:", availableSpace, "GB")
		space, err := shared.ReadFromConsole()
		if err != nil {
			return fmt.Errorf("%s %w", logInfo, err)
		}

		if state == State.Update && space == "" {
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

		if intSpace < int(dFileConf.UsedStorageSpace) || intSpace >= availableSpace {
			fmt.Println("Passed value is out of avaliable space range, please try again")
			continue
		}

		dFileConf.StorageLimit = intSpace
		break
	}

	return nil
}

func SetIpAddr(dFileConf *SecondaryNodeConfig, state string) ([]string, error) {
	const logInfo = "config.SetIpAddr->"
	regIp := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)

	var splitIPAddr []string

	for {

		ipAddr, err := shared.ReadFromConsole()
		if err != nil {
			return nil, fmt.Errorf("%s %w", logInfo, err)
		}

		if state == State.Update && ipAddr == "" {
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
				return nil, fmt.Errorf("%s %w", logInfo, shared.GetDetailedError(err))
			}

			if secondAddrPart <= reservedSecAddrPart {
				fmt.Println("Address", ipAddr, "can't be used as a public ip address")
				continue
			}
		}

		dFileConf.IpAddress = ipAddr

		break
	}

	return splitIPAddr, nil
}

func SetPort(dFileConf *SecondaryNodeConfig, state string) error {
	const logInfo = "config.SetPort->"
	regPort := regexp.MustCompile("[0-9]+|")

	for {
		fmt.Println("Enter http port number (value from 49152 to 65535) or press enter to use default port number 55050")

		httpPort, err := shared.ReadFromConsole()
		if err != nil {
			return fmt.Errorf("%s %w", logInfo, err)
		}

		if state == State.Create && httpPort == "" {
			dFileConf.HTTPPort = fmt.Sprint(55050)
			break
		}

		if state == State.Update && httpPort == "" {
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

		dFileConf.HTTPPort = fmt.Sprint(intHttpPort)
		break
	}

	return nil
}

func ChangeAgreeSendLogs(dFileConf *SecondaryNodeConfig, state string) error {
	const logInfo = "config.ChangeAgreeSendLogs->"
	regPort := regexp.MustCompile("^(?:y|n)$")

	for {
		agree, err := shared.ReadFromConsole()
		if err != nil {
			return fmt.Errorf("%s %w", logInfo, err)
		}

		if state == State.Update && agree == "" {
			break
		}

		if !regPort.MatchString(agree) {
			fmt.Println("Value is incorrect, please try again. [y/n]")
			continue
		}

		if agree == "y" {
			dFileConf.AgreeSendLogs = true
			shared.SendLogs = true
		} else {
			dFileConf.AgreeSendLogs = false
			shared.SendLogs = false
		}

		break
	}

	return nil
}
