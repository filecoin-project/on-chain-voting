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

type VoteHandler struct {
	voteServer service.IVoteService
}

func NewVoteHandler(ps service.IVoteService) *VoteHandler {
	return &VoteHandler{
		voteServer: ps,
	}
}

// GetCountedVotesInfo returns the counted votes info
func (h *VoteHandler) GetCountedVotesInfo(c *constant.Context) {
	var req api.ProposalReq
	if err := c.BindAndValidate(&req); err != nil {
		ParamError(c.Context)
		return
	}

	res, err := h.voteServer.GetCountedVotedList(c.Request.Context(), req.ChainId, req.ProposalId)
	if err != nil {
		SystemError(c.Context)
		return
	}

	SuccessWithData(c.Context, res)
}

func (f *VoteHandler) GetFipEditorGistInfo(c *constant.Context) {
	var req api.AddressReq
	if err := c.BindAndValidate(&req); err != nil {
		ParamError(c.Context)
		return
	}

	gistInfo, err := f.voteServer.GetFipEditorGistInfo(c.Context, req)
	if err != nil {
		Error(c.Context, err)
		return
	}

	SuccessWithData(c.Context, gistInfo)
}

func (f *VoteHandler) VerifyGistValid(c *constant.Context) {
	var req api.VerifyGistReq
	if err := c.BindAndValidate(&req); err != nil {
		ParamError(c.Context)
		return
	}

	obj, err := f.voteServer.VerifyGist(c.Context, req)
	if err != nil {
		Error(c.Context, err)
		return
	}

	SuccessWithData(c.Context, obj)
}
