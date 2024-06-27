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
	"strings"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type SyncRepoImpl struct {
	netIds      []int64
	redisClient *redis.Client
	stream      jetstream.JetStream
	consumer    map[int64]jetstream.Consumer
	publisher   jetstream.Publisher
}

func NewSyncRepoImpl(netIDs []int64, redisClient *redis.Client, stream jetstream.JetStream) (*SyncRepoImpl, error) {
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

	consMap := make(map[int64]jetstream.Consumer)
	for _, netId := range netIDs {
		consumer, err := stream.CreateOrUpdateConsumer(context.Background(), "TASKS",
			jetstream.ConsumerConfig{
				Name:          fmt.Sprintf("processor-%d", netId),
				FilterSubject: fmt.Sprintf("tasks.%d", netId),
				DeliverPolicy: jetstream.DeliverAllPolicy,
				MaxDeliver:    1440,
				AckWait:       1 * time.Minute,
			})
		if err != nil {
			zap.S().Error("failed to create consumer", zap.Error(err))
			return nil, err
		}
		consMap[netId] = consumer
	}

	return &SyncRepoImpl{
		netIds:      netIDs,
		redisClient: redisClient,
		stream:      stream,
		consumer:    consMap,
		publisher:   stream,
	}, nil
}

func (s *SyncRepoImpl) GetAllAddrSyncedDateMap(ctx context.Context, netId int64) (map[string][]string, error) {
	key := fmt.Sprintf(constant.RedisAddrSyncedDate, netId)
	res, err := s.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	m := make(map[string][]string)
	for k, v := range res {
		m[k] = strings.Split(v, ",")
	}

	return m, nil
}

func (s *SyncRepoImpl) GetAddrSyncedDate(ctx context.Context, netId int64, addr string) ([]string, error) {
	key := fmt.Sprintf(constant.RedisAddrSyncedDate, netId)
	res, err := s.redisClient.HGet(ctx, key, addr).Result()
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
	err := s.redisClient.HSet(ctx, key, addr, datesStr).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *SyncRepoImpl) GetAddrPower(ctx context.Context, netId int64, addr string) (map[string]*models.SyncPower, error) {
	key := fmt.Sprintf(constant.RedisAddrPower, netId, addr)
	raw, err := s.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	m := make(map[string]*models.SyncPower)
	for k, v := range raw {
		var temp models.SyncPower
		err := json.Unmarshal([]byte(v), &temp)
		if err != nil {
			return nil, err
		}
		m[k] = &temp
	}

	return m, nil
}

func (s *SyncRepoImpl) SetAddrPower(ctx context.Context, netId int64, addr string, in map[string]*models.SyncPower) error {
	key := fmt.Sprintf(constant.RedisAddrPower, netId, addr)
	m := make(map[string]string)
	for k, power := range in {
		jsonStr, err := json.Marshal(power)
		if err != nil {
			return err
		}
		m[k] = string(jsonStr)
	}

	err := s.redisClient.HSet(ctx, key, m).Err()
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
	cons, ok := s.consumer[netID]
	if !ok {
		return nil, fmt.Errorf("not found consumer by netID(%d)", netID)
	}
	mc, err := cons.Fetch(10)
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
	err = s.redisClient.HSet(ctx, key, dateStr, inJson).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *SyncRepoImpl) GetUserDeveloperWeights(ctx context.Context, dateStr string, username string) (int64, error) {
	key := constant.RedisDeveloperPower
	resStr, err := s.redisClient.HGet(ctx, key, dateStr).Result()
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
	resStr, err := s.redisClient.HGet(ctx, key, dateStr).Result()
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
	exist, err := s.redisClient.HExists(ctx, key, dateStr).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}

	return exist, nil
}
