package test

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"git.denetwork.xyz/dfile/dfile-secondary-node/account"
	"git.denetwork.xyz/dfile/dfile-secondary-node/config"
	dnetsignature "git.denetwork.xyz/dfile/dfile-secondary-node/dnet_signature"
	"git.denetwork.xyz/dfile/dfile-secondary-node/encryption"
	meminfo "git.denetwork.xyz/dfile/dfile-secondary-node/mem_info"
	nodefile "git.denetwork.xyz/dfile/dfile-secondary-node/node_file"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

var (
	accountPassword      = "123"
	accountAddress       string
	nodeAddress          []byte
	ErrorInvalidPassword = errors.New(" could not decrypt key with given password")
	configPath           string
	storagePath          string
)

func TestEmptyAccountListBeforeCreating(t *testing.T) {
	accs := account.List()
	want := 0
	get := len(accs)

	require.Equal(t, want, get)
}

func TestSetIpAddrWhenCreateConfig(t *testing.T) {
	get := config.SecondaryNodeConfig{}
	ip, err := config.SetIpAddr(&get, config.CreateStatus)
	if err != nil {
		t.Error(err)
	}

	if len(ip) != 4 {
		t.Errorf("len of ip must be 4 instead of %v", len(ip))
	}

	want := config.SecondaryNodeConfig{
		IpAddress: shared.TestAddress,
	}

	require.Equal(t, want, get)
}

func TestSetPortWhenCreateConfig(t *testing.T) {
	get := config.SecondaryNodeConfig{}
	err := config.SetPort(&get, config.CreateStatus)
	if err != nil {
		t.Error(err)
	}

	want := config.SecondaryNodeConfig{
		HTTPPort: shared.TestPort,
	}

	require.Equal(t, want, get)
}

func TestSetStorageLimitWhenCreateConfig(t *testing.T) {
	get := config.SecondaryNodeConfig{}
	err := config.SetStorageLimit("", config.CreateStatus, &get)
	if err != nil {
		t.Error(err)
	}

	want := config.SecondaryNodeConfig{
		StorageLimit: shared.TestLimit,
	}

	require.Equal(t, want, get)
}

func TestCreateAccount(t *testing.T) {
	address, config, err := account.Create(accountPassword)
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

	storagePath = filepath.Join(paths.AccsDirPath, accountAddress, paths.StorageDirName)
	configPath = filepath.Join(paths.AccsDirPath, accountAddress, paths.ConfDirName, paths.ConfFileName)

	if _, err := os.Stat(storagePath); err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(configPath); err != nil {
		t.Error(err)
	}

	if config.Address != address || !config.AgreeSendLogs || config.HTTPPort != shared.TestPort || config.IpAddress != shared.TestAddress || config.StorageLimit != shared.TestLimit {
		t.Error("config is invalid")
	}

	nodeAddress = shared.NodeAddr.Bytes()
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

	splitErr := strings.Split(err.Error(), "->")

	require.EqualError(t, want, splitErr[len(splitErr)-1])
}

func TestLoginAccountWithUnknownAddress(t *testing.T) {
	unknownAddress := "accountAddress"
	_, err := account.Login(unknownAddress, accountPassword)
	want := errors.New(" accountAddress address is not found")
	splitErr := strings.Split(err.Error(), "->")

	require.EqualError(t, want, splitErr[len(splitErr)-1])
}

func TestCheckRightPassword(t *testing.T) {
	err := account.CheckPassword(accountPassword, accountAddress)
	if err != nil {
		t.Error(err)
	}
}

func TestImportAccount(t *testing.T) {
	accountAddress, c, err := account.Import()
	if err != nil {
		t.Error(err)
	}

	if accountAddress == "" {
		t.Errorf("import account address must not to be empty")
	}

	wantConfig := config.SecondaryNodeConfig{
		Address:       accountAddress,
		HTTPPort:      shared.TestPort,
		StorageLimit:  shared.TestLimit,
		IpAddress:     shared.TestAddress,
		AgreeSendLogs: true,
	}

	require.Equal(t, wantConfig, c)
}

func TestCheckSignature(t *testing.T) {
	macAddress, err := encryption.GetDeviceMacAddr()
	if err != nil {
		t.Error(err)
	}

	encrForKey := sha256.Sum256([]byte(macAddress))
	privateKeyBytes, err := encryption.DecryptAES(encrForKey[:], encryption.PrivateKey)
	if err != nil {
		t.Error(err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		t.Error(err)
	}

	data := make([]byte, 100)
	rand.Seed(time.Now().Unix())
	rand.Read(data)

	hashData := sha256.Sum256(data)

	signedData, err := crypto.Sign(hashData[:], privateKey)
	if err != nil {
		t.Error(err)
	}

	err = dnetsignature.Check(accountAddress, signedData, hashData)
	if err != nil {
		t.Error(encrForKey)
	}
}

func TestRestoreNodeMemory(t *testing.T) {
	fileSize := 1024 * 1024

	confFile, nodeConfig, err := getConfig()
	if err != nil {
		t.Error(err)
	}

	want := nodeConfig.UsedStorageSpace

	nodeConfig.UsedStorageSpace += int64(fileSize)

	err = config.Save(confFile, *nodeConfig)
	if err != nil {
		confFile.Close()
		t.Error(err)
	}

	confFile.Close()

	meminfo.Restore(configPath, fileSize)

	confFile, nodeConfig, err = getConfig()
	if err != nil {
		t.Error(err)
	}

	confFile.Close()

	require.Equal(t, want, nodeConfig.UsedStorageSpace)
}

func getConfig() (*os.File, *config.SecondaryNodeConfig, error) {
	confFile, fileBytes, err := nodefile.Read(configPath)
	if err != nil {
		return nil, nil, err
	}

	var nodeConfig *config.SecondaryNodeConfig

	err = json.Unmarshal(fileBytes, &nodeConfig)
	if err != nil {
		confFile.Close()
		return nil, nil, err
	}

	return confFile, nodeConfig, nil
}
