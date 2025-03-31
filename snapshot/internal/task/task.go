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

	"power-snapshot/internal/data"
	"power-snapshot/internal/service"
)

type Safejob struct {
	syncService                   *service.SyncService
	isRunningSyncPowerTask        int32
	isRunningDevWeightStepDayTask int32
	isRunningUploadIPFSTask       int32
}

func (j *Safejob) RunSyncPower() {
	if atomic.CompareAndSwapInt32(&j.isRunningSyncPowerTask, 0, 1) {
		defer atomic.StoreInt32(&j.isRunningSyncPowerTask, 0)

		zap.L().Info("start sync power ")
		j.SyncPower()
		zap.L().Info("sync power finished, end time:", zap.Int64("end time", time.Now().Unix()))
	} else {
		zap.L().Info("sync power task is running, continue")
	}
}

func (j *Safejob) RunSyncDevWeightStepDay() {
	if atomic.CompareAndSwapInt32(&j.isRunningDevWeightStepDayTask, 0, 1) {
		defer atomic.StoreInt32(&j.isRunningDevWeightStepDayTask, 0)

		zap.L().Info("start sync dev weight stop day")
		j.SyncDevWeightStepDay()
		zap.L().Info("sync weight stop day finished, end time:", zap.Int64("end time", time.Now().Unix()))
	} else {
		zap.L().Info("sync weight stop day task is running, continue")
	}
}

func (j *Safejob) RunUploadPowerToIPFS() {
	if atomic.CompareAndSwapInt32(&j.isRunningUploadIPFSTask, 0, 1) {
		defer atomic.StoreInt32(&j.isRunningUploadIPFSTask, 0)

		zap.L().Info("start upload address power to ipfs")
		j.UploadPowerToIPFS( data.NewW3Client())
		zap.L().Info("sync upload address power to ipfs finished, end time:", zap.Int64("end time", time.Now().Unix()))
	} else {
		zap.L().Info("sync upload address power to ipfs task is running, continue")
	}
}
