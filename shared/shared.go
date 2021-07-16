package shared

import (
	"bufio"
	"crypto/sha256"

	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"github.com/ricochet2200/go-disk-usage/du"
)

type StorageProviderFs struct {
	Nonce        string     `json:"nonce"`
	SignedFsRoot string     `json:"signedFsRoot"`
	Tree         [][][]byte `json:"tree"`
}

var (
	MU           sync.Mutex
	TestMode     = false
	TestPassword = "test"
	TestLimit    = 1
	TestAddress  = "127.0.0.1"
	TestPort     = "8081"
)

func GetAvailableSpace(storagePath string) int {
	var KB = uint64(1024)
	usage := du.NewDiskUsage(storagePath)
	return int(usage.Free() / (KB * KB * KB))
}

// ====================================================================================

func InitPaths() error {
	const logInfo = "shared.InitPaths->"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return logger.CreateDetails(logInfo, err)
	}

	paths.WorkDirPath = filepath.Join(homeDir, paths.WorkDirName)

	paths.AccsDirPath = filepath.Join(paths.WorkDirPath, "accounts")

	return nil
}

// ====================================================================================

func CreateIfNotExistAccDirs() error {
	const logInfo = "shared.CreateIfNotExistAccDirs->"
	statWDP, err := os.Stat(paths.WorkDirPath)
	err = CheckStatErr(err)
	if err != nil {
		return logger.CreateDetails(logInfo, err)
	}

	if statWDP == nil {
		err = os.MkdirAll(paths.WorkDirPath, os.ModePerm|os.ModeDir)
		if err != nil {
			return logger.CreateDetails(logInfo, err)
		}
	}

	statADP, err := os.Stat(paths.AccsDirPath)
	err = CheckStatErr(err)
	if err != nil {
		return logger.CreateDetails(logInfo, err)
	}

	if statADP == nil {
		err = os.MkdirAll(paths.AccsDirPath, os.ModePerm|os.ModeDir)
		if err != nil {
			return logger.CreateDetails(logInfo, err)
		}
	}

	return nil
}

// ====================================================================================

func CheckStatErr(statErr error) error {
	if statErr == nil {
		return nil
	}

	errParts := strings.Split(statErr.Error(), ":")

	if len(errParts) == 3 && strings.Trim(errParts[2], " ") == "The system cannot find the file specified." {
		return nil
	}

	if len(errParts) == 2 && strings.Trim(errParts[1], " ") == "no such file or directory" {
		return nil
	}

	return statErr
}

// ====================================================================================

func ContainsAccount(accounts []string, address string) bool {
	for _, a := range accounts {
		if a == address {
			return true
		}
	}
	return false
}

// ====================================================================================

func ReadFromConsole() (string, error) {
	const logInfo = "shared.ReadFromConsole->"
	fmt.Print("Enter value here: ")
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", logger.CreateDetails(logInfo, err)
	}

	// remove the delimiter from the string
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")

	return input, nil
}

// ====================================================================================

func CalcRootHash(hashArr []string) (string, [][][]byte, error) {
	const logInfo = "shared.CalcRootHash->"

	arrLen := len(hashArr)

	if arrLen == 0 {
		return "", nil, logger.CreateDetails(logInfo, errors.New("hash array is empty"))
	}

	base := make([][]byte, 0, arrLen+1)

	emptyValue, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		return "", nil, logger.CreateDetails(logInfo, err)
	}

	for _, v := range hashArr {
		decoded, err := hex.DecodeString(v)
		if err != nil {
			return "", nil, logger.CreateDetails(logInfo, err)
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

func GetHashPassword(password string) string {
	pBytes := sha256.Sum256([]byte(password))
	return hex.EncodeToString(pBytes[:])
}

// ====================================================================================
