package account

import (
	"crypto/ecdsa"
	"dfile-secondary-node/server"
	"dfile-secondary-node/shared"
	"errors"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	commonEtherium "github.com/ethereum/go-ethereum/common"
)

//GetAllAccounts go to the folder ~/dfile/accounts and return all accounts addresses in string format
func GetAllAccounts() []string {
	var blockchainAccounts []string

	ks := keystore.NewKeyStore(shared.AccDir, keystore.StandardScryptN, keystore.StandardScryptP)
	etherAccounts := ks.Accounts()

	blockchainAccounts = make([]string, 0, 1)

	for _, a := range etherAccounts {
		blockchainAccounts = append(blockchainAccounts, a.Address.String())
	}

	return blockchainAccounts
}

// CreateAccount creates account and keystore file with encryption with password
func AccountCreate(password string) (string, error) {

	ks := keystore.NewKeyStore(shared.AccDir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Join(shared.AccDir, account.Address.String(), "storage"), 0700)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Join(shared.AccDir, account.Address.String(), "config"), 0700)
	if err != nil {
		return "", err
	}

	return account.Address.String(), nil
}

//DFileAccount is simple structure with main fields for working with smart contracts and blockchain
type DFileAccount struct {
	Address    commonEtherium.Address
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

var DfileAcc DFileAccount

//LoadAccount load in memory keystore file and decrypt it for further use
func AccountLogin(blockchainAccountString, password string) error {

	DfileAcc = DFileAccount{}

	ks := keystore.NewKeyStore(shared.AccDir, keystore.StandardScryptN, keystore.StandardScryptP)
	etherAccounts := ks.Accounts()

	var etherAccount *accounts.Account

	etherAccount = nil

	for _, a := range etherAccounts {
		if blockchainAccountString == a.Address.String() {
			etherAccount = &a
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

	DfileAcc.PrivateKey = key.PrivateKey
	DfileAcc.Address = (*etherAccount).Address
	publicKey := DfileAcc.PrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		errAccountPublicKey := errors.New("account Public Key error: unable to cast from crypto to ecdsa")
		return errAccountPublicKey
	}
	server.AccountAddress = etherAccount.Address.String()
	DfileAcc.PublicKey = publicKeyECDSA

	return nil
}

func CheckPassword(password string, address string) error {

	ks := keystore.NewKeyStore(shared.AccDir, keystore.StandardScryptN, keystore.StandardScryptP)
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
