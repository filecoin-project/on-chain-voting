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

package db

import (
	"fmt"
	"go.uber.org/zap"
	"powervoting-server/config"
	"powervoting-server/model"
	"testing"
)

func TestVoteResult(t *testing.T) {
	config.InitConfig("../")
	InitMysql()
	var voteHistoryList []model.VoteHistory
	var voteResultList []model.VoteResult
	voteHistoryList = append(voteHistoryList, model.VoteHistory{ProposalId: 1, OptionId: 0, Votes: 100, Address: "0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307"})
	voteHistoryList = append(voteHistoryList, model.VoteHistory{ProposalId: 1, OptionId: 1, Votes: 200, Address: "0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307"})
	voteHistoryList = append(voteHistoryList, model.VoteHistory{ProposalId: 1, OptionId: 2, Votes: 300, Address: "0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307"})
	voteResultList = append(voteResultList, model.VoteResult{ProposalId: 1, OptionId: 0, Votes: 1000})
	voteResultList = append(voteResultList, model.VoteResult{ProposalId: 1, OptionId: 1, Votes: 2000})
	voteResultList = append(voteResultList, model.VoteResult{ProposalId: 1, OptionId: 2, Votes: 3000})

	VoteResult(3, voteHistoryList, voteResultList)
}

func TestGetProposalList(t *testing.T) {
	config.InitConfig("../")
	InitMysql()
	list, err := GetProposalList(314159, 5)
	if err != nil {
		zap.L().Error("get proposal list error: ", zap.Error(err))
		return
	}
	fmt.Println(list)
}

func TestGetVoteList(t *testing.T) {
	config.InitConfig("../")
	InitMysql()
	list, err := GetVoteList(314159, 1)
	if err != nil {
		zap.L().Error("get vote list error: ", zap.Error(err))
		return
	}
	fmt.Println(list)
}
