package wallet

import (
	"context"
	"fil-vote/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// AddCmd returns a Cobra command to import a wallet and display the status.
func AddCmd(client *service.RPCClient) *cobra.Command {
	return &cobra.Command{
		Use:   "add [walletType] [privateKey]", // Command name, with wallet type and private key as arguments
		Short: "Import a wallet",               // Short description of the command
		Args:  cobra.ExactArgs(2),              // Ensures exactly two arguments are provided
		RunE: func(cmd *cobra.Command, args []string) error { // Use RunE to return error for better error handling
			// Get the wallet type and private key from arguments
			walletType := args[0]
			privateKey := args[1]

			// Call the WalletImport function to import the wallet
			address, err := client.WalletImport(context.Background(), walletType, privateKey)
			if err != nil {
				// Enhanced logging with more context, including walletType and privateKey
				zap.L().Error("Failed to import wallet",
					zap.String("WalletType", walletType),
					zap.String("PrivateKey", privateKey),
					zap.Error(err))
				return err // Return the error to ensure proper handling
			}

			// Log successful wallet import
			zap.L().Info("Wallet imported successfully",
				zap.String("WalletType", walletType),
				zap.String("Address", address))

			return nil // Return nil to indicate successful completion
		},
	}
}
