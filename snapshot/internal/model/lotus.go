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

// LotusMinerPower is the response of lotus-miner power
type LotusMinerPower struct {
	MinerPower  MinerPower `json:"MinerPower"`
	TotalPower  TotalPower `json:"TotalPower"`
	HasMinPower bool       `json:"HasMinPower"`
}

type MinerPower struct {
	RawBytePower    string `json:"RawBytePower"`
	QualityAdjPower string `json:"QualityAdjPower"`
}

type TotalPower struct {
	RawBytePower    string `json:"RawBytePower"`
	QualityAdjPower string `json:"QualityAdjPower"`
}

type BalanceInfo struct {
	Height  string `json:"Height"`
	Balance string `json:"Balance"`
}

type PieceCID struct {
	CID string `json:"/"`
}

type Proposal struct {
	Client               string   `json:"Client"`
	ClientCollateral     string   `json:"ClientCollateral"`
	EndEpoch             int64    `json:"EndEpoch"`
	Label                string   `json:"Label"`
	PieceCID             PieceCID `json:"PieceCID"`
	PieceSize            int64    `json:"PieceSize"`
	Provider             string   `json:"Provider"`
	ProviderCollateral   string   `json:"ProviderCollateral"`
	StartEpoch           int64    `json:"StartEpoch"`
	StoragePricePerEpoch string   `json:"StoragePricePerEpoch"`
	VerifiedDeal         bool     `json:"VerifiedDeal"`
}

type State struct {
	LastUpdatedEpoch int64 `json:"LastUpdatedEpoch"`
	SectorNumber     int64 `json:"SectorNumber"`
	SectorStartEpoch int64 `json:"SectorStartEpoch"`
	SlashEpoch       int64 `json:"SlashEpoch"`
}


type BlockHeader struct {
	Height    int64 `json:"Height"`
	Timestamp int64 `json:"Timestamp"`
}

type GetBlockParam struct {
	Empty string `json:"/"`
}
