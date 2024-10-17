package service

import (
	"context"
	"log"
	"power-snapshot/config"
	"power-snapshot/internal/data"
	"power-snapshot/internal/repo"
	"power-snapshot/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetAddressPower(t *testing.T) {
	ctx := context.Background()

	// Initialize the logger
	config.InitLogger()

	// Load the configuration from the specified path
	err := config.InitConfig("../../")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	// Initialize the client manager
	manager, err := utils.NewGoEthClientManager(config.Client.Network)
	if err != nil {
		log.Fatalf("Failed to initialize client manager: %v", err)
		return
	}

	redisClient, err := data.NewRedisClient()
	if err != nil {
		panic(err)
	}

	jetstreamClient, err := data.NewJetstreamClient()
	if err != nil {
		panic(err)
	}

	// init repo
	syncRepo, err := repo.NewSyncRepoImpl(manager.ListClientNetWork(), redisClient, jetstreamClient)
	if err != nil {
		panic(err)
	}
	queryRepo, err := repo.NewQueryRepoImpl(redisClient, manager)
	if err != nil {
		panic(err)
	}
	baseRepo, err := repo.NewBaseRepoImpl(manager, redisClient)
	if err != nil {
		panic(err)
	}

	// init service
	syncSrv := NewSyncService(baseRepo, syncRepo)

	// init service
	service := NewQueryService(baseRepo, queryRepo, syncSrv)

	res, err := service.GetAddressPower(ctx, 314159, "0x4fda4174D5D07C906395bfB77806287cc65Fd129", 60)
	assert.Nil(t, err)

	zap.L().Info("result", zap.Any("res", res))
}
