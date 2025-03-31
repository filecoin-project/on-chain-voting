// SPDX-License-Identifier: Apache-2.0
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

pragma solidity ^0.8.20;

import {FipEditorProposalCreateInfo, FipEditorProposalVoteInfo} from "../types.sol";

interface IPowerVotingFipEvent {
    /**
     * @notice Emitted when a new proposal is created.
     * @param proposal The details of the created proposal.
     */
    event FipEditorProposalCreateEvent(FipEditorProposalCreateInfo proposal);

    /**
     * @notice Emitted when  proposal is vote.
     * @param voteInfo The details of the created proposal.
     */
    event FipEditorProposalVoteEvent(FipEditorProposalVoteInfo voteInfo);

    /**
     * @notice Emitted when  proposal is passed.
     * @param proposalId proposal id
     */
    event FipEditorProposalPassedEvent(uint256 proposalId);

}
