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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetDeveloperWeights(t *testing.T) {
	config.InitLogger()
	config.InitConfig("../")
	totalWeights := GetDeveloperWeights()
	assert.NotEmpty(t, totalWeights, 1)
}

func TestGetContributors(t *testing.T) {
	config.InitLogger()
	config.InitConfig("../")

	result := addMonths(time.Now(), -6)
	since := result.Format("2006-01-02T15:04:05Z")

	nodes, err := getContributors("filecoin-project", "venus", since)
	assert.Nil(t, err)

	assert.NotEmpty(t, nodes)
}

func TestAddMonths(t *testing.T) {
	config.InitConfig("../")
	rsp := addMonths(time.Date(2024, 05, 04, 00, 00, 00, 00, time.Local),
		-2)
	res := time.Date(2024, 03, 04, 00, 00, 00, 00, time.Local)

	assert.Equal(t, rsp, res)
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

	result := map[string]int{
		"test1": 5,
		"test2": 2,
		"test3": 5,
	}
	rsp := addMerge(map1, map2)
	assert.Equal(t, rsp, result)
}
