package eth

import (
	"bytes"
	"math"
	"math/big"
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
	defer blockchain.Stop()
	if err != nil {
		t.Fatal(err)
	}

	test := []TestData{
		{
			"proposalPeriod",
			"getUint(bytes32)",
			common.Hex2Bytes("0176c578c6f1c60c1adf0eaaf885a402ec027dab9df2b4835205d54193184560"),
			common.BigToHash(new(big.Int).SetUint64(params.ProposalPeriod)).Bytes(),
			core.MainSystemContractAddr,
		},
		{
			"depositForProducer",
			"getUint(bytes32)",
			common.Hex2Bytes("0799480e52251d0e862f7cc87f3c9f1d3a2c8e57974ada4c8b8106a8054614f3"),
			common.BigToHash(new(big.Int).SetUint64(params.DepositForProducer)).Bytes(),
			core.MainSystemContractAddr,
		},
		{
			"lockTimeForDeposit",
			"getUint(bytes32)",
			common.Hex2Bytes("e85caaee110f28cbddd4d109b1db3c8368298ebf45e06b4db45f3e69fc48b271"),
			common.BigToHash(new(big.Int).SetUint64(params.LockTimeForDeposit)).Bytes(),
			core.MainSystemContractAddr,
		},
		{
			"producerSize",
			"getUint(bytes32)",
			common.Hex2Bytes("22841227cd558f07474bfad22019d60e7b4a93ddf6394b5f1352e24d795e0fcb"),
			common.BigToHash(new(big.Int).SetUint64(params.ProducerSize)).Bytes(),
			core.MainSystemContractAddr,
		},
		{
			"maxProducerSize",
			"getUint(bytes32)",
			common.Hex2Bytes("43b3adf5b12f6bf287de43b7625ab6911c0852ff05b0ae70f416fa07ec19f634"),
			common.BigToHash(new(big.Int).SetUint64(params.MaxProducersSize)).Bytes(),
			core.MainSystemContractAddr,
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
			"stakeForVote",
			"getUint(bytes32)",
			common.Hex2Bytes("9e2c71bfb246dedd6b6c96d1800812ab07b83677f11defcec47899ebf39bf28e"),
			common.BigToHash(new(big.Int).SetUint64(params.StakeForVote)).Bytes(),
			core.MainSystemContractAddr,
		},
		{
			"lockTimeForStake",
			"getUint(bytes32)",
			common.Hex2Bytes("06fc176e8ac7c01c651cb919fca19296a8428ba1545a80886819b1e012f2e173"),
			common.BigToHash(new(big.Int).SetUint64(params.LockTimeForStake)).Bytes(),
			core.MainSystemContractAddr,
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
			"getSystemContract(string)",
			common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001373797374656d2e766f7465436f6e747261637400000000000000000000000000"),
			core.VoteSystemContractAddr.Hash().Bytes(),
			core.MainSystemContractAddr,
		},
		{
			"regSystemContract address",
			"getSystemContract(string)",
			common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001273797374656d2e726567436f6e74726163740000000000000000000000000000"),
			core.RegSystemContractAddr.Hash().Bytes(),
			core.MainSystemContractAddr,
		},
		{
			"check vote system contract",
			"isSystemContract(address)",
			common.BytesToHash([]byte{102}).Bytes(),
			core.VoteContractBoolStorageValue.Bytes(),
			core.MainSystemContractAddr,
		},
		{
			"check reg system contract",
			"isSystemContract(address)",
			common.BytesToHash([]byte{101}).Bytes(),
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
