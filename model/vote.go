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

type Vote struct {
	Id         int64  `json:"id"`
	ProposalId int64  `gorm:"constraint:off" gorm:"not null"`
	Address    string `json:"address" gorm:"not null"`
	VoteInfo   string `json:"voteInfo" gorm:"not null"`
	Network    int64  `json:"network" gorm:"not null"`
}

type VoteResult struct {
	Id         int64   `json:"id"`
	ProposalId int64   `gorm:"constraint:off;not null"`
	OptionId   int64   `json:"optionId" gorm:"not null"`
	Votes      float64 `json:"votes" gorm:"not null"`
	Network    int64   `json:"network" gorm:"not null"`
}

type VoteHistory struct {
	Id         int64   `json:"id"`
	ProposalId int64   `gorm:"constraint:off" gorm:"not null"`
	Address    string  `json:"address" gorm:"not null"`
	OptionId   int64   `json:"optionId" gorm:"not null"`
	Votes      float64 `json:"votes" gorm:"not null"`
	Network    int64   `json:"network" gorm:"not null"`
}

type ContractVote struct {
	Voter    common.Address `json:"voter"`
	VoteInfo string         `json:"voteInfo"`
}

type Vote4Counting struct {
	OptionId int64
	Votes    int64
	Address  string
}
