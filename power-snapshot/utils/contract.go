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
	"encoding/json"
	models "power-snapshot/internal/model"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

// GetVoterAddresses retrieves the addresses of voters stored in the Ethereum smart contract.
func GetVoterAddresses(ethClient models.GoEthClient) ([]common.Address, error) {
	data, err := ethClient.OracleAbi.Pack("getVoterAddresses")
	if err != nil {
		zap.L().Error("pack method and param error", zap.Error(err))
		return nil, err
	}
	msg := ethereum.CallMsg{
		To:   &ethClient.OracleContract,
		Data: data,
	}
	result, err := ethClient.ContractClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		zap.L().Error("call contract error", zap.Error(err))
		return nil, err
	}

	unpack, err := ethClient.OracleAbi.Unpack("getVoterAddresses", result)
	if err != nil {
		zap.L().Error("unpack return data to interface error", zap.Error(err))
		return nil, err
	}

	return unpack[0].([]common.Address), nil
}

// GetVoterInfo retrieves the information of a voter from the Ethereum smart contract.
func GetVoterInfo(address string, ethClient models.GoEthClient) (models.VoterInfo, error) {
	data, err := ethClient.OracleAbi.Pack("getVoterInfo", common.HexToAddress(address))
	if err != nil {
		zap.L().Error("pack method and param error", zap.Error(err))
		return models.VoterInfo{}, err
	}
	msg := ethereum.CallMsg{
		To:   &ethClient.OracleContract,
		Data: data,
	}
	result, err := ethClient.ContractClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		zap.L().Error("call contract error", zap.Error(err))
		return models.VoterInfo{}, err
	}

	unpack, err := ethClient.OracleAbi.Unpack("getVoterInfo", result)
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
