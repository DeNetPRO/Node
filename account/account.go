package account

import (
	"dfile-secondary-node/shared"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/crypto"
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

	keyJson, err := ks.Export(etherAccount, password, password) // TODO remove
	if err != nil {
		return "", err
	}

	key, err := keystore.DecryptKey(keyJson, password) // TODO remove
	if err != nil {
		return "", err
	}

	fmt.Println("Private Key:", hex.EncodeToString(key.PrivateKey.D.Bytes())) // TODO remove

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

func Import(privKey, password string) (string, error) {

	shared.CreateIfNotExistAccDirs()

	ks := keystore.NewKeyStore(shared.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)

	ecdsaPrivKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return "", err
	}

	etherAccount, err := ks.ImportECDSA(ecdsaPrivKey, password)
	if err != nil {
		fmt.Println(err)
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

//LoadAccount load in memory keystore file and decrypt it for further use
func Login(blockchainAccountString, password string) error {

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
		return errors.New("Account Not Found Error: cannot find account for " + blockchainAccountString)
	}

	keyJson, err := ks.Export(*etherAccount, password, password)
	if err != nil {
		return err
	}

	key, err := keystore.DecryptKey(keyJson, password)
	if err != nil {
		return err
	}

	encryptedAddr, err := shared.EncryptNodeAddr(key.Address)
	if err != nil {
		return err
	}

	shared.NodeAddr = encryptedAddr

	return nil
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
