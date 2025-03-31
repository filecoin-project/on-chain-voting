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

type FipProposalVoted struct {
	FipProposalTbl
	Voters string `json:"voters" gorm:"column:voters"`
}
// FipProposalTbl is the table for FIP proposals
type FipProposalTbl struct {
	BaseField
	ProposalId       int64  `json:"proposal_id" gorm:"uniqueIndex:idx_id"`
	ChainId          int64  `json:"chain_id" gorm:"not null;uniqueIndex:idx_id"`
	ProposalType     int    `json:"proposal_type" gorm:"not null;comment:0: revoke, 1. approve"`
	Status           int    `json:"status" gorm:"not null,default:0;comment:0: pending, 1: approved"`
	Creator          string `json:"creator" gorm:"not null"`
	CandidateAddress string `json:"candidate_address" gorm:"not null"`
	CandidateInfo    string `json:"candidate_info" gorm:"type:longtext;not null"`
	BlockNumber      int64  `json:"block_number" gorm:"not null"`
	Timestamp        int64  `json:"timestamp" gorm:"not null"`
}


// FipProposalVoteTbl is the table for FIP proposal voters
type FipProposalVoteTbl struct {
	BaseField
	ChainId     int64  `json:"chain_id" gorm:"not null"`
	ProposalId  int64  `json:"proposal_id" gorm:"not null;uniqueIndex:idx_id_voter"`
	Voter       string `json:"voter" gorm:"not null;uniqueIndex:idx_id_voter"`
	IsRemove    int    `json:"is_remove" gorm:"not null,default:0;comment:0: not remove, 1: remove"`
	BlockNumber int64  `json:"block_number" gorm:"not null"`
	Timestamp   int64  `json:"timestamp" gorm:"not null"`
}


// FipEditorTbl is the table for FIP voters
type FipEditorTbl struct {
	BaseField
	ChainId  int64  `json:"chain_id" gorm:"not null"`
	Editor   string `json:"editor" gorm:"not null;uniqueIndex"`
	IsRemove int    `json:"is_remove" gorm:"not null,default:0"`
}
