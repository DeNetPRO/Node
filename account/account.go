package account

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	blckChain "git.denetwork.xyz/DeNet/dfile-secondary-node/blockchain_provider"
	nodeFile "git.denetwork.xyz/DeNet/dfile-secondary-node/node_file"
	tstpkg "git.denetwork.xyz/DeNet/dfile-secondary-node/tst_pkg"

	termEmul "git.denetwork.xyz/DeNet/dfile-secondary-node/term_emul"
	"github.com/howeyc/gopass"
	"github.com/minio/sha256-simd"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/config"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/encryption"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/hash"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/shared"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

//List returns list of user's created/imported wallet adresses, that are used as user accounts.
func List() []string {
	var blockchainAccounts []string

	scryptN, scryptP := encryption.GetScryptParams()

	ks := keystore.NewKeyStore(paths.AccsDirPath, scryptN, scryptP)
	nodeAccounts := ks.Accounts()

	blockchainAccounts = make([]string, 0)

	for _, a := range nodeAccounts {
		blockchainAccounts = append(blockchainAccounts, a.Address.String())
	}

	return blockchainAccounts
}

// Create is used for creating a new crypto wallet with keystore file.
func Create(password string) (string, config.NodeConfig, error) {
	const location = "account.Create->"
	var nodeConf config.NodeConfig

	err := paths.CreateAccDirs()
	if err != nil {
		return "", nodeConf, logger.CreateDetails(location, err)
	}

	scryptN, scryptP := encryption.GetScryptParams()

	ks := keystore.NewKeyStore(paths.AccsDirPath, scryptN, scryptP)

	nodeAccount, err := ks.NewAccount(password)
	if err != nil {
		return "", nodeConf, logger.CreateDetails(location, err)
	}

	nodeConf, err = initAccount(ks, &nodeAccount, password)
	if err != nil {
		return "", nodeConf, logger.CreateDetails(location, err)
	}

	return nodeAccount.Address.String(), nodeConf, nil
}

//Import is used for importing crypto wallet. Private key is needed.
func Import() (string, config.NodeConfig, error) {
	const location = "account.Import->"
	var nodeConfig config.NodeConfig

	var privKey string
	var originalPassword string
	var err error

	if tstpkg.TestMode {
		privKey = tstpkg.TestPrivateKey
		originalPassword = tstpkg.TestPassword
	} else {
		fmt.Println("Please enter private key of the account you want to import:")

		bytesPrivKey, err := gopass.GetPasswdMasked()
		if err != nil {
			return "", nodeConfig, logger.CreateDetails(location, err)
		}

		privKey = string(bytesPrivKey)

		fmt.Println("Please enter your password:")

		for {
			bytePassword, err := gopass.GetPasswdMasked()
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
	}

	password := hash.Password(originalPassword)
	originalPassword = ""

	err = paths.CreateAccDirs()
	if err != nil {
		return "", nodeConfig, logger.CreateDetails(location, err)
	}

	scryptN, scryptP := encryption.GetScryptParams()

	ks := keystore.NewKeyStore(paths.AccsDirPath, scryptN, scryptP)

	ecdsaPrivKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return "", nodeConfig, logger.CreateDetails(location, err)
	}

	nodeAccount, err := ks.ImportECDSA(ecdsaPrivKey, password)
	if err != nil {
		return "", nodeConfig, logger.CreateDetails(location, err)
	}

	nodeConfig, err = initAccount(ks, &nodeAccount, password)
	if err != nil {
		return "", nodeConfig, logger.CreateDetails(location, err)
	}

	return nodeAccount.Address.String(), nodeConfig, nil
}

//Login checks wallet's address and user's password that was used for crypto wallet creation.
func Login(accountAddress, password string) (*accounts.Account, error) {
	const location = "account.Login->"

	defer debug.FreeOSMemory()

	scryptN, scryptP := encryption.GetScryptParams()

	ks := keystore.NewKeyStore(paths.AccsDirPath, scryptN, scryptP)
	nodeAccounts := ks.Accounts()

	var account *accounts.Account

	for _, a := range nodeAccounts {
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

	secrKey := make([]byte, 32)
	rand.Read(secrKey)

	encryption.SecretKey = secrKey

	secretKeyHash := sha256.Sum256(secrKey)
	encryptedKey, err := encryption.EncryptAES(secretKeyHash[:], key.PrivateKey.D.Bytes())
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	encryption.EncryptedPK = encryptedKey

	return account, nil
}

//ValidateUser asks user for password and checks it.
func ValidateUser() (*accounts.Account, string, error) {
	const location = "account.ValidateUser->"
	var accountAddress, password string
	var nodeAccount *accounts.Account

	accounts := List()

	if len(accounts) > 1 {
		fmt.Println("Please choose an account number")
		for i, a := range accounts {
			fmt.Println(i+1, a)
		}
	}

	for {
		if len(accounts) == 1 {
			accountAddress = accounts[0]
		} else {
			number, err := termEmul.ReadInput()
			if err != nil {
				return nil, "", logger.CreateDetails(location, err)
			}

			accNum, err := strconv.Atoi(number)
			if err != nil {
				fmt.Println("Incorrect value, try again")
				for i, a := range accounts {
					fmt.Println(i+1, a)
				}
				continue
			}

			if accNum < 1 || accNum > len(accounts) {
				fmt.Println("Incorrect value, try again")
				for i, a := range accounts {
					fmt.Println(i+1, a)
				}
				continue
			}

			accountAddress = accounts[accNum-1]
		}

		if !AccExists(accounts, accountAddress) {
			fmt.Println("There is no such account address:")
			for i, a := range accounts {
				fmt.Println(i+1, a)
			}
			continue
		}

		break
	}

	loggedIn := false
	attempts := 3
	for i := 0; i < attempts; i++ {
		fmt.Println("Please enter your password:")
		fmt.Println("Attempts left:", attempts-i)

		bytePassword, err := gopass.GetPasswdMasked()
		if err != nil {
			return nil, "", logger.CreateDetails(location, err)
		}

		originalPassword := string(bytePassword)
		if strings.Trim(originalPassword, " ") == "" {
			fmt.Println("Empty string can't be used as a password. Please enter passwords again")
			continue
		}

		password = hash.Password(originalPassword)
		originalPassword = ""
		bytePassword = nil

		nodeAccount, err = Login(accountAddress, password)
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

	return nodeAccount, password, nil
}

//InitAccount creates directories and files needed for correct work.
func initAccount(ks *keystore.KeyStore, account *accounts.Account, password string) (config.NodeConfig, error) {
	const location = "account.initAccount->"
	var nodeConf config.NodeConfig

	defer debug.FreeOSMemory()

	addressString := account.Address.String()

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

	secrKey := make([]byte, 32)
	rand.Read(secrKey)

	encryption.SecretKey = secrKey

	secretKeyHash := sha256.Sum256(secrKey)
	encryptedKey, err := encryption.EncryptAES(secretKeyHash[:], key.PrivateKey.D.Bytes())
	if err != nil {
		return nodeConf, logger.CreateDetails(location, err)
	}

	encryption.EncryptedPK = encryptedKey

	nodeConf, err = config.Create(addressString)
	if err != nil {
		return nodeConf, logger.CreateDetails(location, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if !tstpkg.TestMode {
		fmt.Println("Registering node...")

		err = blckChain.RegisterNode(ctx, addressString, password, nodeConf.IpAddress, nodeConf.HTTPPort)
		if err != nil {
			return nodeConf, logger.CreateDetails(location, err)
		}

		nodeConf.RegisteredInNetworks[blckChain.CurrentNetwork] = true

		pathToConfigFile := filepath.Join(paths.AccsDirPath, addressString, paths.ConfDirName, paths.ConfFileName)

		confFile, _, err := nodeFile.Read(pathToConfigFile)
		if err != nil {
			return nodeConf, logger.CreateDetails(location, err)
		}
		defer confFile.Close()

		err = config.Save(confFile, nodeConf)
		if err != nil {
			return nodeConf, logger.CreateDetails(location, err)
		}
	}

	err = os.MkdirAll(filepath.Join(paths.StoragePaths[0]), 0700)
	if err != nil {
		return nodeConf, logger.CreateDetails(location, err)
	}

	return nodeConf, nil
}

// ====================================================================================

func AccExists(accounts []string, address string) bool {
	for _, a := range accounts {
		if a == address {
			return true
		}
	}
	return false
}
