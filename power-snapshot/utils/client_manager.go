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
	"errors"
	"os"
	models "power-snapshot/internal/model"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type GoEthClientManager struct {
	lock        sync.Mutex
	instanceMap map[int64]models.GoEthClient
}

func NewGoEthClientManager(network []models.Network) (*GoEthClientManager, error) {
	instMap := make(map[int64]models.GoEthClient)
	for _, net := range network {
		client, err := getGoEthClient(net)
		if err != nil {
			zap.L().Error("init eth client failed", zap.Error(err))
			return nil, err
		}
		instMap[net.Id] = client
	}

	return &GoEthClientManager{
		lock:        sync.Mutex{},
		instanceMap: instMap,
	}, nil
}

func (g *GoEthClientManager) GetClient(netId int64) (models.GoEthClient, error) {
	g.lock.Lock()
	defer g.lock.Unlock()
	client, ok := g.instanceMap[netId]
	if !ok {
		zap.L().Error("client not exist", zap.Int64("netId", netId))
		return models.GoEthClient{}, errors.New("client not exist")
	}

	return client, nil
}

func (g *GoEthClientManager) ListClientNetWork() []int64 {
	g.lock.Lock()
	defer g.lock.Unlock()
	result := make([]int64, 0, len(g.instanceMap))
	for netId := range g.instanceMap {
		result = append(result, netId)
	}

	return result
}

// getGoEthClient initializes a Go-ethereum client with the provided configuration.
func getGoEthClient(network models.Network) (models.GoEthClient, error) {
	contractClient, err := ethclient.Dial(network.ContractRpc)
	if err != nil {
		zap.L().Error("ethClient.Dial error: ",
			zap.String("network", network.Name),
			zap.String("rpc", network.ContractRpc), zap.Error(err))
		return models.GoEthClient{}, err
	}

	rpcs := make([]*ethclient.Client, 0)
	for _, rpc := range network.QueryRpc {
		queryClient, err := ethclient.Dial(rpc)
		if err != nil {
			zap.L().Error("ethClient.Dial error: ",
				zap.String("network", network.Name),
				zap.String("rpc", rpc), zap.Error(err))
			return models.GoEthClient{}, err
		}
		rpcs = append(rpcs, queryClient)
	}

	// contract address, wallet private key , wallet address
	oracleContract := common.HexToAddress(network.OracleContract)

	// open abi file and parse json
	oracleFile, err := os.Open(network.OracleAbi)
	if err != nil {
		zap.L().Error("open abi file error: ", zap.Error(err))
		return models.GoEthClient{}, err
	}
	oracleAbi, err := abi.JSON(oracleFile)
	if err != nil {
		zap.L().Error("abi.JSON error: ", zap.Error(err))
		return models.GoEthClient{}, err
	}

	// generate goEthClient
	goEthClient := models.GoEthClient{
		Id:                network.Id,
		Name:              network.Name,
		IdPrefix:          network.IdPrefix,
		QueryClient:       rpcs,
		QueryRpc:          network.QueryRpc,
		ContractClient:    contractClient,
		OracleStartHeight: network.OracleStartHeight,
		ContractRpc:       network.ContractRpc,
		OracleAbi:         oracleAbi,
		OracleContract:    oracleContract,
	}
	return goEthClient, nil
}
