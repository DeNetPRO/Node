package config_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
		t.Fatal(err)
	}

	defer r.Close()
	defer w.Close()

	_, err = w.WriteString("1\n")
	if err != nil {
		t.Fatal(err)
	}

	os.Stdin = r

	got, err := config.SelectNetwork()
	if err != nil {
		fmt.Println(err)
		t.Fatal()
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
			t.Fatal(err)
		}

		defer r.Close()
		defer w.Close()

		_, err = w.WriteString("5\n")
		if err != nil {
			t.Fatal(err)
		}

		os.Stdin = r

		err = config.SetStorageLimit(&configStruct, config.UpdateStatus)
		if err != nil {
			fmt.Println(err)
			t.Fatal(err)
		}

		require.Equal(t, 5, configStruct.StorageLimit)
	})

	t.Run("negative value", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		defer r.Close()
		defer w.Close()

		_, err = w.WriteString("-1\n")
		if err != nil {
			t.Fatal(err)
		}

		os.Stdin = r

		err = config.SetStorageLimit(&configStruct, config.UpdateStatus)

		require.NotEmpty(t, err)
	})

	t.Run("zero value", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		defer r.Close()
		defer w.Close()

		_, err = w.WriteString("0\n")
		if err != nil {
			t.Fatal(err)
		}

		os.Stdin = r

		err = config.SetStorageLimit(&configStruct, config.UpdateStatus)

		require.NotEmpty(t, err)
	})

	t.Run("too big value", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		defer r.Close()
		defer w.Close()

		_, err = w.WriteString("10000\n")
		if err != nil {
			t.Fatal(err)
		}

		os.Stdin = r

		err = config.SetStorageLimit(&configStruct, config.UpdateStatus)

		require.NotEmpty(t, err)
	})

}

func TestConfigSetIP(t *testing.T) {
	configStruct := config.TestConfig

	require.Equal(t, "127.0.0.1", configStruct.IpAddress)

	t.Run("valid ip address", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		defer r.Close()
		defer w.Close()

		_, err = w.WriteString("91.123.123.123\n")
		if err != nil {
			t.Fatal(err)
		}

		os.Stdin = r

		err = config.SetIpAddr(&configStruct, config.UpdateStatus)
		if err != nil {
			t.Fatal(err)
		}

		require.Equal(t, "91.123.123.123", configStruct.IpAddress)
	})

	t.Run("local ip address", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		defer r.Close()
		defer w.Close()

		_, err = w.WriteString("192.168.1.1\n")
		if err != nil {
			t.Fatal(err)
		}

		os.Stdin = r

		err = config.SetIpAddr(&configStruct, config.UpdateStatus)

		require.NotEmpty(t, err)
	})

	t.Run("wrong ip address", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		defer r.Close()
		defer w.Close()

		_, err = w.WriteString("11292.168.1.1\n")
		if err != nil {
			t.Fatal(err)
		}

		os.Stdin = r

		err = config.SetIpAddr(&configStruct, config.UpdateStatus)

		require.NotEmpty(t, err)
	})

}

func TestConfigSetPort(t *testing.T) {
	configStruct := config.TestConfig

	require.Equal(t, "55050", configStruct.HTTPPort)

	t.Run("valid port", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		defer r.Close()
		defer w.Close()

		_, err = w.WriteString("55051\n")
		if err != nil {
			t.Fatal(err)
		}

		os.Stdin = r

		err = config.SetPort(&configStruct, config.UpdateStatus)
		if err != nil {
			t.Fatal(err)
		}

		require.Equal(t, "55051", configStruct.HTTPPort)
	})

	t.Run("invalid port", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		defer r.Close()
		defer w.Close()

		_, err = w.WriteString("-1\n")
		if err != nil {
			t.Fatal(err)
		}

		os.Stdin = r

		err = config.SetPort(&configStruct, config.UpdateStatus)

		require.NotEmpty(t, err)
	})

}

func TestConfigSave(t *testing.T) {

	configStruct := config.TestConfig

	configStruct.Network = "kovan"
	configStruct.Address = "0x0000000000000000000000000000000000000000"
	configStruct.HTTPPort = "66056"
	configStruct.IpAddress = "102.103.104.105"

	configFile, err := os.OpenFile(filepath.Join(paths.ConfigDirPath, paths.ConfFileName), os.O_RDWR, 0777)
	if err != nil {
		t.Fatal(err)
	}
	defer configFile.Close()

	err = config.Save(configFile, configStruct)
	if err != nil {
		configFile.Close()
		t.Fatal(err)
	}

	var updatedConfig config.NodeConfig
	data, err := os.ReadFile(filepath.Join(paths.ConfigDirPath, paths.ConfFileName))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(data, &updatedConfig)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, configStruct, updatedConfig)
}
