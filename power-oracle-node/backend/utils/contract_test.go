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

package utils

import (
	"backend/config"
	"backend/contract"
	"backend/models"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTasks(t *testing.T) {
	config.InitConfig("../")
	ethClient, err := contract.GetClient(314159)
	if err != nil {
		fmt.Println("get client error: ", err)
		return
	}

	taskList, err := GetTasks(ethClient)
	assert.Nil(t, err)

	testTasks := []*big.Int{}
	assert.Equal(t, taskList, testTasks)

	fmt.Printf("task list: %+v\n", taskList)

}

func TestGetF4Tasks(t *testing.T) {
	config.InitConfig("../")
	ethClient, err := contract.GetClient(314159)
	assert.Error(t, err)

	taskList, err := GetF4Tasks(ethClient)
	assert.Nil(t, err)

	testTasks := []*big.Int{}
	assert.Equal(t, taskList, testTasks)

	fmt.Printf("f4 task list: %+v\n", taskList)

}

func TestGetVoterAddresses(t *testing.T) {
	config.InitConfig("../")
	ethClient, err := contract.GetClient(314159)
	assert.Nil(t, err)

	ethAddressList, err := GetVoterAddresses(ethClient)
	assert.Nil(t, err)

	assert.NotEmpty(t, ethAddressList)
	fmt.Printf("eth address list: %+v\n", ethAddressList)
}

func TestGetVoterInfo(t *testing.T) {
	config.InitConfig("../")
	ethClient, err := contract.GetClient(314159)
	if err != nil {
		fmt.Println("get client error: ", err)
		return
	}

	voterInfo, err := GetVoterInfo("0x763D410594a24048537990dde6ca81c38CfF566a", ethClient)
	assert.Nil(t, err)

	testVoterInfo := models.VoterInfo{
		ActorIds:      []uint64{35363},
		MinerIds:      []uint64{},
		GithubAccount: "",
		EthAddress:    common.HexToAddress("0x763D410594a24048537990dde6ca81c38CfF566a"),
		UcanCid:       "",
	}
	assert.Equal(t, testVoterInfo, voterInfo)
	fmt.Printf("voter info: %+v\n", voterInfo)

}

func TestGetActorIdFromEthAddress(t *testing.T) {
	config.InitConfig("../")
	ethClient, err := contract.GetClient(314159)
	assert.Nil(t, err)

	actorId, err := GetActorIdFromEthAddress("0x763D410594a24048537990dde6ca81c38CfF566a", ethClient)
	assert.Equal(t, "t035363", actorId)
	fmt.Printf("actor id: %+v\n", actorId)
}
