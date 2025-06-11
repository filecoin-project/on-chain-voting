package service_test

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"powervoting-server/model"
	"powervoting-server/model/api"
)

var _ = Describe("Proposal", func() {
	BeforeEach(func() {
		mockProposal.EXPECT().
			GetProposalListWithPagination(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, req api.ProposalListReq) ([]model.ProposalWithVoted, int64, error) {
				if req.Address != "" {
					return []model.ProposalWithVoted{}, 0, fmt.Errorf("not found address: %s", req.Address)
				}

				return []model.ProposalWithVoted{
					{
						ProposalTbl: model.ProposalTbl{
							ProposalId:          1,
							ChainId:             config.Network.ChainId,
							Creator:             "creator1",
							Content:             "content1",
							StartTime:           1,
							EndTime:             2,
							Timestamp:           0,
							Title:               "proposal1",
							Counted:             0,
							BlockNumber:         500_000,
							SnapshotDay:         "20250501",
							SnapshotBlockHeight: 500_000,
							SnapshotTimestamp:   5,
							Percentage: model.Percentage{
								SpPercentage:          2500,
								ClientPercentage:      2500,
								DeveloperPercentage:   2500,
								TokenHolderPercentage: 2500,
							},
						},
					},
				}, 1, nil
			})
		mockProposal.EXPECT().
			GetGitHubNameByCreaters(gomock.Any(), gomock.Any()).
			Return(map[string]model.GiuthubInfo{
				"creator1": {
					GithubName:   "github1",
					GithubAvatar: "avatar1",
				},
				"creator2": {
					GithubName:   "github2",
					GithubAvatar: "avatar2",
				},
			})

		Describe("ProposalList", func() {
			It("should return all proposal list", func() {
				res, err := proposalService.ProposalList(context.Background(), api.ProposalListReq{})
				Expect(err).To(BeNil())
				Expect(res).NotTo(BeNil())
			})
			It("should return error when address is not found", func() {
				_, err := proposalService.ProposalList(context.Background(), api.ProposalListReq{AddressReq: api.AddressReq{
					Address: "address2",
				}})
				Expect(err).NotTo(BeNil())
			})
		})

		Describe("ProposalDetail", func() {
			BeforeEach(func() {
				mockProposal.EXPECT().
					GetProposalById(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req api.ProposalReq) (*model.ProposalTbl, error) {
						if req.ProposalId == 2 {
							return nil, fmt.Errorf("not found proposal id: %d", req.ProposalId)
						}

						return &model.ProposalTbl{
							ProposalId: 1,
							Creator:    "creator1",
						}, nil
					})
			})
			It("should return proposal detail", func() {
				res, err := proposalService.ProposalDetail(context.Background(), api.ProposalReq{ProposalId: 1})
				Expect(err).To(BeNil())
				Expect(res).NotTo(BeNil())
				Expect(res.ProposalId).To(Equal(int64(1)))
				Expect(res.GithubName).To(Equal("github1"))
			})

			It("should return error when proposal id is not found", func() {
				_, err := proposalService.ProposalDetail(context.Background(), api.ProposalReq{ProposalId: 2})
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
