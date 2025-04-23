package repo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"power-snapshot/config"
	"power-snapshot/internal/data"
	"power-snapshot/internal/repo"
)

func getLotus(t *testing.T) *repo.LotusRPCRepo {
	config.InitConfig("../../")
	config.InitLogger()
	redis, err := data.NewRedisClient()
	assert.NoError(t, err)
	return repo.NewLotusRPCRepo(redis)
}

func TestGetWalletBalanceByHeight(t *testing.T) {
	res, err := getLotus(t).GetWalletBalanceByHeight(
		context.Background(),
		"t0161747",
		314159,
		2599578,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}
