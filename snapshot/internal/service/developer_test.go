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

package service_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"power-snapshot/constant"
	models "power-snapshot/internal/model"
	"power-snapshot/internal/service"
	"power-snapshot/utils"
)

var _ = Describe("Developer", func() {
	date, err := time.Parse("2006-01-02", "2022-01-01")
	Expect(err).To(BeNil())
	var mockFunc = func(ctx context.Context, repos []string, w int, fromTime time.Time, t *service.GitHubTokenManager) (map[string]int64, []models.Nodes, error) {
		queryStart := utils.AddMonths(fromTime, constant.GithubDataWithinXMonths)
		since := queryStart.Format(time.RFC3339)
		Expect(since).NotTo(BeEmpty())
		return map[string]int64{
				"test1": 1,
				"test2": 2,
			}, []models.Nodes{
				{
					CommittedDate: date,
					Author: models.Author{
						User: models.User{
							Login: "testUser1",
						},
					},
					Committer: models.Committer{
						User: models.User{
							Login: "testUser1",
						},
					},
				},
			}, nil
	}

	BeforeEach(func() {
		service.GetDeveloperWeights = mockFunc
	})

	Describe("FetchDeveloperWeights", func() {
		It("should success", func() {

			Expect(err).To(BeNil())
			tokenManager := service.NewGitHubTokenManager(conf.Github.Token, &mockgithubRateLimit{})
			res, nodes, err := syncService.FetchDeveloperWeights(date, &mockGithubRepo{}, tokenManager)
			Expect(err).To(BeNil())
			Expect(res).To(Equal(map[string]int64{
				"test1": 1,
				"test2": 2,
			}))
			Expect(nodes).ToNot(BeEmpty())
		})
	})

	Describe("addMerge", func() {
		It("should equal", func() {
			map1 := map[string]int64{
				"test1": 1,
				"test2": 2,
			}
			map2 := map[string]int64{
				"test1": 1,
				"test3": 3,
			}
			expected := map[string]int64{
				"test1": 1,
				"test2": 2,
				"test3": 3,
			}
			Expect(service.AddMerge(map1, map2)).To(Equal(expected))
		})
	})
})

type mockGithubRepo struct{}

func (m *mockGithubRepo) GetRepoNames(orgs []string, users []string, tokenManager *service.GitHubTokenManager) []string {
	return []string{"testOrg/testRepo"}
}

type mockgithubRateLimit struct{}

func (m *mockgithubRateLimit) CheckRateLimitBeforeRequest(token string) (int32, int32) {
	return 20_000, 20_000
}
