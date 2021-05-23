// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package store

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

// StoreABI is the input ABI used to generate the binding from.
const StoreABI = "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_baseDifficulty\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"wrong_hash\",\"type\":\"bytes32\"}],\"name\":\"wrongError\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"base_difficulty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_n\",\"type\":\"uint32\"}],\"name\":\"getBlockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockNumber\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_proof\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_targetDifficulty\",\"type\":\"uint256\"}],\"name\":\"isMatchDifficulty\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_root_hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"}],\"name\":\"isValidMerkleTreeProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isValidSign\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// Store is an auto generated Go binding around an Ethereum contract.
type Store struct {
	StoreCaller     // Read-only binding to the contract
	StoreTransactor // Write-only binding to the contract
	StoreFilterer   // Log filterer for contract events
}

// StoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type StoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StoreSession struct {
	Contract     *Store            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StoreCallerSession struct {
	Contract *StoreCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// StoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StoreTransactorSession struct {
	Contract     *StoreTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type StoreRaw struct {
	Contract *Store // Generic contract binding to access the raw methods on
}

// StoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StoreCallerRaw struct {
	Contract *StoreCaller // Generic read-only contract binding to access the raw methods on
}

// StoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StoreTransactorRaw struct {
	Contract *StoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStore creates a new instance of Store, bound to a specific deployed contract.
func NewStore(address common.Address, backend bind.ContractBackend) (*Store, error) {
	contract, err := bindStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Store{StoreCaller: StoreCaller{contract: contract}, StoreTransactor: StoreTransactor{contract: contract}, StoreFilterer: StoreFilterer{contract: contract}}, nil
}

// NewStoreCaller creates a new read-only instance of Store, bound to a specific deployed contract.
func NewStoreCaller(address common.Address, caller bind.ContractCaller) (*StoreCaller, error) {
	contract, err := bindStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StoreCaller{contract: contract}, nil
}

// NewStoreTransactor creates a new write-only instance of Store, bound to a specific deployed contract.
func NewStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*StoreTransactor, error) {
	contract, err := bindStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StoreTransactor{contract: contract}, nil
}

// NewStoreFilterer creates a new log filterer instance of Store, bound to a specific deployed contract.
func NewStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*StoreFilterer, error) {
	contract, err := bindStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StoreFilterer{contract: contract}, nil
}

// bindStore binds a generic wrapper to an already deployed contract.
func bindStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Store *StoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Store.Contract.StoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Store *StoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Store.Contract.StoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Store *StoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Store.Contract.StoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Store *StoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Store.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Store *StoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Store.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Store *StoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Store.Contract.contract.Transact(opts, method, params...)
}

// BaseDifficulty is a free data retrieval call binding the contract method 0x2d36de7a.
//
// Solidity: function base_difficulty() view returns(uint256)
func (_Store *StoreCaller) BaseDifficulty(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "base_difficulty")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseDifficulty is a free data retrieval call binding the contract method 0x2d36de7a.
//
// Solidity: function base_difficulty() view returns(uint256)
func (_Store *StoreSession) BaseDifficulty() (*big.Int, error) {
	return _Store.Contract.BaseDifficulty(&_Store.CallOpts)
}

// BaseDifficulty is a free data retrieval call binding the contract method 0x2d36de7a.
//
// Solidity: function base_difficulty() view returns(uint256)
func (_Store *StoreCallerSession) BaseDifficulty() (*big.Int, error) {
	return _Store.Contract.BaseDifficulty(&_Store.CallOpts)
}

// GetBlockHash is a free data retrieval call binding the contract method 0xd2b210a1.
//
// Solidity: function getBlockHash(uint32 _n) view returns(bytes32)
func (_Store *StoreCaller) GetBlockHash(opts *bind.CallOpts, _n uint32) ([32]byte, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "getBlockHash", _n)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBlockHash is a free data retrieval call binding the contract method 0xd2b210a1.
//
// Solidity: function getBlockHash(uint32 _n) view returns(bytes32)
func (_Store *StoreSession) GetBlockHash(_n uint32) ([32]byte, error) {
	return _Store.Contract.GetBlockHash(&_Store.CallOpts, _n)
}

// GetBlockHash is a free data retrieval call binding the contract method 0xd2b210a1.
//
// Solidity: function getBlockHash(uint32 _n) view returns(bytes32)
func (_Store *StoreCallerSession) GetBlockHash(_n uint32) ([32]byte, error) {
	return _Store.Contract.GetBlockHash(&_Store.CallOpts, _n)
}

// GetBlockNumber is a free data retrieval call binding the contract method 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint32)
func (_Store *StoreCaller) GetBlockNumber(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "getBlockNumber")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetBlockNumber is a free data retrieval call binding the contract method 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint32)
func (_Store *StoreSession) GetBlockNumber() (uint32, error) {
	return _Store.Contract.GetBlockNumber(&_Store.CallOpts)
}

// GetBlockNumber is a free data retrieval call binding the contract method 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint32)
func (_Store *StoreCallerSession) GetBlockNumber() (uint32, error) {
	return _Store.Contract.GetBlockNumber(&_Store.CallOpts)
}

// IsMatchDifficulty is a free data retrieval call binding the contract method 0x812adbd8.
//
// Solidity: function isMatchDifficulty(uint256 _proof, uint256 _targetDifficulty) view returns(bool)
func (_Store *StoreCaller) IsMatchDifficulty(opts *bind.CallOpts, _proof *big.Int, _targetDifficulty *big.Int) (bool, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "isMatchDifficulty", _proof, _targetDifficulty)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsMatchDifficulty is a free data retrieval call binding the contract method 0x812adbd8.
//
// Solidity: function isMatchDifficulty(uint256 _proof, uint256 _targetDifficulty) view returns(bool)
func (_Store *StoreSession) IsMatchDifficulty(_proof *big.Int, _targetDifficulty *big.Int) (bool, error) {
	return _Store.Contract.IsMatchDifficulty(&_Store.CallOpts, _proof, _targetDifficulty)
}

// IsMatchDifficulty is a free data retrieval call binding the contract method 0x812adbd8.
//
// Solidity: function isMatchDifficulty(uint256 _proof, uint256 _targetDifficulty) view returns(bool)
func (_Store *StoreCallerSession) IsMatchDifficulty(_proof *big.Int, _targetDifficulty *big.Int) (bool, error) {
	return _Store.Contract.IsMatchDifficulty(&_Store.CallOpts, _proof, _targetDifficulty)
}

// IsValidMerkleTreeProof is a free data retrieval call binding the contract method 0x059dbb30.
//
// Solidity: function isValidMerkleTreeProof(bytes32 _root_hash, bytes32[] proof) pure returns(bool)
func (_Store *StoreCaller) IsValidMerkleTreeProof(opts *bind.CallOpts, _root_hash [32]byte, proof [][32]byte) (bool, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "isValidMerkleTreeProof", _root_hash, proof)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidMerkleTreeProof is a free data retrieval call binding the contract method 0x059dbb30.
//
// Solidity: function isValidMerkleTreeProof(bytes32 _root_hash, bytes32[] proof) pure returns(bool)
func (_Store *StoreSession) IsValidMerkleTreeProof(_root_hash [32]byte, proof [][32]byte) (bool, error) {
	return _Store.Contract.IsValidMerkleTreeProof(&_Store.CallOpts, _root_hash, proof)
}

// IsValidMerkleTreeProof is a free data retrieval call binding the contract method 0x059dbb30.
//
// Solidity: function isValidMerkleTreeProof(bytes32 _root_hash, bytes32[] proof) pure returns(bool)
func (_Store *StoreCallerSession) IsValidMerkleTreeProof(_root_hash [32]byte, proof [][32]byte) (bool, error) {
	return _Store.Contract.IsValidMerkleTreeProof(&_Store.CallOpts, _root_hash, proof)
}

// IsValidSign is a free data retrieval call binding the contract method 0x16ecc0e6.
//
// Solidity: function isValidSign(address _signer, bytes message, bytes signature) pure returns(bool)
func (_Store *StoreCaller) IsValidSign(opts *bind.CallOpts, _signer common.Address, message []byte, signature []byte) (bool, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "isValidSign", _signer, message, signature)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidSign is a free data retrieval call binding the contract method 0x16ecc0e6.
//
// Solidity: function isValidSign(address _signer, bytes message, bytes signature) pure returns(bool)
func (_Store *StoreSession) IsValidSign(_signer common.Address, message []byte, signature []byte) (bool, error) {
	return _Store.Contract.IsValidSign(&_Store.CallOpts, _signer, message, signature)
}

// IsValidSign is a free data retrieval call binding the contract method 0x16ecc0e6.
//
// Solidity: function isValidSign(address _signer, bytes message, bytes signature) pure returns(bool)
func (_Store *StoreCallerSession) IsValidSign(_signer common.Address, message []byte, signature []byte) (bool, error) {
	return _Store.Contract.IsValidSign(&_Store.CallOpts, _signer, message, signature)
}

// StoreWrongErrorIterator is returned from FilterWrongError and is used to iterate over the raw logs and unpacked data for WrongError events raised by the Store contract.
type StoreWrongErrorIterator struct {
	Event *StoreWrongError // Event containing the contract specifics and raw log

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
func (it *StoreWrongErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StoreWrongError)
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
		it.Event = new(StoreWrongError)
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
func (it *StoreWrongErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StoreWrongErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StoreWrongError represents a WrongError event raised by the Store contract.
type StoreWrongError struct {
	WrongHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWrongError is a free log retrieval operation binding the contract event 0x93d7023aa9725055249a8ea4bce7cde913f73c0fa7730fb4a097de254fc815fc.
//
// Solidity: event wrongError(bytes32 wrong_hash)
func (_Store *StoreFilterer) FilterWrongError(opts *bind.FilterOpts) (*StoreWrongErrorIterator, error) {

	logs, sub, err := _Store.contract.FilterLogs(opts, "wrongError")
	if err != nil {
		return nil, err
	}
	return &StoreWrongErrorIterator{contract: _Store.contract, event: "wrongError", logs: logs, sub: sub}, nil
}

// WatchWrongError is a free log subscription operation binding the contract event 0x93d7023aa9725055249a8ea4bce7cde913f73c0fa7730fb4a097de254fc815fc.
//
// Solidity: event wrongError(bytes32 wrong_hash)
func (_Store *StoreFilterer) WatchWrongError(opts *bind.WatchOpts, sink chan<- *StoreWrongError) (event.Subscription, error) {

	logs, sub, err := _Store.contract.WatchLogs(opts, "wrongError")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StoreWrongError)
				if err := _Store.contract.UnpackLog(event, "wrongError", log); err != nil {
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

// ParseWrongError is a log parse operation binding the contract event 0x93d7023aa9725055249a8ea4bce7cde913f73c0fa7730fb4a097de254fc815fc.
//
// Solidity: event wrongError(bytes32 wrong_hash)
func (_Store *StoreFilterer) ParseWrongError(log types.Log) (*StoreWrongError, error) {
	event := new(StoreWrongError)
	if err := _Store.contract.UnpackLog(event, "wrongError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
