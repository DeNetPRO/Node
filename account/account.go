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
func Create(password string) (string, config.SecondaryNodeConfig, error) {
	const info = "account.Create"
	var nodeConf config.SecondaryNodeConfig

	err := shared.CreateIfNotExistAccDirs()
	if err != nil {
		return "", nodeConf, fmt.Errorf("%s %w", info, err)
	}

	ks := keystore.NewKeyStore(shared.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)

	etherAccount, err := ks.NewAccount(password)
	if err != nil {
		return "", nodeConf, fmt.Errorf("%s %w", info, err)
	}

	nodeConf, err = initAccount(&etherAccount)
	if err != nil {
		return "", nodeConf, fmt.Errorf("%s %w", info, err)
	}

	return etherAccount.Address.String(), nodeConf, nil
}

func Import() (string, config.SecondaryNodeConfig, error) {
	const info = "account.Import"
	var nodeConfig config.SecondaryNodeConfig

	fmt.Println("Please enter private key of the account you want to import:")

	privKey, err := shared.ReadFromConsole()
	if err != nil {
		return "", nodeConfig, fmt.Errorf("%s %w", info, err)
	}

	ecdsaPrivKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return "", nodeConfig, fmt.Errorf("%s %w", info, err)
	}

	fmt.Println("Please enter new password:")

	var password string

	for {
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", nodeConfig, fmt.Errorf("%s %w", info, err)
		}

		password = string(bytePassword)
		if strings.Trim(password, " ") == "" {
			fmt.Println("Empty string can't be used as a password. Please try again")
			continue
		}

		break
	}

	err = shared.CreateIfNotExistAccDirs()
	if err != nil {
		return "", nodeConfig, fmt.Errorf("%s %w", info, err)
	}

	ks := keystore.NewKeyStore(shared.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)

	etherAccount, err := ks.ImportECDSA(ecdsaPrivKey, password)
	if err != nil {
		fmt.Println(err)
		return "", nodeConfig, fmt.Errorf("%s %w", info, err)
	}

	nodeConfig, err = initAccount(&etherAccount)
	if err != nil {
		return "", nodeConfig, fmt.Errorf("%s %w", info, err)
	}

	return etherAccount.Address.String(), nodeConfig, nil
}

//LoadAccount load in memory keystore file and decrypt it for further use
func Login(blockchainAccountString, password string) (*accounts.Account, error) {
	const info = "account.Login"
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
		err := errors.New("Account Not Found Error: cannot find account for " + blockchainAccountString)
		return nil, fmt.Errorf("%s %w", info, err)
	}

	keyJson, err := ks.Export(*etherAccount, password, password)
	if err != nil {
		fmt.Println("Wrong password")
		return nil, fmt.Errorf("%s %w", info, err)
	}

	key, err := keystore.DecryptKey(keyJson, password)
	if err != nil {
		return nil, fmt.Errorf("%s %w", info, err)
	}

	encryptedAddr, err := shared.EncryptNodeAddr(key.Address)
	if err != nil {
		return nil, fmt.Errorf("%s %w", info, err)
	}

	shared.NodeAddr = encryptedAddr

	return etherAccount, nil
}

func CheckPassword(password, address string) error {
	const info = "account.CheckPassword"
	ks := keystore.NewKeyStore(shared.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := utils.MakeAddress(ks, address)
	if err != nil {
		return fmt.Errorf("%s %w", info, err)
	}
	key, err := ks.Export(acc, password, password)
	if err != nil {
		return fmt.Errorf("%s %w", info, err)
	}
	_, err = keystore.DecryptKey(key, password)
	if err != nil {
		return fmt.Errorf("%s %w", info, err)
	}
	return nil
}

func ValidateUser() (*accounts.Account, string, error) {
	const info = "account.ValidateUser"
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
				return nil, "", fmt.Errorf("%s %w", info, err)
			}
		}

		if !shared.ContainsAccount(accounts, accountAddress) {
			fmt.Println("There is no such account address:")
			for i, a := range accounts {
				fmt.Println(i+1, a)
			}
			continue
		}

		fmt.Println("Please enter your password:")

		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return nil, "", fmt.Errorf("%s %w", info, err)
		}
		password = string(bytePassword)
		if strings.Trim(password, " ") == "" {
			fmt.Println("Empty string can't be used as a password. Please enter passwords again")
			continue
		}

		etherAccount, err = Login(accountAddress, password)
		if err != nil {
			return nil, "", fmt.Errorf("%s %w", info, err)
		}

		break
	}

	return etherAccount, password, nil
}

func initAccount(account *accounts.Account) (config.SecondaryNodeConfig, error) {
	const info = "account.initAccount"
	var nodeConf config.SecondaryNodeConfig

	addressString := account.Address.String()

	err := os.MkdirAll(filepath.Join(shared.AccsDirPath, addressString, shared.StorageDirName), 0700)
	if err != nil {
		return nodeConf, fmt.Errorf("%s %w", info, err)
	}

	err = os.MkdirAll(filepath.Join(shared.AccsDirPath, addressString, shared.ConfDirName), 0700)
	if err != nil {
		return nodeConf, fmt.Errorf("%s %w", info, err)
	}

	encryptedAddr, err := shared.EncryptNodeAddr(account.Address)
	if err != nil {
		return nodeConf, fmt.Errorf("%s %w", info, err)
	}

	shared.NodeAddr = encryptedAddr

	nodeConf, err = config.Create(addressString, config.State.Create)
	if err != nil {
		return nodeConf, fmt.Errorf("%s %w", info, err)
	}

	return nodeConf, nil
}
