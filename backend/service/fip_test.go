package service_test

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"powervoting-server/model"
	"powervoting-server/model/api"
)

var _ = Describe("Fip", func() {

	BeforeEach(func() {
		mockFip.EXPECT().
			GetFipProposalListWithPagination(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, req api.FipProposalListReq) ([]model.FipProposalVoted, int64, error) {
				if req.ProposalType == 1 {
					return []model.FipProposalVoted{}, 0, nil
				} else {
					return proposalListData, 1, nil
				}
			}).AnyTimes()

		mockFip.EXPECT().
			GetFipEditorCount(gomock.Any(), gomock.Eq(config.Network.ChainId)).
			Return(int64(1), nil).AnyTimes()

		mockFip.EXPECT().
			GetFipProposalVoteCount(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(1), nil).AnyTimes()

		mockFip.EXPECT().
			GetValidFipEditorList(gomock.Any(), gomock.Any()).
			Return([]model.FipEditorTbl{
				{
					BaseField: model.BaseField{
						UpdatedAt: time.Unix(0, 0),
					},
					ChainId:  config.Network.ChainId,
					Editor:   "editor1",
					IsRemove: 0,
				},
			}, nil).AnyTimes()
	})
	Describe("GetFipProposalList", func() {
		It("should return a list of fip proposals", func() {
			fipList, err := fipService.GetFipProposalList(context.Background(), api.FipProposalListReq{
				ChainId: config.Network.ChainId,
			})
			Expect(err).To(BeNil())
			Expect(fipList.List).ToNot(BeNil())
		})

		It("should return list of nil", func() {
			fipList, err := fipService.GetFipProposalList(context.Background(), api.FipProposalListReq{
				ChainId:      config.Network.ChainId,
				ProposalType: 1,
			})
			Expect(err).To(BeNil())
			Expect(fipList.List).To(BeNil())
		})
	})

	Describe("GetFipEditorList", func() {
		It("should return a list of fip editors", func() {
			fipList, err := fipService.GetFipEditorList(context.Background(), api.FipEditorListReq{
				ChainId: config.Network.ChainId,
			})
			Expect(err).To(BeNil())
			Expect(fipList).To(Equal([]api.FipEditorRep{
				{
					Editor:  "editor1",
					ChainId: config.Network.ChainId,
				},
			}))
		})
	})
})
