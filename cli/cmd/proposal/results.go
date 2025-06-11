package proposal

import (
	"fil-vote/config"
	"fil-vote/model"
	"fil-vote/service"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"math/big"
	"os"
	"strconv"
)

// ResultsCmd returns a command to list proposals with pagination.
func ResultsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "results [proposal ID]",
		Short: "Display the details of a proposal along with its voting results (if completed)",
		Run: func(cmd *cobra.Command, args []string) {
			// Check if proposal ID is provided
			if len(args) == 0 {
				zap.L().Info("No proposal ID provided")
				return
			}

			proposalIDStr := args[0]
			proposalID, err := strconv.ParseInt(proposalIDStr, 10, 64)
			if err != nil {
				logError("Invalid proposal ID", err, proposalIDStr)
				return
			}

			// Fetch proposal content by ID
			proposal, err := service.GetProposalByID(proposalID)
			if err != nil {
				logError("Error fetching proposal content", err, proposalID)
				return
			}

			// Validate proposal ID
			if proposal.ProposalId == 0 {
				zap.L().Info("Invalid proposal ID", zap.Int64("proposalID", proposalID))
				return
			}

			// Fetch votes if proposal status is completed
			if proposal.Status == 4 {
				votes, err := service.GetVotesByProposalID(proposalID)
				if err != nil {
					logError("Error fetching votes", err, proposalID)
					return
				}
				if len(votes) == 0 {
					fmt.Println("No voting")
					return
				}

				// Print the proposal and votes
				printProposalContents(proposal, votes)
			} else {
				// Log proposal status info
				zap.L().Info("Proposal is not completed, no voting results available", zap.Int64("proposalID", proposalID), zap.String("status", proposal.Status.String()))
			}
		},
	}
	return cmd
}

func logError(message string, err error, proposalID interface{}) {
	zap.L().Error(message, zap.Any("proposalID", proposalID), zap.Error(err))
}

// printProposalContents prints the proposal and votes in a formatted table
func printProposalContents(proposal model.Proposal, votes []model.Vote) {
	voteTable := tablewriter.NewWriter(os.Stdout)
	voteTable.SetHeader([]string{"Voter Address", "SP Power", "Client Power", "Developer Power", "Token Holder Power", "Power Percentage", "Result"})
	voteTable.SetBorder(true)
	voteTable.SetRowLine(true)
	voteTable.SetAutoFormatHeaders(true)
	voteTable.SetAutoWrapText(true)
	voteTable.SetColumnSeparator("|")
	voteTable.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})

	totalPercent := 0
	if proposal.TotalPower.TokenHolderPower != "0" {
		totalPercent += proposal.Percentage.TokenHolderPercentage
	}
	if proposal.TotalPower.ClientPower != "0" {
		totalPercent += proposal.Percentage.ClientPercentage
	}
	if proposal.TotalPower.DeveloperPower != "0" {
		totalPercent += proposal.Percentage.DeveloperPercentage
	}
	if proposal.TotalPower.SpPower != "0" {
		totalPercent += proposal.Percentage.SpPercentage
	}

	// Iterate over the votes and calculate the power percentage
	for _, v := range votes {
		tokenPower, err := convertStringToBigInt(v.TokenHolderPower)
		if err != nil {
			logError("Failed to convert TokenHolderPower", err, v.TokenHolderPower)
			return
		}
		convertedTokenHolderPower := model.ConvertToFIL(tokenPower, config.Client.Network.ChainID)

		spPower, _ := new(big.Int).SetString(v.SpPower, 10)
		clientPower, _ := new(big.Int).SetString(v.ClientPower, 10)
		developerPower, _ := new(big.Int).SetString(v.DeveloperPower, 10)
		tokenHolderPower, _ := new(big.Int).SetString(v.TokenHolderPower, 10)

		totalSpPower, _ := new(big.Int).SetString(proposal.TotalPower.SpPower, 10)
		totalClientPower, _ := new(big.Int).SetString(proposal.TotalPower.ClientPower, 10)
		totalDeveloperPower, _ := new(big.Int).SetString(proposal.TotalPower.DeveloperPower, 10)
		totalTokenHolderPower, _ := new(big.Int).SetString(proposal.TotalPower.TokenHolderPower, 10)

		powerPercent := new(big.Float)
		if tokenHolderPower.Cmp(big.NewInt(0)) != 0 && totalTokenHolderPower.Cmp(big.NewInt(0)) != 0 {
			tokenHolderPowerPercentage := new(big.Float).SetInt(tokenHolderPower)
			tokenHolderPowerPercentage.Quo(tokenHolderPowerPercentage, new(big.Float).SetInt(totalTokenHolderPower))
			tokenHolderPercentage := new(big.Float).SetInt64(int64(proposal.Percentage.TokenHolderPercentage))
			tokenHolderPowerPercentage.Mul(tokenHolderPowerPercentage, tokenHolderPercentage)
			powerPercent.Add(powerPercent, tokenHolderPowerPercentage)
		}

		if spPower.Cmp(big.NewInt(0)) != 0 && totalSpPower.Cmp(big.NewInt(0)) != 0 {
			spPowerPercentage := new(big.Float).SetInt(spPower)
			spPowerPercentage.Quo(spPowerPercentage, new(big.Float).SetInt(totalSpPower))
			spPercentage := new(big.Float).SetInt64(int64(proposal.Percentage.SpPercentage))
			spPowerPercentage.Mul(spPowerPercentage, spPercentage)
			powerPercent.Add(powerPercent, spPowerPercentage)
		}

		if clientPower.Cmp(big.NewInt(0)) != 0 && totalClientPower.Cmp(big.NewInt(0)) != 0 {
			clientPowerPercentage := new(big.Float).SetInt(clientPower)
			clientPowerPercentage.Quo(clientPowerPercentage, new(big.Float).SetInt(totalClientPower))
			clientPercentage := new(big.Float).SetInt64(int64(proposal.Percentage.ClientPercentage))
			clientPowerPercentage.Mul(clientPowerPercentage, clientPercentage)
			powerPercent.Add(powerPercent, clientPowerPercentage)
		}

		if developerPower.Cmp(big.NewInt(0)) != 0 && totalDeveloperPower.Cmp(big.NewInt(0)) != 0 {
			developerPowerPercentage := new(big.Float).SetInt(developerPower)
			developerPowerPercentage.Quo(developerPowerPercentage, new(big.Float).SetInt(totalDeveloperPower))
			developerPercentage := new(big.Float).SetInt64(int64(proposal.Percentage.DeveloperPercentage))
			developerPowerPercentage.Mul(developerPowerPercentage, developerPercentage)
			powerPercent.Add(powerPercent, developerPowerPercentage)
		}

		if powerPercent.Cmp(big.NewFloat(0)) != 0 {
			totalPercentFloat := new(big.Float).SetInt64(int64(totalPercent))
			powerPercent.Quo(powerPercent, totalPercentFloat)
			powerPercent.Mul(powerPercent, big.NewFloat(100))
		} else {
			powerPercent = new(big.Float)
		}

		// Determine the result of the vote (Approve/Reject)
		votedResult := getVotedResult(v.VotedResult)

		// Add the vote data to the vote table
		voteTable.Append([]string{
			v.VoterAddress,
			v.SpPower,
			v.ClientPower,
			v.DeveloperPower,
			convertedTokenHolderPower,
			fmt.Sprintf("%.2f%%", powerPercent),
			votedResult,
		})
	}

	voteTable.Render()
}

// getVotedResult converts the vote result into a string (Approve/Reject)
func getVotedResult(result string) string {
	switch result {
	case model.Approve:
		return model.Approved
	case model.Reject:
		return model.Rejected
	default:
		return "Unknown"
	}
}

// convertStringToBigInt safely converts a string to a big.Int and returns an error if it fails
func convertStringToBigInt(value string) (*big.Int, error) {
	result := new(big.Int)
	_, success := result.SetString(value, 10)
	if !success {
		return nil, fmt.Errorf("invalid number format for value: %s", value)
	}
	return result, nil
}
