package server

import (
	"bytes"
	"crypto/sha256"
	blockchainprovider "dfile-secondary-node/blockchain_provider"
	"dfile-secondary-node/config"
	"dfile-secondary-node/shared"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const gbBytes = int64(1024 * 1024 * 1024)
const oneHunderdMBBytes = int64(1024 * 1024 * 100)

func Start(address, port string) {

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

	nodeAddr, err := shared.DecryptNodeAddr()
	if err != nil {
		http.Error(w, "Couldn't parse address", 500)
		return
	}

	pathToConfig := filepath.Join(shared.AccsDirPath, nodeAddr.String(), shared.ConfDirName, "config.json")

	confFile, err := os.OpenFile(pathToConfig, os.O_RDWR, 0755)
	if err != nil {
		http.Error(w, "Account config problem", 500)
		return
	}
	defer confFile.Close()

	fileBytes, err := io.ReadAll(confFile)
	if err != nil {
		http.Error(w, "Account config problem", 500)
		return
	}

	var dFileConf config.SecondaryNodeConfig

	err = json.Unmarshal(fileBytes, &dFileConf)
	if err != nil {
		http.Error(w, "Account config problem", 500)
		return
	}

	sharedSpaceInBytes := int64(dFileConf.StorageLimit) * gbBytes

	reqFiles := req.MultipartForm.File["files"]

	var filesTotalSize int64

	for _, reqFile := range reqFiles {
		filesTotalSize += reqFile.Size
	}

	dFileConf.UsedStorageSpace += filesTotalSize

	if dFileConf.UsedStorageSpace > sharedSpaceInBytes {
		if err != nil {
			fmt.Println("Insufficient memory avaliable")
			http.Error(w, "Insufficient memory avaliable", 400)
			return
		}
	}

	avaliableSpaceLeft := sharedSpaceInBytes - dFileConf.UsedStorageSpace

	if avaliableSpaceLeft < oneHunderdMBBytes {
		fmt.Println("Shared storage memory is running low", avaliableSpaceLeft/(1024*1024), "MB of space is avaliable")
		fmt.Println("You may need additional space for mining. Total shared space can be changed in account configuration")
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

	nonceInt, err := strconv.Atoi(nonce[0])
	if err != nil {
		http.Error(w, "File saving problem", 400)
		return
	}

	nonceHex := strconv.FormatInt(int64(nonceInt), 16)

	nonceBytes, err := hex.DecodeString(nonceHex)
	if err != nil {
		http.Error(w, "File saving problem", 400)
		return
	}

	nonce32 := make([]byte, 32-len(nonceBytes))
	nonce32 = append(nonce32, nonceBytes...)

	fsRootBytes, err := hex.DecodeString(fsRootHash)
	if err != nil {
		http.Error(w, "File saving problem", 400)
		return
	}

	fsRootNonceBytes := append(fsRootBytes, nonce32...)

	hash := sha256.Sum256(fsRootNonceBytes)

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

	addressPath := filepath.Join(shared.AccsDirPath, nodeAddr.String(), shared.StorageDirName, storageProviderAddress[0])

	stat, err := os.Stat(addressPath)
	if err != nil {
		errPart := strings.Split(err.Error(), ":")

		if strings.Trim(errPart[1], " ") != "no such file or directory" {
			fmt.Println(err)

			http.Error(w, "File saving problem", 500)
			return
		}

	}

	if stat == nil {
		err = os.Mkdir(addressPath, 0700)
		if err != nil {
			fmt.Println(err)

			http.Error(w, "File saving problem", 500)
			return
		}
	}

	treeFile, err := os.Create(filepath.Join(addressPath, "tree.json"))
	if err != nil {
		fmt.Println(err)

		http.Error(w, "File saving problem", 500)
		return
	}
	defer treeFile.Close()

	tree := blockchainprovider.StorageInfo{
		Nonce:        nonce[0],
		SignedFsRoot: signedFsRootHash[0],
		Tree:         fsTree,
	}

	js, err := json.Marshal(tree)
	if err != nil {
		http.Error(w, "File saving problem", 500)
		return
	}

	_, err = treeFile.Write(js)
	if err != nil {
		http.Error(w, "File saving problem", 500)
		return
	}

	treeFile.Sync()

	const eightKB = 8192

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

	configJson, err := json.Marshal(dFileConf)
	if err != nil {
		http.Error(w, "Couldn't update config file", 500)
		return
	}

	err = confFile.Truncate(0)
	if err != nil {
		http.Error(w, "Couldn't update config file", 500)
		return
	}

	_, err = confFile.Seek(0, 0)
	if err != nil {
		http.Error(w, "Couldn't update config file", 500)
		return
	}

	_, err = confFile.Write(configJson)
	if err != nil {
		http.Error(w, "Couldn't update config file", 500)
		return
	}

	confFile.Sync()

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

// ====================================================================================

func ServeFiles(w http.ResponseWriter, req *http.Request) {

	nodeAddr, err := shared.DecryptNodeAddr()
	if err != nil {
		http.Error(w, "Couldn't parse address", 500)
		return
	}

	vars := mux.Vars(req)
	storageProviderAddress := vars["address"]
	fileName := vars["fileName"]

	signatureFromReq := vars["signature"]

	signature, err := hex.DecodeString(signatureFromReq)
	if err != nil {
		http.Error(w, "File saving problem", 400)
		return
	}

	hash := sha256.Sum256([]byte(fileName + storageProviderAddress))

	sigPublicKey, err := crypto.SigToPub(hash[:], signature)
	if err != nil {
		http.Error(w, "File saving problem", 400)
		return
	}

	signatureAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if storageProviderAddress != signatureAddress.String() {
		http.Error(w, "Wrong signature", http.StatusForbidden)
		return
	}

	pathToFile := filepath.Join(shared.AccsDirPath, nodeAddr.String(), shared.StorageDirName, storageProviderAddress, fileName)
	http.ServeFile(w, req, pathToFile)
}
