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
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"powervoting-server/model"
)

// GetBlockTime retrieves the current timestamp from the Ethereum blockchain using the provided Ethereum client.
// It fetches the latest block number and then retrieves the block information to obtain the timestamp.
// The function returns the current timestamp as an int64 or an error if the operation fails.
func GetBlockTime(client *model.GoEthClient, syncedHeight int64) (int64, error) {
	ctx := context.Background()

	block, err := client.Client.BlockByNumber(ctx, big.NewInt(syncedHeight))
	if err != nil {
		return 0, err
	}

	return int64(block.Time()), nil
}

func GetOwnerIdByOracle(client *model.GoEthClient,actorID uint64, minerId []uint64) []uint64 {
	var res []uint64
	for _, miner := range minerId {
		data, err := client.ABI.OraclePowersAbi.Pack("getOwner", miner)
		if err != nil {
			zap.L().Error("Error packing getOwner", zap.Uint64("minerId", miner))
			continue
		}

		msg := ethereum.CallMsg{
			To:   &client.OraclePowersContract,
			Data: data,
		}

		result, err := client.Client.CallContract(context.Background(), msg, nil)
		if err != nil {
			zap.L().Error("Error calling getOwner", zap.Uint64("minerId", miner))
			continue
		}

		upPacked, err := client.ABI.OraclePowersAbi.Unpack("getOwner", result)
		if err != nil {
			zap.L().Error("Error unpacking getOwner", zap.Uint64("minerId", miner))
			continue
		}

		if len(upPacked) == 0 {
			zap.L().Error("getOwner returned empty result", zap.Uint64("minerId", miner))
			continue
		}

		if upPacked[0].(uint64) == uint64(actorID) {
			res = append(res, uint64(miner))
		}
	}

	return res
}

func GetActorIdByAddress(client *model.GoEthClient, address string) (uint64, error) {
	data, err := client.ABI.OraclePowersAbi.Pack("resolveEthAddress", common.HexToAddress(address))
	if err != nil {
		return 0, fmt.Errorf("error packing resolveEthAddress: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &client.OraclePowersContract,
		Data: data,
	}

	result, err := client.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return 0, fmt.Errorf("error calling resolveEthAddress: %v", err)
	}

	if len(result) == 0 {
		return 0, fmt.Errorf("resolveEthAddress returned empty result")
	}

	upPacked, err := client.ABI.OraclePowersAbi.Unpack("resolveEthAddress", result)
	if err != nil {
		return 0, fmt.Errorf("error unpacking resolveEthAddress: %v", err)
	}

	if len(upPacked) == 0 {
		return 0, fmt.Errorf("resolveEthAddress returned empty result")
	}

	return upPacked[0].(uint64), err
}
