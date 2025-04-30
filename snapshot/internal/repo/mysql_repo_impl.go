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

	"gorm.io/gorm/clause"

	"power-snapshot/constant"
	"power-snapshot/internal/data"
	models "power-snapshot/internal/model"
)

type MysqlRepoImpl struct {
	db *data.Mysql
}

func NewMysqlRepoImpl(db *data.Mysql) *MysqlRepoImpl {
	return &MysqlRepoImpl{db: db}
}

// CreateSnapshotBackup creates a backup entry for a snapshot in the database.
// It takes a context and a SnapshotBackupTbl struct as input and returns an error.
// The function attempts to insert the snapshot backup record into the database.
// If a conflict occurs on the combination of `cid` and `day` columns, the operation is ignored (no error is returned).
// If any other database error occurs during the insertion, it is returned to the caller.
// On success, it returns nil.
func (m *MysqlRepoImpl) CreateSnapshotBackup(ctx context.Context, in models.SnapshotBackupTbl) error {
	if err := m.db.Model(in).
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "day"},
			},
			DoNothing: true,
		}).
		Create(&in).Error; err != nil {
		return err
	}

	return nil
}

// GetSnapshotBackupList retrieves a list of snapshot backups for a specific chain ID.
// It takes a context and a chain ID as input and returns a slice of SnapshotBackupTbl or an error.
// The function queries the database for snapshot backups where the chain ID matches the provided value
// and the status is less than the constant `SnapshotBackupRetry`.
// If the query fails, it returns the corresponding error.
// On success, it returns the list of snapshot backups.
func (m *MysqlRepoImpl) GetSnapshotBackupList(ctx context.Context, chainId int64) ([]models.SnapshotBackupTbl, error) {
	var snapshotBackups []models.SnapshotBackupTbl
	if err := m.db.Model(models.SnapshotBackupTbl{}).
		WithContext(ctx).
		Where("chain_id = ? and status < ?", chainId, constant.RetryCount).
		Find(&snapshotBackups).Error; err != nil {
		return nil, err
	}

	return snapshotBackups, nil
}



// UpdateSnapshotBackup updates the status and CID of a snapshot backup entry in the database.
// It takes a context and a SnapshotBackupTbl struct as input and returns an error.
// The function updates the snapshot backup record in the database where the `day` matches the provided value.
// It sets the `status` and `cid` fields to the values provided in the input struct.
// If the update operation fails, it returns the corresponding error.
// On success, it returns nil.
func (m *MysqlRepoImpl) UpdateSnapshotBackup(ctx context.Context, in models.SnapshotBackupTbl) error {
	if err := m.db.Model(in).
		WithContext(ctx).
		Where("day = ?", in.Day).
		Updates(models.SnapshotBackupTbl{Status: in.Status, Cid: in.Cid}).Error; err != nil {
		return err
	}

	return nil
}
