package hash

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"mime/multipart"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
)

const eightKB = 8192

//Hash password with SHA-256
func Password(password string) string {
	pBytes := sha256.Sum256([]byte(password))
	return hex.EncodeToString(pBytes[:])
}

// ====================================================================================

// OneMbParts calculates and returns array of file part's root hash info.
func OneMbParts(reqFileParts []*multipart.FileHeader) ([]string, error) {
	const location = "hash.GetOneMbHashes->"
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

		oneMBHash, _, err := CalcRoot(eightKBHashes)
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

// ====================================================================================

//CalcRoot calculates root hash by building merkle tree
func CalcRoot(hashArr []string) (string, [][][]byte, error) {
	const location = "shared.CalcRootHash->"

	arrLen := len(hashArr)

	if arrLen == 0 {
		return "", nil, logger.CreateDetails(location, errors.New("hash array is empty"))
	}

	base := make([][]byte, 0, arrLen+1)

	emptyValue, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		return "", nil, logger.CreateDetails(location, err)
	}

	for _, v := range hashArr {
		decoded, err := hex.DecodeString(v)
		if err != nil {
			return "", nil, logger.CreateDetails(location, err)
		}
		base = append(base, decoded)
	}

	if len(base)%2 != 0 {
		base = append(base, emptyValue)
	}

	resByte := make([][][]byte, 0, len(base)*2-1)

	resByte = append(resByte, base)

	for len(resByte[len(resByte)-1]) != 1 {
		prevList := resByte[len(resByte)-1]
		resByte = append(resByte, [][]byte{})
		r := len(prevList) / 2

		for i := 0; i < r; i++ {
			a := prevList[i*2]
			b := prevList[i*2+1]

			concatBytes := append(a, b...)
			hSum := sha256.Sum256(concatBytes)

			resByte[len(resByte)-1] = append(resByte[len(resByte)-1], hSum[:])
		}

		if len(resByte[len(resByte)-1])%2 != 0 && len(prevList) > 2 {
			resByte[len(resByte)-1] = append(resByte[len(resByte)-1], emptyValue)
		}
	}

	return hex.EncodeToString(resByte[len(resByte)-1][0]), resByte, nil
}

// ====================================================================================
