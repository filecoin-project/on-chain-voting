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

package contract

import (
	"backend/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClient(t *testing.T) {
	config.InitLogger()
	config.InitConfig("../")
	client, err := GetClient(314159)
	assert.Nil(t, err)
	assert.NotNil(t, client.Client)
	assert.Equal(t, client.Id, int64(314159))
}
