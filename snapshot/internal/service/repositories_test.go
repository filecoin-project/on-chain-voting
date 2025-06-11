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
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"power-snapshot/internal/service"
	mocks "power-snapshot/mock"
)

var _ = Describe("Repositories", func() {
	var mockIgithub *mocks.MockIGithub
	var tokenManager *service.GitHubTokenManager
	BeforeEach(func() {
		mockIgithub = mocks.NewMockIGithub(mockCtrl)
		tokenManager = service.NewGitHubTokenManager(
			[]string{"token1"},
			mockGithubLimit,
		)
		mockIgithub.EXPECT().
			GetRepoNames(gomock.Any(), gomock.Any(), gomock.Any()).
			Return([]string{"repo1", "repo2"})
	})
	Describe("GetRepoNames", func() {
		It("should return repo names", func() {
			repoNames := mockIgithub.GetRepoNames([]string{"org1"}, []string{"user1"}, tokenManager)
			Expect(repoNames).To(Equal([]string{"repo1", "repo2"}))
		})
	})
})
