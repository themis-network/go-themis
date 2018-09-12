package dpos

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/themis-network/go-themis/common"
	"github.com/themis-network/go-themis/core"
	"github.com/themis-network/go-themis/core/types"
	"github.com/themis-network/go-themis/core/vm"
	"github.com/themis-network/go-themis/crypto"
	"github.com/themis-network/go-themis/ethdb"
	"github.com/themis-network/go-themis/params"
)

// testerAccountPool is a pool to maintain currently active tester accounts,
// mapped from textual names used in the tests below to actual Themis private
// keys capable of signing transactions.
type testerAccountPool struct {
	accounts map[string]*ecdsa.PrivateKey
}

func newTesterAccountPool() *testerAccountPool {
	return &testerAccountPool{
		accounts: make(map[string]*ecdsa.PrivateKey),
	}
}

func (ap *testerAccountPool) sign(header *types.Header, signer string) {
	// Ensure we have a persistent key for the signer
	if ap.accounts[signer] == nil {
		ap.accounts[signer], _ = crypto.GenerateKey()
	}
	// Sign the header and embed the signature in extra data
	sig, _ := crypto.Sign(sigHash(header).Bytes(), ap.accounts[signer])
	copy(header.Extra[len(header.Extra)-65:], sig)
}

func (ap *testerAccountPool) address(account string) common.Address {
	// Ensure we have a persistent key for the account
	if ap.accounts[account] == nil {
		ap.accounts[account], _ = crypto.GenerateKey()
	}
	// Resolve and return the Themis address
	return crypto.PubkeyToAddress(ap.accounts[account].PublicKey)
}

func (ap *testerAccountPool) addresses(accounts []string) []common.Address {
	if len(accounts) == 0 {
		return []common.Address{}
	}

	res := make([]common.Address, 0, len(accounts))
	for _, account := range accounts {
		res = append(res, ap.address(account))
	}

	return res
}

type exceptIBM struct {
	ProposedIBM uint64
	DposIBM     uint64
}

// Essential field of parent for time calculation
type parent struct {
	time      uint64
	number    uint64
	producers []common.Address
	version   uint64
	miner     common.Address
}

// Essential field of grand parent for time calculation
type grand struct {
	producers []common.Address
	version   uint64
}

func TestCalculateBlockTime(t *testing.T) {
	producerA := common.BytesToAddress([]byte{1})
	producerB := common.BytesToAddress([]byte{2})
	producerC := common.BytesToAddress([]byte{3})
	producerD := common.BytesToAddress([]byte{4})
	producerE := common.BytesToAddress([]byte{5})
	producersV1 := []common.Address{
		producerA,
		producerB,
		producerC,
		producerD,
	}
	producersV2 := []common.Address{
		producerA,
		producerB,
		producerC,
		producerE,
	}

	tests := []struct {
		producer     common.Address
		parentHeader parent
		grandHeader  grand
		exceptError  error
	}{
		{
			// Calculate on genesis block
			producer: producerA,
			parentHeader: parent{
				time:      1,
				number:    0,
				producers: producersV1,
				version:   0,
			},
			grandHeader: grand{},
			exceptError: nil,
		},
		{
			// Calculate on same version
			producer: producerB,
			parentHeader: parent{
				time:      2,
				number:    1,
				producers: producersV1,
				version:   0,
				miner:     producerB,
			},
			grandHeader: grand{
				producers: producersV1,
				version:   0,
			},
			exceptError: nil,
		},
		{
			// Calculate on different version
			producer: producerE,
			parentHeader: parent{
				time:      2,
				number:    1,
				producers: producersV2,
				version:   1,
				miner:     producerB,
			},
			grandHeader: grand{
				producers: producersV1,
				version:   0,
			},
			exceptError: nil,
		},
		{
			// Current signer not authorized
			producer: producerD,
			parentHeader: parent{
				time:      2,
				number:    1,
				producers: producersV2,
				version:   1,
				miner:     producerB,
			},
			grandHeader: grand{
				producers: producersV1,
				version:   0,
			},
			exceptError: errUnauthorized,
		},
		{
			// Parent signer not authorized
			producer: producerA,
			parentHeader: parent{
				time:      2,
				number:    1,
				producers: producersV1,
				version:   0,
				miner:     producerE,
			},
			grandHeader: grand{
				producers: producersV1,
				version:   0,
			},
			exceptError: errUnauthorized,
		},
	}

	dposEngine := New(&params.DposConfig{})
	for _, tt := range tests {
		parentHeader := &types.Header{
			Time:            new(big.Int).SetUint64(tt.parentHeader.time),
			Number:          new(big.Int).SetUint64(tt.parentHeader.number),
			ActiveProducers: tt.parentHeader.producers,
			ActiveVersion:   tt.parentHeader.version,
			Coinbase:        tt.parentHeader.miner,
		}
		var grandHeader *types.Header
		if len(tt.grandHeader.producers) > 0 {
			grandHeader = &types.Header{
				ActiveProducers: tt.grandHeader.producers,
				ActiveVersion:   tt.grandHeader.version,
			}
		}

		producer := tt.producer

		blockTime, err := dposEngine.getNextBlockTime(grandHeader, parentHeader, producer)
		if err != tt.exceptError {
			t.Errorf("error not fix, have %v want %v", err, tt.exceptError)
		}

		if int64(blockTime) < time.Now().Unix() && err == nil {
			t.Errorf("error block time, block time %d less than now %d", blockTime, time.Now().Unix())
		}

		currentHeader := &types.Header{
			Time:     new(big.Int).SetUint64(blockTime),
			Coinbase: producer,
		}

		if err = dposEngine.verifyBlockTime(grandHeader, parentHeader, currentHeader); err != tt.exceptError {
			t.Errorf("verify error, have %v want %v", err, tt.exceptError)
		}
	}
}

func TestIBMCalculate(t *testing.T) {
	// Create the account pool and generate the initial set of signers
	accounts := newTesterAccountPool()
	signersID := []string{"A", "B", "C", "D"}
	producers := accounts.addresses(signersID)

	tests := []struct {
		singers []string
		wants   []exceptIBM
	}{
		{
			// Every producer just works fine
			singers: []string{"A", "B", "C", "D", "A", "B", "C", "D", "A", "B", "C", "D"},
			wants: []exceptIBM{
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 1, DposIBM: 0},
				{ProposedIBM: 2, DposIBM: 0},
				{ProposedIBM: 3, DposIBM: 0},
				{ProposedIBM: 4, DposIBM: 1},
				{ProposedIBM: 5, DposIBM: 2},
				{ProposedIBM: 6, DposIBM: 3},
				{ProposedIBM: 7, DposIBM: 4},
				{ProposedIBM: 8, DposIBM: 5},
				{ProposedIBM: 9, DposIBM: 6},
			},
		},
		{
			// D producer offline
			singers: []string{"A", "B", "C", "A", "B", "C", "A", "B", "C"},
			wants: []exceptIBM{
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 1, DposIBM: 0},
				{ProposedIBM: 2, DposIBM: 0},
				{ProposedIBM: 3, DposIBM: 0},
				{ProposedIBM: 4, DposIBM: 1},
				{ProposedIBM: 5, DposIBM: 2},
				{ProposedIBM: 6, DposIBM: 3},
			},
		},
		{
			// C,D producer offline
			singers: []string{"A", "B", "A", "B", "A", "B", "A", "B", "A", "B", "A", "B"},
			wants: []exceptIBM{
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
			},
		},
		{
			// B,C,D producer offline
			singers: []string{"A", "A", "A", "A", "A", "A", "A", "A", "A", "A", "A", "A"},
			wants: []exceptIBM{
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
			},
		},
		{
			// Network fragmentation
			singers: []string{"A", "B", "D", "A", "A", "B", "C", "C", "B", "D", "D", "A"},
			wants: []exceptIBM{
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 0, DposIBM: 0},
				{ProposedIBM: 1, DposIBM: 0},
				{ProposedIBM: 2, DposIBM: 0},
				{ProposedIBM: 2, DposIBM: 0},
				{ProposedIBM: 3, DposIBM: 0},
				{ProposedIBM: 5, DposIBM: 2},
				{ProposedIBM: 5, DposIBM: 2},
				{ProposedIBM: 5, DposIBM: 2},
				{ProposedIBM: 8, DposIBM: 5},
				{ProposedIBM: 8, DposIBM: 5},
			},
		},
	}

	gspec := core.DefaultGenesisBlock()
	gspec.Config.Dpos.Producers = producers

	for j, tt := range tests {
		// Start with a new memory db
		db := ethdb.NewMemDatabase()
		gspec.MustCommit(db)
		dpos := NewTest()
		blockchain, _ := core.NewBlockChain(db, nil, gspec.Config, dpos, vm.Config{})
		defer blockchain.Stop()

		chain, err := generateDposChain(accounts, tt.singers, 0, gspec.Config, db, blockchain, func(i int, gen *core.BlockGen) {
			gen.SetCoinbase(accounts.address(tt.singers[i]))
		})
		if err != nil {
			t.Error(err)
			continue
		}
		// Check chain can be inserted into blockchain
		if i, err := blockchain.InsertChain(chain); err != nil {
			t.Errorf("test %d insert block %d: %v\n", j, chain[i].NumberU64(), err)
			continue
		}
		// Check ibm calculation
		for i, block := range chain {
			generateHeader := block.Header()
			if tt.wants[i].DposIBM != generateHeader.DposIBM.Uint64() {
				t.Errorf("test %d block %d dposIBM error: have %d, except %d", j, generateHeader.Number, generateHeader.DposIBM, tt.wants[i].DposIBM)
			}

			if tt.wants[i].ProposedIBM != generateHeader.ProposedIBM.Uint64() {
				t.Errorf("test %d block %d ProposedIBM error: have %d, except %d", j, generateHeader.Number, generateHeader.ProposedIBM, tt.wants[i].ProposedIBM)
			}
		}
	}
}

// For dpos prepare purpose
type dposChain struct {
	blocks     []*types.Block
	blockMap   map[uint64]*types.Block
	startBlock uint64
}

func NewDposTestChain(blocks []*types.Block, parent uint64) *dposChain {
	res := &dposChain{}
	res.blocks = make([]*types.Block, len(blocks))
	copy(res.blocks, blocks)
	res.startBlock = parent + 1

	res.blockMap = make(map[uint64]*types.Block, 0)
	for _, block := range blocks {
		if block == nil {
			continue
		}
		res.blockMap[block.NumberU64()] = block
	}

	return res
}

func (d *dposChain) append(block *types.Block) {
	d.blocks = append(d.blocks, block)
	d.blockMap[block.NumberU64()] = block
}

func (d *dposChain) Blocks() []*types.Block {
	return d.blocks[d.startBlock:]
}

func (d *dposChain) Config() *params.ChainConfig {
	panic("not supported")
}

func (d *dposChain) CurrentHeader() *types.Header {
	return d.blocks[len(d.blocks)-1].Header()
}

func (d *dposChain) GetHeader(hash common.Hash, number uint64) *types.Header {
	panic("not supported")
}

func (d *dposChain) GetHeaderByNumber(number uint64) *types.Header {
	if _, ok := d.blockMap[number]; !ok {
		return nil
	}
	return d.blockMap[number].Header()
}

func (d *dposChain) GetHeaderByHash(hash common.Hash) *types.Header {
	panic("not supported")
}

func (d *dposChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	panic("not supported")
}

// generateDposChain generates blocks fit with dpos rules except pending producers which will call
// evm to get producers from contract.
func generateDposChain(accounts *testerAccountPool, signers []string, parent uint64, config *params.ChainConfig, db ethdb.Database, blockchain *core.BlockChain, gen func(i int, gen *core.BlockGen)) ([]*types.Block, error) {
	testChain := make([]*types.Block, 0)
	for i := uint64(0); i <= parent; i++ {
		block := blockchain.GetBlockByNumber(i)
		if block == nil {
			return nil, errors.New("can not get ancestor from chain")
		}
		testChain = append(testChain, block)
	}

	if len(testChain) < 1 {
		return nil, errors.New("can not get ancestor from chain")
	}

	dposEngine := NewTest()
	blocks, _ := core.GenerateChain(config, testChain[len(testChain)-1], dposEngine, db, len(signers), gen)

	testDposChain := NewDposTestChain(testChain, parent)

	for i, block := range blocks {
		dposHeader := block.Header()
		if err := dposEngine.Prepare(testDposChain, dposHeader); err != nil {
			return nil, err
		}

		if i != 0 {
			dposHeader.ParentHash = blocks[i-1].Hash()
		}

		accounts.sign(dposHeader, signers[i])
		blocks[i] = block.WithSeal(dposHeader)
		testDposChain.append(blocks[i])
	}

	return testDposChain.Blocks(), nil
}

func TestDposIBMRule(t *testing.T) {
	// Create the account pool and generate the initial set of signers
	accounts := newTesterAccountPool()
	signersID := []string{"A", "B", "C", "D"}
	producers := accounts.addresses(signersID)

	gspec := core.DefaultGenesisBlock()
	gspec.Config.Dpos.Producers = producers
	db := ethdb.NewMemDatabase()
	gspec.MustCommit(db)
	dpos := NewTest()
	blockchain, _ := core.NewBlockChain(db, nil, gspec.Config, dpos, vm.Config{})
	defer blockchain.Stop()

	// Generate basic block for blockchain, dposIBM of current header is 6 after insert.
	basicSigners := []string{"A", "B", "C", "D", "A", "B", "C", "D", "A", "B", "C", "D"}
	basicBlocks, err := generateDposChain(accounts, basicSigners, 0, gspec.Config, db, blockchain, func(i int, gen *core.BlockGen) {
		gen.SetCoinbase(accounts.address(basicSigners[i]))
	})
	if err != nil {
		t.Error(err)
		return
	}
	if _, err := blockchain.InsertChain(basicBlocks); err != nil {
		t.Error(err)
		return
	}

	testes := []struct {
		buildOn   uint64
		prev      uint64
		signers   []string
		exceptErr error
	}{
		// Try to replace original chain with longer chain
		{
			buildOn:   0,
			prev:      0,
			signers:   []string{"A", "B", "C", "D", "A", "B", "C", "D", "A", "B", "C", "D"},
			exceptErr: errInvalidBlockBeforeDposIBM,
		},
		{
			buildOn:   1,
			prev:      0,
			signers:   []string{"A", "B", "C", "D", "A", "B", "C", "D", "A", "B", "C", "D"},
			exceptErr: errInvalidBlockBeforeDposIBM,
		},
		{
			buildOn:   2,
			prev:      0,
			signers:   []string{"C", "D", "A", "B", "C", "D", "A", "B", "C", "D"},
			exceptErr: errInvalidBlockBeforeDposIBM,
		},
		// Try to insert a fork block before dpos ibm.
		{
			buildOn:   3,
			prev:      0,
			signers:   []string{"C"},
			exceptErr: errInvalidBlockBeforeDposIBM,
		},
		// Try to insert some block(already know and all on canonical chain) and some new blocks after dposIBM
		{
			buildOn:   6,
			prev:      2,
			signers:   []string{"A", "B", "C", "D"},
			exceptErr: nil,
		},
		// Try to insert some old and new blocks, but the parent of new blocks is before dposIBM
		{
			buildOn:   5,
			prev:      2,
			signers:   []string{"A", "B", "C", "D"},
			exceptErr: errInvalidBlockBeforeDposIBM,
		},
	}

	extra := []byte("make hash diff")
	for i, test := range testes {
		finalChain := make([]*types.Block, 0)
		for i := test.buildOn - 1; i >= 0 && test.buildOn-i <= test.prev; i++ {
			finalChain = append(finalChain, blockchain.GetBlockByNumber(i))
		}

		chain, err := generateDposChain(accounts, test.signers, test.buildOn, gspec.Config, db, blockchain, func(i int, gen *core.BlockGen) {
			gen.SetCoinbase(accounts.address(test.signers[i]))
			gen.SetExtra(extra)
		})
		if err != nil {
			t.Error(err)
			continue
		}

		finalChain = append(finalChain, chain...)

		if _, err := blockchain.InsertChain(finalChain); err != test.exceptErr {
			t.Error("at test ", i, ", except ", test.exceptErr, ", get ", err)
		}

	}
}

// TODO
func TestPendingProducer(t *testing.T) {
}
