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

import "math/big"

type SyncPower struct {
	Address          string   `json:"address"`
	DateStr          string   `json:"dateStr"`
	GithubAccount    string   `json:"githubAccount"`
	DeveloperPower   *big.Int `json:"developerPower"`   // Developer power
	SpPower          *big.Int `json:"spPower"`          // SP power
	ClientPower      *big.Int `json:"clientPower"`      // Client power
	TokenHolderPower *big.Int `json:"tokenHolderPower"` // Token holder power
	BlockHeight      int64    `json:"blockHeight"`      // Block height
}

type AddrInfo struct {
	Addr          string   `json:"addr"`
	IdPrefix      string   `json:"addrWithPrefix"`
	ActionIDs     []uint64 `json:"actionIDs"`
	MinerIDs      []uint64 `json:"minerIDs"`
	GithubAccount string   `json:"githubAccount"`
}
