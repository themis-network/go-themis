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

// Mode defines the type and amount of Dpos verification an dpos engine makes.
type Mode uint

const (
	ModeNormal Mode = iota
	ModeTest
)

// Dpos delegated-proof-of-stake protocol constants.
var (
	extraVanity = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal

	epochLength         = uint64(24 * 60 * 6) // Default blocks after which try to propose a new pending producers scheme
	blockPeriod         = uint64(10)          // Default minimum difference between two consecutive block's timestamps
	proposeBlocksLength = uint64(21)          // Default block length lasts trying to propose a new pending version if failed

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

	// errInvalidNonce if a block's nonce is non-zero.
	errInvalidNonce = errors.New("non-zero nonce")

	// errInvalidPendingProducersVersion is returned if pending producers version bigger
	// than parent.pendingProducersVersion+1 or less than parent.endingProducersVersion
	errInvalidPendingProducersVersion = errors.New("invalid pending producers version")

	// errInvalidActiveProducersVersion is returned if active producers version bigger
	// than parent.activeProducersVersion+1 or less than parent.activeProducersVersion
	errInvalidActiveProducersVersion = errors.New("invalid active producers version")

	// errInvalidPendingProducerBlock is returned if pending producer list isn't proposed
	// second time during one block epoch.
	errInvalidPendingProducerBlock = errors.New("invalid pending producer block")

	// errInvalidIBM is returned if irreversible block num isn't compliance with the rules
	errInvalidIBM = errors.New("invalid dpos irreversible block num")

	// errTooFewProducers is returned if can't get enough producers from system contract.
	errTooFewProducers = errors.New("too few producers")

	// errInvalidBlockBeforeDposIBM is returned when dpos receive a block before or equal dpos
	// ibm.
	errInvalidBlockBeforeDposIBM = errors.New("invalid block before dpos ibm and should be ignored")
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
	config   *params.DposConfig // Consensus engine configuration parameters
	DposMode Mode               // Mode used in dpos

	signer common.Address // Themis address of the signing key
	signFn SignerFn       // Signer function to authorize hashes with

	lock       sync.RWMutex  // Protects the signer fields
	signatures *lru.ARCCache // Signatures of recent blocks to speed up mining

	Call           CallContractFunc           // CallContractFunc is a message call func
	systemContract *core.SystemContractCaller // System contract caller for dpos to get producers' info
	rand           *Random                    // Shuffle the pending producers
}

// New creates a Dpos delegated-proof-of-stake consensus engine with the initial
// signers set to the ones provided by the user.
func New(config *params.DposConfig) *Dpos {
	// Copy config
	conf := *config
	// Set default config when not set
	if conf.BlockPeriod == 0 {
		conf.BlockPeriod = blockPeriod
	}
	if conf.Epoch == 0 {
		conf.Epoch = epochLength
	}
	if conf.ProposeBlocksLength == 0 {
		conf.ProposeBlocksLength = proposeBlocksLength
	}
	// Allocate caches and create the engine
	signatures, _ := lru.NewARC(inmemorySignatures)

	return &Dpos{
		config:         &conf,
		systemContract: core.NewSystemContractCaller(),
		signatures:     signatures,
		rand:           NewRandom(0),
		DposMode:       ModeNormal,
	}
}

func NewTest() *Dpos {
	conf := &params.DposConfig{}
	conf.BlockPeriod = blockPeriod
	conf.Epoch = epochLength
	conf.ProposeBlocksLength = proposeBlocksLength
	// Allocate caches and create the engine
	signatures, _ := lru.NewARC(inmemorySignatures)

	return &Dpos{
		config:         conf,
		systemContract: core.NewSystemContractCaller(),
		signatures:     signatures,
		rand:           NewRandom(0),
		DposMode:       ModeTest,
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

func (d *Dpos) verifyDposField(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}

	parent := getHeader(chain, parents, header.Number.Uint64()-1, header)
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}

	// Reject any block before current dpos ibm
	currentHeader := chain.CurrentHeader()
	// The block should on canonical chain if number <= current dpos IBM.
	if number <= currentHeader.DposIBM.Uint64() {
		localHeader := chain.GetHeaderByNumber(number)
		if localHeader == nil || localHeader.Hash() != header.Hash() {
			return errInvalidBlockBeforeDposIBM
		}
	}

	if err := d.verifyVersion(chain, parents, parent, header); err != nil {
		return err
	}

	if err := d.verifyIBM(chain, parent, header, parents); err != nil {
		return err
	}

	return d.verifySeal(chain, header, parents)
}

// VerifyPendingProducer check pending producers after parent block state have been written
// cause it have to call evm to get producers
func (d *Dpos) VerifyPendingProducers(chain consensus.ChainReader, header *types.Header) error {
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	if header.PendingVersion != parent.PendingVersion+1 {
		return nil
	}

	newProducers, err := d.getPendingProducers(chain, parent)
	if err != nil || !compareProducers(header.PendingProducers, newProducers) {
		return errInvalidPendingProducerList
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
	parent := getHeader(chain, parents, number-1, header)
	grand := getHeader(chain, parents, number-2, header)
	// Ensure signer signs at his time; Also check authority of signer
	if err := d.verifyBlockTime(grand, parent, header); err != nil {
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

	// Set default field
	lastHeader := chain.CurrentHeader()
	// Just check signer can seal.
	_, err := getSignerIndex(lastHeader, signer)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(signer[:], header.Coinbase[:]) {
		return nil, errInvalidCoinbase
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
	currentHeader := chain.CurrentHeader()
	d.prepareDefaultField(currentHeader, header)

	d.tryActiveNewProducersVersion(currentHeader, header)

	// Try to update proposeIBM/dposIBM by calculating latest confirmed block
	if err := d.tryUpdateIBM(chain, currentHeader, header); err != nil {
		return err
	}

	// Ignore error since it's not mandatory to propose new producers
	d.tryProposePendingProducers(chain, currentHeader, header)

	// Get next block time, will return err is header.coinBase is unauthorized.
	myBlockTime, err := d.getNextBlockTime(chain.GetHeaderByNumber(currentHeader.Number.Uint64()-1), currentHeader, header.Coinbase)
	if err != nil {
		return nil
	}
	// Set block time
	header.Time.SetUint64(myBlockTime)

	return nil
}

func (d *Dpos) getPendingProducers(chain consensus.ChainReader, lastHeader *types.Header) ([]common.Address, error) {
	// Get all pending producers info
	topProducersInfo, _, amount, err := d.GetAllSortedProducers(chain, lastHeader)
	if err != nil {
		return nil, err
	}

	if len(topProducersInfo) < int(amount.Int64()) {
		return nil, errTooFewProducers
	}

	topProducers := make([]common.Address, 0)
	for i := 0; i < int(amount.Int64()); i++ {
		topProducers = append(topProducers, topProducersInfo[i])
	}

	// Get pseudo-random order
	d.rand.ResetSeed(lastHeader.Number.Uint64())
	d.rand.Shuffle(topProducers)
	return topProducers, nil
}

func (d *Dpos) GetAllSortedProducers(chain consensus.ChainReader, header *types.Header) ([]common.Address, []*big.Int, *big.Int, error) {
	methodString := "getAllProducersInfo"

	if header == nil {
		return nil, nil, nil, errUnknownBlock
	}

	// Get system contract address
	sysAddress, err := NewAPI(chain, d).GetSystemContract(regContract)
	if err != nil {
		return nil, nil, nil, errors.New("can't get reg contract address")
	}

	caller := core.NewSystemContractCaller()
	inputData, err := caller.RegABI().Pack(methodString)
	if err != nil {
		return nil, nil, nil, err
	}

	call := core.NewCallMsg(sysAddress, inputData, header.Number.Uint64())
	data, err := d.Call(call)
	if err != nil {
		return nil, nil, nil, err
	}

	var (
		ret0 = new([]common.Address)
		ret1 = new([]*big.Int)
		ret2 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}

	err = caller.RegABI().Unpack(out, methodString, data)
	if err != nil {
		return nil, nil, nil, err
	}

	// Get all producers info
	producersAddr := *ret0
	weight := *ret1
	amount := *ret2

	// Sort all weight of producers
	var i uint64
	sortTable := sortNumSlice{}
	for i, voteWeight := range weight {
		sortTable = append(sortTable, &sortNum{i, voteWeight})
	}

	// Sort All producers
	allLen := uint64(len(weight))
	sortedWeight := sortTable.GetTop(allLen)
	producers := make([]common.Address, 0)
	weights := make([]*big.Int, 0)
	for i = 0; i < allLen; i++ {
		producers = append(producers, producersAddr[sortedWeight[i].serial])
		weights = append(weights, sortedWeight[i].num)

	}

	return producers, weights, amount, nil
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

// Close implements consensus.Engine. It's a noop for dpos as there is are no background threads.
func (d *Dpos) Close() error {
	return nil
}

func (d *Dpos) getNextBlockTime(grandParent *types.Header, parent *types.Header, signer common.Address) (uint64, error) {
	blockTime, err := d.calculateNextBlockTime(grandParent, parent, signer)
	if err != nil {
		return 0, err
	}

	// Get block time smaller than local time, add blockPeriod * producerSize to reach localTime
	// Only used in normal mode
	if blockTime < uint64(time.Now().Unix()) && d.DposMode == ModeNormal {
		blockTime = d.reachAfterLocalTime(blockTime, uint64(time.Now().Unix()), len(parent.ActiveProducers))
	}

	return blockTime, nil
}

// calculateNextBlockTime returns next block time
func (d *Dpos) calculateNextBlockTime(grandParent *types.Header, parent *types.Header, signer common.Address) (uint64, error) {
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
		waitBlockTime := uint64(waitBlock) * d.config.BlockPeriod

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

	waitBlockTime := uint64(waitBlock) * d.config.BlockPeriod

	return parent.Time.Uint64() + waitBlockTime, nil
}

// verifyBlockTime
func (d *Dpos) verifyBlockTime(grandParent *types.Header, parent *types.Header, header *types.Header) error {
	// Assume grandParent and parent have been verified.
	currentSignerIndex, err := getSignerIndex(parent, header.Coinbase)
	if err != nil {
		return err
	}

	if grandParent == nil && parent.Number.Uint64() > 0 {
		return errInvalidGrandParent
	}

	producerSize := uint64(len(parent.ActiveProducers))

	// Verify first block of a new active producer version sealed specially
	if parent.Number.Uint64() == 0 || grandParent.ActiveVersion != parent.ActiveVersion {
		waitBlock := currentSignerIndex + 1
		waitBlockTime := uint64(waitBlock) * d.config.BlockPeriod
		if (header.Time.Uint64()-waitBlockTime-parent.Time.Uint64())%(producerSize*d.config.BlockPeriod) != 0 {
			return errInvalidBlockTime
		}
		return nil
	}

	// Calculate block time based on same producer version
	parentSignerIndex, err := getSignerIndex(grandParent, parent.Coinbase)
	if err != nil {
		return err
	}

	waitBlock := currentSignerIndex - parentSignerIndex
	// At the next block epoch
	if waitBlock <= 0 {
		waitBlock += int64(len(parent.ActiveProducers))
	}
	waitBlockTime := uint64(waitBlock) * d.config.BlockPeriod
	if (header.Time.Uint64()-waitBlockTime-parent.Time.Uint64())%(producerSize*d.config.BlockPeriod) != 0 {
		return errInvalidBlockTime
	}

	return nil
}

func (d *Dpos) reachAfterLocalTime(originalBlockTime, localTime uint64, producerSize int) uint64 {
	distance := (localTime-originalBlockTime)/d.config.BlockPeriod/uint64(producerSize) + 1
	return distance*d.config.BlockPeriod*uint64(producerSize) + originalBlockTime
}

// getSignerIndex returns index of signer in header.activeProducers, otherwise an error
// will be returned
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

// getHeader try to get header from parents, then try to get from chain
func getHeader(chain consensus.ChainReader, parents []*types.Header, number uint64, header *types.Header) *types.Header {
	parentsSize := uint64(len(parents))
	distance := header.Number.Uint64() - number
	if parentsSize >= distance {
		return parents[parentsSize-distance]
	}

	// Reach to farthest ancestor
	ancestor := types.CopyHeader(header)
	if parentsSize > 0 {
		ancestor = parents[0]
		distance -= parentsSize
	}

	for i := uint64(0); i < distance; i++ {
		if ancestor == nil {
			return nil
		}
		ancestor = chain.GetHeader(ancestor.ParentHash, ancestor.Number.Uint64()-1)
	}

	return ancestor
}
