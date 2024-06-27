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
	"errors"
	"fmt"
	"power-snapshot/constant"
	models "power-snapshot/internal/model"
	"power-snapshot/utils"

	"github.com/redis/go-redis/v9"
)

type QueryRepoImpl struct {
	redisCli *redis.Client
	manager  *utils.GoEthClientManager
}

func NewQueryRepoImpl(client *redis.Client, manager *utils.GoEthClientManager) (*QueryRepoImpl, error) {
	return &QueryRepoImpl{
		redisCli: client,
		manager:  manager,
	}, nil
}

func (q *QueryRepoImpl) GetAddressPower(ctx context.Context, netId int64, address string, dayStr string) (*models.SyncPower, error) {
	key := fmt.Sprintf(constant.RedisAddrPower, netId, address)
	raw, err := q.redisCli.HGet(ctx, key, dayStr).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var m models.SyncPower
	err = json.Unmarshal([]byte(raw), &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (q *QueryRepoImpl) GetDeveloperWeights(ctx context.Context, dateStr string) (map[string]int64, error) {
	key := constant.RedisDeveloperPower
	resStr, err := q.redisCli.HGet(ctx, key, dateStr).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var m map[string]int64
	err = json.Unmarshal([]byte(resStr), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
