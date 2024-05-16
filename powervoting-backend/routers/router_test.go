package routers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"powervoting-server/model"
	"powervoting-server/response"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	powerVotingRouter := r.Group("/power_voting/api/")
	powerVotingRouter.GET("/health_check", func(c *gin.Context) {
		response.Success(c)
		return
	})

	powerVotingRouter.GET("/proposal/result", func(context *gin.Context) {
		var result []model.VoteResult
		result = append(result, model.VoteResult{
			Id:         1,
			ProposalId: 1,
			OptionId:   1,
			Votes:      1,
			Network:    1,
		})
		context.JSON(200, result)
	})
	powerVotingRouter.GET("/proposal/history", func(context *gin.Context) {
		var result []model.VoteCompleteHistory
		result = append(result, model.VoteCompleteHistory{
			Id:                    1,
			ProposalId:            1,
			Network:               1,
			TotalSpPower:          "1",
			TotalClientPower:      "1",
			TotalTokenHolderPower: "1",
			TotalDeveloperPower:   "1",
			VotePowers: []model.VotePower{
				{
					Id:                      1,
					HistoryId:               1,
					Address:                 "test",
					OptionId:                1,
					Votes:                   1,
					SpPower:                 "1",
					ClientPower:             "1",
					TokenHolderPower:        "1",
					DeveloperPower:          "1",
					SpPowerPercent:          1,
					ClientPowerPercent:      1,
					TokenHolderPowerPercent: 1,
					DeveloperPowerPercent:   1,
					PowerBlockHeight:        1,
				},
			},
		})
		context.JSON(200, result)
	})

	return r
}

func TestHealthCheck(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/power_voting/api/health_check", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestVoteResult(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/power_voting/api/proposal/result", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	voteResult := []model.VoteResult{
		{
			Id:         1,
			ProposalId: 1,
			OptionId:   1,
			Votes:      1,
			Network:    1,
		},
	}
	voteJson, _ := json.Marshal(voteResult)
	assert.Equal(t, string(voteJson), w.Body.String())
}

func TestVoteHistory(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/power_voting/api/proposal/history", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	historyResult := []model.VoteCompleteHistory{
		{
			Id:                    1,
			ProposalId:            1,
			Network:               1,
			TotalSpPower:          "1",
			TotalClientPower:      "1",
			TotalTokenHolderPower: "1",
			TotalDeveloperPower:   "1",
			VotePowers: []model.VotePower{
				{
					Id:                      1,
					HistoryId:               1,
					Address:                 "test",
					OptionId:                1,
					Votes:                   1,
					SpPower:                 "1",
					ClientPower:             "1",
					TokenHolderPower:        "1",
					DeveloperPower:          "1",
					SpPowerPercent:          1,
					ClientPowerPercent:      1,
					TokenHolderPowerPercent: 1,
					DeveloperPowerPercent:   1,
					PowerBlockHeight:        1,
				},
			},
		},
	}
	historyJson, _ := json.Marshal(historyResult)
	assert.Equal(t, string(historyJson), w.Body.String())
}
