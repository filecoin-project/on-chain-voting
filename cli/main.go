package main

import (
	"fil-vote/cmd"
	"fil-vote/cmd/proposal"
	"fil-vote/cmd/wallet"
	"fil-vote/config"
	"fil-vote/service"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	// Initialize the logger
	config.InitLogger()

	// Load configuration from the config file
	if err := config.InitConfig("./"); err != nil {
		log.Fatalf("Error initializing configuration: %v", err)
	}

	// Initialize the RPC client
	client := service.NewRPCClient()

	// Initialize the root command
	rootCmd := cmd.RootCmd

	// Add the "wallet" command to the root command
	addCommandToRoot(rootCmd, wallet.NewWalletCmd(client))

	// Add the "proposal" command to the root command
	addCommandToRoot(rootCmd, proposal.NewProposalCmd(client))

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// addCommandToRoot simplifies adding commands to the root command, avoiding repetitive code
func addCommandToRoot(rootCmd *cobra.Command, cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
