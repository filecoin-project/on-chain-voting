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
	"golang.org/x/sync/errgroup"
	"log"
	"power-snapshot/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetDeveloperWeights(t *testing.T) {
	// Initialize the logger
	config.InitLogger()

	// Load the configuration from the specified path
	err := config.InitConfig("../../")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	totalWeights, err := GetDeveloperWeights(time.Date(2024, 6, 21, 12, 00, 00, 00, time.Local))
	if err != nil {
		log.Fatalf("Failed to get developer weights: %v", err)
	}
	fmt.Println("len(totalWeights):", len(totalWeights))
	assert.NotEmpty(t, totalWeights, 1)
}

func TestGetContributors(t *testing.T) {
	err := config.InitConfig("../../")
	if err != nil {
		return
	}
	config.InitLogger()
	since := addMonths(time.Now(), -6).Format("2006-01-02T15:04:05Z")
	ctx := context.Background()
	tokenManager := NewGitHubTokenManager(config.Client.Github.Token)

	token := tokenManager.GetCoreAvailableToken()
	res, err := getContributors(ctx, "filecoin-project", "venus", since, token)
	assert.Nil(t, err)
	fmt.Printf("%v", res)
	assert.NotEmpty(t, res)
}

func TestAddMonths(t *testing.T) {
	res := addMonths(time.Date(2024, 05, 04, 00, 00, 00, 00, time.Local),
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
		"test1": 5,
		"test2": 2,
		"test3": 5,
	}

	res := addMerge(map1, map2)
	assert.Equal(t, res, expected)
}

func worker(ctx context.Context, id int) error {
	fmt.Printf("Worker %d is processing business logic\n", id)

	// 例如，处理一些数据或调用外部API
	// 模拟可能发生错误的情况
	if id == 60 {
		println("60出错了!!!")
		return fmt.Errorf("business logic error in worker %d", id)
	}

	// 模拟每次 io操作时候的检查 err
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		if ctx.Err() != nil {
			println("中奖了!", id)
			return ctx.Err()
		}
	}

	println("Worker done in fact ", id)
	// 检查上下文是否被取消
	select {
	case <-ctx.Done():
		fmt.Printf("Worker %d cancel\n", id)
	default:
		fmt.Printf("Worker %d finished work\n", id)
	}

	return nil
}

// nestedWorker 启动嵌套的 worker 协程
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
