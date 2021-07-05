package shared

import (
	"bufio"
	"crypto/sha256"
	"dfile-secondary-node/logger"
	"dfile-secondary-node/paths"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ricochet2200/go-disk-usage/du"
)

type StorageInfo struct {
	Nonce        string     `json:"nonce"`
	SignedFsRoot string     `json:"signedFsRoot"`
	Tree         [][][]byte `json:"tree"`
}

var (
	MU sync.Mutex
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
		return fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
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
		return fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
	}

	if statWDP == nil {
		err = os.MkdirAll(paths.WorkDirPath, os.ModePerm|os.ModeDir)
		if err != nil {
			return fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
		}
	}

	statADP, err := os.Stat(paths.AccsDirPath)
	err = CheckStatErr(err)
	if err != nil {
		return fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
	}

	if statADP == nil {
		err = os.MkdirAll(paths.AccsDirPath, os.ModePerm|os.ModeDir)
		if err != nil {
			return fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
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
		return "", fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
	}

	// remove the delimiter from the string
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")

	return input, nil
}

// ====================================================================================

func CalcRootHash(hashArr []string) (string, [][][]byte, error) {
	const logInfo = "shared.CalcRootHash->"
	resByte := [][][]byte{}
	base := [][]byte{}

	emptyValue, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		return "", resByte, fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
	}

	for _, v := range hashArr {
		decoded, err := hex.DecodeString(v)
		if err != nil {
			return "", resByte, fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
		}
		base = append(base, decoded)
	}

	if len(base)%2 != 0 {
		base = append(base, emptyValue)
	}

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
