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
	"power-snapshot/utils"
	"strconv"

	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"
)

// GetPower synchronizes the power snapshot for the given actor and miner at a specific height.
func GetPower(ctx context.Context, lotusRpcClient jsonrpc.RPCClient, actorId, minerId string, height int64) {
	// Fetch the wallet balance
	walletBalance, err := utils.GetWalletBalanceByHeight(ctx, lotusRpcClient, actorId, height)
	if err != nil {
		zap.L().Error("failed to get wallet balance", zap.Error(err))
		return
	}

	// Fetch the miner power if minerId is provided
	var power utils.MinerPower
	if minerId != "" {
		power, err = utils.GetMinerPowerByHeight(ctx, lotusRpcClient, minerId, height)
		if err != nil {
			zap.L().Error("failed to get miner power", zap.Error(err))
			return
		}
	}

	// Fetch the client balance
	clientBalance, err := utils.GetDealSumByHeightAndActorId(ctx, lotusRpcClient, actorId, height)
	if err != nil {
		zap.L().Error("failed to get client balance", zap.Error(err))
		return
	}

	zap.L().Info("fetched balances and power",
		zap.String("walletBalance", walletBalance),
		zap.String("minerPower", power.MinerPower.QualityAdjPower),
		zap.String("clientBalance", strconv.FormatInt(clientBalance, 10)),
	)

}

func GetActorPower(ctx context.Context, lotusRpcClient jsonrpc.RPCClient, actorId string, height int64) (string, string, error) {
	// Fetch the wallet balance
	walletBalance, err := utils.GetWalletBalanceByHeight(ctx, lotusRpcClient, actorId, height)
	if err != nil {
		zap.L().Error("failed to get wallet balance", zap.Error(err))
		return "", "", err
	}

	// Fetch the client balance
	clientBalance, err := utils.GetDealSumByHeightAndActorId(ctx, lotusRpcClient, actorId, height)
	if err != nil {
		zap.L().Error("failed to get client balance", zap.Error(err))
		return "", "", err
	}

	return walletBalance, strconv.FormatInt(clientBalance, 10), nil
}

func GetMinerPower(ctx context.Context, lotusRpcClient jsonrpc.RPCClient, minerId string, height int64) (string, error) {
	// Fetch the miner power if minerId is provided
	power, err := utils.GetMinerPowerByHeight(ctx, lotusRpcClient, minerId, height)
	if err != nil {
		zap.L().Error("failed to get miner power", zap.Error(err))
		return "", nil
	}
	return power.MinerPower.RawBytePower, nil
}
