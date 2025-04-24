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
	"math/big"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	snapshot "powervoting-server/api/rpc"
	"powervoting-server/config"
	"powervoting-server/constant"
	"powervoting-server/data"
	"powervoting-server/model"
	"powervoting-server/service"
	"powervoting-server/utils"
)

var oneHundred = decimal.NewFromInt(100)

type VoteCount struct {
	EthClient   *model.GoEthClient
	SyncService service.ISyncService
}

// VotingCountHandler initiates the voting count process.
// It iterates through the network configurations and retrieves the Ethereum client for each network.
// It then launches a goroutine to handle the voting count for each network.
// Any errors encountered during the retrieval of the Ethereum client are logged.
func VotingCountHandler(syncService service.ISyncService) {
	network := config.Client.Network
	ethClient, err := data.GetClient(syncService, network.ChainId)
	if err != nil {
		zap.L().Error("get go-eth client error:", zap.Error(err))
		return
	}

	voteCount := &VoteCount{
		EthClient:   ethClient,
		SyncService: syncService,
	}

	syncEventInfo, err := syncService.GetSyncEventInfo(context.Background(), network.PowerVotingContract)
	if err != nil {
		zap.L().Error("get sync event info error:", zap.Error(err))
		return
	}

	if err := voteCount.voteCounting(syncEventInfo.SyncedHeight); err != nil {
		zap.L().Error("vote count with err:", zap.Error(err))
		return
	}

	zap.L().Info(
		"voting count finished",
		zap.String("network name", network.Name),
		zap.Int64("network id", network.ChainId),
	)
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
		// Retrieve the snapshot of voting power for all addresses at the specified block height
		allPowers, err := snapshot.GetAllAddressPowerByDay(vc.EthClient.ChainId, p.SnapshotDay)
		if err != nil {
			errList = append(errList, fmt.Errorf("get address power for proposal %d: %w", p.ProposalId, err))
		}

		if err := vc.processCounting(vc.EthClient, p, allPowers); err != nil {
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
func (vc *VoteCount) processCounting(ethClient *model.GoEthClient, proposal model.ProposalTbl, allPowers model.SnapshotAllPower) error {
	// Retrieve all uncounted votes for the proposal
	votesInfo, err := vc.SyncService.GetUncountedVotedList(context.Background(), ethClient.ChainId, proposal.ProposalId)
	if err != nil {
		return fmt.Errorf("get vote info for proposal %d: %w", proposal.ProposalId, err)
	}

	// Convert the address power information to a map for quick lookup
	powersMap := utils.PowersInfoToMap(allPowers.AddrPower)
	// powersMap := make(map[string]model.AddrPower)
	// Calculate the total voting power and update the vote list with weights
	creditsMap, totalCredits, voteList := vc.countWeightCredits(proposal.ProposalId, powersMap, votesInfo, ethClient.ChainId)

	approvePercentage := vc.calculateVotesPercentage(creditsMap[constant.VoteApprove], totalCredits, proposal.Percentage)
	rejectPercentage := vc.calculateVotesPercentage(creditsMap[constant.VoteReject], totalCredits, proposal.Percentage)
	// Calculate the final approval and rejection percentages
	resultPercent := vc.calculateFinalPercentages(approvePercentage, rejectPercentage, len(voteList))
	zap.L().Info("Final percentages",
		zap.Int64("proposal ID", proposal.ProposalId),
		zap.Float64(constant.VoteApprove, resultPercent[constant.VoteApprove]),
		zap.Float64(constant.VoteReject, resultPercent[constant.VoteReject]))

	// Update the proposal with the calculated results
	proposal.Counted = constant.ProposalCounted
	proposal.ProposalResult.ApprovePercentage = resultPercent[constant.VoteApprove]
	proposal.ProposalResult.RejectPercentage = resultPercent[constant.VoteReject]
	proposal.TotalSpPower = totalCredits.SpPower.String()
	proposal.TotalClientPower = totalCredits.ClientPower.String()
	proposal.TotalTokenHolderPower = totalCredits.TokenPower.String()
	proposal.TotalDeveloperPower = totalCredits.DeveloperPower.String()

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

// countWeightCredits calculates the voting weight and percentages for a given proposal.
// It aggregates the voting power of all participants and computes the approval and rejection percentages.
// The function processes each vote record, retrieves the corresponding voting power, and updates the results.
//
// Parameters:
//   - proposalId: The ID of the proposal being calculated.
//   - powersMap: A map containing the voting power of each address, indexed by address.
//   - voteInfos: A list of vote records for the proposal, including the voter's address and encrypted vote result.
//   - chainId: The chain ID associated with the proposal, used for logging and context.
//
// Returns:
//   - map[string]model.VoterPowerCount: A map containing the aggregated voting power for each vote result (e.g., "approve", "reject").
//   - model.VoterPowerCount: A struct containing the total voting power across all vote results.
//   - []model.VoteTbl: A list of updated vote records with calculated voting weights.
func (vc *VoteCount) countWeightCredits(proposalId int64, powersMap map[string]model.AddrPower, voteInfos []model.VoteTbl, chainId int64) (map[string]model.VoterPowerCount, model.VoterPowerCount, []model.VoteTbl) {
	// Initialize the vote power struct with zero values
	var (
		creditsMap = map[string]model.VoterPowerCount{
			constant.VoteApprove: {SpPower: decimal.Zero, ClientPower: decimal.Zero, TokenPower: decimal.Zero, DeveloperPower: decimal.Zero},
			constant.VoteReject:  {SpPower: decimal.Zero, ClientPower: decimal.Zero, TokenPower: decimal.Zero, DeveloperPower: decimal.Zero},
		} // Map to store vote results and their percentages
		voteList []model.VoteTbl
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
			// If the address votes at the time,
			// there may be a situation where the snapshot is synchronized with the address computing power.
			// Here is a protection
			power = model.AddrPower{
				Address:          voteInfo.Address,
				SpPower:          big.NewInt(0),
				ClientPower:      big.NewInt(0),
				TokenHolderPower: big.NewInt(0),
				DeveloperPower:   big.NewInt(0),
				BlockHeight:      0,
			}
		}

		// Decode the encrypted vote result
		voteResult, err := utils.DecodeVoteResult(voteInfo.VoteEncrypted)
		if err != nil {
			zap.L().Error("Decode vote info error", zap.Error(err))
			continue
		}

		// Update the aggregated voting power for the current vote result
		sp := creditsMap[voteResult].SpPower.Add(decimal.NewFromBigInt(power.SpPower, 0))
		cp := creditsMap[voteResult].ClientPower.Add(decimal.NewFromBigInt(power.ClientPower, 0))
		tp := creditsMap[voteResult].TokenPower.Add(decimal.NewFromBigInt(power.TokenHolderPower, 0))
		dp := creditsMap[voteResult].DeveloperPower.Add(decimal.NewFromBigInt(power.DeveloperPower, 0))

		creditsMap[voteResult] = model.VoterPowerCount{
			SpPower:        sp,
			ClientPower:    cp,
			TokenPower:     tp,
			DeveloperPower: dp,
		}

		// Update the vote record with the calculated values
		voteList = append(voteList, model.VoteTbl{
			ProposalId:       proposalId,
			Address:          voteInfo.Address,
			VoteResult:       voteResult,
			SpPower:          power.SpPower.String(),
			ClientPower:      power.ClientPower.String(),
			TokenHolderPower: power.TokenHolderPower.String(),
			DeveloperPower:   power.DeveloperPower.String(),
		})
	}

	// Calculate the total voting power across all vote results
	totalCredits := model.VoterPowerCount{
		SpPower:        creditsMap[constant.VoteApprove].SpPower.Add(creditsMap[constant.VoteReject].SpPower),
		ClientPower:    creditsMap[constant.VoteApprove].ClientPower.Add(creditsMap[constant.VoteReject].ClientPower),
		TokenPower:     creditsMap[constant.VoteApprove].TokenPower.Add(creditsMap[constant.VoteReject].TokenPower),
		DeveloperPower: creditsMap[constant.VoteApprove].DeveloperPower.Add(creditsMap[constant.VoteReject].DeveloperPower),
	}

	return creditsMap, totalCredits, voteList
}

// calculateVotesPercentage calculates the weighted voting percentage for a given voter based on their power contributions.
// It computes the weighted average of the voter's power in different categories (e.g., SP, Client, Token, Developer)
// relative to the total power in each category, and applies the respective percentage weights.
//
// Parameters:
//   - votesPower: A struct containing the voter's power in each category (SP, Client, Token, Developer).
//   - totalPower: A struct containing the total power in each category for all voters.
//   - percentage: A struct containing the percentage weights for each category.
//
// Returns:
//   - decimal.Decimal: The calculated weighted voting percentage, scaled to a value between 0 and 100.
func (vc *VoteCount) calculateVotesPercentage(votesPower model.VoterPowerCount, totalPower model.VoterPowerCount, percentage model.Percentage) decimal.Decimal {
	var (
		totalPercentage = uint16(0)             // Accumulates the total percentage weights used in the calculation
		totalWeight     = decimal.NewFromInt(0) // Accumulates the weighted sum of the voter's power contributions
	)

	// Calculate the weighted contribution for SP (Storage Provider) power
	if !totalPower.SpPower.IsZero() {
		if !votesPower.SpPower.IsZero() {
			totalWeight = totalWeight.Add(votesPower.SpPower.Div(totalPower.SpPower).Mul(decimal.NewFromInt(int64(percentage.SpPercentage))))
		}
		totalPercentage += percentage.SpPercentage
		zap.L().Debug("SP Percentage", zap.Uint16("percentage", percentage.SpPercentage), zap.String("totalWeight", totalWeight.String()))
	}

	// Calculate the weighted contribution for Client power
	if !totalPower.ClientPower.IsZero() {
		if !votesPower.ClientPower.IsZero() {
			totalWeight = totalWeight.Add(votesPower.ClientPower.Div(totalPower.ClientPower).Mul(decimal.NewFromInt(int64(percentage.ClientPercentage))))
		}
		totalPercentage += percentage.ClientPercentage
		zap.L().Debug("Client Percentage", zap.Uint16("percentage", percentage.ClientPercentage), zap.String("totalWeight", totalWeight.String()))
	}

	// Calculate the weighted contribution for Token Holder power
	if !totalPower.TokenPower.IsZero() {
		if !votesPower.TokenPower.IsZero() {
			totalWeight = totalWeight.Add(votesPower.TokenPower.Div(totalPower.TokenPower).Mul(decimal.NewFromInt(int64(percentage.TokenHolderPercentage))))
		}
		totalPercentage += percentage.TokenHolderPercentage
		zap.L().Debug("Token Percentage", zap.Uint16("percentage", percentage.TokenHolderPercentage), zap.String("totalWeight", totalWeight.String()))
	}

	// Calculate the weighted contribution for Developer power
	if !totalPower.DeveloperPower.IsZero() {
		if !votesPower.DeveloperPower.IsZero() {
			totalWeight = totalWeight.Add(votesPower.DeveloperPower.Div(totalPower.DeveloperPower).Mul(decimal.NewFromInt(int64(percentage.DeveloperPercentage))))
		}
		totalPercentage += percentage.DeveloperPercentage
		zap.L().Debug("Developer Percentage", zap.Uint16("percentage", percentage.DeveloperPercentage), zap.String("totalWeight", totalWeight.String()))
	}

	// If no percentage weights were applied, return 0
	if totalPercentage == 0 {
		return decimal.NewFromInt(0)
	}

	zap.L().Info("calculateVotesPercentage", zap.String("totalWeight", totalWeight.String()), zap.Uint16("totalPercentage", totalPercentage))
	// Normalize the weighted sum to a percentage value (0-100)
	return totalWeight.Div(decimal.NewFromInt(int64(totalPercentage))).Mul(oneHundred)
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

	offset := decimal.NewFromInt(100).Sub(approve).Sub(reject)

	zap.L().Info("calculateFinalPercentages", zap.Float64("offset", offset.InexactFloat64()))
	// Process remainder allocation
	switch {
	case approve.GreaterThan(reject):
		approve = approve.Add(offset)
	case reject.GreaterThan(approve):
		reject = reject.Add(offset)
	default:
		approve = decimal.NewFromFloat(50)
		reject = decimal.NewFromFloat(50)
	}

	zap.L().Info("Final percentages", zap.Any("approve", approve), zap.Any("reject", reject))
	// Rounding process
	return map[string]float64{
		constant.VoteApprove: approve.Round(2).InexactFloat64(),
		constant.VoteReject:  reject.Round(2).InexactFloat64(),
	}
}
