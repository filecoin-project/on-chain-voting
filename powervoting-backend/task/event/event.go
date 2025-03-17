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
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang-module/carbon"
	"go.uber.org/zap"

	"powervoting-server/config"
	"powervoting-server/constant"
	"powervoting-server/model"
	"powervoting-server/service"
	"powervoting-server/snapshot"
	"powervoting-server/utils"
)

type Event struct {
	Client      *model.GoEthClient
	SyncService service.ISyncService
	Network     *config.Network
}

func (ev *Event) SubscribeEvent() error {
	ctx := context.Background()
	// Collect log information
	logs, syncHeight, err := ev.FetchMatchingEventLogs(ctx)
	zap.L().Info("fetch matching event logs success", zap.Int64("current sync block height", syncHeight))
	if err != nil {
		if errors.Is(err, constant.ErrAlreadySyncHeight) {
			return nil
		}

		zap.L().Error("fetch matching event logs error", zap.Error(err))
		return err
	}

	// Process the log information
	ev.ProcessingEventLogs(ctx, logs)

	// Update the latest processed block height to the database
	if err = ev.SyncService.UpdateSyncEventInfo(ctx, ev.Network.PowerVotingContract, syncHeight); err != nil {
		zap.L().Error("Update sync event info error: %v", zap.Error(err))
		return err
	}

	return nil
}

// Get the event of the power voting contract according to the specified height interval of the chain.
// The starting height of synchronization is obtained from the configuration file.
// If synchronization fails, the next synchronization will continue to start from the height of the last synchronization.
func (ev *Event) FetchMatchingEventLogs(ctx context.Context) ([]types.Log, int64, error) {
	// Get the latest processed block info
	syncInfo, err := ev.SyncService.GetSyncEventInfo(ctx, ev.Network.PowerVotingContract)
	if err != nil {
		zap.L().Error("Get sync event info error: %v", zap.Error(err))
		return nil, 0, err
	}

	// Get the latest block height from the chain
	header, err := ev.Client.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		zap.L().Error("Get block header error: %v", zap.Error(err))
		return nil, syncInfo.SyncedHeight, err
	}

	endBlock := header.Number.Int64()
	if endBlock <= syncInfo.SyncedHeight {
		if endBlock == syncInfo.SyncedHeight {
			zap.L().Info("It has been synchronized to the latest block height", zap.Int64("latest block height", endBlock))
			return nil, syncInfo.SyncedHeight, constant.ErrAlreadySyncHeight
		}
		zap.L().Info("current block height is less than or equal to the latest processed block height", zap.Int64("latest block height", endBlock))
		return nil, syncInfo.SyncedHeight, fmt.Errorf("invalid block height: %d", endBlock)
	}
	// Limit the number of blocks per processing
	if endBlock-syncInfo.SyncedHeight > constant.SyncBlockLimit {
		endBlock = syncInfo.SyncedHeight + constant.SyncBlockLimit
	}

	// Query the block scope log
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(syncInfo.SyncedHeight + 1),
		ToBlock:   big.NewInt(endBlock),
		Addresses: []common.Address{common.HexToAddress(syncInfo.ContractAddress)},
		Topics: [][]common.Hash{
			{
				ev.Client.ABI.PowerVotingAbi.Events[constant.ProposalEvt].ID,
				ev.Client.ABI.PowerVotingAbi.Events[constant.VoteEvt].ID,
			},
		},
	}
	// Get specified Logs using rpc of getth call chain
	logs, err := ev.Client.Client.FilterLogs(ctx, query)
	if err != nil {
		zap.L().Error("Get event logs error", zap.Error(err), zap.Any("start height", syncInfo.SyncedHeight+1), zap.Any("end height", endBlock))
		return nil, syncInfo.SyncedHeight, err
	}

	return logs, endBlock, nil
}

// ProcessingEventLogs processes a list of event logs using a provided Ethereum client.
func (ev *Event) ProcessingEventLogs(ctx context.Context, logs []types.Log) {
	for _, vLog := range logs {
		// Attempt to parse the event log using the client's PowerVotingAbi and the log data.
		err := ev.parseEvent(ctx, vLog)
		if err != nil {
			zap.L().Error("Parse event error: %v", zap.Error(err))
			continue
		}

		zap.L().Info("Event parsed result:", zap.Any("event", vLog.Topics[0].Hex()))
	}
}

// parses the event log for the forum and stores the parsed results to the database.
func (ev *Event) parseEvent(ctx context.Context, vLog types.Log) error {
	// Parse the event log using the client's PowerVotingAbi and the log data.
	parseFunc := func(event any, unpackName string) error {
		// Decode non-index parameters
		return ev.Client.ABI.PowerVotingAbi.UnpackIntoInterface(event, unpackName, vLog.Data)
	}

	// Get the block header of the block containing the event log
	blockHeader, err := ev.Client.Client.HeaderByNumber(ctx, big.NewInt(int64(vLog.BlockNumber)))
	if err != nil {
		return fmt.Errorf("get block header error: %v", err)
	}

	switch vLog.Topics[0].Hex() {
	// Parse the proposal event
	case ev.Client.ABI.PowerVotingAbi.Events[constant.ProposalEvt].ID.Hex():
		var event ProposalCreateEvent
		if err := parseFunc(&event, constant.ProposalEvt); err != nil {
			return fmt.Errorf("unpack proposal event error: %v", err)
		}

		// Get the proposal creator's information from the Oracle contract
		// Includes the voter's GitHub account and other information
		voterInfo, err := utils.GetProposalCreatorByOracle(event.Proposal.Creator.Hex(), ev.Client)
		if err != nil {
			zap.L().Error("Get proposal creator by oracle error", zap.Error(err))
			voterInfo = model.VoterInfo{}
		}
		zap.L().Info(
			"Get proposal creator infomation by oracle result",
			zap.Int64("proposal id", event.Id.Int64()),
			zap.String("creator", event.Proposal.Creator.Hex()),
			zap.String("github name", voterInfo.GithubAccount),
		)

		data := model.ProposalTbl{
			ProposalId:        event.Id.Int64(),
			Creator:           event.Proposal.Creator.Hex(),
			GithubName:        voterInfo.GithubAccount,
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

	case ev.Client.ABI.PowerVotingAbi.Events[constant.VoteEvt].ID.Hex():
		var event VoteEvent
		if err := parseFunc(&event, constant.VoteEvt); err != nil {
			return err
		}

		data := model.VoteTbl{
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

		zap.L().Info(
			"Sync vote event parsed result",
			zap.Int64("proposal id", event.Id.Int64()),
			zap.String("voter", event.Voter.Hex()),
		)

		if err = ev.SyncService.AddVote(ctx, &data); err != nil {
			return fmt.Errorf("parse %s event error: %w", constant.VoteEvt, err)
		}

		zap.L().Info("Sync vote event success", zap.Int64("proposal id", event.Id.Int64()))
	default:
		return fmt.Errorf("unknown event")
	}

	return nil
}
