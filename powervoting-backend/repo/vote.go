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
				"block_number":       vote.BlockNumber,
				"timestamp":          vote.Timestamp,
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
			DoUpdates: clause.AssignmentColumns([]string{"vote_encrypted"}),
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
