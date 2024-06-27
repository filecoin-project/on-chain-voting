package service

import (
	"context"
	"log"
	"power-snapshot/config"
	"power-snapshot/utils"
	"testing"
)

func TestGetPower(t *testing.T) {

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

	// Get the client for the specified network
	client, err := manager.GetClient(314159)
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
		return
	}

	// Create a new Lotus RPC client
	lotusRpcClient := utils.NewClient(client.QueryRpc[0])
	GetPower(context.Background(), lotusRpcClient, "t017592", "t03751", 1490767)
}
