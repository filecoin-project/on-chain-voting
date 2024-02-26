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
	"go.uber.org/zap"
	"math"
	"math/big"
	"powervoting-server/config"
	"powervoting-server/contract"
	"powervoting-server/db"
	"powervoting-server/model"
	"powervoting-server/utils"
	"time"
)

// VotingCountHandler vote count
func VotingCountHandler() {
	networkList := config.Client.Network
	for _, network := range networkList {
		ethClient, err := contract.GetClient(network.Id)
		if err != nil {
			zap.L().Error("get go-eth client error:", zap.Error(err))
			continue
		}
		go VotingCount(ethClient)
	}
}

func VotingCount(ethClient model.GoEthClient) {
	now, err := utils.GetTimestamp(ethClient)
	if err != nil {
		zap.L().Error("get timestamp on chain error: ", zap.Error(err))
		now = time.Now().Unix()
	}
	proposals, err := db.GetProposalList(ethClient.Id, now)
	zap.L().Info("proposal list: %+v\n", zap.Reflect("proposals", proposals))
	if err != nil {
		zap.L().Error("get proposal from db error:", zap.Error(err))
		return
	}
	for _, proposal := range proposals {
		SyncVote(ethClient, proposal.ProposalId)
		voteInfos, err := db.GetVoteList(ethClient.Id, proposal.ProposalId)
		if err != nil {
			zap.L().Error("get vote info from db error:", zap.Error(err))
			continue
		}
		var voteList []model.Vote4Counting
		zap.L().Info("voteInfos: %+v\n", zap.Reflect("voteInfos", voteInfos))
		for _, voteInfo := range voteInfos {
			list, err := utils.DecodeVoteList(voteInfo)
			if err != nil {
				zap.L().Error("get vote info from IPFS or decrypt error: ", zap.Error(err))
				return
			}
			voteList = append(voteList, list...)
		}
		zap.L().Info("voteList: %+v\n", zap.Reflect("voteList", voteList))
		var voteHistoryList []model.VoteHistory
		// calc total power
		totalSpPower := new(big.Int)
		totalClientPower := new(big.Int)
		totalTokenPower := new(big.Int)
		totalDeveloperPower := new(big.Int)
		powerMap := make(map[string]model.Power)
		addressIsCount := make(map[string]bool)
		for _, vote := range voteList {
			power, err := utils.GetPower(vote.Address, ethClient)
			if err != nil {
				zap.L().Error("get power error: ", zap.Error(err))
				return
			}
			zap.L().Info("address: %s, power: %+v\n", zap.Reflect("address", vote.Address), zap.Reflect("power", power))
			if !addressIsCount[vote.Address] {
				powerMap[vote.Address] = power
				totalSpPower.Add(totalSpPower, power.SpPower)
				totalClientPower.Add(totalClientPower, power.ClientPower)
				totalTokenPower.Add(totalTokenPower, power.TokenHolderPower)
				totalDeveloperPower.Add(totalDeveloperPower, power.DeveloperPower)
				addressIsCount[vote.Address] = true
			}
		}
		// vote counting
		var result = make(map[int64]float64, 5) // max 5 options
		for _, vote := range voteList {
			power := powerMap[vote.Address]
			var votes float64
			if vote.Votes != 0 {
				var spPercent float64
				var clientPercent float64
				var tokenPercent float64
				var developerPercent float64
				if totalSpPower.Int64() != 0 {
					spPercent = float64(power.SpPower.Int64()) / float64(totalSpPower.Int64())
				}
				if totalClientPower.Int64() != 0 {
					clientPercent = float64(power.ClientPower.Int64()) / float64(totalClientPower.Int64())
				}
				if totalTokenPower.Int64() != 0 {
					tokenPercent = float64(power.TokenHolderPower.Int64()) / float64(totalTokenPower.Int64())
				}
				if totalDeveloperPower.Int64() != 0 {
					developerPercent = float64(power.DeveloperPower.Int64()) / float64(totalDeveloperPower.Int64())
				}
				var votePercent = float64(vote.Votes) / 100
				votes = ((spPercent * 25) + (clientPercent * 25) + (tokenPercent * 25) + (developerPercent * 25)) * votePercent
			}
			voteHistory := model.VoteHistory{
				ProposalId: proposal.ProposalId,
				OptionId:   vote.OptionId,
				Votes:      math.Round(votes*100) / 100,
				Address:    vote.Address,
				Network:    ethClient.Id,
			}
			voteHistoryList = append(voteHistoryList, voteHistory)
			if _, ok := result[vote.OptionId]; ok {
				result[vote.OptionId] += votes
			} else {
				result[vote.OptionId] = votes
			}
		}
		var voteResultList []model.VoteResult
		options, err := utils.GetOptions(proposal.Cid)
		if err != nil {
			zap.L().Error("get options error: ", zap.Error(err))
			continue
		}
		for i := 0; i < len(options); i++ {
			voteResult := model.VoteResult{
				ProposalId: proposal.ProposalId,
				OptionId:   int64(i),
				Votes:      math.Round(result[int64(i)]*100) / 100,
				Network:    ethClient.Id,
			}
			voteResultList = append(voteResultList, voteResult)
		}
		// Save vote history and vote result to database and update status
		db.VoteResult(proposal.Id, voteHistoryList, voteResultList)
	}
}
