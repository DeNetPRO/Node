package shared

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ricochet2200/go-disk-usage/du"
)

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
