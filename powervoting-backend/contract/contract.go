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

package contract

import (
	"os"
	"powervoting-server/config"
	"powervoting-server/model"
	"strconv"
	"sync"

	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	lock        sync.Mutex
	instanceMap map[int64]model.GoEthClient
)

func init() {
	instanceMap = make(map[int64]model.GoEthClient)
}

// GetClient retrieves a GoEthClient instance associated with the specified ID.
// it initializes a new client instance with configuration from the network list and returns it.
func GetClient(id int64) (model.GoEthClient, error) {
	lock.Lock()
	defer lock.Unlock()
	client, ok := instanceMap[id]
	if ok {
		return client, nil
	}
	networkList := config.Client.Network
	var clientConfig model.ClientConfig
	for _, network := range networkList {
		if network.Id == id {
			clientConfig = model.ClientConfig{
				Id:                  network.Id,
				Name:                network.Name,
				Rpc:                 network.Rpc,
				PowerVotingContract: network.PowerVotingContract,
				OracleContract:      network.OracleContract,
				PowerVotingAbi:      network.PowerVotingAbi,
				OracleAbi:           network.OracleAbi,
				GasLimit:            network.GasLimit,
				PrivateKey:          network.PrivateKey,
				WalletAddress:       network.WalletAddress,
			}
			break
		}
	}
	ethClient, err := getGoEthClient(clientConfig)
	if err != nil {
		return ethClient, err
	}
	instanceMap[id] = ethClient
	zap.L().Info("network init", zap.String("network id", strconv.FormatInt(id, 10)))
	return ethClient, nil
}

// getGoEthClient initializes a Go-ethereum client with the provided configuration.
func getGoEthClient(clientConfig model.ClientConfig) (model.GoEthClient, error) {
	client, err := ethclient.Dial(clientConfig.Rpc)
	if err != nil {
		zap.L().Error("ethclient.Dial error: ", zap.Error(err))
		return model.GoEthClient{}, err
	}

	// contract address, wallet private key , wallet address
	powerVotingContract := common.HexToAddress(clientConfig.PowerVotingContract)
	oracleContract := common.HexToAddress(clientConfig.OracleContract)
	privateKey, err := crypto.HexToECDSA(clientConfig.PrivateKey)

	// open abi file and parse json
	powerVotingFile, err := os.Open(clientConfig.PowerVotingAbi)
	if err != nil {
		zap.L().Error("open abi file error: ", zap.Error(err))
		return model.GoEthClient{}, err
	}
	powerVotingAbi, err := abi.JSON(powerVotingFile)
	if err != nil {
		zap.L().Error("abi.JSON error: ", zap.Error(err))
		return model.GoEthClient{}, err
	}

	// open abi file and parse json
	oracleFile, err := os.Open(clientConfig.OracleAbi)
	if err != nil {
		zap.L().Error("open abi file error: ", zap.Error(err))
		return model.GoEthClient{}, err
	}
	oracleAbi, err := abi.JSON(oracleFile)
	if err != nil {
		zap.L().Error("abi.JSON error: ", zap.Error(err))
		return model.GoEthClient{}, err
	}

	// generate goEthClient
	goEthClient := model.GoEthClient{
		Id:                  clientConfig.Id,
		Name:                clientConfig.Name,
		Client:              client,
		PowerVotingAbi:      powerVotingAbi,
		OracleAbi:           oracleAbi,
		PowerVotingContract: powerVotingContract,
		OracleContract:      oracleContract,
		PrivateKey:          privateKey,
		WalletAddress:       common.HexToAddress(clientConfig.WalletAddress),
		GasLimit:            uint64(clientConfig.GasLimit),
	}
	return goEthClient, nil
}
