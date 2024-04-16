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
	"backend/contract"
	"backend/models"
	"backend/utils"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"
)

// F4Address retrieves F4 tasks from the smart contract and processes each task asynchronously.
func F4Address(ethClient models.GoEthClient, lotusRpcClient jsonrpc.RPCClient) {
	f4TaskList, err := utils.GetF4Tasks(ethClient)
	if err != nil {
		zap.L().Error("failed to get task id", zap.Error(err))
		return
	}

	if len(f4TaskList) == 0 {
		zap.L().Info("f4 task list is empty")
		return
	}

	for _, taskId := range f4TaskList {
		go ProcessingF4AddressTaskId(taskId, ethClient, lotusRpcClient)
	}
}

// ProcessingF4AddressTaskId processes a specific F4 task ID retrieved from the smart contract.
func ProcessingF4AddressTaskId(taskId *big.Int, ethClient models.GoEthClient, lotusRpcClient jsonrpc.RPCClient) {
	zap.L().Info("processing f4 address task", zap.String("task id", taskId.String()))

	args := []interface{}{taskId}
	ethAddressMap, err := GetContractMapping(models.F4TaskIdToAddress, ethClient, args)
	if err != nil {
		zap.L().Error("failed to get eth address", zap.Error(err))
		if err := DeleteTask(taskId, ethClient); err != nil {
			zap.L().Error("delete task failed", zap.Error(err))
			return
		}
		return
	}

	ethAddress, ok := ethAddressMap[""].(common.Address)
	if !ok {
		zap.L().Info("eth address is not", zap.String("eth address", ethAddress.String()))
		if err := DeleteTask(taskId, ethClient); err != nil {
			zap.L().Error("delete task failed", zap.Error(err))
			return
		}
		return
	}
	zap.L().Info("get the eth address", zap.String("eth address", ethAddress.String()))

	ethToVoterInfo, err := utils.GetVoterInfo(ethAddress.String(), ethClient)
	if len(ethToVoterInfo.ActorIds) != 0 {
		zap.L().Info("eth address already exists", zap.String("eth address", ethAddress.String()))
		if err := DeleteTask(taskId, ethClient); err != nil {
			zap.L().Error("delete task failed", zap.Error(err))
			return
		}
		return
	}

	ethToId, err := utils.GetActorIdFromEthAddress(ethAddress.String(), ethClient)
	if err != nil {
		zap.L().Error("get actor id from eth address failed", zap.Error(err))
		return
	}

	ethToIdIdUint64, err := strconv.ParseUint(ethToId[2:], 10, 64)
	if err != nil {
		zap.L().Error("failed to convert id to uint64", zap.Error(err))
		return
	}

	ethAddressWalletBalance, err := GetWalletBalance(ethToId, lotusRpcClient)
	if err != nil {
		zap.L().Error("failed to get balance", zap.Error(err))
		return
	}
	power := models.Power{
		DeveloperPower:   big.NewInt(0),
		TokenHolderPower: ethAddressWalletBalance,
		BlockHeight:      big.NewInt(0),
	}

	voterInfo := models.VoterInfo{
		EthAddress: ethAddress,
		ActorIds:   []uint64{ethToIdIdUint64},
	}

	if err := contract.TaskCallbackContract(voterInfo, *taskId, power, ethClient); err != nil {
		zap.L().Error("failed to execute task call back contract", zap.Error(err))
		return
	}

}
