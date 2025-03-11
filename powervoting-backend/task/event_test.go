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
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"

	"powervoting-server/config"
	"powervoting-server/constant"
	"powervoting-server/data"
	"powervoting-server/model"
	"powervoting-server/repo"
	"powervoting-server/service"
)

func TestSubscriptEvent(t *testing.T) {
	initConfig()
	mydb := data.NewMysql()
	for _, conf := range config.Client.Network {
		syncService := service.NewSyncService(
			repo.NewSyncRepo(mydb),
			repo.NewVoteRepo(mydb),
			repo.NewProposalRepo(mydb),
		)
		ethClient, err := data.GetClient(syncService, conf.ChainId)
		assert.NoError(t, err)
		ev := &Event{
			Client:      ethClient,
			SyncService: syncService,
			Network:     &conf,
		}

		ev.SubscribeEvent()
	}
}

func TestFetchMatchingEventLogs(t *testing.T) {
	initConfig()

	mydb := data.NewMysql()
	syncService := service.NewSyncService(
		repo.NewSyncRepo(mydb),
		repo.NewVoteRepo(mydb),
		repo.NewProposalRepo(mydb),
	)
	client, err := data.GetClient(syncService, 314159)
	assert.NoError(t, err)

	logs, err := mockLogs(client)
	assert.NoError(t, err)
	assert.NotNil(t, logs)
}

func TestParseEvent(t *testing.T) {
	initConfig()

	mydb := data.NewMysql()
	syncService := service.NewSyncService(
		repo.NewSyncRepo(mydb),
		repo.NewVoteRepo(mydb),
		repo.NewProposalRepo(mydb),
	)

	client, err := data.GetClient(syncService, 314159)
	assert.NoError(t, err)
	logs, err := mockLogs(client)
	assert.NoError(t, err)

	ev := &Event{
		Client:      client,
		SyncService: syncService,
		Network:     nil,
	}
	for _, log := range logs {
		err = ev.parseEvent(context.Background(), log)
		assert.NoError(t, err)
	}

}

func mockLogs(client *model.GoEthClient) ([]types.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(2435250),
		ToBlock:   big.NewInt(2437250),
		Addresses: []common.Address{common.HexToAddress("0xB1bD6785540F4e704b41912d65BB828d1008baa3")},
		Topics: [][]common.Hash{
			{
				client.ABI.PowerVotingAbi.Events[constant.ProposalEvt].ID,
				client.ABI.PowerVotingAbi.Events[constant.VoteEvt].ID,
			},
		},
	}

	logs, err := client.Client.FilterLogs(context.Background(), query)

	return logs, err
}
