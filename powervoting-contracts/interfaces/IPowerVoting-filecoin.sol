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


interface IPowerVoting is IPowerVotingEvent, IPowerVotingError {
    /**
     * @notice Adds a new FIP address.
     * @param fipAddress The address of the new FIP.
     */
    function addFIP(address fipAddress) external;

    /**
     * @notice Removes the specified FIP address.
     * @param fipAddress The address of the FIP to be removed.
     */
    function removeFIP(address fipAddress) external;

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
