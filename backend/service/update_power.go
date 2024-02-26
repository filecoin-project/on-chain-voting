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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"
	"math/big"
	"strconv"
)

// UpdatePower updates the power of voters based on their GitHub accounts and Filecoin balances.
// It retrieves voter addresses, processes each address asynchronously to update their power,
// and saves the updated power to the smart contract.
func UpdatePower(ethClient models.GoEthClient, lotusRpcClient jsonrpc.RPCClient, totalWeights map[string]int) {
	ethAddressList, err := utils.GetVoterAddresses(ethClient)
	if err != nil {
		zap.L().Error("failed to obtain eth address list", zap.Error(err))
		return
	}

	zap.L().Info("up date the list of votes and power", zap.Int("list", len(ethAddressList)))

	if len(ethAddressList) != 0 {
		for _, voterAddress := range ethAddressList {
			if voterAddress != common.HexToAddress("0x0000000000000000000000000000000000000000") {
				go ProcessingUpdatePowerTaskId(voterAddress, totalWeights, ethClient, lotusRpcClient)
			}
		}
	}
}

// ProcessingUpdatePowerTaskId updates the power of a specific voter address based on their GitHub account and Filecoin balance.
func ProcessingUpdatePowerTaskId(ethAddress common.Address, totalWeights map[string]int, ethClient models.GoEthClient, lotusRpcClient jsonrpc.RPCClient) {
	zap.L().Info("update power", zap.String("eth address", ethAddress.String()))

	voterToInfo, err := utils.GetVoterInfo(ethAddress.String(), ethClient)
	if err != nil {
		zap.L().Error("failed to obtain vote info", zap.Error(err))
		return
	}

	tokenHolderPower := new(big.Int)
	for _, id := range voterToInfo.ActorIds {
		var balance *big.Int
		if ethClient.Id == 314159 {
			balance, err = GetWalletBalance("t0"+strconv.FormatUint(id, 10), lotusRpcClient)
		} else {
			balance, err = GetWalletBalance("f0"+strconv.FormatUint(id, 10), lotusRpcClient)
		}
		tokenHolderPower.Add(tokenHolderPower, balance)
	}

	if voterToInfo.GithubAccount != "" {
		weight := totalWeights[voterToInfo.GithubAccount]
		developerWeight := big.NewInt(int64(weight))

		power := models.Power{
			FipEditorPower:   big.NewInt(0),
			DeveloperPower:   developerWeight,
			TokenHolderPower: tokenHolderPower,
		}

		if err := contract.SavePower(ethAddress, power, ethClient); err != nil {
			zap.L().Error("failed to execute save power", zap.Error(err))
			return
		}
		zap.L().Info("developer updated successfully", zap.String("developer", voterToInfo.GithubAccount))
		return
	}

	power := models.Power{
		FipEditorPower:   big.NewInt(0),
		DeveloperPower:   big.NewInt(0),
		TokenHolderPower: tokenHolderPower,
	}

	if err := contract.SavePower(ethAddress, power, ethClient); err != nil {
		zap.L().Error("failed to execute save power", zap.Error(err))
		return
	}
	zap.L().Info("eth address updated successfully", zap.String("address", voterToInfo.EthAddress.String()))
	return

}
