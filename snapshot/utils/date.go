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

	"power-snapshot/config"
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

func CalMissDates(dates []string, durationDays int) []string {
	var startTime = time.Now().Add(-24 * time.Hour)

	if config.Client.SyncStartDate != "" {
		carbonDate, err := time.Parse("20060102", config.Client.SyncStartDate)
		if err == nil {
			elapsed := time.Since(carbonDate)
			elapsedDays := int(elapsed.Hours() / 24)

			if elapsedDays < durationDays {
				durationDays = elapsedDays
			}
		}
	} else {
		startTime = time.Now().Add(-24 * time.Hour)
	}

	return CalDateList(startTime, durationDays, dates)
}

// addMonths adds a specified number of months to a given date.
func AddMonths(input time.Time, months int) time.Time {
	date := time.Date(input.Year(), input.Month(), 1, 0, 0, 0, 0, input.Location())
	date = date.AddDate(0, months, 0)

	lastDay := getLastDayOfMonth(date.Year(), int(date.Month()))
	date = time.Date(date.Year(), date.Month(), lastDay, 0, 0, 0, 0, input.Location())

	if input.Day() < lastDay {
		date = time.Date(date.Year(), date.Month(), input.Day(), 0, 0, 0, 0, input.Location())
	}

	return date
}

// getLastDayOfMonth calculates the last day of the specified month in the given year.
func getLastDayOfMonth(year, month int) int {
	lastDay := 31
	nextMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)
	lastDay = nextMonth.Day()
	return lastDay
}
