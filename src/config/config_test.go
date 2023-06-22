package config_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/DeNet/config"
	"github.com/DeNet/networks"
	nodeTypes "github.com/DeNet/node_types"
	tstpkg "github.com/DeNet/tst_pkg"

	"github.com/stretchr/testify/require"

	"github.com/DeNet/paths"
)

var testConfig nodeTypes.Config

func TestMain(m *testing.M) {
	tstpkg.TestModeOn()
	defer tstpkg.TestModeOff()

	err := paths.Init()
	if err != nil {
		log.Fatal(err)
	}

	testConfig, err = config.Create(tstpkg.Data().AccAddr)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(paths.List().Storages[0], 0700)
	if err != nil {
		log.Fatal(err)
	}

	exitVal := m.Run()

	err = os.RemoveAll(paths.List().WorkDir)
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

	_, err = w.WriteString("2\n")
	if err != nil {
		t.Fatal(err)
	}

	os.Stdin = r

	err = config.SetNetwork(&testConfig)
	if err != nil {
		t.Fatal(err)
	}

	want := networks.List()

	require.Contains(t, want, testConfig.Network)
}

func TestConfigSetStorageLimit(t *testing.T) {

	configStruct := testConfig

	require.Equal(t, 1, configStruct.StorageLimit)

	t.Run("correct value", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		defer r.Close()
		defer w.Close()

		_, err = w.WriteString("2\n")
		if err != nil {
			t.Fatal(err)
		}

		os.Stdin = r

		err = config.SetStorageLimit(&configStruct, config.Stats().Update)
		if err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 3, configStruct.StorageLimit)
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

		err = config.SetStorageLimit(&configStruct, config.Stats().Update)

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

		err = config.SetStorageLimit(&configStruct, config.Stats().Update)

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

		err = config.SetStorageLimit(&configStruct, config.Stats().Update)

		require.NotEmpty(t, err)
	})

}

func TestConfigSetIP(t *testing.T) {
	configStruct := testConfig

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

		err = config.SetIpAddr(&configStruct, config.Stats().Update)
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

		err = config.SetIpAddr(&configStruct, config.Stats().Update)

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

		err = config.SetIpAddr(&configStruct, config.Stats().Update)

		require.NotEmpty(t, err)
	})

}

func TestConfigSetPort(t *testing.T) {
	configStruct := testConfig

	require.Equal(t, ":55050", configStruct.HTTPPort)

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

		err = config.SetPort(&configStruct, config.Stats().Update)
		if err != nil {
			t.Fatal(err)
		}

		require.Equal(t, ":55051", configStruct.HTTPPort)
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

		err = config.SetPort(&configStruct, config.Stats().Update)

		require.NotEmpty(t, err)
	})

}

func TestConfigSave(t *testing.T) {

	configStruct := testConfig

	configStruct.Network = "kovan"
	configStruct.Address = "0x0000000000000000000000000000000000000000"
	configStruct.HTTPPort = "66056"
	configStruct.IpAddress = "102.103.104.105"

	configFile, err := os.OpenFile(paths.List().ConfigFile, os.O_RDWR, 0777)
	if err != nil {
		t.Fatal(err)
	}
	defer configFile.Close()

	err = config.Save(configFile, configStruct)
	if err != nil {
		configFile.Close()
		t.Fatal(err)
	}

	var updatedConfig nodeTypes.Config
	data, err := os.ReadFile(paths.List().ConfigFile)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(data, &updatedConfig)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, configStruct, updatedConfig)
}
