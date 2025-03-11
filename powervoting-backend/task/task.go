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

	"go.uber.org/zap"

	"powervoting-server/service"
)

type Safejob struct {
	syncService               *service.SyncService
	isRunningSyncEventTask    int32
	isRunningVoteCountingTask int32
}

// RunSyncEventTask Secure synchronization contract event
func (j *Safejob) RunSyncEventTask() {
	if atomic.CompareAndSwapInt32(&j.isRunningSyncEventTask, 0, 1) {
		defer atomic.StoreInt32(&j.isRunningSyncEventTask, 0)

		zap.L().Info("start sync event ")
		SyncEventHandler(j.syncService)
		zap.L().Info("sync event finished, end time:", zap.Int64("end time", time.Now().Unix()))
	} else {
		zap.L().Info("sync event task is running, continue")
	}
}

// RunVoteCountingTask Secure synchronization contract voting count
func (j *Safejob) RunVoteCountingTask() {
	if atomic.CompareAndSwapInt32(&j.isRunningVoteCountingTask, 0, 1) {
		defer atomic.StoreInt32(&j.isRunningVoteCountingTask, 0)

		zap.L().Info("start sync voting count ")
		VotingCountHandler(j.syncService)
		zap.L().Info("sync voting count finished, end time: ", zap.Int64("end time", time.Now().Unix()))
	} else {
		zap.L().Info("sync event task is running, continue")
	}
}
