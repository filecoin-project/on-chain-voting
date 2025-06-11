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
	"time"

	"github.com/golang-module/carbon"
	"go.uber.org/zap"

	"power-snapshot/config"
	"power-snapshot/internal/data"
)

// SyncPower is a function that returns a closure for syncing power data across different networks.
func (j *Safejob) SyncPower() {
	// The returned function encapsulates the logic for syncing power data.

	// Create a background context for the operations.
	ctx := context.Background()
	// Iterate over each network configuration in the client's network list.

	// sync date height
	err := j.syncService.SyncDateHeight(ctx, config.Client.Network.ChainId)
	if err != nil {
		zap.L().Error("failed to sync date height, it will skipped ", zap.Error(err), zap.Int64("network_id", config.Client.Network.ChainId))
		return
	}

	err = j.syncService.SyncAllAddrPower(ctx, config.Client.Network.ChainId)
	if err != nil {
		zap.L().Error("failed to sync all addr power, it will skipped ", zap.Error(err), zap.Int64("network_id", config.Client.Network.ChainId))
	}
}

// SyncDevWeightStepDay returns a function that synchronizes developer weights for each day within a specified range.
// It takes a pointer to a SyncService as an argument.
func (j *Safejob) SyncDevWeightStepDay() {
	// Return an anonymous function that performs the synchronization.

	// Create a background context for the operation.
	ctx := context.Background()
	// Calculate the start date as the current date minus the data expiration duration, and set it to the end of the day.
	start := carbon.Now().SubDays(j.syncService.GetExpirationData()).EndOfDay()
	// Calculate the end date as yesterday and set it to the end of the day.
	end := carbon.Now().Yesterday().EndOfDay()

	// find latest index

	for end.Gte(start) {
		exist, err := j.syncService.ExistDeveloperWeight(ctx, end.ToShortDateString())
		if err != nil {
			zap.L().Error("SyncDevWeightStepDay", zap.String("date", end.ToShortDateString()))
			return
		}
		if !exist {
			err := j.syncService.SyncDeveloperWeight(ctx, end.ToShortDateString())
			if err != nil {
				return
			}

			break
		}
		end = end.SubDay()
	}
}

// UploadPowerToIPFS returns a function that uploads power data to IPFS.
func (j *Safejob) UploadPowerToIPFS(w3client *data.W3Client) {

	zap.L().Info("backup power start: ", zap.Int64("timestamp", time.Now().Unix()))

	// Iterate over networks and upload power data to IPFS concurrently.
	ctx := context.Background()
	// Upload power data for the current network.
	if err := j.syncService.UploadPowerToIPFS(ctx, config.Client.Network.ChainId, w3client); err != nil {
		zap.L().Error("backup power finished with err:", zap.Error(err))
	}

}
