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
	"fmt"
	"testing"
)

func TestGetTasks(t *testing.T) {
	config.InitConfig("../")
	ethClient, err := contract.GetClient(314159)
	if err != nil {
		fmt.Println("get client error: ", err)
		return
	}

	taskList, err := GetTasks(ethClient)
	if err != nil {
		fmt.Println("get task list error: ", err)
		return
	}
	fmt.Printf("task list: %+v\n", taskList)

}

func TestGetF4Tasks(t *testing.T) {
	config.InitConfig("../")
	ethClient, err := contract.GetClient(314159)
	if err != nil {
		fmt.Println("get client error: ", err)
		return
	}

	taskList, err := GetF4Tasks(ethClient)
	if err != nil {
		fmt.Println("get f4 task list error: ", err)
		return
	}
	fmt.Printf("f4 task list: %+v\n", taskList)

}

func TestGetVoterAddresses(t *testing.T) {
	config.InitConfig("../")
	ethClient, err := contract.GetClient(314159)
	if err != nil {
		fmt.Println("get client error: ", err)
		return
	}

	ethAddressList, err := GetVoterAddresses(ethClient)
	if err != nil {
		fmt.Println("get eth address list error: ", err)
		return
	}
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
	if err != nil {
		fmt.Println("get voter info error: ", err)
		return
	}
	fmt.Printf("voter info: %+v\n", voterInfo)

}

func TestGetActorIdFromEthAddress(t *testing.T) {
	config.InitConfig("../")
	ethClient, err := contract.GetClient(314159)
	if err != nil {
		fmt.Println("get client error: ", err)
		return
	}

	actorId, err := GetActorIdFromEthAddress("0x763D410594a24048537990dde6ca81c38CfF566a", ethClient)
	if err != nil {
		fmt.Println("get actor id error: ", err)
		return
	}
	fmt.Printf("actor id: %+v\n", actorId)

}
