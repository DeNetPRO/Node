package account

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	blckChain "git.denetwork.xyz/DeNet/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/cleaner"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/networks"
	nodeFile "git.denetwork.xyz/DeNet/dfile-secondary-node/node_file"
	tstpkg "git.denetwork.xyz/DeNet/dfile-secondary-node/tst_pkg"

	termEmul "git.denetwork.xyz/DeNet/dfile-secondary-node/term_emul"
	"github.com/howeyc/gopass"
	"github.com/minio/sha256-simd"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/config"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/encryption"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/hash"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	nodeTypes "git.denetwork.xyz/DeNet/dfile-secondary-node/node_types"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	encryptedPK []byte
	secretKey   []byte
)

//List returns list of user's created/imported wallet adresses, that are used as user accounts.
func List() []string {
	var blockchainAccounts []string

	scryptN, scryptP := encryption.GetScryptParams()

	ks := keystore.NewKeyStore(paths.List().AccsDir, scryptN, scryptP)
	nodeAccounts := ks.Accounts()

	blockchainAccounts = make([]string, 0)

	for _, a := range nodeAccounts {
		blockchainAccounts = append(blockchainAccounts, a.Address.String())
	}

	return blockchainAccounts
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

// Create is used for creating a new crypto wallet with keystore file.
func Create(password string) (string, nodeTypes.Config, error) {
	const location = "account.Create->"
	var nodeConf nodeTypes.Config

	err := paths.CreateAccDirs()
	if err != nil {
		return "", nodeConf, logger.MarkLocation(location, err)
	}

	scryptN, scryptP := encryption.GetScryptParams()

	ks := keystore.NewKeyStore(paths.List().AccsDir, scryptN, scryptP)

	nodeAccount, err := ks.NewAccount(password)
	if err != nil {
		return "", nodeConf, logger.MarkLocation(location, err)
	}

	nodeConf, err = makeAccount(ks, &nodeAccount, password)
	if err != nil {
		return "", nodeConf, logger.MarkLocation(location, err)
	}

	return nodeAccount.Address.String(), nodeConf, nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

//Import is used for importing crypto wallet. Private key is needed.
func Import() (string, nodeTypes.Config, error) {
	const location = "account.Import->"
	var nodeConfig nodeTypes.Config

	var privKey string
	var originalPassword string
	var err error

	if tstpkg.Data().TestMode {
		privKey = tstpkg.Data().PrivateKey
		originalPassword = tstpkg.Data().Password
	} else {
		fmt.Println("Please enter private key of the account you want to import:")

		bytesPrivKey, err := gopass.GetPasswdMasked()
		if err != nil {
			return "", nodeConfig, logger.MarkLocation(location, err)
		}

		privKey = string(bytesPrivKey)

		fmt.Println("\nPlease enter your password:")

		for {
			bytePassword, err := gopass.GetPasswdMasked()
			if err != nil {
				return "", nodeConfig, logger.MarkLocation(location, err)
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
		return "", nodeConfig, logger.MarkLocation(location, err)
	}

	scryptN, scryptP := encryption.GetScryptParams()

	ks := keystore.NewKeyStore(paths.List().AccsDir, scryptN, scryptP)

	ecdsaPrivKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return "", nodeConfig, logger.MarkLocation(location, err)
	}

	nodeAccount, err := ks.ImportECDSA(ecdsaPrivKey, password)
	if err != nil {
		return "", nodeConfig, logger.MarkLocation(location, err)
	}

	nodeConfig, err = makeAccount(ks, &nodeAccount, password)
	if err != nil {
		return "", nodeConfig, logger.MarkLocation(location, err)
	}

	go blckChain.StartMakingProofs(nodeAccount.Address, password, nodeConfig)
	go cleaner.Start()

	return nodeAccount.Address.String(), nodeConfig, nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

//Login checks wallet's address and user's password that was used for crypto wallet creation.
func Login(nodeAddr, password string) (*accounts.Account, error) {
	const location = "account.Login->"

	defer debug.FreeOSMemory()

	scryptN, scryptP := encryption.GetScryptParams()

	ks := keystore.NewKeyStore(paths.List().AccsDir, scryptN, scryptP)
	nodeAccounts := ks.Accounts()

	var account *accounts.Account

	for _, a := range nodeAccounts {
		if nodeAddr == a.Address.String() {
			account = &a
			break
		}
	}

	if account == nil {
		err := errors.New(nodeAddr + " address is not found")
		return nil, logger.MarkLocation(location, err)
	}

	keyJson, err := ks.Export(*account, password, password)
	if err != nil {
		fmt.Println("Wrong password")
		return nil, logger.MarkLocation(location, err)
	}

	key, err := keystore.DecryptKey(keyJson, password)
	if err != nil {
		return nil, logger.MarkLocation(location, err)
	}

	secrKey := make([]byte, 32)
	rand.Read(secrKey)

	secretKey = secrKey

	secretKeyHash := sha256.Sum256(secrKey)
	encryptedKey, err := encryption.EncryptAES(secretKeyHash[:], key.PrivateKey.D.Bytes())
	if err != nil {
		return nil, logger.MarkLocation(location, err)
	}

	encryptedPK = encryptedKey

	return account, nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

//Unlock asks user for password and checks it.
func Unlock() (*accounts.Account, string, error) {
	const location = "account.Unlock->"
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
				return nil, "", logger.MarkLocation(location, err)
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
		fmt.Println("\nPlease enter your password (attempts left:", fmt.Sprint(attempts-i)+")")

		bytePassword, err := gopass.GetPasswdMasked()
		if err != nil {
			return nil, "", logger.MarkLocation(location, err)
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
			logger.MarkLocation(location, err)
			continue
		}

		loggedIn = true
		break
	}

	if !loggedIn {
		return nil, "", logger.MarkLocation(location, errors.New("couldn't log in in 3 attempts"))
	}

	return nodeAccount, password, nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

//makeAccount creates directories and files needed for correct work.
func makeAccount(ks *keystore.KeyStore, account *accounts.Account, password string) (nodeTypes.Config, error) {
	const location = "account.makeAccount->"
	var nodeConf nodeTypes.Config

	defer debug.FreeOSMemory()

	addressString := account.Address.String()

	keyJson, err := ks.Export(*account, password, password)
	if err != nil {
		fmt.Println("Wrong password")
		return nodeConf, logger.MarkLocation(location, err)
	}

	key, err := keystore.DecryptKey(keyJson, password)
	if err != nil {
		return nodeConf, logger.MarkLocation(location, err)
	}

	secrKey := make([]byte, 32)
	rand.Read(secrKey)

	secretKey = secrKey

	secretKeyHash := sha256.Sum256(secrKey)
	encryptedKey, err := encryption.EncryptAES(secretKeyHash[:], key.PrivateKey.D.Bytes())
	if err != nil {
		return nodeConf, logger.MarkLocation(location, err)
	}

	encryptedPK = encryptedKey

	nodeConf, err = config.Create(addressString)
	if err != nil {
		return nodeConf, logger.MarkLocation(location, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if !tstpkg.Data().TestMode {
		fmt.Println("Registering node...")

		err = blckChain.RegisterNode(ctx, account.Address, password, nodeConf)
		if err != nil {
			return nodeConf, logger.MarkLocation(location, err)
		}

		nodeConf.RegisteredInNetworks[networks.Current()] = true

		confFile, _, err := nodeFile.Read(paths.List().ConfigFile)
		if err != nil {
			return nodeConf, logger.MarkLocation(location, err)
		}
		defer confFile.Close()

		err = config.Save(confFile, nodeConf)
		if err != nil {
			return nodeConf, logger.MarkLocation(location, err)
		}
	}

	return nodeConf, nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func AccExists(accounts []string, address string) bool {
	for _, a := range accounts {
		if a == address {
			return true
		}
	}
	return false
}
