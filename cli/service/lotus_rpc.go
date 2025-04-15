package service

import (
	"context"
	"fil-vote/config"
	"fil-vote/model"
	"fmt"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"
)

// RPCClient is a service struct that exposes a JSON-RPC client
type RPCClient struct {
	client jsonrpc.RPCClient
}

// NewRPCClient creates a new RPCClient instance
// It initializes the JSON-RPC client with the configured RPC URL and custom headers.
func NewRPCClient() *RPCClient {
	// Create a new JSON-RPC client with the configured RPC URL and custom headers
	client := jsonrpc.NewClientWithOpts(config.Client.Network.RPC, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Bearer " + config.Client.Network.Token, // Add Bearer token for authorization
		},
	})
	return &RPCClient{client: client}
}

// ListWallets retrieves the list of wallets connected to the Lotus node
// It calls the "Filecoin.WalletList" method on the RPC client to fetch the wallets.
func (rpc *RPCClient) ListWallets(ctx context.Context) ([]string, error) {
	var wallets []string
	// Call the "Filecoin.WalletList" method to get the list of wallets
	err := rpc.client.CallFor(ctx, &wallets, "Filecoin.WalletList")
	if err != nil {
		zap.L().Error("error retrieving wallet list", zap.Error(err))
		return nil, err // Return the error to be handled by caller
	}
	return wallets, nil
}

func (rpc *RPCClient) WalletDefaultAddress(ctx context.Context) (string, error) {
	var walletDefaultAddress string
	// Call the Filecoin API method to retrieve the default wallet address
	err := rpc.client.CallFor(ctx, &walletDefaultAddress, "Filecoin.WalletDefaultAddress")
	if err != nil {
		// Enhanced logging with context for debugging
		zap.L().Error("Failed to retrieve default wallet address", zap.Error(err))
		return "", fmt.Errorf("failed to retrieve default wallet address: %w", err)
	}
	return walletDefaultAddress, nil
}

func (rpc *RPCClient) WalletSetDefault(ctx context.Context, walletAddress string) error {
	// Call the Filecoin API method to set the default wallet address
	_, err := rpc.client.Call(ctx, "Filecoin.WalletSetDefault", walletAddress)
	if err != nil {
		// Enhanced logging with context for debugging
		zap.L().Error("Failed to set default wallet address", zap.String("walletAddress", walletAddress), zap.Error(err))
		return fmt.Errorf("failed to set wallet address as default: %w", err)
	}
	return nil
}

func (rpc *RPCClient) WalletImport(ctx context.Context, walletType, privateKey string) (string, error) {
	var walletAddress string

	// Call the method to import the wallet
	err := rpc.client.CallFor(ctx, &walletAddress, "Filecoin.WalletImport",
		[]interface{}{
			map[string]interface{}{
				"Type":       walletType, // Wallet type (e.g., "ecdsa", "secp256k1")
				"PrivateKey": privateKey, // The private key to import the wallet
			},
		},
	)

	if err != nil {
		// Return a wrapped error with more detailed context
		return "", err
	}

	// Return the wallet address upon success
	return walletAddress, nil
}

// GetNonce retrieves the nonce for a given address
// The nonce is used to prevent replay attacks by ensuring the uniqueness of each transaction.
func (rpc *RPCClient) GetNonce(ctx context.Context, address string) (int, error) {
	var nonce int
	// Call the "Filecoin.MpoolGetNonce" method to fetch the nonce for the provided address
	if err := rpc.client.CallFor(ctx, &nonce, "Filecoin.MpoolGetNonce", address); err != nil {
		zap.L().Error("failed to get nonce for address", zap.String("address", address), zap.Error(err))
		return 0, err // Return the error to be handled by caller
	}
	return nonce, nil
}

// EstimateGas estimates the gas usage for a message
// It calls the "Filecoin.GasEstimateMessageGas" method with the message and returns the estimated gas message.
func (rpc *RPCClient) EstimateGas(ctx context.Context, msg model.Message) (types.Message, error) {
	params := []interface{}{
		msg, // The message to be estimated
		map[string]interface{}{
			"MaxFee": "0", // Maximum fee (set to 0 in this example)
		},
		types.EmptyTSK, // A dummy value for the CID
	}

	var message types.Message
	// Call the "Filecoin.GasEstimateMessageGas" method to estimate the gas for the message
	if err := rpc.client.CallFor(ctx, &message, "Filecoin.GasEstimateMessageGas", params); err != nil {
		zap.L().Error("failed to estimate gas", zap.Error(err))
		return types.Message{}, err // Return the error to be handled by caller
	}
	return message, nil
}

// SendMessage sends a signed message to the Lotus node's message pool
// It signs the message and then pushes it to the network.
func (rpc *RPCClient) SendMessage(ctx context.Context, msg types.Message) (string, error) {
	var signedMsg types.SignedMessage
	// Call the "Filecoin.WalletSignMessage" method to sign the message
	if err := rpc.client.CallFor(ctx, &signedMsg, "Filecoin.WalletSignMessage", []interface{}{
		msg.From, // The address signing the message
		msg,      // The message to be signed
	}); err != nil {
		zap.L().Error("failed to sign message", zap.Error(err))
		return "", err // Return the error to be handled by caller
	}

	// Call the "Filecoin.MpoolPush" method to push the signed message to the message pool
	resp, err := rpc.client.Call(ctx, "Filecoin.MpoolPush", []interface{}{signedMsg})
	if err != nil {
		zap.L().Error("failed to push message", zap.Error(err))
		return "", err // Return the error to be handled by caller
	}

	// Assuming the response contains the CID of the message (you may need to adjust this based on actual response structure)
	messageCID := resp.Result.(map[string]interface{})["/"].(string)
	// Return the message CID (assuming this is the response format)
	return messageCID, nil
}
