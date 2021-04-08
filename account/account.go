package account

import (
	"crypto/ecdsa"
	"dfile-secondary-node/common"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	commonEtherium "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

//TODO: write comment
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

// TODO: write comment
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

	return account.Address.String(), nil
}

type DFileAccount struct {
	Address    commonEtherium.Address
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

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

	fmt.Println(hexutil.Encode(crypto.FromECDSA(dfileAccount.PrivateKey)))
	fmt.Println(hexutil.Encode(crypto.FromECDSAPub(dfileAccount.PublicKey)))
	fmt.Println(crypto.PubkeyToAddress(*dfileAccount.PublicKey).String())

	return nil
}
