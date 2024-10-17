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
	"log"
	"power-snapshot/config"
	"power-snapshot/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetPower(t *testing.T) {

	// Initialize the logger
	config.InitLogger()

	// Load the configuration from the specified path
	err := config.InitConfig("../../")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	// Initialize the client manager
	manager, err := utils.NewGoEthClientManager(config.Client.Network)
	if err != nil {
		log.Fatalf("Failed to initialize client manager: %v", err)
		return
	}

	// Get the client for the specified network
	client, err := manager.GetClient(314159)
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
		return
	}

	// Create a new Lotus RPC client
	lotusRpcClient := utils.NewClient(client.QueryRpc[0])
	GetPower(context.Background(), lotusRpcClient, "t017592", "t03751", 1490767)
}

func TestGetActorPower(t *testing.T) {
	ctx := context.Background()

	// Initialize the logger
	config.InitLogger()

	// Load the configuration from the specified path
	err := config.InitConfig("../../")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	// Initialize the client manager
	manager, err := utils.NewGoEthClientManager(config.Client.Network)
	if err != nil {
		log.Fatalf("Failed to initialize client manager: %v", err)
		return
	}

	// Get the client for the specified network
	client, err := manager.GetClient(314159)
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
		return
	}

	lotusRpcClient := utils.NewClient(client.QueryRpc[0])

	walletBalance, clientBalance, err := GetActorPower(ctx, lotusRpcClient, "t099523", 2058000)

	assert.Nil(t, err)

	zap.L().Info("result", zap.Any("walletBalance", walletBalance), zap.Any("clientBalance", clientBalance))
}

func TestGetMinerPower(t *testing.T) {
	ctx := context.Background()

	// Initialize the logger
	config.InitLogger()

	// Load the configuration from the specified path
	err := config.InitConfig("../../")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	// Initialize the client manager
	manager, err := utils.NewGoEthClientManager(config.Client.Network)
	if err != nil {
		log.Fatalf("Failed to initialize client manager: %v", err)
		return
	}

	// Get the client for the specified network
	client, err := manager.GetClient(314159)
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
		return
	}

	lotusRpcClient := utils.NewClient(client.QueryRpc[0])
	GetMinerPower(ctx, lotusRpcClient, "t03751", 1)
}
