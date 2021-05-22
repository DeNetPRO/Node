package account

import (
	"bytes"
	"crypto/sha256"
	"dfile-secondary-node/shared"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type treeInfo struct {
	Nonce string     `json:"Nonce"`
	Tree  [][][]byte `json:"Tree"`
}

const eightKB = 8192

func SendProof() {
	file, err := os.Open("/home/r/dfile/accounts/0xA1c06ba6c5D0845727E7A2204FcF4f4C9636D8F4/storage/0xfC73A8Fe0eBcA03AE1481479C9132de97d757963/970c279ab111e55191d4e953013541e38b5f8f80e7c94ad13a9cd207ae2a5f31")
	if err != nil {
		log.Fatal("Fatal error")
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Fatal error")
	}

	fileFsTree, err := os.Open("/home/r/dfile/accounts/0xA1c06ba6c5D0845727E7A2204FcF4f4C9636D8F4/storage/0xfC73A8Fe0eBcA03AE1481479C9132de97d757963/tree.json")
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

	for i := 0; i < len(fileBytes); i += eightKB {
		hSum := sha256.Sum256(fileBytes[i : i+eightKB])
		eightKBHashes = append(eightKBHashes, hex.EncodeToString(hSum[:]))
	}

	_, fileTree, err := shared.CalcRootHash(eightKBHashes)
	if err != nil {
		log.Fatal("Fatal error")
	}

	fmt.Println(hex.EncodeToString(fileTree[len(fileTree)-1][0]))

	for _, v := range fsTreeStruct.Tree {
		for _, l := range v {
			fmt.Println(hex.EncodeToString(l))
		}
	}

	// proof := makeProof(tree[2][0], tree)

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

func makeProof(start []byte, tree [][][]byte) [][]byte {
	stage := 0
	proof := [][]byte{}

	var aPos int
	var bPos int

	for stage < len(tree)-1 {
		pos := getPos(start, tree[stage])
		if pos == -1 {
			stage++
			continue
		}

		if pos%2 != 0 {
			aPos = pos - 1
			bPos = pos
		} else {
			aPos = pos
			bPos = pos + 1
		}

		proof = append(proof, tree[stage][aPos])
		proof = append(proof, tree[stage][bPos])

		concatBytes := append(tree[stage][aPos], tree[stage][bPos]...)
		hSum := sha256.Sum256(concatBytes)

		start = hSum[:]
		stage++

	}

	proof = append(proof, start)

	return proof
}
