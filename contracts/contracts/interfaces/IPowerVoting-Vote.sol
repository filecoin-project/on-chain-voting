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

import {IPowerVotingEvent} from "./IPowerVotingEvent.sol";
import {IPowerVotingError} from "./IPowerVotingError.sol";
import {FipEditorProposal} from "../types.sol";

interface IPowerVoting is IPowerVotingEvent, IPowerVotingError {
    /**
     * @notice Set the maximum number of random offset days for a snapshot.
     *         The random offset days will be in the range from 1 to the specified value.
     * @dev This function allows the contract administrator to configure the upper limit
     *      of the random offset days used to determine snapshot timestamps.
     * @param _snapshotMaxRandomOffsetDays The maximum number of random offset days for a snapshot.
     *                                     It must be a positive integer.
     */
    function setSnapshotMaxRandomOffsetDays(
        uint16 _snapshotMaxRandomOffsetDays
    ) external;

    /**
     * @notice set proposal content and title max length
     * @param _titleMaxLength  max length of title
     * @param _contentMaxLength max length of content
     */
    function setLengthLimits(
        uint256 _titleMaxLength,
        uint256 _contentMaxLength
    ) external;
    /**
     * @notice Creates a new proposal.
     * @param startTime The start time of the proposal. It is expected to be a valid Unix timestamp representing when the proposal officially begins to accept votes or be considered active.
     * @param endTime The expiration time of the proposal. This Unix timestamp indicates when the proposal will no longer be open for voting or further actions.
     * @param tokenHolderPercentage The token's weight when the proposal is counted. This value, within 0 - 100.
     * @param spPercentage  The sp's weight when the proposal is counted. This value, within 0 - 100.
     * @param clientPercentage The client's weight when the proposal is counted. This value, within 0 - 100.
     * @param developerPercentage The developer's weight when the proposal is counted. This value, within 0 - 100.
     * @param content A detailed description of the proposal.
     * @param title A short and descriptive title for the proposal.
     */
    function createProposal(
        uint256 startTime,
        uint256 endTime,
        uint16 tokenHolderPercentage,
        uint16 spPercentage,
        uint16 clientPercentage,
        uint16 developerPercentage,
        string memory content,
        string memory title
    ) external;

    /**
     * @notice Voting rights for proposals.
     * @param id The ID of the proposal.
     * @param info Additional information related to the vote.
     */
    function vote(uint256 id, string calldata info) external;

    /**
     * @notice Updates the address of the FipEditor contract.
     * @param fipEditorAddress The new address of the FipEditor contract.
     */
    function updateFipEditorContract(address fipEditorAddress) external;

}
