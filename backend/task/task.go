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
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"powervoting-server/service"
)

type Safejob struct {
	syncService               *service.SyncService
	isRunningSyncEventTask    int32
	isRunningVoteCountingTask int32
}

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

// RunSyncEventTask Secure synchronization contract event
func (j *Safejob) RunSyncEventTask() {
	if atomic.CompareAndSwapInt32(&j.isRunningSyncEventTask, 0, 1) {
		defer atomic.StoreInt32(&j.isRunningSyncEventTask, 0)

		zap.L().Debug("start sync event ")
		SyncEventHandler(j.syncService)
		zap.L().Debug("sync event finished, end time:", zap.Int64("end time", time.Now().Unix()))
	} else {
		zap.L().Warn("sync event task is running, continue")
	}
}

// RunVoteCountingTask Secure synchronization contract voting count
func (j *Safejob) RunVoteCountingTask() {
	if atomic.CompareAndSwapInt32(&j.isRunningVoteCountingTask, 0, 1) {
		defer atomic.StoreInt32(&j.isRunningVoteCountingTask, 0)

		zap.L().Debug("start sync voting count ")
		VotingCountHandler(j.syncService)
		zap.L().Debug("sync voting count finished, end time: ", zap.Int64("end time", time.Now().Unix()))
	} else {
		zap.L().Warn("sync voting count task is running, continue")
	}
}
