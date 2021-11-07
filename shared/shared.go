package shared

import (
	"fmt"
	"runtime"
	"sync"

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
	NodeAddr         common.Address
	MU               sync.Mutex
	TestMode         = false
	TestWorkDirName  = "denet-node-test"
	TestIP           = "127.0.0.1"
	TestPort         = "55050"
	TestStorageLimit = 1
	TestNetwork      = "kovan"
	TestPrivateKey   = "16f98d96422dd7f21965755bd64c9dcd9cfc5d36e029002d9cc579f42511c7ed"
	TestPassword     = "123"
)

//Return nodes available space in GB
func GetAvailableSpace(storagePath string) int {
	var KB = uint64(1024)
	usage := du.NewDiskUsage(storagePath)
	return int(usage.Free() / (KB * KB * KB))
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

func TestModeOn() {
	TestMode = true
}

func TestModeOff() {
	TestMode = false
}
