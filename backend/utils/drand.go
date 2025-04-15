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
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/drand/tlock"
	drandhttp "github.com/drand/tlock/networks/http"

	"go.uber.org/zap"

	"powervoting-server/config"
)

// It decrypts the encrypted data, unmarshals it into a structured format,
// and constructs a list of vote counts for each option.
// The function returns the decoded vote list or an error if the decoding fails.
func DecodeVoteResult(voteInfo string) (string, error) {
	var (
		decrypt []byte
		err     error
	)
	retry_times := 5

	for i := 0; i < retry_times; i++ {
		decrypt, err = Decrypt(voteInfo)
		if i == retry_times-1 && err != nil {
			zap.L().Error("decrypt error:", zap.Error(err))
			return "", err
		}
		if err != nil {
			zap.L().Warn(fmt.Sprintf("Decrypt failed: %v, retry times: %d\n", err, i))
			continue
		}
		break
	}
	var parsedResult [][]string
	err = json.Unmarshal(decrypt, &parsedResult)
	if err != nil {
		zap.L().Error("unmarshal error：", zap.Error(err))
		return "", err
	}

	if len(parsedResult) == 0 && len(parsedResult[0]) == 0 {
		return "", fmt.Errorf("no vote information")
	}

	return parsedResult[0][0], nil
}

// Decrypt decrypts the IPFS data using the T-lock encryption scheme.
// It replaces escape characters in the IPFS string, constructs a drand network,
// and decrypts the data using T-lock encryption.
// The function returns the decrypted data or an error if decryption fails.
func Decrypt(decStr string) ([]byte, error) {
	// Construct a network that can talk to a drand network. Example using the mainnet fastnet network.
	replace := strings.ReplaceAll(decStr, "\\n", "\n")
	replace2 := strings.ReplaceAll(replace, "\"", "")

	var network *drandhttp.Network
	var err error
	for _, url := range config.Client.Drand.Url {
		network, err = drandhttp.NewNetwork(url, config.Client.Drand.ChainHash)
		if err == nil {
			break
		}
	}
	reader := strings.NewReader(replace2)
	// Write the encrypted file data to this buffer.
	var cipherData bytes.Buffer
	// Encrypt the data for the given round.
	if err := tlock.New(network).Decrypt(&cipherData, reader); err != nil {
		zap.L().Error("Decrypt: %v", zap.Error(err))
		return nil, err
	}

	data := make([]byte, cipherData.Len())
	_, err = cipherData.Read(data)
	if err != nil {
		zap.L().Error("read: %v", zap.Error(err))
		return nil, err
	}
	return data, nil
}

