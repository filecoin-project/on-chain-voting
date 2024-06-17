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
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

// GetUcanFromIpfs retrieves a UCAN (User-Centric Access Network) from IPFS (InterPlanetary File System).
func GetUcanFromIpfs(ucanCid string) (string, error) {
	url := fmt.Sprintf("https://%s.ipfs.w3s.link/", ucanCid)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ucan := strings.Trim(string(body), "\"")
	if ucan[:5] == "https" {
		url := ucan
		parts := strings.Split(url, "/")
		owner := parts[4]
		repo := parts[5]
		sha := parts[8]

		blob, _, err := config.GithubClient.Git.GetBlob(context.Background(), owner, repo, sha)
		if err != nil {
			zap.L().Error("get blob error:", zap.Error(err))
			return "", err
		}

		content := blob.GetContent()

		oneLineString := strings.ReplaceAll(content, "\n", "")
		oneLineString = strings.ReplaceAll(oneLineString, " ", "")
		decodedData, err := base64.StdEncoding.DecodeString(oneLineString)
		if err != nil {
			zap.L().Error("base64 decoding error:", zap.Error(err))
			return "", err
		}

		_, payload, _, _, err := ucanSplit(string(decodedData))
		if err != nil {
			zap.L().Error("ucan eth split error", zap.Error(err))

			return "", err
		}
		if payload.Aud == owner {
			return string(decodedData), nil
		}
		return "", errors.New("the owner in the url is different from the aud in the carrier")
	}

	return ucan, nil
}
