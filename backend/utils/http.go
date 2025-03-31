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
	"io"
	"net/http"

	"go.uber.org/zap"

	"powervoting-server/constant"
	"powervoting-server/model"
)

func GetGistInfoByGistId(gistId string) model.GistInfo {
	client := &http.Client{
		Timeout: constant.RequestTimeout,
	}

	url := constant.GistApiPrefix + gistId
	resp, err := client.Get(url)
	if err != nil {
		zap.L().Error("Failed to get gist info, http request failed", zap.Error(err))
		return model.GistInfo{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		zap.L().Error("Failed to get gist info", zap.String("status code", resp.Status), zap.String("body", string(body)), zap.Error(err))
		return model.GistInfo{}
	}

	var result model.GistInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		zap.L().Error("Failed to get gist info, json decode failed", zap.Error(err))
		return model.GistInfo{}
	}

	return result

}
