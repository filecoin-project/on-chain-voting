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

package constant

import (
	"time"
)

const (
	DataExpiredDuration     = 60
	GithubDataWithinXMonths = 6
	SnapshotBackupSync      = 0
	RetryCount              = 3
	SnapshotBackupSyncd     = 4
	TwoHoursBlockNumber     = 2 * 3600 / 30
	TaskActionActor         = "actor"
	TaskActionMiner         = "miner"
	DeveloperWeightsFile    = "%d_developer_weights_%s"
	DealsFile               = "%d_deals_%s"
	FileSuffix              = ".json"
	SavedHeightDuration     = -DataExpiredDuration * 2880

	MinimumTokenCapacity = 18000
	MinimumTokenNum      = 4
)

var (
	TimeoutWith15s = 15 * time.Second
	TimeoutWith3M  = 15 * time.Second * 12
)

var (
	DEALS     = "deals"
	DEVELOPER = "developer"
)

const (
	GithubCoreOrg = iota
	GithubEcosystemOrg
	GithubUser
)
