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
	"math/big"
	"time"
)

// Power represents the power information.
type Power struct {
	DeveloperPower   *big.Int `json:"developerPower"`   // Developer power
	SpPower          *big.Int `json:"spPower"`          // SP power
	ClientPower      *big.Int `json:"clientPower"`      // Client power
	TokenHolderPower *big.Int `json:"tokenHolderPower"` // Token holder power
	BlockHeight      *big.Int `json:"blockHeight"`      // Block height
}

// ContractPower represents the contract power information.
type ContractPower struct {
	DeveloperPower   *big.Int `json:"developerPower"`   // Developer power
	SpPower          [][]byte `json:"spPower"`          // SP power
	ClientPower      [][]byte `json:"clientPower"`      // Client power
	TokenHolderPower *big.Int `json:"tokenHolderPower"` // Token holder power
	BlockHeight      *big.Int `json:"blockHeight"`      // Block height
}

// VoterToPowerStatus represents the voter's power status.
type VoterToPowerStatus struct {
	DayId        *big.Int `json:"dayId"`        // Day ID
	HasFullRound *big.Int `json:"hasFullRound"` // Has full round
}

type SnapshotByDay struct {
	Id        int64     `json:"id"`
	Day       string    `gorm:"not null" json:"day"`
	PowerInfo string    `gorm:"type:text" json:"powerInfo"` //The power set of all addresses for the day
	NetId     int64     `gorm:"not null" json:"netId"`      // Block id
	Cid       string    `gorm:"" json:"cid"`                //File ID uploaded to W3Storage
	Height    int64     `gorm:"not null" json:"height"`     // Block height
	CreatedAt time.Time `gorm:"not null,autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null,autoUpdateTime" json:"updatedAt"`
}

type AllPowerByDay struct {
	Id        int64  `json:"id"`
	Day       string `json:"day"`
	PowerInfo string `json:"powerInfo"`
}
