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

)

// func TestParseGistContent(t *testing.T) {
// 	config.Client.Github.Token = []string{
// 		"",
// 	}
// 	config.Client.Network.Rpc = "http://192.168.11.139:1235/rpc/v1"
// 	gist, err := FetchGistInfoByGistId("9e55f044e364baae82eef6d038a1b5df")
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, gist)
// 	voterInfo, err := ParseGistContent(gist.Files)
// 	assert.NoError(t, err)
// 	fmt.Printf("match res: %v\n", voterInfo)

// 	// sigBytes, _ := json.Marshal(res.SigObject)
// 	verifyRes, err := VerifySignature(
// 		"t1bh2fekhhi3c4rcynxvah6hqei2s4geylxizvzfa",
// 		voterInfo.Signature,
// 		[]byte(voterInfo.SigObjectStr),
// 	)

// 	assert.NoError(t, err)
// 	assert.True(t, verifyRes)
// }

// func TestVerifySignature(t *testing.T) {
// 	// sigBytes, _ := json.Marshal(res.SigObject)
// 	msg := `{"walletAddress":"","githubName":"1","timestamp":1744861584}`
// 	sig := ""
// 	verifyRes, err := VerifySignature(
// 		"",
// 		sig,
// 		[]byte(msg),
// 	)

// 	assert.NoError(t, err)
// 	assert.True(t, verifyRes)
// }

// func TestVerifyFilecoinAddrSignature(t *testing.T) {
// 	config.GetDefaultConfig()
// 	res, err := VerifyFilecoinAddrSignature(
// 		`0x4fda4174D5D07C906395bfB77806287cc65Fd129`,
// 		"0xa6824f7b6ec4a308476c3414cad8495ac283359ff9eb7235637063dcb0213ce7169ab205991408980a9452ecc909acaf5d5c103ad5f9fc0a99682142ba84cf8a1b",
// 		[]byte(`{"walletAddress":"0x4fda4174D5D07C906395bfB77806287cc65Fd129","githubName":"liuzeming1","timestamp":1745327875}`),
// 	)

// 	assert.NoError(t, err)

// 	assert.True(t, res)
// }

// func TestVerifySignature(t *testing.T) {
// 	res, err := VerifySignature(
// 		"0x4fda4174D5D07C906395bfB77806287cc65Fd129",
// 		"0xa6824f7b6ec4a308476c3414cad8495ac283359ff9eb7235637063dcb0213ce7169ab205991408980a9452ecc909acaf5d5c103ad5f9fc0a99682142ba84cf8a1b",
// 		[]byte(`{"walletAddress":"0x4fda4174D5D07C906395bfB77806287cc65Fd129","githubName":"liuzeming1","timestamp":1745327875}`),
// 	)

// 	assert.NoError(t, err)
// 	assert.True(t, res)
// }
