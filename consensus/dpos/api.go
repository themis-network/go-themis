package dpos

import (
	"encoding/binary"
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

// Get all producers info of json
type ProducersInfo struct {
	Producers []common.Address `json:"producers"               gencodec:"required"`
	Weight    []*big.Int       `json:"weight"                  gencodec:"required"`
	Size      *big.Int         `json:"size"                    gencodec:"required"`
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
	ErrInvalidInput = errors.New("invalid input")
)

// Get active producers of the giving block number
func (api *API) GetActiveProducers(blockNumber *big.Int) ([]common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	header := api.chain.CurrentHeader()
	if blockNumber == nil || (*blockNumber).Cmp(big.NewInt(0)) < 0 || (*blockNumber).Cmp(header.Number) > 0 {
		return nil, ErrInvalidInput
	} else {
		header = api.chain.GetHeaderByNumber(blockNumber.Uint64())
		return (*header).ActiveProducers, nil
	}
}

// Get pending producer of the giving block number
func (api *API) GetPendingProducer(blockNumber *big.Int) ([]common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	header := api.chain.CurrentHeader()
	if blockNumber == nil || (*blockNumber).Cmp(big.NewInt(0)) < 0 || (*blockNumber).Cmp(header.Number) > 0 {
		return nil, ErrInvalidInput
	} else {
		header = api.chain.GetHeaderByNumber(blockNumber.Uint64())
		return (*header).PendingProducers, nil
	}
}

// Get all producers info by evm
func (api *API) GetAllProducers(blockNumber *big.Int) (*ProducersInfo, error) {
	// Retrieve the requested block number (or current if none requested)
	header := api.chain.CurrentHeader()
	if blockNumber == nil || (*blockNumber).Cmp(big.NewInt(0)) < 0 || (*blockNumber).Cmp(header.Number) > 0 {
		return nil, ErrInvalidInput
	} else {
		header = api.chain.GetHeaderByNumber(blockNumber.Uint64())
	}

	// Get all producers info for system contract
	regContractAddrBytes, err := api.dpos.Call(api.dpos.systemContract.GetRegSystemContractCall(header))
	if err != nil {
		return nil, err
	}
	regContractAddr := api.dpos.systemContract.GetRegSystemContractAddress(regContractAddrBytes)
	producerInfo, err := api.dpos.Call(api.dpos.systemContract.GetAllProducersInfoCall(header, &regContractAddr))
	if err != nil {
		return nil, err
	}
	producersAddr, weight, size, err := api.dpos.systemContract.GetAllProducersInfo(producerInfo)
	if err != nil {
		return nil, err
	}

	res := &ProducersInfo{
		Producers: producersAddr,
		Weight:    weight,
		Size:      size,
	}

	return res, nil
}

func (api *API) GetVoteInfo(addr *common.Address, blockNumber *big.Int) (*Voteinfo, error) {
	if addr == nil {
		return nil, ErrInvalidInput
	}
	// Retrieve the requested block number (or current if none requested)
	header := api.chain.CurrentHeader()
	if blockNumber == nil || (*blockNumber).Cmp(big.NewInt(0)) < 0 || (*blockNumber).Cmp(header.Number) > 0 {
		return nil, ErrInvalidInput
	} else {
		header = api.chain.GetHeaderByNumber(blockNumber.Uint64())
	}

	voteAddress, err := api.GetSystemContract("system.voteContract")
	if err != nil {
		return nil, errors.New("can't get vote contract address")
	}
	methodId := "dc1e30da" //web3.sha3("getVoteInfo(address)")[:4]
	inputData := common.Hex2Bytes(methodId)
	inputData = append(inputData, make([]byte, 12)...)
	inputData = append(inputData, addr.Bytes()...)

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
	caller := core.NewSystemContractCaller()
	caller.VoteABI().Unpack(out, "getVoteInfo", data)

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
	if blockNumber == nil || (*blockNumber).Cmp(big.NewInt(0)) < 0 || (*blockNumber).Cmp(header.Number) > 0 {
		return nil, ErrInvalidInput
	} else {
		header = api.chain.GetHeaderByNumber(blockNumber.Uint64())
	}

	methodId := "b9e2bea0" //web3.sha3("getProposal()")
	inputData := common.Hex2Bytes(methodId)

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
	caller := core.NewSystemContractCaller()
	caller.MainABI().Unpack(out, "getProposal", data)

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
	methodId := "79e41595" //web3.sha3("getSystemContract(string)")
	inputData := common.Hex2Bytes(methodId)
	inputData = append(inputData, abiEncodeOfOneString(contractName)...)

	// Get address for system contract
	call := core.NewCallMsg(&core.MainSystemContractAddr, inputData, header.Number.Uint64())
	data, err := api.dpos.Call(call)
	if err != nil {
		return nil, err
	}

	var res = new(common.Address)
	caller := core.NewSystemContractCaller()
	caller.MainABI().Unpack(res, "getSystemContract", data)

	return res, nil
}

func abiEncodeOfOneString(name string) []byte {
	lenOfPaddedTo := 32
	// part1
	part1 := common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000020")
	// part2
	stringLen := len(name)
	part2 := make([]byte, lenOfPaddedTo)
	binary.LittleEndian.PutUint64(part2, uint64(stringLen))
	i := 0
	j := lenOfPaddedTo - 1
	for i < len(part2)/2 {
		part2[i], part2[j] = part2[j], part2[i]
		i++
		j--
	}
	// part3
	var needLen int
	if stringLen%lenOfPaddedTo == 0 {
		needLen = 0
	} else {
		needLen = (stringLen/lenOfPaddedTo+1)*lenOfPaddedTo - stringLen
	}
	prat3TailZero := make([]byte, needLen)
	part3 := append([]byte(name), prat3TailZero...)
	// Return result
	return append(part1, append(part2, part3...)...)
}
