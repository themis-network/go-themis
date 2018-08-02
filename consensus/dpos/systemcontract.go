package dpos

import (
	"math/big"
	"strings"

	"github.com/themis-network/go-themis/accounts/abi"
	"github.com/themis-network/go-themis/common"
	"github.com/themis-network/go-themis/core/types"
)

const (
	CURRENTCONTRACTABI          = "[{\"constant\":true,\"inputs\":[],\"name\":\"getTopProducers\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"
	MAINCONTRACTABI             = "[{\"constant\":true,\"inputs\":[],\"name\":\"getCurrentSystemContractAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"
	get_top_producers_methodid  = "0xdf339b08"
	get_current_system_contract = "0x114176e6"
)

var (
	main_system_contract_address    = common.BytesToAddress([]byte{9})
	current_system_contract_address = common.BytesToAddress([]byte{10})
)

// Upgradeable system contract for dpos
// Producer/Voters can send tx to system contract for reg, unreg, vote and so on.
var HardcodedContractsDpos = []HardcodedContract{
	NewMainContract(),
	NewCurrentSystemContract(),
}

type HardcodedContract interface {
	GetCode() string
	GetContractAddr() common.Address
	GetStorage() map[common.Hash]common.Hash
}

func NewMainContract() *MainContract {
	// Ignore error
	parsed, _ := abi.JSON(strings.NewReader(MAINCONTRACTABI))
	return &MainContract{
		Addr: main_system_contract_address,
		Code: "0x608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063114176e6146044575b600080fd5b348015604f57600080fd5b5060566098565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6000600a9050905600a165627a7a723058208e5a1fda24895d86a4099e256fb8798b6b6868f80d96afd0803c716eb943e8ba0029",
		ABI:  parsed,
	}
}

// Main dpos upgradeable system contract entry.
// Producer can send tx to this contract for proposing upgrade system
// contract, the storage data will leave on this contract.
type MainContract struct {
	Addr       common.Address
	Code       string
	StorageSet map[common.Hash]common.Hash
	ABI        abi.ABI
}

func (m *MainContract) GetCode() string {
	return m.Code
}

func (m *MainContract) GetContractAddr() common.Address {
	return m.Addr
}

func (m *MainContract) GetStorage() map[common.Hash]common.Hash {
	return m.StorageSet
}

func (m *MainContract) GetCurrentSystemContract(data []byte) common.Address {
	var res = new(common.Address)
	m.ABI.Unpack(res, "getCurrentSystemContractAddress", data)
	return *res
}

func NewCurrentSystemContract() *CurrentSystemContract {
	// Ignore error
	parsed, _ := abi.JSON(strings.NewReader(CURRENTCONTRACTABI))
	return &CurrentSystemContract{
		Addr: current_system_contract_address,
		Code: "0x60806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680637006e01c14610072578063c989eb66146100c9578063df339b0814610120578063f36707021461018c578063fdbbfaf4146101e3575b600080fd5b34801561007e57600080fd5b5061008761023a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156100d557600080fd5b506100de610263565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561012c57600080fd5b5061013561028d565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561017857808201518184015260208101905061015d565b505050509050019250505060405180910390f35b34801561019857600080fd5b506101a1610481565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156101ef57600080fd5b506101f86104ab565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60608060046040519080825280602002602001820160405280156102c05781602001602082028038833980820191505090505b5090506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff168160008151811015156102f457fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681600181518110151561036257fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168160028151811015156103d057fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681600381518110151561043e57fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508091505090565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050905600a165627a7a723058201ed0e4295b54e26d2f90341b69b3fb60c1f354d466f0c3b3e334e1d075deb4f40029",
		// This is just for test
		StorageSet: map[common.Hash]common.Hash{
			common.BytesToHash([]byte{0}): common.BytesToAddress([]byte{1}).Hash(),
			common.BytesToHash([]byte{1}): common.BytesToAddress([]byte{2}).Hash(),
			common.BytesToHash([]byte{2}): common.BytesToAddress([]byte{3}).Hash(),
			common.BytesToHash([]byte{3}): common.BytesToAddress([]byte{4}).Hash(),
		},
		ABI: parsed,
	}
}

// Current active system contract.
// Producer reg, unreg; Voter votes, unvotes will send tx to this contract.
type CurrentSystemContract struct {
	Addr       common.Address
	Code       string
	StorageSet map[common.Hash]common.Hash
	ABI        abi.ABI
}

func (c *CurrentSystemContract) GetCode() string {
	return c.Code
}

func (c *CurrentSystemContract) GetContractAddr() common.Address {
	return c.Addr
}

func (c *CurrentSystemContract) GetStorage() map[common.Hash]common.Hash {
	return c.StorageSet
}

func (c *CurrentSystemContract) GetTopProducers(data []byte) []common.Address {
	var res = new([]common.Address)
	c.ABI.Unpack(res, "getTopProducers", data)
	return *res
}

func GetTopProducersCall(header *types.Header, contract *common.Address) SystemCall {
	return NewCallMsg(contract, common.Hex2Bytes(get_top_producers_methodid), header.Number.Uint64())
}

func GetCurrentSystemContractCall(header *types.Header) SystemCall {
	copyAddr := main_system_contract_address
	return NewCallMsg(&copyAddr, common.Hex2Bytes(get_current_system_contract), header.Number.Uint64())
}

func NewCallMsg(addr *common.Address, data []byte, blockNum uint64) SystemCall {
	return SystemCall{
		ToAddr:    addr,
		DataField: data,
		GasSupply: 50000000,
		AtBlock:   blockNum,
	}
}

type SystemCall struct {
	ToAddr    *common.Address
	DataField []byte
	GasSupply uint64
	AtBlock   uint64
}

func (s SystemCall) From() common.Address { return common.Address{} }
func (s SystemCall) Nonce() uint64        { return 0 }
func (s SystemCall) CheckNonce() bool     { return false }
func (s SystemCall) To() *common.Address  { return s.ToAddr }
func (s SystemCall) GasPrice() *big.Int   { return new(big.Int).SetUint64(0) }
func (s SystemCall) Gas() uint64          { return s.GasSupply }
func (s SystemCall) Value() *big.Int      { return new(big.Int).SetUint64(0) }
func (s SystemCall) Data() []byte         { return s.DataField }
