package config

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"testing"

// 	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
// 	"git.denetwork.xyz/DeNet/dfile-secondary-node/shared"
// 	"github.com/stretchr/testify/require"
// )

// var (
// 	WorkDir     = "tmp"
// 	AccountsDir = "accounts"
// )

// func TestMain(m *testing.M) {
// 	shared.TestModeOn()
// 	defer shared.TestModeOff()

// 	os.RemoveAll(WorkDir)

// 	err := os.Mkdir(WorkDir, 0777)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	paths.WorkDirName = WorkDir
// 	paths.WorkDirPath = WorkDir
// 	paths.AccsDirPath = filepath.Join(WorkDir, AccountsDir)

// 	exitVal := m.Run()

// 	err = os.RemoveAll(WorkDir)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	os.Exit(exitVal)
// }

// func TestConfigCreate(t *testing.T) {
// 	address := "some_address"

// 	pathToConfig := filepath.Join(paths.AccsDirPath, address, paths.ConfDirName)
// 	os.MkdirAll(pathToConfig, 0777)

// 	config, err := Create(address)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	require.NotEmpty(t, config, "config didn't create")
// 	require.Equal(t, address, config.Address, "account address don't match")
// 	require.Equal(t, shared.TestIP, config.IpAddress, "ip address is incorrect, want: ", shared.TestIP, " got: ", config.IpAddress)
// 	require.Equal(t, shared.TestPort, config.HTTPPort, "port is incorrect, want: ", shared.TestPort, " got: ", config.HTTPPort)
// 	require.Equal(t, shared.TestNetwork, config.Network, "network is incorrect, want: ", shared.TestNetwork, " got: ", config.Network)
// 	require.Equal(t, shared.TestStorageLimit, config.StorageLimit, "storage limit is incorrect, want: ", shared.TestStorageLimit, " got: ", config.StorageLimit)
// 	require.Empty(t, config.UsedStorageSpace, "used storage space must be 0 instead ", config.UsedStorageSpace)
// }

// func TestSelectNetwork(t *testing.T) {
// 	path := filepath.Join(WorkDir, "stdin")
// 	file, err := os.Create(path)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	file.WriteString("1\n")
// 	file.Sync()
// 	file.Seek(0, 0)

// 	os.Stdin = file

// 	got, err := SelectNetwork()
// 	if err != nil {
// 		fmt.Println(err)
// 		t.Error()
// 	}

// 	file.Close()
// 	os.Remove(path)

// 	want := []string{"polygon", "mumbai", "kovan"}

// 	require.Contains(t, want, got)
// }

// func TestConfigSetStorageLimit(t *testing.T) {
// 	address := "some_address"

// 	pathToConfig := filepath.Join(paths.AccsDirPath, address, paths.ConfDirName)
// 	os.RemoveAll(pathToConfig)
// 	os.MkdirAll(pathToConfig, 0777)

// 	config, err := Create(address)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	want := 2

// 	path := filepath.Join(WorkDir, "stdin")
// 	file, err := os.Create(path)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	file.WriteString("2\n")
// 	file.Sync()
// 	file.Seek(0, 0)

// 	os.Stdin = file

// 	err = SetStorageLimit(pathToConfig, UpdateStatus, &config)
// 	if err != nil {
// 		fmt.Println(err)
// 		t.Error(err)
// 	}

// 	file.Close()
// 	os.Remove(path)

// 	got := config.StorageLimit

// 	require.Equal(t, want, got)
// }

// func TestConfigSetNegativeStorageLimit(t *testing.T) {
// 	address := "some_address"

// 	pathToConfig := filepath.Join(paths.AccsDirPath, address, paths.ConfDirName)
// 	os.RemoveAll(pathToConfig)
// 	os.MkdirAll(pathToConfig, 0777)

// 	config, err := Create(address)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	path := filepath.Join(WorkDir, "stdin")
// 	file, err := os.Create(path)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	file.WriteString("-1\n")
// 	file.Sync()
// 	file.Seek(0, 0)

// 	os.Stdin = file

// 	defer file.Close()
// 	defer os.Remove(path)

// 	err = SetStorageLimit(pathToConfig, UpdateStatus, &config)

// 	require.NotEmpty(t, err)
// }

// func TestConfigSetIP(t *testing.T) {
// 	address := "some_address"

// 	pathToConfig := filepath.Join(paths.AccsDirPath, address, paths.ConfDirName)
// 	os.RemoveAll(pathToConfig)
// 	os.MkdirAll(pathToConfig, 0777)

// 	config, err := Create(address)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	path := filepath.Join(WorkDir, "stdin")
// 	file, err := os.Create(path)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	ip := "123.123.123.123"
// 	file.WriteString(ip + "\n")
// 	file.Sync()
// 	file.Seek(0, 0)

// 	os.Stdin = file

// 	defer file.Close()
// 	defer os.Remove(path)

// 	err = SetIpAddr(&config, UpdateStatus)
// 	if err != nil {
// 		fmt.Println(err)
// 		t.Error(err)
// 	}

// 	require.Equal(t, ip, config.IpAddress)
// }

// func TestConfigSetLocalIP(t *testing.T) {
// 	address := "some_address"

// 	pathToConfig := filepath.Join(paths.AccsDirPath, address, paths.ConfDirName)
// 	os.RemoveAll(pathToConfig)
// 	os.MkdirAll(pathToConfig, 0777)

// 	config, err := Create(address)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	path := filepath.Join(WorkDir, "stdin")
// 	file, err := os.Create(path)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	ip := "127.0.0.1"
// 	file.WriteString(ip + "\n")
// 	file.Sync()
// 	file.Seek(0, 0)

// 	os.Stdin = file

// 	defer file.Close()
// 	defer os.Remove(path)

// 	err = SetIpAddr(&config, UpdateStatus)

// 	require.NotEmpty(t, err)
// }

// func TestConfigSetWrongIP(t *testing.T) {
// 	address := "some_address"

// 	pathToConfig := filepath.Join(paths.AccsDirPath, address, paths.ConfDirName)
// 	os.RemoveAll(pathToConfig)
// 	os.MkdirAll(pathToConfig, 0777)

// 	config, err := Create(address)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	path := filepath.Join(WorkDir, "stdin")
// 	file, err := os.Create(path)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	ip := "1227.0.0.1"
// 	file.WriteString(ip + "\n")
// 	file.Sync()
// 	file.Seek(0, 0)

// 	os.Stdin = file

// 	defer file.Close()
// 	defer os.Remove(path)

// 	err = SetIpAddr(&config, UpdateStatus)

// 	require.NotEmpty(t, err)
// }

// func TestConfigSetPort(t *testing.T) {
// 	address := "some_address"

// 	pathToConfig := filepath.Join(paths.AccsDirPath, address, paths.ConfDirName)
// 	os.RemoveAll(pathToConfig)
// 	os.MkdirAll(pathToConfig, 0777)

// 	config, err := Create(address)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	path := filepath.Join(WorkDir, "stdin")
// 	file, err := os.Create(path)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	port := "55051"
// 	file.WriteString(port + "\n")
// 	file.Sync()
// 	file.Seek(0, 0)

// 	os.Stdin = file

// 	defer file.Close()
// 	defer os.Remove(path)

// 	err = SetPort(&config, UpdateStatus)
// 	if err != nil {
// 		fmt.Println(err)
// 		t.Error(err)
// 	}

// 	require.Equal(t, port, config.HTTPPort)
// }

// func TestConfigSeWrongtPort(t *testing.T) {
// 	address := "some_address"

// 	pathToConfig := filepath.Join(paths.AccsDirPath, address, paths.ConfDirName)
// 	os.RemoveAll(pathToConfig)
// 	os.MkdirAll(pathToConfig, 0777)

// 	config, err := Create(address)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	path := filepath.Join(WorkDir, "stdin")
// 	file, err := os.Create(path)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	port := "-1"
// 	file.WriteString(port + "\n")
// 	file.Sync()
// 	file.Seek(0, 0)

// 	os.Stdin = file

// 	defer file.Close()
// 	defer os.Remove(path)

// 	err = SetPort(&config, UpdateStatus)
// 	require.NotEmpty(t, err)
// }

// func TestConfigSave(t *testing.T) {
// 	address := "some_address"

// 	pathToConfig := filepath.Join(paths.AccsDirPath, address, paths.ConfDirName)
// 	os.RemoveAll(pathToConfig)
// 	os.MkdirAll(pathToConfig, 0777)

// 	config, err := Create(address)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	config.Network = "kovan"
// 	config.Address = "0x"
// 	config.HTTPPort = "66056"
// 	config.IpAddress = "102.103.104.105"

// 	path := filepath.Join(pathToConfig, paths.ConfFileName)

// 	configFile, err := os.OpenFile(path, os.O_RDWR, 0777)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	err = Save(configFile, config)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	var got NodeConfig
// 	data, err := os.ReadFile(path)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	err = json.Unmarshal(data, &got)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	require.Equal(t, config, got)
// }
