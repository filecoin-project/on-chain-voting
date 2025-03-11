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

// GetProposalCreatorByOracle retrieves the voter information for a given address using an oracle contract.
func GetProposalCreatorByOracle(address string, client *model.GoEthClient) (model.VoterInfo, error) {
	// Pack the "getVoterInfo" method and the address parameter into a byte array using the oracle ABI.
	data, err := client.ABI.OracleAbi.Pack("getVoterInfo", common.HexToAddress(address))
	if err != nil {
		// Log an error if packing the method and parameter fails.
		zap.L().Error("Pack method and param error: ", zap.Error(err))
		return model.VoterInfo{}, err
	}

	// Call the contract using the packed data.
	result, err := client.Client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &client.OracleContract, // The address of the oracle contract.
		Data: data,                   // The packed data containing the method and parameters.
	}, nil)
	if err != nil {
		return model.VoterInfo{}, err
	}

	// Define a variable to hold the unpacked voter information.
	var voterInfo model.VoterInfo
	// Unpack the result into the voterInfo variable using the "getVoterInfo" method.
	err = client.ABI.OracleAbi.UnpackIntoInterface(&voterInfo, "getVoterInfo", result)
	if err != nil {
		zap.L().Error("Unpack return data to interface error: ", zap.Error(err))
		return model.VoterInfo{}, err
	}

	return voterInfo, nil
}
