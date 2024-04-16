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
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"powervoting-server/model"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

// GetPower get power
func GetPower(address string, client model.GoEthClient) (model.Power, error) {
	num, err := rand.Int(rand.Reader, big.NewInt(1440))
	if err != nil {
		zap.L().Error("Generate random number error: ", zap.Error(err))
		return model.Power{}, err
	}
	data, err := client.OracleAbi.Pack("getPower", common.HexToAddress(address), big.NewInt(1))
	if err != nil {
		zap.L().Error("Pack method and param error: ", zap.Error(err))
		return model.Power{}, err
	}
	zap.L().Info(fmt.Sprintf("Get power random number: %d\n", num))

	msg := ethereum.CallMsg{
		To:   &client.OracleContract,
		Data: data,
	}
	result, err := client.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		zap.L().Error("Call contract error: ", zap.Error(err))
		return model.Power{}, err
	}
	unpack, err := client.OracleAbi.Unpack("getPower", result)
	if err != nil {
		zap.L().Error("Unpack return data to interface error: ", zap.Error(err))
		return model.Power{}, err
	}
	powerInterface := unpack[0]
	marshal, err := json.Marshal(powerInterface)
	if err != nil {
		zap.L().Error("marshal error: ", zap.Error(err))
		return model.Power{}, err
	}
	var contractPower model.ContractPower
	err = json.Unmarshal(marshal, &contractPower)
	if err != nil {
		zap.L().Error("unmarshal error: ", zap.Error(err))
		return model.Power{}, err
	}
	var power model.Power
	power.TokenHolderPower = contractPower.TokenHolderPower
	power.DeveloperPower = contractPower.DeveloperPower
	power.BlockHeight = contractPower.BlockHeight
	totalClientPower := new(big.Int)
	for _, clientPower := range contractPower.ClientPower {
		power := new(big.Int)
		power.SetBytes(clientPower)
		totalClientPower.Add(totalClientPower, power)
	}
	power.ClientPower = totalClientPower
	totalSpPower := new(big.Int)
	for _, spPower := range contractPower.SpPower {
		power := new(big.Int)
		power.SetBytes(spPower)
		totalSpPower.Add(totalSpPower, power)
	}
	power.SpPower = totalSpPower
	return power, nil
}

func GetTimestamp(client model.GoEthClient) (int64, error) {
	ctx := context.Background()
	number, err := client.Client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	block, err := client.Client.BlockByNumber(ctx, big.NewInt(int64(number)))
	if err != nil {
		return 0, err
	}
	now := int64(block.Time())
	return now, nil
}

// GetVote Get vote info
func GetVote(client model.GoEthClient, proposalId int64, voteId int64) (model.ContractVote, error) {
	data, err := client.PowerVotingAbi.Pack("proposalToVote", big.NewInt(proposalId), big.NewInt(voteId))
	if err != nil {
		zap.L().Error("Pack method and param error: ", zap.Error(err))
		return model.ContractVote{}, err
	}
	msg := ethereum.CallMsg{
		To:   &client.PowerVotingContract,
		Data: data,
	}
	result, err := client.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		zap.L().Error("Call contract error: ", zap.Error(err))
		return model.ContractVote{}, err
	}
	var voteInfo model.ContractVote
	err = client.PowerVotingAbi.UnpackIntoInterface(&voteInfo, "proposalToVote", result)
	if err != nil {
		zap.L().Error("Unpack return data to interface error: ", zap.Error(err))
		return model.ContractVote{}, err
	}
	return voteInfo, nil
}

// GetProposal Get proposal
func GetProposal(client model.GoEthClient, proposalId int64) (model.ContractProposal, error) {
	data, err := client.PowerVotingAbi.Pack("idToProposal", big.NewInt(proposalId))
	if err != nil {
		zap.L().Error("Pack method and param error: ", zap.Error(err))
		return model.ContractProposal{}, err
	}
	msg := ethereum.CallMsg{
		To:   &client.PowerVotingContract,
		Data: data,
	}
	result, err := client.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		zap.L().Error("Call contract error: ", zap.Error(err))
		return model.ContractProposal{}, err
	}
	var proposal model.ContractProposal
	err = client.PowerVotingAbi.UnpackIntoInterface(&proposal, "idToProposal", result)
	if err != nil {
		zap.L().Error("Unpack return data to interface error: ", zap.Error(err))
		return model.ContractProposal{}, err
	}
	return proposal, nil
}

func GetProposalLatestId(client model.GoEthClient) (int, error) {
	data, err := client.PowerVotingAbi.Pack("proposalId")
	if err != nil {
		zap.L().Error("Pack method and param error: ", zap.Error(err))
		return 0, err
	}
	msg := ethereum.CallMsg{
		To:   &client.PowerVotingContract,
		Data: data,
	}
	result, err := client.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		zap.L().Error("Call contract error: ", zap.Error(err))
		return 0, err
	}
	unpack, err := client.PowerVotingAbi.Unpack("proposalId", result)
	if err != nil {
		zap.L().Error("unpack proposal id error: ", zap.Error(err))
		return 0, err
	}
	idStr := unpack[0]
	marshal, err := json.Marshal(idStr)
	if err != nil {
		zap.L().Error("json marshal error: ", zap.Error(err))
		return 0, err
	}
	id, err := strconv.Atoi(string(marshal))
	if err != nil {
		zap.L().Error("string to int error: ", zap.Error(err))
		return 0, err
	}
	return id, nil
}
