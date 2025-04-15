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
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-resty/resty/v2"
	"github.com/storyicon/sigverify"
	"go.uber.org/zap"

	"powervoting-server/constant"
	"powervoting-server/model"
)

var (
	objRex       = regexp.MustCompile(`(?i)signing this object\s*({.*?})\s*`)
	signatureRex = regexp.MustCompile(`(?i)signature:\s*(0x[\da-fA-F]+)`)
)

func FetchGistInfoByGistId(gistId string) (model.Gist, error) {
	client := resty.New().
		SetTimeout(constant.RequestTimeout).
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second)

	url := constant.GistApiPrefix + gistId
	var result model.Gist

	resp, err := client.R().
		SetHeader("Accept", "application/vnd.github.v3+json").
		SetResult(&result).
		Get(url)

	if err != nil {
		zap.L().Error("HTTP request failed",
			zap.Error(err),
			zap.String("url", url))
		return model.Gist{}, err
	}

	switch resp.StatusCode() {
	case http.StatusNotFound:
		zap.L().Error("Gist not found",
			zap.String("gist_id", gistId))
		return model.Gist{}, constant.ErrGistNotFound

	case http.StatusOK:
		return result, nil

	default:
		zap.L().Error("Unexpected status code",
			zap.Int("status_code", resp.StatusCode()),
			zap.String("response_body", string(resp.Body())),
			zap.String("gist_id", gistId))
		return model.Gist{},
			fmt.Errorf("API returned %d: %s", resp.StatusCode(), resp.Body())
	}
}

func ParseGistContent(files map[string]model.GistFiles) (model.GistVoterInfo, error) {
	var gistFile model.GistFiles
	for _, file := range files {
		gistFile = file
	}

	matchFunc := func(re *regexp.Regexp) string {
		matches := re.FindStringSubmatch(gistFile.Content)
		if len(matches) != 2 {
			return ""
		}

		return matches[1]
	}

	var sigObj model.SigObject

	objStr := matchFunc(objRex)
	if err := json.Unmarshal([]byte(objStr), &sigObj); err != nil {
		return model.GistVoterInfo{}, fmt.Errorf("invalid signature object: %s", gistFile.Content)
	}

	return model.GistVoterInfo{
		SigObjectStr: objStr,
		SigObject:    sigObj,
		Signature:    matchFunc(signatureRex),
	}, nil
}

func VerifyAuthorizeAllow(address, githubName string, gist model.Gist) bool {
	sigObj, err := ParseGistContent(gist.Files)
	if err != nil {
		zap.L().Error("ParseGistContent error",  zap.Error(err))
		return false
	}

	isValid, err := VerifySignature(address, sigObj.Signature, []byte(sigObj.SigObjectStr))
	if !isValid || err != nil {
		zap.L().Warn(
			"VerifySignature error",
			zap.String("address", address),
			zap.String("signature", sigObj.Signature),
			zap.String("sig msg", sigObj.SigObjectStr),
			zap.Error(err),
		)
		return false
	}

	if gist.Owner.Login != sigObj.SigObject.GitHubName || githubName != sigObj.SigObject.GitHubName {
		zap.L().Warn(
			"Invalid Github name",
			zap.String("expected github name", githubName),
			zap.String("actual github name", sigObj.SigObject.GitHubName),
		)
		return false
	}

	if sigObj.SigObject.WalletAddress != address {
		zap.L().Warn(
			"Invalid address",
			zap.String("expected address", address),
			zap.String("actual address", sigObj.SigObject.WalletAddress),
		)
		return false
	}

	return true
}

func VerifySignature(address string, signature string, msgData []byte) (bool, error) {
	isValid, err := sigverify.VerifyEllipticCurveHexSignatureEx(
		common.HexToAddress(address),
		msgData,
		signature,
	)

	if err != nil {
		return isValid, err
	}

	return isValid, nil
}
