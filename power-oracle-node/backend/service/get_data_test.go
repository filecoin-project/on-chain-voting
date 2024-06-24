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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUcanIpfs(t *testing.T) {
	var err error
	config.InitConfig("../")
	contract.GoEthClient, err = contract.GetClient(314159)
	assert.Nil(t, err)

	rsp, err := GetUcanFromIpfs("bafkreidnrzkmshm36e4uzj5mxl6gmwv2b6uhkkwvvpxo3anjcxepbskn3q")
	assert.Nil(t, err)

	expected := `eyJhbGciOiJlY2RzYSIsInR5cGUiOiJKV1QiLCJ2ZXJzaW9uIjoiMC4wLjEifQ.eyJpc3MiOiIweDc2M0Q0MTA1OTRhMjQwNDg1Mzc5OTBkZGU2Y2E4MWMzOENmRjU2NmEiLCJhdWQiOiJ0MXdhNGd2eWVlazRvaDV6ZzM3NW9vNmx3aGNtZHdneHdzNXJnc2x5eSIsInByZiI6ImV5SmhiR2NpT2lKelpXTndNalUyYXpFaUxDSjBlWEJsSWpvaVNsZFVJaXdpZG1WeWMybHZiaUk2SWpBdU1TSjkuZXlKcGMzTWlPaUowTVhkaE5HZDJlV1ZsYXpSdmFEVjZaek0zTlc5dk5teDNhR050WkhkbmVIZHpOWEpuYzJ4NWVTSXNJbUYxWkNJNklqQjROell6UkRReE1EVTVOR0V5TkRBME9EVXpOems1TUdSa1pUWmpZVGd4WXpNNFEyWkdOVFkyWVNJc0ltRmpkQ0k2SW1Ga1pDSXNJbkJ5WmlJNklpSjkuMWdpSE1mX01HQ2drWFJVRE1Ed2VpeTRWeHpnTzVTeV95a2haYVgteUFZQlpndElQdnIyRUFzYkQ1bTNlYjlWOGJRSEJ0SFVXOXowMVhkRkJzY05tVWdFIiwiYWN0IjoiYWRkIn0.MHgzY2FkNTg0YmE5YzEwODNiMGQwZjBiY2VhOGUwM2VjYjM5NzM3NTlmMTkwNjllYjI0OGUxZmZiNTEwMjIzZjBjMmYyYWYxODYwNzZkM2I0NTkyMzhlMTkxNzE2NWIyYjI3NTEwOTQ2N2QzZWY0ZTcyNTI3MTk1MzMzNmRjN2E4ZTFj`

	assert.NotEmpty(t, rsp)
	assert.Equal(t, rsp, expected)
}
