package service_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	configs "powervoting-server/config"
	"powervoting-server/mocks"
	"powervoting-server/model"
	"powervoting-server/service"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var (
	mockFip         *mocks.MockFipRepo
	mockLotus       *mocks.MockLotusRepo
	mockVote        *mocks.MockVoteRepo
	mockSync        *mocks.MockSyncRepo
	mockProposal    *mocks.MockProposalRepo
	mockCtrl        = gomock.NewController(GinkgoT())
	fipService      service.IFipService
	proposalService service.IProposalService
	rpcService      service.RpcService
	syncService     service.ISyncService
	config          configs.Config

	proposalListData = []model.FipProposalVoted{
		{
			FipProposalTbl: model.FipProposalTbl{
				ProposalId:       1,
				ChainId:          config.Network.ChainId,
				ProposalType:     1,
				Creator:          "f01",
				CandidateAddress: "f00",
				CandidateInfo:    "info",
				Timestamp:        time.Now().Unix(),
				BlockNumber:      500_000,
				Status:           0,
			},

			Voters: "f01",
		},
	}
)

var _ = BeforeSuite(func() {
	config = configs.Config{
		Network: configs.Network{
			ChainId: 314159,
		},
	}
	mockFip = mocks.NewMockFipRepo(mockCtrl)
	mockProposal = mocks.NewMockProposalRepo(mockCtrl)
	mockVote = mocks.NewMockVoteRepo(mockCtrl)
	mockLotus = mocks.NewMockLotusRepo(mockCtrl)
	mockSync = mocks.NewMockSyncRepo(mockCtrl)
	fipService = service.NewFipService(mockFip)
	proposalService = service.NewProposalService(mockProposal)
	syncService = service.NewSyncService(mockSync, mockVote, mockProposal, mockFip, mockLotus)
	rpcService = *service.NewRpcService(mockVote, mockSync,mockLotus)
})
