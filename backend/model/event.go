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

package model

// sync event table
type SyncEventTbl struct {
	Id                         int64  `json:"id"`
	ChainId                    int64  `json:"chain_id" gorm:"not null"`
	ChainName                  string `json:"chain_name" gorm:"not null"`
	PowerVotingContractAddress string `json:"power_voting_contract_address" gorm:"not null;uniqueIndex"`
	FipProposalContractAddress string `json:"fip_proposal_contract_address" gorm:"not null"`
	SyncedHeight               int64  `json:"synced_height" gorm:"not null,default:0"`
}

type GithubRepos struct {
	Id           int64  `json:"id"`
	RepoName     string `json:"repo_name" gorm:"not null;uniqueIndex:repo_name_org_type"`
	OrgId        int64  `json:"org_id" gorm:"not null;uniqueIndex:repo_name_org_type"`
	OrgType      int    `json:"org_type" gorm:"not null"`
	IsNotDeleted int    `json:"is_not_deleted" gorm:"not null"`
	AddTime      int64  `json:"add_time" gorm:"not null"`
}
