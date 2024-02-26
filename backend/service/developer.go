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
	"backend/config"
	"backend/models"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

var CoreFilecoinRepos = []string{
	"filecoin-project/go-hamt-ipld",
	"filecoin-project/venus",
	"filecoin-project/assets",
	"filecoin-project/consensus",
	"filecoin-project/filecoin-network-viz",
	"filecoin-project/filecoin-network-sim",
	"filecoin-project/rust-fil-proofs",
	"filecoin-project/merkletree",
	"filecoin-project/go-ipld-cbor",
	"filecoin-project/go-leb128",
	"filecoin-project/sample-data",
	"filecoin-project/bls-signatures",
	"filecoin-project/specs",
	"filecoin-project/community",
	"filecoin-project/sapling-crypto",
	"filecoin-project/filecoin-explorer",
	"filecoin-project/replication-game",
	"filecoin-project/research",
	"filecoin-project/fff",
	"filecoin-project/paired",
	"filecoin-project/infra",
	"filecoin-project/designdocs",
	"filecoin-project/replication-game-leaderboard",
	"filecoin-project/bitsets",
	"filecoin-project/fil-calculations",
	"filecoin-project/filecoin",
	"filecoin-project/rust-fil-secp256k1",
	"filecoin-project/lua-filecoin",
	"filecoin-project/go-filecoin-badges",
	"filecoin-project/bellperson",
	"filecoin-project/phase2",
	"filecoin-project/FIPs",
	"filecoin-project/filecoin-contributors",
	"filecoin-project/orient",
	"filecoin-project/rust-filbase",
	"filecoin-project/starling",
	"filecoin-project/rust-fil-sector-builder",
	"filecoin-project/rust-fil-proofs-ffi",
	"filecoin-project/rust-fil-ffi-toolkit",
	"filecoin-project/lotus",
	"filecoin-project/go-bls-sigs",
	"filecoin-project/go-sectorbuilder",
	"filecoin-project/venus-docs",
	"filecoin-project/cryptolab",
	"filecoin-project/zexe",
	"filecoin-project/PebblingAndDepthReductionAttacks",
	"filecoin-project/devgrants",
	"filecoin-project/filecoin-http-api",
	"filecoin-project/go-amt-ipld",
	"filecoin-project/drg-attacks",
	"filecoin-project/cpp-filecoin",
	"filecoin-project/go-data-transfer",
	"filecoin-project/security-analysis",
	"filecoin-project/chain-validation",
	"filecoin-project/go-http-api",
	"filecoin-project/fil-ocl",
	"filecoin-project/go-storage-miner",
	"filecoin-project/group",
	"filecoin-project/rust-fil-logger",
	"filecoin-project/lotus-docs",
	"filecoin-project/filecoin-ffi",
	"filecoin-project/go-fil-filestore",
	"filecoin-project/go-retrieval-market-project",
	"filecoin-project/go-fil-markets",
	"filecoin-project/go-shared-types",
	"filecoin-project/go-address",
	"filecoin-project/go-statestore",
	"filecoin-project/go-cbor-util",
	"filecoin-project/go-crypto",
	"filecoin-project/go-paramfetch",
	"filecoin-project/powersoftau",
	"filecoin-project/rust-filecoin-proofs-api",
	"filecoin-project/specs-actors",
	"filecoin-project/serialization-vectors",
	"filecoin-project/go-statemachine",
	"filecoin-project/go-padreader",
	// "filecoin-project/lotus-badges", // Error
	"filecoin-project/go-bitfield",
	"filecoin-project/go-fil-commcid",
	"filecoin-project/slate",
	"filecoin-project/specs-storage",
	"filecoin-project/fungi",
	"filecoin-project/filecoin-docs",
	"filecoin-project/sector-storage",
	"filecoin-project/storage-fsm",
	"filecoin-project/benchmarks",
	"filecoin-project/go-storedcounter",
	"filecoin-project/rust-sha2ni",
	"filecoin-project/ec-gpu",
	"filecoin-project/rust-fil-nse-gpu",
	"filecoin-project/sentinel-drone",
	"filecoin-project/sentinel",
	"filecoin-project/lotus-archived",
	"filecoin-project/go-jsonrpc",
	"filecoin-project/oni",
	"filecoin-project/mapr",
	"filecoin-project/phase2-attestations",
	"filecoin-project/filecoin-client-tutorial",
	"filecoin-project/slate-react-system",
	"filecoin-project/rust-gpu-tools",
	"filecoin-project/go-multistore",
	"filecoin-project/network-info",
	"filecoin-project/slate-stats",
	"filecoin-project/blstrs",
	"filecoin-project/test-vectors",
	"filecoin-project/statediff",
	// "filecoin-project/phase2-attestations-internal-verification", // Error
	"filecoin-project/fil-blst",
	"filecoin-project/blst",
	"filecoin-project/go-state-types",
	"filecoin-project/fuzzing-lotus",
	"filecoin-project/storage.filecoin.io",
	"filecoin-project/go-ds-versioning",
	"filecoin-project/core-devs",
	"filecoin-project/lily",
	"filecoin-project/helm-charts",
	"filecoin-project/notary-governance",
	"filecoin-project/website-security.filecoin.io",
	"filecoin-project/website-bounty.filecoin.io",
	"filecoin-project/product",
	"filecoin-project/ent",
	"filecoin-project/filecoin-discover-validator",
	"filecoin-project/go-commp-utils",
	"filecoin-project/go-bs-lmdb",
	"filecoin-project/filecoin-phase2",
	"filecoin-project/filecoin-plus-client-onboarding",
	"filecoin-project/sentinel-tick",
	"filecoin-project/taupipp",
	"filecoin-project/venus-sealer",
	"filecoin-project/community-china",
	"filecoin-project/venus-wallet",
	"filecoin-project/go-bs-postgres-chainnotated",
	"filecoin-project/retrieval-market-spec",
	"filecoin-project/go-fil-commp-hashhash",
	"filecoin-project/ffi-stub",
	"filecoin-project/dealbot",
	"filecoin-project/filecoin-plus-large-datasets",
	"filecoin-project/go-legs",
	"filecoin-project/sentinel-locations",
	"filecoin-project/go-dagaggregator-unixfs",
	"filecoin-project/dagstore",
	"filecoin-project/homebrew-lotus",
	"filecoin-project/eudico",
	"filecoin-project/sturdy-journey",
	"filecoin-project/ecodash",
	"filecoin-project/lily-archiver",
	"filecoin-project/fvm-runtime-experiment",
	"filecoin-project/system-test-matrix",
	"filecoin-project/fvm-specs",
	"filecoin-project/filecoin-discover-dealer",
	"filecoin-project/fil-rustacuda",
	"filecoin-project/retrieval-market",
	"filecoin-project/cgo-blockstore",
	"filecoin-project/venus-common-types",
	"filecoin-project/boost",
	"filecoin-project/bellperson-gadgets",
	"filecoin-project/ref-fvm",
	// "filecoin-project/lightning-planning", // Error
	// "filecoin-project/go-data-transfer-bus", // Error
	"filecoin-project/athena",
	"filecoin-project/fvm-test-vectors",
	"filecoin-project/data-transfer-benchmark",
	"filecoin-project/builtin-actors",
	"filecoin-project/filecoin-retrieval-market-website",
	"filecoin-project/boost-docs",
	"filecoin-project/sp-mentorship-grants",
	"filecoin-project/ref-fvm-fuzz-corpora",
	"filecoin-project/builtin-actors-bundler",
	"filecoin-project/fil_pasta_curves",
	"filecoin-project/filecoin-plus-leaderboard",
	"filecoin-project/go-cid-tools",
	"filecoin-project/client-growth",
	"filecoin-project/fvm-wasm-instrument",
	"filecoin-project/mir",
	"filecoin-project/filecoin-chain-archiver",
	"filecoin-project/gitops-profile-catalog",
	"filecoin-project/pubsub",
	"filecoin-project/fvm-docs",
	"filecoin-project/fvm-evm",
	"filecoin-project/filecoin-plus-leaderboard-data",
	"filecoin-project/go-ulimit",
	"filecoin-project/wasmtime",
	"filecoin-project/testnet-wallaby",
	"filecoin-project/cache-action",
	"filecoin-project/fvm-example-actors",
	"filecoin-project/tornado-deploy",
	"filecoin-project/halo2",
	"filecoin-project/near-blake2",
	"filecoin-project/parity-wasm",
	"filecoin-project/chain-love-website",
	"filecoin-project/homebrew-exporter",
	"filecoin-project/snapcraft-exporter",
	"filecoin-project/filet",
	"filecoin-project/release-metrics-gitops",
	"filecoin-project/docker-hub-exporter",
	"filecoin-project/github-exporter",
	"filecoin-project/wge-post-deployment-template-renderer",
	"filecoin-project/fevm-hardhat-kit",
	"filecoin-project/filecoin-actor-utils",
	"filecoin-project/wasm-tools",
	"filecoin-project/lassie",
	"filecoin-project/fvm-bench",
	"filecoin-project/circleci-exporter",
	"filecoin-project/fevm-foundry-kit",
	"filecoin-project/testnet-hyperspace",
	"filecoin-project/gov_docs",
	"filecoin-project/awesome-filecoin",
	"filecoin-project/go-data-segment",
	"filecoin-project/retrieval-load-testing",
	"filecoin-project/filcryo",
	"filecoin-project/fevm-contract-tests",
	"filecoin-project/openzeppelin-contracts",
	"filecoin-project/terraform-google-lily-starter",
	"filecoin-project/go-retrieval-types",
	"filecoin-project/lassie-event-recorder",
	"filecoin-project/rhea-load-testing",
	"filecoin-project/fil-sppark",
	"filecoin-project/fvm-starter-kit-deal-making",
	"filecoin-project/data-prep-tools",
	"filecoin-project/solidity-data-segment",
	"filecoin-project/TODOnotes",
	"filecoin-project/boost-graphsync",
	"filecoin-project/boost-gfm",
	"filecoin-project/filecoin-fvm-localnet",
	"filecoin-project/filecoin-docs-CAT",
	"filecoin-project/filecoin-docs-stefnotes",
	"filecoin-project/fevm-data-dao-kit",
	"filecoin-project/pc2_cuda",
	"filecoin-project/FIPs-1",
	"filecoin-project/cod-starter-kit",
	"filecoin-project/testnet-calibration",
	"filecoin-project/raas-starter-kit",
	"filecoin-project/filsnap",
	"filecoin-project/kubo-api-client",
	"filecoin-project/filecoin-solidity",
	"filecoin-project/motion",
	"filecoin-project/filecoin-plus-application-json-restructure",
	"filecoin-project/fil-docs-bot",
	"filecoin-project/sp-automation",
	"filecoin-project/gb-filecoin-docs",
	"filecoin-project/motion-s3-connector",
	"filecoin-project/motion-cloudserver",
	"filecoin-project/motion-arsenal",
	"filecoin-project/filplus-backend",
	"filecoin-project/filplus-tooling-backend-test",
	"filecoin-project/filplus-registry",
	"filecoin-project/filplus-utils",
	"filecoin-project/filplus-ssa-bot",
	"filecoin-project/filecoin-plus-roadmap",
	"filecoin-project/biscuit-dao",
	"filecoin-project/roadmap",
	"filecoin-project/meta-aggregator-frontend",
	"filecoin-project/filecoin-plus-falcon",
	"filecoin-project/motion-sp-test",
	"filecoin-project/tla-f3",
	"filecoin-project/go-f3",
	"filecoin-project/state-storage-starter-kit",
}

var EcosystemRepos = []string{
	"libp2p/go-libp2p",
	"libp2p/rust-libp2p",
	"libp2p/js-libp2p",
}

// GetDeveloperWeights calculates the weights of developers based on their contributions to specified repositories.
// It combines the weights from CoreFilecoinRepos and EcosystemRepos.
func GetDeveloperWeights() map[string]int {
	coreFilecoin := getDeveloperWeights(CoreFilecoinRepos, 2)
	ecosystem := getDeveloperWeights(EcosystemRepos, 1)
	totalWeights := addMerge(coreFilecoin, ecosystem)
	return totalWeights
}

// getDeveloperWeights calculates the weights of developers based on their contributions to the specified repositories.
func getDeveloperWeights(repository []string, weight int) map[string]int {
	var wg sync.WaitGroup
	weights := make(chan map[string]int)

	for _, fullRepoName := range repository {
		parts := strings.Split(fullRepoName, "/")
		if len(parts) == 2 {
			organization := parts[0]
			repository := parts[1]

			result := addMonths(time.Now(), -6)

			since := result.Format("2006-01-02T15:04:05Z")
			time.Sleep(2 * time.Millisecond)

			wg.Add(1)
			go getDeveloperWeightsByRepo(organization, repository, since, weights, &wg)
		} else {
			zap.L().Info("invalid repository format", zap.String("repository", fullRepoName))
		}
	}

	go func() {
		wg.Wait()
		close(weights)
	}()

	finalWeights := make(map[string]int)
	for weightMap := range weights {
		for user, weights := range weightMap {
			if weights >= 2 {
				finalWeights[user] = weight
			}
		}
	}
	return finalWeights
}

// getDeveloperWeightsByRepo calculates the weights of developers for a specific repository.
func getDeveloperWeightsByRepo(organization, repository, since string, weights chan<- map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()

	commits, err := getContributorsWithRetry(organization, repository, since, 3, 5*time.Second)
	if err != nil {
		zap.L().Error("failed to get contributors", zap.Error(err))
		return
	}

	weightMap := make(map[string]int)
	for _, v := range commits {
		if v.Committer.User.Login != "" {
			weightMap[v.Committer.User.Login]++
		}
		if v.Author.User.Login != "" {
			weightMap[v.Author.User.Login]++
		}
	}
	weights <- weightMap
}

// getContributors retrieves the commit history of a repository from GitHub GraphQL API.
func getContributors(owner, name, since string) ([]models.Nodes, error) {
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
		zap.L().Error("failed to create  http request", zap.Error(err))
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", config.Client.Github.GithubToken))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "api.github.com")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		zap.L().Error("failed to make  http request", zap.Error(err))
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
		zap.L().Error("failed to unmarshal json response", zap.Error(err))
		return nil, err
	}
	time.Sleep(1000 * time.Millisecond)
	return calculate.Data.Repository.DefaultBranchRef.Target.History.Nodes, nil
}

// getContributorsWithRetry retries fetching contributors for a repository with backoff and retry logic.
func getContributorsWithRetry(organization, repository, since string, maxRetries int, retryInterval time.Duration) ([]models.Nodes, error) {
	var commits []models.Nodes
	var err error

	for i := 0; i < maxRetries; i++ {
		commits, err = getContributors(organization, repository, since)
		if err == nil {
			break
		}
		zap.L().Error(fmt.Sprintf("Retry %d ", i+1), zap.Error(err))
		time.Sleep(retryInterval)
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

// addMerge merges two maps of strings to integers by summing up their values for common keys.
func addMerge(obj1, obj2 map[string]int) map[string]int {
	result := make(map[string]int)

	for key, value := range obj1 {
		result[key] = value
	}

	for key, value := range obj2 {
		if existingValue, ok := result[key]; ok {
			result[key] = existingValue + value
		} else {
			result[key] = value
		}
	}

	return result
}
