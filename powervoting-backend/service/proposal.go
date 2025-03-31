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
	"fmt"
	"time"

	"go.uber.org/zap"

	"powervoting-server/constant"
	"powervoting-server/model"
	"powervoting-server/model/api"
)

// ProposalRepo defines the interface for managing proposal-related data.
// It provides methods for creating, retrieving, and updating proposals and proposal drafts.
type ProposalRepo interface {
	// GetProposalListWithPagination retrieves a paginated list of proposals.
	// The results can be filtered by status and searched by keywords in the title.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout.
	//   - req: Contains pagination details, status filter, and search keywords.
	//
	// Returns:
	//   - []model.ProposalTbl: A list of proposals matching the criteria.
	//   - int64: The total number of proposals matching the criteria (for pagination).
	//   - error: An error if the query operation fails; otherwise, nil.
	GetProposalListWithPagination(ctx context.Context, req api.ProposalListReq) ([]model.ProposalWithVoted, int64, error)

	// GetProposalById retrieves the details of a single proposal by its ID.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout.
	//   - req: Contains the proposal ID to query.
	//
	// Returns:
	//   - *model.ProposalTbl: The proposal details if found.
	//   - error: An error if the query operation fails; otherwise, nil.
	GetProposalById(ctx context.Context, req api.ProposalReq) (*model.ProposalTbl, error)

	// CreateProposalDraft creates a new proposal draft.
	// A creator can only have one draft proposal at a time; creating a new draft overwrites the existing one.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout.
	//   - in: The proposal draft data to be created.
	//
	// Returns:
	//   - int64: The ID of the created draft if successful.
	//   - error: An error if the creation operation fails; otherwise, nil.
	CreateProposalDraft(ctx context.Context, in *model.ProposalDraftTbl) (int64, error)

	// GetProposalDraftByAddress retrieves a proposal draft by the creator's address.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout.
	//   - req: Contains the creator's address to query.
	//
	// Returns:
	//   - *model.ProposalDraftTbl: The proposal draft if found.
	//   - error: An error if the query operation fails; otherwise, nil.
	GetProposalDraftByAddress(ctx context.Context, req api.GetDraftReq) (*model.ProposalDraftTbl, error)

	// CreateProposal creates a new proposal record in the repository.
	// This is typically called when a ProposalCreate event is parsed during the execution of a synchronous contract event.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout.
	//   - in: The proposal data to be created.
	//
	// Returns:
	//   - int64: The ID of the created proposal if successful.
	//   - error: An error if the creation operation fails; otherwise, nil.
	CreateProposal(ctx context.Context, in *model.ProposalTbl) (int64, error)

	// UpdateProposal updates the counting information of a proposal in the repository.
	// This is typically called during the vote counting process when the proposal reaches the time to count votes.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout.
	//   - in: The proposal data to be updated.
	//
	// Returns:
	//   - error: An error if the update operation fails; otherwise, nil.
	UpdateProposal(ctx context.Context, in *model.ProposalTbl) error

	// GetUncountedProposalList retrieves a list of proposals that have not been counted yet, filtered by chain ID and timestamp.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout.
	//   - chainId: The chain ID to filter proposals.
	//   - timestamp: The timestamp to filter proposals (e.g., proposals before this timestamp).
	//
	// Returns:
	//   - []model.ProposalTbl: A list of uncounted proposals if the query is successful.
	//   - error: An error if the query operation fails; otherwise, nil.
	GetUncountedProposalList(ctx context.Context, chainId int64, timestamp int64) ([]model.ProposalTbl, error)

	UpdateProposalGitHubName(ctx context.Context, createrAddress, githubName string) error
}

// IProposalService defines the interface for managing proposal-related operations.
// It provides methods for creating, retrieving, and listing proposals and proposal drafts.
type IProposalService interface {
	AddDraft(ctx context.Context, req *api.AddProposalDraftReq) error

	GetDraft(ctx context.Context, req api.GetDraftReq) (*api.ProposalDraftRep, error)

	ProposalDetail(ctx context.Context, req api.ProposalReq) (*api.ProposalRep, error)

	ProposalList(ctx context.Context, req api.ProposalListReq) (*api.CountListRep, error)
}

var _ IProposalService = (*ProposalService)(nil)

// ProposalService provides functionality for managing proposal-related operations.
// It encapsulates the dependencies required for interacting with proposal data,
// such as creating, retrieving, and updating proposals and proposal drafts.
type ProposalService struct {
	repo ProposalRepo // repo provides access to the underlying proposal repository
}

func NewProposalService(repo ProposalRepo) *ProposalService {
	return &ProposalService{
		repo: repo,
	}
}

// ProposalList retrieves a paginated list of proposals and transforms the data into a response format.
// It queries the underlying proposal repository for the data, handles status calculations, and logs any errors encountered.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - req: Contains pagination details, status filter, and search keywords.
//
// Returns:
//   - *api.CountListRep: A response containing the list of proposals and the total count.
//   - error: An error if the query operation fails; otherwise, nil.
func (p *ProposalService) ProposalList(ctx context.Context, req api.ProposalListReq) (*api.CountListRep, error) {
	// Fetch proposals from the repository with pagination
	proposals, count, err := p.repo.GetProposalListWithPagination(ctx, req)
	if err != nil {
		zap.L().Error("GetProposalListWithPagination error", zap.Error(err))
		return nil, errors.New("fail to get proposal list")
	}

	// Transform proposals into the response format
	var list []api.ProposalRep
	for _, proposal := range proposals {
		temp := api.ProposalRep{
			ProposalId: proposal.ProposalId,
			ChainId:    proposal.ChainId,
			GithubName: proposal.GithubName,
			Content:    proposal.Content,
			StartTime:  proposal.StartTime,
			EndTime:    proposal.EndTime,
			Voted:      proposal.Voted != nil,
			Creator:    proposal.Creator,
			Title:      proposal.Title,
			Status:     constant.ProposalStatusPending,
			CreatedAt:  proposal.CreatedAt.Unix(),
			UpdatedAt:  proposal.UpdatedAt.Unix(),
			VotePercentage: api.ProposalVotePercentage{
				Reject:  proposal.RejectPercentage,
				Approve: proposal.ApprovePercentage,
			},
		}

		// Calculate the proposal status based on the current time and counted flag
		if proposal.Counted == constant.ProposalCreate {
			if temp.StartTime < time.Now().Unix() && temp.EndTime > time.Now().Unix() {
				temp.Status = constant.ProposalStatusInProgress
			} else if temp.EndTime < time.Now().Unix() {
				temp.Status = constant.ProposalStatusCounting
			}
		} else {
			temp.Status = constant.ProposalStatusCompleted
		}

		list = append(list, temp)
	}

	// Return the response with the list of proposals and total count
	return &api.CountListRep{
		Total: count,
		List:  list,
	}, nil
}

// ProposalDetail retrieves the details of a specific proposal by its ID.
// It queries the underlying proposal repository for the data and logs any errors encountered.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - req: Contains the proposal ID to query.
//
// Returns:
//   - *api.ProposalRep: The proposal details if found.
//   - error: An error if the query operation fails or if no proposal is found; otherwise, nil.
func (p *ProposalService) ProposalDetail(ctx context.Context, req api.ProposalReq) (*api.ProposalRep, error) {
	// Fetch the proposal details from the repository
	proposal, err := p.repo.GetProposalById(ctx, req)
	if err != nil {
		zap.L().Error("GetProposalById error", zap.Error(err))
		return nil, errors.New("fail to get proposal detail")
	}

	// Check if the proposal exists
	if proposal == nil {
		return nil, fmt.Errorf("no proposal found by id: %d", req.ProposalId)
	}

	status := constant.ProposalStatusPending
	if proposal.Counted == constant.ProposalCreate {
		if proposal.StartTime < time.Now().Unix() && proposal.EndTime > time.Now().Unix() {
			status = constant.ProposalStatusInProgress
		} else if proposal.EndTime < time.Now().Unix() {
			status = constant.ProposalStatusCounting
		}
	} else {
		status = constant.ProposalStatusCompleted
	}

	// Transform the proposal data into the response format
	return &api.ProposalRep{
		ProposalId: proposal.ProposalId,
		ChainId:    proposal.ChainId,
		Content:    proposal.Content,
		GithubName: proposal.GithubName,
		StartTime:  proposal.StartTime,
		EndTime:    proposal.EndTime,
		Creator:    proposal.Creator,
		Title:      proposal.Title,
		Status:     status,
		CreatedAt:  proposal.CreatedAt.Unix(),
		UpdatedAt:  proposal.UpdatedAt.Unix(),
		VotePercentage: api.ProposalVotePercentage{
			Approve: proposal.ApprovePercentage,
			Reject:  proposal.RejectPercentage,
		},
		Percentage: api.ProposalPercentage{
			SpPercentage:          proposal.SpPercentage,
			TokenHolderPercentage: proposal.TokenHolderPercentage,
			ClientPercentage:      proposal.ClientPercentage,
			DeveloperPercentage:   proposal.DeveloperPercentage,
		},
		TotalPower: api.TotalPower{
			SpPower:          proposal.TotalSpPower,
			TokenHolderPower: proposal.TotalTokenHolderPower,
			ClientPower:      proposal.TotalClientPower,
			DeveloperPower:   proposal.TotalDeveloperPower,
		},
		SnapshotInfo: api.SnapshotInfo{
			SnapshotHeight: proposal.SnapshotBlockHeight,
			SnapshotDay:    proposal.SnapshotDay,
		},
	}, nil
}

// AddDraft creates a new proposal draft and persists it to the repository.
// It transforms the request data into the draft model and logs any errors encountered.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - req: Contains the draft proposal data to be created.
//
// Returns:
//   - error: An error if the creation operation fails; otherwise, nil.
func (p *ProposalService) AddDraft(ctx context.Context, req *api.AddProposalDraftReq) error {
	// Create a new proposal draft in the repository
	_, err := p.repo.CreateProposalDraft(ctx, &model.ProposalDraftTbl{
		Creator:   req.Creator,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Timezone:  req.Timezone,
		ChainId:   req.ChainId,
		Title:     req.Title,
		Content:   req.Content,
		Percentage: model.Percentage{
			SpPercentage:          req.SpPercentage,
			TokenHolderPercentage: req.TokenHolderPercentage,
			ClientPercentage:      req.ClientPercentage,
			DeveloperPercentage:   req.DeveloperPercentage,
		},
	})

	// Handle errors during draft creation
	if err != nil {
		zap.L().Error("CreateProposalDraft error", zap.Error(err))
		return errors.New("fail to create proposal draft")
	}

	return nil
}

// GetDraft retrieves a proposal draft by the creator's address.
// It queries the underlying repository for the draft data and logs any errors encountered.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - req: Contains the creator's address to query.
//
// Returns:
//   - *api.ProposalDraftRep: The proposal draft details if found.
//   - error: An error if the query operation fails or if no draft is found; otherwise, nil.
func (p *ProposalService) GetDraft(ctx context.Context, req api.GetDraftReq) (*api.ProposalDraftRep, error) {
	// Fetch the proposal draft from the repository
	res, err := p.repo.GetProposalDraftByAddress(ctx, req)
	if err != nil {
		zap.L().Error("GetProposalByCreatorAddress error", zap.Error(err))
		return nil, errors.New("fail to get proposal draft")
	}

	// Check if the draft exists
	if res == nil {
		return nil, errors.New("no draft found")
	}

	// Transform the draft data into the response format
	return &api.ProposalDraftRep{
		Content:               res.Content,
		StartTime:             res.StartTime,
		EndTime:               res.EndTime,
		Timezone:              res.Timezone,
		Title:                 res.Title,
		SpPercentage:          res.SpPercentage,
		TokenHolderPercentage: res.TokenHolderPercentage,
		ClientPercentage:      res.ClientPercentage,
		DeveloperPercentage:   res.DeveloperPercentage,
	}, nil
}
