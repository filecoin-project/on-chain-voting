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
	"errors"
	"fmt"
	"github.com/nats-io/nats.go/jetstream"
	"math/big"
	"power-snapshot/constant"
	models "power-snapshot/internal/model"
	"power-snapshot/utils"
	"slices"
	"strconv"
	"time"

	"github.com/golang-module/carbon"
	"github.com/samber/lo"
	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"
)

type SyncRepo interface {
	// GetAllAddrSyncedDateMap KEY:ADDR:[]SYNCED_STRING
	GetAllAddrSyncedDateMap(ctx context.Context, netId int64) (map[string][]string, error)
	GetAddrSyncedDate(ctx context.Context, netId int64, addr string) ([]string, error)
	SetAddrSyncedDate(ctx context.Context, netId int64, addr string, dates []string) error

	// GetAddrPower KEY ADDR:DATE:POWER
	GetAddrPower(ctx context.Context, netId int64, addr string) (map[string]*models.SyncPower, error)
	// SetAddrPower KEY ADDR:DATE:POWER
	SetAddrPower(ctx context.Context, netId int64, addr string, in map[string]*models.SyncPower) error

	AddTask(ctx context.Context, netId int64, task *models.Task) error
	GetTask(ctx context.Context, netId int64) (jetstream.MessageBatch, error)

	SetDeveloperWeights(ctx context.Context, dateStr string, in map[string]int64) error
	GetDeveloperWeights(ctx context.Context, dateStr string) (map[string]int64, error)
	GetUserDeveloperWeights(ctx context.Context, dateStr string, username string) (int64, error)
	ExistDeveloperWeights(ctx context.Context, dateStr string) (bool, error)
}

type SyncService struct {
	baseRepo BaseRepo
	syncRepo SyncRepo
}

func NewSyncService(baseRepo BaseRepo, syncRepo SyncRepo) *SyncService {
	return &SyncService{
		baseRepo: baseRepo,
		syncRepo: syncRepo,
	}
}

func (s *SyncService) GetAllAddrInfoList(ctx context.Context, netID int64, IdPrefix string) ([]models.AddrInfo, error) {
	list, err := s.baseRepo.ListVoterAddr(ctx, netID)
	if err != nil {
		zap.L().Error("failed to pending sync addr list", zap.Error(err))
		return nil, err
	}

	pendingSyncAddrList := make([]models.AddrInfo, 0)
	for _, addr := range list {
		voteInfo, err := s.baseRepo.GetVoteInfo(ctx, netID, addr)
		if err != nil {
			zap.L().Error("failed to get vote info, skip this addr", zap.String("addr", addr), zap.Error(err))
			continue
		}
		pendingSyncAddrList = append(pendingSyncAddrList, models.AddrInfo{
			Addr:          addr,
			IdPrefix:      IdPrefix,
			ActionIDs:     voteInfo.ActorIds,
			MinerIDs:      voteInfo.MinerIds,
			GithubAccount: voteInfo.GithubAccount,
		})
	}

	return pendingSyncAddrList, nil
}

func (s *SyncService) GetAddrInfo(ctx context.Context, netID int64, addr string) (*models.AddrInfo, error) {
	ethClient, err := s.baseRepo.GetEthClient(ctx, netID)
	if err != nil {
		zap.L().Error("failed to get ethClient", zap.Error(err))
		return nil, err
	}

	list, err := s.baseRepo.ListVoterAddr(ctx, netID)
	if err != nil {
		zap.L().Error("failed to pending sync addr list", zap.Error(err))
		return nil, err
	}
	if !slices.Contains(list, addr) {
		return nil, fmt.Errorf("addr %s not exist in voter", addr)
	}
	voteInfo, err := s.baseRepo.GetVoteInfo(ctx, netID, addr)
	if err != nil {
		zap.L().Error("failed to get vote info, skip this addr", zap.String("addr", addr), zap.Error(err))
		return nil, err
	}
	m := &models.AddrInfo{
		Addr:          addr,
		IdPrefix:      ethClient.IdPrefix,
		ActionIDs:     voteInfo.ActorIds,
		MinerIDs:      voteInfo.MinerIds,
		GithubAccount: voteInfo.GithubAccount,
	}

	return m, nil
}

func (s *SyncService) SyncAllAddrPower(ctx context.Context, netID int64) error {
	ethClient, err := s.baseRepo.GetEthClient(ctx, netID)
	if err != nil {
		zap.L().Error("failed to get ethClient", zap.Error(err))
		return err
	}
	dhMap, err := s.baseRepo.GetDateHeightMap(ctx, netID)
	if err != nil {
		zap.L().Error("failed to get dates-height map", zap.Error(err))
		return err
	}
	pendingSyncedAddr, err := s.GetAllAddrInfoList(ctx, netID, ethClient.IdPrefix)
	if err != nil {
		zap.L().Error("failed to get GetAllAddrInfoList", zap.Error(err))
		return err
	}
	zap.L().Info("pendingSyncedAddr", zap.Any("count", len(pendingSyncedAddr)))
	dateMap, err := s.syncRepo.GetAllAddrSyncedDateMap(ctx, netID)
	if err != nil {
		return err
	}

	taskList := make([]models.Task, 0, len(pendingSyncedAddr)*3)
	// make task meta info
	for _, info := range pendingSyncedAddr {
		task := models.Task{
			UID:           fmt.Sprintf("%s", info.Addr),
			Address:       info.Addr,
			SubTasks:      nil,
			GithubAccount: info.GithubAccount,
		}
		// cal miss data date
		dm, ok := dateMap[info.Addr]
		if !ok {
			dm = []string{}
		}
		missDates := CalMissDates(dm)
		if len(missDates) == 0 {
			zap.L().Info("address no miss date to sync", zap.String("addr", info.Addr))
		}

		subTaskList := make([]models.SubTask, 0)
		for _, date := range missDates {
			for _, actorID := range info.ActionIDs {
				actorIDStr := info.IdPrefix + strconv.FormatUint(actorID, 10)
				subTaskList = append(subTaskList, models.SubTask{
					UID:         fmt.Sprintf("%s-%s-%s", info.Addr, date, actorIDStr),
					Address:     info.Addr,
					DateStr:     date,
					BlockHeight: dhMap[date],
					Typ:         "actor",
					IDStr:       actorIDStr,
				})
			}
			for _, minerID := range info.MinerIDs {
				minerIDStr := info.IdPrefix + strconv.FormatUint(minerID, 10)
				subTaskList = append(subTaskList, models.SubTask{
					UID:         fmt.Sprintf("%s-%s-%s", info.Addr, date, minerIDStr),
					Address:     info.Addr,
					DateStr:     date,
					BlockHeight: dhMap[date],
					Typ:         "miner",
					IDStr:       minerIDStr,
				})
			}
		}
		task.SubTasks = subTaskList
		taskList = append(taskList, task)
	}

	zap.L().Info("task", zap.Any("count", len(pendingSyncedAddr)))
	for _, task := range taskList {
		err := s.syncRepo.AddTask(ctx, netID, &task)
		if err != nil {
			zap.S().Error("failed to add task, skip this addr", zap.String("addr", task.Address), zap.Error(err))
			continue
		}
	}

	return nil
}

func (s *SyncService) SyncAddrPower(ctx context.Context, netID int64, addr string) error {
	dhMap, err := s.baseRepo.GetDateHeightMap(ctx, netID)
	if err != nil {
		zap.L().Error("failed to get dates-height map", zap.Error(err))
		return err
	}
	info, err := s.GetAddrInfo(ctx, netID, addr)
	if err != nil {
		zap.L().Error("failed to get GetAllAddrInfoList", zap.Error(err))
		return err
	}
	dateMap, err := s.syncRepo.GetAllAddrSyncedDateMap(ctx, netID)
	if err != nil {
		return err
	}

	// cal miss data date
	dm, ok := dateMap[info.Addr]
	if !ok {
		dm = []string{}
	}
	missDate := CalMissDates(dm)
	if len(missDate) == 0 {
		zap.L().Info("address no miss date to sync", zap.String("addr", addr))
	}

	subTaskList := make([]models.SubTask, 0, len(missDate)*3)
	for _, date := range CalMissDates(dm) {
		for _, actorID := range info.ActionIDs {
			actorIDStr := info.IdPrefix + strconv.FormatUint(actorID, 10)
			subTaskList = append(subTaskList, models.SubTask{
				UID:         fmt.Sprintf("%s-%s-%s", info.Addr, date, actorIDStr),
				Address:     info.Addr,
				DateStr:     date,
				BlockHeight: dhMap[date],
				Typ:         "actor",
				IDStr:       actorIDStr,
			})
		}
		for _, minerID := range info.MinerIDs {
			minerIDStr := info.IdPrefix + strconv.FormatUint(minerID, 10)
			subTaskList = append(subTaskList, models.SubTask{
				UID:         fmt.Sprintf("%s-%s-%s", info.Addr, date, minerIDStr),
				Address:     info.Addr,
				DateStr:     date,
				BlockHeight: dhMap[date],
				Typ:         "miner",
				IDStr:       minerIDStr,
			})
		}
	}

	task := models.Task{
		UID:           fmt.Sprintf("%s", info.Addr),
		Address:       info.Addr,
		SubTasks:      subTaskList,
		GithubAccount: info.GithubAccount,
	}
	err = s.syncRepo.AddTask(ctx, netID, &task)
	if err != nil {
		zap.S().Error("failed to add task, skip this addr", zap.String("addr", task.Address), zap.Error(err))
		return err
	}

	return nil
}

func calDateList(startTime time.Time, duration int, dates []string) []string {
	base := carbon.FromStdTime(startTime)
	allDatesList := make([]string, 0, duration)
	for i := 0; i < duration; i++ {
		allDatesList = append(allDatesList, base.ToShortDateString())
		base = base.SubDay()
	}

	diff, _ := lo.Difference(allDatesList, dates)
	return diff
}

func CalMissDates(dates []string) []string {
	return calDateList(time.Now().Add(-(24 * time.Hour)), constant.DataExpiredDuration, dates)
}

func (s *SyncService) SyncDateHeight(ctx context.Context, netID int64) error {
	rpcClient, err := s.baseRepo.GetLotusClientByHashKey(ctx, netID, carbon.Now().SubDay().EndOfDay().ToDateString())
	if err != nil {
		zap.L().Error("failed to get ethClient", zap.Error(err))
		return err
	}

	dhMap, err := getDateHeight(ctx,
		rpcClient,
		carbon.Now().SubDay().EndOfDay().ToStdTime(),
		constant.DataExpiredDuration)
	if err != nil {
		return err
	}

	err = s.baseRepo.SetDateHeightMap(ctx, netID, dhMap)
	if err != nil {
		zap.L().Error("failed to set dates-height map", zap.Error(err))
		return err
	}
	return nil
}

func getDateHeight(ctx context.Context, rpc jsonrpc.RPCClient, startTime time.Time, duration int) (map[string]int64, error) {
	nowHeight, err := utils.GetNewestHeight(ctx, rpc)
	if err != nil {
		zap.L().Error("failed to get newest height", zap.Error(err))
		return nil, err
	}
	nowHeightInfo, err := utils.GetBlockHeader(ctx, rpc, nowHeight)
	if err != nil {
		zap.L().Error("failed to get newest height info", zap.Error(err))
		return nil, err
	}
	if nowHeightInfo.Timestamp < startTime.Unix() {
		return nil, errors.New("the start time is later than the latest block time")
	}

	dh := make(map[string]int64)
	needSyncDates := calDateList(startTime, duration, []string{})
	needSyncDatesLength := len(needSyncDates)
	if needSyncDatesLength == 0 {
		return dh, nil
	}

	curDatePos := 0
	// Assume the block time is 30 seconds and subtract the number of blocks equivalent to two hours each time.
	for i := nowHeight; i > 0; i = i - (2 * 3600 / 30) {
		time.Sleep(50 * time.Millisecond)
		iter, err := utils.GetBlockHeader(ctx, rpc, i)
		if err != nil {
			zap.L().Error("failed to get newest height info", zap.Error(err))
			continue
		}
		// format time
		syncDate := carbon.ParseByLayout(needSyncDates[curDatePos], carbon.ShortDateLayout).StartOfDay()
		blockDate := carbon.CreateFromTimestamp(iter.Timestamp).StartOfDay()
		zap.L().Info("now", zap.Int64("block_height", i), zap.Any("syncDate", syncDate.ToShortDateString()), zap.Any("blockDate", blockDate.ToShortDateString()))

		if blockDate.Lt(syncDate) {
			return nil, errors.New(fmt.Sprintf("The sync date(%s) is later than the block date(%s), this sync date will never be synchronized.",
				syncDate.ToShortDateString(),
				blockDate.ToShortDateString()))
		}
		// If the current block time is later than the start time, skip these blocks.
		if blockDate.Gt(syncDate) {
			continue
		}
		if blockDate.Eq(syncDate) {
			dh[syncDate.ToShortDateString()] = iter.Height
			curDatePos++
		}
		// Break the loop when all dates are synchronized.
		if curDatePos == needSyncDatesLength {
			break
		}
	}

	return dh, nil
}

func (s *SyncService) SyncAllDeveloperWeight(ctx context.Context) error {
	base := carbon.Now().SubDay().EndOfDay()
	for i := 0; i < constant.DataExpiredDuration; i++ {
		m, err := GetDeveloperWeights(base.ToStdTime())
		if err != nil {
			return err
		}
		err = s.syncRepo.SetDeveloperWeights(ctx, base.ToShortDateString(), m)
		if err != nil {
			zap.S().Error("failed to set developer power", zap.String("date", base.ToShortDateString()), zap.Error(err))
			return err
		}
		base = base.SubDay()
	}

	return nil

}

func (s *SyncService) SyncDeveloperWeight(ctx context.Context, dayStr string) error {
	base := carbon.ParseByLayout(dayStr, carbon.ShortDateLayout).EndOfDay()
	m, err := GetDeveloperWeights(base.ToStdTime())
	if err != nil {
		return err
	}
	err = s.syncRepo.SetDeveloperWeights(ctx, base.ToShortDateString(), m)
	if err != nil {
		zap.S().Error("failed to set developer power", zap.String("date", base.ToShortDateString()), zap.Error(err))
		return err
	}

	return nil
}

func (s *SyncService) ExistDeveloperWeight(ctx context.Context, dayStr string) (bool, error) {
	base := carbon.ParseByLayout(dayStr, carbon.ShortDateLayout).EndOfDay()
	exist, err := s.syncRepo.ExistDeveloperWeights(ctx, base.ToShortDateString())
	if err != nil {
		zap.S().Error("failed to exist developer weight", zap.String("date", base.ToShortDateString()), zap.Error(err))
		return false, err
	}
	return exist, nil
}

func (s *SyncService) StartSyncWorker(ctx context.Context, netID int64) error {
	zap.S().Info("starting sync worker ", netID)

	for {
		taskMsg, err := s.syncRepo.GetTask(ctx, netID)
		if err != nil {
			zap.S().Error("failed to get task", err)
			return err
		}
		for msg := range taskMsg.Messages() {
			var task models.Task
			err := json.Unmarshal(msg.Data(), &task)
			if err != nil {
				zap.S().Error("failed to unmarshal task", err)
				return err
			}

			jsonRpcClient, err := s.baseRepo.GetLotusClientByHashKey(ctx, netID, task.Address)
			if err != nil {
				zap.L().Error("failed to get ethClient", zap.Error(err))
				return err
			}

			go func() {
				zap.L().Info("start sync address", zap.Any("task_uid", task.UID))
				power, err := s.syncRepo.GetAddrPower(ctx, netID, task.Address)
				if err != nil {
					zap.S().Error("failed to get addr power ", zap.Error(err))
					return
				}
				result := make(map[string]*models.SyncPower)
				for _, st := range task.SubTasks {
					temp := &models.SyncPower{
						Address:          st.Address,
						DateStr:          st.DateStr,
						GithubAccount:    task.GithubAccount,
						DeveloperPower:   big.NewInt(0),
						SpPower:          big.NewInt(0),
						ClientPower:      big.NewInt(0),
						TokenHolderPower: big.NewInt(0),
						BlockHeight:      st.BlockHeight,
					}
					if st.Typ == "actor" {
						walletBalance, clientBalance, err := GetActorPower(ctx, jsonRpcClient, st.IDStr, st.BlockHeight)
						if err != nil {
							zap.L().Error("failed to get actor power", zap.Error(err))
							return
						}
						if len(walletBalance) != 0 {
							wl, ok := big.NewInt(0).SetString(walletBalance, 10)
							if !ok {
								zap.L().Error("failed to parse wallet balance", zap.Error(err), zap.String("wallet_balance", walletBalance))
								return
							}
							temp.TokenHolderPower = temp.TokenHolderPower.Add(temp.TokenHolderPower, wl)
						}
						if len(clientBalance) != 0 {
							wl, ok := big.NewInt(0).SetString(clientBalance, 10)
							if !ok {
								zap.L().Error("failed to parse client balance", zap.Error(err), zap.String("client_balance", clientBalance))
								return
							}
							temp.ClientPower = temp.ClientPower.Add(temp.ClientPower, wl)
						}
					}
					if st.Typ == "miner" {
						minerBalance, err := GetMinerPower(ctx, jsonRpcClient, st.IDStr, st.BlockHeight)
						if err != nil {
							zap.L().Error("failed to get miner power", zap.Error(err))
							return
						}

						if len(minerBalance) != 0 {
							ml, ok := big.NewInt(0).SetString(minerBalance, 10)
							if !ok {
								zap.L().Error("failed to parse miner power", zap.Error(err))
								return
							}

							temp.SpPower = temp.SpPower.Add(temp.SpPower, ml)
						}
					}
					if _, exists := result[st.DateStr]; !exists {
						result[st.DateStr] = temp
					} else {
						result[st.DateStr].SpPower.Add(result[st.DateStr].SpPower, temp.SpPower)
						result[st.DateStr].TokenHolderPower.Add(result[st.DateStr].TokenHolderPower, temp.TokenHolderPower)
						result[st.DateStr].ClientPower.Add(result[st.DateStr].ClientPower, temp.ClientPower)
					}
				}

				dates := make([]string, 0, len(result))
				for dateStr, syncPower := range result {
					power[dateStr] = syncPower
					dates = append(dates, dateStr)
				}
				err = s.syncRepo.SetAddrPower(ctx, netID, task.Address, power)
				if err != nil {
					zap.S().Error("failed to set addr power", zap.Error(err))
					return
				}

				oldDates, err := s.syncRepo.GetAddrSyncedDate(ctx, netID, task.Address)
				if err != nil {
					zap.S().Error("failed to get addr synced date", zap.Error(err))
					return
				}
				newDates := append(oldDates, dates...)
				slices.Sort(newDates)
				newDates = lo.Uniq(newDates)

				err = s.syncRepo.SetAddrSyncedDate(ctx, netID, task.Address, newDates)
				if err != nil {
					zap.S().Error("failed to set addr synced", zap.Error(err))
					return
				}

				err = msg.Ack()
				if err != nil {
					zap.S().Error("failed to ack task", zap.Error(err))
					return
				}
			}()
		}
	}
}
