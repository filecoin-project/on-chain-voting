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

import "../src/Oracle.sol";
import "../src/interfaces/IOracle.sol";

contract OracleTest {

    Oracle public oracle;

    address public oracleAddress;

    constructor() {
        oracle = new Oracle();
        oracle.initialize();
        oracleAddress = address(oracle);
    }

    function test_update_power_voting_contract() external {
        oracle.updatePowerVotingContract(address(this));
        address powerVotingAddress = oracle.powerVotingContract();
        require(powerVotingAddress == address(this), "update power voting contract error");
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

    function test_remove_voter() external {
        oracle.updateNodeAllowList(address(this), true);

        IOracle.Power memory power = IOracle.Power(100, new bytes[](258), new bytes[](672), 200, block.number);
        oracle.savePower(address(this), power);
        IOracle.Power memory getPower = oracle.getPower(address(this), 1);
        require(getPower.developerPower != 0, "remove voter failed, developerPower error");
        require(getPower.tokenHolderPower != 0, "remove voter failed, tokenHolderPower error");

        oracle.removeVoter(address(this), 1);
        IOracle.Power memory getPower1 = oracle.getPower(address(this), 1);
        require(getPower1.developerPower == 0, "remove voter failed, developerPower error");
        require(getPower1.tokenHolderPower == 0, "remove voter failed, tokenHolderPower error");
    }

    function test_get_power() external {
        oracle.updateNodeAllowList(address(this), true);

        IOracle.Power memory power = IOracle.Power(100, new bytes[](258), new bytes[](672), 200, block.number);
        oracle.savePower(address(this), power);

        IOracle.Power memory getPower = oracle.getPower(address(this), 1);
    }

    function test_update_node_allow_list() external {
        oracle.updateNodeAllowList(address(this), true);
        bool node = oracle.nodeAllowList(address(this));
        require(node, "update node allow list error");
    }

    function test_save_power() external {
        oracle.updateNodeAllowList(address(this), true);
        IOracle.Power memory power = IOracle.Power(100, new bytes[](1579), new bytes[](2578), 200, 1631497);
        oracle.savePower(address(this), power);

        IOracle.Power memory getPower = oracle.getPower(address(this), 1);
        require(getPower.developerPower == 100, "save power failed, developerPower error");
        require(getPower.tokenHolderPower == 200, "save power failed, tokenHolderPower error");
    }
}
