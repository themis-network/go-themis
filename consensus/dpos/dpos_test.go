package dpos

import (
    "testing"
    "fmt"
    "github.com/themis-network/go-themis/common"
)

func TestDpos_Author(t *testing.T) {
    fmt.Print(common.BytesToAddress([]byte{9}).Hex())
}