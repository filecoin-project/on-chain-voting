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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"power-snapshot/config"
	models "power-snapshot/internal/model"
)

var CoreFilecoinOrg = []string{
	"filecoin-project",
}

var EcosystemOrg = []string{
	"team-acaisia",
	"AlgoveraAI",
	"ArchlyFi",
	"ArtifactLabsOfficial",
	"AudiusProject",
	"bacalhau-project",
	"banyancomputer",
	"CIDgravity",
	"collectif-dao",
	"CoopHive",
	"DALN-app",
	"DataHaus",
	"Defluencer",
	"desci-labs",
	"fidlabs",
	"filecoin-saturn",
	"filecoin-station",
	"filedrive-team",
	"Filemarket-xyz",
	"fileverse",
	"filfi",
	"fission-codes",
	"FleekHQ",
	"fluencelabs",
	"foxwallet",
	"froghub-io",
	"functionland",
	"gainforest",
	"glifio",
	"Huddle01",
	"internet-development",
	"lagrangedao",
	"Leto-gg",
	"lighthouse-web3",
	"livepeer",
	"marlinprotocol",
	"Murmuration-Labs",
	"Nebula-Block-Data",
	"nftstorage",
	"fluencelabs",
	"NuLink-network",
	"numbersprotocol",
	"oceanprotocol",
	"okcontract",
	"opscientia",
	"portraitgg",
	"guardianproject",
	"pyth-network",
	"ReppoLabs",
	"SFT-project",
	"sinsoio",
	"starboard-ventures",
	"stfil-io",
	"supranational",
	"filswan",
	"tablelandnetwork",
	"TRANSFERArchive",
	"VCCRI",
	"vshareCloud-Project",
	"bacalhau-project",
	"WeatherXM",
	"web3-storage",
	"ZentaChain",
	"Zondax",
}

var GithubUser = []string{
	"Ghost-Drive",
	"lordshashank",
	"Filenova",
	"3Pad",
	"OpenGateLab",
	"CORGIFILE",
	"xbankingorg",
}

// GetDeveloperWeights Calculates the weight of a developer based on their contributions to a given repository
func GetDeveloperWeights(fromTime time.Time) (map[string]int64, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	var coreFilecoin, ecosystem map[string]int64
	tokenManager := NewGitHubTokenManager(config.Client.Github.Token)

	eg.Go(func() error {
		var err error
		coreFilecoinRepos := GetRepoNames(CoreFilecoinOrg, nil, tokenManager)
		coreFilecoin, err = getDeveloperWeights(ctx, coreFilecoinRepos, 2, fromTime, tokenManager)
		return err
	})

	eg.Go(func() error {
		var err error
		ecosystemRepos := GetRepoNames(EcosystemOrg, GithubUser, tokenManager)
		ecosystem, err = getDeveloperWeights(ctx, ecosystemRepos, 1, fromTime, tokenManager)
		return err
	})

	err := eg.Wait()
	if err != nil {
		zap.L().Error("!!!!!!Error getting developer weights", zap.Error(err))
		return nil, err
	}

	totalWeights := addMerge(coreFilecoin, ecosystem)
	return totalWeights, nil
}

// getDeveloperWeights calculates the weights of developers based on their contributions to the specified repositories.
func getDeveloperWeights(ctx context.Context, repository []string, weight int, fromTime time.Time, tokenManager *GitHubTokenManager) (map[string]int64, error) {
	var wg sync.WaitGroup
	weights := make(chan map[string]int64, len(repository))
	eg, ctx := errgroup.WithContext(ctx)

	for _, fullRepoName := range repository {
		parts := strings.Split(fullRepoName, "/")
		if len(parts) == 2 {
			organization := parts[0]
			repo := parts[1]
			result := addMonths(fromTime, -6)

			since := result.Format("2006-01-02T15:04:05Z")
			time.Sleep(10 * time.Millisecond)

			wg.Add(1)
			eg.Go(func() error {
				token := tokenManager.GetAvailableToken()
				if token == "" {
					return fmt.Errorf("no available token")
				}

				// Use the token
				weightMap, err := getDeveloperWeightsByRepo(ctx, organization, repo, since, token)
				if err != nil {
					// If an error occurs, release the token immediately
					tokenManager.DecreaseUsage(token)
					return err
				}

				// Successfully used the token, send results
				weights <- weightMap

				// Release the token immediately after usage
				tokenManager.DecreaseUsage(token)

				// Done with this goroutine
				wg.Done()

				return nil
			})

		} else {
			zap.L().Info("invalid repository format", zap.String("repository", fullRepoName))
		}
	}

	go func() {
		wg.Wait()
		close(weights)
	}()

	finalWeights := make(map[string]int64)
	for weightMap := range weights {
		for user, weights := range weightMap {
			if weights >= 2 {
				finalWeights[user] = int64(weight)
			}
		}
	}

	err := eg.Wait()
	if err != nil {
		return nil, err
	}

	return finalWeights, nil
}

// getDeveloperWeightsByRepo calculates the weights of developers for a specific repository.
func getDeveloperWeightsByRepo(ctx context.Context, organization, repository, since, token string) (map[string]int64, error) {
	commits, err := getContributorsWithRetry(ctx, organization, repository, since, token, 3, 5*time.Second)
	if err != nil {
		zap.L().Error("failed to get contributors", zap.Error(err))
		return nil, err
	}

	weightMap := make(map[string]int64)
	for _, v := range commits {
		if v.Committer.User.Login != "" {
			weightMap[v.Committer.User.Login]++
		}
		if v.Author.User.Login != "" {
			weightMap[v.Author.User.Login]++
		}
	}

	return weightMap, nil
}

//TODO: 原始算力的计算方式_; UploadIpfs上传的bytes
// getContributors retrieves the commit history of a repository from GitHub GraphQL API.
func getContributors(ctx context.Context, owner, name, since, token string) ([]models.Nodes, error) {
	payload := strings.NewReader(fmt.Sprintf("{\"query\":\"  "+
		"query paginatedCommits($cursor: String) {\\n "+
		"     repository(owner: \\\"%s\\\", name: \\\"%s\\\") {\\n"+
		"        defaultBranchRef {\\n"+
		"          target {\\n"+
		"            ... on Commit {\\n"+
		"              history(first:100, since: \\\"%s\\\"\\n,"+
		" after: $cursor) {\\n"+
		"                totalCount\\n"+
		"                pageInfo {\\n"+
		"                  endCursor\\n"+
		"                  hasNextPage\\n"+
		"                }\\n"+
		"                nodes {\\n"+
		"                  ... on Commit {\\n "+
		"                     committedDate\\n "+
		"                     author {\\n"+
		"                        user {\\n"+
		"                          login\\n "+
		"                       }\\n"+
		"                      }\\n"+
		"                      committer {\\n "+
		"                       user {\\n"+
		"                          login\\n"+
		"                        }\\n"+
		"                      }\\n"+
		"                  }\\n"+
		"                }\\n"+
		"              }\\n"+
		"            }\\n"+
		"          }\\n"+
		"        }\\n"+
		"      }\\n"+
		"    }\",\"variables\":{}}", owner, name, since))

	client := &http.Client{}

	req, err := http.NewRequest("POST", config.Client.Github.GraphQl, payload)
	if err != nil {
		zap.L().Error("failed to create http request", zap.Error(err))
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "api.github.com")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		zap.L().Error("failed to make http request", zap.Error(err))
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("failed to read http response body", zap.Error(err))
		return nil, err
	}
	var calculate models.Developer
	err = json.Unmarshal(body, &calculate)
	if err != nil {
		zap.L().Error("failed to unmarshal json response",
			zap.Error(err),
			zap.String("body", string(body)))
		return nil, err
	}
	time.Sleep(1000 * time.Millisecond)
	return calculate.Data.Repository.DefaultBranchRef.Target.History.Nodes, nil
}

// getContributorsWithRetry retries fetching contributors for a repository with backoff and retry logic.
func getContributorsWithRetry(ctx context.Context, organization, repository, since, token string, maxRetries int, retryInterval time.Duration) ([]models.Nodes, error) {
	var commits []models.Nodes
	var err error

	for i := 0; i < maxRetries; i++ {
		commits, err = getContributors(ctx, organization, repository, since, token)
		if err == nil {
			break
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		zap.L().Error(fmt.Sprintf("Retry %d ", i+1), zap.Error(err))
		time.Sleep(retryInterval)
	}

	select {
	case <-ctx.Done():
		zap.L().Info("ctx done", zap.Error(ctx.Err()))
		return nil, ctx.Err()
	default:
		zap.L().Info("Retrieving contributors", zap.Int("count", len(commits)))
	}

	return commits, err
}

// addMonths adds a specified number of months to a given date.
func addMonths(input time.Time, months int) time.Time {
	date := time.Date(input.Year(), input.Month(), 1, 0, 0, 0, 0, input.Location())
	date = date.AddDate(0, months, 0)

	lastDay := getLastDayOfMonth(date.Year(), int(date.Month()))
	date = time.Date(date.Year(), date.Month(), lastDay, 0, 0, 0, 0, input.Location())

	if input.Day() < lastDay {
		date = time.Date(date.Year(), date.Month(), input.Day(), 0, 0, 0, 0, input.Location())
	}

	return date
}

// getLastDayOfMonth calculates the last day of the specified month in the given year.
func getLastDayOfMonth(year, month int) int {
	lastDay := 31
	nextMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)
	lastDay = nextMonth.Day()
	return lastDay
}

// addMerge Merge weights.
func addMerge(coreFilecoin, ecosystem map[string]int64) map[string]int64 {
	result := make(map[string]int64)

	for coreFilecoinAccount, coreFilecoinWeight := range coreFilecoin {
		result[coreFilecoinAccount] = coreFilecoinWeight
	}

	for ecosystemAccount, ecosystemWeight := range ecosystem {
		if _, ok := result[ecosystemAccount]; !ok {
			result[ecosystemAccount] = ecosystemWeight
		}
	}

	return result
}
