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
	NodeAddr             common.Address
	MU                   sync.Mutex
	TestMode             = false
	TestWorkDirName      = "denet-node-test"
	TestIP               = "127.0.0.1"
	TestPort             = "55050"
	TestStorageLimit     = 1
	TestNetwork          = "kovan"
	TestPassword         = "123"
	TestNodeConfDir      = "./"
	TestAccAddr          = "0x3D4dFC62B4F0BDD7730fCB68FaC150E71D76bc24"
	TestPrivateKey       = "16f98d96422dd7f21965755bd64c9dcd9cfc5d36e029002d9cc579f42511c7ed"
	TestSecretKey        = []byte{66, 180, 56, 47, 96, 21, 163, 67, 241, 114, 1, 225, 108, 61, 241, 226, 250, 28, 194, 158, 234, 62, 230, 223, 251, 50, 73, 76, 245, 218, 143, 115}
	TestPKHash           = []byte{166, 130, 151, 14, 57, 26, 220, 249, 192, 230, 178, 57, 69, 112, 95, 215, 238, 209, 203, 160, 153, 131, 179, 84, 254, 192, 244, 101, 221, 161, 9, 80, 78, 2, 215, 181, 73, 131, 244, 221, 204, 127, 249, 128, 178, 53, 213, 80, 200, 75, 110, 83, 92, 171, 184, 10, 242, 169, 37, 21}
	TestUsedStorageSpace = 10000
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
