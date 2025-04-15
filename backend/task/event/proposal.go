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
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang-module/carbon"
	"go.uber.org/zap"

	snapshot "powervoting-server/api/rpc"
	"powervoting-server/constant"
	"powervoting-server/model"
)

func (ev *Event) HandleProposalCreate(ctx context.Context, event ProposalCreateEvent, blockHeader *types.Header) error {
	data := model.ProposalTbl{
		ProposalId:        event.Id.Int64(),
		Creator:           event.Proposal.Creator.Hex(),
		StartTime:         event.Proposal.StartTime.Int64(),
		EndTime:           event.Proposal.EndTime.Int64(),
		Content:           event.Proposal.Content,
		Title:             event.Proposal.Title,
		SnapshotTimestamp: event.Proposal.SnapshotTimestamp.Int64(),
		Percentage: model.Percentage{
			TokenHolderPercentage: event.Proposal.TokenHolderPercentage,
			SpPercentage:          event.Proposal.SpPercentage,
			ClientPercentage:      event.Proposal.ClientPercentage,
			DeveloperPercentage:   event.Proposal.DeveloperPercentage,
		},
		Timestamp:   int64(blockHeader.Time),
		BlockNumber: blockHeader.Number.Int64(),
		ChainId:     ev.Client.ChainId,
		BaseField: model.BaseField{
			CreatedAt: time.Unix(int64(blockHeader.Time), 0),
		},
	}

	zap.L().Info(
		"Sync proposal event parsed result",
		zap.Int64("proposal id", event.Id.Int64()),
		zap.String("creator", event.Proposal.Creator.Hex()),
		zap.Int64("start time", event.Proposal.StartTime.Int64()),
		zap.Int64("end time", event.Proposal.EndTime.Int64()),
		zap.String("title", event.Proposal.Title),
		zap.String("content", event.Proposal.Content),
		zap.Int64("snapshot timestamp", event.Proposal.SnapshotTimestamp.Int64()),
	)

	snapshotday := carbon.CreateFromTimestamp(event.Proposal.SnapshotTimestamp.Int64()).ToShortDateString()
	snapshotInfo, err := snapshot.UploadSnapshotInfo(ev.Client.ChainId, snapshotday)
	if err != nil {
		zap.L().Error("record snapshot error", zap.Error(err))
	}
	zap.L().Info(
		"record snapshot info",
		zap.Int64("proposal id", event.Id.Int64()),
		zap.Int64("chain id", ev.Client.ChainId),
		zap.Int64("snapshot block height", snapshotInfo.Height),
		zap.String("snapshot day", snapshotInfo.Day),
	)

	data.SnapshotBlockHeight = snapshotInfo.Height
	data.SnapshotDay = snapshotday

	if err = ev.SyncService.AddProposal(ctx, &data); err != nil {
		return fmt.Errorf("parse %s event error: %w", constant.ProposalEvt, err)
	}
	zap.L().Info("Sync proposal event success", zap.Int64("proposal id", event.Id.Int64()))

	return nil
}
