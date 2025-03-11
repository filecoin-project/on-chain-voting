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
package mock

import (
	"context"
	"powervoting-server/model"
	"powervoting-server/service"
)

type MockSyncService struct {
}

// CreateSyncEventInfo implements service.SyncRepo.
func (m *MockSyncService) CreateSyncEventInfo(ctx context.Context, in *model.SyncEventTbl) (int64, error) {
	panic("unimplemented")
}

// GetSyncEventInfo implements service.SyncRepo.
func (m *MockSyncService) GetSyncEventInfo(ctx context.Context, addr string) (*model.SyncEventTbl, error) {
	panic("unimplemented")
}

// UpdateSyncEventInfo implements service.SyncRepo.
func (m *MockSyncService) UpdateSyncEventInfo(ctx context.Context, addr string, height int64) error {
	panic("unimplemented")
}

var _ service.SyncRepo = (*MockSyncService)(nil)
