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

package task

import (
	"go.uber.org/zap"

	"powervoting-server/config"
	"powervoting-server/data"
	"powervoting-server/service"
	"powervoting-server/task/event"
)

// SyncEventHandler handles the synchronization of contract events across multiple networks.
// It initializes a goroutine for each network to subscribe to and process contract events.
// Errors encountered during synchronization are collected and logged at the end.
//
// Parameters:
//   - syncService: The sync service used to manage synchronization operations.
func SyncEventHandler(syncService *service.SyncService) {
	network := config.Client.Network

	// Get the Ethereum client for the current network
	ethClient, err := data.GetClient(syncService, network.ChainId)
	if err != nil {
		zap.L().Error("get go-eth client error:", zap.Error(err))
		return
	}

	// Increment the WaitGroup counter for the new goroutine

	var syncEvent = &event.Event{
		SyncService: syncService,
		Network:     &network,
		Client:      ethClient,
	}

	// Subscribe to contract events for the current network
	if err := syncEvent.SubscribeEvent(); err != nil {
		zap.L().Error("sync finished with err:", zap.Error(err))
	}
}
