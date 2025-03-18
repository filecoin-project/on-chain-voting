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

	"powervoting-server/data"
	"powervoting-server/model"
	"powervoting-server/service"
)

var _ service.SyncRepo = (*SyncRepoImpl)(nil)

type SyncRepoImpl struct {
	mydb *gorm.DB
}

func NewSyncRepo(mydb *gorm.DB) *SyncRepoImpl {
	return &SyncRepoImpl{mydb: mydb}
}

// CreateSyncEventInfo creates a new synchronization event in the database.
// It takes a context and a pointer to a SyncEventTbl struct as input.
// It returns the ID of the created event and an error if any occurred.
func (s *SyncRepoImpl) CreateSyncEventInfo(ctx context.Context, in *model.SyncEventTbl) (int64, error) {
	if err := s.mydb.Model(model.SyncEventTbl{}).
		WithContext(ctx).
		Create(in).Error; err != nil {
		// Check if the error is due to a duplicate entry.
		if data.IsDuplicateEntryError(err) {
			// If it's a duplicate entry, return 0 for the ID and no error.
			return 0, nil
		}

		return 0, fmt.Errorf("create sync event error: %w", err)
	}

	return in.Id, nil
}

// GetSyncEventInfo retrieves the synchronization event information for a given contract address.
// It returns a pointer to the SyncEventTbl struct containing the event information and an error if any occurred.
func (s *SyncRepoImpl) GetSyncEventInfo(ctx context.Context, addr string) (*model.SyncEventTbl, error) {
	var syncEvent model.SyncEventTbl
	tx := s.mydb.Model(model.SyncEventTbl{}).
		WithContext(ctx).
		Where("contract_address = ?", addr).
		Take(&syncEvent)

	return &syncEvent, tx.Error
}

// UpdateSyncEventInfo updates the synced height of a sync event in the database.
func (s *SyncRepoImpl) UpdateSyncEventInfo(ctx context.Context, addr string, height int64) error {
	err := s.mydb.Model(model.SyncEventTbl{}).
		WithContext(ctx).
		Where("contract_address = ?", addr).
		Updates(model.SyncEventTbl{SyncedHeight: height}).Error
	if err != nil {
		return fmt.Errorf("update sync event error: %w", err)
	}

	return nil
}
