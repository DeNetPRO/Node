package account

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"git.denetwork.xyz/dfile/dfile-secondary-node/config"
	"git.denetwork.xyz/dfile/dfile-secondary-node/encryption"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/term"
)

var (
	IpAddr string
)

//List go to the folder ~/dfile-node/accounts and return all accounts addresses in string format
func List() []string {
	var blockchainAccounts []string

	ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	etherAccounts := ks.Accounts()

	blockchainAccounts = make([]string, 0)

	for _, a := range etherAccounts {
		blockchainAccounts = append(blockchainAccounts, a.Address.String())
	}

	return blockchainAccounts
}

// Create account and keystore file with encryption with password
func Create(password string) (string, config.SecondaryNodeConfig, error) {
	const logLoc = "account.Create->"
	var nodeConf config.SecondaryNodeConfig

	err := shared.CreateIfNotExistAccDirs()
	if err != nil {
		return "", nodeConf, logger.CreateDetails(logLoc, err)
	}

	ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)

	etherAccount, err := ks.NewAccount(password)
	if err != nil {
		return "", nodeConf, logger.CreateDetails(logLoc, err)
	}

	nodeConf, err = initAccount(ks, &etherAccount, password)
	if err != nil {
		return "", nodeConf, logger.CreateDetails(logLoc, err)
	}

	return etherAccount.Address.String(), nodeConf, nil
}

//Import account with private key
func Import() (string, config.SecondaryNodeConfig, error) {
	const logLoc = "account.Import->"
	var nodeConfig config.SecondaryNodeConfig

	fmt.Println("Please enter private key of the account you want to import:")

	privKey, err := shared.ReadFromConsole()
	if err != nil {
		return "", nodeConfig, logger.CreateDetails(logLoc, err)
	}

	ecdsaPrivKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return "", nodeConfig, logger.CreateDetails(logLoc, err)
	}

	fmt.Println("Please enter your password:")

	var originalPassword string

	for {
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", nodeConfig, logger.CreateDetails(logLoc, err)
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
		return "", nodeConfig, logger.CreateDetails(logLoc, err)
	}

	ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)

	etherAccount, err := ks.ImportECDSA(ecdsaPrivKey, password)
	if err != nil {
		fmt.Println(err)
		return "", nodeConfig, logger.CreateDetails(logLoc, err)
	}

	nodeConfig, err = initAccount(ks, &etherAccount, password)
	if err != nil {
		return "", nodeConfig, logger.CreateDetails(logLoc, err)
	}

	return etherAccount.Address.String(), nodeConfig, nil
}

//Login to account with password
func Login(blockchainAccountString, password string) (*accounts.Account, error) {
	const logLoc = "account.Login->"
	ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	etherAccounts := ks.Accounts()

	var account *accounts.Account

	for _, a := range etherAccounts {
		if blockchainAccountString == a.Address.String() {
			account = &a
			break
		}
	}

	if account == nil {
		err := errors.New("Account Not Found Error: cannot find account for " + blockchainAccountString)
		return nil, logger.CreateDetails(logLoc, err)
	}

	keyJson, err := ks.Export(*account, password, password)
	if err != nil {
		fmt.Println("Wrong password")
		return nil, logger.CreateDetails(logLoc, err)
	}

	key, err := keystore.DecryptKey(keyJson, password)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	shared.NodeAddr = account.Address

	macAddr, err := encryption.GetDeviceMacAddr()
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	encrForKey := sha256.Sum256([]byte(macAddr))
	encryptedKey, err := encryption.EncryptAES(encrForKey[:], key.PrivateKey.D.Bytes())
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	encryption.PrivateKey = encryptedKey

	return account, nil
}

//CheckPassword is go to KeyStore and try to decrypt store with incoming account and password.
//Return nil if password is valid
func CheckPassword(password, address string) error {
	const logLoc = "account.CheckPassword->"
	ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := utils.MakeAddress(ks, address)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}
	key, err := ks.Export(acc, password, password)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}
	_, err = keystore.DecryptKey(key, password)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}
	return nil
}

//ValidateUser provides entry into the node.
//It gives 3 attemps to enter into the choosen account
func ValidateUser() (*accounts.Account, string, error) {
	const logLoc = "account.ValidateUser->"
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
				return nil, "", logger.CreateDetails(logLoc, err)
			}

			accountNumber, err := strconv.Atoi(number)
			if err != nil {
				return nil, "", logger.CreateDetails(logLoc, err)
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
			return nil, "", logger.CreateDetails(logLoc, err)
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
			logger.CreateDetails(logLoc, err)
			continue
		}

		loggedIn = true
		break
	}

	if !loggedIn {
		return nil, "", logger.CreateDetails(logLoc, errors.New("couldn't log in in 3 attempts"))
	}

	return etherAccount, password, nil
}

//initAccount creates the necessary environment for the account.
//It creates dir "storage" and config file. And also encrypts private key
func initAccount(ks *keystore.KeyStore, account *accounts.Account, password string) (config.SecondaryNodeConfig, error) {
	const logLoc = "account.initAccount->"
	var nodeConf config.SecondaryNodeConfig

	addressString := account.Address.String()

	err := os.MkdirAll(filepath.Join(paths.AccsDirPath, addressString, paths.StorageDirName), 0700)
	if err != nil {
		return nodeConf, logger.CreateDetails(logLoc, err)
	}

	err = os.MkdirAll(filepath.Join(paths.AccsDirPath, addressString, paths.ConfDirName), 0700)
	if err != nil {
		return nodeConf, logger.CreateDetails(logLoc, err)
	}

	keyJson, err := ks.Export(*account, password, password)
	if err != nil {
		fmt.Println("Wrong password")
		return nodeConf, logger.CreateDetails(logLoc, err)
	}

	key, err := keystore.DecryptKey(keyJson, password)
	if err != nil {
		return nodeConf, logger.CreateDetails(logLoc, err)
	}

	shared.NodeAddr = account.Address

	macAddr, err := encryption.GetDeviceMacAddr()
	if err != nil {
		return nodeConf, logger.CreateDetails(logLoc, err)
	}

	encrForKey := sha256.Sum256([]byte(macAddr))
	encryptedKey, err := encryption.EncryptAES(encrForKey[:], key.PrivateKey.D.Bytes())
	if err != nil {
		return nodeConf, logger.CreateDetails(logLoc, err)
	}

	encryption.PrivateKey = encryptedKey

	nodeConf, err = config.Create(addressString, password)
	if err != nil {
		return nodeConf, logger.CreateDetails(logLoc, err)
	}

	return nodeConf, nil
}
