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
package service_test

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"power-snapshot/constant"
	models "power-snapshot/internal/model"
)

var _ = Describe("Query", func() {

	Describe("GetAddressPower", func() {
		// BeforeEach(func() {
		// 	mockQuery.EXPECT().
		// 		GetAddressPower(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		// 		Return(&syncPower, nil)
		// })
		It("should return address power", func() {

			power, err := queryService.GetAddressPower(context.Background(), conf.Network.ChainId, "f01", 1)
			Expect(err).To(BeNil())
			Expect(power).To(Equal(&syncPower))
		})
		It("should return error", func() {
			power, err := queryService.GetAddressPower(context.Background(), conf.Network.ChainId, "f01", 61)
			Expect(err.Error()).To(Equal("day count is too long"))
			Expect(power).To(BeNil())
		})

	})

	Describe("GetAddressPowerByDay", func() {
		It("should return address powe of day", func() {
			res, err := queryService.GetAddressPowerByDay(context.Background(), conf.Network.ChainId, "f02", "20250101", time.Now())
			Expect(err).To(BeNil())
			Expect(res).To(Equal(&syncPower))
		})
		Context("power is nil", func() {
			BeforeEach(func() {
				t, _ := time.Parse("20060102", "20240130")
				mockQuery.EXPECT().
					GetAddressPower(gomock.Any(), gomock.Eq(int64(0)), gomock.Any(), gomock.Eq("20250531")).
					Return(nil, nil).AnyTimes()

				mockLotus.EXPECT().
					GetBlockHeader(gomock.Any(), gomock.Eq(int64(0)), gomock.Any()).
					Return(models.BlockHeader{Height: 1, Timestamp: t.Unix()}, nil).AnyTimes()

				mockBase.EXPECT().
					GetDateHeightMap(gomock.Any(), gomock.Eq(int64(0))).
					Return(map[string]int64{"20250501": 1}, nil).AnyTimes()
			})

			When("error not nil", func() {
				It("latest block time is earlier than the sync time", func() {
					res, err := queryService.GetAddressPowerByDay(context.Background(), int64(0), "f02", "20250531", time.Now())
					Expect(err.Error()).To(Equal(constant.ErrorEarlierBlockTime.Error()))
					Expect(res).To(BeNil())
				})

			})
		})
	})
	Describe("GetDataHeight", func() {
		It("should return data height", func() {
			height, err := queryService.GetDataHeight(context.Background(), conf.Network.ChainId, "20250101")
			Expect(err).To(BeNil())
			Expect(height).To(Equal(int64(1)))
		})
	})

	Describe("GetAllAddressPowerByDay", func() {
		It("should return all address power of day", func() {
			res, err := queryService.GetAllAddressPowerByDay(context.Background(), conf.Network.ChainId, "20250101")
			Expect(err).To(BeNil())
			Expect(res).To(Equal(map[string]any{
				"addrPower": []models.SyncPower{
					syncPower,
				},
				"devPower": "testDayPower",
			}))

		})

	})
})
