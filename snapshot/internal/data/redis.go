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

package data

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"power-snapshot/config"
)

type RedisClient struct {
	client  *redis.Client
	expDays int
}

// writeOptions controls expiration behavior for write operations
type writeOptions struct {
	expire bool
}

// WriteOption configures write operation behavior
type WriteOption func(*writeOptions)

// WithExpire enables automatic expiration (default for write operations)
func WithExpire() WriteOption {
	return func(o *writeOptions) {
		o.expire = true
	}
}

// WithoutExpire disables automatic expiration
func WithoutExpire() WriteOption {
	return func(o *writeOptions) {
		o.expire = false
	}
}

func NewRedisClient(expDays int) (*RedisClient, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Client.Redis.URI,
		Username: config.Client.Redis.User,
		Password: config.Client.Redis.Password,
		DB:       config.Client.Redis.DB,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		zap.S().Error("failed to connect to redis", zap.Error(err))
		return nil, err
	}

	return &RedisClient{
		client:  redisClient,
		expDays: expDays,
	}, nil
}

// Set stores key-value pair with expiration control
// Default uses expiration days from config, use WithoutExpire() to disable
func (l *RedisClient) Set(ctx context.Context, key string, value interface{}, opts ...WriteOption) error {
	cfg := &writeOptions{expire: true}
	for _, opt := range opts {
		opt(cfg)
	}

	var exp time.Duration
	if cfg.expire {
		exp = l.expiration()
	}
	return l.client.Set(ctx, key, value, exp).Err()
}

func (l *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return l.client.Get(ctx, key).Result()
}

// HSet stores hash field-value pairs with expiration control
// Uses transaction pipeline when expiration is enabled for atomic operation
func (l *RedisClient) HSet(ctx context.Context, key string, values []interface{}, opts ...WriteOption) error {
	cfg := &writeOptions{expire: true}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.expire {
		_, err := l.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HSet(ctx, key, values...)
			pipe.Expire(ctx, key, l.expiration())
			return nil
		})
		return err
	}
	return l.client.HSet(ctx, key, values...).Err()
}

func (l *RedisClient) HGet(ctx context.Context, key, field string) (string, error) {
	return l.client.HGet(ctx, key, field).Result()
}

func (l *RedisClient) Del(ctx context.Context, keys ...string) error {
	return l.client.Del(ctx, keys...).Err()
}

func (l *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	result := l.client.Exists(ctx, key)
	return result.Val() > 0, result.Err()
}

func (l *RedisClient) ExpireTime(ctx context.Context, key string) (time.Duration, error) {
	return l.client.TTL(ctx, key).Result()
}

func (l *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return l.client.HGetAll(ctx, key).Result()
}
func (l *RedisClient) HExists(ctx context.Context, key, field string) (bool, error) {
	result := l.client.HExists(ctx, key, field)
	if result.Err() != nil {
		return false, fmt.Errorf("HEXISTS command error: %w", result.Err())
	}
	return result.Val(), nil
}

func (l *RedisClient) Close() error {
	return l.client.Close()
}

func (l *RedisClient) Keys(ctx context.Context, pattern string) ([]string, error) {
	return l.client.Keys(ctx, pattern).Result()
}

// expiration converts expiration days to duration
func (l *RedisClient) expiration() time.Duration {
	if l.expDays <= 0 {
		return 0
	}
	return time.Duration(l.expDays) * 24 * time.Hour
}
