// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package pos

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

// PosMetaData contains all meta data concerning the Pos contract.
var PosMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_storage_address\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_payments\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_baseDifficulty\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"wrong_hash\",\"type\":\"bytes32\"}],\"name\":\"WrongError\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"base_difficulty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_storage_address\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_payments_address\",\"type\":\"address\"}],\"name\":\"changeSystemAddresses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"closeDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"debug_mode\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"_curDate\",\"type\":\"uint32\"}],\"name\":\"getAvailableDeposit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_n\",\"type\":\"uint32\"}],\"name\":\"getBlockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockNumber\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDifficulty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_file\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_block_number\",\"type\":\"uint256\"}],\"name\":\"getProof\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_user_storage_size\",\"type\":\"uint256\"}],\"name\":\"getUserRewardInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"getUserRootHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"invisibleMintGasToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_proof\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_targetDifficulty\",\"type\":\"uint256\"}],\"name\":\"isMatchDifficulty\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_root_hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"}],\"name\":\"isValidMerkleTreeProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"name\":\"limitReached\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"makeDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxDepositPerUser\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"min_storage_require\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"node_nft_address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paymentsAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user_address\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_block_number\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"_user_root_hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"_user_storage_size\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_user_root_hash_nonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_user_signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_file\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"merkleProof\",\"type\":\"bytes32[]\"}],\"name\":\"sendProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_node_address\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_user_address\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_block_number\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"_user_root_hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"_user_storage_size\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_user_root_hash_nonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_user_signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_file\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"merkleProof\",\"type\":\"bytes32[]\"}],\"name\":\"sendProofFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_size\",\"type\":\"uint256\"}],\"name\":\"setMinStorage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_new\",\"type\":\"address\"}],\"name\":\"setNodeNFTAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"turnDebugMode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_new_difficulty\",\"type\":\"uint256\"}],\"name\":\"updateBaseDifficulty\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"user_storage_address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_file\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"_block_number\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"_time_passed\",\"type\":\"uint256\"}],\"name\":\"verifyFileProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// PosABI is the input ABI used to generate the binding from.
// Deprecated: Use PosMetaData.ABI instead.
var PosABI = PosMetaData.ABI

// Pos is an auto generated Go binding around an Ethereum contract.
type Pos struct {
	PosCaller     // Read-only binding to the contract
	PosTransactor // Write-only binding to the contract
	PosFilterer   // Log filterer for contract events
}

// PosCaller is an auto generated read-only Go binding around an Ethereum contract.
type PosCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PosTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PosTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PosFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PosFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PosSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PosSession struct {
	Contract     *Pos              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PosCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PosCallerSession struct {
	Contract *PosCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// PosTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PosTransactorSession struct {
	Contract     *PosTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PosRaw is an auto generated low-level Go binding around an Ethereum contract.
type PosRaw struct {
	Contract *Pos // Generic contract binding to access the raw methods on
}

// PosCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PosCallerRaw struct {
	Contract *PosCaller // Generic read-only contract binding to access the raw methods on
}

// PosTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PosTransactorRaw struct {
	Contract *PosTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPos creates a new instance of Pos, bound to a specific deployed contract.
func NewPos(address common.Address, backend bind.ContractBackend) (*Pos, error) {
	contract, err := bindPos(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Pos{PosCaller: PosCaller{contract: contract}, PosTransactor: PosTransactor{contract: contract}, PosFilterer: PosFilterer{contract: contract}}, nil
}

// NewPosCaller creates a new read-only instance of Pos, bound to a specific deployed contract.
func NewPosCaller(address common.Address, caller bind.ContractCaller) (*PosCaller, error) {
	contract, err := bindPos(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PosCaller{contract: contract}, nil
}

// NewPosTransactor creates a new write-only instance of Pos, bound to a specific deployed contract.
func NewPosTransactor(address common.Address, transactor bind.ContractTransactor) (*PosTransactor, error) {
	contract, err := bindPos(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PosTransactor{contract: contract}, nil
}

// NewPosFilterer creates a new log filterer instance of Pos, bound to a specific deployed contract.
func NewPosFilterer(address common.Address, filterer bind.ContractFilterer) (*PosFilterer, error) {
	contract, err := bindPos(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PosFilterer{contract: contract}, nil
}

// bindPos binds a generic wrapper to an already deployed contract.
func bindPos(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PosABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pos *PosRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Pos.Contract.PosCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pos *PosRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pos.Contract.PosTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pos *PosRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pos.Contract.PosTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pos *PosCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Pos.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pos *PosTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pos.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pos *PosTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pos.Contract.contract.Transact(opts, method, params...)
}

// BaseDifficulty is a free data retrieval call binding the contract method 0x2d36de7a.
//
// Solidity: function base_difficulty() view returns(uint256)
func (_Pos *PosCaller) BaseDifficulty(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "base_difficulty")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseDifficulty is a free data retrieval call binding the contract method 0x2d36de7a.
//
// Solidity: function base_difficulty() view returns(uint256)
func (_Pos *PosSession) BaseDifficulty() (*big.Int, error) {
	return _Pos.Contract.BaseDifficulty(&_Pos.CallOpts)
}

// BaseDifficulty is a free data retrieval call binding the contract method 0x2d36de7a.
//
// Solidity: function base_difficulty() view returns(uint256)
func (_Pos *PosCallerSession) BaseDifficulty() (*big.Int, error) {
	return _Pos.Contract.BaseDifficulty(&_Pos.CallOpts)
}

// DebugMode is a free data retrieval call binding the contract method 0x4631dd94.
//
// Solidity: function debug_mode() view returns(bool)
func (_Pos *PosCaller) DebugMode(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "debug_mode")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// DebugMode is a free data retrieval call binding the contract method 0x4631dd94.
//
// Solidity: function debug_mode() view returns(bool)
func (_Pos *PosSession) DebugMode() (bool, error) {
	return _Pos.Contract.DebugMode(&_Pos.CallOpts)
}

// DebugMode is a free data retrieval call binding the contract method 0x4631dd94.
//
// Solidity: function debug_mode() view returns(bool)
func (_Pos *PosCallerSession) DebugMode() (bool, error) {
	return _Pos.Contract.DebugMode(&_Pos.CallOpts)
}

// GetAvailableDeposit is a free data retrieval call binding the contract method 0x2355d34c.
//
// Solidity: function getAvailableDeposit(address _user, uint256 _amount, uint32 _curDate) view returns(uint256)
func (_Pos *PosCaller) GetAvailableDeposit(opts *bind.CallOpts, _user common.Address, _amount *big.Int, _curDate uint32) (*big.Int, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "getAvailableDeposit", _user, _amount, _curDate)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAvailableDeposit is a free data retrieval call binding the contract method 0x2355d34c.
//
// Solidity: function getAvailableDeposit(address _user, uint256 _amount, uint32 _curDate) view returns(uint256)
func (_Pos *PosSession) GetAvailableDeposit(_user common.Address, _amount *big.Int, _curDate uint32) (*big.Int, error) {
	return _Pos.Contract.GetAvailableDeposit(&_Pos.CallOpts, _user, _amount, _curDate)
}

// GetAvailableDeposit is a free data retrieval call binding the contract method 0x2355d34c.
//
// Solidity: function getAvailableDeposit(address _user, uint256 _amount, uint32 _curDate) view returns(uint256)
func (_Pos *PosCallerSession) GetAvailableDeposit(_user common.Address, _amount *big.Int, _curDate uint32) (*big.Int, error) {
	return _Pos.Contract.GetAvailableDeposit(&_Pos.CallOpts, _user, _amount, _curDate)
}

// GetBlockHash is a free data retrieval call binding the contract method 0xd2b210a1.
//
// Solidity: function getBlockHash(uint32 _n) view returns(bytes32)
func (_Pos *PosCaller) GetBlockHash(opts *bind.CallOpts, _n uint32) ([32]byte, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "getBlockHash", _n)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBlockHash is a free data retrieval call binding the contract method 0xd2b210a1.
//
// Solidity: function getBlockHash(uint32 _n) view returns(bytes32)
func (_Pos *PosSession) GetBlockHash(_n uint32) ([32]byte, error) {
	return _Pos.Contract.GetBlockHash(&_Pos.CallOpts, _n)
}

// GetBlockHash is a free data retrieval call binding the contract method 0xd2b210a1.
//
// Solidity: function getBlockHash(uint32 _n) view returns(bytes32)
func (_Pos *PosCallerSession) GetBlockHash(_n uint32) ([32]byte, error) {
	return _Pos.Contract.GetBlockHash(&_Pos.CallOpts, _n)
}

// GetBlockNumber is a free data retrieval call binding the contract method 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint32)
func (_Pos *PosCaller) GetBlockNumber(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "getBlockNumber")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetBlockNumber is a free data retrieval call binding the contract method 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint32)
func (_Pos *PosSession) GetBlockNumber() (uint32, error) {
	return _Pos.Contract.GetBlockNumber(&_Pos.CallOpts)
}

// GetBlockNumber is a free data retrieval call binding the contract method 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint32)
func (_Pos *PosCallerSession) GetBlockNumber() (uint32, error) {
	return _Pos.Contract.GetBlockNumber(&_Pos.CallOpts)
}

// GetDifficulty is a free data retrieval call binding the contract method 0xb6baffe3.
//
// Solidity: function getDifficulty() view returns(uint256)
func (_Pos *PosCaller) GetDifficulty(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "getDifficulty")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDifficulty is a free data retrieval call binding the contract method 0xb6baffe3.
//
// Solidity: function getDifficulty() view returns(uint256)
func (_Pos *PosSession) GetDifficulty() (*big.Int, error) {
	return _Pos.Contract.GetDifficulty(&_Pos.CallOpts)
}

// GetDifficulty is a free data retrieval call binding the contract method 0xb6baffe3.
//
// Solidity: function getDifficulty() view returns(uint256)
func (_Pos *PosCallerSession) GetDifficulty() (*big.Int, error) {
	return _Pos.Contract.GetDifficulty(&_Pos.CallOpts)
}

// GetProof is a free data retrieval call binding the contract method 0x3f6df031.
//
// Solidity: function getProof(bytes _file, address _sender, uint256 _block_number) view returns(bytes, bytes32)
func (_Pos *PosCaller) GetProof(opts *bind.CallOpts, _file []byte, _sender common.Address, _block_number *big.Int) ([]byte, [32]byte, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "getProof", _file, _sender, _block_number)

	if err != nil {
		return *new([]byte), *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	out1 := *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return out0, out1, err

}

// GetProof is a free data retrieval call binding the contract method 0x3f6df031.
//
// Solidity: function getProof(bytes _file, address _sender, uint256 _block_number) view returns(bytes, bytes32)
func (_Pos *PosSession) GetProof(_file []byte, _sender common.Address, _block_number *big.Int) ([]byte, [32]byte, error) {
	return _Pos.Contract.GetProof(&_Pos.CallOpts, _file, _sender, _block_number)
}

// GetProof is a free data retrieval call binding the contract method 0x3f6df031.
//
// Solidity: function getProof(bytes _file, address _sender, uint256 _block_number) view returns(bytes, bytes32)
func (_Pos *PosCallerSession) GetProof(_file []byte, _sender common.Address, _block_number *big.Int) ([]byte, [32]byte, error) {
	return _Pos.Contract.GetProof(&_Pos.CallOpts, _file, _sender, _block_number)
}

// GetUserRewardInfo is a free data retrieval call binding the contract method 0x1ce41cb3.
//
// Solidity: function getUserRewardInfo(address _user, uint256 _user_storage_size) view returns(uint256, uint256)
func (_Pos *PosCaller) GetUserRewardInfo(opts *bind.CallOpts, _user common.Address, _user_storage_size *big.Int) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "getUserRewardInfo", _user, _user_storage_size)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetUserRewardInfo is a free data retrieval call binding the contract method 0x1ce41cb3.
//
// Solidity: function getUserRewardInfo(address _user, uint256 _user_storage_size) view returns(uint256, uint256)
func (_Pos *PosSession) GetUserRewardInfo(_user common.Address, _user_storage_size *big.Int) (*big.Int, *big.Int, error) {
	return _Pos.Contract.GetUserRewardInfo(&_Pos.CallOpts, _user, _user_storage_size)
}

// GetUserRewardInfo is a free data retrieval call binding the contract method 0x1ce41cb3.
//
// Solidity: function getUserRewardInfo(address _user, uint256 _user_storage_size) view returns(uint256, uint256)
func (_Pos *PosCallerSession) GetUserRewardInfo(_user common.Address, _user_storage_size *big.Int) (*big.Int, *big.Int, error) {
	return _Pos.Contract.GetUserRewardInfo(&_Pos.CallOpts, _user, _user_storage_size)
}

// GetUserRootHash is a free data retrieval call binding the contract method 0xf9a76fde.
//
// Solidity: function getUserRootHash(address _user) view returns(bytes32, uint256)
func (_Pos *PosCaller) GetUserRootHash(opts *bind.CallOpts, _user common.Address) ([32]byte, *big.Int, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "getUserRootHash", _user)

	if err != nil {
		return *new([32]byte), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetUserRootHash is a free data retrieval call binding the contract method 0xf9a76fde.
//
// Solidity: function getUserRootHash(address _user) view returns(bytes32, uint256)
func (_Pos *PosSession) GetUserRootHash(_user common.Address) ([32]byte, *big.Int, error) {
	return _Pos.Contract.GetUserRootHash(&_Pos.CallOpts, _user)
}

// GetUserRootHash is a free data retrieval call binding the contract method 0xf9a76fde.
//
// Solidity: function getUserRootHash(address _user) view returns(bytes32, uint256)
func (_Pos *PosCallerSession) GetUserRootHash(_user common.Address) ([32]byte, *big.Int, error) {
	return _Pos.Contract.GetUserRootHash(&_Pos.CallOpts, _user)
}

// IsMatchDifficulty is a free data retrieval call binding the contract method 0x812adbd8.
//
// Solidity: function isMatchDifficulty(uint256 _proof, uint256 _targetDifficulty) view returns(bool)
func (_Pos *PosCaller) IsMatchDifficulty(opts *bind.CallOpts, _proof *big.Int, _targetDifficulty *big.Int) (bool, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "isMatchDifficulty", _proof, _targetDifficulty)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsMatchDifficulty is a free data retrieval call binding the contract method 0x812adbd8.
//
// Solidity: function isMatchDifficulty(uint256 _proof, uint256 _targetDifficulty) view returns(bool)
func (_Pos *PosSession) IsMatchDifficulty(_proof *big.Int, _targetDifficulty *big.Int) (bool, error) {
	return _Pos.Contract.IsMatchDifficulty(&_Pos.CallOpts, _proof, _targetDifficulty)
}

// IsMatchDifficulty is a free data retrieval call binding the contract method 0x812adbd8.
//
// Solidity: function isMatchDifficulty(uint256 _proof, uint256 _targetDifficulty) view returns(bool)
func (_Pos *PosCallerSession) IsMatchDifficulty(_proof *big.Int, _targetDifficulty *big.Int) (bool, error) {
	return _Pos.Contract.IsMatchDifficulty(&_Pos.CallOpts, _proof, _targetDifficulty)
}

// IsValidMerkleTreeProof is a free data retrieval call binding the contract method 0x059dbb30.
//
// Solidity: function isValidMerkleTreeProof(bytes32 _root_hash, bytes32[] proof) pure returns(bool)
func (_Pos *PosCaller) IsValidMerkleTreeProof(opts *bind.CallOpts, _root_hash [32]byte, proof [][32]byte) (bool, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "isValidMerkleTreeProof", _root_hash, proof)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidMerkleTreeProof is a free data retrieval call binding the contract method 0x059dbb30.
//
// Solidity: function isValidMerkleTreeProof(bytes32 _root_hash, bytes32[] proof) pure returns(bool)
func (_Pos *PosSession) IsValidMerkleTreeProof(_root_hash [32]byte, proof [][32]byte) (bool, error) {
	return _Pos.Contract.IsValidMerkleTreeProof(&_Pos.CallOpts, _root_hash, proof)
}

// IsValidMerkleTreeProof is a free data retrieval call binding the contract method 0x059dbb30.
//
// Solidity: function isValidMerkleTreeProof(bytes32 _root_hash, bytes32[] proof) pure returns(bool)
func (_Pos *PosCallerSession) IsValidMerkleTreeProof(_root_hash [32]byte, proof [][32]byte) (bool, error) {
	return _Pos.Contract.IsValidMerkleTreeProof(&_Pos.CallOpts, _root_hash, proof)
}

// LimitReached is a free data retrieval call binding the contract method 0x509395c5.
//
// Solidity: function limitReached(address , uint32 ) view returns(uint256)
func (_Pos *PosCaller) LimitReached(opts *bind.CallOpts, arg0 common.Address, arg1 uint32) (*big.Int, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "limitReached", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LimitReached is a free data retrieval call binding the contract method 0x509395c5.
//
// Solidity: function limitReached(address , uint32 ) view returns(uint256)
func (_Pos *PosSession) LimitReached(arg0 common.Address, arg1 uint32) (*big.Int, error) {
	return _Pos.Contract.LimitReached(&_Pos.CallOpts, arg0, arg1)
}

// LimitReached is a free data retrieval call binding the contract method 0x509395c5.
//
// Solidity: function limitReached(address , uint32 ) view returns(uint256)
func (_Pos *PosCallerSession) LimitReached(arg0 common.Address, arg1 uint32) (*big.Int, error) {
	return _Pos.Contract.LimitReached(&_Pos.CallOpts, arg0, arg1)
}

// MaxDepositPerUser is a free data retrieval call binding the contract method 0x82e5310c.
//
// Solidity: function maxDepositPerUser() view returns(uint256)
func (_Pos *PosCaller) MaxDepositPerUser(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "maxDepositPerUser")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxDepositPerUser is a free data retrieval call binding the contract method 0x82e5310c.
//
// Solidity: function maxDepositPerUser() view returns(uint256)
func (_Pos *PosSession) MaxDepositPerUser() (*big.Int, error) {
	return _Pos.Contract.MaxDepositPerUser(&_Pos.CallOpts)
}

// MaxDepositPerUser is a free data retrieval call binding the contract method 0x82e5310c.
//
// Solidity: function maxDepositPerUser() view returns(uint256)
func (_Pos *PosCallerSession) MaxDepositPerUser() (*big.Int, error) {
	return _Pos.Contract.MaxDepositPerUser(&_Pos.CallOpts)
}

// MinStorageRequire is a free data retrieval call binding the contract method 0x173a8fbc.
//
// Solidity: function min_storage_require() view returns(uint256)
func (_Pos *PosCaller) MinStorageRequire(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "min_storage_require")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinStorageRequire is a free data retrieval call binding the contract method 0x173a8fbc.
//
// Solidity: function min_storage_require() view returns(uint256)
func (_Pos *PosSession) MinStorageRequire() (*big.Int, error) {
	return _Pos.Contract.MinStorageRequire(&_Pos.CallOpts)
}

// MinStorageRequire is a free data retrieval call binding the contract method 0x173a8fbc.
//
// Solidity: function min_storage_require() view returns(uint256)
func (_Pos *PosCallerSession) MinStorageRequire() (*big.Int, error) {
	return _Pos.Contract.MinStorageRequire(&_Pos.CallOpts)
}

// NodeNftAddress is a free data retrieval call binding the contract method 0x81cfd709.
//
// Solidity: function node_nft_address() view returns(address)
func (_Pos *PosCaller) NodeNftAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "node_nft_address")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NodeNftAddress is a free data retrieval call binding the contract method 0x81cfd709.
//
// Solidity: function node_nft_address() view returns(address)
func (_Pos *PosSession) NodeNftAddress() (common.Address, error) {
	return _Pos.Contract.NodeNftAddress(&_Pos.CallOpts)
}

// NodeNftAddress is a free data retrieval call binding the contract method 0x81cfd709.
//
// Solidity: function node_nft_address() view returns(address)
func (_Pos *PosCallerSession) NodeNftAddress() (common.Address, error) {
	return _Pos.Contract.NodeNftAddress(&_Pos.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Pos *PosCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Pos *PosSession) Owner() (common.Address, error) {
	return _Pos.Contract.Owner(&_Pos.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Pos *PosCallerSession) Owner() (common.Address, error) {
	return _Pos.Contract.Owner(&_Pos.CallOpts)
}

// PaymentsAddress is a free data retrieval call binding the contract method 0x694d0473.
//
// Solidity: function paymentsAddress() view returns(address)
func (_Pos *PosCaller) PaymentsAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "paymentsAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PaymentsAddress is a free data retrieval call binding the contract method 0x694d0473.
//
// Solidity: function paymentsAddress() view returns(address)
func (_Pos *PosSession) PaymentsAddress() (common.Address, error) {
	return _Pos.Contract.PaymentsAddress(&_Pos.CallOpts)
}

// PaymentsAddress is a free data retrieval call binding the contract method 0x694d0473.
//
// Solidity: function paymentsAddress() view returns(address)
func (_Pos *PosCallerSession) PaymentsAddress() (common.Address, error) {
	return _Pos.Contract.PaymentsAddress(&_Pos.CallOpts)
}

// TimeLimit is a free data retrieval call binding the contract method 0xc08d1fe5.
//
// Solidity: function timeLimit() view returns(uint256)
func (_Pos *PosCaller) TimeLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "timeLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TimeLimit is a free data retrieval call binding the contract method 0xc08d1fe5.
//
// Solidity: function timeLimit() view returns(uint256)
func (_Pos *PosSession) TimeLimit() (*big.Int, error) {
	return _Pos.Contract.TimeLimit(&_Pos.CallOpts)
}

// TimeLimit is a free data retrieval call binding the contract method 0xc08d1fe5.
//
// Solidity: function timeLimit() view returns(uint256)
func (_Pos *PosCallerSession) TimeLimit() (*big.Int, error) {
	return _Pos.Contract.TimeLimit(&_Pos.CallOpts)
}

// UserStorageAddress is a free data retrieval call binding the contract method 0x1079a326.
//
// Solidity: function user_storage_address() view returns(address)
func (_Pos *PosCaller) UserStorageAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "user_storage_address")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UserStorageAddress is a free data retrieval call binding the contract method 0x1079a326.
//
// Solidity: function user_storage_address() view returns(address)
func (_Pos *PosSession) UserStorageAddress() (common.Address, error) {
	return _Pos.Contract.UserStorageAddress(&_Pos.CallOpts)
}

// UserStorageAddress is a free data retrieval call binding the contract method 0x1079a326.
//
// Solidity: function user_storage_address() view returns(address)
func (_Pos *PosCallerSession) UserStorageAddress() (common.Address, error) {
	return _Pos.Contract.UserStorageAddress(&_Pos.CallOpts)
}

// VerifyFileProof is a free data retrieval call binding the contract method 0x1d285f0a.
//
// Solidity: function verifyFileProof(address _sender, bytes _file, uint32 _block_number, uint256 _time_passed) view returns(bool)
func (_Pos *PosCaller) VerifyFileProof(opts *bind.CallOpts, _sender common.Address, _file []byte, _block_number uint32, _time_passed *big.Int) (bool, error) {
	var out []interface{}
	err := _Pos.contract.Call(opts, &out, "verifyFileProof", _sender, _file, _block_number, _time_passed)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyFileProof is a free data retrieval call binding the contract method 0x1d285f0a.
//
// Solidity: function verifyFileProof(address _sender, bytes _file, uint32 _block_number, uint256 _time_passed) view returns(bool)
func (_Pos *PosSession) VerifyFileProof(_sender common.Address, _file []byte, _block_number uint32, _time_passed *big.Int) (bool, error) {
	return _Pos.Contract.VerifyFileProof(&_Pos.CallOpts, _sender, _file, _block_number, _time_passed)
}

// VerifyFileProof is a free data retrieval call binding the contract method 0x1d285f0a.
//
// Solidity: function verifyFileProof(address _sender, bytes _file, uint32 _block_number, uint256 _time_passed) view returns(bool)
func (_Pos *PosCallerSession) VerifyFileProof(_sender common.Address, _file []byte, _block_number uint32, _time_passed *big.Int) (bool, error) {
	return _Pos.Contract.VerifyFileProof(&_Pos.CallOpts, _sender, _file, _block_number, _time_passed)
}

// ChangeSystemAddresses is a paid mutator transaction binding the contract method 0xeb5439dc.
//
// Solidity: function changeSystemAddresses(address _storage_address, address _payments_address) returns()
func (_Pos *PosTransactor) ChangeSystemAddresses(opts *bind.TransactOpts, _storage_address common.Address, _payments_address common.Address) (*types.Transaction, error) {
	return _Pos.contract.Transact(opts, "changeSystemAddresses", _storage_address, _payments_address)
}

// ChangeSystemAddresses is a paid mutator transaction binding the contract method 0xeb5439dc.
//
// Solidity: function changeSystemAddresses(address _storage_address, address _payments_address) returns()
func (_Pos *PosSession) ChangeSystemAddresses(_storage_address common.Address, _payments_address common.Address) (*types.Transaction, error) {
	return _Pos.Contract.ChangeSystemAddresses(&_Pos.TransactOpts, _storage_address, _payments_address)
}

// ChangeSystemAddresses is a paid mutator transaction binding the contract method 0xeb5439dc.
//
// Solidity: function changeSystemAddresses(address _storage_address, address _payments_address) returns()
func (_Pos *PosTransactorSession) ChangeSystemAddresses(_storage_address common.Address, _payments_address common.Address) (*types.Transaction, error) {
	return _Pos.Contract.ChangeSystemAddresses(&_Pos.TransactOpts, _storage_address, _payments_address)
}

// CloseDeposit is a paid mutator transaction binding the contract method 0xdc8d0a40.
//
// Solidity: function closeDeposit(address _token) returns()
func (_Pos *PosTransactor) CloseDeposit(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _Pos.contract.Transact(opts, "closeDeposit", _token)
}

// CloseDeposit is a paid mutator transaction binding the contract method 0xdc8d0a40.
//
// Solidity: function closeDeposit(address _token) returns()
func (_Pos *PosSession) CloseDeposit(_token common.Address) (*types.Transaction, error) {
	return _Pos.Contract.CloseDeposit(&_Pos.TransactOpts, _token)
}

// CloseDeposit is a paid mutator transaction binding the contract method 0xdc8d0a40.
//
// Solidity: function closeDeposit(address _token) returns()
func (_Pos *PosTransactorSession) CloseDeposit(_token common.Address) (*types.Transaction, error) {
	return _Pos.Contract.CloseDeposit(&_Pos.TransactOpts, _token)
}

// InvisibleMintGasToken is a paid mutator transaction binding the contract method 0x77ac86a1.
//
// Solidity: function invisibleMintGasToken(address _from, address _user, uint256 _amount) returns()
func (_Pos *PosTransactor) InvisibleMintGasToken(opts *bind.TransactOpts, _from common.Address, _user common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Pos.contract.Transact(opts, "invisibleMintGasToken", _from, _user, _amount)
}

// InvisibleMintGasToken is a paid mutator transaction binding the contract method 0x77ac86a1.
//
// Solidity: function invisibleMintGasToken(address _from, address _user, uint256 _amount) returns()
func (_Pos *PosSession) InvisibleMintGasToken(_from common.Address, _user common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Pos.Contract.InvisibleMintGasToken(&_Pos.TransactOpts, _from, _user, _amount)
}

// InvisibleMintGasToken is a paid mutator transaction binding the contract method 0x77ac86a1.
//
// Solidity: function invisibleMintGasToken(address _from, address _user, uint256 _amount) returns()
func (_Pos *PosTransactorSession) InvisibleMintGasToken(_from common.Address, _user common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Pos.Contract.InvisibleMintGasToken(&_Pos.TransactOpts, _from, _user, _amount)
}

// MakeDeposit is a paid mutator transaction binding the contract method 0xbf5d0a00.
//
// Solidity: function makeDeposit(address _token, uint256 _amount) returns()
func (_Pos *PosTransactor) MakeDeposit(opts *bind.TransactOpts, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Pos.contract.Transact(opts, "makeDeposit", _token, _amount)
}

// MakeDeposit is a paid mutator transaction binding the contract method 0xbf5d0a00.
//
// Solidity: function makeDeposit(address _token, uint256 _amount) returns()
func (_Pos *PosSession) MakeDeposit(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Pos.Contract.MakeDeposit(&_Pos.TransactOpts, _token, _amount)
}

// MakeDeposit is a paid mutator transaction binding the contract method 0xbf5d0a00.
//
// Solidity: function makeDeposit(address _token, uint256 _amount) returns()
func (_Pos *PosTransactorSession) MakeDeposit(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Pos.Contract.MakeDeposit(&_Pos.TransactOpts, _token, _amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Pos *PosTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pos.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Pos *PosSession) RenounceOwnership() (*types.Transaction, error) {
	return _Pos.Contract.RenounceOwnership(&_Pos.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Pos *PosTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Pos.Contract.RenounceOwnership(&_Pos.TransactOpts)
}

// SendProof is a paid mutator transaction binding the contract method 0x888c2a85.
//
// Solidity: function sendProof(address _user_address, uint32 _block_number, bytes32 _user_root_hash, uint64 _user_storage_size, uint64 _user_root_hash_nonce, bytes _user_signature, bytes _file, bytes32[] merkleProof) returns()
func (_Pos *PosTransactor) SendProof(opts *bind.TransactOpts, _user_address common.Address, _block_number uint32, _user_root_hash [32]byte, _user_storage_size uint64, _user_root_hash_nonce uint64, _user_signature []byte, _file []byte, merkleProof [][32]byte) (*types.Transaction, error) {
	return _Pos.contract.Transact(opts, "sendProof", _user_address, _block_number, _user_root_hash, _user_storage_size, _user_root_hash_nonce, _user_signature, _file, merkleProof)
}

// SendProof is a paid mutator transaction binding the contract method 0x888c2a85.
//
// Solidity: function sendProof(address _user_address, uint32 _block_number, bytes32 _user_root_hash, uint64 _user_storage_size, uint64 _user_root_hash_nonce, bytes _user_signature, bytes _file, bytes32[] merkleProof) returns()
func (_Pos *PosSession) SendProof(_user_address common.Address, _block_number uint32, _user_root_hash [32]byte, _user_storage_size uint64, _user_root_hash_nonce uint64, _user_signature []byte, _file []byte, merkleProof [][32]byte) (*types.Transaction, error) {
	return _Pos.Contract.SendProof(&_Pos.TransactOpts, _user_address, _block_number, _user_root_hash, _user_storage_size, _user_root_hash_nonce, _user_signature, _file, merkleProof)
}

// SendProof is a paid mutator transaction binding the contract method 0x888c2a85.
//
// Solidity: function sendProof(address _user_address, uint32 _block_number, bytes32 _user_root_hash, uint64 _user_storage_size, uint64 _user_root_hash_nonce, bytes _user_signature, bytes _file, bytes32[] merkleProof) returns()
func (_Pos *PosTransactorSession) SendProof(_user_address common.Address, _block_number uint32, _user_root_hash [32]byte, _user_storage_size uint64, _user_root_hash_nonce uint64, _user_signature []byte, _file []byte, merkleProof [][32]byte) (*types.Transaction, error) {
	return _Pos.Contract.SendProof(&_Pos.TransactOpts, _user_address, _block_number, _user_root_hash, _user_storage_size, _user_root_hash_nonce, _user_signature, _file, merkleProof)
}

// SendProofFrom is a paid mutator transaction binding the contract method 0x7cd5e5a9.
//
// Solidity: function sendProofFrom(address _node_address, address _user_address, uint32 _block_number, bytes32 _user_root_hash, uint64 _user_storage_size, uint64 _user_root_hash_nonce, bytes _user_signature, bytes _file, bytes32[] merkleProof) returns()
func (_Pos *PosTransactor) SendProofFrom(opts *bind.TransactOpts, _node_address common.Address, _user_address common.Address, _block_number uint32, _user_root_hash [32]byte, _user_storage_size uint64, _user_root_hash_nonce uint64, _user_signature []byte, _file []byte, merkleProof [][32]byte) (*types.Transaction, error) {
	return _Pos.contract.Transact(opts, "sendProofFrom", _node_address, _user_address, _block_number, _user_root_hash, _user_storage_size, _user_root_hash_nonce, _user_signature, _file, merkleProof)
}

// SendProofFrom is a paid mutator transaction binding the contract method 0x7cd5e5a9.
//
// Solidity: function sendProofFrom(address _node_address, address _user_address, uint32 _block_number, bytes32 _user_root_hash, uint64 _user_storage_size, uint64 _user_root_hash_nonce, bytes _user_signature, bytes _file, bytes32[] merkleProof) returns()
func (_Pos *PosSession) SendProofFrom(_node_address common.Address, _user_address common.Address, _block_number uint32, _user_root_hash [32]byte, _user_storage_size uint64, _user_root_hash_nonce uint64, _user_signature []byte, _file []byte, merkleProof [][32]byte) (*types.Transaction, error) {
	return _Pos.Contract.SendProofFrom(&_Pos.TransactOpts, _node_address, _user_address, _block_number, _user_root_hash, _user_storage_size, _user_root_hash_nonce, _user_signature, _file, merkleProof)
}

// SendProofFrom is a paid mutator transaction binding the contract method 0x7cd5e5a9.
//
// Solidity: function sendProofFrom(address _node_address, address _user_address, uint32 _block_number, bytes32 _user_root_hash, uint64 _user_storage_size, uint64 _user_root_hash_nonce, bytes _user_signature, bytes _file, bytes32[] merkleProof) returns()
func (_Pos *PosTransactorSession) SendProofFrom(_node_address common.Address, _user_address common.Address, _block_number uint32, _user_root_hash [32]byte, _user_storage_size uint64, _user_root_hash_nonce uint64, _user_signature []byte, _file []byte, merkleProof [][32]byte) (*types.Transaction, error) {
	return _Pos.Contract.SendProofFrom(&_Pos.TransactOpts, _node_address, _user_address, _block_number, _user_root_hash, _user_storage_size, _user_root_hash_nonce, _user_signature, _file, merkleProof)
}

// SetMinStorage is a paid mutator transaction binding the contract method 0x785d6917.
//
// Solidity: function setMinStorage(uint256 _size) returns()
func (_Pos *PosTransactor) SetMinStorage(opts *bind.TransactOpts, _size *big.Int) (*types.Transaction, error) {
	return _Pos.contract.Transact(opts, "setMinStorage", _size)
}

// SetMinStorage is a paid mutator transaction binding the contract method 0x785d6917.
//
// Solidity: function setMinStorage(uint256 _size) returns()
func (_Pos *PosSession) SetMinStorage(_size *big.Int) (*types.Transaction, error) {
	return _Pos.Contract.SetMinStorage(&_Pos.TransactOpts, _size)
}

// SetMinStorage is a paid mutator transaction binding the contract method 0x785d6917.
//
// Solidity: function setMinStorage(uint256 _size) returns()
func (_Pos *PosTransactorSession) SetMinStorage(_size *big.Int) (*types.Transaction, error) {
	return _Pos.Contract.SetMinStorage(&_Pos.TransactOpts, _size)
}

// SetNodeNFTAddress is a paid mutator transaction binding the contract method 0xca811507.
//
// Solidity: function setNodeNFTAddress(address _new) returns()
func (_Pos *PosTransactor) SetNodeNFTAddress(opts *bind.TransactOpts, _new common.Address) (*types.Transaction, error) {
	return _Pos.contract.Transact(opts, "setNodeNFTAddress", _new)
}

// SetNodeNFTAddress is a paid mutator transaction binding the contract method 0xca811507.
//
// Solidity: function setNodeNFTAddress(address _new) returns()
func (_Pos *PosSession) SetNodeNFTAddress(_new common.Address) (*types.Transaction, error) {
	return _Pos.Contract.SetNodeNFTAddress(&_Pos.TransactOpts, _new)
}

// SetNodeNFTAddress is a paid mutator transaction binding the contract method 0xca811507.
//
// Solidity: function setNodeNFTAddress(address _new) returns()
func (_Pos *PosTransactorSession) SetNodeNFTAddress(_new common.Address) (*types.Transaction, error) {
	return _Pos.Contract.SetNodeNFTAddress(&_Pos.TransactOpts, _new)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Pos *PosTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Pos.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Pos *PosSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Pos.Contract.TransferOwnership(&_Pos.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Pos *PosTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Pos.Contract.TransferOwnership(&_Pos.TransactOpts, newOwner)
}

// TurnDebugMode is a paid mutator transaction binding the contract method 0xfc57bfcd.
//
// Solidity: function turnDebugMode() returns()
func (_Pos *PosTransactor) TurnDebugMode(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pos.contract.Transact(opts, "turnDebugMode")
}

// TurnDebugMode is a paid mutator transaction binding the contract method 0xfc57bfcd.
//
// Solidity: function turnDebugMode() returns()
func (_Pos *PosSession) TurnDebugMode() (*types.Transaction, error) {
	return _Pos.Contract.TurnDebugMode(&_Pos.TransactOpts)
}

// TurnDebugMode is a paid mutator transaction binding the contract method 0xfc57bfcd.
//
// Solidity: function turnDebugMode() returns()
func (_Pos *PosTransactorSession) TurnDebugMode() (*types.Transaction, error) {
	return _Pos.Contract.TurnDebugMode(&_Pos.TransactOpts)
}

// UpdateBaseDifficulty is a paid mutator transaction binding the contract method 0xb5c6da27.
//
// Solidity: function updateBaseDifficulty(uint256 _new_difficulty) returns()
func (_Pos *PosTransactor) UpdateBaseDifficulty(opts *bind.TransactOpts, _new_difficulty *big.Int) (*types.Transaction, error) {
	return _Pos.contract.Transact(opts, "updateBaseDifficulty", _new_difficulty)
}

// UpdateBaseDifficulty is a paid mutator transaction binding the contract method 0xb5c6da27.
//
// Solidity: function updateBaseDifficulty(uint256 _new_difficulty) returns()
func (_Pos *PosSession) UpdateBaseDifficulty(_new_difficulty *big.Int) (*types.Transaction, error) {
	return _Pos.Contract.UpdateBaseDifficulty(&_Pos.TransactOpts, _new_difficulty)
}

// UpdateBaseDifficulty is a paid mutator transaction binding the contract method 0xb5c6da27.
//
// Solidity: function updateBaseDifficulty(uint256 _new_difficulty) returns()
func (_Pos *PosTransactorSession) UpdateBaseDifficulty(_new_difficulty *big.Int) (*types.Transaction, error) {
	return _Pos.Contract.UpdateBaseDifficulty(&_Pos.TransactOpts, _new_difficulty)
}

// PosOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Pos contract.
type PosOwnershipTransferredIterator struct {
	Event *PosOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PosOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PosOwnershipTransferred)
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
		it.Event = new(PosOwnershipTransferred)
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
func (it *PosOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PosOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PosOwnershipTransferred represents a OwnershipTransferred event raised by the Pos contract.
type PosOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Pos *PosFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PosOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Pos.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PosOwnershipTransferredIterator{contract: _Pos.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Pos *PosFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PosOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Pos.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PosOwnershipTransferred)
				if err := _Pos.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Pos *PosFilterer) ParseOwnershipTransferred(log types.Log) (*PosOwnershipTransferred, error) {
	event := new(PosOwnershipTransferred)
	if err := _Pos.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PosWrongErrorIterator is returned from FilterWrongError and is used to iterate over the raw logs and unpacked data for WrongError events raised by the Pos contract.
type PosWrongErrorIterator struct {
	Event *PosWrongError // Event containing the contract specifics and raw log

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
func (it *PosWrongErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PosWrongError)
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
		it.Event = new(PosWrongError)
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
func (it *PosWrongErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PosWrongErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PosWrongError represents a WrongError event raised by the Pos contract.
type PosWrongError struct {
	WrongHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWrongError is a free log retrieval operation binding the contract event 0xba80455d9a4978272e947d7083631bd5cc9203cc9526eea9b5f32f21c1e90c57.
//
// Solidity: event WrongError(bytes32 wrong_hash)
func (_Pos *PosFilterer) FilterWrongError(opts *bind.FilterOpts) (*PosWrongErrorIterator, error) {

	logs, sub, err := _Pos.contract.FilterLogs(opts, "WrongError")
	if err != nil {
		return nil, err
	}
	return &PosWrongErrorIterator{contract: _Pos.contract, event: "WrongError", logs: logs, sub: sub}, nil
}

// WatchWrongError is a free log subscription operation binding the contract event 0xba80455d9a4978272e947d7083631bd5cc9203cc9526eea9b5f32f21c1e90c57.
//
// Solidity: event WrongError(bytes32 wrong_hash)
func (_Pos *PosFilterer) WatchWrongError(opts *bind.WatchOpts, sink chan<- *PosWrongError) (event.Subscription, error) {

	logs, sub, err := _Pos.contract.WatchLogs(opts, "WrongError")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PosWrongError)
				if err := _Pos.contract.UnpackLog(event, "WrongError", log); err != nil {
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

// ParseWrongError is a log parse operation binding the contract event 0xba80455d9a4978272e947d7083631bd5cc9203cc9526eea9b5f32f21c1e90c57.
//
// Solidity: event WrongError(bytes32 wrong_hash)
func (_Pos *PosFilterer) ParseWrongError(log types.Log) (*PosWrongError, error) {
	event := new(PosWrongError)
	if err := _Pos.contract.UnpackLog(event, "WrongError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
