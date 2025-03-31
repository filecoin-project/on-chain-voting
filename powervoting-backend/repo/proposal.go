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

package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"powervoting-server/constant"
	"powervoting-server/data"
	"powervoting-server/model"
	"powervoting-server/model/api"
	"powervoting-server/service"
)

var _ service.ProposalRepo = (*ProposalRepoImpl)(nil)

type ProposalRepoImpl struct {
	mydb *gorm.DB
}

func NewProposalRepo(mydb *gorm.DB) *ProposalRepoImpl {
	return &ProposalRepoImpl{
		mydb: mydb,
	}
}

// GetProposalListWithPagination retrieves a paginated list of proposals based on the given request.
// It returns a slice of ProposalTbl, the total count of proposals, and an error if any occurred.
func (p *ProposalRepoImpl) GetProposalListWithPagination(ctx context.Context, req api.ProposalListReq) ([]model.ProposalWithVoted, int64, error) {
	// Build the base SQL queries for counting and listing proposals based on the request parameters.
	queryCount, queryList := p.buildProposalBaseQuery(req)
	// Apply status filter to both the count and list queries if a status is specified in the request.
	p.applyStatusFilter(queryCount, req.Status)
	p.applyStatusFilter(queryList, req.Status)

	// Execute the count query to get the total number of proposals that match the criteria.
	count, err := data.ExecuteCountQuery(ctx, queryCount)
	if err != nil {
		return nil, 0, fmt.Errorf("get proposal list count error: %w", err)
	}

	subQuery := p.mydb.Model(&model.VoteTbl{}).
		Select("1").
		Where("proposal_id = proposal_tbl.proposal_id").
		Where("chain_id = proposal_tbl.chain_id").
		Where("address = ?", req.Addr)

	queryList.Select("proposal_tbl.*, (?) AS voted", subQuery)

	var proposals []model.ProposalWithVoted
	if err := queryList.WithContext(ctx).Find(&proposals).Error; err != nil {
		return nil, 0, fmt.Errorf("get proposal list error: %w", err)
	}

	return proposals, count, nil
}

// GetProposalById retrieves a proposal from the database based on the provided proposal ID and chain ID.
func (p *ProposalRepoImpl) GetProposalById(ctx context.Context, req api.ProposalReq) (*model.ProposalTbl, error) {
	var proposal model.ProposalTbl
	// Query the database for a proposal with the specified ID and chain ID.
	// The WithContext method ensures that the query is cancellable and can be timed out.
	// The Where method specifies the conditions for the query.
	// The First method retrieves the first record that matches the conditions.
	query := p.mydb.Model(&model.ProposalTbl{}).
		Where("proposal_id = ? AND chain_id = ?", req.ProposalId, req.ChainId)

	if err := query.WithContext(ctx).First(&proposal).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("get proposal by id error: %w", err)
	}

	return &proposal, nil
}

// CreateProposalDraft creates a new proposal draft in the database.
// If a proposal draft with the same creator already exists, it updates the existing draft.
func (p *ProposalRepoImpl) CreateProposalDraft(ctx context.Context, in *model.ProposalDraftTbl) (int64, error) {
	if err := p.mydb.Model(model.ProposalDraftTbl{}).
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "creator"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"start_time",
				"end_time",
				"title",
				"timezone",
				"content",
				"token_holder_percentage",
				"sp_percentage",
				"client_percentage",
				"developer_percentage",
			}),
		}).Create(in).Error; err != nil {
		return 0, fmt.Errorf("create proposal draft error: %w", err)
	}

	return in.ID, nil
}

// GetProposalDraftByAddress retrieves a proposal draft from the database based on the creator's address.
func (p *ProposalRepoImpl) GetProposalDraftByAddress(ctx context.Context, req api.GetDraftReq) (*model.ProposalDraftTbl, error) {
	var proposalDraft model.ProposalDraftTbl
	// Query the database for a proposal draft where the creator matches the provided address.
	// The query uses GORM's Model method to specify the model type and WithContext to pass the context.
	// The Where clause filters records based on the creator field, and First retrieves the first matching record.
	if err := p.mydb.Model(model.ProposalDraftTbl{}).
		WithContext(ctx).
		Where("creator = ?", req.Address).
		First(&proposalDraft).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("get proposal draft error: %w", err)
	}

	return &proposalDraft, nil
}

// CreateProposal creates a new proposal in the database.
// If a proposal with the same proposal_id already exists, it updates the existing proposal.
func (p *ProposalRepoImpl) CreateProposal(ctx context.Context, in *model.ProposalTbl) (int64, error) {
	if err := p.mydb.Model(model.ProposalTbl{}).
		WithContext(ctx).
		Clauses(clause.OnConflict{
			// Specify the columns that should trigger the conflict (in this case, proposal_id).
			Columns: []clause.Column{
				{Name: "proposal_id"},
			},
			// Define the update behavior when a conflict occurs.
			// Here, it specifies the columns to update with the new values from the input.
			DoUpdates: clause.AssignmentColumns([]string{
				"creator",
				"github_name",
				"start_time",
				"end_time",
				"timestamp",
				"chain_id",
				"title",
				"content",
				"block_number",
				"snapshot_day",
				"token_holder_percentage",
				"snapshot_block_height",
				"sp_percentage",
				"client_percentage",
				"developer_percentage",
			}),
		}).Create(&in).Error; err != nil {
		return 0, fmt.Errorf("create proposal error: %w", err)
	}

	return in.ID, nil
}

// UpdateProposal updates the specified proposal in the database.
func (p *ProposalRepoImpl) UpdateProposal(ctx context.Context, in *model.ProposalTbl) error {
	// Start a new database transaction and set the context for it.
	err := p.mydb.Model(model.ProposalTbl{}).
		WithContext(ctx).
		// Specify the condition to find the proposal by its ID.
		Where("id = ?", in.ID).
		// Update the specified columns with new values from the input proposal.
		UpdateColumns(map[string]any{
			"counted":                  in.Counted,
			"approve_percentage":       in.ProposalResult.ApprovePercentage,
			"reject_percentage":        in.ProposalResult.RejectPercentage,
			"total_sp_power":           in.TotalSpPower,
			"total_token_holder_power": in.TotalTokenHolderPower,
			"total_client_power":       in.TotalClientPower,
			"total_developer_power":    in.TotalDeveloperPower,
			"updated_at":               time.Now(),
		}).Error

	return err
}

// UpdateProposalGitHubName implements service.ProposalRepo.
func (p *ProposalRepoImpl) UpdateProposalGitHubName(ctx context.Context, createrAddress, githubName string) error {
	if err := p.mydb.Model(model.ProposalTbl{}).
		WithContext(ctx).
		Where("creator = ?", createrAddress).
		UpdateColumns(map[string]any{
			"github_name": githubName,
			"updated_at":  time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("update proposal github name error: %w", err)
	}

	return nil
}

// GetUnCountedProposalList retrieves a list of proposals based on the provided network ID and timestamp.
// It queries the database for proposals with the following conditions:
// 1. Matching network ID.
// 2. Expiration time before or equal to the provided timestamp.
// 3. Uncounted vote counting proposals.
// It returns the list of proposals and any error encountered during the database query.
func (p *ProposalRepoImpl) GetUncountedProposalList(ctx context.Context, chainId int64, timestamp int64) ([]model.ProposalTbl, error) {
	var proposalList []model.ProposalTbl
	tx := p.mydb.Model(model.ProposalTbl{}).
		WithContext(ctx).
		Where("chain_id = ?", chainId).
		Where("end_time <= ?", timestamp).
		Where("counted = ?", constant.ProposalCreate).
		Order("id desc").Find(&proposalList)

	return proposalList, tx.Error
}

// buildProposalBaseQuery constructs the base query for counting and listing proposals based on the given request.
func (p *ProposalRepoImpl) buildProposalBaseQuery(req api.ProposalListReq) (queryCount *gorm.DB, queryList *gorm.DB) {
	// Define a function to apply common conditions to the query.
	baseCondition := func(query *gorm.DB) {
		// If a search key is provided, add a condition to filter by title.
		if req.SearchKey != "" {
			safeSearchKey := "%" + strings.ReplaceAll(
				strings.ReplaceAll(req.SearchKey, "%", "\\%"),
				"_", "\\_",
			) + "%"
			query.Where("title LIKE ?", safeSearchKey)
		}

		// If a chain ID is provided, add a condition to filter by chain ID.
		if req.ChainId != 0 {
			query.Where("chain_id = ?", req.ChainId)
		}
	}

	// Initialize the query for counting proposals.
	queryCount = p.mydb.Model(model.ProposalTbl{})

	// Initialize the query for listing proposals with ordering, pagination, and common conditions.
	queryList = p.mydb.Model(model.ProposalTbl{}).
		Order("end_time desc").
		Limit(req.PageSize).
		Offset(int(req.Offset()))

	baseCondition(queryCount)
	baseCondition(queryList)

	// Return the constructed queries for counting and listing proposals.
	return queryCount, queryList
}

// applyStatusFilter applies a status filter to the given GORM query based on the provided status.
func (p *ProposalRepoImpl) applyStatusFilter(query *gorm.DB, status int) {
	// Get the current Unix timestamp to compare with proposal times.
	now := time.Now().Unix()
	// Use a switch statement to handle different status cases.
	switch status {
	case constant.ProposalStatusPending:
		// For pending proposals, filter where 'counted' is not counted or valid and 'start_time' is in the future.
		query.Where("counted = ? AND start_time > ?", constant.ProposalCreate, now)
	case constant.ProposalStatusInProgress:
		// For in-progress proposals, filter where 'counted' is not counted or valid, 'start_time' is in the past, and 'end_time' is in the future.
		query.Where("counted = ? AND start_time < ? AND end_time > ?",
			constant.ProposalCreate, now, now)
	case constant.ProposalStatusCounting:
		// For counting proposals, filter where 'counted' is not counted or valid and 'end_time' is in the past.
		query.Where("counted = ? AND end_time < ?", constant.ProposalCreate, now)
	case constant.ProposalStatusCompleted:
		// For completed proposals, filter where 'counted' is counted or invalid.
		query.Where("counted = ?", constant.ProposalCounted)
	}
}
