package blockchainprovider

import (
	"bytes"
	"context"
	"crypto/sha256"
	abiPOS "dfile-secondary-node/POS_abi"
	nodeNFT "dfile-secondary-node/node_nft_abi"
	"dfile-secondary-node/shared"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type StorageInfo struct {
	Nonce        string     `json:"nonce"`
	SignedFsRoot string     `json:"signedFsRoot"`
	Tree         [][][]byte `json:"tree"`
}

const eightKB = 8192

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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

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

func StartMining() {

	nodeAddr, err := shared.DecryptNodeAddr()
	if err != nil {
		log.Fatal(err)
	}

	pathToAccStorage := filepath.Join(shared.AccsDirPath, nodeAddr.String(), shared.StorageDirName)

	regAddr := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	regFileName := regexp.MustCompile("[0-9A-Za-z_]")

	client, err := ethclient.Dial("https://kovan.infura.io/v3/a4a45777ca65485d983c278291e322f2")
	if err != nil {
		log.Fatal(err)
	}

	tokenAddress := common.HexToAddress("0x2E8630780A231E8bCf12Ba1172bEB9055deEBF8B")
	instance, err := abiPOS.NewStore(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	baseDfficulty, err := instance.BaseDifficulty(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Minute)

	for {

		blockNum, err := client.BlockNumber(ctx)
		if err != nil {
			log.Fatal(err)
		}

		nodeBalance, err := client.BalanceAt(ctx, nodeAddr, big.NewInt(int64(blockNum)))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Balance", nodeBalance)

		storageProviderAddresses := []string{}

		err = filepath.WalkDir(pathToAccStorage,
			func(path string, info fs.DirEntry, err error) error {
				if err != nil {
					log.Fatal("Fatal error")
				}

				if regAddr.MatchString(info.Name()) {
					storageProviderAddresses = append(storageProviderAddresses, info.Name())
				}

				return nil
			})
		if err != nil {
			log.Fatal("Fatal error")
		}

		if len(storageProviderAddresses) == 0 {
			continue
		}

		for _, address := range storageProviderAddresses {
			storageProviderAddr := common.HexToAddress(address)
			paymentToken, rew, userDifficulty, err := instance.GetUserRewardInfo(&bind.CallOpts{}, storageProviderAddr)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(paymentToken, rew, userDifficulty)

			fileNames := []string{}

			pathToStorProviderFiles := filepath.Join(pathToAccStorage, storageProviderAddr.String())

			err = filepath.WalkDir(pathToStorProviderFiles,
				func(path string, info fs.DirEntry, err error) error {
					if err != nil {
						log.Fatal(err)
					}

					if regFileName.MatchString(info.Name()) && len(info.Name()) == 64 {
						fileNames = append(fileNames, info.Name())

					}

					return nil
				})
			if err != nil {
				log.Fatal(err)
			}

			for _, fileName := range fileNames {
				storedFile, err := os.Open(filepath.Join(pathToStorProviderFiles, fileName))
				if err != nil {
					log.Fatal(err)
				}

				storedFileBytes, err := io.ReadAll(storedFile)
				if err != nil {
					log.Fatal(err)
				}

				storedFile.Close()

				blockNum, err := client.BlockNumber(ctx)
				if err != nil {
					log.Fatal(err)
				}

				blockHash, err := instance.GetBlockHash(&bind.CallOpts{}, uint32(blockNum))
				if err != nil {
					log.Fatal(err)
				}

				fileBytesAddrBlockHash := append(storedFileBytes, nodeAddr.Bytes()...)
				fileBytesAddrBlockHash = append(fileBytesAddrBlockHash, blockHash[:]...)

				hashedFileAddrBlock := sha256.Sum256(fileBytesAddrBlockHash)

				hexFileAddrBlock := "0x" + hex.EncodeToString(hashedFileAddrBlock[:])

				decodedBigInt, err := hexutil.DecodeBig(hexFileAddrBlock)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(decodedBigInt)
				fmt.Println(baseDfficulty)

				remainder := decodedBigInt.Rem(decodedBigInt, baseDfficulty)

				lessThanUserDifficulty := -1

				compareResult := remainder.CmpAbs(userDifficulty)

				if compareResult == lessThanUserDifficulty {

				}

			}

		}

		time.Sleep(time.Second * 30)

	}

}

func SendProof(password string) {

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Minute)

	nodeAddr, err := shared.DecryptNodeAddr()
	if err != nil {
		log.Fatal(err)
	}

	pathToAcc := filepath.Join(shared.AccsDirPath, nodeAddr.String())

	pathToFile := filepath.Join(pathToAcc, shared.StorageDirName, "0x537F6af3A07e58986Bb5041c304e9Eb2283396CD", "1ab3a0828b5e6b50ac9a3e76b1b33f49587ecf8ea5a58e2fde0429cc11f02342")

	file, err := os.Open(pathToFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	pathToFsTree := filepath.Join(pathToAcc, shared.StorageDirName, "0x537F6af3A07e58986Bb5041c304e9Eb2283396CD", "tree.json")

	fileFsTree, err := os.Open(pathToFsTree)
	if err != nil {
		log.Fatal("Fatal error")
	}
	defer fileFsTree.Close()

	treeBytes, err := io.ReadAll(fileFsTree)
	if err != nil {
		log.Fatal("Fatal error")
	}

	var storageFsStruct StorageInfo

	err = json.Unmarshal(treeBytes, &storageFsStruct)
	if err != nil {
		log.Fatal("Fatal error")
	}

	eightKBHashes := []string{}

	bytesToProve := fileBytes[:eightKB]

	for i := 0; i < len(fileBytes); i += eightKB {
		hSum := sha256.Sum256(fileBytes[i : i+eightKB])
		eightKBHashes = append(eightKBHashes, hex.EncodeToString(hSum[:]))
	}

	_, fileTree, err := shared.CalcRootHash(eightKBHashes)
	if err != nil {
		log.Fatal("Fatal error")
	}

	hashFileRoot := fileTree[len(fileTree)-1][0]

	treeToFsRoot := [][][]byte{}

	for _, baseHash := range storageFsStruct.Tree[0] {
		diff := bytes.Compare(hashFileRoot, baseHash)
		if diff == 0 {
			treeToFsRoot = append(treeToFsRoot, fileTree[:len(fileTree)-1]...)
			treeToFsRoot = append(treeToFsRoot, storageFsStruct.Tree...)
		}
	}

	proof := makeProof(fileTree[0][0], treeToFsRoot)

	client, err := ethclient.Dial("https://kovan.infura.io/v3/a4a45777ca65485d983c278291e322f2")
	if err != nil {
		log.Fatal(err)
	}

	tokenAddress := common.HexToAddress("0x2E8630780A231E8bCf12Ba1172bEB9055deEBF8B")
	instance, err := abiPOS.NewStore(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	blockNum, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatal(err)
	}

	signedFSRootHash, err := hex.DecodeString(storageFsStruct.SignedFsRoot)
	if err != nil {
		log.Fatal(err)
	}

	chnID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
	}

	transactNonce, err := client.NonceAt(ctx, nodeAddr, big.NewInt(int64(blockNum)))
	if err != nil {
		log.Fatal(err)
	}

	opts := &bind.TransactOpts{
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

	intNonce, err := strconv.Atoi(storageFsStruct.Nonce)
	if err != nil {
		log.Fatal(err)
	}

	dif, err := instance.SendProof(opts, common.HexToAddress("0x537F6af3A07e58986Bb5041c304e9Eb2283396CD"), uint32(blockNum), proof[len(proof)-1], uint64(intNonce), signedFSRootHash[:64], bytesToProve, proof)
	if err != nil {
		log.Fatal(err)
	}

	password = ""

	fmt.Println(dif)

}

func getPos(hash []byte, list [][]byte) int {
	for i, v := range list {
		diff := bytes.Compare(v, hash)
		if diff == 0 {
			return i
		}
	}

	return -1

}

// Builds merkle tree proof
func makeProof(start []byte, tree [][][]byte) [][32]byte { // returns slice of 32 bytes array because smart contract awaits this type
	stage := 0
	proof := [][32]byte{}

	var firstNodePosition int
	var secondNodePosition int

	for stage < len(tree) {
		pos := getPos(start, tree[stage])
		if pos == -1 {
			break
		}

		if pos%2 != 0 {
			firstNodePosition = pos - 1
			secondNodePosition = pos
		} else {
			firstNodePosition = pos
			secondNodePosition = pos + 1
		}

		if len(tree[stage]) == 1 {
			root := [32]byte{}
			for i, v := range tree[stage][0] {
				root[i] = v
			}

			proof = append(proof, root)

			return proof
		}

		firstNode := [32]byte{}
		for i, v := range tree[stage][firstNodePosition] {
			firstNode[i] = v
		}

		proof = append(proof, firstNode)

		secondNode := [32]byte{}
		for i, v := range tree[stage][secondNodePosition] {
			secondNode[i] = v
		}

		proof = append(proof, secondNode)

		concatBytes := append(tree[stage][firstNodePosition], tree[stage][secondNodePosition]...)
		hSum := sha256.Sum256(concatBytes)

		start = hSum[:]
		stage++

	}

	return proof
}
