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

import "../Oracle.sol";
import "../interfaces/IOracle.sol";

contract OracleTest {

    Oracle public oracle;

    address public oracleAddress;

    constructor() {
        oracle = new Oracle();
        oracle.initialize();
        oracleAddress = address(oracle);
    }

    function test_update_power_voting_contract() external {
        oracle.updatePowerVotingContract(0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307);
        address powerVotingAddress = oracle.powerVotingContract();
        require(powerVotingAddress == 0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307, "update power voting contract error");
    }

    function test_add_task() external {
        oracle.updatePowerVotingContract(address(this));
        string memory ucanCid = "bafkreibqn3lahzdjlg4aly7iinc4r7qvgz7hcqm36cbdtezsybuoopsrgm";
        oracle.addTask(ucanCid);
        uint256[] memory taskIdList = oracle.getTasks();
        require(taskIdList[0] == 1, "add task failed, task id error");
        string memory ucan = oracle.taskIdToUcanCid(1);
        require(keccak256(abi.encodePacked(ucanCid)) == keccak256(abi.encodePacked(ucan)), "add task failed, ucan cid error");
    }

    function test_task_callback() external {
        oracle.updateNodeAllowList(address(this), true);
        IOracle.VoterInfo memory voterInfo = IOracle.VoterInfo(new uint64[](35150), new uint64[](0), "github", address(0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307), "ucan cid");
        IOracle.Power memory power = IOracle.Power(3,new bytes[](0),new bytes[](0),200,1631497);
        oracle.taskCallback(voterInfo, 1, power);

        (string memory githubAccount, address  ethAddress, string memory ucanCid) = oracle.voterToInfo(address(0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307));
        require(ethAddress == 0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307, "task callback failed, ethAddress error");
        require(keccak256(abi.encodePacked(githubAccount)) == keccak256(abi.encodePacked("github")), "task callback failed, github error");
        require(keccak256(abi.encodePacked(ucanCid)) == keccak256(abi.encodePacked("ucan cid")), "task callback failed, ucan id error");

        IOracle.Power memory getPower = oracle.getPower(0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307, 1);
        require(getPower.developerPower == 3, "task callback failed, developerPower error");
        require(getPower.tokenHolderPower == 200, "task callback failed, tokenHolderPower error");


        uint256 blockHeight = oracle.voterAddressToBlockHeight(0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307);
        require(blockHeight != 0, "task callback failed, block height error");

        address[] memory voterList = oracle.getVoterAddresses();
        require(voterList[0] == 0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307, "task callback failed, voter list error");

        string memory ucanCID = oracle.taskIdToUcanCid(1);
        require(bytes(ucanCID).length == 0, "task callback failed, taskId to ucan cid error");

        uint256[] memory taskIdList = oracle.getTasks();
        require(taskIdList.length == 0, "task callback failed, taskId list error");
    }

    function testRemoveVoter() external {
        oracle.updateNodeAllowList(address(this), true);
        oracle.removeVoter(0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307, 1);

        (string memory githubAccount, address  ethAddress, string memory ucanCid) = oracle.voterToInfo(address(0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307));
        require(ethAddress == address(0), "remove voter failed, ethAddress error");
        require(bytes(githubAccount).length == 0, "remove voter failed, github error");
        require(bytes(ucanCid).length == 0, "remove voter failed, ucan id id error");

        IOracle.Power memory getPower = oracle.getPower(0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307, 1);
        require(getPower.developerPower == 0, "remove voter failed, developerPower error");
        require(getPower.tokenHolderPower == 0, "remove voter failed, tokenHolderPower error");
    }

    function testGetPower() external view returns(uint256, uint256) {
        IOracle.Power memory getPower = oracle.getPower(0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307, 1);
        return (getPower.developerPower,getPower.tokenHolderPower);
    }


    function testUpdateNodeAllowList() external {
        oracle.updateNodeAllowList(address(this), true);
        bool node = oracle.nodeAllowList(address(this));
        require(node, "update node allow list error");
    }

    function testSavePower() external {
        IOracle.Power memory power = IOracle.Power(100,new bytes[](0),new bytes[](0),200,1631497);
        oracle.savePower(0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307, power);
        IOracle.Power memory getPower = oracle.getPower(0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307, 1);
        require(getPower.developerPower == 3, "save power failed, developerPower error");
        require(getPower.tokenHolderPower == 200, "save power failed, tokenHolderPower error");
    }


}
