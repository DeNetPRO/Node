package meminfo_test

// func TestMain(m *testing.M) {
// 	shared.TestModeOn()

// 	defer shared.TestModeOff()

// 	err := paths.Init()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	_, err = config.Create(shared.TestAccAddr)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	exitVal := m.Run()

// 	err = os.RemoveAll(paths.WorkDirPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	os.Exit(exitVal)
// }

// func TestRestoreNodeMemory(t *testing.T) {

// 	const fileSize = 1024

// 	confFilePath := filepath.Join(paths.ConfigDirPath, paths.ConfFileName)

// 	memInfo.Restore(confFilePath, fileSize)

// 	configFileBytes, err := os.ReadFile(confFilePath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	var config config.NodeConfig

// 	err = json.Unmarshal(configFileBytes, &config)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	require.Equal(t, int64(8976), config.UsedStorageSpace)

// }
