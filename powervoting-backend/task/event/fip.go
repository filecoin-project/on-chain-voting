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

	"powervoting-server/constant"
	"powervoting-server/model"
)

func (ev *Event) HandleFipEditorProposalCreateEvent(ctx context.Context, event FipEditorProposalCreateEvent, blockHeader *types.Header) error {
	if err := ev.SyncService.CreateFipProposal(ctx, &model.FipProposalTbl{
		ProposalId:       event.Proposal.ProposalId.Int64(),
		ChainId:          ev.Client.ChainId,
		ProposalType:     int(event.Proposal.ProposalType),
		Creator:          event.Proposal.Creator.Hex(),
		CandidateAddress: event.Proposal.CandidateAddress.Hex(),
		CandidateInfo:    event.Proposal.CandidateInfo,
		BlockNumber:      blockHeader.Number.Int64(),
		Timestamp:        int64(blockHeader.Time),
	}); err != nil {
		zap.L().Error(
			"failed to create fip proposal",
			zap.Int64("proposal id", event.Proposal.ProposalId.Int64()),
			zap.String("creator", event.Proposal.Creator.Hex()),
			zap.String("candidate address", event.Proposal.CandidateAddress.Hex()),
			zap.String("candidate info", event.Proposal.CandidateInfo),
		)
		return err
	}

	zap.L().Info(
		"fip proposal created",
		zap.Any("event", event),
	)

	return nil
}

func (ev *Event) HandleFipEditorProposalPassedEvent(ctx context.Context, event FipEditorProposalPassedEvent, blockHeader *types.Header) error {
	zap.L().Info("fip proposal passed", zap.Int64("proposal id", event.ProposalId.Int64()))

	fipProposal, err := ev.SyncService.UpdateStatusAndGetFipProposal(ctx, event.ProposalId.Int64(), ev.Client.ChainId)
	if err != nil {
		zap.L().Error(
			"failed to update fip proposal status",
			zap.Int64("proposal id", event.ProposalId.Int64()),
		)
		return err
	}

	if fipProposal.ProposalType == constant.FipProposalRevoke {
		zap.L().Info("fip editor revoke", zap.Int64("proposal id", event.ProposalId.Int64()))
		if err := ev.SyncService.RevokeFipProposal(ctx, fipProposal.ChainId, fipProposal.CandidateAddress); err != nil {
			zap.L().Error(
				"failed to handle revoke proposal",
				zap.String("revoke address", fipProposal.CandidateAddress),
				zap.Int64("proposal id", event.ProposalId.Int64()),
			)

			return err
		}

		zap.L().Info("fip editor revoked successfully", zap.Int64("proposal id", event.ProposalId.Int64()), zap.String("revoke address", fipProposal.CandidateAddress))
	} else {
		zap.L().Info("fip editor approve", zap.Int64("proposal id", event.ProposalId.Int64()))
		if err := ev.SyncService.CreateFipEditor(ctx, &model.FipEditorTbl{
			ChainId:  ev.Client.ChainId,
			Editor:   fipProposal.CandidateAddress,
			IsRemove: constant.FipEditorValid,
		}); err != nil {
			zap.L().Error(
				"failed to create fip voter",
				zap.String("voter", fipProposal.CandidateAddress),
			)
			return err
		}

		zap.L().Info("fip editor approved successfully", zap.Int64("proposal id", event.ProposalId.Int64()), zap.String("approve address", fipProposal.CandidateAddress))
	}

	return nil
}

func (ev *Event) HandleFipEditorProposalVoteEvent(ctx context.Context, event FipEditorProposalVoteEvent, blockHeader *types.Header) error {
	if err := ev.SyncService.CreateFipProposalVotedInfo(ctx, &model.FipProposalVoteTbl{
		ProposalId:  event.VoteInfo.ProposalId.Int64(),
		ChainId:     ev.Client.ChainId,
		Voter:       event.VoteInfo.Voter.Hex(),
		BlockNumber: blockHeader.Number.Int64(),
		Timestamp:   int64(blockHeader.Time),
	}); err != nil {
		zap.L().Error(
			"failed to create fip proposal voted info",
			zap.Int64("proposal id", event.VoteInfo.ProposalId.Int64()),
			zap.String("voter", event.VoteInfo.Voter.Hex()),
		)
		return err
	}

	zap.L().Info("fip proposal voted", zap.Any("event", event))
	return nil
}
