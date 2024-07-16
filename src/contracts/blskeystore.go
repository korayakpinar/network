// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// BLSKeystoreMetaData contains all meta data concerning the BLSKeystore contract.
var BLSKeystoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"initialPubKey\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"GetPubKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PubKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// BLSKeystoreABI is the input ABI used to generate the binding from.
// Deprecated: Use BLSKeystoreMetaData.ABI instead.
var BLSKeystoreABI = BLSKeystoreMetaData.ABI

// BLSKeystore is an auto generated Go binding around an Ethereum contract.
type BLSKeystore struct {
	BLSKeystoreCaller     // Read-only binding to the contract
	BLSKeystoreTransactor // Write-only binding to the contract
	BLSKeystoreFilterer   // Log filterer for contract events
}

// BLSKeystoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type BLSKeystoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BLSKeystoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BLSKeystoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BLSKeystoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BLSKeystoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BLSKeystoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BLSKeystoreSession struct {
	Contract     *BLSKeystore      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BLSKeystoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BLSKeystoreCallerSession struct {
	Contract *BLSKeystoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// BLSKeystoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BLSKeystoreTransactorSession struct {
	Contract     *BLSKeystoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BLSKeystoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type BLSKeystoreRaw struct {
	Contract *BLSKeystore // Generic contract binding to access the raw methods on
}

// BLSKeystoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BLSKeystoreCallerRaw struct {
	Contract *BLSKeystoreCaller // Generic read-only contract binding to access the raw methods on
}

// BLSKeystoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BLSKeystoreTransactorRaw struct {
	Contract *BLSKeystoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBLSKeystore creates a new instance of BLSKeystore, bound to a specific deployed contract.
func NewBLSKeystore(address common.Address, backend bind.ContractBackend) (*BLSKeystore, error) {
	contract, err := bindBLSKeystore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BLSKeystore{BLSKeystoreCaller: BLSKeystoreCaller{contract: contract}, BLSKeystoreTransactor: BLSKeystoreTransactor{contract: contract}, BLSKeystoreFilterer: BLSKeystoreFilterer{contract: contract}}, nil
}

// NewBLSKeystoreCaller creates a new read-only instance of BLSKeystore, bound to a specific deployed contract.
func NewBLSKeystoreCaller(address common.Address, caller bind.ContractCaller) (*BLSKeystoreCaller, error) {
	contract, err := bindBLSKeystore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BLSKeystoreCaller{contract: contract}, nil
}

// NewBLSKeystoreTransactor creates a new write-only instance of BLSKeystore, bound to a specific deployed contract.
func NewBLSKeystoreTransactor(address common.Address, transactor bind.ContractTransactor) (*BLSKeystoreTransactor, error) {
	contract, err := bindBLSKeystore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BLSKeystoreTransactor{contract: contract}, nil
}

// NewBLSKeystoreFilterer creates a new log filterer instance of BLSKeystore, bound to a specific deployed contract.
func NewBLSKeystoreFilterer(address common.Address, filterer bind.ContractFilterer) (*BLSKeystoreFilterer, error) {
	contract, err := bindBLSKeystore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BLSKeystoreFilterer{contract: contract}, nil
}

// bindBLSKeystore binds a generic wrapper to an already deployed contract.
func bindBLSKeystore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BLSKeystoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BLSKeystore *BLSKeystoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BLSKeystore.Contract.BLSKeystoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BLSKeystore *BLSKeystoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BLSKeystore.Contract.BLSKeystoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BLSKeystore *BLSKeystoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BLSKeystore.Contract.BLSKeystoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BLSKeystore *BLSKeystoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BLSKeystore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BLSKeystore *BLSKeystoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BLSKeystore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BLSKeystore *BLSKeystoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BLSKeystore.Contract.contract.Transact(opts, method, params...)
}

// GetPubKey is a free data retrieval call binding the contract method 0x57c8a80b.
//
// Solidity: function GetPubKey() view returns(bytes)
func (_BLSKeystore *BLSKeystoreCaller) GetPubKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _BLSKeystore.contract.Call(opts, &out, "GetPubKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetPubKey is a free data retrieval call binding the contract method 0x57c8a80b.
//
// Solidity: function GetPubKey() view returns(bytes)
func (_BLSKeystore *BLSKeystoreSession) GetPubKey() ([]byte, error) {
	return _BLSKeystore.Contract.GetPubKey(&_BLSKeystore.CallOpts)
}

// GetPubKey is a free data retrieval call binding the contract method 0x57c8a80b.
//
// Solidity: function GetPubKey() view returns(bytes)
func (_BLSKeystore *BLSKeystoreCallerSession) GetPubKey() ([]byte, error) {
	return _BLSKeystore.Contract.GetPubKey(&_BLSKeystore.CallOpts)
}

// PubKey is a free data retrieval call binding the contract method 0xd5bf8e5f.
//
// Solidity: function PubKey() view returns(bytes)
func (_BLSKeystore *BLSKeystoreCaller) PubKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _BLSKeystore.contract.Call(opts, &out, "PubKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// PubKey is a free data retrieval call binding the contract method 0xd5bf8e5f.
//
// Solidity: function PubKey() view returns(bytes)
func (_BLSKeystore *BLSKeystoreSession) PubKey() ([]byte, error) {
	return _BLSKeystore.Contract.PubKey(&_BLSKeystore.CallOpts)
}

// PubKey is a free data retrieval call binding the contract method 0xd5bf8e5f.
//
// Solidity: function PubKey() view returns(bytes)
func (_BLSKeystore *BLSKeystoreCallerSession) PubKey() ([]byte, error) {
	return _BLSKeystore.Contract.PubKey(&_BLSKeystore.CallOpts)
}
