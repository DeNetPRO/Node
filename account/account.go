package account

import (
	"crypto/sha256"
	"dfile-secondary-node/shared"
	"errors"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
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

	keyJson, err := ks.Export(etherAccount, password, password)
	if err != nil {
		return "", err
	}

	key, err := keystore.DecryptKey(keyJson, password)
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

	macAddr, err := shared.GetDeviceMacAddr()
	if err != nil {
		return "", err
	}

	encrForAddr := sha256.Sum256([]byte(macAddr))
	encryptedAddr, err := shared.EncryptAES(encrForAddr[:], key.Address.Bytes())
	if err != nil {
		return "", err
	}

	shared.DfileAcc.Address = encryptedAddr

	return addressString, nil
}

//LoadAccount load in memory keystore file and decrypt it for further use
func Login(blockchainAccountString, password string) error {

	ks := keystore.NewKeyStore(shared.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	etherAccounts := ks.Accounts()

	var etherAccount *accounts.Account

	etherAccount = nil

	for _, a := range etherAccounts {
		if blockchainAccountString == a.Address.String() {
			etherAccount = &a
			break
		}
	}

	if etherAccount == nil {
		errAccountNotFound := errors.New("Account Not Found Error: cannot find account for " + blockchainAccountString)
		return errAccountNotFound
	}

	keyJson, err := ks.Export(*etherAccount, password, password)
	if err != nil {
		return err
	}

	key, err := keystore.DecryptKey(keyJson, password)
	if err != nil {
		return err
	}

	macAddr, err := shared.GetDeviceMacAddr()
	if err != nil {
		return err
	}

	encrForAddr := sha256.Sum256([]byte(macAddr))
	encryptedAddr, err := shared.EncryptAES(encrForAddr[:], key.Address.Bytes())
	if err != nil {
		return err
	}

	shared.DfileAcc.Address = encryptedAddr

	return nil
}

func CheckPassword(password string, address string) error {

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
