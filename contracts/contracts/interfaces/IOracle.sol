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

interface IOracle  {
    /**
     * @notice update a list of miner IDs for a specific voter.
     * @param minerIds List of miner IDs to be added.
     */
    function updateMinerIds(uint64[] memory minerIds) external;

    /**
     * @notice Update a voter's github authorization information
     * @param gistId  github gist id with authorization information 
     */
    function updateGistId(string memory gistId) external;

    /**
     * @notice Event emitted when a gist id is updated for a delegate.
     * @param voterAddress The address of the voter who is update the gist id.
     * @param  gistId github gist id with authorization information
     */
    event UpdateGistIdsEvent(address voterAddress, string gistId);

    /**
     * @notice Event emitted when a miner ID is updated for a delegate.
     * @param voterAddress The address of the voter who is update the miner ID.
     * @param minerIds An array of new miner IDs associated with the delegate.
     */
    event UpdateMinerIdsEvent(address voterAddress, uint64[] minerIds);
}
