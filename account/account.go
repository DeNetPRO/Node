package account

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/minio/sha256-simd"

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

//List returns list of user's created/imported wallet adresses, that are used as user accounts.
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

// Create is used for creating a new crypto wallet with keystore file.
func Create(password string) (string, config.SecondaryNodeConfig, error) {
	const location = "account.Create->"
	var nodeConf config.SecondaryNodeConfig

	err := paths.CreateAccDirs()
	if err != nil {
		return "", nodeConf, logger.CreateDetails(location, err)
	}

	ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)

	etherAccount, err := ks.NewAccount(password)
	if err != nil {
		return "", nodeConf, logger.CreateDetails(location, err)
	}

	nodeConf, err = initAccount(ks, &etherAccount, password)
	if err != nil {
		return "", nodeConf, logger.CreateDetails(location, err)
	}

	return etherAccount.Address.String(), nodeConf, nil
}

//Import is used for importing crypto wallet. Private key is needed.
func Import() (string, config.SecondaryNodeConfig, error) {
	const location = "account.Import->"
	var nodeConfig config.SecondaryNodeConfig

	var privKey string
	var err error

	if !shared.TestMode {
		fmt.Println("Please enter private key of the account you want to import:")

		privKey, err = shared.ReadFromConsole()
		if err != nil {
			return "", nodeConfig, logger.CreateDetails(location, err)
		}
	} else {
		privKey = shared.TestPrivateKey
	}

	ecdsaPrivKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return "", nodeConfig, logger.CreateDetails(location, err)
	}

	var password string
	if !shared.TestMode {
		fmt.Println("Please enter your password:")

		var originalPassword string

		for {
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return "", nodeConfig, logger.CreateDetails(location, err)
			}

			originalPassword = string(bytePassword)
			if strings.Trim(originalPassword, " ") == "" {
				fmt.Println("Empty string can't be used as a password. Please try again")
				continue
			}

			break
		}

		password = shared.GetHashPassword(originalPassword)
		originalPassword = ""
	} else {
		password = shared.TestPassword
	}

	err = paths.CreateAccDirs()
	if err != nil {
		return "", nodeConfig, logger.CreateDetails(location, err)
	}

	ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)

	etherAccount, err := ks.ImportECDSA(ecdsaPrivKey, password)
	if err != nil {
		return "", nodeConfig, logger.CreateDetails(location, err)
	}

	nodeConfig, err = initAccount(ks, &etherAccount, password)
	if err != nil {
		return "", nodeConfig, logger.CreateDetails(location, err)
	}

	return etherAccount.Address.String(), nodeConfig, nil
}

//Login checks wallet's address and user's password that was used for crypto wallet creation.
func Login(accountAddress, password string) (*accounts.Account, error) {
	const location = "account.Login->"
	ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	etherAccounts := ks.Accounts()

	var account *accounts.Account

	for _, a := range etherAccounts {
		if accountAddress == a.Address.String() {
			account = &a
			break
		}
	}

	if account == nil {
		err := errors.New(accountAddress + " address is not found")
		return nil, logger.CreateDetails(location, err)
	}

	keyJson, err := ks.Export(*account, password, password)
	if err != nil {
		fmt.Println("Wrong password")
		return nil, logger.CreateDetails(location, err)
	}

	key, err := keystore.DecryptKey(keyJson, password)
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	shared.NodeAddr = account.Address

	macAddr, err := encryption.GetDeviceMacAddr()
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	encrForKey := sha256.Sum256([]byte(macAddr))
	encryptedKey, err := encryption.EncryptAES(encrForKey[:], key.PrivateKey.D.Bytes())
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	encryption.PrivateKey = encryptedKey

	return account, nil
}

//CheckPassword checks crypto wallet's password.
func CheckPassword(password, address string) error {
	const location = "account.CheckPassword->"
	ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := utils.MakeAddress(ks, address)
	if err != nil {
		return logger.CreateDetails(location, err)
	}
	key, err := ks.Export(acc, password, password)
	if err != nil {
		return logger.CreateDetails(location, err)
	}
	_, err = keystore.DecryptKey(key, password)
	if err != nil {
		return logger.CreateDetails(location, err)
	}
	return nil
}

//ValidateUser asks user for password and checks it.
func ValidateUser() (*accounts.Account, string, error) {
	const location = "account.ValidateUser->"
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
				return nil, "", logger.CreateDetails(location, err)
			}

			accountNumber, err := strconv.Atoi(number)
			if err != nil {
				return nil, "", logger.CreateDetails(location, err)
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

		if !accExists(accounts, accountAddress) {
			fmt.Println("There is no such account address:")
			for i, a := range accounts {
				fmt.Println(i+1, a)
			}
			continue
		}

		fmt.Println("Please enter your password:")

		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return nil, "", logger.CreateDetails(location, err)
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
			logger.CreateDetails(location, err)
			continue
		}

		loggedIn = true
		break
	}

	if !loggedIn {
		return nil, "", logger.CreateDetails(location, errors.New("couldn't log in in 3 attempts"))
	}

	return etherAccount, password, nil
}

//InitAccount creates directories and files needed for correct work.
func initAccount(ks *keystore.KeyStore, account *accounts.Account, password string) (config.SecondaryNodeConfig, error) {
	const location = "account.initAccount->"
	var nodeConf config.SecondaryNodeConfig

	addressString := account.Address.String()

	err := os.MkdirAll(filepath.Join(paths.AccsDirPath, addressString, paths.StorageDirName), 0700)
	if err != nil {
		return nodeConf, logger.CreateDetails(location, err)
	}

	err = os.MkdirAll(filepath.Join(paths.AccsDirPath, addressString, paths.ConfDirName), 0700)
	if err != nil {
		return nodeConf, logger.CreateDetails(location, err)
	}

	keyJson, err := ks.Export(*account, password, password)
	if err != nil {
		fmt.Println("Wrong password")
		return nodeConf, logger.CreateDetails(location, err)
	}

	key, err := keystore.DecryptKey(keyJson, password)
	if err != nil {
		return nodeConf, logger.CreateDetails(location, err)
	}

	shared.NodeAddr = account.Address

	macAddr, err := encryption.GetDeviceMacAddr()
	if err != nil {
		return nodeConf, logger.CreateDetails(location, err)
	}

	encrForKey := sha256.Sum256([]byte(macAddr))
	encryptedKey, err := encryption.EncryptAES(encrForKey[:], key.PrivateKey.D.Bytes())
	if err != nil {
		return nodeConf, logger.CreateDetails(location, err)
	}

	encryption.PrivateKey = encryptedKey

	nodeConf, err = config.Create(addressString, password)
	if err != nil {
		return nodeConf, logger.CreateDetails(location, err)
	}

	return nodeConf, nil
}

// ====================================================================================

func accExists(accounts []string, address string) bool {
	for _, a := range accounts {
		if a == address {
			return true
		}
	}
	return false
}
