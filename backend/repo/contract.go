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

package repo

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"powervoting-server/model"
)


type IContractRepo interface {
	GetBlockTime(syncedHeight int64) (int64, error)
	GetOwnerIdByOracle(actorID uint64, minerId []uint64) []uint64
	GetActorIdByAddress(address string) (uint64, error)
	GetVotedAlgorithm() (string, error)
}
type ContractRepo struct {
	client *model.GoEthClient
}

func NewContractRepo(client *model.GoEthClient) *ContractRepo {
	return &ContractRepo{
		client: client,
	}
}

// GetBlockTime retrieves the timestamp of a specific block by its height
// Parameters:
//   - client: Ethereum client instance for blockchain access
//   - syncedHeight: Block number to query
//
// Returns:
//   - int64: Unix timestamp of the block
//   - error: Any error that occurred during the query
func (c *ContractRepo) GetBlockTime(syncedHeight int64) (int64, error) {
	// Create background context for the block query
	ctx := context.Background()

	// Get block information by block number
	block, err := c.client.Client.BlockByNumber(ctx, big.NewInt(syncedHeight))
	if err != nil {
		// Return 0 and error if block query fails
		return 0, err
	}

	// Return the block's timestamp (converted to int64)
	return int64(block.Time()), nil
}

func (c *ContractRepo) GetOwnerIdByOracle(actorID uint64, minerId []uint64) []uint64 {
	var res []uint64
	for _, miner := range minerId {
		callData, err := c.client.ABI.OraclePowersAbi.Pack("getOwner", miner)
		if err != nil {
			zap.L().Error("Error packing getOwner", zap.Uint64("minerId", miner))
			continue
		}

		upPacked, err := c.callContractMethod(c.client.OraclePowersContract, c.client.ABI.OraclePowersAbi, "getOwner", callData)
		if err != nil {
			zap.L().Error("Error calling getOwner", zap.Uint64("minerId", miner), zap.Error(err))
			continue
		}

		if upPacked[0].(uint64) == uint64(actorID) {
			res = append(res, uint64(miner))
		}
	}

	return res
}

func (c *ContractRepo) GetActorIdByAddress(address string) (uint64, error) {
	callData, err := c.client.ABI.OraclePowersAbi.Pack("resolveEthAddress", common.HexToAddress(address))
	if err != nil {
		return 0, fmt.Errorf("error packing resolveEthAddress: %v", err)
	}

	upPacked, err := c.callContractMethod(c.client.OraclePowersContract, c.client.ABI.OraclePowersAbi, "resolveEthAddress", callData)
	if err != nil {
		return 0, err
	}

	return upPacked[0].(uint64), err
}

func (c *ContractRepo) GetVotedAlgorithm() (string, error) {
	callData, err := c.client.ABI.PowerVotingConfAbi.Pack("votingCountingAlgorithm")
	if err != nil {
		return "", fmt.Errorf("error packing votingCountingAlgorithm: %v", err)
	}

	upPacked, err := c.callContractMethod(c.client.PowerVotingConfContract, c.client.ABI.PowerVotingConfAbi, "votingCountingAlgorithm", callData)
	if err != nil {
		return "", err
	}

	return upPacked[0].(string), nil
}

// callContractMethod calls a contract method and returns the unpacked results
// Parameters:
//   - client: Ethereum client instance
//   - contractAddr: Address of the target contract
//   - abi: Contract ABI for method encoding/decoding
//   - method: Name of the contract method to call
//   - callData: Pre-encoded method call data
//
// Returns:
//   - []any: Unpacked return values from the contract call
//   - error: Any error that occurred during execution
func (c *ContractRepo) callContractMethod(contractAddr common.Address, abi *abi.ABI, method string, callData []byte) ([]any, error) {
	// Prepare the contract call message
	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: callData,
	}

	// Execute the contract call
	result, err := c.client.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("error calling %s: %v", method, err)
	}

	// Check for empty response
	if len(result) == 0 {
		return nil, fmt.Errorf("%s returned empty result", method)
	}

	// Unpack the raw return data using ABI
	unPacked, err := abi.Unpack(method, result)
	if err != nil {
		return nil, fmt.Errorf("error unpacking %s: %v", method, err)
	}

	return unPacked, nil
}
