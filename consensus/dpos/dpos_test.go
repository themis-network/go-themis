package dpos

import (
	"math/big"
	"testing"
	"time"

	"github.com/themis-network/go-themis/common"
	"github.com/themis-network/go-themis/core/types"
	"github.com/themis-network/go-themis/params"
)

var (
	producer0 = common.BytesToAddress([]byte{1})
	producer1 = common.BytesToAddress([]byte{2})
	producer2 = common.BytesToAddress([]byte{3})
	producer3 = common.BytesToAddress([]byte{4})
)

func TestCalculateBlockTime(t *testing.T) {
	producers := []common.Address{
		producer0,
		producer1,
		producer2,
		producer3,
	}
	genesisHeader := &types.Header{
		Time:            new(big.Int).SetUint64(1534327996),
		Number:          new(big.Int).SetUint64(0),
		ActiveProducers: producers,
	}
	
	dposEngine := New(&params.DposConfig{})

	blockTime1, err := dposEngine.calculateNextBlockTime(nil, genesisHeader, producer0)
	if err != nil {
		t.Errorf(err.Error())
	}

	newHeader := &types.Header{
		Time:     new(big.Int).SetUint64(blockTime1),
		Coinbase: producer0,
		Number:   new(big.Int).SetUint64(1),
	}

	// Check at a different time point
	time.Sleep(time.Second * 3)

	err = dposEngine.verifyBlockTime(nil, genesisHeader, newHeader)
	if err != nil {
		t.Errorf(err.Error())
	}

}
