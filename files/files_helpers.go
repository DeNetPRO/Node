package files

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/minio/sha256-simd"

	"git.denetwork.xyz/dfile/dfile-secondary-node/config"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	nodeAbi "git.denetwork.xyz/dfile/dfile-secondary-node/node_abi"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"github.com/ethereum/go-ethereum/crypto"
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
	const location = "files.UpdateFileSystemInfo->"

	addressPath := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spAddress)

	_, err := os.Stat(addressPath)
	if err != nil {
		return logger.CreateDetails(location, errors.New("no files of "+spAddress))
	}

	shared.MU.Lock()
	spFsFile, fileBytes, err := shared.ReadFile(filepath.Join(addressPath, paths.SpFsFilename))
	if err != nil {
		shared.MU.Unlock()
		return logger.CreateDetails(location, err)
	}

	defer spFsFile.Close()

	var spFs shared.StorageProviderData

	err = json.Unmarshal(fileBytes, &spFs)
	if err != nil {
		shared.MU.Unlock()
		return logger.CreateDetails(location, err)
	}

	spFsFile.Close()
	shared.MU.Unlock()

	newNonceInt, err := strconv.Atoi(updatedFs.Nonce)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	currentNonceInt, err := strconv.Atoi(spFs.Nonce)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	if newNonceInt < currentNonceInt {
		return logger.CreateDetails(location, fmt.Errorf("%v fs info is up to date", spAddress))
	}

	nonceHex := strconv.FormatInt(int64(newNonceInt), 16)

	nonceBytes, err := hex.DecodeString(nonceHex)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	nonce32 := make([]byte, 32-len(nonceBytes))
	nonce32 = append(nonce32, nonceBytes...)

	sort.Strings(updatedFs.NewFs)

	concatFsHashesBuilder := strings.Builder{}

	for _, hash := range updatedFs.NewFs {
		concatFsHashesBuilder.WriteString(hash)
	}

	fsTreeNonceBytes := append([]byte(concatFsHashesBuilder.String()), nonce32...)
	fsTreeNonceHash := sha256.Sum256(fsTreeNonceBytes)

	signedFsys, err := hex.DecodeString(signedFileSystem)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	err = checkDataSign(spAddress, signedFsys, fsTreeNonceHash)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	fsRootHash, fsTree, err := shared.CalcRootHash(updatedFs.NewFs)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	signedRootHash, err := hex.DecodeString(updatedFs.SignedFsRootHash)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	fsRootBytes, err := hex.DecodeString(fsRootHash)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	fsRootNonceBytes := append(fsRootBytes, nonce32...)

	hash := sha256.Sum256(fsRootNonceBytes)

	err = checkDataSign(spAddress, signedRootHash, hash)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	shared.MU.Lock()

	spFsFile, err = os.Create(filepath.Join(addressPath, paths.SpFsFilename))
	if err != nil {
		shared.MU.Unlock()
		return logger.CreateDetails(location, err)
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
		return logger.CreateDetails(location, err)
	}

	spFsFile.Sync()
	shared.MU.Unlock()

	return nil
}

// ====================================================================================

//Restore certain file size
func RestoreMemoryInfo(pathToConfig string, intFileSize int) {
	location := "files.restoreMemoryInfo->"

	shared.MU.Lock()
	confFile, fileBytes, err := shared.ReadFile(pathToConfig)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(location, err))
		return
	}
	defer confFile.Close()

	var nodeConfig config.SecondaryNodeConfig

	err = json.Unmarshal(fileBytes, &nodeConfig)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(location, err))
		return
	}

	nodeConfig.UsedStorageSpace -= int64(intFileSize)

	err = config.Save(confFile, nodeConfig)
	if err != nil {
		shared.MU.Unlock()
		logger.Log(logger.CreateDetails(location, err))
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

func checkDataSign(spAddress string, signature []byte, hash [32]byte) error {
	sigPublicKey, err := crypto.SigToPub(hash[:], signature)
	if err != nil {
		return err
	}

	signatureAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if spAddress != signatureAddress.String() {
		return err
	}

	return nil
}

// ====================================================================================

func deleteParts(addressPath string, fileHashes []string) {
	logger.Log("deleting file parts after error...")

	for _, hash := range fileHashes {
		pathToFile := filepath.Join(addressPath, hash)

		os.Remove(pathToFile)
	}
}

// ====================================================================================

func saveSpFsInfo(addressPath string, spData *shared.StorageProviderData) error {
	const location = "files.saveSpFsInfo->"

	stat, err := os.Stat(addressPath)
	err = shared.CheckStatErr(err)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	if stat == nil {
		err = os.Mkdir(addressPath, 0700)
		if err != nil {
			return logger.CreateDetails(location, err)
		}
	}

	shared.MU.Lock()
	defer shared.MU.Unlock()

	spFsFile, err := os.Create(filepath.Join(addressPath, paths.SpFsFilename))
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	defer spFsFile.Close()

	err = shared.WriteFile(spFsFile, spData)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	spFsFile.Sync()

	return nil
}

// ====================================================================================

func getOneMbHashes(reqFileParts []*multipart.FileHeader) ([]string, error) {
	const location = "files.GetOneMbHashes->"
	eightKBHashes := make([]string, 0, 128)
	oneMBHashes := make([]string, 0, len(reqFileParts))

	for _, reqFilePart := range reqFileParts {

		var buf bytes.Buffer

		rqFile, err := reqFilePart.Open()
		if err != nil {
			return nil, logger.CreateDetails(location, err)
		}

		_, err = io.Copy(&buf, rqFile)
		if err != nil {
			rqFile.Close()
			return nil, logger.CreateDetails(location, err)
		}

		rqFile.Close()

		bufBytes := buf.Bytes()
		eightKBHashes = eightKBHashes[:0]

		for i := 0; i < len(bufBytes); i += eightKB {
			hSum := sha256.Sum256(bufBytes[i : i+eightKB])
			eightKBHashes = append(eightKBHashes, hex.EncodeToString(hSum[:]))
		}

		oneMBHash, _, err := shared.CalcRootHash(eightKBHashes)
		if err != nil {
			return nil, logger.CreateDetails(location, err)
		}

		if reqFilePart.Filename != oneMBHash {
			return nil, logger.CreateDetails(location, err)
		}

		oneMBHashes = append(oneMBHashes, oneMBHash)
	}

	return oneMBHashes, nil
}
