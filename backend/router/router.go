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
	"reflect"

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
func InitRouters(r *gin.Engine, proposalService service.IProposalService, voteService service.IVoteService, fipService service.IFipService) {

	proposalHandler := api.NewProposalHandler(proposalService)
	voteHandler := api.NewVoteHandler(voteService)
	fipHandler := api.NewFipHandle(fipService)
	powerVotingRouter := r.Group(constant.PowerVotingApiPrefix)
	r.GET(constant.PowerVotingApiPrefix+"/health_check", func(c *gin.Context) {
		api.Success(c)
	})

	proposalRouter(powerVotingRouter, proposalHandler, voteHandler)
	powerRouter(powerVotingRouter)
	fipEditor(powerVotingRouter, fipHandler, voteHandler)
}

// proposalRouter defines routes related to proposal management.
func proposalRouter(rg *gin.RouterGroup, ph *api.ProposalHandler, vh *api.VoteHandler) {
	rg.GET("/proposal/votes", wrap(vh.GetCountedVotesInfo)) // Get counted votes for a proposal

	rg.GET("/proposal/list", wrap(ph.GetProposalList))        // Get a list of proposals
	rg.GET("/proposal/details", wrap(ph.GetProposalDetail))   // Get details of a specific proposal
	rg.POST("/proposal/draft/add", wrap(ph.PostDraft))        // Add a new proposal draft
	rg.DELETE("/proposal/draft/delete", wrap(ph.DeleteDraft)) // Delete a specific proposal draft
	rg.GET("/proposal/draft/get", wrap(ph.GetDraft))          // Get a specific proposal draft
}

// powerRouter defines routes related to power distribution and management.
func powerRouter(rg *gin.RouterGroup) {
	rg.GET("/power/getPower", wrap(api.GetAddressPower)) // Get power distribution for a specific address
}

// fipEditor sets up routes for handling FIP (Federated Identity Proposal) related operations.
func fipEditor(rg *gin.RouterGroup, fh *api.FipHandle, vh *api.VoteHandler) {
	rg.GET("/fipProposal/list", wrap(fh.GetFipProposalList)) // Get a list of fipProposals
	rg.GET("/fipEditor/list", wrap(fh.GetFipEditorList))     // Get a list of fipEditors
	rg.GET("/voter/info", wrap(vh.GetFipEditorGistInfo))     // Get FIP editor gist info
	rg.GET("/fipEditor/checkGist", wrap(vh.VerifyGistValid))
	// The wrap function is used to handle the request and response, passing the fh.GetFipEditorList handler.
}

// wrap is a utility function to wrap handlers with additional context and validation.
func wrap(h func(c *constant.Context)) gin.HandlerFunc {
	validate := validator.New()
	validate.RegisterValidation("is-integer", func(fl validator.FieldLevel) bool {
		field := fl.Field()
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return true
		case reflect.Float32, reflect.Float64:
			floatValue := field.Float()
			return floatValue == float64(int64(floatValue))
		default:
			return false
		}
	})

	return func(c *gin.Context) {
		h(&constant.Context{Context: c, Validate: validate}) // Pass context and validator to the handler
	}
}
