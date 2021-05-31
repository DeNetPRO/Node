package mining

import (
	"bytes"
	"context"
	"crypto/sha256"
	POFstorage "dfile-secondary-node/POF_storage"
	"dfile-secondary-node/account"
	"dfile-secondary-node/shared"
	"io/fs"
	"regexp"
	"time"

	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type StorageInfo struct {
	Nonce        string     `json:"nonce"`
	SignedFsRoot string     `json:"signedFsRoot"`
	Tree         [][][]byte `json:"tree"`
}

const eightKB = 8192

func StartMining() {

	for {
		time.Sleep(time.Second * 1)
		pathToAccStorage := filepath.Join(shared.AccsDirPath, account.DfileAcc.Address.String(), shared.StorageDirName)

		storageAddresses := []string{}

		re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

		err := filepath.WalkDir(pathToAccStorage,
			func(path string, info fs.DirEntry, err error) error {
				if err != nil {
					log.Fatal("Fatal error")
				}

				if re.MatchString(info.Name()) {
					storageAddresses = append(storageAddresses, info.Name())
				}

				return nil
			})
		if err != nil {
			log.Fatal("Fatal error")
		}

		if len(storageAddresses) == 0 {
			continue
		}

		client, err := ethclient.Dial("https://kovan.infura.io/v3/a4a45777ca65485d983c278291e322f2")
		if err != nil {
			log.Fatal(err)
		}

		tokenAddress := common.HexToAddress("0x2E8630780A231E8bCf12Ba1172bEB9055deEBF8B")
		instance, err := POFstorage.NewStore(tokenAddress, client)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range storageAddresses {
			commonAddr := common.HexToAddress(v)
			address, rew, rew1, err := instance.GetUserRewardInfo(&bind.CallOpts{}, commonAddr)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(address, rew, rew1)
		}

	}

}

func SendProof() {

	pathToAcc := filepath.Join(shared.AccsDirPath, account.DfileAcc.Address.String())

	pathToFile := filepath.Join(pathToAcc, shared.StorageDirName, "0x9c20A547Ea5347e8a9AaC1A8f3e81D9C6600E4E0", "338b83e118db0891ede737fc791dab8c0e95761404b9f5376cf2e70094979cb5")

	file, err := os.Open(pathToFile)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Fatal error")
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Fatal error")
	}

	pathToFsTree := filepath.Join(pathToAcc, shared.StorageDirName, "0x9c20A547Ea5347e8a9AaC1A8f3e81D9C6600E4E0", "tree.json")

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
	instance, err := POFstorage.NewStore(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	blockNum, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	signedFSRootHash, err := hex.DecodeString(storageFsStruct.SignedFsRoot)
	if err != nil {
		log.Fatal(err)
	}

	dif, err := instance.SendProof(&bind.TransactOpts{}, account.DfileAcc.Address, uint32(blockNum.Size()), proof[len(proof)-1], 1621758724, signedFSRootHash, bytesToProve, proof)
	if err != nil {
		log.Fatal(err)
	}

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
