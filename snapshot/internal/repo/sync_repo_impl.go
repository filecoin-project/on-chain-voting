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
	"strings"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"power-snapshot/constant"
	"power-snapshot/internal/data"
	models "power-snapshot/internal/model"
)

type SyncRepoImpl struct {
	netId       int64
	redisClient *data.RedisClient
	stream      jetstream.JetStream
	consumer    jetstream.Consumer
	publisher   jetstream.Publisher
}

func NewSyncRepoImpl(netID int64, redisClient *data.RedisClient, stream jetstream.JetStream) (*SyncRepoImpl, error) {
	// init mq
	cfg := jetstream.StreamConfig{
		Name:      "TASKS",
		Retention: jetstream.WorkQueuePolicy,
		Subjects:  []string{"tasks.>"},
		Storage:   jetstream.FileStorage,
	}

	_, err := stream.CreateOrUpdateStream(context.Background(), cfg)
	if err != nil {
		zap.S().Error("Failed to create stream", zap.Error(err))
		return nil, err
	}

	consumer, err := stream.CreateOrUpdateConsumer(context.Background(), "TASKS",
		jetstream.ConsumerConfig{
			Name:          fmt.Sprintf("processor-%d", netID),
			FilterSubject: fmt.Sprintf("tasks.%d", netID),
			DeliverPolicy: jetstream.DeliverAllPolicy,
			MaxDeliver:    1440,
			AckWait:       1 * time.Minute,
		})
	if err != nil {
		zap.S().Error("failed to create consumer", zap.Error(err))
		return nil, err
	}

	return &SyncRepoImpl{
		netId:       netID,
		redisClient: redisClient,
		stream:      stream,
		consumer:    consumer,
		publisher:   stream,
	}, nil
}

// GetAllAddrSyncedDateMap retrieves synchronization dates mapping for all addresses under specified network
//
// Parameters:
//
//	ctx   : Context for request lifecycle control
//	netId : Network identifier used to construct Redis key
//
// Returns:
//
//	map[string][]string: Address-to-dates mapping (e.g. {"addr1": ["20250301", "20250302"]})
//	error             : Returns Redis operation errors
func (s *SyncRepoImpl) GetAllAddrSyncedDateMap(ctx context.Context, netId int64) (map[string][]string, error) {
	// Construct Redis hash key using format pattern from constants
	key := fmt.Sprintf(constant.RedisAddrSyncedDate, netId)

	// Retrieve all fields/values from Redis hash
	res, err := s.redisClient.HGetAll(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	// Initialize result map with pre-allocated space
	m := make(map[string][]string)

	// Convert hash values to string slices using comma delimiter
	for k, v := range res {
		m[k] = strings.Split(v, ",")
	}

	return m, nil
}

func (s *SyncRepoImpl) GetAddrSyncedDate(ctx context.Context, netId int64, addr string) ([]string, error) {
	key := fmt.Sprintf(constant.RedisAddrSyncedDate, netId)
	res, err := s.redisClient.HGet(ctx, key, addr)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	m := strings.Split(res, ",")
	return m, nil
}

func (s *SyncRepoImpl) SetAddrSyncedDate(ctx context.Context, netId int64, addr string, dates []string) error {
	key := fmt.Sprintf(constant.RedisAddrSyncedDate, netId)
	datesStr := strings.Join(dates, ",")
	err := s.redisClient.HSet(ctx, key, []any{
		addr, datesStr,
	}, data.WithoutExpire())
	if err != nil {
		return err
	}

	return nil
}

func (s *SyncRepoImpl) GetAddrPower(ctx context.Context, netId int64, addr string) (map[string]models.SyncPower, error) {
	key := fmt.Sprintf(constant.RedisAddrPower, netId, addr)
	raw, err := s.redisClient.HGetAll(ctx, key)
	if err != nil {
		return nil, err
	}

	m := make(map[string]models.SyncPower)
	for k, v := range raw {
		var temp models.SyncPower
		err := json.Unmarshal([]byte(v), &temp)
		if err != nil {
			return nil, err
		}
		m[k] = temp
	}

	return m, nil
}

func (s *SyncRepoImpl) SetAddrPower(ctx context.Context, netId int64, addr string, in map[string]models.SyncPower) error {
	if len(in) == 0 {
		return nil
	}

	key := fmt.Sprintf(constant.RedisAddrPower, netId, addr)
	m := make(map[string]string)
	for k, power := range in {
		jsonStr, err := json.Marshal(power)
		if err != nil {
			return err
		}
		m[k] = string(jsonStr)
	}

	err := s.redisClient.HSet(ctx, key, []any{m})
	if err != nil {
		return err
	}

	return nil
}

func (s *SyncRepoImpl) AddTask(ctx context.Context, netID int64, task *models.Task) error {
	key := fmt.Sprintf("tasks.%d", netID)
	jsonStr, err := json.Marshal(task)
	if err != nil {
		return err
	}
	_, err = s.stream.Publish(ctx, key, jsonStr)
	if err != nil {
		return err
	}

	return nil
}

func (s *SyncRepoImpl) GetTask(ctx context.Context, netID int64) (jetstream.MessageBatch, error) {
	mc, err := s.consumer.Fetch(10)
	if err != nil {
		return nil, err
	}

	return mc, nil
}

func (s *SyncRepoImpl) SetDeveloperWeights(ctx context.Context, dateStr string, in map[string]int64) error {
	key := constant.RedisDeveloperPower
	inJson, err := json.Marshal(in)
	if err != nil {
		return err
	}
	err = s.redisClient.HSet(ctx, key, []any{dateStr, inJson})
	if err != nil {
		return err
	}

	return nil
}

func (s *SyncRepoImpl) GetUserDeveloperWeights(ctx context.Context, dateStr string, username string) (int64, error) {
	key := constant.RedisDeveloperPower
	resStr, err := s.redisClient.HGet(ctx, key, dateStr)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}

	var m map[string]int64
	err = json.Unmarshal([]byte(resStr), &m)
	if err != nil {
		return 0, err
	}

	res, ok := m[username]
	if !ok {
		return 0, nil
	}

	return res, nil
}

func (s *SyncRepoImpl) GetDeveloperWeights(ctx context.Context, dateStr string) (map[string]int64, error) {
	key := constant.RedisDeveloperPower
	resStr, err := s.redisClient.HGet(ctx, key, dateStr)
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

func (s *SyncRepoImpl) ExistDeveloperWeights(ctx context.Context, dateStr string) (bool, error) {
	key := constant.RedisDeveloperPower
	exist, err := s.redisClient.HExists(ctx, key, dateStr)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}

	return exist, nil
}
