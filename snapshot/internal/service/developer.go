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
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"power-snapshot/config"
	"power-snapshot/constant"
	models "power-snapshot/internal/model"
	"power-snapshot/utils"
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

// FetchDeveloperWeights Calculates the weight of a developer based on their contributions to a given repository
func FetchDeveloperWeights(fromTime time.Time) (map[string]int64, []models.Nodes, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	var coreFilecoin, ecosystem map[string]int64
	var coreCommits, ecosystemCommits, commitsResult []models.Nodes
	tokenManager := NewGitHubTokenManager(config.Client.Github.Token, GithubRateLimit{})

	eg.Go(func() error {
		var err error
		coreFilecoinRepos := GetRepoNames(CoreFilecoinOrg, nil, tokenManager)
		coreFilecoin, coreCommits, err = getDeveloperWeights(ctx, coreFilecoinRepos, 2, fromTime, tokenManager)
		return err
	})

	eg.Go(func() error {
		var err error
		ecosystemRepos := GetRepoNames(EcosystemOrg, GithubUser, tokenManager)
		ecosystem, ecosystemCommits, err = getDeveloperWeights(ctx, ecosystemRepos, 1, fromTime, tokenManager)
		return err
	})

	err := eg.Wait()
	if err != nil {
		zap.L().Error("!!!!!!Error getting developer weights", zap.Error(err))
		return nil, nil, err
	}

	commitsResult = append(commitsResult, coreCommits...)
	commitsResult = append(commitsResult, ecosystemCommits...)

	totalWeights := addMerge(coreFilecoin, ecosystem)
	return totalWeights, commitsResult, nil
}

// getDeveloperWeights calculates the weights of developers based on their contributions to the specified repositories.
func getDeveloperWeights(ctx context.Context, repositories []string, weight int, fromTime time.Time, tokenManager *GitHubTokenManager) (map[string]int64, []models.Nodes, error) {
	const maxConcurrency = 10
	var (
		wg          sync.WaitGroup
		mu          sync.Mutex
		errCh       = make(chan error, len(repositories))
		weightsCh   = make(chan map[string]int64, len(repositories))
		commitsPool []models.Nodes
		sem         = make(chan struct{}, maxConcurrency)
	)

	for index, repo := range repositories {
		parts := strings.Split(repo, "/")
		if len(parts) != 2 {
			zap.L().Warn("invalid repository format", zap.String("repo", repo))
			continue
		}
		org, repoName := parts[0], parts[1]

		wg.Add(1)
		go func(org, repo string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			weightMap, commits, err := getRepoData(ctx, index, len(repositories), org, repo, fromTime, tokenManager)
			if err != nil {
				errCh <- fmt.Errorf("repo %s error: %w", repo, err)
				return
			}

			mu.Lock()
			commitsPool = append(commitsPool, commits...)
			mu.Unlock()

			select {
			case weightsCh <- weightMap:
			case <-ctx.Done():
			}
		}(org, repoName)
	}

	go func() {
		wg.Wait()
		close(weightsCh)
		close(errCh)
	}()

	finalWeights := make(map[string]int64)
	for w := range weightsCh {
		for user, cnt := range w {
			if cnt >= 2 {
				finalWeights[user] = int64(weight)
			}
		}
	}

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return finalWeights, commitsPool, fmt.Errorf("%d errors occurred: %v", len(errs), errors.Join(errs...))
	}

	return finalWeights, commitsPool, nil
}

func getRepoData(ctx context.Context, index, repoCount int, org, repo string, fromTime time.Time, tokenManager *GitHubTokenManager) (map[string]int64, []models.Nodes, error) {
	queryStart := utils.AddMonths(fromTime, -constant.GithubDataWithinXMonths)
	since := queryStart.Format(time.RFC3339)

	return getDeveloperWeightsByRepo(ctx, index, repoCount, org, repo, since, tokenManager)
}

// getDeveloperWeightsByRepo calculates the weights of developers for a specific repository.
func getDeveloperWeightsByRepo(ctx context.Context, index, repoCount int, organization, repository, since string, tokenManager *GitHubTokenManager) (map[string]int64, []models.Nodes, error) {
	commits, err := getContributorsWithRetry(ctx, index, repoCount, organization, repository, since, tokenManager, 5*time.Second)
	if err != nil {
		zap.L().Error("failed to get contributors", zap.Error(err))
		return nil, nil, err
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

	return weightMap, commits, nil
}

// getContributors retrieves the commit history of a repository from GitHub GraphQL API.
func getContributors(ctx context.Context, owner, name, since, token string) ([]models.Nodes, error) {
	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return r.StatusCode() >= 500 || err != nil
		})

	var allNodes []models.Nodes
	var endCursor *string
	for {
		var response models.Developer
		resp, err := client.R().
			SetAuthToken(token).
			SetContext(ctx).
			SetHeaders(map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/vnd.github.v3+json",
				"User-Agent":   "Golang-Resty-Client",
			}).
			SetBody(map[string]any{
				"query": `
                    query paginatedCommits($cursor: String) {
                        repository(owner: "` + owner + `", name: "` + name + `") {
                            defaultBranchRef {
                                target {
                                    ... on Commit {
                                        history(first: 100, since: "` + since + `", after: $cursor) {
                                            nodes {
                                                ... on Commit {
                                                    committedDate
                                                    author { user { login } }
                                                    committer { user { login } }
                                                }
                                            }
                                            pageInfo {
                                                endCursor
                                                hasNextPage
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }`,
				"variables": map[string]interface{}{
					"cursor": endCursor,
				},
			}).
			SetResult(&response).
			Post(config.Client.Github.GraphQl)

		if err != nil || resp.StatusCode() != 200 {
			zap.L().Error("Request failed",
				zap.Error(err),
				zap.Int("status", resp.StatusCode()),
				zap.String("response", resp.String()))
			return nil, fmt.Errorf("request failed: %v", err)
		}

		if len(response.Errors) > 0 {
			errorMessages := make([]string, len(response.Errors))
			for i, e := range response.Errors {
				errorMessages[i] = e.Message
			}
			return nil, fmt.Errorf("GraphQL errors: %v", errorMessages)
		}

		history := response.Data.Repository.DefaultBranchRef.Target.History
		allNodes = append(allNodes, history.Nodes...)

		if !history.PageInfo.HasNextPage {
			break
		}

		endCursor = &history.PageInfo.EndCursor
		time.Sleep(1 * time.Second)
	}

	return allNodes, nil
}

// getContributorsWithRetry retries fetching contributors for a repository with backoff and retry logic.
func getContributorsWithRetry(ctx context.Context, index, repoCount int, organization, repository, since string, tokenManager *GitHubTokenManager, retryInterval time.Duration) ([]models.Nodes, error) {
	reqToken := tokenManager.GetAvailableToken()
	if reqToken == "" {
		zap.L().Warn("All tokens are reached with the usage rate limit, wait for an hour and try again")
		return nil, fmt.Errorf("all tokens are reached with the usage rate limit, wait for an hour and try again")
	}

	defer tokenManager.DecreaseUsage(reqToken)

	var commits []models.Nodes
	var lastErr error
	for attempt := 0; attempt < constant.MaxGithubGraphRequetRetries; attempt++ {
		select {
		case <-ctx.Done():
			zap.L().Warn("Operation canceled",
				zap.Error(ctx.Err()),
				zap.Int("attempt", attempt))
			continue
		default:
			commits, lastErr = getContributors(ctx, organization, repository, since, reqToken)
			if lastErr == nil {
				break
			}

			if strings.Contains(lastErr.Error(), "API rate limit exceeded") {
				zap.L().Warn("Rate limit exceeded, refreshing token")
				tokenManager.RefreshToken()
				reqToken = tokenManager.GetAvailableToken()
			}

			gCap, _ := tokenManager.GetApiCap(reqToken)
			if gCap == 0 {
				reqToken = tokenManager.GetAvailableToken()
			}

			waitTime := time.Duration(math.Pow(2, float64(attempt))) * retryInterval
			zap.L().Warn("Retrying...",
				zap.Int("attempt", attempt+1),
				zap.Duration("wait", waitTime),
				zap.Error(lastErr))

			time.Sleep(waitTime)
		}
	}

	if lastErr != nil {
		zap.L().Error("Permanent failure after retries",
			zap.Int("max_retries", constant.MaxGithubGraphRequetRetries),
			zap.Error(lastErr))
		return nil, fmt.Errorf("after %d retries: %w", constant.MaxGithubGraphRequetRetries, lastErr)
	}

	zap.L().Info("Finalizing data processing",
		zap.Int("count", len(commits)),
		zap.String("organization", organization),
		zap.String("repository", repository),
		zap.String("Syncd progress", fmt.Sprintf("%d/%d", index+1, repoCount)),
	)

	return commits, nil
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
