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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"powervoting-server/config"
	"powervoting-server/constant"
)

func initConfig() {
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
		Network: config.Network{
			ChainId:              314159,
			Name:                 "FileCoin-Calibration",
			Rpc:                  "https://filecoin-calibration.chainup.net/rpc/v1",
			PowerVotingContract:  "0x4fe1B0D71FBFe97458D5c29D47928e1EA3b4466b",
			SyncEventStartHeight: 240000,
			OracleContract:       "0x974e0AffA36Ef25ad3F99Edda6a0f9Cc09D354Ff",
		},
		Snapshot: config.Snapshot{},
	}
}
func TestDecodeVoteResult(t *testing.T) {
	initConfig()
	decStr := `-----BEGIN AGE ENCRYPTED FILE-----
YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IHRsb2NrIDE2MzI1ODEyIDUyZGI5YmE3
MGUwY2MwZjZlYWY3ODAzZGQwNzQ0N2ExZjU0Nzc3MzVmZDNmNjYxNzkyYmE5NDYw
MGM4NGU5NzEKdFVxWng1VVpiSW1pRjk2ZmJSU2dXRVJrblRjeEFuWjdkblg5VHNx
djNaZkdzYlh1eVBUSUh4NHZabkJWWmFwVwpCa09MNGpZMUVtNkQ2cjdGK0Z6eEE3
Y0VBR3F3SCtabFMyenE5TjRtNEF4d0FoMkFBS1NhbkxXRVozRzBMSlViCnQ4ZEF1
MHhCbHhmTW1LQnJKRXN3NkFIV3h6VjNSNXlHdnlobk03SFY3b3cKLS0tIHFwbHJi
Qnhxd1JvZHVKaElYSVZOa3dMZkt0SHdIc2hMVzA0NHVmSmsvcTQKBVM4IBEuFQcP
l0YZKbPlAmnEhcp3EAwQ84BvVSibhTmIzq/MYdsHnTX/1O8=
-----END AGE ENCRYPTED FILE-----
`
	res, err := DecodeVoteResult(decStr)
	assert.NoError(t, err, "decode error：", err)
	assert.Equal(t, constant.VoteReject, res)
}

func TestDecrypt(t *testing.T) {
	initConfig()
	decStr := `-----BEGIN AGE ENCRYPTED FILE-----
YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IHRsb2NrIDE2MzI1ODEyIDUyZGI5YmE3
MGUwY2MwZjZlYWY3ODAzZGQwNzQ0N2ExZjU0Nzc3MzVmZDNmNjYxNzkyYmE5NDYw
MGM4NGU5NzEKcExwOUJkK2ZWbFpEMEx5OGZ3OVhCa09pWldRWnpHSGZtWGErcSsw
b3pzdEFON2VNTjU3T3RFbTRFQ2d2bElpMApBV2ZvU3A5cy9Ydjl0YkRrb3p1Tk5v
RjdaRUo0c2RtdWRWQ2hEOENQcm1FN1VBYUxtaXBLR29ScVFKdFFFbHUrCitDT2Vr
SENnYWE4eHNGNUxuVzdwWEhkd25sRzlsSmxXSy9wUVJIaEM2S2cKLS0tIHNUV2Yr
RVdEZlJQWm9yaS9ITnNqcGNuQmxNcUlob3VKYkErRUJDb243eUkK2rehvaQY2kad
04WmD3ZNvrbcQRZwWonF+Ww4UhpwUApwbCjpEx80W5jfIo+p
-----END AGE ENCRYPTED FILE-----
`

	decrypt, err := Decrypt(decStr)
	assert.NoError(t, err, "decrypt error：", err)

	var data [][]string
	err = json.Unmarshal(decrypt, &data)
	assert.NoError(t, err, "unmarshal error：", err)
	assert.Len(t, data, 1)
	assert.Len(t, data[0], 1)
	assert.Equal(t, constant.VoteApprove, data[0][0])
}
