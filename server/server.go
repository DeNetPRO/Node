package server

import (
	"bytes"
	"crypto/sha256"
	"dfile-secondary-node/shared"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type treeInfo struct {
	Nonce string     `json:"Nonce"`
	Tree  [][][]byte `json:"Tree"`
}

var AccountAddress string

func Start(address, port string) {

	AccountAddress = address

	r := mux.NewRouter()

	r.HandleFunc("/upload", SaveFiles).Methods("POST")
	r.HandleFunc("/download/{fileKey}", ServeFiles).Methods("GET")

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodOptions,
		},

		AllowedHeaders: []string{
			"Accept",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
		},
	})

	fmt.Println("Dfile node is ready and started listening on port: " + port)

	err := http.ListenAndServe(":"+port, corsOpts.Handler(checkSignature(r)))
	if err != nil {
		panic(err)
	}

}

// ====================================================================================

func checkSignature(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// splittedPath := strings.Split(r.URL.Path, "/")
		// signature := splittedPath[len(splittedPath)-1]
		// splittedPath = splittedPath[:len(splittedPath)-1]
		// reqURL := strings.Join(splittedPath, "/")

		// verified, err := verifySignature(sessionKeyBytes, reqURL, signature)
		// if err != nil {
		// 	http.Error(w, "session key verification error", 500)
		// 	return
		// }

		// if !verified {
		// 	http.Error(w, "wrong session key", http.StatusForbidden)
		// }

		h.ServeHTTP(w, r)
	})
}

// ========================================================================================================

func SaveFiles(w http.ResponseWriter, req *http.Request) {

	err := req.ParseMultipartForm(1 << 20) // maxMemory 32MB
	if err != nil {
		http.Error(w, "Parse multiform problem", 400)
		return
	}

	fs := req.MultipartForm.Value["fs"]

	fsHashes := make([]string, len(fs))
	copy(fsHashes, fs)

	sort.Strings(fsHashes)

	fsRootHash, fsTree, err := shared.CalcRootHash(fsHashes)
	if err != nil {
		http.Error(w, "File saving problem", 400)
		return
	}

	signedFsRootHash := req.MultipartForm.Value["fsRootHash"]

	signature, err := hex.DecodeString(signedFsRootHash[0])
	if err != nil {
		http.Error(w, "File saving problem", 400)
		return
	}

	nonce := req.MultipartForm.Value["nonce"]

	hash := sha256.Sum256([]byte(fsRootHash + nonce[0]))

	sigPublicKey, err := crypto.SigToPub(hash[:], signature)
	if err != nil {
		http.Error(w, "File saving problem", 400)
		return
	}

	storageProviderAddress := req.MultipartForm.Value["address"]

	senderAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if storageProviderAddress[0] != fmt.Sprint(senderAddress) {
		http.Error(w, "Wrong signature", http.StatusForbidden)
		return
	}

	addressPath := filepath.Join(shared.AccDir, AccountAddress, "storage", storageProviderAddress[0])

	stat, err := os.Stat(addressPath)
	if err != nil {
		errPart := strings.Split(err.Error(), ":")

		if strings.Trim(errPart[1], " ") != "no such file or directory" {
			http.Error(w, "File saving problem", 500)
			return
		}

	}

	if stat == nil {
		err = os.Mkdir(addressPath, 0700)
		if err != nil {
			http.Error(w, "File saving problem", 500)
			return
		}
	}

	treeFile, err := os.Create(filepath.Join(addressPath, "tree.json"))
	if err != nil {
		http.Error(w, "File saving problem", 500)
		return
	}
	defer treeFile.Close()

	tree := treeInfo{
		Nonce: nonce[0],
		Tree:  fsTree,
	}

	js, err := json.Marshal(tree)
	if err != nil {
		http.Error(w, "File saving problem", 500)
		return
	}

	treeFile.Write(js)
	treeFile.Sync()

	const eightKB = 8192

	reqFiles := req.MultipartForm.File["files"]

	oneMBHashes := []string{}

	for _, reqFile := range reqFiles {
		eightKBHashes := []string{}

		var buf bytes.Buffer

		rqFile, err := reqFile.Open()
		if err != nil {
			http.Error(w, "File check problem", 500)
			return
		}

		_, err = io.Copy(&buf, rqFile)
		if err != nil {
			rqFile.Close()
			http.Error(w, "File check problem", 500)
			return
		}

		rqFile.Close()

		bufBytes := buf.Bytes()

		for i := 0; i < len(bufBytes); i += eightKB {
			hSum := sha256.Sum256(bufBytes[i : i+eightKB])
			eightKBHashes = append(eightKBHashes, hex.EncodeToString(hSum[:]))
		}

		oneMBHash, _, err := shared.CalcRootHash(eightKBHashes)
		if err != nil {
			http.Error(w, "Wrong file", 400)
			return
		}

		if reqFile.Filename != oneMBHash {
			http.Error(w, "Wrong file", 400)
			return
		}

		oneMBHashes = append(oneMBHashes, oneMBHash)

	}

	fsContainsFile := false

	var fileRootHash string

	if len(oneMBHashes) == 1 {
		fileRootHash = oneMBHashes[0]
	} else {
		sort.Strings(oneMBHashes)
		fileRootHash, _, err = shared.CalcRootHash(oneMBHashes)
		if err != nil {
			http.Error(w, "Wrong file", 400)
			return
		}
	}

	for _, k := range fs {
		if k == fileRootHash {
			fsContainsFile = true
		}
	}

	if !fsContainsFile {
		http.Error(w, "Wrong file", 400)
		return
	}

	for _, reqFile := range reqFiles {

		rqFile, err := reqFile.Open()
		if err != nil {
			http.Error(w, "File saving problem", 500)
			return
		}
		defer rqFile.Close()

		pathToFile := filepath.Join(addressPath, reqFile.Filename)

		newFile, err := os.Create(pathToFile)
		if err != nil {
			http.Error(w, "File saving problem", 500)
			return
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, rqFile)
		if err != nil {
			http.Error(w, "File saving problem", 500)
			return
		}

		newFile.Sync()
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

// ====================================================================================

func ServeFiles(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	storageProviderAddress := vars["address"]
	name := vars["fileName"]

	pathToFile := filepath.Join(shared.AccDir, AccountAddress, "storage", storageProviderAddress, name)
	http.ServeFile(w, req, pathToFile)
}
