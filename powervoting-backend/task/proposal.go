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
			zap.L().Info("sync proposal start:", zap.Int64("networkId", network.Id))
			if err := SyncProposal(ethClient); err != nil {
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
func SyncProposal(ethClient model.GoEthClient) error {
	var dict model.Dict
	if err := db.Engine.Model(model.Dict{}).Where("name", constant.ProposalStartKey).Find(&dict).Error; err != nil {
		zap.L().Error("Get proposal start index error: ", zap.Error(err))
		return err
	}
	start, err := strconv.Atoi(dict.Value)
	if err != nil {
		zap.L().Error("Translate string to int error: ", zap.Error(err))
		return err
	}
	end, err := utils.GetProposalLatestId(ethClient)
	if err != nil {
		zap.L().Error("get proposal latest id error: ", zap.Error(err))
		return err
	}
	for start <= end {
		contractProposal, err := utils.GetProposal(ethClient, int64(start))
		if err != nil {
			zap.L().Error("Get proposal error: ", zap.Error(err))
			start++
			break
		}
		if len(contractProposal.Cid) == 0 {
			start++
			continue
		}
		var count int64
		if err = db.Engine.Model(model.Proposal{}).Where("cid", contractProposal.Cid).Count(&count).Error; err != nil {
			zap.L().Error("get proposal count error: ", zap.Error(err))
			return err
		}
		if count > 0 {
			start++
			continue
		}
		proposal := model.Proposal{
			Cid:          contractProposal.Cid,
			ProposalId:   int64(start),
			ProposalType: contractProposal.ProposalType.Int64(),
			Creator:      contractProposal.Creator.String(),
			StartTime:    contractProposal.StartTime.Int64(),
			ExpTime:      contractProposal.ExpTime.Int64(),
			VoteCount:    contractProposal.VotesCount.Int64(),
			Network:      ethClient.Id,
		}
		if err = db.Engine.Model(model.Proposal{}).Create(&proposal).Error; err != nil {
			zap.L().Error("create proposal error: ", zap.Error(err))
			return err
		}
		if err = db.Engine.Model(model.Dict{}).Create(&model.Dict{
			Name:  fmt.Sprintf("%s-%d", constant.VoteStartKey, proposal.ProposalId),
			Value: "1",
		}).Error; err != nil {
			zap.L().Error("create vote dict error: ", zap.Error(err))
			return err
		}
		start++
	}
	if err = db.Engine.Model(model.Dict{}).Where("name", constant.ProposalStartKey).Update("value", start).Error; err != nil {
		zap.L().Error("update proposal start key error: ", zap.Error(err))
		return err
	}

	return nil
}
