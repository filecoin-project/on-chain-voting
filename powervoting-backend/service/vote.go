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

package service

import (
	"context"

	"go.uber.org/zap"

	"powervoting-server/model"
	"powervoting-server/model/api"
)

// VoteRepo defines the interface for managing vote-related data.
// It provides methods for creating, retrieving, and updating vote records.
type VoteRepo interface {
	// CreateVote creates a new vote record in the repository.
	// This is typically called when a Vote event is parsed during the execution of a synchronous contract event.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout.
	//   - in: The vote data to be created.
	//
	// Returns:
	//   - int64: The ID of the created record if successful.
	//   - error: An error if the creation operation fails; otherwise, nil.
	CreateVote(ctx context.Context, in *model.VoteTbl) (int64, error)

	// GetVoteList retrieves a list of votes for a specific proposal, optionally filtered by counted status.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout.
	//   - chainId: The chain ID associated with the proposal.
	//   - proposalId: The proposal ID to filter votes.
	//   - counted: A flag to filter votes based on whether they have been counted.
	//
	// Returns:
	//   - []model.VoteTbl: A list of votes matching the criteria.
	//   - error: An error if the query operation fails; otherwise, nil.
	GetVoteList(ctx context.Context, chainId, proposalId int64, counted bool) ([]model.VoteTbl, error)

	// BatchUpdateVotes updates multiple vote records in the repository in a single operation.
	// This is typically used during the vote counting process, where the weight of each vote
	// is calculated and the results are updated in the database.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout.
	//   - votes: A slice of vote records to be updated.
	//
	// Returns:
	//   - error: An error if the batch update operation fails; otherwise, nil.
	BatchUpdateVotes(ctx context.Context, votes []model.VoteTbl) error

	// CreateVoterAddress creates or updates a voter address record in the database.
	// If a record with the same address already exists, it updates the `update_height` field.
	// Otherwise, it inserts a new record.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
	//   - in: A pointer to the `VoterAddressTbl` struct containing the voter address data to be created or updated.
	//
	// Returns:
	//   - int64: The ID of the created or updated voter address record.
	//   - error: An error object if the operation fails.
	CreateVoterAddress(ctx context.Context, in *model.VoterAddressTbl) (int64, error)

	// GetNewVoterAddresss retrieves a list of voter addresses that were created after a specified block height.
	// It queries the database for voter addresses with an `init_created_height` greater than the provided height,
	// orders the results in descending order by `init_created_height`, and returns the list along with the highest height found.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
	//   - height: The block height threshold. Only voter addresses created after this height are returned.
	//
	// Returns:
	//   - []model.VoterAddressTbl: A list of voter addresses that meet the criteria.
	//   - int64: The highest `init_created_height` from the retrieved voter addresses. Returns 0 if no addresses are found.
	//   - error: An error object if the query fails.
	GetNewVoterAddresss(ctx context.Context, height int64) ([]model.VoterAddressTbl, error)
}

type IVoteService interface {
	GetCountedVotedList(ctx context.Context, chainId, proposalId int64) ([]api.Voted, error)
}

var _ IVoteService = (*VoteService)(nil)

// VoteService provides functionality for managing vote-related operations.
// It encapsulates the dependencies required for interacting with vote data,
// such as creating, retrieving, and updating vote records.
type VoteService struct {
	repo VoteRepo // repo provides access to the underlying vote repository
}

func NewVoteService(repo VoteRepo) *VoteService {
	return &VoteService{
		repo: repo,
	}
}

// GetCountedVotedList retrieves a list of counted votes for a specific chain and proposal.
// It queries the underlying vote repository for the data and logs any errors encountered.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - chainId: The chain ID to filter votes.
//   - proposalId: The proposal ID to filter votes.
//
// Returns:
//   - []model.VoteTbl: A list of counted votes if the query is successful.
//   - error: An error if the query operation fails; otherwise, nil.
func (v *VoteService) GetCountedVotedList(ctx context.Context, chainId, proposalId int64) ([]api.Voted, error) {
	res, err := v.repo.GetVoteList(ctx, chainId, proposalId, true)
	if err != nil {
		zap.L().Error("GetVoteList failed", zap.Error(err))
		return nil, err
	}

	var votes []api.Voted
	for _, v := range res {
		votes = append(votes, api.Voted{
			ChainId:      v.ChainId,
			ProposalId:   v.ProposalId,
			VoterAddress: v.Address,
			VotedResult:  v.VoteResult,
			VotedTime:    v.CreatedAt.Unix(),
			PowerRep: api.PowerRep{
				SpPower:          v.SpPower,
				TokenHolderPower: v.TokenHolderPower,
				ClientPower:      v.ClientPower,
				DeveloperPower:   v.DeveloperPower,
			},
		})
	}
	return votes, nil
}
