package account

import (
	"bytes"
	"context"
	"crypto/sha256"
	POFstorage "dfile-secondary-node/POF_storage"
	"dfile-secondary-node/shared"
	"strconv"

	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type treeInfo struct {
	Nonce string     `json:"Nonce"`
	Tree  [][][]byte `json:"Tree"`
}

const eightKB = 8192

func SendProof() {

	pathToAcc := filepath.Join(shared.AccDir, DfileAcc.Address.String())

	pathToFile := filepath.Join(pathToAcc, "storage", "0x9c20A547Ea5347e8a9AaC1A8f3e81D9C6600E4E0", "338b83e118db0891ede737fc791dab8c0e95761404b9f5376cf2e70094979cb5")

	file, err := os.Open(pathToFile)
	if err != nil {
		log.Fatal("Fatal error")
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Fatal error")
	}

	pathToFsTree := filepath.Join(pathToAcc, "storage", "0x9c20A547Ea5347e8a9AaC1A8f3e81D9C6600E4E0", "tree.json")

	fileFsTree, err := os.Open(pathToFsTree)
	if err != nil {
		log.Fatal("Fatal error")
	}
	defer fileFsTree.Close()

	treeBytes, err := io.ReadAll(fileFsTree)
	if err != nil {
		log.Fatal("Fatal error")
	}

	var fsTreeStruct treeInfo

	err = json.Unmarshal(treeBytes, &fsTreeStruct)
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

	fsRootHash := fsTreeStruct.Tree[len(fsTreeStruct.Tree)-1][0]
	fmt.Println("fsRootHash", fsRootHash)

	treeToFsRoot := [][][]byte{}

	for _, baseHash := range fsTreeStruct.Tree[0] {
		diff := bytes.Compare(hashFileRoot, baseHash)
		if diff == 0 {
			treeToFsRoot = append(treeToFsRoot, fileTree[:len(fileTree)-1]...)
			treeToFsRoot = append(treeToFsRoot, fsTreeStruct.Tree...)
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

	nonceHex := strconv.FormatInt(1621758724, 16)
	nonceBytes, err := hex.DecodeString(nonceHex)
	if err != nil {
		log.Fatal(err)
	}

	fsRHashNonce := append(fsRootHash, nonceBytes...)

	encrKey := sha256.Sum256(DfileAcc.Address.Bytes())

	decryptedData, err := shared.DecryptAES(encrKey[:], DfileAcc.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	accPrivKey, err := crypto.HexToECDSA(hex.EncodeToString(decryptedData))
	if err != nil {
		log.Fatal(err)
	}

	hash := sha256.Sum256(fsRHashNonce)

	signedFSRootHash, err := crypto.Sign(hash[:], accPrivKey)
	if err != nil {
		log.Fatal(err)
	}

	dif, err := instance.SendProof(&bind.TransactOpts{}, DfileAcc.Address, uint32(blockNum.Size()), proof[len(proof)-1], 1621758724, signedFSRootHash, bytesToProve, proof)
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

func makeProof(start []byte, tree [][][]byte) [][32]byte {
	stage := 0
	proof := [][32]byte{}

	var aPos int
	var bPos int

	for stage < len(tree) {
		pos := getPos(start, tree[stage])
		if pos == -1 {
			break
		}

		if pos%2 != 0 {
			aPos = pos - 1
			bPos = pos
		} else {
			aPos = pos
			bPos = pos + 1
		}

		if len(tree[stage]) == 1 {
			tmp := [32]byte{}

			for i, v := range tree[stage][0] {
				tmp[i] = v
			}

			proof = append(proof, tmp)

			return proof
		}

		tmp1 := [32]byte{}
		for i, v := range tree[stage][aPos] {
			tmp1[i] = v
		}

		proof = append(proof, tmp1)

		tmp2 := [32]byte{}
		for i, v := range tree[stage][bPos] {
			tmp2[i] = v
		}

		proof = append(proof, tmp2)

		concatBytes := append(tree[stage][aPos], tree[stage][bPos]...)
		hSum := sha256.Sum256(concatBytes)

		start = hSum[:]
		stage++

	}

	return proof
}
