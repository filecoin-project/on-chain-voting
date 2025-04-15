package wallet

import (
	"fil-vote/service"
	"github.com/spf13/cobra"
)

func NewWalletCmd(client *service.RPCClient) *cobra.Command {
	walletsCmd := &cobra.Command{
		Use:   "wallet",
		Short: "Manage wallet",
	}
	walletsCmd.AddCommand(LsCmd(client))
	walletsCmd.AddCommand(UseCmd(client))
	walletsCmd.AddCommand(AddCmd(client))

	return walletsCmd
}
