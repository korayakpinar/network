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

// OperatorsMetaData contains all meta data concerning the Operators contract.
var OperatorsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"GetOperatorCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"}],\"name\":\"GetOperatorIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"}],\"name\":\"IsRegistered\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"Operators\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"blsPubKeyCID\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"}],\"name\":\"getBLSPubKeyCID\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getBLSPubKeyCIDByIndex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"}],\"name\":\"hasSubmittedBLSKey\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"hasSubmittedKey\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registerOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"blsPubKeyCID\",\"type\":\"string\"}],\"name\":\"submitBlsKeyCID\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// OperatorsABI is the input ABI used to generate the binding from.
// Deprecated: Use OperatorsMetaData.ABI instead.
var OperatorsABI = OperatorsMetaData.ABI

// Operators is an auto generated Go binding around an Ethereum contract.
type Operators struct {
	OperatorsCaller     // Read-only binding to the contract
	OperatorsTransactor // Write-only binding to the contract
	OperatorsFilterer   // Log filterer for contract events
}

// OperatorsCaller is an auto generated read-only Go binding around an Ethereum contract.
type OperatorsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OperatorsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OperatorsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OperatorsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OperatorsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OperatorsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OperatorsSession struct {
	Contract     *Operators        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OperatorsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OperatorsCallerSession struct {
	Contract *OperatorsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// OperatorsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OperatorsTransactorSession struct {
	Contract     *OperatorsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// OperatorsRaw is an auto generated low-level Go binding around an Ethereum contract.
type OperatorsRaw struct {
	Contract *Operators // Generic contract binding to access the raw methods on
}

// OperatorsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OperatorsCallerRaw struct {
	Contract *OperatorsCaller // Generic read-only contract binding to access the raw methods on
}

// OperatorsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OperatorsTransactorRaw struct {
	Contract *OperatorsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOperators creates a new instance of Operators, bound to a specific deployed contract.
func NewOperators(address common.Address, backend bind.ContractBackend) (*Operators, error) {
	contract, err := bindOperators(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Operators{OperatorsCaller: OperatorsCaller{contract: contract}, OperatorsTransactor: OperatorsTransactor{contract: contract}, OperatorsFilterer: OperatorsFilterer{contract: contract}}, nil
}

// NewOperatorsCaller creates a new read-only instance of Operators, bound to a specific deployed contract.
func NewOperatorsCaller(address common.Address, caller bind.ContractCaller) (*OperatorsCaller, error) {
	contract, err := bindOperators(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OperatorsCaller{contract: contract}, nil
}

// NewOperatorsTransactor creates a new write-only instance of Operators, bound to a specific deployed contract.
func NewOperatorsTransactor(address common.Address, transactor bind.ContractTransactor) (*OperatorsTransactor, error) {
	contract, err := bindOperators(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OperatorsTransactor{contract: contract}, nil
}

// NewOperatorsFilterer creates a new log filterer instance of Operators, bound to a specific deployed contract.
func NewOperatorsFilterer(address common.Address, filterer bind.ContractFilterer) (*OperatorsFilterer, error) {
	contract, err := bindOperators(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OperatorsFilterer{contract: contract}, nil
}

// bindOperators binds a generic wrapper to an already deployed contract.
func bindOperators(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OperatorsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Operators *OperatorsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Operators.Contract.OperatorsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Operators *OperatorsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Operators.Contract.OperatorsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Operators *OperatorsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Operators.Contract.OperatorsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Operators *OperatorsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Operators.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Operators *OperatorsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Operators.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Operators *OperatorsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Operators.Contract.contract.Transact(opts, method, params...)
}

// GetOperatorCount is a free data retrieval call binding the contract method 0x9a5c350d.
//
// Solidity: function GetOperatorCount() view returns(uint256)
func (_Operators *OperatorsCaller) GetOperatorCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Operators.contract.Call(opts, &out, "GetOperatorCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOperatorCount is a free data retrieval call binding the contract method 0x9a5c350d.
//
// Solidity: function GetOperatorCount() view returns(uint256)
func (_Operators *OperatorsSession) GetOperatorCount() (*big.Int, error) {
	return _Operators.Contract.GetOperatorCount(&_Operators.CallOpts)
}

// GetOperatorCount is a free data retrieval call binding the contract method 0x9a5c350d.
//
// Solidity: function GetOperatorCount() view returns(uint256)
func (_Operators *OperatorsCallerSession) GetOperatorCount() (*big.Int, error) {
	return _Operators.Contract.GetOperatorCount(&_Operators.CallOpts)
}

// GetOperatorIndex is a free data retrieval call binding the contract method 0x90748cc8.
//
// Solidity: function GetOperatorIndex(address operatorAddress) view returns(uint256)
func (_Operators *OperatorsCaller) GetOperatorIndex(opts *bind.CallOpts, operatorAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Operators.contract.Call(opts, &out, "GetOperatorIndex", operatorAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOperatorIndex is a free data retrieval call binding the contract method 0x90748cc8.
//
// Solidity: function GetOperatorIndex(address operatorAddress) view returns(uint256)
func (_Operators *OperatorsSession) GetOperatorIndex(operatorAddress common.Address) (*big.Int, error) {
	return _Operators.Contract.GetOperatorIndex(&_Operators.CallOpts, operatorAddress)
}

// GetOperatorIndex is a free data retrieval call binding the contract method 0x90748cc8.
//
// Solidity: function GetOperatorIndex(address operatorAddress) view returns(uint256)
func (_Operators *OperatorsCallerSession) GetOperatorIndex(operatorAddress common.Address) (*big.Int, error) {
	return _Operators.Contract.GetOperatorIndex(&_Operators.CallOpts, operatorAddress)
}

// IsRegistered is a free data retrieval call binding the contract method 0xc224122d.
//
// Solidity: function IsRegistered(address operatorAddress) view returns(bool)
func (_Operators *OperatorsCaller) IsRegistered(opts *bind.CallOpts, operatorAddress common.Address) (bool, error) {
	var out []interface{}
	err := _Operators.contract.Call(opts, &out, "IsRegistered", operatorAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRegistered is a free data retrieval call binding the contract method 0xc224122d.
//
// Solidity: function IsRegistered(address operatorAddress) view returns(bool)
func (_Operators *OperatorsSession) IsRegistered(operatorAddress common.Address) (bool, error) {
	return _Operators.Contract.IsRegistered(&_Operators.CallOpts, operatorAddress)
}

// IsRegistered is a free data retrieval call binding the contract method 0xc224122d.
//
// Solidity: function IsRegistered(address operatorAddress) view returns(bool)
func (_Operators *OperatorsCallerSession) IsRegistered(operatorAddress common.Address) (bool, error) {
	return _Operators.Contract.IsRegistered(&_Operators.CallOpts, operatorAddress)
}

// Operators is a free data retrieval call binding the contract method 0x6acbb296.
//
// Solidity: function Operators(uint256 ) view returns(address operator, string blsPubKeyCID)
func (_Operators *OperatorsCaller) Operators(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Operator     common.Address
	BlsPubKeyCID string
}, error) {
	var out []interface{}
	err := _Operators.contract.Call(opts, &out, "Operators", arg0)

	outstruct := new(struct {
		Operator     common.Address
		BlsPubKeyCID string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Operator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.BlsPubKeyCID = *abi.ConvertType(out[1], new(string)).(*string)

	return *outstruct, err

}

// Operators is a free data retrieval call binding the contract method 0x6acbb296.
//
// Solidity: function Operators(uint256 ) view returns(address operator, string blsPubKeyCID)
func (_Operators *OperatorsSession) Operators(arg0 *big.Int) (struct {
	Operator     common.Address
	BlsPubKeyCID string
}, error) {
	return _Operators.Contract.Operators(&_Operators.CallOpts, arg0)
}

// Operators is a free data retrieval call binding the contract method 0x6acbb296.
//
// Solidity: function Operators(uint256 ) view returns(address operator, string blsPubKeyCID)
func (_Operators *OperatorsCallerSession) Operators(arg0 *big.Int) (struct {
	Operator     common.Address
	BlsPubKeyCID string
}, error) {
	return _Operators.Contract.Operators(&_Operators.CallOpts, arg0)
}

// GetBLSPubKeyCID is a free data retrieval call binding the contract method 0x5e83ae33.
//
// Solidity: function getBLSPubKeyCID(address operatorAddress) view returns(string)
func (_Operators *OperatorsCaller) GetBLSPubKeyCID(opts *bind.CallOpts, operatorAddress common.Address) (string, error) {
	var out []interface{}
	err := _Operators.contract.Call(opts, &out, "getBLSPubKeyCID", operatorAddress)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetBLSPubKeyCID is a free data retrieval call binding the contract method 0x5e83ae33.
//
// Solidity: function getBLSPubKeyCID(address operatorAddress) view returns(string)
func (_Operators *OperatorsSession) GetBLSPubKeyCID(operatorAddress common.Address) (string, error) {
	return _Operators.Contract.GetBLSPubKeyCID(&_Operators.CallOpts, operatorAddress)
}

// GetBLSPubKeyCID is a free data retrieval call binding the contract method 0x5e83ae33.
//
// Solidity: function getBLSPubKeyCID(address operatorAddress) view returns(string)
func (_Operators *OperatorsCallerSession) GetBLSPubKeyCID(operatorAddress common.Address) (string, error) {
	return _Operators.Contract.GetBLSPubKeyCID(&_Operators.CallOpts, operatorAddress)
}

// GetBLSPubKeyCIDByIndex is a free data retrieval call binding the contract method 0xc292cfb5.
//
// Solidity: function getBLSPubKeyCIDByIndex(uint256 index) view returns(string)
func (_Operators *OperatorsCaller) GetBLSPubKeyCIDByIndex(opts *bind.CallOpts, index *big.Int) (string, error) {
	var out []interface{}
	err := _Operators.contract.Call(opts, &out, "getBLSPubKeyCIDByIndex", index)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetBLSPubKeyCIDByIndex is a free data retrieval call binding the contract method 0xc292cfb5.
//
// Solidity: function getBLSPubKeyCIDByIndex(uint256 index) view returns(string)
func (_Operators *OperatorsSession) GetBLSPubKeyCIDByIndex(index *big.Int) (string, error) {
	return _Operators.Contract.GetBLSPubKeyCIDByIndex(&_Operators.CallOpts, index)
}

// GetBLSPubKeyCIDByIndex is a free data retrieval call binding the contract method 0xc292cfb5.
//
// Solidity: function getBLSPubKeyCIDByIndex(uint256 index) view returns(string)
func (_Operators *OperatorsCallerSession) GetBLSPubKeyCIDByIndex(index *big.Int) (string, error) {
	return _Operators.Contract.GetBLSPubKeyCIDByIndex(&_Operators.CallOpts, index)
}

// HasSubmittedBLSKey is a free data retrieval call binding the contract method 0xcd3613b5.
//
// Solidity: function hasSubmittedBLSKey(address operatorAddress) view returns(bool)
func (_Operators *OperatorsCaller) HasSubmittedBLSKey(opts *bind.CallOpts, operatorAddress common.Address) (bool, error) {
	var out []interface{}
	err := _Operators.contract.Call(opts, &out, "hasSubmittedBLSKey", operatorAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasSubmittedBLSKey is a free data retrieval call binding the contract method 0xcd3613b5.
//
// Solidity: function hasSubmittedBLSKey(address operatorAddress) view returns(bool)
func (_Operators *OperatorsSession) HasSubmittedBLSKey(operatorAddress common.Address) (bool, error) {
	return _Operators.Contract.HasSubmittedBLSKey(&_Operators.CallOpts, operatorAddress)
}

// HasSubmittedBLSKey is a free data retrieval call binding the contract method 0xcd3613b5.
//
// Solidity: function hasSubmittedBLSKey(address operatorAddress) view returns(bool)
func (_Operators *OperatorsCallerSession) HasSubmittedBLSKey(operatorAddress common.Address) (bool, error) {
	return _Operators.Contract.HasSubmittedBLSKey(&_Operators.CallOpts, operatorAddress)
}

// HasSubmittedKey is a free data retrieval call binding the contract method 0xabd58956.
//
// Solidity: function hasSubmittedKey(address ) view returns(bool)
func (_Operators *OperatorsCaller) HasSubmittedKey(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Operators.contract.Call(opts, &out, "hasSubmittedKey", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasSubmittedKey is a free data retrieval call binding the contract method 0xabd58956.
//
// Solidity: function hasSubmittedKey(address ) view returns(bool)
func (_Operators *OperatorsSession) HasSubmittedKey(arg0 common.Address) (bool, error) {
	return _Operators.Contract.HasSubmittedKey(&_Operators.CallOpts, arg0)
}

// HasSubmittedKey is a free data retrieval call binding the contract method 0xabd58956.
//
// Solidity: function hasSubmittedKey(address ) view returns(bool)
func (_Operators *OperatorsCallerSession) HasSubmittedKey(arg0 common.Address) (bool, error) {
	return _Operators.Contract.HasSubmittedKey(&_Operators.CallOpts, arg0)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x2acde098.
//
// Solidity: function registerOperator() returns()
func (_Operators *OperatorsTransactor) RegisterOperator(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Operators.contract.Transact(opts, "registerOperator")
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x2acde098.
//
// Solidity: function registerOperator() returns()
func (_Operators *OperatorsSession) RegisterOperator() (*types.Transaction, error) {
	return _Operators.Contract.RegisterOperator(&_Operators.TransactOpts)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x2acde098.
//
// Solidity: function registerOperator() returns()
func (_Operators *OperatorsTransactorSession) RegisterOperator() (*types.Transaction, error) {
	return _Operators.Contract.RegisterOperator(&_Operators.TransactOpts)
}

// SubmitBlsKeyCID is a paid mutator transaction binding the contract method 0xb3706e1f.
//
// Solidity: function submitBlsKeyCID(string blsPubKeyCID) returns()
func (_Operators *OperatorsTransactor) SubmitBlsKeyCID(opts *bind.TransactOpts, blsPubKeyCID string) (*types.Transaction, error) {
	return _Operators.contract.Transact(opts, "submitBlsKeyCID", blsPubKeyCID)
}

// SubmitBlsKeyCID is a paid mutator transaction binding the contract method 0xb3706e1f.
//
// Solidity: function submitBlsKeyCID(string blsPubKeyCID) returns()
func (_Operators *OperatorsSession) SubmitBlsKeyCID(blsPubKeyCID string) (*types.Transaction, error) {
	return _Operators.Contract.SubmitBlsKeyCID(&_Operators.TransactOpts, blsPubKeyCID)
}

// SubmitBlsKeyCID is a paid mutator transaction binding the contract method 0xb3706e1f.
//
// Solidity: function submitBlsKeyCID(string blsPubKeyCID) returns()
func (_Operators *OperatorsTransactorSession) SubmitBlsKeyCID(blsPubKeyCID string) (*types.Transaction, error) {
	return _Operators.Contract.SubmitBlsKeyCID(&_Operators.TransactOpts, blsPubKeyCID)
}
