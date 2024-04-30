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
	"fmt"
	"testing"
	"time"
)

func TestGetDeveloperWeights(t *testing.T) {
	config.InitLogger()
	config.InitConfig("../")
	totalWeights := GetDeveloperWeights()
	fmt.Println(totalWeights)
}

func TestGetContributors(t *testing.T) {
	config.InitLogger()
	config.InitConfig("../")

	result := addMonths(time.Now(), -6)
	since := result.Format("2006-01-02T15:04:05Z")

	nodes, err := getContributors("filecoin-project", "venus", since)
	if err != nil {
		t.Error(err)
	}
	for _, v := range nodes {
		fmt.Println(v.Committer.User.Login)
		fmt.Println(v.Author.User.Login)

	}
}

func TestAddMonths(t *testing.T) {
	config.InitConfig("../")
	rsp := addMonths(time.Now(), -6)
	fmt.Println(rsp)
}

func TestAddMerge(t *testing.T) {
	map1 := map[string]int{
		"test1": 1,
		"test2": 2,
	}

	map2 := map[string]int{
		"test1": 4,
		"test3": 5,
	}

	rsp := addMerge(map1, map2)
	fmt.Println(rsp)
}
