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

	"git.denetwork.xyz/DeNet/dfile-secondary-node/networks"
	"github.com/ricochet2200/go-disk-usage/du"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	nodeTypes "git.denetwork.xyz/DeNet/dfile-secondary-node/node_types"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"

	termEmul "git.denetwork.xyz/DeNet/dfile-secondary-node/term_emul"
	tstpkg "git.denetwork.xyz/DeNet/dfile-secondary-node/tst_pkg"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/upnp"
)

type Statuses struct {
	Create string
	Update string
}

var stats = Statuses{
	Create: "creare",
	Update: "update",
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

var RPC string

func Stats() Statuses {
	return stats
}

//Creates a new acount configuration
func Create(address string) (nodeTypes.Config, error) {
	const location = "config.Create->"

	paths.SetConfigPath(address)

	var nodeConfig nodeTypes.Config

	if tstpkg.Data().TestMode {
		nodeConfig = tstpkg.TestConfig()
		RPC = nodeConfig.RPC[nodeConfig.Network]
		setStorage(address, &nodeConfig)
		networks.Set(nodeConfig.Network)
	} else {
		nodeConfig = nodeTypes.Config{
			Address:              address,
			StoragePaths:         []string{},
			SendBugReports:       true,
			RegisteredInNetworks: map[string]bool{},
			RPC: map[string]string{"kovan": "https://kovan.infura.io/v3/45b81222fded4427b3a6589e0396c596",
				"polygon": "https://polygon-rpc.com"},
		}

		RPC = nodeConfig.RPC[nodeConfig.Network]

		err := setStorage(address, &nodeConfig)
		if err != nil {
			return nodeConfig, logger.MarkLocation(location, err)
		}

		fmt.Println("\nHow much GB are you going to share? (should be positive number)")

		err = SetStorageLimit(&nodeConfig, stats.Create)
		if err != nil {
			return nodeConfig, logger.MarkLocation(location, err)
		}

		err = SetNetwork(&nodeConfig)
		if err != nil {
			return nodeConfig, logger.MarkLocation(location, err)
		}

		if upnp.InternetDevice != nil {
			ip, err := upnp.InternetDevice.PublicIP()
			if err != nil {
				return nodeConfig, logger.MarkLocation(location, err)
			}
			nodeConfig.IpAddress = ip
			fmt.Println("Your public IP address", ip, "is added to config")
		} else {
			fmt.Println("\nPlease enter your public ip address")
			err = SetIpAddr(&nodeConfig, stats.Create)
			if err != nil {
				return nodeConfig, logger.MarkLocation(location, err)
			}
		}

		err = SetPort(&nodeConfig, stats.Create)
		if err != nil {
			return nodeConfig, logger.MarkLocation(location, err)
		}
	}

	err := os.MkdirAll(paths.List().ConfigDir, 0700)
	if err != nil {
		return nodeConfig, logger.MarkLocation(location, err)
	}

	confFile, err := os.Create(paths.List().ConfigFile)
	if err != nil {
		return nodeConfig, logger.MarkLocation(location, err)
	}
	defer confFile.Close()

	confJSON, err := json.Marshal(nodeConfig)
	if err != nil {
		return nodeConfig, logger.MarkLocation(location, err)
	}

	_, err = confFile.Write(confJSON)
	if err != nil {
		return nodeConfig, logger.MarkLocation(location, err)
	}

	fmt.Println("Saving config...")

	confFile.Sync()

	return nodeConfig, nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

// Creates storage directory and adds paths to storage dir/dirs
func setStorage(address string, nodeConfig *nodeTypes.Config) error {
	const location = "config.SetStoragePath"

	strgPathAndNodeAddr := filepath.Join("storage", address)
	defaultPath := filepath.Join(paths.List().WorkDir, strgPathAndNodeAddr)

	// -------
	nodeConfig.StoragePaths = append(nodeConfig.StoragePaths, defaultPath) //TODO select paths
	paths.SetStoragePaths(nodeConfig.StoragePaths)

	err := os.MkdirAll(paths.List().Storages[0], 0700)
	if err != nil {
		return logger.MarkLocation(location, err)
	}
	return nil
	// -------

	mountPoints, err := paths.GetMountPoints()
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	if tstpkg.Data().TestMode || len(mountPoints) == 0 {

		err := paths.CreateStorage(defaultPath)
		if err != nil {
			return logger.MarkLocation(location, err)
		}

		nodeConfig.StoragePaths = append(nodeConfig.StoragePaths, defaultPath)
		paths.SetStoragePaths(nodeConfig.StoragePaths)
		fmt.Println("Storage is set to default path")
		return nil
	}

	for {
		fmt.Println("\nPlease select path to storage. You can type several numbers by splitting them with space.")
		fmt.Println("Type * to select all paths, or type another full path to storage." + "\n")

		pointNum := 0

		fmt.Println(fmt.Sprint(pointNum, ". ", defaultPath, " (default path)"))

		for _, point := range mountPoints {
			pointNum++
			fmt.Println(fmt.Sprint(pointNum, ". ", filepath.Join(point, strgPathAndNodeAddr)))
		}

		validatedPaths, err := validatePaths(strgPathAndNodeAddr, defaultPath, mountPoints)
		if err != nil {
			return logger.MarkLocation(location, err)
		}

		for _, path := range validatedPaths {
			err := paths.CreateStorage(path)
			if err != nil {
				continue
			}

			nodeConfig.StoragePaths = append(nodeConfig.StoragePaths, path)
		}

		if len(nodeConfig.StoragePaths) == 0 {
			continue
		}

		paths.SetStoragePaths(nodeConfig.StoragePaths)

		return nil
	}

}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

// Validates storage path info from the input
func validatePaths(strgPathAndNodeAddr, defaultPath string, mountPoints []string) ([]string, error) {

	const location = "config.validateStoragePaths ->"

	inputPaths := []string{}

	reg := regexp.MustCompile(("^[a-zA-Z0-9_.-]*$"))

	for {

		input, err := termEmul.ReadInput()
		if err != nil {
			return inputPaths, logger.MarkLocation(location, err)
		}

		if strings.Trim(input, " ") == "*" {
			inputPaths = []string{defaultPath}
			for _, point := range mountPoints {
				inputPaths = append(inputPaths, filepath.Join(point, strgPathAndNodeAddr))
			}

			return inputPaths, nil
		}

		splitInput := strings.Split(input, " ")

		alreadySelected := map[int]bool{}

		for _, inputPart := range splitInput {

			if strings.Contains(strings.Trim(inputPart, " "), "/") {
				inputPaths = append(inputPaths, inputPart)
				continue
			}

			indxNum, err := strconv.Atoi(inputPart)
			if err != nil {
				continue
			}

			if alreadySelected[indxNum] {
				continue
			}

			if indxNum < 0 || indxNum > len(mountPoints) {
				continue
			}

			if indxNum == 0 {
				inputPaths = append(inputPaths, defaultPath)
				continue
			}

			inputPaths = append(inputPaths, filepath.Join(mountPoints[indxNum-1], strgPathAndNodeAddr))

			alreadySelected[indxNum] = true

		}

		if len(inputPaths) == 0 {
			fmt.Println("no valid value, please try again")
			continue
		}

		pathExists := map[string]bool{}

		validatedPaths := []string{}

		for _, path := range inputPaths {

			if pathExists[path] {
				continue
			}

			splitPath := strings.Split(path, "")

			onlySymbols := true

			for _, char := range splitPath {
				if reg.MatchString(char) {
					onlySymbols = false
				}
			}

			if splitPath[0] != "/" || onlySymbols {
				pathExists[path] = true
				continue
			}

			validatedPaths = append(validatedPaths, path)

		}

		if len(validatedPaths) == 0 {
			fmt.Println("please try again")
			continue
		}

		return validatedPaths, nil

	}
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

// Adds selected network name to config
func SetNetwork(nodeConfig *nodeTypes.Config) error {
	const location = "config.SetNetwork"

	fmt.Println("\nChoose a network")

	supportedNets := networks.List()

	indx := 1

	for _, network := range supportedNets {
		phrase := network + " "

		// TODO remove
		if network == "polygon" || network == "mumbai" {
			phrase += "(not available in current version)"
		}

		fmt.Println(indx, phrase)
		indx++
	}

	for {
		index, err := termEmul.ReadInput()
		if err != nil {
			return logger.MarkLocation(location, err)
		}

		netIndxNum, err := strconv.Atoi(index)
		if err != nil {
			fmt.Println("Incorrect value, try again")
			continue
		}

		if netIndxNum < 1 || netIndxNum > len(supportedNets) {
			fmt.Println("Incorrect value, try again")
			continue
		}

		err = networks.Set(supportedNets[netIndxNum-1])
		if err != nil {
			fmt.Println("Incorrect value, try again")
			continue
		}

		// TODO remove
		if supportedNets[netIndxNum-1] == "polygon" || supportedNets[netIndxNum-1] == "mumbai" {
			fmt.Println("Network is not supported, please try again")
			continue
		}

		nodeConfig.Network = supportedNets[netIndxNum-1]

		return nil
	}

}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

// Adds allocated disk space info to config
func SetStorageLimit(nodeConfig *nodeTypes.Config, state string) error {
	const location = "config.SetStorageLimit->"

	regNum := regexp.MustCompile(("[0-9]+"))

	for _, path := range nodeConfig.StoragePaths {
		for {

			availableSpace, err := getAvailableSpace(path)
			if err != nil {
				return logger.MarkLocation(location, err)
			}

			fmt.Println("Current available space:", availableSpace, "GB")
			space, err := termEmul.ReadInput()
			if err != nil {
				return logger.MarkLocation(location, err)
			}

			if state == stats.Update && space == "" {
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

				if tstpkg.Data().TestMode {
					return logger.MarkLocation(location, errors.New("out of range"))
				}

				fmt.Println("Passed value is out of avaliable space range, please try again")
				continue
			}

			nodeConfig.StorageLimit += intSpace // TODO check if paths are in one partition
			break
		}
	}

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

//Adds ip address info to config
func SetIpAddr(nodeConfig *nodeTypes.Config, state string) error {
	const location = "config.SetIpAddr->"

	regIp := regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`) // check regex

	for {
		ipAddr, err := termEmul.ReadInput()
		if err != nil {
			return logger.MarkLocation(location, err)
		}

		if state == stats.Update && ipAddr == "" {
			break
		}

		match := regIp.MatchString(ipAddr)

		if !match {
			if tstpkg.Data().TestMode {
				return errors.New("incorrect value")
			}

			fmt.Println("Value is incorrect, please try again")
			continue
		}

		splitIPAddr := strings.Split(ipAddr, ".")

		if fullyReservedIPs[splitIPAddr[0]] {

			if tstpkg.Data().TestMode {
				return errors.New("can't be used as a public ip address")
			}

			fmt.Println("Address", ipAddr, "can't be used as a public ip address")
			continue
		}

		reservedSecAddrPart, partiallyReserved := partiallyReservedIPs[splitIPAddr[0]]

		if partiallyReserved {
			secondAddrPart, err := strconv.Atoi(splitIPAddr[1])
			if err != nil {
				return logger.MarkLocation(location, err)
			}

			if secondAddrPart <= reservedSecAddrPart {

				if tstpkg.Data().TestMode {
					return errors.New("can't be used as a public ip address")
				}

				fmt.Println("Address", ipAddr, "can't be used as a public ip address")
				continue
			}
		}

		nodeConfig.IpAddress = ipAddr
		break
	}

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

//Adds port info to config
func SetPort(nodeConfig *nodeTypes.Config, state string) error {
	const location = "config.SetPort->"

	regPort := regexp.MustCompile("[0-9]+|")

	for {
		fmt.Println("Enter http port number (value from 49152 to 65535) or press enter to use default port number 55050")

		httpPort, err := termEmul.ReadInput()
		if err != nil {
			return logger.MarkLocation(location, err)
		}

		if state == stats.Create && httpPort == "" {
			nodeConfig.HTTPPort = fmt.Sprint(":", 55050)
			break
		}

		if state == stats.Update && httpPort == "" {
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

			if tstpkg.Data().TestMode {
				return errors.New("incorrect value")
			}

			fmt.Println("Value is incorrect, please try again")
			continue
		}

		nodeConfig.HTTPPort = fmt.Sprint(":", intHttpPort)
		break
	}

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

// Turn on/off bug report sending
func SwitchReports(nodeConfig *nodeTypes.Config, state string) error {
	const location = "config.SetBugReportsFlag->"
	regPort := regexp.MustCompile("^(?:y|n)$")

	for {
		agree, err := termEmul.ReadInput()
		if err != nil {
			return logger.MarkLocation(location, err)
		}

		if state == stats.Update && agree == "" {
			break
		}

		if !regPort.MatchString(agree) {
			fmt.Println("Value is incorrect, please try again. [y/n]")
			continue
		}

		if agree == "y" {
			nodeConfig.SendBugReports = true
			logger.SendReports = true
		} else {
			nodeConfig.SendBugReports = false
			logger.SendReports = false
		}

		break
	}

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

//Saves the configuration file
func Save(configFile *os.File, Config nodeTypes.Config) error {
	confJSON, err := json.Marshal(Config)
	if err != nil {
		return err
	}

	err = configFile.Truncate(0)
	if err != nil {
		return err
	}

	_, err = configFile.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = configFile.Write(confJSON)
	if err != nil {
		return err
	}

	err = configFile.Sync()
	if err != nil {
		return err
	}

	configFile.Close()

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

//Return nodes available space in GB
func getAvailableSpace(path string) (int, error) {
	const location = "config.GetAvailableSpace ->"

	const KB = uint64(1024)

	_, err := os.Stat(path)
	if err != nil {
		return 0, logger.MarkLocation(location, err)
	}

	usage := du.NewDiskUsage(path)
	return int(usage.Available() / (KB * KB * KB)), nil
}
