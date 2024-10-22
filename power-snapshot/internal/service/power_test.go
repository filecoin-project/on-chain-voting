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
	"power-snapshot/utils"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"
)

func TestGetPower(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lotusRpcClient := utils.NewMockRPCClient(gomock.NewController(t))

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
			"Test":   5,
		},
		JSONRPC: "2.0",
		Error:   nil,
		ID:      0,
	}, nil)

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
			"Test":   6,
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

	GetPower(context.Background(), lotusRpcClient, "t017592", "t03751", 2057965)
}

func TestGetActorPower(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lotusRpcClient := utils.NewMockRPCClient(gomock.NewController(t))

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
			"Test":   3,
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

	walletBalance, clientBalance, err := GetActorPower(context.Background(), lotusRpcClient, "t017592", 2057965)

	assert.Nil(t, err)

	zap.L().Info("result", zap.Any("walletBalance", walletBalance), zap.Any("clientBalance", clientBalance))
}

func TestGetMinerPower(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lotusRpcClient := utils.NewMockRPCClient(gomock.NewController(t))

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

	GetMinerPower(context.Background(), lotusRpcClient, "t03751", 2057965)
}
