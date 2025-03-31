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

package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

// GoEthClient represents the structure for interacting with the Ethereum client.
type GoEthClient struct {
	ChainId           int64               // Unique identifier for the client
	Name              string              // Name of the client
	IdPrefix          string              // id prefix
	QueryClient       []*ethclient.Client // Ethereum client instance list
	QueryRpc          []string            // RPC endpoint for the query client
	Amount            *big.Int            // Amount for transactions.
}
