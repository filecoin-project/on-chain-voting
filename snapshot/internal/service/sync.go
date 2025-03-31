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
	"math/big"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/golang-module/carbon"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"power-snapshot/api"
	"power-snapshot/constant"
	"power-snapshot/internal/data"
	models "power-snapshot/internal/model"
	"power-snapshot/utils"
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

	GetDict(ctx context.Context, netId int64) (int64, error)
	SetDelegateEvent(ctx context.Context, netId int64, createDelegateEvent []models.CreateDelegateEvent, deleteDelegateEvent []models.DeleteDelegateEvent, endBlock int64) error
	GetDelegateEvent(ctx context.Context, netId int64, addr string, maxBlockHeight int64) (models.CreateDelegateEvent, models.DeleteDelegateEvent, error)
}

type SyncService struct {
	baseRepo  BaseRepo
	syncRepo  SyncRepo
	mysqlRepo MysqlRepo
	lotusRepo LotusRepo
}

func NewSyncService(baseRepo BaseRepo, syncRepo SyncRepo, mysqlRepo MysqlRepo, lotusRepo LotusRepo) *SyncService {
	return &SyncService{
		baseRepo:  baseRepo,
		syncRepo:  syncRepo,
		mysqlRepo: mysqlRepo,
		lotusRepo: lotusRepo,
	}
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
	missDate := utils.CalMissDates(dm)
	if len(missDate) == 0 {
		zap.L().Info("address no miss date to sync", zap.String("addr", addr))
	}

	subTaskList := make([]models.SubTask, 0, len(missDate)*3)
	for _, date := range utils.CalMissDates(dm) {
		for _, actorID := range info.ActionIDs {
			actorIDStr := info.IdPrefix + strconv.FormatUint(actorID, 10)
			subTaskList = append(subTaskList, models.SubTask{
				UID:         fmt.Sprintf("%s-%s-%s", info.Addr, date, actorIDStr),
				Address:     info.Addr,
				DateStr:     date,
				BlockHeight: dhMap[date],
				Typ:         constant.TaskActionActor,
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
				Typ:         constant.TaskActionMiner,
				IDStr:       minerIDStr,
			})
		}
	}

	task := models.Task{
		UID:           info.Addr,
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

/**
 * @Description: Sync the height of the block
 * @param ctx  context.Context
 * @param netID int64
 * @return error
 */
func (s *SyncService) SyncDateHeight(ctx context.Context, netID int64) error {
	dhMap, err := s.getDateHeight(ctx,
		netID,
		carbon.Now().SubDay().EndOfDay().ToStdTime(),
		constant.DataExpiredDuration)
	if err != nil {
		return err
	}

	zap.L().Info("get date-height map", zap.Any("date height map", dhMap))

	err = s.baseRepo.SetDateHeightMap(ctx, netID, dhMap)
	if err != nil {
		zap.L().Error("failed to set dates-height map", zap.Error(err))
		return err
	}

	zap.L().Info("sync date height success", zap.Int64("chain id", netID))
	return nil
}

func (s *SyncService) getDateHeight(ctx context.Context, netId int64, syncEndTime time.Time, syncCountedDays int) (map[string]int64, error) {
	newestHeight, err := s.lotusRepo.GetNewestHeight(ctx, netId)
	if err != nil {
		zap.L().Error("failed to get newest height", zap.Error(err))
		return nil, err
	}

	newestHeightInfo, err := s.lotusRepo.GetBlockHeader(ctx, netId, newestHeight)
	if err != nil {
		zap.L().Error("failed to get newest height info", zap.Error(err))
		return nil, err
	}

	if newestHeightInfo.Timestamp < syncEndTime.Unix() {
		return nil, errors.New("the latest block time is earlier than the sync time, please check the chain network")
	}

	dh := make(map[string]int64)
	needSyncDates :=  utils.CalDateList(syncEndTime, syncCountedDays, []string{})

	if len(needSyncDates) == 0 {
		return dh, nil
	}

	curDatePos := 0
	// Assume the block time is 30 seconds and subtract the number of blocks equivalent to two hours each time.
	for height := newestHeight; height > 0; height = height - (constant.TwoHoursBlockNumber) {
		// If the current block time is earlier than the start time, skip these blocks.
		time.Sleep(50 * time.Millisecond)
		blockHeader, err := s.lotusRepo.GetBlockHeader(ctx, netId, height)
		if err != nil {
			zap.L().Error("failed to get height info", zap.Int64("height", height), zap.Error(err))
			continue
		}
		// format time
		syncDate := carbon.ParseByLayout(needSyncDates[curDatePos], carbon.ShortDateLayout).StartOfDay()
		blockDate := carbon.CreateFromTimestamp(blockHeader.Timestamp).StartOfDay()
		zap.L().Info("now", zap.Int64("block_height", height), zap.Any("syncDate", syncDate.ToShortDateString()), zap.Any("blockDate", blockDate.ToShortDateString()))

		if blockDate.Lt(syncDate) {
			return nil, fmt.Errorf(
				"the sync date(%s) is later than the block date(%s), this sync date will never be synchronized",
				syncDate.ToShortDateString(),
				blockDate.ToShortDateString(),
			)
		}
		// If the current block time is later than the start time, skip these blocks.
		if blockDate.Gt(syncDate) {
			continue
		}
		if blockDate.Eq(syncDate) {
			dh[syncDate.ToShortDateString()] = blockHeader.Height
			curDatePos++
		}
		// Break the loop when all dates are synchronized.
		if curDatePos == len(needSyncDates) {
			break
		}
	}

	return dh, nil
}

func (s *SyncService) GetAllAddrInfoList(ctx context.Context, netID int64, IdPrefix string) ([]models.AddrInfo, error) {
	list, err := api.GetAllVoterAddresss(netID)
	if err != nil {
		zap.L().Error("failed to pending sync addr list", zap.Error(err))
		return nil, err
	}

	pendingSyncAddrList := make([]models.AddrInfo, 0)
	for _, addr := range list {
		voteInfo, err := api.GetVoterInfo(addr)
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

	list, err := api.GetAllVoterAddresss(netID)
	if err != nil {
		zap.L().Error("failed to pending sync addr list", zap.Error(err))
		return nil, err
	}
	if !slices.Contains(list, addr) {
		return nil, fmt.Errorf("addr %s not exist in voter", addr)
	}
	voteInfo, err := api.GetVoterInfo(addr)
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
			UID:           info.Addr,
			Address:       info.Addr,
			SubTasks:      nil,
			GithubAccount: info.GithubAccount,
		}
		// cal miss data date
		dm, ok := dateMap[info.Addr]
		if !ok {
			dm = []string{}
		}
		missDates := utils.CalMissDates(dm)
		if len(missDates) == 0 {
			zap.L().Info("address no miss date to sync", zap.String("addr", info.Addr))
		} else {
			zap.L().Info("address add task", zap.String("addr", info.Addr))
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
					Typ:         constant.TaskActionActor,
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
					Typ:         constant.TaskActionMiner,
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
		zap.L().Info("address add task", zap.String("addr", task.Address))
		if err != nil {
			zap.S().Error("failed to add task, skip this addr", zap.String("addr", task.Address), zap.Error(err))
			continue
		}
	}

	return nil
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

	zap.L().Info("Sync developer weight success", zap.String("date", dayStr))
	return nil
}

func (s *SyncService) ExistDeveloperWeight(ctx context.Context, dayStr string) (bool, error) {
	base := carbon.ParseByLayout(dayStr, carbon.ShortDateLayout).EndOfDay()
	exist, err := s.syncRepo.ExistDeveloperWeights(ctx, base.ToShortDateString())
	if err != nil {
		zap.S().Error("failed to exist developer weight", zap.String("date", base.ToShortDateString()), zap.Error(err))
		return false, err
	}

	zap.L().Info("exist developer weight", zap.String("date", base.ToShortDateString()), zap.Bool("exist", exist))
	return exist, nil
}

// StartSyncWorker starts a worker to process synchronization tasks for a specific network ID.
// It continuously fetches tasks from the sync repository and processes them concurrently.
// Each task involves fetching power data for an address, calculating power metrics (e.g., token holder power, client power, SP power),
// and updating the results in the sync repository.
// The worker handles subtasks of type "actor" and "miner" to calculate specific power metrics.
// If any step fails, an error is logged, and the task is skipped.
// On successful completion of a task, the worker acknowledges the message to mark it as processed.
// The worker runs indefinitely until the context is canceled or an unrecoverable error occurs.
func (s *SyncService) StartSyncWorker(ctx context.Context, netID int64) error {
	for {
		// Get nats message queue object
		taskMsg, err := s.syncRepo.GetTask(ctx, netID)
		zap.S().Info("loop task")
		if err != nil {
			zap.S().Error("failed to get task", err)
			time.Sleep(5 * time.Second)
			continue
		}

		var wg sync.WaitGroup
		// Process each message in the task concurrently.
		for taskMsg := range taskMsg.Messages() {
			zap.L().Info("get task", zap.Any("task", taskMsg))
			wg.Add(1)
			go func(msg jetstream.Msg) {
				defer wg.Done()

				// Unmarshal the task message into a Task struct.
				var task models.Task
				err := json.Unmarshal(msg.Data(), &task)
				if err != nil {
					zap.S().Error("failed to unmarshal task", err)

					if err := taskMsg.Ack(); err != nil {
						zap.S().Error("failed to ack task", err)
					}

					return
				}

				zap.L().Info("start sync address", zap.Any("task_uid", task.UID))

				// Fetch the existing power data form redis for the address.
				power, err := s.syncRepo.GetAddrPower(ctx, netID, task.Address)
				if err != nil {
					zap.S().Error("failed to get addr power ", zap.Error(err))
					return
				}

				// Initialize a map to store the results of power calculations.
				result := make(map[string]*models.SyncPower)
				for _, st := range task.SubTasks {
					// Initialize a SyncPower struct for the subtask.
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

					// Handle subtasks of type "actor".
					if st.Typ == constant.TaskActionActor {

						walletBalance, clientBalance, err := s.GetActorPower(ctx, st.IDStr, netID, st.BlockHeight)
						if err != nil {
							zap.L().Error("failed to get actor power", zap.Error(err))
							return
						}

						// Parse and add wallet balance to token holder power.
						if len(walletBalance) != 0 {
							wl, ok := big.NewInt(0).SetString(walletBalance, 10)
							if !ok {
								zap.L().Error("failed to parse wallet balance", zap.Error(err), zap.String("wallet_balance", walletBalance))
								return
							}
							temp.TokenHolderPower = temp.TokenHolderPower.Add(temp.TokenHolderPower, wl)
						}

						// Parse and add client balance to client power.
						if len(clientBalance) != 0 {
							wl, ok := big.NewInt(0).SetString(clientBalance, 10)
							if !ok {
								zap.L().Error("failed to parse client balance", zap.Error(err), zap.String("client_balance", clientBalance))
								return
							}
							temp.ClientPower = temp.ClientPower.Add(temp.ClientPower, wl)
						}
					}

					// Handle subtasks of type "miner".
					if st.Typ == constant.TaskActionMiner {
						tipsetKey, err := s.lotusRepo.GetTipSetByHeight(ctx, netID, st.BlockHeight)
						if err != nil {
							zap.L().Error("failed to get tipset key", zap.Error(err))
							return
						}
						minerPower, err := s.lotusRepo.GetMinerPowerByHeight(ctx, netID, st.IDStr, tipsetKey)
						if err != nil {
							zap.L().Error("failed to get miner power", zap.Error(err))
							return
						}

						// Parse and add miner balance to SP power.
						if len(minerPower.MinerPower.RawBytePower) != 0 {
							ml, ok := big.NewInt(0).SetString(minerPower.MinerPower.RawBytePower, 10)
							if !ok {
								zap.L().Error("failed to parse miner power", zap.Error(err))
								return
							}
							temp.SpPower = temp.SpPower.Add(temp.SpPower, ml)
						}
					}

					// Merge results for the same date.
					if _, exists := result[st.DateStr]; !exists {
						result[st.DateStr] = temp
					} else {
						result[st.DateStr].SpPower.Add(result[st.DateStr].SpPower, temp.SpPower)
						result[st.DateStr].TokenHolderPower.Add(result[st.DateStr].TokenHolderPower, temp.TokenHolderPower)
						result[st.DateStr].ClientPower.Add(result[st.DateStr].ClientPower, temp.ClientPower)
					}
				}

				// Update the power data for the address.
				dates := make([]string, 0, len(result))
				for dateStr, syncPower := range result {
					power[dateStr] = syncPower
					dates = append(dates, dateStr)
				}

				// Save the updated power data to the sync repository.
				err = s.syncRepo.SetAddrPower(ctx, netID, task.Address, power)
				if err != nil {
					zap.S().Error("failed to set addr power", zap.Error(err))
					return
				}

				// Update the list of synced dates for the address.
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

				zap.L().Info("start sync address success", zap.Any("task_uid", task.UID))

				// Acknowledge the message to mark it as processed.
				err = msg.Ack()
				if err != nil {
					zap.S().Error("failed to ack task", zap.Error(err))
					return
				}
			}(taskMsg)
		}

		zap.L().Info("The sync worker task is running finished")
		wg.Wait()
	}
}

func (s *SyncService) UploadSnapshotInfoByDay(ctx context.Context, allPower map[string]any, day string, chainId int64) (int64, error) {
	dateHeight, err := s.baseRepo.GetDateHeightMap(ctx, chainId)
	if err != nil {
		zap.L().Error("failed to get snapshot height", zap.Error(err))
		return 0, err
	}

	snapshotHeight, exist := dateHeight[day]
	if !exist {
		zap.L().Error("snapshot height not exist", zap.Int64("chainId", chainId), zap.String("day", day), zap.Error(err))
		return 0, errors.New("snapshot height not exist")
	}

	jsonStr, err := json.Marshal(allPower)
	if err != nil {
		zap.L().Error("failed to marshal snapshot", zap.Error(err))
		return 0, err
	}

	if err := s.mysqlRepo.CreateSnapshotBackup(ctx, models.SnapshotBackupTbl{
		Day:     day,
		ChainId: chainId,
		Height:  snapshotHeight,
		RawData: string(jsonStr),
		Status:  constant.SnapshotBackupSync,
	}); err != nil {
		zap.L().Error("failed to create snapshot", zap.Error(err))
		return 0, err
	}

	return snapshotHeight, nil
}

func (s *SyncService) UploadPowerToIPFS(ctx context.Context, chainId int64, w3client *data.W3Client) error {
	zap.L().Info("start to upload power to ipfs", zap.Int64("chainId", chainId))
	snapshotList, err := s.mysqlRepo.GetSnapshotBackupList(ctx, chainId)
	if err != nil {
		zap.L().Error("failed to get snapshot list", zap.Error(err))
		return err
	}

	zap.L().Info("upload power to ipfs", zap.Int("count snapshots", len(snapshotList)))
	for _, snapshot := range snapshotList {
		if snapshot.Status < constant.SnapshotBackupRetry {
			// upload to ipfs
			rawBytes := []byte(snapshot.RawData)
			if len(rawBytes) == 0 {
				zap.L().Error("invalid bytes", zap.Error(err))
				return err
			}

			cid, err := w3client.UploadByte(rawBytes)
			if err != nil {
				zap.L().Error("failed to upload power to ipfs", zap.Error(err))
				snapshot.Status += 1
			} else {
				snapshot.Cid = cid
				snapshot.Status = constant.SnapshotBackupSyncd
			}

			if err := s.mysqlRepo.UpdateSnapshotBackup(ctx, snapshot); err != nil {
				zap.L().Error("failed to update snapshot", zap.Error(err))
				return err
			}

			zap.L().Info(
				"upload power to ipfs success",
				zap.Int64("chainId", chainId),
				zap.String("cid", cid),
				zap.Int("status", snapshot.Status),
			)
		}
	}
	return nil
}

func (s *SyncService) GetActorPower(ctx context.Context, actorId string, netId, height int64) (string, string, error) {
	// Fetch the wallet balance
	walletBalance, err := s.lotusRepo.GetWalletBalanceByHeight(ctx, actorId, netId, height)
	if err != nil {
		zap.L().Error("failed to get wallet balance", zap.Error(err))
		return "", "", err
	}

	t, err := s.lotusRepo.GetClientBalanceByHeight(ctx, netId, height)
	if err != nil {
		return walletBalance, "", err
	}

	var clientBalance int64
	for _, v := range t {
		if v.Proposal.Client == actorId && v.Proposal.EndEpoch > height && v.Proposal.VerifiedDeal {
			clientBalance += v.Proposal.PieceSize
		}
	}

	return walletBalance, strconv.FormatInt(clientBalance, 10), nil
}
