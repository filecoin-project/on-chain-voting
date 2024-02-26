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
	"backend/config"
	"backend/contract"
	"backend/models"
	"backend/utils"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"math/big"
)

// TaskScheduler schedules and runs various tasks using cron jobs.
func TaskScheduler() {

	crontab := cron.New(cron.WithSeconds())

	HandleUpdatePowerTask := HandleUpdatePower

	contractCallBackTask := HandleContractCallBack

	HandleF4AddressTask := HandleF4Address

	HandleUpdatePowerSpec := "0 0 * * * ? "

	HandleContractCallBackSpec := "0/30 * * * * ? "

	HandleF4AddressSpec := "0/30 * * * * ? "

	_, err := crontab.AddFunc(HandleUpdatePowerSpec, HandleUpdatePowerTask)
	if err != nil {
		zap.L().Error("add update power count task error", zap.Error(err))
	}

	_, err = crontab.AddFunc(HandleContractCallBackSpec, contractCallBackTask)
	if err != nil {
		zap.L().Error("add contract call back count  error", zap.Error(err))
	}

	_, err = crontab.AddFunc(HandleF4AddressSpec, HandleF4AddressTask)
	if err != nil {
		zap.L().Error("add f4 address task error", zap.Error(err))
	}

	crontab.Start()

	select {}
}

// HandleContractCallBack Contract Call Back
func HandleContractCallBack() {
	for _, network := range config.Client.Network {
		ethClient, err := contract.GetClient(network.Id)
		if err != nil {
			zap.L().Error("get go-eth client error:", zap.Error(err))
			continue
		}
		lotusRpcClient := utils.NewClient(ethClient.Rpc)
		go ContractCallBack(ethClient, lotusRpcClient)
	}
}

// HandleUpdatePower Update Power
func HandleUpdatePower() {
	totalWeights := GetDeveloperWeights()
	for _, network := range config.Client.Network {
		ethClient, err := contract.GetClient(network.Id)
		if err != nil {
			zap.L().Error("get go-eth client error:", zap.Error(err))
			continue
		}
		lotusRpcClient := utils.NewClient(ethClient.Rpc)
		go UpdatePower(ethClient, lotusRpcClient, totalWeights)
	}
}

// HandleF4Address Update F4address
func HandleF4Address() {
	for _, network := range config.Client.Network {
		ethClient, err := contract.GetClient(network.Id)
		if err != nil {
			zap.L().Error("get go-eth client error:", zap.Error(err))
			continue
		}
		lotusRpcClient := utils.NewClient(ethClient.Rpc)
		go F4Address(ethClient, lotusRpcClient)
	}
}

// DeleteTask deletes a task by calling the TaskCallbackContract with zeroed power values.
func DeleteTask(taskId *big.Int, ethClient models.GoEthClient) error {
	power := models.Power{
		FipEditorPower:   big.NewInt(0),
		DeveloperPower:   big.NewInt(0),
		TokenHolderPower: big.NewInt(0),
	}
	if err := contract.TaskCallbackContract(models.VoterInfo{}, *taskId, power, ethClient); err != nil {
		return fmt.Errorf("failed to execute TaskCallbackContract: %w", err)
	}
	return nil
}
