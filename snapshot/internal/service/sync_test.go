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

package service_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	models "power-snapshot/internal/model"
)

var _ = Describe("Sync", func() {
	Describe("GetAllAddrInfoList", func() {
		It("should return the all address info list", func() {
			res, err := syncService.GetAllAddrInfoList(context.Background(), conf.Network.ChainId)
			Expect(err).To(BeNil())
			Expect(res).To(Equal([]models.AddrInfo{
				{
					Addr:          "f01",
					ActionIDs:     []string{"f01"},
					GithubAccount: "test",
					MinerIDs:      []string{"f01"},
				},
			}))
		})
	})

	Describe("UploadSnapshotInfoByDay", func() {
		It("should create snapshot info to db", func() {
			res, err := syncService.UploadSnapshotInfoByDay(context.Background(), snapshotData, "20250101", conf.Network.ChainId)
			Expect(err).To(BeNil())
			Expect(res).To(Equal(int64(1)))

		})
	})
})

// func TestFetchDeals(t *testing.T) {
// 	config.InitConfig("../../")
// 	config.InitLogger()
// 	manager, err := data.NewGoEthClientManager()
// 	assert.NoError(t, err)
// 	redisClient, err := data.NewRedisClient(60)
// 	assert.NoError(t, err)
// 	baseRepo := repo.NewBaseRepoImpl(manager, redisClient)
// 	lotusRepo := repo.NewLotusRPCRepo(redisClient)
// 	syncService := service.NewSyncService(baseRepo, nil, nil, lotusRepo, nil, nil)
// 	err = syncService.FetchDeals(context.Background(), config.Client.Network.ChainId)
// 	assert.NoError(t, err)
// }

// func TestSubTaskWorker(t *testing.T) {
// 	config.InitConfig("../../")
// 	config.InitLogger()
// 	manager, err := data.NewGoEthClientManager()
// 	assert.NoError(t, err)
// 	redisClient, err := data.NewRedisClient(60)
// 	assert.NoError(t, err)
// 	baseRepo := repo.NewBaseRepoImpl(manager, redisClient)
// 	lotusRepo := repo.NewLotusRPCRepo(redisClient)

// 	jetstreamClient, err := data.NewJetstreamClient()
// 	assert.NoError(t, err)
// 	defer func() {
// 		err := jetstreamClient.Drain()
// 		assert.NoError(t, err)
// 	}()
// 	syncRepo, err := repo.NewSyncRepoImpl(config.Client.Network.ChainId, redisClient, jetstreamClient)
// 	assert.NoError(t, err)
// 	syncService := service.NewSyncService(baseRepo, syncRepo, nil, lotusRepo, nil, nil)
// 	res, err := syncService.SubTaskWorker(context.Background(), config.Client.Network.ChainId, models.Task{
// 		UID:           "test",
// 		Address:       "0xfF000000000000000000000000000000000278bc",
// 		GithubAccount: "liuzeming1",
// 		SubTasks: []models.SubTask{
// 			{
// 				UID:         "test",
// 				Address:     "0xfF000000000000000000000000000000000278bc",
// 				Typ:         constant.TaskActionActor,
// 				ShortID:     "t0161980",
// 				DateStr:     "20250519",
// 				BlockHeight: 2677894,
// 			},
// 		},
// 	})
// 	assert.NoError(t, err)
// 	assert.NotNil(t, res)
// 	zap.L().Info("",zap.Any("res", res))
// }
