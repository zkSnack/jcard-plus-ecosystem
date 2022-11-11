// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package state

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

// StateMetaData contains all meta data concerning the State contract.
var StateMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"blockN\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"}],\"name\":\"StateUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getState\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"blockN\",\"type\":\"uint64\"}],\"name\":\"getStateDataByBlock\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getStateDataById\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"getStateDataByTime\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"}],\"name\":\"getTransitionInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"identities\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"BlockN\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"BlockTimestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"State\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIVerifier\",\"name\":\"_verifierContractAddr\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newVerifier\",\"type\":\"address\"}],\"name\":\"setVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldState\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newState\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isOldStateGenesis\",\"type\":\"bool\"},{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"}],\"name\":\"transitState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"transitions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"replacedAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"replacedAtBlock\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"createdAtBlock\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"replacedBy\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifier\",\"outputs\":[{\"internalType\":\"contractIVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// StateABI is the input ABI used to generate the binding from.
// Deprecated: Use StateMetaData.ABI instead.
var StateABI = StateMetaData.ABI

// State is an auto generated Go binding around an Ethereum contract.
type State struct {
	StateCaller     // Read-only binding to the contract
	StateTransactor // Write-only binding to the contract
	StateFilterer   // Log filterer for contract events
}

// StateCaller is an auto generated read-only Go binding around an Ethereum contract.
type StateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StateSession struct {
	Contract     *State            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StateCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StateCallerSession struct {
	Contract *StateCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// StateTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StateTransactorSession struct {
	Contract     *StateTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StateRaw is an auto generated low-level Go binding around an Ethereum contract.
type StateRaw struct {
	Contract *State // Generic contract binding to access the raw methods on
}

// StateCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StateCallerRaw struct {
	Contract *StateCaller // Generic read-only contract binding to access the raw methods on
}

// StateTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StateTransactorRaw struct {
	Contract *StateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewState creates a new instance of State, bound to a specific deployed contract.
func NewState(address common.Address, backend bind.ContractBackend) (*State, error) {
	contract, err := bindState(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &State{StateCaller: StateCaller{contract: contract}, StateTransactor: StateTransactor{contract: contract}, StateFilterer: StateFilterer{contract: contract}}, nil
}

// NewStateCaller creates a new read-only instance of State, bound to a specific deployed contract.
func NewStateCaller(address common.Address, caller bind.ContractCaller) (*StateCaller, error) {
	contract, err := bindState(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StateCaller{contract: contract}, nil
}

// NewStateTransactor creates a new write-only instance of State, bound to a specific deployed contract.
func NewStateTransactor(address common.Address, transactor bind.ContractTransactor) (*StateTransactor, error) {
	contract, err := bindState(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StateTransactor{contract: contract}, nil
}

// NewStateFilterer creates a new log filterer instance of State, bound to a specific deployed contract.
func NewStateFilterer(address common.Address, filterer bind.ContractFilterer) (*StateFilterer, error) {
	contract, err := bindState(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StateFilterer{contract: contract}, nil
}

// bindState binds a generic wrapper to an already deployed contract.
func bindState(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StateABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_State *StateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _State.Contract.StateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_State *StateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _State.Contract.StateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_State *StateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _State.Contract.StateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_State *StateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _State.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_State *StateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _State.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_State *StateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _State.Contract.contract.Transact(opts, method, params...)
}

// GetState is a free data retrieval call binding the contract method 0x44c9af28.
//
// Solidity: function getState(uint256 id) view returns(uint256)
func (_State *StateCaller) GetState(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getState", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetState is a free data retrieval call binding the contract method 0x44c9af28.
//
// Solidity: function getState(uint256 id) view returns(uint256)
func (_State *StateSession) GetState(id *big.Int) (*big.Int, error) {
	return _State.Contract.GetState(&_State.CallOpts, id)
}

// GetState is a free data retrieval call binding the contract method 0x44c9af28.
//
// Solidity: function getState(uint256 id) view returns(uint256)
func (_State *StateCallerSession) GetState(id *big.Int) (*big.Int, error) {
	return _State.Contract.GetState(&_State.CallOpts, id)
}

// GetStateDataByBlock is a free data retrieval call binding the contract method 0xd8dcd971.
//
// Solidity: function getStateDataByBlock(uint256 id, uint64 blockN) view returns(uint64, uint64, uint256)
func (_State *StateCaller) GetStateDataByBlock(opts *bind.CallOpts, id *big.Int, blockN uint64) (uint64, uint64, *big.Int, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getStateDataByBlock", id, blockN)

	if err != nil {
		return *new(uint64), *new(uint64), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	out1 := *abi.ConvertType(out[1], new(uint64)).(*uint64)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

// GetStateDataByBlock is a free data retrieval call binding the contract method 0xd8dcd971.
//
// Solidity: function getStateDataByBlock(uint256 id, uint64 blockN) view returns(uint64, uint64, uint256)
func (_State *StateSession) GetStateDataByBlock(id *big.Int, blockN uint64) (uint64, uint64, *big.Int, error) {
	return _State.Contract.GetStateDataByBlock(&_State.CallOpts, id, blockN)
}

// GetStateDataByBlock is a free data retrieval call binding the contract method 0xd8dcd971.
//
// Solidity: function getStateDataByBlock(uint256 id, uint64 blockN) view returns(uint64, uint64, uint256)
func (_State *StateCallerSession) GetStateDataByBlock(id *big.Int, blockN uint64) (uint64, uint64, *big.Int, error) {
	return _State.Contract.GetStateDataByBlock(&_State.CallOpts, id, blockN)
}

// GetStateDataById is a free data retrieval call binding the contract method 0xc8d1e53e.
//
// Solidity: function getStateDataById(uint256 id) view returns(uint64, uint64, uint256)
func (_State *StateCaller) GetStateDataById(opts *bind.CallOpts, id *big.Int) (uint64, uint64, *big.Int, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getStateDataById", id)

	if err != nil {
		return *new(uint64), *new(uint64), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	out1 := *abi.ConvertType(out[1], new(uint64)).(*uint64)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

// GetStateDataById is a free data retrieval call binding the contract method 0xc8d1e53e.
//
// Solidity: function getStateDataById(uint256 id) view returns(uint64, uint64, uint256)
func (_State *StateSession) GetStateDataById(id *big.Int) (uint64, uint64, *big.Int, error) {
	return _State.Contract.GetStateDataById(&_State.CallOpts, id)
}

// GetStateDataById is a free data retrieval call binding the contract method 0xc8d1e53e.
//
// Solidity: function getStateDataById(uint256 id) view returns(uint64, uint64, uint256)
func (_State *StateCallerSession) GetStateDataById(id *big.Int) (uint64, uint64, *big.Int, error) {
	return _State.Contract.GetStateDataById(&_State.CallOpts, id)
}

// GetStateDataByTime is a free data retrieval call binding the contract method 0x0281bec2.
//
// Solidity: function getStateDataByTime(uint256 id, uint64 timestamp) view returns(uint64, uint64, uint256)
func (_State *StateCaller) GetStateDataByTime(opts *bind.CallOpts, id *big.Int, timestamp uint64) (uint64, uint64, *big.Int, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getStateDataByTime", id, timestamp)

	if err != nil {
		return *new(uint64), *new(uint64), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	out1 := *abi.ConvertType(out[1], new(uint64)).(*uint64)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

// GetStateDataByTime is a free data retrieval call binding the contract method 0x0281bec2.
//
// Solidity: function getStateDataByTime(uint256 id, uint64 timestamp) view returns(uint64, uint64, uint256)
func (_State *StateSession) GetStateDataByTime(id *big.Int, timestamp uint64) (uint64, uint64, *big.Int, error) {
	return _State.Contract.GetStateDataByTime(&_State.CallOpts, id, timestamp)
}

// GetStateDataByTime is a free data retrieval call binding the contract method 0x0281bec2.
//
// Solidity: function getStateDataByTime(uint256 id, uint64 timestamp) view returns(uint64, uint64, uint256)
func (_State *StateCallerSession) GetStateDataByTime(id *big.Int, timestamp uint64) (uint64, uint64, *big.Int, error) {
	return _State.Contract.GetStateDataByTime(&_State.CallOpts, id, timestamp)
}

// GetTransitionInfo is a free data retrieval call binding the contract method 0xbb795715.
//
// Solidity: function getTransitionInfo(uint256 state) view returns(uint256, uint256, uint64, uint64, uint256, uint256)
func (_State *StateCaller) GetTransitionInfo(opts *bind.CallOpts, state *big.Int) (*big.Int, *big.Int, uint64, uint64, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getTransitionInfo", state)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(uint64), *new(uint64), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(uint64)).(*uint64)
	out3 := *abi.ConvertType(out[3], new(uint64)).(*uint64)
	out4 := *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	out5 := *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return out0, out1, out2, out3, out4, out5, err

}

// GetTransitionInfo is a free data retrieval call binding the contract method 0xbb795715.
//
// Solidity: function getTransitionInfo(uint256 state) view returns(uint256, uint256, uint64, uint64, uint256, uint256)
func (_State *StateSession) GetTransitionInfo(state *big.Int) (*big.Int, *big.Int, uint64, uint64, *big.Int, *big.Int, error) {
	return _State.Contract.GetTransitionInfo(&_State.CallOpts, state)
}

// GetTransitionInfo is a free data retrieval call binding the contract method 0xbb795715.
//
// Solidity: function getTransitionInfo(uint256 state) view returns(uint256, uint256, uint64, uint64, uint256, uint256)
func (_State *StateCallerSession) GetTransitionInfo(state *big.Int) (*big.Int, *big.Int, uint64, uint64, *big.Int, *big.Int, error) {
	return _State.Contract.GetTransitionInfo(&_State.CallOpts, state)
}

// Identities is a free data retrieval call binding the contract method 0xe20490b5.
//
// Solidity: function identities(uint256 , uint256 ) view returns(uint64 BlockN, uint64 BlockTimestamp, uint256 State)
func (_State *StateCaller) Identities(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (struct {
	BlockN         uint64
	BlockTimestamp uint64
	State          *big.Int
}, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "identities", arg0, arg1)

	outstruct := new(struct {
		BlockN         uint64
		BlockTimestamp uint64
		State          *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BlockN = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.BlockTimestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.State = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Identities is a free data retrieval call binding the contract method 0xe20490b5.
//
// Solidity: function identities(uint256 , uint256 ) view returns(uint64 BlockN, uint64 BlockTimestamp, uint256 State)
func (_State *StateSession) Identities(arg0 *big.Int, arg1 *big.Int) (struct {
	BlockN         uint64
	BlockTimestamp uint64
	State          *big.Int
}, error) {
	return _State.Contract.Identities(&_State.CallOpts, arg0, arg1)
}

// Identities is a free data retrieval call binding the contract method 0xe20490b5.
//
// Solidity: function identities(uint256 , uint256 ) view returns(uint64 BlockN, uint64 BlockTimestamp, uint256 State)
func (_State *StateCallerSession) Identities(arg0 *big.Int, arg1 *big.Int) (struct {
	BlockN         uint64
	BlockTimestamp uint64
	State          *big.Int
}, error) {
	return _State.Contract.Identities(&_State.CallOpts, arg0, arg1)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_State *StateCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_State *StateSession) Owner() (common.Address, error) {
	return _State.Contract.Owner(&_State.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_State *StateCallerSession) Owner() (common.Address, error) {
	return _State.Contract.Owner(&_State.CallOpts)
}

// Transitions is a free data retrieval call binding the contract method 0x683ace65.
//
// Solidity: function transitions(uint256 ) view returns(uint256 replacedAtTimestamp, uint256 createdAtTimestamp, uint64 replacedAtBlock, uint64 createdAtBlock, uint256 replacedBy, uint256 id)
func (_State *StateCaller) Transitions(opts *bind.CallOpts, arg0 *big.Int) (struct {
	ReplacedAtTimestamp *big.Int
	CreatedAtTimestamp  *big.Int
	ReplacedAtBlock     uint64
	CreatedAtBlock      uint64
	ReplacedBy          *big.Int
	Id                  *big.Int
}, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "transitions", arg0)

	outstruct := new(struct {
		ReplacedAtTimestamp *big.Int
		CreatedAtTimestamp  *big.Int
		ReplacedAtBlock     uint64
		CreatedAtBlock      uint64
		ReplacedBy          *big.Int
		Id                  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ReplacedAtTimestamp = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.CreatedAtTimestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ReplacedAtBlock = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.CreatedAtBlock = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.ReplacedBy = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Id = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Transitions is a free data retrieval call binding the contract method 0x683ace65.
//
// Solidity: function transitions(uint256 ) view returns(uint256 replacedAtTimestamp, uint256 createdAtTimestamp, uint64 replacedAtBlock, uint64 createdAtBlock, uint256 replacedBy, uint256 id)
func (_State *StateSession) Transitions(arg0 *big.Int) (struct {
	ReplacedAtTimestamp *big.Int
	CreatedAtTimestamp  *big.Int
	ReplacedAtBlock     uint64
	CreatedAtBlock      uint64
	ReplacedBy          *big.Int
	Id                  *big.Int
}, error) {
	return _State.Contract.Transitions(&_State.CallOpts, arg0)
}

// Transitions is a free data retrieval call binding the contract method 0x683ace65.
//
// Solidity: function transitions(uint256 ) view returns(uint256 replacedAtTimestamp, uint256 createdAtTimestamp, uint64 replacedAtBlock, uint64 createdAtBlock, uint256 replacedBy, uint256 id)
func (_State *StateCallerSession) Transitions(arg0 *big.Int) (struct {
	ReplacedAtTimestamp *big.Int
	CreatedAtTimestamp  *big.Int
	ReplacedAtBlock     uint64
	CreatedAtBlock      uint64
	ReplacedBy          *big.Int
	Id                  *big.Int
}, error) {
	return _State.Contract.Transitions(&_State.CallOpts, arg0)
}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_State *StateCaller) Verifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "verifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_State *StateSession) Verifier() (common.Address, error) {
	return _State.Contract.Verifier(&_State.CallOpts)
}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_State *StateCallerSession) Verifier() (common.Address, error) {
	return _State.Contract.Verifier(&_State.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _verifierContractAddr) returns()
func (_State *StateTransactor) Initialize(opts *bind.TransactOpts, _verifierContractAddr common.Address) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "initialize", _verifierContractAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _verifierContractAddr) returns()
func (_State *StateSession) Initialize(_verifierContractAddr common.Address) (*types.Transaction, error) {
	return _State.Contract.Initialize(&_State.TransactOpts, _verifierContractAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _verifierContractAddr) returns()
func (_State *StateTransactorSession) Initialize(_verifierContractAddr common.Address) (*types.Transaction, error) {
	return _State.Contract.Initialize(&_State.TransactOpts, _verifierContractAddr)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_State *StateTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_State *StateSession) RenounceOwnership() (*types.Transaction, error) {
	return _State.Contract.RenounceOwnership(&_State.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_State *StateTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _State.Contract.RenounceOwnership(&_State.TransactOpts)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x5437988d.
//
// Solidity: function setVerifier(address newVerifier) returns()
func (_State *StateTransactor) SetVerifier(opts *bind.TransactOpts, newVerifier common.Address) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "setVerifier", newVerifier)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x5437988d.
//
// Solidity: function setVerifier(address newVerifier) returns()
func (_State *StateSession) SetVerifier(newVerifier common.Address) (*types.Transaction, error) {
	return _State.Contract.SetVerifier(&_State.TransactOpts, newVerifier)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x5437988d.
//
// Solidity: function setVerifier(address newVerifier) returns()
func (_State *StateTransactorSession) SetVerifier(newVerifier common.Address) (*types.Transaction, error) {
	return _State.Contract.SetVerifier(&_State.TransactOpts, newVerifier)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_State *StateTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_State *StateSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _State.Contract.TransferOwnership(&_State.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_State *StateTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _State.Contract.TransferOwnership(&_State.TransactOpts, newOwner)
}

// TransitState is a paid mutator transaction binding the contract method 0x28f88a65.
//
// Solidity: function transitState(uint256 id, uint256 oldState, uint256 newState, bool isOldStateGenesis, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_State *StateTransactor) TransitState(opts *bind.TransactOpts, id *big.Int, oldState *big.Int, newState *big.Int, isOldStateGenesis bool, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "transitState", id, oldState, newState, isOldStateGenesis, a, b, c)
}

// TransitState is a paid mutator transaction binding the contract method 0x28f88a65.
//
// Solidity: function transitState(uint256 id, uint256 oldState, uint256 newState, bool isOldStateGenesis, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_State *StateSession) TransitState(id *big.Int, oldState *big.Int, newState *big.Int, isOldStateGenesis bool, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _State.Contract.TransitState(&_State.TransactOpts, id, oldState, newState, isOldStateGenesis, a, b, c)
}

// TransitState is a paid mutator transaction binding the contract method 0x28f88a65.
//
// Solidity: function transitState(uint256 id, uint256 oldState, uint256 newState, bool isOldStateGenesis, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_State *StateTransactorSession) TransitState(id *big.Int, oldState *big.Int, newState *big.Int, isOldStateGenesis bool, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _State.Contract.TransitState(&_State.TransactOpts, id, oldState, newState, isOldStateGenesis, a, b, c)
}

// StateInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the State contract.
type StateInitializedIterator struct {
	Event *StateInitialized // Event containing the contract specifics and raw log

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
func (it *StateInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateInitialized)
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
		it.Event = new(StateInitialized)
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
func (it *StateInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateInitialized represents a Initialized event raised by the State contract.
type StateInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_State *StateFilterer) FilterInitialized(opts *bind.FilterOpts) (*StateInitializedIterator, error) {

	logs, sub, err := _State.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &StateInitializedIterator{contract: _State.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_State *StateFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *StateInitialized) (event.Subscription, error) {

	logs, sub, err := _State.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateInitialized)
				if err := _State.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_State *StateFilterer) ParseInitialized(log types.Log) (*StateInitialized, error) {
	event := new(StateInitialized)
	if err := _State.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StateOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the State contract.
type StateOwnershipTransferredIterator struct {
	Event *StateOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *StateOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateOwnershipTransferred)
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
		it.Event = new(StateOwnershipTransferred)
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
func (it *StateOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateOwnershipTransferred represents a OwnershipTransferred event raised by the State contract.
type StateOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_State *StateFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StateOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _State.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &StateOwnershipTransferredIterator{contract: _State.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_State *StateFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StateOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _State.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateOwnershipTransferred)
				if err := _State.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_State *StateFilterer) ParseOwnershipTransferred(log types.Log) (*StateOwnershipTransferred, error) {
	event := new(StateOwnershipTransferred)
	if err := _State.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StateStateUpdatedIterator is returned from FilterStateUpdated and is used to iterate over the raw logs and unpacked data for StateUpdated events raised by the State contract.
type StateStateUpdatedIterator struct {
	Event *StateStateUpdated // Event containing the contract specifics and raw log

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
func (it *StateStateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateStateUpdated)
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
		it.Event = new(StateStateUpdated)
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
func (it *StateStateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateStateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateStateUpdated represents a StateUpdated event raised by the State contract.
type StateStateUpdated struct {
	Id        *big.Int
	BlockN    uint64
	Timestamp uint64
	State     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStateUpdated is a free log retrieval operation binding the contract event 0x81c6f328b24014ef550c34a433275b52f3a8a0f32aa871adec069ab526a02390.
//
// Solidity: event StateUpdated(uint256 id, uint64 blockN, uint64 timestamp, uint256 state)
func (_State *StateFilterer) FilterStateUpdated(opts *bind.FilterOpts) (*StateStateUpdatedIterator, error) {

	logs, sub, err := _State.contract.FilterLogs(opts, "StateUpdated")
	if err != nil {
		return nil, err
	}
	return &StateStateUpdatedIterator{contract: _State.contract, event: "StateUpdated", logs: logs, sub: sub}, nil
}

// WatchStateUpdated is a free log subscription operation binding the contract event 0x81c6f328b24014ef550c34a433275b52f3a8a0f32aa871adec069ab526a02390.
//
// Solidity: event StateUpdated(uint256 id, uint64 blockN, uint64 timestamp, uint256 state)
func (_State *StateFilterer) WatchStateUpdated(opts *bind.WatchOpts, sink chan<- *StateStateUpdated) (event.Subscription, error) {

	logs, sub, err := _State.contract.WatchLogs(opts, "StateUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateStateUpdated)
				if err := _State.contract.UnpackLog(event, "StateUpdated", log); err != nil {
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

// ParseStateUpdated is a log parse operation binding the contract event 0x81c6f328b24014ef550c34a433275b52f3a8a0f32aa871adec069ab526a02390.
//
// Solidity: event StateUpdated(uint256 id, uint64 blockN, uint64 timestamp, uint256 state)
func (_State *StateFilterer) ParseStateUpdated(log types.Log) (*StateStateUpdated, error) {
	event := new(StateStateUpdated)
	if err := _State.contract.UnpackLog(event, "StateUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
