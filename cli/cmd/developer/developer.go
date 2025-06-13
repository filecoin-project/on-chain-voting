package developer

import (
	"fil-vote/service"
	"github.com/spf13/cobra"
)

func NewDeveloperCmd(client *service.RPCClient) *cobra.Command {
	developerCmd := &cobra.Command{
		Use:   "developer",
		Short: "Manager developer",
	}

	developerCmd.AddCommand(GenerateCmd(client))
	developerCmd.AddCommand(UpdateGistIdCmd(client))

	return developerCmd
}
