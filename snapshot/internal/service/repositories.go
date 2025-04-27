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
	"fmt"
	"power-snapshot/utils"
	"sync"

	"go.uber.org/zap"
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

	for {
		var repos []Repository
		linkHeader, err := utils.FetchGithubDeveloper(url, token, map[string]string{"per_page": "100"}, &repos)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		nextPageURL := utils.GetNextPageURL(linkHeader)
		if nextPageURL == "" {
			break
		}

		url = nextPageURL
	}

	return allRepos, nil
}
