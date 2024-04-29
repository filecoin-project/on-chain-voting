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

package model

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GoEthClient go-ethereum client
type GoEthClient struct {
	Id                  int64
	Name                string
	Client              *ethclient.Client
	PowerVotingAbi      abi.ABI
	OracleAbi           abi.ABI
	PowerVotingContract common.Address
	OracleContract      common.Address
}

// ClientConfig config for get go-ethereum client
type ClientConfig struct {
	Id                  int64
	Name                string
	Rpc                 string
	PowerVotingContract string
	OracleContract      string
	PowerVotingAbi      string
	OracleAbi           string
}
