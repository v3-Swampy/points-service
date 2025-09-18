// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// SwappiFactoryMetaData contains all meta data concerning the SwappiFactory contract.
var SwappiFactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_feeToSetter\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"PairCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"feeTo\",\"type\":\"address\"}],\"name\":\"feeToChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"feeToSetter\",\"type\":\"address\"}],\"name\":\"feeToSetterChanged\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"INIT_CODE_PAIR_HASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"allPairs\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"allPairsLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"}],\"name\":\"createPair\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeToSetter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"getPair\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_feeTo\",\"type\":\"address\"}],\"name\":\"setFeeTo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_feeToSetter\",\"type\":\"address\"}],\"name\":\"setFeeToSetter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// SwappiFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use SwappiFactoryMetaData.ABI instead.
var SwappiFactoryABI = SwappiFactoryMetaData.ABI

// SwappiFactory is an auto generated Go binding around an Ethereum contract.
type SwappiFactory struct {
	SwappiFactoryCaller     // Read-only binding to the contract
	SwappiFactoryTransactor // Write-only binding to the contract
	SwappiFactoryFilterer   // Log filterer for contract events
}

// SwappiFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type SwappiFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwappiFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SwappiFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwappiFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SwappiFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwappiFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SwappiFactorySession struct {
	Contract     *SwappiFactory    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SwappiFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SwappiFactoryCallerSession struct {
	Contract *SwappiFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// SwappiFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SwappiFactoryTransactorSession struct {
	Contract     *SwappiFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SwappiFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type SwappiFactoryRaw struct {
	Contract *SwappiFactory // Generic contract binding to access the raw methods on
}

// SwappiFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SwappiFactoryCallerRaw struct {
	Contract *SwappiFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// SwappiFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SwappiFactoryTransactorRaw struct {
	Contract *SwappiFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSwappiFactory creates a new instance of SwappiFactory, bound to a specific deployed contract.
func NewSwappiFactory(address common.Address, backend bind.ContractBackend) (*SwappiFactory, error) {
	contract, err := bindSwappiFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SwappiFactory{SwappiFactoryCaller: SwappiFactoryCaller{contract: contract}, SwappiFactoryTransactor: SwappiFactoryTransactor{contract: contract}, SwappiFactoryFilterer: SwappiFactoryFilterer{contract: contract}}, nil
}

// NewSwappiFactoryCaller creates a new read-only instance of SwappiFactory, bound to a specific deployed contract.
func NewSwappiFactoryCaller(address common.Address, caller bind.ContractCaller) (*SwappiFactoryCaller, error) {
	contract, err := bindSwappiFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SwappiFactoryCaller{contract: contract}, nil
}

// NewSwappiFactoryTransactor creates a new write-only instance of SwappiFactory, bound to a specific deployed contract.
func NewSwappiFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*SwappiFactoryTransactor, error) {
	contract, err := bindSwappiFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SwappiFactoryTransactor{contract: contract}, nil
}

// NewSwappiFactoryFilterer creates a new log filterer instance of SwappiFactory, bound to a specific deployed contract.
func NewSwappiFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*SwappiFactoryFilterer, error) {
	contract, err := bindSwappiFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SwappiFactoryFilterer{contract: contract}, nil
}

// bindSwappiFactory binds a generic wrapper to an already deployed contract.
func bindSwappiFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SwappiFactoryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwappiFactory *SwappiFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SwappiFactory.Contract.SwappiFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwappiFactory *SwappiFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwappiFactory.Contract.SwappiFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwappiFactory *SwappiFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwappiFactory.Contract.SwappiFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwappiFactory *SwappiFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SwappiFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwappiFactory *SwappiFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwappiFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwappiFactory *SwappiFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwappiFactory.Contract.contract.Transact(opts, method, params...)
}

// INITCODEPAIRHASH is a free data retrieval call binding the contract method 0x5855a25a.
//
// Solidity: function INIT_CODE_PAIR_HASH() view returns(bytes32)
func (_SwappiFactory *SwappiFactoryCaller) INITCODEPAIRHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SwappiFactory.contract.Call(opts, &out, "INIT_CODE_PAIR_HASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// INITCODEPAIRHASH is a free data retrieval call binding the contract method 0x5855a25a.
//
// Solidity: function INIT_CODE_PAIR_HASH() view returns(bytes32)
func (_SwappiFactory *SwappiFactorySession) INITCODEPAIRHASH() ([32]byte, error) {
	return _SwappiFactory.Contract.INITCODEPAIRHASH(&_SwappiFactory.CallOpts)
}

// INITCODEPAIRHASH is a free data retrieval call binding the contract method 0x5855a25a.
//
// Solidity: function INIT_CODE_PAIR_HASH() view returns(bytes32)
func (_SwappiFactory *SwappiFactoryCallerSession) INITCODEPAIRHASH() ([32]byte, error) {
	return _SwappiFactory.Contract.INITCODEPAIRHASH(&_SwappiFactory.CallOpts)
}

// AllPairs is a free data retrieval call binding the contract method 0x1e3dd18b.
//
// Solidity: function allPairs(uint256 ) view returns(address)
func (_SwappiFactory *SwappiFactoryCaller) AllPairs(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SwappiFactory.contract.Call(opts, &out, "allPairs", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AllPairs is a free data retrieval call binding the contract method 0x1e3dd18b.
//
// Solidity: function allPairs(uint256 ) view returns(address)
func (_SwappiFactory *SwappiFactorySession) AllPairs(arg0 *big.Int) (common.Address, error) {
	return _SwappiFactory.Contract.AllPairs(&_SwappiFactory.CallOpts, arg0)
}

// AllPairs is a free data retrieval call binding the contract method 0x1e3dd18b.
//
// Solidity: function allPairs(uint256 ) view returns(address)
func (_SwappiFactory *SwappiFactoryCallerSession) AllPairs(arg0 *big.Int) (common.Address, error) {
	return _SwappiFactory.Contract.AllPairs(&_SwappiFactory.CallOpts, arg0)
}

// AllPairsLength is a free data retrieval call binding the contract method 0x574f2ba3.
//
// Solidity: function allPairsLength() view returns(uint256)
func (_SwappiFactory *SwappiFactoryCaller) AllPairsLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SwappiFactory.contract.Call(opts, &out, "allPairsLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AllPairsLength is a free data retrieval call binding the contract method 0x574f2ba3.
//
// Solidity: function allPairsLength() view returns(uint256)
func (_SwappiFactory *SwappiFactorySession) AllPairsLength() (*big.Int, error) {
	return _SwappiFactory.Contract.AllPairsLength(&_SwappiFactory.CallOpts)
}

// AllPairsLength is a free data retrieval call binding the contract method 0x574f2ba3.
//
// Solidity: function allPairsLength() view returns(uint256)
func (_SwappiFactory *SwappiFactoryCallerSession) AllPairsLength() (*big.Int, error) {
	return _SwappiFactory.Contract.AllPairsLength(&_SwappiFactory.CallOpts)
}

// FeeTo is a free data retrieval call binding the contract method 0x017e7e58.
//
// Solidity: function feeTo() view returns(address)
func (_SwappiFactory *SwappiFactoryCaller) FeeTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SwappiFactory.contract.Call(opts, &out, "feeTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeTo is a free data retrieval call binding the contract method 0x017e7e58.
//
// Solidity: function feeTo() view returns(address)
func (_SwappiFactory *SwappiFactorySession) FeeTo() (common.Address, error) {
	return _SwappiFactory.Contract.FeeTo(&_SwappiFactory.CallOpts)
}

// FeeTo is a free data retrieval call binding the contract method 0x017e7e58.
//
// Solidity: function feeTo() view returns(address)
func (_SwappiFactory *SwappiFactoryCallerSession) FeeTo() (common.Address, error) {
	return _SwappiFactory.Contract.FeeTo(&_SwappiFactory.CallOpts)
}

// FeeToSetter is a free data retrieval call binding the contract method 0x094b7415.
//
// Solidity: function feeToSetter() view returns(address)
func (_SwappiFactory *SwappiFactoryCaller) FeeToSetter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SwappiFactory.contract.Call(opts, &out, "feeToSetter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeToSetter is a free data retrieval call binding the contract method 0x094b7415.
//
// Solidity: function feeToSetter() view returns(address)
func (_SwappiFactory *SwappiFactorySession) FeeToSetter() (common.Address, error) {
	return _SwappiFactory.Contract.FeeToSetter(&_SwappiFactory.CallOpts)
}

// FeeToSetter is a free data retrieval call binding the contract method 0x094b7415.
//
// Solidity: function feeToSetter() view returns(address)
func (_SwappiFactory *SwappiFactoryCallerSession) FeeToSetter() (common.Address, error) {
	return _SwappiFactory.Contract.FeeToSetter(&_SwappiFactory.CallOpts)
}

// GetPair is a free data retrieval call binding the contract method 0xe6a43905.
//
// Solidity: function getPair(address , address ) view returns(address)
func (_SwappiFactory *SwappiFactoryCaller) GetPair(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (common.Address, error) {
	var out []interface{}
	err := _SwappiFactory.contract.Call(opts, &out, "getPair", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPair is a free data retrieval call binding the contract method 0xe6a43905.
//
// Solidity: function getPair(address , address ) view returns(address)
func (_SwappiFactory *SwappiFactorySession) GetPair(arg0 common.Address, arg1 common.Address) (common.Address, error) {
	return _SwappiFactory.Contract.GetPair(&_SwappiFactory.CallOpts, arg0, arg1)
}

// GetPair is a free data retrieval call binding the contract method 0xe6a43905.
//
// Solidity: function getPair(address , address ) view returns(address)
func (_SwappiFactory *SwappiFactoryCallerSession) GetPair(arg0 common.Address, arg1 common.Address) (common.Address, error) {
	return _SwappiFactory.Contract.GetPair(&_SwappiFactory.CallOpts, arg0, arg1)
}

// CreatePair is a paid mutator transaction binding the contract method 0xc9c65396.
//
// Solidity: function createPair(address tokenA, address tokenB) returns(address pair)
func (_SwappiFactory *SwappiFactoryTransactor) CreatePair(opts *bind.TransactOpts, tokenA common.Address, tokenB common.Address) (*types.Transaction, error) {
	return _SwappiFactory.contract.Transact(opts, "createPair", tokenA, tokenB)
}

// CreatePair is a paid mutator transaction binding the contract method 0xc9c65396.
//
// Solidity: function createPair(address tokenA, address tokenB) returns(address pair)
func (_SwappiFactory *SwappiFactorySession) CreatePair(tokenA common.Address, tokenB common.Address) (*types.Transaction, error) {
	return _SwappiFactory.Contract.CreatePair(&_SwappiFactory.TransactOpts, tokenA, tokenB)
}

// CreatePair is a paid mutator transaction binding the contract method 0xc9c65396.
//
// Solidity: function createPair(address tokenA, address tokenB) returns(address pair)
func (_SwappiFactory *SwappiFactoryTransactorSession) CreatePair(tokenA common.Address, tokenB common.Address) (*types.Transaction, error) {
	return _SwappiFactory.Contract.CreatePair(&_SwappiFactory.TransactOpts, tokenA, tokenB)
}

// SetFeeTo is a paid mutator transaction binding the contract method 0xf46901ed.
//
// Solidity: function setFeeTo(address _feeTo) returns()
func (_SwappiFactory *SwappiFactoryTransactor) SetFeeTo(opts *bind.TransactOpts, _feeTo common.Address) (*types.Transaction, error) {
	return _SwappiFactory.contract.Transact(opts, "setFeeTo", _feeTo)
}

// SetFeeTo is a paid mutator transaction binding the contract method 0xf46901ed.
//
// Solidity: function setFeeTo(address _feeTo) returns()
func (_SwappiFactory *SwappiFactorySession) SetFeeTo(_feeTo common.Address) (*types.Transaction, error) {
	return _SwappiFactory.Contract.SetFeeTo(&_SwappiFactory.TransactOpts, _feeTo)
}

// SetFeeTo is a paid mutator transaction binding the contract method 0xf46901ed.
//
// Solidity: function setFeeTo(address _feeTo) returns()
func (_SwappiFactory *SwappiFactoryTransactorSession) SetFeeTo(_feeTo common.Address) (*types.Transaction, error) {
	return _SwappiFactory.Contract.SetFeeTo(&_SwappiFactory.TransactOpts, _feeTo)
}

// SetFeeToSetter is a paid mutator transaction binding the contract method 0xa2e74af6.
//
// Solidity: function setFeeToSetter(address _feeToSetter) returns()
func (_SwappiFactory *SwappiFactoryTransactor) SetFeeToSetter(opts *bind.TransactOpts, _feeToSetter common.Address) (*types.Transaction, error) {
	return _SwappiFactory.contract.Transact(opts, "setFeeToSetter", _feeToSetter)
}

// SetFeeToSetter is a paid mutator transaction binding the contract method 0xa2e74af6.
//
// Solidity: function setFeeToSetter(address _feeToSetter) returns()
func (_SwappiFactory *SwappiFactorySession) SetFeeToSetter(_feeToSetter common.Address) (*types.Transaction, error) {
	return _SwappiFactory.Contract.SetFeeToSetter(&_SwappiFactory.TransactOpts, _feeToSetter)
}

// SetFeeToSetter is a paid mutator transaction binding the contract method 0xa2e74af6.
//
// Solidity: function setFeeToSetter(address _feeToSetter) returns()
func (_SwappiFactory *SwappiFactoryTransactorSession) SetFeeToSetter(_feeToSetter common.Address) (*types.Transaction, error) {
	return _SwappiFactory.Contract.SetFeeToSetter(&_SwappiFactory.TransactOpts, _feeToSetter)
}

// SwappiFactoryPairCreatedIterator is returned from FilterPairCreated and is used to iterate over the raw logs and unpacked data for PairCreated events raised by the SwappiFactory contract.
type SwappiFactoryPairCreatedIterator struct {
	Event *SwappiFactoryPairCreated // Event containing the contract specifics and raw log

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
func (it *SwappiFactoryPairCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwappiFactoryPairCreated)
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
		it.Event = new(SwappiFactoryPairCreated)
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
func (it *SwappiFactoryPairCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwappiFactoryPairCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwappiFactoryPairCreated represents a PairCreated event raised by the SwappiFactory contract.
type SwappiFactoryPairCreated struct {
	Token0 common.Address
	Token1 common.Address
	Pair   common.Address
	Arg3   *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPairCreated is a free log retrieval operation binding the contract event 0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9.
//
// Solidity: event PairCreated(address indexed token0, address indexed token1, address pair, uint256 arg3)
func (_SwappiFactory *SwappiFactoryFilterer) FilterPairCreated(opts *bind.FilterOpts, token0 []common.Address, token1 []common.Address) (*SwappiFactoryPairCreatedIterator, error) {

	var token0Rule []interface{}
	for _, token0Item := range token0 {
		token0Rule = append(token0Rule, token0Item)
	}
	var token1Rule []interface{}
	for _, token1Item := range token1 {
		token1Rule = append(token1Rule, token1Item)
	}

	logs, sub, err := _SwappiFactory.contract.FilterLogs(opts, "PairCreated", token0Rule, token1Rule)
	if err != nil {
		return nil, err
	}
	return &SwappiFactoryPairCreatedIterator{contract: _SwappiFactory.contract, event: "PairCreated", logs: logs, sub: sub}, nil
}

// WatchPairCreated is a free log subscription operation binding the contract event 0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9.
//
// Solidity: event PairCreated(address indexed token0, address indexed token1, address pair, uint256 arg3)
func (_SwappiFactory *SwappiFactoryFilterer) WatchPairCreated(opts *bind.WatchOpts, sink chan<- *SwappiFactoryPairCreated, token0 []common.Address, token1 []common.Address) (event.Subscription, error) {

	var token0Rule []interface{}
	for _, token0Item := range token0 {
		token0Rule = append(token0Rule, token0Item)
	}
	var token1Rule []interface{}
	for _, token1Item := range token1 {
		token1Rule = append(token1Rule, token1Item)
	}

	logs, sub, err := _SwappiFactory.contract.WatchLogs(opts, "PairCreated", token0Rule, token1Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwappiFactoryPairCreated)
				if err := _SwappiFactory.contract.UnpackLog(event, "PairCreated", log); err != nil {
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

// ParsePairCreated is a log parse operation binding the contract event 0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9.
//
// Solidity: event PairCreated(address indexed token0, address indexed token1, address pair, uint256 arg3)
func (_SwappiFactory *SwappiFactoryFilterer) ParsePairCreated(log types.Log) (*SwappiFactoryPairCreated, error) {
	event := new(SwappiFactoryPairCreated)
	if err := _SwappiFactory.contract.UnpackLog(event, "PairCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SwappiFactoryFeeToChangedIterator is returned from FilterFeeToChanged and is used to iterate over the raw logs and unpacked data for FeeToChanged events raised by the SwappiFactory contract.
type SwappiFactoryFeeToChangedIterator struct {
	Event *SwappiFactoryFeeToChanged // Event containing the contract specifics and raw log

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
func (it *SwappiFactoryFeeToChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwappiFactoryFeeToChanged)
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
		it.Event = new(SwappiFactoryFeeToChanged)
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
func (it *SwappiFactoryFeeToChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwappiFactoryFeeToChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwappiFactoryFeeToChanged represents a FeeToChanged event raised by the SwappiFactory contract.
type SwappiFactoryFeeToChanged struct {
	FeeTo common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterFeeToChanged is a free log retrieval operation binding the contract event 0xa0c77798bfd987cb29b84c7dc00b611739240b52e0c4192fdbb51e98e7f26ddb.
//
// Solidity: event feeToChanged(address feeTo)
func (_SwappiFactory *SwappiFactoryFilterer) FilterFeeToChanged(opts *bind.FilterOpts) (*SwappiFactoryFeeToChangedIterator, error) {

	logs, sub, err := _SwappiFactory.contract.FilterLogs(opts, "feeToChanged")
	if err != nil {
		return nil, err
	}
	return &SwappiFactoryFeeToChangedIterator{contract: _SwappiFactory.contract, event: "feeToChanged", logs: logs, sub: sub}, nil
}

// WatchFeeToChanged is a free log subscription operation binding the contract event 0xa0c77798bfd987cb29b84c7dc00b611739240b52e0c4192fdbb51e98e7f26ddb.
//
// Solidity: event feeToChanged(address feeTo)
func (_SwappiFactory *SwappiFactoryFilterer) WatchFeeToChanged(opts *bind.WatchOpts, sink chan<- *SwappiFactoryFeeToChanged) (event.Subscription, error) {

	logs, sub, err := _SwappiFactory.contract.WatchLogs(opts, "feeToChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwappiFactoryFeeToChanged)
				if err := _SwappiFactory.contract.UnpackLog(event, "feeToChanged", log); err != nil {
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

// ParseFeeToChanged is a log parse operation binding the contract event 0xa0c77798bfd987cb29b84c7dc00b611739240b52e0c4192fdbb51e98e7f26ddb.
//
// Solidity: event feeToChanged(address feeTo)
func (_SwappiFactory *SwappiFactoryFilterer) ParseFeeToChanged(log types.Log) (*SwappiFactoryFeeToChanged, error) {
	event := new(SwappiFactoryFeeToChanged)
	if err := _SwappiFactory.contract.UnpackLog(event, "feeToChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SwappiFactoryFeeToSetterChangedIterator is returned from FilterFeeToSetterChanged and is used to iterate over the raw logs and unpacked data for FeeToSetterChanged events raised by the SwappiFactory contract.
type SwappiFactoryFeeToSetterChangedIterator struct {
	Event *SwappiFactoryFeeToSetterChanged // Event containing the contract specifics and raw log

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
func (it *SwappiFactoryFeeToSetterChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwappiFactoryFeeToSetterChanged)
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
		it.Event = new(SwappiFactoryFeeToSetterChanged)
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
func (it *SwappiFactoryFeeToSetterChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwappiFactoryFeeToSetterChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwappiFactoryFeeToSetterChanged represents a FeeToSetterChanged event raised by the SwappiFactory contract.
type SwappiFactoryFeeToSetterChanged struct {
	FeeToSetter common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterFeeToSetterChanged is a free log retrieval operation binding the contract event 0x21acc73ed9ca7f270fce10277221147f56ad2fa56e5cae218f65618a845e484d.
//
// Solidity: event feeToSetterChanged(address feeToSetter)
func (_SwappiFactory *SwappiFactoryFilterer) FilterFeeToSetterChanged(opts *bind.FilterOpts) (*SwappiFactoryFeeToSetterChangedIterator, error) {

	logs, sub, err := _SwappiFactory.contract.FilterLogs(opts, "feeToSetterChanged")
	if err != nil {
		return nil, err
	}
	return &SwappiFactoryFeeToSetterChangedIterator{contract: _SwappiFactory.contract, event: "feeToSetterChanged", logs: logs, sub: sub}, nil
}

// WatchFeeToSetterChanged is a free log subscription operation binding the contract event 0x21acc73ed9ca7f270fce10277221147f56ad2fa56e5cae218f65618a845e484d.
//
// Solidity: event feeToSetterChanged(address feeToSetter)
func (_SwappiFactory *SwappiFactoryFilterer) WatchFeeToSetterChanged(opts *bind.WatchOpts, sink chan<- *SwappiFactoryFeeToSetterChanged) (event.Subscription, error) {

	logs, sub, err := _SwappiFactory.contract.WatchLogs(opts, "feeToSetterChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwappiFactoryFeeToSetterChanged)
				if err := _SwappiFactory.contract.UnpackLog(event, "feeToSetterChanged", log); err != nil {
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

// ParseFeeToSetterChanged is a log parse operation binding the contract event 0x21acc73ed9ca7f270fce10277221147f56ad2fa56e5cae218f65618a845e484d.
//
// Solidity: event feeToSetterChanged(address feeToSetter)
func (_SwappiFactory *SwappiFactoryFilterer) ParseFeeToSetterChanged(log types.Log) (*SwappiFactoryFeeToSetterChanged, error) {
	event := new(SwappiFactoryFeeToSetterChanged)
	if err := _SwappiFactory.contract.UnpackLog(event, "feeToSetterChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
