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
	"math/big"
	"powervoting-server/client"
	"powervoting-server/config"
	"powervoting-server/constant"
	"powervoting-server/contract"
	"powervoting-server/db"
	"powervoting-server/model"
	"powervoting-server/utils"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// VotingCountHandler initiates the voting count process.
// It iterates through the network configurations and retrieves the Ethereum client for each network.
// It then launches a goroutine to handle the voting count for each network.
// Any errors encountered during the retrieval of the Ethereum client are logged.
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
			if err := VotingCount(ethClient, db.Engine); err != nil {
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

// bigIntDiv performs division operation on big integers and returns the result as a float64.
// Parameters x and y are pointers to big integers representing the dividend and divisor respectively.
// The return value z represents the division result as a float64.
// If the divisor is zero, it returns 0 to avoid division by zero error.
func bigIntDiv(x *big.Int, y *big.Int) decimal.Decimal {
	xd := decimal.NewFromBigInt(x, 0)
	yd := decimal.NewFromBigInt(y, 0)
	if yd.IsZero() {
		return decimal.Zero
	}

	// align precision
	return xd.Div(yd).Round(5)
}

// VotingCount performs the vote counting process for proposals.
// It retrieves the current timestamp from the blockchain or the local system in case of an error.
// Then, it fetches the list of proposals from the database based on the provided Ethereum client ID and timestamp.
// For each proposal, it synchronizes the votes from the blockchain, calculates the voting power, and aggregates the results.
// Finally, it saves the voting history, vote results, and updates the status of the proposals in the database.
func VotingCount(ethClient model.GoEthClient, db db.DataRepo) error {
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
		err := SyncVote(ethClient, proposal.ProposalId, db)
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

		num, err := rand.Int(rand.Reader, big.NewInt(61))
		if err != nil {
			zap.L().Error("Generate random number error: ", zap.Error(err))
			return err
		}

		var votePowerList []model.VotePower
		// calc total power
		totalSpPower := new(big.Int)
		totalClientPower := new(big.Int)
		totalTokenPower := new(big.Int)
		totalDeveloperPower := new(big.Int)
		powerMap := make(map[string]model.Power)
		addressIsCount := make(map[string]bool)
		for _, vote := range voteList {
			power, err := client.GetAddressPower(ethClient.Id, vote.Address, int32(num.Int64()))
			if err != nil {
				zap.L().Error("get power error: ", zap.Error(err))
				return err
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

		zap.L().Info("total data",
			zap.String("totalSpPower", totalSpPower.String()),
			zap.String("totalClientPower", totalClientPower.String()),
			zap.String("totalTokenPower", totalTokenPower.String()),
			zap.String("totalDeveloperPower", totalDeveloperPower.String()))

		// vote counting
		var result = make(map[int64]decimal.Decimal, 5) // max 5 options
		var resultPercent = make(map[int64]float64, 5)
		decimalOneHundred := decimal.NewFromFloat(100.0)

		for _, vote := range voteList {
			power := powerMap[vote.Address]
			var votes decimal.Decimal
			if vote.Votes != 0 {
				var spPercent decimal.Decimal
				var clientPercent decimal.Decimal
				var tokenPercent decimal.Decimal
				var developerPercent decimal.Decimal

				// The vote consists of 4 parts,
				// If any value of part is zero, it will not be included in the votePercent calculation
				validOption := int64(0)
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

				zap.L().Info("percent data",
					zap.String("spPercent", spPercent.StringFixed(7)),
					zap.String("clientPercent", clientPercent.StringFixed(7)),
					zap.String("tokenPercent", tokenPercent.StringFixed(7)),
					zap.String("developerPercent", developerPercent.StringFixed(7)))

				votes = spPercent.
					Add(clientPercent).
					Add(tokenPercent).
					Add(developerPercent).
					Div(decimal.NewFromInt(validOption))

				votePower := model.VotePower{
					HistoryId:               proposal.ProposalId,
					Address:                 vote.Address,
					OptionId:                vote.OptionId,
					Votes:                   votes.Mul(decimalOneHundred).Round(2).InexactFloat64(),
					SpPower:                 power.SpPower.String(),
					ClientPower:             power.ClientPower.String(),
					TokenHolderPower:        power.TokenHolderPower.String(),
					DeveloperPower:          power.DeveloperPower.String(),
					SpPowerPercent:          spPercent.Mul(decimalOneHundred).Round(2).InexactFloat64(),
					ClientPowerPercent:      clientPercent.Mul(decimalOneHundred).Round(2).InexactFloat64(),
					TokenHolderPowerPercent: tokenPercent.Mul(decimalOneHundred).Round(2).InexactFloat64(),
					DeveloperPowerPercent:   developerPercent.Mul(decimalOneHundred).Round(2).InexactFloat64(),
					PowerBlockHeight:        int64(power.BlockHeight.Uint64()),
				}
				votePowerList = append(votePowerList, votePower)
			}
			result[vote.OptionId] = result[vote.OptionId].Add(votes)
		}

		// Ensure totals add up to 100% , we need to add an offset = (1 - approveVotes - rejectVotes)
		// And ensure that this operation does not disrupt the final result.
		// it would be: if approveVotes > rejectVotes then approveVotes + offset > rejectVotes
		// so:
		// approveVotes > rejectVotes -> approveVotes + offset
		// approveVotes == rejectVotes -> approveVotes = rejectVotes = 0.5
		// approveVotes < rejectVotes -> rejectVotes + offset
		if len(voteList) != 0 {
			approveVotes := result[constant.VoteApprove]
			rejectVotes := result[constant.VoteReject]
			offset := decimal.NewFromFloat(1.0).Sub(approveVotes).Sub(rejectVotes)
			if offset.IsPositive() && !offset.IsZero() {
				if approveVotes.GreaterThan(rejectVotes) {
					approveVotes = approveVotes.Add(offset)
				} else if approveVotes.LessThan(rejectVotes) {
					rejectVotes = rejectVotes.Add(offset)
				} else {
					approveVotes = decimal.NewFromFloat(0.5)
					rejectVotes = decimal.NewFromFloat(0.5)
				}

				zap.L().Info("approve and reject",
					zap.String("approve", approveVotes.StringFixed(5)),
					zap.String("reject", rejectVotes.StringFixed(5)))
			}
			// round five decimal place to four decimal places and recalculate
			if approveVotes.GreaterThan(rejectVotes) {
				approveVotes = approveVotes.Round(4)
				rejectVotes = decimal.NewFromFloat(1.0).Sub(approveVotes)
			} else if approveVotes.LessThan(rejectVotes) {
				rejectVotes = rejectVotes.Round(4)
				approveVotes = decimal.NewFromFloat(1.0).Sub(rejectVotes)
			}
			resultPercent[constant.VoteApprove] = approveVotes.Mul(decimalOneHundred).InexactFloat64()
			resultPercent[constant.VoteReject] = rejectVotes.Mul(decimalOneHundred).InexactFloat64()

			zap.L().Info("result percent",
				zap.Float64("approve", resultPercent[constant.VoteApprove]),
				zap.Float64("reject", resultPercent[constant.VoteReject]))
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
				Votes:      resultPercent[int64(i)],
				Network:    ethClient.Id,
			}
			voteResultList = append(voteResultList, voteResult)
		}
		// Save vote history and vote result to database and update status
		db.VoteResult(proposal.Id, voteHistory, voteResultList)
	}
	return nil
}
