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

	"powervoting-server/utils"
)

func (ev *Event) HandleOracleUpdateGistId(ctx context.Context, event OracleUpdateGistIdsEvent, blockHeader *types.Header) error {
	zap.L().Info("oracle update gist ids event handled", zap.String("voter", event.VoterAddress.Hex()), zap.String("gist id", event.GistIds))
	if err := ev.SyncService.UpdateVoterAndProposalGithubNameByGistInfo(ctx, event.VoterAddress.Hex(), event.GistIds); err != nil {
		zap.L().Error("failed to update voter by gist info", zap.Error(err))
		return err
	}

	zap.L().Info("oracle update gist ids event handled successfully")
	return nil
}

func (ev *Event) HandleOracleUpdateMinerIds(ctx context.Context, event OracleUpdateMinerIdsEvent, blockHeader *types.Header) error {
	zap.L().Info("oracle update miner ids event handled", zap.String("voter", event.VoterAddress.Hex()), zap.Any("miner ids", event.MinerIds))

	ownerId := utils.GetOwnerIdByOracle(ev.Client, event.MinerIds)
	actionId, err := utils.GetActorIdByAddress(ev.Client, event.VoterAddress.Hex())
	if err != nil {
		zap.L().Error("failed to get actor id by address", zap.String("voter", event.VoterAddress.Hex()), zap.Error(err))
		return err
	}

	if ownerId != actionId {
		zap.L().Error("oracle update miner ids event failed: owner id != action id", zap.String("voter", event.VoterAddress.Hex()), zap.Any("owner id", ownerId), zap.Any("action id", actionId))
		return nil
	}

	if err := ev.SyncService.UpdateVoterByMinerIds(ctx, event.VoterAddress.Hex(), event.MinerIds, ownerId); err != nil {
		zap.L().Error("failed to update voter by miner ids", zap.Error(err))
		return err
	}

	zap.L().Info("oracle update miner ids event handled successfully")
	return nil
}
