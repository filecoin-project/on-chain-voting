package power

import (
	"fil-vote/service"
	"github.com/spf13/cobra"
)

func NewPowerCmd(client *service.RPCClient) *cobra.Command {
	powerCmd := &cobra.Command{
		Use:   "power",
		Short: "Manager power",
	}

	powerCmd.AddCommand(LsCmd(client))

	return powerCmd
}
