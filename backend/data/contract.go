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

package data

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	"powervoting-server/config"
	"powervoting-server/model"
	"powervoting-server/service"
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
func GetClient(syncService service.ISyncService, chainId int64) (*model.GoEthClient, error) {
	lock.Lock()
	defer lock.Unlock()
	client, ok := instanceMap[chainId]
	if ok {
		return &client, nil
	}
	network := config.Client.Network

	clientConfig := model.ClientConfig{
		ChainId:              network.ChainId,
		Name:                 network.Name,
		Rpc:                  network.Rpc,
		PowerVotingContract:  network.PowerVotingContract,
		OracleContract:       network.OracleContract,
		SyncEventStartHeight: network.SyncEventStartHeight,
		FipContract:          network.FipContract,
	}

	if err := syncService.CreateFipEditor(context.Background(), &model.FipEditorTbl{
		ChainId: chainId,
		Editor:  network.FipInitEditor,
	}); err != nil {
		zap.L().Error("CreateFipEditor error: ", zap.Error(err))
		return nil, err
	}

	if err := syncService.CreateSyncEventInfo(context.Background(), &model.SyncEventTbl{
		ChainId:                    chainId,
		ChainName:                  clientConfig.Name,
		PowerVotingContractAddress: clientConfig.PowerVotingContract,
		FipProposalContractAddress: clientConfig.FipContract,
		SyncedHeight:               clientConfig.SyncEventStartHeight,
	}); err != nil {
		zap.L().Error("CreateSyncEventInfo error: ", zap.Error(err))
		return nil, err
	}

	zap.L().Info("network init", zap.String("network id", strconv.FormatInt(chainId, 10)))

	ethClient, err := getGoEthClient(clientConfig, config.Client.ABIPath)
	if err != nil {
		return nil, err
	}
	instanceMap[chainId] = ethClient
	zap.L().Info("network init", zap.String("network id", strconv.FormatInt(chainId, 10)))

	return &ethClient, nil
}

// getGoEthClient initializes a Go-ethereum client with the provided configuration.
func getGoEthClient(clientConfig model.ClientConfig, abiPath config.ABIPath) (model.GoEthClient, error) {
	client, err := ethclient.Dial(clientConfig.Rpc)
	if err != nil {
		zap.L().Error("ethclient.Dial error: ", zap.Error(err))
		return model.GoEthClient{}, err
	}

	// contract address, wallet private key , wallet address
	powerVotingContract := common.HexToAddress(clientConfig.PowerVotingContract)
	oracleContract := common.HexToAddress(clientConfig.OracleContract)
	orcalePowersContract := common.HexToAddress(clientConfig.OraclePowersContract)
	fipContract := common.HexToAddress(clientConfig.FipContract)
	// generate goEthClient
	goEthClient := model.GoEthClient{
		ChainId:              clientConfig.ChainId,
		Name:                 clientConfig.Name,
		Client:               client,
		PowerVotingContract:  powerVotingContract,
		OracleContract:       oracleContract,
		OraclePowersContract: orcalePowersContract,
		FipContract:          fipContract,
		ABI: &model.ABI{
			PowerVotingAbi:  GetAbiFromLocalFile(abiPath.PowerVotingAbi),
			OracleAbi:       GetAbiFromLocalFile(abiPath.OracleAbi),
			FipAbi:          GetAbiFromLocalFile(abiPath.FipAbi),
		},
	}
	return goEthClient, nil
}

func GetAbiFromLocalFile(abiPath string) *abi.ABI {
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
