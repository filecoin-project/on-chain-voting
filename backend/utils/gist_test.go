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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGistInfoByGistId(t *testing.T) {
	res, err := FetchGistInfoByGistId("c8a001be0c90e8c616e60100c1af54bf")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestParseGistContent(t *testing.T) {
	gist, err := FetchGistInfoByGistId("0d753f709c11735e7598ae6cf2657c60")
	assert.NoError(t, err)
	assert.NotEmpty(t, gist)
	voterInfo, err := ParseGistContent(gist.Files)
	assert.NoError(t, err)
	fmt.Printf("match res: %v\n", voterInfo)

	// sigBytes, _ := json.Marshal(res.SigObject)
	verifyRes, err := VerifySignature(
		"0x763D410594a24048537990dde6ca81c38CfF566a",
		voterInfo.Signature,
		[]byte(voterInfo.SigObjectStr),
	)

	assert.NoError(t, err)
	assert.True(t, verifyRes)

}
