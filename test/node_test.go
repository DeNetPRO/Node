package test

import (
	"dfile-secondary-node/account"
	"dfile-secondary-node/config"
	"dfile-secondary-node/shared"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	accountPassword      = "123"
	accountAddress       string
	nodeAddress          []byte
	ErrorInvalidPassword = errors.New("could not decrypt key with given password")
	configName           = "config.json"
	configPort           = "55051"
	configAddress        string
	configStorageLimit   = 1
	ip                   = "185.140.19.95"
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

func TestEmptyAccountListBeforeCreating(t *testing.T) {
	accs := account.List()
	want := 0
	get := len(accs)

	require.Equal(t, want, get)
}

func TestCreateAccount(t *testing.T) {
	address, err := account.Create(accountPassword)
	if err != nil {
		t.Error(err)
	}

	if address == "" {
		t.Error("Address is empty")
	}

	accountAddress = address

	accs := account.List()
	want := 1
	get := len(accs)

	require.Equal(t, want, get)

	storagePath := filepath.Join(shared.AccsDirPath, accountAddress, shared.StorageDirName)
	configPath := filepath.Join(shared.AccsDirPath, accountAddress, shared.ConfDirName)

	if _, err := os.Stat(storagePath); err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(configPath); err != nil {
		t.Error(err)
	}

	nodeAddress = shared.NodeAddr
}

func TestLoginAccountWithCorrectAddressAndPassword(t *testing.T) {
	account, err := account.Login(accountAddress, accountPassword)
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, accountAddress, account.Address.String())
}

func TestLoginAccountWithInvalidPassword(t *testing.T) {
	_, err := account.Login(accountAddress, "invalid")
	want := ErrorInvalidPassword

	require.EqualError(t, want, err.Error())
}

func TestLoginAccountWithUnknownAddress(t *testing.T) {
	unknownAddress := "accountAddress"
	_, err := account.Login(unknownAddress, accountPassword)
	want := errors.New("Account Not Found Error: cannot find account for " + unknownAddress)

	require.EqualError(t, want, err.Error())
}

func TestCheckRightPassword(t *testing.T) {
	err := account.CheckPassword(accountPassword, accountAddress)
	if err != nil {
		t.Error(err)
	}
}

func TestCheckInvalidPassword(t *testing.T) {
	err := account.CheckPassword("accountPassword", accountAddress)
	want := ErrorInvalidPassword
	require.EqualError(t, want, err.Error())
}

func TestCreateConfig(t *testing.T) {
	storageLimit := "1"
	config, err := createConfigForTests(accountAddress, accountPassword, ip, storageLimit, configPort)
	if err != nil {
		t.Error(err)
	}

	configPath := filepath.Join(shared.AccsDirPath, accountAddress, shared.ConfDirName, configName)
	if _, err := os.Stat(configPath); err != nil {
		t.Error(err)
	}

	configAddress = config.Address

	require.Equal(t, configPort, config.HTTPPort)
	require.Equal(t, configStorageLimit, config.StorageLimit)
	require.Equal(t, int64(0), config.UsedStorageSpace)
}

func createConfigForTests(address, password, ipAddress, storageLimit, port string) (*config.SecondaryNodeConfig, error) {
	dFileConf := &config.SecondaryNodeConfig{Address: address}
	pathToConfig := filepath.Join(shared.AccsDirPath, address, shared.ConfDirName)
	regNum := regexp.MustCompile(("[0-9]+"))

	availableSpace := shared.GetAvailableSpace(pathToConfig)
	space := storageLimit

	match := regNum.MatchString(space)

	if !match {
		return nil, fmt.Errorf("storage limit is incorrect")
	}

	intSpace, err := strconv.Atoi(space)
	if err != nil {
		return nil, fmt.Errorf("storage limit failed parse")
	}

	if intSpace < 0 || intSpace >= availableSpace {
		return nil, fmt.Errorf("storage limit is incorrect")
	}

	dFileConf.StorageLimit = intSpace

	regIp := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)

	var splitIPAddr []string
	ipAddr := ipAddress
	match = regIp.MatchString(ipAddr)

	if !match {
		return nil, fmt.Errorf("ip is incorrect, please try again")
	}

	splitIPAddr = strings.Split(ipAddr, ".")

	if fullyReservedIPs[splitIPAddr[0]] {
		return nil, errors.New("Address" + ipAddr + "can't be used as a public ip address")
	}

	reservedSecAddrPart, isReserved := partiallyReservedIPs[splitIPAddr[0]]

	if isReserved {
		secondAddrPart, err := strconv.Atoi(splitIPAddr[1])
		if err != nil {
			return dFileConf, fmt.Errorf("ip  part is incorrect, please try again")
		}

		if secondAddrPart <= reservedSecAddrPart {
			return nil, errors.New("Address" + ipAddr + "can't be used as a public ip address")
		}
	}

	regPort := regexp.MustCompile("[0-9]+|")

	httpPort := port

	if httpPort == "" {
		dFileConf.HTTPPort = fmt.Sprint(55050)
	} else {
		match = regPort.MatchString(httpPort)
		if !match {
			return nil, fmt.Errorf("port is incorrect, please try again")
		}

		intHttpPort, err := strconv.Atoi(httpPort)
		if err != nil {
			return nil, fmt.Errorf("port is incorrect, please try again")
		}

		if intHttpPort < 49152 || intHttpPort > 65535 {
			return nil, fmt.Errorf("port is incorrect, please try again")
		}

		dFileConf.HTTPPort = fmt.Sprint(intHttpPort)
	}

	//TODO add blockchain request
	/*
		err := blockchainprovider.RegisterNode(address, password, splittedAddr, dFileConf.HTTPPort)
		if err != nil {
			log.Fatal("Couldn't register node in network")
		}
	*/

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
