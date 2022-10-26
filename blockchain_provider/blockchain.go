package blckChain

import (
	"bytes"
	"context"
	"errors"
	"log"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/minio/sha256-simd"

	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/config"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/encryption"
	erc20 "git.denetwork.xyz/DeNet/dfile-secondary-node/erc20"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/hash"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/networks"

	nodeFile "git.denetwork.xyz/DeNet/dfile-secondary-node/node_file"
	nodeNftAbi "git.denetwork.xyz/DeNet/dfile-secondary-node/node_nft_abi"
	nodeTypes "git.denetwork.xyz/DeNet/dfile-secondary-node/node_types"
	PoS "git.denetwork.xyz/DeNet/dfile-secondary-node/pos"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/sign"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const eightKB = 8192

var (
	mutex     sync.Mutex
	proofOpts *bind.TransactOpts
)

//RegisterNode registers a node in the ethereum network.
//Node's balance should have more than 200000000000000 wei to pay transaction comission.
func RegisterNode(ctx context.Context, nodeAddr common.Address, password string, nodeConfig nodeTypes.Config) error {
	const location = "blckChain.RegisterNode->"
	ipAddr := [4]uint8{}

	splitIPAddr := strings.Split(nodeConfig.IpAddress, ".")

	for i, v := range splitIPAddr {

		intIPPart, err := strconv.Atoi(v)
		if err != nil {
			return logger.MarkLocation(location, err)
		}

		ipAddr[i] = uint8(intIPPart)
	}

	port := strings.TrimPrefix(nodeConfig.HTTPPort, ":")

	intPort, err := strconv.Atoi(port)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	client, err := ethclient.Dial(nodeConfig.RPC[nodeConfig.Network])
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	defer client.Close()

	blockNum, err := client.BlockNumber(ctx)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	_, err = checkBalance(client, nodeAddr, blockNum, true)
	if err != nil {
		logger.Log(logger.MarkLocation(location, err))
		log.Fatal("couldn't check balance")
	}

	nodeNft, err := nodeNftAbi.NewNodeNft(common.HexToAddress(networks.Fields().NODE), client)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	opts, err := initTrxOpts(ctx, client, nodeAddr, password, blockNum)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	_, err = nodeNft.CreateNode(opts, ipAddr, uint16(intPort))
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

//UpdateNodeInfo updates node's ip address or port info.
func UpdateNodeInfo(ctx context.Context, nodeAddr common.Address, password, newIP, newPort string) error {
	const location = "blckChain.UpdateNodeInfo->"
	ipInfo := [4]uint8{}

	splitIPAddr := strings.Split(newIP, ".")

	for i, v := range splitIPAddr {
		intPart, err := strconv.Atoi(v)
		if err != nil {
			return logger.MarkLocation(location, err)
		}

		ipInfo[i] = uint8(intPart)
	}

	port := strings.TrimLeft(newPort, ":")

	intPort, err := strconv.Atoi(port)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	client, err := ethclient.Dial(config.RPC)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	defer client.Close()

	nodeNft, err := nodeNftAbi.NewNodeNft(common.HexToAddress(networks.Fields().NODE), client)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	blockNum, err := client.BlockNumber(ctx)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	opts, err := initTrxOpts(ctx, client, nodeAddr, password, blockNum)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	nodeId, err := nodeNft.GetNodeIDByAddress(&bind.CallOpts{BlockNumber: big.NewInt(int64(blockNum))}, nodeAddr)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	_, err = nodeNft.UpdateNode(opts, nodeId, ipInfo, uint16(intPort))
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

//StartMakingProofs checks reward value for stored file part and sends proof to smart contract if reward is enough.
func StartMakingProofs(nodeAddr common.Address, password string, nodeConfig nodeTypes.Config) {
	const location = "blckChain.StartMakingProofs->"

	regAddr := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	regFileName := regexp.MustCompile("[0-9A-Za-z_]")

	client, err := ethclient.Dial(nodeConfig.RPC[nodeConfig.Network])
	if err != nil {
		logger.Log(logger.MarkLocation(location, err))
		log.Fatal("couldn't set up a new network client")
	}
	defer client.Close()

	posInstance, err := PoS.NewPos(common.HexToAddress(networks.Fields().PoS), client)
	if err != nil {
		logger.Log(logger.MarkLocation(location, err))
		log.Fatal("couldn't set up new proof of storage instance")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

	blockNum, err := client.BlockNumber(ctx)
	if err != nil {
		cancel()
		logger.Log(logger.MarkLocation(location, err))
	}

	_, err = checkBalance(client, nodeAddr, blockNum, true)
	if err != nil {
		cancel()
		logger.Log(logger.MarkLocation(location, err))
		log.Fatal("couldn't check balance")
	}

	baseDiff, err := posInstance.BaseDifficulty(&bind.CallOpts{BlockNumber: big.NewInt(int64(blockNum))})
	if err != nil {
		cancel()
		logger.Log(logger.MarkLocation(location, err))
		log.Fatal("couldn't get base difficulty")
	}

	proofOpts, err = initTrxOpts(ctx, client, nodeAddr, password, blockNum)
	if err != nil {
		cancel()
		logger.Log(logger.MarkLocation(location, err))
		log.Fatal("couldn't initialize transaction options")
	}

	debug.FreeOSMemory()

	transactNonce, err := client.NonceAt(ctx, nodeAddr, big.NewInt(int64(blockNum)))
	if err != nil {
		cancel()
		logger.Log(logger.MarkLocation(location, err))
		log.Fatal("couldn't get transaction nonce")
	}

	cancel()

	proofOpts.Nonce = big.NewInt(int64(transactNonce))

	fmt.Println(networks.Current(), "network selected")

	pathToAccStorage := filepath.Join(paths.List().Storages[0], networks.Current())

	for {

		time.Sleep(time.Second * 20)

		stat, err := os.Stat(pathToAccStorage)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			logger.Log(logger.MarkLocation(location, err))
			log.Fatal(err)
		}

		if stat == nil {
			fmt.Println("no files from", networks.Current(), "to proof")
			time.Sleep(time.Minute * 1)
			continue
		}

		dirFiles, err := nodeFile.ReadDirFiles(pathToAccStorage)
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
			continue
		}

		storageProviderAddresses := []string{}

		for _, f := range dirFiles {
			if regAddr.MatchString(f.Name()) {
				storageProviderAddresses = append(storageProviderAddresses, f.Name())
			}
		}

		if len(storageProviderAddresses) == 0 {
			err = os.Remove(pathToAccStorage)
			if err != nil {
				logger.Log(logger.MarkLocation(location, err))
			}
			continue
		}

		for _, spAddress := range storageProviderAddresses {

			time.Sleep(time.Second * 10)

			pathToFsTree := filepath.Join(paths.List().Storages[0], networks.Current(), spAddress, paths.List().SpFsFilename)

			mutex.Lock()

			spFsFile, spFsBytes, err := nodeFile.Read(pathToFsTree)
			if err != nil {
				mutex.Unlock()
				logger.Log(logger.MarkLocation(location, err))
				continue
			}

			spFsFile.Close()
			mutex.Unlock()

			var spFs nodeTypes.StorageProviderData

			err = json.Unmarshal(spFsBytes, &spFs)
			if err != nil {
				logger.Log(logger.MarkLocation(location, err))
				continue
			}

			reward, userDifficulty, err := posInstance.GetUserRewardInfo(&bind.CallOpts{}, common.HexToAddress(spAddress), big.NewInt(int64(spFs.Storage))) // first value is paymentToken
			if err != nil {
				logger.Log(logger.MarkLocation(location, err))
				continue
			}

			erc, err := getERCContract(client)
			if err != nil {
				logger.Log(logger.MarkLocation(location, err))
				continue
			}

			weiBalance, err := erc.BalanceOf(&bind.CallOpts{}, common.HexToAddress(spAddress))
			if err != nil {
				logger.Log(logger.MarkLocation(location, err))
				continue
			}

			rewardToGbY := float64(reward.Int64()) / 1000000000000000000 * 1000

			if rewardToGbY < 0.3 {
				fmt.Println("reward for", spAddress, "files is not enough")

				continue
			}

			balanceIsEnough := reward.Cmp(weiBalance) == -1

			if !balanceIsEnough {
				fmt.Println(spAddress, "balance is not enough")
				// if weiBalance.Cmp(big.NewInt(0)) == 0 {
				// cleaner.MarkUnused(spAddress)
				// }
				continue
			}

			pathToStorProviderFiles := filepath.Join(pathToAccStorage, spAddress)

			dirFiles, err := nodeFile.ReadDirFiles(pathToStorProviderFiles)
			if err != nil {
				logger.Log(logger.MarkLocation(location, err))
				continue
			}

			fileNames := []string{}

			for _, f := range dirFiles {
				if len(f.Name()) == 64 && regFileName.MatchString(f.Name()) {
					fileNames = append(fileNames, f.Name())
				}
			}

			if len(fileNames) == 0 {
				err = os.Remove(pathToStorProviderFiles)
				if err != nil {
					logger.Log(logger.MarkLocation(location, err))
				}
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

			blockNum, err = client.BlockNumber(ctx)
			if err != nil {
				cancel()
				logger.Log(logger.MarkLocation(location, err))
				continue
			}

			cancel()

			blockHash, err := posInstance.GetBlockHash(&bind.CallOpts{}, uint32(blockNum-10)) // checking older blocknum to guarantee valid result
			if err != nil {
				logger.Log(logger.MarkLocation(location, err))
				continue
			}

			for _, fileName := range fileNames {
				mutex.Lock()

				storedFile, storedFileBytes, err := nodeFile.Read(filepath.Join(pathToStorProviderFiles, fileName))
				if err != nil {
					mutex.Unlock()
					logger.Log(logger.MarkLocation(location, err))
					continue
				}

				storedFile.Close()
				mutex.Unlock()

				fileEightKB := make([]byte, eightKB)

				copy(fileEightKB, storedFileBytes[:eightKB])

				fileBytesAddrBlockHash := append(fileEightKB, nodeAddr.Bytes()...)
				fileProof := append(fileBytesAddrBlockHash, blockHash[:]...)

				fileProofSha := sha256.Sum256(fileProof)

				stringFileProof := hex.EncodeToString(fileProofSha[:])

				stringFileProof = strings.TrimLeft(stringFileProof, "0") // leading zeroes lead to decoding errors

				bigIntFromProof, err := hexutil.DecodeBig("0x" + stringFileProof)
				if err != nil {
					logger.Log(logger.MarkLocation(location, err))
				}

				remainder := bigIntFromProof.Rem(bigIntFromProof, baseDiff)

				difficultyIsEnough := remainder.CmpAbs(userDifficulty) == -1

				if !difficultyIsEnough {
					// fmt.Println("difficulty is not enough")
					continue
				}

				fmt.Println("Trying proof", fileName, "for reward:", reward)

				err = sendProof(client, spFs, storedFileBytes, nodeAddr, common.HexToAddress(spAddress), blockNum-10, posInstance) // sending blocknum that we used for verifying proof
				if err != nil {
					logger.Log(logger.MarkLocation(location, err))
					continue
				} else {

					fmt.Println("proof is sent")

					break
				}

			}
		}

	}
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

//SendProof checks Storage Providers's file system root hash and nounce info and sends proof to smart contract.
func sendProof(client *ethclient.Client, spFs nodeTypes.StorageProviderData, fileBytes []byte, nodeAddr common.Address, spAddress common.Address,
	blockNum uint64, posInstance *PoS.Pos) error {

	const location = "blckChain.sendProof->"

	balanceIsLow, err := checkBalance(client, nodeAddr, blockNum, false)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	if balanceIsLow {
		return logger.MarkLocation(location, errors.New("not sufficient funds for transactions"))
	}

	eightKBHashes := []string{}

	for i := 0; i < len(fileBytes); i += eightKB {
		hSum := sha256.Sum256(fileBytes[i : i+eightKB])
		eightKBHashes = append(eightKBHashes, hex.EncodeToString(hSum[:]))
	}

	_, fileTree, err := hash.CalcRoot(eightKBHashes)
	if err != nil {
		return logger.MarkLocation(location, err)
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

	path := makePath(fileTree[0][0], treeToFsRoot)

	if len(path) == 0 {
		return logger.MarkLocation(location, errors.New("proof is empty"))
	}

	fsRootHashBytes := path[len(path)-1]

	contractRootHash, contractNonce, err := posInstance.GetUserRootHash(&bind.CallOpts{}, spAddress)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	var zeroHash [32]byte

	contractFsHashIsZero := bytes.Equal(zeroHash[:], contractRootHash[:])

	curentNonceIsZero := contractNonce.Cmp(big.NewInt(int64(0))) == 0

	firstProof := curentNonceIsZero && contractFsHashIsZero

	if !firstProof {
		contractNonceIsBigger := contractNonce.Cmp(big.NewInt(int64(spFs.Nonce))) == 1

		rootHashesEqual := bytes.Equal(fsRootHashBytes[:], contractRootHash[:])

		if contractNonceIsBigger && !rootHashesEqual {
			if contractNonceIsBigger {
				fmt.Println("Contract nonce is bigger, contract nonce:", contractNonce, "provider nonce:", spFs.Nonce)
			}

			if !rootHashesEqual {
				fmt.Println("contract root hash is not equal to provider root hash")
			}

			return logger.MarkLocation(location, errors.New("fs root hash info is not valid"))

		}
	}

	nonceBytes := make([]byte, 4)

	binary.BigEndian.PutUint32(nonceBytes, spFs.Nonce)

	alignBytes := make([]byte, 28) // need to align nonce and storage info to 32 bytes for smart contract

	alignedNonceBytes := append(alignBytes, nonceBytes...)

	storageBytes := make([]byte, 4)

	binary.BigEndian.PutUint32(storageBytes, spFs.Storage)

	alignedStorage := append(alignBytes, storageBytes...)

	fsRootStorageBytes := append(fsRootHashBytes[:], alignedStorage...)

	fsRootStorageNonceBytes := append(fsRootStorageBytes, alignedNonceBytes...)

	err = sign.Check(spAddress.String(), spFs.SignedFsInfo, sha256.Sum256(fsRootStorageNonceBytes))
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	signedFSRootNonceStorage, err := hex.DecodeString(spFs.SignedFsInfo)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	if signedFSRootNonceStorage[len(signedFSRootNonceStorage)-1] == 1 { //ecdsa version fix
		signedFSRootNonceStorage[len(signedFSRootNonceStorage)-1] = 28
	} else {
		signedFSRootNonceStorage = signedFSRootNonceStorage[:64]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	proofOpts.Context = ctx

	trx, err := posInstance.SendProof(proofOpts, common.HexToAddress(spAddress.String()), uint32(blockNum), fsRootHashBytes, uint64(spFs.Storage), uint64(spFs.Nonce), signedFSRootNonceStorage, fileBytes[:eightKB], path)
	if err != nil {

		debug.FreeOSMemory()

		if err.Error() == "Transaction nonce is too low. Try incrementing the nonce." {
			proofOpts.Nonce = proofOpts.Nonce.Add(proofOpts.Nonce, big.NewInt(int64(1)))

			fmt.Println("Trying to prove with incremented nonce")

			trx, err = posInstance.SendProof(proofOpts, common.HexToAddress(spAddress.String()), uint32(blockNum), fsRootHashBytes, uint64(spFs.Storage), uint64(spFs.Nonce), signedFSRootNonceStorage, fileBytes[:eightKB], path)
			if err != nil {
				debug.FreeOSMemory()
				return logger.MarkLocation(location, err)
			}

		} else {
			debug.FreeOSMemory()
			proofOpts.Nonce = proofOpts.Nonce.Add(proofOpts.Nonce, big.NewInt(int64(1)))
			return logger.MarkLocation(location, err)
		}
	}

	fmt.Printf("transaction hash: %v\n", fmt.Sprint(networks.Fields().TRX, trx.Hash()))

	debug.FreeOSMemory()
	proofOpts.Nonce = proofOpts.Nonce.Add(proofOpts.Nonce, big.NewInt(int64(1)))

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func checkBalance(client *ethclient.Client, nodeAddr common.Address, blockNum uint64, exitOnLow bool) (bool, error) {

	const location = "blckChain.checkBalance->"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	defer cancel()

	nodeBalance, err := client.BalanceAt(ctx, nodeAddr, big.NewInt(int64(blockNum)))
	if err != nil {
		return false, logger.MarkLocation(location, err)
	}

	nodeBalanceIsLow := nodeBalance.Cmp(big.NewInt(1500000000000000)) == -1

	if nodeBalanceIsLow {
		fmt.Println("Insufficient funds for paying", networks.Current(), "transaction fees. Balance:", nodeBalance)
		fmt.Println("Please top up your balance")

		if exitOnLow {
			fmt.Println("Exited")
			os.Exit(0)
		} else {
			return true, nil
		}

	}

	return false, nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
// InitTrxOpts makes transaction options that are needed when sending request to smart contract.
func initTrxOpts(ctx context.Context, client *ethclient.Client, nodeAddr common.Address, password string, blockNum uint64) (*bind.TransactOpts, error) {
	const location = "blckChain.initTrxOpts->"

	transactNonce, err := client.NonceAt(ctx, nodeAddr, big.NewInt(int64(blockNum)))
	if err != nil {
		return nil, logger.MarkLocation(location, err)
	}

	chnID, err := client.ChainID(ctx)
	if err != nil {
		return nil, logger.MarkLocation(location, err)
	}

	opts := &bind.TransactOpts{
		From:  nodeAddr,
		Nonce: big.NewInt(int64(transactNonce)),
		Signer: func(a common.Address, t *types.Transaction) (*types.Transaction, error) {
			scryptN, scryptP := encryption.GetScryptParams()

			ks := keystore.NewKeyStore(paths.List().AccsDir, scryptN, scryptP)
			acs := ks.Accounts()
			for _, ac := range acs {
				if ac.Address == a {
					err := ks.TimedUnlock(ac, password, 3)
					if err != nil {
						return t, err
					}
					return ks.SignTx(ac, t, chnID)
				}
			}
			return t, nil
		},
		Value:    big.NewInt(0),
		GasPrice: big.NewInt(35000000000), // 35 Gwei
		GasLimit: 1000000,
		Context:  ctx,
		NoSend:   false,
	}

	return opts, nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

// MakeProof builds merkle tree proof. Passed tree value is an array of file hashes that are located on different levels of merkle tree.
// Returns slice of 32 bytes array for passing it to smart contract.
func makePath(start []byte, tree [][][]byte) [][32]byte {
	stage := 0
	path := [][32]byte{}

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

			path = append(path, root)

			return path
		}

		firstNode := [32]byte{}
		for i, v := range tree[stage][firstNodePosition] {
			firstNode[i] = v
		}

		path = append(path, firstNode)

		secondNode := [32]byte{}
		for i, v := range tree[stage][secondNodePosition] {
			secondNode[i] = v
		}

		path = append(path, secondNode)

		concatBytes := append(tree[stage][firstNodePosition], tree[stage][secondNodePosition]...)
		hSum := sha256.Sum256(concatBytes)

		start = hSum[:]
		stage++
	}

	return path
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

// GetPos returns element's position in the merkle tree's checked level.
func getPos(hash []byte, list [][]byte) int {
	for i, v := range list {
		diff := bytes.Compare(v, hash)
		if diff == 0 {
			return i
		}
	}

	return -1
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func getERCContract(client *ethclient.Client) (*erc20.Erc20, error) {
	const location = "blckChain.GetERCContract->"

	ercAddr := common.HexToAddress(networks.Fields().ERC)

	erc, err := erc20.NewErc20(ercAddr, client)
	if err != nil {
		return nil, logger.MarkLocation(location, err)
	}

	return erc, err
}
