package power

import (
	"context"
	"fil-vote/service"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"math/big"
	"os"
)

// Conversion constants for FIL denominations
const (
	attoFIL  = 1
	femtoFIL = attoFIL * 1_000
	picoFIL  = femtoFIL * 1_000
	nanoFIL  = picoFIL * 1_000
	microFIL = nanoFIL * 1_000
	milliFIL = microFIL * 1_000
	FIL      = milliFIL * 1_000
)

// Convert to FIL denomination
func convertToFIL(value *big.Int) string {
	if value.Cmp(big.NewInt(FIL)) >= 0 {
		// Convert to FIL
		return fmt.Sprintf("%s FIL", new(big.Float).Quo(new(big.Float).SetInt(value), big.NewFloat(FIL)).Text('f', 2))
	} else if value.Cmp(big.NewInt(milliFIL)) >= 0 {
		// Convert to milliFIL
		return fmt.Sprintf("%s milliFIL", new(big.Float).Quo(new(big.Float).SetInt(value), big.NewFloat(milliFIL)).Text('f', 2))
	} else if value.Cmp(big.NewInt(microFIL)) >= 0 {
		// Convert to microFIL
		return fmt.Sprintf("%s microFIL", new(big.Float).Quo(new(big.Float).SetInt(value), big.NewFloat(microFIL)).Text('f', 2))
	} else if value.Cmp(big.NewInt(nanoFIL)) >= 0 {
		// Convert to nanoFIL
		return fmt.Sprintf("%s nanoFIL", new(big.Float).Quo(new(big.Float).SetInt(value), big.NewFloat(nanoFIL)).Text('f', 2))
	} else if value.Cmp(big.NewInt(picoFIL)) >= 0 {
		// Convert to picoFIL
		return fmt.Sprintf("%s picoFIL", new(big.Float).Quo(new(big.Float).SetInt(value), big.NewFloat(picoFIL)).Text('f', 2))
	} else if value.Cmp(big.NewInt(femtoFIL)) >= 0 {
		// Convert to femtoFIL
		return fmt.Sprintf("%s femtoFIL", new(big.Float).Quo(new(big.Float).SetInt(value), big.NewFloat(femtoFIL)).Text('f', 2))
	} else {
		// Convert to attoFIL
		return fmt.Sprintf("%s attoFIL", value.String())
	}
}

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
			table.SetHeader([]string{"Wallet", "Developer Power", "SP Power", "Client Power", "Token Holder Power"})
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
				convertedTokenHolderPower := convertToFIL(tokenHolderPower)

				// Add row to the table with the converted Token Holder Power
				table.Append([]string{
					v,
					power.Data.DeveloperPower,
					power.Data.SpPower,
					power.Data.ClientPower,
					convertedTokenHolderPower,
				})
			}

			// Render the table
			table.Render()
		},
	}

	// Add the 'day' flag to the command
	cmd.Flags().String("day", "", "The day to retrieve the power for (required)")

	return cmd
}
