// Copyright (C) 2023-2024 StorSwift Inc.
// This file is part of the PowerVoting library.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

// GoEthClient represents a client for interacting with the go-ethereum library.
type GoEthClient struct {
	Id              int64             // Identifier for the client.
	Name            string            // Name of the client.
	Rpc             string            // RPC endpoint for the client.
	RpcToken        string            // RPC token for authentication.
	Client          *ethclient.Client // Ethereum client instance.
	Abi             abi.ABI           // ABI instance for smart contract interaction.
	Amount          *big.Int          // Amount for transactions.
	GasLimit        uint64            // Gas limit for transactions.
	ChainID         *big.Int          // Chain ID for Ethereum network.
	ContractAddress common.Address    // Address of the smart contract.
	PrivateKey      *ecdsa.PrivateKey // Private key for the wallet.
	WalletAddress   common.Address    // Address of the wallet.
}

// ClientConfig represents the configuration for obtaining a go-ethereum client.
type ClientConfig struct {
	Id              int64  // Identifier for the client.
	Name            string // Name of the client.
	Rpc             string // RPC endpoint for the client.
	ContractAddress string // Address of the smart contract.
	PrivateKey      string // Private key for the wallet.
	WalletAddress   string // Address of the wallet.
	AbiPath         string // Path to the ABI file.
	GasLimit        int64  // Gas limit for transactions.
}
