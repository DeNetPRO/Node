package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

func deleteFileParts(addressPath string, fileHashes []string) {
	logger.Log("deleting file parts after error...")

	for _, hash := range fileHashes {
		pathToFile := filepath.Join(addressPath, hash)

		os.Remove(pathToFile)
	}
}

//Provides back up to "node Address" using old multipart form. Returning ip address node if successful
func backUpCopy(nodeAddress, addressPath string, multiForm *multipart.Form, fileSize int) (string, error) {
	const logLoc = "server.backUpCopy->"

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

//Restore certain file size
func RestoreMemoryInfo(pathToConfig string, intFileSize int) {
	logLoc := "server.restoreMemoryInfo->"

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

func BackUpFileKeys(nodeAddress, addressPath string, spData *shared.StorageProviderData, fileSize int, fileKeys []string) ([]string, error) {
	const logLoc = "server.BackUpFileKeys->"

	pipeConns := fasthttputil.NewPipeConns()
	pr := pipeConns.Conn1()
	pw := pipeConns.Conn2()

	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()

		err := prepareMultipartForm(writer, spData.Address, spData.Nonce, spData.SignedFsRoot)
		if err != nil {
			return
		}

		for _, key := range fileKeys {
			path := filepath.Join(addressPath, key)
			file, err := os.Open(path)
			if err != nil {
				return
			}

			defer file.Close()

			filePart, err := writer.CreateFormFile("files", key)
			if err != nil {
				return
			}

			_, err = io.Copy(filePart, file)
			if err != nil {
				return
			}
		}

		writer.Close()
	}()

	req, err := http.NewRequest("POST", "http://"+nodeAddress+"/backup/new/"+strconv.Itoa(fileSize), pr)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, logger.CreateDetails(logLoc, shared.ErrFileSaving)
	}

	nodesResp := NodesResponse{}
	err = json.Unmarshal(body, &nodesResp)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	return nodesResp.Nodes, nil
}

func BackUpFileParts(nodeAddress string, spData *shared.StorageProviderData, fileSize int, fileParts []*multipart.FileHeader) ([]string, error) {
	const logLoc = "server.BackUpFileParts->"

	pipeConns := fasthttputil.NewPipeConns()
	pr := pipeConns.Conn1()
	pw := pipeConns.Conn2()

	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		err := prepareMultipartForm(writer, spData.Address, spData.Nonce, spData.SignedFsRoot)
		if err != nil {
			return
		}

		for _, part := range fileParts {
			name := part.Filename
			filePart, err := writer.CreateFormFile("files", name)
			if err != nil {
				return
			}

			temp, err := part.Open()
			if err != nil {
				return
			}

			io.Copy(filePart, temp)
		}

		writer.Close()
	}()

	req, err := http.NewRequest("POST", "http://"+nodeAddress+"/backup/new/"+strconv.Itoa(fileSize), pr)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, logger.CreateDetails(logLoc, shared.ErrFileSaving)
	}

	nodesResp := NodesResponse{}
	err = json.Unmarshal(body, &nodesResp)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	return nodesResp.Nodes, nil
}

func VerifyStorageProviderAddress(spAddress, fileSize, signatureFromReq string, fileKeys []string) error {
	const logLoc = "server.VerifyStorageProviderAddress->"
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

func initSPFile(addressPath string, spData *shared.StorageProviderData) error {
	const logLoc = "server.initSPFile->"

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

func prepareMultipartForm(writer *multipart.Writer, spAddress, nonce, fsRootHash string) error {
	const logLoc = "server.initSPFile->"
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
