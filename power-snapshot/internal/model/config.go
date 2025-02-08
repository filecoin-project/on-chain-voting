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
	Server  Server
	Nats    Nats
	Network []Network // Network configuration details.
	Github  GitHub    // Github configuration details.
	Redis   Redis     // Redis configuration details.
}

// Network  configuration for a blockchain network.
type Network struct {
	Id                int64    // Identifier for the network.
	Name              string   // Name of the network.
	IdPrefix          string   // Address prefix.
	QueryRpc          []string // Query RPC endpoint for the network.
	ContractRpc       string   // Contract RPC endpoint for the network.
	OracleAbi         string   // Path to the ABI file.
	OracleContract    string   // Address of the smart contract.
	OracleStartHeight int64    // Start height
}

// GitHub represents the configuration for GitHub integration.
type GitHub struct {
	Token   []string // GitHub token for authentication.
	GraphQl string   // GraphQL endpoint for GitHub API.
}

// Server represents the server configuration.
type Server struct {
	Port string // Port number for the server
}

type Redis struct {
	URI      string
	User     string
	Password string
	DB       int
}

type Nats struct {
	URI string
}
