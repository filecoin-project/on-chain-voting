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
	"go.uber.org/zap"
	"power-snapshot/config"
	"power-snapshot/constant"
	"power-snapshot/internal/service"

	"github.com/golang-module/carbon"
)

func SyncPower(service *service.SyncService) func() {
	return func() {
		ctx := context.Background()
		for _, network := range config.Client.Network {
			// sync date height
			err := service.SyncDateHeight(ctx, network.Id)
			if err != nil {
				zap.L().Error("failed to sync date height, it will skipped ", zap.Error(err), zap.Int64("network_id", network.Id))
				continue
			}
			err = service.SyncAllAddrPower(ctx, network.Id)
			if err != nil {
				zap.L().Error("failed to sync all addr power, it will skipped ", zap.Error(err), zap.Int64("network_id", network.Id))
				continue
			}
		}
	}
}

func SyncDevWeightStepDay(service *service.SyncService) func() {
	return func() {
		ctx := context.Background()
		start := carbon.Now().SubDays(constant.DataExpiredDuration).EndOfDay()
		end := carbon.Now().Yesterday().EndOfDay()

		// find latest index
		for l := 0; l < len(config.Client.Github.Token)*2; {
			for i := start; i.Timestamp() <= end.Timestamp(); i = i.AddDay() {
				exist, err := service.ExistDeveloperWeight(ctx, i.ToShortDateString())
				if err != nil {
					zap.L().Error("SyncDevWeightStepDay", zap.String("date", i.ToShortDateString()))
					return
				}
				if !exist {
					err := service.SyncDeveloperWeight(ctx, i.ToShortDateString())
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

func SyncDelegateEvent(service *service.SyncService) func() {
	return func() {
		ctx := context.Background()
		for _, network := range config.Client.Network {
			// sync date height
			err := service.SyncDelegateEvent(ctx, network.Id)
			if err != nil {
				zap.L().Error("failed to sync delegate event, it will skipped ", zap.Error(err), zap.Int64("network_id", network.Id))
				continue
			}
		}
	}
}
