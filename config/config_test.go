package config_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	blckChain "git.denetwork.xyz/DeNet/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/config"
	tstpkg "git.denetwork.xyz/DeNet/dfile-secondary-node/tst_pkg"
	"github.com/stretchr/testify/require"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
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

	err = os.MkdirAll(paths.StoragePaths[0], 0700)
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

func TestSelectNetwork(t *testing.T) {

	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}

	defer r.Close()
	defer w.Close()

	_, err = w.WriteString("1\n")
	if err != nil {
		t.Error(err)
	}

	os.Stdin = r

	got, err := config.SelectNetwork()
	if err != nil {
		fmt.Println(err)
		t.Error()
	}

	want := []string{}

	for net := range blckChain.Networks {
		want = append(want, net)
	}

	require.Contains(t, want, got)
}

func TestConfigSetStorageLimit(t *testing.T) {

	configStruct := config.TestConfig

	require.Equal(t, 1, configStruct.StorageLimit)

	t.Run("correct value", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Error(err)
		}

		defer r.Close()
		defer w.Close()

		_, err = w.WriteString("5\n")
		if err != nil {
			t.Error(err)
		}

		os.Stdin = r

		err = config.SetStorageLimit(&configStruct, config.UpdateStatus)
		if err != nil {
			fmt.Println(err)
			t.Error(err)
		}

		require.Equal(t, 5, configStruct.StorageLimit)
	})

	t.Run("negative value", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Error(err)
		}

		defer r.Close()
		defer w.Close()

		_, err = w.WriteString("-1\n")
		if err != nil {
			t.Error(err)
		}

		os.Stdin = r

		err = config.SetStorageLimit(&configStruct, config.UpdateStatus)

		require.NotEmpty(t, err)
	})

	t.Run("zero value", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Error(err)
		}

		defer r.Close()
		defer w.Close()

		_, err = w.WriteString("0\n")
		if err != nil {
			t.Error(err)
		}

		os.Stdin = r

		err = config.SetStorageLimit(&configStruct, config.UpdateStatus)

		require.NotEmpty(t, err)
	})

}

func TestConfigSetIP(t *testing.T) {
	configStruct := config.TestConfig

	require.Equal(t, "127.0.0.1", configStruct.IpAddress)

	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}

	defer r.Close()
	defer w.Close()

	_, err = w.WriteString("123.123.123.123\n")
	if err != nil {
		t.Error(err)
	}

	os.Stdin = r

	err = config.SetIpAddr(&configStruct, config.UpdateStatus)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	require.Equal(t, "123.123.123.123", configStruct.IpAddress)

}

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
