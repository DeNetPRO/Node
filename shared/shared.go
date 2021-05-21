package shared

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ricochet2200/go-disk-usage/du"
)

func GetAvailableSpace(storagePath string) int {
	var KB = uint64(1024)
	usage := du.NewDiskUsage(storagePath)
	return int(usage.Free() / (KB * KB * KB))
}

var (
	WorkDir string
	AccDir  string
)

// GetHomeDirectory return path to the home directory of dfile
func CreateIfNotExistAccDirs() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Fatal error")
	}

	workDir := filepath.Join(homeDir, "dfileTest")

	_, err = os.Stat(workDir)
	if err != nil {
		errPart := strings.Split(err.Error(), ":")

		if strings.Trim(errPart[1], " ") != "no such file or directory" {
			log.Fatal("Fatal error")
		}

		err = os.MkdirAll(workDir, os.ModePerm|os.ModeDir)
		if err != nil {
			log.Fatal("Fatal error")
		}
	}

	WorkDir = workDir

	accDir := filepath.Join(WorkDir, "accounts")

	_, err = os.Stat(accDir)
	if err != nil {
		errPart := strings.Split(err.Error(), ":")

		if strings.Trim(errPart[1], " ") != "no such file or directory" {
			log.Fatal("Fatal error")
		}

		err = os.MkdirAll(accDir, os.ModePerm|os.ModeDir)
		if err != nil {
			log.Fatal("Fatal error")
		}
	}

	AccDir = accDir

}

func ContainsAccount(accounts []string, address string) bool {
	for _, a := range accounts {
		if a == address {
			return true
		}
	}
	return false
}

func ReadFromConsole() (string, error) {
	fmt.Print("Enter value here: ")
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// remove the delimiter from the string
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")

	return input, err
}

func CalcRootHash(hashArr []string) (string, error) {
	hashArrLen := len(hashArr)

	i := 0
	j := i + 1

	lvlCount := 2
	upperLvl := hashArrLen + hashArrLen/lvlCount

	var decodedJ []byte

	for len(hashArr) < hashArrLen*2-1 {

		decodedI, err := hex.DecodeString(hashArr[i])
		if err != nil {
			return "", err
		}

		if upperLvl < hashArrLen*2 && len(hashArr) == upperLvl {

			if upperLvl%2 != 0 {
				hashArr = append(hashArr, "0000000000000000000000000000000000000000000000000000000000000000")

				decodedJ, err = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
				if err != nil {
					return "", err
				}

				hashArrLen += 1

			} else {
				decodedJ, err = hex.DecodeString(hashArr[j])
				if err != nil {
					return "", err
				}
			}

			lvlCount *= 2
			upperLvl = upperLvl + hashArrLen/lvlCount
		} else {
			decodedJ, err = hex.DecodeString(hashArr[j])
			if err != nil {
				return "", err
			}
		}

		concatBytes := append(decodedI, decodedJ...)

		hSum := sha256.Sum256(concatBytes)
		hashArr = append(hashArr, hex.EncodeToString(hSum[:]))

		i++
	}

	return hashArr[len(hashArr)-1], nil

}
