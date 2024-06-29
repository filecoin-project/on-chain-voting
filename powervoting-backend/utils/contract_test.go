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
	"fmt"
	"go.uber.org/zap"
	"powervoting-server/config"
	"powervoting-server/contract"
	"testing"
)

func TestGetProposal(t *testing.T) {
	d, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(d)
	config.InitConfig("../")
	ethClient, err := contract.GetClient(314159)
	if err != nil {
		zap.L().Error("get go-eth client error:", zap.Error(err))
	}
	vote, err := GetProposal(ethClient, 1)
	if err != nil {
		zap.L().Error("get proposal error: ", zap.Error(err))
	}
	fmt.Printf("proposal info: %+v\n", vote)
}

func TestGetVote(t *testing.T) {
	d, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(d)

	config.InitConfig("../")
	ethClient, err := contract.GetClient(314159)
	if err != nil {
		zap.L().Error("get go-eth client error:", zap.Error(err))
	}
	vote, err := GetVote(ethClient, 1, 1)
	if err != nil {
		zap.L().Error("get vote error: ", zap.Error(err))
	}
	fmt.Printf("vote info: %+v\n", vote)
}

func TestGetProposalLatestId(t *testing.T) {
	d, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(d)

	config.InitConfig("../")
	ethClient, err := contract.GetClient(314159)
	if err != nil {
		zap.L().Error("get go-eth client error:", zap.Error(err))
	}
	id, err := GetProposalLatestId(ethClient)
	if err != nil {
		zap.L().Error("get proposal latest id error: ", zap.Error(err))
	}
	fmt.Println("latest id: ", id)
}

func TestGetVoterToPowerStatus(t *testing.T) {
	d, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(d)

	config.InitConfig("../")
	client, err := contract.GetClient(314159)

	if err != nil {
		zap.L().Error("Get client error: ", zap.Error(err))
		return
	}
	voterToPowerStatus, err := GetVoterToPowerStatus("0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307", client)

	if err != nil {
		zap.L().Error("Get VoterToPowerStatus error: ", zap.Error(err))
		return
	}
	fmt.Printf("power: %+v\n", voterToPowerStatus)

}
