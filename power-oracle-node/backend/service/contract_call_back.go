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
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"
	"math/big"
	"strconv"
)

// ContractCallBack is a function to handle contract callbacks.
// It retrieves the task list and starts goroutines to concurrently process tasks.
func ContractCallBack(ethClient models.GoEthClient, lotusRpcClient jsonrpc.RPCClient) {
	taskList, err := utils.GetTasks(ethClient)
	if err != nil {
		zap.L().Error("failed to get task id", zap.Error(err))
		return
	}

	if len(taskList) == 0 {
		zap.L().Info("task list is empty")
		return
	}

	for _, taskId := range taskList {
		go ProcessingContractCallBackTaskId(taskId, ethClient, lotusRpcClient)
	}

}

// ProcessingContractCallBackTaskId processes a specific task identified by its ID.
func ProcessingContractCallBackTaskId(taskId *big.Int, ethClient models.GoEthClient, lotusRpcClient jsonrpc.RPCClient) {
	zap.L().Info("processing task", zap.String("id", taskId.String()))

	args := []interface{}{taskId}
	ucanCid, err := GetContractMapping(models.TaskIdToUcanCid, ethClient, args)
	if err != nil {
		zap.L().Error("failed to get ucan cid", zap.Error(err))
		if err := DeleteTask(taskId, ethClient); err != nil {
			zap.L().Error("delete task failed", zap.Error(err))
			return
		}
		return
	}

	ucanCidStr, ok := ucanCid[""].(string)
	if !ok {
		zap.L().Error("ucan cid is not a string")

		if err := DeleteTask(taskId, ethClient); err != nil {
			zap.L().Error("delete task failed", zap.Error(err))
			return
		}
		return
	}
	zap.L().Info("get ucan cid", zap.String("ucan cid", ucanCidStr))

	ucan, err := GetUcanFromIpfs(ucanCidStr)
	if err != nil {
		zap.L().Error("failed to get ucan string from ipfs", zap.Error(err))
		if err := DeleteTask(taskId, ethClient); err != nil {
			zap.L().Error("delete task failed", zap.Error(err))
			return
		}
		return
	}
	zap.L().Info("get ucan string", zap.String("ucan string", ucan))

	iss, aud, act, isGitHub, err := VerifyUCAN(ucan, lotusRpcClient)
	if err != nil {
		zap.L().Error("failed to verify ucan", zap.Error(err))
		if err := DeleteTask(taskId, ethClient); err != nil {
			zap.L().Error("delete task failed", zap.Error(err))
			return
		}
		return
	}
	zap.L().Info("ucan verification successful,permission is", zap.String("ucan string", ucan))
	zap.L().Info("get ucan string", zap.String("action", act))

	ethToVoterInfo, err := utils.GetVoterInfo(iss, ethClient)
	if ethToVoterInfo.UcanCid == ucanCidStr {
		zap.L().Info("ucan cid already exists", zap.String("eth address", iss))
		if err := DeleteTask(taskId, ethClient); err != nil {
			zap.L().Error("delete task failed", zap.Error(err))
			return
		}
		return
	}

	switch act {

	case "add":
		if err := ActionIsAdd(iss, aud, ucanCidStr, isGitHub, taskId, ethClient, lotusRpcClient); err != nil {
			zap.L().Error("failed to perform callback add", zap.Error(err))
			if err := DeleteTask(taskId, ethClient); err != nil {
				zap.L().Error("delete task failed", zap.Error(err))
				return
			}
			return
		}
		zap.L().Info("task call back add successful")

	case "del":
		if err := contract.RemoveVoter(common.HexToAddress(iss), *taskId, ethClient); err != nil {
			zap.L().Error("task call back executing delete failed", zap.Error(err))
			if err := DeleteTask(taskId, ethClient); err != nil {
				zap.L().Error("delete task failed", zap.Error(err))
				return
			}
			return
		}
		zap.L().Info("task callback del successful")
	}
	return
}

// ActionIsAdd performs the 'add' action for a given UCAN.
func ActionIsAdd(iss, aud, ucanCid string, isGitHub bool, taskId *big.Int, ethClient models.GoEthClient, lotusRpcClient jsonrpc.RPCClient) error {
	voterInfo, err := utils.GetVoterInfo(iss, ethClient)
	if err != nil {
		return err
	}

	if isGitHub {
		if len(voterInfo.ActorIds) < 2 {
			return GitHubAccount(iss, aud, ucanCid, taskId, ethClient, lotusRpcClient)
		}

	} else {
		if voterInfo.GithubAccount == "" {
			return FileCoinAccount(iss, aud, ucanCid, taskId, ethClient, lotusRpcClient)
		}
	}
	return errors.New("already exists")

}

// GitHubAccount handles the 'add' action for a GitHub account.
func GitHubAccount(iss, aud, ucanCid string, taskId *big.Int, ethClient models.GoEthClient, lotusRpcClient jsonrpc.RPCClient) error {
	ethToId, err := utils.GetActorIdFromEthAddress(iss, ethClient)
	if err != nil {
		zap.L().Error("failed to convert ethereum address to filecoin address", zap.Error(err))
		return err

	}

	ethToIdIdUint64, err := strconv.ParseUint(ethToId[2:], 10, 64)
	if err != nil {
		zap.L().Error("failed to convert id to uint64", zap.Error(err))
		return err
	}

	voterInfo := models.VoterInfo{
		ActorIds:      []uint64{ethToIdIdUint64},
		UcanCid:       ucanCid,
		EthAddress:    common.HexToAddress(iss),
		GithubAccount: aud,
	}

	totalWeights := GetDeveloperWeights()
	weight := totalWeights[aud]
	ethAddressWalletBalance, err := GetWalletBalance(ethToId, lotusRpcClient)
	if err != nil {
		zap.L().Error("failed to get balance", zap.Error(err))
		return err
	}

	power := models.Power{
		DeveloperPower:   big.NewInt(int64(weight)),
		TokenHolderPower: ethAddressWalletBalance,
		BlockHeight:      big.NewInt(0),
	}

	if err := contract.TaskCallbackContract(voterInfo, *taskId, power, ethClient); err != nil {
		zap.L().Error("failed to execute task call back contract", zap.Error(err))
		return err
	}

	return nil
}

// FileCoinAccount handles the 'add' action for a Filecoin account.
func FileCoinAccount(iss, aud, ucanCid string, taskId *big.Int, ethClient models.GoEthClient, lotusRpcClient jsonrpc.RPCClient) error {
	filecoinId, err := utils.IDFromAddress(context.Background(), lotusRpcClient, aud)
	if err != nil {
		zap.L().Error("failed to get id from address", zap.Error(err))
		return err
	}

	filecoinIdUint64, err := strconv.ParseUint(filecoinId[2:], 10, 64)
	if err != nil {
		zap.L().Error("failed to convert id to uint64", zap.Error(err))
		return err
	}

	ethToId, err := utils.GetActorIdFromEthAddress(iss, ethClient)
	if err != nil {
		zap.L().Error("failed to convert ethereum address to filecoin address", zap.Error(err))
		return err

	}

	ethToIdIdUint64, err := strconv.ParseUint(ethToId[2:], 10, 64)
	if err != nil {
		zap.L().Error("failed to convert id to uint64", zap.Error(err))
		return err
	}

	voterInfo := models.VoterInfo{
		ActorIds:   []uint64{filecoinIdUint64, ethToIdIdUint64},
		UcanCid:    ucanCid,
		EthAddress: common.HexToAddress(iss),
	}

	power := models.Power{
		DeveloperPower:   big.NewInt(0),
		TokenHolderPower: big.NewInt(0),
		BlockHeight:      big.NewInt(0),
	}

	filecoinAddressWalletBalance, err := GetWalletBalance(filecoinId, lotusRpcClient)
	if err != nil {
		zap.L().Error("failed to get balance", zap.Error(err))

		return err
	}

	ethAddressWalletBalance, err := GetWalletBalance(ethToId, lotusRpcClient)
	if err != nil {
		zap.L().Error("failed to get balance", zap.Error(err))
		return err
	}

	power.TokenHolderPower.Add(ethAddressWalletBalance, filecoinAddressWalletBalance)

	if err := contract.TaskCallbackContract(voterInfo, *taskId, power, ethClient); err != nil {
		zap.L().Error("failed to execute task call back contract", zap.Error(err))
		return err
	}

	return nil
}
