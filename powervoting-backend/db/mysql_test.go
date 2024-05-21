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

package db

import (
	"fmt"
	"powervoting-server/config"
	"powervoting-server/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestVoteResult(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	config.InitConfig("../")
	InitMysql()

	var voteHistory model.VoteCompleteHistory
	var voteResultList []model.VoteResult
	var votePower []model.VotePower

	votePower = append(votePower, model.VotePower{
		HistoryId:               66, // VoteCompleteHistory Indicates the id
		OptionId:                0,
		Votes:                   2.5,
		Address:                 "0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307",
		SpPower:                 "",
		ClientPower:             "",
		TokenHolderPower:        "",
		DeveloperPower:          "",
		SpPowerPercent:          0,
		ClientPowerPercent:      0,
		TokenHolderPowerPercent: 0.25,
		DeveloperPowerPercent:   0,
		PowerBlockHeight:        3212312,
	})
	votePower = append(votePower, model.VotePower{
		HistoryId:               66, // VoteCompleteHistory Indicates the id
		OptionId:                1,
		Votes:                   2.5,
		Address:                 "0xf58cC34cf80BDF9D3aD82E7AC57aCd02cA592193",
		SpPower:                 "",
		ClientPower:             "",
		TokenHolderPower:        "18446744073709551616",
		DeveloperPower:          "",
		SpPowerPercent:          0,
		ClientPowerPercent:      0,
		TokenHolderPowerPercent: 0.25,
		DeveloperPowerPercent:   0,
		PowerBlockHeight:        3212312,
	})
	voteHistory = model.VoteCompleteHistory{
		ProposalId:            65,
		Network:               314159,
		TotalSpPower:          "0",
		TotalClientPower:      "0",
		TotalTokenHolderPower: "18446744073709551616",
		TotalDeveloperPower:   "0",
		VotePowers:            votePower,
	}
	voteResultList = append(voteResultList, model.VoteResult{ProposalId: 78, OptionId: 0, Votes: 1000, Network: 314159})
	voteResultList = append(voteResultList, model.VoteResult{ProposalId: 78, OptionId: 1, Votes: 2000, Network: 314159})
	voteResultList = append(voteResultList, model.VoteResult{ProposalId: 78, OptionId: 2, Votes: 3000, Network: 314159})
	id, err := Engine.VoteResult(1, voteHistory, voteResultList)
	assert.Nil(t, err)

	// assert voteHistory
	var count int64
	err = Engine.Model(&voteHistory).
		Where("id = ?", id).
		Where("proposal_id = ?", 65).
		Where("network = ?", 314159).
		Count(&count).Error
	assert.Nil(t, err)

	assert.Equal(t, int64(1), count)
}

func TestDatabaseInit(t *testing.T) {
	config.InitConfig("../")
	InitMysql()
	var voteHistory model.VoteCompleteHistory
	var votePower []model.VotePower
	Engine.AutoMigrate(&voteHistory, &votePower)
}

func TestVoteResultQuery(t *testing.T) {
	config.InitConfig("../")
	InitMysql()

	expectedHistory := model.VoteCompleteHistory{
		Id:                    1,
		ProposalId:            65,
		Network:               314159,
		TotalSpPower:          "0",
		TotalClientPower:      "0",
		TotalTokenHolderPower: "18446744073709551616",
		TotalDeveloperPower:   "0",
		VotePowers: []model.VotePower{
			{
				Id:                      1,
				HistoryId:               1,
				Address:                 "0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307",
				OptionId:                0,
				Votes:                   2.5,
				SpPower:                 "",
				ClientPower:             "",
				TokenHolderPower:        "",
				DeveloperPower:          "",
				SpPowerPercent:          0,
				ClientPowerPercent:      0,
				TokenHolderPowerPercent: 0.25,
				DeveloperPowerPercent:   0,
				PowerBlockHeight:        3212312,
			},
			{
				Id:                      2,
				HistoryId:               1,
				Address:                 "0xf58cC34cf80BDF9D3aD82E7AC57aCd02cA592193",
				OptionId:                1,
				Votes:                   2.5,
				SpPower:                 "",
				ClientPower:             "",
				TokenHolderPower:        "18446744073709551616",
				DeveloperPower:          "",
				SpPowerPercent:          0,
				ClientPowerPercent:      0,
				TokenHolderPowerPercent: 0.25,
				DeveloperPowerPercent:   0,
				PowerBlockHeight:        3212312,
			},
		},
	}

	proposalId := 65
	network := 314159
	var history model.VoteCompleteHistory
	tx := Engine.Preload("VotePowers").Where("proposal_id", proposalId).Where("network", network).Find(&history)
	assert.Nil(t, tx.Error)

	assert.Equal(t, history, expectedHistory)
	fmt.Printf("%+v\n", history)
}

func TestGetProposalList(t *testing.T) {
	config.InitConfig("../")
	InitMysql()

	expectedProposalList := []model.Proposal{
		{
			Id:           1,
			ProposalId:   1,
			Cid:          "1",
			ProposalType: 1,
			Creator:      "1",
			StartTime:    1,
			ExpTime:      2,
			VoteCount:    1,
			Status:       0,
			Network:      314159,
		},
	}

	res, err := Engine.GetProposalList(314159, 2)
	assert.Nil(t, err)

	assert.Equal(t, res, expectedProposalList)

	fmt.Println(res)
}

func TestGetVoteList(t *testing.T) {
	config.InitConfig("../")
	InitMysql()

	expectedVoteList := []model.Vote{
		{
			Id:         1,
			ProposalId: 1,
			Address:    "test_addr",
			VoteInfo:   "test_info",
			Network:    314159,
		},
	}

	res, err := Engine.GetVoteList(314159, 1)
	assert.Nil(t, err)

	assert.Equal(t, res, expectedVoteList)
	fmt.Println(res)
}
