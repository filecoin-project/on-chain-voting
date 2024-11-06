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

package routers

import (
	"powervoting-server/api"
	"powervoting-server/response"

	"github.com/gin-gonic/gin"
)

// InitRouters initializes the routers for the power voting API endpoints.
// It defines routes for health check, proposal result, and proposal history.
// The health check route returns a success response.
// The proposal result route is mapped to the VoteResult handler function.
// The proposal history route is mapped to the VoteHistory handler function.
func InitRouters(r *gin.Engine) {
	powerVotingRouter := r.Group("/power_voting/api/")
	powerVotingRouter.GET("/health_check", func(c *gin.Context) {
		response.Success(c)
		return
	})

	powerVotingRouter.GET("/proposal/result", api.VoteResult)
	powerVotingRouter.GET("/proposal/history", api.VoteHistory)
	powerVotingRouter.GET("/proposal/list", api.ProposalList)
	powerVotingRouter.POST("/proposal/add", api.AddProposal)

	powerVotingRouter.POST("/proposal/draft/add", api.AddDraft)
	powerVotingRouter.GET("/proposal/draft/get", api.GetDraft)

	powerVotingRouter.POST("/w3storage/upload", api.W3Upload)

	powerVotingRouter.GET("/filecoin/height", api.GetHeight)

	powerVotingRouter.GET("/power/getPower", api.GetPower)
}
