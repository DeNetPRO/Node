package blckChain

import (
	"bytes"
	"context"
	"errors"
	"log"
	"runtime/debug"
	"strings"

	"github.com/minio/sha256-simd"

	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	proofOfStAbi "git.denetwork.xyz/DeNet/dfile-secondary-node/POS_abi"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/encryption"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/hash"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	nodeFile "git.denetwork.xyz/DeNet/dfile-secondary-node/node_file"
	nodeNftAbi "git.denetwork.xyz/DeNet/dfile-secondary-node/node_nft_abi"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/sign"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/shared"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type NtwrkParams struct {
	RPC string
	NFT string
	PoS string
}

var Networks = map[string]NtwrkParams{
	"polygon": {
		RPC: "https://polygon-rpc.com",
		NFT: "0xfe1f5CB22cF4972584c6a0938FEAF90c597b567b",
		PoS: "0x70c478be3d87ab921e0168137f5abe53b5812fc8",
	},
	"kovan": {
		RPC: "https://kovan.infura.io/v3/6433ee0efa38494a85541b00cd377c5f",
		NFT: "0x8De6417e4738a41619d0D13ef0661563f1A918ec",
		PoS: "0x60828cfBBFbcB474c913FaDE151AD4AFa9a07919",
	},
	"mumbai": {
		RPC: "https://rpc-mumbai.maticvigil.com",
		NFT: "0xBb86dcf291419d3F5b4B2211122D0E6fCB693777",
		PoS: "0x389E8fE67c73551043184F740126C91866c0fB78",
	},
}

const eightKB = 8192

var (
	proofOpts      *bind.TransactOpts
	CurrentNetwork string
)

//RegisterNode registers a node in the ethereum network.
//Node's balance should have more than 200000000000000 wei to pay transaction comission.
func RegisterNode(ctx context.Context, address, password, ip, port string) error {
	const location = "blckChain.RegisterNode->"
	ipAddr := [4]uint8{}

	splitIPAddr := strings.Split(ip, ".")

	for i, v := range splitIPAddr {
		intIPPart, err := strconv.Atoi(v)
		if err != nil {
			return logger.CreateDetails(location, err)
		}

		ipAddr[i] = uint8(intIPPart)
	}

	intPort, err := strconv.Atoi(port)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	client, err := ethclient.Dial(Networks[CurrentNetwork].RPC)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	defer client.Close()

	blockNum, err := client.BlockNumber(ctx)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	_, err = checkBalance(client, blockNum, true)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		log.Fatal("couldn't check balance")
	}

	nodeNft, err := nodeNftAbi.NewNodeNft(common.HexToAddress(Networks[CurrentNetwork].NFT), client)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	opts, err := initTrxOpts(ctx, client, shared.NodeAddr, password, blockNum)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	_, err = nodeNft.CreateNode(opts, ipAddr, uint16(intPort))
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	return nil
}

// ====================================================================================

//UpdateNodeInfo updates node's ip address or port info.
func UpdateNodeInfo(ctx context.Context, nodeAddr common.Address, password, newIP, newPort string) error {
	const location = "blckChain.UpdateNodeInfo->"
	ipInfo := [4]uint8{}

	splitIPAddr := strings.Split(newIP, ".")

	for i, v := range splitIPAddr {
		intPart, err := strconv.Atoi(v)
		if err != nil {
			return logger.CreateDetails(location, err)
		}

		ipInfo[i] = uint8(intPart)
	}

	intPort, err := strconv.Atoi(newPort)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	client, err := ethclient.Dial(Networks[CurrentNetwork].RPC)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	defer client.Close()

	nodeNft, err := nodeNftAbi.NewNodeNft(common.HexToAddress(Networks[CurrentNetwork].NFT), client)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	blockNum, err := client.BlockNumber(ctx)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	opts, err := initTrxOpts(ctx, client, nodeAddr, password, blockNum)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	nodeId, err := nodeNft.GetNodeIDByAddress(&bind.CallOpts{BlockNumber: big.NewInt(int64(blockNum))}, shared.NodeAddr)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	_, err = nodeNft.UpdateNode(opts, nodeId, ipInfo, uint16(intPort))
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	return nil
}

// ====================================================================================

//StartMakingProofs checks reward value for stored file part and sends proof to smart contract if reward is enough.
func StartMakingProofs(password string) {
	const location = "blckChain.StartMakingProofs->"

	regAddr := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	regFileName := regexp.MustCompile("[0-9A-Za-z_]")

	client, err := ethclient.Dial(Networks[CurrentNetwork].RPC)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		log.Fatal("couldn't set up a new network client")
	}
	defer client.Close()

	posInstance, err := proofOfStAbi.NewProofOfStorage(common.HexToAddress(Networks[CurrentNetwork].PoS), client)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		log.Fatal("couldn't set up new proof of storage instance")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)

	blockNum, err := client.BlockNumber(ctx)
	if err != nil {
		cancel()
		logger.Log(logger.CreateDetails(location, err))
	}

	_, err = checkBalance(client, blockNum, true)
	if err != nil {
		cancel()
		logger.Log(logger.CreateDetails(location, err))
		log.Fatal("couldn't check balance")
	}

	baseDiff, err := posInstance.BaseDifficulty(&bind.CallOpts{BlockNumber: big.NewInt(int64(blockNum))})
	if err != nil {
		cancel()
		logger.Log(logger.CreateDetails(location, err))
		log.Fatal("couldn't get base difficulty")
	}

	proofOpts, err = initTrxOpts(ctx, client, shared.NodeAddr, password, blockNum)
	if err != nil {
		cancel()
		logger.Log(logger.CreateDetails(location, err))
		log.Fatal("couldn't initialize transaction options")
	}

	debug.FreeOSMemory()

	transactNonce, err := client.NonceAt(ctx, shared.NodeAddr, big.NewInt(int64(blockNum)))
	if err != nil {
		cancel()
		logger.Log(logger.CreateDetails(location, err))
		log.Fatal("couldn't get transaction nonce")
	}

	cancel()

	proofOpts.Nonce = big.NewInt(int64(transactNonce))

	fmt.Println(CurrentNetwork, "network selected")

	pathToAccStorage := filepath.Join(paths.StoragePaths[0], CurrentNetwork)

	for {
		stat, err := os.Stat(pathToAccStorage)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal(err)
		}

		if stat == nil {
			fmt.Println("no files from", CurrentNetwork, "to proof")
			time.Sleep(time.Minute * 1)
			continue
		}

		dirFiles, err := nodeFile.ReadDirFiles(pathToAccStorage)
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
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
				logger.Log(logger.CreateDetails(location, err))
			}
			continue
		}

		for _, spAddress := range storageProviderAddresses {

			time.Sleep(time.Second * 10)

			_, reward, userDifficulty, err := posInstance.GetUserRewardInfo(&bind.CallOpts{}, common.HexToAddress(spAddress)) // first value is paymentToken
			if err != nil {
				logger.Log(logger.CreateDetails(location, err))
				continue
			}

			fmt.Println("reward for", spAddress, "files is", reward) //TODO remove
			fmt.Println("Min reward value:", 200000000000000)

			rewardIsEnough := reward.Cmp(big.NewInt(200000000000000)) == 1

			if !rewardIsEnough {
				continue
			}

			pathToStorProviderFiles := filepath.Join(pathToAccStorage, spAddress)

			dirFiles, err := nodeFile.ReadDirFiles(pathToStorProviderFiles)
			if err != nil {
				logger.Log(logger.CreateDetails(location, err))
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
					logger.Log(logger.CreateDetails(location, err))
				}
				continue
			}

			rand.Seed(time.Now().UnixNano())
			randomFilePos := rand.Intn(len(fileNames))

			quarter := len(fileNames) / 4

			fileNames = fileNames[randomFilePos:]
			if quarter > 0 && randomFilePos+quarter < len(fileNames) {
				fileNames = fileNames[randomFilePos : randomFilePos+quarter]
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

			blockNum, err = client.BlockNumber(ctx)
			if err != nil {
				cancel()
				logger.Log(logger.CreateDetails(location, err))
				continue
			}

			cancel()

			blockHash, err := posInstance.GetBlockHash(&bind.CallOpts{}, uint32(blockNum-10)) // checking older blocknum to guarantee valid result
			if err != nil {
				logger.Log(logger.CreateDetails(location, err))
				continue
			}

			for _, fileName := range fileNames {
				shared.MU.Lock()

				storedFile, storedFileBytes, err := nodeFile.Read(filepath.Join(pathToStorProviderFiles, fileName))
				if err != nil {
					shared.MU.Unlock()
					logger.Log(logger.CreateDetails(location, err))
					continue
				}

				storedFile.Close()
				shared.MU.Unlock()

				fileEightKB := make([]byte, eightKB)

				copy(fileEightKB, storedFileBytes[:eightKB])

				fileBytesAddrBlockHash := append(fileEightKB, shared.NodeAddr.Bytes()...)
				fileProof := append(fileBytesAddrBlockHash, blockHash[:]...)

				fileProofSha := sha256.Sum256(fileProof)

				stringFileProof := hex.EncodeToString(fileProofSha[:])

				stringFileProof = strings.TrimLeft(stringFileProof, "0") // leading zeroes lead to decoding errors

				bigIntFromProof, err := hexutil.DecodeBig("0x" + stringFileProof)
				if err != nil {
					logger.Log(logger.CreateDetails(location, err))
				}

				remainder := bigIntFromProof.Rem(bigIntFromProof, baseDiff)

				difficultyIsEnough := remainder.CmpAbs(userDifficulty) == -1

				if !difficultyIsEnough {
					fmt.Println("difficulty is not enough")
					continue
				}

				fmt.Println("Trying proof", fileName, "for reward:", reward)

				err = sendProof(client, storedFileBytes, shared.NodeAddr, common.HexToAddress(spAddress), blockNum-10, posInstance) // sending blocknum that we used for verifying proof
				if err != nil {
					logger.Log(logger.CreateDetails(location, err))
					continue
				} else {

					fmt.Println("proof is sent")

					break
				}

			}
		}

	}
}

// ====================================================================================

//SendProof checks Storage Providers's file system root hash and nounce info and sends proof to smart contract.
func sendProof(client *ethclient.Client, fileBytes []byte, nodeAddr common.Address, spAddress common.Address,
	blockNum uint64, posInstance *proofOfStAbi.ProofOfStorage) error {

	const location = "blckChain.sendProof->"

	balanceIsLow, err := checkBalance(client, blockNum, false)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	if balanceIsLow {
		return logger.CreateDetails(location, errors.New("not sufficient funds for transactions"))
	}

	pathToFsTree := filepath.Join(paths.StoragePaths[0], CurrentNetwork, spAddress.String(), paths.SpFsFilename)

	shared.MU.Lock()

	spFsFile, spFsBytes, err := nodeFile.Read(pathToFsTree)
	if err != nil {
		shared.MU.Unlock()
		return logger.CreateDetails(location, err)
	}

	spFsFile.Close()
	shared.MU.Unlock()

	var spFs shared.StorageProviderData

	err = json.Unmarshal(spFsBytes, &spFs)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	eightKBHashes := []string{}

	for i := 0; i < len(fileBytes); i += eightKB {
		hSum := sha256.Sum256(fileBytes[i : i+eightKB])
		eightKBHashes = append(eightKBHashes, hex.EncodeToString(hSum[:]))
	}

	_, fileTree, err := hash.CalcRoot(eightKBHashes)
	if err != nil {
		return logger.CreateDetails(location, err)
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

	contractRootHash, contractNonce, err := posInstance.GetUserRootHash(&bind.CallOpts{}, spAddress)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	var zeroHash [32]byte

	??ontractFsHashIsZero := bytes.Equal(zeroHash[:], contractRootHash[:])

	curentNonceIsZero := contractNonce.Cmp(big.NewInt(int64(0))) == 0

	firstProof := curentNonceIsZero && ??ontractFsHashIsZero

	nonceInt, err := strconv.Atoi(spFs.Nonce)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	if !firstProof {
		contractNonceIsBigger := contractNonce.Cmp(big.NewInt(int64(nonceInt))) == 1

		rootHashesEqual := bytes.Equal(fsRootHashBytes[:], contractRootHash[:])

		if contractNonceIsBigger && !rootHashesEqual {
			fmt.Println("fs root hash info is not valid!!!")
			return logger.CreateDetails(location, errors.New("fs root hash info is not valid"))
		}
	}

	treeToFsRoot = nil

	nonceHex := strconv.FormatInt(int64(nonceInt), 16)

	nonceBytes, err := hex.DecodeString(nonceHex)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	nonce32 := make([]byte, 32-len(nonceBytes))
	nonce32 = append(nonce32, nonceBytes...)

	fsRootNonceBytes := append(fsRootHashBytes[:], nonce32...)

	err = sign.Check(spAddress.String(), spFs.SignedFsRootNonceHash, sha256.Sum256(fsRootNonceBytes))
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	signedFSRootNonceHash, err := hex.DecodeString(spFs.SignedFsRootNonceHash)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	if signedFSRootNonceHash[len(signedFSRootNonceHash)-1] == 1 { //ecdsa version fix
		signedFSRootNonceHash[len(signedFSRootNonceHash)-1] = 28
	} else {
		signedFSRootNonceHash = signedFSRootNonceHash[:64]
	}

	fmt.Println("transactNonce", proofOpts.Nonce)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	proofOpts.Context = ctx

	_, err = posInstance.SendProof(proofOpts, common.HexToAddress(spAddress.String()), uint32(blockNum), fsRootHashBytes, uint64(nonceInt), signedFSRootNonceHash, fileBytes[:eightKB], proof)
	if err != nil {

		debug.FreeOSMemory()

		if err.Error() == "Transaction nonce is too low. Try incrementing the nonce." {
			proofOpts.Nonce = proofOpts.Nonce.Add(proofOpts.Nonce, big.NewInt(int64(1)))

			fmt.Println("Trying to prove with incremented nonce")

			_, err = posInstance.SendProof(proofOpts, common.HexToAddress(spAddress.String()), uint32(blockNum), fsRootHashBytes, uint64(nonceInt), signedFSRootNonceHash, fileBytes[:eightKB], proof)
			if err != nil {
				debug.FreeOSMemory()
				return logger.CreateDetails(location, err)
			}

		} else {
			debug.FreeOSMemory()
			proofOpts.Nonce = proofOpts.Nonce.Add(proofOpts.Nonce, big.NewInt(int64(1)))
			return logger.CreateDetails(location, err)
		}
	}

	debug.FreeOSMemory()
	proofOpts.Nonce = proofOpts.Nonce.Add(proofOpts.Nonce, big.NewInt(int64(1)))

	return nil
}

// ====================================================================================

func checkBalance(client *ethclient.Client, blockNum uint64, exitOnLow bool) (bool, error) {

	const location = "blckChain.checkBalance->"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	defer cancel()

	nodeBalance, err := client.BalanceAt(ctx, shared.NodeAddr, big.NewInt(int64(blockNum)))
	if err != nil {
		return false, logger.CreateDetails(location, err)
	}

	nodeBalanceIsLow := nodeBalance.Cmp(big.NewInt(1500000000000000)) == -1

	if nodeBalanceIsLow {
		fmt.Println("Insufficient funds for paying", CurrentNetwork, "transaction fees. Balance:", nodeBalance)
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

// ====================================================================================
// InitTrxOpts makes transaction options that are needed when sending request to smart contract.
func initTrxOpts(ctx context.Context, client *ethclient.Client, nodeAddr common.Address, password string, blockNum uint64) (*bind.TransactOpts, error) {
	const location = "blckChain.initTrxOpts->"

	transactNonce, err := client.NonceAt(ctx, nodeAddr, big.NewInt(int64(blockNum)))
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	chnID, err := client.ChainID(ctx)
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	opts := &bind.TransactOpts{
		From:  nodeAddr,
		Nonce: big.NewInt(int64(transactNonce)),
		Signer: func(a common.Address, t *types.Transaction) (*types.Transaction, error) {
			scryptN, scryptP := encryption.GetScryptParams()

			ks := keystore.NewKeyStore(paths.AccsDirPath, scryptN, scryptP)
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
		GasPrice: big.NewInt(20000000000), // 20 Gwei
		GasLimit: 1000000,
		Context:  ctx,
		NoSend:   false,
	}

	return opts, nil
}

// ====================================================================================

// MakeProof builds merkle tree proof. Passed tree value is an array of file hashes that are located on different levels of merkle tree.
// Returns slice of 32 bytes array for passing it to smart contract.
func makeProof(start []byte, tree [][][]byte) [][32]byte {
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

// ====================================================================================

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
