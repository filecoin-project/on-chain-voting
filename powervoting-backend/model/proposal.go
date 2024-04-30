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
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Proposal struct {
	Id           int64  `json:"id"`
	ProposalId   int64  `json:"proposalId"`
	Cid          string `json:"cid" gorm:"not null"`
	ProposalType int64  `json:"proposalType" gorm:"not null"`
	Creator      string `json:"creator" gorm:"not null"`
	ExpTime      int64  `json:"expTime" gorm:"not null"`
	VoteCount    int64  `json:"voteCount" gorm:"not null"`
	Status       int    `json:"status" gorm:"not null"`
	Network      int64  `json:"network" gorm:"not null"`
}

type ContractProposal struct {
	Cid          string         `json:"cid"`
	ProposalType *big.Int       `json:"proposalType"`
	Creator      common.Address `json:"creator"`
	ExpTime      *big.Int       `json:"expTime"`
	VotesCount   *big.Int       `json:"votesCount"`
}

type ProposalDetail struct {
	ProposalType int64    `json:"ProposalType"`
	VoteType     int64    `json:"VoteType"`
	Timezone     string   `json:"Timezone"`
	Time         int64    `json:"Time"`
	Name         string   `json:"Name"`
	Descriptions string   `json:"Descriptions"`
	Option       []string `json:"option"`
	GMTOffset    []string `json:"GMTOffset"`
	ShowTime     string   `json:"showTime"`
	Address      string   `json:"Address"`
	ChainId      int64    `json:"chainId"`
	CurrentTime  int64    `json:"currentTime"`
}
