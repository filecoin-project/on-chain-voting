package wallet

import (
	"context"
	"fil-vote/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

// UseCmd returns a Cobra command to set a wallet as the default and display the status in a table.
func UseCmd(client *service.RPCClient) *cobra.Command {
	return &cobra.Command{
		Use:   "use [walletAddress]",                                   // Command name with an argument for wallet address
		Short: "Set a wallet as the default wallet for the Lotus node", // Brief description
		Args:  cobra.ExactArgs(1),                                      // Ensure exactly one argument is provided
		Run: func(cmd *cobra.Command, args []string) {
			// Get the wallet address from the argument
			walletAddress := args[0]

			// Set the provided wallet as the default
			err := client.WalletSetDefault(context.Background(), walletAddress)
			if err != nil {
				// Display error in the table if something goes wrong
				createWalletUseTable(walletAddress, "Failed", err.Error())
			} else {
				// Display success in the table
				createWalletUseTable(walletAddress, "Success", "")
			}
		},
	}
}

// displayWalletStatus formats and displays the wallet address and status in a table.
func createWalletUseTable(walletAddress, status, errMsg string) {
	// Initialize a table for displaying the result
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Wallet Address", "Status", "Error Message"})
	table.SetAutoFormatHeaders(false)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})

	// Add the row with wallet address and status
	if status == "Success" {
		table.Append([]string{walletAddress, status, ""})
	} else {
		table.Append([]string{walletAddress, status, errMsg})
	}

	// Render the table
	table.Render()
}
