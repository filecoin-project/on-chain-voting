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
	"go.uber.org/zap"

	"powervoting-server/constant"
	"powervoting-server/model"
)

func (ev *Event) HandleVote(ctx context.Context, event VoteEvent, blockHeader *types.Header) error {
	voteData := model.VoteTbl{
		ProposalId:    event.Id.Int64(),
		Address:       event.Voter.Hex(),
		VoteEncrypted: event.VoteInfo,
		ChainId:       ev.Client.ChainId,
		Timestamp:     int64(blockHeader.Time),
		BlockNumber:   blockHeader.Number.Int64(),
		BaseField: model.BaseField{
			CreatedAt: time.Unix(int64(blockHeader.Time), 0),
		},
	}

	voterAddressData := model.VoterInfoTbl{
		Address:     event.Voter.Hex(),
		ChainId:     ev.Client.ChainId,
		BlockNumber: blockHeader.Number.Int64(),
		Timestamp:   int64(blockHeader.Time),
	}

	zap.L().Info(
		"Sync vote event parsed result",
		zap.Int64("proposal id", event.Id.Int64()),
		zap.String("voter", event.Voter.Hex()),
	)

	if err := ev.SyncService.AddVote(ctx, &voteData); err != nil {
		return fmt.Errorf("parse %s event error: %w", constant.VoteEvt, err)
	}

	if err := ev.SyncService.AddVoterAddress(ctx, &voterAddressData); err != nil {
		return fmt.Errorf("parse %s event error: %w", constant.VoteEvt, err)
	}

	zap.L().Info("Sync vote event success", zap.Int64("proposal id", event.Id.Int64()))
	return nil
}
