package developer

import (
	"context"
	"fil-vote/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

// UpdateGistIdCmd creates a new command that updates the Gist ID for the GitHub proof
func UpdateGistIdCmd(client *service.RPCClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-github-proof ",
		Short: "Claims a GitHub proof with a specific Gist ID",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString("from")
			gistId, err := cmd.Flags().GetString("gistId")

			if err != nil || gistId == "" {
				zap.L().Error("Invalid gistId", zap.Error(err))
				return
			}

			// If the "from" flag is empty, retrieve the default wallet address
			if from == "" {
				var err error
				from, err = client.WalletDefaultAddress(context.Background())
				if err != nil {
					zap.L().Error("Failed to retrieve default wallet address", zap.Error(err))
					return
				}
			}

			// Update the Gist ID with the proof
			messageHash, err := service.UpdateGistId(client, from, gistId)
			if err != nil {
				zap.L().Error("Failed to update Gist ID", zap.Error(err))
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Message Hash"})
			table.SetRowLine(true)

			table.Append([]string{messageHash})

			table.Render()
		},
	}

	cmd.Flags().String("from", "", "Wallet address to generate the proof for (optional)")
	cmd.Flags().String("gistId", "", "Wallet address to generate the proof for (required)")
	cmd.MarkFlagRequired("gistId")
	return cmd
}
