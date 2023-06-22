package tstpkg

import (
	nodeTypes "github.com/DeNetPRO/src/node_types"
)

type TstData struct {
	TestMode    bool
	WorkDirName string
	AccAddr     string
	Password    string
	PrivateKey  string
	EncrKey     []byte
	PKHash      []byte
}

var testData = TstData{
	TestMode:    false,
	WorkDirName: "denet-node-test",
	AccAddr:     "0x5Cae405D9A28B51D1bDfFdF6B22c97D2E78dc527",
	Password:    "123",
	PrivateKey:  "e9c9abe3bb861c9d393e369683bc321690a6a7496e81858269d34e141272c4f7",
	EncrKey:     []byte{34, 54, 133, 210, 117, 225, 63, 166, 227, 161, 232, 86, 255, 62, 147, 233, 65, 187, 194, 144, 138, 170, 98, 49, 253, 187, 26, 244, 233, 194, 172, 45},
	PKHash:      []byte{51, 136, 239, 160, 130, 219, 203, 136, 55, 18, 22, 72, 57, 208, 79, 187, 186, 86, 203, 140, 33, 245, 198, 96, 97, 234, 61, 140, 204, 114, 120, 64, 203, 170, 59, 173, 111, 237, 233, 237, 254, 82, 190, 80, 42, 182, 13, 42, 86, 150, 209, 74, 136, 70, 187, 249, 229, 99, 204, 74},
}

func Data() TstData {
	return testData
}

var testConfig = nodeTypes.Config{
	Address:              testData.AccAddr,
	SendBugReports:       false,
	RegisteredInNetworks: map[string]bool{},
	IpAddress:            "127.0.0.1",
	HTTPPort:             ":55050",
	Network:              "kovan",
	StorageLimit:         1,
	StoragePaths:         []string{},
	UsedStorageSpace:     int64(10000),
	RPC:                  map[string]string{"kovan": "https://kovan.infura.io/v3/45b81222fded4427b3a6589e0396c596"},
}

func TestModeOn() {
	testData.TestMode = true
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func TestModeOff() {
	testData.TestMode = false
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func TestConfig() nodeTypes.Config {
	return testConfig
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
