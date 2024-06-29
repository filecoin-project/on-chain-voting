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
import "./IOracleError.sol";

interface IOracle is IOracleError {

    // power struct
    struct Power {
        uint256 developerPower;
        bytes[] spPower;
        bytes[] clientPower;
        uint256 tokenHolderPower;
        uint256 blockHeight;
    }

    // voter info
    struct VoterInfo {
        uint64[] actorIds;
        uint64[] minerIds;
        string githubAccount;
        address ethAddress;
        string ucanCid;
    }

    // power status
    struct PowerStatus{
        uint256 dayId;
        uint256 hasFullRound;
    }

    /**
     * @notice Updates the PowerVoting contract address.
     * @param powerVotingAddress The address of the PowerVoting contract.
     */
    function updatePowerVotingContract(address powerVotingAddress) external;

    /**
     * @notice Adds a list of miner IDs for a specific voter.
     * @param minerIds List of miner IDs to be added.
     * @param voter Address of the voter.
     */
    function addMinerIds(uint64[] memory minerIds, address voter) external;

    /**
     * @notice Adds a task with the specified UCAN CID.
     * @param ucanCid The UCAN CID associated with the task.
     */
    function addTask(string calldata ucanCid) external;

    /**
     * @notice Retrieves the list of task IDs.
     * @return Task[] An array containing the task IDs.
     */
    function getTasks() external view returns(uint256[] memory);

    /**
     * @notice Adds an F4 task for the specified voter.
     * @param voter The address of the voter.
     */
    function addF4Task(address voter) external;

    /**
     * @notice Retrieves the list of F4 task IDs.
     * @return uint256[] An array containing the F4 task IDs.
     */
    function getF4Tasks() external view returns(uint256[] memory);

    /**
     * @notice Callback function for updating task information.
     * @param voterInfoParam Voter information containing the Ethereum address and other details.
     * @param taskId The ID of the task being updated.
     */
    function taskCallback(VoterInfo calldata voterInfoParam, uint256 taskId) external;

    /**
     * @notice Updates the node allow list by adding or removing a node.
     * @param nodeAddress Address of the node to be added or removed.
     * @param allow Boolean indicating whether to allow or disallow the node.
     */
    function updateNodeAllowList(address nodeAddress, bool allow) external;

    /**
     * @notice Removes a voter along with associated task information.
     * @param voterAddress Address of the voter to be removed.
     */
    function removeVoter(address voterAddress, uint256 taskId) external;

    /**
     * @notice Retrieves the list of voter addresses.
     * @return address[] An array containing the addresses of all voters.
     */
    function getVoterAddresses() external view returns(address[] memory);

    /**
     * @notice Retrieves the information associated with a specific voter.
     * @param voter The address of the voter.
     */
    function getVoterInfo(address voter) external view returns(VoterInfo memory);

    /**
     * @notice Event emitted when a delegate is created.
     * @param voterAddress The address of the voter who is delegating.
     * @param actorIds An array of actor IDs associated with the delegate.
     * @param github The GitHub username of the delegate.
     */
    event CreateDelegate(address voterAddress,uint64[] actorIds,string github);

    /**
     * @notice Event emitted when a delegate is deleted.
     * @param voterAddress The address of the voter who is removing the delegate.
     * @param actorIds An array of actor IDs that are being removed.
     * @param minerIds An array of miner IDs associated with the delegate.
     * @param github The GitHub username of the delegate.
     */
    event DeleteDelegate(address voterAddress,uint64[] actorIds,uint64[] minerIds,string github);

    /**
     * @notice Event emitted when a miner ID is updated for a delegate.
     * @param voterAddress The address of the voter who is updating the miner ID.
     * @param minerIds An array of new miner IDs associated with the delegate.
     */
    event UpdateMinerId(address voterAddress,uint64[] minerIds);
}
