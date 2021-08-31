package fsysinfo

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

	dnetsignature "git.denetwork.xyz/dfile/dfile-secondary-node/dnet_signature"
	"git.denetwork.xyz/dfile/dfile-secondary-node/errs"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	nodeFile "git.denetwork.xyz/dfile/dfile-secondary-node/node_file"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
)

type UpdatedFsInfo struct {
	NewFs            []string
	Nonce            string
	SignedFsRootHash string
}

// UpdateFileSystemInfo updates Storage Provider's nounce and file system's root hash info.
func Update(updatedFs *UpdatedFsInfo, spAddress, signedFileSystem string) error {
	const location = "files.UpdateFileSystemInfo->"

	addressPath := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spAddress)

	_, err := os.Stat(addressPath)
	if err != nil {
		return logger.CreateDetails(location, errors.New("no files of "+spAddress))
	}

	shared.MU.Lock()
	spFsFile, fileBytes, err := nodeFile.Read(filepath.Join(addressPath, paths.SpFsFilename))
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

	err = dnetsignature.Check(spAddress, signedFsys, fsTreeNonceHash)
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

	err = dnetsignature.Check(spAddress, signedRootHash, hash)
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

	err = nodeFile.Write(spFsFile, spFs)
	if err != nil {
		shared.MU.Unlock()
		return logger.CreateDetails(location, err)
	}

	spFsFile.Sync()
	shared.MU.Unlock()

	return nil
}

// ====================================================================================

// SaveSpFsInfo saves Storage Provider file system and nounce info from the request.
func Save(addressPath string, spData *shared.StorageProviderData) error {
	const location = "files.saveSpFsInfo->"

	stat, err := os.Stat(addressPath)
	err = errs.CheckStatErr(err)
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

	err = nodeFile.Write(spFsFile, spData)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	spFsFile.Sync()

	return nil
}
