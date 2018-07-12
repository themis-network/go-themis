// Copyright 2017 The go-themis Authors
// This file is part of the go-themis library.
//
// The go-themis library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-themis library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-themis library. If not, see <http://www.gnu.org/licenses/>.

package clique

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/themis-network/go-themis/common"
	"github.com/themis-network/go-themis/core"
	"github.com/themis-network/go-themis/core/rawdb"
	"github.com/themis-network/go-themis/core/types"
	"github.com/themis-network/go-themis/crypto"
	"github.com/themis-network/go-themis/ethdb"
	"github.com/themis-network/go-themis/params"
)

type testerVote struct {
	signer string
	voted  string
	auth   bool
}

// testerAccountPool is a pool to maintain currently active tester accounts,
// mapped from textual names used in the tests below to actual Ethereum private
// keys capable of signing transactions.
type testerAccountPool struct {
	accounts map[string]*ecdsa.PrivateKey
}

func newTesterAccountPool() *testerAccountPool {
	return &testerAccountPool{
		accounts: make(map[string]*ecdsa.PrivateKey),
	}
}

func (ap *testerAccountPool) sign(header *types.Header, signer string) {
	// Ensure we have a persistent key for the signer
	if ap.accounts[signer] == nil {
		ap.accounts[signer], _ = crypto.GenerateKey()
	}
	// Sign the header and embed the signature in extra data
	sig, _ := crypto.Sign(sigHash(header).Bytes(), ap.accounts[signer])
	copy(header.Extra[len(header.Extra)-65:], sig)
}

func (ap *testerAccountPool) address(account string) common.Address {
	// Ensure we have a persistent key for the account
	if ap.accounts[account] == nil {
		ap.accounts[account], _ = crypto.GenerateKey()
	}
	// Resolve and return the Ethereum address
	return crypto.PubkeyToAddress(ap.accounts[account].PublicKey)
}

// testerChainReader implements consensus.ChainReader to access the genesis
// block. All other methods and requests will panic.
type testerChainReader struct {
	db ethdb.Database
}

func (r *testerChainReader) Config() *params.ChainConfig                 { return params.AllCliqueProtocolChanges }
func (r *testerChainReader) CurrentHeader() *types.Header                { panic("not supported") }
func (r *testerChainReader) GetHeader(common.Hash, uint64) *types.Header { panic("not supported") }
func (r *testerChainReader) GetBlock(common.Hash, uint64) *types.Block   { panic("not supported") }
func (r *testerChainReader) GetHeaderByHash(common.Hash) *types.Header   { panic("not supported") }
func (r *testerChainReader) GetHeaderByNumber(number uint64) *types.Header {
	if number == 0 {
		return rawdb.ReadHeader(r.db, rawdb.ReadCanonicalHash(r.db, 0), 0)
	}
	panic("not supported")
}

// Tests that voting is evaluated correctly for various simple and complex scenarios.
func TestVoting(t *testing.T) {
	// Define the various voting scenarios to test
	tests := []struct {
		epoch   uint64
		signers []string
		votes   []testerVote
		results []string
	}{
		{
			// Single signer, no votes cast
			signers: []string{"A"},
			votes:   []testerVote{{signer: "A"}},
			results: []string{"A"},
		}, {
			// Three signer(Supper signer will not vote), self vote out will be accepted.
			signers: []string{"A", "B"},
			votes:   []testerVote{
				{signer: "A", voted: "A", auth:false},
			},
			results: []string{"B"},
		}, {
			// Three signer(Supper signer will not vote), voting to add two others (only accept first, second needs 3 votes)
			// Supper signer will be added automatically
			signers: []string{"A", "B"},
			votes: []testerVote{
				{signer: "A", voted: "C", auth: true},
				{signer: "B", voted: "C", auth: true},
				{signer: "C"},
				{signer: "A", voted: "D", auth: true},
				{signer: "B", voted: "D", auth: true},
			},
			results: []string{"A", "B", "C"},
		}, {
			// Three signers(Supper signer will not vote), voting to add three others (only accept first two, third needs 3 votes already)
			// Supper signer will be added automatically
			signers: []string{"A", "B", "C"},
			votes: []testerVote{
				{signer: "A", voted: "D", auth: true},
				{signer: "B", voted: "D", auth: true},
				{signer: "C", voted: "D", auth: true},
				{signer: "A", voted: "E", auth: true},
				{signer: "B", voted: "E", auth: true},
				{signer: "C", voted: "E", auth: true},
				{signer: "D"},
				{signer: "E"},
				{signer: "A", voted: "F", auth: true},
				{signer: "B", voted: "F", auth: true},
				{signer: "C", voted: "F", auth: true},
			},
			results: []string{"A", "B", "C", "D", "E"},
		}, {
			// Single signer, dropping itself (weird, but one less cornercase by explicitly allowing this)
			signers: []string{"A"},
			votes: []testerVote{
				{signer: "A", voted: "A", auth: false},
			},
			results: []string{},
		}, {
			// Three signers(Supper signer will not vote), actually needing mutual consent to drop either of them (not fulfilled)
			signers: []string{"A", "B"},
			votes: []testerVote{
				{signer: "A", voted: "B", auth: false},
			},
			results: []string{"A", "B"},
		}, {
			// Three signers(Supper signer will not vote), actually needing mutual consent to drop either of them (fulfilled)
			signers: []string{"A", "B"},
			votes: []testerVote{
				{signer: "A", voted: "B", auth: false},
				{signer: "B", voted: "B", auth: false},
			},
			results: []string{"A"},
		}, {
			// Five signers(Supper signer will not vote), three of them deciding to drop the fourth
			signers: []string{"A", "B", "C", "D"},
			votes: []testerVote{
				{signer: "A", voted: "D", auth: false},
				{signer: "B", voted: "D", auth: false},
				{signer: "C", voted: "D", auth: false},
			},
			results: []string{"A", "B", "C"},
		}, {
			// Six signers(Supper signer will not vote), consensus of three not being enough to drop anyone
			signers: []string{"A", "B", "C", "D", "E"},
			votes: []testerVote{
				{signer: "A", voted: "D", auth: false},
				{signer: "B", voted: "D", auth: false},
				{signer: "C", voted: "D", auth: false},
			},
			results: []string{"A", "B", "C", "D", "E"},
		}, {
			// Five signers(Supper signer will not vote), consensus of three already being enough to drop someone
			signers: []string{"A", "B", "C", "D"},
			votes: []testerVote{
				{signer: "A", voted: "D", auth: false},
				{signer: "B", voted: "D", auth: false},
				{signer: "C", voted: "D", auth: false},
			},
			results: []string{"A", "B", "C"},
		}, {
			// Authorizations are counted once per signer per target
			signers: []string{"A", "B"},
			votes: []testerVote{
				{signer: "A", voted: "C", auth: true},
				{signer: "B"},
				{signer: "A", voted: "C", auth: true},
				{signer: "B"},
				{signer: "A", voted: "C", auth: true},
			},
			results: []string{"A", "B"},
		}, {
			// Authorizing multiple accounts concurrently is permitted
			// Supper signer will not vote
			signers: []string{"A", "B", "C"},
			votes: []testerVote{
				{signer: "A", voted: "D", auth: true},
				{signer: "B"},
				{signer: "C"},
				{signer: "A", voted: "E", auth: true},
				{signer: "B"},
				{signer: "C", voted: "D", auth: true},
				{signer: "A"},
				{signer: "B"},
				{signer: "C", voted: "E", auth: true},
				{signer: "A"},
				{signer: "B", voted: "E", auth: true},
				{signer: "C"},
				{signer: "A"},
				{signer: "B", voted: "D", auth: true},
			},
			results: []string{"A", "B", "C", "D", "E"},
		}, {
			// Deauthorizations are counted once per signer per target
			signers: []string{"A", "B"},
			votes: []testerVote{
				{signer: "A", voted: "B", auth: false},
				{signer: "B"},
				{signer: "A", voted: "B", auth: false},
				{signer: "B"},
				{signer: "A", voted: "B", auth: false},
			},
			results: []string{"A", "B"},
		}, {
			// Deauthorizing multiple accounts concurrently is permitted
			// Supper signer will not vote
			signers: []string{"A", "B", "C", "D", "E"},
			votes: []testerVote{
				{signer: "A", voted: "D", auth: false},
				{signer: "B"},
				{signer: "C"},
				{signer: "D"},
				{signer: "A", voted: "E", auth: false},
				{signer: "B"},
				{signer: "C"},
				{signer: "D"},
				{signer: "A"},
				{signer: "B", voted: "E", auth: false},
				{signer: "C", voted: "E", auth: false},
				{signer: "D", voted: "E", auth: false},
				{signer: "B", voted: "D", auth: false},
				{signer: "C", voted: "D", auth: false},
			},
			results: []string{"A", "B", "C"},
		}, {
			// Votes from deauthorized signers are discarded immediately (deauth votes)
			// Supper signer not vote
			signers: []string{"A", "B", "C", "D", "E", "F"},
			votes: []testerVote{
				{signer: "F", voted: "B", auth: false},
				{signer: "A", voted: "F", auth: false},
				{signer: "B", voted: "F", auth: false},
				{signer: "C", voted: "F", auth: false},
				{signer: "D", voted: "F", auth: false},
				{signer: "E", voted: "F", auth: false},
				{signer: "A", voted: "B", auth: false},
				{signer: "C", voted: "B", auth: false},
				{signer: "D", voted: "B", auth: false},
			},
			results: []string{"A", "B", "C", "D", "E"},
		}, {
			// Votes from deauthorized signers are discarded immediately (auth votes)
			signers: []string{"A", "B", "C", "D", "E", "F"},
			votes: []testerVote{
				{signer: "F", voted: "B", auth: false},
				{signer: "A", voted: "F", auth: false},
				{signer: "B", voted: "F", auth: false},
				{signer: "C", voted: "F", auth: false},
				{signer: "D", voted: "F", auth: false},
				{signer: "A", voted: "B", auth: false},
				{signer: "E", voted: "B", auth: false},
				{signer: "C", voted: "B", auth: false},
			},
			results: []string{"A", "B", "C", "D", "E"},
		}, {
			// Cascading changes are not allowed, only the account being voted on may change
			signers: []string{"A", "B", "C", "D"},
			votes: []testerVote{
				{signer: "A", voted: "C", auth: false},
				{signer: "B"},
				{signer: "C"},
				{signer: "A", voted: "D", auth: false},
				{signer: "B", voted: "C", auth: false},
				{signer: "C"},
				{signer: "A"},
				{signer: "B", voted: "D", auth: false},
				{signer: "C", voted: "D", auth: false},
			},
			results: []string{"A", "B", "C"},
		}, {
			// Changes reaching consensus out of bounds (via a deauth) execute on touch
			signers: []string{"A", "B", "C", "D", "E"},
			votes: []testerVote{
				{signer: "A", voted: "C", auth: false},
				{signer: "B", voted: "C", auth: false},
				{signer: "C"},
				{signer: "E"},
				{signer: "A", voted: "D", auth: false},
				{signer: "B"},
				{signer: "C"},
				{signer: "E", voted: "C", auth: false},
				{signer: "A"},
				{signer: "B", voted: "D", auth: false},
				{signer: "C", voted: "D", auth: false},
				{signer: "E", voted: "D", auth: false},
				{signer: "A"},
				{signer: "B"},
				{signer: "C", voted: "C", auth: true},
			},
			results: []string{"A", "B", "E"},
		}, {
			// Changes reaching consensus out of bounds (via a deauth) may go out of consensus on first touch
			signers: []string{"A", "B", "C", "D"},
			votes: []testerVote{
				{signer: "A", voted: "C", auth: false},
				{signer: "B"},
				{signer: "C"},
				{signer: "A", voted: "D", auth: false},
				{signer: "B", voted: "C", auth: false},
				{signer: "C"},
				{signer: "A"},
				{signer: "B", voted: "D", auth: false},
				{signer: "C", voted: "D", auth: false},
				{signer: "A"},
				{signer: "B", voted: "C", auth: true},
			},
			results: []string{"A", "B", "C"},
		}, {
			// Ensure that pending votes don't survive authorization status changes. This
			// corner case can only appear if a signer is quickly added, removed and then
			// readded (or the inverse), while one of the original voters dropped. If a
			// past vote is left cached in the system somewhere, this will interfere with
			// the final signer outcome.(Vote will be discarded once it have been sealed in
			// a block)
			// Supper signer not vote.
			signers: []string{"A", "B", "C", "D", "E", "F"},
			votes: []testerVote{
				{signer: "A", voted: "G", auth: true}, // Authorize G, 4 votes needed
				{signer: "B", voted: "G", auth: true},
				{signer: "C", voted: "G", auth: true},
				{signer: "D", voted: "G", auth: true},
				{signer: "E", voted: "G", auth: false}, // Deauthorize G, 5 votes needed (leave A's previous vote "unchanged")
				{signer: "F", voted: "G", auth: false},
				{signer: "A", voted: "G", auth: false},
				{signer: "B", voted: "G", auth: false},
				{signer: "C", voted: "G", auth: false},
				{signer: "D", voted: "G", auth: true}, // Almost authorize G, 2/3 votes needed
				{signer: "E", voted: "G", auth: true},
				{signer: "F", voted: "G", auth: true},
				{signer: "B", voted: "A", auth: false}, // Deauthorize A, 4 votes needed
				{signer: "C", voted: "A", auth: false},
				{signer: "D", voted: "A", auth: false},
				{signer: "E", voted: "A", auth: false},
				{signer: "B", voted: "G", auth: true}, // Finish authorizing G, 3/3 votes needed
			},
			results: []string{"B", "C", "D", "E", "F", "G"},
		}, {
			// Epoch transitions reset all votes to allow chain checkpointing
			epoch:   3,
			signers: []string{"A", "B"},
			votes: []testerVote{
				{signer: "A", voted: "C", auth: true},
				{signer: "B"},
				{signer: "A"}, // Checkpoint block, (don't vote here, it's validated outside of snapshots)
				{signer: "B", voted: "C", auth: true},
			},
			results: []string{"A", "B"},
		},
	}
	// Run through the scenarios and test them
	for i, tt := range tests {
		// Create the account pool and generate the initial set of signers
		accounts := newTesterAccountPool()

		signers := make([]common.Address, len(tt.signers))
		for j, signer := range tt.signers {
			signers[j] = accounts.address(signer)
		}
		for j := 0; j < len(signers); j++ {
			for k := j + 1; k < len(signers); k++ {
				if bytes.Compare(signers[j][:], signers[k][:]) > 0 {
					signers[j], signers[k] = signers[k], signers[j]
				}
			}
		}
		// Create the genesis block with the initial set of signers
		genesis := &core.Genesis{
			ExtraData: make([]byte, extraVanity+common.AddressLength*len(signers)+extraSeal),
		}
		for j, signer := range signers {
			copy(genesis.ExtraData[extraVanity+j*common.AddressLength:], signer[:])
		}
		// Create a pristine blockchain with the genesis injected
		db := ethdb.NewMemDatabase()
		genesis.Commit(db)

		// Assemble a chain of headers from the cast votes
		headers := make([]*types.Header, len(tt.votes))
		for j, vote := range tt.votes {
			headers[j] = &types.Header{
				Number:   big.NewInt(int64(j) + 1),
				Time:     big.NewInt(int64(j) * int64(blockPeriod)),
				Coinbase: accounts.address(vote.voted),
				Extra:    make([]byte, extraVanity+extraSeal),
			}
			if j > 0 {
				headers[j].ParentHash = headers[j-1].Hash()
			}
			if vote.auth {
				copy(headers[j].Nonce[:], nonceAuthVote)
			}
			accounts.sign(headers[j], vote.signer)
		}
		// Pass all the headers through clique and ensure tallying succeeds
		head := headers[len(headers)-1]

		snap, err := New(&params.CliqueConfig{Epoch: tt.epoch}, db).snapshot(&testerChainReader{db: db}, head.Number.Uint64(), head.Hash(), headers)
		if err != nil {
			t.Errorf("test %d: failed to create voting snapshot: %v", i, err)
			continue
		}
		// Verify the final list of signers against the expected ones
		signers = make([]common.Address, len(tt.results) + 1)
		for j, signer := range tt.results {
			signers[j] = accounts.address(signer)
		}
		// Supper will be added automatically
		signers[len(tt.results)] = SupperSigner
		for j := 0; j < len(signers); j++ {
			for k := j + 1; k < len(signers); k++ {
				if bytes.Compare(signers[j][:], signers[k][:]) > 0 {
					signers[j], signers[k] = signers[k], signers[j]
				}
			}
		}
		result := snap.signers()
		if len(result) != len(signers) {
			t.Errorf("test %d: signers mismatch: have %x, want %x", i, result, signers)
			continue
		}
		for j := 0; j < len(result); j++ {
			if !bytes.Equal(result[j][:], signers[j][:]) {
				t.Errorf("test %d, signer %d: signer mismatch: have %x, want %x", i, j, result[j], signers[j])
			}
		}
	}
}