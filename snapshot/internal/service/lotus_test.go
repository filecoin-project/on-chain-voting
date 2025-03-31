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

package service

import (
	"context"
	"encoding/json"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"

	"power-snapshot/internal/mock"
	"power-snapshot/internal/repo"
)

var (
	address = "t410fpqsmux52n4pcfcssbeiuoz2g2jn6l3n6zf6tmii"
	ethAddr = "0x7C24ca5FBA6f1E228a520911476746D25Be5EdbE"
	id      = "t099523"
)

func TestWalletBalanceByHeight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lotusRpcClient := mock.NewMockRPCClient(gomock.NewController(t))
	lotusRpcClient.EXPECT().Call(gomock.Any(), "Filecoin.ChainGetTipSetByHeight", int64(2057965), gomock.Any()).Return(&jsonrpc.RPCResponse{
		Result: map[string]interface{}{
			"Cids": []interface{}{
				map[string]interface{}{
					"/": "bafy2bzacebes5mgufnyilg5e5p2mvhftqvoplzo65zxvcjy3j62n6w4er3hmm",
				},
				map[string]interface{}{
					"/": "bafy2bzacecqarqyklux426of26bwmcvd3rjjjlsqprd3umvbjtcn6hcbjoypa",
				},
			},
			"Blocks": 1,
			"Height": json.Number("2057965"),
			"Test":   1,
		},
		JSONRPC: "2.0",
		Error:   nil,
		ID:      0,
	}, nil)

	lotusRpcClient.EXPECT().Call(gomock.Any(), "Filecoin.StateGetActor", gomock.Any()).Return(&jsonrpc.RPCResponse{
		Result: map[string]interface{}{
			"Balance": "100999995981726481390",
			"Test":    2,
		},
		JSONRPC: "2.0",
		Error:   nil,
		ID:      0,
	}, nil)

	lotusRepo := repo.NewLotusRPCRepo(redis.NewClient(&redis.Options{
		ClientName: "SanpshotTest",
	}))
	lotusService := NewLotusService(lotusRepo)
	rsp, err := lotusService.GetWalletBalanceByHeight(context.Background(), "t03751", 314159, 2057965)
	assert.Nil(t, err)
	// expectedBalance := "467051509071593317720"
	// assert.Equal(t, expectedBalance, rsp)

	zap.L().Info("result", zap.Any("balance", rsp))
}

func TestGetMinerPowerByHeight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lotusRpcClient := mock.NewMockRPCClient(gomock.NewController(t))

	lotusRpcClient.EXPECT().Call(gomock.Any(), "Filecoin.ChainGetTipSetByHeight", int64(2057965), gomock.Any()).Return(&jsonrpc.RPCResponse{
		Result: map[string]interface{}{
			"Cids": []interface{}{
				map[string]interface{}{
					"/": "bafy2bzacebes5mgufnyilg5e5p2mvhftqvoplzo65zxvcjy3j62n6w4er3hmm",
				},
				map[string]interface{}{
					"/": "bafy2bzacecqarqyklux426of26bwmcvd3rjjjlsqprd3umvbjtcn6hcbjoypa",
				},
			},
			"Blocks": 1,
			"Height": json.Number("2057965"),
			"Test":   1,
		},
		JSONRPC: "2.0",
		Error:   nil,
		ID:      0,
	}, nil)

	lotusRpcClient.EXPECT().Call(gomock.Any(), "Filecoin.StateMinerPower", gomock.Any()).Return(&jsonrpc.RPCResponse{
		Result: map[string]interface{}{
			"MinerPower": map[string]interface{}{
				"RawBytePower":    "0",
				"QualityAdjPower": "0",
			},
			"TotalPower": map[string]interface{}{
				"RawBytePower":    "676818126372864",
				"QualityAdjPower": "1954804254375936",
			},
			"HasMinPower": false,
			"test":        3,
		},
		JSONRPC: "2.0",
		Error:   nil,
		ID:      0,
	}, nil)
	lotusRepo := repo.NewLotusRPCRepo(redis.NewClient(&redis.Options{
		ClientName: "SanpshotTest",
	}))
	lotusService := NewLotusService(lotusRepo)
	tipsetKey, err := lotusService.GetTipSetByHeight(context.Background(), 314159, 2057965)
	assert.Nil(t, err)
	rsp, err := lotusService.GetMinerPowerByHeight(context.Background(), 314159, "t03751", tipsetKey)
	assert.Nil(t, err)

	zap.L().Info("result", zap.Any("miner", rsp))
}

func TestGetClientBalanceByHeight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lotusRpcClient := mock.NewMockRPCClient(gomock.NewController(t))

	lotusRpcClient.EXPECT().Call(gomock.Any(), "Filecoin.ChainGetTipSetByHeight", int64(2057965), gomock.Any()).Return(&jsonrpc.RPCResponse{
		Result: map[string]interface{}{
			"Cids": []interface{}{
				map[string]interface{}{
					"/": "bafy2bzacebes5mgufnyilg5e5p2mvhftqvoplzo65zxvcjy3j62n6w4er3hmm",
				},
				map[string]interface{}{
					"/": "bafy2bzacecqarqyklux426of26bwmcvd3rjjjlsqprd3umvbjtcn6hcbjoypa",
				},
			},
			"Blocks": 1,
			"Height": json.Number("2057965"),
			"Test":   1,
		},
		JSONRPC: "2.0",
		Error:   nil,
		ID:      0,
	}, nil)

	lotusRpcClient.EXPECT().Call(gomock.Any(), "Filecoin.StateMarketDeals", gomock.Any()).Return(&jsonrpc.RPCResponse{
		Result: map[string]interface{}{
			"Deals": map[string]interface{}{
				"116140": map[string]interface{}{
					"Proposal": map[string]interface{}{
						"VerifiedDeal": true,
						"EndEpoch":     json.Number("1729065330"),
						"PieceSize":    json.Number("1048576"),
						"Client":       "t03751",
					},
					"State": "",
				},
				"test": 4,
			},
		},
		JSONRPC: "2.0",
		Error:   nil,
		ID:      0,
	}, nil)
	lotusRepo := repo.NewLotusRPCRepo(redis.NewClient(&redis.Options{
		ClientName: "SanpshotTest",
	}))
	lotusService := NewLotusService(lotusRepo)
	_, err := lotusService.GetClientBalanceByHeight(context.Background(), 314159, 2057965)
	assert.Nil(t, err)
}

func TestGetNewestHeightAndTipset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lotusRpcClient := mock.NewMockRPCClient(gomock.NewController(t))

	lotusRpcClient.EXPECT().Call(gomock.Any(), "Filecoin.ChainHead").Return(&jsonrpc.RPCResponse{
		Result: map[string]interface{}{
			"Cids": []interface{}{
				map[string]interface{}{
					"/": "bafy2bzacebes5mgufnyilg5e5p2mvhftqvoplzo65zxvcjy3j62n6w4er3hmm",
				},
				map[string]interface{}{
					"/": "bafy2bzacecqarqyklux426of26bwmcvd3rjjjlsqprd3umvbjtcn6hcbjoypa",
				},
			},
			"Blocks": 1,
			"Height": json.Number("2072508"),
			"Test":   5,
		},
		JSONRPC: "2.0",
		Error:   nil,
		ID:      0,
	}, nil)

	lotusRepo := repo.NewLotusRPCRepo(redis.NewClient(&redis.Options{
		ClientName: "SanpshotTest",
	}))
	lotusService := NewLotusService(lotusRepo)
	_, err := lotusService.GetNewestHeight(context.Background(), 314159)
	assert.Nil(t, err)
}

func TestGetBlockHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lotusRpcClient := mock.NewMockRPCClient(gomock.NewController(t))
	lotusRpcClient.EXPECT().Call(gomock.Any(), "Filecoin.ChainGetTipSetByHeight", int64(2057965), gomock.Any()).Return(&jsonrpc.RPCResponse{
		Result: map[string]interface{}{
			"Cids": []interface{}{
				map[string]interface{}{
					"/": "bafy2bzacebes5mgufnyilg5e5p2mvhftqvoplzo65zxvcjy3j62n6w4er3hmm",
				},
				map[string]interface{}{
					"/": "bafy2bzacecqarqyklux426of26bwmcvd3rjjjlsqprd3umvbjtcn6hcbjoypa",
				},
			},
			"Blocks": 1,
			"Height": json.Number("2057965"),
			"Test":   1,
		},
		JSONRPC: "2.0",
		Error:   nil,
		ID:      0,
	}, nil)

	lotusRpcClient.EXPECT().Call(gomock.Any(), "Filecoin.ChainGetBlock", gomock.Any()).Return(&jsonrpc.RPCResponse{
		Result: map[string]interface{}{
			"Blocks":    1,
			"Height":    json.Number("2057965"),
			"Timestamp": json.Number("1729065330"),
			"Test":      2,
		},
		JSONRPC: "2.0",
		Error:   nil,
		ID:      0,
	}, nil)

	lotusRepo := repo.NewLotusRPCRepo(redis.NewClient(&redis.Options{
		ClientName: "SanpshotTest",
	}))
	lotusService := NewLotusService(lotusRepo)

	rsp, err := lotusService.GetBlockHeader(context.Background(), 314159, 2057965)
	assert.Nil(t, err)
	assert.Equal(t, int64(2057965), rsp.Height)
	assert.Equal(t, int64(1729065330), rsp.Timestamp)

	zap.L().Info("result", zap.Any("block", rsp))
}
