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
	"backend/config"
	"backend/models"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"
	"math/big"
	"os"
	"sync"
)

var (
	lock        sync.Mutex
	instanceMap map[int64]models.GoEthClient

	LotusRpcClient jsonrpc.RPCClient
	GoEthClient    models.GoEthClient
)

func init() {
	instanceMap = make(map[int64]models.GoEthClient)
}

// GetClient get go eth client
func GetClient(id int64) (models.GoEthClient, error) {
	lock.Lock()
	defer lock.Unlock()

	client, ok := instanceMap[id]
	if ok {
		return client, nil
	}

	networkList := config.Client.Network
	var clientConfig models.ClientConfig
	for _, network := range networkList {
		if network.Id == id {
			clientConfig = models.ClientConfig{
				Id:              network.Id,
				Name:            network.Name,
				Rpc:             network.Rpc,
				ContractAddress: network.ContractAddress,
				PrivateKey:      network.PrivateKey,
				WalletAddress:   network.WalletAddress,
				AbiPath:         network.AbiPath,
				GasLimit:        network.GasLimit,
			}
			break
		}
	}

	ethClient, err := getGoEthClient(clientConfig)
	if err != nil {
		return ethClient, err
	}

	instanceMap[id] = ethClient
	zap.L().Info("network init net id", zap.Int64("id", id))

	return ethClient, nil
}

// getGoEthClient get go-ethereum client
func getGoEthClient(clientConfig models.ClientConfig) (models.GoEthClient, error) {
	client, err := ethclient.Dial(clientConfig.Rpc)
	if err != nil {
		zap.L().Error("eth client  dial error", zap.Error(err))
		return models.GoEthClient{}, err
	}

	contractAddress := common.HexToAddress(clientConfig.ContractAddress)
	privateKey, err := crypto.HexToECDSA(clientConfig.PrivateKey)
	walletAddress := common.HexToAddress(clientConfig.WalletAddress)

	open, err := os.Open(clientConfig.AbiPath)
	if err != nil {
		zap.L().Error("open abi file error", zap.Error(err))
		return models.GoEthClient{}, err
	}
	contractAbi, err := abi.JSON(open)
	if err != nil {
		zap.L().Error("abi json error", zap.Error(err))
		return models.GoEthClient{}, err
	}
	gasLimit := uint64(clientConfig.GasLimit)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		zap.L().Error("client network id error", zap.Error(err))
		return models.GoEthClient{}, err
	}

	goEthClient := models.GoEthClient{
		Id:              clientConfig.Id,
		Name:            clientConfig.Name,
		Rpc:             clientConfig.Rpc,
		Client:          client,
		Abi:             contractAbi,
		GasLimit:        gasLimit,
		ChainID:         chainID,
		ContractAddress: contractAddress,
		PrivateKey:      privateKey,
		WalletAddress:   walletAddress,
	}
	return goEthClient, nil
}

// TaskCallbackContract contract task call back
func TaskCallbackContract(voterInfo models.VoterInfo, taskId big.Int, power models.Power, ethClient models.GoEthClient) error {
	lock.Lock()
	defer lock.Unlock()

	var data []byte
	var err error
	data, err = ethClient.Abi.Pack("taskCallback", voterInfo, &taskId, power)
	if err != nil {
		zap.L().Error("contract abi pack error", zap.Error(err))
		return err
	}

	nonce, err := ethClient.Client.PendingNonceAt(context.Background(), ethClient.WalletAddress)
	if err != nil {
		zap.L().Error("pending nonce at abi  pack error", zap.Error(err))
		return err
	}

	gasPrice, err := ethClient.Client.SuggestGasPrice(context.Background())
	if err != nil {
		zap.L().Error("client suggest gas price  error", zap.Error(err))
		return err
	}
	zap.L().Info("nonce", zap.Uint64("nonce", nonce))
	zap.L().Info("gas price", zap.String("gas price", gasPrice.String()))

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   ethClient.ChainID,
		Nonce:     nonce,
		GasTipCap: gasPrice,
		GasFeeCap: gasPrice,
		Gas:       ethClient.GasLimit,
		To:        &ethClient.ContractAddress,
		Value:     ethClient.Amount,
		Data:      data,
	})

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(ethClient.ChainID), ethClient.PrivateKey)
	if err != nil {
		zap.L().Error("types  sign tx  error", zap.Error(err))
		return err
	}

	err = ethClient.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		zap.L().Error("client  send transaction  error", zap.Error(err))
		return err
	}

	zap.L().Info(fmt.Sprintf("net work: %s, transaction id: %s", ethClient.Name, signedTx.Hash().Hex()))

	return nil
}

// SavePower contract up data power
func SavePower(address common.Address, power models.Power, ethClient models.GoEthClient) error {
	lock.Lock()
	defer lock.Unlock()

	var data []byte
	var err error
	data, err = ethClient.Abi.Pack("savePower", address, power)
	if err != nil {
		zap.L().Error("contract abi pack error", zap.Error(err))
		return err
	}

	nonce, err := ethClient.Client.PendingNonceAt(context.Background(), ethClient.WalletAddress)
	if err != nil {
		zap.L().Error("pending nonce at abi  pack error", zap.Error(err))
		return err
	}

	gasPrice, err := ethClient.Client.SuggestGasPrice(context.Background())
	if err != nil {
		zap.L().Error("client suggest gas price  error", zap.Error(err))
		return err
	}
	zap.L().Info("nonce", zap.Uint64("nonce", nonce))
	zap.L().Info("gas price", zap.String("gas price", gasPrice.String()))

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   ethClient.ChainID,
		Nonce:     nonce,
		GasTipCap: gasPrice,
		GasFeeCap: gasPrice,
		Gas:       ethClient.GasLimit,
		To:        &ethClient.ContractAddress,
		Value:     ethClient.Amount,
		Data:      data,
	})

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(ethClient.ChainID), ethClient.PrivateKey)
	if err != nil {
		zap.L().Error("types  sign tx  error", zap.Error(err))
		return err
	}

	err = ethClient.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		zap.L().Error("client  send transaction  error", zap.Error(err))
		return err
	}

	zap.L().Info(fmt.Sprintf("net work: %s, transaction id: %s", ethClient.Name, signedTx.Hash().Hex()))

	return nil
}

// RemoveVoter remove voter
func RemoveVoter(address common.Address, taskId big.Int, ethClient models.GoEthClient) error {
	lock.Lock()
	defer lock.Unlock()

	var data []byte
	var err error
	data, err = ethClient.Abi.Pack("removeVoter", address, &taskId)
	if err != nil {
		zap.L().Error("contract abi pack error", zap.Error(err))
		return err
	}

	nonce, err := ethClient.Client.PendingNonceAt(context.Background(), ethClient.WalletAddress)
	if err != nil {
		zap.L().Error("pending nonce at abi  pack error", zap.Error(err))
		return err
	}

	gasPrice, err := ethClient.Client.SuggestGasPrice(context.Background())
	if err != nil {
		zap.L().Error("client suggest gas price  error", zap.Error(err))
		return err
	}
	zap.L().Info("nonce", zap.Uint64("nonce", nonce))
	zap.L().Info("gas price", zap.String("gas price", gasPrice.String()))

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   ethClient.ChainID,
		Nonce:     nonce,
		GasTipCap: gasPrice,
		GasFeeCap: gasPrice,
		Gas:       ethClient.GasLimit,
		To:        &ethClient.ContractAddress,
		Value:     ethClient.Amount,
		Data:      data,
	})

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(ethClient.ChainID), ethClient.PrivateKey)
	if err != nil {
		zap.L().Error("types  sign tx  error", zap.Error(err))
		return err
	}

	err = ethClient.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		zap.L().Error("client  send transaction  error", zap.Error(err))
		return err
	}

	zap.L().Info(fmt.Sprintf("net work: %s, transaction id: %s", ethClient.Name, signedTx.Hash().Hex()))

	return nil
}
