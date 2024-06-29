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
)

// Repository represents a GitHub repository.
type Repository struct {
	Name string `json:"name"`
}

// GetRepoNames retrieves names of repositories for given organizations or users.
func GetRepoNames(orgs []string, users []string, tokenManager *GitHubTokenManager) []string {
	var wg sync.WaitGroup
	wg.Add(len(orgs) + len(users))

	reposChan := make(chan []string, len(orgs)+len(users))

	fetchRepos := func(entity, name string, fetchFunc func(string, string) ([]Repository, error)) {
		defer wg.Done()
		token := tokenManager.GetCoreAvailableToken()
		if token == "" {
			return
		}
		defer tokenManager.CoreDecreaseUsage(token)

		repos, err := fetchFunc(entity, token)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Error fetching %s repositories", name), zap.String("entity", entity), zap.Error(err))
			return
		}

		var formattedRepos []string
		for _, repo := range repos {
			formattedRepos = append(formattedRepos, fmt.Sprintf("%s/%s", entity, repo.Name))
		}
		zap.L().Info(fmt.Sprintf("Obtained %s repositories", name), zap.String("entity", entity))
		reposChan <- formattedRepos
	}

	// Fetch repositories for organizations
	for _, org := range orgs {
		go fetchRepos(org, "organization", getAllOrgRepos)
	}

	// Fetch repositories for users
	for _, user := range users {
		go fetchRepos(user, "user", getAllUserRepos)
	}

	go func() {
		wg.Wait()
		close(reposChan)
	}()

	var allRepos []string
	for repos := range reposChan {
		allRepos = append(allRepos, repos...)
	}
	return allRepos
}

// getAllOrgRepos fetches all repositories for a given organization from GitHub.
func getAllOrgRepos(org, token string) ([]Repository, error) {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos", org)
	return fetchRepositories(url, token)
}

// getAllUserRepos fetches all repositories for a given user from GitHub.
func getAllUserRepos(user, token string) ([]Repository, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", user)
	return fetchRepositories(url, token)
}

// fetchRepositories performs the common logic for fetching repositories from GitHub.
func fetchRepositories(url, token string) ([]Repository, error) {
	var allRepos []Repository
	client := http.Client{}

	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		req.Header.Set("Authorization", "token "+token)
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		q := req.URL.Query()
		q.Add("per_page", "100")
		req.URL.RawQuery = q.Encode()

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %v", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("not found (404)")
		}

		var repos []Repository
		err = json.NewDecoder(resp.Body).Decode(&repos)
		if err != nil {
			return nil, fmt.Errorf("error unmarshal JSON: %v", err)
		}

		allRepos = append(allRepos, repos...)

		nextPageURL := getNextPageURL(resp.Header.Get("Link"))
		if nextPageURL == "" {
			break
		}

		url = nextPageURL
	}

	return allRepos, nil
}

// getNextPageURL extracts the next page URL from the Link header.
func getNextPageURL(linkHeader string) string {
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
