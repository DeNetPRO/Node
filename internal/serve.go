package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"

	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"github.com/ethereum/go-ethereum/crypto"
)

func ServeFile(spAddress, fileKey, signatureFromReq string) (string, error) {
	const logLoc = "internal.ServeFile->"

	signature, err := hex.DecodeString(signatureFromReq)
	if err != nil {
		return "", logger.CreateDetails(logLoc, err)
	}

	hash := sha256.Sum256([]byte(fileKey + spAddress))

	sigPublicKey, err := crypto.SigToPub(hash[:], signature)
	if err != nil {
		return "", logger.CreateDetails(logLoc, err)
	}

	signatureAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if spAddress != signatureAddress.String() {
		return "", logger.CreateDetails(logLoc, shared.ErrWrongSignature)
	}

	pathToFile := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, spAddress, fileKey)

	_, err = os.Stat(pathToFile)
	if err != nil {
		return "", logger.CreateDetails(logLoc, err)
	}

	return pathToFile, nil
}
