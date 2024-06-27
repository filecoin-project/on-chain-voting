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
	"power-snapshot/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewRedisClient() (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Client.Redis.URI,
		Username: config.Client.Redis.User,
		Password: config.Client.Redis.Password,
		DB:       config.Client.Redis.DB,
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		zap.S().Error("failed to connect to redis", zap.Error(err))
		return nil, err
	}

	return redisClient, nil
}
