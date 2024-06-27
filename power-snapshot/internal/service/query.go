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
	"errors"
	"fmt"
	"math/big"
	"power-snapshot/constant"
	models "power-snapshot/internal/model"

	"github.com/golang-module/carbon"
	"go.uber.org/zap"
)

type QueryRepo interface {
	GetAddressPower(ctx context.Context, netId int64, address string, dayStr string) (*models.SyncPower, error)
	GetDeveloperWeights(ctx context.Context, dateStr string) (map[string]int64, error)
}

type QueryService struct {
	baseRepo  BaseRepo
	queryRepo QueryRepo
	syncSrv   *SyncService
}

func NewQueryService(baseRepo BaseRepo, queryRepo QueryRepo, sync *SyncService) *QueryService {
	return &QueryService{
		baseRepo:  baseRepo,
		queryRepo: queryRepo,
		syncSrv:   sync,
	}
}

func (q *QueryService) GetAddressPower(ctx context.Context, netId int64, address string, dayCount int32) (*models.SyncPower, error) {
	if dayCount > constant.DataExpiredDuration {
		return nil, errors.New("day count is too long")
	}
	dayStr := carbon.Now().SubDays(int(dayCount)).EndOfDay().ToShortDateString()
	dayTime := carbon.Now().SubDays(int(dayCount)).EndOfDay().ToStdTime()
	power, err := q.queryRepo.GetAddressPower(ctx, netId, address, dayStr)
	if err != nil {
		zap.L().Error("error getting address power ", zap.Error(err))
		return nil, err
	}
	if power == nil {
		// get power from chain now and sync all power sync
		err := q.syncSrv.SyncAddrPower(ctx, netId, address)
		if err != nil {
			zap.L().Error("error getting address power", zap.Error(err))
			return nil, err
		}
		// start get power
		jsonRpc, err := q.baseRepo.GetLotusClientByHashKey(ctx, netId, address)
		if err != nil {
			zap.L().Error("error getting eth client", zap.Error(err))
			return nil, err
		}
		dh, err := q.baseRepo.GetDateHeightMap(ctx, netId)
		if err != nil {
			zap.L().Error("error getting date height map", zap.Error(err))
			return nil, err
		}
		height := dh[dayStr]
		info, err := q.syncSrv.GetAddrInfo(ctx, netId, address)
		if err != nil {
			zap.L().Error("error getting address info", zap.Error(err))
			return nil, err
		}
		power = &models.SyncPower{
			GithubAccount:    info.GithubAccount,
			Address:          info.Addr,
			DateStr:          dayStr,
			DeveloperPower:   big.NewInt(0),
			SpPower:          big.NewInt(0),
			ClientPower:      big.NewInt(0),
			TokenHolderPower: big.NewInt(0),
			BlockHeight:      height,
		}

		for _, actionId := range info.ActionIDs {
			idStr := fmt.Sprintf("%s%d", info.IdPrefix, actionId)

			walletBalance, clientBalance, err := GetActorPower(ctx, jsonRpc, idStr, height)
			if err != nil {
				zap.L().Error("failed to get actor power", zap.Error(err))
				return nil, err
			}
			if len(walletBalance) != 0 {
				wl, ok := big.NewInt(0).SetString(walletBalance, 10)
				if !ok {
					zap.L().Error("failed to parse wallet balance", zap.Error(err), zap.String("wallet_balance", walletBalance))
					return nil, errors.New("failed to parse wallet balance")
				}
				power.TokenHolderPower = power.TokenHolderPower.Add(power.TokenHolderPower, wl)
			}
			if len(clientBalance) != 0 {
				wl, ok := big.NewInt(0).SetString(clientBalance, 10)
				if !ok {
					zap.L().Error("failed to parse client balance", zap.Error(err), zap.String("client_balance", clientBalance))
					return nil, errors.New("failed to parse client balance")
				}
				power.ClientPower = power.ClientPower.Add(power.ClientPower, wl)
			}
		}

		for _, minerId := range info.MinerIDs {
			idStr := fmt.Sprintf("%s%d", info.IdPrefix, minerId)
			if len(info.MinerIDs) != 0 {
				minerBalance, err := GetMinerPower(ctx, jsonRpc, idStr, height)
				if err != nil {
					zap.L().Error("failed to get miner power", zap.Error(err))
					return nil, err
				}

				if len(minerBalance) != 0 {
					ml, ok := big.NewInt(0).SetString(minerBalance, 10)
					if !ok {
						zap.L().Error("failed to parse miner power", zap.Error(err))
						return nil, errors.New("failed to parse miner power")
					}

					power.SpPower = power.SpPower.Add(power.SpPower, ml)
				}
			}
		}
		if len(info.GithubAccount) != 0 {
			dwm, err := GetDeveloperWeights(dayTime)
			if err != nil {
				zap.L().Error("error getting developer weights", zap.Error(err))
				power.DeveloperPower = big.NewInt(0)
			} else {
				if weight, ok := dwm[info.GithubAccount]; ok {
					power.DeveloperPower = big.NewInt(weight)
				}
			}
		}
	}

	if len(power.GithubAccount) != 0 {
		devWeight, err := q.queryRepo.GetDeveloperWeights(ctx, dayStr)
		if err != nil {
			zap.L().Error("error getting developer weight", zap.Error(err))
			return nil, err
		}
		// if this day's history have not synced, return 0 and log error
		if devWeight == nil {
			zap.L().Error("no developer weight synced from github",
				zap.String("date", dayStr),
				zap.String("address", address),
				zap.String("account", power.GithubAccount))
		}

		if w, ok := devWeight[power.GithubAccount]; ok {
			power.DeveloperPower = big.NewInt(w)
		}
	}

	return power, nil
}
