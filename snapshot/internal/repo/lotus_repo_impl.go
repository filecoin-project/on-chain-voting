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
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	filecoinAddress "github.com/filecoin-project/go-address"
	"github.com/redis/go-redis/v9"
	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"

	"power-snapshot/config"
	"power-snapshot/constant"
	models "power-snapshot/internal/model"
	"power-snapshot/utils/types"
)

type LotusRPCRepo struct {
	rpcClient   jsonrpc.RPCClient
	redisClient *redis.Client
}

func NewLotusRPCRepo(redisClient *redis.Client) *LotusRPCRepo {
	return &LotusRPCRepo{
		rpcClient:   jsonrpc.NewClientWithOpts(config.Client.Network.QueryRpc[0], &jsonrpc.RPCClientOpts{}),
		redisClient: redisClient,
	}
}

func (l *LotusRPCRepo) GetTipSetByHeight(ctx context.Context, netId, height int64) ([]any, error) {
	key := fmt.Sprintf(constant.RedisTipset, netId)
	defer func() {
		go l.cleanExpiredHeights(context.Background(), netId)
	}()

	res, err := l.redisClient.HGet(ctx, key, strconv.FormatInt(height, 10)).Result()
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
	}

	if res != "" {
		var resTipSet []any
		if err := json.Unmarshal([]byte(res), &resTipSet); err != nil {
			return nil, err
		}
		return resTipSet, nil
	}

	resp, err := l.rpcClient.Call(ctx, "Filecoin.ChainGetTipSetByHeight", height, types.TipSetKey{})
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		zap.L().Error("get tipset by height error", zap.Int64("height", height), zap.Error(resp.Error))
		return []any{}, nil
	}

	rMap, ok := resp.Result.(map[string]interface{})
	if !ok {
		return nil, resp.Error
	}

	resTipSet := rMap["Cids"].([]interface{})

	if err := l.cacheWithExpiration(ctx, key, height, resTipSet); err != nil {
		zap.L().Warn("cache with expiration failed", zap.Error(err))
	}
	return resTipSet, nil
}

func (l *LotusRPCRepo) cleanExpiredHeights(ctx context.Context, netId int64) {
	lastBlockHeight, err := l.GetNewestHeight(ctx, netId)
	if err != nil {
		zap.L().Error("get newest height failed", zap.Error(err))
		return
	}
	expirationThreshold := lastBlockHeight + constant.SavedHeightDuration
	key := fmt.Sprintf(constant.RedisTipset, netId)

	var cursor uint64
	for {
		fields, nextCursor, err := l.redisClient.HScan(ctx, key, cursor, "", 1000).Result()
		if err != nil {
			zap.L().Error("HSCAN failed", zap.Error(err))
			return
		}

		var toDelete []string
		for i := 0; i < len(fields); i += 2 {
			heightStr := fields[i]
			height, err := strconv.ParseInt(heightStr, 10, 64)
			if err != nil || height < expirationThreshold {
				toDelete = append(toDelete, heightStr)
			}
		}
		if len(toDelete) > 0 {

			zap.L().Info("cleaning expired heights", zap.Int("count", len(toDelete)), zap.Any("cleaned heights", toDelete))
			if err := l.redisClient.HDel(ctx, key, toDelete...).Err(); err != nil {
				zap.L().Error("HDEL failed", zap.Strings("keys", toDelete), zap.Error(err))
			} else {
				zap.L().Info("cleaned expired heights",
					zap.Int("count", len(toDelete)),
					zap.Int64("threshold", expirationThreshold))
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
}

func (l *LotusRPCRepo) cacheWithExpiration(ctx context.Context, key string, height int64, data []any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	pipe := l.redisClient.Pipeline()
	pipe.HSet(ctx, key, strconv.FormatInt(height, 10), string(jsonData))

	pipe.Expire(ctx, key, constant.DataExpiredDuration*24*time.Hour)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("pipeline failed: %w", err)
	}
	return nil
}

func (l *LotusRPCRepo) GetAddrBalanceBySpecialHeight(ctx context.Context, addr string, netId, height int64) (string, error) {
	tipSetKey, err := l.GetTipSetByHeight(ctx, netId, height)
	if err != nil {
		return "", err
	}

	addressStr, err := filecoinAddress.NewFromString(addr)
	if err != nil {
		return "", err
	}

	tipSetList := append([]any{}, addressStr, tipSetKey)
	resp, err := l.rpcClient.Call(ctx, "Filecoin.StateGetActor", &tipSetList)
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

func (l *LotusRPCRepo) GetMinerPowerByHeight(ctx context.Context, netId int64, addr string, tipsetKey []interface{}) (models.LotusMinerPower, error) {
	addressStr, err := filecoinAddress.NewFromString(addr)
	if err != nil {
		return models.LotusMinerPower{}, err
	}

	rpcParams := append([]any{}, addressStr, tipsetKey)
	resp, err := l.rpcClient.Call(ctx, "Filecoin.StateMinerPower", rpcParams)
	if err != nil {
		return models.LotusMinerPower{}, err
	}

	tmp, err := json.Marshal(resp.Result)
	if err != nil {
		return models.LotusMinerPower{}, err
	}

	var minerPower models.LotusMinerPower
	if err := json.Unmarshal(tmp, &minerPower); err != nil {
		return models.LotusMinerPower{}, err
	}

	return minerPower, nil
}
//FIXME: StateMarketDeals is verry large, need to optimize
func (l *LotusRPCRepo) GetClientBalanceBySpecialHeight(ctx context.Context, netId, height int64) (models.StateMarketDeals, error) {
	// tipSetKey, err := l.GetTipSetByHeight(ctx, netId, height)
	// if err != nil {
	// 	return models.StateMarketDeals{}, err
	// }

	// tipSetList := append([]any{}, tipSetKey)

	// var t models.StateMarketDeals
	// resp, err := l.rpcClient.Call(ctx, "Filecoin.StateMarketDeals", tipSetList)
	// if err != nil {
	// 	return models.StateMarketDeals{}, err
	// }

	// tmp, err := json.Marshal(resp.Result)
	// if err != nil {
	// 	return models.StateMarketDeals{}, err
	// }

	// if err := json.Unmarshal(tmp, &t); err != nil {
	// 	return models.StateMarketDeals{}, err
	// }

	return models.StateMarketDeals{}, nil
}

func (l *LotusRPCRepo) GetNewestHeight(ctx context.Context, netId int64) (height int64, err error) {
	resp, err := l.rpcClient.Call(ctx, "Filecoin.ChainHead")
	if err != nil {
		return 0, nil
	}

	height, err = resp.Result.(map[string]interface{})["Height"].(json.Number).Int64()
	if err != nil {
		return 0, nil
	}

	return height, nil
}

func (l *LotusRPCRepo) GetBlockHeader(ctx context.Context, netId, height int64) (models.BlockHeader, error) {
	tipSetKey, err := l.GetTipSetByHeight(ctx, netId, height)
	if err != nil {
		return models.BlockHeader{}, err
	}

	marshal, err := json.Marshal(tipSetKey)
	if err != nil {
		return models.BlockHeader{}, err
	}

	var blockParam []models.GetBlockParam
	err = json.Unmarshal(marshal, &blockParam)
	if err != nil {
		return models.BlockHeader{}, err
	}

	resp, err := l.rpcClient.Call(ctx, "Filecoin.ChainGetBlock", []models.GetBlockParam{{Empty: blockParam[0].Empty}})

	if err != nil {
		return models.BlockHeader{}, err
	}

	if resp.Error != nil {
		return models.BlockHeader{}, resp.Error
	}

	tmp, err := json.Marshal(resp.Result.(map[string]interface{}))
	if err != nil {
		return models.BlockHeader{}, err
	}

	var block models.BlockHeader
	if err := json.Unmarshal(tmp, &block); err != nil {
		return block, err
	}

	return block, nil
}

func (l *LotusRPCRepo) GetWalletBalanceByHeight(ctx context.Context, id string, netId, height int64) (string, error) {
	tipSetKey, err := l.GetTipSetByHeight(ctx, netId, height)
	if err != nil {
		return "0", err
	}

	resp, err := l.rpcClient.Call(ctx, "Filecoin.StateGetActor", id, tipSetKey)
	if err != nil {
		return "0", err
	}

	if resp.Error != nil {
		return "0", resp.Error
	}

	tmp, err := json.Marshal(resp.Result)
	if err != nil {
		return "0", err
	}

	var t types.BalanceInfo
	if err := json.Unmarshal(tmp, &t); err != nil {
		return "0", err
	}

	return t.Balance, nil
}

//FIXME: StateMarketDeals is verry large, need to optimize
func (l *LotusRPCRepo) GetClientBalanceByHeight(ctx context.Context, netId, height int64) (types.StateMarketDeals, error) {
	// tipSetKey, err := l.GetTipSetByHeight(ctx, netId, height)
	// if err != nil {
	// 	return types.StateMarketDeals{}, err
	// }

	// tipSetList := append([]any{}, tipSetKey)

	// var t types.StateMarketDeals
	// resp, err := l.rpcClient.Call(ctx, "Filecoin.StateMarketDeals", tipSetList)
	// if err != nil {
	// 	return types.StateMarketDeals{}, err
	// }

	// tmp, err := json.Marshal(resp.Result)
	// if err != nil {
	// 	return types.StateMarketDeals{}, err
	// }

	// if err := json.Unmarshal(tmp, &t); err != nil {
	// 	return types.StateMarketDeals{}, err
	// }

	return types.StateMarketDeals{}, nil
}
