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
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/golang-module/carbon"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"power-snapshot/api"
	"power-snapshot/constant"
	"power-snapshot/internal/data"
	models "power-snapshot/internal/model"
	"power-snapshot/utils"
)

type SyncRepo interface {
	// GetAllAddrSyncedDateMap retrieves synchronization dates mapping for all addresses under specified network.(e.g. {"addr1": ["20250301", "20250302"]})
	GetAllAddrSyncedDateMap(ctx context.Context, netId int64) (map[string][]string, error)
	GetAddrSyncedDate(ctx context.Context, netId int64, addr string) ([]string, error)
	SetAddrSyncedDate(ctx context.Context, netId int64, addr string, dates []string) error

	// GetAddrPower KEY ADDR:DATE:POWER
	GetAddrPower(ctx context.Context, netId int64, addr string) (map[string]models.SyncPower, error)
	// SetAddrPower KEY ADDR:DATE:POWER
	SetAddrPower(ctx context.Context, netId int64, addr string, in map[string]models.SyncPower) error

	// nats message quue
	// Add stream message
	AddTask(ctx context.Context, netId int64, task *models.Task) error
	// Get stream message
	GetTask(ctx context.Context, netId int64) (jetstream.MessageBatch, error)

	SetDeveloperWeights(ctx context.Context, dateStr string, in map[string]int64) error
	GetDeveloperWeights(ctx context.Context, dateStr string) (map[string]int64, error)
	GetUserDeveloperWeights(ctx context.Context, dateStr string, username string) (int64, error)
	ExistDeveloperWeights(ctx context.Context, dateStr string) (bool, error)
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

//* task
/* -------------------------------------------------------------------------- */
/*                                	SyncPower                                 */
/* -------------------------------------------------------------------------- */

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

	dm, err := s.baseRepo.GetDateHeightMap(ctx, netId)
	if err != nil {
		zap.L().Error("failed to get date height map", zap.Error(err))
		return nil, err
	}

	var alreadySyncDay []string
	for k := range dm {
		alreadySyncDay = append(alreadySyncDay, k)
	}

	if newestHeightInfo.Timestamp < syncEndTime.Unix() {
		return nil, errors.New("the latest block time is earlier than the sync time, please check the chain network")
	}

	needSyncDates := utils.CalDateList(syncEndTime, syncCountedDays, alreadySyncDay)

	if len(needSyncDates) == 0 {
		return dm, nil
	}

	curDatePos := 0
	retryCount := 0
	// Assume the block time is 30 seconds and subtract the number of blocks equivalent to two hours each time.
	for height := newestHeight; height > 0; height = height - (constant.TwoHoursBlockNumber) {
		// If the current block time is earlier than the start time, skip these blocks.
		time.Sleep(50 * time.Millisecond)
		blockHeader, err := s.lotusRepo.GetBlockHeader(ctx, netId, height)
		if err != nil {
			zap.L().Error("failed to get height info", zap.Int64("height", height), zap.Error(err))
			time.Sleep(1 * time.Second)
			retryCount++
			if retryCount > 5 {
				return nil, err
			}

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
			dm[syncDate.ToShortDateString()] = blockHeader.Height
			curDatePos++
		}
		// Break the loop when all dates are synchronized.
		if curDatePos == len(needSyncDates) {
			break
		}
	}

	return dm, nil
}

func (s *SyncService) SyncAllAddrPower(ctx context.Context, netID int64) error {
	dhMap, err := s.baseRepo.GetDateHeightMap(ctx, netID)
	if err != nil {
		zap.L().Error("failed to get dates-height map", zap.Error(err))
		return err
	}

	pendingSyncedAddr, err := s.GetAllAddrInfoList(ctx, netID)
	if err != nil {
		zap.L().Error("failed to get GetAllAddrInfoList", zap.Error(err))
		return err
	}
	zap.L().Info("pendingSyncedAddr", zap.Any("count", len(pendingSyncedAddr)))

	dateMap, err := s.syncRepo.GetAllAddrSyncedDateMap(ctx, netID)
	if err != nil {
		zap.L().Error("failed to get all address synced date map", zap.Error(err))
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
			continue
		} else {
			zap.L().Info("address add task", zap.String("addr", info.Addr))
		}

		subTaskList := make([]models.SubTask, 0)
		for _, date := range missDates {
			for _, actorID := range info.ActionIDs {
				subTaskList = append(subTaskList, models.SubTask{
					UID:         fmt.Sprintf("%s-%s-%s", info.Addr, date, actorID),
					Address:     info.Addr,
					DateStr:     date,
					BlockHeight: dhMap[date],
					Typ:         constant.TaskActionActor,
					IDStr:       actorID,
				})
			}
		}

		for _, minerID := range info.MinerIDs {
			date := carbon.Now().SubDay().ToShortDateString()
			subTaskList = append(subTaskList, models.SubTask{
				UID:         fmt.Sprintf("%s-%s-%s", info.Addr, date, minerID),
				Address:     info.Addr,
				DateStr:     date,
				BlockHeight: dhMap[date],
				Typ:         constant.TaskActionMiner,
				IDStr:       minerID,
			})
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

func (s *SyncService) GetAllAddrInfoList(ctx context.Context, netID int64) ([]models.AddrInfo, error) {
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
			ActionIDs:     voteInfo.ActorIds,
			MinerIDs:      voteInfo.MinerIds,
			GithubAccount: voteInfo.GithubAccount,
		})
	}

	return pendingSyncAddrList, nil
}

/* ------------------------------ SyncPower END ----------------------------- */

//* task
/* -------------------------------------------------------------------------- */
/*                            SyncDevWeightStepDay                            */
/* -------------------------------------------------------------------------- */

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

func (s *SyncService) SyncDeveloperWeight(ctx context.Context, dayStr string) error {
	dayEndTime := carbon.ParseByLayout(dayStr, carbon.ShortDateLayout).EndOfDay()
	m, commits, err := FetchDeveloperWeights(dayEndTime.ToStdTime())
	if err != nil {
		return err
	}

	if len(m) == 0 {
		zap.L().Info("no developer weight to sync", zap.String("date", dayEndTime.ToShortDateString()))
		return nil
	}

	err = s.syncRepo.SetDeveloperWeights(ctx, dayEndTime.ToShortDateString(), m)
	if err != nil {
		zap.S().Error("failed to set developer power", zap.String("date", dayEndTime.ToShortDateString()), zap.Error(err))
		return err
	}

	err = s.baseRepo.SaveDeveloperWeightsToFile(ctx, dayStr, commits)
	if err != nil {
		zap.S().Error("failed to set developer power", zap.String("date", dayStr), zap.Error(err))
		return err
	}
	zap.L().Info("Sync developer weight success", zap.String("date", dayStr))
	return nil
}

/* ------------------------ SyncDevWeightStepDay END ------------------------ */

//* task
/* -------------------------------------------------------------------------- */
/*                              UploadPowerToIPFS                             */
/* -------------------------------------------------------------------------- */
func (s *SyncService) UploadPowerToIPFS(ctx context.Context, chainId int64, w3client *data.W3Client) error {
	zap.L().Info("start to upload power to ipfs", zap.Int64("chainId", chainId))
	snapshotList, err := s.mysqlRepo.GetSnapshotBackupList(ctx, chainId)
	if err != nil {
		zap.L().Error("failed to get snapshot list", zap.Error(err))
		return err
	}

	zap.L().Info("upload power to ipfs", zap.Int("count snapshots", len(snapshotList)))
	for _, snapshot := range snapshotList {
		if snapshot.Status < constant.RetryCount {
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

/* -------------------------- UploadPowerToIPFS END ------------------------- */

//*main - sync all voter power
/* -------------------------------------------------------------------------- */
/*                               StartSyncWorker                              */
/* -------------------------------------------------------------------------- */

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
		if err != nil {
			zap.S().Error("failed to get task", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if taskMsg == nil {
			return fmt.Errorf("not found task")
		}

		var eg errgroup.Group
		// Process each message in the task concurrently.
		eg.SetLimit(10)
		for taskMsg := range taskMsg.Messages() {
			eg.Go(func() error {
				// Unmarshal the task message into a Task struct.
				var task models.Task
				err := json.Unmarshal(taskMsg.Data(), &task)
				if err != nil {
					zap.S().Error("failed to unmarshal task", err)

					if err := taskMsg.Ack(); err != nil {
						zap.S().Error("failed to ack task", err)
					}

					return err
				}

				zap.L().Info("start sync address", zap.Any("task_uid", task.UID))

				// Fetch the existing power data form redis for the address.
				power, err := s.syncRepo.GetAddrPower(ctx, netID, task.Address)
				if err != nil {
					zap.S().Error("failed to get addr power, task not ack", zap.Error(err))
					return err
				}

				// Initialize a map to store the results of power calculations.
				result := make(map[string]models.SyncPower)
				for _, subTask := range task.SubTasks {
					zap.L().Info(
						"start sync subtask",
						zap.String("subTask uid", subTask.UID),
						zap.String("sync date", subTask.DateStr),
						zap.Int64("block height", subTask.BlockHeight),
						zap.Int64("retry count", subTask.RetryCount),
						zap.String("sub task type", subTask.Typ),
					)
					// Initialize a SyncPower struct for the subtask.
					temp := models.SyncPower{
						Address:          subTask.Address,
						DateStr:          subTask.DateStr,
						GithubAccount:    task.GithubAccount,
						DeveloperPower:   big.NewInt(0),
						SpPower:          big.NewInt(0),
						ClientPower:      big.NewInt(0),
						TokenHolderPower: big.NewInt(0),
						BlockHeight:      subTask.BlockHeight,
					}

					// Handle subtasks of type "actor".
					if subTask.Typ == constant.TaskActionActor {
						walletBalance, clientBalance, err := s.GetActorBalance(ctx, subTask.IDStr, netID, subTask.BlockHeight)
						if err != nil {
							if strings.Contains(err.Error(), constant.ActorNotFound) {
								zap.L().Warn(
									"actor not found, continue",
									zap.String("subTask uid", subTask.UID),
									zap.Int64("height", subTask.BlockHeight),
									zap.String("actor id", subTask.IDStr),
								)

								continue
							}
							zap.L().Error(
								"failed to get actor power, task not ack",
								zap.String("subTask uid", subTask.UID),
								zap.Int64("height", subTask.BlockHeight),
								zap.String("actor id", subTask.IDStr),
								zap.Error(err),
							)
							return err
						}

						// Parse and add wallet balance to token holder power.
						temp.TokenHolderPower = temp.TokenHolderPower.Add(temp.TokenHolderPower, utils.StringToBigInt(walletBalance))
						// Parse and add client balance to client power.
						temp.ClientPower = temp.ClientPower.Add(temp.ClientPower, utils.StringToBigInt(clientBalance))

					}

					// Handle subtasks of type "miner".
					if subTask.Typ == constant.TaskActionMiner {
						tipsetKey, err := s.lotusRepo.GetTipSetByHeight(ctx, netID, subTask.BlockHeight)
						if err != nil {
							zap.L().Error("failed to get tipset key, task not ack", zap.String("subTask uid", subTask.UID), zap.Error(err))
							return err
						}

						minerPower, err := s.lotusRepo.GetMinerPowerByHeight(ctx, netID, subTask.IDStr, tipsetKey)
						if err != nil {
							if strings.Contains(err.Error(), constant.ActorNotFound) {
								zap.L().Warn(
									"actor not found, continue",
									zap.String("subTask uid", subTask.UID),
									zap.Int64("height", subTask.BlockHeight),
									zap.String("actor id", subTask.IDStr),
								)

								continue
							}

							zap.L().Error("failed to get miner power, task not ack", zap.String("subTask uid", subTask.UID), zap.Error(err))
							return err
						}

						// Parse and add miner balance to SP power.
						if len(minerPower.MinerPower.RawBytePower) != 0 {
							ml, ok := big.NewInt(0).SetString(minerPower.MinerPower.RawBytePower, 10)
							if !ok {
								zap.L().Error("failed to parse miner power, task not ack", zap.String("subTask uid", subTask.UID), zap.Error(err))
								return err
							}
							temp.SpPower = temp.SpPower.Add(temp.SpPower, ml)
						}
					}

					// Merge results for the same date.
					if _, exists := result[subTask.DateStr]; !exists {
						result[subTask.DateStr] = temp
					} else {
						result[subTask.DateStr].SpPower.Add(result[subTask.DateStr].SpPower, temp.SpPower)
						result[subTask.DateStr].TokenHolderPower.Add(result[subTask.DateStr].TokenHolderPower, temp.TokenHolderPower)
						result[subTask.DateStr].ClientPower.Add(result[subTask.DateStr].ClientPower, temp.ClientPower)
					}

					zap.L().Info("finish sync subtask", zap.String("subTask uid", subTask.UID), zap.String("sync date", subTask.DateStr), zap.Int64("block height", subTask.BlockHeight), zap.Int64("retry count", subTask.RetryCount), zap.String("sub task type", subTask.Typ))
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
					zap.S().Error("failed to set addr power, task not ack", zap.Error(err))
					return err
				}

				// Update the list of synced dates for the address.
				oldDates, err := s.syncRepo.GetAddrSyncedDate(ctx, netID, task.Address)
				if err != nil {
					zap.S().Error("failed to get addr synced date, task not ack", zap.Error(err))
					return err
				}

				newDates := append(oldDates, dates...)
				slices.Sort(newDates)
				newDates = lo.Uniq(newDates)

				err = s.syncRepo.SetAddrSyncedDate(ctx, netID, task.Address, newDates)
				if err != nil {
					zap.S().Error("failed to set addr synced, task not ack", zap.Error(err))
					return err
				}

				zap.L().Info("sync address success", zap.Any("task_uid", task.UID))

				// Acknowledge the message to mark it as processed.
				err = taskMsg.Ack()
				if err != nil {
					zap.S().Error("failed to ack task", zap.Error(err))
					return err
				}

				zap.L().Info("The sync worker task is running finished", zap.Any("task", task))
				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			zap.L().Error("failed to sync address", zap.Error(err))
			return err
		}

	}
}

// GetActorBalance get actor balance
func (s *SyncService) GetActorBalance(ctx context.Context, actorId string, netId, height int64) (string, string, error) {
	walletBalance, err := s.lotusRepo.GetWalletBalanceByHeight(ctx, actorId, netId, height)
	if err != nil {
		zap.L().Error("failed to get wallet balance", zap.String("actor id", actorId), zap.Int64("height", height), zap.Error(err))
		return "0", "0", err
	}

	t, err := s.lotusRepo.GetClientBalanceByHeight(ctx, netId, height)
	if err != nil {
		return walletBalance, "0", err
	}

	var clientBalance int64
	for _, v := range t {
		if v.Proposal.Client == actorId && v.Proposal.EndEpoch > height && v.Proposal.VerifiedDeal {
			clientBalance += v.Proposal.PieceSize
		}
	}

	return walletBalance, strconv.FormatInt(clientBalance, 10), nil
}

// add voter address sync power task to message queue
func (s *SyncService) AddAddrPowerTaskToMQ(ctx context.Context, netID int64, addr string) error {
	dhMap, err := s.baseRepo.GetDateHeightMap(ctx, netID)
	if err != nil {
		zap.L().Error("failed to get dates-height map", zap.Error(err))
		return err
	}

	info, err := s.GetAddrInfo(ctx, netID, addr)
	if err != nil {
		zap.L().Error("failed to get addr info", zap.Error(err))
		return err
	}

	addrSyncedDates, err := s.syncRepo.GetAddrSyncedDate(ctx, netID, addr)
	if err != nil {
		return err
	}

	missDate := utils.CalMissDates(addrSyncedDates)
	if len(missDate) == 0 {
		zap.L().Info("address no miss date to sync", zap.String("addr", addr))
		return nil
	}

	subTaskList := make([]models.SubTask, 0, len(missDate)*3)
	for _, date := range missDate {
		blockHeight, exits := dhMap[date]
		if !exits {
			continue
		}

		for _, actorID := range info.ActionIDs {
			subTaskList = append(subTaskList, models.SubTask{
				UID:         fmt.Sprintf("%s-%s-%s", info.Addr, date, actorID),
				Address:     info.Addr,
				DateStr:     date,
				BlockHeight: blockHeight,
				Typ:         constant.TaskActionActor,
				IDStr:       actorID,
			})
		}

		for _, minerID := range info.MinerIDs {
			subTaskList = append(subTaskList, models.SubTask{
				UID:         fmt.Sprintf("%s-%s-%s", info.Addr, date, minerID),
				Address:     info.Addr,
				DateStr:     date,
				BlockHeight: blockHeight,
				Typ:         constant.TaskActionMiner,
				IDStr:       minerID,
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

// GetAddrInfo retrieves address information including voting details and associated accounts
//
// Parameters:
//
//	ctx   : context.Context - Context for request cancellation and timeouts
//	netID : int64           - Network identifier used for chain-specific prefix
//	addr  : string          - Blockchain address to query
//
// Returns:
//
//	*models.AddrInfo - Structured address information containing:
//	                   - Address and network prefix
//	                   - Associated action IDs
//	                   - Miner IDs
//	                   - GitHub account
//	error           - Errors from voter info API or nil if successful
func (s *SyncService) GetAddrInfo(ctx context.Context, netID int64, addr string) (*models.AddrInfo, error) {
	// Retrieve voter information from external API
	voteInfo, err := api.GetVoterInfo(addr)
	if err != nil {
		// Log detailed error including the failing address
		zap.L().Error("failed to get vote info, skip this addr",
			zap.String("addr", addr),
			zap.Error(err))
		return nil, err
	}

	// Construct address information response
	m := &models.AddrInfo{
		Addr:          addr,                   // Original queried address
		ActionIDs:     voteInfo.ActorIds,      // Associated action identifiers
		MinerIDs:      voteInfo.MinerIds,      // Linked miner identifiers
		GithubAccount: voteInfo.GithubAccount, // Connected GitHub account
	}

	return m, nil
}

func (s *SyncService) SyncLatestDeveloperWeight(ctx context.Context) error {
	base := carbon.Now().SubDay().EndOfDay()
	exist, err := s.ExistDeveloperWeight(ctx, base.ToShortDateString())
	if err != nil {
		zap.L().Error("SyncDevWeightStepDay", zap.String("date", base.ToShortDateString()))
		return err
	}

	if exist {
		return nil
	}

	m, commits, err := FetchDeveloperWeights(base.ToStdTime())
	if err != nil {
		return err
	}

	if len(m) == 0 {
		zap.L().Info("no developer weight to sync", zap.String("date", base.ToShortDateString()))
		return nil
	}
	if err = s.syncRepo.SetDeveloperWeights(ctx, base.ToShortDateString(), m); err != nil {
		zap.S().Error("failed to set developer power", zap.String("date", base.ToShortDateString()), zap.Error(err))
		return err
	}

	if err := s.baseRepo.SaveDeveloperWeightsToFile(ctx, base.ToShortDateString(), commits); err != nil {
		zap.S().Error("failed to set developer power", zap.String("date", base.ToShortDateString()), zap.Error(err))
		return err
	}

	zap.L().Info("SyncLatestDeveloperWeight Success", zap.String("date", base.ToShortDateString()))
	return nil
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
		return snapshotHeight, errors.New("snapshot height not exist")
	}

	developerCommitsData, err := s.baseRepo.GetDeveloperWeights(ctx, day)
	if err != nil {
		if os.IsNotExist(err) {
			zap.L().Error("file not found", zap.String("filename", constant.DeveloperWeightsFilePrefix+day))
		} else {
			zap.L().Error("failed to get developer commits", zap.Error(err))
			return snapshotHeight, err
		}
	}

	allPower["devPower"] = developerCommitsData

	jsonStr, err := json.Marshal(allPower)
	if err != nil {
		zap.L().Error("failed to marshal snapshot", zap.Error(err))
		return snapshotHeight, err
	}

	if err := s.mysqlRepo.CreateSnapshotBackup(ctx, models.SnapshotBackupTbl{
		Day:     day,
		ChainId: chainId,
		Height:  snapshotHeight,
		RawData: string(jsonStr),
		Status:  constant.SnapshotBackupSync,
	}); err != nil {
		zap.L().Error("failed to create snapshot", zap.Error(err))
		return snapshotHeight, err
	}

	return snapshotHeight, nil
}
