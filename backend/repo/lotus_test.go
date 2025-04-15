package repo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"powervoting-server/config"
	"powervoting-server/repo"
)

func TestGetActorIdByAddress(t *testing.T) {
	config.InitConfig("../")
	config.InitLogger()

	lotusRepo := repo.NewLotusRPCRepo()

	res, err := lotusRepo.GetActorIdByAddress(context.Background(), "0x763D410594a24048537990dde6ca81c38CfF566a")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestFilecoinAddressToID(t *testing.T) {
	config.InitConfig("../")
	config.InitLogger()

	lotusRepo := repo.NewLotusRPCRepo()
	res, err := lotusRepo.GetValidMinerIds(context.Background(), "", []uint64{})
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

