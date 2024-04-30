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

package service

import (
	"backend/config"
	"backend/contract"
	"backend/utils"
	"fmt"
	"testing"
)

const (
	Secp256k1Ucan = "eyJhbGciOiJlY2RzYSIsInR5cGUiOiJKV1QiLCJ2ZXJzaW9uIjoiMC4wLjEifQ.eyJpc3MiOiIweDc2M0Q0MTA1OTRhMjQwNDg1Mzc5OTBkZGU2Y2E4MWMzOENmRjU2NmEiLCJhdWQiOiJ0MXdhNGd2eWVlazRvaDV6ZzM3NW9vNmx3aGNtZHdneHdzNXJnc2x5eSIsInByZiI6ImV5SmhiR2NpT2lKelpXTndNalUyYXpFaUxDSjBlWEJsSWpvaVNsZFVJaXdpZG1WeWMybHZiaUk2SWpBdU1DNHhJbjAuZXlKcGMzTWlPaUowTVhkaE5HZDJlV1ZsYXpSdmFEVjZaek0zTlc5dk5teDNhR050WkhkbmVIZHpOWEpuYzJ4NWVTSXNJbUYxWkNJNklqQjROell6UkRReE1EVTVOR0V5TkRBME9EVXpOems1TUdSa1pUWmpZVGd4WXpNNFEyWkdOVFkyWVNJc0ltRmpkQ0k2SW1Ga1pDSXNJbkJ5WmlJNklpSjkubWpCY3RycEFlSVVpNUQ3S0hwbjBMalRfU2h4UVZmR0dNRGMyYlJlNmdHWldnQjgzbjV2VEREYTlPX1ZaQTRhUlVpaWJOdHRRWkEyUFB5ak1qdkxzTndFIiwiYWN0IjoiYWRkIn0.MHg5NzUwODM4OGYwODhkMDk2MzQ4YjAxOGZkYTlmYWU0NjE0YWQ2ZWVjOWNkMDhkMGZlZDFhZDQwM2U0OTY4Mzk3NDIwMmU3NzI5OWFjOTg4YjNlMGRiNTRhMzM0ZWZjNmI0NjVjZGE3Y2MxZDBlN2FmY2FkZDkzYzIyYTViNzE1OTFj"
	BlsUcan       = "eyJhbGciOiJlY2RzYSIsInR5cGUiOiJKV1QiLCJ2ZXJzaW9uIjoiMC4wLjEifQ.eyJpc3MiOiIweDc2M0Q0MTA1OTRhMjQwNDg1Mzc5OTBkZGU2Y2E4MWMzOENmRjU2NmEiLCJhdWQiOiJ0M3c2bGpmd3J1eGlkNmZ2ZWlyamFtNjJrb3d6M3puMnVsNWY1a2Ntamg2d2RmNzc3Z2U0Z2o1dmRlcDNwcGpieXM1YWF4eWJqN2Rmd3M3YmdyYWt1YSIsInByZiI6ImV5SmhiR2NpT2lKaWJITWlMQ0owZVhCbElqb2lTbGRVSWl3aWRtVnljMmx2YmlJNklqQXVNQzR4SW4wLmV5SnBjM01pT2lKME0zYzJiR3BtZDNKMWVHbGtObVoyWldseWFtRnROakpyYjNkNk0zcHVNblZzTldZMWEyTnRhbWcyZDJSbU56YzNaMlUwWjJvMWRtUmxjRE53Y0dwaWVYTTFZV0Y0ZVdKcU4yUm1kM00zWW1keVlXdDFZU0lzSW1GMVpDSTZJakI0TnpZelJEUXhNRFU1TkdFeU5EQTBPRFV6TnprNU1HUmtaVFpqWVRneFl6TTRRMlpHTlRZMllTSXNJbUZqZENJNkltRmtaQ0lzSW5CeVppSTZJaUo5LmwybGdTeUcyaHVTVTdCMEFYUXRxTGFwSVh6THh5WFpaVXNwSkFMdlJXMzBRRHRFcF9WTHU5LU9rTW9fU25pN3lEbmpTVzVQRjliNnVBVFFWYjFiMnE3MUxYZXdtOGNvdU95U1JLYlZFNExuOFZlVVVVcFNtSEdXbVlPazFoTl9oIiwiYWN0IjoiYWRkIn0.MHhhZDYxY2JjZTI5NDI5NDU2NTVmZDBhYWNhNTcyZGU1MzJiNTM1OTQ0NmEzMjJiNTE5Y2EwOGEwZWIzYzUwMGQwN2ViYzM5NTgzYTJkOTBjMTlmNTRlYmM3NzdjNWU2YzhmODA1NmUzZmQ4OGU4ZTgxYTJmMmMxNTU4ZjVhZjIyMTFi"
)

func TestVerifyUCAN(t *testing.T) {
	var err error
	config.InitLogger()
	config.InitConfig("../")
	contract.GoEthClient, err = contract.GetClient(314159)
	if err != nil {
		t.Error(err)
	}
	lotusRpcClient := utils.NewClient(contract.GoEthClient.Rpc)

	iss, aud, act, isGitHub, err := VerifyUCAN(BlsUcan, lotusRpcClient)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(iss, aud, act, isGitHub)
}
