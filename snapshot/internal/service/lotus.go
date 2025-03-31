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

	"go.uber.org/zap"

	models "power-snapshot/internal/model"
	"power-snapshot/internal/repo"
	"power-snapshot/utils/types"
)

type LotusRepo interface {
	GetTipSetByHeight(ctx context.Context, netId, height int64) ([]any, error)
	GetAddrBalanceBySpecialHeight(ctx context.Context, addr string, netId, height int64) (string, error)
	GetMinerPowerByHeight(ctx context.Context, netId int64, addr string, tipsetKey []interface{}) (models.LotusMinerPower, error)
	GetClientBalanceBySpecialHeight(ctx context.Context, netId, height int64) (models.StateMarketDeals, error)
	GetNewestHeight(ctx context.Context, netId int64) (height int64, err error)
	GetBlockHeader(ctx context.Context, netId, height int64) (models.BlockHeader, error)
	GetWalletBalanceByHeight(ctx context.Context, id string, netId, height int64) (string, error)
	GetClientBalanceByHeight(ctx context.Context, netId, height int64) (types.StateMarketDeals, error)
}

type LotusService struct {
	lotusRepo LotusRepo
	logger    *zap.Logger
}

var _ LotusRepo = (*repo.LotusRPCRepo)(nil)

func NewLotusService(lotusRepo LotusRepo) *LotusService {
	return &LotusService{
		lotusRepo: lotusRepo,
		logger:    zap.L().Named("Lotus Service"),
	}
}

func (l *LotusService) GetTipSetByHeight(ctx context.Context, netId, height int64) ([]any, error) {
	return l.lotusRepo.GetTipSetByHeight(ctx, netId, height)
}

func (l *LotusService) GetAddrBalanceBySpecialHeight(ctx context.Context, addr string, netId, height int64) (string, error) {
	res, err := l.lotusRepo.GetAddrBalanceBySpecialHeight(ctx, addr, netId, height)
	if err != nil {
		l.logger.Error(
			"error when get addr balance",
			zap.String("address", addr),
			zap.Int64("height", height),
			zap.Error(err),
		)

		return "", err
	}

	l.logger.Info("get addr balance success", zap.String("address", addr), zap.String("balance", res))
	return res, err
}

func (l *LotusService) GetMinerPowerByHeight(ctx context.Context, netId int64, addr string, tipsetKey []interface{}) (*models.LotusMinerPower, error) {
	res, err := l.lotusRepo.GetMinerPowerByHeight(ctx, netId, addr, tipsetKey)
	if err != nil {
		l.logger.Error("error when get miner power", zap.String("address", addr), zap.Error(err))
		return nil, err
	}

	l.logger.Info("get miner power success", zap.String("address", addr), zap.Any("power", res))
	return &res, nil
}

func (l *LotusService) GetClientBalanceBySpecialHeight(ctx context.Context, netId, height int64) (*models.StateMarketDeals, error) {
	res, err := l.lotusRepo.GetClientBalanceBySpecialHeight(ctx, netId, height)
	if err != nil {
		l.logger.Error("error when get client balance", zap.Int64("height", height), zap.Error(err))
		return nil, err
	}

	l.logger.Info("get client balance success", zap.Int64("height", height), zap.Any("balance", res))
	return &res, nil
}

func (l *LotusService) GetNewestHeight(ctx context.Context, netId int64) (height int64, err error) {
	res, err := l.lotusRepo.GetNewestHeight(ctx, netId)
	if err != nil {
		l.logger.Error("error when get newest height", zap.Error(err))
		return 0, err
	}

	l.logger.Info("get newest height success", zap.Int64("height", res))
	return res, nil
}

func (l *LotusService) GetBlockHeader(ctx context.Context, netId, height int64) (*models.BlockHeader, error) {
	res, err := l.lotusRepo.GetBlockHeader(ctx, netId, height)
	if err != nil {
		l.logger.Error("error when get block header", zap.Int64("height", height), zap.Error(err))
		return nil, err
	}

	l.logger.Info("get block header success", zap.Int64("height", height), zap.Any("header", res))
	return &res, nil
}

func (l *LotusService) GetWalletBalanceByHeight(ctx context.Context, id string, netId, height int64) (string, error) {
	res, err := l.lotusRepo.GetWalletBalanceByHeight(ctx, id, netId, height)
	if err != nil {
		l.logger.Error("error when get wallet balance", zap.String("address", id), zap.Int64("height", height), zap.Error(err))
		return "", err
	}

	l.logger.Info("get wallet balance success", zap.String("address", id), zap.String("balance", res))
	return res, nil
}

func (l *LotusService) GetClientBalanceByHeight(ctx context.Context, netId, height int64) (*types.StateMarketDeals, error) {
	res, err := l.lotusRepo.GetClientBalanceByHeight(ctx, netId, height)
	if err != nil {
		l.logger.Error("error when get client balance", zap.Int64("height", height), zap.Error(err))
		return nil, err
	}

	l.logger.Info("get client balance success", zap.Int64("height", height), zap.Any("balance", res))
	return &res, nil
}
