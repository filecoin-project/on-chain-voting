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
	"errors"

	"go.uber.org/zap"

	"powervoting-server/constant"
	"powervoting-server/model"
	"powervoting-server/utils"
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

type LotusRepo interface {
	GetActorIdByAddress(ctx context.Context, addr string) (string, error)
	GetValidMinerIds(ctx context.Context, minerId string, minerIds []uint64) (model.StringSlice, error)
	EthAddrToFilcoinAddr(ctx context.Context, addr string) (string, error)
	FilecoinAddressToID(ctx context.Context, addr string) (string, error)
	FilecoinAddrToEthAddr(ctx context.Context, addr string) (string, error) 
}
type ISyncService interface {
	UpdateSyncEventInfo(ctx context.Context, addr string, height int64) error
	GetSyncEventInfo(ctx context.Context, addr string) (*model.SyncEventTbl, error)
	CreateSyncEventInfo(ctx context.Context, in *model.SyncEventTbl) error
	AddProposal(ctx context.Context, in *model.ProposalTbl) error
	UpdateProposal(ctx context.Context, in *model.ProposalTbl) error
	UncountedProposalList(ctx context.Context, chainId, endTime int64) ([]model.ProposalTbl, error)
	BatchUpdateVotes(ctx context.Context, votes []model.VoteTbl) error
	AddVote(ctx context.Context, in *model.VoteTbl) error
	GetUncountedVotedList(ctx context.Context, chainId, proposalId int64) ([]model.VoteTbl, error)
	AddVoterAddress(ctx context.Context, in *model.VoterInfoTbl) error
	CreateFipProposal(ctx context.Context, in *model.FipProposalTbl) error

	UpdateStatusAndGetFipProposal(ctx context.Context, proposalId, chainId int64) (*model.FipProposalTbl, error)
	CreateFipProposalVotedInfo(ctx context.Context, in *model.FipProposalVoteTbl) error
	CreateFipEditor(ctx context.Context, in *model.FipEditorTbl) error
	RevokeFipProposal(ctx context.Context, chainId int64, cadidateAddresss string) error

	UpdateVoterAndProposalGithubNameByGistInfo(ctx context.Context, voterInfo *model.VoterInfoTbl) error
	UpdateVoterByMinerIds(ctx context.Context, voterAddress string, minerIds []uint64) error
}

// SyncService provides functionality for synchronizing data across repositories.
// It encapsulates dependencies required for synchronization operations,
// including access to sync, vote, and proposal repositories.
type SyncService struct {
	repo         SyncRepo     // repo handles synchronization-related data operations
	voteRepo     VoteRepo     // voteRepo manages voting-related data
	proposalRepo ProposalRepo // proposalRepo handles proposal-related data
	fipRepo      FipRepo      // fipRepo handles fip-related data
	lotusRepo    LotusRepo    // lotusRepo handles lotus-related data
}

func NewSyncService(repo SyncRepo, voteRepo VoteRepo, proposalRepo ProposalRepo, fipRepo FipRepo, lotusrepo LotusRepo) *SyncService {
	return &SyncService{
		repo:         repo,
		voteRepo:     voteRepo,
		proposalRepo: proposalRepo,
		fipRepo:      fipRepo,
		lotusRepo:    lotusrepo,
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
	if in == nil {
		return errors.New("sync event is nil")
	}

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
	if in == nil {
		return errors.New("proposal is nil")
	}

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
	if in == nil {
		return errors.New("proposal is nil")
	}

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
	if in == nil {
		return errors.New("vote data is nil")
	}

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

// AddVoterAddress adds a new voter address record to the database or updates an existing one.
// It delegates the operation to the `CreateVoterAddress` method of the `voteRepo` and handles any errors that may occur.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - in: A pointer to the `VoterAddressTbl` struct containing the voter address data to be added or updated.
//
// Returns:
//   - error: An error object if the operation fails. Returns nil on success.
func (s *SyncService) AddVoterAddress(ctx context.Context, in *model.VoterInfoTbl) error {
	if in == nil {
		return errors.New("voter address is nil")
	}

	actorId, err := s.lotusRepo.GetActorIdByAddress(ctx, in.Address)
	if err != nil {
		zap.L().Error("GetActorIdByAddress failed", zap.String("address", in.Address), zap.Error(err))
		return err
	}

	in.OwnerId = actorId
	_, err = s.voteRepo.CreateVoterAddress(ctx, in)
	if err != nil {
		zap.L().Error("CreateVoterAddress failed", zap.Error(err))
		return err
	}

	return nil
}

// CreateFipProposal creates a new FipProposal record in the database by delegating the operation
// to the underlying repository (`fipRepo`). It performs basic validation to ensure the input is not nil.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - in: A pointer to the `FipProposalTbl` struct containing the FipProposal data to be created.
//
// Returns:
//   - error: An error object if the operation fails. Returns `nil` if the operation is successful.
func (s *SyncService) CreateFipProposal(ctx context.Context, in *model.FipProposalTbl) error {
	// Validate the input to ensure it is not nil.
	if in == nil {
		return errors.New("fip proposal is nil") // Return an error if the input is nil
	}

	// Delegate the creation of the FipProposal to the underlying repository (`fipRepo`).
	_, err := s.fipRepo.CreateFipProposal(ctx, in)
	if err != nil {
		// Log the error using the Zap logger for debugging and monitoring purposes.
		zap.L().Error("CreateFipProposal failed", zap.Error(err))
		return err // Return the error if the creation fails
	}

	// Return nil to indicate the operation was successful.
	return nil
}

// CreateFipProposalVotedInfo creates a new FipProposalVoter record in the database by delegating the operation
// to the underlying repository (`fipRepo`). It performs basic validation to ensure the input is not nil.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - in: A pointer to the `FipProposalVoteTbl` struct containing the FipProposalVoter data to be created.
//
// Returns:
//   - error: An error object if the operation fails. Returns `nil` if the operation is successful.
func (s *SyncService) CreateFipProposalVotedInfo(ctx context.Context, in *model.FipProposalVoteTbl) error {
	// Validate the input to ensure it is not nil.
	if in == nil {
		return errors.New("fip proposal voted info is nil") // Return an error if the input is nil
	}

	// Delegate the creation of the FipProposalVoter to the underlying repository (`fipRepo`).
	_, err := s.fipRepo.CreateFipProposalVote(ctx, in)
	if err != nil {
		// Log the error using the Zap logger for debugging and monitoring purposes.
		zap.L().Error("CreateFipProposalVotedInfo failed", zap.Error(err))
		return err // Return the error if the creation fails
	}

	// Return nil to indicate the operation was successful.
	return nil
}

// CreateFipEditor creates a new FipVoter record in the database by delegating the operation
// to the underlying repository (`fipRepo`). It performs basic validation to ensure the input is not nil.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - in: A pointer to the `FipEditorTbl` struct containing the FipVoter data to be created.
//
// Returns:
//   - error: An error object if the operation fails. Returns `nil` if the operation is successful.
func (s *SyncService) CreateFipEditor(ctx context.Context, in *model.FipEditorTbl) error {
	// Validate the input to ensure it is not nil.
	if in == nil {
		return errors.New("fip voter is nil") // Return an error if the input is nil
	}

	// Delegate the creation of the FipVoter to the underlying repository (`fipRepo`).
	_, err := s.fipRepo.CreateFipEditor(ctx, in)
	if err != nil {
		// Log the error using the Zap logger for debugging and monitoring purposes.
		zap.L().Error("CreateFipEditor failed", zap.Error(err))
		return err // Return the error if the creation fails
	}

	// Return nil to indicate the operation was successful.
	return nil
}

// UpdateStatusAndGetFipProposal updates the status of a FipProposal by delegating the operation
// to the underlying repository (`fipRepo`). It handles errors by logging them and returning them to the caller.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - proposalId: The unique identifier of the proposal for which the status is to be updated.
//   - chainId: The unique identifier of the blockchain network associated with the proposal.
//
// Returns:
//   - error: An error object if the operation fails. Returns `nil` if the operation is successful.
func (s *SyncService) UpdateStatusAndGetFipProposal(ctx context.Context, proposalId, chainId int64) (*model.FipProposalTbl, error) {
	// Delegate the update operation to the underlying repository (`fipRepo`).
	res, err := s.fipRepo.UpdateStatusAndGetFipProposal(ctx, proposalId, chainId)
	if err != nil {
		// Log the error using the Zap logger for debugging and monitoring purposes.
		zap.L().Error("UpdateStatusAndGetFipProposal failed", zap.Error(err))
		return nil, err // Return the error if the update fails
	}

	// Return nil to indicate the operation was successful.
	return res, nil
}

func (s *SyncService) RevokeFipProposal(ctx context.Context, chainId int64, cadidateAddresss string) error {
	if err := s.fipRepo.UpdateFipEditorByAddress(ctx, cadidateAddresss); err != nil {
		zap.L().Error(
			"UpdateFipEditorByAddress failed",
			zap.String("cadidateAddresss", cadidateAddresss),
			zap.Error(err),
		)

		return err
	}

	unpassFipProposals, err := s.fipRepo.GetUnpassFipProposalList(ctx, chainId)
	if err != nil {
		zap.L().Error("GetUnpassFipList failed", zap.Error(err))
		return err
	}

	zap.L().Info("RevokeFipProposal unpassFipProposals", zap.Int("length", len(unpassFipProposals)))

	for _, proposal := range unpassFipProposals {
		zap.L().Info("unpassFipProposals", zap.Any("fip proposal", proposal))
		if err := s.fipRepo.UpdateFipProposalVoteByAddress(ctx, proposal.ProposalId, cadidateAddresss); err != nil {
			zap.L().Error(
				"UpdateFipProposalVoteByAddress failed",
				zap.String("cadidateAddresss", cadidateAddresss),
				zap.Int64("proposalId", proposal.ProposalId),
				zap.Error(err),
			)
			continue
		}
	}

	return nil
}

func (s *SyncService) UpdateVoterAndProposalGithubNameByGistInfo(ctx context.Context, voterInfo *model.VoterInfoTbl) error {
	actorId, err := s.lotusRepo.GetActorIdByAddress(ctx, voterInfo.Address)
	if err != nil {
		zap.L().Error("FilecoinAddressToID failed", zap.Error(err))
		return err
	}

	gist, err := utils.FetchGistInfoByGistId(voterInfo.GistId)
	if err != nil {
		zap.L().Error("GetGistInfoByGistId failed", zap.Error(err))
		return err
	}

	gistFile := model.GistFiles{}
	for _, value := range gist.Files {
		gistFile = value
	}

	gistInfo := ""
	if gistFile.Size <= constant.MaxFileSize {
		gistInfo = gistFile.Content
	}

	zap.L().Info(
		"UpdateVoterByGistInfo",
		zap.String("gistId", voterInfo.GistId),
		zap.String("github id", gist.Owner.Login),
	)

	voterInfo.GithubName = gist.Owner.Login
	voterInfo.GithubAvatar = gist.Owner.AvatarUrl
	voterInfo.GistInfo = gistInfo
	voterInfo.OwnerId = actorId

	isValid := utils.VerifyAuthorizeAllow(gist.Owner.Login, gist, func(gistAddr string) bool {
		if gistAddr == voterInfo.Address {
			return true
		}

		reqAddrActorId, err := s.lotusRepo.GetActorIdByAddress(ctx, voterInfo.Address)
		if err != nil {
			zap.L().Error("GetActorIdByAddress failed by VerifyGist", zap.String("address", voterInfo.Address), zap.Error(err))
			return false
		}

		gistAddrActorId, err := s.lotusRepo.GetActorIdByAddress(ctx, gistAddr)
		if err != nil {
			zap.L().Error("GetActorIdByAddress failed by VerifyGist", zap.String("address", gistAddr), zap.Error(err))
			return false
		}

		return reqAddrActorId == gistAddrActorId
	})
	if !isValid {
		voterInfo.GithubName = ""
		voterInfo.GistInfo = ""
		voterInfo.GistId = ""
		voterInfo.GithubAvatar = ""
	}

	if err := s.voteRepo.UpdateVoterByGistInfo(ctx, voterInfo); err != nil {
		zap.L().Error("UpdateVoterByGistInfo failed", zap.Error(err))
		return err
	}

	return nil
}

func (s *SyncService) UpdateVoterByMinerIds(ctx context.Context, voterAddress string, minerIds []uint64) error {
	actorId, err := s.lotusRepo.GetActorIdByAddress(ctx, voterAddress)
	if err != nil {
		zap.L().Error("FilecoinAddressToID failed", zap.String("Func", "UpdateVoterByMinerIds"), zap.String("voterAddress", voterAddress), zap.Error(err))
		return err
	}

	validMinerIds, err := s.lotusRepo.GetValidMinerIds(ctx, actorId, minerIds)

	if err := s.voteRepo.UpdateVoterByMinerInfo(ctx, &model.VoterInfoTbl{
		Address:  voterAddress,
		MinerIds: validMinerIds,
		OwnerId:  actorId,
	}); err != nil {

		zap.L().Error("UpdateVoterByMinerIds failed", zap.Error(err))

		return err
	}

	return nil
}
