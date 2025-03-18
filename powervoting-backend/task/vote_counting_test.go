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

package task

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"powervoting-server/config"
	"powervoting-server/constant"
	"powervoting-server/data"
	"powervoting-server/mock"
	"powervoting-server/model"
)

func newVoteCount(t *testing.T) *VoteCount {
	mock.InifMockConfig()

	syncService := &mock.MockSyncService{}
	client, err := data.GetClient(syncService, 314159)
	assert.Nil(t, err)

	vc := &VoteCount{
		EthClient:   client,
		SyncService: syncService,
	}
	return vc
}

func TestVotingCountHandler(t *testing.T) {
	vc := newVoteCount(t)

	proposals, err := vc.SyncService.UncountedProposalList(context.Background(), 314159, 4)
	assert.Nil(t, err)
	assert.Len(t, proposals, 1)

	assert.NoError(t, err)

	vc.processCounting(vc.EthClient, proposals[0], model.SnapshotAllPower{
		AddrPower: mockPower(),
	})
}
func TestCountWeightCredits(t *testing.T) {
	vc := newVoteCount(t)

	votes, err := vc.SyncService.GetUncountedVotedList(context.Background(), 314159, 1)
	assert.NoError(t, err)
	votePower, totalPower, votesList := vc.countWeightCredits(
		1,
		map[string]model.AddrPower{
			"0x1234567890123456789012345678901234567890": {
				Address:          "0x1234567890123456789012345678901234567890",
				DeveloperPower:   big.NewInt(0),
				ClientPower:      big.NewInt(0),
				SpPower:          big.NewInt(0),
				TokenHolderPower: big.NewInt(1000),
			},

			"0x1234567890123456789012345678901234567891": {
				Address:          "0x1234567890123456789012345678901234567891",
				DeveloperPower:   big.NewInt(0),
				ClientPower:      big.NewInt(0),
				SpPower:          big.NewInt(0),
				TokenHolderPower: big.NewInt(9000),
			},
		},
		votes,
		model.Percentage{
			SpPercentage:          uint16(2500),
			DeveloperPercentage:   uint16(2500),
			ClientPercentage:      uint16(2500),
			TokenHolderPercentage: uint16(2500),
		},
		314159,
	)

	assert.Len(t, votePower, 2)
	assert.Equal(t, decimal.NewFromInt(1000), votePower[constant.VoteApprove].TokenPower)
	assert.Equal(t, decimal.NewFromInt(9000), votePower[constant.VoteReject].TokenPower)
	assert.Equal(t, decimal.NewFromInt(10000), totalPower.TokenPower)
	assert.NotNil(t, votePower)
	assert.NotNil(t, totalPower)
	assert.NotNil(t, votesList)
}

func TestCalculateVotesPercentage(t *testing.T) {
	vc := newVoteCount(t)
	config.InitLogger()
	votesPower := map[string]model.VoterPowerCount{
		constant.VoteApprove: {
			SpPower:        decimal.NewFromInt(1000),
			DeveloperPower: decimal.NewFromInt(2000),
			ClientPower:    decimal.NewFromInt(1000),
			TokenPower:     decimal.NewFromInt(1000),
		},
		constant.VoteReject: {
			SpPower:        decimal.NewFromInt(9000),
			DeveloperPower: decimal.NewFromInt(8000),
			ClientPower:    decimal.NewFromInt(1000),
			TokenPower:     decimal.NewFromInt(1000),
		},
	}

	totalPower := model.VoterPowerCount{
		SpPower:        decimal.NewFromInt(10000),
		DeveloperPower: decimal.NewFromInt(10000),
		ClientPower:    decimal.NewFromInt(2000),
		TokenPower:     decimal.NewFromInt(2000),
	}

	// Calculation formula:
	// (SpPower / totalPower) * SpPercentage +
	// (DeveloperPower / totalPower) * DeveloperPercentage +
	// (ClientPower / totalPower) * ClientPercentage +
	// (TokenPower / totalPower) * TokenHolderPercentage
	// = weight
	//
	// percentage = 0
	// if spPower != 0 -> percentage += SpPercentage
	// if developerPower != 0 -> percentage += DeveloperPercentage
	// if clientPower != 0 -> percentage += ClientPercentage
	// if tokenPower != 0 -> percentage += TokenHolderPercentage
	//
	// res := weight / percentage * 100%
	res := vc.calculateVotesPercentage(votesPower[constant.VoteApprove], totalPower, model.Percentage{
		SpPercentage:          uint16(2500),
		DeveloperPercentage:   uint16(2500),
		ClientPercentage:      uint16(2500),
		TokenHolderPercentage: uint16(2500),
	})
	fmt.Printf("res: %v\n", res.String())
	assert.NotNil(t, res)
}
func mockPower() []model.AddrPower {
	return []model.AddrPower{
		{
			Address:          "0x1234567890123456789012345678901234567890",
			DeveloperPower:   big.NewInt(0),
			ClientPower:      big.NewInt(0),
			SpPower:          big.NewInt(0),
			TokenHolderPower: big.NewInt(1000),
		},
		{
			Address:          "0x1234567890123456789012345678901234567891",
			DeveloperPower:   big.NewInt(0),
			ClientPower:      big.NewInt(0),
			SpPower:          big.NewInt(0),
			TokenHolderPower: big.NewInt(9000),
		},
	}
}
