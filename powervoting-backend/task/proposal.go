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

package task

import (
	"fmt"
	"powervoting-server/config"
	"powervoting-server/constant"
	"powervoting-server/contract"
	"powervoting-server/db"
	"powervoting-server/model"
	"powervoting-server/utils"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

// SyncProposalHandler asynchronously synchronizes proposals across multiple networks.
func SyncProposalHandler() {
	wg := sync.WaitGroup{}
	errList := make([]error, 0, len(config.Client.Network))
	mu := &sync.Mutex{}

	for _, network := range config.Client.Network {
		network := network
		ethClient, err := contract.GetClient(network.Id)
		if err != nil {
			zap.L().Error("get go-eth client error:", zap.Error(err))
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := SyncProposal(ethClient, db.Engine); err != nil {
				mu.Lock()
				errList = append(errList, err)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	if len(errList) != 0 {
		zap.L().Error("sync finished with err:", zap.Errors("errors", errList))
	}
	zap.L().Info("sync proposal finished: ", zap.Int64("timestamp", time.Now().Unix()))
}

// SyncProposal syncs proposals for a given Ethereum client.
func SyncProposal(ethClient model.GoEthClient, db db.DataRepo) error {
	dict, err := db.GetDict(constant.ProposalStartKey)
	if err != nil {
		return err
	}
	start, err := strconv.Atoi(dict.Value)
	if err != nil {
		zap.L().Error("Translate string to int error: ", zap.Error(err))
		return err
	}

	zap.L().Info("Sync proposal start: ", zap.Int("start", start))

	end, err := utils.GetProposalLatestId(ethClient)
	if err != nil {
		zap.L().Error("get proposal latest id error: ", zap.Error(err))
		return err
	}

	zap.L().Info("Sync proposal end: ", zap.Int("end", end))
	for start <= end {
		contractProposal, err := utils.GetProposal(ethClient, int64(start))
		if err != nil {
			zap.L().Error("Get proposal error: ", zap.Error(err))
			start++
			break
		}

		zap.L().Info("contract proposal: ", zap.Any("contractProposal", contractProposal))
		if len(contractProposal.Cid) == 0 {
			start++
			continue
		}
		proposal := model.Proposal{
			ProposalId:   int64(start),
			Cid:          contractProposal.Cid,
			ProposalType: contractProposal.ProposalType.Int64(),
			Creator:      contractProposal.Creator.String(),
			StartTime:    contractProposal.StartTime.Int64(),
			ExpTime:      contractProposal.ExpTime.Int64(),
			VoteCount:    contractProposal.VotesCount.Int64(),
			Network:      ethClient.Id,
			Status:       constant.ProposalStatusPending,
		}

		zap.L().Info("update proposal info: ", zap.Any("proposal", proposal))

		_, err = db.UpdateProposal(&proposal)
		if err != nil {
			return err
		}
		inDict := &model.Dict{
			Name:  fmt.Sprintf("%s-%d", constant.VoteStartKey, proposal.ProposalId),
			Value: "1",
		}
		_, err = db.CreateDict(inDict)
		if err != nil {
			return err
		}
		start++
	}
	err = db.UpdateDict(constant.ProposalStartKey, strconv.Itoa(start))
	if err != nil {
		zap.L().Error("update proposal start key error: ", zap.Error(err))
		return err
	}

	return nil
}
