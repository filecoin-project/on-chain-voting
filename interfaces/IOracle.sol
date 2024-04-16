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
        uint256 hourId;
        uint256 hasFullRound;
    }

    /**
     * updatePowerVotingContract: update powerVoting contract
     * @param powerVotingAddress: powerVoting contract address
     */
    function updatePowerVotingContract(address powerVotingAddress) external;

    /**
     * addMinerIds: add miner id list
     * @param minerIds: miner id list
     * @param voter: voter address
     */
    function addMinerIds(uint64[] memory minerIds, address voter) external;

    /**
     * addTask: add task
     * @param ucanCid: ucan cid
     */
    function addTask(string calldata ucanCid) external;

    /**
     * getTasks: get task list
     * @return Task[]: task list
     */
    function getTasks() external view returns(uint256[] memory);

    /**
     * addF4Task: add f4 task
     * @param voter: voter
     */
    function addF4Task(address voter) external;

    /**
     * getF4Tasks: get f4 task id list
     * @return uint256[]: task id list
     */
    function getF4Tasks() external view returns(uint256[] memory);

    /**
     * taskCallback: task callback function
     * @param voterInfoParam: voter info
     * @param taskId: task id
     * @param power: power
     */
    function taskCallback(VoterInfo calldata voterInfoParam, uint256 taskId, Power calldata power) external;

    /**
     * getPower: get voting power
     * @param voterAddress: voter address
     * @param id: id
     */
    function getPower(address voterAddress, uint256 id) external returns(Power memory);

    /**
     * updateNodeAllowList: update node allowlist
     * @param nodeAddress: node address
     * @param allow:
     */
    function updateNodeAllowList(address nodeAddress, bool allow) external;

    /**
     * removeVoter: remove voter
     * @param voterAddress: voter address
     */
    function removeVoter(address voterAddress, uint256 taskId) external;

    /**
     * getVoterAddresses: get voter list
     * @return address[]: voter address list
     */
    function getVoterAddresses() external view returns(address[] memory);

    /**
     * getVoterInfo: get voter info
     * @param voter: voter address
     */
    function getVoterInfo(address voter) external view returns(VoterInfo memory);

    /**
     * savePower: save voter power
     * @param voterAddress: voter address
     * @param powerParam: power
     */
    function savePower(address voterAddress, Power calldata powerParam) external;

}
