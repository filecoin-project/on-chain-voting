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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"power-snapshot/internal/service"
)

var _ = Describe("RateLimit", func() {
	var tokenManager *service.GitHubTokenManager
	BeforeEach(func() {
		tokenManager = service.NewGitHubTokenManager(
			[]string{"token1"},
			mockGithubLimit,
		)
	})
	Context("GitHubTokenManager", func() {
		When("token not used", func() {
			It("token cap should be full", func() {
				cap := tokenManager.GetAllTokenCap()
				Expect(cap).To(Equal(int32(5000)))
			})
			It("should return gap", func() {
				gCap, cCap := tokenManager.GetApiCap("token1")
				Expect(gCap).To(Equal(int32(5000)))
				Expect(cCap).To(Equal(int32(5000)))
			})

			It("should return token", func() {
				token := tokenManager.GetCoreAvailableToken()
				Expect(token).To(Equal("token1"))
			})
		})

		When("token used", func() {
			BeforeEach(func() {
				for i := 0; i < 100; i++ {
					tokenManager.GetAvailableToken()
				}
			})
			It("should graph token must be 4900", func() {
				gCap, cCap := tokenManager.CheckTokenUsedNum("token1")
				Expect(gCap).To(Equal(int32(100)))
				Expect(cCap).To(Equal(int32(0)))
			})

		})
	})
})
