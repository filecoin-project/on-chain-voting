package proposal

import (
	"fil-vote/service"
	"github.com/spf13/cobra"
)

func NewProposalCmd(client *service.RPCClient) *cobra.Command {
	proposalCmd := &cobra.Command{
		Use:   "proposal",
		Short: "Manage proposal",
	}

	proposalCmd.AddCommand(ListProposalsCmd())
	proposalCmd.AddCommand(ApproveCmd(client))
	proposalCmd.AddCommand(RejectCmd(client))
	proposalCmd.AddCommand(ResultsCmd())

	return proposalCmd
}
