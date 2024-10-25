package model

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

import (
	"time"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Proposal represents the structure of a proposal.
type Proposal struct {
	Id           int64     `json:"id"`                                                    // Unique identifier
	ProposalId   int64     `json:"proposalId" gorm:"default:0;"`                          // Proposal ID
	Cid          string    `json:"cid" gorm:"not null;index:uniq_cid,unique"`             // CID
	ProposalType int64     `json:"proposalType" gorm:"not null,default:0"`                // Proposal type
	Creator      string    `json:"creator" gorm:"not null"`                               // Creator address
	StartTime    int64     `json:"startTime" gorm:"not null"`                             // Start time
	ExpTime      int64     `json:"expTime" gorm:"not null"`                               // Expiry time
	VoteCount    int64     `json:"voteCount" gorm:"not null,default:0"`                   // Vote count
	Status       int       `json:"status" gorm:"not null,default:0"`                      // Proposal status
	Network      int64     `json:"network" gorm:"not null"`                               // Network ID
	Name         string    `json:"name" gorm:"not null,default:''"`                       // Name
	Timezone     string    `json:"timezone" gorm:"not null"`                              // Timezone
	Descriptions string    `json:"descriptions" gorm:"not null,default:'';type:longtext"` // Descriptions
	GithubName   string    `json:"githubName" gorm:"not null,default:''"`                 // Github name
	GithubAvatar string    `json:"githubAvatar" gorm:"not null,default:''"`               // Github avatar
	GMTOffset    string    `json:"gmtOffset" gorm:"not null,default:''"`                  // GMT offset
	CurrentTime  int64     `json:"currentTime" gorm:"not null,default:0"`                 // Current time
	CreatedAt    time.Time `json:"createdAt" gorm:"not null,autoCreateTime"`              // Created time
	UpdatedAt    time.Time `json:"updatedAt" gorm:"not null,autoUpdateTime"`              // Updated time
}

type ProposalList struct {
	Id           int64        `json:"id"`                                                    // Unique identifier
	ProposalId   int64        `json:"proposalId" gorm:"default:0"`                           // Proposal ID
	Cid          string       `json:"cid" gorm:"not null"`                                   // CID
	ProposalType int64        `json:"proposalType" gorm:"not null,default:0"`                // Proposal type
	Creator      string       `json:"creator" gorm:"not null"`                               // Creator address
	StartTime    int64        `json:"startTime" gorm:"not null"`                             // Start time
	ExpTime      int64        `json:"expTime" gorm:"not null"`                               // Expiry time
	VoteCount    int64        `json:"voteCount" gorm:"not null,default:0"`                   // Vote count
	Status       int          `json:"status" gorm:"not null,default:0"`                      // Proposal status
	Network      int64        `json:"network" gorm:"not null"`                               // Network ID
	Name         string       `json:"name" gorm:"not null,default:''"`                       // Name
	Timezone     string       `json:"timezone" gorm:"not null"`                              // Timezone
	Descriptions string       `json:"descriptions" gorm:"not null,default:'';type:longtext"` // Descriptions
	GithubName   string       `json:"githubName" gorm:"not null,default:''"`                 // Github name
	GithubAvatar string       `json:"githubAvatar" gorm:"not null,default:''"`               // Github avatar
	GMTOffset    string       `json:"gmtOffset" gorm:"not null,default:''"`                  // GMT offset
	CurrentTime  int64        `json:"currentTime" gorm:"not null,default:0"`                 // Current time
	CreatedAt    time.Time    `json:"createdAt" gorm:"not null,autoCreateTime"`              // Created time
	UpdatedAt    time.Time    `json:"updatedAt" gorm:"not null,autoUpdateTime"`              // Updated time
	VoteResult   []VoteResult `json:"voteResult"`
}

// ContractProposal represents the structure of a proposal stored on the blockchain.
type ContractProposal struct {
	Cid          string         `json:"cid"`          // CID
	ProposalType *big.Int       `json:"proposalType"` // Proposal type
	Creator      common.Address `json:"creator"`      // Creator address
	StartTime    *big.Int       `json:"startTime"`    // Start time
	ExpTime      *big.Int       `json:"expTime"`      // Expiry time
	VotesCount   *big.Int       `json:"votesCount"`   // Votes count
}

// ProposalDetail represents the detailed information about a proposal.
type ProposalDetail struct {
	ProposalType int64    `json:"ProposalType"` // Proposal type
	VoteType     int64    `json:"VoteType"`     // Vote type
	Timezone     string   `json:"Timezone"`     // Timezone
	Time         []string `json:"Time"`         // Time
	Name         string   `json:"Name"`         // Name
	Descriptions string   `json:"Descriptions"` // Descriptions
	Option       []string `json:"option"`       // Options
	GMTOffset    []string `json:"GMTOffset"`    // GMT offset
	ShowTime     []string `json:"showTime"`     // Show time
	Address      string   `json:"Address"`      // Address
	ChainId      int64    `json:"chainId"`      // Chain ID
	CurrentTime  int64    `json:"currentTime"`  // Current time
}

type ProposalDraft struct {
	Id           int64     `json:"id"`                                         // Unique identifier
	Time         string    `json:"Time"`                                       // Time
	Option       string    `json:"Option"`                                     // Options
	Address      string    `json:"Address" binding:"required"`                 // Address
	ChainId      int64     `json:"ChainId" binding:"required"`                 // Chain ID
	ProposalId   int64     `json:"proposalId"`                                 // Proposal ID
	Cid          string    `json:"cid" gorm:"not null"`                        // CID
	ProposalType int64     `json:"proposalType" gorm:"not null,default:0"`     // Proposal type
	Creator      string    `json:"creator" gorm:"not null"`                    // Creator address
	StartTime    int64     `json:"startTime" gorm:"not null"`                  // Start time
	ExpTime      int64     `json:"expTime" gorm:"not null"`                    // Expiry time
	VoteCount    int64     `json:"voteCount" gorm:"not null"`                  // Vote count
	Status       int       `json:"status" gorm:"not null"`                     // Proposal status
	Network      int64     `json:"network" gorm:"not null"`                    // Network ID
	Name         string    `json:"name" gorm:"not null"`                       // Name
	Timezone     string    `json:"timezone" gorm:"not null"`                   // Timezone
	Descriptions string    `json:"descriptions" gorm:"not null;type:longtext"` // Descriptions
	GithubName   string    `json:"githubName" gorm:"not null,default:''"`      // Github name
	GithubAvatar string    `json:"githubAvatar" gorm:"not null,default:''"`    // Github avatar
	GMTOffset    string    `json:"gmtOffset" gorm:"not null,default:''"`       // GMT offset
	CurrentTime  int64     `json:"currentTime" gorm:"not null"`                // Current time
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`            // Created time
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime"`            // Updated time
}
