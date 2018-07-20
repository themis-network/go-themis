// Copyright 2015 The go-themis Authors
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

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Themis network.
var MainnetBootnodes = []string{}

// ThemisTestnetBootnodes are the enode URLS of the P2P bootstrap nodes running on the
// Themis test network.
var ThemisTestnetBootnodes = []string{
	"enode://ff2f149dfd070c194ae427af55df7b3ee3c3c719c9d800391c90d2622811ea3d28b0718ddd9a88be4e69fa232ee9f2f0cca6cf7542b31dd7fcc7d5d1e8e602e9@47.93.163.113:30303",
	"enode://0d1e7b39d6b8a917f66eee24c682e3a67837aac404a0df20a4ec2beb461d53c9f12c0c67f433b78795bc351ed19d49455caaa8fedc328f09d390163542220194@45.249.245.140:30303",
	"enode://f2d7bc51877064f148560aee1a02bdaadfa6299a99c79633181597375d6cd3911570937445984c75ee458fd9d06d7db2c06b2a91ae4f49359d98d0c7a366f1fb@103.14.34.124:30303",
	"enode://d6d7fa6655bcacfb28dd2838a84d8261fb438949845497878d9baa0bccbcee78b74272aa7e747e0ca13488372f2f03bcda0614a356cc70de0f37326ef5731797@54.206.18.140:30303",
	"enode://55b60bc3df4d336b10b625e9e2ed5f86217fd6f8208cf98e13c6d96db142f9773167e2c4410561171fe5fc93783dd82f18c7e3359ecaced2e3ced69fbd49ea02@104.199.188.174:30303",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{}
