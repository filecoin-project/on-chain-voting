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
	"go.uber.org/zap"
	"powervoting-server/config"
	"powervoting-server/constant"
	"powervoting-server/contract"
	"powervoting-server/db"
	"powervoting-server/model"
	"powervoting-server/utils"
	"strconv"
	"sync"
	"time"
)

// SyncVoteHandler asynchronously synchronizes votes for proposals across multiple networks.
func SyncVoteHandler() {
	wg := sync.WaitGroup{}
	errList := make([]error, 0)
	mu := &sync.Mutex{}

	for _, network := range config.Client.Network {
		var proposalList []model.Proposal
		if err := db.Engine.Model(model.Proposal{}).Where("status", 0).Where("network", network.Id).Find(&proposalList).Error; err != nil {
			zap.L().Error("get proposal list error: ", zap.Error(err))
		}
		ethClient, err := contract.GetClient(network.Id)
		if err != nil {
			zap.L().Error("get go-eth client error:", zap.Error(err))
			continue
		}
		for _, proposal := range proposalList {
			proposal := proposal
			wg.Add(1)
			go func() {
				defer wg.Done()
				zap.L().Info("sync vote start",
					zap.Int64("networkId", network.Id),
					zap.Int64("proposalId", proposal.ProposalId),
					zap.Int64("timestamp", time.Now().Unix()))
				err := SyncVote(ethClient, proposal.ProposalId, db.Engine)
				if err != nil {
					mu.Lock()
					errList = append(errList, err)
					mu.Unlock()
				}
			}()
		}
	}
	wg.Wait()

	if len(errList) != 0 {
		zap.L().Error("sync vote with err:", zap.Errors("errors", errList))
	}
	zap.L().Info("sync vote finished: ", zap.Int64("timestamp", time.Now().Unix()))
}

// SyncVote syncs votes for a given proposal and Ethereum client.
func SyncVote(ethClient model.GoEthClient, proposalId int64, db db.DataRepo) error {
	dictName := fmt.Sprintf("%s-%d", constant.VoteStartKey, proposalId)
	dict, err := db.GetDict(dictName)
	if err != nil {
		return err
	}
	start, err := strconv.Atoi(dict.Value)
	if err != nil {
		zap.L().Error("Translate string to int error: ", zap.Error(err))
		return err
	}
	contractProposal, err := utils.GetProposal(ethClient, proposalId)
	if err != nil {
		zap.L().Error("get proposal error: ", zap.Error(err))
		return err
	}
	end := int(contractProposal.VotesCount.Int64())
	for start <= end {
		contractVote, err := utils.GetVote(ethClient, proposalId, int64(start))
		if err != nil {
			zap.L().Error("Get vote error: ", zap.Error(err))
			start++
			break
		}
		if len(contractVote.VoteInfo) == 0 {
			start++
			continue
		}
		count, err := db.CountVotes(map[string]any{
			"network":     ethClient.Id,
			"proposal_id": proposalId,
			"address":     contractVote.Voter.String()})
		if err != nil {
			return err
		}
		if count > 0 {
			_ = db.UpdateVoteInfo(map[string]any{
				"network":     ethClient.Id,
				"proposal_id": proposalId,
				"address":     contractVote.Voter.String()}, contractVote.VoteInfo)
			start++
			continue
		}
		vote := model.Vote{
			ProposalId: proposalId,
			Address:    contractVote.Voter.String(),
			VoteInfo:   contractVote.VoteInfo,
			Network:    ethClient.Id,
		}
		_, err = db.CreateVote(&vote)
		if err != nil {
			return err
		}
		start++
	}

	err = db.UpdateDict(dictName, strconv.Itoa(start))
	if err != nil {
		return err
	}

	return nil
}
