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
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	models "power-snapshot/internal/model"
	"power-snapshot/internal/service"
	mocks "power-snapshot/mock"
)

var _ = Describe("Lotus", func() {
	var (
		lotusService *service.LotusService
		mockLotus    *mocks.MockLotusRepo
		mockCtrl     *gomock.Controller
	)
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockLotus = mocks.NewMockLotusRepo(mockCtrl)
		lotusService = service.NewLotusService(mockLotus)
	})
	Describe("GetTipSetByHeight", func() {
		It("should return tipset", func() {
			mockLotus.EXPECT().
				GetTipSetByHeight(gomock.Any(), conf.Network.ChainId, gomock.Any()).
				Return(tipset, nil)

			res, err := lotusService.GetTipSetByHeight(context.Background(), conf.Network.ChainId, 1)
			Expect(err).To(BeNil())
			Expect(res).NotTo(BeNil())
			Expect(res).To(Equal(tipset))
		})
	})
	Describe("GetAddrBalanceBySpecialHeight", func() {
		var (
			amount string
			err    error
		)
		JustBeforeEach(func() {
			mockLotus.EXPECT().
				GetAddrBalanceBySpecialHeight(gomock.Any(), gomock.Any(), conf.Network.ChainId, gomock.Any()).
				Return(amount, err)
		})
		When("get balance ", func() {
			BeforeEach(func() {
				amount = "100"
			})
			It("should return addr balance", func() {
				res, err := lotusService.GetAddrBalanceBySpecialHeight(context.Background(), "f01", conf.Network.ChainId, 1)
				Expect(err).To(BeNil())
				Expect(res).To(Equal("100"))
			})
		})
		When("get balance error", func() {
			BeforeEach(func() {
				err = errors.New("error")
			})
			It("should return error", func() {
				_, err := lotusService.GetAddrBalanceBySpecialHeight(context.Background(), "f01", conf.Network.ChainId, 1)
				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("GetMinerPowerByHeight", func() {
		It("should return miner power", func() {
			d := models.LotusMinerPower{
				MinerPower: models.MinerPower{
					RawBytePower:    "100",
					QualityAdjPower: "100",
				},
				TotalPower: models.TotalPower{
					RawBytePower:    "100",
					QualityAdjPower: "100",
				},
				HasMinPower: true,
			}
			mockLotus.EXPECT().
				GetMinerPowerByHeight(gomock.Any(), conf.Network.ChainId, gomock.Any(), gomock.Any()).
				Return(d, nil)

			res, err := lotusService.GetMinerPowerByHeight(context.Background(), conf.Network.ChainId, "f01", tipset)
			Expect(err).To(BeNil())
			Expect(res).NotTo(BeNil())
			Expect(res).To(Equal(&d))
		})
	})

	Describe("GetNewestHeight", func() {
		It("should return newest height", func() {
			mockLotus.EXPECT().
				GetNewestHeight(gomock.Any(), conf.Network.ChainId).
				Return(int64(1), nil)

			res, err := lotusService.GetNewestHeight(context.Background(), conf.Network.ChainId)
			Expect(err).To(BeNil())
			Expect(res).To(Equal(int64(1)))

		})
	})
	Describe("GetBlockHeader", func() {
		It("should return block header", func() {
			blockHeader := models.BlockHeader{
				Height:    1,
				Timestamp: 1,
			}
			mockLotus.EXPECT().
				GetBlockHeader(gomock.Any(), conf.Network.ChainId, gomock.Any()).
				Return(blockHeader, nil)
			res, err := lotusService.GetBlockHeader(context.Background(), conf.Network.ChainId, 1)
			Expect(err).To(BeNil())
			Expect(res).NotTo(BeNil())
			Expect(res).To(Equal(&blockHeader))
		})
	})

	Describe("GetWalletBalanceByHeight", func() {
		It("should return wallet balance", func() {
			mockLotus.EXPECT().
				GetWalletBalanceByHeight(gomock.Any(), "f01", conf.Network.ChainId, gomock.Any()).
				Return("100", nil)
			res, err := lotusService.GetWalletBalanceByHeight(context.Background(), "f01", conf.Network.ChainId, 1)
			Expect(err).To(BeNil())
			Expect(res).To(Equal("100"))
		})
	})
})
