package meminfo_test

import (
	"log"
	"os"
	"testing"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/shared"
)

func TestMain(m *testing.M) {
	shared.TestModeOn()
	defer shared.TestModeOff()

	err := paths.Init()
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

// func TestRestoreNodeMemory(t *testing.T) {
// 	fileSize := 1024 * 1024

// 	confFile, nodeConfig, err := getConfig()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	want := nodeConfig.UsedStorageSpace

// 	nodeConfig.UsedStorageSpace += int64(fileSize)

// 	err = config.Save(confFile, *nodeConfig)
// 	if err != nil {
// 		confFile.Close()
// 		t.Fatal(err)
// 	}

// 	confFile.Close()

// 	meminfo.Restore(configPath, fileSize)

// 	confFile, nodeConfig, err = getConfig()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	confFile.Close()

// 	require.Equal(t, want, nodeConfig.UsedStorageSpace)
// }
