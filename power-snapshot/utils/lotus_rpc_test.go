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
package utils

import (
	"context"
	"power-snapshot/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	address = "0x7C24ca5FBA6f1E228a520911476746D25Be5EdbE"
	id      = "t099523"
)

func TestIDFromAddress(t *testing.T) {
	config.InitLogger()
	err := config.InitConfig("../")
	assert.Nil(t, err)
	manager, err := NewGoEthClientManager(config.Client.Network)
	assert.Nil(t, err)
	client, err := manager.GetClient(314159)
	assert.Nil(t, err)
	lotusRpcClient := NewClient(client.QueryRpc[0])

	res, err := IDFromAddress(context.Background(), lotusRpcClient, address)
	assert.Nil(t, err)

	expectedID := "t065744"
	assert.Equal(t, res, expectedID)

	zap.L().Info("result", zap.Any("id", res))
}

func TestWalletBalanceByHeight(t *testing.T) {
	config.InitLogger()
	err := config.InitConfig("../")
	assert.Nil(t, err)
	manager, err := NewGoEthClientManager(config.Client.Network)
	assert.Nil(t, err)
	client, err := manager.GetClient(314159)
	assert.Nil(t, err)

	lotusRpcClient := NewClient(client.QueryRpc[0])
	rsp, err := GetWalletBalanceByHeight(context.Background(), lotusRpcClient, id, 1730418)
	assert.Nil(t, err)
	// expectedBalance := "467051509071593317720"
	// assert.Equal(t, expectedBalance, rsp)

	zap.L().Info("result", zap.Any("balance", rsp))
}

func TestGetMinerPowerByHeight(t *testing.T) {
	config.InitLogger()
	err := config.InitConfig("../")
	assert.Nil(t, err)
	manager, err := NewGoEthClientManager(config.Client.Network)
	assert.Nil(t, err)
	client, err := manager.GetClient(314159)
	assert.Nil(t, err)

	lotusRpcClient := NewClient(client.QueryRpc[0])
	rsp, err := GetMinerPowerByHeight(context.Background(), lotusRpcClient, "t03751", 1669397)
	assert.Nil(t, err)

	assert.True(t, rsp.HasMinPower)
	expectedRawBytePower := "1060925641588736"
	assert.Equal(t, expectedRawBytePower, rsp.MinerPower.RawBytePower)

	zap.L().Info("result", zap.Any("miner", rsp))
}

func TestGetClientBalanceByHeight(t *testing.T) {
	config.InitLogger()
	err := config.InitConfig("../")
	assert.Nil(t, err)
	manager, err := NewGoEthClientManager(config.Client.Network)
	assert.Nil(t, err)
	client, err := manager.GetClient(314159)
	assert.Nil(t, err)

	lotusRpcClient := NewClient(client.QueryRpc[0])
	rsp, err := GetClientBalanceByHeight(context.Background(), lotusRpcClient, id, 1669397)
	assert.Nil(t, err)

	zap.L().Info("result", zap.Any("client", rsp))
}

func TestGetNewestHeightAndTipset(t *testing.T) {
	config.InitLogger()
	err := config.InitConfig("../")
	assert.Nil(t, err)
	manager, err := NewGoEthClientManager(config.Client.Network)
	assert.Nil(t, err)
	client, err := manager.GetClient(314159)
	assert.Nil(t, err)

	lotusRpcClient := NewClient(client.QueryRpc[0])
	rsp, err := GetNewestHeight(context.Background(), lotusRpcClient)
	println(rsp)
	assert.Nil(t, err)
}

func TestGetBlockHeader(t *testing.T) {
	config.InitLogger()
	err := config.InitConfig("../")
	assert.Nil(t, err)
	manager, err := NewGoEthClientManager(config.Client.Network)
	assert.Nil(t, err)
	client, err := manager.GetClient(314159)
	assert.Nil(t, err)

	lotusRpcClient := NewClient(client.QueryRpc[0])
	rsp, err := GetBlockHeader(context.Background(), lotusRpcClient, 1560498)
	assert.Nil(t, err)

	zap.L().Info("result", zap.Any("block", rsp))
}
