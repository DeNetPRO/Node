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
	file, err := os.Open("/home/r/dfile/accounts/0x546bf14Ba029D21359608182d0B9a4c9FacD7ed5/storage/0x9c20A547Ea5347e8a9AaC1A8f3e81D9C6600E4E0/338b83e118db0891ede737fc791dab8c0e95761404b9f5376cf2e70094979cb5")
	if err != nil {
		fmt.Println(err)
		log.Fatal("Fatal error")
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Fatal error")
	}

	fileFsTree, err := os.Open("/home/r/dfile/accounts/0x546bf14Ba029D21359608182d0B9a4c9FacD7ed5/storage/0x9c20A547Ea5347e8a9AaC1A8f3e81D9C6600E4E0/tree.json")
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

	hashFileRoot := fileTree[len(fileTree)-1][0]

	treeToFsRoot := [][][]byte{}

	for _, baseHash := range fsTreeStruct.Tree[0] {
		diff := bytes.Compare(hashFileRoot, baseHash)
		if diff == 0 {
			treeToFsRoot = append(treeToFsRoot, fileTree[:len(fileTree)-1]...)
			treeToFsRoot = append(treeToFsRoot, fsTreeStruct.Tree...)
		}
	}

	proof := makeProof(fileTree[2][0], treeToFsRoot)

	fmt.Println("proof", proof[len(proof)-1])

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
