// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package node_nft

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// IMetaDataDeNetNode is an auto generated low-level Go binding around an user-defined struct.
type IMetaDataDeNetNode struct {
	IpAddress    [4]uint8
	Port         uint16
	CreatedAt    *big.Int
	UpdatedAt    *big.Int
	UpdatesCount *big.Int
	Rank         *big.Int
}

// NodeNftMetaData contains all meta data concerning the NodeNft contract.
var NodeNftMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_pos\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nodeLimit\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newPoSAddress\",\"type\":\"address\"}],\"name\":\"ChangePoSAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8[4]\",\"name\":\"ipAddress\",\"type\":\"uint8[4]\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"port\",\"type\":\"uint16\"}],\"name\":\"UpdateNodeStatus\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeOwner\",\"type\":\"address\"}],\"name\":\"addSuccessProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newAddress\",\"type\":\"address\"}],\"name\":\"changePoS\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8[4]\",\"name\":\"ip\",\"type\":\"uint8[4]\"},{\"internalType\":\"uint16\",\"name\":\"port\",\"type\":\"uint16\"}],\"name\":\"createNode\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_node\",\"type\":\"address\"}],\"name\":\"getNodeIDByAddress\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAlivePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxNodeID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextNodeID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nodeByAddress\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"nodeInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8[4]\",\"name\":\"ipAddress\",\"type\":\"uint8[4]\"},{\"internalType\":\"uint16\",\"name\":\"port\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatesCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rank\",\"type\":\"uint256\"}],\"internalType\":\"structIMetaData.DeNetNode\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nodesAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proofOfStorageAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_nodeID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"stealNode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nodeID\",\"type\":\"uint256\"},{\"internalType\":\"uint8[4]\",\"name\":\"ip\",\"type\":\"uint8[4]\"},{\"internalType\":\"uint16\",\"name\":\"port\",\"type\":\"uint16\"}],\"name\":\"updateNode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newLimit\",\"type\":\"uint256\"}],\"name\":\"updateNodesLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// NodeNftABI is the input ABI used to generate the binding from.
// Deprecated: Use NodeNftMetaData.ABI instead.
var NodeNftABI = NodeNftMetaData.ABI

// NodeNft is an auto generated Go binding around an Ethereum contract.
type NodeNft struct {
	NodeNftCaller     // Read-only binding to the contract
	NodeNftTransactor // Write-only binding to the contract
	NodeNftFilterer   // Log filterer for contract events
}

// NodeNftCaller is an auto generated read-only Go binding around an Ethereum contract.
type NodeNftCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeNftTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NodeNftTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeNftFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NodeNftFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeNftSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NodeNftSession struct {
	Contract     *NodeNft          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodeNftCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NodeNftCallerSession struct {
	Contract *NodeNftCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// NodeNftTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NodeNftTransactorSession struct {
	Contract     *NodeNftTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// NodeNftRaw is an auto generated low-level Go binding around an Ethereum contract.
type NodeNftRaw struct {
	Contract *NodeNft // Generic contract binding to access the raw methods on
}

// NodeNftCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NodeNftCallerRaw struct {
	Contract *NodeNftCaller // Generic read-only contract binding to access the raw methods on
}

// NodeNftTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NodeNftTransactorRaw struct {
	Contract *NodeNftTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNodeNft creates a new instance of NodeNft, bound to a specific deployed contract.
func NewNodeNft(address common.Address, backend bind.ContractBackend) (*NodeNft, error) {
	contract, err := bindNodeNft(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NodeNft{NodeNftCaller: NodeNftCaller{contract: contract}, NodeNftTransactor: NodeNftTransactor{contract: contract}, NodeNftFilterer: NodeNftFilterer{contract: contract}}, nil
}

// NewNodeNftCaller creates a new read-only instance of NodeNft, bound to a specific deployed contract.
func NewNodeNftCaller(address common.Address, caller bind.ContractCaller) (*NodeNftCaller, error) {
	contract, err := bindNodeNft(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NodeNftCaller{contract: contract}, nil
}

// NewNodeNftTransactor creates a new write-only instance of NodeNft, bound to a specific deployed contract.
func NewNodeNftTransactor(address common.Address, transactor bind.ContractTransactor) (*NodeNftTransactor, error) {
	contract, err := bindNodeNft(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NodeNftTransactor{contract: contract}, nil
}

// NewNodeNftFilterer creates a new log filterer instance of NodeNft, bound to a specific deployed contract.
func NewNodeNftFilterer(address common.Address, filterer bind.ContractFilterer) (*NodeNftFilterer, error) {
	contract, err := bindNodeNft(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NodeNftFilterer{contract: contract}, nil
}

// bindNodeNft binds a generic wrapper to an already deployed contract.
func bindNodeNft(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NodeNftABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeNft *NodeNftRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeNft.Contract.NodeNftCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeNft *NodeNftRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeNft.Contract.NodeNftTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeNft *NodeNftRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeNft.Contract.NodeNftTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeNft *NodeNftCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeNft.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeNft *NodeNftTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeNft.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeNft *NodeNftTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeNft.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NodeNft *NodeNftCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NodeNft *NodeNftSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _NodeNft.Contract.BalanceOf(&_NodeNft.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NodeNft *NodeNftCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _NodeNft.Contract.BalanceOf(&_NodeNft.CallOpts, owner)
}

// GetNodeIDByAddress is a free data retrieval call binding the contract method 0x482f6ce2.
//
// Solidity: function getNodeIDByAddress(address _node) view returns(uint256)
func (_NodeNft *NodeNftCaller) GetNodeIDByAddress(opts *bind.CallOpts, _node common.Address) (*big.Int, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "getNodeIDByAddress", _node)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeIDByAddress is a free data retrieval call binding the contract method 0x482f6ce2.
//
// Solidity: function getNodeIDByAddress(address _node) view returns(uint256)
func (_NodeNft *NodeNftSession) GetNodeIDByAddress(_node common.Address) (*big.Int, error) {
	return _NodeNft.Contract.GetNodeIDByAddress(&_NodeNft.CallOpts, _node)
}

// GetNodeIDByAddress is a free data retrieval call binding the contract method 0x482f6ce2.
//
// Solidity: function getNodeIDByAddress(address _node) view returns(uint256)
func (_NodeNft *NodeNftCallerSession) GetNodeIDByAddress(_node common.Address) (*big.Int, error) {
	return _NodeNft.Contract.GetNodeIDByAddress(&_NodeNft.CallOpts, _node)
}

// MaxAlivePeriod is a free data retrieval call binding the contract method 0xc60f86a4.
//
// Solidity: function maxAlivePeriod() view returns(uint256)
func (_NodeNft *NodeNftCaller) MaxAlivePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "maxAlivePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxAlivePeriod is a free data retrieval call binding the contract method 0xc60f86a4.
//
// Solidity: function maxAlivePeriod() view returns(uint256)
func (_NodeNft *NodeNftSession) MaxAlivePeriod() (*big.Int, error) {
	return _NodeNft.Contract.MaxAlivePeriod(&_NodeNft.CallOpts)
}

// MaxAlivePeriod is a free data retrieval call binding the contract method 0xc60f86a4.
//
// Solidity: function maxAlivePeriod() view returns(uint256)
func (_NodeNft *NodeNftCallerSession) MaxAlivePeriod() (*big.Int, error) {
	return _NodeNft.Contract.MaxAlivePeriod(&_NodeNft.CallOpts)
}

// MaxNodeID is a free data retrieval call binding the contract method 0xbcc7ee4a.
//
// Solidity: function maxNodeID() view returns(uint256)
func (_NodeNft *NodeNftCaller) MaxNodeID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "maxNodeID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxNodeID is a free data retrieval call binding the contract method 0xbcc7ee4a.
//
// Solidity: function maxNodeID() view returns(uint256)
func (_NodeNft *NodeNftSession) MaxNodeID() (*big.Int, error) {
	return _NodeNft.Contract.MaxNodeID(&_NodeNft.CallOpts)
}

// MaxNodeID is a free data retrieval call binding the contract method 0xbcc7ee4a.
//
// Solidity: function maxNodeID() view returns(uint256)
func (_NodeNft *NodeNftCallerSession) MaxNodeID() (*big.Int, error) {
	return _NodeNft.Contract.MaxNodeID(&_NodeNft.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NodeNft *NodeNftCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NodeNft *NodeNftSession) Name() (string, error) {
	return _NodeNft.Contract.Name(&_NodeNft.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NodeNft *NodeNftCallerSession) Name() (string, error) {
	return _NodeNft.Contract.Name(&_NodeNft.CallOpts)
}

// NextNodeID is a free data retrieval call binding the contract method 0x49c10e26.
//
// Solidity: function nextNodeID() view returns(uint256)
func (_NodeNft *NodeNftCaller) NextNodeID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "nextNodeID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextNodeID is a free data retrieval call binding the contract method 0x49c10e26.
//
// Solidity: function nextNodeID() view returns(uint256)
func (_NodeNft *NodeNftSession) NextNodeID() (*big.Int, error) {
	return _NodeNft.Contract.NextNodeID(&_NodeNft.CallOpts)
}

// NextNodeID is a free data retrieval call binding the contract method 0x49c10e26.
//
// Solidity: function nextNodeID() view returns(uint256)
func (_NodeNft *NodeNftCallerSession) NextNodeID() (*big.Int, error) {
	return _NodeNft.Contract.NextNodeID(&_NodeNft.CallOpts)
}

// NodeByAddress is a free data retrieval call binding the contract method 0x3fd23654.
//
// Solidity: function nodeByAddress(address ) view returns(uint256)
func (_NodeNft *NodeNftCaller) NodeByAddress(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "nodeByAddress", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NodeByAddress is a free data retrieval call binding the contract method 0x3fd23654.
//
// Solidity: function nodeByAddress(address ) view returns(uint256)
func (_NodeNft *NodeNftSession) NodeByAddress(arg0 common.Address) (*big.Int, error) {
	return _NodeNft.Contract.NodeByAddress(&_NodeNft.CallOpts, arg0)
}

// NodeByAddress is a free data retrieval call binding the contract method 0x3fd23654.
//
// Solidity: function nodeByAddress(address ) view returns(uint256)
func (_NodeNft *NodeNftCallerSession) NodeByAddress(arg0 common.Address) (*big.Int, error) {
	return _NodeNft.Contract.NodeByAddress(&_NodeNft.CallOpts, arg0)
}

// NodeInfo is a free data retrieval call binding the contract method 0xb02439ae.
//
// Solidity: function nodeInfo(uint256 tokenId) view returns((uint8[4],uint16,uint256,uint256,uint256,uint256))
func (_NodeNft *NodeNftCaller) NodeInfo(opts *bind.CallOpts, tokenId *big.Int) (IMetaDataDeNetNode, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "nodeInfo", tokenId)

	if err != nil {
		return *new(IMetaDataDeNetNode), err
	}

	out0 := *abi.ConvertType(out[0], new(IMetaDataDeNetNode)).(*IMetaDataDeNetNode)

	return out0, err

}

// NodeInfo is a free data retrieval call binding the contract method 0xb02439ae.
//
// Solidity: function nodeInfo(uint256 tokenId) view returns((uint8[4],uint16,uint256,uint256,uint256,uint256))
func (_NodeNft *NodeNftSession) NodeInfo(tokenId *big.Int) (IMetaDataDeNetNode, error) {
	return _NodeNft.Contract.NodeInfo(&_NodeNft.CallOpts, tokenId)
}

// NodeInfo is a free data retrieval call binding the contract method 0xb02439ae.
//
// Solidity: function nodeInfo(uint256 tokenId) view returns((uint8[4],uint16,uint256,uint256,uint256,uint256))
func (_NodeNft *NodeNftCallerSession) NodeInfo(tokenId *big.Int) (IMetaDataDeNetNode, error) {
	return _NodeNft.Contract.NodeInfo(&_NodeNft.CallOpts, tokenId)
}

// NodesAvailable is a free data retrieval call binding the contract method 0x0431a525.
//
// Solidity: function nodesAvailable() view returns(uint256)
func (_NodeNft *NodeNftCaller) NodesAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "nodesAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NodesAvailable is a free data retrieval call binding the contract method 0x0431a525.
//
// Solidity: function nodesAvailable() view returns(uint256)
func (_NodeNft *NodeNftSession) NodesAvailable() (*big.Int, error) {
	return _NodeNft.Contract.NodesAvailable(&_NodeNft.CallOpts)
}

// NodesAvailable is a free data retrieval call binding the contract method 0x0431a525.
//
// Solidity: function nodesAvailable() view returns(uint256)
func (_NodeNft *NodeNftCallerSession) NodesAvailable() (*big.Int, error) {
	return _NodeNft.Contract.NodesAvailable(&_NodeNft.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NodeNft *NodeNftCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NodeNft *NodeNftSession) Owner() (common.Address, error) {
	return _NodeNft.Contract.Owner(&_NodeNft.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NodeNft *NodeNftCallerSession) Owner() (common.Address, error) {
	return _NodeNft.Contract.Owner(&_NodeNft.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NodeNft *NodeNftCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NodeNft *NodeNftSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _NodeNft.Contract.OwnerOf(&_NodeNft.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NodeNft *NodeNftCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _NodeNft.Contract.OwnerOf(&_NodeNft.CallOpts, tokenId)
}

// ProofOfStorageAddress is a free data retrieval call binding the contract method 0x4ef98de7.
//
// Solidity: function proofOfStorageAddress() view returns(address)
func (_NodeNft *NodeNftCaller) ProofOfStorageAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "proofOfStorageAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProofOfStorageAddress is a free data retrieval call binding the contract method 0x4ef98de7.
//
// Solidity: function proofOfStorageAddress() view returns(address)
func (_NodeNft *NodeNftSession) ProofOfStorageAddress() (common.Address, error) {
	return _NodeNft.Contract.ProofOfStorageAddress(&_NodeNft.CallOpts)
}

// ProofOfStorageAddress is a free data retrieval call binding the contract method 0x4ef98de7.
//
// Solidity: function proofOfStorageAddress() view returns(address)
func (_NodeNft *NodeNftCallerSession) ProofOfStorageAddress() (common.Address, error) {
	return _NodeNft.Contract.ProofOfStorageAddress(&_NodeNft.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NodeNft *NodeNftCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NodeNft *NodeNftSession) Symbol() (string, error) {
	return _NodeNft.Contract.Symbol(&_NodeNft.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NodeNft *NodeNftCallerSession) Symbol() (string, error) {
	return _NodeNft.Contract.Symbol(&_NodeNft.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_NodeNft *NodeNftCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_NodeNft *NodeNftSession) TotalSupply() (*big.Int, error) {
	return _NodeNft.Contract.TotalSupply(&_NodeNft.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_NodeNft *NodeNftCallerSession) TotalSupply() (*big.Int, error) {
	return _NodeNft.Contract.TotalSupply(&_NodeNft.CallOpts)
}

// AddSuccessProof is a paid mutator transaction binding the contract method 0x96b9445c.
//
// Solidity: function addSuccessProof(address _nodeOwner) returns()
func (_NodeNft *NodeNftTransactor) AddSuccessProof(opts *bind.TransactOpts, _nodeOwner common.Address) (*types.Transaction, error) {
	return _NodeNft.contract.Transact(opts, "addSuccessProof", _nodeOwner)
}

// AddSuccessProof is a paid mutator transaction binding the contract method 0x96b9445c.
//
// Solidity: function addSuccessProof(address _nodeOwner) returns()
func (_NodeNft *NodeNftSession) AddSuccessProof(_nodeOwner common.Address) (*types.Transaction, error) {
	return _NodeNft.Contract.AddSuccessProof(&_NodeNft.TransactOpts, _nodeOwner)
}

// AddSuccessProof is a paid mutator transaction binding the contract method 0x96b9445c.
//
// Solidity: function addSuccessProof(address _nodeOwner) returns()
func (_NodeNft *NodeNftTransactorSession) AddSuccessProof(_nodeOwner common.Address) (*types.Transaction, error) {
	return _NodeNft.Contract.AddSuccessProof(&_NodeNft.TransactOpts, _nodeOwner)
}

// ChangePoS is a paid mutator transaction binding the contract method 0x98b16026.
//
// Solidity: function changePoS(address _newAddress) returns()
func (_NodeNft *NodeNftTransactor) ChangePoS(opts *bind.TransactOpts, _newAddress common.Address) (*types.Transaction, error) {
	return _NodeNft.contract.Transact(opts, "changePoS", _newAddress)
}

// ChangePoS is a paid mutator transaction binding the contract method 0x98b16026.
//
// Solidity: function changePoS(address _newAddress) returns()
func (_NodeNft *NodeNftSession) ChangePoS(_newAddress common.Address) (*types.Transaction, error) {
	return _NodeNft.Contract.ChangePoS(&_NodeNft.TransactOpts, _newAddress)
}

// ChangePoS is a paid mutator transaction binding the contract method 0x98b16026.
//
// Solidity: function changePoS(address _newAddress) returns()
func (_NodeNft *NodeNftTransactorSession) ChangePoS(_newAddress common.Address) (*types.Transaction, error) {
	return _NodeNft.Contract.ChangePoS(&_NodeNft.TransactOpts, _newAddress)
}

// CreateNode is a paid mutator transaction binding the contract method 0x9e66580d.
//
// Solidity: function createNode(uint8[4] ip, uint16 port) returns(uint256)
func (_NodeNft *NodeNftTransactor) CreateNode(opts *bind.TransactOpts, ip [4]uint8, port uint16) (*types.Transaction, error) {
	return _NodeNft.contract.Transact(opts, "createNode", ip, port)
}

// CreateNode is a paid mutator transaction binding the contract method 0x9e66580d.
//
// Solidity: function createNode(uint8[4] ip, uint16 port) returns(uint256)
func (_NodeNft *NodeNftSession) CreateNode(ip [4]uint8, port uint16) (*types.Transaction, error) {
	return _NodeNft.Contract.CreateNode(&_NodeNft.TransactOpts, ip, port)
}

// CreateNode is a paid mutator transaction binding the contract method 0x9e66580d.
//
// Solidity: function createNode(uint8[4] ip, uint16 port) returns(uint256)
func (_NodeNft *NodeNftTransactorSession) CreateNode(ip [4]uint8, port uint16) (*types.Transaction, error) {
	return _NodeNft.Contract.CreateNode(&_NodeNft.TransactOpts, ip, port)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NodeNft *NodeNftTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeNft.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NodeNft *NodeNftSession) RenounceOwnership() (*types.Transaction, error) {
	return _NodeNft.Contract.RenounceOwnership(&_NodeNft.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NodeNft *NodeNftTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _NodeNft.Contract.RenounceOwnership(&_NodeNft.TransactOpts)
}

// StealNode is a paid mutator transaction binding the contract method 0x517c0d57.
//
// Solidity: function stealNode(uint256 _nodeID, address _to) returns()
func (_NodeNft *NodeNftTransactor) StealNode(opts *bind.TransactOpts, _nodeID *big.Int, _to common.Address) (*types.Transaction, error) {
	return _NodeNft.contract.Transact(opts, "stealNode", _nodeID, _to)
}

// StealNode is a paid mutator transaction binding the contract method 0x517c0d57.
//
// Solidity: function stealNode(uint256 _nodeID, address _to) returns()
func (_NodeNft *NodeNftSession) StealNode(_nodeID *big.Int, _to common.Address) (*types.Transaction, error) {
	return _NodeNft.Contract.StealNode(&_NodeNft.TransactOpts, _nodeID, _to)
}

// StealNode is a paid mutator transaction binding the contract method 0x517c0d57.
//
// Solidity: function stealNode(uint256 _nodeID, address _to) returns()
func (_NodeNft *NodeNftTransactorSession) StealNode(_nodeID *big.Int, _to common.Address) (*types.Transaction, error) {
	return _NodeNft.Contract.StealNode(&_NodeNft.TransactOpts, _nodeID, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NodeNft *NodeNftTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _NodeNft.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NodeNft *NodeNftSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NodeNft.Contract.TransferOwnership(&_NodeNft.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NodeNft *NodeNftTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NodeNft.Contract.TransferOwnership(&_NodeNft.TransactOpts, newOwner)
}

// UpdateNode is a paid mutator transaction binding the contract method 0x8978b8f7.
//
// Solidity: function updateNode(uint256 nodeID, uint8[4] ip, uint16 port) returns()
func (_NodeNft *NodeNftTransactor) UpdateNode(opts *bind.TransactOpts, nodeID *big.Int, ip [4]uint8, port uint16) (*types.Transaction, error) {
	return _NodeNft.contract.Transact(opts, "updateNode", nodeID, ip, port)
}

// UpdateNode is a paid mutator transaction binding the contract method 0x8978b8f7.
//
// Solidity: function updateNode(uint256 nodeID, uint8[4] ip, uint16 port) returns()
func (_NodeNft *NodeNftSession) UpdateNode(nodeID *big.Int, ip [4]uint8, port uint16) (*types.Transaction, error) {
	return _NodeNft.Contract.UpdateNode(&_NodeNft.TransactOpts, nodeID, ip, port)
}

// UpdateNode is a paid mutator transaction binding the contract method 0x8978b8f7.
//
// Solidity: function updateNode(uint256 nodeID, uint8[4] ip, uint16 port) returns()
func (_NodeNft *NodeNftTransactorSession) UpdateNode(nodeID *big.Int, ip [4]uint8, port uint16) (*types.Transaction, error) {
	return _NodeNft.Contract.UpdateNode(&_NodeNft.TransactOpts, nodeID, ip, port)
}

// UpdateNodesLimit is a paid mutator transaction binding the contract method 0x6a0b80ac.
//
// Solidity: function updateNodesLimit(uint256 _newLimit) returns()
func (_NodeNft *NodeNftTransactor) UpdateNodesLimit(opts *bind.TransactOpts, _newLimit *big.Int) (*types.Transaction, error) {
	return _NodeNft.contract.Transact(opts, "updateNodesLimit", _newLimit)
}

// UpdateNodesLimit is a paid mutator transaction binding the contract method 0x6a0b80ac.
//
// Solidity: function updateNodesLimit(uint256 _newLimit) returns()
func (_NodeNft *NodeNftSession) UpdateNodesLimit(_newLimit *big.Int) (*types.Transaction, error) {
	return _NodeNft.Contract.UpdateNodesLimit(&_NodeNft.TransactOpts, _newLimit)
}

// UpdateNodesLimit is a paid mutator transaction binding the contract method 0x6a0b80ac.
//
// Solidity: function updateNodesLimit(uint256 _newLimit) returns()
func (_NodeNft *NodeNftTransactorSession) UpdateNodesLimit(_newLimit *big.Int) (*types.Transaction, error) {
	return _NodeNft.Contract.UpdateNodesLimit(&_NodeNft.TransactOpts, _newLimit)
}

// NodeNftChangePoSAddressIterator is returned from FilterChangePoSAddress and is used to iterate over the raw logs and unpacked data for ChangePoSAddress events raised by the NodeNft contract.
type NodeNftChangePoSAddressIterator struct {
	Event *NodeNftChangePoSAddress // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeNftChangePoSAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeNftChangePoSAddress)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeNftChangePoSAddress)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeNftChangePoSAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeNftChangePoSAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeNftChangePoSAddress represents a ChangePoSAddress event raised by the NodeNft contract.
type NodeNftChangePoSAddress struct {
	NewPoSAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterChangePoSAddress is a free log retrieval operation binding the contract event 0x7c5f6c5c3ca6ba61d3c0a70572313385c14a071e872bc4451980ddedbc9e8840.
//
// Solidity: event ChangePoSAddress(address indexed newPoSAddress)
func (_NodeNft *NodeNftFilterer) FilterChangePoSAddress(opts *bind.FilterOpts, newPoSAddress []common.Address) (*NodeNftChangePoSAddressIterator, error) {

	var newPoSAddressRule []interface{}
	for _, newPoSAddressItem := range newPoSAddress {
		newPoSAddressRule = append(newPoSAddressRule, newPoSAddressItem)
	}

	logs, sub, err := _NodeNft.contract.FilterLogs(opts, "ChangePoSAddress", newPoSAddressRule)
	if err != nil {
		return nil, err
	}
	return &NodeNftChangePoSAddressIterator{contract: _NodeNft.contract, event: "ChangePoSAddress", logs: logs, sub: sub}, nil
}

// WatchChangePoSAddress is a free log subscription operation binding the contract event 0x7c5f6c5c3ca6ba61d3c0a70572313385c14a071e872bc4451980ddedbc9e8840.
//
// Solidity: event ChangePoSAddress(address indexed newPoSAddress)
func (_NodeNft *NodeNftFilterer) WatchChangePoSAddress(opts *bind.WatchOpts, sink chan<- *NodeNftChangePoSAddress, newPoSAddress []common.Address) (event.Subscription, error) {

	var newPoSAddressRule []interface{}
	for _, newPoSAddressItem := range newPoSAddress {
		newPoSAddressRule = append(newPoSAddressRule, newPoSAddressItem)
	}

	logs, sub, err := _NodeNft.contract.WatchLogs(opts, "ChangePoSAddress", newPoSAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeNftChangePoSAddress)
				if err := _NodeNft.contract.UnpackLog(event, "ChangePoSAddress", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChangePoSAddress is a log parse operation binding the contract event 0x7c5f6c5c3ca6ba61d3c0a70572313385c14a071e872bc4451980ddedbc9e8840.
//
// Solidity: event ChangePoSAddress(address indexed newPoSAddress)
func (_NodeNft *NodeNftFilterer) ParseChangePoSAddress(log types.Log) (*NodeNftChangePoSAddress, error) {
	event := new(NodeNftChangePoSAddress)
	if err := _NodeNft.contract.UnpackLog(event, "ChangePoSAddress", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeNftOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the NodeNft contract.
type NodeNftOwnershipTransferredIterator struct {
	Event *NodeNftOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeNftOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeNftOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeNftOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeNftOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeNftOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeNftOwnershipTransferred represents a OwnershipTransferred event raised by the NodeNft contract.
type NodeNftOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NodeNft *NodeNftFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*NodeNftOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NodeNft.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NodeNftOwnershipTransferredIterator{contract: _NodeNft.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NodeNft *NodeNftFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NodeNftOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NodeNft.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeNftOwnershipTransferred)
				if err := _NodeNft.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NodeNft *NodeNftFilterer) ParseOwnershipTransferred(log types.Log) (*NodeNftOwnershipTransferred, error) {
	event := new(NodeNftOwnershipTransferred)
	if err := _NodeNft.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeNftTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the NodeNft contract.
type NodeNftTransferIterator struct {
	Event *NodeNftTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeNftTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeNftTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeNftTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeNftTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeNftTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeNftTransfer represents a Transfer event raised by the NodeNft contract.
type NodeNftTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_NodeNft *NodeNftFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*NodeNftTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _NodeNft.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeNftTransferIterator{contract: _NodeNft.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_NodeNft *NodeNftFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *NodeNftTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _NodeNft.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeNftTransfer)
				if err := _NodeNft.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_NodeNft *NodeNftFilterer) ParseTransfer(log types.Log) (*NodeNftTransfer, error) {
	event := new(NodeNftTransfer)
	if err := _NodeNft.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeNftUpdateNodeStatusIterator is returned from FilterUpdateNodeStatus and is used to iterate over the raw logs and unpacked data for UpdateNodeStatus events raised by the NodeNft contract.
type NodeNftUpdateNodeStatusIterator struct {
	Event *NodeNftUpdateNodeStatus // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeNftUpdateNodeStatusIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeNftUpdateNodeStatus)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeNftUpdateNodeStatus)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeNftUpdateNodeStatusIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeNftUpdateNodeStatusIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeNftUpdateNodeStatus represents a UpdateNodeStatus event raised by the NodeNft contract.
type NodeNftUpdateNodeStatus struct {
	From      common.Address
	TokenId   *big.Int
	IpAddress [4]uint8
	Port      uint16
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUpdateNodeStatus is a free log retrieval operation binding the contract event 0x8d7f699bed84a209ce7d1b706d08428cc21bcc5fccc7a632fea7bea99e9bba68.
//
// Solidity: event UpdateNodeStatus(address indexed from, uint256 indexed tokenId, uint8[4] ipAddress, uint16 port)
func (_NodeNft *NodeNftFilterer) FilterUpdateNodeStatus(opts *bind.FilterOpts, from []common.Address, tokenId []*big.Int) (*NodeNftUpdateNodeStatusIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _NodeNft.contract.FilterLogs(opts, "UpdateNodeStatus", fromRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeNftUpdateNodeStatusIterator{contract: _NodeNft.contract, event: "UpdateNodeStatus", logs: logs, sub: sub}, nil
}

// WatchUpdateNodeStatus is a free log subscription operation binding the contract event 0x8d7f699bed84a209ce7d1b706d08428cc21bcc5fccc7a632fea7bea99e9bba68.
//
// Solidity: event UpdateNodeStatus(address indexed from, uint256 indexed tokenId, uint8[4] ipAddress, uint16 port)
func (_NodeNft *NodeNftFilterer) WatchUpdateNodeStatus(opts *bind.WatchOpts, sink chan<- *NodeNftUpdateNodeStatus, from []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _NodeNft.contract.WatchLogs(opts, "UpdateNodeStatus", fromRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeNftUpdateNodeStatus)
				if err := _NodeNft.contract.UnpackLog(event, "UpdateNodeStatus", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdateNodeStatus is a log parse operation binding the contract event 0x8d7f699bed84a209ce7d1b706d08428cc21bcc5fccc7a632fea7bea99e9bba68.
//
// Solidity: event UpdateNodeStatus(address indexed from, uint256 indexed tokenId, uint8[4] ipAddress, uint16 port)
func (_NodeNft *NodeNftFilterer) ParseUpdateNodeStatus(log types.Log) (*NodeNftUpdateNodeStatus, error) {
	event := new(NodeNftUpdateNodeStatus)
	if err := _NodeNft.contract.UnpackLog(event, "UpdateNodeStatus", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
