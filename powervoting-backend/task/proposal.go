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
	"sync"

	"go.uber.org/zap"

	"powervoting-server/config"
	"powervoting-server/data"
	"powervoting-server/service"
)

// SyncEventHandler handles the synchronization of contract events across multiple networks.
// It initializes a goroutine for each network to subscribe to and process contract events.
// Errors encountered during synchronization are collected and logged at the end.
//
// Parameters:
//   - syncService: The sync service used to manage synchronization operations.
func SyncEventHandler(syncService *service.SyncService) {
	// Use a WaitGroup to wait for all goroutines to finish
	wg := sync.WaitGroup{}
	// Use a slice to collect errors from all goroutines
	errList := make([]error, 0, len(config.Client.Network))
	// Use a mutex to safely append errors from multiple goroutines
	mu := &sync.Mutex{}

	// Iterate over each network configuration
	for _, network := range config.Client.Network {
		network := network // Create a local copy of the network variable for the goroutine

		// Get the Ethereum client for the current network
		ethClient, err := data.GetClient(syncService, network.ChainId)
		if err != nil {
			zap.L().Error("get go-eth client error:", zap.Error(err))
			continue // Skip this network if the client cannot be initialized
		}

		// Increment the WaitGroup counter for the new goroutine
		wg.Add(1)
		go func(network config.Network) {
			defer wg.Done() // Decrement the WaitGroup counter when the goroutine completes

			// Create an event handler for the current network
			var syncEvent = &Event{
				SyncService: syncService,
				Network:     &network,
				Client:      ethClient,
			}

			// Subscribe to contract events for the current network
			if err := syncEvent.SubscribeEvent(); err != nil {
				mu.Lock()
				errList = append(errList, err) // Append the error to the error list
				mu.Unlock()
			}
		}(network) // Pass the local copy of the network to the goroutine
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Log any errors encountered during synchronization
	if len(errList) != 0 {
		zap.L().Error("sync finished with err:", zap.Errors("errors", errList))
	}
}
