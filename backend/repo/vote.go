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

	"powervoting-server/model"
	"powervoting-server/service"
)

var _ service.VoteRepo = (*VoteRepoImpl)(nil)

type VoteRepoImpl struct {
	mydb *gorm.DB
}

func NewVoteRepo(mydb *gorm.DB) *VoteRepoImpl {
	return &VoteRepoImpl{
		mydb: mydb,
	}
}

func (v *VoteRepoImpl) BatchUpdateVotes(ctx context.Context, votes []model.VoteTbl) error {
	if len(votes) == 0 {
		return nil
	}

	tx := v.mydb.Begin().WithContext(ctx)
	if tx.Error != nil {
		return fmt.Errorf("begin tx error: %w", tx.Error)
	}

	for _, vote := range votes {
		if err := tx.Model(model.VoteTbl{}).
			WithContext(ctx).
			Where("proposal_id = ? and address = ?", vote.ProposalId, vote.Address).
			UpdateColumns(map[string]any{
				"vote_result":        vote.VoteResult,
				"sp_power":           vote.SpPower,
				"client_power":       vote.ClientPower,
				"developer_power":    vote.DeveloperPower,
				"token_holder_power": vote.TokenHolderPower,
				"updated_at":         time.Now(),
			}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("update vote error: %w", err)
		}
	}

	return tx.Commit().Error
}

// CreateVote creates a new vote record in the database.
func (v VoteRepoImpl) CreateVote(ctx context.Context, in *model.VoteTbl) (int64, error) {
	// Use OnConflict to implement the existence of the OnConflict implementation,
	// update if it does not exist, insert it if it does not exist.
	err := v.mydb.Model(model.VoteTbl{}).
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "proposal_id"},
				{Name: "address"},
				{Name: "chain_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"vote_encrypted",
				"timestamp",
				"block_number",
				"updated_at",
			}),
		}).
		Create(in).Error
	if err != nil {
		return 0, fmt.Errorf("create or update vote error: %w", err)
	}

	return in.BaseField.ID, nil
}

// GetVoteList retrieves a list of votes based on the provided network ID and proposal ID.
// It queries the database for votes with the following conditions:
// 1. Matching network ID.
// 2. Matching proposal ID.
// It returns the list of votes and any error encountered during the database query.
func (v *VoteRepoImpl) GetVoteList(ctx context.Context, chainId, proposalId int64, counted bool) ([]model.VoteTbl, error) {
	var proposalList []model.VoteTbl
	tx := v.mydb.Model(model.VoteTbl{}).
		WithContext(ctx).
		Where("chain_id = ? and proposal_id = ?", chainId, proposalId)
	if counted {
		tx = tx.Where("vote_result != ''")
	} else {
		tx = tx.Where("vote_result = ''")
	}

	tx.Find(&proposalList)
	return proposalList, tx.Error
}

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
func (v *VoteRepoImpl) CreateVoterAddress(ctx context.Context, in *model.VoterInfoTbl) (int64, error) {
	// Use the `OnConflict` clause to handle duplicate addresses
	err := v.mydb.Model(model.VoterInfoTbl{}).
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "address"}, // Conflict resolution based on the `address` column
				{Name: "chain_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"updated_at"}), // Update `update_height` if the address exists
		}).
		Create(in).Error // Perform the create or update operation

	// Handle errors during the operation
	if err != nil {
		return 0, fmt.Errorf("create or update voter address error: %w", err)
	}

	// Return the ID of the created or updated record
	return in.BaseField.ID, nil
}

// GetAllVoterAddresss retrieves a list of voter addresses that were created after a specified block height.
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
func (v *VoteRepoImpl) GetAllVoterAddresss(ctx context.Context, chainId int64) ([]model.VoterInfoTbl, error) {
	var proposalList []model.VoterInfoTbl

	// Query the database for voter addresses created after the specified height
	if err := v.mydb.Model(model.VoterInfoTbl{}).
		WithContext(ctx).
		Where("chain_id = ?", chainId).
		Find(&proposalList).Error; err != nil {
		return proposalList, err
	}

	return proposalList, nil
}

// UpdateVoterByGistInfo implements service.VoteRepo.
func (v *VoteRepoImpl) UpdateVoterByGistInfo(ctx context.Context, in *model.VoterInfoTbl) error {
	if err := v.mydb.Model(model.VoterInfoTbl{}).
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "address"},
				{Name: "chain_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"gist_id",
				"github_id",
				"gist_info",
				"block_number",
				"timestamp",
				"updated_at",
			}),
		}).Create(in).Error; err != nil {
		return err
	}

	if err := v.mydb.Model(model.VoterInfoTbl{}).
		WithContext(ctx).
		Where("address <> ? and github_id = ?", in.Address, in.GithubId).
		UpdateColumns(map[string]any{
			"gist_id":    "",
			"github_id":  "",
			"gist_info":  "",
			"updated_at": time.Now(),
		}).Error; err != nil {
		return err
	}

	return nil
}

// UpdateVoterByMinerInfo implements service.VoteRepo.
func (v *VoteRepoImpl) UpdateVoterByMinerInfo(ctx context.Context, in *model.VoterInfoTbl) error {
	if err := v.mydb.Model(model.VoterInfoTbl{}).
		WithContext(ctx).
		Where("address = ?", in.Address).UpdateColumns(map[string]any{
		"miner_ids":  in.MinerIds,
		"owner_id":   in.OwnerId,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return err
	}

	return nil
}

// GetVoterInfoByAddress implements service.VoteRepo.
func (v *VoteRepoImpl) GetVoterInfoByAddress(ctx context.Context, address string) (*model.VoterInfoTbl, error) {
	var voterInfo model.VoterInfoTbl
	if err := v.mydb.Model(model.VoterInfoTbl{}).
		WithContext(ctx).
		Where("address = ?", address).
		First(&voterInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("voter address not found: %s", address)
		}
	}

	return &voterInfo, nil
}
