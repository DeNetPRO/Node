// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package proof_of_storage

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

// ProofOfStorageMetaData contains all meta data concerning the ProofOfStorage contract.
var ProofOfStorageMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_storage_address\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_payments_address\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_baseDifficulty\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"wrong_hash\",\"type\":\"bytes32\"}],\"name\":\"WrongError\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"REWARD_DIFFICULTY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"admin_set_user_data\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"base_difficulty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_storage_address\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_payments_address\",\"type\":\"address\"}],\"name\":\"changeSystemAddresses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"closeDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_n\",\"type\":\"uint32\"}],\"name\":\"getBlockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockNumber\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDifficulty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"getUserRewardInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"getUserRootHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_proof\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_targetDifficulty\",\"type\":\"uint256\"}],\"name\":\"isMatchDifficulty\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_root_hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"}],\"name\":\"isValidMerkleTreeProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isValidSign\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"makeDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"node_nft_address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"payments_address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user_address\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_block_number\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"_user_root_hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"_user_root_hash_nonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_user_signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_file\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"merkleProof\",\"type\":\"bytes32[]\"}],\"name\":\"sendProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_node_address\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_user_address\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_block_number\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"_user_root_hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"_user_root_hash_nonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_user_signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_file\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"merkleProof\",\"type\":\"bytes32[]\"}],\"name\":\"sendProofFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_new\",\"type\":\"address\"}],\"name\":\"setNodeNFTAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"setUserPlan\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_new_difficulty\",\"type\":\"uint256\"}],\"name\":\"updateBaseDifficulty\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"user_storage_address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_file\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"_block_number\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"_blocks_complited\",\"type\":\"uint256\"}],\"name\":\"verifyFileProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ProofOfStorageABI is the input ABI used to generate the binding from.
// Deprecated: Use ProofOfStorageMetaData.ABI instead.
var ProofOfStorageABI = ProofOfStorageMetaData.ABI

// ProofOfStorage is an auto generated Go binding around an Ethereum contract.
type ProofOfStorage struct {
	ProofOfStorageCaller     // Read-only binding to the contract
	ProofOfStorageTransactor // Write-only binding to the contract
	ProofOfStorageFilterer   // Log filterer for contract events
}

// ProofOfStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProofOfStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofOfStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProofOfStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofOfStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProofOfStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofOfStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProofOfStorageSession struct {
	Contract     *ProofOfStorage   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProofOfStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProofOfStorageCallerSession struct {
	Contract *ProofOfStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ProofOfStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProofOfStorageTransactorSession struct {
	Contract     *ProofOfStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ProofOfStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProofOfStorageRaw struct {
	Contract *ProofOfStorage // Generic contract binding to access the raw methods on
}

// ProofOfStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProofOfStorageCallerRaw struct {
	Contract *ProofOfStorageCaller // Generic read-only contract binding to access the raw methods on
}

// ProofOfStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProofOfStorageTransactorRaw struct {
	Contract *ProofOfStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProofOfStorage creates a new instance of ProofOfStorage, bound to a specific deployed contract.
func NewProofOfStorage(address common.Address, backend bind.ContractBackend) (*ProofOfStorage, error) {
	contract, err := bindProofOfStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ProofOfStorage{ProofOfStorageCaller: ProofOfStorageCaller{contract: contract}, ProofOfStorageTransactor: ProofOfStorageTransactor{contract: contract}, ProofOfStorageFilterer: ProofOfStorageFilterer{contract: contract}}, nil
}

// NewProofOfStorageCaller creates a new read-only instance of ProofOfStorage, bound to a specific deployed contract.
func NewProofOfStorageCaller(address common.Address, caller bind.ContractCaller) (*ProofOfStorageCaller, error) {
	contract, err := bindProofOfStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProofOfStorageCaller{contract: contract}, nil
}

// NewProofOfStorageTransactor creates a new write-only instance of ProofOfStorage, bound to a specific deployed contract.
func NewProofOfStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*ProofOfStorageTransactor, error) {
	contract, err := bindProofOfStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProofOfStorageTransactor{contract: contract}, nil
}

// NewProofOfStorageFilterer creates a new log filterer instance of ProofOfStorage, bound to a specific deployed contract.
func NewProofOfStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*ProofOfStorageFilterer, error) {
	contract, err := bindProofOfStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProofOfStorageFilterer{contract: contract}, nil
}

// bindProofOfStorage binds a generic wrapper to an already deployed contract.
func bindProofOfStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ProofOfStorageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProofOfStorage *ProofOfStorageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProofOfStorage.Contract.ProofOfStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProofOfStorage *ProofOfStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.ProofOfStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProofOfStorage *ProofOfStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.ProofOfStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProofOfStorage *ProofOfStorageCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProofOfStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProofOfStorage *ProofOfStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProofOfStorage *ProofOfStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.contract.Transact(opts, method, params...)
}

// REWARDDIFFICULTY is a free data retrieval call binding the contract method 0x60ab35e5.
//
// Solidity: function REWARD_DIFFICULTY() view returns(uint256)
func (_ProofOfStorage *ProofOfStorageCaller) REWARDDIFFICULTY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ProofOfStorage.contract.Call(opts, &out, "REWARD_DIFFICULTY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// REWARDDIFFICULTY is a free data retrieval call binding the contract method 0x60ab35e5.
//
// Solidity: function REWARD_DIFFICULTY() view returns(uint256)
func (_ProofOfStorage *ProofOfStorageSession) REWARDDIFFICULTY() (*big.Int, error) {
	return _ProofOfStorage.Contract.REWARDDIFFICULTY(&_ProofOfStorage.CallOpts)
}

// REWARDDIFFICULTY is a free data retrieval call binding the contract method 0x60ab35e5.
//
// Solidity: function REWARD_DIFFICULTY() view returns(uint256)
func (_ProofOfStorage *ProofOfStorageCallerSession) REWARDDIFFICULTY() (*big.Int, error) {
	return _ProofOfStorage.Contract.REWARDDIFFICULTY(&_ProofOfStorage.CallOpts)
}

// BaseDifficulty is a free data retrieval call binding the contract method 0x2d36de7a.
//
// Solidity: function base_difficulty() view returns(uint256)
func (_ProofOfStorage *ProofOfStorageCaller) BaseDifficulty(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ProofOfStorage.contract.Call(opts, &out, "base_difficulty")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseDifficulty is a free data retrieval call binding the contract method 0x2d36de7a.
//
// Solidity: function base_difficulty() view returns(uint256)
func (_ProofOfStorage *ProofOfStorageSession) BaseDifficulty() (*big.Int, error) {
	return _ProofOfStorage.Contract.BaseDifficulty(&_ProofOfStorage.CallOpts)
}

// BaseDifficulty is a free data retrieval call binding the contract method 0x2d36de7a.
//
// Solidity: function base_difficulty() view returns(uint256)
func (_ProofOfStorage *ProofOfStorageCallerSession) BaseDifficulty() (*big.Int, error) {
	return _ProofOfStorage.Contract.BaseDifficulty(&_ProofOfStorage.CallOpts)
}

// GetBlockHash is a free data retrieval call binding the contract method 0xd2b210a1.
//
// Solidity: function getBlockHash(uint32 _n) view returns(bytes32)
func (_ProofOfStorage *ProofOfStorageCaller) GetBlockHash(opts *bind.CallOpts, _n uint32) ([32]byte, error) {
	var out []interface{}
	err := _ProofOfStorage.contract.Call(opts, &out, "getBlockHash", _n)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBlockHash is a free data retrieval call binding the contract method 0xd2b210a1.
//
// Solidity: function getBlockHash(uint32 _n) view returns(bytes32)
func (_ProofOfStorage *ProofOfStorageSession) GetBlockHash(_n uint32) ([32]byte, error) {
	return _ProofOfStorage.Contract.GetBlockHash(&_ProofOfStorage.CallOpts, _n)
}

// GetBlockHash is a free data retrieval call binding the contract method 0xd2b210a1.
//
// Solidity: function getBlockHash(uint32 _n) view returns(bytes32)
func (_ProofOfStorage *ProofOfStorageCallerSession) GetBlockHash(_n uint32) ([32]byte, error) {
	return _ProofOfStorage.Contract.GetBlockHash(&_ProofOfStorage.CallOpts, _n)
}

// GetBlockNumber is a free data retrieval call binding the contract method 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint32)
func (_ProofOfStorage *ProofOfStorageCaller) GetBlockNumber(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ProofOfStorage.contract.Call(opts, &out, "getBlockNumber")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetBlockNumber is a free data retrieval call binding the contract method 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint32)
func (_ProofOfStorage *ProofOfStorageSession) GetBlockNumber() (uint32, error) {
	return _ProofOfStorage.Contract.GetBlockNumber(&_ProofOfStorage.CallOpts)
}

// GetBlockNumber is a free data retrieval call binding the contract method 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint32)
func (_ProofOfStorage *ProofOfStorageCallerSession) GetBlockNumber() (uint32, error) {
	return _ProofOfStorage.Contract.GetBlockNumber(&_ProofOfStorage.CallOpts)
}

// GetDifficulty is a free data retrieval call binding the contract method 0xb6baffe3.
//
// Solidity: function getDifficulty() view returns(uint256)
func (_ProofOfStorage *ProofOfStorageCaller) GetDifficulty(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ProofOfStorage.contract.Call(opts, &out, "getDifficulty")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDifficulty is a free data retrieval call binding the contract method 0xb6baffe3.
//
// Solidity: function getDifficulty() view returns(uint256)
func (_ProofOfStorage *ProofOfStorageSession) GetDifficulty() (*big.Int, error) {
	return _ProofOfStorage.Contract.GetDifficulty(&_ProofOfStorage.CallOpts)
}

// GetDifficulty is a free data retrieval call binding the contract method 0xb6baffe3.
//
// Solidity: function getDifficulty() view returns(uint256)
func (_ProofOfStorage *ProofOfStorageCallerSession) GetDifficulty() (*big.Int, error) {
	return _ProofOfStorage.Contract.GetDifficulty(&_ProofOfStorage.CallOpts)
}

// GetUserRewardInfo is a free data retrieval call binding the contract method 0x69517310.
//
// Solidity: function getUserRewardInfo(address _user) view returns(address, uint256, uint256)
func (_ProofOfStorage *ProofOfStorageCaller) GetUserRewardInfo(opts *bind.CallOpts, _user common.Address) (common.Address, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _ProofOfStorage.contract.Call(opts, &out, "getUserRewardInfo", _user)

	if err != nil {
		return *new(common.Address), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

// GetUserRewardInfo is a free data retrieval call binding the contract method 0x69517310.
//
// Solidity: function getUserRewardInfo(address _user) view returns(address, uint256, uint256)
func (_ProofOfStorage *ProofOfStorageSession) GetUserRewardInfo(_user common.Address) (common.Address, *big.Int, *big.Int, error) {
	return _ProofOfStorage.Contract.GetUserRewardInfo(&_ProofOfStorage.CallOpts, _user)
}

// GetUserRewardInfo is a free data retrieval call binding the contract method 0x69517310.
//
// Solidity: function getUserRewardInfo(address _user) view returns(address, uint256, uint256)
func (_ProofOfStorage *ProofOfStorageCallerSession) GetUserRewardInfo(_user common.Address) (common.Address, *big.Int, *big.Int, error) {
	return _ProofOfStorage.Contract.GetUserRewardInfo(&_ProofOfStorage.CallOpts, _user)
}

// GetUserRootHash is a free data retrieval call binding the contract method 0xf9a76fde.
//
// Solidity: function getUserRootHash(address _user) view returns(bytes32, uint256)
func (_ProofOfStorage *ProofOfStorageCaller) GetUserRootHash(opts *bind.CallOpts, _user common.Address) ([32]byte, *big.Int, error) {
	var out []interface{}
	err := _ProofOfStorage.contract.Call(opts, &out, "getUserRootHash", _user)

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
func (_ProofOfStorage *ProofOfStorageSession) GetUserRootHash(_user common.Address) ([32]byte, *big.Int, error) {
	return _ProofOfStorage.Contract.GetUserRootHash(&_ProofOfStorage.CallOpts, _user)
}

// GetUserRootHash is a free data retrieval call binding the contract method 0xf9a76fde.
//
// Solidity: function getUserRootHash(address _user) view returns(bytes32, uint256)
func (_ProofOfStorage *ProofOfStorageCallerSession) GetUserRootHash(_user common.Address) ([32]byte, *big.Int, error) {
	return _ProofOfStorage.Contract.GetUserRootHash(&_ProofOfStorage.CallOpts, _user)
}

// IsMatchDifficulty is a free data retrieval call binding the contract method 0x812adbd8.
//
// Solidity: function isMatchDifficulty(uint256 _proof, uint256 _targetDifficulty) view returns(bool)
func (_ProofOfStorage *ProofOfStorageCaller) IsMatchDifficulty(opts *bind.CallOpts, _proof *big.Int, _targetDifficulty *big.Int) (bool, error) {
	var out []interface{}
	err := _ProofOfStorage.contract.Call(opts, &out, "isMatchDifficulty", _proof, _targetDifficulty)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsMatchDifficulty is a free data retrieval call binding the contract method 0x812adbd8.
//
// Solidity: function isMatchDifficulty(uint256 _proof, uint256 _targetDifficulty) view returns(bool)
func (_ProofOfStorage *ProofOfStorageSession) IsMatchDifficulty(_proof *big.Int, _targetDifficulty *big.Int) (bool, error) {
	return _ProofOfStorage.Contract.IsMatchDifficulty(&_ProofOfStorage.CallOpts, _proof, _targetDifficulty)
}

// IsMatchDifficulty is a free data retrieval call binding the contract method 0x812adbd8.
//
// Solidity: function isMatchDifficulty(uint256 _proof, uint256 _targetDifficulty) view returns(bool)
func (_ProofOfStorage *ProofOfStorageCallerSession) IsMatchDifficulty(_proof *big.Int, _targetDifficulty *big.Int) (bool, error) {
	return _ProofOfStorage.Contract.IsMatchDifficulty(&_ProofOfStorage.CallOpts, _proof, _targetDifficulty)
}

// IsValidMerkleTreeProof is a free data retrieval call binding the contract method 0x059dbb30.
//
// Solidity: function isValidMerkleTreeProof(bytes32 _root_hash, bytes32[] proof) pure returns(bool)
func (_ProofOfStorage *ProofOfStorageCaller) IsValidMerkleTreeProof(opts *bind.CallOpts, _root_hash [32]byte, proof [][32]byte) (bool, error) {
	var out []interface{}
	err := _ProofOfStorage.contract.Call(opts, &out, "isValidMerkleTreeProof", _root_hash, proof)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidMerkleTreeProof is a free data retrieval call binding the contract method 0x059dbb30.
//
// Solidity: function isValidMerkleTreeProof(bytes32 _root_hash, bytes32[] proof) pure returns(bool)
func (_ProofOfStorage *ProofOfStorageSession) IsValidMerkleTreeProof(_root_hash [32]byte, proof [][32]byte) (bool, error) {
	return _ProofOfStorage.Contract.IsValidMerkleTreeProof(&_ProofOfStorage.CallOpts, _root_hash, proof)
}

// IsValidMerkleTreeProof is a free data retrieval call binding the contract method 0x059dbb30.
//
// Solidity: function isValidMerkleTreeProof(bytes32 _root_hash, bytes32[] proof) pure returns(bool)
func (_ProofOfStorage *ProofOfStorageCallerSession) IsValidMerkleTreeProof(_root_hash [32]byte, proof [][32]byte) (bool, error) {
	return _ProofOfStorage.Contract.IsValidMerkleTreeProof(&_ProofOfStorage.CallOpts, _root_hash, proof)
}

// IsValidSign is a free data retrieval call binding the contract method 0x16ecc0e6.
//
// Solidity: function isValidSign(address _signer, bytes message, bytes signature) pure returns(bool)
func (_ProofOfStorage *ProofOfStorageCaller) IsValidSign(opts *bind.CallOpts, _signer common.Address, message []byte, signature []byte) (bool, error) {
	var out []interface{}
	err := _ProofOfStorage.contract.Call(opts, &out, "isValidSign", _signer, message, signature)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidSign is a free data retrieval call binding the contract method 0x16ecc0e6.
//
// Solidity: function isValidSign(address _signer, bytes message, bytes signature) pure returns(bool)
func (_ProofOfStorage *ProofOfStorageSession) IsValidSign(_signer common.Address, message []byte, signature []byte) (bool, error) {
	return _ProofOfStorage.Contract.IsValidSign(&_ProofOfStorage.CallOpts, _signer, message, signature)
}

// IsValidSign is a free data retrieval call binding the contract method 0x16ecc0e6.
//
// Solidity: function isValidSign(address _signer, bytes message, bytes signature) pure returns(bool)
func (_ProofOfStorage *ProofOfStorageCallerSession) IsValidSign(_signer common.Address, message []byte, signature []byte) (bool, error) {
	return _ProofOfStorage.Contract.IsValidSign(&_ProofOfStorage.CallOpts, _signer, message, signature)
}

// NodeNftAddress is a free data retrieval call binding the contract method 0x81cfd709.
//
// Solidity: function node_nft_address() view returns(address)
func (_ProofOfStorage *ProofOfStorageCaller) NodeNftAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ProofOfStorage.contract.Call(opts, &out, "node_nft_address")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NodeNftAddress is a free data retrieval call binding the contract method 0x81cfd709.
//
// Solidity: function node_nft_address() view returns(address)
func (_ProofOfStorage *ProofOfStorageSession) NodeNftAddress() (common.Address, error) {
	return _ProofOfStorage.Contract.NodeNftAddress(&_ProofOfStorage.CallOpts)
}

// NodeNftAddress is a free data retrieval call binding the contract method 0x81cfd709.
//
// Solidity: function node_nft_address() view returns(address)
func (_ProofOfStorage *ProofOfStorageCallerSession) NodeNftAddress() (common.Address, error) {
	return _ProofOfStorage.Contract.NodeNftAddress(&_ProofOfStorage.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ProofOfStorage *ProofOfStorageCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ProofOfStorage.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ProofOfStorage *ProofOfStorageSession) Owner() (common.Address, error) {
	return _ProofOfStorage.Contract.Owner(&_ProofOfStorage.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ProofOfStorage *ProofOfStorageCallerSession) Owner() (common.Address, error) {
	return _ProofOfStorage.Contract.Owner(&_ProofOfStorage.CallOpts)
}

// PaymentsAddress is a free data retrieval call binding the contract method 0x95747678.
//
// Solidity: function payments_address() view returns(address)
func (_ProofOfStorage *ProofOfStorageCaller) PaymentsAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ProofOfStorage.contract.Call(opts, &out, "payments_address")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PaymentsAddress is a free data retrieval call binding the contract method 0x95747678.
//
// Solidity: function payments_address() view returns(address)
func (_ProofOfStorage *ProofOfStorageSession) PaymentsAddress() (common.Address, error) {
	return _ProofOfStorage.Contract.PaymentsAddress(&_ProofOfStorage.CallOpts)
}

// PaymentsAddress is a free data retrieval call binding the contract method 0x95747678.
//
// Solidity: function payments_address() view returns(address)
func (_ProofOfStorage *ProofOfStorageCallerSession) PaymentsAddress() (common.Address, error) {
	return _ProofOfStorage.Contract.PaymentsAddress(&_ProofOfStorage.CallOpts)
}

// UserStorageAddress is a free data retrieval call binding the contract method 0x1079a326.
//
// Solidity: function user_storage_address() view returns(address)
func (_ProofOfStorage *ProofOfStorageCaller) UserStorageAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ProofOfStorage.contract.Call(opts, &out, "user_storage_address")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UserStorageAddress is a free data retrieval call binding the contract method 0x1079a326.
//
// Solidity: function user_storage_address() view returns(address)
func (_ProofOfStorage *ProofOfStorageSession) UserStorageAddress() (common.Address, error) {
	return _ProofOfStorage.Contract.UserStorageAddress(&_ProofOfStorage.CallOpts)
}

// UserStorageAddress is a free data retrieval call binding the contract method 0x1079a326.
//
// Solidity: function user_storage_address() view returns(address)
func (_ProofOfStorage *ProofOfStorageCallerSession) UserStorageAddress() (common.Address, error) {
	return _ProofOfStorage.Contract.UserStorageAddress(&_ProofOfStorage.CallOpts)
}

// VerifyFileProof is a free data retrieval call binding the contract method 0x1d285f0a.
//
// Solidity: function verifyFileProof(address _sender, bytes _file, uint32 _block_number, uint256 _blocks_complited) view returns(bool)
func (_ProofOfStorage *ProofOfStorageCaller) VerifyFileProof(opts *bind.CallOpts, _sender common.Address, _file []byte, _block_number uint32, _blocks_complited *big.Int) (bool, error) {
	var out []interface{}
	err := _ProofOfStorage.contract.Call(opts, &out, "verifyFileProof", _sender, _file, _block_number, _blocks_complited)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyFileProof is a free data retrieval call binding the contract method 0x1d285f0a.
//
// Solidity: function verifyFileProof(address _sender, bytes _file, uint32 _block_number, uint256 _blocks_complited) view returns(bool)
func (_ProofOfStorage *ProofOfStorageSession) VerifyFileProof(_sender common.Address, _file []byte, _block_number uint32, _blocks_complited *big.Int) (bool, error) {
	return _ProofOfStorage.Contract.VerifyFileProof(&_ProofOfStorage.CallOpts, _sender, _file, _block_number, _blocks_complited)
}

// VerifyFileProof is a free data retrieval call binding the contract method 0x1d285f0a.
//
// Solidity: function verifyFileProof(address _sender, bytes _file, uint32 _block_number, uint256 _blocks_complited) view returns(bool)
func (_ProofOfStorage *ProofOfStorageCallerSession) VerifyFileProof(_sender common.Address, _file []byte, _block_number uint32, _blocks_complited *big.Int) (bool, error) {
	return _ProofOfStorage.Contract.VerifyFileProof(&_ProofOfStorage.CallOpts, _sender, _file, _block_number, _blocks_complited)
}

// AdminSetUserData is a paid mutator transaction binding the contract method 0x06db598a.
//
// Solidity: function admin_set_user_data(address _from, address _user, address _token, uint256 _amount) returns()
func (_ProofOfStorage *ProofOfStorageTransactor) AdminSetUserData(opts *bind.TransactOpts, _from common.Address, _user common.Address, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ProofOfStorage.contract.Transact(opts, "admin_set_user_data", _from, _user, _token, _amount)
}

// AdminSetUserData is a paid mutator transaction binding the contract method 0x06db598a.
//
// Solidity: function admin_set_user_data(address _from, address _user, address _token, uint256 _amount) returns()
func (_ProofOfStorage *ProofOfStorageSession) AdminSetUserData(_from common.Address, _user common.Address, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.AdminSetUserData(&_ProofOfStorage.TransactOpts, _from, _user, _token, _amount)
}

// AdminSetUserData is a paid mutator transaction binding the contract method 0x06db598a.
//
// Solidity: function admin_set_user_data(address _from, address _user, address _token, uint256 _amount) returns()
func (_ProofOfStorage *ProofOfStorageTransactorSession) AdminSetUserData(_from common.Address, _user common.Address, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.AdminSetUserData(&_ProofOfStorage.TransactOpts, _from, _user, _token, _amount)
}

// ChangeSystemAddresses is a paid mutator transaction binding the contract method 0xeb5439dc.
//
// Solidity: function changeSystemAddresses(address _storage_address, address _payments_address) returns()
func (_ProofOfStorage *ProofOfStorageTransactor) ChangeSystemAddresses(opts *bind.TransactOpts, _storage_address common.Address, _payments_address common.Address) (*types.Transaction, error) {
	return _ProofOfStorage.contract.Transact(opts, "changeSystemAddresses", _storage_address, _payments_address)
}

// ChangeSystemAddresses is a paid mutator transaction binding the contract method 0xeb5439dc.
//
// Solidity: function changeSystemAddresses(address _storage_address, address _payments_address) returns()
func (_ProofOfStorage *ProofOfStorageSession) ChangeSystemAddresses(_storage_address common.Address, _payments_address common.Address) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.ChangeSystemAddresses(&_ProofOfStorage.TransactOpts, _storage_address, _payments_address)
}

// ChangeSystemAddresses is a paid mutator transaction binding the contract method 0xeb5439dc.
//
// Solidity: function changeSystemAddresses(address _storage_address, address _payments_address) returns()
func (_ProofOfStorage *ProofOfStorageTransactorSession) ChangeSystemAddresses(_storage_address common.Address, _payments_address common.Address) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.ChangeSystemAddresses(&_ProofOfStorage.TransactOpts, _storage_address, _payments_address)
}

// CloseDeposit is a paid mutator transaction binding the contract method 0xdc8d0a40.
//
// Solidity: function closeDeposit(address _token) returns()
func (_ProofOfStorage *ProofOfStorageTransactor) CloseDeposit(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _ProofOfStorage.contract.Transact(opts, "closeDeposit", _token)
}

// CloseDeposit is a paid mutator transaction binding the contract method 0xdc8d0a40.
//
// Solidity: function closeDeposit(address _token) returns()
func (_ProofOfStorage *ProofOfStorageSession) CloseDeposit(_token common.Address) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.CloseDeposit(&_ProofOfStorage.TransactOpts, _token)
}

// CloseDeposit is a paid mutator transaction binding the contract method 0xdc8d0a40.
//
// Solidity: function closeDeposit(address _token) returns()
func (_ProofOfStorage *ProofOfStorageTransactorSession) CloseDeposit(_token common.Address) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.CloseDeposit(&_ProofOfStorage.TransactOpts, _token)
}

// MakeDeposit is a paid mutator transaction binding the contract method 0xbf5d0a00.
//
// Solidity: function makeDeposit(address _token, uint256 _amount) returns()
func (_ProofOfStorage *ProofOfStorageTransactor) MakeDeposit(opts *bind.TransactOpts, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ProofOfStorage.contract.Transact(opts, "makeDeposit", _token, _amount)
}

// MakeDeposit is a paid mutator transaction binding the contract method 0xbf5d0a00.
//
// Solidity: function makeDeposit(address _token, uint256 _amount) returns()
func (_ProofOfStorage *ProofOfStorageSession) MakeDeposit(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.MakeDeposit(&_ProofOfStorage.TransactOpts, _token, _amount)
}

// MakeDeposit is a paid mutator transaction binding the contract method 0xbf5d0a00.
//
// Solidity: function makeDeposit(address _token, uint256 _amount) returns()
func (_ProofOfStorage *ProofOfStorageTransactorSession) MakeDeposit(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.MakeDeposit(&_ProofOfStorage.TransactOpts, _token, _amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ProofOfStorage *ProofOfStorageTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProofOfStorage.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ProofOfStorage *ProofOfStorageSession) RenounceOwnership() (*types.Transaction, error) {
	return _ProofOfStorage.Contract.RenounceOwnership(&_ProofOfStorage.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ProofOfStorage *ProofOfStorageTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ProofOfStorage.Contract.RenounceOwnership(&_ProofOfStorage.TransactOpts)
}

// SendProof is a paid mutator transaction binding the contract method 0x74556ca9.
//
// Solidity: function sendProof(address _user_address, uint32 _block_number, bytes32 _user_root_hash, uint64 _user_root_hash_nonce, bytes _user_signature, bytes _file, bytes32[] merkleProof) returns()
func (_ProofOfStorage *ProofOfStorageTransactor) SendProof(opts *bind.TransactOpts, _user_address common.Address, _block_number uint32, _user_root_hash [32]byte, _user_root_hash_nonce uint64, _user_signature []byte, _file []byte, merkleProof [][32]byte) (*types.Transaction, error) {
	return _ProofOfStorage.contract.Transact(opts, "sendProof", _user_address, _block_number, _user_root_hash, _user_root_hash_nonce, _user_signature, _file, merkleProof)
}

// SendProof is a paid mutator transaction binding the contract method 0x74556ca9.
//
// Solidity: function sendProof(address _user_address, uint32 _block_number, bytes32 _user_root_hash, uint64 _user_root_hash_nonce, bytes _user_signature, bytes _file, bytes32[] merkleProof) returns()
func (_ProofOfStorage *ProofOfStorageSession) SendProof(_user_address common.Address, _block_number uint32, _user_root_hash [32]byte, _user_root_hash_nonce uint64, _user_signature []byte, _file []byte, merkleProof [][32]byte) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.SendProof(&_ProofOfStorage.TransactOpts, _user_address, _block_number, _user_root_hash, _user_root_hash_nonce, _user_signature, _file, merkleProof)
}

// SendProof is a paid mutator transaction binding the contract method 0x74556ca9.
//
// Solidity: function sendProof(address _user_address, uint32 _block_number, bytes32 _user_root_hash, uint64 _user_root_hash_nonce, bytes _user_signature, bytes _file, bytes32[] merkleProof) returns()
func (_ProofOfStorage *ProofOfStorageTransactorSession) SendProof(_user_address common.Address, _block_number uint32, _user_root_hash [32]byte, _user_root_hash_nonce uint64, _user_signature []byte, _file []byte, merkleProof [][32]byte) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.SendProof(&_ProofOfStorage.TransactOpts, _user_address, _block_number, _user_root_hash, _user_root_hash_nonce, _user_signature, _file, merkleProof)
}

// SendProofFrom is a paid mutator transaction binding the contract method 0x46f03dd5.
//
// Solidity: function sendProofFrom(address _node_address, address _user_address, uint32 _block_number, bytes32 _user_root_hash, uint64 _user_root_hash_nonce, bytes _user_signature, bytes _file, bytes32[] merkleProof) returns()
func (_ProofOfStorage *ProofOfStorageTransactor) SendProofFrom(opts *bind.TransactOpts, _node_address common.Address, _user_address common.Address, _block_number uint32, _user_root_hash [32]byte, _user_root_hash_nonce uint64, _user_signature []byte, _file []byte, merkleProof [][32]byte) (*types.Transaction, error) {
	return _ProofOfStorage.contract.Transact(opts, "sendProofFrom", _node_address, _user_address, _block_number, _user_root_hash, _user_root_hash_nonce, _user_signature, _file, merkleProof)
}

// SendProofFrom is a paid mutator transaction binding the contract method 0x46f03dd5.
//
// Solidity: function sendProofFrom(address _node_address, address _user_address, uint32 _block_number, bytes32 _user_root_hash, uint64 _user_root_hash_nonce, bytes _user_signature, bytes _file, bytes32[] merkleProof) returns()
func (_ProofOfStorage *ProofOfStorageSession) SendProofFrom(_node_address common.Address, _user_address common.Address, _block_number uint32, _user_root_hash [32]byte, _user_root_hash_nonce uint64, _user_signature []byte, _file []byte, merkleProof [][32]byte) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.SendProofFrom(&_ProofOfStorage.TransactOpts, _node_address, _user_address, _block_number, _user_root_hash, _user_root_hash_nonce, _user_signature, _file, merkleProof)
}

// SendProofFrom is a paid mutator transaction binding the contract method 0x46f03dd5.
//
// Solidity: function sendProofFrom(address _node_address, address _user_address, uint32 _block_number, bytes32 _user_root_hash, uint64 _user_root_hash_nonce, bytes _user_signature, bytes _file, bytes32[] merkleProof) returns()
func (_ProofOfStorage *ProofOfStorageTransactorSession) SendProofFrom(_node_address common.Address, _user_address common.Address, _block_number uint32, _user_root_hash [32]byte, _user_root_hash_nonce uint64, _user_signature []byte, _file []byte, merkleProof [][32]byte) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.SendProofFrom(&_ProofOfStorage.TransactOpts, _node_address, _user_address, _block_number, _user_root_hash, _user_root_hash_nonce, _user_signature, _file, merkleProof)
}

// SetNodeNFTAddress is a paid mutator transaction binding the contract method 0xca811507.
//
// Solidity: function setNodeNFTAddress(address _new) returns()
func (_ProofOfStorage *ProofOfStorageTransactor) SetNodeNFTAddress(opts *bind.TransactOpts, _new common.Address) (*types.Transaction, error) {
	return _ProofOfStorage.contract.Transact(opts, "setNodeNFTAddress", _new)
}

// SetNodeNFTAddress is a paid mutator transaction binding the contract method 0xca811507.
//
// Solidity: function setNodeNFTAddress(address _new) returns()
func (_ProofOfStorage *ProofOfStorageSession) SetNodeNFTAddress(_new common.Address) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.SetNodeNFTAddress(&_ProofOfStorage.TransactOpts, _new)
}

// SetNodeNFTAddress is a paid mutator transaction binding the contract method 0xca811507.
//
// Solidity: function setNodeNFTAddress(address _new) returns()
func (_ProofOfStorage *ProofOfStorageTransactorSession) SetNodeNFTAddress(_new common.Address) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.SetNodeNFTAddress(&_ProofOfStorage.TransactOpts, _new)
}

// SetUserPlan is a paid mutator transaction binding the contract method 0xe067ae4c.
//
// Solidity: function setUserPlan(address _token) returns()
func (_ProofOfStorage *ProofOfStorageTransactor) SetUserPlan(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _ProofOfStorage.contract.Transact(opts, "setUserPlan", _token)
}

// SetUserPlan is a paid mutator transaction binding the contract method 0xe067ae4c.
//
// Solidity: function setUserPlan(address _token) returns()
func (_ProofOfStorage *ProofOfStorageSession) SetUserPlan(_token common.Address) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.SetUserPlan(&_ProofOfStorage.TransactOpts, _token)
}

// SetUserPlan is a paid mutator transaction binding the contract method 0xe067ae4c.
//
// Solidity: function setUserPlan(address _token) returns()
func (_ProofOfStorage *ProofOfStorageTransactorSession) SetUserPlan(_token common.Address) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.SetUserPlan(&_ProofOfStorage.TransactOpts, _token)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ProofOfStorage *ProofOfStorageTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ProofOfStorage.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ProofOfStorage *ProofOfStorageSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.TransferOwnership(&_ProofOfStorage.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ProofOfStorage *ProofOfStorageTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.TransferOwnership(&_ProofOfStorage.TransactOpts, newOwner)
}

// UpdateBaseDifficulty is a paid mutator transaction binding the contract method 0xb5c6da27.
//
// Solidity: function updateBaseDifficulty(uint256 _new_difficulty) returns()
func (_ProofOfStorage *ProofOfStorageTransactor) UpdateBaseDifficulty(opts *bind.TransactOpts, _new_difficulty *big.Int) (*types.Transaction, error) {
	return _ProofOfStorage.contract.Transact(opts, "updateBaseDifficulty", _new_difficulty)
}

// UpdateBaseDifficulty is a paid mutator transaction binding the contract method 0xb5c6da27.
//
// Solidity: function updateBaseDifficulty(uint256 _new_difficulty) returns()
func (_ProofOfStorage *ProofOfStorageSession) UpdateBaseDifficulty(_new_difficulty *big.Int) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.UpdateBaseDifficulty(&_ProofOfStorage.TransactOpts, _new_difficulty)
}

// UpdateBaseDifficulty is a paid mutator transaction binding the contract method 0xb5c6da27.
//
// Solidity: function updateBaseDifficulty(uint256 _new_difficulty) returns()
func (_ProofOfStorage *ProofOfStorageTransactorSession) UpdateBaseDifficulty(_new_difficulty *big.Int) (*types.Transaction, error) {
	return _ProofOfStorage.Contract.UpdateBaseDifficulty(&_ProofOfStorage.TransactOpts, _new_difficulty)
}

// ProofOfStorageOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ProofOfStorage contract.
type ProofOfStorageOwnershipTransferredIterator struct {
	Event *ProofOfStorageOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ProofOfStorageOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofOfStorageOwnershipTransferred)
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
		it.Event = new(ProofOfStorageOwnershipTransferred)
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
func (it *ProofOfStorageOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofOfStorageOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofOfStorageOwnershipTransferred represents a OwnershipTransferred event raised by the ProofOfStorage contract.
type ProofOfStorageOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ProofOfStorage *ProofOfStorageFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ProofOfStorageOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ProofOfStorage.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ProofOfStorageOwnershipTransferredIterator{contract: _ProofOfStorage.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ProofOfStorage *ProofOfStorageFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ProofOfStorageOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ProofOfStorage.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofOfStorageOwnershipTransferred)
				if err := _ProofOfStorage.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ProofOfStorage *ProofOfStorageFilterer) ParseOwnershipTransferred(log types.Log) (*ProofOfStorageOwnershipTransferred, error) {
	event := new(ProofOfStorageOwnershipTransferred)
	if err := _ProofOfStorage.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofOfStorageWrongErrorIterator is returned from FilterWrongError and is used to iterate over the raw logs and unpacked data for WrongError events raised by the ProofOfStorage contract.
type ProofOfStorageWrongErrorIterator struct {
	Event *ProofOfStorageWrongError // Event containing the contract specifics and raw log

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
func (it *ProofOfStorageWrongErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofOfStorageWrongError)
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
		it.Event = new(ProofOfStorageWrongError)
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
func (it *ProofOfStorageWrongErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofOfStorageWrongErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofOfStorageWrongError represents a WrongError event raised by the ProofOfStorage contract.
type ProofOfStorageWrongError struct {
	WrongHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWrongError is a free log retrieval operation binding the contract event 0xba80455d9a4978272e947d7083631bd5cc9203cc9526eea9b5f32f21c1e90c57.
//
// Solidity: event WrongError(bytes32 wrong_hash)
func (_ProofOfStorage *ProofOfStorageFilterer) FilterWrongError(opts *bind.FilterOpts) (*ProofOfStorageWrongErrorIterator, error) {

	logs, sub, err := _ProofOfStorage.contract.FilterLogs(opts, "WrongError")
	if err != nil {
		return nil, err
	}
	return &ProofOfStorageWrongErrorIterator{contract: _ProofOfStorage.contract, event: "WrongError", logs: logs, sub: sub}, nil
}

// WatchWrongError is a free log subscription operation binding the contract event 0xba80455d9a4978272e947d7083631bd5cc9203cc9526eea9b5f32f21c1e90c57.
//
// Solidity: event WrongError(bytes32 wrong_hash)
func (_ProofOfStorage *ProofOfStorageFilterer) WatchWrongError(opts *bind.WatchOpts, sink chan<- *ProofOfStorageWrongError) (event.Subscription, error) {

	logs, sub, err := _ProofOfStorage.contract.WatchLogs(opts, "WrongError")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofOfStorageWrongError)
				if err := _ProofOfStorage.contract.UnpackLog(event, "WrongError", log); err != nil {
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
func (_ProofOfStorage *ProofOfStorageFilterer) ParseWrongError(log types.Log) (*ProofOfStorageWrongError, error) {
	event := new(ProofOfStorageWrongError)
	if err := _ProofOfStorage.contract.UnpackLog(event, "WrongError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
