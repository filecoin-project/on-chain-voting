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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Power represents the power information.
type Power struct {
	DeveloperPower   *big.Int `json:"developerPower"`   // Developer power
	SpPower          *big.Int `json:"spPower"`          // SP power
	ClientPower      *big.Int `json:"clientPower"`      // Client power
	TokenHolderPower *big.Int `json:"tokenHolderPower"` // Token holder power
	BlockHeight      *big.Int `json:"blockHeight"`      // Block height
}

type VoterInfo struct {
	ActorIds      interface{}    `json:"actorIds"`      // Actor IDs
	MinerIds      interface{}    `json:"minerIds"`      // Miner IDs
	GithubAccount string         `json:"githubAccount"` // Github account name
	EthAddress    common.Address `json:"ethAddress"`    // Ethereum address
}

type AddrPower struct {
	Address          string   `json:"address"`
	DateStr          string   `json:"dateStr"`
	GithubAccount    string   `json:"githubAccount"`    // Github account name
	DeveloperPower   *big.Int `json:"developerPower"`   // Developer power
	SpPower          *big.Int `json:"spPower"`          // SP power
	ClientPower      *big.Int `json:"clientPower"`      // Client power
	TokenHolderPower *big.Int `json:"tokenHolderPower"` // Token holder power
	BlockHeight      int64    `json:"blockHeight"`      // Block height
}

// All computing power data structures obtained from snapshot service
type SnapshotAllPower struct {
	AddrPower []AddrPower `json:"addrPower"`
	DevPower  any         `json:"devPower"`
}

type SnapshotHeight struct {
	Height int64  `json:"height"`
	Day    string `json:"day"`
}
