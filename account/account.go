package account

import (
	"dfile-secondary-node/config"
	"dfile-secondary-node/shared"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/term"
)

//GetAllAccounts go to the folder ~/dfile/accounts and return all accounts addresses in string format
func List() []string {
	var blockchainAccounts []string

	ks := keystore.NewKeyStore(shared.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	etherAccounts := ks.Accounts()

	blockchainAccounts = make([]string, 0, 1)

	for _, a := range etherAccounts {
		blockchainAccounts = append(blockchainAccounts, a.Address.String())
	}

	return blockchainAccounts
}

// CreateAccount creates account and keystore file with encryption with password
func Create(password string) (string, error) {

	shared.CreateIfNotExistAccDirs()

	ks := keystore.NewKeyStore(shared.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)

	etherAccount, err := ks.NewAccount(password)
	if err != nil {
		return "", err
	}

	addressString := etherAccount.Address.String()

	err = os.MkdirAll(filepath.Join(shared.AccsDirPath, addressString, shared.StorageDirName), 0700)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Join(shared.AccsDirPath, addressString, shared.ConfDirName), 0700)
	if err != nil {
		return "", err
	}

	encryptedAddr, err := shared.EncryptNodeAddr(etherAccount.Address)
	if err != nil {
		return "", err
	}

	shared.NodeAddr = encryptedAddr

	return addressString, nil
}

func Import() error {

	fmt.Println("Please enter private key of the account you want to import:")

	privKey, err := shared.ReadFromConsole()
	if err != nil {
		return err
	}

	ecdsaPrivKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return err
	}

	fmt.Println("Please enter new password:")

	var password string

	for {
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}

		password = string(bytePassword)
		if strings.Trim(password, " ") == "" {
			fmt.Println("Empty string can't be used as a password. Please try again")
			continue
		}

		break
	}

	shared.CreateIfNotExistAccDirs()

	ks := keystore.NewKeyStore(shared.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)

	etherAccount, err := ks.ImportECDSA(ecdsaPrivKey, password)
	if err != nil {
		fmt.Println(err)
		return err
	}

	addressString := etherAccount.Address.String()

	err = os.MkdirAll(filepath.Join(shared.AccsDirPath, addressString, shared.StorageDirName), 0700)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(shared.AccsDirPath, addressString, shared.ConfDirName), 0700)
	if err != nil {
		return err
	}

	encryptedAddr, err := shared.EncryptNodeAddr(etherAccount.Address)
	if err != nil {
		return err
	}

	shared.NodeAddr = encryptedAddr

	config.Create(addressString, config.State.Create)

	return nil
}

//LoadAccount load in memory keystore file and decrypt it for further use
func Login(blockchainAccountString, password string) (*accounts.Account, error) {

	ks := keystore.NewKeyStore(shared.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	etherAccounts := ks.Accounts()

	var etherAccount *accounts.Account

	for _, a := range etherAccounts {
		if blockchainAccountString == a.Address.String() {
			etherAccount = &a
			break
		}
	}

	if etherAccount == nil {
		return nil, errors.New("Account Not Found Error: cannot find account for " + blockchainAccountString)
	}

	keyJson, err := ks.Export(*etherAccount, password, password)
	if err != nil {
		fmt.Println("Wrong password")
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

func CheckPassword(password, address string) error {

	ks := keystore.NewKeyStore(shared.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := utils.MakeAddress(ks, address)
	if err != nil {
		return err
	}
	key, err := ks.Export(acc, password, password)
	if err != nil {
		return err
	}
	_, err = keystore.DecryptKey(key, password)
	if err != nil {
		return err
	}
	return nil
}

func ValidateUser() (*accounts.Account, string, error) {
	var accountAddress, password string
	var etherAccount *accounts.Account
	var err error

	accounts := List()

	if len(accounts) > 1 {
		fmt.Println("Please choose an account")
		for i, a := range accounts {
			fmt.Println(i+1, a)
		}
	}

	for {

		if len(accounts) == 1 {
			accountAddress = accounts[0]
		} else {
			accountAddress, err = shared.ReadFromConsole()
			if err != nil {
				return nil, "", err
			}
		}

		addressMatches := shared.ContainsAccount(accounts, accountAddress)

		if !addressMatches {
			fmt.Println("There is no such account address:")
			for i, a := range accounts {
				fmt.Println(i+1, a)
			}
			continue
		}

		fmt.Println("Please enter your password:")

		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return nil, "", err
		}
		password = string(bytePassword)
		if strings.Trim(password, " ") == "" {
			fmt.Println("Empty string can't be used as a password. Please enter passwords again")
			continue
		}

		etherAccount, err = Login(accountAddress, password)
		if err != nil {
			return nil, "", err
		}

		break
	}

	return etherAccount, password, nil
}
