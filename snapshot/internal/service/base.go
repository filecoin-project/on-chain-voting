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
)

type BaseRepo interface {
	GetEthClient(ctx context.Context, netID int64) (*models.GoEthClient, error)

	GetLotusClient(ctx context.Context, netId int64) (jsonrpc.RPCClient, error)

	GetDateHeightMap(ctx context.Context, netId int64) (map[string]int64, error)
	SetDateHeightMap(ctx context.Context, netId int64, height map[string]int64) error
}

type MysqlRepo interface {
	CreateSnapshotBackup(ctx context.Context, in models.SnapshotBackupTbl) error
	GetSnapshotBackupList(ctx context.Context, chainId int64) ([]models.SnapshotBackupTbl, error)
	UpdateSnapshotBackup(ctx context.Context, in models.SnapshotBackupTbl) error
}
