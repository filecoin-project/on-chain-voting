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

interface IPowerVotingFipEditor is IPowerVotingFipError {
    /**
     * @notice Creates a new FIP editor proposal.
     * @param fipEditorAddress The address of the FIP editor.
     * @param voterInfo the voter's information.
     * @param fipEditorProposalType The type of FIP editor proposal.
     */
    function createFipEditorProposal(
        address fipEditorAddress,
        string calldata voterInfo,
        int8 fipEditorProposalType
    ) external;

    /**
     * @notice Approves a FIP editor proposal.
     * @param fipEditorAddress The address of the FIP editor.
     * @param proposalId The ID of the proposal.
     */
    function approveFipEditor(
        address fipEditorAddress,
        uint256 proposalId
    ) external;

    /**
     * @notice Revokes a FIP editor proposal.
     * @param fipEditorAddress The address of the FIP editor.
     * @param proposalId The ID of the proposal.
     */
    function revokeFipEditor(
        address fipEditorAddress,
        uint256 proposalId
    ) external;

    /**
     * @notice Gets the list of FIP editor addresses.
     * @return An array containing the addresses of all FIP editors.
     */
    function getFipAddressList() external returns (address[] memory);

    /**
     * @notice Gets the details of a FIP editor proposal based on its ID.
     * @param id The ID of the proposal.
     * @return The details of the proposal.
     */
    function getFipEditorProposal(
        uint256 id
    ) external returns (uint256, address, string memory, address[] memory);

    /**
     * @notice Gets the list of IDs of all approved proposals.
     * @return An array containing the IDs of all approved proposals.
     */
    function getApproveProposalId() external view returns (uint256[] memory);

    /**
     * @notice Gets the list of IDs of all revoked proposals.
     * @return An array containing the IDs of all revoked proposals.
     */
    function getRevokeProposalId() external view returns (uint256[] memory);

      /**
     * @notice Check whether the specified address is FipEditor
     * @param sender: the check address
     * @return  true:fip editor
     */
    function isFipEditor(address sender) external view returns (bool);


    /**
    * @notice  update the length 
     * @param _voterInfoMaxLength  max length of voter info
     */
    function setLengthLimits(uint256 _voterInfoMaxLength) external;
}
