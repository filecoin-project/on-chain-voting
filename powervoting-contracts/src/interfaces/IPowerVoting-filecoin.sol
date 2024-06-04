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

pragma solidity ^0.8.19;


import {IPowerVotingEvent} from "./IPowerVotingEvent.sol";
import {IPowerVotingError} from "./IPowerVotingError.sol";
import {FipEditorProposal} from "../types.sol";


interface IPowerVoting is IPowerVotingEvent, IPowerVotingError {
    /**
    * @notice Creates a new FIP editor proposal.
    * @param fipEditorAddress The address of the FIP editor.
    * @param voterInfoCid The CID (Content Identifier) of the voter's information.
    * @param fipEditorProposalType The type of FIP editor proposal.
    */
    function createFipEditorProposal(address fipEditorAddress, string calldata voterInfoCid, int8 fipEditorProposalType) external;

    /**
    * @notice Approves a FIP editor proposal.
    * @param fipEditorAddress The address of the FIP editor.
    * @param proposalId The ID of the proposal.
    */
    function approveFipEditor(address fipEditorAddress, uint256 proposalId) external;

    /**
    * @notice Revokes a FIP editor proposal.
    * @param fipEditorAddress The address of the FIP editor.
    * @param proposalId The ID of the proposal.
    */
    function revokeFipEditor(address fipEditorAddress, uint256 proposalId) external;

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
    function getFipEditorProposal(uint256 id) external returns (FipEditorProposal memory);

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
     * @notice Creates a new proposal.
     * @param proposalCid The CID of the proposal.
     * @param startTime The start time of the proposal.
     * @param expTime The expiration time of the proposal.
     * @param proposalType The type of the proposal.
     */
    function createProposal(string calldata proposalCid, uint248 startTime, uint248 expTime, uint256 proposalType) external;

    /**
     * @notice Voting rights for proposals.
     * @param id The ID of the proposal.
     * @param info Additional information related to the vote.
     */
    function vote(uint256 id, string calldata info) external;

    /**
     * @notice Delegates the specified UCAN CID to the  Oracle for processing.
     * @param ucanCid The UCAN CID to be delegated.
     */
    function ucanDelegate(string calldata ucanCid) external;

    /**
     * @notice Updates the address of the Oracle contract.
     * @param oracleAddress The new address of the Oracle contract.
     */
    function updateOracleContract(address oracleAddress) external;

    /**
     * @notice Adds miner IDs to the Oracle contract.
     * @param minerIds An array containing the miner IDs to be added.
     */
    function addMinerId(uint64[] memory minerIds) external;
}
