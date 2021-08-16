package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"github.com/ethereum/go-ethereum/crypto"
)

type UpdatedFsInfo struct {
	NewFs            []string
	Nonce            string
	SignedFsRootHash string
}

func UpdateFileSystemInfo(updatedFs *UpdatedFsInfo, spAddress, signedFileSystem string) error {
	const logLoc = "internal.UpdateFileSystemInfo->"

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
