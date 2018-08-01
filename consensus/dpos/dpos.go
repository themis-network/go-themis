package dpos

import (
	"math/big"
	"sync"

	"github.com/themis-network/go-themis/common"
	"github.com/themis-network/go-themis/consensus"
	"github.com/themis-network/go-themis/core/types"
	"github.com/themis-network/go-themis/core/state"
	"github.com/themis-network/go-themis/ethdb"
	"github.com/themis-network/go-themis/params"
)

var (
	blockReward    *big.Int = big.NewInt(5e+18) // Block reward in wei for successfully mining a block
)

// Dpos is the proof-of-stack consensus engine.
type Dpos struct {
	config *params.CliqueConfig // Consensus engine configuration parameters
	db     ethdb.Database       // Database to store and retrieve snapshot checkpoints

	proposals map[common.Address]bool // Current list of proposals we are pushing

	signer common.Address // Ethereum address of the signing key
	//signFn SignerFn       // Signer function to authorize hashes with
	lock sync.RWMutex // Protects the signer fields
}

//Author return the coinbase of the header.
func (d *Dpos) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (d *Dpos) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	return d.verifyHeader(chain, header, nil)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers.
func (d *Dpos) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))

	go func() {
		for i, header := range headers {
			err := d.verifyHeader(chain, header, headers[:i])

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

func (d *Dpos) verifyHeader(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	return nil
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (d *Dpos) Prepare(chain consensus.ChainReader, header *types.Header) error {
	// Set default field
	// Try to propose a new pending producers scheme when epoch start
	inturn := false
	lastheader := chain.CurrentHeader()
	for _, active := range lastheader.ActiveProducers {
		if active == header.Coinbase {
			inturn = true
			break
		}
	}
	if !inturn {
		return nil
	}

	// Try to propose a new active producers scheme when pending producers'block become IBM
	if lastheader.ProposePendingProducersBlock.Cmp(lastheader.DposIBM) <= 0 {
		header.ActiveProducers = chain.GetHeaderByNumber(lastheader.ProposePendingProducersBlock.Uint64()).PendingProducers
		header.ActiveVersion = lastheader.ActiveVersion + 1
	} else {
		header.ActiveProducers = lastheader.ActiveProducers
		header.ActiveVersion = lastheader.ActiveVersion
	}

	proposed := false
	proposedList := set{0, make([]list, 21)}
	// Try to propose a new proposedIBM block(set proposedIBM block num)
	for i := header.ProposedIBM.Uint64(); i != header.Number.Uint64(); i++ {

		iheader := chain.GetHeaderByNumber(i)
		if len(iheader.ActiveProducers) != len(header.ActiveProducers) {
			break
		}

		if proposedList.find(iheader.Coinbase[:]) {
			continue
		} else {
			proposedList.insert(iheader.Coinbase[:])
		}
		if proposedList.size > (len(iheader.ActiveProducers)*2/3 + 1) {
			proposed = true
			break
		}
	}
	// Try to propose a new dposIBM block(set dposIBM block num)
	if proposed {
		header.DposIBM = header.ProposedIBM
		header.ProposedIBM.Add(header.ProposedIBM, big.NewInt(1))
	}
	return nil
}

func (d *Dpos) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// Accumulate any block and uncle rewards and commit the final state root
	accumulateRewards(state, header)
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))

	// Header seems complete, assemble into a block and return
	return types.NewBlock(header, txs, nil, receipts), nil
}

// AccumulateRewards credits the coinbase of the given block with the mining
// reward. The total reward consists of the static block reward and rewards for
// included uncles. The coinbase of each uncle block is also rewarded.
func accumulateRewards(state *state.StateDB, header *types.Header) {
	state.AddBalance(header.Coinbase, blockReward)
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
func (d *Dpos) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return big.NewInt(1)
}

// APIs implements consensus.Engine, returning the user facing RPC API to allow
// controlling the signer voting.
func (d *Dpos) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "dpos",
		Version:   "1.0",
		Service:   &API{chain: chain, dpos: d},
		Public:    true,
	}}
}