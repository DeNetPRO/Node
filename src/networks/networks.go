package networks

import (
	"github.com/DeNetPRO/src/errs"
	"github.com/DeNetPRO/src/logger"
	nodeTypes "github.com/DeNetPRO/src/node_types"
)

var currentNetwork string

var networks = map[string]nodeTypes.NtwrkParams{
	"polygon": {
		RPC:  "https://polygon-rpc.com",
		NODE: "0xfe1f5CB22cF4972584c6a0938FEAF90c597b567b",
		PoS:  "0x70c478be3d87ab921e0168137f5abe53b5812fc8",
		ERC:  "0xB27FAF7d98590Af6Ac38548edFBf05EEc0c18164",
		TRX:  "https://polygonscan.com/tx/",
	},
	"kovan": {
		RPC:  "https://kovan.infura.io/v3/45b81222fded4427b3a6589e0396c596",
		NODE: "0x805F977832959F39Cd7C2D9cDf5B30cE5A560d16",
		PoS:  "0x9F685bb981dD1b93668F0632eD99bAbf5171c747",
		ERC:  "0x98329d51486C0A942fCb3fAE5A0a18E05708cdc0",
		TRX:  "https://kovan.etherscan.io/tx/",
	},
	"mumbai": {
		RPC:  "https://rpc-mumbai.maticvigil.com",
		NODE: "0xBb86dcf291419d3F5b4B2211122D0E6fCB693777",
		PoS:  "0x389E8fE67c73551043184F740126C91866c0fB78",
		ERC:  "0xbAFBE687B0bD5D6fb7e87BB5Fc3E5f140394bC01",
		TRX:  "https://mumbai.polygonscan.com/tx/",
	},
}

func List() []string {
	nets := make([]string, 0, len(networks))

	for net := range networks {
		nets = append(nets, net)
	}

	return nets
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func Set(net string) error {
	const location = "networks.Set ->"

	_, supportedNet := networks[net]

	if !supportedNet {
		return logger.MarkLocation(location, errs.List().Network)
	}

	currentNetwork = net

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func Check(net string) error {
	const location = "networks.Check ->"

	_, supportedNet := networks[net]

	if !supportedNet {
		return logger.MarkLocation(location, errs.List().Network)
	}

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func Fields() nodeTypes.NtwrkParams {
	return networks[currentNetwork]
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func Current() string {
	return currentNetwork
}
