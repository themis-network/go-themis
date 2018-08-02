package dpos

import (
    "github.com/themis-network/go-themis/consensus"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the proof-of-authority scheme.
type API struct {
    chain  consensus.ChainReader
    dpos   *Dpos
}
