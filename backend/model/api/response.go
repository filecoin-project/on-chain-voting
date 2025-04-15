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

package api

import "powervoting-server/model"

// Response represents the structure of a generic response.
type Response struct {
	Code    int    `json:"code"`    // Response code
	Message string `json:"message"` // Response message
	Data    any    `json:"data"`    // Response data
}

// CountListRep represents a response containing a total count and a list of items.
type CountListRep struct {
	Total int64 `json:"total"` // Total count of items
	List  any   `json:"list"`  // List of items
}

// PowerRep represents the power distribution among different roles in the system.
type PowerRep struct {
	DeveloperPower   string `json:"developerPower"`   // Developer power
	SpPower          string `json:"spPower"`          // SP power
	ClientPower      string `json:"clientPower"`      // Client power
	TokenHolderPower string `json:"tokenHolderPower"` // Token holder power
}

// DataHeightRep represents the block height data for a specific day and chain.
type DataHeightRep struct {
	Day     string `json:"day"`         // Day of the data
	Height  int64  `json:"blockHeight"` // Block height
	ChainId int64  `json:"chainId"`     // Chain ID
}

// ProposalRep represents the details of a proposal.
type ProposalRep struct {
	ProposalId     int64                  `json:"proposalId"`               // Proposal ID
	Creator        string                 `json:"address"`                  // Creator address
	GithubName     string                 `json:"githubName"`               // Github name
	StartTime      int64                  `json:"startTime"`                // Start time
	EndTime        int64                  `json:"endTime"`                  // End time
	ChainId        int64                  `json:"chainId"`                  // Chain ID
	Title          string                 `json:"title"`                    // Proposal title
	Content        string                 `json:"content"`                  // Proposal content
	CreatedAt      int64                  `json:"createdAt"`                // Created time
	UpdatedAt      int64                  `json:"updatedAt"`                // Updated time
	Voted          bool                   `json:"voted"`                    // Whether the proposal has been voted
	Status         int                    `json:"status"`                   // Proposal status
	VotePercentage ProposalVotePercentage `json:"votePercentage,omitempty"` // Voting result percentages
	SnapshotInfo   SnapshotInfo           `json:"snapshotInfo,omitempty"`   // Snapshot information
	Percentage     ProposalPercentage     `json:"percentage,omitempty"`     // Proposal percentage
	TotalPower     TotalPower             `json:"totalPower,omitempty"`     // Total power
}

type SnapshotInfo struct {
	SnapshotDay    string `json:"snapshotDay"`    // Snapshot day
	SnapshotHeight int64  `json:"snapshotHeight"` // Snapshot height
}

type TotalPower struct {
	SpPower          string `json:"spPower"`          // SP power
	TokenHolderPower string `json:"tokenHolderPower"` // Token holder power
	DeveloperPower   string `json:"developerPower"`   // Developer power
	ClientPower      string `json:"clientPower"`      // Client power
}

// ProposalVotePercentage represents the voting percentages for a proposal.
type ProposalVotePercentage struct {
	Approve float64 `json:"approve"` // Approve percentage
	Reject  float64 `json:"reject"`  // Reject percentage
}

// ProposalDraftRep represents the draft details of a proposal.
type ProposalDraftRep struct {
	Title                 string `json:"title"`                 // Proposal title
	Content               string `json:"content"`               // Proposal content
	StartTime             int64  `json:"startTime"`             // Start time
	EndTime               int64  `json:"endTime"`               // End time
	Timezone              string `json:"timezone"`              // Timezone
	TokenHolderPercentage uint16 `json:"tokenHolderPercentage"` // Token holder percentage
	SpPercentage          uint16 `json:"spPercentage"`          // SP percentage
	DeveloperPercentage   uint16 `json:"developerPercentage"`   // Developer percentage
	ClientPercentage      uint16 `json:"clientPercentage"`      // Client percentage
}

// Voted represents the details of a vote cast on a proposal.
type Voted struct {
	ProposalId   int64  `json:"proposalId"`   // Proposal ID
	ChainId      int64  `json:"chainId"`      // Chain ID
	VoterAddress string `json:"voterAddress"` // Voter address
	VotedResult  string `json:"votedResult"`  // Voted result [approve, reject]
	Percentage   string `json:"percentage"`   // Single vote as a percentage of the entire proposal
	VotedTime    int64  `json:"votedTime"`    // Voted time
	PowerRep            // Voter power information
}

type FipProposalRep struct {
	ProposalId       int64    `json:"proposalId"`       // Proposal ID
	ChainId          int64    `json:"chainId"`          // Chain ID
	ProposalType     int      `json:"proposalType"`     // FIP proposal type
	Creator          string   `json:"creator"`          // FIP proposal creator
	CandidateAddress string   `json:"candidateAddress"` // FIP proposal candidate address
	CandidateInfo    string   `json:"candidateInfo"`    // FIP proposal candidate info
	Timestamp        int64    `json:"timestamp"`        // FIP proposal timestamp
	VotedCount       int64    `json:"votedCount"`       // FIP proposal voted count
	EditorCount      int64    `json:"editorCount"`      // FIP proposal editor count
	Status           int      `json:"status"`           // FIP proposal status
	VotedAddresss    []string `json:"votedAddresss"`    // FIP proposal voted addresss
}

type FipEditorRep struct {
	Editor    string `json:"editor"`
	ChainId   int64  `json:"chainId"`
	Timestamp int64  `json:"timestamp"`
}

type FipEditorGistInfoRep struct {
	GistId     string          `json:"gistId"`
	GistSigObj model.SigObject `json:"gistSigObj"`
}
