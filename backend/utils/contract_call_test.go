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

package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"powervoting-server/config"
	"powervoting-server/data"
	"powervoting-server/repo"
	"powervoting-server/service"
	"powervoting-server/utils"
)

func getSyncService() *service.SyncService {
	config.InitConfig("../")

	config.Client.ABIPath.PowerVotingAbi = "../abi/power-voting.json"
	config.Client.ABIPath.FipAbi = "../abi/power-voting-fip.json"
	config.Client.ABIPath.OraclePowersAbi = "../abi/oracle-powers.json"
	config.Client.ABIPath.OracleAbi = "../abi/oracle.json"
	config.InitLogger()

	syncSyrvuce := service.NewSyncService(
		repo.NewSyncRepo(data.NewMysql()),
		repo.NewVoteRepo(data.NewMysql()),
		repo.NewProposalRepo(data.NewMysql()),
		repo.NewFipRepo(data.NewMysql()),
		repo.NewLotusRPCRepo(),
	)

	return syncSyrvuce
}
func TestGetActorIdByAddress(t *testing.T) {
	client, err := data.GetClient(getSyncService(), 314159)
	assert.NoError(t, err)
	id, err := utils.GetActorIdByAddress(client, "0x8cDc8c7a027f18503f4A7C24e4b7488B08A56223")
	assert.NoError(t, err)
	assert.Equal(t, 17855, id)
}

func TestGetOwnerIdByOracle(t *testing.T) {
	client, err := data.GetClient(getSyncService(), 314159)
	assert.NoError(t, err)
	id := utils.GetOwnerIdByOracle(client, 17829, []uint64{144416})
	assert.Equal(t, "5", id)
}
