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
	"github.com/redis/go-redis/v9"
	"github.com/ybbus/jsonrpc/v3"
	"hash/fnv"
	"power-snapshot/constant"
	models "power-snapshot/internal/model"
	"power-snapshot/utils"
	"time"
)

type BaseRepoImpl struct {
	manager     *utils.GoEthClientManager
	redisClient *redis.Client
}

func NewBaseRepoImpl(manager *utils.GoEthClientManager, redisClient *redis.Client) (*BaseRepoImpl, error) {
	return &BaseRepoImpl{
		manager:     manager,
		redisClient: redisClient,
	}, nil
}

func (s *BaseRepoImpl) ListVoterAddr(ctx context.Context, netID int64) ([]string, error) {
	client, err := s.manager.GetClient(netID)
	if err != nil {
		return nil, err
	}

	ethAddressList, err := utils.GetVoterAddresses(client)
	if err != nil {
		return []string{}, err
	}

	ethStrAddrList := make([]string, 0)
	for _, addr := range ethAddressList {
		if addr.Hex() == "0x0000000000000000000000000000000000000000" {
			continue
		}
		ethStrAddrList = append(ethStrAddrList, addr.String())
	}

	return ethStrAddrList, nil
}

func (s *BaseRepoImpl) GetVoteInfo(ctx context.Context, netId int64, addr string) (*models.VoterInfo, error) {
	// TODO get this from redis and sync every day
	client, err := s.manager.GetClient(netId)
	if err != nil {
		return nil, err
	}

	voterInfo, err := utils.GetVoterInfo(addr, client)
	if err != nil {
		return nil, err
	}

	return &voterInfo, nil
}

func (s *BaseRepoImpl) GetEthClient(ctx context.Context, netId int64) (*models.GoEthClient, error) {
	client, err := s.manager.GetClient(netId)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (s *BaseRepoImpl) GetLotusClient(ctx context.Context, netId int64) (jsonrpc.RPCClient, error) {
	client, err := s.manager.GetClient(netId)
	if err != nil {
		return nil, err
	}

	return jsonrpc.NewClient(client.QueryRpc[0]), nil
}

func (s *BaseRepoImpl) GetLotusClientByHashKey(ctx context.Context, netId int64, key string) (jsonrpc.RPCClient, error) {
	client, err := s.manager.GetClient(netId)
	if err != nil {
		return nil, err
	}
	h := fnv.New32a()
	data := fmt.Sprintf("%s_%d", key, time.Now().UnixNano())
	_, err = h.Write([]byte(data))
	if err != nil {
		return nil, err
	}
	index := h.Sum32() % uint32(len(client.QueryRpc))

	return jsonrpc.NewClient(client.QueryRpc[index]), nil
}

func (s *BaseRepoImpl) GetDateHeightMap(ctx context.Context, netId int64) (map[string]int64, error) {
	key := fmt.Sprintf(constant.RedisDateHeight, netId)
	jsonStr, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	m := make(map[string]int64)
	err = json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *BaseRepoImpl) SetDateHeightMap(ctx context.Context, netId int64, height map[string]int64) error {
	key := fmt.Sprintf(constant.RedisDateHeight, netId)
	jsonStr, err := json.Marshal(height)
	if err != nil {
		return err
	}
	err = s.redisClient.Set(ctx, key, jsonStr, redis.KeepTTL).Err()
	if err != nil {
		return err
	}

	return nil
}
