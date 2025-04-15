package wallet

import (
	"context"
	"fil-vote/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"go.uber.org/zap" // 引入zap日志库
	"os"
)

// LsCmd returns a Cobra command to list wallets connected to the Lotus node.
func LsCmd(client *service.RPCClient) *cobra.Command {
	return &cobra.Command{
		Use:   "ls",                                   // Command name
		Short: "List wallets connected to Lotus node", // Brief description of the command
		Run: func(cmd *cobra.Command, args []string) {
			// Fetch wallets list
			wallets, err := client.ListWallets(context.Background())
			if err != nil {
				zap.L().Error("Error fetching wallets", zap.Error(err))
				return
			}

			// Fetch default wallet address
			defaultAddress, err := client.WalletDefaultAddress(context.Background())
			if err != nil {
				zap.L().Error("Error fetching default wallet address", zap.Error(err))
				return
			}

			// Create and render the wallet list table
			createWalletLsTable(wallets, defaultAddress)
		},
	}
}

// createWalletLsTable initializes a table writer with appropriate formatting for wallet display.
func createWalletLsTable(wallets []string, defaultAddress string) {
	table := tablewriter.NewWriter(os.Stdout)

	// Table header setup
	table.SetHeader([]string{"Address", "Default"})
	table.SetAutoFormatHeaders(false)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER})

	// Iterate over wallets and get their details (e.g., address and default status)
	for _, wallet := range wallets {
		defaultStatus := ""
		if wallet == defaultAddress {
			defaultStatus = "X"
		}

		// Append wallet information to the table
		table.Append([]string{wallet, defaultStatus})
	}

	// Render the table to stdout
	table.Render()
}
