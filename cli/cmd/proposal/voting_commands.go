package proposal

import (
	"context"
	"errors"
	"fil-vote/model"
	"fil-vote/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

// ApproveCmd creates a command to handle voting on proposals (approve).
func ApproveCmd(client *service.RPCClient) *cobra.Command {
	return createVoteCommand(client, model.Approve, "Vote on a proposal to approve")
}

// RejectCmd creates a command to handle voting on proposals (reject).
func RejectCmd(client *service.RPCClient) *cobra.Command {
	return createVoteCommand(client, model.Reject, "Vote on a proposal to reject")
}

// createVoteCommand is a helper function that creates both approve and reject commands.
func createVoteCommand(client *service.RPCClient, voteType, description string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   voteType,
		Short: description,
		Run: func(cmd *cobra.Command, args []string) {
			from, proposalId, err := retrieveVoteParameters(cmd, client)
			if err != nil {
				zap.L().Error("Failed to retrieve voting parameters", zap.Error(err))
				return
			}

			// Vote on the proposal and handle errors
			messageHash, err := service.Vote(client, from, voteType, proposalId)
			if err != nil {
				zap.L().Error("Failed to submit vote", zap.Error(err))
				return
			}

			// Display the message hash in a table format
			displayMessageHash(messageHash)
		},
	}

	// Add flags for the command
	AddVoteFlags(cmd)
	return cmd
}

// retrieveVoteParameters retrieves voting parameters from the command or defaults if not provided.
func retrieveVoteParameters(cmd *cobra.Command, client *service.RPCClient) (string, int64, error) {
	from, _ := cmd.Flags().GetString("from")
	proposalId, _ := cmd.Flags().GetInt64("proposalId")

	// If the 'from' address is empty, use the default wallet address
	if from == "" {
		var err error
		from, err = client.WalletDefaultAddress(context.Background())
		if err != nil {
			zap.L().Error("Error fetching default wallet address", zap.Error(err))
			return "", 0, err
		}
	}

	// Validate that proposalId is provided
	if proposalId == 0 {
		err := errors.New("proposal ID is required and cannot be zero")
		zap.L().Error("Proposal ID is required and cannot be zero", zap.Int64("proposalId", proposalId))
		return "", 0, err
	}

	return from, proposalId, nil
}

// displayMessageHash formats and displays the message hash in a table.
func displayMessageHash(messageHash string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Message Hash"}) // Table header with only Message Hash
	table.SetRowLine(true)

	// Add the message hash to the table
	table.Append([]string{messageHash})

	// Render the table
	table.Render()
}

// AddVoteFlags defines the command-line flags for the 'vote' command.
func AddVoteFlags(cmd *cobra.Command) {
	cmd.Flags().String("from", "", "The address from which to send the vote (optional)")
	cmd.Flags().Int64("proposalId", 0, "ID of the proposal to vote on (required)")

	// Mark proposalId as required
	cmd.MarkFlagRequired("proposalId")
}
