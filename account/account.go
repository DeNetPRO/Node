package account

import (
	"dfile-secondary-node/config"
	"dfile-secondary-node/encryption"
	"dfile-secondary-node/logger"
	"dfile-secondary-node/paths"
	"dfile-secondary-node/shared"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

	ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	etherAccounts := ks.Accounts()

	blockchainAccounts = make([]string, 0, 1)

	for _, a := range etherAccounts {
		blockchainAccounts = append(blockchainAccounts, a.Address.String())
	}

	return blockchainAccounts
}

// CreateAccount creates account and keystore file with encryption with password
func Create(password string) (string, config.SecondaryNodeConfig, error) {
	const logInfo = "account.Create->"
	var nodeConf config.SecondaryNodeConfig

	err := shared.CreateIfNotExistAccDirs()
	if err != nil {
		return "", nodeConf, logger.CreateDetails(logInfo, err)
	}

	ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)

	etherAccount, err := ks.NewAccount(password)
	if err != nil {
		return "", nodeConf, logger.CreateDetails(logInfo, err)
	}

	nodeConf, err = initAccount(&etherAccount, password)
	if err != nil {
		return "", nodeConf, logger.CreateDetails(logInfo, err)
	}

	return etherAccount.Address.String(), nodeConf, nil
}

func Import() (string, config.SecondaryNodeConfig, error) {
	const logInfo = "account.Import->"
	var nodeConfig config.SecondaryNodeConfig

	fmt.Println("Please enter private key of the account you want to import:")

	privKey, err := shared.ReadFromConsole()
	if err != nil {
		return "", nodeConfig, logger.CreateDetails(logInfo, err)
	}

	ecdsaPrivKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return "", nodeConfig, logger.CreateDetails(logInfo, err)
	}

	fmt.Println("Please enter your password:")

	var originalPassword string

	for {
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", nodeConfig, logger.CreateDetails(logInfo, err)
		}

		originalPassword = string(bytePassword)
		if strings.Trim(originalPassword, " ") == "" {
			fmt.Println("Empty string can't be used as a password. Please try again")
			continue
		}

		break
	}

	password := shared.GetHashPassword(originalPassword)
	originalPassword = ""

	err = shared.CreateIfNotExistAccDirs()
	if err != nil {
		return "", nodeConfig, logger.CreateDetails(logInfo, err)
	}

	ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)

	etherAccount, err := ks.ImportECDSA(ecdsaPrivKey, password)
	if err != nil {
		fmt.Println(err)
		return "", nodeConfig, logger.CreateDetails(logInfo, err)
	}

	nodeConfig, err = initAccount(&etherAccount, password)
	if err != nil {
		return "", nodeConfig, logger.CreateDetails(logInfo, err)
	}

	return etherAccount.Address.String(), nodeConfig, nil
}

//LoadAccount load in memory keystore file and decrypt it for further use
func Login(blockchainAccountString, password string) (*accounts.Account, error) {
	const logInfo = "account.Login->"
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
		return nil, logger.CreateDetails(logInfo, err)
	}

	keyJson, err := ks.Export(*etherAccount, password, password)
	if err != nil {
		fmt.Println("Wrong password")
		return nil, logger.CreateDetails(logInfo, err)
	}

	key, err := keystore.DecryptKey(keyJson, password)
	if err != nil {
		return nil, logger.CreateDetails(logInfo, err)
	}

	encryptedAddr, err := encryption.EncryptNodeAddr(key.Address)
	if err != nil {
		return nil, logger.CreateDetails(logInfo, err)
	}

	encryption.NodeAddr = encryptedAddr

	return etherAccount, nil
}

func CheckPassword(password, address string) error {
	const logInfo = "account.CheckPassword->"
	ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := utils.MakeAddress(ks, address)
	if err != nil {
		return logger.CreateDetails(logInfo, err)
	}
	key, err := ks.Export(acc, password, password)
	if err != nil {
		return logger.CreateDetails(logInfo, err)
	}
	_, err = keystore.DecryptKey(key, password)
	if err != nil {
		return logger.CreateDetails(logInfo, err)
	}
	return nil
}

func ValidateUser() (*accounts.Account, string, error) {
	const logInfo = "account.ValidateUser->"
	var accountAddress, password string
	var etherAccount *accounts.Account

	accounts := List()

	if len(accounts) > 1 {
		fmt.Println("Please choose an account number")
		for i, a := range accounts {
			fmt.Println(i+1, a)
		}
	}

	loggedIn := false

	for i := 0; i < 3; i++ {
		if len(accounts) == 1 {
			accountAddress = accounts[0]
		} else {
			number, err := shared.ReadFromConsole()
			if err != nil {
				return nil, "", logger.CreateDetails(logInfo, err)
			}

			accountNumber, err := strconv.Atoi(number)
			if err != nil {
				return nil, "", logger.CreateDetails(logInfo, err)
			}

			if accountNumber < 1 || accountNumber > len(accounts) {
				fmt.Println("Number is incorrect")
				for i, a := range accounts {
					fmt.Println(i+1, a)
				}
				continue
			}

			accountAddress = accounts[accountNumber-1]
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
			return nil, "", logger.CreateDetails(logInfo, err)
		}

		originalPassword := string(bytePassword)
		if strings.Trim(originalPassword, " ") == "" {
			fmt.Println("Empty string can't be used as a password. Please enter passwords again")
			continue
		}

		password = shared.GetHashPassword(originalPassword)
		originalPassword = ""
		bytePassword = nil

		etherAccount, err = Login(accountAddress, password)
		if err != nil {
			logger.CreateDetails(logInfo, err)
			continue
		}

		loggedIn = true
		break
	}

	if !loggedIn {
		return nil, "", logger.CreateDetails(logInfo, errors.New("couldn't log in in 3 attempts"))
	}

	return etherAccount, password, nil
}

func initAccount(account *accounts.Account, password string) (config.SecondaryNodeConfig, error) {
	const logInfo = "account.initAccount->"
	var nodeConf config.SecondaryNodeConfig

	addressString := account.Address.String()

	err := os.MkdirAll(filepath.Join(paths.AccsDirPath, addressString, paths.StorageDirName), 0700)
	if err != nil {
		return nodeConf, logger.CreateDetails(logInfo, err)
	}

	err = os.MkdirAll(filepath.Join(paths.AccsDirPath, addressString, paths.ConfDirName), 0700)
	if err != nil {
		return nodeConf, logger.CreateDetails(logInfo, err)
	}

	encryptedAddr, err := encryption.EncryptNodeAddr(account.Address)
	if err != nil {
		return nodeConf, logger.CreateDetails(logInfo, err)
	}

	encryption.NodeAddr = encryptedAddr

	nodeConf, err = config.Create(addressString, password)
	if err != nil {
		return nodeConf, logger.CreateDetails(logInfo, err)
	}

	return nodeConf, nil
}
