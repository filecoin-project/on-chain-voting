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
	"backend/models"
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"math/big"
	"strconv"
)

// GetTasks retrieves the list of tasks from the Ethereum smart contract.
func GetTasks(ethClient models.GoEthClient) ([]*big.Int, error) {
	data, err := ethClient.Abi.Pack("getTasks")
	if err != nil {
		zap.L().Error("pack method and param error", zap.Error(err))
		return nil, err
	}
	msg := ethereum.CallMsg{
		To:   &ethClient.ContractAddress,
		Data: data,
	}
	result, err := ethClient.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		zap.L().Error("call contract error", zap.Error(err))
		return nil, err
	}

	unpack, err := ethClient.Abi.Unpack("getTasks", result)
	if err != nil {
		zap.L().Error("unpack return data to interface error", zap.Error(err))
		return nil, err
	}
	return unpack[0].([]*big.Int), nil
}

// GetF4Tasks retrieves the list of F4 tasks from the Ethereum smart contract.
func GetF4Tasks(ethClient models.GoEthClient) ([]*big.Int, error) {
	data, err := ethClient.Abi.Pack("getF4Tasks")
	if err != nil {
		zap.L().Error("pack method and param error", zap.Error(err))
		return nil, err
	}
	msg := ethereum.CallMsg{
		To:   &ethClient.ContractAddress,
		Data: data,
	}
	result, err := ethClient.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		zap.L().Error("call contract error", zap.Error(err))
		return nil, err
	}

	unpack, err := ethClient.Abi.Unpack("getF4Tasks", result)
	if err != nil {
		zap.L().Error("unpack return data to interface error", zap.Error(err))
		return nil, err
	}
	return unpack[0].([]*big.Int), nil
}

// GetVoterInfo retrieves the information of a voter from the Ethereum smart contract.
func GetVoterInfo(address string, ethClient models.GoEthClient) (models.VoterInfo, error) {
	data, err := ethClient.Abi.Pack("getVoterInfo", common.HexToAddress(address))
	if err != nil {
		zap.L().Error("pack method and param error", zap.Error(err))
		return models.VoterInfo{}, err
	}
	msg := ethereum.CallMsg{
		To:   &ethClient.ContractAddress,
		Data: data,
	}
	result, err := ethClient.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		zap.L().Error("call contract error", zap.Error(err))
		return models.VoterInfo{}, err
	}

	unpack, err := ethClient.Abi.Unpack("getVoterInfo", result)
	if err != nil {
		zap.L().Error("unpack return data to interface error", zap.Error(err))
		return models.VoterInfo{}, err
	}

	voterInfoInterface := unpack[0]
	marshal, err := json.Marshal(voterInfoInterface)
	if err != nil {
		zap.L().Error("marshal error", zap.Error(err))
		return models.VoterInfo{}, err
	}

	var voterInfo models.VoterInfo
	err = json.Unmarshal(marshal, &voterInfo)
	if err != nil {
		zap.L().Error("unmarshal error", zap.Error(err))
		return models.VoterInfo{}, err
	}

	return voterInfo, nil
}

// GetActorIdFromEthAddress retrieves the actor ID associated with the given Ethereum address.
func GetActorIdFromEthAddress(address string, ethClient models.GoEthClient) (string, error) {
	data, err := ethClient.Abi.Pack("resolveEthAddress", common.HexToAddress(address))
	if err != nil {
		zap.L().Error("pack method and param error", zap.Error(err))
		return "", nil
	}
	msg := ethereum.CallMsg{
		To:   &ethClient.ContractAddress,
		Data: data,
	}
	result, err := ethClient.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		zap.L().Error("call contract error", zap.Error(err))
		return "", nil
	}

	unpack, err := ethClient.Abi.Unpack("resolveEthAddress", result)
	if err != nil {
		zap.L().Error("unpack return data to interface error", zap.Error(err))
		return "", nil
	}

	var ethId string
	if ethClient.Name == "Calibration" {
		ethId = "t0" + strconv.FormatUint(unpack[0].(uint64), 10)

	} else {
		ethId = "f0" + strconv.FormatUint(unpack[0].(uint64), 10)

	}

	return ethId, nil
}
