package account

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"testing"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/config"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/encryption"
	nodeFile "git.denetwork.xyz/DeNet/dfile-secondary-node/node_file"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/shared"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/stretchr/testify/require"
)

var (
	WorkDir     = "tmp"
	AccountsDir = "accounts"
)

func TestMain(m *testing.M) {
	shared.TestModeOn()
	defer shared.TestModeOff()

	os.RemoveAll(WorkDir)

	err := os.Mkdir(WorkDir, 0777)
	if err != nil {
		log.Fatal(err)
	}

	paths.WorkDirName = WorkDir
	paths.WorkDirPath = WorkDir
	paths.AccsDirPath = filepath.Join(WorkDir, AccountsDir)

	exitVal := m.Run()

	err = os.RemoveAll(WorkDir)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitVal)
}

func TestListAfterInitMustBeEmpty(t *testing.T) {
	list := List()

	want := 0
	got := len(list)

	require.Equal(t, want, got)
}

func TestListMustNotBeEmptyIfAccountsExists(t *testing.T) {
	for i := 0; i < 5; i++ {
		ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.LightScryptN, keystore.LightScryptP)

		_, err := ks.NewAccount("test")
		if err != nil {
			t.Error(err)
		}
	}

	list := List()

	want := 5
	got := len(list)

	require.Equal(t, want, got)
}

func TestAccountCreate(t *testing.T) {
	address, config, err := Create("test")
	if err != nil {
		t.Error(err)
	}

	require.NotEmpty(t, address, "account address is empty")
	require.NotEmpty(t, config, "config didn't create")
	require.Equal(t, shared.TestIP, config.IpAddress, "ip address is incorrect, want: ", shared.TestIP, " got: ", config.IpAddress)
	require.Equal(t, shared.TestPort, config.HTTPPort, "port is incorrect, want: ", shared.TestPort, " got: ", config.HTTPPort)
	require.Equal(t, shared.TestNetwork, config.Network, "network is incorrect, want: ", shared.TestNetwork, " got: ", config.Network)
	require.Equal(t, shared.TestStorageLimit, config.StorageLimit, "storage limit is incorrect, want: ", shared.TestStorageLimit, " got: ", config.StorageLimit)
	require.Empty(t, config.UsedStorageSpace, "used storage space must be 0 instead ", config.UsedStorageSpace)
	require.NotEmpty(t, encryption.EncryptedPK, "account private key is empty")
	require.Equal(t, shared.NodeAddr.String(), address, "node address don't equal")

	list := List()
	require.Contains(t, list, address, "account address is not in the account list")
}

func TestAccountLoginWithCurrectPassword(t *testing.T) {
	password := "test"
	address, createdConfig, err := Create(password)
	if err != nil {
		t.Error(err)
	}

	account, err := Login(address, password)
	if err != nil {
		t.Error(err)
	}

	require.NotEmpty(t, account, "account is empty after login")
	require.Equal(t, address, account.Address.String(), "account address is different, want: ", address, " got: ", account.Address.String())

	pathToConfigDir := filepath.Join(paths.AccsDirPath, account.Address.String(), paths.ConfDirName)
	pathToConfigFile := filepath.Join(pathToConfigDir, paths.ConfFileName)

	var nodeConfig config.NodeConfig
	confFile, fileBytes, err := nodeFile.Read(pathToConfigFile)
	if err != nil {
		t.Error(err)
	}

	defer confFile.Close()

	err = json.Unmarshal(fileBytes, &nodeConfig)
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, createdConfig, nodeConfig, "configs are different")
}

func TestAccountLoginWithWrongPassword(t *testing.T) {
	password := "test"
	address, _, err := Create(password)
	if err != nil {
		t.Error(err)
	}

	wrongPassword := "wrong"
	account, err := Login(address, wrongPassword)

	require.Empty(t, account, "account must be empty after fail login")
	require.NotEmpty(t, err, "must be error instead ", err)
}

func TestAccountLoginWithEmptyAddress(t *testing.T) {
	password := "test"
	_, _, err := Create(password)
	if err != nil {
		t.Error(err)
	}

	account, err := Login("", password)
	require.Empty(t, account, "account must be empty after fail login")
	require.NotEmpty(t, err, "must be error instead ", err)
}

func TestAccountImport(t *testing.T) {
	address, config, err := Import()
	if err != nil {
		t.Error(err)
	}

	require.NotEmpty(t, address, "account address is empty")
	require.NotEmpty(t, config, "config didn't create")
	require.Equal(t, shared.TestIP, config.IpAddress, "ip address is incorrect, want: ", shared.TestIP, " got: ", config.IpAddress)
	require.Equal(t, shared.TestPort, config.HTTPPort, "port is incorrect, want: ", shared.TestPort, " got: ", config.HTTPPort)
	require.Equal(t, shared.TestNetwork, config.Network, "network is incorrect, want: ", shared.TestNetwork, " got: ", config.Network)
	require.Equal(t, shared.TestStorageLimit, config.StorageLimit, "storage limit is incorrect, want: ", shared.TestStorageLimit, " got: ", config.StorageLimit)
	require.Empty(t, config.UsedStorageSpace, "used storage space must be 0 instead ", config.UsedStorageSpace)
	require.NotEmpty(t, encryption.EncryptedPK, "account private key is empty")
	require.Equal(t, shared.NodeAddr.String(), address, "node address don't equal")

	list := List()
	require.Contains(t, list, address, "account address is not in the account list")
}

func TestCheckCorrectPassword(t *testing.T) {
	password := "test"
	address, _, err := Create(password)
	if err != nil {
		t.Error(err)
	}

	err = CheckPassword(password, address)
	require.Empty(t, err, "error must be nil instead ", err)
}

func TestCheckWrongPassword(t *testing.T) {
	password := "test"
	address, _, err := Create(password)
	if err != nil {
		t.Error(err)
	}

	wrongPassword := "wrong"
	err = CheckPassword(wrongPassword, address)
	require.NotEmpty(t, err, "error must not be empty")
}
