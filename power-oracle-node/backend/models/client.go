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

// GoEthClient go-ethereum client
type GoEthClient struct {
	Id              int64
	Name            string
	Rpc             string
	RpcToken        string
	Client          *ethclient.Client
	Abi             abi.ABI
	Amount          *big.Int
	GasLimit        uint64
	ChainID         *big.Int
	ContractAddress common.Address
	PrivateKey      *ecdsa.PrivateKey
	WalletAddress   common.Address
}

// ClientConfig config for get go-ethereum client
type ClientConfig struct {
	Id              int64
	Name            string
	Rpc             string
	ContractAddress string
	PrivateKey      string
	WalletAddress   string
	AbiPath         string
	GasLimit        int64
}
