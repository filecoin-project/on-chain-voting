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

package api

import (
	"powervoting-server/constant"
	"powervoting-server/model/api"
	"powervoting-server/service"
)

type ProposalHandler struct {
	proposqlService service.IProposalService
}

func NewProposalHandler(ps service.IProposalService) *ProposalHandler {
	return &ProposalHandler{
		proposqlService: ps,
	}
}

// AddDraft function handles an HTTP request to add a draft proposal to the database.
func (p *ProposalHandler) PostDraft(c *constant.Context) {
	var draft api.AddProposalDraftReq
	if err := c.BindAndValidate(&draft); err != nil {
		ParamError(c.Context)
		return
	}

	if err := p.proposqlService.AddDraft(c.Request.Context(), &draft); err != nil {
		Error(c.Context, err)
		return
	}

	Success(c.Context)
}

// GetDraft function handles an HTTP request to retrieve a draft proposal from the database.
func (p *ProposalHandler) GetDraft(c *constant.Context) {
	var req api.GetDraftReq
	if err := c.BindAndValidate(&req); err != nil {
		ParamError(c.Context)
		return
	}

	res, err := p.proposqlService.GetDraft(c.Context, req)
	if err != nil {
		Error(c.Context, err)
		return
	}

	SuccessWithData(c.Context, res)
}

// GetProposalDetails is a function that handles the request to get details of a proposal.
func (p *ProposalHandler) GetProposalDetail(c *constant.Context) {
	var req api.ProposalReq
	if err := c.BindAndValidate(&req); err != nil {
		ParamError(c.Context)
		return
	}

	res, err := p.proposqlService.ProposalDetail(c.Context, req)
	if err != nil {
		Error(c.Context, err)
		return
	}

	SuccessWithData(c.Context, res)
}

// GetProposalList is a function that handles the request to get a list of proposals.
func (p *ProposalHandler) GetProposalList(c *constant.Context) {
	var req api.ProposalListReq

	if err := c.BindAndValidate(&req); err != nil {
		ParamError(c.Context)
		return
	}

	res, err := p.proposqlService.ProposalList(c.Request.Context(), req)
	if err != nil {
		Error(c.Context, err)
		return
	}

	SuccessWithData(c.Context, res)
}
