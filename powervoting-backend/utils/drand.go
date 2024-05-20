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
	"io"
	"net/http"
	"powervoting-server/config"
	"powervoting-server/model"
	"strconv"
	"strings"

	"github.com/drand/tlock"
	drandhttp "github.com/drand/tlock/networks/http"
	"go.uber.org/zap"
)

// DecodeVoteList decodes the vote information retrieved from IPFS.
// It decrypts the encrypted data, unmarshals it into a structured format,
// and constructs a list of vote counts for each option.
// The function returns the decoded vote list or an error if the decoding fails.
func DecodeVoteList(voteInfo model.Vote) ([]model.Vote4Counting, error) {
	var voteList []model.Vote4Counting
	ipfs, err := GetIpfs(voteInfo.VoteInfo)
	if err != nil {
		zap.L().Error("get vote info from IPFS error: ", zap.Error(err))
		return voteList, err
	}
	retry_times := 5
	var decrypt []byte
	for i := 0; i < retry_times; i++ {
		decrypt, err = Decrypt(ipfs)
		if i == retry_times-1 && err != nil {
			zap.L().Error("decrypt error:", zap.Error(err))
			return voteList, err
		}
		if err != nil {
			zap.L().Warn(fmt.Sprintf("Decrypt failed: %v, retry times: %d\n", err, i))
			continue
		}
		break
	}
	var mapData [][]string
	err = json.Unmarshal(decrypt, &mapData)
	if err != nil {
		zap.L().Error("unmarshal error：", zap.Error(err))
		return voteList, err
	}

	for _, value := range mapData {
		optionId, err := strconv.Atoi(value[0])
		if err != nil {
			zap.L().Error("string to int error: ", zap.Error(err))
			return voteList, err
		}
		optionValue, err := strconv.Atoi(value[1])
		if err != nil {
			zap.L().Error("string to int error: ", zap.Error(err))
			return voteList, err
		}
		vote := model.Vote4Counting{
			OptionId: int64(optionId),
			Votes:    int64(optionValue),
			Address:  voteInfo.Address,
		}
		voteList = append(voteList, vote)
	}
	return voteList, nil
}

// GetIpfs retrieves data from IPFS using the provided votInfo.
// It constructs the IPFS URL and sends an HTTP GET request to fetch the data.
// The function returns the retrieved data or an error if the operation fails.
func GetIpfs(votInfo string) (string, error) {
	url := fmt.Sprintf("https://%s.ipfs.w3s.link/", votInfo)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err

	}

	return string(body), err
}

// Decrypt decrypts the IPFS data using the T-lock encryption scheme.
// It replaces escape characters in the IPFS string, constructs a drand network,
// and decrypts the data using T-lock encryption.
// The function returns the decrypted data or an error if decryption fails.
func Decrypt(ipfs string) ([]byte, error) {
	// Construct a network that can talk to a drand network. Example using the mainnet fastnet network.
	replace := strings.ReplaceAll(ipfs, "\\n", "\n")
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

// GetOptions retrieves the options of a proposal from IPFS.
// It takes the CID of the proposal as input and returns the options or an error if retrieval fails.
func GetOptions(cid string) ([]string, error) {
	ipfs, err := GetIpfs(cid)
	if err != nil {
		zap.L().Error("get ipfs error: ", zap.Error(err))
		return nil, err
	}
	var proposalDetail model.ProposalDetail
	err = json.Unmarshal([]byte(ipfs), &proposalDetail)
	if err != nil {
		zap.L().Error("json unmarshal error: ", zap.Error(err))
		return nil, err
	}
	return proposalDetail.Option, nil
}
