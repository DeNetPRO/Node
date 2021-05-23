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
	file, err := os.Open("/home/r/dfile/accounts/0x0F3eb0a4881F542a511154B2DF6334aB7d545753/storage/0x2839fE5865CcbB28326e6aD053af76E255B98B28/c137e04cb9b00d93a97ba50f8edfa2f9b61cf1b7f5af8517c17d2d155e5d1e1b")
	if err != nil {
		fmt.Println(err)
		log.Fatal("Fatal error")
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Fatal error")
	}

	fileFsTree, err := os.Open("/home/r/dfile/accounts/0x0F3eb0a4881F542a511154B2DF6334aB7d545753/storage/0x2839fE5865CcbB28326e6aD053af76E255B98B28/tree.json")
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

	fmt.Println("fileRoot", hex.EncodeToString(fileTree[len(fileTree)-1][0]))

	// for _, v := range fsTreeStruct.Tree {
	// 	for _, l := range v {
	// 		fmt.Println(hex.EncodeToString(l))
	// 	}
	// }

	proof := makeProof(fileTree[2][0], fileTree)

	fmt.Println("proof", hex.EncodeToString(proof[len(proof)-1]))

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

	for stage < len(tree) {
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

		if len(tree[stage]) == 1 {
			proof = append(proof, tree[stage][0])

			return proof
		}

		proof = append(proof, tree[stage][aPos])
		proof = append(proof, tree[stage][bPos])

		concatBytes := append(tree[stage][aPos], tree[stage][bPos]...)
		hSum := sha256.Sum256(concatBytes)

		start = hSum[:]
		stage++

	}

	return proof
}
