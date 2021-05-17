package account

import (
	"crypto/ecdsa"
	"dfile-secondary-node/common"
	"errors"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	commonEtherium "github.com/ethereum/go-ethereum/common"
)

//GetAllAccounts go to the folder ~/dfile/accounts and return all accounts addresses in string format
func GetAllAccounts() ([]string, error) {
	var blockchainAccounts []string

	accountDir, err := common.GetAccountDirectory()
	if err != nil {
		return nil, err
	}

	ks := keystore.NewKeyStore(accountDir, keystore.StandardScryptN, keystore.StandardScryptP)
	etherAccounts := ks.Accounts()

	blockchainAccounts = make([]string, 0, 1)

	for _, a := range etherAccounts {
		blockchainAccounts = append(blockchainAccounts, a.Address.String())
	}

	return blockchainAccounts, err
}

// CreateAccount creates account and keystore file with encryption with password
func CreateAccount(password string) (string, error) {
	accountDir, err := common.GetAccountDirectory()
	if err != nil {
		return "", err
	}

	ks := keystore.NewKeyStore(accountDir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		return "", err
	}

	// err = os.MkdirAll(filepath.Join(accountDir, account.Address.String(), "storage"), 0700)
	// if err != nil {
	// 	return "", err
	// }

	return account.Address.String(), nil
}

//DFileAccount is simple structure with main fields for working with smart contracts and blockchain
type DFileAccount struct {
	Address    commonEtherium.Address
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

//LoadAccount load in memory keystore file and decrypt it for further use
func (dfileAccount *DFileAccount) LoadAccount(blockchainAccountString, password string) error {
	accountDir, err := common.GetAccountDirectory()
	if err != nil {
		return err
	}

	ks := keystore.NewKeyStore(accountDir, keystore.StandardScryptN, keystore.StandardScryptP)
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

	dfileAccount.PrivateKey = key.PrivateKey
	dfileAccount.Address = (*etherAccount).Address
	publicKey := dfileAccount.PrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		errAccountPublicKey := errors.New("Account Public Key error: unable to cast from crypto to ecdsa.")
		return errAccountPublicKey
	}
	dfileAccount.PublicKey = publicKeyECDSA

	return nil
}

func CheckPassword(password string, address string) error {

	accountDir, err := common.GetAccountDirectory()
	if err != nil {
		return err
	}

	ks := keystore.NewKeyStore(accountDir, keystore.StandardScryptN, keystore.StandardScryptP)
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
