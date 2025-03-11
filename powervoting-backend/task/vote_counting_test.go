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
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"powervoting-server/config"
	"powervoting-server/data"
	"powervoting-server/model"
	"powervoting-server/repo"
	"powervoting-server/service"
)

func initConfig() {
	config.InitConfig("../")
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	config.Client.ABIPath.PowerVotingAbi = "../abi/power-voting.json"
	config.Client.ABIPath.OracleAbi = "../abi/oracle.json"
}

func newVoteCount(t *testing.T) *VoteCount {
	initConfig()
	mydb := data.NewMysql()
	syncService := service.NewSyncService(
		repo.NewSyncRepo(mydb),
		repo.NewVoteRepo(mydb),
		repo.NewProposalRepo(mydb),
	)
	client, err := data.GetClient(syncService, 314159)
	assert.Nil(t, err)

	vc := &VoteCount{
		EthClient:   client,
		SyncService: syncService,
	}
	return vc
}

func TestCalculateVoteWeight(t *testing.T) {
	vc := newVoteCount(t)
	totalPower, voltList := vc.calculateVoteWeight(
		1,
		mockPowersMap(),
		mockVote(),
		model.Percentage{
			SpPercentage:          25,
			ClientPercentage:      25,
			TokenHolderPercentage: 40,
			DeveloperPercentage:   10,
		},
		314159,
	)

	assert.NotNil(t, totalPower)

	assert.Equal(t, 2, len(voltList))
}
func mockProposal() model.ProposalTbl {
	return model.ProposalTbl{
		ProposalId: 999999,
		ChainId:    314159,
		Percentage: model.Percentage{
			SpPercentage:          25,
			ClientPercentage:      25,
			TokenHolderPercentage: 40,
			DeveloperPercentage:   10,
		},
	}
}

func mockVote() []model.VoteTbl {
	return []model.VoteTbl{
		{
			Address:          "vote1",
			ProposalId:       1,
			ChainId:          314159,
			ClientPower:      "1000",
			TokenHolderPower: "1000",
			DeveloperPower:   "2000",
			SpPower:          "1000",
			VoteEncrypted:    "-----BEGIN AGE ENCRYPTED FILE-----YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IHRsb2NrIDE1OTIwMjEyIDUyZGI5YmE3MGUwY2MwZjZlYWY3ODAzZGQwNzQ0N2ExZjU0Nzc3MzVmZDNmNjYxNzkyYmE5NDYwMGM4NGU5NzEKa0JUbU00U3VWcndSMUVHUkF1MHJ0VWlyNGRrMnFVVWFCQnVtM0Q0cnNvQUN5SUJMdmNzV0cvY3ByT21yekxSbApESHRNanVzblZwWFptMU5JS2YyUnlXYkRYZ01FeVQwSjNqLzlLeXJpSSt5UkNOUFp1ak5NNmJJbDdqbHorUmkrCkV6RVBKcU12NVRUckx3bHRxOUNCRklqV01KaUp5bWJndksyRERNMlN1ZWMKLS0tIFVRWGZhQmozY3hqMHIxR3pVSDlram9MRWgvR2J0RVhZbGM3MTU2clVlaGsKxfAndo7NGBy8V5vTkXikn1BRSOQ+I0fW0RJgwYCAsGFc7Dqnkp9/qOkexzKY-----END AGE ENCRYPTED FILE-----",
		},
		{
			Address:          "vote2",
			ProposalId:       1,
			ChainId:          314159,
			ClientPower:      "1000",
			TokenHolderPower: "1000",
			DeveloperPower:   "2000",
			SpPower:          "2000",
			VoteEncrypted:    "-----BEGIN AGE ENCRYPTED FILE-----YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IHRsb2NrIDE1OTIwMjEyIDUyZGI5YmE3MGUwY2MwZjZlYWY3ODAzZGQwNzQ0N2ExZjU0Nzc3MzVmZDNmNjYxNzkyYmE5NDYwMGM4NGU5NzEKaE9VdTZ3Y0U1UGY4L3NuMzhpNjJ0MjUvcmVHeDV6Q0pINzRwYWp2SE1RRVdQTkx1aHpkb3d3amV4QlZ4M1k5RwpCYWtEc2NXenNoQis2TzJXNlhFS0VjN0pyMTRwTUZzUm81cy9jVjhuWnFxZWdHMDZkOG93b1RkaE8xOFJiY2VZCnlWWHl6R3VaazlnRXJ1VEY0YStmblVvTW9rVG5IWnU1MlBYTEZ6ZHFhdnMKLS0tIHVHbFVHdDNmREZKOG0wNUFMN1AvdGROY1NIeThiOGJzcTFGWVRoeVF1SzQKq1Cogk65rCvDECoKIFq8Z6bAvOtZWQ/fMvWstqXcaApSM42vnM8iBprD+d8=-----END AGE ENCRYPTED FILE-----",
		},
	}
}

func mockPowersMap() map[string]model.AddrPower {
	return map[string]model.AddrPower{
		"vote1": {
			Address:          "vote1",
			DeveloperPower:   big.NewInt(2000),
			TokenHolderPower: big.NewInt(1000),
			ClientPower:      big.NewInt(1000),
			SpPower:          big.NewInt(1000),
		},
		"vote2": {
			Address:          "vote2",
			DeveloperPower:   big.NewInt(2000),
			TokenHolderPower: big.NewInt(1000),
			ClientPower:      big.NewInt(1000),
			SpPower:          big.NewInt(2000),
		},
	}
}
