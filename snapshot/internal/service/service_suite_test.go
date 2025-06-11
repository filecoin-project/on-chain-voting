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
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"power-snapshot/constant"
	models "power-snapshot/internal/model"
	"power-snapshot/internal/service"
	mocks "power-snapshot/mock"
	"power-snapshot/utils/types"
)

var conf models.Config
var tipset []any
var (
	queryService    *service.QueryService
	syncService     *service.SyncService
	mockLotus       *mocks.MockLotusRepo
	mockSync        *mocks.MockSyncRepo
	mockBase        *mocks.MockBaseRepo
	mockQuery       *mocks.MockQueryRepo
	mockMysql       *mocks.MockMysqlRepo
	mockContract    *mocks.MockContractRepo
	mockCtrl        *gomock.Controller
	mockBackendGrpc *mocks.MockIBackendGRPC
	mockGithubLimit *mocks.MockGithubLimit
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var _ = BeforeSuite(func() {
	conf = models.Config{
		Network: models.Network{
			ChainId: 314159,
			Name:    "TestNet",
		},
		Github: models.GitHub{
			Token: []string{"testToken"},
		},
	}
	tipset = []any{
		map[string]string{"/": "baby1"},
		map[string]string{"/": "baby2"},
	}

	mockCtrl = gomock.NewController(GinkgoT())
	mockLotus = mocks.NewMockLotusRepo(mockCtrl)
	mockBase = mocks.NewMockBaseRepo(mockCtrl)
	mockSync = mocks.NewMockSyncRepo(mockCtrl)
	mockQuery = mocks.NewMockQueryRepo(mockCtrl)
	mockMysql = mocks.NewMockMysqlRepo(mockCtrl)
	mockContract = mocks.NewMockContractRepo(mockCtrl)
	mockBackendGrpc = mocks.NewMockIBackendGRPC(mockCtrl)
	mockGithubLimit = mocks.NewMockGithubLimit(mockCtrl)
	syncService = service.NewSyncService(mockBase, mockSync, mockMysql, mockLotus, mockContract, mockBackendGrpc)
	queryService = service.NewQueryService(mockBase, mockQuery, syncService, mockLotus)
	mockBaseFunc()
	mockQueryFunc()
	mockLotusFunc()

	mockSyncFunc()
	mockBackendGrpcFunc()
	mockGithubLimitFunc()
	mockMysqlFunc()
	mockContractFunc()
})

var syncPower = models.SyncPower{
	Address:          "f01",
	DateStr:          "20250101",
	GithubAccount:    "test",
	DeveloperPower:   big.NewInt(1),
	SpPower:          big.NewInt(100),
	ClientPower:      big.NewInt(100),
	TokenHolderPower: big.NewInt(10_000_000),
	BlockHeight:      1,
}

var snapshotData = map[string]any{
	"day": "20250101",
	"info": []models.AddrInfo{
		{
			Addr:          "f01",
			ActionIDs:     []string{"f01"},
			GithubAccount: "test",
			MinerIDs:      []string{"f01"},
		},
	},
}

func mockQueryFunc() {
	mockQuery.EXPECT().
		GetAddressPower(gomock.Any(), gomock.Eq(conf.Network.ChainId), gomock.Any(), gomock.Any()).
		Return(&syncPower, nil).AnyTimes()
	mockQuery.EXPECT().
		GetDeveloperWeights(gomock.Any(), gomock.Any()).
		Return(map[string]int64{"test": 1}, nil).AnyTimes()
	mockQuery.EXPECT().
		GetAddressPowerByDay(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]models.SyncPower{syncPower}, nil).AnyTimes()
	mockQuery.EXPECT().
		GetDevPowerByDay(gomock.Any(), gomock.Any()).
		Return("testDayPower", nil).AnyTimes()
}

func mockBaseFunc() {
	mockBase.EXPECT().
		GetDateHeightMap(gomock.Any(), gomock.Eq(conf.Network.ChainId)).
		Return(map[string]int64{"20250101": 1}, nil).AnyTimes()
	mockBase.EXPECT().
		SetDateHeightMap(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).AnyTimes()
	mockBase.EXPECT().
		GetDealsFromLocal(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(types.StateMarketDeals{}, nil).AnyTimes()
	mockBase.EXPECT().
		GetDeveloperWeights(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]models.Nodes{}, nil).AnyTimes()
}

func mockSyncFunc() {
	mockSync.EXPECT().
		GetAddrSyncedDate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]string{"20250101"}, nil).AnyTimes()
	mockSync.EXPECT().
		AddTask(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).AnyTimes()
}
func mockLotusFunc() {
	t, _ := time.Parse("20060102", "20250501")
	mockLotus.EXPECT().
		GetNewestHeight(gomock.Any(), gomock.Any()).
		Return(int64(500000), nil).AnyTimes()
	mockLotus.EXPECT().
		GetBlockHeader(gomock.Any(), gomock.Eq(conf.Network.ChainId), gomock.Any()).
		Return(models.BlockHeader{Height: 500000, Timestamp: t.Unix()}, nil).AnyTimes()
	mockLotus.EXPECT().
		GetWalletBalanceByHeight(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("100", nil).AnyTimes()
	mockLotus.EXPECT().
		GetTipSetByHeight(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(tipset, nil).AnyTimes()
	mockLotus.EXPECT().
		GetMinerPowerByHeight(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(models.LotusMinerPower{MinerPower: models.MinerPower{
			RawBytePower:    "100",
			QualityAdjPower: "100",
		}, TotalPower: models.TotalPower{
			RawBytePower:    "100",
			QualityAdjPower: "100",
		}}, nil).AnyTimes()
}
func mockBackendGrpcFunc() {

	mockBackendGrpc.EXPECT().GetAllVoterAddresss(gomock.Eq(conf.Network.ChainId)).
		Return([]string{"f01"}, nil).AnyTimes()
	mockBackendGrpc.EXPECT().GetVoterInfo(gomock.Any()).Return(models.VoterInfo{
		ActorIds:      []string{"f01"},
		MinerIds:      []string{"f01"},
		GithubAccount: "test",
		EthAddress:    common.Address{},
	}, nil).AnyTimes()
}

func mockGithubLimitFunc() {
	mockGithubLimit.EXPECT().
		CheckRateLimitBeforeRequest(gomock.Eq("token1")).
		Return(int32(5000), int32(5000)).AnyTimes()
}

func mockMysqlFunc() {
	mockMysql.EXPECT().
		CreateSnapshotBackup(gomock.Any(), gomock.Any()).
		Return(nil).AnyTimes()
}

func mockContractFunc() {
	mockContract.EXPECT().
		GetGithubRepoInfo(gomock.Any()).
		Return([]string{"test"}, nil).AnyTimes()
	mockContract.EXPECT().
		GetExpirationData().Return(constant.DataExpiredDuration, nil).AnyTimes()
}
