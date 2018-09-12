package dpos

import (
	"bytes"
	"math/big"

	"github.com/deckarep/golang-set"
	"github.com/themis-network/go-themis/common"
	"github.com/themis-network/go-themis/consensus"
	"github.com/themis-network/go-themis/core/types"
	"github.com/themis-network/go-themis/log"
)

// prepareDefaultField set header fields to default values.
func (d *Dpos) prepareDefaultField(currentHeader, header *types.Header) {
	// Set the correct difficulty
	header.Difficulty = new(big.Int).Set(difficulty)

	// Ensure the extra data has all it's components
	if len(header.Extra) < extraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, extraVanity-len(header.Extra))...)
	}
	header.Extra = header.Extra[:extraVanity]
	header.Extra = append(header.Extra, make([]byte, extraSeal)...)

	// Mix digest is reserved for now, set to empty
	header.MixDigest = common.Hash{}

	// Set nonce
	copy(header.Nonce[:], nonce[:])

	// Set default pending producers and version
	header.PendingVersion = currentHeader.PendingVersion
	header.PendingProducers = make([]common.Address, len(currentHeader.PendingProducers))
	copy(header.PendingProducers, currentHeader.PendingProducers)

	// Set default active producers and version
	header.ActiveVersion = currentHeader.ActiveVersion
	header.ActiveProducers = make([]common.Address, len(currentHeader.ActiveProducers))
	copy(header.ActiveProducers, currentHeader.ActiveProducers)

	// Set proposing blcok
	header.ProposePendingProducersBlock = new(big.Int).Set(currentHeader.ProposePendingProducersBlock)
}

func (d *Dpos) tryActiveNewProducersVersion(currentHeader, header *types.Header) {
	// Try to propose a new active producers scheme when pending producers'block become IBM
	if currentHeader.ProposePendingProducersBlock.Uint64() != 0 && currentHeader.ProposePendingProducersBlock.Cmp(currentHeader.DposIBM) == 0 {

		log.Info("active new producers", "number", header.Number.Uint64(), "producers", common.PrettyAddresses(currentHeader.PendingProducers))

		// Set new producers version
		header.ActiveProducers = make([]common.Address, len(currentHeader.PendingProducers))
		copy(header.ActiveProducers, currentHeader.PendingProducers)
		header.ActiveVersion = currentHeader.ActiveVersion + 1

		// Set pending info to zero
		header.ProposePendingProducersBlock = new(big.Int).SetUint64(0)
		header.PendingProducers = []common.Address{}
	}
}

func (d *Dpos) tryUpdateIBM(chain consensus.ChainReader, currentHeader, header *types.Header) error {
	// Calculate header confirmed by 2/3 * producers + 1
	// This will be called by miner, so it is always on canonical chain
	proposedIBMHeader, updated, err := d.getLatestConfirmedHeader(chain, currentHeader, header, nil, true)
	if err != nil {
		return err
	}
	if updated {
		header.DposIBM = new(big.Int).Set(proposedIBMHeader.ProposedIBM)
		header.ProposedIBM = new(big.Int).Set(proposedIBMHeader.Number)
		return nil
	}

	// Keep same with parent
	header.DposIBM = new(big.Int).Set(currentHeader.DposIBM)
	header.ProposedIBM = new(big.Int).Set(currentHeader.ProposedIBM)
	return nil
}

func (d *Dpos) getLatestConfirmedHeader(chain consensus.ChainReader, currentHeader, header *types.Header, parents []*types.Header, onCanonical bool) (*types.Header, bool, error) {
	producerSize := len(currentHeader.ActiveProducers)
	if producerSize <= 0 {
		return nil, false, errInvalidActiveProducerList
	}

	confirmCount := mapset.NewSet()
	// Try to propose a new proposedIBM block(set proposedIBM block num)
	iheader := currentHeader
	proposedIBM := currentHeader.Number.Uint64()
	for i := proposedIBM; i != currentHeader.ProposedIBM.Uint64(); i-- {

		if onCanonical {
			iheader = chain.GetHeaderByNumber(i)
		} else {
			iheader = getHeader(chain, parents, i, header)
		}

		if iheader == nil || iheader.ActiveVersion != currentHeader.ActiveVersion {
			return currentHeader, false, nil
		}

		confirmCount.Add(iheader.Coinbase)
		if confirmCount.Cardinality() >= (producerSize*2/3 + 1) {
			return iheader, true, nil
		}
	}

	return currentHeader, false, nil
}

func (d *Dpos) tryProposePendingProducers(chain consensus.ChainReader, currentHeader, header *types.Header) error {
	// The Newest block header of last epoch
	lastEpochEndBlockNum := header.Number.Uint64()/d.config.Epoch*d.config.Epoch - 1
	lastEpochEndBlock := chain.GetHeaderByNumber(lastEpochEndBlockNum)
	// Try to propose a new pending producers scheme when epoch start or pending version not update on the first block of current epoch
	pendingVersionNotUpdated := lastEpochEndBlock != nil && (currentHeader.PendingVersion-lastEpochEndBlock.PendingVersion) < 1
	updateAllowed := header.Number.Uint64()-lastEpochEndBlockNum <= d.config.ProposeBlocksLength

	if (header.Number.Uint64()%d.config.Epoch == 0 || pendingVersionNotUpdated) && updateAllowed {
		topProducers, err := d.getPendingProducers(currentHeader)
		if err != nil {
			log.Debug("get pendingProducers failed", "number", header.Number.Uint64(), "error", err)
			return err
		}

		// Set pending producers and version
		header.PendingProducers = make([]common.Address, len(topProducers))
		copy(header.PendingProducers, topProducers)
		header.PendingVersion = currentHeader.PendingVersion + 1
		header.ProposePendingProducersBlock.Set(header.Number)
	}

	return nil
}

func (d *Dpos) verifyVersion(chain consensus.ChainReader, parents []*types.Header, parent, header *types.Header) error {
	// Verify pending producer list
	if header.PendingVersion == parent.PendingVersion && header.ActiveVersion == parent.ActiveVersion && !compareProducers(header.PendingProducers, parent.PendingProducers) {
		return errInvalidPendingProducerList
	}
	if header.PendingVersion == parent.PendingVersion && header.ActiveVersion == parent.ActiveVersion+1 && (header.ProposePendingProducersBlock.Uint64() != 0 || len(header.PendingProducers) != 0) {
		return errInvalidPendingProducerList
	}

	// Validate propose block number, the list of pending producers will be checked by VerifyPendingProducers
	if header.PendingVersion == parent.PendingVersion+1 {
		// Valid propose pending block except pending producers, which will be validated in body
		lastEpochEndBlockNum := header.Number.Uint64()/d.config.Epoch*d.config.Epoch - 1
		lastEpochEndBlock := getHeader(chain, parents, lastEpochEndBlockNum, header)
		updateAllowed := header.Number.Uint64()-lastEpochEndBlockNum <= d.config.ProposeBlocksLength
		if lastEpochEndBlock == nil || parent.PendingVersion-lastEpochEndBlock.PendingVersion != 0 || !updateAllowed {
			return errInvalidPendingProducerBlock
		}
	}
	// header.pendingVersion can not bigger than parent.PendingVersion+1 or smaller than parent.PendingVersion
	if header.PendingVersion > parent.PendingVersion+1 || header.PendingVersion < parent.PendingVersion {
		return errInvalidPendingProducersVersion
	}

	// Producer list must be same as long as version is the same
	if header.ActiveVersion == parent.ActiveVersion && !compareProducers(header.ActiveProducers, parent.ActiveProducers) {
		return errInvalidActiveProducerList
	}
	// Check active producers list
	if header.ActiveVersion == parent.ActiveVersion+1 {
		// Check parent.dpos-parent.ProposePendingProducersBlock == 0
		if parent.ProposePendingProducersBlock.Cmp(parent.DposIBM) != 0 {
			return errInvalidActiveProducerList
		}

		if !compareProducers(header.ActiveProducers, parent.PendingProducers) {
			return errInvalidActiveProducerList
		}
	}
	if header.ActiveVersion < parent.ActiveVersion || header.ActiveVersion > parent.ActiveVersion+1 {
		return errInvalidActiveProducersVersion
	}

	return nil
}

func (d *Dpos) verifyIBM(chain consensus.ChainReader, parent, header *types.Header, parents []*types.Header) error {
	// Verify ProposedIBM and DposIBM
	producerSize := len(parent.ActiveProducers)
	if producerSize <= 0 {
		return errInvalidActiveProducerList
	}

	proposedIBMHeader, updated, err := d.getLatestConfirmedHeader(chain, parent, header, parents, false)
	// Header.dposIBM and header.proposedIBM should be same with parent if not proposed
	if err != nil {
		return err
	}
	if !updated && (parent.DposIBM.Cmp(header.DposIBM) != 0 || parent.ProposedIBM.Cmp(header.ProposedIBM) != 0) {
		return errInvalidIBM
	}
	// Header.dposIBM and header.proposedIBM should be the latest confirmed block related.
	if updated {
		if proposedIBMHeader == nil || proposedIBMHeader.ProposedIBM.Cmp(header.DposIBM) != 0 || header.ProposedIBM.Uint64()-proposedIBMHeader.Number.Uint64() != 0 {
			return errInvalidIBM
		}
	}

	return nil
}
