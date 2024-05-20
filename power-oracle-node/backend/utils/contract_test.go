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

	testList := []common.Address{
		common.HexToAddress("0xd064B424f265EAfFf6400cEaBF11bB979F0b484d"),
		common.HexToAddress("0x0000000000000000000000000000000000000000"),
		common.HexToAddress("0x5315eDfd9cF69a46d382E01c5F6fD0f533881d9c"),
		common.HexToAddress("0xD1584862B753be094771523A0584308C608b3D7B"),
		common.HexToAddress("0x7C24ca5FBA6f1E228a520911476746D25Be5EdbE"),
		common.HexToAddress("0xdD52CA4bE75B7B89f30569aF55C3F8361E8f431c"),
		common.HexToAddress("0xAD7B311Cc1cfa104B2Cbd7F9a90b6520EF79cdC5"),
		common.HexToAddress("0xcDBaDc54727976faD12b3FAbce5C3C57629DF96e"),
		common.HexToAddress("0x696FEf6cd9D2c243A607cc3ba055bdEfc9464a41"),
		common.HexToAddress("0xf58cC34cf80BDF9D3aD82E7AC57aCd02cA592193"),
		common.HexToAddress("0x7652b16C9386290906e1FFC8FDC6346D1eEB76A3"),
		common.HexToAddress("0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307"),
		common.HexToAddress("0x85D4e31D5cD7D6dEFE6db9945F20b61a179b1949"),
		common.HexToAddress("0xe95C3DBbb10583B0524f2619BA2FBB51a9FA0249"),
		common.HexToAddress("0x763D410594a24048537990dde6ca81c38CfF566a"),
		common.HexToAddress("0x31c0600B18b8Fe9e5BF3F112205d36fE4fbCc552"),
		common.HexToAddress("0xe4c7b2bb1d600bCD0A9af60dda3874e369C37bc4"),
		common.HexToAddress("0x4fda4174D5D07C906395bfB77806287cc65Fd129"),
	}

	assert.Equal(t, testList, ethAddressList)
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
