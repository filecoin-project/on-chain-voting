package task_test

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"powervoting-server/api/rpc"
	pb "powervoting-server/api/rpc/proto"
	"powervoting-server/mocks"
	"powervoting-server/model"
)

var (
	mockISyncService   *mocks.MockISyncService
	mockContractRepo   *mocks.MockIContractRepo
	mockSnapshotClient *mocks.MockSnapshotClient
	mockDrand          *mocks.MockIEncrypted
)

func TestTask(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Task Suite")
}

var _ = BeforeSuite(func() {
	ctrl := gomock.NewController(GinkgoT())
	mockISyncService = mocks.NewMockISyncService(ctrl)
	mockContractRepo = mocks.NewMockIContractRepo(ctrl)
	mockSnapshotClient = mocks.NewMockSnapshotClient(ctrl)
	mockDrand = mocks.NewMockIEncrypted(ctrl)
	mockISyncServiceFunc()
	mockContractRepoFunc()
	mockQueryGrpcFunc()
	mockDrandFunc()
	rpc.SetSnapshotClient(mockSnapshotClient)
})

func mockContractRepoFunc() {
	mockContractRepo.EXPECT().
		GetBlockTime(gomock.Eq(int64(1))).
		Return(int64(1), nil).AnyTimes()

}

func mockQueryGrpcFunc() {
	allAddrPower := model.SnapshotAllPower{
		AddrPower: []model.AddrPower{
			{
				Address:          "0x1",
				DateStr:          "20250101",
				GithubAccount:    "test1",
				DeveloperPower:   big.NewInt(50),
				SpPower:          big.NewInt(0),
				ClientPower:      big.NewInt(100),
				TokenHolderPower: big.NewInt(150),
			},
			{
				Address:          "0x2",
				DateStr:          "20250101",
				GithubAccount:    "test2",
				DeveloperPower:   big.NewInt(300),
				SpPower:          big.NewInt(300),
				ClientPower:      big.NewInt(100),
				TokenHolderPower: big.NewInt(200),
			},
			{
				Address:          "0x3",
				DateStr:          "20250101",
				GithubAccount:    "test3",
				DeveloperPower:   big.NewInt(50),
				SpPower:          big.NewInt(100),
				ClientPower:      big.NewInt(200),
				TokenHolderPower: big.NewInt(50),
			},
		},
		DevPower: nil,
	}

	allAddrPowerJsonStr, _ := json.Marshal(allAddrPower)
	mockSnapshotClient.EXPECT().
		GetAllAddrPowerByDay(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(
			&pb.GetAllAddrPowerByDayResponse{
				Info:  string(allAddrPowerJsonStr),
				Day:   "20250501",
				NetId: 314159,
			}, nil,
		).AnyTimes()

}

func mockISyncServiceFunc() {
	mockISyncService.EXPECT().
		CreateFipEditor(gomock.Any(), gomock.Any()).
		Return(nil).AnyTimes()

	mockISyncService.EXPECT().
		CreateSyncEventInfo(gomock.Any(), gomock.Any()).
		Return(nil).AnyTimes()

	mockISyncService.EXPECT().
		GetUncountedVotedList(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]model.VoteTbl{
			{
				ProposalId:    1,
				Address:       "0x1",
				VoteEncrypted: "approve",
			},
			{
				ProposalId:    1,
				Address:       "0x2",
				VoteEncrypted: "reject",
			},
			{
				ProposalId:    1,
				Address:       "0x3",
				VoteEncrypted: "approve",
			},
		}, nil).AnyTimes()

	mockISyncService.EXPECT().
		UpdateProposal(gomock.Any(), gomock.Any()).
		DoAndReturn(func (ctx context.Context,p *model.ProposalTbl) error {
			Expect(p.Counted).To(Equal(1))
			Expect(p.ProposalResult.ApprovePercentage).To(Equal(43.75))
			Expect(p.ProposalResult.RejectPercentage).To(Equal(56.25))
			return nil
		}).AnyTimes()

	mockISyncService.EXPECT().
		BatchUpdateVotes(gomock.Any(), gomock.Any()).
		Return(nil).AnyTimes()
}

func mockDrandFunc() {
	mockDrand.EXPECT().
		DecodeVoteResult(gomock.Any()).DoAndReturn(func(v string) (string, error) {
		if v == "approve" {
			return "approve", nil
		}
		return "reject", nil
	}).AnyTimes()

}
