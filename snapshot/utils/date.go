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

package utils

import (
	"time"

	"github.com/golang-module/carbon"
	"github.com/samber/lo"

	"power-snapshot/constant"
)

/**
 * @Description: Refactora calculates a list of dates that need to be synchronized
 * @param syncEndTime The end time of the synchronization
 * @param syncCountedDays The number of days to be synchronized.
 * @param dates
 * @return []string
 */
func CalDateList(syncEndTime time.Time, syncCountedDays int, dates []string) []string {
	base := carbon.FromStdTime(syncEndTime)
	allDatesList := make([]string, 0, syncCountedDays)
	// Calculate the number of days to be synchronized
	for range syncCountedDays {
		allDatesList = append(allDatesList, base.ToShortDateString())
		base = base.SubDay()
	}

	diff, _ := lo.Difference(allDatesList, dates)
	return diff
}

func CalMissDates(dates []string) []string {
	return CalDateList(time.Now().Add(-(24 * time.Hour)), constant.DataExpiredDuration, dates)
}
