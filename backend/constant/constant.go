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

package constant

import "time"

const (
	PowerVotingApiPrefix = "/power_voting/api"
	GistApiPrefix        = "https://api.github.com/gists/"
	// VoteApprove represents the approval vote status.
	VoteApprove = "approve"

	// VoteReject represents the rejection vote status.
	VoteReject = "reject"

	// http request timeout time
	RequestTimeout = time.Second * 15
	MaxFileSize = 1024*2

	// geth The maximum supported event parsing block limit
	SyncBlockLimit = 2880

	// ProposalStatusPending represents the pending proposal status.
	ProposalStatusPending    = 1
	ProposalStatusInProgress = 2
	ProposalStatusCounting   = 3
	ProposalStatusCompleted  = 4

	ProposalCreate  = 0 // proposal created
	ProposalCounted = 1 // proposal counted

	FipProposalRevoke = 0 // fip proposal revoked
	FipProposalUnpass = 0 // fip proposal unpassed
	FipProposalPass   = 1 // fip proposal passed
	FipEditorValid    = 0 // fip editor unremoved
	FipEditorInvalid  = 1 // fip editor removed
	// contract event logs name
	ProposalEvt             = "ProposalCreate"
	VoteEvt                 = "Vote"
	FipCreateEvt            = "FipEditorProposalCreateEvent"
	FipPassedEvt            = "FipEditorProposalPassedEvent"
	FipVoteEvt              = "FipEditorProposalVoteEvent"
	OracleUpdateGistIdsEvt  = "UpdateGistIdsEvent"
	OracleUpdateMinerIdsEvt = "UpdateMinerIdsEvent"
	// mysql duplicate error code
	MysqlDuplicateEntryErrorCode = 1062
)
