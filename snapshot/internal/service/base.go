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

package service

import (
	"context"

	"github.com/ybbus/jsonrpc/v3"

	models "power-snapshot/internal/model"
	"power-snapshot/utils/types"
)

type BaseRepo interface {
	GetLotusClient(ctx context.Context, netId int64) (jsonrpc.RPCClient, error)
	// Mapping of dates (YYYYMMDD) to block heights. (e.g. {"20250301": 123456})
	GetDateHeightMap(ctx context.Context, netId int64) (map[string]int64, error)
	SetDateHeightMap(ctx context.Context, netId int64, height map[string]int64) error
	SaveToLocalFile(ctx context.Context, netId int64, dayStr string, t string, data any) error
	GetDeveloperWeights(ctx context.Context, netId int64, dayStr string) ([]models.Nodes, error)
	GetDealsFromLocal(ctx context.Context, netId int64, dayStr string) (types.StateMarketDeals, error)
}

type MysqlRepo interface {
	CreateSnapshotBackup(ctx context.Context, in models.SnapshotBackupTbl) error
	GetSnapshotBackupList(ctx context.Context, chainId int64) ([]models.SnapshotBackupTbl, error)
	UpdateSnapshotBackup(ctx context.Context, in models.SnapshotBackupTbl) error
}

type ContractRepo interface {
	SetSnapshotHeight(dateStr string, height uint64) error
	GetSnapshotHeight(dateStr string) (int64, error)
	GetExpirationData() (int, error)
}
