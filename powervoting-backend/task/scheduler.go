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

	"powervoting-server/service"
)

// TaskScheduler initializes and starts the task scheduler.
// It creates a new cron scheduler with seconds precision.
// It defines task functions for voting count, proposal synchronization, and vote synchronization.
// It schedules the tasks to run at specific intervals:
//   - Voting count task runs every 5 minutes.
//   - Proposal synchronization task runs every 30 seconds.
//   - Vote synchronization task runs every 30 seconds.
//
// Any error encountered during task scheduling is logged.
func TaskScheduler(syncService *service.SyncService) {
	// create a new scheduler
	crontab := cron.New(cron.WithSeconds())
	// stop the scheduler when the program exits
	defer crontab.Stop()

	job := Safejob{
		syncService: syncService,
	}

	_, err := crontab.AddFunc("0/10 * * * * ?", job.RunSyncEventTask)
	if err != nil {
		zap.L().Error("add proposal sync task failed: ", zap.Error(err))
		return
	}

	_, err = crontab.AddFunc("0 0/5 * * * ?", job.RunVoteCountingTask)
	if err != nil {
		zap.L().Error("add voting count task failed: ", zap.Error(err))
		return
	}

	// start
	crontab.Start()

	select {}
}
