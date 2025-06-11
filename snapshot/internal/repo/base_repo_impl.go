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
	"hash/fnv"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ybbus/jsonrpc/v3"

	"power-snapshot/config"
	"power-snapshot/constant"
	"power-snapshot/internal/data"
	models "power-snapshot/internal/model"
	"power-snapshot/utils/types"
)

type BaseRepoImpl struct {
	ethClient     *data.GoEthClientManager
	redisClient   *data.RedisClient
}

func NewBaseRepoImpl(manager *data.GoEthClientManager, redisClient *data.RedisClient) *BaseRepoImpl {
	return &BaseRepoImpl{
		ethClient:   manager,
		redisClient: redisClient,
	}
}

func (s *BaseRepoImpl) GetLotusClient(ctx context.Context, netId int64) (jsonrpc.RPCClient, error) {
	client := s.ethClient.GetClient()

	return jsonrpc.NewClient(client.QueryRpc[0]), nil
}

func (s *BaseRepoImpl) GetLotusClientByHashKey(ctx context.Context, netId int64, key string) (jsonrpc.RPCClient, error) {
	client := s.ethClient.GetClient()

	h := fnv.New32a()
	data := fmt.Sprintf("%s_%d", key, time.Now().UnixNano())
	_, err := h.Write([]byte(data))
	if err != nil {
		return nil, err
	}

	index := h.Sum32() % uint32(len(client.QueryRpc))

	return jsonrpc.NewClient(client.QueryRpc[index]), nil
}

// GetDateHeightMap retrieves date-to-block-height mapping from Redis storage
//
// Parameters:
//
//	ctx   : context.Context - Request context for cancellation/timeout control
//	netId : int64           - Network identifier used for Redis key construction
//
// Returns:
//
//	map[string]int64 - Mapping of dates (YYYYMMDD) to block heights. (e.g. {"20250301": 123456})
//	error           - Returns Redis connection errors or JSON unmarshal failures
func (s *BaseRepoImpl) GetDateHeightMap(ctx context.Context, netId int64) (map[string]int64, error) {
	// Construct Redis key using configured pattern (e.g. "date:height:%d")
	key := fmt.Sprintf(constant.RedisDateHeight, netId)

	// Fetch JSON-encoded data from Redis
	jsonStr, err := s.redisClient.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return make(map[string]int64), nil
		}
		return nil, err
	}

	// Initialize target map structure
	m := make(map[string]int64)

	// Deserialize JSON string to map
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

	err = s.redisClient.Set(ctx, key, jsonStr)
	if err != nil {
		return err
	}

	return nil
}

func (s *BaseRepoImpl) SaveToLocalFile(ctx context.Context, netId int64, dayStr, t string, data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var path string
	var dir = config.Client.FilePath
	switch t {
	case constant.DEVELOPER:
		path = filepath.Join(dir, fmt.Sprintf(constant.DeveloperWeightsFile, netId, dayStr))
	case constant.DEALS:
		path = filepath.Join(dir, fmt.Sprintf(constant.DealsFile, netId, dayStr))
	default:
		return fmt.Errorf("SaveToLocalFile invalid type: %s", t)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := os.WriteFile(path+constant.FileSuffix, jsonData, 0644); err != nil {
		return err
	}

	return s.cleanupOldFiles(netId, t)
}

func (s *BaseRepoImpl) GetDeveloperWeights(ctx context.Context, netId int64, dayStr string) ([]models.Nodes, error) {
	filename := filepath.Join(config.Client.FilePath, fmt.Sprintf(constant.DeveloperWeightsFile, netId, dayStr))

	data, err := os.ReadFile(filename + constant.FileSuffix)
	if err != nil {
		return nil, err
	}

	var commits []models.Nodes
	if err := json.Unmarshal(data, &commits); err != nil {
		return nil, err
	}

	return commits, nil
}

func (s *BaseRepoImpl) GetDealsFromLocal(ctx context.Context, netId int64, dayStr string) (types.StateMarketDeals, error) {
	filename := filepath.Join(config.Client.FilePath, fmt.Sprintf(constant.DealsFile, netId, dayStr))
	data, err := os.ReadFile(filename + constant.FileSuffix)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var deals types.StateMarketDeals
	if err := json.Unmarshal(data, &deals); err != nil {
		return nil, err
	}

	return deals, nil
}

func (s *BaseRepoImpl) cleanupOldFiles(netId int64, t string) error {
	var pattern string

	switch t {
	case constant.DEVELOPER:
		pattern = filepath.Join(config.Client.FilePath, fmt.Sprintf(constant.DeveloperWeightsFile, netId, ""))
	case constant.DEALS:
		pattern = filepath.Join(config.Client.FilePath, fmt.Sprintf(constant.DealsFile, netId, ""))
	default:
		return fmt.Errorf("cleanupOldFiles invalid type: %s", t)
	}

	pattern += "*"
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -60)

	for _, file := range files {
		dateStr := filepath.Base(file)
		group := strings.Split(dateStr, "_")
		dateStr = group[2]

		fileDate, err := time.Parse("20060102", dateStr)
		if err != nil {
			continue
		}

		if fileDate.Before(cutoff) {
			if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}
	return nil
}
