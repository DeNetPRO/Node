package shared

import (
	"github.com/ricochet2200/go-disk-usage/du"
)

func GetAvailableSpace(storagePath string) int {
	var KB = uint64(1024)
	usage := du.NewDiskUsage(storagePath)
	return int(usage.Free() / (KB * KB * KB))
}
