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
	"power-snapshot/utils/types"

	filecoinAddress "github.com/filecoin-project/go-address"
	"github.com/ybbus/jsonrpc/v3"
)

// NewClient create lotus rcp client.
func NewClient(endpoint string) jsonrpc.RPCClient {
	return jsonrpc.NewClientWithOpts(endpoint, &jsonrpc.RPCClientOpts{})
}

func GetTipSetByHeight(lotusRpcClient jsonrpc.RPCClient, height int64) (int64, []interface{}, error) {
	resp, err := lotusRpcClient.Call(context.Background(), "Filecoin.ChainGetTipSetByHeight", height, types.TipSetKey{})
	if err != nil {
		return 0, nil, err
	}

	resHeight, err := resp.Result.(map[string]interface{})["Height"].(json.Number).Int64()
	if err != nil {
		return 0, nil, err
	}

	resTipSet := resp.Result.(map[string]interface{})["Cids"].([]interface{})

	return resHeight, resTipSet, nil
}

// IDFromAddress retrieves the identifier associated with a Filecoin address.
func IDFromAddress(ctx context.Context, lotusRpcClient jsonrpc.RPCClient, address string) (string, error) {
	resp, err := lotusRpcClient.Call(ctx, "Filecoin.StateLookupID", address, types.TipSetKey{})
	if err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", resp.Error
	}

	return resp.Result.(string), nil
}

func GetWalletBalanceByHeight(ctx context.Context, lotusRpcClient jsonrpc.RPCClient, id string, height int64) (string, error) {
	height, tipSetKey, err := GetTipSetByHeight(lotusRpcClient, height)
	if err != nil {
		return "", err
	}

	addressStr, err := filecoinAddress.NewFromString(id)
	if err != nil {
		return "", err
	}

	tipSetList := append([]any{}, addressStr, tipSetKey)

	resp, err := lotusRpcClient.Call(ctx, "Filecoin.StateGetActor", &tipSetList)
	if err != nil {
		return "", err
	}

	tmp, err := json.Marshal(resp.Result)
	if err != nil {
		return "", err
	}
	var t types.BalanceInfo
	if err := json.Unmarshal(tmp, &t); err != nil {
		return "", err
	}

	return t.Balance, nil
}

type MinerPower struct {
	MinerPower struct {
		RawBytePower    string `json:"RawBytePower"`
		QualityAdjPower string `json:"QualityAdjPower"`
	} `json:"MinerPower"`
	TotalPower struct {
		RawBytePower    string `json:"RawBytePower"`
		QualityAdjPower string `json:"QualityAdjPower"`
	} `json:"TotalPower"`
	HasMinPower bool `json:"HasMinPower"`
}

func GetMinerPowerByHeight(ctx context.Context, lotusRpcClient jsonrpc.RPCClient, id string, height int64) (MinerPower, error) {
	_, tipSetKey, err := GetTipSetByHeight(lotusRpcClient, height)
	if err != nil {
		return MinerPower{}, err
	}

	addressStr, err := filecoinAddress.NewFromString(id)
	if err != nil {
		return MinerPower{}, err
	}

	tipSetList := append([]any{}, addressStr, tipSetKey)

	resp, err := lotusRpcClient.Call(ctx, "Filecoin.StateMinerPower", tipSetList)
	if err != nil {
		return MinerPower{}, err
	}

	tmp, err := json.Marshal(resp.Result)
	if err != nil {
		return MinerPower{}, err
	}
	var minerPower MinerPower
	if err := json.Unmarshal(tmp, &minerPower); err != nil {
		return MinerPower{}, err
	}

	return minerPower, nil
}

func GetClientBalanceByHeight(ctx context.Context, lotusRpcClient jsonrpc.RPCClient, id string, height int64) (string, error) {
	height, tipSetKey, err := GetTipSetByHeight(lotusRpcClient, height)
	if err != nil {
		return "", err
	}

	addressStr, err := filecoinAddress.NewFromString(id)
	if err != nil {
		return "", err
	}

	tipSetList := append([]any{}, addressStr, tipSetKey)

	var t types.StorageMarket
	resp, err := lotusRpcClient.Call(ctx, "Filecoin.StateMarketBalance", tipSetList)
	if err != nil {
		return "", err
	}

	tmp, err := json.Marshal(resp.Result)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(tmp, &t); err != nil {
		return "", err
	}

	return t.Locked, nil
}

func GetNewestHeight(ctx context.Context, lotusRpcClient jsonrpc.RPCClient) (height int64, err error) {
	resp, err := lotusRpcClient.Call(ctx, "Filecoin.ChainHead")
	if err != nil {
		return 0, nil
	}

	height, err = resp.Result.(map[string]interface{})["Height"].(json.Number).Int64()
	if err != nil {
		return 0, nil
	}

	return height, nil
}

type BlockHeader struct {
	Height    int64 `json:"Height"`
	Timestamp int64 `json:"Timestamp"`
}

type GetBlockParam struct {
	Empty string `json:"/"`
}

func GetBlockHeader(ctx context.Context, lotusRpcClient jsonrpc.RPCClient, height int64) (BlockHeader, error) {
	height, tipSetKey, err := GetTipSetByHeight(lotusRpcClient, height)

	marshal, err := json.Marshal(tipSetKey)
	if err != nil {
		return BlockHeader{}, err
	}

	var blockParam []GetBlockParam
	err = json.Unmarshal(marshal, &blockParam)
	if err != nil {
		return BlockHeader{}, err
	}

	resp, err := lotusRpcClient.Call(ctx, "Filecoin.ChainGetBlock", []GetBlockParam{{Empty: blockParam[0].Empty}})

	if err != nil {
		return BlockHeader{}, err
	}

	if resp.Error != nil {
		return BlockHeader{}, resp.Error
	}

	tmp, err := json.Marshal(resp.Result.(map[string]interface{}))
	if err != nil {
		return BlockHeader{}, err
	}

	var block BlockHeader
	if err := json.Unmarshal(tmp, &block); err != nil {
		return block, err
	}
	return block, nil
}
