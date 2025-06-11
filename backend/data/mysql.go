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

package data

import (
	"context"
	"errors"
	"fmt"

	goMysql "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"powervoting-server/config"
	"powervoting-server/constant"
	"powervoting-server/model"
)

// InitMysql initializes the MySQL database connection and performs necessary migrations.
// It configures the database connection using the provided MySQL configuration.
// After establishing the connection, it auto-migrates the required database tables:
// Proposal, Vote, VoteResult, VoteCompleteHistory, Dict, and VotePowerTbl.
// Additionally, it checks if the Dict table contains an entry for proposal start key,
// and creates one if not found.
// Any error encountered during the initialization process is logged.
func NewMysql() *gorm.DB {
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
			SingularTable: true, // Singular table
		},
	})
	if err != nil {
		zap.L().Error("connect mysql error: ", zap.Error(err))
		panic(err)
	}

	db.AutoMigrate(&model.ProposalTbl{})
	db.AutoMigrate(&model.VoteTbl{})
	db.AutoMigrate(&model.ProposalDraftTbl{})
	db.AutoMigrate(&model.SyncEventTbl{})
	db.AutoMigrate(&model.VoterInfoTbl{})
	db.AutoMigrate(&model.GithubRepos{})
	db.AutoMigrate(&model.FipProposalTbl{})
	db.AutoMigrate(&model.FipProposalVoteTbl{})
	db.AutoMigrate(&model.FipEditorTbl{})

	return db
}

// IsDuplicateEntryError returns true if the error is a duplicate entry error
func IsDuplicateEntryError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	var mysqlErr *goMysql.MySQLError
	if errors.As(err, &mysqlErr) {
		if mysqlErr.Number == constant.MysqlDuplicateEntryErrorCode {
			return true
		}
	}

	return false
}

// General mysql counter function, counting the number of rows that meet the conditions.
// quer y is a DB instance that has already built the query conditions.
func ExecuteCountQuery(ctx context.Context, query *gorm.DB) (int64, error) {
	var count int64
	if err := query.WithContext(ctx).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
