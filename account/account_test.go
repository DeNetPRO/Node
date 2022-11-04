package account_test

import (
	"log"
	"os"
	"testing"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/account"
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

	err = os.RemoveAll(paths.List().WorkDir)
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

	_, err = os.Stat(paths.List().AccsDir)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(paths.List().Storages[0])
	if err != nil {
		t.Fatal(err)
	}

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

func TestImportAccount(t *testing.T) { //TODO add test checks
	accountAddress, _, err := account.Import()
	if err != nil {
		t.Fatal(err)
	}

	if accountAddress == "" {
		t.Errorf("import account address must not be empty")
	}

}
