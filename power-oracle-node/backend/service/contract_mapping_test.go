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
	"backend/contract"
	"backend/models"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestGetContractMapping(t *testing.T) {
	var err error
	config.InitLogger()
	config.InitConfig("../")
	contract.GoEthClient, err = contract.GetClient(314159)
	if err != nil {
		assert.Error(t, err)
	}

	ethAddress := common.HexToAddress("0x763D410594a24048537990dde6ca81c38CfF566a")

	voterPower, err := GetContractMapping(models.VoterToPower, contract.GoEthClient, []interface{}{ethAddress})
	if err != nil {
		assert.Error(t, err)
	}
	assert.NotNil(t, voterPower)

	voterInfo, err := GetContractMapping(models.VoterToInfo, contract.GoEthClient, []interface{}{ethAddress})
	if err != nil {
		assert.Error(t, err)
	}
	assert.NotNil(t, voterInfo)

	fmt.Println("voter power:", voterPower)
	fmt.Println("voter info:", voterInfo)
}
