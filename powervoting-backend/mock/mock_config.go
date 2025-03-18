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

package mock

import "powervoting-server/config"

func InifMockConfig() {
	config.Client = config.Config{
		ABIPath: config.ABIPath{
			PowerVotingAbi: "../abi/power-voting.json",
			OracleAbi:      "../abi/oracle.json",
		},
		Drand: config.Drand{
			Url: []string{
				"https://api2.drand.sh/",
				"https://api.drand.secureweb3.com:6875",
				"https://api.drand.sh/",
				"https://api3.drand.sh/",
				"https://drand.cloudflare.com/",
			},
			ChainHash: "52db9ba70e0cc0f6eaf7803dd07447a1f5477735fd3f661792ba94600c84e971",
		},
		Network: []config.Network{
			{
				ChainId:                         314159,
				Name:                            "FileCoin-Calibration",
				Rpc:                             "https://filecoin-calibration.chainup.net/rpc/v1",
				PowerVotingContract:             "0x4fe1B0D71FBFe97458D5c29D47928e1EA3b4466b",
				PowerVotingContractDeployHeight: 240000,
				OracleContract:                  "0x974e0AffA36Ef25ad3F99Edda6a0f9Cc09D354Ff",
			},
		},
		Snapshot: config.Snapshot{},
	}
}
