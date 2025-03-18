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

	"powervoting-server/model"
	"powervoting-server/service"
)

type MockVoteService struct {
}



var _ service.VoteRepo = (*MockVoteService)(nil)

// BatchUpdateVotes implements service.VoteRepo.
func (m *MockVoteService) BatchUpdateVotes(ctx context.Context, votes []model.VoteTbl) error {
	panic("unimplemented")
}

// CreateVote implements service.VoteRepo.
func (m *MockVoteService) CreateVote(ctx context.Context, in *model.VoteTbl) (int64, error) {
	panic("unimplemented")
}

// GetVoteList implements service.VoteRepo.
func (m *MockVoteService) GetVoteList(ctx context.Context, chainId int64, proposalId int64, counted bool) ([]model.VoteTbl, error) {
	if chainId != 1 || proposalId != 1 {
		return nil, errors.New("invalid chainId or proposalId")
	}
	return []model.VoteTbl{
		{
			ChainId:    chainId,
			ProposalId: proposalId,
			Address:    "voter",
		},
	}, nil
}

// CreateVoterAddress implements service.VoteRepo.
func (m *MockVoteService) CreateVoterAddress(ctx context.Context, in *model.VoterAddressTbl) (int64, error) {
	panic("unimplemented")
}

// GetNewVoterAddresss implements service.VoteRepo.
func (m *MockVoteService) GetNewVoterAddresss(ctx context.Context, height int64) ([]model.VoterAddressTbl, error) {
	panic("unimplemented")
}