package server

import (
	"context"
	"os/signal"

	"github.com/minio/sha256-simd"

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
	"git.denetwork.xyz/dfile/dfile-secondary-node/errs"
	fsysInfo "git.denetwork.xyz/dfile/dfile-secondary-node/fsys_info"
	"git.denetwork.xyz/dfile/dfile-secondary-node/hash"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	memInfo "git.denetwork.xyz/dfile/dfile-secondary-node/mem_info"
	nodeFile "git.denetwork.xyz/dfile/dfile-secondary-node/node_file"

	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	spFiles "git.denetwork.xyz/dfile/dfile-secondary-node/sp_files"
	"git.denetwork.xyz/dfile/dfile-secondary-node/upnp"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type NodeAddressResponse struct {
	NodeAddress string `json:"node_address"`
}

const gbBytes = int64(1024 * 1024 * 1024)
const oneHunderdMBBytes = int64(1024 * 1024 * 100)
const serverStartFatalMessage = "Couldn't start server"

func Start(port string) {
	const location = "server.Start->"
	r := mux.NewRouter()

	r.HandleFunc("/upload/{size}", SaveFiles).Methods("POST")
	r.HandleFunc("/download/{spAddress}/{fileKey}/{signature}", ServeFiles).Methods("GET")
	r.HandleFunc("/update_fs/{spAddress}/{signedFsys}", UpdateFsInfo).Methods("POST")
	r.HandleFunc("/copy/{size}", CopyFile).Methods("POST")

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
		logger.Log(logger.CreateDetails(location, err))
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
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal(serverStartFatalMessage)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	err = server.Shutdown(context.Background())
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
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
		// 	http.Error(w, "session key verification error", http.StatusInternalServerError)
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
	const location = "server.SaveFiles->"

	pathToConfig := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.ConfDirName, paths.ConfFileName)

	intFileSize, _, _, err := checkSpace(req, pathToConfig)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		http.Error(w, errs.SpaceCheck.Error(), http.StatusInternalServerError)
		return
	}

	spData, err := parseRequest(req)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		memInfo.Restore(pathToConfig, intFileSize)
		http.Error(w, errs.ParseMultipartForm.Error(), http.StatusBadRequest)
		return
	}

	err = spFiles.Save(req, spData, pathToConfig, intFileSize)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		memInfo.Restore(pathToConfig, intFileSize)
		http.Error(w, errs.Internal.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

// ====================================================================================

func ServeFiles(w http.ResponseWriter, req *http.Request) {
	const location = "server.ServeFiles->"

	vars := mux.Vars(req)
	spAddress := vars["spAddress"]
	fileKey := vars["fileKey"]
	signatureFromReq := vars["signature"]

	if spAddress == "" || fileKey == "" || signatureFromReq == "" {
		logger.Log(logger.CreateDetails(location, errs.InvalidArgument))
		http.Error(w, errs.InvalidArgument.Error(), http.StatusBadRequest)
		return
	}

	pathToFile, err := spFiles.Serve(spAddress, fileKey, signatureFromReq)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		http.Error(w, errs.Internal.Error(), http.StatusInternalServerError)
		return
	}

	if !shared.TestMode {
		logger.Log("serving file: " + fileKey)
	}

	http.ServeFile(w, req, pathToFile)
}

// ====================================================================================

func UpdateFsInfo(w http.ResponseWriter, req *http.Request) {
	const location = "server.UpdateFsInfo->"

	vars := mux.Vars(req)
	spAddress := vars["spAddress"]
	signedFsys := vars["signedFsys"]

	if spAddress == "" || signedFsys == "" {
		logger.Log(logger.CreateDetails(location, errs.InvalidArgument))
		http.Error(w, errs.InvalidArgument.Error(), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		http.Error(w, errs.InvalidArgument.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	updatedFs := &fsysInfo.UpdatedFsInfo{}
	err = json.Unmarshal(body, &updatedFs)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		http.Error(w, errs.UpdateFsInfo.Error(), http.StatusInternalServerError)
		return
	}

	err = fsysInfo.Update(updatedFs, spAddress, signedFsys)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		http.Error(w, errs.Internal.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Updated!")
}

// ====================================================================================

func CopyFile(w http.ResponseWriter, req *http.Request) {
	location := "server.CopyFile->"

	pathToConfig := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.ConfDirName, paths.ConfFileName)

	intFileSize, enoughSpace, nodeConfig, err := checkSpace(req, pathToConfig)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		http.Error(w, errs.SpaceCheck.Error(), http.StatusInternalServerError)
		return
	}

	spData, err := parseRequest(req)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		memInfo.Restore(pathToConfig, intFileSize)
		http.Error(w, errs.ParseMultipartForm.Error(), http.StatusBadRequest)
		return
	}

	nodeAddress, err := spFiles.Copy(req, spData, &nodeConfig, pathToConfig, intFileSize, enoughSpace)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		memInfo.Restore(pathToConfig, intFileSize)
		http.Error(w, errs.Internal.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(nodeAddress)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// ====================================================================================

//Checks space on the node.
//Returns the size of the input file, true -> if there is enough space and false -> if otherwise.
//And also node's config.
func checkSpace(r *http.Request, pathToConfig string) (int, bool, config.SecondaryNodeConfig, error) {
	const location = "server.checkSpace"

	var nodeConfig config.SecondaryNodeConfig

	vars := mux.Vars(r)
	fileSize := vars["size"]

	intFileSize, err := strconv.Atoi(fileSize)
	if err != nil {
		return 0, false, nodeConfig, logger.CreateDetails(location, err)
	}

	if intFileSize == 0 {
		return 0, false, nodeConfig, logger.CreateDetails(location, err)
	}

	shared.MU.Lock()
	confFile, fileBytes, err := nodeFile.Read(pathToConfig)
	if err != nil {
		return 0, false, nodeConfig, logger.CreateDetails(location, err)
	}
	defer confFile.Close()

	err = json.Unmarshal(fileBytes, &nodeConfig)
	if err != nil {
		return 0, false, nodeConfig, logger.CreateDetails(location, err)
	}

	sharedSpaceInBytes := int64(nodeConfig.StorageLimit) * gbBytes

	nodeConfig.UsedStorageSpace += int64(intFileSize)

	if nodeConfig.UsedStorageSpace > sharedSpaceInBytes {
		return 0, false, nodeConfig, logger.CreateDetails(location, errors.New("not enough space"))
	}

	avaliableSpaceLeft := sharedSpaceInBytes - nodeConfig.UsedStorageSpace

	if avaliableSpaceLeft < oneHunderdMBBytes {
		fmt.Println("Shared storage memory is running low,", avaliableSpaceLeft/(1024*1024), "MB of space is avaliable")
		fmt.Println("You may need additional space for storing data. Total shared space can be changed in account configuration")
	}

	err = config.Save(confFile, nodeConfig)
	if err != nil {
		return 0, false, nodeConfig, logger.CreateDetails(location, err)
	}
	confFile.Close()
	shared.MU.Unlock()

	return intFileSize, true, nodeConfig, nil
}

// ====================================================================================

//Parse the request multipartForm
func parseRequest(r *http.Request) (*shared.StorageProviderData, error) {
	const location = "server.parseRequest"

	err := r.ParseMultipartForm(1 << 20) // maxMemory 32MB
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	fs := r.MultipartForm.Value["fs"]

	sort.Strings(fs)

	fsRootHash, fsTree, err := hash.CalcRoot(fs)
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	signedFsRootHash := r.MultipartForm.Value["fsRootHash"]

	signature, err := hex.DecodeString(signedFsRootHash[0])
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	nonce := r.MultipartForm.Value["nonce"]

	nonceInt, err := strconv.Atoi(nonce[0])
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	nonceHex := strconv.FormatInt(int64(nonceInt), 16)

	nonceBytes, err := hex.DecodeString(nonceHex)
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	nonce32 := make([]byte, 32-len(nonceBytes))
	nonce32 = append(nonce32, nonceBytes...)

	fsRootBytes, err := hex.DecodeString(fsRootHash)
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	fsRootNonceBytes := append(fsRootBytes, nonce32...)

	hash := sha256.Sum256(fsRootNonceBytes)

	sigPublicKey, err := crypto.SigToPub(hash[:], signature)
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	storageProviderAddress := r.MultipartForm.Value["address"]

	senderAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if storageProviderAddress[0] != fmt.Sprint(senderAddress) {
		fmt.Println(storageProviderAddress[0], fmt.Sprint(senderAddress))
		return nil, logger.CreateDetails(location, errs.WrongSignature)
	}

	return &shared.StorageProviderData{
		Address:      storageProviderAddress[0],
		Fs:           fs,
		Nonce:        nonce[0],
		SignedFsRoot: signedFsRootHash[0],
		Tree:         fsTree,
	}, nil
}
