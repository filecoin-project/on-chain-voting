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
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"powervoting-server/constant"
	"powervoting-server/data"
	"powervoting-server/model"
	"powervoting-server/model/api"
	"powervoting-server/service"
)

type FipRepoImpl struct {
	mydb *gorm.DB
}

var _ service.FipRepo = (*FipRepoImpl)(nil)

func NewFipRepo(mydb *gorm.DB) *FipRepoImpl {
	return &FipRepoImpl{mydb: mydb}
}

// CreateFipProposal creates or updates a Fip record in the database.
// If a record with the same `proposal_id` already exists, it updates the specified fields.
// Otherwise, it inserts a new record.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - in: A pointer to the `FipTbl` struct containing the Fip data to be created or updated.
//
// Returns:
//   - int64: The ID of the created or updated Fip record.
//   - error: An error object if the operation fails.
func (f *FipRepoImpl) CreateFipProposal(ctx context.Context, in *model.FipProposalTbl) (int64, error) {
	// Use the database context to execute the operation with the provided context.
	// The `Clauses` method is used to specify the behavior on conflict.
	// In this case, if a record with the same `proposal_id` exists, the specified fields will be updated.
	if err := f.mydb.Model(&model.FipProposalTbl{}).WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "proposal_id"}, {Name: "chain_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"proposal_type",
				"creator",
				"candidate_address",
				"candidate_info",
				"updated_at",
				"block_number",
				"timestamp",
			}),
		},
	).Create(in).Error; err != nil {
		// If an error occurs during the operation, return 0 and the error.
		return 0, err
	}

	// If the operation is successful, return the ID of the created or updated record.
	return in.ID, nil
}

// GetFipProposalListWithPagination retrieves a paginated list of FipEditor records from the database based on the provided request parameters.
// The list is filtered by `proposal_type`, ordered by `proposal_id` in descending order, and paginated using `PageSize` and `Offset`.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: The request object containing filtering and pagination parameters (e.g., `FipType`, `PageSize`, and `Offset`).
//
// Returns:
//   - []model.FipProposalTbl: A slice of `FipProposalTbl` structs representing the retrieved records.
//   - error: An error object if the operation fails.
func (f *FipRepoImpl) GetFipProposalListWithPagination(ctx context.Context, req api.FipProposalListReq) ([]model.FipProposalVoted, int64, error) {
	var list []model.FipProposalVoted

	baseQuery := f.mydb.Model(&model.FipProposalTbl{}).
		Where("proposal_type = ?", req.ProposalType). // Filter by `proposal_type`
		// The code is using a method called `Where` to filter data based on the condition `chain_id = req.ChainId`. This is likely part of a database query or data manipulation operation in a Go program. The `req.ChainId` variable is being used as a parameter to filter the data.
		Where("chain_id = ?", req.ChainId).
		Where("status = ?", constant.FipProposalUnpass)
	count, err := data.ExecuteCountQuery(ctx, baseQuery)
	if err != nil {
		return nil, 0, err
	}
	// Use the database context to execute the query with the provided context.
	// Filter the records by `proposal_type`, order them by `proposal_id` in descending order,
	// and apply pagination using `Limit` and `Offset`.
	subQuery := f.mydb.Model(&model.FipProposalVoteTbl{}).
		WithContext(ctx).
		Select("COALESCE(JSON_ARRAYAGG(voter), JSON_ARRAY())").
		Where("proposal_id = ?", gorm.Expr("fip_proposal_tbl.proposal_id")).
		Where("is_remove = ?", constant.FipEditorValid)
	if err := baseQuery. // Filter by `chain_id`
				Select("fip_proposal_tbl.*, (?) AS voters", subQuery).
				Order("proposal_id desc").      // Order by `proposal_id` in descending order
				Limit(req.PageSize).            // Limit the number of records returned to `PageSize`
				Offset(int(req.Offset())).      // Skip records based on the calculated offset
				Find(&list).Error; err != nil { // Execute the query and store the results in `list`
		return nil, 0, err // Return nil and the error if the query fails
	}

	// Return the retrieved list of records and nil error if the operation is successful
	return list, count, nil
}

// CreateFipProposalVote creates a new FipProposalVoter record in the database. If a record with the same `voter` and `proposal_id`
// already exists, the operation does nothing (no update or insertion is performed).
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - in: A pointer to the `FipProposalVoteTbl` struct containing the FipProposalVoter data to be created.
//
// Returns:
//   - int64: The ID of the created or existing FipProposalVoter record.
//   - error: An error object if the operation fails. Returns `nil` if the operation is successful.
func (f *FipRepoImpl) CreateFipProposalVote(ctx context.Context, in *model.FipProposalVoteTbl) (int64, error) {
	// Use the database context to execute the operation with the provided context.
	// The `Clauses` method specifies the behavior on conflict. In this case, if a record with the same
	// `voter` and `proposal_id` already exists, the operation does nothing (`DoNothing: true`).
	if err := f.mydb.Model(&model.FipProposalVoteTbl{}).
		WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "voter"}, {Name: "proposal_id"}}, // Conflict columns
				DoUpdates: clause.AssignmentColumns([]string{
					"is_remove",
					"block_number",
					"timestamp",
				}),
			},
		).
		Create(in).Error; err != nil { // Attempt to create the record
		return 0, err // Return 0 and the error if the operation fails
	}

	// Return the ID of the created or existing record and nil error if the operation is successful
	return in.ID, nil
}

// CreateFipEditor creates a new FipVoter record in the database. This function inserts the provided
// FipVoter data into the database and returns the ID of the newly created record.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - in: A pointer to the `FipEditorTbl` struct containing the FipVoter data to be created.
//
// Returns:
//   - int64: The ID of the newly created FipVoter record.
//   - error: An error object if the operation fails. Returns `nil` if the operation is successful.
func (f *FipRepoImpl) CreateFipEditor(ctx context.Context, in *model.FipEditorTbl) (int64, error) {
	// Use the database context to execute the create operation with the provided context.
	// The `Create` method inserts the provided `FipEditorTbl` data into the database.
	if err := f.mydb.Model(&model.FipEditorTbl{}).
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "editor"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"is_remove",
				"updated_at",
			}),
		}).
		Create(in).Error; err != nil {
		return 0, err // Return 0 and the error if the creation fails
	}

	// Return the ID of the newly created record and nil error if the operation is successful
	return in.ID, nil
}

// GetFipEditorCount retrieves the total count of FipVoter records in the database for a specific `chain_id`.
// This function is typically used to determine the number of votes cast on a specific blockchain network.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - chainId: The unique identifier of the blockchain network for which the vote count is to be retrieved.
//
// Returns:
//   - int64: The total count of FipVoter records for the specified `chain_id`.
//   - error: An error object if the operation fails. Returns `nil` if the operation is successful.
func (f *FipRepoImpl) GetFipEditorCount(ctx context.Context, chainId int64) (int64, error) {
	// Build the query to filter FipVoter records by the specified `chain_id`.
	query := f.mydb.Model(&model.FipEditorTbl{}).
		WithContext(ctx).
		Where("chain_id = ? and is_remove = ?", chainId, constant.FipEditorValid)

	// Execute the count query using a helper function `ExecuteCountQuery`.
	return data.ExecuteCountQuery(ctx, query)
}

// GetFipProposalVoteCount implements service.FipRepo.
func (f *FipRepoImpl) GetFipProposalVoteCount(ctx context.Context, chainId, proposalId int64) (int64, error) {
	query := f.mydb.Model(&model.FipProposalVoteTbl{}).
		WithContext(ctx).
		Where("chain_id = ? and proposal_id = ? and is_remove = ?", chainId, proposalId, constant.FipEditorValid)

	// Execute the count query using a helper function `ExecuteCountQuery`.
	return data.ExecuteCountQuery(ctx, query)
}

// UpdateStatusAndGetFipProposal updates the status of a FipProposal in the database. It sets the `status` field
// to 1 for the record matching the specified `proposal_id` and `chain_id`.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - proposalId: The unique identifier of the proposal for which the status is to be updated.
//   - chainId: The unique identifier of the blockchain network associated with the proposal.
//
// Returns:
//   - error: An error object if the operation fails. Returns `nil` if the operation is successful.
func (f *FipRepoImpl) UpdateStatusAndGetFipProposal(ctx context.Context, proposalId int64, chainId int64) (*model.FipProposalTbl, error) {
	var updatedProposal model.FipProposalTbl

	err := f.mydb.Transaction(func(tx *gorm.DB) error {
		baseQuery := tx.Model(&model.FipProposalTbl{}).
			WithContext(ctx).
			Where("proposal_id = ? AND chain_id = ?", proposalId, chainId)

		if err := baseQuery.
			Update("status", constant.FipProposalPass).
			Error; err != nil {
			return fmt.Errorf("update status failed: %w", err)
		}

		if err := baseQuery.
			First(&updatedProposal).
			Error; err != nil {
			return fmt.Errorf("fetch updated proposal failed: %w", err)
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("proposal not found after update: %w", err)
		}
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return &updatedProposal, nil
}

// UpdateFipProposalVoteByAddress updates the vote status of a specific proposal for a given voter address.
// It sets the `is_remove` field to a predefined invalid state for the record matching the specified
// `proposal_id` and voter `address`.
//
// Parameters:
//   - ctx:        The context for managing request-scoped values, cancellation signals, and deadlines.
//   - proposalId: The unique identifier of the proposal for which the vote status is to be updated.
//   - address:    The blockchain address of the voter whose voting record should be updated.
//
// Returns:
//   - error:      An error object if the database operation fails. Returns `nil` if the update is successful.
func (f *FipRepoImpl) UpdateFipProposalVoteByAddress(ctx context.Context, proposalId int64, address string) error {
	if err := f.mydb.Model(&model.FipProposalVoteTbl{}).
		WithContext(ctx).
		Where("voter = ? AND proposal_id = ?", address, proposalId).
		UpdateColumns(map[string]any{
			"is_remove":  constant.FipEditorInvalid,
			"updated_at": time.Now(),
		}).
		Error; err != nil {
		return err
	}

	return nil
}

// UpdateFipEditorByAddress updates the status of a FipEditor in the database. It sets the `is_remove` field
// to a predefined invalid state for the record matching the specified `editor` address.
//
// Parameters:
//   - ctx:     The context for managing request-scoped values, cancellation signals, and deadlines.
//   - address: The blockchain address of the editor whose status is to be updated.
//
// Returns:
//   - error:   An error object if the database operation fails. Returns `nil` if the update is successful.
func (f *FipRepoImpl) UpdateFipEditorByAddress(ctx context.Context, address string) error {
	if err := f.mydb.Model(&model.FipEditorTbl{}).
		WithContext(ctx).
		Where("editor = ?", address).
		UpdateColumns(map[string]any{
			"is_remove": constant.FipEditorInvalid,
		}).
		Error; err != nil {
		return err
	}

	return nil
}

// GetUnpassFipProposalList retrieves a list of FipProposals that have not yet passed for a specific blockchain network.
// It queries the database for records matching the specified `chain_id` and a status indicating the proposal
// has not passed (as defined by `constant.FipProposalPass`).
//
// Parameters:
//   - ctx:    The context for managing request-scoped values, cancellation signals, and deadlines.
//   - chainId: The unique identifier of the blockchain network for which to retrieve the proposals.
//
// Returns:
//   - []model.FipProposalTbl: A slice of FipProposal records that match the query criteria.
//   - error:                  An error object if the database query fails. Returns `nil` if the query is successful.
func (f *FipRepoImpl) GetUnpassFipProposalList(ctx context.Context, chainId int64) ([]model.FipProposalTbl, error) {
	var fipProposals []model.FipProposalTbl
	if err := f.mydb.Model(&model.FipProposalTbl{}).
		WithContext(ctx).
		Where("chain_id = ? AND status = ?", chainId, constant.FipProposalUnpass).
		Find(&fipProposals).Error; err != nil {
		return nil, err
	}

	return fipProposals, nil
}

func (f *FipRepoImpl) GetValidFipEditorList(ctx context.Context, req api.FipEditorListReq) ([]model.FipEditorTbl, error) {
	var list []model.FipEditorTbl
	if err := f.mydb.Model(&model.FipEditorTbl{}).Where("is_remove = ? and chain_id = ?", constant.FipEditorValid, req.ChainId).WithContext(ctx).Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}
