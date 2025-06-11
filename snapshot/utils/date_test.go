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


package utils_test

import (
	"fmt"
	"time"

	"github.com/golang-module/carbon"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"power-snapshot/utils"
)

var _ = Describe("Date", func() {
	syncEndTime, _ := time.Parse("20060102", "20250501")
	dates := []string{
		"20250423",
		"20250424",
		"20250425",
	}
	Context("CalDateList", func() {
		It("should return the deduplicated list and fill in the missing date strings", func() {
			res := utils.CalDateList(syncEndTime, 60, dates)
			Expect(len(res)).To(Equal(57))
			Expect(res[0]).To(Equal("20250501"))
		})
	})
	Context("CalMissDates", func() {
		It("should return a list of dates from the previous day up to a specified number of days ago", func() {
			res := utils.CalMissDates(dates, 30)
			Expect(len(res)).To(Equal(30))
			Expect(res[0]).To(Equal(carbon.Now().SubDay().EndOfDay().ToShortDateString()))
		})
	})

	Context("AddMonths", func() {
		It("should add the specified number of months to the given date", func() {
			res := utils.AddMonths(syncEndTime, -3)
			Expect(res.Format("yymmdd")).To(Equal(carbon.Parse("20250201").ToStdTime().Format("yymmdd")))
			fmt.Println(res)
		})
	})
})
