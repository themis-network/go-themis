package dpos

import (
	"bytes"
	"errors"
	"math/big"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/themis-network/go-themis/accounts"
	"github.com/themis-network/go-themis/common"
	"github.com/themis-network/go-themis/consensus"
	"github.com/themis-network/go-themis/core/state"
	"github.com/themis-network/go-themis/core/types"
	"github.com/themis-network/go-themis/crypto"
	"github.com/themis-network/go-themis/crypto/sha3"
	"github.com/themis-network/go-themis/log"
	"github.com/themis-network/go-themis/params"
	"github.com/themis-network/go-themis/rlp"
	"github.com/themis-network/go-themis/rpc"
)

const (
	inmemorySignatures = 4096 // Number of recent block signatures to keep in memory
)

// Dpos delegated-proof-of-stake protocol constants.
var (
	extraVanity = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal

	epochLength = uint64(24 * 60 * 60) // Default seconds after which try to propose a new pending producers scheme
	blockPeriod = uint64(15)           // Default minimum difference between two consecutive block's timestamps

	uncleHash = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
)

var (
	blockReward = big.NewInt(5e+18) // Block reward in wei for successfully mining a block
)

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	// errUnknownBlock is returned when the list of signers is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")

	// errMissingSignature is returned if a block's extra-data section doesn't seem
	// to contain a 65 byte secp256k1 signature.
	errMissingSignature = errors.New("extra-data 65 byte suffix signature missing")

	// errUnauthorized is returned if a header is signed by a non-authorized entity.
	errUnauthorized = errors.New("unauthorized")

	// errInvalidSigner is returned if coinbase is not same with signer
	errInvalidCoinbase = errors.New("coinbase not same with signer")

	// errInvalidSignerAtTimestamp is returned if signer of a block
	errInvalidSignerAtTimestamp = errors.New("invalid singer at timestamp")
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
	config *params.DposConfig // Consensus engine configuration parameters

	signer common.Address // Themis address of the signing key
	signFn SignerFn       // Signer function to authorize hashes with

	lock       sync.RWMutex  // Protects the signer fields
	signatures *lru.ARCCache // Signatures of recent blocks to speed up mining

	Call           CallContractFunc      // CallContractFunc is a message call func
	systemContract *SystemContractCaller // System contract caller for dpos to get producers' info
}

// New creates a Dpos delegated-proof-of-stake consensus engine with the initial
// signers set to the ones provided by the user.
func New(config *params.DposConfig) *Dpos {
	// Copy config
	conf := *config
	signatures, _ := lru.NewARC(inmemorySignatures)

	return &Dpos{
		config:         &conf,
		systemContract: NewSystemContractCaller(mainSystemContractABI, regSystemContractABI),
		signatures:     signatures,
	}
}

//Author return the coinbase of the header.
func (d *Dpos) Author(header *types.Header) (common.Address, error) {
	return ecrecover(header, d.signatures)
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
	//TODO verify header is produced by inturn node
	var lastheader *types.Header
	if parents != nil {
		lastheader = parents[len(parents)-1]
	} else {
		lastheader = chain.CurrentHeader()
	}

	//verify the ActiveProducers
	if header.ActiveVersion == lastheader.ActiveVersion {
		if len(header.ActiveProducers) != len(lastheader.ActiveProducers) {
			return errors.New("wrong ActiveProducers List")
		}
		for i := 0; i != len(header.ActiveProducers); i++ {
			if header.ActiveProducers[i] != lastheader.ActiveProducers[i] {
				return errors.New("wrong ActiveProducers List")
			}
		}
	} else if header.ActiveVersion == lastheader.ActiveVersion+1 {
		temp := chain.GetHeaderByNumber(lastheader.ProposePendingProducersBlock.Uint64())
		if temp == nil {
			for _, h := range parents {
				if lastheader.ProposePendingProducersBlock.Uint64() == h.Number.Uint64() {
					temp = h
				}
			}
		}
		if temp == nil {
			return errors.New("cannot find the ProposePendingProducersBlock")
		}
		if len(header.ActiveProducers) != len(temp.PendingProducers) {
			return errors.New("wrong ActiveProducers List")
		}
		for i := 0; i != len(header.ActiveProducers); i++ {
			if header.ActiveProducers[i] != temp.PendingProducers[i] {
				return errors.New("wrong ActiveProducers List")
			}
		}
	} else {
		return errors.New("wrong ActiveProducers List")
	}

	//verify DposIBM
	if header.DposIBM == lastheader.DposIBM {

	} else if header.DposIBM == lastheader.ProposedIBM {

	} else {
		return errors.New("wrong DposIBM")
	}

	//verify ProposedIBM
	proposed := false
	proposedList := set{0, make([]list, 21)}
	i := header.Number.Uint64()
	for ; i != header.ProposedIBM.Uint64(); i-- {

		iheader := chain.GetHeaderByNumber(i)
		if iheader == nil {
			for _, h := range parents {
				if i == h.Number.Uint64() {
					iheader = h
				}
			}
		}
		if iheader == nil {
			continue
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
	if proposed {

	} else {
		return errors.New("wrong ProposedIBM")
	}
	return nil
}

// VerifyUncles implements consensus.Engine, always returning an error for any
// uncles as this consensus mechanism doesn't permit uncles.
func (d *Dpos) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errors.New("uncles not allowed")
	}
	return nil
}

func (d *Dpos) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}

	// Resolve the authorization key and check against signers
	signer, err := ecrecover(header, d.signatures)
	if err != nil {
		return err
	}
	if header.Coinbase != signer {
		return errInvalidCoinbase
	}

	// Ensure signer signs at his time
	// Also check authority of signer
	parent := chain.CurrentHeader()
	var grandParent *types.Header
	if parent.Number.Uint64() != 0 {
		grandParent = chain.GetHeaderByNumber(parent.Number.Uint64() - 1)
	}
	if err := verifyBlockTime(grandParent, parent, header.Coinbase); err != nil {
		return err
	}

	return nil
}

func (d *Dpos) Seal(chain consensus.ChainReader, block *types.Block, stop <-chan struct{}) (*types.Block, error) {
	header := block.Header()

	// Sealing the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return nil, errUnknownBlock
	}

	// Don't hold the signer fields for the entire sealing procedure
	d.lock.RLock()
	signer, signFn := d.signer, d.signFn
	d.lock.RUnlock()

	if !bytes.Equal(signer[:], header.Coinbase[:]) {
		return nil, errInvalidCoinbase
	}

	if _, err := getSignerIndex(chain.CurrentHeader(), header.Coinbase); err != nil {
		return nil, err
	}

	// Sweet, the protocol permits us to sign the block, wait for our time
	delay := time.Unix(header.Time.Int64(), 0).Sub(time.Now()) // nolint: gosimple
	log.Trace("Waiting for slot to sign and propagate", "delay", common.PrettyDuration(delay))

	select {
	case <-stop:
		return nil, nil
	case <-time.After(delay):
	}

	// Sign all the things!
	sighash, err := signFn(accounts.Account{Address: signer}, sigHash(header).Bytes())
	if err != nil {
		return nil, err
	}
	copy(header.Extra[len(header.Extra)-extraSeal:], sighash)

	return block.WithSeal(header), nil
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (d *Dpos) Prepare(chain consensus.ChainReader, header *types.Header) error {
	// Set default field
	// Try to propose a new pending producers scheme when epoch start
	lastHeader := chain.CurrentHeader()
	if _, err := getSignerIndex(lastHeader, header.Coinbase); err != nil {
		return err
	}

	// Try to propose a new active producers scheme when pending producers'block become IBM
	if lastHeader.ProposePendingProducersBlock.Cmp(lastHeader.DposIBM) <= 0 {
		header.ActiveProducers = chain.GetHeaderByNumber(lastHeader.ProposePendingProducersBlock.Uint64()).PendingProducers
		header.ActiveVersion = lastHeader.ActiveVersion + 1
	} else {
		header.ActiveProducers = lastHeader.ActiveProducers
		header.ActiveVersion = lastHeader.ActiveVersion
	}

	proposed := false
	proposedList := set{0, make([]list, 21)}
	// Try to propose a new proposedIBM block(set proposedIBM block num)
	i := header.Number.Uint64()
	for ; i != header.ProposedIBM.Uint64(); i-- {

		iheader := chain.GetHeaderByNumber(i)
		if iheader.ActiveVersion != header.ActiveVersion {
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
		header.ProposedIBM.SetUint64(i)
	}
	return nil
}

// Finalize implements consensus.Engine, ensuring no uncles are set and returns
// the final block.
func (d *Dpos) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// Accumulate any block and uncle rewards and commit the final state root
	accumulateRewards(state, header)
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))

	// Header seems complete, assemble into a block and return
	return types.NewBlock(header, txs, nil, receipts), nil
}

// AccumulateRewards credits the coinbase of the given block with the producing
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

// calcualteNextBlockTime returns next block time
func calcualteNextBlockTime(grandParent *types.Header, parent *types.Header, signer common.Address) (*big.Int, error) {
	// Assume grandParent and parent have been verified.
	currentSignerIndex, err := getSignerIndex(parent, signer)
	if err != nil {
		return nil, err
	}

	// Start a new active producers or first sealed block
	if parent.Number.Uint64() == 0 || grandParent.ActiveVersion != parent.ActiveVersion {
		// This block is the first block applying new active producers and will start a new epoch
		// NextBlockTime = parent.time + blockPeriod * (index + 1)
		waitBlock := currentSignerIndex + 1
		waitBlockTime := uint64(waitBlock) * blockPeriod

		return parent.Time.Add(parent.Time, new(big.Int).SetUint64(waitBlockTime)), nil
	}

	// Calculate block time based on same producer version
	parentSignerIndex, err := getSignerIndex(grandParent, parent.Coinbase)
	if err != nil {
		return nil, err
	}

	waitBlock := currentSignerIndex - parentSignerIndex
	// At the next block epoch
	if waitBlock <= 0 {
		waitBlock += int64(len(parent.ActiveProducers))
	}

	waitBlockTime := uint64(waitBlock) * blockPeriod

	return parent.Time.Add(parent.Time, new(big.Int).SetUint64(waitBlockTime)), nil
}

// verifyBlockTime
func verifyBlockTime(grandParent *types.Header, parent *types.Header, signer common.Address) error {
	return nil
}

// getSignerIndex returns index of signer in header.activeProducers, otherwise an error
// will return
func getSignerIndex(header *types.Header, signer common.Address) (int64, error) {
	var signerIndex int64
	inturn := false
	for i, active := range header.ActiveProducers {
		if active == signer {
			inturn = true
			signerIndex = int64(i)
			break
		}
	}
	if !inturn {
		return 0, errUnauthorized
	}

	return signerIndex, nil
}
