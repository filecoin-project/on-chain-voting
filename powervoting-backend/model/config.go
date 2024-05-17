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

package model

// Config represents the configuration structure for the PowerVoting application.
type Config struct {
	Server  Server    // Server configuration
	Mysql   Mysql     // MySQL database configuration
	Drand   Drand     // Drand network configuration
	Network []Network // List of network configurations
}

// Server represents the server configuration.
type Server struct {
	Port string // Port number for the server
}

// Mysql represents the MySQL database configuration.
type Mysql struct {
	Url      string // URL of the MySQL database
	Username string // Username for accessing the MySQL database
	Password string // Password for accessing the MySQL database
}

// Drand represents the Drand network configuration.
type Drand struct {
	Url       []string // List of URLs for the Drand network
	ChainHash string   // Chain hash for the Drand network
}

// Network represents the configuration for a specific network.
type Network struct {
	Id                  int64  // Unique identifier for the network
	Name                string // Name of the network
	Rpc                 string // RPC endpoint for the network
	PowerVotingAbi      string // ABI (Application Binary Interface) for PowerVoting contract
	OracleAbi           string // ABI for Oracle contract
	PowerVotingContract string // Contract address for PowerVoting
	OracleContract      string // Contract address for Oracle
}
