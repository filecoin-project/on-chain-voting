package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"powervoting-server/config"
	"powervoting-server/constant"
	"powervoting-server/model"
	"powervoting-server/model/api"
	"powervoting-server/service"
	"powervoting-server/utils"
)

var _ service.IProposalService = (*MockProposalService)(nil)
var _ service.IFipService = (*MockFipService)(nil)
var _ service.IVoteService = (*MockVoteService)(nil)

type MockProposalService struct {
	mock.Mock
}

type MockFipService struct {
	mock.Mock
}

// GetFipEditorList implements service.IFipService.
func (m *MockFipService) GetFipEditorList(ctx context.Context, req api.FipEditorListReq) ([]api.FipEditorRep, error) {
	panic("unimplemented")
}

// GetFipProposalList implements service.IFipService.
func (m *MockFipService) GetFipProposalList(ctx context.Context, req api.FipProposalListReq) (*api.CountListRep, error) {
	panic("unimplemented")
}

// AddDraft implements service.IProposalService.
func (m *MockProposalService) AddDraft(ctx context.Context, req *api.AddProposalDraftReq) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// GetDraft implements service.IProposalService.
func (m *MockProposalService) GetDraft(ctx context.Context, req api.AddressReq) (*api.ProposalDraftRep, error) {
	m.Called(ctx, req)
	return nil, nil
}

// ProposalDetail implements service.IProposalService.
func (m *MockProposalService) ProposalDetail(ctx context.Context, req api.ProposalReq) (*api.ProposalRep, error) {
	m.Called(ctx, req)
	return nil, nil
}

// ProposalList implements service.IProposalService.
func (m *MockProposalService) ProposalList(ctx context.Context, req api.ProposalListReq) (*api.CountListRep, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*api.CountListRep), args.Error(1)
}

type MockVoteService struct {
	mock.Mock
}

// VerifyGist implements service.IVoteService.
func (m *MockVoteService) VerifyGist(ctx context.Context, req api.VerifyGistReq) (*model.SigObject, error) {
	panic("unimplemented")
}

// GetFipEditorGistInfo implements service.IVoteService.
func (m *MockVoteService) GetFipEditorGistInfo(ctx context.Context, req api.AddressReq) (*api.FipEditorGistInfoRep, error) {
	panic("unimplemented")
}

// GetCountedVotedList implements service.IVoteService.
func (m *MockVoteService) GetCountedVotedList(ctx context.Context, chainId int64, proposalId int64) ([]api.Voted, error) {
	args := m.Called(ctx, chainId, proposalId)
	return args.Get(0).([]api.Voted), args.Error(1)
}

// 测试初始化
func setupRouter(p *MockProposalService, v *MockVoteService, f *MockFipService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	InitRouters(r, p, v, f)
	return r
}

func TestHealthCheck(t *testing.T) {
	router := setupRouter(nil, nil, nil)
	req, _ := http.NewRequest("GET", constant.PowerVotingApiPrefix+"/health_check", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"code":0`)
}

func TestGetCountedVotesInfo_Success(t *testing.T) {
	proposalService := new(MockProposalService)
	voteService := new(MockVoteService)
	fipService := new(MockFipService)
	router := setupRouter(proposalService, voteService, fipService)

	data := api.Voted{VoterAddress: "addr1", ProposalId: 100}
	voteService.On("GetCountedVotedList",
		mock.Anything, // context
		int64(1),      // chainId
		int64(123),    // proposalId
	).Return([]api.Voted{
		data,
	}, nil)

	req := httptest.NewRequest(
		"GET",
		constant.PowerVotingApiPrefix+"/proposal/votes?chainId=1&proposalId=123",
		nil,
	)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, utils.ObjToString(api.Response{
		Code:    0,
		Message: constant.CodeOKStr,
		Data:    []api.Voted{data},
	}), resp.Body.String())
}

func TestGetCountedVotesInfo_MissingParams(t *testing.T) {
	proposalService := new(MockProposalService)
	voteService := new(MockVoteService)
	fipService := new(MockFipService)
	router := setupRouter(proposalService, voteService, fipService)
	testCases := []struct {
		query       string
		expectedErr string
	}{
		{"chainId=1", constant.CodeParamErrorStr},
		{"proposalId=123", constant.CodeParamErrorStr},
		{"", constant.CodeParamErrorStr},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("GET", constant.PowerVotingApiPrefix+"/proposal/votes?"+tc.query, nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Contains(t, resp.Body.String(), `"code":1000`)
		assert.Contains(t, resp.Body.String(), tc.expectedErr)
	}
}

func TestGetProposalDetail_InvalidChainId(t *testing.T) {
	proposalService := new(MockProposalService)
	voteService := new(MockVoteService)
	fipService := new(MockFipService)
	router := setupRouter(proposalService, voteService, fipService)
	req, _ := http.NewRequest("GET", constant.PowerVotingApiPrefix+"/proposal/details?proposalId=123&chainId=invalid", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Contains(t, resp.Body.String(), `"code":1000`)
	assert.Contains(t, resp.Body.String(), constant.CodeParamErrorStr)
}

func TestGetPower_InvalidAddress(t *testing.T) {
	config.GetDefaultConfig()
	proposalService := new(MockProposalService)
	voteService := new(MockVoteService)
	fipService := new(MockFipService)
	router := setupRouter(proposalService, voteService, fipService)

	req, _ := http.NewRequest("GET", constant.PowerVotingApiPrefix+"/power/getPower?address=invalid&chainId=314159&day=20250105", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Contains(t, resp.Body.String(), `"code":1000`)
	assert.Contains(t, resp.Body.String(), constant.CodeParamErrorStr)
}

func TestPostDraft_Success(t *testing.T) {
	config.GetDefaultConfig()
	proposalService := new(MockProposalService)
	voteService := new(MockVoteService)
	fipService := new(MockFipService)
	router := setupRouter(proposalService, voteService, fipService)
	body := `{"title": "Test Proposal"}`

	proposalService.On("AddDraft",
		mock.Anything,
		&api.AddProposalDraftReq{
			Title: "Test Proposal",
		}).Return(nil)
	req, _ := http.NewRequest("POST", constant.PowerVotingApiPrefix+"/proposal/draft/add", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Contains(t, resp.Body.String(), `"code":0`)
	assert.Contains(t, resp.Body.String(), constant.CodeOKStr)
}

func TestToEthAddr(t *testing.T) {
	config.GetDefaultConfig()
	addr := api.AddressReq{
		Address: "t3qamo6mesb5ppiqqbmsytvadjlbld7k6xji2wkn55hb4w2c5pa4gyyehtxjl6gtxbzs7kcul3z744frid4oqq",
	}

	res, err := addr.ToEthAddr()
	assert.Nil(t, err)
	assert.NotEqual(t, addr.Address, res)
}
