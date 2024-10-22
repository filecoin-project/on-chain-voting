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

const (
	// ProposalStartKey is the key used to record the start of proposals.
	ProposalStartKey = "ProposalStartKey"

	// VoteStartKey is the key used to record the start of votes.
	VoteStartKey = "VoteStartKey"

	// VoteApprove represents the approval vote status.
	VoteApprove = 0

	// VoteReject represents the rejection vote status.
	VoteReject = 1

	ProposalStatusStoring    = 0
	ProposalStatusPending    = 1
	ProposalStatusInProgress = 2
	ProposalStatusCounting   = 3
	ProposalStatusCompleted  = 4
	ProposalStatusRejected   = 5
	ProposalStatusPassed     = 6

	Period = 60
)
