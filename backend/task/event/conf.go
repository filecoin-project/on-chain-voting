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

package event

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

func (ev *Event) HandleConfRepoAddedEvent(ctx context.Context, event ConfRepoAddedEvent, blockHeader *types.Header) error {
	zap.L().Info("handle conf repo added event", zap.Any("added repo", event))
	if err := ev.SyncService.AddGithubRepoName(ctx, int(event.RepoInfo.OrgType), event.RepoInfo.RepoName, event.Id.Int64(), int64(blockHeader.Time)); err != nil {
		zap.L().Error("failed to add github repo name", zap.Error(err))
		return err
	}

	return nil
}

func (ev *Event) HandleConfRepoRemovedEvent(ctx context.Context, event ConfRepoRemovedEvent, blockHeader *types.Header) error {
	zap.L().Info("handle conf repo removed event", zap.Any("removed repo", event))
	if err := ev.SyncService.DeleteGithubRepoName(ctx, event.Id.Int64()); err != nil {
		zap.L().Error("failed to delete github repo name", zap.Error(err))
		return err
	}

	return nil
}
