package dpos

import (
	"math/big"
	"sync"

	"github.com/themis-network/go-themis/common"
	"github.com/themis-network/go-themis/consensus"
	"github.com/themis-network/go-themis/core/types"
	"github.com/themis-network/go-themis/ethdb"
	"github.com/themis-network/go-themis/params"
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
