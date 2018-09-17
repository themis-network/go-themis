package dpos

import (
	"errors"
	"math/big"

	"github.com/themis-network/go-themis/common"
	"github.com/themis-network/go-themis/consensus"
	"github.com/themis-network/go-themis/core"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the delegated-proof-of-stake scheme.
type API struct {
	chain consensus.ChainReader
	dpos  *Dpos
}

type ProducerInfo struct {
	Addr   common.Address `json:"addr"               gencodec:"required"`
	Weight *big.Int       `json:"weight"                  gencodec:"required"`
}

// Get all producers info of json
type ProducersInfo struct {
	Producers []ProducerInfo `json:"producers"               gencodec:"required"`
	Size      *big.Int       `json:"size"                    gencodec:"required"`
}

// Get vote info of json
type Voteinfo struct {
	Proxy     common.Address   `json:"proxy"                   gencodec:"required"`
	Producers []common.Address `json:"producers"               gencodec:"required"`
	Staked    *big.Int         `json:"staked"                  gencodec:"required"`
	Weight    *big.Int         `json:"weight"                  gencodec:"required"`
}

type ProposalInfo struct {
	Id               *big.Int       `json:"id"                      gencodec:"required"`
	Status           bool           `json:"status"                  gencodec:"required"`
	Proposer         common.Address `json:"proposer"                gencodec:"required"`
	ProposeTime      *big.Int       `json:"proposeTime"             gencodec:"required"`
	MaliciousBP      common.Address `json:"maliciousBP"             gencodec:"required"`
	Keys             [][32]byte     `json:"keys"                    gencodec:"required"`
	Values           []*big.Int     `json:"values"                  gencodec:"required"`
	Flag             uint8          `json:"flag"                    gencodec:"required"`
	ApproveVoteCount *big.Int       `json:"approveVoteCount"        gencodec:"required"`
	DisapproveCount  *big.Int       `json:"disapproveCount"         gencodec:"required"`
}

var (
	// Error info
	errInvalidInput = errors.New("invalid input")

	// Contract name
	regContract  = "system.regContract"
	voteContract = "system.voteContract"
)

func NewAPI(chain consensus.ChainReader, dpos *Dpos) *API {
	return &API{
		chain: chain,
		dpos:  dpos,
	}
}

// Get active producers of the giving block number
func (api *API) GetActiveProducers(blockNumber *big.Int) ([]common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	header := api.chain.CurrentHeader()
	if blockNumber != nil && (*blockNumber).Cmp(big.NewInt(0)) > 0 && (*blockNumber).Cmp(header.Number) < 0 {
		header = api.chain.GetHeaderByNumber(blockNumber.Uint64())
	} else if blockNumber != nil && ((*blockNumber).Cmp(big.NewInt(0)) < 0 || (*blockNumber).Cmp(header.Number) > 0) {
		return nil, errInvalidInput
	}

	// Ensure we have an actually valid block
	if header == nil {
		return nil, errUnknownBlock
	}

	return (*header).ActiveProducers, nil
}

// Get pending producer of the giving block number
func (api *API) GetPendingProducer(blockNumber *big.Int) ([]common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	header := api.chain.CurrentHeader()
	if blockNumber != nil && (*blockNumber).Cmp(big.NewInt(0)) > 0 && (*blockNumber).Cmp(header.Number) < 0 {
		header = api.chain.GetHeaderByNumber(blockNumber.Uint64())
	} else if blockNumber != nil && ((*blockNumber).Cmp(big.NewInt(0)) < 0 || (*blockNumber).Cmp(header.Number) > 0) {
		return nil, errInvalidInput
	}

	// Ensure we have an actually valid block
	if header == nil {
		return nil, errUnknownBlock
	}

	return (*header).PendingProducers, nil
}

// Get all producers info by evm
func (api *API) GetAllProducers(blockNumber *big.Int, sizeNumber *big.Int) (*ProducersInfo, error) {
	methodString := "getAllProducersInfo"
	// Retrieve the requested block number (or current if none requested)
	header := api.chain.CurrentHeader()
	if blockNumber != nil && (*blockNumber).Cmp(big.NewInt(0)) > 0 && (*blockNumber).Cmp(header.Number) < 0 {
		header = api.chain.GetHeaderByNumber(blockNumber.Uint64())
	} else if blockNumber != nil && ((*blockNumber).Cmp(big.NewInt(0)) < 0 || (*blockNumber).Cmp(header.Number) > 0) {
		return nil, errInvalidInput
	}

	// Ensure we have an actually valid block
	if header == nil {
		return nil, errUnknownBlock
	}

	// Get system contract address
	sysAddress, err := api.GetSystemContract(regContract)
	if err != nil {
		return nil, errors.New("can't get reg contract address")
	}

	caller := core.NewSystemContractCaller()
	inputData, err := caller.RegABI().Pack(methodString)
	if err != nil {
		return nil, err
	}

	call := core.NewCallMsg(sysAddress, inputData, header.Number.Uint64())
	data, err := api.dpos.Call(call)
	if err != nil {
		return nil, err
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
		return nil, err
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

	// Negative number for all producers, zero for pending producers, positive number for producers of given number.
	var getNumber uint64
	len := uint64(len(weight))
	size := (*sizeNumber).Int64()
	if size < 0 {
		getNumber = len
	} else if size == 0 {
		// Cancel proposal if can not get enough producers
		if amount.Uint64() > len {
			return nil, errTooFewProducers
		}
		getNumber = amount.Uint64()
	} else {
		if uint64(size) > len {
			getNumber = len
		} else {
			getNumber = uint64(size)
		}
	}

	// Get top producers of given number
	sortedProducers := sortTable.GetTop(getNumber)
	topProducers := make([]ProducerInfo, 0)
	for i = 0; i < getNumber; i++ {
		topProducers = append(topProducers, ProducerInfo{
			Addr:   producersAddr[sortedProducers[i].serial],
			Weight: sortedProducers[i].num,
		})
	}

	topProducersInfo := &ProducersInfo{
		Producers: topProducers,
		Size:      amount,
	}

	return topProducersInfo, nil
}

func (api *API) GetVoteInfo(addr *common.Address, blockNumber *big.Int) (*Voteinfo, error) {
	// Retrieve the requested block number (or current if none requested)
	header := api.chain.CurrentHeader()
	if blockNumber != nil && (*blockNumber).Cmp(big.NewInt(0)) > 0 && (*blockNumber).Cmp(header.Number) < 0 {
		header = api.chain.GetHeaderByNumber(blockNumber.Uint64())
	} else if blockNumber != nil && ((*blockNumber).Cmp(big.NewInt(0)) < 0 || (*blockNumber).Cmp(header.Number) > 0) {
		return nil, errInvalidInput
	}

	// Ensure we have an actually valid block
	if header == nil {
		return nil, errUnknownBlock
	}

	voteAddress, err := api.GetSystemContract(voteContract)
	if err != nil {
		return nil, errors.New("can't get vote contract address")
	}

	caller := core.NewSystemContractCaller()
	inputData, err := caller.VoteABI().Pack("getVoteInfo", addr)
	if err != nil {
		return nil, err
	}

	call := core.NewCallMsg(voteAddress, inputData, header.Number.Uint64())
	data, err := api.dpos.Call(call)
	if err != nil {
		return nil, err
	}

	var (
		ret0 = new(common.Address)
		ret1 = new([]common.Address)
		ret2 = new(*big.Int)
		ret3 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
	}

	err = caller.VoteABI().Unpack(out, "getVoteInfo", data)
	if err != nil {
		return nil, err
	}

	res := &Voteinfo{
		*ret0,
		*ret1,
		*ret2,
		*ret3,
	}

	return res, nil
}

func (api *API) GetProposal(blockNumber *big.Int) (*ProposalInfo, error) {
	// Retrieve the requested block number (or current if none requested)
	header := api.chain.CurrentHeader()
	if blockNumber != nil && (*blockNumber).Cmp(big.NewInt(0)) > 0 && (*blockNumber).Cmp(header.Number) < 0 {
		header = api.chain.GetHeaderByNumber(blockNumber.Uint64())
	} else if blockNumber != nil && ((*blockNumber).Cmp(big.NewInt(0)) < 0 || (*blockNumber).Cmp(header.Number) > 0) {
		return nil, errInvalidInput
	}

	// Ensure we have an actually valid block
	if header == nil {
		return nil, errUnknownBlock
	}

	caller := core.NewSystemContractCaller()
	inputData, err := caller.MainABI().Pack("getProposal")
	if err != nil {
		return nil, err
	}

	call := core.NewCallMsg(&core.MainSystemContractAddr, inputData, header.Number.Uint64())
	data, err := api.dpos.Call(call)
	if err != nil {
		return nil, err
	}

	var (
		ret0 = new(*big.Int)
		ret1 = new(bool)
		ret2 = new(common.Address)
		ret3 = new(*big.Int)
		ret4 = new(common.Address)
		ret5 = new([][32]byte)
		ret6 = new([]*big.Int)
		ret7 = new(uint8)
		ret8 = new(*big.Int)
		ret9 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
		ret5,
		ret6,
		ret7,
		ret8,
		ret9,
	}

	err = caller.MainABI().Unpack(out, "getProposal", data)
	if err != nil {
		return nil, err
	}

	res := &ProposalInfo{
		*ret0,
		*ret1,
		*ret2,
		*ret3,
		*ret4,
		*ret5,
		*ret6,
		*ret7,
		*ret8,
		*ret9,
	}

	return res, nil
}

func (api *API) GetSystemContract(contractName string) (*common.Address, error) {
	if contractName == "" {
		return nil, errors.New("null string")
	}
	// Get contract address from current block header
	header := api.chain.CurrentHeader()

	// Get input data for system call
	inputData, err := api.dpos.systemContract.MainABI().Pack("getSystemContract", contractName)
	if err != nil {
		return nil, err
	}

	// Get address for system contract
	call := core.NewCallMsg(&core.MainSystemContractAddr, inputData, header.Number.Uint64())
	data, err := api.dpos.Call(call)
	if err != nil {
		return nil, err
	}

	var res = new(common.Address)
	caller := core.NewSystemContractCaller()
	err = caller.MainABI().Unpack(res, "getSystemContract", data)
	if err != nil {
		return nil, err
	}

	return res, nil
}
