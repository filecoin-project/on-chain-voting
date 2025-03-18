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

type ProposalWithVoted struct {
	ProposalTbl
	Voted *bool `json:"voted,omitempty" gorm:"column:voted"` // Virtual field indicating whether to vote
}

// proposal table
type ProposalTbl struct {
	BaseField
	ProposalId          int64  `json:"proposal_id" gorm:"not null,default:0;uniqueIndex:idx_proposal_chain_id"` // Proposal ID
	Creator             string `json:"creator" gorm:"not null"`                                                 // Creator address
	GithubName          string `json:"github_name" gorm:"not null"`                                             // Creator's github name
	StartTime           int64  `json:"start_time" gorm:"not null"`                                              // Start time
	EndTime             int64  `json:"end_time" gorm:"not null"`                                                // Expiry time
	Timestamp           int64  `json:"timestamp" gorm:"not null"`                                               // Proposal create time
	Counted             int    `json:"counted" gorm:"not null,default:0"`                                       // Whether the proposal has been counted. [0: false, 1: true] 0 is not counted, 1 is counted
	ChainId             int64  `json:"chain_id" gorm:"not null;uniqueIndex:idx_proposal_chain_id"`              // Chain ID
	Title               string `json:"title" gorm:"type:longtext;not null,default:''"`                          // Name
	Content             string `json:"content" gorm:"type:longtext;not null,default:''"`                        // Descriptions
	BlockNumber         int64  `json:"block_number" gorm:"not null"`                                            // Proposal created block number
	SnapshotDay         string `json:"snapshot_day" gorm:"not null"`                                            // Snapshot day
	SnapshotBlockHeight int64  `json:"snapshot_block_height" gorm:"not null,default:0"`                         // Snapshot height
	SnapshotTimestamp   int64  `json:"snapshot_timestamp" gorm:"not null,default:0"`
	ProposalResult
	Percentage
	TotalPower
}

// Percent voting used to count proposals
type ProposalResult struct {
	ApprovePercentage float64 `json:"approve_percentage" gorm:"not null"` // Percentage of votes for the proposal
	RejectPercentage  float64 `json:"reject_percentage" gorm:"not null"`  // Percentage of votes against the proposal
}

// Percentage represents the percentage of votes for different groups.
type Percentage struct {
	TokenHolderPercentage uint16 `json:"token_holder_percentage" gorm:"not null,default:0"`  // Token holder percentage
	SpPercentage          uint16 `json:"sp_percentage" gorm:"not null,default:0"`            // SP percentage
	ClientPercentage      uint16 `json:"client_percentage" gorm:"not null,default:0"`        // Client percentage
	DeveloperPercentage   uint16 `json:"developer_percentage" gorm:"developer_percentage:0"` //	Developer percentage
}

// All voters used to count proposals
type TotalPower struct {
	TotalSpPower          string `json:"total_sp_power" gorm:"not null"`           // Total SP power
	TotalClientPower      string `json:"total_client_power" gorm:"not null"`       // Total client power
	TotalTokenHolderPower string `json:"total_token_holder_power" gorm:"not null"` // Total token holder power
	TotalDeveloperPower   string `json:"total_developer_power" gorm:"not null"`    // Total developer power
}

// Drafts of proposals are used to save proposals published by users
type ProposalDraftTbl struct {
	BaseField
	Creator   string `json:"creator" gorm:"not null;unique"`     // Creator address
	StartTime int64  `json:"start_time" gorm:"not null"`         // Start time
	EndTime   int64  `json:"end_time" gorm:"not null"`           // Expiry time
	Timezone  string `json:"timezone" gorm:"not null"`           // Proposal timezone
	ChainId   int64  `json:"chain_id" gorm:"not null"`           // Network ID
	Title     string `json:"title" gorm:"not null,default:''"`   // Name
	Content   string `json:"content" gorm:"not null,default:''"` // Descriptions
	Percentage
}
