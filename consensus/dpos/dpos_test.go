package dpos

import (
    "testing"
    "fmt"
    "time"
    "github.com/themis-network/go-themis/common"
)

func TestDpos_Author(t *testing.T) {
    fmt.Print(common.BytesToAddress([]byte{9}).Hex())
}

func TestRandom((t *testing.T)) {
	seed := uint64(time.Now().Unix())
	if seed <= 0 {
		t.Errorf("seed less than zero.")
	}

	rand := NewRandom(seed)

	for i := 0; i<5 ; i++ {
		fmt.Println(rand.GenRandom())
	}
}