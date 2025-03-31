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
	"github.com/ethereum/go-ethereum/common"
)

// VoterInfo struct represents information about a voter stored in the Ethereum smart contract.
type VoterInfo struct {
	ActorIds      []uint64       `json:"actorIds"`      // List of actor IDs associated with the voter.
	MinerIds      []uint64       `json:"minerIds"`      // List of miner IDs associated with the voter.
	GithubAccount string         `json:"githubAccount"` // GitHub account linked to the voter.
	EthAddress    common.Address `json:"ethAddress"`    // Ethereum address of the voter.
	UcanCid       string         `json:"ucanCid"`       // CID (Content Identifier) of the UCAN (UnixFS) associated with the voter.
}
