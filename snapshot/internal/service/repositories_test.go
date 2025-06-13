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
	"log"
	"power-snapshot/config"
	"testing"

	"github.com/stretchr/testify/assert"

)

func TestGetRepoNames(t *testing.T) {
	// Initialize the logger
	config.InitLogger()

	// Load the configuration from the specified path
	err := config.InitConfig("../../")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}
	tokenManager := NewGitHubTokenManager(config.Client.Github.Token, GithubRateLimit{})
	allRepos := GetRepoNames(EcosystemOrg, GithubUser, tokenManager)
	fmt.Println(len(allRepos))

}

func TestFetchRepositories(t *testing.T) {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos", "ArchlyFi")
	res, err := fetchRepositories(
		url,
		"",
	)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
