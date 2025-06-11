package power

import (
	"context"
	"fil-vote/config"
	"fil-vote/service"
	"fil-vote/utils"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"math/big"
	"os"
	"strconv"
)

func LsCmd(client *service.RPCClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List wallet powers for a specific day",
		Run: func(cmd *cobra.Command, args []string) {
			// Retrieve the 'day' flag value
			day, err := cmd.Flags().GetString("day")
			if err != nil || day == "" {
				// Handle invalid day input and print clear error
				zap.L().Error("Invalid day", zap.Error(err))
				return
			}

			// Fetch the list of wallets
			wallets, err := client.ListWallets(context.Background())
			if err != nil {
				return
			}

			// Check if no wallets were found
			if len(wallets) == 0 {
				zap.L().Error("No wallets found", zap.Error(err))
				return
			}

			// Prepare the table for output
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Wallet", "SP Power", "Client Power", "Developer Power", "Token Holder Power"})
			table.SetBorder(true)
			table.SetRowLine(true)
			table.SetAutoFormatHeaders(true)
			table.SetAutoWrapText(true)
			table.SetColumnSeparator("|")

			// Display the power for each wallet
			fmt.Printf("\nListing wallet powers for day: %s\n", day)
			for _, v := range wallets {
				// Get the power for each wallet
				power, err := service.GetPower(day, v)
				if err != nil {
					// Handle individual wallet errors and continue with the next wallet
					zap.L().Error("Failed to retrieve power", zap.String("wallet", v), zap.Error(err))
					fmt.Printf("Error: Failed to retrieve power for wallet %s\n", v)
					continue // Continue processing the remaining wallets
				}

				// Convert the Token Holder Power from attoFIL to a readable format
				tokenHolderPower := new(big.Int)
				// Assume power.Data.TokenHolderPower is in attoFIL, and convert it
				tokenHolderPower.SetString(power.Data.TokenHolderPower, 10)
				convertedTokenHolderPower := utils.ConvertToFIL(tokenHolderPower, config.Client.Network.ChainID)

				spPower, err := strconv.ParseInt(power.Data.SpPower, 10, 64)
				convertedSpPower := utils.ConvertSize(spPower)

				table.Append([]string{
					v,
					convertedSpPower,
					power.Data.ClientPower,
					power.Data.DeveloperPower,
					convertedTokenHolderPower,
				})
			}

			// Render the table
			table.Render()
		},
	}

	// Add the 'day' flag to the command
	cmd.Flags().String("day", "", "The day to retrieve the power for (required)")
	cmd.MarkFlagRequired("day")

	return cmd
}
