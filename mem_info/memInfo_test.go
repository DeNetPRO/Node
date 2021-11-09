package meminfo_test

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"testing"

	memInfo "git.denetwork.xyz/DeNet/dfile-secondary-node/mem_info"
	tstpkg "git.denetwork.xyz/DeNet/dfile-secondary-node/tst_pkg"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/config"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	tstpkg.TestModeOn()

	defer tstpkg.TestModeOff()

	err := paths.Init()
	if err != nil {
		log.Fatal(err)
	}

	_, err = config.Create(tstpkg.TestAccAddr)
	if err != nil {
		log.Fatal(err)
	}

	exitVal := m.Run()

	err = os.RemoveAll(paths.WorkDirPath)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitVal)
}

func TestRestoreNodeMemory(t *testing.T) {

	const fileSize = 1024

	confFilePath := filepath.Join(paths.ConfigDirPath, paths.ConfFileName)

	memInfo.Restore(confFilePath, fileSize)

	configFileBytes, err := os.ReadFile(confFilePath)
	if err != nil {
		t.Fatal(err)
	}

	var config config.NodeConfig

	err = json.Unmarshal(configFileBytes, &config)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, int64(8976), config.UsedStorageSpace)

}
