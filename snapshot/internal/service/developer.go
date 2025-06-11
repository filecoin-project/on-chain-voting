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
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"power-snapshot/config"
	"power-snapshot/constant"
	models "power-snapshot/internal/model"
	"power-snapshot/utils"
)

var GetDeveloperWeights = getDeveloperWeights

// FetchDeveloperWeights Calculates the weight of a developer based on their contributions to a given repository
func (s *SyncService) FetchDeveloperWeights(fromTime time.Time, githubRepo IGithub, tokenManager *GitHubTokenManager) (map[string]int64, []models.Nodes, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg, errCtx := errgroup.WithContext(ctx)

	var coreFilecoin, ecosystem map[string]int64
	var coreCommits, ecosystemCommits, commitsResult []models.Nodes

	if tokenManager.GetAllTokenCap() < constant.MinimumTokenCapacity {
		zap.L().Warn(
			"Current token capacity is less than minimum token capacity",
			zap.Int32("current token capacity", tokenManager.GetAllTokenCap()),
			zap.Int32("minimum token capacity", constant.MinimumTokenCapacity),
		)
		return nil, nil, constant.ErrorNoTokenAvailable
	}

	eg.Go(func() error {
		var err error
		coreOrg, err := s.backendRpcClient.GetGighubRepoInfo(constant.GithubCoreOrg)
		if err != nil {
			return err
		}
		coreFilecoinRepos := githubRepo.GetRepoNames(coreOrg, nil, tokenManager)
		coreFilecoin, coreCommits, err = GetDeveloperWeights(errCtx, coreFilecoinRepos, 2, fromTime, tokenManager)
		return err
	})

	eg.Go(func() error {
		var err error
		ecosystemOrg, err := s.backendRpcClient.GetGighubRepoInfo(constant.GithubEcosystemOrg)
		if err != nil {
			return err
		}
		githubUser, err := s.backendRpcClient.GetGighubRepoInfo(constant.GithubUser)
		if err != nil {
			return err
		}
		ecosystemRepos := githubRepo.GetRepoNames(ecosystemOrg, githubUser, tokenManager)
		ecosystem, ecosystemCommits, err = GetDeveloperWeights(errCtx, ecosystemRepos, 1, fromTime, tokenManager)
		return err
	})

	err := eg.Wait()
	if err != nil {
		return nil, nil, err
	}

	commitsResult = append(commitsResult, coreCommits...)
	commitsResult = append(commitsResult, ecosystemCommits...)

	totalWeights := AddMerge(coreFilecoin, ecosystem)
	return totalWeights, commitsResult, nil
}

// getDeveloperWeights calculates the weights of developers based on their contributions to the specified repositories.
func getDeveloperWeights(ctx context.Context, repositories []string, weight int, fromTime time.Time, tokenManager *GitHubTokenManager) (map[string]int64, []models.Nodes, error) {
	const maxConcurrency = 10
	var (
		mu           sync.Mutex
		stopTag      int32
		finalWeights = make(map[string]int64)
		commitsPool  []models.Nodes
		errs         []error
		eg           errgroup.Group
	)
	eg.SetLimit(maxConcurrency)
	for index, repo := range repositories {
		if atomic.LoadInt32(&stopTag) == 1 {
			zap.L().Error("!!!!!! Failed to process repos, break loop !!!!!", zap.Error(constant.ErrorNoTokenAvailable))
			return nil, nil, constant.ErrorNoTokenAvailable
		}

		parts := strings.Split(repo, "/")
		if len(parts) != 2 {
			zap.L().Warn("invalid repository format", zap.String("repo", repo))
			continue
		}
		org, repoName := parts[0], parts[1]

		eg.Go(func() error {
			weightMap, commits, err := GetRepoData(ctx, index, len(repositories), org, repoName, fromTime, tokenManager)
			if err != nil {
				if errors.Is(err, constant.ErrorNoTokenAvailable) {
					atomic.StoreInt32(&stopTag, 1)
				}

				mu.Lock()
				errs = append(errs, fmt.Errorf("repo %s/%s error: %w", org, repoName, err))
				mu.Unlock()

				return nil
			}

			mu.Lock()
			commitsPool = append(commitsPool, commits...)
			for user, cnt := range weightMap {
				if cnt >= 2 {
					finalWeights[user] = int64(weight)
				}
			}
			mu.Unlock()
			return nil
		})
	}

	_ = eg.Wait()
	if len(errs) > 0 {
		zap.L().Error("!!!!!! Failed to process repos !!!!!!", zap.Int("err count", len(errs)), zap.Errors("errors", errs))
		if len(errs) > 3 {
			return finalWeights, commitsPool, fmt.Errorf("%d errors occurred: %v", len(errs), errors.Join(errs...))
		}
	}

	return finalWeights, commitsPool, nil
}

type IGithubRepoDetails interface {
	getDeveloperWeights(ctx context.Context, repositories []string, weight int, fromTime time.Time, tokenManager *GitHubTokenManager)
}
type GithubRepoDetails struct {
}

func GetRepoData(ctx context.Context, index, repoCount int, org, repo string, fromTime time.Time, tokenManager *GitHubTokenManager) (map[string]int64, []models.Nodes, error) {
	queryStart := utils.AddMonths(fromTime, -constant.GithubDataWithinXMonths)
	since := queryStart.Format(time.RFC3339)

	return getDeveloperWeightsByRepo(ctx, index, repoCount, org, repo, since, tokenManager)
}

// getDeveloperWeightsByRepo calculates the weights of developers for a specific repository.
func getDeveloperWeightsByRepo(ctx context.Context, index, repoCount int, organization, repository, since string, tokenManager *GitHubTokenManager) (map[string]int64, []models.Nodes, error) {
	commits, err := getContributorsWithRetry(ctx, index, repoCount, organization, repository, since, tokenManager, time.Second)
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

// getContributorsWithRetry retries fetching contributors for a repository with backoff and retry logic.
func getContributorsWithRetry(ctx context.Context, index, repoCount int, organization, repository, since string, tokenManager *GitHubTokenManager, retryInterval time.Duration) ([]models.Nodes, error) {
	retries := len(config.Client.Github.Token) + 1

	var commits []models.Nodes
	var lastErr error
	for attempt := 0; attempt < retries; attempt++ {
		select {
		case <-ctx.Done():
			zap.L().Warn("Operation canceled",
				zap.Error(ctx.Err()),
				zap.Int("attempt", attempt))
			return nil, ctx.Err()
		default:
			reqToken := tokenManager.GetAvailableToken()
			if reqToken == "" {
				zap.L().Warn("All tokens are reached with the usage rate limit, wait for an hour and try again")
				return nil, constant.ErrorNoTokenAvailable
			}

			commits, lastErr = getContributors(organization, repository, since, reqToken)
			if lastErr == nil {
				break
			}

			if strings.Contains(lastErr.Error(), "API rate limit exceeded") {
				zap.L().Warn("Rate limit exceeded, refreshing token")
				tokenManager.RefreshToken()
			}

			zap.L().Warn("Retrying...",
				zap.Int("fetch index", index),
				zap.String("repo", fmt.Sprintf("%s/%s", organization, repository)),
				zap.Int("attempt", attempt+1),
				zap.Duration("wait", retryInterval),
				zap.Error(lastErr))

			time.Sleep(retryInterval)
		}
	}

	if lastErr != nil {
		zap.L().Error("Permanent failure after retries",
			zap.Int("max_retries", retries),
			zap.Error(lastErr))
		return nil, fmt.Errorf("after %d retries: %w", retries, lastErr)
	}

	zap.L().Info("Finalizing data processing",
		zap.Int("count", len(commits)),
		zap.String("organization", organization),
		zap.String("repository", repository),
		zap.String("Syncd progress", fmt.Sprintf("%d/%d", index+1, repoCount)),
	)

	return commits, nil
}

// AddMerge Merge weights.
func AddMerge(coreFilecoin, ecosystem map[string]int64) map[string]int64 {
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

// getContributors retrieves the commit history of a repository from GitHub GraphQL API.
func getContributors(owner, name, since, token string) ([]models.Nodes, error) {
	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			if r.StatusCode() == 404 {
				return false
			}
			return r.StatusCode() >= 500 || err != nil
		})

	var allNodes []models.Nodes
	var endCursor *string
	for {
		var response models.Developer
		resp, err := client.R().
			SetAuthToken(token).
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
		time.Sleep(time.Second)
	}

	return allNodes, nil
}
