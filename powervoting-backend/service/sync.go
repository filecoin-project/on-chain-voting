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
)

// SyncRepo defines the interface for managing synchronization event data.
// It provides methods for creating, updating, and retrieving synchronization event information.
type SyncRepo interface {
	// CreateSyncEventInfo creates a new synchronization event record in the repository.
	// If the synchronization information already exists in the database, it will be ignored.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout.
	//   - in: The synchronization event data to be created.
	//
	// Returns:
	//   - int64: The ID of the created record if successful.
	//   - error: An error if the creation operation fails; otherwise, nil.
	CreateSyncEventInfo(ctx context.Context, in *model.SyncEventTbl) (int64, error)

	// UpdateSyncEventInfo updates the synchronization event information for a given address and block height.
	// This is typically called after a timed task completes the synchronization of contract events.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout.
	//   - addr: The address associated with the sync event.
	//   - height: The block height to update.
	//
	// Returns:
	//   - error: An error if the update operation fails; otherwise, nil.
	UpdateSyncEventInfo(ctx context.Context, addr string, height int64) error

	// GetSyncEventInfo retrieves synchronization event information for a given address.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout.
	//   - addr: The address associated with the sync event.
	//
	// Returns:
	//   - *model.SyncEventTbl: The synchronization event data if found.
	//   - error: An error if the query operation fails; otherwise, nil.
	GetSyncEventInfo(ctx context.Context, addr string) (*model.SyncEventTbl, error)
}

// SyncService provides functionality for synchronizing data across repositories.
// It encapsulates dependencies required for synchronization operations,
// including access to sync, vote, and proposal repositories.
type SyncService struct {
	repo         SyncRepo     // repo handles synchronization-related data operations
	voteRepo     VoteRepo     // voteRepo manages voting-related data
	proposalRepo ProposalRepo // proposalRepo handles proposal-related data
}

func NewSyncService(repo SyncRepo, voteRepo VoteRepo, proposalRepo ProposalRepo) *SyncService {
	return &SyncService{
		repo:         repo,
		voteRepo:     voteRepo,
		proposalRepo: proposalRepo,
	}
}

// UpdateSyncEventInfo updates the synchronization event information for a given address and block height.
// It delegates the update operation to the underlying repository and logs any errors encountered.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - addr: The address associated with the sync event.
//   - height: Block height that has been synchronized.
//
// Returns:
//   - error: An error if the update operation fails; otherwise, nil.
func (s *SyncService) UpdateSyncEventInfo(ctx context.Context, addr string, height int64) error {
	if err := s.repo.UpdateSyncEventInfo(ctx, addr, height); err != nil {
		zap.L().Error("UpdateSyncEventInfo failed", zap.Error(err))
		return err
	}

	return nil
}

// GetSyncEventInfo retrieves synchronization event information for a given address.
// It queries the underlying repository for the data and logs any errors encountered.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - addr: The address is the power voting contract address.
//
// Returns:
//   - *model.SyncEventTbl: The synchronization event data if found.
//   - error: An error if the query operation fails; otherwise, nil.
func (s *SyncService) GetSyncEventInfo(ctx context.Context, addr string) (*model.SyncEventTbl, error) {
	res, err := s.repo.GetSyncEventInfo(ctx, addr)
	if err != nil {
		zap.L().Error("GetSyncEventInfo failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}

// CreateSyncEventInfo creates a new synchronization event record in the repository.
// It delegates the creation operation to the underlying repository and logs any errors encountered.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - in: The synchronization event data to be created.
//
// Returns:
//   - error: An error if the creation operation fails; otherwise, nil.
func (s *SyncService) CreateSyncEventInfo(ctx context.Context, in *model.SyncEventTbl) error {
	_, err := s.repo.CreateSyncEventInfo(ctx, in)
	if err != nil {
		zap.L().Error("CreateSyncEventInfo failed", zap.Error(err))
	}

	return nil
}

// AddProposal adds a new proposal record to the repository.
// It delegates the creation operation to the underlying proposal repository and logs any errors encountered.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - in: The proposal data to be added.
//
// Returns:
//   - error: An error if the creation operation fails; otherwise, nil.
func (s *SyncService) AddProposal(ctx context.Context, in *model.ProposalTbl) error {
	_, err := s.proposalRepo.CreateProposal(ctx, in)
	if err != nil {
		zap.L().Error("AddProposal failed", zap.Error(err))
		return err
	}

	return nil
}

// UpdateProposal updates an existing proposal record in the repository.
// It delegates the update operation to the underlying proposal repository and logs any errors encountered.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - in: The proposal data to be updated.
//
// Returns:
//   - error: An error if the update operation fails; otherwise, nil.
func (s *SyncService) UpdateProposal(ctx context.Context, in *model.ProposalTbl) error {
	if err := s.proposalRepo.UpdateProposal(ctx, in); err != nil {
		zap.L().Error("UpdateProposal failed", zap.Error(err))
		return err
	}

	return nil
}

// UncountedProposalList retrieves a list of proposals that have not been counted, filtered by chain ID and end time.
// It queries the underlying proposal repository for the data and logs any errors encountered.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - chainId: The chain ID to filter proposals.
//   - endTime: The end time to filter proposals (e.g., proposals before this timestamp).
//
// Returns:
//   - []model.ProposalTbl: A list of uncounted proposals if the query is successful.
//   - error: An error if the query operation fails; otherwise, nil.
func (s *SyncService) UncountedProposalList(ctx context.Context, chainId, endTime int64) ([]model.ProposalTbl, error) {
	proposals, err := s.proposalRepo.GetUncountedProposalList(ctx, chainId, endTime)
	if err != nil {
		zap.L().Error("GetUncountedProposalList error", zap.Error(err))
		return nil, err
	}

	return proposals, nil
}

// BatchUpdateVotes updates multiple vote records in the repository in a single operation.
// It delegates the batch update operation to the underlying vote repository and logs any errors encountered.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - votes: A slice of vote records to be updated.
//
// Returns:
//   - error: An error if the batch update operation fails; otherwise, nil.
func (s *SyncService) BatchUpdateVotes(ctx context.Context, votes []model.VoteTbl) error {
	if err := s.voteRepo.BatchUpdateVotes(ctx, votes); err != nil {
		zap.L().Error("BatchUpdateVotes failed", zap.Error(err))
		return err
	}

	return nil
}

// AddVote adds a new vote record to the repository.
// It delegates the creation operation to the underlying vote repository and logs any errors encountered.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - in: The vote data to be added.
//
// Returns:
//   - error: An error if the creation operation fails; otherwise, nil.
func (s *SyncService) AddVote(ctx context.Context, in *model.VoteTbl) error {
	_, err := s.voteRepo.CreateVote(ctx, in)
	if err != nil {
		zap.L().Error("CreateVote failed", zap.Error(err))
		return err
	}

	return nil
}

// GetUncountedVotedList retrieves a list of uncounted votes for a specific chain and proposal.
// It queries the underlying vote repository for the data and logs any errors encountered.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - chainId: The chain ID to filter votes.
//   - proposalId: The proposal ID to filter votes.
//
// Returns:
//   - []model.VoteTbl: A list of uncounted votes if the query is successful.
//   - error: An error if the query operation fails; otherwise, nil.
func (s *SyncService) GetUncountedVotedList(ctx context.Context, chainId, proposalId int64) ([]model.VoteTbl, error) {
	res, err := s.voteRepo.GetVoteList(ctx, chainId, proposalId, false)
	if err != nil {
		zap.L().Error("GetVoteList failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}

