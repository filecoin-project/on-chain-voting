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

package models

import "time"

type SnapshotBackupTbl struct {
	Id        int64     `json:"id"`
	Day       string    `json:"day" gorm:"not null;unique:uniq_day"`    // Day is the backup day of the snapshot
	Height    int64     `json:"height" gorm:"not null"`                 // Height is the block height of the snapshot info
	Cid       string    `json:"cid" gorm:"not null;"`                   // Cid is the cid that the proposal is stored in ipfs
	ChainId   int64     `json:"chain_id" gorm:"not null"`               // Chain id
	Status    int       `json:"status" gorm:"not null"`                 // Status is the status of the sync snapshot, 0 means not sync, 3 means failed, 4 means synced
	RawData   string    `json:"raw_data" gorm:"type:longtext;not null"` // RawData is the raw data of the snapshot
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
