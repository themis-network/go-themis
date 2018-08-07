package dpos

import (
    "github.com/themis-network/go-themis/consensus"
    "github.com/themis-network/go-themis/common"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the delegated-proof-of-stake scheme.
type API struct {
    chain  consensus.ChainReader
    dpos   *Dpos
}

// TODO
func (a *API) GetActiveProducers() []common.Address {
    return []common.Address{}
}
