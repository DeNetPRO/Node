package blockchainprovider

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"

	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	abiPOS "git.denetwork.xyz/dfile/dfile-secondary-node/POS_abi"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	nodeAbi "git.denetwork.xyz/dfile/dfile-secondary-node/node_abi"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const eightKB = 8192
const NFT = "0xBfAfdaE6B77a02A4684D39D1528c873961528342"

//const ethClientAddr = "https://kovan.infura.io/v3/a4a45777ca65485d983c278291e322f2"

const ethClientAddr = "https://kovan.infura.io/v3/6433ee0efa38494a85541b00cd377c5f"

func RegisterNode(ctx context.Context, address, password string, ip []string, port string) error {
	const logLoc = "blockchainprovider.RegisterNode->"
	ipAddr := [4]uint8{}

	for i, v := range ip {
		intIPPart, err := strconv.Atoi(v)
		if err != nil {
			return logger.CreateDetails(logLoc, err)
		}

		ipAddr[i] = uint8(intIPPart)
	}

	intPort, err := strconv.Atoi(port)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	client, err := ethclient.Dial(ethClientAddr)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	defer client.Close()

	blockNum, err := client.BlockNumber(ctx)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	balance, err := client.BalanceAt(ctx, common.HexToAddress(address), big.NewInt(int64(blockNum-1)))
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	balanceIsInsufficient := balance.Cmp(big.NewInt(200000000000000)) == -1

	if balanceIsInsufficient {
		fmt.Println("Your account has insufficient funds for registering in net. Balance:", balance, "wei")
		fmt.Println("Please top up your balance")
		os.Exit(0)
	}

	node, err := nodeAbi.NewNodeNft(common.HexToAddress(NFT), client)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	opts, err := initTrxOpts(ctx, client, shared.NodeAddr, password, blockNum)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	_, err = node.CreateNode(opts, ipAddr, uint16(intPort))
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	return nil
}

// ====================================================================================

func GetNodeInfoByID() (nodeAbi.SimpleMetaDataDeNetNode, error) {
	const logLoc = "blockchainprovider.GetNodeInfoByID->"
	var nodeInfo nodeAbi.SimpleMetaDataDeNetNode

	client, err := ethclient.Dial(ethClientAddr)
	if err != nil {
		return nodeInfo, logger.CreateDetails(logLoc, err)
	}

	defer client.Close()

	node, err := nodeAbi.NewNodeNft(common.HexToAddress(NFT), client)
	if err != nil {
		return nodeInfo, logger.CreateDetails(logLoc, err)
	}

	nodeInfo, err = node.GetNodeById(&bind.CallOpts{}, big.NewInt(2))
	if err != nil {
		return nodeInfo, logger.CreateDetails(logLoc, err)
	}

	return nodeInfo, nil
}

// ====================================================================================

func GetNodeNFT() (*nodeAbi.NodeNft, error) {
	const logLoc = "blockchainprovider.getNodeNFT->"

	nftAddr := common.HexToAddress("0xBfAfdaE6B77a02A4684D39D1528c873961528342")

	// https://kovan.infura.io/v3/a4a45777ca65485d983c278291e322f2

	client, err := ethclient.Dial("https://kovan.infura.io/v3/6433ee0efa38494a85541b00cd377c5f")
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	defer client.Close()

	node, err := nodeAbi.NewNodeNft(nftAddr, client)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	return node, err
}

// ====================================================================================

func UpdateNodeInfo(ctx context.Context, nodeAddr common.Address, password, newPort string, newIP []string) error {
	const logLoc = "blockchainprovider.UpdateNodeInfo->"
	ipInfo := [4]uint8{}

	for i, v := range newIP {
		intPart, err := strconv.Atoi(v)
		if err != nil {
			return logger.CreateDetails(logLoc, err)
		}

		ipInfo[i] = uint8(intPart)
	}

	intPort, err := strconv.Atoi(newPort)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	client, err := ethclient.Dial(ethClientAddr)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	defer client.Close()

	node, err := nodeAbi.NewNodeNft(common.HexToAddress(NFT), client)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	blockNum, err := client.BlockNumber(ctx)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	opts, err := initTrxOpts(ctx, client, nodeAddr, password, blockNum)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	_, err = node.UpdateNode(opts, big.NewInt(2), ipInfo, uint16(intPort))
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	return nil
}

// ====================================================================================

func StartMining(password string) {
	const logLoc = "blockchainprovider.StartMining->"

	pathToAccStorage := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName)

	regAddr := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	regFileName := regexp.MustCompile("[0-9A-Za-z_]")

	client, err := ethclient.Dial(ethClientAddr)
	if err != nil {
		logger.Log(logger.CreateDetails(logLoc, err))
	}
	defer client.Close()

	tokenAddress := common.HexToAddress("0x2E8630780A231E8bCf12Ba1172bEB9055deEBF8B")
	instance, err := abiPOS.NewStore(tokenAddress, client)
	if err != nil {
		logger.Log(logger.CreateDetails(logLoc, err))
	}

	for {
		fmt.Println("Sleeping...")
		time.Sleep(time.Minute * 10)
		storageProviderAddresses := []string{}
		err = filepath.WalkDir(pathToAccStorage,
			func(path string, info fs.DirEntry, err error) error {
				if err != nil {
					logger.Log(logger.CreateDetails(logLoc, err))
				}

				if regAddr.MatchString(info.Name()) {
					storageProviderAddresses = append(storageProviderAddresses, info.Name())
				}

				return nil
			})

		if err != nil {
			logger.Log(logger.CreateDetails(logLoc, err))
			continue
		}

		if len(storageProviderAddresses) == 0 {
			continue
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Minute*1)

		blockNum, err := client.BlockNumber(ctx)
		if err != nil {
			logger.Log(logger.CreateDetails(logLoc, err))
			continue
		}

		nodeBalance, err := client.BalanceAt(ctx, shared.NodeAddr, big.NewInt(int64(blockNum-1)))
		if err != nil {
			logger.Log(logger.CreateDetails(logLoc, err))
			continue
		}

		nodeBalanceIsLow := nodeBalance.Cmp(big.NewInt(1500000000000000)) == -1

		if nodeBalanceIsLow {
			fmt.Println("Your account has insufficient funds for paying transaction fee. Balance:", nodeBalance, "wei")
			fmt.Println("Please top up your balance")
			continue
		}

		for _, spAddress := range storageProviderAddresses {

			time.Sleep(time.Second * 5)

			storageProviderAddr := common.HexToAddress(spAddress)
			_, reward, userDifficulty, err := instance.GetUserRewardInfo(&bind.CallOpts{}, storageProviderAddr) // first value is paymentToken
			if err != nil {
				logger.Log(logger.CreateDetails(logLoc, err))
			}

			fileNames := []string{}

			pathToStorProviderFiles := filepath.Join(pathToAccStorage, storageProviderAddr.String())

			err = filepath.WalkDir(pathToStorProviderFiles,
				func(path string, info fs.DirEntry, err error) error {
					if err != nil {
						logger.Log(logger.CreateDetails(logLoc, err))
					}

					if regFileName.MatchString(info.Name()) && len(info.Name()) == 64 {
						fileNames = append(fileNames, info.Name())
					}

					return nil
				})
			if err != nil {
				logger.Log(logger.CreateDetails(logLoc, err))
				continue
			}

			if len(fileNames) == 0 {
				err = os.RemoveAll(pathToStorProviderFiles)
				if err != nil {
					logger.Log(logger.CreateDetails(logLoc, err))
				}
				continue
			}

			fmt.Println("reward for", spAddress, "files is", reward) //TODO remove
			fmt.Println("Min reward value:", 350000000000000000)

			rewardisEnough := reward.Cmp(big.NewInt(350000000000000000)) == 1

			if !rewardisEnough {
				continue
			}

			rand.Seed(time.Now().UnixNano())
			randomFilePos := rand.Intn(len(fileNames))

			fileName := fileNames[randomFilePos]

			shared.MU.Lock()

			storedFile, storedFileBytes, err := shared.ReadFile(filepath.Join(pathToStorProviderFiles, fileName))
			if err != nil {
				logger.Log(logger.CreateDetails(logLoc, err))
				continue
			}

			storedFile.Close()
			shared.MU.Unlock()

			ctx, _ := context.WithTimeout(context.Background(), time.Minute*1)

			blockNum, err := client.BlockNumber(ctx)
			if err != nil {
				logger.Log(logger.CreateDetails(logLoc, err))
				continue
			}

			blockHash, err := instance.GetBlockHash(&bind.CallOpts{}, uint32(blockNum-1))
			if err != nil {
				logger.Log(logger.CreateDetails(logLoc, err))
				continue
			}

			fileBytesAddrBlockHash := append(storedFileBytes, shared.NodeAddr.Bytes()...)
			fileBytesAddrBlockHash = append(fileBytesAddrBlockHash, blockHash[:]...)

			hashedFileAddrBlock := sha256.Sum256(fileBytesAddrBlockHash)

			stringFileAddrBlock := hex.EncodeToString(hashedFileAddrBlock[:])

			stringFileAddrBlock = strings.TrimLeft(stringFileAddrBlock, "0")

			decodedBigInt, err := hexutil.DecodeBig("0x" + stringFileAddrBlock)
			if err != nil {
				logger.Log(logger.CreateDetails(logLoc, err))
				continue
			}

			baseDfficulty, err := instance.BaseDifficulty(&bind.CallOpts{})
			if err != nil {
				logger.Log(logger.CreateDetails(logLoc, err))
				continue
			}

			// fmt.Println("decodedBigInt", decodedBigInt)
			// fmt.Println("baseDfficulty", baseDfficulty)
			// fmt.Println("userDifficulty", userDifficulty)

			// diffIsMuch, err := instance.IsMatchDifficulty(&bind.CallOpts{}, decodedBigInt, userDifficulty)
			// if err != nil {
			// 	logger.Log(logger.CreateDetails(logLoc, err))
			// 	continue
			// }

			// fmt.Println("diffIsMuch", diffIsMuch)

			remainder := decodedBigInt.Rem(decodedBigInt, baseDfficulty)

			remainderIsLessUserDifficulty := remainder.CmpAbs(userDifficulty) == -1

			if remainderIsLessUserDifficulty {
				fmt.Println("checking file:", fileName)
				fmt.Println("Trying proof", fileName, "for reward:", reward)

				err := sendProof(ctx, client, password, storedFileBytes, shared.NodeAddr, spAddress, blockNum, instance)
				if err != nil {
					logger.Log(logger.CreateDetails(logLoc, err))
					continue
				}
			}

		}
	}
}

// ====================================================================================

func sendProof(ctx context.Context, client *ethclient.Client, password string, fileBytes []byte,
	nodeAddr common.Address, spAddress string, blockNum uint64, instance *abiPOS.Store) error {
	const logLoc = "blockchainprovider.sendProof->"
	pathToFsTree := filepath.Join(paths.AccsDirPath, nodeAddr.String(), paths.StorageDirName, spAddress, "tree.json")

	shared.MU.Lock()

	spFsFile, spFsBytes, err := shared.ReadFile(pathToFsTree)
	if err != nil {
		shared.MU.Unlock()
		return logger.CreateDetails(logLoc, err)
	}

	spFsFile.Close()
	shared.MU.Unlock()

	var spFs shared.StorageProviderFs

	err = json.Unmarshal(spFsBytes, &spFs)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	eightKBHashes := []string{}

	bytesToProve := fileBytes[:eightKB]

	for i := 0; i < len(fileBytes); i += eightKB {
		hSum := sha256.Sum256(fileBytes[i : i+eightKB])
		eightKBHashes = append(eightKBHashes, hex.EncodeToString(hSum[:]))
	}

	_, fileTree, err := shared.CalcRootHash(eightKBHashes)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	eightKBHashes = nil

	hashFileRoot := fileTree[len(fileTree)-1][0]

	treeToFsRoot := [][][]byte{}

	for _, baseHash := range spFs.Tree[0] {
		diff := bytes.Compare(hashFileRoot, baseHash)
		if diff == 0 {
			treeToFsRoot = append(treeToFsRoot, fileTree[:len(fileTree)-1]...)
			treeToFsRoot = append(treeToFsRoot, spFs.Tree...)
		}
	}

	proof := makeProof(fileTree[0][0], treeToFsRoot)

	fsRootHashBytes := proof[len(proof)-1]

	treeToFsRoot = nil

	opts, err := initTrxOpts(ctx, client, nodeAddr, password, blockNum)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	nonceInt, err := strconv.Atoi(spFs.Nonce)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	nonceHex := strconv.FormatInt(int64(nonceInt), 16)

	nonceBytes, err := hex.DecodeString(nonceHex)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	nonce32 := make([]byte, 32-len(nonceBytes))
	nonce32 = append(nonce32, nonceBytes...)

	fsRootNonceBytes := append(fsRootHashBytes[:], nonce32...)

	signedFSRootHash, err := hex.DecodeString(spFs.SignedFsRoot)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	if signedFSRootHash[len(signedFSRootHash)-1] == 1 { //ecdsa version fix
		signedFSRootHash[len(signedFSRootHash)-1] = 28
	} else {
		signedFSRootHash = signedFSRootHash[:64]
	}

	signatureIsValid, err := instance.IsValidSign(&bind.CallOpts{}, common.HexToAddress(spAddress), fsRootNonceBytes, signedFSRootHash)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	if !signatureIsValid {
		return logger.CreateDetails(logLoc, errors.New(spAddress+" signature is not valid"))
	}

	_, err = instance.SendProof(opts, common.HexToAddress(spAddress), uint32(blockNum), fsRootHashBytes, uint64(nonceInt), signedFSRootHash, bytesToProve, proof)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	proof = nil

	return nil
}

// ====================================================================================

func initTrxOpts(ctx context.Context, client *ethclient.Client, nodeAddr common.Address, password string, blockNum uint64) (*bind.TransactOpts, error) {
	const logLoc = "blockchainprovider.initTrxOpts->"

	transactNonce, err := client.NonceAt(ctx, nodeAddr, big.NewInt(int64(blockNum)))
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	chnID, err := client.ChainID(ctx)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	opts := &bind.TransactOpts{
		From:  nodeAddr,
		Nonce: big.NewInt(int64(transactNonce)),
		Signer: func(a common.Address, t *types.Transaction) (*types.Transaction, error) {
			ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
			acs := ks.Accounts()
			for _, ac := range acs {
				if ac.Address == a {
					ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
					err := ks.TimedUnlock(ac, password, 1)
					if err != nil {
						return t, err
					}
					return ks.SignTx(ac, t, chnID)
				}
			}
			return t, nil
		},
		Value:    big.NewInt(0),
		GasPrice: big.NewInt(1000000000),
		GasLimit: 1000000,
		Context:  ctx,
		NoSend:   false,
	}

	return opts, nil
}

// ====================================================================================

func getPos(hash []byte, list [][]byte) int {
	for i, v := range list {
		diff := bytes.Compare(v, hash)
		if diff == 0 {
			return i
		}
	}

	return -1
}

// ====================================================================================

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
