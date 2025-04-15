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

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"testing"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"

	"power-snapshot/api"
	"power-snapshot/config"
	"power-snapshot/constant"
	"power-snapshot/internal/data"
	models "power-snapshot/internal/model"
	"power-snapshot/internal/repo"
	"power-snapshot/utils"
	"power-snapshot/utils/types"
)

var (
	testNetID int64 = 314159
)

type mockBaseRepo struct {
	AddrSyncedDateMap map[string][]string
	DhMap             map[string]int64
}

func (m *mockBaseRepo) GetLotusClientByHashKey(ctx context.Context, netID int64, key string) (jsonrpc.RPCClient, error) {
	manager, err := data.NewGoEthClientManager(config.Client.Network)
	if err != nil {
		return nil, err
	}

	client := manager.GetClient()

	h := fnv.New32a()
	key = fmt.Sprintf("%s_%d", key, time.Now().UnixNano())
	_, err = h.Write([]byte(key))
	if err != nil {
		return nil, err
	}

	index := h.Sum32() % uint32(len(client.QueryRpc))

	return jsonrpc.NewClient(client.QueryRpc[index]), nil
}

func (s *mockBaseRepo) SetDeveloperWeights(ctx context.Context, dayStr string, commits []models.Nodes) error {
	return nil
}

func (s *mockBaseRepo) GetDeveloperWeights(ctx context.Context, dayStr string) ([]models.Nodes, error) {
	return []models.Nodes{}, nil
}

func (m *mockBaseRepo) GetEthClient(ctx context.Context, netID int64) (*models.GoEthClient, error) {
	manager, err := data.NewGoEthClientManager(config.Client.Network)
	if err != nil {
		return nil, err
	}
	client := manager.GetClient()

	return &client, nil
}

func (m *mockBaseRepo) GetLotusClient(ctx context.Context, netId int64) (jsonrpc.RPCClient, error) {
	return nil, nil

}
func (m *mockBaseRepo) CreateSnapshot(ctx context.Context, netId int64, cid, day string) error {
	panic("implement me")
}
func (m *mockBaseRepo) SetDateHeightMap(ctx context.Context, netId int64, height map[string]int64) error {
	return nil
}

func (m *mockBaseRepo) GetVoteInfo(ctx context.Context, netID int64, addr string) (*models.VoterInfo, error) {
	config.InitLogger()
	err := config.InitConfig("../../")
	if err != nil {
		return nil, err
	}

	voterInfo, err := api.GetVoterInfo(addr)
	if err != nil {
		return nil, err
	}

	return &voterInfo, nil
}

type mockLotusRepo struct {
	logger *zap.Logger
}

// GetAddrBalanceBySpecialHeight implements LotusRepo.
func (m *mockLotusRepo) GetAddrBalanceBySpecialHeight(ctx context.Context, addr string, netId int64, height int64) (string, error) {
	panic("unimplemented")
}

// GetBlockHeader implements LotusRepo.
func (m *mockLotusRepo) GetBlockHeader(ctx context.Context, netId int64, height int64) (models.BlockHeader, error) {
	return models.BlockHeader{
		Height:    height,
		Timestamp: time.Now().Unix(),
	}, nil
}

// GetClientBalanceByHeight implements LotusRepo.
func (m *mockLotusRepo) GetClientBalanceByHeight(ctx context.Context, netId int64, height int64) (types.StateMarketDeals, error) {
	m.logger.Debug("GetClientBalanceByHeight", zap.Any("netId", netId), zap.Any("height", height))
	return types.StateMarketDeals{
		"0": types.Deal{
			Proposal: types.Proposal{},
		},
	}, nil
}

// GetClientBalanceBySpecialHeight implements LotusRepo.
func (m *mockLotusRepo) GetClientBalanceBySpecialHeight(ctx context.Context, netId int64, height int64) (models.StateMarketDeals, error) {
	panic("unimplemented")
}

// GetMinerPowerByHeight implements LotusRepo.
func (m *mockLotusRepo) GetMinerPowerByHeight(ctx context.Context, netId int64, addr string, tipsetKey []interface{}) (models.LotusMinerPower, error) {
	panic("unimplemented")
}

// GetNewestHeight implements LotusRepo.
func (m *mockLotusRepo) GetNewestHeight(ctx context.Context, netId int64) (height int64, err error) {
	return 1, nil
}

// GetTipSetByHeight implements LotusRepo.
func (m *mockLotusRepo) GetTipSetByHeight(ctx context.Context, netId int64, height int64) ([]any, error) {
	panic("unimplemented")
}

// GetWalletBalanceByHeight implements LotusRepo.
func (m *mockLotusRepo) GetWalletBalanceByHeight(ctx context.Context, id string, netId int64, height int64) (string, error) {
	m.logger.Debug("GetWalletBalanceByHeight", zap.Any("id", id), zap.Any("netId", netId), zap.Any("height", height))
	return "1000", nil
}

type mockSyncRepo struct {
	AddrSyncedDateMap map[string][]string
	AddrPowerMap      map[string]map[string]models.SyncPower
	logger            *zap.Logger
}

func (m *mockSyncRepo) CreateSnapshot(ctx context.Context, netId int64, cid, day string) error {
	panic("implement me")
}

func (m *mockSyncRepo) SetDeveloperWeights(ctx context.Context, dateStr string, in map[string]int64) error {
	m.logger.Debug("SetDeveloperWeights", zap.String("dateStr", dateStr), zap.Any("in", in))
	return nil
}

func (m *mockSyncRepo) GetDeveloperWeights(ctx context.Context, dateStr string) (map[string]int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockSyncRepo) GetUserDeveloperWeights(ctx context.Context, dateStr string, username string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockSyncRepo) ExistDeveloperWeights(ctx context.Context, dateStr string) (bool, error) {
	m.logger.Debug("ExistDeveloperWeights", zap.String("dateStr", dateStr))
	return false, nil
}

func (m *mockSyncRepo) GetAddrSyncedDate(ctx context.Context, netId int64, addr string) ([]string, error) {
	m.logger.Debug("GetAddrSyncedDate", zap.Any("netId", netId), zap.Any("addr", addr))
	return nil, nil
}

func (m *mockSyncRepo) SetAddrSyncedDate(ctx context.Context, netId int64, addr string, dates []string) error {
	m.logger.Debug("SetAddrSyncedDate", zap.Any("netId", netId), zap.Any("addr", addr), zap.Any("dates", dates))
	return nil
}

func (m *mockSyncRepo) GetAddrPower(ctx context.Context, netId int64, addr string) (map[string]models.SyncPower, error) {
	m.logger.Debug("GetAddrPower", zap.Any("netId", netId), zap.Any("addr", addr))
	return nil, nil
}

func (m *mockSyncRepo) GetTask(ctx context.Context, netID int64) (jetstream.MessageBatch, error) {
	m.logger.Debug("GetTask")
	return nil, nil
}

func (m *mockSyncRepo) AddTask(ctx context.Context, netID int64, task *models.Task) error {
	m.logger.Debug("AddTask", zap.Any("task", task))
	return nil
}

func (m *mockSyncRepo) SetAddrPower(ctx context.Context, netId int64, addr string, in map[string]models.SyncPower) error {
	config.InitLogger()
	err := config.InitConfig("../../")
	if err != nil {
		return err
	}
	client, err := data.NewRedisClient()
	if err != nil {
		return err
	}

	key := fmt.Sprintf(constant.RedisAddrPower, netId, addr)
	mm := make(map[string]string)
	for k, v := range in {
		j, err := json.Marshal(v)
		if err != nil {
			return err
		}
		mm[k] = string(j)
	}
	err = client.HSet(ctx, key, mm).Err()
	if err != nil {
		return err
	}
	zap.L().Info("SetAddrPower", zap.Any("in", in))

	return nil
}

func (m *mockBaseRepo) SetAllAddrSyncedDateMap(ctx context.Context, netId int64, addrSyncedDate map[string][]string) error {
	m.AddrSyncedDateMap = addrSyncedDate
	return nil
}

func (m *mockSyncRepo) GetAllAddrSyncedDateMap(ctx context.Context, netId int64) (map[string][]string, error) {
	return m.AddrSyncedDateMap, nil
}

func (m *mockBaseRepo) GetDateHeightMap(ctx context.Context, netId int64) (map[string]int64, error) {
	raw := `{"20240415":1528817,"20240416":1531697,"20240417":1534577,"20240418":1537457,"20240419":1540337,"20240420":1543217,"20240421":1546097,"20240422":1548977,"20240423":1551857,"20240424":1554737,"20240425":1557617,"20240426":1560497,"20240427":1563377,"20240428":1566257,"20240429":1569137,"20240430":1572017,"20240501":1574897,"20240502":1577777,"20240503":1580657,"20240504":1583537,"20240505":1586417,"20240506":1589297,"20240507":1592177,"20240508":1595057,"20240509":1597937,"20240510":1600817,"20240511":1603697,"20240512":1606577,"20240513":1609457,"20240514":1612337,"20240515":1615217,"20240516":1618097,"20240517":1620977,"20240518":1623857,"20240519":1626737,"20240520":1629617,"20240521":1632497,"20240522":1635377,"20240523":1638257,"20240524":1641137,"20240525":1644017,"20240526":1646897,"20240527":1649777,"20240528":1652657,"20240529":1655537,"20240530":1658417,"20240531":1661297,"20240601":1664177,"20240602":1667057,"20240603":1669937,"20240604":1672817,"20240605":1675697,"20240606":1678577,"20240607":1681457,"20240608":1684337,"20240609":1687217,"20240610":1690097,"20240611":1692977,"20240612":1695857,"20240613":1698017}`
	err := json.Unmarshal([]byte(raw), &m.DhMap)
	if err != nil {
		return nil, err
	}

	return m.DhMap, nil
}

// GetSnapshotBackupList implements MysqlRepo.
func (m *mockmMysqlRepo) GetSnapshotBackupList(ctx context.Context, chainId int64) ([]models.SnapshotBackupTbl, error) {
	return []models.SnapshotBackupTbl{
		{
			RawData: "test",
		},
	}, nil
}

// UpdateSnapshotBackup implements MysqlRepo.
func (m *mockmMysqlRepo) UpdateSnapshotBackup(ctx context.Context, in models.SnapshotBackupTbl) error {
	m.logger.Debug("UpdateSnapshotBackup", zap.Any("in", in))
	return nil
}

// CreateSnapshotBackup implements MysqlRepo.
func (m *mockmMysqlRepo) CreateSnapshotBackup(ctx context.Context, in models.SnapshotBackupTbl) error {
	m.logger.Debug("CreateSnapshotBackup", zap.Any("in", in))
	return nil
}

var _ MysqlRepo = (*mockmMysqlRepo)(nil)

type mockmMysqlRepo struct {
	logger *zap.Logger
}

func getMockSyncService(t *testing.T) *SyncService {
	config.InitLogger()
	err := config.InitConfig("../..")
	config.Client.W3Client.Proof = "../../proof.ucan"
	assert.NoError(t, err)

	b := &mockBaseRepo{}
	r := &mockSyncRepo{
		AddrPowerMap:      make(map[string]map[string]models.SyncPower),
		AddrSyncedDateMap: map[string][]string{},
		logger:            zap.L().Named("Testing Mock SyncRepo"),
	}
	m := &mockmMysqlRepo{
		logger: zap.L().Named("Testing Mock MysqlRepo"),
	}
	l := &mockLotusRepo{
		logger: zap.L().Named("Testing Mock LotusRepo"),
	}

	return &SyncService{
		baseRepo:  b,
		syncRepo:  r,
		mysqlRepo: m,
		lotusRepo: l,
	}
}
func TestDiffAddrList(t *testing.T) {
	l1, l2 := []string{"a", "b"}, []string{"a", "b", "c", "d"}
	_, d2 := lo.Difference(l1, l2)

	expected := []string{"c", "d"}
	assert.Equal(t, expected, d2)
}

func TestCalMissDates(t *testing.T) {
	config.InitLogger()

	d1 := []string{"20251111", "20251110"}
	d2 := []string{""}
	d3 := []string{"20230101,20230102"}
	d4 := []string{"20240604", "20240602", "20240530"}

	r1 := utils.CalDateList(time.Date(2025, 11, 11, 0, 0, 0, 0, time.Local), 5, d1)
	r2 := utils.CalDateList(time.Date(2024, 6, 4, 0, 0, 0, 0, time.Local), 60, d2)
	r3 := utils.CalDateList(time.Date(2024, 6, 4, 0, 0, 0, 0, time.Local), 60, d3)
	r4 := utils.CalDateList(time.Date(2024, 6, 5, 0, 0, 0, 0, time.Local), 7, d4)

	excepted1 := []string{"20251109", "20251108", "20251107"}
	assert.Equal(t, excepted1, r1)

	excepted2 := []string{"20240604", "20240603", "20240602", "20240601", "20240531", "20240530", "20240529", "20240528", "20240527", "20240526", "20240525", "20240524", "20240523", "20240522", "20240521", "20240520", "20240519", "20240518", "20240517", "20240516", "20240515", "20240514", "20240513", "20240512", "20240511", "20240510", "20240509", "20240508", "20240507", "20240506", "20240505", "20240504", "20240503", "20240502", "20240501", "20240430", "20240429", "20240428", "20240427", "20240426", "20240425", "20240424", "20240423", "20240422", "20240421", "20240420", "20240419", "20240418", "20240417", "20240416", "20240415", "20240414", "20240413", "20240412", "20240411", "20240410", "20240409", "20240408", "20240407", "20240406"}
	assert.Equal(t, excepted2, r2)

	excepted3 := []string{"20240604", "20240603", "20240602", "20240601", "20240531", "20240530", "20240529", "20240528", "20240527", "20240526", "20240525", "20240524", "20240523", "20240522", "20240521", "20240520", "20240519", "20240518", "20240517", "20240516", "20240515", "20240514", "20240513", "20240512", "20240511", "20240510", "20240509", "20240508", "20240507", "20240506", "20240505", "20240504", "20240503", "20240502", "20240501", "20240430", "20240429", "20240428", "20240427", "20240426", "20240425", "20240424", "20240423", "20240422", "20240421", "20240420", "20240419", "20240418", "20240417", "20240416", "20240415", "20240414", "20240413", "20240412", "20240411", "20240410", "20240409", "20240408", "20240407", "20240406"}
	assert.Equal(t, excepted3, r3)

	excepted4 := []string{"20240605", "20240603", "20240601", "20240531"}
	assert.Equal(t, excepted4, r4)
}

func TestSyncAllAddrPower(t *testing.T) {
	sync := getMockSyncService(t)
	_, err := sync.GetAllAddrInfoList(context.Background(), 314159)
	assert.NoError(t, err)

	err = sync.SyncAllAddrPower(context.Background(), testNetID)
	assert.NoError(t, err)
}

func getSyncService(t *testing.T) *SyncService {
	config.InitConfig("../../")

	config.InitLogger()
	config.Client.W3Client.Proof = "../../proof.ucan"
	manager, err := data.NewGoEthClientManager(config.Client.Network)
	assert.NoError(t, err)
	// init datasource
	redisClient, err := data.NewRedisClient()
	assert.NoError(t, err)

	jetstreamClient, err := data.NewJetstreamClient()
	assert.NoError(t, err)

	baseRepo := repo.NewBaseRepoImpl(manager, redisClient)
	assert.NoError(t, err)
	syncRepo, err := repo.NewSyncRepoImpl(314159, redisClient, jetstreamClient)
	assert.NoError(t, err)

	mysalRepo := repo.NewMysqlRepoImpl(data.NewMysql())
	lotusRepo := repo.NewLotusRPCRepo(redisClient)
	return NewSyncService(baseRepo, syncRepo, mysalRepo, lotusRepo)

}

func TestUploadPowerToIPFS(t *testing.T) {

	syncService := getMockSyncService(t)

	err := syncService.UploadPowerToIPFS(context.Background(), 314159, data.NewW3Client())
	assert.NoError(t, err)
}

func TestExistDeveloperWeight(t *testing.T) {
	sync := getMockSyncService(t)
	res, err := sync.ExistDeveloperWeight(context.Background(), "20250331")
	assert.NoError(t, err)
	assert.Equal(t, false, res)
}

func TestSyncDeveloperWeight(t *testing.T) {
	sync := getMockSyncService(t)
	err := sync.SyncDeveloperWeight(context.Background(), "20250311")
	assert.NoError(t, err)
}

func TestSyncDateHeight(t *testing.T) {
	err := getSyncService(t).SyncDateHeight(context.Background(), 314159)
	assert.NoError(t, err)
}

func TestGetActorBalance(t *testing.T) {
	walletBalance, clientBalance, err := getMockSyncService(t).GetActorBalance(context.Background(), "f0114153", 314159, 2539001)
	assert.NoError(t, err)
	assert.Equal(t, "1000", walletBalance)
	assert.Equal(t, "0", clientBalance)
}

func TestAddAddrPowerTaskToMQ(t *testing.T) {
	err := getMockSyncService(t).AddAddrPowerTaskToMQ(context.Background(), 314159, "20250331")
	assert.NoError(t, err)
}

func TestUploadSnapshotInfoByDay(t *testing.T) {
	allPower := map[string]any{
		"total": "1000",
		"miners": map[string]any{
			"t01000": map[string]any{
				"total": "1000",
				"power": "1000",
			},
		},
	}
	height, err := getMockSyncService(t).UploadSnapshotInfoByDay(context.Background(), allPower, "20240511", 314159)
	assert.NoError(t, err)
	assert.Equal(t, int64(1603697), height)
}
