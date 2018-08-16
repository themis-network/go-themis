package eth

import (
	"bytes"
    "math"
    "testing"
    
	"github.com/themis-network/go-themis/common"
	"github.com/themis-network/go-themis/consensus/dpos"
	"github.com/themis-network/go-themis/core"
	"github.com/themis-network/go-themis/core/types"
	"github.com/themis-network/go-themis/core/vm"
	"github.com/themis-network/go-themis/crypto"
	"github.com/themis-network/go-themis/ethdb"
	"github.com/themis-network/go-themis/params"
)

type TestData struct {
	name   string
	method string
	input  []byte
	want   []byte
	callTo common.Address
}

// Get evm for test
func getNewEvm(msg core.Message, genesisBlock *types.Block, chainConfig *params.ChainConfig, blockchain *core.BlockChain) *vm.EVM {
	context := core.NewEVMContext(msg, genesisBlock.Header(), blockchain, nil)
	statedb, _ := blockchain.StateAt(genesisBlock.Root())
	evm := vm.NewEVM(context, statedb, chainConfig, vm.Config{})
	return evm
}

func TestSystemContract(t *testing.T) {
	// Create memory database for test
	genesis := core.DefaultGenesisBlock()
	db := ethdb.NewMemDatabase()
	genesisBlock, err := genesis.Commit(db)
	if err != nil {
		t.Fatal(err)
	}

	// Create blockchain for test
	chainConfig := genesis.Config
	engine := dpos.New(chainConfig.Dpos)
	blockchain, err := core.NewBlockChain(db, &core.CacheConfig{}, chainConfig, engine, vm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	test := []TestData{
		{
			"depositForJoin",
			"depositForJoin()",
			[]byte{},
			core.DepositForJoinStorageValue.Bytes(),
			core.RegSystemContractAddr,
		},
		{
			"lockTimeForDeposit",
			"lockTimeForDeposit()",
			[]byte{},
			core.LockTimeForDepositStorageValue.Bytes(),
			core.RegSystemContractAddr,
		},
		{
			"lengthOFEpoch",
			"lengthOFEpoch()",
			[]byte{},
			core.LengthOFEpochStorageValue.Bytes(),
			core.RegSystemContractAddr,
		},
		{
			"initOutTime",
			"initOutTime()",
			[]byte{},
			core.InitOutTimeStorageValue.Bytes(),
			core.RegSystemContractAddr,
		},
		{
			"systemStorage",
			"systemStorage()",
			[]byte{},
			core.SystemContractStorageValue.Bytes(),
			core.RegSystemContractAddr,
		},
		{
			"producerOp",
			"producerOp()",
			[]byte{},
			core.ProducerOpStorageValue.Bytes(),
			core.RegSystemContractAddr,
		},
		{
			"initOutTime",
			"initOutTime()",
			[]byte{},
			core.OutTimeStorageValue.Bytes(),
			core.VoteSystemContractAddr,
		},
		{
			"leastDepositForVote",
			"leastDepositForVote()",
			[]byte{},
			core.LeastDepositForVoteStorageValue.Bytes(),
			core.VoteSystemContractAddr,
		},
		{
			"lockTimeForVote",
			"lockTimeForVote()",
			[]byte{},
			core.LockTimeForVoteStorageValue.Bytes(),
			core.VoteSystemContractAddr,
		},
		{
			"systemStorage",
			"systemStorage()",
			[]byte{},
			core.MainContractStorageValue.Bytes(),
			core.VoteSystemContractAddr,
		},
		{
			"voteSystemContract address",
			"getVoteSystemContract()",
			[]byte{},
			core.VoteSystemContractAddr.Hash().Bytes(),
			core.MainSystemContractAddr,
		},
		{
			"regSystemContract address",
			"getRegSystemContract()",
			[]byte{},
			core.RegSystemContractAddr.Hash().Bytes(),
			core.MainSystemContractAddr,
		},
		{
			"check vote system contract",
			"isSystemContract(address)",
			common.BytesToHash([]byte{11}).Bytes(),
			core.VoteContractBoolStorageValue.Bytes(),
			core.MainSystemContractAddr,
		},
		{
			"check reg system contract",
			"isSystemContract(address)",
			common.BytesToHash([]byte{10}).Bytes(),
			core.RegContractBoolStorageValue.Bytes(),
			core.MainSystemContractAddr,
		},
	}

	for _, v := range test {
		methodHash := crypto.Keccak256([]byte(v.method))
		input := methodHash[0:4]
		if v.input != nil {
			input = append(input, v.input...)
		}
		// Build call msg
		msg := core.NewCallMsg(&v.callTo, input, 0)
		evm := getNewEvm(msg, genesisBlock, genesis.Config, blockchain)

		gp := new(core.GasPool).AddGas(math.MaxUint64)
		// Ignore error
		res, _, _, _ := core.ApplyMessage(evm, msg, gp)

		if !bytes.Equal(res, v.want) {
			t.Errorf("variables %s init failed: value mismatch: got %x, want %x", v.name, res, v.want)
			continue
		}
	}
}
