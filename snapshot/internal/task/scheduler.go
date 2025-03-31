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
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"power-snapshot/internal/service"
)

type JobInfo struct {
	Spec string
	Job  func()
}

// TaskScheduler initializes and starts the task scheduler.
// It creates a new cron scheduler with seconds precision.
// Any error encountered during task scheduling is logged.

func TaskScheduler(syncService *service.SyncService) {
	// create a new scheduler
	crontab := cron.New(cron.WithSeconds())
	defer crontab.Stop()
	// task function
	// 5 minutes
	job := Safejob{
		syncService: syncService,
	}

	_, err := crontab.AddFunc("0 0 1 * * ?", job.RunSyncPower)
	if err != nil {
		zap.L().Error("failed to add RunSyncPower task to scheduler", zap.Error(err))
	}

	_, err = crontab.AddFunc("0 0 0/1 * * ?", job.RunSyncDevWeightStepDay)
	if err != nil {
		zap.L().Error("failed to add RunSyncDevWeightStopDay task to scheduler", zap.Error(err))
	}

	_, err = crontab.AddFunc("0 0/10 * * * ?", job.RunUploadPowerToIPFS)
	if err != nil {
		zap.L().Error("failed to add RunUploadPowerToIPFS task to scheduler", zap.Error(err))
	}
	// start
	crontab.Start()

	select {}
}
