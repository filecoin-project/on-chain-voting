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
	"crypto/rand"
	"math"
	"math/big"
	"powervoting-server/config"
	"powervoting-server/constant"
	"powervoting-server/contract"
	"powervoting-server/db"
	"powervoting-server/model"
	"powervoting-server/utils"
	"sync"
	"time"

	"go.uber.org/zap"
)

// VotingCountHandler vote count
func VotingCountHandler() {
	networkList := config.Client.Network
	wg := sync.WaitGroup{}
	errList := make([]error, 0, len(config.Client.Network))
	mu := &sync.Mutex{}

	for _, network := range networkList {
		network := network
		ethClient, err := contract.GetClient(network.Id)
		if err != nil {
			zap.L().Error("get go-eth client error:", zap.Error(err))
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			zap.L().Info("vote count start:", zap.Int64("networkId", network.Id))
			if err := VotingCount(ethClient); err != nil {
				mu.Lock()
				errList = append(errList, err)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	if len(errList) != 0 {
		zap.L().Error("vote count with err:", zap.Errors("errors", errList))
	}
	zap.L().Info("vote count finished: ", zap.Int64("timestamp", time.Now().Unix()))
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

func VotingCount(ethClient model.GoEthClient) error {
	now, err := utils.GetTimestamp(ethClient)
	if err != nil {
		zap.L().Error("get timestamp on chain error: ", zap.Error(err))
		now = time.Now().Unix()
		return err
	}
	proposals, err := db.GetProposalList(ethClient.Id, now)
	zap.L().Info("proposal list: ", zap.Reflect("proposals", proposals))
	if err != nil {
		zap.L().Error("get proposal from db error:", zap.Error(err))
		return err
	}
	for _, proposal := range proposals {
		err := SyncVote(ethClient, proposal.ProposalId)
		if err != nil {
			zap.L().Error("sync vote error:", zap.Error(err))
		}
		voteInfos, err := db.GetVoteList(ethClient.Id, proposal.ProposalId)
		if err != nil {
			zap.L().Error("get vote info from db error:", zap.Error(err))
			continue
		}
		var voteList []model.Vote4Counting
		zap.L().Info("voteInfos: ", zap.Reflect("voteInfos", voteInfos))
		for _, voteInfo := range voteInfos {
			list, err := utils.DecodeVoteList(voteInfo)
			if err != nil {
				zap.L().Error("get vote info from IPFS or decrypt error: ", zap.Error(err))
				return err
			}
			voteList = append(voteList, list...)
		}
		zap.L().Info("voteList: ", zap.Reflect("voteList", voteList))
		var votePowerList []model.VotePower
		// calc total power
		totalSpPower := new(big.Int)
		totalClientPower := new(big.Int)
		totalTokenPower := new(big.Int)
		totalDeveloperPower := new(big.Int)
		powerMap := make(map[string]model.Power)
		addressIsCount := make(map[string]bool)
		for _, vote := range voteList {
			num, err := rand.Int(rand.Reader, big.NewInt(60))
			if err != nil {
				zap.L().Error("Generate random number error: ", zap.Error(err))
				return err
			}
			num.Add(num, big.NewInt(1))
			power, err := utils.GetPower(vote.Address, num, ethClient)
			if err != nil {
				zap.L().Error("get power error: ", zap.Error(err))
				return err
			}
			if power.BlockHeight.Uint64() == 0 {
				voterToPowerStatus, err := utils.GetVoterToPowerStatus(vote.Address, ethClient)
				if err != nil {
					zap.L().Error("get voter to power status error: ", zap.Error(err))
					return err
				}
				if voterToPowerStatus.DayId.Int64() != 0 {
					num, err := rand.Int(rand.Reader, voterToPowerStatus.DayId)
					if err != nil {
						zap.L().Error("Generate random number error: ", zap.Error(err))
						return err
					}
					num.Add(num, big.NewInt(1))
					power, err = utils.GetPower(vote.Address, num, ethClient)
					if err != nil {
						zap.L().Error("get power error: ", zap.Error(err))
						return err
					}
				}
			}

			zap.L().Info("address and power", zap.Reflect("address", vote.Address), zap.Reflect("power", power))
			if !addressIsCount[vote.Address] {
				powerMap[vote.Address] = power
				totalSpPower.Add(totalSpPower, power.SpPower)
				totalClientPower.Add(totalClientPower, power.ClientPower)
				totalTokenPower.Add(totalTokenPower, power.TokenHolderPower)
				totalDeveloperPower.Add(totalDeveloperPower, power.DeveloperPower)
				addressIsCount[vote.Address] = true
			}
		}

		// zap.L().Info("total data",
		// 	zap.String("totalSpPower", totalSpPower.String()),
		// 	zap.String("totalClientPower", totalClientPower.String()),
		// 	zap.String("totalTokenPower", totalTokenPower.String()),
		// 	zap.String("totalDeveloperPower", totalDeveloperPower.String()))

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

				// The vote consists of 4 parts,
				// If any value of part is zero, it will not be included in the votePercent calculation
				validOption := float64(0)
				if totalSpPower.Cmp(big.NewInt(0)) != 0 {
					spPercent = bigIntDiv(power.SpPower, totalSpPower)
					validOption += 1
				}
				if totalClientPower.Cmp(big.NewInt(0)) != 0 {
					clientPercent = bigIntDiv(power.ClientPower, totalClientPower)
					validOption += 1
				}
				if totalTokenPower.Cmp(big.NewInt(0)) != 0 {
					tokenPercent = bigIntDiv(power.TokenHolderPower, totalTokenPower)
					validOption += 1
				}
				if totalDeveloperPower.Cmp(big.NewInt(0)) != 0 {
					developerPercent = bigIntDiv(power.DeveloperPower, totalDeveloperPower)
					validOption += 1
				}

				// zap.L().Info("percent data",
				// 	zap.Float64("spPercent", spPercent),
				// 	zap.Float64("clientPercent", clientPercent),
				// 	zap.Float64("tokenPercent", tokenPercent),
				// 	zap.Float64("developerPercent", developerPercent))

				votes = (spPercent + clientPercent + tokenPercent + developerPercent) / validOption
				votePower := model.VotePower{
					HistoryId:               proposal.ProposalId,
					Address:                 vote.Address,
					OptionId:                vote.OptionId,
					Votes:                   math.Round(votes*10000) / 100,
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

		// Ensure totals add up to 100% , we need to add an offset: bv = (1 - approveVotes - rejectVotes)
		// And ensure that this operation does not disrupt the final result.
		// it would be: if approveVotes > rejectVotes then approveVotes + bv > rejectVotes
		// so:
		// approveVotes > rejectVotes -> approveVotes + bv
		// approveVotes == rejectVotes -> approveVotes = rejectVotes = 0.5
		// approveVotes < rejectVotes -> rejectVotes + bv
		if len(voteList) != 0 {
			approveVotes := result[constant.VoteApprove]
			rejectVotes := result[constant.VoteReject]
			bv := 1 - approveVotes - rejectVotes
			if bv > 0 && math.Abs(bv) > math.SmallestNonzeroFloat64 {
				if approveVotes > rejectVotes {
					approveVotes += bv
				}

				if approveVotes < rejectVotes {
					rejectVotes += bv
				}

				if math.Abs(approveVotes-rejectVotes) < math.SmallestNonzeroFloat64 {
					approveVotes = 0.5
					rejectVotes = 0.5
				}
				result[constant.VoteApprove] = approveVotes
				result[constant.VoteReject] = rejectVotes
				zap.L().Info("approve and reject",
					zap.Float64("approve", approveVotes),
					zap.Float64("reject", rejectVotes))
			}
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
				Votes:      math.Round(result[int64(i)]*10000) / 100,
				Network:    ethClient.Id,
			}
			voteResultList = append(voteResultList, voteResult)
		}
		// Save vote history and vote result to database and update status
		db.VoteResult(proposal.Id, voteHistory, voteResultList)
	}
	return nil
}
