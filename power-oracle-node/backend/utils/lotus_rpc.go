// Copyright (C) 2023-2024 StorSwift Inc.
// This file is part of the PowerVoting library.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"backend/utils/types"
	"context"
	filecoinAddress "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/ybbus/jsonrpc/v3"
)

// NewClient create lotus rcp client.
func NewClient(endpoint string) jsonrpc.RPCClient {
	return jsonrpc.NewClientWithOpts(endpoint, &jsonrpc.RPCClientOpts{})
}

// IDFromAddress retrieves the identifier associated with a Filecoin address.
func IDFromAddress(ctx context.Context, lotusRpcClient jsonrpc.RPCClient, address string) (string, error) {
	resp, err := lotusRpcClient.Call(ctx, "Filecoin.StateLookupID", address, types.TipSetKey{})
	if err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", resp.Error
	}

	return resp.Result.(string), nil
}

// WalletVerify verifies a signature against a specified address and data.
func WalletVerify(ctx context.Context, lotusRpcClient jsonrpc.RPCClient, address string, signature crypto.Signature, data []byte) (bool, error) {
	addressStr, err := filecoinAddress.NewFromString(address)
	if err != nil {
		return false, err
	}
	resp, err := lotusRpcClient.Call(ctx, "Filecoin.WalletVerify", addressStr, data, signature)
	if err != nil {
		return false, err
	}

	if resp.Error != nil {
		return false, resp.Error
	}

	getBool, err := resp.GetBool()
	if err != nil {
		return false, err
	}
	return getBool, err
}
