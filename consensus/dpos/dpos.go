package dpos

import (
	"errors"
	"math/big"
	"sync"

	"github.com/hashicorp/golang-lru"
	"github.com/themis-network/go-themis/accounts"
	"github.com/themis-network/go-themis/common"
	"github.com/themis-network/go-themis/consensus"
	"github.com/themis-network/go-themis/core/state"
	"github.com/themis-network/go-themis/core/types"
	"github.com/themis-network/go-themis/crypto"
	"github.com/themis-network/go-themis/crypto/sha3"
	"github.com/themis-network/go-themis/params"
	"github.com/themis-network/go-themis/rlp"
	"github.com/themis-network/go-themis/rpc"
)

var (
	extraVanity          = 32                // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal            = 65                // Fixed number of extra-data suffix bytes reserved for signer seal
	blockReward          = big.NewInt(5e+18) // Block reward in wei for successfully mining a block
)

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	// errMissingSignature is returned if a block's extra-data section doesn't seem
	// to contain a 65 byte secp256k1 signature.
	errMissingSignature = errors.New("extra-data 65 byte suffix signature missing")
)

// SignerFn is a signer callback function to request a hash to be signed by a
// backing account.
type SignerFn func(accounts.Account, []byte) ([]byte, error)

// sigHash returns the hash which is used as input for the delegated-proof-of-stake
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	rlp.Encode(hasher, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-65], // Yes, this will panic if extra is too short
		header.MixDigest,
		header.Nonce,
		header.PendingProducers,
		header.PendingVersion,
		header.ProposePendingProducersBlock,
		header.ActiveProducers,
		header.ActiveVersion,
		header.ProposedIBM,
		header.DposIBM,
	})
	hasher.Sum(hash[:0])
	return hash
}

// ecrecover extracts the Themis account address from a signed header.
func ecrecover(header *types.Header, sigcache *lru.ARCCache) (common.Address, error) {
	// If the signature's already cached, return that
	hash := header.Hash()
	if address, known := sigcache.Get(hash); known {
		return address.(common.Address), nil
	}
	// Retrieve the signature from the header extra-data
	if len(header.Extra) < extraSeal {
		return common.Address{}, errMissingSignature
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]

	// Recover the public key and the Ethereum address
	pubkey, err := crypto.Ecrecover(sigHash(header).Bytes(), signature)
	if err != nil {
		return common.Address{}, err
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	sigcache.Add(hash, signer)
	return signer, nil
}

type CallContractFunc func(SystemCall) ([]byte, error)

// Dpos is the delegated-proof-of-stake consensus engine.
type Dpos struct {
	config          *params.DposConfig     // Consensus engine configuration parameters
	signer          common.Address         // Themis address of the signing key
	signFn          SignerFn               // Signer function to authorize hashes with
	lock            sync.RWMutex           // Protects the signer fields
	Call            CallContractFunc       // CallContractFunc is a message call func
	mainContract    *MainContract          // Main system contract for dpos
	currentContract *CurrentSystemContract // Current system contract for dpos
}

// New creates a Dpos delegated-proof-of-stake consensus engine with the initial
// signers set to the ones provided by the user.
func New(config *params.DposConfig) *Dpos {
	// Copy config
	conf := *config

	return &Dpos{
		config:          &conf,
		mainContract:    NewMainContract(),
		currentContract: NewCurrentSystemContract(),
	}
}

//Author return the coinbase of the header.
func (d *Dpos) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

// Authorize injects a private key into the consensus engine to mint new blocks
// with.
func (d *Dpos) Authorize(signer common.Address, signFn SignerFn) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.signer = signer
	d.signFn = signFn
}

// SetCallFunc set the function to execute call in evm.
func (d *Dpos) SetCallFunc(callFn CallContractFunc) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.Call = callFn
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

func (d *Dpos) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	return nil
}

func (d *Dpos) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	return nil
}

func (d *Dpos) Seal(chain consensus.ChainReader, block *types.Block, stop <-chan struct{}) (*types.Block, error) {
	return nil, nil
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
// reward. The total reward consists of the static block reward.
func accumulateRewards(state *state.StateDB, header *types.Header) {
	state.AddBalance(header.Coinbase, blockReward)
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
func (d *Dpos) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return big.NewInt(1)
}

// APIs implements consensus.Engine, returning the user facing RPC API to get system
// contract info.
func (d *Dpos) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "dpos",
		Version:   "1.0",
		Service:   &API{chain: chain, dpos: d},
		Public:    true,
	}}
}
