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

// Config represents the overall configuration structure.
type Config struct {
	Server   Server
	Nats     Nats
	Network  Network // Network configuration details.
	Github   GitHub    // Github configuration details.
	Redis    Redis     // Redis configuration details.
	Mysql    Mysql     // Mysql configuration details.
	W3Client W3Client  // W3Client configuration details.
	Rate     Rate      // Rate configuration details.
	DataPath DataPath    // Data path for storing files.
}

// Network  configuration for a blockchain network.
type Network struct {
	ChainId  int64    // Identifier for the network.
	Name     string   // Name of the network.
	QueryRpc []string // Query RPC endpoint for the network.
}

// GitHub represents the configuration for GitHub integration.
type GitHub struct {
	Token   []string // GitHub token for authentication.
	GraphQl string   // GraphQL endpoint for GitHub API.
}

// Server represents the server configuration.
type Server struct {
	Port   string // Port number for the server
	RpcUri string // RPC URI for the server
}

type Redis struct {
	URI      string // URI for the Redis server
	User     string // Username for accessing the Redis database
	Password string // Password for accessing the Redis database
	DB       int    // Database number
}

type Nats struct {
	URI string // URI for the NATS server
}

type Mysql struct {
	Url      string // URL of the MySQL database
	Username string // Username for accessing the MySQL database
	Password string // Password for accessing the MySQL database
}

type W3Client struct {
	Space      string // IPFS space
	Proof      string // IPFS proof
	PrivateKey string // Private key for the wallet
}

type Rate struct {
	GithubRequestLimit int64 // Limit for GitHub requests
}

type DataPath struct {
	DeveloperWeights string // Path to the developer weights file
}

type ChainPrefix map[int64]string

func (c *ChainPrefix) GetPrefix(netId int64) string {
	value, exist := (*c)[netId]
	if !exist {
		return ""
	}

	return value
}
