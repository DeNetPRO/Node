package server

import (
	"bytes"
	"context"
	"crypto/sha256"
	"math/big"
	"mime/multipart"
	"os/signal"
	"strings"

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

	blockchainprovider "git.denetwork.xyz/dfile/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/dfile/dfile-secondary-node/config"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	nodeAbi "git.denetwork.xyz/dfile/dfile-secondary-node/node_abi"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"git.denetwork.xyz/dfile/dfile-secondary-node/upnp"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

type updatedFsInfo struct {
	NewFs            []string
	Nonce            string
	SignedFsRootHash string
}

type NodeAddressResponse struct {
	NodeAddress string `json:"node_address"`
}

const gbBytes = int64(1024 * 1024 * 1024)
const oneHunderdMBBytes = int64(1024 * 1024 * 100)
const serverStartFatalMessage = "Couldn't start server"

func Start(port string) {
	const logInfo = "server.Start->"
	r := mux.NewRouter()

	r.HandleFunc("/upload/{size}", SaveFiles).Methods("POST")
	r.HandleFunc("/download/{spAddress}/{fileKey}/{signature}", ServeFiles).Methods("GET")
	r.HandleFunc("/update_fs/{spAddress}/{signedFsys}", updateFsInfo).Methods("POST")
	r.HandleFunc("/copy/{size}", CopyFile).Methods("POST")
	r.HandleFunc("/backup/{size}", BackUp).Methods("POST")

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
		upnp.InternetDevice.Forward(intPort, "node")
		defer upnp.InternetDevice.Close(intPort)
	}

	fmt.Println("Dfile node is ready and started listening on port: " + port)

	server := http.Server{
		Addr:    ":" + port,
		Handler: corsOpts.Handler(checkSignature(r)),
	}

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			log.Fatal(serverStartFatalMessage)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	err = server.Shutdown(context.Background())
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		log.Fatal(err)
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

	pathToConfig := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.ConfDirName, "config.json")

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

	var nodeConfig config.SecondaryNodeConfig

	err = json.Unmarshal(fileBytes, &nodeConfig)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Account config problem", 500)
		return
	}

	sharedSpaceInBytes := int64(nodeConfig.StorageLimit) * gbBytes

	nodeConfig.UsedStorageSpace += int64(intFileSize)

	if nodeConfig.UsedStorageSpace > sharedSpaceInBytes {
		shared.MU.Unlock()
		err := errors.New("insufficient memory avaliable")
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, err.Error(), 400)
		return
	}

	avaliableSpaceLeft := sharedSpaceInBytes - nodeConfig.UsedStorageSpace

	if avaliableSpaceLeft < oneHunderdMBBytes {
		fmt.Println("Shared storage memory is running low,", avaliableSpaceLeft/(1024*1024), "MB of space is avaliable")
		fmt.Println("You may need additional space for mining. Total shared space can be changed in account configuration")
	}

	err = config.Save(confFile, nodeConfig)
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

	addressPath := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, storageProviderAddress[0])

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
	spFsFile, err := os.Create(filepath.Join(addressPath, "tree.json"))
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}
	defer spFsFile.Close()

	spFs := shared.StorageProviderFs{
		Nonce:        nonce[0],
		SignedFsRoot: signedFsRootHash[0],
		Tree:         fsTree,
	}

	js, err := json.Marshal(spFs)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}

	_, err = spFsFile.Write(js)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}

	spFsFile.Sync()
	spFsFile.Close()
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

	vars := mux.Vars(req)
	spAddress := vars["spAddress"]
	fileKey := vars["fileKey"]
	signatureFromReq := vars["signature"]

	signature, err := hex.DecodeString(signatureFromReq)
	if err != nil {
		http.Error(w, "File serving problem", 400)
		return
	}

	hash := sha256.Sum256([]byte(fileKey + spAddress))

	sigPublicKey, err := crypto.SigToPub(hash[:], signature)
	if err != nil {
		http.Error(w, "File serving problem", 400)
		return
	}

	signatureAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if spAddress != signatureAddress.String() {
		http.Error(w, "Wrong signature", http.StatusForbidden)
		return
	}

	pathToFile := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spAddress, fileKey)

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

	const httpErrorMsg = "Fs info update problem"

	vars := mux.Vars(req)
	spAddress := vars["spAddress"]
	signedFsys := vars["signedFsys"]

	addressPath := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spAddress)

	_, err := os.Stat(addressPath)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, errors.New("no files of "+spAddress)))
		return
	}

	shared.MU.Lock()
	spFsFile, err := os.Open(filepath.Join(addressPath, "tree.json"))
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, httpErrorMsg, 500)
		return
	}
	defer spFsFile.Close()

	var spFs shared.StorageProviderFs

	fileBytes, err := io.ReadAll(spFsFile)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, httpErrorMsg, 500)
		return
	}

	err = json.Unmarshal(fileBytes, &spFs)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, httpErrorMsg, 500)
		return
	}

	spFsFile.Close()
	shared.MU.Unlock()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, httpErrorMsg, 500)
		return
	}
	defer req.Body.Close()

	var updatedFs updatedFsInfo

	err = json.Unmarshal(body, &updatedFs)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, httpErrorMsg, 500)
		return
	}

	newNonceInt, err := strconv.Atoi(updatedFs.Nonce)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, httpErrorMsg, 400)
		return
	}

	currentNonceInt, err := strconv.Atoi(spFs.Nonce)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, httpErrorMsg, 400)
		return
	}

	if newNonceInt < currentNonceInt {
		logger.Log(spAddress + " fs info is up to date")
		http.Error(w, httpErrorMsg, 400)
		return
	}

	nonceHex := strconv.FormatInt(int64(newNonceInt), 16)

	nonceBytes, err := hex.DecodeString(nonceHex)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, httpErrorMsg, 400)
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

	fsysSignature, err := hex.DecodeString(signedFsys)
	if err != nil {
		http.Error(w, httpErrorMsg, 500)
		return
	}

	sigPublicKey, err := crypto.SigToPub(fsTreeNonceSha[:], fsysSignature)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Wrong signature", 400)
		return
	}

	signatureAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if spAddress != signatureAddress.String() {
		logger.Log(logger.CreateDetails(logInfo, errors.New("wrong signature")))
		http.Error(w, "Wrong signature", http.StatusForbidden)
		return
	}

	fsRootHash, fsTree, err := shared.CalcRootHash(updatedFs.NewFs)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, httpErrorMsg, 500)
		return
	}

	rootSignature, err := hex.DecodeString(updatedFs.SignedFsRootHash)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, httpErrorMsg, 500)
		return
	}

	fsRootBytes, err := hex.DecodeString(fsRootHash)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, httpErrorMsg, 500)
		return
	}

	fsRootNonceBytes := append(fsRootBytes, nonce32...)

	hash := sha256.Sum256(fsRootNonceBytes)

	sigPublicKey, err = crypto.SigToPub(hash[:], rootSignature)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, httpErrorMsg, 500)
		return
	}

	signatureAddress = crypto.PubkeyToAddress(*sigPublicKey)

	if spAddress != signatureAddress.String() {
		logger.Log(logger.CreateDetails(logInfo, errors.New("wrong signature")))
		http.Error(w, "Wrong signature", http.StatusForbidden)
		return
	}

	shared.MU.Lock()

	spFsFile, err = os.Create(filepath.Join(addressPath, "tree.json"))
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, httpErrorMsg, 500)
		return
	}
	defer spFsFile.Close()

	spFs = shared.StorageProviderFs{
		Nonce:        updatedFs.Nonce,
		SignedFsRoot: updatedFs.SignedFsRootHash,
		Tree:         fsTree,
	}

	js, err := json.Marshal(spFs)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, httpErrorMsg, 500)
		return
	}

	_, err = spFsFile.Write(js)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, httpErrorMsg, 500)
		return
	}

	spFsFile.Sync()

	logger.Log("Updated fs info")

	shared.MU.Unlock()

}

func getNodeIP(nodeInfo nodeAbi.SimpleMetaDataDeNetNode) string {
	ipBuilder := strings.Builder{}
	for i, v := range nodeInfo.IpAddress {
		stringPart := strconv.Itoa(int(v))
		ipBuilder.WriteString(stringPart)

		if i < 3 {
			ipBuilder.WriteString(".")
		}
	}

	stringPort := strconv.Itoa(int(nodeInfo.Port))
	ipBuilder.WriteString(":")
	ipBuilder.WriteString(stringPort)

	return ipBuilder.String()
}

// ====================================================================================

func CopyFile(w http.ResponseWriter, r *http.Request) {
	logInfo := "server.CopyFile->"
	fmt.Println("Copy file")

	vars := mux.Vars(r)
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

	pathToConfig := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.ConfDirName, "config.json")

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

	var nodeConfig config.SecondaryNodeConfig

	err = json.Unmarshal(fileBytes, &nodeConfig)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Account config problem", 500)
		return
	}

	sharedSpaceInBytes := int64(nodeConfig.StorageLimit) * gbBytes

	nodeConfig.UsedStorageSpace += int64(intFileSize)

	noMemory := false

	if nodeConfig.UsedStorageSpace > sharedSpaceInBytes {
		noMemory = true
		nodeConfig.UsedStorageSpace -= int64(intFileSize)
	}

	avaliableSpaceLeft := sharedSpaceInBytes - nodeConfig.UsedStorageSpace

	if avaliableSpaceLeft < oneHunderdMBBytes {
		fmt.Println("Shared storage memory is running low,", avaliableSpaceLeft/(1024*1024), "MB of space is avaliable")
		fmt.Println("You may need additional space for mining. Total shared space can be changed in account configuration")
	}

	err = config.Save(confFile, nodeConfig)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Couldn't update config file", 500)
		return
	}
	confFile.Close()
	shared.MU.Unlock()

	err = r.ParseMultipartForm(1 << 20) // maxMemory 32MB
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "Parse multiform problem", 400)
		return
	}

	fs := r.MultipartForm.Value["fs"]

	sort.Strings(fs)

	fsRootHash, fsTree, err := shared.CalcRootHash(fs)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 400)
		return
	}

	signedFsRootHash := r.MultipartForm.Value["fsRootHash"]

	signature, err := hex.DecodeString(signedFsRootHash[0])
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 400)
		return
	}

	nonce := r.MultipartForm.Value["nonce"]

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

	storageProviderAddress := r.MultipartForm.Value["address"]

	senderAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if storageProviderAddress[0] != fmt.Sprint(senderAddress) {
		err := errors.New("wrong signature")
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	addressPath := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, storageProviderAddress[0])

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

	if noMemory {
		nftNode, err := blockchainprovider.GetNodeNFT()
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			http.Error(w, err.Error(), 400)
			return
		}

		total, err := nftNode.TotalSupply(&bind.CallOpts{})
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			http.Error(w, err.Error(), 400)
			return
		}

		intTotal := total.Int64()

		fastReq := fasthttp.AcquireRequest()
		fastResp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(fastReq)
		defer fasthttp.ReleaseResponse(fastResp)

		for i := int64(0); i < intTotal; i++ {
			node, err := nftNode.GetNodeById(&bind.CallOpts{}, big.NewInt(i))
			if err != nil {
				fmt.Println(err)
				continue
			}

			nodeIP := getNodeIP(node)

			if nodeIP == nodeConfig.IpAddress+":"+nodeConfig.HTTPPort {
				continue
			}

			url := "http://" + nodeIP
			fastReq.Reset()
			fastResp.Reset()

			fastReq.Header.SetRequestURI(url)
			fastReq.Header.SetMethod("GET")
			fastReq.Header.Set("Connection", "close")

			err = fasthttp.Do(fastReq, fastResp)
			if err != nil {
				continue
			}

			nodeAddress, err := backUpTo(nodeIP, addressPath, r.MultipartForm, intFileSize)
			if err != nil {
				continue
			}

			resp := NodeAddressResponse{
				NodeAddress: nodeAddress,
			}

			js, err := json.Marshal(resp)
			if err != nil {
				continue
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}

		http.Error(w, "no available nodes", 500)
		return
	}

	shared.MU.Lock()
	spFsFile, err := os.Create(filepath.Join(addressPath, "tree.json"))
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}
	defer spFsFile.Close()

	spFs := shared.StorageProviderFs{
		Nonce:        nonce[0],
		SignedFsRoot: signedFsRootHash[0],
		Tree:         fsTree,
	}

	js, err := json.Marshal(spFs)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}

	_, err = spFsFile.Write(js)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}

	spFsFile.Sync()
	spFsFile.Close()
	shared.MU.Unlock()

	hashes := r.MultipartForm.File["hashes"]
	hashesFile, err := hashes[0].Open()
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}

	hashesBody, err := io.ReadAll(hashesFile)
	if err != nil {
		hashesFile.Close()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}

	hashDif := make(map[string]string)
	err = json.Unmarshal(hashesBody, &hashDif)
	if err != nil {
		hashesFile.Close()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}

	hashesFile.Close()

	for old, new := range hashDif {
		path := filepath.Join(addressPath, old)
		file, err := os.Open(path)
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			restoreMemoryInfo(pathToConfig, intFileSize)
			http.Error(w, "File saving problem", 500)
			return
		}

		defer file.Close()

		newPath := filepath.Join(addressPath, new)
		newFile, err := os.Create(newPath)
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			restoreMemoryInfo(pathToConfig, intFileSize)
			http.Error(w, "File saving problem", 500)
			return
		}

		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			restoreMemoryInfo(pathToConfig, intFileSize)
			http.Error(w, "File saving problem", 500)
			return
		}

		newFile.Sync()
		newFile.Close()
	}

	resp := NodeAddressResponse{
		NodeAddress: nodeConfig.IpAddress + ":" + nodeConfig.HTTPPort,
	}

	js, err = json.Marshal(resp)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// ====================================================================================

func BackUp(w http.ResponseWriter, r *http.Request) {
	logInfo := "server.BackUp->"
	vars := mux.Vars(r)
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

	pathToConfig := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.ConfDirName, "config.json")

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

	var nodeConfig config.SecondaryNodeConfig

	err = json.Unmarshal(fileBytes, &nodeConfig)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Account config problem", 500)
		return
	}

	sharedSpaceInBytes := int64(nodeConfig.StorageLimit) * gbBytes

	nodeConfig.UsedStorageSpace += int64(intFileSize)

	if nodeConfig.UsedStorageSpace > sharedSpaceInBytes {
		shared.MU.Unlock()
		err := errors.New("insufficient memory avaliable")
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, err.Error(), 400)
		return
	}

	avaliableSpaceLeft := sharedSpaceInBytes - nodeConfig.UsedStorageSpace

	if avaliableSpaceLeft < oneHunderdMBBytes {
		fmt.Println("Shared storage memory is running low,", avaliableSpaceLeft/(1024*1024), "MB of space is avaliable")
		fmt.Println("You may need additional space for mining. Total shared space can be changed in account configuration")
	}

	err = config.Save(confFile, nodeConfig)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		http.Error(w, "Couldn't update config file", 500)
		return
	}
	confFile.Close()
	shared.MU.Unlock()

	err = r.ParseMultipartForm(1 << 20) // maxMemory 32MB
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "Parse multiform problem", 400)
		return
	}

	fs := r.MultipartForm.Value["fs"]

	sort.Strings(fs)

	fsRootHash, fsTree, err := shared.CalcRootHash(fs)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 400)
		return
	}

	signedFsRootHash := r.MultipartForm.Value["fsRootHash"]

	signature, err := hex.DecodeString(signedFsRootHash[0])
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 400)
		return
	}

	nonce := r.MultipartForm.Value["nonce"]

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

	storageProviderAddress := r.MultipartForm.Value["address"]

	senderAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if storageProviderAddress[0] != fmt.Sprint(senderAddress) {
		err := errors.New("wrong signature")
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	addressPath := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, storageProviderAddress[0])

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
	spFsFile, err := os.Create(filepath.Join(addressPath, "tree.json"))
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}
	defer spFsFile.Close()

	spFs := shared.StorageProviderFs{
		Nonce:        nonce[0],
		SignedFsRoot: signedFsRootHash[0],
		Tree:         fsTree,
	}

	js, err := json.Marshal(spFs)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}

	_, err = spFsFile.Write(js)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}

	spFsFile.Sync()
	spFsFile.Close()
	shared.MU.Unlock()

	reqFileParts := r.MultipartForm.File["files"]

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

	hashes := r.MultipartForm.File["hashes"]
	hashesFile, err := hashes[0].Open()
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}

	hashesBody, err := io.ReadAll(hashesFile)
	if err != nil {
		hashesFile.Close()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}

	hashDif := make(map[string]string)
	err = json.Unmarshal(hashesBody, &hashDif)
	if err != nil {
		hashesFile.Close()
		logger.Log(logger.CreateDetails(logInfo, err))
		restoreMemoryInfo(pathToConfig, intFileSize)
		http.Error(w, "File saving problem", 500)
		return
	}

	hashesFile.Close()

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

		pathToFile := filepath.Join(addressPath, hashDif[reqFilePart.Filename])

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

		logger.Log("Saved file " + hashDif[reqFilePart.Filename] + " (" + fmt.Sprint(count) + "/" + fmt.Sprint(total) + ")" + " from " + storageProviderAddress[0]) //TODO remove

		newFile.Sync()
		rqFile.Close()
		newFile.Close()

		count++
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

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

	var nodeConfig config.SecondaryNodeConfig

	err = json.Unmarshal(fileBytes, &nodeConfig)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		return
	}

	nodeConfig.UsedStorageSpace -= int64(intFileSize)

	err = config.Save(confFile, nodeConfig)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logInfo, err))
		return
	}
	shared.MU.Unlock()
}

func backUpTo(nodeAddress, addressPath string, multiForm *multipart.Form, fileSize int) (string, error) {
	const logInfo = "server.logInfo->"

	pipeConns := fasthttputil.NewPipeConns()
	pr := pipeConns.Conn1()
	pw := pipeConns.Conn2()

	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()

		address := multiForm.Value["address"]
		err := writer.WriteField("address", address[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		wholeFileHashes := multiForm.Value["fs"]
		for _, wholeHash := range wholeFileHashes {
			err = writer.WriteField("fs", wholeHash)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		nonce := multiForm.Value["nonce"]
		err = writer.WriteField("nonce", nonce[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		fsRootHash := multiForm.Value["fsRootHash"]
		err = writer.WriteField("fsRootHash", fsRootHash[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		hashes := multiForm.File["hashes"]
		hashesFile, err := hashes[0].Open()
		if err != nil {
			fmt.Println(err)
			return
		}

		defer hashesFile.Close()

		hashesBody, err := io.ReadAll(hashesFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		h, err := writer.CreateFormFile("hashes", "hashes")
		if err != nil {
			fmt.Println(err)
			return
		}

		h.Write(hashesBody)

		hashDif := make(map[string]string)
		err = json.Unmarshal(hashesBody, &hashDif)
		if err != nil {
			fmt.Println(err)
			return
		}

		for old := range hashDif {
			path := filepath.Join(addressPath, old)
			file, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
				return
			}

			defer file.Close()

			filePart, err := writer.CreateFormFile("files", old)
			if err != nil {
				fmt.Println(err)
				return
			}

			_, err = io.Copy(filePart, file)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		writer.Close()
	}()

	req, err := http.NewRequest("POST", "http://"+nodeAddress+"/backup/"+strconv.Itoa(fileSize), pr)
	if err != nil {
		return "", logger.CreateDetails(logInfo, err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", logger.CreateDetails(logInfo, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", logger.CreateDetails(logInfo, err)
	}

	defer resp.Body.Close()

	fmt.Println(string(body))
	if string(body) != "OK" {

		return "", logger.CreateDetails(logInfo, errors.New("saving problem"))
	}

	return nodeAddress, nil
}
