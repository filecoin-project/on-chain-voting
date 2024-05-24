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

import "../src/PowerVoting-filecoin.sol";
import "power-oracle-contracts/Oracle.sol";
import "forge-std/Test.sol";

contract TestPowerVoting is Test {
    PowerVoting public powerVoting;
    Oracle public oracle;

    address public oracleAddress;
    address public powerAddress;
    constructor() {
        oracle = new Oracle();
        oracleAddress = address(oracle);

        powerVoting = new PowerVoting();
        powerAddress = address(powerVoting);

        oracle.initialize();
        oracle.updatePowerVotingContract(powerAddress);

        powerVoting.initialize(oracleAddress);
    }

    function test_initialization() public {
        require(powerVoting.oracleContract() == oracleAddress, "Oracle contract address mismatch");
    }

    function test_fip_management() public {
        powerVoting.addFIP(address(this));
        require(powerVoting.fipMap(address(this)) == true, "FIP address not added correctly");

        powerVoting.removeFIP(address(this));
        require(powerVoting.fipMap(address(this)) == false, "FIP address not removed correctly");
    }

    function test_proposal_creation() public {
        powerVoting.addFIP(address(this));

        string memory proposalCid = "ProposalCID";
        uint248 startTime = uint248(block.timestamp + 60);
        uint248 expTime = uint248(startTime + 3600);
        uint256 proposalType = 1;

        powerVoting.createProposal(proposalCid, startTime, expTime, proposalType);

        (string memory cid,uint256 pType,address creator,uint248 sTime,uint248 eTime,uint256 votesCount) = powerVoting.idToProposal(1);
        assert(keccak256(abi.encodePacked(cid)) == keccak256(abi.encodePacked(proposalCid)));
        assert(creator == address(this));
        assert(sTime == startTime);
        assert(eTime == expTime);
        assert(pType == proposalType);
        assert(votesCount == 0);
    }

    function test_voting() public {
        powerVoting.addFIP(address(this));

        string memory proposalCid = "ProposalCID";
        uint248 startTime = uint248(block.timestamp);
        uint248 expTime = uint248(startTime + 1);
        uint256 proposalType = 1;
        powerVoting.createProposal(proposalCid, startTime, expTime, proposalType);

        string memory voteInfo = "I support this proposal!";
        powerVoting.vote(1, voteInfo);

        (string memory info, address voter) = powerVoting.proposalToVote(1, 1);
        require(voter == address(this), "Voter address mismatch");
        require(keccak256(bytes(info)) == keccak256(bytes(voteInfo)), "Vote info mismatch");
    }
}
