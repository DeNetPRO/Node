package server

import (
	"bytes"
	"crypto/sha256"

	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"git.denetwork.xyz/dfile/dfile-secondary-node/config"
	"git.denetwork.xyz/dfile/dfile-secondary-node/encryption"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"git.denetwork.xyz/dfile/dfile-secondary-node/update"
	"git.denetwork.xyz/dfile/dfile-secondary-node/upnp"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type updatedFsInfo struct {
	NewFs            []string
	Nonce            string
	SignedFsRootHash string
}

const gbBytes = int64(1024 * 1024 * 1024)
const oneHunderdMBBytes = int64(1024 * 1024 * 100)
const serverStartFatalMessage = "Couldn't start server"

func Start(address, port string) {
	const logInfo = "server.Start->"
	r := mux.NewRouter()

	r.HandleFunc("/upload/{size}", SaveFiles).Methods("POST")
	r.HandleFunc("/download/{address}/{fileKey}/{signature}", ServeFiles).Methods("GET")
	r.HandleFunc("/update_fs/{address}/{signedFsys}", updateFsInfo).Methods("POST")

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

	intPort, err := strconv.Atoi(port)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		log.Fatal(serverStartFatalMessage)
	}

	if upnp.InternetDevice != nil {
		err = upnp.InternetDevice.Forward(uint16(intPort), "node")
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
		}
		defer upnp.InternetDevice.Clear(uint16(intPort))
	}

	fmt.Println("Dfile node is ready and started listening on port: " + port)

	err = http.ListenAndServe(":"+port, corsOpts.Handler(checkSignature(r)))
	if err != nil {
		log.Fatal(serverStartFatalMessage)
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
	const logInfo = "server.SaveFiles->"

	nodeAddr, err := encryption.DecryptNodeAddr()
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Couldn't parse address", 500)
		return
	}

	vars := mux.Vars(req)
	fileSize := vars["size"]

	intFileSize, err := strconv.Atoi(fileSize)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Couldn't parse address", 500)
		return
	}

	if intFileSize == 0 {
		logger.Log(logger.CreateDetails(logInfo, errors.New("file size is 0")))
		http.Error(w, "file size is 0", 400)
		return
	}

	pathToConfig := filepath.Join(paths.AccsDirPath, nodeAddr.String(), paths.ConfDirName, "config.json")

	shared.MU.Lock()
	confFile, err := os.OpenFile(pathToConfig, os.O_RDWR, 0755)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Account config problem", 500)
		return
	}
	defer confFile.Close()

	fileBytes, err := io.ReadAll(confFile)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Account config problem", 500)
		return
	}

	var dFileConf config.SecondaryNodeConfig

	err = json.Unmarshal(fileBytes, &dFileConf)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Account config problem", 500)
		return
	}

	sharedSpaceInBytes := int64(dFileConf.StorageLimit) * gbBytes

	dFileConf.UsedStorageSpace += int64(intFileSize)

	if dFileConf.UsedStorageSpace > sharedSpaceInBytes {
		shared.MU.Unlock()
		err := errors.New("insufficient memory avaliable")
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, err.Error(), 400)
		return
	}

	avaliableSpaceLeft := sharedSpaceInBytes - dFileConf.UsedStorageSpace

	if avaliableSpaceLeft < oneHunderdMBBytes {
		fmt.Println("Shared storage memory is running low,", avaliableSpaceLeft/(1024*1024), "MB of space is avaliable")
		fmt.Println("You may need additional space for mining. Total shared space can be changed in account configuration")
	}

	err = config.Save(confFile, dFileConf)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Couldn't update config file", 500)
		return
	}
	confFile.Close()
	shared.MU.Unlock()

	err = req.ParseMultipartForm(1 << 20) // maxMemory 32MB
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "Parse multiform problem", 400)
		return
	}

	fs := req.MultipartForm.Value["fs"]

	sort.Strings(fs)

	fsRootHash, fsTree, err := shared.CalcRootHash(fs)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 400)
		return
	}

	signedFsRootHash := req.MultipartForm.Value["fsRootHash"]

	signature, err := hex.DecodeString(signedFsRootHash[0])
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 400)
		return
	}

	nonce := req.MultipartForm.Value["nonce"]

	nonceInt, err := strconv.Atoi(nonce[0])
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 400)
		return
	}

	nonceHex := strconv.FormatInt(int64(nonceInt), 16)

	nonceBytes, err := hex.DecodeString(nonceHex)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 400)
		return
	}

	nonce32 := make([]byte, 32-len(nonceBytes))
	nonce32 = append(nonce32, nonceBytes...)

	fsRootBytes, err := hex.DecodeString(fsRootHash)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 400)
		return
	}

	fsRootNonceBytes := append(fsRootBytes, nonce32...)

	hash := sha256.Sum256(fsRootNonceBytes)

	sigPublicKey, err := crypto.SigToPub(hash[:], signature)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 400)
		return
	}

	storageProviderAddress := req.MultipartForm.Value["address"]

	senderAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if storageProviderAddress[0] != fmt.Sprint(senderAddress) {
		err := errors.New("wrong signature")
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	addressPath := filepath.Join(paths.AccsDirPath, nodeAddr.String(), paths.StorageDirName, storageProviderAddress[0])

	stat, err := os.Stat(addressPath)
	err = shared.CheckStatErr(err)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}

	if stat == nil {
		err = os.Mkdir(addressPath, 0700)
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			restoreMemoryInfo(pathToConfig, intFileSize)
			http.Error(w, "File saving problem", 500)
			return
		}
	}

	shared.MU.Lock()
	treeFile, err := os.Create(filepath.Join(addressPath, "tree.json"))
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}
	defer treeFile.Close()

	tree := shared.StorageInfo{
		Nonce:        nonce[0],
		SignedFsRoot: signedFsRootHash[0],
		Tree:         fsTree,
	}

	js, err := json.Marshal(tree)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}

	_, err = treeFile.Write(js)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}

	treeFile.Sync()
	treeFile.Close()
	shared.MU.Unlock()

	reqFileParts := req.MultipartForm.File["files"]

	const eightKB = 8192

	oneMBHashes := make([]string, 0, len(reqFileParts))

	for _, reqFilePart := range reqFileParts {

		eightKBHashes := make([]string, 0, 128)

		var buf bytes.Buffer

		rqFile, err := reqFilePart.Open()
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			restoreMemoryInfo(pathToConfig, intFileSize)
			http.Error(w, "File check problem", 500)
			return
		}

		_, err = io.Copy(&buf, rqFile)
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			restoreMemoryInfo(pathToConfig, intFileSize)
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
			logger.Log(logger.CreateDetails(logInfo, err))
			restoreMemoryInfo(pathToConfig, intFileSize)
			http.Error(w, "Wrong file", 400)
			return
		}

		if reqFilePart.Filename != oneMBHash {
			err := errors.New("wrong file")
			logger.Log(logger.CreateDetails(logInfo, err))
			restoreMemoryInfo(pathToConfig, intFileSize)
			http.Error(w, err.Error(), 400)
			return
		}

		oneMBHashes = append(oneMBHashes, oneMBHash)

	}

	fsContainsFile := false

	var wholeFileHash string

	if len(oneMBHashes) == 1 {
		wholeFileHash = oneMBHashes[0]
	} else {
		sort.Strings(oneMBHashes)
		wholeFileHash, _, err = shared.CalcRootHash(oneMBHashes)
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			restoreMemoryInfo(pathToConfig, intFileSize)
			http.Error(w, "Wrong file", 400)
			return
		}
	}

	for _, fileHash := range fs {
		if fileHash == wholeFileHash {
			fsContainsFile = true
		}
	}

	if !fsContainsFile {
		err := errors.New("wrong file")
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, err.Error(), 400)
		return
	}

	count := 1
	total := len(oneMBHashes)

	for _, reqFilePart := range reqFileParts {
		rqFile, err := reqFilePart.Open()
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			deleteFileParts(addressPath, oneMBHashes)
			restoreMemoryInfo(pathToConfig, intFileSize)
			http.Error(w, "File saving problem", 500)
			return
		}
		defer rqFile.Close()

		pathToFile := filepath.Join(addressPath, reqFilePart.Filename)

		newFile, err := os.Create(pathToFile)
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			deleteFileParts(addressPath, oneMBHashes)
			restoreMemoryInfo(pathToConfig, intFileSize)
			http.Error(w, "File saving problem", 500)
			return
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, rqFile)
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			deleteFileParts(addressPath, oneMBHashes)
			restoreMemoryInfo(pathToConfig, intFileSize)
			http.Error(w, "File saving problem", 500)
			return
		}

		logger.Log("Saved file " + reqFilePart.Filename + " (" + fmt.Sprint(count) + "/" + fmt.Sprint(total) + ")" + " from " + storageProviderAddress[0]) //TODO remove

		newFile.Sync()
		rqFile.Close()
		newFile.Close()

		count++
	}

	go update.FsInfo(nodeAddr.String(), storageProviderAddress[0], fsRootHash, nonce[0], fs, nonce32, fsRootNonceBytes)

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

// ====================================================================================

func deleteFileParts(addressPath string, fileHashes []string) {

	logger.Log("deleting file parts after error...")

	for _, hash := range fileHashes {
		pathToFile := filepath.Join(addressPath, hash)

		os.Remove(pathToFile)
	}
}

// ====================================================================================

func ServeFiles(w http.ResponseWriter, req *http.Request) {
	const logInfo = "server.ServeFiles->"

	nodeAddr, err := encryption.DecryptNodeAddr()
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Couldn't parse address", 500)
		return
	}

	vars := mux.Vars(req)
	addressFromReq := vars["address"]
	fileKey := vars["fileKey"]
	signatureFromReq := vars["signature"]

	signature, err := hex.DecodeString(signatureFromReq)
	if err != nil {
		http.Error(w, "File serving problem", 400)
		return
	}

	hash := sha256.Sum256([]byte(fileKey + addressFromReq))

	sigPublicKey, err := crypto.SigToPub(hash[:], signature)
	if err != nil {
		http.Error(w, "File serving problem", 400)
		return
	}

	signatureAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if addressFromReq != signatureAddress.String() {
		http.Error(w, "Wrong signature", http.StatusForbidden)
		return
	}

	pathToFile := filepath.Join(paths.AccsDirPath, nodeAddr.String(), paths.StorageDirName, addressFromReq, fileKey)

	_, err = os.Stat(pathToFile)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	logger.Log("serving file: " + fileKey)

	http.ServeFile(w, req, pathToFile)
}

// ====================================================================================

func updateFsInfo(w http.ResponseWriter, req *http.Request) {
	const logInfo = "server.UpdateFsInfo->"

	nodeAddr, err := encryption.DecryptNodeAddr()
	if err != nil {
		http.Error(w, "Internal node error", 500)
		return
	}

	vars := mux.Vars(req)
	addressFromReq := vars["address"]
	signedFsys := vars["signedFsys"]

	addressPath := filepath.Join(paths.AccsDirPath, nodeAddr.String(), paths.StorageDirName, addressFromReq)

	_, err = os.Stat(addressPath)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, errors.New("account not found")))
		return
	}

	fsysSignature, err := hex.DecodeString(signedFsys)
	if err != nil {
		http.Error(w, "Wrong signature", 400)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Internal error", 500)
		return
	}
	defer req.Body.Close()

	var updatedFs updatedFsInfo

	err = json.Unmarshal(body, &updatedFs)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Internal error", 500)
		return
	}

	nonceInt, err := strconv.Atoi(updatedFs.Nonce)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Fs info update problem", 400)
		return
	}

	nonceHex := strconv.FormatInt(int64(nonceInt), 16)

	nonceBytes, err := hex.DecodeString(nonceHex)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Fs info update problem", 400)
		return
	}

	nonce32 := make([]byte, 32-len(nonceBytes))
	nonce32 = append(nonce32, nonceBytes...)

	sort.Strings(updatedFs.NewFs)

	concatFsHashes := ""

	for _, hash := range updatedFs.NewFs {
		concatFsHashes += hash
	}

	fsTreeNonceBytes := append([]byte(concatFsHashes), nonce32...)

	fsTreeNonceSha := sha256.Sum256(fsTreeNonceBytes)

	sigPublicKey, err := crypto.SigToPub(fsTreeNonceSha[:], fsysSignature)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Wrong signature", 400)
		return
	}

	signatureAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if addressFromReq != signatureAddress.String() {
		logger.Log(logger.CreateDetails(logInfo, errors.New("wrong signature")))
		http.Error(w, "Wrong signature", http.StatusForbidden)
		return
	}

	fsRootHash, fsTree, err := shared.CalcRootHash(updatedFs.NewFs)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Fs info update problem", 400)
		return
	}

	rootSignature, err := hex.DecodeString(updatedFs.SignedFsRootHash)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Fs info update problem", 400)
		return
	}

	fsRootBytes, err := hex.DecodeString(fsRootHash)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Fs info update problem", 400)
		return
	}

	fsRootNonceBytes := append(fsRootBytes, nonce32...)

	hash := sha256.Sum256(fsRootNonceBytes)

	sigPublicKey, err = crypto.SigToPub(hash[:], rootSignature)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Fs info update problem", 400)
		return
	}

	signatureAddress = crypto.PubkeyToAddress(*sigPublicKey)

	if addressFromReq != signatureAddress.String() {
		logger.Log(logger.CreateDetails(logInfo, errors.New("wrong signature")))
		http.Error(w, "Wrong signature", http.StatusForbidden)
		return
	}

	shared.MU.Lock()

	treeFile, err := os.Create(filepath.Join(addressPath, "tree.json"))
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, errors.New("wrong signature")))
		http.Error(w, "Fs info update problem", 500)
		return
	}
	defer treeFile.Close()

	tree := shared.StorageInfo{
		Nonce:        updatedFs.Nonce,
		SignedFsRoot: updatedFs.SignedFsRootHash,
		Tree:         fsTree,
	}

	js, err := json.Marshal(tree)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Fs info update problem", 500)
		return
	}

	_, err = treeFile.Write(js)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Fs info update problem", 500)
		return
	}

	treeFile.Sync()

	logger.Log("Updated fs info")

	shared.MU.Unlock()

}

// ====================================================================================

func restoreMemoryInfo(pathToConfig string, intFileSize int) {
	logInfo := "server.restoreMemoryInfo->"

	shared.MU.Lock()
	confFile, err := os.OpenFile(pathToConfig, os.O_RDWR, 0755)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		return
	}
	defer confFile.Close()

	fileBytes, err := io.ReadAll(confFile)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		return
	}

	var dFileConf config.SecondaryNodeConfig

	err = json.Unmarshal(fileBytes, &dFileConf)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		return
	}

	dFileConf.UsedStorageSpace -= int64(intFileSize)

	err = config.Save(confFile, dFileConf)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		return
	}
	shared.MU.Unlock()
}
