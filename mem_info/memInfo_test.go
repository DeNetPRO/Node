package meminfo_test

import (
	"log"
	"os"
	"testing"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/config"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/shared"
)

func TestMain(m *testing.M) {
	shared.TestModeOn()
	defer shared.TestModeOff()

	_, err := config.Create(shared.TestAccAddr)
	if err != nil {
		log.Fatal(err)
	}

	exitVal := m.Run()

	// err = os.RemoveAll(paths.WorkDirPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	os.Exit(exitVal)
}

// func TestRestoreNodeMemory(t *testing.T) {

// 	memInfo.Restore(configPath, fileSize)

// 	confFile, nodeConfig, err = getConfig()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	confFile.Close()

// 	require.Equal(t, want, nodeConfig.UsedStorageSpace)
// }
