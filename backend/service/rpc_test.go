package service_test

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"powervoting-server/model"
)

var _ = Describe("Rpc", func() {
	BeforeEach(func() {
		mockVote.EXPECT().
			GetAllVoterAddresss(gomock.Any(), gomock.Any()).
			Return([]model.VoterInfoTbl{
				{
					Address: "vote1",
				},
				{
					Address: "vote2",
				},
			}, nil)

	})
	Describe("GetAllVoterAddresss", func() {
		It("should return all voter address", func() {
			voter, err := rpcService.GetAllVoterAddresss(context.Background(), config.Network.ChainId)
			Expect(err).To(BeNil())
			Expect(voter).To(Equal([]string{
				"vote1",
				"vote2",
			}))
		})
	})

	Describe("GetVoterInfoByAddress", func() {
		var miners []string
		JustBeforeEach(func() {
			mockVote.EXPECT().
				GetVoterInfoByAddress(gomock.Any(), gomock.Any()).
				Return(&model.VoterInfoTbl{
					Address:  "vote1",
					ChainId:  config.Network.ChainId,
					MinerIds: miners,
				}, nil)
			mockLotus.EXPECT().GetValidMinerIds(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, chainId string, minerIds []string) ([]string, error) {
					if len(minerIds) == 0 {
						return nil, fmt.Errorf("no miner id")
					}

					for _, minerId := range minerIds {
						if minerId == "miner1" {
							return []string{"miner1"}, nil
						}
					}

					return nil, fmt.Errorf("miner is not valid")
				})

			When("has valid miner id", func() {
				BeforeEach(func() {
					miners = []string{"miner1"}
				})
				It("should return voter info", func() {
					voter, err := rpcService.GetVoterInfoByAddress(context.Background(), "vote1")
					Expect(err).To(BeNil())
					Expect(voter).To(Equal(model.VoterInfoTbl{
						Address:  "vote1",
						ChainId:  config.Network.ChainId,
						MinerIds: []string{"miner1"},
					}))
				})
			})

			When("no valid miner id",func ()  {
				BeforeEach(func() {
					miners = []string{"miner2"}
				})

				It("should return error", func() {
					_, err := rpcService.GetVoterInfoByAddress(context.Background(), "vote1")
					Expect(err).ToNot(BeNil())
				})
				
			})

			It("should return error when no miner id", func() {
				_, err := rpcService.GetVoterInfoByAddress(context.Background(), "vote2")
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
