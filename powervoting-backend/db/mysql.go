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

package db

import (
	"fmt"
	"powervoting-server/config"
	"powervoting-server/constant"
	"powervoting-server/model"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Mysql struct {
	*gorm.DB
}

var Engine *Mysql

// InitMysql initializes the MySQL database connection and performs necessary migrations.
// It configures the database connection using the provided MySQL configuration.
// After establishing the connection, it auto-migrates the required database tables:
// Proposal, Vote, VoteResult, VoteCompleteHistory, Dict, and VotePower.
// Additionally, it checks if the Dict table contains an entry for proposal start key,
// and creates one if not found.
// Any error encountered during the initialization process is logged.
func InitMysql() {
	var err error
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       fmt.Sprintf("%s:%s@tcp(%s)/power-voting-filecoin?charset=utf8&parseTime=True&loc=Local", config.Client.Mysql.Username, config.Client.Mysql.Password, config.Client.Mysql.Url),
		DefaultStringSize:         256,   // string size
		DisableDatetimePrecision:  true,  // datetime precision is disabled. Databases earlier than MySQL 5.6 do not support dateTime precision
		DontSupportRenameIndex:    true,  // Rename indexes by deleting and creating new indexes. Databases before MySQL 5.7 and MariaDB do not support rename indexes
		DontSupportRenameColumn:   true,  // Rename columns with 'change'. Databases prior to MySQL 8 and MariaDB do not support renaming columns
		SkipInitializeWithVersion: false, // This parameter is automatically configured based on the current MySQL version
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tbl_", // Add tbl_ to the table prefix
			SingularTable: true,   // Singular table
		},
	})
	if err != nil {
		zap.L().Error("connect mysql error: ", zap.Error(err))
		return
	}

	db.AutoMigrate(&model.Proposal{})
	db.AutoMigrate(&model.Vote{})
	db.AutoMigrate(&model.VoteResult{})
	db.AutoMigrate(&model.VoteCompleteHistory{})
	db.AutoMigrate(&model.Dict{})
	db.AutoMigrate(&model.VotePower{})
	db.AutoMigrate(&model.ProposalDraft{})

	var count int64
	db.Model(model.Dict{}).Where("name", constant.ProposalStartKey).Count(&count)
	if count == 0 {
		db.Model(model.Dict{}).Create(&model.Dict{
			Name:  constant.ProposalStartKey,
			Value: "1",
		})
	}

	Engine = &Mysql{db}
}

// GetProposalList retrieves a list of proposals based on the provided network ID and timestamp.
// It queries the database for proposals with the following conditions:
// 1. Matching network ID.
// 2. Expiration time before or equal to the provided timestamp.
// 3. Status set to 0 (active).
// It returns the list of proposals and any error encountered during the database query.
func (m *Mysql) GetProposalList(network int64, timestamp int64) ([]model.Proposal, error) {
	var proposalList []model.Proposal
	tx := m.Model(model.Proposal{}).Where("network = ? and exp_time <= ? and status = ?", network, timestamp, constant.ProposalStatusPending).Order("id desc").Find(&proposalList)
	return proposalList, tx.Error
}

// GetVoteList retrieves a list of votes based on the provided network ID and proposal ID.
// It queries the database for votes with the following conditions:
// 1. Matching network ID.
// 2. Matching proposal ID.
// It returns the list of votes and any error encountered during the database query.
func (m *Mysql) GetVoteList(network, proposalId int64) ([]model.Vote, error) {
	var proposalList []model.Vote
	tx := m.Model(model.Vote{}).Where("network", network).Where("proposal_id", proposalId).Find(&proposalList)
	return proposalList, tx.Error
}

// VoteResult updates the database with the provided vote result and associated history.
// It begins a transaction to ensure atomicity and consistency of database operations.
// It creates a new record in the VoteCompleteHistory table for the given history.
// It creates multiple records in the VoteResult table in batches.
// It updates the status of the proposal with the provided ID to indicate that it has been voted on.
// Any error encountered during the transaction is logged.
func (m *Mysql) VoteResult(proposalId int64, history model.VoteCompleteHistory, result []model.VoteResult) (int64, error) {
	err := m.Transaction(func(tx *gorm.DB) error {
		create := tx.Model(model.VoteCompleteHistory{}).Create(&history)
		if create.Error != nil {
			zap.L().Error("batch create error: ", zap.Error(create.Error))
			return create.Error
		}
		create = tx.Model(model.VoteResult{}).CreateInBatches(result, len(result))
		if create.Error != nil {
			zap.L().Error("batch create error: ", zap.Error(create.Error))
			return create.Error
		}
		update := tx.Model(model.Proposal{}).Where("id", proposalId).Update("status", constant.ProposalStatusCompleted)
		if update.Error != nil {
			zap.L().Error("update proposal status error: ", zap.Error(update.Error))
			return update.Error
		}
		return nil
	})

	return history.Id, err
}

func (m *Mysql) GetDict(key string) (*model.Dict, error) {
	var dict model.Dict
	if err := m.Model(model.Dict{}).Where("name", key).Find(&dict).Error; err != nil {
		zap.L().Error("Get vote start index error: ", zap.Error(err))
		return nil, err
	}
	return &dict, nil
}

func (m *Mysql) CreateDict(in *model.Dict) (int64, error) {
	if err := m.Model(model.Dict{}).Create(in).Error; err != nil {
		zap.L().Error("create vote dict error: ", zap.Error(err))
		return 0, err
	}

	return in.Id, nil
}

func (m *Mysql) CountVotes(filter map[string]any) (int64, error) {
	var count int64
	if err := m.Model(model.Vote{}).Where(filter).Count(&count).Error; err != nil {
		zap.L().Error("get vote count error: ", zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (m *Mysql) UpdateVoteInfo(filter map[string]any, in string) error {
	if err := m.Model(model.Vote{}).Where(filter).Update("vote_info", in).Error; err != nil {
		zap.L().Error("update vote error", zap.Error(err))
		return err
	}

	return nil
}

func (m *Mysql) CreateVote(in *model.Vote) (int64, error) {
	if err := m.Model(model.Vote{}).Create(in).Error; err != nil {
		zap.L().Error("create vote error: ", zap.Error(err))
		return 0, err
	}

	return in.Id, nil
}

func (m *Mysql) UpdateDict(key string, value string) error {
	if err := m.Model(model.Dict{}).Where("name", key).Update("value", value).Error; err != nil {
		zap.L().Error("update vote start key error: ", zap.Error(err))
		return err
	}

	return nil
}

func (m *Mysql) CountProposal(filter map[string]any) (int64, error) {
	var count int64
	if err := m.Model(model.Proposal{}).Where(filter).Count(&count).Error; err != nil {
		zap.L().Error("get proposal count error: ", zap.Error(err))
		return 0, err
	}

	return count, nil
}

func (m *Mysql) CreateProposal(in *model.Proposal) (int64, error) {
	if err := m.Model(model.Proposal{}).Create(in).Error; err != nil {
		zap.L().Error("create proposal error: ", zap.Error(err))
		return 0, err
	}
	return in.Id, nil
}

func (m *Mysql) UpdateProposal(in *model.Proposal) (int64, error) {
	var proposal model.Proposal
	m.Model(model.Proposal{}).Where("cid", in.Cid).Take(&proposal)
	if proposal.Id == 0 {
		if _, err := m.CreateProposal(in); err != nil {
			zap.L().Error("create proposal error: ", zap.Error(err))
			return 0, err
		}
	} else {
		if proposal.Status == constant.ProposalStatusStoring && in.Status == constant.ProposalStatusPending {
			in.Status = constant.ProposalStatusPending
		}
		if err := m.Model(model.Proposal{}).Where("cid", in.Cid).Updates(in).Error; err != nil {
			zap.L().Error("create proposal error: ", zap.Error(err))
			return 0, err
		}
	}
	return in.Id, nil
}
