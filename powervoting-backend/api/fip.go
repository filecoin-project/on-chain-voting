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

type FipHandle struct {
	fipService service.IFipService
}

func NewFipHandle(fipService service.IFipService) *FipHandle {
	return &FipHandle{
		fipService: fipService,
	}
}

// GetFipProposalList is a method of the FipHandle struct that retrieves a list of FIPs (Feature Improvement Proposals).
func (f *FipHandle) GetFipProposalList(c *constant.Context) {
	// Declare a variable of type api.FipEditorListReq to hold the request data.
	var req api.FipProposalListReq
	// Attempt to bind and validate the incoming request data to the req variable.
	// If there is an error during binding or validation, call the ParamError function to handle it.
	if err := c.BindAndValidate(&req); err != nil {
		ParamError(c.Context)
		return
	}

	// Call the GetFipProposalList method of the fipService field to retrieve the list of FIPs based on the request data.
	// Pass the context from the incoming request and the validated request data.
	fipProposalList, err := f.fipService.GetFipProposalList(c.Context, req)
	// If there is an error during the retrieval, call the SystemError function to handle it.
	if err != nil {
		SystemError(c.Context)
		return
	}

	// If the retrieval is successful, call the SuccessWithData function to send the FIP list back in the response.
	SuccessWithData(c.Context, fipProposalList)
}

func (f *FipHandle) GetFipEditorList(c *constant.Context) {
	var req api.FipEditorListReq
	if err := c.BindAndValidate(&req); err != nil {
		ParamError(c.Context)
		return
	}

	fipEditorList, err := f.fipService.GetFipEditorList(c.Context, req)
	if err != nil {
		SystemError(c.Context)
		return
	}

	SuccessWithData(c.Context, fipEditorList)
}
