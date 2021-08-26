package files

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"git.denetwork.xyz/dfile/dfile-secondary-node/config"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	nodeAbi "git.denetwork.xyz/dfile/dfile-secondary-node/node_abi"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/valyala/fasthttp/fasthttputil"
)

type NodesResponse struct {
	Nodes []string `json:"nodes"`
}

type UpdatedFsInfo struct {
	NewFs            []string
	Nonce            string
	SignedFsRootHash string
}

func UpdateFileSystemInfo(updatedFs *UpdatedFsInfo, spAddress, signedFileSystem string) error {
	const logLoc = "files.UpdateFileSystemInfo->"

	addressPath := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spAddress)

	_, err := os.Stat(addressPath)
	if err != nil {
		return logger.CreateDetails(logLoc, errors.New("no files of "+spAddress))
	}

	shared.MU.Lock()
	spFsFile, fileBytes, err := shared.ReadFile(filepath.Join(addressPath, paths.SpFsFilename))
	if err != nil {
		shared.MU.Unlock()
		return logger.CreateDetails(logLoc, err)
	}

	defer spFsFile.Close()

	var spFs shared.StorageProviderData

	err = json.Unmarshal(fileBytes, &spFs)
	if err != nil {
		shared.MU.Unlock()
		return logger.CreateDetails(logLoc, err)
	}

	spFsFile.Close()
	shared.MU.Unlock()

	newNonceInt, err := strconv.Atoi(updatedFs.Nonce)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	currentNonceInt, err := strconv.Atoi(spFs.Nonce)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	if newNonceInt < currentNonceInt {
		return logger.CreateDetails(logLoc, fmt.Errorf("%v fs info is up to date", spAddress))
	}

	nonceHex := strconv.FormatInt(int64(newNonceInt), 16)

	nonceBytes, err := hex.DecodeString(nonceHex)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	nonce32 := make([]byte, 32-len(nonceBytes))
	nonce32 = append(nonce32, nonceBytes...)

	sort.Strings(updatedFs.NewFs)

	concatFsHashesBuilder := strings.Builder{}

	for _, hash := range updatedFs.NewFs {
		concatFsHashesBuilder.WriteString(hash)
	}

	fsTreeNonceBytes := append([]byte(concatFsHashesBuilder.String()), nonce32...)
	fsTreeNonceSha := sha256.Sum256(fsTreeNonceBytes)

	fsysSignature, err := hex.DecodeString(signedFileSystem)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	sigPublicKey, err := crypto.SigToPub(fsTreeNonceSha[:], fsysSignature)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	signatureAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if spAddress != signatureAddress.String() {
		return logger.CreateDetails(logLoc, err)
	}

	fsRootHash, fsTree, err := shared.CalcRootHash(updatedFs.NewFs)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	rootSignature, err := hex.DecodeString(updatedFs.SignedFsRootHash)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	fsRootBytes, err := hex.DecodeString(fsRootHash)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	fsRootNonceBytes := append(fsRootBytes, nonce32...)

	hash := sha256.Sum256(fsRootNonceBytes)

	sigPublicKey, err = crypto.SigToPub(hash[:], rootSignature)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	signatureAddress = crypto.PubkeyToAddress(*sigPublicKey)

	if spAddress != signatureAddress.String() {
		return logger.CreateDetails(logLoc, err)
	}

	shared.MU.Lock()

	spFsFile, err = os.Create(filepath.Join(addressPath, paths.SpFsFilename))
	if err != nil {
		shared.MU.Unlock()
		return logger.CreateDetails(logLoc, err)
	}

	defer spFsFile.Close()

	spFs = shared.StorageProviderData{
		Nonce:        updatedFs.Nonce,
		SignedFsRoot: updatedFs.SignedFsRootHash,
		Tree:         fsTree,
	}

	err = shared.WriteFile(spFsFile, spFs)
	if err != nil {
		shared.MU.Unlock()
		return logger.CreateDetails(logLoc, err)
	}

	spFsFile.Sync()
	shared.MU.Unlock()

	return nil
}

// ====================================================================================

//Provides back up to "node Address" using old multipart form. Returning ip address node if successful
func backUpCopy(nodeAddress, addressPath string, multiForm *multipart.Form, fileSize int) (string, error) {
	const logLoc = "files_helpers.backUpCopy->"

	pipeConns := fasthttputil.NewPipeConns()
	pr := pipeConns.Conn1()
	pw := pipeConns.Conn2()

	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()

		address := multiForm.Value["address"]
		nonce := multiForm.Value["nonce"]
		fsRootHash := multiForm.Value["fsRootHash"]

		err := prepareMultipartForm(writer, address[0], nonce[0], fsRootHash[0])
		if err != nil {
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

	req, err := http.NewRequest("POST", "http://"+nodeAddress+"/backup/copy/"+strconv.Itoa(fileSize), pr)
	if err != nil {
		return "", logger.CreateDetails(logLoc, err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", logger.CreateDetails(logLoc, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", logger.CreateDetails(logLoc, err)
	}

	defer resp.Body.Close()

	if string(body) != "OK" {
		return "", logger.CreateDetails(logLoc, shared.ErrFileSaving)
	}

	return nodeAddress, nil
}

// ====================================================================================

//Restore certain file size
func RestoreMemoryInfo(pathToConfig string, intFileSize int) {
	logLoc := "files.restoreMemoryInfo->"

	shared.MU.Lock()
	confFile, fileBytes, err := shared.ReadFile(pathToConfig)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logLoc, err))
		return
	}
	defer confFile.Close()

	var nodeConfig config.SecondaryNodeConfig

	err = json.Unmarshal(fileBytes, &nodeConfig)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logLoc, err))
		return
	}

	nodeConfig.UsedStorageSpace -= int64(intFileSize)

	err = config.Save(confFile, nodeConfig)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(logLoc, err))
		return
	}
	shared.MU.Unlock()
}

// ====================================================================================

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

func VerifyStorageProviderAddress(spAddress, fileSize, signatureFromReq string, fileKeys []string) error {
	const logLoc = "files.VerifyStorageProviderAddress->"
	var wholeFileHash string
	var err error

	if len(fileKeys) == 1 {
		wholeFileHash = fileKeys[0]
	} else {
		sort.Strings(fileKeys)
		wholeFileHash, _, err = shared.CalcRootHash(fileKeys)
		if err != nil {
			return logger.CreateDetails(logLoc, err)
		}
	}

	signature, err := hex.DecodeString(signatureFromReq)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	hash := sha256.Sum256([]byte(spAddress + fileSize + wholeFileHash))

	sigPublicKey, err := crypto.SigToPub(hash[:], signature)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	signatureAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if spAddress != signatureAddress.String() {
		return logger.CreateDetails(logLoc, shared.ErrWrongSignature)
	}

	return nil
}

// ====================================================================================

func DeleteParts(addressPath string, fileHashes []string) {
	logger.Log("deleting file parts after error...")

	for _, hash := range fileHashes {
		pathToFile := filepath.Join(addressPath, hash)

		os.Remove(pathToFile)
	}
}

// ====================================================================================

func initSPFile(addressPath string, spData *shared.StorageProviderData) error {
	const logLoc = "files.initSPFile->"

	stat, err := os.Stat(addressPath)
	err = shared.CheckStatErr(err)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	if stat == nil {
		err = os.Mkdir(addressPath, 0700)
		if err != nil {
			return logger.CreateDetails(logLoc, err)
		}
	}

	shared.MU.Lock()
	defer shared.MU.Unlock()

	spFsFile, err := os.Create(filepath.Join(addressPath, paths.SpFsFilename))
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	defer spFsFile.Close()

	err = shared.WriteFile(spFsFile, spData)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	spFsFile.Sync()

	return nil
}

// ====================================================================================

func prepareMultipartForm(writer *multipart.Writer, spAddress, nonce, fsRootHash string) error {
	const logLoc = "files.initSPFile->"
	err := writer.WriteField("address", spAddress)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	err = writer.WriteField("nonce", nonce)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	err = writer.WriteField("fsRootHash", fsRootHash)
	if err != nil {
		return logger.CreateDetails(logLoc, err)
	}

	return nil
}
