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

import {IPowerVotingFipError} from "./IPowerVotingError.sol";
import {IPowerVotingFipEvent} from "./IPowerVotingFipEvent.sol";

interface IPowerVotingFipEditor is IPowerVotingFipError, IPowerVotingFipEvent {
    /**
     * create a proposal to add or remove fipeditor
     * @param candidateAddress candidate address
     * @param candidateInfo  candidate info
     * @param fipEditorProposalType fipeditor proposal type:PROPOSAL_TYPE_APPROVE or PROPOSAL_TYPE_REVOKE
     */
    function createFipEditorProposal(
        address candidateAddress,
        string calldata candidateInfo,
        int8 fipEditorProposalType
    ) external;

    /**
     * vote on the proposal
     * @param proposalId proposal
     */
    function voteFipEditorProposal(uint256 proposalId) external;

    /**
     * @notice Check whether the specified address is FipEditor
     * @param sender: the check address
     * @return  true:fip editor
     */
    function isFipEditor(address sender) external view returns (bool);

    /**
     * @notice  update the length
     * @param _candidateInfoMaxLength  max length of candidate info
     */
    function setLengthLimits(uint256 _candidateInfoMaxLength) external;
}
