package dpos

import (
	"crypto/ecdsa"
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
	signers := make(map[string]common.Address, 0)
	producers := make([]common.Address, len(signersID))
	for i, id := range signersID {
		signers[id] = accounts.address(id)
		producers[i] = signers[id]
	}

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
	dpos := NewTest()

	for j, tt := range tests {
		// Start with a new memory db
		db := ethdb.NewMemDatabase()
		genesis := gspec.MustCommit(db)
		blockchain, _ := core.NewBlockChain(db, nil, gspec.Config, dpos, vm.Config{})
		defer blockchain.Stop()

		chain, _ := core.GenerateChain(gspec.Config, genesis, dpos, db, len(tt.singers), func(i int, gen *core.BlockGen) {
			index := i % len(tt.singers)
			gen.SetCoinbase(signers[tt.singers[index]])
		})

		for i, block := range chain {
			newHeader := block.Header()
			if err := dpos.Prepare(blockchain, newHeader); err != nil {
				t.Errorf(err.Error())
			}
			if i != 0 {
				newHeader.ParentHash = chain[i-1].Hash()
			}
			index := i % len(tt.singers)
			accounts.sign(newHeader, tt.singers[index])
			chain[i] = block.WithSeal(newHeader)

			if i, err := blockchain.InsertChain([]*types.Block{chain[i]}); err != nil {
				t.Errorf("test %d insert block %d: %v\n", j, chain[i].NumberU64(), err)
				return
			}
			if tt.wants[i].DposIBM != newHeader.DposIBM.Uint64() {
				t.Errorf("test %d block %d dposIBM error: have %d, except %d", j, newHeader.Number, newHeader.DposIBM, tt.wants[i].DposIBM)
			}

			if tt.wants[i].ProposedIBM != newHeader.ProposedIBM.Uint64() {
				t.Errorf("test %d block %d ProposedIBM error: have %d, except %d", j, newHeader.Number, newHeader.ProposedIBM, tt.wants[i].ProposedIBM)
			}
		}
	}
}

// TODO
func TestPendingProducer(t *testing.T) {

}
