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
	"context"
	"sync"
	"time"

	"github.com/golang-module/carbon"
	"go.uber.org/zap"

	"power-snapshot/config"
	"power-snapshot/constant"
	"power-snapshot/internal/data"
	models "power-snapshot/internal/model"
)

// SyncPower is a function that returns a closure for syncing power data across different networks.
func (j *Safejob) SyncPower() func() {
	// The returned function encapsulates the logic for syncing power data.
	return func() {
		// Create a background context for the operations.
		ctx := context.Background()
		// Iterate over each network configuration in the client's network list.
		for _, network := range config.Client.Network {
			// sync date height
			err := j.syncService.SyncDateHeight(ctx, network.ChainId)
			if err != nil {
				zap.L().Error("failed to sync date height, it will skipped ", zap.Error(err), zap.Int64("network_id", network.ChainId))
				continue
			}

			err = j.syncService.SyncAllAddrPower(ctx, network.ChainId)
			if err != nil {
				zap.L().Error("failed to sync all addr power, it will skipped ", zap.Error(err), zap.Int64("network_id", network.ChainId))
				continue
			}
		}
	}
}

// SyncDevWeightStepDay returns a function that synchronizes developer weights for each day within a specified range.
// It takes a pointer to a SyncService as an argument.
func (j *Safejob) SyncDevWeightStepDay() func() {
	// Return an anonymous function that performs the synchronization.
	return func() {
		// Create a background context for the operation.
		ctx := context.Background()
		// Calculate the start date as the current date minus the data expiration duration, and set it to the end of the day.
		start := carbon.Now().SubDays(constant.DataExpiredDuration).EndOfDay()
		// Calculate the end date as yesterday and set it to the end of the day.
		end := carbon.Now().Yesterday().EndOfDay()

		// find latest index
		for l := 0; l < len(config.Client.Github.Token)*2; {
			for i := start; i.Timestamp() <= end.Timestamp(); i = i.AddDay() {
				exist, err := j.syncService.ExistDeveloperWeight(ctx, i.ToShortDateString())
				if err != nil {
					zap.L().Error("SyncDevWeightStepDay", zap.String("date", i.ToShortDateString()))
					return
				}
				if !exist {
					err := j.syncService.SyncDeveloperWeight(ctx, i.ToShortDateString())
					if err != nil {
						return
					}
					break
				}
			}
			l++
		}
	}
}

// UploadPowerToIPFS returns a function that uploads power data to IPFS.
func (j *Safejob) UploadPowerToIPFS(w3client *data.W3Client) func() {
	return func() {
		zap.L().Info("backup power start: ", zap.Int64("timestamp", time.Now().Unix()))
		wg := sync.WaitGroup{}
		errList := make([]error, 0, len(config.Client.Network))
		mu := &sync.Mutex{}

		// Iterate over networks and upload power data to IPFS concurrently.
		for _, network := range config.Client.Network {
			wg.Add(1)
			go func(network models.Network) {
				defer wg.Done()
				ctx := context.Background()

				// Upload power data for the current network.
				if err := j.syncService.UploadPowerToIPFS(ctx, network.ChainId, w3client); err != nil {
					mu.Lock()
					errList = append(errList, err)
					mu.Unlock()
				}
			}(network)
		}

		// Wait for all goroutines to finish.
		wg.Wait()

		// Log errors if any occurred during the upload process.
		if len(errList) != 0 {
			zap.L().Error("backup power finished with err:", zap.Errors("errors", errList))
		}
	}
}
