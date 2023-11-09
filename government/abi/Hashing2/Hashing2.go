// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Hashing2

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
	_ = abi.ConvertType
)

// TypesOutputRootProof is an auto generated low-level Go binding around an user-defined struct.
type TypesOutputRootProof struct {
	Version                  [32]byte
	StateRoot                [32]byte
	MessagePasserStorageRoot [32]byte
	LatestBlockhash          [32]byte
}

// TypesWithdrawalTransaction is an auto generated low-level Go binding around an user-defined struct.
type TypesWithdrawalTransaction struct {
	Nonce    *big.Int
	Sender   common.Address
	Target   common.Address
	Value    *big.Int
	GasLimit *big.Int
	Data     []byte
}

// Hashing2MetaData contains all meta data concerning the Hashing2 contract.
var Hashing2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"messagePasserStorageRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"latestBlockhash\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.OutputRootProof\",\"name\":\"_outputRootProof\",\"type\":\"tuple\"}],\"name\":\"hashOutputRootProof\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"messagePasserStorageRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"latestBlockhash\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.OutputRootProof\",\"name\":\"_outputRootProof\",\"type\":\"tuple\"}],\"name\":\"hashOutputRootProof2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.WithdrawalTransaction\",\"name\":\"_tx\",\"type\":\"tuple\"}],\"name\":\"hashWithdrawal\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestBlockhash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messagePasserStorageRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"outputRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// Hashing2ABI is the input ABI used to generate the binding from.
// Deprecated: Use Hashing2MetaData.ABI instead.
var Hashing2ABI = Hashing2MetaData.ABI

// Hashing2 is an auto generated Go binding around an Ethereum contract.
type Hashing2 struct {
	Hashing2Caller     // Read-only binding to the contract
	Hashing2Transactor // Write-only binding to the contract
	Hashing2Filterer   // Log filterer for contract events
}

// Hashing2Caller is an auto generated read-only Go binding around an Ethereum contract.
type Hashing2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Hashing2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Hashing2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Hashing2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Hashing2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Hashing2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Hashing2Session struct {
	Contract     *Hashing2         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Hashing2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Hashing2CallerSession struct {
	Contract *Hashing2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// Hashing2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Hashing2TransactorSession struct {
	Contract     *Hashing2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// Hashing2Raw is an auto generated low-level Go binding around an Ethereum contract.
type Hashing2Raw struct {
	Contract *Hashing2 // Generic contract binding to access the raw methods on
}

// Hashing2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Hashing2CallerRaw struct {
	Contract *Hashing2Caller // Generic read-only contract binding to access the raw methods on
}

// Hashing2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Hashing2TransactorRaw struct {
	Contract *Hashing2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewHashing2 creates a new instance of Hashing2, bound to a specific deployed contract.
func NewHashing2(address common.Address, backend bind.ContractBackend) (*Hashing2, error) {
	contract, err := bindHashing2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Hashing2{Hashing2Caller: Hashing2Caller{contract: contract}, Hashing2Transactor: Hashing2Transactor{contract: contract}, Hashing2Filterer: Hashing2Filterer{contract: contract}}, nil
}

// NewHashing2Caller creates a new read-only instance of Hashing2, bound to a specific deployed contract.
func NewHashing2Caller(address common.Address, caller bind.ContractCaller) (*Hashing2Caller, error) {
	contract, err := bindHashing2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Hashing2Caller{contract: contract}, nil
}

// NewHashing2Transactor creates a new write-only instance of Hashing2, bound to a specific deployed contract.
func NewHashing2Transactor(address common.Address, transactor bind.ContractTransactor) (*Hashing2Transactor, error) {
	contract, err := bindHashing2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Hashing2Transactor{contract: contract}, nil
}

// NewHashing2Filterer creates a new log filterer instance of Hashing2, bound to a specific deployed contract.
func NewHashing2Filterer(address common.Address, filterer bind.ContractFilterer) (*Hashing2Filterer, error) {
	contract, err := bindHashing2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Hashing2Filterer{contract: contract}, nil
}

// bindHashing2 binds a generic wrapper to an already deployed contract.
func bindHashing2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Hashing2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Hashing2 *Hashing2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Hashing2.Contract.Hashing2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Hashing2 *Hashing2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Hashing2.Contract.Hashing2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Hashing2 *Hashing2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Hashing2.Contract.Hashing2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Hashing2 *Hashing2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Hashing2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Hashing2 *Hashing2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Hashing2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Hashing2 *Hashing2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Hashing2.Contract.contract.Transact(opts, method, params...)
}

// HashOutputRootProof is a free data retrieval call binding the contract method 0x24d2bbb0.
//
// Solidity: function hashOutputRootProof((bytes32,bytes32,bytes32,bytes32) _outputRootProof) pure returns(bytes32)
func (_Hashing2 *Hashing2Caller) HashOutputRootProof(opts *bind.CallOpts, _outputRootProof TypesOutputRootProof) ([32]byte, error) {
	var out []interface{}
	err := _Hashing2.contract.Call(opts, &out, "hashOutputRootProof", _outputRootProof)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashOutputRootProof is a free data retrieval call binding the contract method 0x24d2bbb0.
//
// Solidity: function hashOutputRootProof((bytes32,bytes32,bytes32,bytes32) _outputRootProof) pure returns(bytes32)
func (_Hashing2 *Hashing2Session) HashOutputRootProof(_outputRootProof TypesOutputRootProof) ([32]byte, error) {
	return _Hashing2.Contract.HashOutputRootProof(&_Hashing2.CallOpts, _outputRootProof)
}

// HashOutputRootProof is a free data retrieval call binding the contract method 0x24d2bbb0.
//
// Solidity: function hashOutputRootProof((bytes32,bytes32,bytes32,bytes32) _outputRootProof) pure returns(bytes32)
func (_Hashing2 *Hashing2CallerSession) HashOutputRootProof(_outputRootProof TypesOutputRootProof) ([32]byte, error) {
	return _Hashing2.Contract.HashOutputRootProof(&_Hashing2.CallOpts, _outputRootProof)
}

// HashWithdrawal is a free data retrieval call binding the contract method 0x7d4395ac.
//
// Solidity: function hashWithdrawal((uint256,address,address,uint256,uint256,bytes) _tx) pure returns(bytes32)
func (_Hashing2 *Hashing2Caller) HashWithdrawal(opts *bind.CallOpts, _tx TypesWithdrawalTransaction) ([32]byte, error) {
	var out []interface{}
	err := _Hashing2.contract.Call(opts, &out, "hashWithdrawal", _tx)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashWithdrawal is a free data retrieval call binding the contract method 0x7d4395ac.
//
// Solidity: function hashWithdrawal((uint256,address,address,uint256,uint256,bytes) _tx) pure returns(bytes32)
func (_Hashing2 *Hashing2Session) HashWithdrawal(_tx TypesWithdrawalTransaction) ([32]byte, error) {
	return _Hashing2.Contract.HashWithdrawal(&_Hashing2.CallOpts, _tx)
}

// HashWithdrawal is a free data retrieval call binding the contract method 0x7d4395ac.
//
// Solidity: function hashWithdrawal((uint256,address,address,uint256,uint256,bytes) _tx) pure returns(bytes32)
func (_Hashing2 *Hashing2CallerSession) HashWithdrawal(_tx TypesWithdrawalTransaction) ([32]byte, error) {
	return _Hashing2.Contract.HashWithdrawal(&_Hashing2.CallOpts, _tx)
}

// LatestBlockhash is a free data retrieval call binding the contract method 0xeda743d3.
//
// Solidity: function latestBlockhash() view returns(bytes32)
func (_Hashing2 *Hashing2Caller) LatestBlockhash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Hashing2.contract.Call(opts, &out, "latestBlockhash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LatestBlockhash is a free data retrieval call binding the contract method 0xeda743d3.
//
// Solidity: function latestBlockhash() view returns(bytes32)
func (_Hashing2 *Hashing2Session) LatestBlockhash() ([32]byte, error) {
	return _Hashing2.Contract.LatestBlockhash(&_Hashing2.CallOpts)
}

// LatestBlockhash is a free data retrieval call binding the contract method 0xeda743d3.
//
// Solidity: function latestBlockhash() view returns(bytes32)
func (_Hashing2 *Hashing2CallerSession) LatestBlockhash() ([32]byte, error) {
	return _Hashing2.Contract.LatestBlockhash(&_Hashing2.CallOpts)
}

// MessagePasserStorageRoot is a free data retrieval call binding the contract method 0x109ddb75.
//
// Solidity: function messagePasserStorageRoot() view returns(bytes32)
func (_Hashing2 *Hashing2Caller) MessagePasserStorageRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Hashing2.contract.Call(opts, &out, "messagePasserStorageRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MessagePasserStorageRoot is a free data retrieval call binding the contract method 0x109ddb75.
//
// Solidity: function messagePasserStorageRoot() view returns(bytes32)
func (_Hashing2 *Hashing2Session) MessagePasserStorageRoot() ([32]byte, error) {
	return _Hashing2.Contract.MessagePasserStorageRoot(&_Hashing2.CallOpts)
}

// MessagePasserStorageRoot is a free data retrieval call binding the contract method 0x109ddb75.
//
// Solidity: function messagePasserStorageRoot() view returns(bytes32)
func (_Hashing2 *Hashing2CallerSession) MessagePasserStorageRoot() ([32]byte, error) {
	return _Hashing2.Contract.MessagePasserStorageRoot(&_Hashing2.CallOpts)
}

// OutputRoot is a free data retrieval call binding the contract method 0xa78deacf.
//
// Solidity: function outputRoot() view returns(bytes32)
func (_Hashing2 *Hashing2Caller) OutputRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Hashing2.contract.Call(opts, &out, "outputRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// OutputRoot is a free data retrieval call binding the contract method 0xa78deacf.
//
// Solidity: function outputRoot() view returns(bytes32)
func (_Hashing2 *Hashing2Session) OutputRoot() ([32]byte, error) {
	return _Hashing2.Contract.OutputRoot(&_Hashing2.CallOpts)
}

// OutputRoot is a free data retrieval call binding the contract method 0xa78deacf.
//
// Solidity: function outputRoot() view returns(bytes32)
func (_Hashing2 *Hashing2CallerSession) OutputRoot() ([32]byte, error) {
	return _Hashing2.Contract.OutputRoot(&_Hashing2.CallOpts)
}

// StateRoot is a free data retrieval call binding the contract method 0x9588eca2.
//
// Solidity: function stateRoot() view returns(bytes32)
func (_Hashing2 *Hashing2Caller) StateRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Hashing2.contract.Call(opts, &out, "stateRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// StateRoot is a free data retrieval call binding the contract method 0x9588eca2.
//
// Solidity: function stateRoot() view returns(bytes32)
func (_Hashing2 *Hashing2Session) StateRoot() ([32]byte, error) {
	return _Hashing2.Contract.StateRoot(&_Hashing2.CallOpts)
}

// StateRoot is a free data retrieval call binding the contract method 0x9588eca2.
//
// Solidity: function stateRoot() view returns(bytes32)
func (_Hashing2 *Hashing2CallerSession) StateRoot() ([32]byte, error) {
	return _Hashing2.Contract.StateRoot(&_Hashing2.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(bytes32)
func (_Hashing2 *Hashing2Caller) Version(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Hashing2.contract.Call(opts, &out, "version")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(bytes32)
func (_Hashing2 *Hashing2Session) Version() ([32]byte, error) {
	return _Hashing2.Contract.Version(&_Hashing2.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(bytes32)
func (_Hashing2 *Hashing2CallerSession) Version() ([32]byte, error) {
	return _Hashing2.Contract.Version(&_Hashing2.CallOpts)
}

// HashOutputRootProof2 is a paid mutator transaction binding the contract method 0x441d7303.
//
// Solidity: function hashOutputRootProof2((bytes32,bytes32,bytes32,bytes32) _outputRootProof) returns()
func (_Hashing2 *Hashing2Transactor) HashOutputRootProof2(opts *bind.TransactOpts, _outputRootProof TypesOutputRootProof) (*types.Transaction, error) {
	return _Hashing2.contract.Transact(opts, "hashOutputRootProof2", _outputRootProof)
}

// HashOutputRootProof2 is a paid mutator transaction binding the contract method 0x441d7303.
//
// Solidity: function hashOutputRootProof2((bytes32,bytes32,bytes32,bytes32) _outputRootProof) returns()
func (_Hashing2 *Hashing2Session) HashOutputRootProof2(_outputRootProof TypesOutputRootProof) (*types.Transaction, error) {
	return _Hashing2.Contract.HashOutputRootProof2(&_Hashing2.TransactOpts, _outputRootProof)
}

// HashOutputRootProof2 is a paid mutator transaction binding the contract method 0x441d7303.
//
// Solidity: function hashOutputRootProof2((bytes32,bytes32,bytes32,bytes32) _outputRootProof) returns()
func (_Hashing2 *Hashing2TransactorSession) HashOutputRootProof2(_outputRootProof TypesOutputRootProof) (*types.Transaction, error) {
	return _Hashing2.Contract.HashOutputRootProof2(&_Hashing2.TransactOpts, _outputRootProof)
}
