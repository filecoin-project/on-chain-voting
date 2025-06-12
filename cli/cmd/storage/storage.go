package storage

import (
	"fil-vote/service"
	"github.com/spf13/cobra"
)

func NewStorageCmd(client *service.RPCClient) *cobra.Command {
	storageCmd := &cobra.Command{
		Use:   "storage-provider",
		Short: "Manage Storage Provider",
	}

	storageCmd.AddCommand(ClaimMinerActorIdsCmd(client))

	return storageCmd
}
