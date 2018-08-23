// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stub

import (
	"math/big"
	"strings"

	ethereum "github.com/themis-network/go-themis"
	"github.com/themis-network/go-themis/accounts/abi"
	"github.com/themis-network/go-themis/accounts/abi/bind"
	"github.com/themis-network/go-themis/common"
	"github.com/themis-network/go-themis/core/types"
	"github.com/themis-network/go-themis/event"
)

// TradeABI is the input ABI used to generate the binding from.
const TradeABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"addArbitrator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"trusteeNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"},{\"name\":\"user\",\"type\":\"address\"}],\"name\":\"isOrderTrustee\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"},{\"name\":\"secrets\",\"type\":\"string\"},{\"name\":\"userID\",\"type\":\"uint32\"},{\"name\":\"verifyData\",\"type\":\"string\"}],\"name\":\"uploadSecret\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"}],\"name\":\"getRequester\",\"outputs\":[{\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"}],\"name\":\"getOrderTrustees\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"}],\"name\":\"getOrderBuyer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"}],\"name\":\"finishOrder\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"},{\"name\":\"user\",\"type\":\"uint32\"}],\"name\":\"getVerifyData\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"},{\"name\":\"userID\",\"type\":\"uint32\"}],\"name\":\"confirmTradeOrder\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"removeArbitrator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"},{\"name\":\"winner\",\"type\":\"uint32\"}],\"name\":\"judge\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"isArbitrator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"}],\"name\":\"getWinner\",\"outputs\":[{\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_trustee\",\"type\":\"address\"}],\"name\":\"updateTrusteeContract\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"trusteeContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"},{\"name\":\"createUserID\",\"type\":\"uint32\"}],\"name\":\"cancelTrade\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"}],\"name\":\"getOrderSeller\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_trusteeNumber\",\"type\":\"uint8\"}],\"name\":\"updateDefaultTrusteeNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"}],\"name\":\"getOrderStatus\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawFee\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"}],\"name\":\"getPerFeeOfOrder\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"},{\"name\":\"user\",\"type\":\"uint32\"}],\"name\":\"arbitrate\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"},{\"name\":\"trusteeID\",\"type\":\"address\"},{\"name\":\"user\",\"type\":\"uint32\"}],\"name\":\"getSecret\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderID\",\"type\":\"uint80\"},{\"name\":\"userID\",\"type\":\"uint32\"},{\"name\":\"userType\",\"type\":\"uint8\"}],\"name\":\"createNewTradeOrder\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"orderID\",\"type\":\"uint80\"},{\"indexed\":true,\"name\":\"user\",\"type\":\"uint32\"},{\"indexed\":false,\"name\":\"userType\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"feePayed\",\"type\":\"uint256\"}],\"name\":\"LogCreateOrder\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"orderID\",\"type\":\"uint80\"},{\"indexed\":true,\"name\":\"creator\",\"type\":\"address\"}],\"name\":\"LogCancelTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"orderID\",\"type\":\"uint80\"},{\"indexed\":true,\"name\":\"user\",\"type\":\"uint32\"},{\"indexed\":false,\"name\":\"trustees\",\"type\":\"address[]\"},{\"indexed\":false,\"name\":\"feePayed\",\"type\":\"uint256\"}],\"name\":\"LogConfirmTradeOrder\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"orderID\",\"type\":\"uint80\"},{\"indexed\":true,\"name\":\"user\",\"type\":\"uint32\"},{\"indexed\":false,\"name\":\"secrets\",\"type\":\"string\"}],\"name\":\"LogUploadSecret\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"orderID\",\"type\":\"uint80\"}],\"name\":\"LogFinishOrder\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"trustee\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LogWithdrawFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"orderID\",\"type\":\"uint80\"},{\"indexed\":true,\"name\":\"user\",\"type\":\"uint32\"}],\"name\":\"Arbitrate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"orderID\",\"type\":\"uint80\"},{\"indexed\":true,\"name\":\"winner\",\"type\":\"uint32\"},{\"indexed\":true,\"name\":\"judge\",\"type\":\"address\"}],\"name\":\"Judge\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"who\",\"type\":\"address\"}],\"name\":\"AddArbitrator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"who\",\"type\":\"address\"}],\"name\":\"RemoveArbitrator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newNumber\",\"type\":\"uint256\"}],\"name\":\"LogUpdateDefaultTrusteeNumber\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"LogUpdateTrusteeContract\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// Trade is an auto generated Go binding around an Ethereum contract.
type Trade struct {
	TradeCaller     // Read-only binding to the contract
	TradeTransactor // Write-only binding to the contract
	TradeFilterer   // Log filterer for contract events
}

// TradeCaller is an auto generated read-only Go binding around an Ethereum contract.
type TradeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TradeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TradeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TradeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TradeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TradeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TradeSession struct {
	Contract     *Trade            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TradeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TradeCallerSession struct {
	Contract *TradeCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TradeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TradeTransactorSession struct {
	Contract     *TradeTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TradeRaw is an auto generated low-level Go binding around an Ethereum contract.
type TradeRaw struct {
	Contract *Trade // Generic contract binding to access the raw methods on
}

// TradeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TradeCallerRaw struct {
	Contract *TradeCaller // Generic read-only contract binding to access the raw methods on
}

// TradeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TradeTransactorRaw struct {
	Contract *TradeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTrade creates a new instance of Trade, bound to a specific deployed contract.
func NewTrade(address common.Address, backend bind.ContractBackend) (*Trade, error) {
	contract, err := bindTrade(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Trade{TradeCaller: TradeCaller{contract: contract}, TradeTransactor: TradeTransactor{contract: contract}, TradeFilterer: TradeFilterer{contract: contract}}, nil
}

// NewTradeCaller creates a new read-only instance of Trade, bound to a specific deployed contract.
func NewTradeCaller(address common.Address, caller bind.ContractCaller) (*TradeCaller, error) {
	contract, err := bindTrade(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TradeCaller{contract: contract}, nil
}

// NewTradeTransactor creates a new write-only instance of Trade, bound to a specific deployed contract.
func NewTradeTransactor(address common.Address, transactor bind.ContractTransactor) (*TradeTransactor, error) {
	contract, err := bindTrade(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TradeTransactor{contract: contract}, nil
}

// NewTradeFilterer creates a new log filterer instance of Trade, bound to a specific deployed contract.
func NewTradeFilterer(address common.Address, filterer bind.ContractFilterer) (*TradeFilterer, error) {
	contract, err := bindTrade(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TradeFilterer{contract: contract}, nil
}

// bindTrade binds a generic wrapper to an already deployed contract.
func bindTrade(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TradeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Trade *TradeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Trade.Contract.TradeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Trade *TradeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Trade.Contract.TradeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Trade *TradeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Trade.Contract.TradeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Trade *TradeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Trade.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Trade *TradeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Trade.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Trade *TradeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Trade.Contract.contract.Transact(opts, method, params...)
}

// GetOrderBuyer is a free data retrieval call binding the contract method 0x2666b767.
//
// Solidity: function getOrderBuyer(orderID uint80) constant returns(uint256)
func (_Trade *TradeCaller) GetOrderBuyer(opts *bind.CallOpts, orderID *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Trade.contract.Call(opts, out, "getOrderBuyer", orderID)
	return *ret0, err
}

// GetOrderBuyer is a free data retrieval call binding the contract method 0x2666b767.
//
// Solidity: function getOrderBuyer(orderID uint80) constant returns(uint256)
func (_Trade *TradeSession) GetOrderBuyer(orderID *big.Int) (*big.Int, error) {
	return _Trade.Contract.GetOrderBuyer(&_Trade.CallOpts, orderID)
}

// GetOrderBuyer is a free data retrieval call binding the contract method 0x2666b767.
//
// Solidity: function getOrderBuyer(orderID uint80) constant returns(uint256)
func (_Trade *TradeCallerSession) GetOrderBuyer(orderID *big.Int) (*big.Int, error) {
	return _Trade.Contract.GetOrderBuyer(&_Trade.CallOpts, orderID)
}

// GetOrderSeller is a free data retrieval call binding the contract method 0xdcaba7a3.
//
// Solidity: function getOrderSeller(orderID uint80) constant returns(uint256)
func (_Trade *TradeCaller) GetOrderSeller(opts *bind.CallOpts, orderID *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Trade.contract.Call(opts, out, "getOrderSeller", orderID)
	return *ret0, err
}

// GetOrderSeller is a free data retrieval call binding the contract method 0xdcaba7a3.
//
// Solidity: function getOrderSeller(orderID uint80) constant returns(uint256)
func (_Trade *TradeSession) GetOrderSeller(orderID *big.Int) (*big.Int, error) {
	return _Trade.Contract.GetOrderSeller(&_Trade.CallOpts, orderID)
}

// GetOrderSeller is a free data retrieval call binding the contract method 0xdcaba7a3.
//
// Solidity: function getOrderSeller(orderID uint80) constant returns(uint256)
func (_Trade *TradeCallerSession) GetOrderSeller(orderID *big.Int) (*big.Int, error) {
	return _Trade.Contract.GetOrderSeller(&_Trade.CallOpts, orderID)
}

// GetOrderStatus is a free data retrieval call binding the contract method 0xde757e93.
//
// Solidity: function getOrderStatus(orderID uint80) constant returns(uint8)
func (_Trade *TradeCaller) GetOrderStatus(opts *bind.CallOpts, orderID *big.Int) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Trade.contract.Call(opts, out, "getOrderStatus", orderID)
	return *ret0, err
}

// GetOrderStatus is a free data retrieval call binding the contract method 0xde757e93.
//
// Solidity: function getOrderStatus(orderID uint80) constant returns(uint8)
func (_Trade *TradeSession) GetOrderStatus(orderID *big.Int) (uint8, error) {
	return _Trade.Contract.GetOrderStatus(&_Trade.CallOpts, orderID)
}

// GetOrderStatus is a free data retrieval call binding the contract method 0xde757e93.
//
// Solidity: function getOrderStatus(orderID uint80) constant returns(uint8)
func (_Trade *TradeCallerSession) GetOrderStatus(orderID *big.Int) (uint8, error) {
	return _Trade.Contract.GetOrderStatus(&_Trade.CallOpts, orderID)
}

// GetOrderTrustees is a free data retrieval call binding the contract method 0x1821bb82.
//
// Solidity: function getOrderTrustees(orderID uint80) constant returns(address[])
func (_Trade *TradeCaller) GetOrderTrustees(opts *bind.CallOpts, orderID *big.Int) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Trade.contract.Call(opts, out, "getOrderTrustees", orderID)
	return *ret0, err
}

// GetOrderTrustees is a free data retrieval call binding the contract method 0x1821bb82.
//
// Solidity: function getOrderTrustees(orderID uint80) constant returns(address[])
func (_Trade *TradeSession) GetOrderTrustees(orderID *big.Int) ([]common.Address, error) {
	return _Trade.Contract.GetOrderTrustees(&_Trade.CallOpts, orderID)
}

// GetOrderTrustees is a free data retrieval call binding the contract method 0x1821bb82.
//
// Solidity: function getOrderTrustees(orderID uint80) constant returns(address[])
func (_Trade *TradeCallerSession) GetOrderTrustees(orderID *big.Int) ([]common.Address, error) {
	return _Trade.Contract.GetOrderTrustees(&_Trade.CallOpts, orderID)
}

// GetPerFeeOfOrder is a free data retrieval call binding the contract method 0xf3241e96.
//
// Solidity: function getPerFeeOfOrder(orderID uint80) constant returns(uint256)
func (_Trade *TradeCaller) GetPerFeeOfOrder(opts *bind.CallOpts, orderID *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Trade.contract.Call(opts, out, "getPerFeeOfOrder", orderID)
	return *ret0, err
}

// GetPerFeeOfOrder is a free data retrieval call binding the contract method 0xf3241e96.
//
// Solidity: function getPerFeeOfOrder(orderID uint80) constant returns(uint256)
func (_Trade *TradeSession) GetPerFeeOfOrder(orderID *big.Int) (*big.Int, error) {
	return _Trade.Contract.GetPerFeeOfOrder(&_Trade.CallOpts, orderID)
}

// GetPerFeeOfOrder is a free data retrieval call binding the contract method 0xf3241e96.
//
// Solidity: function getPerFeeOfOrder(orderID uint80) constant returns(uint256)
func (_Trade *TradeCallerSession) GetPerFeeOfOrder(orderID *big.Int) (*big.Int, error) {
	return _Trade.Contract.GetPerFeeOfOrder(&_Trade.CallOpts, orderID)
}

// GetRequester is a free data retrieval call binding the contract method 0x0fc9afaf.
//
// Solidity: function getRequester(orderID uint80) constant returns(uint32)
func (_Trade *TradeCaller) GetRequester(opts *bind.CallOpts, orderID *big.Int) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _Trade.contract.Call(opts, out, "getRequester", orderID)
	return *ret0, err
}

// GetRequester is a free data retrieval call binding the contract method 0x0fc9afaf.
//
// Solidity: function getRequester(orderID uint80) constant returns(uint32)
func (_Trade *TradeSession) GetRequester(orderID *big.Int) (uint32, error) {
	return _Trade.Contract.GetRequester(&_Trade.CallOpts, orderID)
}

// GetRequester is a free data retrieval call binding the contract method 0x0fc9afaf.
//
// Solidity: function getRequester(orderID uint80) constant returns(uint32)
func (_Trade *TradeCallerSession) GetRequester(orderID *big.Int) (uint32, error) {
	return _Trade.Contract.GetRequester(&_Trade.CallOpts, orderID)
}

// GetSecret is a free data retrieval call binding the contract method 0xfd748b90.
//
// Solidity: function getSecret(orderID uint80, trusteeID address, user uint32) constant returns(string)
func (_Trade *TradeCaller) GetSecret(opts *bind.CallOpts, orderID *big.Int, trusteeID common.Address, user uint32) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Trade.contract.Call(opts, out, "getSecret", orderID, trusteeID, user)
	return *ret0, err
}

// GetSecret is a free data retrieval call binding the contract method 0xfd748b90.
//
// Solidity: function getSecret(orderID uint80, trusteeID address, user uint32) constant returns(string)
func (_Trade *TradeSession) GetSecret(orderID *big.Int, trusteeID common.Address, user uint32) (string, error) {
	return _Trade.Contract.GetSecret(&_Trade.CallOpts, orderID, trusteeID, user)
}

// GetSecret is a free data retrieval call binding the contract method 0xfd748b90.
//
// Solidity: function getSecret(orderID uint80, trusteeID address, user uint32) constant returns(string)
func (_Trade *TradeCallerSession) GetSecret(orderID *big.Int, trusteeID common.Address, user uint32) (string, error) {
	return _Trade.Contract.GetSecret(&_Trade.CallOpts, orderID, trusteeID, user)
}

// GetVerifyData is a free data retrieval call binding the contract method 0x465be783.
//
// Solidity: function getVerifyData(orderID uint80, user uint32) constant returns(string)
func (_Trade *TradeCaller) GetVerifyData(opts *bind.CallOpts, orderID *big.Int, user uint32) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Trade.contract.Call(opts, out, "getVerifyData", orderID, user)
	return *ret0, err
}

// GetVerifyData is a free data retrieval call binding the contract method 0x465be783.
//
// Solidity: function getVerifyData(orderID uint80, user uint32) constant returns(string)
func (_Trade *TradeSession) GetVerifyData(orderID *big.Int, user uint32) (string, error) {
	return _Trade.Contract.GetVerifyData(&_Trade.CallOpts, orderID, user)
}

// GetVerifyData is a free data retrieval call binding the contract method 0x465be783.
//
// Solidity: function getVerifyData(orderID uint80, user uint32) constant returns(string)
func (_Trade *TradeCallerSession) GetVerifyData(orderID *big.Int, user uint32) (string, error) {
	return _Trade.Contract.GetVerifyData(&_Trade.CallOpts, orderID, user)
}

// GetWinner is a free data retrieval call binding the contract method 0xac364eae.
//
// Solidity: function getWinner(orderID uint80) constant returns(uint32)
func (_Trade *TradeCaller) GetWinner(opts *bind.CallOpts, orderID *big.Int) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _Trade.contract.Call(opts, out, "getWinner", orderID)
	return *ret0, err
}

// GetWinner is a free data retrieval call binding the contract method 0xac364eae.
//
// Solidity: function getWinner(orderID uint80) constant returns(uint32)
func (_Trade *TradeSession) GetWinner(orderID *big.Int) (uint32, error) {
	return _Trade.Contract.GetWinner(&_Trade.CallOpts, orderID)
}

// GetWinner is a free data retrieval call binding the contract method 0xac364eae.
//
// Solidity: function getWinner(orderID uint80) constant returns(uint32)
func (_Trade *TradeCallerSession) GetWinner(orderID *big.Int) (uint32, error) {
	return _Trade.Contract.GetWinner(&_Trade.CallOpts, orderID)
}

// IsArbitrator is a free data retrieval call binding the contract method 0x9f6bd2a9.
//
// Solidity: function isArbitrator(who address) constant returns(bool)
func (_Trade *TradeCaller) IsArbitrator(opts *bind.CallOpts, who common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Trade.contract.Call(opts, out, "isArbitrator", who)
	return *ret0, err
}

// IsArbitrator is a free data retrieval call binding the contract method 0x9f6bd2a9.
//
// Solidity: function isArbitrator(who address) constant returns(bool)
func (_Trade *TradeSession) IsArbitrator(who common.Address) (bool, error) {
	return _Trade.Contract.IsArbitrator(&_Trade.CallOpts, who)
}

// IsArbitrator is a free data retrieval call binding the contract method 0x9f6bd2a9.
//
// Solidity: function isArbitrator(who address) constant returns(bool)
func (_Trade *TradeCallerSession) IsArbitrator(who common.Address) (bool, error) {
	return _Trade.Contract.IsArbitrator(&_Trade.CallOpts, who)
}

// IsOrderTrustee is a free data retrieval call binding the contract method 0x08903cf8.
//
// Solidity: function isOrderTrustee(orderID uint80, user address) constant returns(bool)
func (_Trade *TradeCaller) IsOrderTrustee(opts *bind.CallOpts, orderID *big.Int, user common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Trade.contract.Call(opts, out, "isOrderTrustee", orderID, user)
	return *ret0, err
}

// IsOrderTrustee is a free data retrieval call binding the contract method 0x08903cf8.
//
// Solidity: function isOrderTrustee(orderID uint80, user address) constant returns(bool)
func (_Trade *TradeSession) IsOrderTrustee(orderID *big.Int, user common.Address) (bool, error) {
	return _Trade.Contract.IsOrderTrustee(&_Trade.CallOpts, orderID, user)
}

// IsOrderTrustee is a free data retrieval call binding the contract method 0x08903cf8.
//
// Solidity: function isOrderTrustee(orderID uint80, user address) constant returns(bool)
func (_Trade *TradeCallerSession) IsOrderTrustee(orderID *big.Int, user common.Address) (bool, error) {
	return _Trade.Contract.IsOrderTrustee(&_Trade.CallOpts, orderID, user)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Trade *TradeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Trade.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Trade *TradeSession) Owner() (common.Address, error) {
	return _Trade.Contract.Owner(&_Trade.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Trade *TradeCallerSession) Owner() (common.Address, error) {
	return _Trade.Contract.Owner(&_Trade.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_Trade *TradeCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Trade.contract.Call(opts, out, "paused")
	return *ret0, err
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_Trade *TradeSession) Paused() (bool, error) {
	return _Trade.Contract.Paused(&_Trade.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_Trade *TradeCallerSession) Paused() (bool, error) {
	return _Trade.Contract.Paused(&_Trade.CallOpts)
}

// TrusteeContract is a free data retrieval call binding the contract method 0xbb133331.
//
// Solidity: function trusteeContract() constant returns(address)
func (_Trade *TradeCaller) TrusteeContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Trade.contract.Call(opts, out, "trusteeContract")
	return *ret0, err
}

// TrusteeContract is a free data retrieval call binding the contract method 0xbb133331.
//
// Solidity: function trusteeContract() constant returns(address)
func (_Trade *TradeSession) TrusteeContract() (common.Address, error) {
	return _Trade.Contract.TrusteeContract(&_Trade.CallOpts)
}

// TrusteeContract is a free data retrieval call binding the contract method 0xbb133331.
//
// Solidity: function trusteeContract() constant returns(address)
func (_Trade *TradeCallerSession) TrusteeContract() (common.Address, error) {
	return _Trade.Contract.TrusteeContract(&_Trade.CallOpts)
}

// TrusteeNumber is a free data retrieval call binding the contract method 0x044f4a4e.
//
// Solidity: function trusteeNumber() constant returns(uint8)
func (_Trade *TradeCaller) TrusteeNumber(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Trade.contract.Call(opts, out, "trusteeNumber")
	return *ret0, err
}

// TrusteeNumber is a free data retrieval call binding the contract method 0x044f4a4e.
//
// Solidity: function trusteeNumber() constant returns(uint8)
func (_Trade *TradeSession) TrusteeNumber() (uint8, error) {
	return _Trade.Contract.TrusteeNumber(&_Trade.CallOpts)
}

// TrusteeNumber is a free data retrieval call binding the contract method 0x044f4a4e.
//
// Solidity: function trusteeNumber() constant returns(uint8)
func (_Trade *TradeCallerSession) TrusteeNumber() (uint8, error) {
	return _Trade.Contract.TrusteeNumber(&_Trade.CallOpts)
}

// AddArbitrator is a paid mutator transaction binding the contract method 0x01fabd75.
//
// Solidity: function addArbitrator(who address) returns(bool)
func (_Trade *TradeTransactor) AddArbitrator(opts *bind.TransactOpts, who common.Address) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "addArbitrator", who)
}

// AddArbitrator is a paid mutator transaction binding the contract method 0x01fabd75.
//
// Solidity: function addArbitrator(who address) returns(bool)
func (_Trade *TradeSession) AddArbitrator(who common.Address) (*types.Transaction, error) {
	return _Trade.Contract.AddArbitrator(&_Trade.TransactOpts, who)
}

// AddArbitrator is a paid mutator transaction binding the contract method 0x01fabd75.
//
// Solidity: function addArbitrator(who address) returns(bool)
func (_Trade *TradeTransactorSession) AddArbitrator(who common.Address) (*types.Transaction, error) {
	return _Trade.Contract.AddArbitrator(&_Trade.TransactOpts, who)
}

// Arbitrate is a paid mutator transaction binding the contract method 0xfacf2dd7.
//
// Solidity: function arbitrate(orderID uint80, user uint32) returns(bool)
func (_Trade *TradeTransactor) Arbitrate(opts *bind.TransactOpts, orderID *big.Int, user uint32) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "arbitrate", orderID, user)
}

// Arbitrate is a paid mutator transaction binding the contract method 0xfacf2dd7.
//
// Solidity: function arbitrate(orderID uint80, user uint32) returns(bool)
func (_Trade *TradeSession) Arbitrate(orderID *big.Int, user uint32) (*types.Transaction, error) {
	return _Trade.Contract.Arbitrate(&_Trade.TransactOpts, orderID, user)
}

// Arbitrate is a paid mutator transaction binding the contract method 0xfacf2dd7.
//
// Solidity: function arbitrate(orderID uint80, user uint32) returns(bool)
func (_Trade *TradeTransactorSession) Arbitrate(orderID *big.Int, user uint32) (*types.Transaction, error) {
	return _Trade.Contract.Arbitrate(&_Trade.TransactOpts, orderID, user)
}

// CancelTrade is a paid mutator transaction binding the contract method 0xce36d79b.
//
// Solidity: function cancelTrade(orderID uint80, createUserID uint32) returns(bool)
func (_Trade *TradeTransactor) CancelTrade(opts *bind.TransactOpts, orderID *big.Int, createUserID uint32) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "cancelTrade", orderID, createUserID)
}

// CancelTrade is a paid mutator transaction binding the contract method 0xce36d79b.
//
// Solidity: function cancelTrade(orderID uint80, createUserID uint32) returns(bool)
func (_Trade *TradeSession) CancelTrade(orderID *big.Int, createUserID uint32) (*types.Transaction, error) {
	return _Trade.Contract.CancelTrade(&_Trade.TransactOpts, orderID, createUserID)
}

// CancelTrade is a paid mutator transaction binding the contract method 0xce36d79b.
//
// Solidity: function cancelTrade(orderID uint80, createUserID uint32) returns(bool)
func (_Trade *TradeTransactorSession) CancelTrade(orderID *big.Int, createUserID uint32) (*types.Transaction, error) {
	return _Trade.Contract.CancelTrade(&_Trade.TransactOpts, orderID, createUserID)
}

// ConfirmTradeOrder is a paid mutator transaction binding the contract method 0x6f9e0b40.
//
// Solidity: function confirmTradeOrder(orderID uint80, userID uint32) returns(bool)
func (_Trade *TradeTransactor) ConfirmTradeOrder(opts *bind.TransactOpts, orderID *big.Int, userID uint32) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "confirmTradeOrder", orderID, userID)
}

// ConfirmTradeOrder is a paid mutator transaction binding the contract method 0x6f9e0b40.
//
// Solidity: function confirmTradeOrder(orderID uint80, userID uint32) returns(bool)
func (_Trade *TradeSession) ConfirmTradeOrder(orderID *big.Int, userID uint32) (*types.Transaction, error) {
	return _Trade.Contract.ConfirmTradeOrder(&_Trade.TransactOpts, orderID, userID)
}

// ConfirmTradeOrder is a paid mutator transaction binding the contract method 0x6f9e0b40.
//
// Solidity: function confirmTradeOrder(orderID uint80, userID uint32) returns(bool)
func (_Trade *TradeTransactorSession) ConfirmTradeOrder(orderID *big.Int, userID uint32) (*types.Transaction, error) {
	return _Trade.Contract.ConfirmTradeOrder(&_Trade.TransactOpts, orderID, userID)
}

// CreateNewTradeOrder is a paid mutator transaction binding the contract method 0xff06a28c.
//
// Solidity: function createNewTradeOrder(orderID uint80, userID uint32, userType uint8) returns(bool)
func (_Trade *TradeTransactor) CreateNewTradeOrder(opts *bind.TransactOpts, orderID *big.Int, userID uint32, userType uint8) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "createNewTradeOrder", orderID, userID, userType)
}

// CreateNewTradeOrder is a paid mutator transaction binding the contract method 0xff06a28c.
//
// Solidity: function createNewTradeOrder(orderID uint80, userID uint32, userType uint8) returns(bool)
func (_Trade *TradeSession) CreateNewTradeOrder(orderID *big.Int, userID uint32, userType uint8) (*types.Transaction, error) {
	return _Trade.Contract.CreateNewTradeOrder(&_Trade.TransactOpts, orderID, userID, userType)
}

// CreateNewTradeOrder is a paid mutator transaction binding the contract method 0xff06a28c.
//
// Solidity: function createNewTradeOrder(orderID uint80, userID uint32, userType uint8) returns(bool)
func (_Trade *TradeTransactorSession) CreateNewTradeOrder(orderID *big.Int, userID uint32, userType uint8) (*types.Transaction, error) {
	return _Trade.Contract.CreateNewTradeOrder(&_Trade.TransactOpts, orderID, userID, userType)
}

// FinishOrder is a paid mutator transaction binding the contract method 0x27acb487.
//
// Solidity: function finishOrder(orderID uint80) returns(bool)
func (_Trade *TradeTransactor) FinishOrder(opts *bind.TransactOpts, orderID *big.Int) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "finishOrder", orderID)
}

// FinishOrder is a paid mutator transaction binding the contract method 0x27acb487.
//
// Solidity: function finishOrder(orderID uint80) returns(bool)
func (_Trade *TradeSession) FinishOrder(orderID *big.Int) (*types.Transaction, error) {
	return _Trade.Contract.FinishOrder(&_Trade.TransactOpts, orderID)
}

// FinishOrder is a paid mutator transaction binding the contract method 0x27acb487.
//
// Solidity: function finishOrder(orderID uint80) returns(bool)
func (_Trade *TradeTransactorSession) FinishOrder(orderID *big.Int) (*types.Transaction, error) {
	return _Trade.Contract.FinishOrder(&_Trade.TransactOpts, orderID)
}

// Judge is a paid mutator transaction binding the contract method 0x98f3fafa.
//
// Solidity: function judge(orderID uint80, winner uint32) returns(bool)
func (_Trade *TradeTransactor) Judge(opts *bind.TransactOpts, orderID *big.Int, winner uint32) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "judge", orderID, winner)
}

// Judge is a paid mutator transaction binding the contract method 0x98f3fafa.
//
// Solidity: function judge(orderID uint80, winner uint32) returns(bool)
func (_Trade *TradeSession) Judge(orderID *big.Int, winner uint32) (*types.Transaction, error) {
	return _Trade.Contract.Judge(&_Trade.TransactOpts, orderID, winner)
}

// Judge is a paid mutator transaction binding the contract method 0x98f3fafa.
//
// Solidity: function judge(orderID uint80, winner uint32) returns(bool)
func (_Trade *TradeTransactorSession) Judge(orderID *big.Int, winner uint32) (*types.Transaction, error) {
	return _Trade.Contract.Judge(&_Trade.TransactOpts, orderID, winner)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Trade *TradeTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Trade *TradeSession) Pause() (*types.Transaction, error) {
	return _Trade.Contract.Pause(&_Trade.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Trade *TradeTransactorSession) Pause() (*types.Transaction, error) {
	return _Trade.Contract.Pause(&_Trade.TransactOpts)
}

// RemoveArbitrator is a paid mutator transaction binding the contract method 0x973ad270.
//
// Solidity: function removeArbitrator(who address) returns(bool)
func (_Trade *TradeTransactor) RemoveArbitrator(opts *bind.TransactOpts, who common.Address) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "removeArbitrator", who)
}

// RemoveArbitrator is a paid mutator transaction binding the contract method 0x973ad270.
//
// Solidity: function removeArbitrator(who address) returns(bool)
func (_Trade *TradeSession) RemoveArbitrator(who common.Address) (*types.Transaction, error) {
	return _Trade.Contract.RemoveArbitrator(&_Trade.TransactOpts, who)
}

// RemoveArbitrator is a paid mutator transaction binding the contract method 0x973ad270.
//
// Solidity: function removeArbitrator(who address) returns(bool)
func (_Trade *TradeTransactorSession) RemoveArbitrator(who common.Address) (*types.Transaction, error) {
	return _Trade.Contract.RemoveArbitrator(&_Trade.TransactOpts, who)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Trade *TradeTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Trade *TradeSession) RenounceOwnership() (*types.Transaction, error) {
	return _Trade.Contract.RenounceOwnership(&_Trade.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Trade *TradeTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Trade.Contract.RenounceOwnership(&_Trade.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Trade *TradeTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Trade *TradeSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Trade.Contract.TransferOwnership(&_Trade.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Trade *TradeTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Trade.Contract.TransferOwnership(&_Trade.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Trade *TradeTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Trade *TradeSession) Unpause() (*types.Transaction, error) {
	return _Trade.Contract.Unpause(&_Trade.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Trade *TradeTransactorSession) Unpause() (*types.Transaction, error) {
	return _Trade.Contract.Unpause(&_Trade.TransactOpts)
}

// UpdateDefaultTrusteeNumber is a paid mutator transaction binding the contract method 0xdd080e87.
//
// Solidity: function updateDefaultTrusteeNumber(_trusteeNumber uint8) returns(bool)
func (_Trade *TradeTransactor) UpdateDefaultTrusteeNumber(opts *bind.TransactOpts, _trusteeNumber uint8) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "updateDefaultTrusteeNumber", _trusteeNumber)
}

// UpdateDefaultTrusteeNumber is a paid mutator transaction binding the contract method 0xdd080e87.
//
// Solidity: function updateDefaultTrusteeNumber(_trusteeNumber uint8) returns(bool)
func (_Trade *TradeSession) UpdateDefaultTrusteeNumber(_trusteeNumber uint8) (*types.Transaction, error) {
	return _Trade.Contract.UpdateDefaultTrusteeNumber(&_Trade.TransactOpts, _trusteeNumber)
}

// UpdateDefaultTrusteeNumber is a paid mutator transaction binding the contract method 0xdd080e87.
//
// Solidity: function updateDefaultTrusteeNumber(_trusteeNumber uint8) returns(bool)
func (_Trade *TradeTransactorSession) UpdateDefaultTrusteeNumber(_trusteeNumber uint8) (*types.Transaction, error) {
	return _Trade.Contract.UpdateDefaultTrusteeNumber(&_Trade.TransactOpts, _trusteeNumber)
}

// UpdateTrusteeContract is a paid mutator transaction binding the contract method 0xad7ce550.
//
// Solidity: function updateTrusteeContract(_trustee address) returns(bool)
func (_Trade *TradeTransactor) UpdateTrusteeContract(opts *bind.TransactOpts, _trustee common.Address) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "updateTrusteeContract", _trustee)
}

// UpdateTrusteeContract is a paid mutator transaction binding the contract method 0xad7ce550.
//
// Solidity: function updateTrusteeContract(_trustee address) returns(bool)
func (_Trade *TradeSession) UpdateTrusteeContract(_trustee common.Address) (*types.Transaction, error) {
	return _Trade.Contract.UpdateTrusteeContract(&_Trade.TransactOpts, _trustee)
}

// UpdateTrusteeContract is a paid mutator transaction binding the contract method 0xad7ce550.
//
// Solidity: function updateTrusteeContract(_trustee address) returns(bool)
func (_Trade *TradeTransactorSession) UpdateTrusteeContract(_trustee common.Address) (*types.Transaction, error) {
	return _Trade.Contract.UpdateTrusteeContract(&_Trade.TransactOpts, _trustee)
}

// UploadSecret is a paid mutator transaction binding the contract method 0x0a7db658.
//
// Solidity: function uploadSecret(orderID uint80, secrets string, userID uint32, verifyData string) returns(bool)
func (_Trade *TradeTransactor) UploadSecret(opts *bind.TransactOpts, orderID *big.Int, secrets string, userID uint32, verifyData string) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "uploadSecret", orderID, secrets, userID, verifyData)
}

// UploadSecret is a paid mutator transaction binding the contract method 0x0a7db658.
//
// Solidity: function uploadSecret(orderID uint80, secrets string, userID uint32, verifyData string) returns(bool)
func (_Trade *TradeSession) UploadSecret(orderID *big.Int, secrets string, userID uint32, verifyData string) (*types.Transaction, error) {
	return _Trade.Contract.UploadSecret(&_Trade.TransactOpts, orderID, secrets, userID, verifyData)
}

// UploadSecret is a paid mutator transaction binding the contract method 0x0a7db658.
//
// Solidity: function uploadSecret(orderID uint80, secrets string, userID uint32, verifyData string) returns(bool)
func (_Trade *TradeTransactorSession) UploadSecret(orderID *big.Int, secrets string, userID uint32, verifyData string) (*types.Transaction, error) {
	return _Trade.Contract.UploadSecret(&_Trade.TransactOpts, orderID, secrets, userID, verifyData)
}

// WithdrawFee is a paid mutator transaction binding the contract method 0xe941fa78.
//
// Solidity: function withdrawFee() returns(bool)
func (_Trade *TradeTransactor) WithdrawFee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "withdrawFee")
}

// WithdrawFee is a paid mutator transaction binding the contract method 0xe941fa78.
//
// Solidity: function withdrawFee() returns(bool)
func (_Trade *TradeSession) WithdrawFee() (*types.Transaction, error) {
	return _Trade.Contract.WithdrawFee(&_Trade.TransactOpts)
}

// WithdrawFee is a paid mutator transaction binding the contract method 0xe941fa78.
//
// Solidity: function withdrawFee() returns(bool)
func (_Trade *TradeTransactorSession) WithdrawFee() (*types.Transaction, error) {
	return _Trade.Contract.WithdrawFee(&_Trade.TransactOpts)
}

// TradeAddArbitratorIterator is returned from FilterAddArbitrator and is used to iterate over the raw logs and unpacked data for AddArbitrator events raised by the Trade contract.
type TradeAddArbitratorIterator struct {
	Event *TradeAddArbitrator // Event containing the contract specifics and raw log

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
func (it *TradeAddArbitratorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeAddArbitrator)
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
		it.Event = new(TradeAddArbitrator)
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
func (it *TradeAddArbitratorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeAddArbitratorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeAddArbitrator represents a AddArbitrator event raised by the Trade contract.
type TradeAddArbitrator struct {
	Who common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterAddArbitrator is a free log retrieval operation binding the contract event 0x1cfba79c837dd282b5affd88ad85c693d8f3fc6abb9999b92849f776f499045b.
//
// Solidity: e AddArbitrator(who indexed address)
func (_Trade *TradeFilterer) FilterAddArbitrator(opts *bind.FilterOpts, who []common.Address) (*TradeAddArbitratorIterator, error) {

	var whoRule []interface{}
	for _, whoItem := range who {
		whoRule = append(whoRule, whoItem)
	}

	logs, sub, err := _Trade.contract.FilterLogs(opts, "AddArbitrator", whoRule)
	if err != nil {
		return nil, err
	}
	return &TradeAddArbitratorIterator{contract: _Trade.contract, event: "AddArbitrator", logs: logs, sub: sub}, nil
}

// WatchAddArbitrator is a free log subscription operation binding the contract event 0x1cfba79c837dd282b5affd88ad85c693d8f3fc6abb9999b92849f776f499045b.
//
// Solidity: e AddArbitrator(who indexed address)
func (_Trade *TradeFilterer) WatchAddArbitrator(opts *bind.WatchOpts, sink chan<- *TradeAddArbitrator, who []common.Address) (event.Subscription, error) {

	var whoRule []interface{}
	for _, whoItem := range who {
		whoRule = append(whoRule, whoItem)
	}

	logs, sub, err := _Trade.contract.WatchLogs(opts, "AddArbitrator", whoRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeAddArbitrator)
				if err := _Trade.contract.UnpackLog(event, "AddArbitrator", log); err != nil {
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

// TradeArbitrateIterator is returned from FilterArbitrate and is used to iterate over the raw logs and unpacked data for Arbitrate events raised by the Trade contract.
type TradeArbitrateIterator struct {
	Event *TradeArbitrate // Event containing the contract specifics and raw log

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
func (it *TradeArbitrateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeArbitrate)
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
		it.Event = new(TradeArbitrate)
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
func (it *TradeArbitrateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeArbitrateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeArbitrate represents a Arbitrate event raised by the Trade contract.
type TradeArbitrate struct {
	OrderID *big.Int
	User    uint32
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterArbitrate is a free log retrieval operation binding the contract event 0x1091687dccfb26b1f39c83752f98f95d7220fd8aeedcc9cea7d8c312af711cc1.
//
// Solidity: e Arbitrate(orderID uint80, user indexed uint32)
func (_Trade *TradeFilterer) FilterArbitrate(opts *bind.FilterOpts, user []uint32) (*TradeArbitrateIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Trade.contract.FilterLogs(opts, "Arbitrate", userRule)
	if err != nil {
		return nil, err
	}
	return &TradeArbitrateIterator{contract: _Trade.contract, event: "Arbitrate", logs: logs, sub: sub}, nil
}

// WatchArbitrate is a free log subscription operation binding the contract event 0x1091687dccfb26b1f39c83752f98f95d7220fd8aeedcc9cea7d8c312af711cc1.
//
// Solidity: e Arbitrate(orderID uint80, user indexed uint32)
func (_Trade *TradeFilterer) WatchArbitrate(opts *bind.WatchOpts, sink chan<- *TradeArbitrate, user []uint32) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Trade.contract.WatchLogs(opts, "Arbitrate", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeArbitrate)
				if err := _Trade.contract.UnpackLog(event, "Arbitrate", log); err != nil {
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

// TradeJudgeIterator is returned from FilterJudge and is used to iterate over the raw logs and unpacked data for Judge events raised by the Trade contract.
type TradeJudgeIterator struct {
	Event *TradeJudge // Event containing the contract specifics and raw log

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
func (it *TradeJudgeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeJudge)
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
		it.Event = new(TradeJudge)
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
func (it *TradeJudgeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeJudgeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeJudge represents a Judge event raised by the Trade contract.
type TradeJudge struct {
	OrderID *big.Int
	Winner  uint32
	Judge   common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterJudge is a free log retrieval operation binding the contract event 0x91cd4f1183722432c51b6df8de80eef3598ebc4ef71ce76150fd2c8e4ed9b092.
//
// Solidity: e Judge(orderID uint80, winner indexed uint32, judge indexed address)
func (_Trade *TradeFilterer) FilterJudge(opts *bind.FilterOpts, winner []uint32, judge []common.Address) (*TradeJudgeIterator, error) {

	var winnerRule []interface{}
	for _, winnerItem := range winner {
		winnerRule = append(winnerRule, winnerItem)
	}
	var judgeRule []interface{}
	for _, judgeItem := range judge {
		judgeRule = append(judgeRule, judgeItem)
	}

	logs, sub, err := _Trade.contract.FilterLogs(opts, "Judge", winnerRule, judgeRule)
	if err != nil {
		return nil, err
	}
	return &TradeJudgeIterator{contract: _Trade.contract, event: "Judge", logs: logs, sub: sub}, nil
}

// WatchJudge is a free log subscription operation binding the contract event 0x91cd4f1183722432c51b6df8de80eef3598ebc4ef71ce76150fd2c8e4ed9b092.
//
// Solidity: e Judge(orderID uint80, winner indexed uint32, judge indexed address)
func (_Trade *TradeFilterer) WatchJudge(opts *bind.WatchOpts, sink chan<- *TradeJudge, winner []uint32, judge []common.Address) (event.Subscription, error) {

	var winnerRule []interface{}
	for _, winnerItem := range winner {
		winnerRule = append(winnerRule, winnerItem)
	}
	var judgeRule []interface{}
	for _, judgeItem := range judge {
		judgeRule = append(judgeRule, judgeItem)
	}

	logs, sub, err := _Trade.contract.WatchLogs(opts, "Judge", winnerRule, judgeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeJudge)
				if err := _Trade.contract.UnpackLog(event, "Judge", log); err != nil {
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

// TradeLogCancelTradeIterator is returned from FilterLogCancelTrade and is used to iterate over the raw logs and unpacked data for LogCancelTrade events raised by the Trade contract.
type TradeLogCancelTradeIterator struct {
	Event *TradeLogCancelTrade // Event containing the contract specifics and raw log

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
func (it *TradeLogCancelTradeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeLogCancelTrade)
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
		it.Event = new(TradeLogCancelTrade)
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
func (it *TradeLogCancelTradeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeLogCancelTradeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeLogCancelTrade represents a LogCancelTrade event raised by the Trade contract.
type TradeLogCancelTrade struct {
	OrderID *big.Int
	Creator common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogCancelTrade is a free log retrieval operation binding the contract event 0x6e2a76d5e51755e0d0ead9415eb85f81de3e3024755a804f28728cf4c42adb1d.
//
// Solidity: e LogCancelTrade(orderID indexed uint80, creator indexed address)
func (_Trade *TradeFilterer) FilterLogCancelTrade(opts *bind.FilterOpts, orderID []*big.Int, creator []common.Address) (*TradeLogCancelTradeIterator, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _Trade.contract.FilterLogs(opts, "LogCancelTrade", orderIDRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return &TradeLogCancelTradeIterator{contract: _Trade.contract, event: "LogCancelTrade", logs: logs, sub: sub}, nil
}

// WatchLogCancelTrade is a free log subscription operation binding the contract event 0x6e2a76d5e51755e0d0ead9415eb85f81de3e3024755a804f28728cf4c42adb1d.
//
// Solidity: e LogCancelTrade(orderID indexed uint80, creator indexed address)
func (_Trade *TradeFilterer) WatchLogCancelTrade(opts *bind.WatchOpts, sink chan<- *TradeLogCancelTrade, orderID []*big.Int, creator []common.Address) (event.Subscription, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _Trade.contract.WatchLogs(opts, "LogCancelTrade", orderIDRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeLogCancelTrade)
				if err := _Trade.contract.UnpackLog(event, "LogCancelTrade", log); err != nil {
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

// TradeLogConfirmTradeOrderIterator is returned from FilterLogConfirmTradeOrder and is used to iterate over the raw logs and unpacked data for LogConfirmTradeOrder events raised by the Trade contract.
type TradeLogConfirmTradeOrderIterator struct {
	Event *TradeLogConfirmTradeOrder // Event containing the contract specifics and raw log

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
func (it *TradeLogConfirmTradeOrderIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeLogConfirmTradeOrder)
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
		it.Event = new(TradeLogConfirmTradeOrder)
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
func (it *TradeLogConfirmTradeOrderIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeLogConfirmTradeOrderIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeLogConfirmTradeOrder represents a LogConfirmTradeOrder event raised by the Trade contract.
type TradeLogConfirmTradeOrder struct {
	OrderID  *big.Int
	User     uint32
	Trustees []common.Address
	FeePayed *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogConfirmTradeOrder is a free log retrieval operation binding the contract event 0xa6115165a89d929fa9c050d5546756b15c16b4cbf0ed66bd893c9209aff2db44.
//
// Solidity: e LogConfirmTradeOrder(orderID indexed uint80, user indexed uint32, trustees address[], feePayed uint256)
func (_Trade *TradeFilterer) FilterLogConfirmTradeOrder(opts *bind.FilterOpts, orderID []*big.Int, user []uint32) (*TradeLogConfirmTradeOrderIterator, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Trade.contract.FilterLogs(opts, "LogConfirmTradeOrder", orderIDRule, userRule)
	if err != nil {
		return nil, err
	}
	return &TradeLogConfirmTradeOrderIterator{contract: _Trade.contract, event: "LogConfirmTradeOrder", logs: logs, sub: sub}, nil
}

// WatchLogConfirmTradeOrder is a free log subscription operation binding the contract event 0xa6115165a89d929fa9c050d5546756b15c16b4cbf0ed66bd893c9209aff2db44.
//
// Solidity: e LogConfirmTradeOrder(orderID indexed uint80, user indexed uint32, trustees address[], feePayed uint256)
func (_Trade *TradeFilterer) WatchLogConfirmTradeOrder(opts *bind.WatchOpts, sink chan<- *TradeLogConfirmTradeOrder, orderID []*big.Int, user []uint32) (event.Subscription, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Trade.contract.WatchLogs(opts, "LogConfirmTradeOrder", orderIDRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeLogConfirmTradeOrder)
				if err := _Trade.contract.UnpackLog(event, "LogConfirmTradeOrder", log); err != nil {
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

// TradeLogCreateOrderIterator is returned from FilterLogCreateOrder and is used to iterate over the raw logs and unpacked data for LogCreateOrder events raised by the Trade contract.
type TradeLogCreateOrderIterator struct {
	Event *TradeLogCreateOrder // Event containing the contract specifics and raw log

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
func (it *TradeLogCreateOrderIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeLogCreateOrder)
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
		it.Event = new(TradeLogCreateOrder)
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
func (it *TradeLogCreateOrderIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeLogCreateOrderIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeLogCreateOrder represents a LogCreateOrder event raised by the Trade contract.
type TradeLogCreateOrder struct {
	OrderID  *big.Int
	User     uint32
	UserType uint8
	FeePayed *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogCreateOrder is a free log retrieval operation binding the contract event 0xe435d37c6b977d33a57d31d222155c537fde79601f1609f780be49c06d03c672.
//
// Solidity: e LogCreateOrder(orderID indexed uint80, user indexed uint32, userType uint8, feePayed uint256)
func (_Trade *TradeFilterer) FilterLogCreateOrder(opts *bind.FilterOpts, orderID []*big.Int, user []uint32) (*TradeLogCreateOrderIterator, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Trade.contract.FilterLogs(opts, "LogCreateOrder", orderIDRule, userRule)
	if err != nil {
		return nil, err
	}
	return &TradeLogCreateOrderIterator{contract: _Trade.contract, event: "LogCreateOrder", logs: logs, sub: sub}, nil
}

// WatchLogCreateOrder is a free log subscription operation binding the contract event 0xe435d37c6b977d33a57d31d222155c537fde79601f1609f780be49c06d03c672.
//
// Solidity: e LogCreateOrder(orderID indexed uint80, user indexed uint32, userType uint8, feePayed uint256)
func (_Trade *TradeFilterer) WatchLogCreateOrder(opts *bind.WatchOpts, sink chan<- *TradeLogCreateOrder, orderID []*big.Int, user []uint32) (event.Subscription, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Trade.contract.WatchLogs(opts, "LogCreateOrder", orderIDRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeLogCreateOrder)
				if err := _Trade.contract.UnpackLog(event, "LogCreateOrder", log); err != nil {
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

// TradeLogFinishOrderIterator is returned from FilterLogFinishOrder and is used to iterate over the raw logs and unpacked data for LogFinishOrder events raised by the Trade contract.
type TradeLogFinishOrderIterator struct {
	Event *TradeLogFinishOrder // Event containing the contract specifics and raw log

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
func (it *TradeLogFinishOrderIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeLogFinishOrder)
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
		it.Event = new(TradeLogFinishOrder)
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
func (it *TradeLogFinishOrderIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeLogFinishOrderIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeLogFinishOrder represents a LogFinishOrder event raised by the Trade contract.
type TradeLogFinishOrder struct {
	OrderID *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogFinishOrder is a free log retrieval operation binding the contract event 0x5c31eea1382076853b0d0e7c4e8c4d8601bb83e9c07aeb582ae9a1fcd2928ebb.
//
// Solidity: e LogFinishOrder(orderID indexed uint80)
func (_Trade *TradeFilterer) FilterLogFinishOrder(opts *bind.FilterOpts, orderID []*big.Int) (*TradeLogFinishOrderIterator, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}

	logs, sub, err := _Trade.contract.FilterLogs(opts, "LogFinishOrder", orderIDRule)
	if err != nil {
		return nil, err
	}
	return &TradeLogFinishOrderIterator{contract: _Trade.contract, event: "LogFinishOrder", logs: logs, sub: sub}, nil
}

// WatchLogFinishOrder is a free log subscription operation binding the contract event 0x5c31eea1382076853b0d0e7c4e8c4d8601bb83e9c07aeb582ae9a1fcd2928ebb.
//
// Solidity: e LogFinishOrder(orderID indexed uint80)
func (_Trade *TradeFilterer) WatchLogFinishOrder(opts *bind.WatchOpts, sink chan<- *TradeLogFinishOrder, orderID []*big.Int) (event.Subscription, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}

	logs, sub, err := _Trade.contract.WatchLogs(opts, "LogFinishOrder", orderIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeLogFinishOrder)
				if err := _Trade.contract.UnpackLog(event, "LogFinishOrder", log); err != nil {
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

// TradeLogUpdateDefaultTrusteeNumberIterator is returned from FilterLogUpdateDefaultTrusteeNumber and is used to iterate over the raw logs and unpacked data for LogUpdateDefaultTrusteeNumber events raised by the Trade contract.
type TradeLogUpdateDefaultTrusteeNumberIterator struct {
	Event *TradeLogUpdateDefaultTrusteeNumber // Event containing the contract specifics and raw log

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
func (it *TradeLogUpdateDefaultTrusteeNumberIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeLogUpdateDefaultTrusteeNumber)
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
		it.Event = new(TradeLogUpdateDefaultTrusteeNumber)
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
func (it *TradeLogUpdateDefaultTrusteeNumberIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeLogUpdateDefaultTrusteeNumberIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeLogUpdateDefaultTrusteeNumber represents a LogUpdateDefaultTrusteeNumber event raised by the Trade contract.
type TradeLogUpdateDefaultTrusteeNumber struct {
	NewNumber *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLogUpdateDefaultTrusteeNumber is a free log retrieval operation binding the contract event 0x9b0b2211ad6eb917787b88e888ce2b894e4e4c5111f1347c1b857a76a1b634ea.
//
// Solidity: e LogUpdateDefaultTrusteeNumber(newNumber uint256)
func (_Trade *TradeFilterer) FilterLogUpdateDefaultTrusteeNumber(opts *bind.FilterOpts) (*TradeLogUpdateDefaultTrusteeNumberIterator, error) {

	logs, sub, err := _Trade.contract.FilterLogs(opts, "LogUpdateDefaultTrusteeNumber")
	if err != nil {
		return nil, err
	}
	return &TradeLogUpdateDefaultTrusteeNumberIterator{contract: _Trade.contract, event: "LogUpdateDefaultTrusteeNumber", logs: logs, sub: sub}, nil
}

// WatchLogUpdateDefaultTrusteeNumber is a free log subscription operation binding the contract event 0x9b0b2211ad6eb917787b88e888ce2b894e4e4c5111f1347c1b857a76a1b634ea.
//
// Solidity: e LogUpdateDefaultTrusteeNumber(newNumber uint256)
func (_Trade *TradeFilterer) WatchLogUpdateDefaultTrusteeNumber(opts *bind.WatchOpts, sink chan<- *TradeLogUpdateDefaultTrusteeNumber) (event.Subscription, error) {

	logs, sub, err := _Trade.contract.WatchLogs(opts, "LogUpdateDefaultTrusteeNumber")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeLogUpdateDefaultTrusteeNumber)
				if err := _Trade.contract.UnpackLog(event, "LogUpdateDefaultTrusteeNumber", log); err != nil {
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

// TradeLogUpdateTrusteeContractIterator is returned from FilterLogUpdateTrusteeContract and is used to iterate over the raw logs and unpacked data for LogUpdateTrusteeContract events raised by the Trade contract.
type TradeLogUpdateTrusteeContractIterator struct {
	Event *TradeLogUpdateTrusteeContract // Event containing the contract specifics and raw log

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
func (it *TradeLogUpdateTrusteeContractIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeLogUpdateTrusteeContract)
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
		it.Event = new(TradeLogUpdateTrusteeContract)
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
func (it *TradeLogUpdateTrusteeContractIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeLogUpdateTrusteeContractIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeLogUpdateTrusteeContract represents a LogUpdateTrusteeContract event raised by the Trade contract.
type TradeLogUpdateTrusteeContract struct {
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogUpdateTrusteeContract is a free log retrieval operation binding the contract event 0x70a92b8c5bb39169e0a947c5d3a1fa4a03256fa46ed634a51a75ef801c4bd01c.
//
// Solidity: e LogUpdateTrusteeContract(newAddress indexed address)
func (_Trade *TradeFilterer) FilterLogUpdateTrusteeContract(opts *bind.FilterOpts, newAddress []common.Address) (*TradeLogUpdateTrusteeContractIterator, error) {

	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _Trade.contract.FilterLogs(opts, "LogUpdateTrusteeContract", newAddressRule)
	if err != nil {
		return nil, err
	}
	return &TradeLogUpdateTrusteeContractIterator{contract: _Trade.contract, event: "LogUpdateTrusteeContract", logs: logs, sub: sub}, nil
}

// WatchLogUpdateTrusteeContract is a free log subscription operation binding the contract event 0x70a92b8c5bb39169e0a947c5d3a1fa4a03256fa46ed634a51a75ef801c4bd01c.
//
// Solidity: e LogUpdateTrusteeContract(newAddress indexed address)
func (_Trade *TradeFilterer) WatchLogUpdateTrusteeContract(opts *bind.WatchOpts, sink chan<- *TradeLogUpdateTrusteeContract, newAddress []common.Address) (event.Subscription, error) {

	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _Trade.contract.WatchLogs(opts, "LogUpdateTrusteeContract", newAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeLogUpdateTrusteeContract)
				if err := _Trade.contract.UnpackLog(event, "LogUpdateTrusteeContract", log); err != nil {
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

// TradeLogUploadSecretIterator is returned from FilterLogUploadSecret and is used to iterate over the raw logs and unpacked data for LogUploadSecret events raised by the Trade contract.
type TradeLogUploadSecretIterator struct {
	Event *TradeLogUploadSecret // Event containing the contract specifics and raw log

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
func (it *TradeLogUploadSecretIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeLogUploadSecret)
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
		it.Event = new(TradeLogUploadSecret)
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
func (it *TradeLogUploadSecretIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeLogUploadSecretIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeLogUploadSecret represents a LogUploadSecret event raised by the Trade contract.
type TradeLogUploadSecret struct {
	OrderID *big.Int
	User    uint32
	Secrets string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogUploadSecret is a free log retrieval operation binding the contract event 0x8a59d01dda427123e224b10a5103435e6a94ce386bd3d81052074263f9defce8.
//
// Solidity: e LogUploadSecret(orderID indexed uint80, user indexed uint32, secrets string)
func (_Trade *TradeFilterer) FilterLogUploadSecret(opts *bind.FilterOpts, orderID []*big.Int, user []uint32) (*TradeLogUploadSecretIterator, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Trade.contract.FilterLogs(opts, "LogUploadSecret", orderIDRule, userRule)
	if err != nil {
		return nil, err
	}
	return &TradeLogUploadSecretIterator{contract: _Trade.contract, event: "LogUploadSecret", logs: logs, sub: sub}, nil
}

// WatchLogUploadSecret is a free log subscription operation binding the contract event 0x8a59d01dda427123e224b10a5103435e6a94ce386bd3d81052074263f9defce8.
//
// Solidity: e LogUploadSecret(orderID indexed uint80, user indexed uint32, secrets string)
func (_Trade *TradeFilterer) WatchLogUploadSecret(opts *bind.WatchOpts, sink chan<- *TradeLogUploadSecret, orderID []*big.Int, user []uint32) (event.Subscription, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Trade.contract.WatchLogs(opts, "LogUploadSecret", orderIDRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeLogUploadSecret)
				if err := _Trade.contract.UnpackLog(event, "LogUploadSecret", log); err != nil {
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

// TradeLogWithdrawFeeIterator is returned from FilterLogWithdrawFee and is used to iterate over the raw logs and unpacked data for LogWithdrawFee events raised by the Trade contract.
type TradeLogWithdrawFeeIterator struct {
	Event *TradeLogWithdrawFee // Event containing the contract specifics and raw log

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
func (it *TradeLogWithdrawFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeLogWithdrawFee)
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
		it.Event = new(TradeLogWithdrawFee)
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
func (it *TradeLogWithdrawFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeLogWithdrawFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeLogWithdrawFee represents a LogWithdrawFee event raised by the Trade contract.
type TradeLogWithdrawFee struct {
	Trustee common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogWithdrawFee is a free log retrieval operation binding the contract event 0xbcacd35e44ebcdaa615013d05335c060187b65a417e24f0714a9a6b629d64137.
//
// Solidity: e LogWithdrawFee(trustee indexed address, amount uint256)
func (_Trade *TradeFilterer) FilterLogWithdrawFee(opts *bind.FilterOpts, trustee []common.Address) (*TradeLogWithdrawFeeIterator, error) {

	var trusteeRule []interface{}
	for _, trusteeItem := range trustee {
		trusteeRule = append(trusteeRule, trusteeItem)
	}

	logs, sub, err := _Trade.contract.FilterLogs(opts, "LogWithdrawFee", trusteeRule)
	if err != nil {
		return nil, err
	}
	return &TradeLogWithdrawFeeIterator{contract: _Trade.contract, event: "LogWithdrawFee", logs: logs, sub: sub}, nil
}

// WatchLogWithdrawFee is a free log subscription operation binding the contract event 0xbcacd35e44ebcdaa615013d05335c060187b65a417e24f0714a9a6b629d64137.
//
// Solidity: e LogWithdrawFee(trustee indexed address, amount uint256)
func (_Trade *TradeFilterer) WatchLogWithdrawFee(opts *bind.WatchOpts, sink chan<- *TradeLogWithdrawFee, trustee []common.Address) (event.Subscription, error) {

	var trusteeRule []interface{}
	for _, trusteeItem := range trustee {
		trusteeRule = append(trusteeRule, trusteeItem)
	}

	logs, sub, err := _Trade.contract.WatchLogs(opts, "LogWithdrawFee", trusteeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeLogWithdrawFee)
				if err := _Trade.contract.UnpackLog(event, "LogWithdrawFee", log); err != nil {
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

// TradeOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the Trade contract.
type TradeOwnershipRenouncedIterator struct {
	Event *TradeOwnershipRenounced // Event containing the contract specifics and raw log

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
func (it *TradeOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeOwnershipRenounced)
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
		it.Event = new(TradeOwnershipRenounced)
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
func (it *TradeOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeOwnershipRenounced represents a OwnershipRenounced event raised by the Trade contract.
type TradeOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Trade *TradeFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*TradeOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Trade.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TradeOwnershipRenouncedIterator{contract: _Trade.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Trade *TradeFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *TradeOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Trade.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeOwnershipRenounced)
				if err := _Trade.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
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

// TradeOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Trade contract.
type TradeOwnershipTransferredIterator struct {
	Event *TradeOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TradeOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeOwnershipTransferred)
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
		it.Event = new(TradeOwnershipTransferred)
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
func (it *TradeOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeOwnershipTransferred represents a OwnershipTransferred event raised by the Trade contract.
type TradeOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Trade *TradeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TradeOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Trade.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TradeOwnershipTransferredIterator{contract: _Trade.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Trade *TradeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TradeOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Trade.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeOwnershipTransferred)
				if err := _Trade.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// TradePauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the Trade contract.
type TradePauseIterator struct {
	Event *TradePause // Event containing the contract specifics and raw log

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
func (it *TradePauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradePause)
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
		it.Event = new(TradePause)
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
func (it *TradePauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradePauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradePause represents a Pause event raised by the Trade contract.
type TradePause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: e Pause()
func (_Trade *TradeFilterer) FilterPause(opts *bind.FilterOpts) (*TradePauseIterator, error) {

	logs, sub, err := _Trade.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &TradePauseIterator{contract: _Trade.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: e Pause()
func (_Trade *TradeFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *TradePause) (event.Subscription, error) {

	logs, sub, err := _Trade.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradePause)
				if err := _Trade.contract.UnpackLog(event, "Pause", log); err != nil {
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

// TradeRemoveArbitratorIterator is returned from FilterRemoveArbitrator and is used to iterate over the raw logs and unpacked data for RemoveArbitrator events raised by the Trade contract.
type TradeRemoveArbitratorIterator struct {
	Event *TradeRemoveArbitrator // Event containing the contract specifics and raw log

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
func (it *TradeRemoveArbitratorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeRemoveArbitrator)
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
		it.Event = new(TradeRemoveArbitrator)
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
func (it *TradeRemoveArbitratorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeRemoveArbitratorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeRemoveArbitrator represents a RemoveArbitrator event raised by the Trade contract.
type TradeRemoveArbitrator struct {
	Who common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRemoveArbitrator is a free log retrieval operation binding the contract event 0xf9f012dbf94ec6f26d3a73fbbae56a56fc7b236c4390d891201576ba3aaeb891.
//
// Solidity: e RemoveArbitrator(who indexed address)
func (_Trade *TradeFilterer) FilterRemoveArbitrator(opts *bind.FilterOpts, who []common.Address) (*TradeRemoveArbitratorIterator, error) {

	var whoRule []interface{}
	for _, whoItem := range who {
		whoRule = append(whoRule, whoItem)
	}

	logs, sub, err := _Trade.contract.FilterLogs(opts, "RemoveArbitrator", whoRule)
	if err != nil {
		return nil, err
	}
	return &TradeRemoveArbitratorIterator{contract: _Trade.contract, event: "RemoveArbitrator", logs: logs, sub: sub}, nil
}

// WatchRemoveArbitrator is a free log subscription operation binding the contract event 0xf9f012dbf94ec6f26d3a73fbbae56a56fc7b236c4390d891201576ba3aaeb891.
//
// Solidity: e RemoveArbitrator(who indexed address)
func (_Trade *TradeFilterer) WatchRemoveArbitrator(opts *bind.WatchOpts, sink chan<- *TradeRemoveArbitrator, who []common.Address) (event.Subscription, error) {

	var whoRule []interface{}
	for _, whoItem := range who {
		whoRule = append(whoRule, whoItem)
	}

	logs, sub, err := _Trade.contract.WatchLogs(opts, "RemoveArbitrator", whoRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeRemoveArbitrator)
				if err := _Trade.contract.UnpackLog(event, "RemoveArbitrator", log); err != nil {
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

// TradeUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the Trade contract.
type TradeUnpauseIterator struct {
	Event *TradeUnpause // Event containing the contract specifics and raw log

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
func (it *TradeUnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeUnpause)
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
		it.Event = new(TradeUnpause)
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
func (it *TradeUnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeUnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeUnpause represents a Unpause event raised by the Trade contract.
type TradeUnpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: e Unpause()
func (_Trade *TradeFilterer) FilterUnpause(opts *bind.FilterOpts) (*TradeUnpauseIterator, error) {

	logs, sub, err := _Trade.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &TradeUnpauseIterator{contract: _Trade.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: e Unpause()
func (_Trade *TradeFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *TradeUnpause) (event.Subscription, error) {

	logs, sub, err := _Trade.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeUnpause)
				if err := _Trade.contract.UnpackLog(event, "Unpause", log); err != nil {
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
