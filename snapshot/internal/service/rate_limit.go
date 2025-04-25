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
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

type AtomicTokenCounter struct {
	count atomic.Int32
	_     [64 - 4]byte
}

func (a *AtomicTokenCounter) Add(n int32) {
	a.count.Add(n) // atomic add
}

func (a *AtomicTokenCounter) Set(n int32) {
	a.count.Store(n) // atomic set
}

func (a *AtomicTokenCounter) Get() int32 {
	return a.count.Load() // atomic get
}

func (a *AtomicTokenCounter) Eq(v int32) bool {
	return a.Get() == v
}

func (a *AtomicTokenCounter) Gt(v int32) bool {
	return a.Get() > v
}

func (a *AtomicTokenCounter) Gte(v int32) bool {
	return a.Get() >= v
}

func (a *AtomicTokenCounter) Lt(v int32) bool {
	return a.Get() < v
}

func (a *AtomicTokenCounter) Lte(v int32) bool {
	return a.Get() <= v
}

type GitHubTokenManager struct {
	tokens       []string
	graphQLUsage map[string]*AtomicTokenCounter
	graphQLCap   map[string]*AtomicTokenCounter
	coreUsage    map[string]*AtomicTokenCounter
	coreCap      map[string]*AtomicTokenCounter
	mu           sync.RWMutex
}

func NewGitHubTokenManager(tokens []string) *GitHubTokenManager {
	manager := &GitHubTokenManager{
		tokens:       tokens,
		graphQLUsage: make(map[string]*AtomicTokenCounter),
		graphQLCap:   make(map[string]*AtomicTokenCounter),
		coreUsage:    make(map[string]*AtomicTokenCounter),
		coreCap:      make(map[string]*AtomicTokenCounter),
	}

	for _, token := range tokens {
		manager.graphQLUsage[token] = new(AtomicTokenCounter)
		manager.coreUsage[token] = new(AtomicTokenCounter)
		manager.graphQLCap[token] = new(AtomicTokenCounter)
		manager.coreCap[token] = new(AtomicTokenCounter)
	}

	manager.RefreshToken()
	return manager
}

func (m *GitHubTokenManager) RefreshToken() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, token := range m.tokens {
		gCap, cCap := CheckRateLimitBeforeRequest(token)

		m.graphQLCap[token].Set(gCap)
		m.coreCap[token].Set(cCap)

		zap.L().Info("TokenRefresh",
			zap.String("token", token),
			zap.Int32("graphQL", m.graphQLCap[token].Get()),
			zap.Int32("core", m.coreCap[token].Get()))
	}

}

func (m *GitHubTokenManager) GetAvailableToken() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var (
		maxToken string
		maxCap   int32
	)

	for token := range m.graphQLCap {
		currentCap := m.graphQLCap[token].Get()
		currentUsage := m.graphQLUsage[token].Get()

		if currentUsage < currentCap && currentCap > maxCap && currentCap > 0 {
			maxToken = token
			maxCap = currentCap
		}
	}

	if maxToken != "" {
		m.graphQLUsage[maxToken].Add(1)
	}

	return maxToken
}

func (m *GitHubTokenManager) DecreaseUsage(token string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if usage := m.graphQLUsage[token]; usage != nil {
		usage.Add(-1)
		if current := usage.Get(); current < 0 {
			usage.Set(0)
		}
	}
}

func (m *GitHubTokenManager) GetCoreAvailableToken() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var maxToken string
	maxCapacity := int32(0)

	for token, capacity := range m.coreCap {
		if m.coreUsage[token].Get() < capacity.Get() && capacity.Get() > maxCapacity {
			maxToken = token
			maxCapacity = capacity.Get()
		}
	}

	if maxToken != "" {
		m.coreUsage[maxToken].Add(1)
	}

	return maxToken
}

func (m *GitHubTokenManager) CoreDecreaseUsage(token string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if m.coreUsage[token].Get() > 0 {
		m.coreUsage[token].Add(-1)
	}
}

type RateLimitResponse struct {
	Resources struct {
		Core struct {
			Limit     int32 `json:"limit"`
			Used      int32 `json:"used"`
			Remaining int32 `json:"remaining"`
			Reset     int32 `json:"reset"`
		} `json:"core"`
		Search struct {
			Limit     int32 `json:"limit"`
			Used      int32 `json:"used"`
			Remaining int32 `json:"remaining"`
			Reset     int32 `json:"reset"`
		} `json:"search"`
		GraphQL struct {
			Limit     int32 `json:"limit"`
			Used      int32 `json:"used"`
			Remaining int32 `json:"remaining"`
			Reset     int32 `json:"reset"`
		} `json:"graphql"`
	} `json:"resources"`
	Rate struct {
		Limit     int32 `json:"limit"`
		Used      int32 `json:"used"`
		Remaining int32 `json:"remaining"`
		Reset     int32 `json:"reset"`
	} `json:"rate"`
}

func CheckRateLimitBeforeRequest(token string) (int32, int32) {
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
