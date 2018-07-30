package dpos

import (
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
