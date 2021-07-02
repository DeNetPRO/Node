package test

import (
	"dfile-secondary-node/account"
	"dfile-secondary-node/config"
	"dfile-secondary-node/paths"
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

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/stretchr/testify/require"
)

var (
	accountPassword      = "123"
	accountAddress       string
	nodeAddress          []byte
	ErrorInvalidPassword = errors.New("could not decrypt key with given password")
	configPort           = "55051"
	configStorageLimit   = "1"
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
	address, config, err := accountCreateTest(accountPassword, ip, configStorageLimit, configPort)
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

	storagePath := filepath.Join(paths.AccsDirPath, accountAddress, paths.StorageDirName)
	configPath := filepath.Join(paths.AccsDirPath, accountAddress, paths.ConfDirName)

	if _, err := os.Stat(storagePath); err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(configPath); err != nil {
		t.Error(err)
	}

	intConfigStorageLimit, _ := strconv.Atoi(configStorageLimit)
	if config.Address != address || !config.AgreeSendLogs || config.HTTPPort != configPort || config.IpAddress != ip || config.StorageLimit != intConfigStorageLimit {
		t.Error("config is invalid")
	}

	nodeAddress = shared.NodeAddr
}

func TestLoginAccountWithCorrectAddressAndPassword(t *testing.T) {
	account, err := testLogin(accountAddress, accountPassword)
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, accountAddress, account.Address.String())
}

func TestLoginAccountWithInvalidPassword(t *testing.T) {
	_, err := testLogin(accountAddress, "invalid")
	want := ErrorInvalidPassword

	require.EqualError(t, want, err.Error())
}

func TestLoginAccountWithUnknownAddress(t *testing.T) {
	unknownAddress := "accountAddress"
	_, err := testLogin(unknownAddress, accountPassword)
	want := errors.New("Account Not Found Error: cannot find account for " + unknownAddress)

	require.EqualError(t, want, err.Error())
}

func TestCheckRightPassword(t *testing.T) {
	err := account.CheckPassword(accountPassword, accountAddress)
	if err != nil {
		t.Error(err)
	}
}

func accountCreateTest(password, ipAddress, storageLimit, port string) (string, *config.SecondaryNodeConfig, error) {
	var nodeConf *config.SecondaryNodeConfig
	err := shared.CreateIfNotExistAccDirs()
	if err != nil {
		return "", nil, err
	}

	ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)

	etherAccount, err := ks.NewAccount(password)
	if err != nil {
		return "", nil, err
	}

	nodeConf, err = initTestAccount(&etherAccount, password, ipAddress, storageLimit, port)
	if err != nil {
		return "", nodeConf, err
	}

	return etherAccount.Address.String(), nodeConf, nil
}

func createConfigForTests(address, password, ipAddress, storageLimit, port string) (*config.SecondaryNodeConfig, error) {
	dFileConf := &config.SecondaryNodeConfig{Address: address, AgreeSendLogs: true}
	pathToConfig := filepath.Join(paths.AccsDirPath, address, paths.ConfDirName)
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

	dFileConf.IpAddress = ipAddr

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

func initTestAccount(account *accounts.Account, password, ipAddress, storageLimit, port string) (*config.SecondaryNodeConfig, error) {
	nodeConf := &config.SecondaryNodeConfig{}
	addressString := account.Address.String()

	err := os.MkdirAll(filepath.Join(paths.AccsDirPath, addressString, paths.StorageDirName), 0700)
	if err != nil {
		return nodeConf, err
	}

	err = os.MkdirAll(filepath.Join(paths.AccsDirPath, addressString, paths.ConfDirName), 0700)
	if err != nil {
		return nodeConf, err
	}

	encryptedAddr, err := shared.EncryptNodeAddr(account.Address)
	if err != nil {
		return nodeConf, err
	}

	shared.NodeAddr = encryptedAddr

	nodeConf, err = createConfigForTests(account.Address.String(), password, ipAddress, storageLimit, port)
	if err != nil {
		return nodeConf, err
	}

	return nodeConf, nil
}

func testLogin(blockchainAccountString, password string) (*accounts.Account, error) {
	ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	etherAccounts := ks.Accounts()

	var etherAccount *accounts.Account

	for _, a := range etherAccounts {
		if blockchainAccountString == a.Address.String() {
			etherAccount = &a
			break
		}
	}

	if etherAccount == nil {
		err := errors.New("Account Not Found Error: cannot find account for " + blockchainAccountString)
		return nil, err
	}

	keyJson, err := ks.Export(*etherAccount, password, password)
	if err != nil {
		return nil, err
	}

	key, err := keystore.DecryptKey(keyJson, password)
	if err != nil {
		return nil, err
	}

	encryptedAddr, err := shared.EncryptNodeAddr(key.Address)
	if err != nil {
		return nil, err
	}

	shared.NodeAddr = encryptedAddr

	return etherAccount, nil
}
