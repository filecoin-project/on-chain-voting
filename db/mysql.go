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

var Engine *gorm.DB

func InitMysql() {
	var err error
	Engine, err = gorm.Open(mysql.New(mysql.Config{
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

	Engine.AutoMigrate(&model.Proposal{})
	Engine.AutoMigrate(&model.Vote{})
	Engine.AutoMigrate(&model.VoteResult{})
	Engine.AutoMigrate(&model.VoteHistory{})
	Engine.AutoMigrate(&model.Dict{})
	Engine.AutoMigrate(&model.VotePower{})

	var count int64
	Engine.Model(model.Dict{}).Where("name", constant.ProposalStartKey).Count(&count)
	if count == 0 {
		Engine.Model(model.Dict{}).Create(&model.Dict{
			Name:  constant.ProposalStartKey,
			Value: "1",
		})
	}
}

func GetProposalList(network int64, timestamp int64) ([]model.Proposal, error) {
	var proposalList []model.Proposal
	tx := Engine.Model(model.Proposal{}).Where("network = ? and exp_time <= ? and status = 0", network, timestamp).Find(&proposalList)
	return proposalList, tx.Error
}

func GetVoteList(network, proposalId int64) ([]model.Vote, error) {
	var proposalList []model.Vote
	tx := Engine.Model(model.Vote{}).Where("network", network).Where("proposal_id", proposalId).Find(&proposalList)
	return proposalList, tx.Error
}

func VoteResult(proposalId int64, history model.VoteCompleteHistory, result []model.VoteResult) {
	Engine.Transaction(func(tx *gorm.DB) error {
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
		update := tx.Model(model.Proposal{}).Where("id", proposalId).Update("status", 1)
		if update.Error != nil {
			zap.L().Error("update proposal status error: ", zap.Error(update.Error))
			return update.Error
		}
		return nil
	})
}
