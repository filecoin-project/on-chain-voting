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
	"github.com/ethereum/go-ethereum/common"
)

// Vote including its ID, proposal ID, voter address, vote information, and network.
type Vote struct {
	Id         int64  `json:"id"`                       // Unique identifier
	ProposalId int64  `gorm:"constraint:off;not null"`  // Proposal ID
	Address    string `json:"address" gorm:"not null"`  // Voter address
	VoteInfo   string `json:"voteInfo" gorm:"not null"` // Vote information
	Network    int64  `json:"network" gorm:"not null"`  // Network ID
}

// VoteResult represents the result of a vote for a particular option of a proposal.
type VoteResult struct {
	Id         int64   `json:"id"`                                                        // Unique identifier
	ProposalId int64   `gorm:"constraint:off;not null;uniqueIndex:uniq_proposal_option"`  // Proposal ID
	OptionId   int64   `gorm:"uniqueIndex:uniq_proposal_option;not null" json:"optionId"` // Option ID
	Votes      float64 `json:"votes" gorm:"not null"`                                     // Number of votes
	Network    int64   `json:"network" gorm:"not null"`                                   // Network ID
}

// VoteHistory records the history of each vote cast for a proposal option.
type VoteHistory struct {
	Id         int64   `json:"id"`                       // Unique identifier
	ProposalId int64   `gorm:"constraint:off;not null"`  // Proposal ID
	Address    string  `json:"address" gorm:"not null"`  // Voter address
	OptionId   int64   `json:"optionId" gorm:"not null"` // Option ID
	Votes      float64 `json:"votes" gorm:"not null"`    // Number of votes
	Network    int64   `json:"network" gorm:"not null"`  // Network ID
}

// VoteCompleteHistory contains the complete history of votes for a proposal, including power information for each voter.
type VoteCompleteHistory struct {
	Id                    int64       `json:"id"`                                     // Unique identifier
	ProposalId            int64       `gorm:"constraint:off; not null"`               // Proposal ID
	Network               int64       `json:"network" gorm:"not null"`                // Network ID
	TotalSpPower          string      `json:"totalSpPower" gorm:"not null"`           // Total SP power
	TotalClientPower      string      `json:"totalClientPower" gorm:"not null"`       // Total client power
	TotalTokenHolderPower string      `json:"totalTokenHolderPower" gorm:"not null"`  // Total token holder power
	TotalDeveloperPower   string      `json:"totalDeveloperPower" gorm:"not null"`    // Total developer power
	VotePowers            []VotePower `json:"votePowers" gorm:"foreignKey:HistoryId"` // Vote powers
}

// ContractVote represents a vote stored on a smart contract.
type ContractVote struct {
	Voter    common.Address `json:"voter"`    // Voter address
	VoteInfo string         `json:"voteInfo"` // Vote information
}

// Vote4Counting represents a vote used for counting purposes.
type Vote4Counting struct {
	OptionId int64  // Option ID
	Votes    int64  // Number of votes
	Address  string // Voter address
}

// VotePower represents the power of a vote.
type VotePower struct {
	Id                      int64   `json:"id"`                                        // Unique identifier
	HistoryId               int64   `json:"historyId" gorm:"constraint:off; not null"` // History ID
	Address                 string  `json:"address" gorm:"not null"`                   // Voter address
	OptionId                int64   `json:"optionId" gorm:"not null"`                  // Option ID
	Votes                   float64 `json:"votes" gorm:"not null"`                     // Number of votes
	SpPower                 string  `json:"spPower" gorm:"not null"`                   // SP power
	ClientPower             string  `json:"clientPower" gorm:"not null"`               // Client power
	TokenHolderPower        string  `json:"tokenHolderPower" gorm:"not null"`          // Token holder power
	DeveloperPower          string  `json:"developerPower" gorm:"not null"`            // Developer power
	SpPowerPercent          float64 `json:"spPowerPercent" gorm:"not null"`            // SP power percentage
	ClientPowerPercent      float64 `json:"clientPowerPercent" gorm:"not null"`        // Client power percentage
	TokenHolderPowerPercent float64 `json:"tokenHolderPowerPercent" gorm:"not null"`   // Token holder power percentage
	DeveloperPowerPercent   float64 `json:"developerPowerPercent" gorm:"not null"`     // Developer power percentage
	PowerBlockHeight        int64   `json:"powerBlockHeight" gorm:"not null"`          // Power block height
}
