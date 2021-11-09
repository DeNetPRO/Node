package account_test

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"testing"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/account"
	blckChain "git.denetwork.xyz/DeNet/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/config"
	tstpkg "git.denetwork.xyz/DeNet/dfile-secondary-node/tst_pkg"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"

	"github.com/stretchr/testify/require"
)

var (
	testPasswd  = "testPasswd"
	testAccAddr string
)

func TestMain(m *testing.M) {
	tstpkg.TestModeOn()
	defer tstpkg.TestModeOff()

	err := paths.Init()
	if err != nil {
		log.Fatal(err)
	}

	exitVal := m.Run()

	err = os.RemoveAll(paths.WorkDirPath)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitVal)
}

func TestListAccsBeforeCreation(t *testing.T) {
	require.Equal(t, 0, len(account.List()))
}

func TestAccCreate(t *testing.T) {
	accountAddress, _, err := account.Create(testPasswd)
	if err != nil {
		t.Fatal(err)
	}

	testAccAddr = accountAddress

	_, err = os.Stat(paths.AccsDirPath)
	if err != nil {
		t.Fatal(err)
	}

	pathToStorage := filepath.Join(paths.StoragePaths[0], blckChain.CurrentNetwork)

	_, err = os.Stat(pathToStorage)
	if err != nil {
		t.Fatal(err)
	}

	pathToConfigFile := filepath.Join(paths.ConfigDirPath, paths.ConfFileName)

	_, err = os.Stat(pathToConfigFile)
	if err != nil {
		t.Fatal(err)
	}

	confFileBytes, err := os.ReadFile(pathToConfigFile)
	if err != nil {
		t.Fatal(err)
	}

	var accConfig config.NodeConfig

	err = json.Unmarshal(confFileBytes, &accConfig)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, config.TestConfig, accConfig)

}

func TestListAccExists(t *testing.T) {
	if !account.AccExists(account.List(), testAccAddr) {
		t.Fatal("account does not exist")
	}
}

func TestAccountLogin(t *testing.T) {

	t.Run("correct password", func(t *testing.T) {
		_, err := account.Login(testAccAddr, testPasswd)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("incorrect password", func(t *testing.T) {
		account, err := account.Login(testAccAddr, "wrongPassword")

		require.Empty(t, account, "account value must be empty when password is wrong")
		require.NotEmpty(t, err, "error value must not be empty when password is wrong", err)
	})

	t.Run("account value is empty", func(t *testing.T) {
		account, err := account.Login("", testPasswd)

		require.Empty(t, account, "account value must be empty")
		require.NotEmpty(t, err, "error value must not be empty when account value is empty", err)
	})

}

func TestImportAccount(t *testing.T) {
	accountAddress, accConfig, err := account.Import()
	if err != nil {
		t.Fatal(err)
	}

	if accountAddress == "" {
		t.Errorf("import account address must not be empty")
	}

	require.Equal(t, config.TestConfig, accConfig)

}
