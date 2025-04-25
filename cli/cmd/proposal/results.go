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
	voteTable.SetHeader([]string{"Voter Address", "Token Power", "Power Percentage", "Result"})
	voteTable.SetBorder(true)
	voteTable.SetRowLine(true)
	voteTable.SetAutoFormatHeaders(true)
	voteTable.SetAutoWrapText(true)
	voteTable.SetColumnSeparator("|")
	voteTable.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})

	// Get total token power once and convert it for percentage calculations
	tokenHolderPower, err := convertStringToBigInt(proposal.TotalPower.TokenHolderPower)
	if err != nil {
		logError("Failed to convert TokenHolderPower", err, proposal.TotalPower.TokenHolderPower)
		return
	}

	// Iterate over the votes and calculate the power percentage
	for _, v := range votes {
		tokenPower, err := convertStringToBigInt(v.TokenHolderPower)
		if err != nil {
			logError("Failed to convert TokenHolderPower", err, v.TokenHolderPower)
			return
		}

		tokenPowerInTFIL := new(big.Float).SetInt(tokenPower)
		tokenPowerInTFIL = tokenPowerInTFIL.Quo(tokenPowerInTFIL, big.NewFloat(1e18))

		var tokenPowerWithUnit string
		if config.Client.Network.ChainID == 314 {
			tokenPowerWithUnit = fmt.Sprintf("%.2f FIL", tokenPowerInTFIL)
		} else if config.Client.Network.ChainID == 314159 {
			tokenPowerWithUnit = fmt.Sprintf("%.2f tFIL", tokenPowerInTFIL)
		}

		// Calculate the percentage of token power
		percentage := calculatePercentage(tokenPower, tokenHolderPower)

		// Determine the result of the vote (Approve/Reject)
		votedResult := getVotedResult(v.VotedResult)

		// Add the vote data to the vote table
		voteTable.Append([]string{
			v.VoterAddress,
			tokenPowerWithUnit,
			percentage,
			votedResult,
		})
	}

	// Render the vote table
	voteTable.Render()
}

// calculatePercentage computes the percentage of token power in relation to total token holder power
func calculatePercentage(tokenPower, tokenHolderPower *big.Int) string {
	// Special case: If tokenPower is zero, return 0%
	if tokenPower.Cmp(big.NewInt(0)) == 0 {
		return "0.00%"
	}

	percentage := new(big.Float).SetInt(tokenPower)
	percentage.Quo(percentage, new(big.Float).SetInt(tokenHolderPower))
	percentage.Mul(percentage, big.NewFloat(100)) // Convert to percentage

	// Format percentage to 2 decimal places
	return fmt.Sprintf("%.2f%%", percentage)
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
