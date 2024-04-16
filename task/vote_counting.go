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
	"math"
	"math/big"
	"powervoting-server/config"
	"powervoting-server/contract"
	"powervoting-server/db"
	"powervoting-server/model"
	"powervoting-server/utils"
	"time"

	"go.uber.org/zap"
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

// bigIntDiv bigInt division, which returns four decimal digits reserved
func bigIntDiv(x *big.Int, y *big.Int) (z float64) {
	var x_x = big.NewInt(0)
	if y.Uint64() == 0 {
		return 0
	}
	x_x.Mul(x, big.NewInt(10000))
	x_x.Div(x_x, y)
	z = float64(x_x.Uint64()) / 10000
	return
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
		var votePowerList []model.VotePower
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
				if totalSpPower.Uint64() != 0 {
					spPercent = bigIntDiv(power.SpPower, totalSpPower)
				}
				if totalClientPower.Uint64() != 0 {
					clientPercent = bigIntDiv(power.ClientPower, totalClientPower)
				}
				if totalTokenPower.Uint64() != 0 {
					tokenPercent = bigIntDiv(power.TokenHolderPower, totalTokenPower)
				}
				if totalDeveloperPower.Uint64() != 0 {
					developerPercent = bigIntDiv(power.DeveloperPower, totalDeveloperPower)
				}
				var votePercent = float64(vote.Votes) / 100
				votes = ((spPercent * 25) + (clientPercent * 25) + (tokenPercent * 25) + (developerPercent * 25)) * votePercent

				votePower := model.VotePower{
					HistoryId:               proposal.ProposalId,
					Address:                 vote.Address,
					OptionId:                vote.OptionId,
					Votes:                   math.Round(votes*100) / 100,
					SpPower:                 power.SpPower.String(),
					ClientPower:             power.ClientPower.String(),
					TokenHolderPower:        power.TokenHolderPower.String(),
					DeveloperPower:          power.DeveloperPower.String(),
					SpPowerPercent:          math.Round(spPercent*10000) / 100,
					ClientPowerPercent:      math.Round(clientPercent*10000) / 100,
					TokenHolderPowerPercent: math.Round(tokenPercent*10000) / 100,
					DeveloperPowerPercent:   math.Round(developerPercent*10000) / 100,
					PowerBlockHeight:        int64(power.BlockHeight.Uint64()),
				}
				votePowerList = append(votePowerList, votePower)
			}
			result[vote.OptionId] += votes
		}

		voteHistory := model.VoteCompleteHistory{
			ProposalId:            proposal.ProposalId,
			Network:               proposal.Network,
			TotalSpPower:          totalSpPower.String(),
			TotalClientPower:      totalClientPower.String(),
			TotalTokenHolderPower: totalTokenPower.String(),
			TotalDeveloperPower:   totalDeveloperPower.String(),
			VotePowers:            votePowerList,
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
		db.VoteResult(proposal.Id, voteHistory, voteResultList)
	}
}
