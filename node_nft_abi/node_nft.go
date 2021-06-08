// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package node_nft

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SimpleMetaDataDeNetNode is an auto generated low-level Go binding around an user-defined struct.
type SimpleMetaDataDeNetNode struct {
	IpAddress    [4]uint8
	Port         uint16
	BlockCreated *big.Int
	LastUpdate   *big.Int
	UpdatesCount *big.Int
}

// NodeNftABI is the input ABI used to generate the binding from.
const NodeNftABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8[4]\",\"name\":\"ip_address\",\"type\":\"uint8[4]\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"port\",\"type\":\"uint16\"}],\"name\":\"UpdateNodeStatus\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"NodeInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8[4]\",\"name\":\"ip_address\",\"type\":\"uint8[4]\"},{\"internalType\":\"uint16\",\"name\":\"port\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"block_created\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"last_update\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updates_count\",\"type\":\"uint256\"}],\"internalType\":\"structSimpleMetaData.DeNetNode\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8[4]\",\"name\":\"ip\",\"type\":\"uint8[4]\"},{\"internalType\":\"uint16\",\"name\":\"port\",\"type\":\"uint16\"}],\"name\":\"createNode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"node_id\",\"type\":\"uint256\"}],\"name\":\"getNodeById\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8[4]\",\"name\":\"ip_address\",\"type\":\"uint8[4]\"},{\"internalType\":\"uint16\",\"name\":\"port\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"block_created\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"last_update\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updates_count\",\"type\":\"uint256\"}],\"internalType\":\"structSimpleMetaData.DeNetNode\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"node_id\",\"type\":\"uint256\"},{\"internalType\":\"uint8[4]\",\"name\":\"ip\",\"type\":\"uint8[4]\"},{\"internalType\":\"uint16\",\"name\":\"port\",\"type\":\"uint16\"}],\"name\":\"updateNode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

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

// NodeInfo is a free data retrieval call binding the contract method 0x7aa210af.
//
// Solidity: function NodeInfo(uint256 tokenId) view returns((uint8[4],uint16,uint256,uint256,uint256))
func (_NodeNft *NodeNftCaller) NodeInfo(opts *bind.CallOpts, tokenId *big.Int) (SimpleMetaDataDeNetNode, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "NodeInfo", tokenId)

	if err != nil {
		return *new(SimpleMetaDataDeNetNode), err
	}

	out0 := *abi.ConvertType(out[0], new(SimpleMetaDataDeNetNode)).(*SimpleMetaDataDeNetNode)

	return out0, err

}

// NodeInfo is a free data retrieval call binding the contract method 0x7aa210af.
//
// Solidity: function NodeInfo(uint256 tokenId) view returns((uint8[4],uint16,uint256,uint256,uint256))
func (_NodeNft *NodeNftSession) NodeInfo(tokenId *big.Int) (SimpleMetaDataDeNetNode, error) {
	return _NodeNft.Contract.NodeInfo(&_NodeNft.CallOpts, tokenId)
}

// NodeInfo is a free data retrieval call binding the contract method 0x7aa210af.
//
// Solidity: function NodeInfo(uint256 tokenId) view returns((uint8[4],uint16,uint256,uint256,uint256))
func (_NodeNft *NodeNftCallerSession) NodeInfo(tokenId *big.Int) (SimpleMetaDataDeNetNode, error) {
	return _NodeNft.Contract.NodeInfo(&_NodeNft.CallOpts, tokenId)
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

// GetNodeById is a free data retrieval call binding the contract method 0x0e6a8496.
//
// Solidity: function getNodeById(uint256 node_id) view returns((uint8[4],uint16,uint256,uint256,uint256))
func (_NodeNft *NodeNftCaller) GetNodeById(opts *bind.CallOpts, node_id *big.Int) (SimpleMetaDataDeNetNode, error) {
	var out []interface{}
	err := _NodeNft.contract.Call(opts, &out, "getNodeById", node_id)

	if err != nil {
		return *new(SimpleMetaDataDeNetNode), err
	}

	out0 := *abi.ConvertType(out[0], new(SimpleMetaDataDeNetNode)).(*SimpleMetaDataDeNetNode)

	return out0, err

}

// GetNodeById is a free data retrieval call binding the contract method 0x0e6a8496.
//
// Solidity: function getNodeById(uint256 node_id) view returns((uint8[4],uint16,uint256,uint256,uint256))
func (_NodeNft *NodeNftSession) GetNodeById(node_id *big.Int) (SimpleMetaDataDeNetNode, error) {
	return _NodeNft.Contract.GetNodeById(&_NodeNft.CallOpts, node_id)
}

// GetNodeById is a free data retrieval call binding the contract method 0x0e6a8496.
//
// Solidity: function getNodeById(uint256 node_id) view returns((uint8[4],uint16,uint256,uint256,uint256))
func (_NodeNft *NodeNftCallerSession) GetNodeById(node_id *big.Int) (SimpleMetaDataDeNetNode, error) {
	return _NodeNft.Contract.GetNodeById(&_NodeNft.CallOpts, node_id)
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

// CreateNode is a paid mutator transaction binding the contract method 0x9e66580d.
//
// Solidity: function createNode(uint8[4] ip, uint16 port) returns()
func (_NodeNft *NodeNftTransactor) CreateNode(opts *bind.TransactOpts, ip [4]uint8, port uint16) (*types.Transaction, error) {
	return _NodeNft.contract.Transact(opts, "createNode", ip, port)
}

// CreateNode is a paid mutator transaction binding the contract method 0x9e66580d.
//
// Solidity: function createNode(uint8[4] ip, uint16 port) returns()
func (_NodeNft *NodeNftSession) CreateNode(ip [4]uint8, port uint16) (*types.Transaction, error) {
	return _NodeNft.Contract.CreateNode(&_NodeNft.TransactOpts, ip, port)
}

// CreateNode is a paid mutator transaction binding the contract method 0x9e66580d.
//
// Solidity: function createNode(uint8[4] ip, uint16 port) returns()
func (_NodeNft *NodeNftTransactorSession) CreateNode(ip [4]uint8, port uint16) (*types.Transaction, error) {
	return _NodeNft.Contract.CreateNode(&_NodeNft.TransactOpts, ip, port)
}

// UpdateNode is a paid mutator transaction binding the contract method 0x8978b8f7.
//
// Solidity: function updateNode(uint256 node_id, uint8[4] ip, uint16 port) returns()
func (_NodeNft *NodeNftTransactor) UpdateNode(opts *bind.TransactOpts, node_id *big.Int, ip [4]uint8, port uint16) (*types.Transaction, error) {
	return _NodeNft.contract.Transact(opts, "updateNode", node_id, ip, port)
}

// UpdateNode is a paid mutator transaction binding the contract method 0x8978b8f7.
//
// Solidity: function updateNode(uint256 node_id, uint8[4] ip, uint16 port) returns()
func (_NodeNft *NodeNftSession) UpdateNode(node_id *big.Int, ip [4]uint8, port uint16) (*types.Transaction, error) {
	return _NodeNft.Contract.UpdateNode(&_NodeNft.TransactOpts, node_id, ip, port)
}

// UpdateNode is a paid mutator transaction binding the contract method 0x8978b8f7.
//
// Solidity: function updateNode(uint256 node_id, uint8[4] ip, uint16 port) returns()
func (_NodeNft *NodeNftTransactorSession) UpdateNode(node_id *big.Int, ip [4]uint8, port uint16) (*types.Transaction, error) {
	return _NodeNft.Contract.UpdateNode(&_NodeNft.TransactOpts, node_id, ip, port)
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
// Solidity: event UpdateNodeStatus(address indexed from, uint256 indexed tokenId, uint8[4] ip_address, uint16 port)
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
// Solidity: event UpdateNodeStatus(address indexed from, uint256 indexed tokenId, uint8[4] ip_address, uint16 port)
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
// Solidity: event UpdateNodeStatus(address indexed from, uint256 indexed tokenId, uint8[4] ip_address, uint16 port)
func (_NodeNft *NodeNftFilterer) ParseUpdateNodeStatus(log types.Log) (*NodeNftUpdateNodeStatus, error) {
	event := new(NodeNftUpdateNodeStatus)
	if err := _NodeNft.contract.UnpackLog(event, "UpdateNodeStatus", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
