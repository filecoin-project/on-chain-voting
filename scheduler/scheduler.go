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
	"powervoting-server/task"
)

func TaskScheduler() {
	// create a new scheduler
	crontab := cron.New(cron.WithSeconds())
	// task function
	votingCountTask := task.VotingCountHandler
	syncProposalTask := task.SyncProposalHandler
	syncVoteTask := task.SyncVoteHandler

	// 5 minutes
	votingCountSpec := "0 0/5 * * * ?"
	// 1 minutes
	syncSpec := "0/30 * * * * ?"

	_, err := crontab.AddFunc(syncSpec, syncProposalTask)
	if err != nil {
		zap.L().Error("add proposal sync task failed: ", zap.Error(err))
		return
	}
	_, err = crontab.AddFunc(syncSpec, syncVoteTask)
	if err != nil {
		zap.L().Error("add vote sync task failed: ", zap.Error(err))
		return
	}
	_, err = crontab.AddFunc(votingCountSpec, votingCountTask)
	if err != nil {
		zap.L().Error("add voting count task failed: ", zap.Error(err))
		return
	}
	// start
	crontab.Start()

	select {}
}
