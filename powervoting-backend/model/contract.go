// Copyright (C) 2023-2024 StorSwift Inc.
// This file is part of the PowerVoting library.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GoEthClient represents the structure for interacting with the Ethereum client.
type GoEthClient struct {
	Id                  int64             // Unique identifier for the client
	Name                string            // Name of the client
	Client              *ethclient.Client // Ethereum client instance
	PowerVotingAbi      abi.ABI           // ABI (Application Binary Interface) for PowerVoting contract
	OracleAbi           abi.ABI           // ABI for Oracle contract
	PowerVotingContract common.Address    // Contract address for PowerVoting
	OracleContract      common.Address    // Contract address for Oracle
}

// ClientConfig represents the configuration for creating a GoEthClient instance.
type ClientConfig struct {
	Id                  int64  // Unique identifier for the client
	Name                string // Name of the client
	Rpc                 string // RPC endpoint for the client
	PowerVotingContract string // Contract address for PowerVoting
	OracleContract      string // Contract address for Oracle
	PowerVotingAbi      string // ABI for PowerVoting contract (as a string)
	OracleAbi           string // ABI for Oracle contract (as a string)
}
