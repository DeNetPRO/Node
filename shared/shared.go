package shared

import (
	"bufio"
	"bytes"
	"io"
	"mime/multipart"

	"github.com/minio/sha256-simd"

	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ricochet2200/go-disk-usage/du"
)

const eightKB = 8192

type StorageProviderData struct {
	Address      string     `json:"address"`
	Nonce        string     `json:"nonce"`
	SignedFsRoot string     `json:"signedFsRoot"`
	Tree         [][][]byte `json:"tree"`
	Fs           []string
}

var (
	NodeAddr common.Address
	MU       sync.Mutex

	//Tests variables used in TestMode
	TestMode       = false
	TestPassword   = "test"
	TestLimit      = 1
	TestAddress    = "127.0.0.1"
	TestPort       = "8081"
	TestPrivateKey = "16f98d96422dd7f21965755bd64c9dcd9cfc5d36e029002d9cc579f42511c7ed"
)

//Return nodes available space in GB
func GetAvailableSpace(storagePath string) int {
	var KB = uint64(1024)
	usage := du.NewDiskUsage(storagePath)
	return int(usage.Free() / (KB * KB * KB))
}

// ====================================================================================

func ReadFromConsole() (string, error) {
	const location = "shared.ReadFromConsole->"
	fmt.Print("Enter value here: ")
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", logger.CreateDetails(location, err)
	}

	// remove the delimiter from the string
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")

	return input, nil
}

// ====================================================================================

//Calculate root hash by making merkle tree
func CalcRootHash(hashArr []string) (string, [][][]byte, error) {
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

//Hash password with SHA-256
func GetHashPassword(password string) string {
	pBytes := sha256.Sum256([]byte(password))
	return hex.EncodeToString(pBytes[:])
}

// ====================================================================================

// GetOneMbHashes calculates and returns array of file part's root hash info.
func GetOneMbHashes(reqFileParts []*multipart.FileHeader) ([]string, error) {
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

		oneMBHash, _, err := CalcRootHash(eightKBHashes)
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
