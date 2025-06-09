package storage

import (
	"context"
	"fil-vote/service"
	"fmt"
	"github.com/spf13/cobra"
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
					return
				}
				actorIds[i] = actorId
			}

			_, err = service.UpdateMinerIds(client, from, actorIds)
			if err != nil {
				return
			}

			fmt.Println("Claim successfully processed for actor IDs:", actorIds)
		},
	}

	cmd.Flags().String("from", "", "Wallet address to claim the miner actor IDs.")

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
