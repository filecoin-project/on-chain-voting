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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"

	"power-snapshot/utils"
)

// func TestGetDeveloperWeights(t *testing.T) {
// 	// Initialize the logger
// 	config.InitLogger()

// 	// Load the configuration from the specified path
// 	err := config.InitConfig("../../")
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 		return
// 	}

// 	totalWeights, commit, err := GetDeveloperWeights(time.Now())
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, commit)
// 	zap.L().Info("result", zap.Any("client", totalWeights))
// }

// func TestGetContributors(t *testing.T) {
// 	err := config.InitConfig("../../")
// 	if err != nil {
// 		return
// 	}
// 	config.InitLogger()
// 	since := utils.AddMonths(time.Now(), -6).Format("2006-01-02T15:04:05Z")
// 	tokenManager := NewGitHubTokenManager(config.Client.Github.Token)

// 	token := tokenManager.GetCoreAvailableToken()
// 	res, err := getContributors(context.Background(), "filecoin-project", "on-chain-voting", since, token)
// 	assert.Nil(t, err)

// 	zap.L().Info("result", zap.Any("res", res))
// }

func TestAddMonths(t *testing.T) {
	res := utils.AddMonths(time.Date(2024, 05, 04, 00, 00, 00, 00, time.Local),
		-2)

	expected := time.Date(2024, 03, 04, 00, 00, 00, 00, time.Local)

	assert.Equal(t, res, expected)
}

func TestAddMerge(t *testing.T) {
	map1 := map[string]int64{
		"test1": 1,
		"test2": 2,
	}

	map2 := map[string]int64{
		"test1": 4,
		"test3": 5,
	}

	expected := map[string]int64{
		"test1": 1,
		"test2": 2,
		"test3": 5,
	}

	res := addMerge(map1, map2)
	assert.Equal(t, res, expected)
}

func worker(ctx context.Context, id int) error {
	fmt.Printf("Worker %d is processing business logic\n", id)

	if id == 60 {
		println("Worker 60 is error")
		return fmt.Errorf("business logic error in worker %d", id)
	}

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	println("Worker done in fact ", id)
	select {
	case <-ctx.Done():
		fmt.Printf("Worker %d cancel\n", id)
	default:
		fmt.Printf("Worker %d finished work\n", id)
	}

	return nil
}

func nestedWorker(ctx context.Context, id int) error {
	g, ctx := errgroup.WithContext(ctx)

	for i := 1; i <= 100; i++ {
		i := i // capture loop variable
		g.Go(func() error {
			return worker(ctx, id*10+i)
		})
	}

	return g.Wait()
}

func TestContext(t *testing.T) {
	// Create a context with cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create an errgroup with the context
	g, ctx := errgroup.WithContext(ctx)

	// Launch top-level nested workers
	for i := 1; i <= 2; i++ {
		i := i // capture loop variable
		g.Go(func() error {
			return nestedWorker(ctx, i)
		})
	}

	// Wait for all nested workers to finish or any error to occur
	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("All workers finish")
}

// func TestGetContributors(t *testing.T) {
// 	// ZearnLabs/DefiLlama-Adapters
// 	config.InitConfig("../../")
// 	res, err := getContributors(
// 		"",
// 		"",
// 		"",
// 		"",
// 	)

// 	assert.Nil(t, err)
// 	assert.NotEmpty(t, res)
// }
