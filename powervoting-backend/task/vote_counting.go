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
	"context"
	"fmt"
	"sync"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"powervoting-server/config"
	"powervoting-server/constant"
	"powervoting-server/data"
	"powervoting-server/model"
	"powervoting-server/service"
	"powervoting-server/snapshot"
	"powervoting-server/utils"
)

type TotalPower struct {
	spPower           decimal.Decimal
	clientPower       decimal.Decimal
	tokenPower        decimal.Decimal
	developerPower    decimal.Decimal
	totalPowerWeight  decimal.Decimal
	approvePercentage decimal.Decimal
	rejectPercentage  decimal.Decimal
}

var oneHundred = decimal.NewFromInt(100)

type VoteCount struct {
	EthClient   *model.GoEthClient
	SyncService *service.SyncService
}

// VotingCountHandler initiates the voting count process.
// It iterates through the network configurations and retrieves the Ethereum client for each network.
// It then launches a goroutine to handle the voting count for each network.
// Any errors encountered during the retrieval of the Ethereum client are logged.
func VotingCountHandler(syncService *service.SyncService) {
	networkList := config.Client.Network
	wg := sync.WaitGroup{}
	errList := make([]error, 0, len(config.Client.Network))
	mu := &sync.Mutex{}

	for _, network := range networkList {
		network := network
		ethClient, err := data.GetClient(syncService, network.ChainId)
		if err != nil {
			zap.L().Error("get go-eth client error:", zap.Error(err))
			continue
		}

		voteCount := &VoteCount{
			EthClient:   ethClient,
			SyncService: syncService,
		}

		syncEventInfo, err := syncService.GetSyncEventInfo(context.Background(), network.PowerVotingContract)
		if err != nil {
			zap.L().Error("get sync event info error:", zap.Error(err))
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := voteCount.voteCounting(syncEventInfo.SyncedHeight); err != nil {
				mu.Lock()
				errList = append(errList, err)
				mu.Unlock()
			}

			zap.L().Info(
				"voting count finished",
				zap.String("network name", network.Name),
				zap.Int64("network id", network.ChainId),
			)
		}()
	}

	wg.Wait()
	if len(errList) != 0 {
		zap.L().Error("vote count with err:", zap.Errors("errors", errList))
	}
}

// voteCounting is the main function that handles the voting count process.
// It retrieves pending proposals and processes them concurrently.
// Returns an error if any proposal processing fails.
func (vc *VoteCount) voteCounting(syncedHeight int64) error {
	syncedTimestamp, err := utils.GetBlockTime(vc.EthClient, syncedHeight)
	if err != nil {
		return fmt.Errorf("failed to get latest block timestamp: %w", err)
	}

	zap.L().Info("synced timestamp", zap.Int64("timestamp", syncedTimestamp))
	// Get pending proposals
	proposals, err := vc.SyncService.UncountedProposalList(context.Background(), vc.EthClient.ChainId, syncedTimestamp)
	if err != nil {
		zap.L().Error("get pending proposals error:", zap.Error(err))
		return fmt.Errorf("failed to get proposal list: %w", err)
	}

	// Parallel handling of proposals
	errList := make([]error, 0, len(proposals))

	for _, p := range proposals {
		if err := vc.processCounting(vc.EthClient, p); err != nil {
			errList = append(errList, fmt.Errorf("count voted failed: %w, proposal ID %d ", err, p.ProposalId))
		}
	}

	if len(errList) > 0 {
		return fmt.Errorf("completed processing but encountered %d error: %+v", len(errList), errList)
	}

	return nil
}

// processCounting handles the vote counting for a single proposal.
// It retrieves all uncounted votes for the proposal, calculates the total voting power,
// processes individual votes, and persists the results to the database.
//
// Parameters:
//   - ethClient: The Ethereum client used to interact with the blockchain.
//   - proposal: The proposal for which votes are being counted.
//
// Returns:
//   - error: An error if any step in the process fails; otherwise, nil.
func (vc *VoteCount) processCounting(ethClient *model.GoEthClient, proposal model.ProposalTbl) error {
	// Retrieve all uncounted votes for the proposal
	votesInfo, err := vc.SyncService.GetUncountedVotedList(context.Background(), ethClient.ChainId, proposal.ProposalId)
	if err != nil {
		return fmt.Errorf("get vote info for proposal %d: %w", proposal.ProposalId, err)
	}

	// Retrieve the snapshot of voting power for all addresses at the specified block height
	allPowers, err := snapshot.GetAllAddressPowerByDay(ethClient.ChainId, proposal.SnapshotDay)
	if err != nil {
		return fmt.Errorf("get address power for proposal %d: %w", proposal.ProposalId, err)
	}

	// Convert the address power information to a map for quick lookup
	powersMap := utils.PowersInfoToMap(allPowers.AddrPower)

	// Calculate the total voting power and update the vote list with weights
	totalPower, voteList := vc.calculateVoteWeight(proposal.ProposalId, powersMap, votesInfo, proposal.Percentage, ethClient.ChainId)
	zap.L().Info("total power calculated",
		zap.String("sp power", totalPower.spPower.String()),
		zap.String("client power", totalPower.clientPower.String()),
		zap.String("token power", totalPower.tokenPower.String()),
		zap.String("developer power", totalPower.developerPower.String()))

	// Calculate the final approval and rejection percentages
	resultPercent := vc.calculateFinalPercentages(totalPower.approvePercentage, totalPower.rejectPercentage, len(voteList))
	zap.L().Info("Final percentages",
		zap.Int64("proposal ID", proposal.ProposalId),
		zap.Float64(constant.VoteApprove, resultPercent[constant.VoteApprove]),
		zap.Float64(constant.VoteReject, resultPercent[constant.VoteReject]))

	// Update the proposal with the calculated results
	proposal.Counted = constant.ProposalCounted
	proposal.ProposalResult.ApprovePercentage = resultPercent[constant.VoteApprove]
	proposal.ProposalResult.RejectPercentage = resultPercent[constant.VoteReject]
	proposal.TotalSpPower = totalPower.spPower.String()
	proposal.TotalClientPower = totalPower.clientPower.String()
	proposal.TotalTokenHolderPower = totalPower.tokenPower.String()
	proposal.TotalDeveloperPower = totalPower.developerPower.String()

	// Persist the updated proposal to the database
	if err := vc.SyncService.UpdateProposal(context.Background(), &proposal); err != nil {
		return fmt.Errorf("update proposal %d: %w", proposal.ProposalId, err)
	}
	zap.L().Info(
		"Proposal voted completed",
		zap.Int64("proposal ID", proposal.ProposalId),
		zap.Any("result", proposal.ProposalResult),
	)

	// Batch update the vote records in the database
	if err := vc.SyncService.BatchUpdateVotes(context.Background(), voteList); err != nil {
		return fmt.Errorf("batch update vote: %w", err)
	}
	zap.L().Info("Batch update vote completed", zap.Any("vote power number", len(voteList)))

	return nil
}

// calculateVoteWeight calculates the voting weight and percentages for a given proposal.
// It aggregates the voting power of all participants and computes the approval and rejection percentages.
//
// Parameters:
//   - proposalId: The ID of the proposal being calculated.
//   - voteInfos: A list of vote records for the proposal.
//   - proposalPercentage: The percentage distribution of voting power among different roles.
//   - chainId: The chain ID associated with the proposal.
//   - snapshotBlockHeight: The block height at which the snapshot of voting power was taken.
//
// Returns:
//   - *TotalPower: A struct containing the total voting power and calculated percentages.
//   - []model.VoteTbl: A list of updated vote records with calculated voting weights.
func (vc *VoteCount) calculateVoteWeight(proposalId int64, powersMap map[string]model.AddrPower, voteInfos []model.VoteTbl, proposalPercentage model.Percentage, chainId int64) (*TotalPower, []model.VoteTbl) {
	// Initialize the total power struct with zero values
	totalPower := &TotalPower{
		spPower:          decimal.Zero,
		clientPower:      decimal.Zero,
		tokenPower:       decimal.Zero,
		developerPower:   decimal.Zero,
		totalPowerWeight: decimal.Zero,
	}

	var (
		voteResultMap = make(map[string]decimal.Decimal) // Map to store vote results and their percentages
		voteList      []model.VoteTbl                    // List to store updated vote records
	)

	// Iterate over each vote record
	for _, voteInfo := range voteInfos {
		// Retrieve the voting power for the current address
		power, exists := powersMap[voteInfo.Address]
		if !exists {
			zap.L().Error(
				"address power not found",
				zap.String("address", voteInfo.Address),
				zap.Int64("chainId", chainId),
				zap.Int64("proposalId", voteInfo.ProposalId),
			)
			continue
		}

		// Accumulate the total voting power for each role
		totalPower.spPower = totalPower.spPower.Add(decimal.NewFromBigInt(power.SpPower, 0))
		totalPower.clientPower = totalPower.clientPower.Add(decimal.NewFromBigInt(power.ClientPower, 0))
		totalPower.tokenPower = totalPower.tokenPower.Add(decimal.NewFromBigInt(power.TokenHolderPower, 0))
		totalPower.developerPower = totalPower.developerPower.Add(decimal.NewFromBigInt(power.DeveloperPower, 0))

		// Decode the encrypted vote result
		voteResult, err := utils.DecodeVoteResult(voteInfo.VoteEncrypted)
		if err != nil {
			zap.L().Error("Decode vote info error", zap.Error(err))
			continue
		}

		// Calculate the vote percentage based on the voting power and proposal percentage distribution
		vote, votePercentage := vc.calculateVotePercentage(power, proposalPercentage)
		voteResultMap[voteResult] = voteResultMap[voteResult].Add(votePercentage)

		// Update the vote record with the calculated values
		vote.ProposalId = proposalId
		vote.Address = voteInfo.Address
		vote.VoteResult = voteResult
		voteList = append(voteList, vote)
	}

	// Calculate the weighted voting power for each role
	spPowerWeight := totalPower.spPower.Mul(decimal.NewFromInt(int64(proposalPercentage.SpPercentage)))
	clientPowerWeight := totalPower.clientPower.Mul(decimal.NewFromInt(int64(proposalPercentage.ClientPercentage)))
	tokenPowerWeight := totalPower.tokenPower.Mul(decimal.NewFromInt(int64(proposalPercentage.TokenHolderPercentage)))
	developerPowerWeight := totalPower.developerPower.Mul(decimal.NewFromInt(int64(proposalPercentage.DeveloperPercentage)))

	// Calculate the total weighted voting power
	totalPower.totalPowerWeight = spPowerWeight.Add(clientPowerWeight).Add(tokenPowerWeight).Add(developerPowerWeight)

	// Calculate the approval and rejection percentages
	totalPower.approvePercentage = voteResultMap[constant.VoteApprove].Div(totalPower.totalPowerWeight).Mul(oneHundred)
	totalPower.rejectPercentage = voteResultMap[constant.VoteReject].Div(totalPower.totalPowerWeight).Mul(oneHundred)

	// Log the calculated percentages
	zap.L().Info(
		"calculate proposal voted percentage",
		zap.Int64("proposalId", proposalId),
		zap.String("approvePercentage", totalPower.approvePercentage.String()),
		zap.String("rejectPercentage", totalPower.rejectPercentage.String()),
	)

	return totalPower, voteList
}

// calculateVotePercentage calculates the contribution of a single vote.
// It considers the voter's power distribution across different dimensions
// and applies the configured weight percentages.
func (vc *VoteCount) calculateVotePercentage(power model.AddrPower, proposalPercentage model.Percentage) (model.VoteTbl, decimal.Decimal) {
	// power * percentage
	spPowerWeight := decimal.NewFromBigInt(power.SpPower, 0).Mul(decimal.NewFromInt(int64(proposalPercentage.SpPercentage)))
	clientPowerWeight := decimal.NewFromBigInt(power.ClientPower, 0).Mul(decimal.NewFromInt(int64(proposalPercentage.ClientPercentage)))
	tokenPowerWeight := decimal.NewFromBigInt(power.TokenHolderPower, 0).Mul(decimal.NewFromInt(int64(proposalPercentage.TokenHolderPercentage)))
	developerPowerWeight := decimal.NewFromBigInt(power.DeveloperPower, 0).Mul(decimal.NewFromInt(int64(proposalPercentage.DeveloperPercentage)))

	voterPowerWeight := spPowerWeight.Add(clientPowerWeight).Add(tokenPowerWeight).Add(developerPowerWeight)

	return model.VoteTbl{
		SpPower:          power.SpPower.String(),
		ClientPower:      power.ClientPower.String(),
		TokenHolderPower: power.TokenHolderPower.String(),
		DeveloperPower:   power.DeveloperPower.String(),
	}, voterPowerWeight
}

// calculateFinalPercentages calculates the final voting percentages and handles precision.
// It ensures the total percentage adds up to 100% by distributing any remainder.
func (vc *VoteCount) calculateFinalPercentages(approve, reject decimal.Decimal, totalVotes int) map[string]float64 {
	if totalVotes == 0 {
		return map[string]float64{
			constant.VoteApprove: 0,
			constant.VoteReject:  0,
		}
	}

	offset := decimal.NewFromInt(1).Sub(approve).Sub(reject)

	// Process remainder allocation
	switch {
	case approve.GreaterThan(reject):
		approve = approve.Add(offset)
	case reject.GreaterThan(approve):
		reject = reject.Add(offset)
	default:
		approve = decimal.NewFromFloat(0.5)
		reject = decimal.NewFromFloat(0.5)
	}

	// Rounding process
	return map[string]float64{
		constant.VoteApprove: approve.Mul(oneHundred).Round(2).InexactFloat64(),
		constant.VoteReject:  reject.Mul(oneHundred).Round(2).InexactFloat64(),
	}
}
