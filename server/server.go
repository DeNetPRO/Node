package server

import (
	"context"
	"errors"
	"os/signal"
	"strings"

	"github.com/minio/sha256-simd"

	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	blckChain "git.denetwork.xyz/DeNet/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/config"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/errs"
	fsysInfo "git.denetwork.xyz/DeNet/dfile-secondary-node/fsys_info"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/hash"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	memInfo "git.denetwork.xyz/DeNet/dfile-secondary-node/mem_info"

	nodeFile "git.denetwork.xyz/DeNet/dfile-secondary-node/node_file"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/upnp"

	_ "git.denetwork.xyz/DeNet/dfile-secondary-node/docs"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/shared"
	spFiles "git.denetwork.xyz/DeNet/dfile-secondary-node/sp_files"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type ReqData struct {
	RequesterAddr string
	FileName      string
	FsTreeHash    string
}

type NodeAddressResponse struct {
	NodeAddress string `json:"node_address"`
}

const gbBytes = int64(1024 * 1024 * 1024)
const oneHunderdMBBytes = int64(1024 * 1024 * 100)
const serverStartFatalMessage = "Couldn't start server"

func Start(port string) {
	const location = "server.Start->"
	r := mux.NewRouter()

	r.HandleFunc("/ping", healthCheck).Methods("GET")

	r.HandleFunc("/upload/{verificationData}/{size}/{network}", SaveFiles).Methods("POST")

	r.HandleFunc("/download/{verificationData}/{access}/{network}", ServeFiles).Methods("GET")

	r.HandleFunc("/update_fs/{verificationData}/{network}", UpdateFsInfo).Methods("POST")

	r.HandleFunc("/backup_fs/{verificationData}", backUpSpSf).Methods("GET", "POST")

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
		Handler: corsOpts.Handler(verifyRequest(r)),
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

var pathExceptions = map[string]bool{
	"ping": true,
}

func verifyRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		address := r.URL.Path
		splitPath := strings.Split(address, "/")

		exceptedPath := pathExceptions[splitPath[1]]

		if !exceptedPath {
			verificationData := strings.Split(splitPath[2], "$")

			requesterAddr := verificationData[0]
			signedData := verificationData[1]
			unsignedData := verificationData[2]

			err := verifySignature(requesterAddr, signedData, unsignedData)
			if err != nil {
				http.Error(w, errors.New("forbidden").Error(), http.StatusForbidden)
				return
			}

			var requestData = ReqData{
				RequesterAddr: requesterAddr,
			}

			ctx := context.WithValue(r.Context(), "requestData", requestData)

			if splitPath[1] == "download" {

				requestData.FileName = unsignedData

				ctx = context.WithValue(r.Context(), "requestData", requestData)

				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if splitPath[1] == "update_fs" {

				requestData.FsTreeHash = unsignedData

				ctx = context.WithValue(r.Context(), "requestData", requestData)

				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

		}

		next.ServeHTTP(w, r)

	})
}

// ====================================================================================

func verifySignature(requesterAddr, signedData, unsignedData string) error {
	signature, err := hex.DecodeString(signedData)
	if err != nil {
		return err
	}

	hash := sha256.Sum256([]byte(unsignedData))

	sigPublicKey, err := crypto.SigToPub(hash[:], signature)
	if err != nil {
		return err
	}

	signatureAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if signatureAddress.String() != requesterAddr {
		return errors.New("wrong signature")
	}

	return nil
}

// ====================================================================================

// Healthcheck godoc
// @Summary Check node status
// @Description Checking node performance
// @Success 200 {string} string "ok"
// @Header 200 {string} Status "OK"
// @Router /ping [get]
func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

// ========================================================================================================

// SaveFiles godoc
// @Summary Save files
// @Description Save files from Storage Provider
// @Accept  multipart/form-data
// @Param size path int true "file size in bytes"
// @Param network path string true "network type"
// @Param address formData string true "Storage Provider address"
// @Param fsRootHash formData string  true "signed file system root hash"
// @Param nonce formData int true "current nonce"
// @Param fs formData []string true "array of hashes of all storage provider files"
// @Param files formData file  true "files parts"
// @Success 200 {string} Status "OK"
// @Router /upload/{size}/{network} [post]
func SaveFiles(w http.ResponseWriter, req *http.Request) {
	const location = "server.SaveFiles->"

	vars := mux.Vars(req)
	network := vars["network"]

	if network == "" {
		network = "kovan"
	}

	_, netExists := blckChain.Networks[network]

	if !netExists {
		http.Error(w, errs.NetworkCheck.Error(), http.StatusBadRequest)
		return
	}

	pathToConfig := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.ConfDirName, paths.ConfFileName)

	fileSize, spaceNotEnough, _, err := checkAndReserveSpace(req, pathToConfig)
	if err != nil {

		logger.Log(logger.CreateDetails(location, err))

		if spaceNotEnough {
			http.Error(w, errs.NoSpace.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, errs.SpaceCheck.Error(), http.StatusBadRequest)
		return
	}

	spData, err := parseRequest(req)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		memInfo.Restore(pathToConfig, fileSize)
		http.Error(w, errs.ParseMultipartForm.Error(), http.StatusInternalServerError)
		return
	}

	pathToSpFiles := filepath.Join(paths.StoragePaths[0], network, spData.Address)

	dirStat, err := os.Stat(pathToSpFiles)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Log(logger.CreateDetails(location, err))
		memInfo.Restore(pathToConfig, fileSize)
		http.Error(w, errs.Internal.Error(), http.StatusInternalServerError)
		return
	}

	var dirFilesCount = 0

	if dirStat != nil {
		dirFiles, err := nodeFile.ReadDirFiles(pathToSpFiles)
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
		}

		dirFilesCount = len(dirFiles)
	}

	if dirFilesCount > 3300 { // max Mb storage per user
		http.Error(w, errs.NoSpace.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("dir contains", dirFilesCount, "files")

	err = spFiles.Save(req, spData, pathToSpFiles)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		memInfo.Restore(pathToConfig, fileSize)
		http.Error(w, errs.Internal.Error(), http.StatusInternalServerError)
		return
	}

	if !shared.TestMode {
		logger.SendStatistic(spData.Address, network, req.RemoteAddr, logger.Upload, int64(fileSize))
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

// ====================================================================================

// ServeFiles godoc
// @Summary Serve file
// @Description Serve file by key
// @Produce octet-stream
// @Param spAddress path string true "Storage Provider address"
// @Param fileKey path string true "file key"
// @Param signature path string true "Storage Provider signature"
// @Param newtork path string  true "network type"
// @Success 200 {file} binary
// @Router /download/{spAddress}/{fileKey}/{signature}/{network} [get]
func ServeFiles(w http.ResponseWriter, r *http.Request) {
	const location = "server.ServeFiles->"

	vars := mux.Vars(r)
	network := vars["network"]

	_, netExists := blckChain.Networks[network]

	if !netExists {
		http.Error(w, errs.NetworkCheck.Error(), http.StatusBadRequest)
		return
	}

	access := vars["access"]

	accessParams := strings.Split(access, "$")

	if len(accessParams) != 3 {
		http.Error(w, errs.InvalidArgument.Error(), http.StatusBadRequest)
		return
	}

	ownerAddr := accessParams[0]
	signedGrant := accessParams[1]
	permittedTo := accessParams[2]

	err := verifySignature(ownerAddr, signedGrant, permittedTo)
	if err != nil {
		http.Error(w, errors.New("forbidden").Error(), http.StatusForbidden)
		return
	}

	rqtData := r.Context().Value("requestData").(ReqData)

	if rqtData.RequesterAddr != permittedTo {
		http.Error(w, errors.New("forbidden").Error(), http.StatusForbidden)
		return
	}

	pathToFile := filepath.Join(paths.StoragePaths[0], network, ownerAddr, rqtData.FileName)

	stat, err := os.Stat(pathToFile)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		http.Error(w, errs.Internal.Error(), http.StatusInternalServerError)
		return
	}

	logger.Log("serving file: " + rqtData.FileName)

	if !shared.TestMode {
		logger.SendStatistic(rqtData.RequesterAddr, network, r.RemoteAddr, logger.Download, stat.Size())
	}

	http.ServeFile(w, r, pathToFile)
}

// ====================================================================================

// UpdateFsInfo godoc
// @Summary Update Storage Provider's filesystem
// @Description Update Storage Provider's filesystem, etc. root hash, nonce, file system
// @Accept  json
// @Param spAddress path string true "Storage Provider address"
// @Param signedFsys path string true "Signed Storage Provider root hash"
// @Param newtork path string  true "network type"
// @Param updatedFsInfo body fsysInfo.UpdatedFsInfo true "updatedFsInfo"
// @Success 200 {string} Status "OK"
// @Router /update_fs/{spAddress}/{signedFsys}/{network} [post]
func UpdateFsInfo(w http.ResponseWriter, r *http.Request) {
	const location = "server.UpdateFsInfo->"

	vars := mux.Vars(r)
	network := vars["network"]

	_, netExists := blckChain.Networks[network]

	if !netExists {
		http.Error(w, errs.NetworkCheck.Error(), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		http.Error(w, errs.InvalidArgument.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updatedFs := &fsysInfo.UpdatedFsInfo{}
	err = json.Unmarshal(body, &updatedFs)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		http.Error(w, errs.UpdateFsInfo.Error(), http.StatusInternalServerError)
		return
	}

	rqtData := r.Context().Value("requestData").(ReqData)

	err = fsysInfo.Update(updatedFs, rqtData.RequesterAddr, rqtData.FsTreeHash, network)
	if err != nil {
		logger.Log(logger.CreateDetails(location, err))
		http.Error(w, errs.Internal.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Updated!")
}

// ====================================================================================

//Checks space on the node.
//Returns the size of the input file, true -> if there is enough space and false -> if otherwise.
//And also node's config.
func checkAndReserveSpace(r *http.Request, pathToConfig string) (int, bool, config.NodeConfig, error) {
	const location = "server.checkSpace"

	var nodeConfig config.NodeConfig

	vars := mux.Vars(r)
	fileSize := vars["size"]

	intFileSize, err := strconv.Atoi(fileSize)
	if err != nil {
		return 0, false, nodeConfig, logger.CreateDetails(location, err)
	}

	if intFileSize == 0 {
		return 0, false, nodeConfig, logger.CreateDetails(location, errors.New("file size is 0"))
	}

	shared.MU.Lock()
	defer shared.MU.Unlock()

	confFile, fileBytes, err := nodeFile.Read(pathToConfig)
	if err != nil {
		shared.MU.Unlock()
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
		return 0, true, nodeConfig, logger.CreateDetails(location, errs.NoSpace)
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

	return intFileSize, false, nodeConfig, nil
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

// StorageSystem godoc
// @Summary Returns Storage Provider filesystem on "GET" request and refreshes filesystem on "POST"
// @Accept multipart/form-data
// @Param spAddress path string true "Storage Provider address"
// @Param signature path string true "Signed Storage Provider address"
// @Router /storage/system/{spAddress}/{signature} [post]
// @Param fs formData file  true "encoded Storage Provider filesystem"
// @Success 200 {string} Status "OK"
// @Router /storage/system/{spAddress}/{signature} [get]
// @Success 200 {file} binary
func backUpSpSf(w http.ResponseWriter, r *http.Request) {
	const location = "server.backUpSpSf"

	rqtData := r.Context().Value("requestData").(ReqData)

	stat, _ := os.Stat(paths.SystemsDirPath)
	if stat == nil {
		os.Mkdir(paths.SystemsDirPath, 0777)
	}

	switch r.Method {
	case http.MethodGet:
		path, exists := spFiles.SearchStorageFilesystem(rqtData.RequesterAddr)
		if !exists {
			logger.Log(logger.CreateDetails(location, errs.StorageSystemNotFound))
			http.Error(w, errs.StorageSystemNotFound.Error(), http.StatusBadRequest)
			return
		}

		http.ServeFile(w, r, path)
	case http.MethodPost:
		err := r.ParseMultipartForm(1 << 20) // maxMemory 32MB
		if err != nil {
			logger.Log(logger.CreateDetails(location, errs.ParseMultipartForm))
			http.Error(w, errs.ParseMultipartForm.Error(), http.StatusBadRequest)
			return
		}

		fileSystemHeader := r.MultipartForm.File["fs"]
		if len(fileSystemHeader) == 0 {
			logger.Log(logger.CreateDetails(location, errs.InvalidArgument))
			http.Error(w, errs.InvalidArgument.Error(), http.StatusBadRequest)
			return
		}

		err = spFiles.UpdateStorageFilesystem(rqtData.RequesterAddr, fileSystemHeader[0])
		if err != nil {
			logger.Log(logger.CreateDetails(location, errs.Internal))
			http.Error(w, errs.Internal.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Storage", rqtData.RequesterAddr, "FS backed up!")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK")
		return
	default:
		err := errors.New("invalid method")
		logger.Log(logger.CreateDetails(location, err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
}
