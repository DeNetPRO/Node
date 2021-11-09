package shared

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ricochet2200/go-disk-usage/du"
)

type StorageProviderData struct {
	Address               string     `json:"address"`
	Nonce                 string     `json:"nonce"`
	SignedFsRootNonceHash string     `json:"signedFsRoot"`
	Tree                  [][][]byte `json:"tree"`
	Fs                    []string
}

var (
	NodeAddr common.Address
	MU       sync.Mutex
)

//Return nodes available space in GB
func GetAvailableSpace() (int, error) {

	const location = "shared.GetAvailableSpace ->"
	const KB = uint64(1024)

	fmt.Println(paths.StoragePaths)

	_, err := os.Stat(paths.WorkDirPath)
	if err != nil {
		return 0, logger.CreateDetails(location, err)
	}

	usage := du.NewDiskUsage(paths.WorkDirPath)
	return int(usage.Available() / (KB * KB * KB)), nil
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("HeapInuse = %v MiB", bToMb(m.HeapInuse))
	fmt.Printf("\tHeapIdle = %v MiB", bToMb(m.HeapIdle))
	fmt.Printf("\tHeapReleased = %v MiB\n", bToMb(m.HeapReleased))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
