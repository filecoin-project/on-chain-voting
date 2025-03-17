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

	"github.com/stretchr/testify/assert"

	"powervoting-server/data"
	"powervoting-server/mock"
	"powervoting-server/model"
	"powervoting-server/utils"
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
func TestCalculateVoteWeight(t *testing.T) {
	vc := newVoteCount(t)
	votesInfo, err := vc.SyncService.GetUncountedVotedList(context.Background(), 314159, 4)
	assert.Nil(t, err)

	powerMap := utils.PowersInfoToMap(mockPower())
	totalPower, voltList := vc.calculateVoteWeight(
		4,
		// mockPowersMap(),
		powerMap,
		// mockVote(),
		votesInfo,
		model.Percentage{
			SpPercentage:          2500,
			ClientPercentage:      2500,
			TokenHolderPercentage: 2500,
			DeveloperPercentage:   2500,
		},
		314159,
	)

	assert.NotNil(t, totalPower)
	resultPercent := vc.calculateFinalPercentages(totalPower.approvePercentage, totalPower.rejectPercentage, len(voltList))
	fmt.Printf("resultPercent: %v\n", resultPercent)
	assert.Equal(t, 2, len(voltList))
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
