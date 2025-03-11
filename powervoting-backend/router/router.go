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

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"powervoting-server/api"
	"powervoting-server/constant"
	"powervoting-server/service"
)

// InitRouters initializes the routers for the power voting API endpoints.
// It defines routes for health check, proposal result, and proposal history.
// The health check route returns a success response.
// The proposal result route is mapped to the VoteResult handler function.
// The proposal history route is mapped to the VoteHistory handler function.
func InitRouters(r *gin.Engine, proposalService service.IProposalService, voteService service.IVoteService) {

	proposalHandler := api.NewProposalHandler(proposalService)
	voteHandle := api.NewVoteHandler(voteService)

	powerVotingRouter := r.Group(constant.PowerVotingApiPrefix)
	r.GET(constant.PowerVotingApiPrefix+"/health_check", func(c *gin.Context) {
		api.Success(c)
	})

	proposalRouter(powerVotingRouter, proposalHandler, voteHandle)
	powerRouter(powerVotingRouter)
}

// proposalRouter defines routes related to proposal management.
func proposalRouter(rg *gin.RouterGroup, ph *api.ProposalHandler, vh *api.VoteHandler) {
	rg.GET("/proposal/votes", wrap(vh.GetCountedVotesInfo)) // Get counted votes for a proposal
	rg.GET("/proposal/list", wrap(ph.GetProposalList))      // Get a list of proposals
	rg.GET("/proposal/details", wrap(ph.GetProposalDetail)) // Get details of a specific proposal
	rg.POST("/proposal/draft/add", wrap(ph.PostDraft))      // Add a new proposal draft
	rg.GET("/proposal/draft/get", wrap(ph.GetDraft))        // Get a specific proposal draft
}

// powerRouter defines routes related to power distribution and management.
func powerRouter(rg *gin.RouterGroup) {
	rg.GET("/power/getPower", wrap(api.GetAddressPower)) // Get power distribution for a specific address
}

// wrap is a utility function to wrap handlers with additional context and validation.
func wrap(h func(c *constant.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(&constant.Context{Context: c, Validate: validator.New()}) // Pass context and validator to the handler
	}
}
