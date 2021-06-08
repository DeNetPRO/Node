package blockchainprovider

import (
	"context"
	nodeNFT "dfile-secondary-node/node_nft_abi"
	"dfile-secondary-node/shared"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func RegisterNode(password string, ip []string, port int) error {

	nodeAddr, err := shared.DecryptNodeAddr()
	if err != nil {
		return err
	}

	nftAddr := common.HexToAddress("0xBfAfdaE6B77a02A4684D39D1528c873961528342")

	client, err := ethclient.Dial("https://kovan.infura.io/v3/a4a45777ca65485d983c278291e322f2")
	if err != nil {
		return err
	}

	node, err := nodeNFT.NewNodeNft(nftAddr, client)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Minute)

	blockNum, err := client.BlockNumber(ctx)
	if err != nil {
		return err
	}

	transactNonce, err := client.NonceAt(ctx, nodeAddr, big.NewInt(int64(blockNum)))
	if err != nil {
		return err
	}

	chnID, err := client.ChainID(ctx)
	if err != nil {
		return err
	}

	var opts = &bind.TransactOpts{
		From:  nodeAddr,
		Nonce: big.NewInt(int64(transactNonce)),
		Signer: func(a common.Address, t *types.Transaction) (*types.Transaction, error) {
			ks := keystore.NewKeyStore(shared.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
			acs := ks.Accounts()
			for _, ac := range acs {
				if ac.Address == a {
					ks := keystore.NewKeyStore(shared.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
					ks.TimedUnlock(ac, password, 1)
					return ks.SignTx(ac, t, chnID)
				}
			}
			return t, nil
		},
		Value:    big.NewInt(0),
		GasPrice: big.NewInt(5000000000),
		GasLimit: 1000000,
		Context:  ctx,
		NoSend:   false,
	}

	ipAddr := [4]uint8{}

	for i, v := range ip {
		intIPPart, err := strconv.Atoi(v)
		if err != nil {
			return err
		}

		ipAddr[i] = uint8(intIPPart)
	}

	_, err = node.CreateNode(opts, ipAddr, uint16(port))
	if err != nil {
		return err
	}

	return nil
}
