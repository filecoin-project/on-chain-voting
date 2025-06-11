package storage

import (
	"context"
	"fil-vote/service"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

func ClaimMinerActorIdsCmd(client *service.RPCClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-miner-actor-ids [ACTOR_ID_1] [ACTOR_ID_2] ...",
		Short: "Claim that the current wallet address is the miner address for the supplied actor IDs",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			from, err := retrieveParameters(cmd, client)
			if err != nil {
				return
			}

			actorIds := make([]uint64, len(args))
			for i, arg := range args {
				var actorId uint64
				_, err := fmt.Sscanf(arg, "%d", &actorId)
				if err != nil {
					zap.L().Error("Invalid actorIds", zap.Error(err))
					return
				}
				actorIds[i] = actorId
			}

			messageHash, err := service.UpdateMinerIds(client, from, actorIds)
			if err != nil {
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Message Hash"})
			table.SetRowLine(true)

			table.Append([]string{messageHash})

			table.Render()
		},
	}

	cmd.Flags().String("from", "", "Wallet address to claim the miner actor IDs(optional)")

	return cmd
}

func retrieveParameters(cmd *cobra.Command, client *service.RPCClient) (string, error) {
	from, _ := cmd.Flags().GetString("from")

	if from == "" {
		var err error
		from, err = client.WalletDefaultAddress(context.Background())
		if err != nil {
			return "", err
		}
	}
	return from, nil
}
