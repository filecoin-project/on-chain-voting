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

package scheduler

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type JobInfo struct {
	Spec string
	Job  func()
}

// TaskScheduler initializes and starts the task scheduler.
// It creates a new cron scheduler with seconds precision.
// Any error encountered during task scheduling is logged.

func TaskScheduler(jobs []JobInfo) {
	// create a new scheduler
	crontab := cron.New(cron.WithSeconds())
	// task function
	// 5 minutes
	for _, job := range jobs {
		_, err := crontab.AddFunc(job.Spec, job.Job)
		if err != nil {
			zap.L().Error("add sync task failed: ", zap.Error(err))
			return
		}
	}
	// start
	crontab.Start()
}
