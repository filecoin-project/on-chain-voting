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

import "github.com/shopspring/decimal"

// vote table
type VoteTbl struct {
	BaseField
	ProposalId       int64  `json:"proposal_id" gorm:"constraint:off;not null;uniqueIndex:idx_proposal_address"` // Proposal ID
	ChainId          int64  `json:"chain_id" gorm:"not null;uniqueIndex:idx_proposal_address"`                   // Chain ID
	Address          string `json:"address" gorm:"not null;uniqueIndex:idx_proposal_address"`                    // Voter address
	VoteEncrypted    string `json:"vote_encrypted" gorm:"type:longtext;not null"`                                // Vote encrypted
	VoteResult       string `json:"vote_result" gorm:"not null,default:''"`                                      // Vote result ["approve", "reject"]
	SpPower          string `json:"sp_power" gorm:"not null"`                                                    // SP power
	ClientPower      string `json:"client_power" gorm:"not null"`                                                // Client power
	TokenHolderPower string `json:"token_holder_power" gorm:"not null"`                                          // Token holder power
	DeveloperPower   string `json:"developer_power" gorm:"not null"`                                             // Developer power
	BlockNumber      int64  `json:"block_number" gorm:"not null"`                                                // Proposal created block number
	Timestamp        int64  `json:"timestamp" gorm:"not null"`                                                   // Proposal create time
}

type VoterInfoTbl struct {
	BaseField
	Address     string   `json:"address" gorm:"not null;uniqueIndex:idx_address_chain"`  // Voter address
	ChainId     int64    `json:"chain_id" gorm:"not null;uniqueIndex:idx_address_chain"` // Chain ID
	MinerIds    []string `json:"miner_ids" gorm:"type:json"`                             // Miner IDs
	GistId      string   `json:"gist_id" gorm:"not null"`                                // Github Gist ID
	GistInfo    string   `json:"gist_info" gorm:"type:longtext;not null"`                // Github Gist Info
	OwnerId     string   `json:"owner_id" gorm:"not null"`                               // Voter Owner ID
	GithubId    string   `json:"github_id" gorm:"not null"`                              // Github Name
	BlockNumber int64    `json:"block_number" gorm:"not null"`                           // Voter create block number
	Timestamp   int64    `json:"timestamp" gorm:"not null"`                              // Voter create time
}

// vote result count
type VoterPowerCount struct {
	SpPower        decimal.Decimal // total SP power
	ClientPower    decimal.Decimal // total client power
	TokenPower     decimal.Decimal // total token holder power
	DeveloperPower decimal.Decimal // total developer power
}

// vote percentage
type PowerPercentage struct {
	ApprovePercentage decimal.Decimal // approve vote percentage
	RejectPercentage  decimal.Decimal // reject vote percentage
}

// github gist response
type Gist struct {
	Id    string               `json:"id"`
	Files map[string]GistFiles `json:"files"`
	Owner GistOwner            `json:"owner"`
}

type GistFiles struct {
	FileName string `json:"filename"`
	Language string `json:"language"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	Size     int64  `json:"size"`
}

type GistOwner struct {
	Login string `json:"login"`
	Id    int64  `json:"id"`
}

type GistVoterInfo struct {
	SigObject    SigObject `json:"sigObject"`
	SigObjectStr string    `json:"sigObjectStr"`
	Signature    string    `json:"signature"`
}

type SigObject struct {
	GitHubName    string `json:"githubName"`
	WalletAddress string `json:"walletAddress"`
	Timestamp     int64  `json:"timestamp"`
}
