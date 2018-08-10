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
	"github.com/themis-network/go-themis/common/hexutil"
	"github.com/themis-network/go-themis/consensus"
	"github.com/themis-network/go-themis/core"
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

	epochLength = uint64(24 * 60 * 6) // Default blocks after which try to propose a new pending producers scheme
	blockPeriod = uint64(10)          // Default minimum difference between two consecutive block's timestamps

	uncleHash  = types.CalcUncleHash(nil)                 // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
	difficulty = big.NewInt(1)                            // Block difficulty for dpos
	nonce      = hexutil.MustDecode("0x0000000000000000") // Nonce number for dpos.
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

	// errMissingVanity is returned if a block's extra-data section is shorter than
	// 32 bytes, which is required to store the signer vanity.
	errMissingVanity = errors.New("extra-data 32 byte vanity prefix missing")

	// errMissingSignature is returned if a block's extra-data section doesn't seem
	// to contain a 65 byte secp256k1 signature.
	errMissingSignature = errors.New("extra-data 65 byte suffix signature missing")

	// errUnauthorized is returned if a header is signed by a non-authorized entity.
	errUnauthorized = errors.New("unauthorized")

	// errInvalidSigner is returned if coinbase is not same with signer
	errInvalidCoinbase = errors.New("coinbase not same with signer")

	// errInvalidSignerAtTimestamp is returned if signer of a block
	errInvalidSignerAtTimestamp = errors.New("invalid singer at timestamp")

	// errInvalidGrandParent is returned if grandParent is nil while parent is not genesis block
	// when calculating block time
	errInvalidGrandParent = errors.New("invalid grand parent for calculating block time")

	// errInvalidBlockTime is returned if timestamp of block is not fit consensus rule
	errInvalidBlockTime = errors.New("invalid block time")

	// errInvalidActiveProducerList is returned if active producer list is not the same
	// but active producer version is same
	errInvalidActiveProducerList = errors.New("invalid active producer list")

	// errInvalidPendingProducerList is returned if pending producer list is not the same
	// but pending producer version is same
	errInvalidPendingProducerList = errors.New("invalid pending producer list")

	// errInvalidMixDigest is returned if a block's mix digest is non-zero.
	errInvalidMixDigest = errors.New("non-zero mix digest")

	// errInvalidUncleHash is returned if a block contains an non-empty uncle list.
	errInvalidUncleHash = errors.New("non empty uncle hash")

	// errInvalidNonce if a block's nonce is non-zero
	errInvalidNonce = errors.New("non-zero nonce")
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

type CallContractFunc func(core.SystemCall) ([]byte, error)

// Dpos is the delegated-proof-of-stake consensus engine.
type Dpos struct {
	config *params.DposConfig // Consensus engine configuration parameters

	signer common.Address // Themis address of the signing key
	signFn SignerFn       // Signer function to authorize hashes with

	lock       sync.RWMutex  // Protects the signer fields
	signatures *lru.ARCCache // Signatures of recent blocks to speed up mining

	Call           CallContractFunc           // CallContractFunc is a message call func
	systemContract *core.SystemContractCaller // System contract caller for dpos to get producers' info
}

// New creates a Dpos delegated-proof-of-stake consensus engine with the initial
// signers set to the ones provided by the user.
func New(config *params.DposConfig) *Dpos {
	// Copy config
	conf := *config
	signatures, _ := lru.NewARC(inmemorySignatures)

	return &Dpos{
		config:         &conf,
		systemContract: core.NewSystemContractCaller(),
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
	if header.Number == nil {
		return errUnknownBlock
	}

	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != (common.Hash{}) {
		return errInvalidMixDigest
	}
	// Ensure that the block doesn't contain any uncles which are meaningless in dpos
	if header.UncleHash != uncleHash {
		return errInvalidUncleHash
	}

	// Don't waste time checking blocks from the future
	if header.Time.Cmp(big.NewInt(time.Now().Unix())) > 0 {
		return consensus.ErrFutureBlock
	}

	// Nonce must be 0x00..0
	if !bytes.Equal(header.Nonce[:], nonce) {
		return errInvalidNonce
	}

	// Check that the extra-data contains both the vanity and signature
	if len(header.Extra) < extraVanity {
		return errMissingVanity
	}
	if len(header.Extra) < extraVanity+extraSeal {
		return errMissingSignature
	}

	return d.verifyDposField(chain, header, parents)
}

func compareAddressList(a []common.Address, b []common.Address) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i != len(b); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (d *Dpos) verifyDposField(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}

	// TODO add verify pending producer list

	// verify the ActiveProducers
	if header.ActiveVersion == parent.ActiveVersion {
		// Producer list must be same as long as version is the same
		if !compareProducers(header.ActiveProducers, parent.ActiveProducers) {
			return errInvalidActiveProducerList
		}
	} else if header.ActiveVersion == parent.ActiveVersion+1 {
		temp := chain.GetHeaderByNumber(parent.ProposePendingProducersBlock.Uint64())
		// TODO necessary
		if temp == nil {
			for h := len(parents); h >= 0; h-- {
				if parent.ProposePendingProducersBlock.Uint64() == parents[h].Number.Uint64() {
					temp = parents[h]
				}
			}
		}
		if temp == nil {
			return errors.New("cannot find the ProposePendingProducersBlock")
		}

		// Check parent.dpos-parent.ProposePendingProducersBlock > 0
		if header.ProposePendingProducersBlock.Cmp(header.DposIBM) > 0 {
			return errors.New("wrong ActiveProducers List")
		}

		if !compareAddressList(header.ActiveProducers, temp.PendingProducers) {
			return errors.New("wrong ActiveProducers List")
		}
	} else {
		return errors.New("wrong ActiveProducers List")
	}

	// verify DposIBM
	// TODO two different situation should be verified different
	if header.DposIBM == parent.DposIBM {

	} else if header.DposIBM == parent.ProposedIBM {

	} else {
		return errors.New("wrong DposIBM")
	}

	//verify ProposedIBM
	proposed := false
	proposedList := set{0, make([]list, 21)}
	i := header.Number.Uint64()
	for ; i != header.ProposedIBM.Uint64(); i-- {

		iheader := chain.GetHeaderByNumber(i)
		// TODO Not necessary
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

		// TODO try mapping and compare version
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
	if !proposed {
		return errors.New("wrong ProposedIBM")
	}
	return d.verifySeal(chain, header, parents)
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
	return d.verifySeal(chain, header, nil)
}

func (d *Dpos) verifySeal(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
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

	// If parent's number is 0(genesis block), grandParent will get nil, so it's ok.
	var parent, grand *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}
	if len(parents) > 1 {
		grand = parents[len(parents)-2]
	} else {
		grand = chain.GetHeader(parent.ParentHash, parent.Number.Uint64()-1)
	}
	// Ensure signer signs at his time; Also check authority of signer
	if err := verifyBlockTime(grand, parent, header); err != nil {
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
	// TODO move this to seal (will cause error at here)
	// Calculate next block time, will return err is header.coinBase is unauthorized.
	myBlockTime, err := calculateNextBlockTime(chain.GetHeaderByNumber(lastHeader.Number.Uint64()-1), lastHeader, header.Coinbase)
	if err != nil {
		return err
	}
	// Set block time
	header.Time.SetUint64(myBlockTime)

	// Set the correct difficulty
	header.Difficulty = difficulty

	// Ensure the extra data has all it's components
	if len(header.Extra) < extraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, extraVanity-len(header.Extra))...)
	}
	header.Extra = header.Extra[:extraVanity]

	// Mix digest is reserved for now, set to empty
	header.MixDigest = common.Hash{}

	// Set nonce
	copy(header.Nonce[:], nonce[:])

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

	//set pendingproducers, pendingversion and ProposePendingProducersBlock
	//first set default value
	copy(header.PendingProducers, lastHeader.PendingProducers)
	header.PendingVersion = lastHeader.PendingVersion
	header.ProposePendingProducersBlock = header.Number
	//the Newest block header of last epoch
	lastEpochNumNewest := header.Number.Uint64()/epochLength*epochLength - 1
	lastEpochBlockNewest := chain.GetHeaderByNumber(lastEpochNumNewest)
	//if the first block of epoch, or pending version not update on the first block of current epoch
	if header.Number.Uint64() % epochLength == 0 || (lastHeader.PendingVersion-lastEpochBlockNewest.PendingVersion) < 1 {
		topProducers, err1 := getPendingProducers(lastHeader, d.systemContract, d.Call, lastHeader.Time.Uint64())
		if  err1 == nil {
			copy(header.PendingProducers, topProducers)
			header.PendingVersion = lastHeader.PendingVersion + 1
			header.ProposePendingProducersBlock.Add(header.ProposePendingProducersBlock, new(big.Int).SetUint64(1))
		}
	}

	return nil
}

func getPendingProducers(lastHeader *types.Header, systemContract *core.SystemContractCaller, Call CallContractFunc, seed uint64) ([]common.Address, error) {
	//get top producers info by system contract
	data, err := Call(systemContract.GetRegSystemContractCall(lastHeader))
	if err != nil {
		return nil, err
	}
	contractAddr := systemContract.GetRegSystemContractAddress(data)
	data1, err1 := Call(systemContract.GetAllProducersInfoCall(lastHeader, &contractAddr))
	if err1 != nil {
		return nil, err1
	}
	producersAddr, weightsBig, amount, err2 := systemContract.GetAllProducersInfo(data1)
	if err2 != nil {
		return nil, err2
	}

	topProducers, err3 := Shuffle(producersAddr, weightsBig, amount, seed)
	if err3 != nil {
		return nil, err3
	}
	return topProducers, nil
}

// Finalize implements consensus.Engine, ensuring no uncles are set and returns
// the final block.
func (d *Dpos) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// Accumulate any block and uncle rewards and commit the final state root
	accumulateRewards(state, header)
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = uncleHash

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
	return difficulty
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

// calculateNextBlockTime returns next block time
func calculateNextBlockTime(grandParent *types.Header, parent *types.Header, signer common.Address) (uint64, error) {
	// Assume grandParent and parent have been verified.
	currentSignerIndex, err := getSignerIndex(parent, signer)
	if err != nil {
		return 0, err
	}
	if grandParent == nil && parent.Number.Uint64() > 0 {
		return 0, errInvalidGrandParent
	}

	// Start a new active producers or first sealed block
	if parent.Number.Uint64() == 0 || grandParent.ActiveVersion != parent.ActiveVersion {
		// This block is the first block applying new active producers and will start a new epoch
		// NextBlockTime = parent.time + blockPeriod * (index + 1)
		waitBlock := currentSignerIndex + 1
		waitBlockTime := uint64(waitBlock) * blockPeriod

		return parent.Time.Uint64() + waitBlockTime, nil
	}

	// Calculate block time based on same producer version
	parentSignerIndex, err := getSignerIndex(grandParent, parent.Coinbase)
	if err != nil {
		return 0, err
	}

	waitBlock := currentSignerIndex - parentSignerIndex
	// At the next block epoch
	if waitBlock <= 0 {
		waitBlock += int64(len(parent.ActiveProducers))
	}

	waitBlockTime := uint64(waitBlock) * blockPeriod

	return parent.Time.Uint64() + waitBlockTime, nil
}

// verifyBlockTime
func verifyBlockTime(grandParent *types.Header, parent *types.Header, header *types.Header) error {
	rightTime, err := calculateNextBlockTime(grandParent, parent, header.Coinbase)
	if err != nil {
		return err
	}
	if rightTime != header.Time.Uint64() {
		return errInvalidBlockTime
	}

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

func compareProducers(src []common.Address, dst []common.Address) bool {
	if len(src) != len(dst) {
		return false
	}

	for i, v := range src {
		if v != dst[i] {
			return false
		}
	}
	return true
}
