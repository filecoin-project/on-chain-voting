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
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"power-snapshot/config"
	"power-snapshot/constant"
)

type ContractRepoImpl struct {
	goEthClient             *ethclient.Client
	powerVotingConfAbi      *abi.ABI
	PowerVotingConfContract common.Address
}

func NewContractRepoImpl(client *ethclient.Client) *ContractRepoImpl {
	network := config.Client.Network
	if len(network.QueryRpc) == 0 {
		panic("no rpc url found")
	}

	return &ContractRepoImpl{
		goEthClient:             client,
		powerVotingConfAbi:      getAbiFromLocalFile(network.PowerVotingAbiPath),
		PowerVotingConfContract: common.HexToAddress(network.PowerVotingConfContract),
	}
}


func (c *ContractRepoImpl) SetSnapshotHeight(dateStr string, height uint64) error {
	callData, err := c.powerVotingConfAbi.Pack("setSnapshotHeight", dateStr, height)
	if err != nil {
		return fmt.Errorf("error packing setSnapshotHeight: %v", err)
	}

	_, err = c.callContractMethod(c.PowerVotingConfContract, c.powerVotingConfAbi, "setSnapshotHeight", callData)
	if err != nil {
		return err
	}

	return nil
}

func (c *ContractRepoImpl) GetSnapshotHeight(dateStr string) (int64, error) {
	callData, err := c.powerVotingConfAbi.Pack("getSnapshotHeight", dateStr)
	if err != nil {
		return 0, fmt.Errorf("error packing getSnapshotHeight: %v", err)
	}

	upPacked, err := c.callContractMethod(c.PowerVotingConfContract, c.powerVotingConfAbi, "getSnapshotHeight", callData)
	if err != nil {
		return 0, err
	}

	return int64(upPacked[0].(uint64)), nil
}

func (c *ContractRepoImpl) GetExpirationData() (int, error) {
	callDate, err := c.powerVotingConfAbi.Pack("expDays")
	if err != nil {
		return 0, fmt.Errorf("error packing expDays: %v", err)
	}

	upPacked, err := c.callContractMethod(c.PowerVotingConfContract, c.powerVotingConfAbi, "expDays", callDate)
	if err != nil {
		return 0, err
	}

	res := int(upPacked[0].(uint64))
	if res == 0 {
		return constant.DataExpiredDuration, nil
	}

	return res, nil
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
func (c *ContractRepoImpl) callContractMethod(contractAddr common.Address, abi *abi.ABI, method string, callData []byte) ([]any, error) {
	// Prepare the contract call message
	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: callData,
	}

	// Execute the contract call
	result, err := c.goEthClient.CallContract(context.Background(), msg, nil)
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

func getAbiFromLocalFile(abiPath string) *abi.ABI {
	// open abi file and parse json
	abiFile, err := os.Open(abiPath)
	if err != nil {
		panic(fmt.Errorf("failed to open abi file: %s, error: %v", abiPath, err))
	}
	abi, err := abi.JSON(abiFile)
	if err != nil {
		panic(err)
	}

	return &abi
}
