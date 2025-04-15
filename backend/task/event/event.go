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
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"powervoting-server/config"
	"powervoting-server/constant"
	"powervoting-server/model"
	"powervoting-server/service"
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

	// Get the latest processed block info
	syncInfo, err := ev.SyncService.GetSyncEventInfo(ctx, ev.Network.PowerVotingContract)
	if err != nil {
		zap.L().Error("Get sync event info error: %v", zap.Error(err))
		return err
	}

	// Get the latest block height from the chain
	header, err := ev.Client.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		zap.L().Error("Get block header error: %v", zap.Error(err))
		return err
	}

	endBlock := header.Number.Int64()
	if endBlock <= syncInfo.SyncedHeight {
		if endBlock == syncInfo.SyncedHeight {
			zap.L().Info("It has been synchronized to the latest block height", zap.Int64("latest block height", endBlock))
			return constant.ErrAlreadySyncHeight
		}

		zap.L().Info("current block height is less than or equal to the latest processed block height", zap.Int64("latest block height", endBlock))
		return fmt.Errorf("invalid block height: %d", endBlock)
	}
	// Limit the number of blocks per processing
	if endBlock-syncInfo.SyncedHeight > constant.SyncBlockLimit {
		endBlock = syncInfo.SyncedHeight + constant.SyncBlockLimit
	}

	logs, err := ev.FetchMatchingEventLogs(ctx, big.NewInt(syncInfo.SyncedHeight+1), big.NewInt(endBlock))
	zap.L().Info(
		"fetch matching event logs success",
		zap.Int("event logs length", len(logs)),
		zap.Int64("current sync block height", syncInfo.SyncedHeight),
	)

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
	if err = ev.SyncService.UpdateSyncEventInfo(ctx, ev.Network.PowerVotingContract, endBlock); err != nil {
		zap.L().Error("Update sync event info error: %v", zap.Error(err))
		return err
	}

	return nil
}

// Get the event of the power voting contract according to the specified height interval of the chain.
// The starting height of synchronization is obtained from the configuration file.
// If synchronization fails, the next synchronization will continue to start from the height of the last synchronization.
func (ev *Event) FetchMatchingEventLogs(ctx context.Context, fromBlock, ToBlock *big.Int) ([]types.Log, error) {

	// Query the block scope log
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   ToBlock,
		Addresses: []common.Address{
			common.HexToAddress(ev.Network.PowerVotingContract),
			common.HexToAddress(ev.Network.FipContract),
			common.HexToAddress(ev.Network.OracleContract),
		},
		Topics: [][]common.Hash{
			{
				ev.Client.ABI.PowerVotingAbi.Events[constant.ProposalEvt].ID,
				ev.Client.ABI.PowerVotingAbi.Events[constant.VoteEvt].ID,
				ev.Client.ABI.FipAbi.Events[constant.FipCreateEvt].ID,
				ev.Client.ABI.FipAbi.Events[constant.FipVoteEvt].ID,
				ev.Client.ABI.FipAbi.Events[constant.FipPassedEvt].ID,
				ev.Client.ABI.OracleAbi.Events[constant.OracleUpdateGistIdsEvt].ID,
				ev.Client.ABI.OracleAbi.Events[constant.OracleUpdateMinerIdsEvt].ID,
			},
		},
	}
	// Get specified Logs using rpc of getth call chain
	logs, err := ev.Client.Client.FilterLogs(ctx, query)
	if err != nil {
		zap.L().Error("Get event logs error", zap.Error(err))
		return nil, err
	}

	return logs, nil
}

// ProcessingEventLogs processes a list of event logs using a provided Ethereum client.
func (ev *Event) ProcessingEventLogs(ctx context.Context, logs []types.Log) {
	for _, vLog := range logs {
		// Attempt to parse the event log using the client's PowerVotingAbi and the log data.
		err := ev.parseEvent(ctx, vLog)
		if err != nil {
			zap.L().Error("Parse event error", zap.String("rpc url", config.Client.Network.Rpc), zap.Error(err))
			continue
		}

		zap.L().Info("Event parsed result:", zap.Any("event", vLog.Topics[0].Hex()))
	}
}

// parses the event log for the forum and stores the parsed results to the database.
func (ev *Event) parseEvent(ctx context.Context, vLog types.Log) error {
	// Get the block header of the block containing the event log
	blockHeader, err := ev.Client.Client.HeaderByNumber(ctx, big.NewInt(int64(vLog.BlockNumber)))
	if err != nil {
		// return fmt.Errorf("get block header error: %v", err)
		zap.L().Warn("get block header error ", zap.Error(err))
		blockHeader = &types.Header{
			Time:   uint64(time.Now().Unix()),
			Number: big.NewInt(int64(vLog.BlockNumber)),
		}
	}

	switch vLog.Topics[0].Hex() {
	// Parse the proposal event
	case ev.Client.ABI.PowerVotingAbi.Events[constant.ProposalEvt].ID.Hex():
		var event ProposalCreateEvent
		if err := ev.parsePowerVotingEvent(&event, constant.ProposalEvt, vLog.Data); err != nil {
			return fmt.Errorf("unpack proposal event error: %v", err)
		}

		if err := ev.HandleProposalCreate(ctx, event, blockHeader); err != nil {
			return fmt.Errorf("handle proposal create error: %v", err)
		}

	case ev.Client.ABI.PowerVotingAbi.Events[constant.VoteEvt].ID.Hex():
		var event VoteEvent
		if err := ev.parsePowerVotingEvent(&event, constant.VoteEvt, vLog.Data); err != nil {
			return err
		}

		if err := ev.HandleVote(ctx, event, blockHeader); err != nil {
			return fmt.Errorf("handle vote error: %v", err)
		}

	case ev.Client.ABI.FipAbi.Events[constant.FipCreateEvt].ID.Hex():
		var event FipEditorProposalCreateEvent
		if err := ev.parseFipEvent(&event, constant.FipCreateEvt, vLog.Data); err != nil {
			return fmt.Errorf("unpack fip create event error: %v", err)
		}

		if err := ev.HandleFipEditorProposalCreateEvent(ctx, event, blockHeader); err != nil {
			return fmt.Errorf("handle fip create error: %v", err)
		}

	case ev.Client.ABI.FipAbi.Events[constant.FipPassedEvt].ID.Hex():
		var event FipEditorProposalPassedEvent
		if err := ev.parseFipEvent(&event, constant.FipPassedEvt, vLog.Data); err != nil {
			return fmt.Errorf("unpack fip passed event error: %v", err)
		}

		if err := ev.HandleFipEditorProposalPassedEvent(ctx, event, blockHeader); err != nil {
			return fmt.Errorf("handle fip passed error: %v", err)
		}

	case ev.Client.ABI.FipAbi.Events[constant.FipVoteEvt].ID.Hex():
		var event FipEditorProposalVoteEvent
		if err := ev.parseFipEvent(&event, constant.FipVoteEvt, vLog.Data); err != nil {
			return fmt.Errorf("unpack fip vote event error: %v", err)
		}

		if err := ev.HandleFipEditorProposalVoteEvent(ctx, event, blockHeader); err != nil {
			return fmt.Errorf("handle fip vote error: %v", err)
		}

	case ev.Client.ABI.OracleAbi.Events[constant.OracleUpdateGistIdsEvt].ID.Hex():
		var event OracleUpdateGistIdsEvent
		if err := ev.parseOracleEvent(&event, constant.OracleUpdateGistIdsEvt, vLog.Data); err != nil {
			return fmt.Errorf("unpack oracle update gist ids event error: %v", err)
		}

		if err := ev.HandleOracleUpdateGistId(ctx, event, ev.Client.ChainId,blockHeader); err != nil {
			return fmt.Errorf("handle oracle update gist ids error: %v", err)
		}

	case ev.Client.ABI.OracleAbi.Events[constant.OracleUpdateMinerIdsEvt].ID.Hex():
		var event OracleUpdateMinerIdsEvent
		if err := ev.parseOracleEvent(&event, constant.OracleUpdateMinerIdsEvt, vLog.Data); err != nil {
			return fmt.Errorf("unpack oracle update miner ids event error: %v", err)
		}

		if err := ev.HandleOracleUpdateMinerIds(ctx, event, blockHeader); err != nil {
			return fmt.Errorf("handle oracle update miner ids error: %v", err)
		}

	default:
		return fmt.Errorf("unknown event")
	}

	return nil
}

// Parse the event log using the client's PowerVotingAbi and the log data.
func (ev *Event) parsePowerVotingEvent(event any, unpackName string, log []byte) error {
	// Decode non-index parameters
	return ev.Client.ABI.PowerVotingAbi.UnpackIntoInterface(event, unpackName, log)
}

// Parse the event log using the client's Fip Abi and the log data.
func (ev *Event) parseFipEvent(event any, unpackName string, log []byte) error {
	// Decode non-index parameters
	return ev.Client.ABI.FipAbi.UnpackIntoInterface(event, unpackName, log)
}

func (ev *Event) parseOracleEvent(event any, unpackName string, log []byte) error {
	// Decode non-index parameters
	return ev.Client.ABI.OracleAbi.UnpackIntoInterface(event, unpackName, log)
}

func FetchEventFromRPC(rpcUrl string) ([]types.Log, error) {
	client := resty.New().
		SetTimeout(constant.RequestTimeout).
		SetRetryCount(3).
		SetRetryWaitTime(time.Second)
	var events FilFoxRPCEvents
	resp, err := client.R().SetResult(&events).Get(rpcUrl)
	if err != nil {
		zap.L().Error("fetch event from rpc error", zap.String("rpcUrl", rpcUrl), zap.Error(err))
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		zap.L().Error(
			"fetch event from rpc error",
			zap.String("rpcUrl", rpcUrl),
			zap.Int("code", resp.StatusCode()),
			zap.String("body", resp.String()),
		)
		return nil, fmt.Errorf("fetch event from rpc error: %s", resp.String())
	}

	var res []types.Log
	for _, log := range events.EventLogs {
		res = append(res, types.Log{
			Address:     common.HexToAddress(log.Address),
			Data:        common.FromHex(log.Data),
			Topics:      utils.ConvertTopics(log.Topics),
			BlockNumber: uint64(log.BlockNumber),
			Removed:     log.Removed,
		})
	}

	return res, nil
}
