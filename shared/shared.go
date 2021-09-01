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

type RatingInfo struct {
	Rating                float32                               `json:"rating"`
	StorageProviders      map[string]map[string]*RatingFileInfo `json:"storage_providers"`
	ConnectedNodes        map[string]int                        `json:"connected_nodes"`
	NumberOfAuthorityConn int                                   `json:"nac"`
}

type RatingFileInfo struct {
	FileKey string   `json:"file_key"`
	Nodes   []string `json:"nodes"`
}

var (
	NodeAddr common.Address
	MU       sync.Mutex

	//Tests variables used in TestMode
	TestMode     = false
	TestPassword = "test"
	TestLimit    = 1
	TestAddress  = "127.0.0.1"
	TestPort     = "8081"
)

//Return nodes available space in GB
func GetAvailableSpace(storagePath string) int {
	var KB = uint64(1024)
	usage := du.NewDiskUsage(storagePath)
	return int(usage.Free() / (KB * KB * KB))
}
