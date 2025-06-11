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
	count *atomic.Int32
	_     [64 - 4]byte
}

func NewAtomicTokenCounter() *AtomicTokenCounter {
	return &AtomicTokenCounter{
		count: new(atomic.Int32),
	}
}

func (a *AtomicTokenCounter) Add(n int32) {
	a.count.Add(n) // atomic add
}

func (a *AtomicTokenCounter) Set(n int32) {
	a.count.Store(n) // atomic set
}

func (a *AtomicTokenCounter) Get() int32 {
	if a == nil {
		a = NewAtomicTokenCounter()
	}

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
	githubApi    GithubLimit
}

func NewGitHubTokenManager(tokens []string, githubApi GithubLimit) *GitHubTokenManager {
	manager := &GitHubTokenManager{
		tokens:       tokens,
		graphQLUsage: make(map[string]*AtomicTokenCounter),
		graphQLCap:   make(map[string]*AtomicTokenCounter),
		coreUsage:    make(map[string]*AtomicTokenCounter),
		coreCap:      make(map[string]*AtomicTokenCounter),
		githubApi:    githubApi,
	}

	for _, token := range tokens {
		manager.graphQLUsage[token] = NewAtomicTokenCounter()
		manager.coreUsage[token] = NewAtomicTokenCounter()
		manager.graphQLCap[token] = NewAtomicTokenCounter()
		manager.coreCap[token] = NewAtomicTokenCounter()

	}

	manager.RefreshToken()
	return manager
}

func (m *GitHubTokenManager) RefreshToken() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	res := ""
	for _, token := range m.tokens {
		zap.L().Info("TokenRefresh Before",
			zap.String("token", token),
			zap.String("graphQL", fmt.Sprintf("limit: %d, used: %d", m.graphQLCap[token].Get(), m.graphQLUsage[token].Get())),
			zap.String("core", fmt.Sprintf("limit: %d, used: %d", m.coreCap[token].Get(), m.coreUsage[token].Get())),
		)

		gCap, cCap := m.githubApi.CheckRateLimitBeforeRequest(token)
		m.graphQLCap[token].Set(gCap)
		m.coreCap[token].Set(cCap)
		m.graphQLUsage[token].Set(0)
		m.coreUsage[token].Set(0)

		zap.L().Info("TokenRefresh After",
			zap.String("token", token),
			zap.String("graphQL", fmt.Sprintf("limit: %d, used: %d", gCap, m.graphQLUsage[token].Get())),
			zap.String("core", fmt.Sprintf("limit: %d, used: %d", cCap, m.coreUsage[token].Get())))
		fmt.Printf("\n\n")
		if gCap > 0 && cCap > 0 {
			res = token
		}
	}

	return res
}

func (m *GitHubTokenManager) CheckTokenUsedNum(token string) (int32, int32) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.graphQLUsage[token].Get(), m.coreUsage[token].Get()
}

func (m *GitHubTokenManager) GetAvailableToken() string {
	m.mu.Lock()

	for token := range m.graphQLCap {
		currentCap := m.graphQLCap[token].Get()
		currentUsage := m.graphQLUsage[token].Get()

		if currentUsage < currentCap && currentCap > 0 && currentCap-currentUsage > 15 {
			m.graphQLUsage[token].Add(1)
			m.mu.Unlock()
			return token
		}
	}
	
	m.mu.Unlock()
	return m.RefreshToken()

}

func (m *GitHubTokenManager) GetAllTokenCap() int32 {
	var result int32
	
	m.RefreshToken()
	for _, token := range m.tokens {
		result += m.graphQLCap[token].Get()
	}

	return result
}

func (m *GitHubTokenManager) GetApiCap(token string) (int32, int32) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.graphQLCap[token].Get(), m.coreCap[token].Get()
}

func (m *GitHubTokenManager) GetCoreAvailableToken() string {
	m.mu.Lock()

	for token := range m.coreCap {
		currentCap := m.coreCap[token].Get()
		currentUsage := m.coreUsage[token].Get()
		if currentUsage < currentCap && currentCap > 0 {
			m.coreUsage[token].Add(1)
			m.mu.Unlock()
			return token
		}
	}

	m.mu.Unlock()
	return m.RefreshToken()
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

type GithubLimit interface {
	CheckRateLimitBeforeRequest(token string) (int32, int32)
}

type GithubRateLimit struct {
}

func (g GithubRateLimit) CheckRateLimitBeforeRequest(token string) (int32, int32) {
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
