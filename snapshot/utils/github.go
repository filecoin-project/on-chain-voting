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
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"power-snapshot/constant"
)

func FetchGithubDeveloper(url, token string, queryParams map[string]string, v any) (string, error) {
	client := resty.New().
		SetTimeout(constant.TimeoutWith15s).
		SetRetryCount(constant.RetryCount).
		SetRetryWaitTime(2 * time.Second).
		SetAuthToken(token)

	resp, err := client.R().
		SetHeader("Accept", "application/vnd.github.v3+json").
		SetResult(v).
		SetQueryParams(queryParams).
		Get(url)

	if err != nil {
		zap.L().Error(
			"failed to request github",
			zap.String("request url", url),
			zap.Any("query params", queryParams),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to request github: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.RawResponse.Header.Get("Link"), nil
	case http.StatusNotFound:
		zap.L().Error("failed to request github: not found")
		return "", errors.New("failed to request github: not found")
	default:
		zap.L().Error("failed to request github", zap.String("response", resp.String()))
		return "", fmt.Errorf("failed to request github: %s", resp.String())
	}

}

// getNextPageURL extracts the next page URL from the Link header.
func GetNextPageURL(linkHeader string) string {
	if linkHeader == "" {
		return ""
	}
	links := parseLinkHeader(linkHeader)
	return links["next"]
}

// parseLinkHeader parses the Link header into a map of link relations and URLs.
func parseLinkHeader(linkHeader string) map[string]string {
	links := make(map[string]string)
	if linkHeader != "" {
		entries := splitByComma(linkHeader)
		for _, entry := range entries {
			parts := splitBySemicolon(entry)
			if len(parts) < 2 {
				continue
			}
			url := extractURL(parts[0])
			rel := extractRel(parts[1])
			if url != "" && rel != "" {
				links[rel] = url
			}
		}
	}
	return links
}

// splitByComma splits a string by commas.
func splitByComma(s string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}

// splitBySemicolon splits a string by semicolons.
func splitBySemicolon(s string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ';' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}

// extractURL extracts a URL enclosed in angle brackets from a string.
func extractURL(s string) string {
	start := 0
	for ; start < len(s) && (s[start] == ' ' || s[start] == '\t'); start++ {
	}
	if start >= len(s) || s[start] != '<' {
		return ""
	}
	start++
	end := start
	for ; end < len(s) && s[end] != '>'; end++ {
	}
	if end >= len(s) || s[end] != '>' {
		return ""
	}
	return s[start:end]
}

// extractRel extracts a relation from a string formatted as 'rel="value"'.
func extractRel(s string) string {
	start := 0
	for ; start < len(s) && (s[start] == ' ' || s[start] == '\t'); start++ {
	}
	if start >= len(s) || s[start] != 'r' || start+3 >= len(s) || s[start+1] != 'e' || s[start+2] != 'l' || s[start+3] != '=' {
		return ""
	}
	start += 4
	if start >= len(s) {
		return ""
	}
	if s[start] == '"' {
		start++
		end := start
		for ; end < len(s) && s[end] != '"'; end++ {
		}
		if end >= len(s) || s[end] != '"' {
			return ""
		}
		return s[start:end]
	}
	end := start
	for ; end < len(s) && s[end] != ',' && s[end] != ';'; end++ {
	}
	return s[start:end]
}
