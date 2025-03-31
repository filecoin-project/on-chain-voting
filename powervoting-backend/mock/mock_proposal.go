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

package mock

import (
	"context"
	"errors"
	"fmt"

	"powervoting-server/model"
	"powervoting-server/model/api"
	"powervoting-server/service"
)

type MockProposalService struct {
}

// UpdateProposalGitHubName implements service.ProposalRepo.
func (m *MockProposalService) UpdateProposalGitHubName(ctx context.Context, createrAddress string, githubName string) error {
	panic("unimplemented")
}

var _ service.ProposalRepo = (*MockProposalService)(nil)

// CreateProposal implements service.ProposalRepo.
func (m *MockProposalService) CreateProposal(ctx context.Context, in *model.ProposalTbl) (int64, error) {
	panic("unimplemented")
}

// CreateProposalDraft implements service.ProposalRepo.
func (m *MockProposalService) CreateProposalDraft(ctx context.Context, in *model.ProposalDraftTbl) (int64, error) {
	fmt.Printf("mock create proposal draft success.  proposal: %v", in)
	return 1, nil
}

// GetProposalById implements service.ProposalRepo.
func (m *MockProposalService) GetProposalById(ctx context.Context, req api.ProposalReq) (*model.ProposalTbl, error) {
	return &model.ProposalTbl{
		ProposalId: 1,
	}, nil
}

// GetProposalDraftByAddress implements service.ProposalRepo.
func (m *MockProposalService) GetProposalDraftByAddress(ctx context.Context, req api.GetDraftReq) (*model.ProposalDraftTbl, error) {
	if req.Address == "test" {
		return &model.ProposalDraftTbl{
			Creator: "test",
		}, nil
	}

	return nil, errors.New("not found")
}

// GetProposalListWithPagination implements service.ProposalRepo.
func (m *MockProposalService) GetProposalListWithPagination(ctx context.Context, req api.ProposalListReq) ([]model.ProposalWithVoted, int64, error) {
	return []model.ProposalWithVoted{
		{
			ProposalTbl: model.ProposalTbl{
				ProposalId: 1,
			},
		},
	}, 1, nil
}

// GetUncountedProposalList implements service.ProposalRepo.
func (m *MockProposalService) GetUncountedProposalList(ctx context.Context, chainId int64, timestamp int64) ([]model.ProposalTbl, error) {
	panic("unimplemented")
}

// UpdateProposal implements service.ProposalRepo.
func (m *MockProposalService) UpdateProposal(ctx context.Context, in *model.ProposalTbl) error {
	panic("unimplemented")
}
