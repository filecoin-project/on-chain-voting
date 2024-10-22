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
	"fmt"
	"math/big"
	"powervoting-server/model"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

var (
	lock sync.Mutex
)

// GetTimestamp retrieves the current timestamp from the Ethereum blockchain using the provided Ethereum client.
// It fetches the latest block number and then retrieves the block information to obtain the timestamp.
// The function returns the current timestamp as an int64 or an error if the operation fails.
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

// GetVote retrieves information about a specific vote associated with a proposal from the PowerVoting smart contract.
// It packs the necessary parameters and calls the PowerVoting smart contract to fetch the vote details.
// The function returns the contract vote information or an error if the operation fails.
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

// GetProposal retrieves information about a specific proposal from the PowerVoting smart contract.
// It packs the necessary parameters and calls the PowerVoting smart contract to fetch the proposal details.
// The function returns the contract proposal information or an error if the operation fails.
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

// GetProposalLatestId retrieves the latest proposal ID from the PowerVoting smart contract.
// It packs the necessary method and parameters and calls the PowerVoting smart contract to fetch the latest proposal ID.
// The function returns the latest proposal ID or an error if the operation fails.
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

// GetVoterToPowerStatus retrieves the power status of a voter from the Oracle contract.
// It packs the necessary method and parameters and calls the Oracle contract to fetch the power status.
// The function returns the power status of the voter or an error if the operation fails.
func GetVoterToPowerStatus(address string, client model.GoEthClient) (model.VoterToPowerStatus, error) {
	data, err := client.OracleAbi.Pack("voterToPowerStatus", common.HexToAddress(address))
	if err != nil {
		zap.L().Error("Pack method and param error: ", zap.Error(err))
		return model.VoterToPowerStatus{}, err
	}

	msg := ethereum.CallMsg{
		To:   &client.OracleContract,
		Data: data,
	}
	result, err := client.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		zap.L().Error("Call contract error: ", zap.Error(err))
		return model.VoterToPowerStatus{}, err
	}
	var voterToPowerStatus model.VoterToPowerStatus
	err = client.OracleAbi.UnpackIntoInterface(&voterToPowerStatus, "voterToPowerStatus", result)
	if err != nil {
		zap.L().Error("Unpack return data to interface error: ", zap.Error(err))
		return model.VoterToPowerStatus{}, err
	}
	return voterToPowerStatus, nil
}

func AddSnapshot(client model.GoEthClient, cid, day string) error {
	lock.Lock()
	defer lock.Unlock()

	data, err := client.OracleAbi.Pack("addSnapshot", day, cid)
	if err != nil {
		zap.L().Error("AddSnapshot Pack method and param error: ", zap.Error(err))
		return err
	}

	nonce, err := client.Client.PendingNonceAt(context.Background(), client.WalletAddress)
	if err != nil {
		zap.L().Error("addSnapshot pending nonce at abi  pack error", zap.Error(err))
		return err
	}

	gasPrice, err := client.Client.SuggestGasPrice(context.Background())
	if err != nil {
		zap.L().Error("addSnapshot client suggest gas price  error", zap.Error(err))
		return err
	}
	zap.L().Info("addSnapshot nonce", zap.Uint64("nonce", nonce))
	zap.L().Info("addSnapshot gas price", zap.String("gas price", gasPrice.String()))

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      client.GasLimit,
		To:       &client.OracleContract,
		Value:    client.Amount,
		Data:     data,
	})

	zap.L().Info("AddSnapshot data: ", zap.String("data", string(day)))

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(client.ChainID), client.PrivateKey)
	if err != nil {
		zap.L().Error("addSnapshot types  sign tx  error", zap.Error(err))
		return err
	}

	err = client.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		zap.L().Error("addSnapshot client send transaction  error", zap.Error(err))
		return err
	}

	zap.S().Info(fmt.Sprintf("addSnapshot Successfully, network: %s, transaction id: %s", client.Name, signedTx.Hash().Hex()))

	return nil
}
