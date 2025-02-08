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
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type GitHubTokenManager struct {
	tokens               []string
	graphQLTokenUsage    map[string]int
	graphQLTokenCapacity map[string]int
	coreTokenUsage       map[string]int
	coreTokenCapacity    map[string]int
	mu                   sync.Mutex
}

func NewGitHubTokenManager(tokens []string) *GitHubTokenManager {
	manager := &GitHubTokenManager{
		tokens:               tokens,
		graphQLTokenUsage:    make(map[string]int),
		graphQLTokenCapacity: make(map[string]int),
		coreTokenUsage:       make(map[string]int),
		coreTokenCapacity:    make(map[string]int),
	}

	for _, token := range tokens {
		manager.graphQLTokenCapacity[token], manager.coreTokenCapacity[token] = CheckRateLimitBeforeRequest(token)
		zap.L().Info("GitHubTokenManager",
			zap.Int("graphQLTokenCapacity", manager.graphQLTokenCapacity[token]),
			zap.Int("graphQLTokenCapacity", manager.coreTokenCapacity[token]))
	}

	return manager
}

func (m *GitHubTokenManager) GetAvailableToken() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var maxToken string
	maxCapacity := 0

	for token, capacity := range m.graphQLTokenCapacity {
		if m.graphQLTokenUsage[token] < capacity && capacity > maxCapacity {
			maxToken = token
			maxCapacity = capacity
		}
	}

	if maxToken != "" {
		m.graphQLTokenUsage[maxToken]++
	}
	return maxToken
}

func (m *GitHubTokenManager) DecreaseUsage(token string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.graphQLTokenUsage[token] > 0 {
		m.graphQLTokenUsage[token]--
	}
}

func (m *GitHubTokenManager) GetCoreAvailableToken() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var maxToken string
	maxCapacity := 0

	for token, capacity := range m.coreTokenCapacity {
		if m.coreTokenUsage[token] < capacity && capacity > maxCapacity {
			maxToken = token
			maxCapacity = capacity
		}
	}

	if maxToken != "" {
		m.coreTokenUsage[maxToken]++
	}
	return maxToken
}

func (m *GitHubTokenManager) CoreDecreaseUsage(token string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.coreTokenUsage[token] > 0 {
		m.coreTokenUsage[token]--
	}
}

type RateLimitResponse struct {
	Resources struct {
		Core struct {
			Limit     int `json:"limit"`
			Used      int `json:"used"`
			Remaining int `json:"remaining"`
			Reset     int `json:"reset"`
		} `json:"core"`
		Search struct {
			Limit     int `json:"limit"`
			Used      int `json:"used"`
			Remaining int `json:"remaining"`
			Reset     int `json:"reset"`
		} `json:"search"`
		GraphQL struct {
			Limit     int `json:"limit"`
			Used      int `json:"used"`
			Remaining int `json:"remaining"`
			Reset     int `json:"reset"`
		} `json:"graphql"`
	} `json:"resources"`
	Rate struct {
		Limit     int `json:"limit"`
		Used      int `json:"used"`
		Remaining int `json:"remaining"`
		Reset     int `json:"reset"`
	} `json:"rate"`
}

func CheckRateLimitBeforeRequest(token string) (int, int) {
	// Perform a GET request to the GitHub rate limit endpoint
	url := "https://api.github.com/rate_limit"
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create HTTP request:", err)
		return 0, 0
	}
	req.Header.Set("Authorization", "token "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to make GET request:", err)
		return 0, 0
	}
	defer resp.Body.Close()

	// Decode the JSON response into RateLimitResponse struct
	var rateLimit RateLimitResponse
	err = json.NewDecoder(resp.Body).Decode(&rateLimit)
	if err != nil {
		fmt.Println("Failed to decode JSON:", err)
		return 0, 0
	}

	// Access and return GraphQL and Core remaining requests
	graphQLRemaining := rateLimit.Resources.GraphQL.Remaining
	coreRemaining := rateLimit.Resources.Core.Remaining
	return graphQLRemaining, coreRemaining
}
