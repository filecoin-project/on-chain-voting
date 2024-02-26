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

import "../PowerVoting-filecoin.sol";

contract PowerVotingTest {

    PowerVoting public powerVotingAddress;

    constructor() {
        powerVotingAddress = new PowerVoting();
    }

    function testCreateProposal() external {
        string memory proposalCid = "bafkreibqn3lahzdjlg4aly7iinc4r7qvgz7hcqm36cbdtezsybuoopsrgm";
        uint248 expTime = 1706667703178;
        uint256 proposalType = 1;
        powerVotingAddress.createProposal(proposalCid, expTime, proposalType);
        (string memory cid, uint256 newType, address creator, uint248 newExpTime, uint256 votesCount) = powerVotingAddress.idToProposal(1);
        require(keccak256(abi.encodePacked(proposalCid)) == keccak256(abi.encodePacked(cid)), "create proposal failed, proposal cid error");
        require(proposalType == newType, "create proposal failed, proposal type error");
        require(creator == address(this), "create proposal failed, proposal creator error");
        require(expTime == newExpTime, "create proposal failed, proposal expTime error");
        require(votesCount == 0, "create proposal failed, proposal votesCount error");
    }

    function testVote() external {
        string memory voteInfo = "bafkreic2hs32eeortzls7utl5bu3yjxieb64k2q3afqn2l7enofeamvqjq";
        powerVotingAddress.vote(1, voteInfo);
        (string memory newVoteInfo, address voter) = powerVotingAddress.proposalToVote(1, 1);
        require(keccak256(abi.encodePacked(voteInfo)) == keccak256(abi.encodePacked(newVoteInfo)), "vote failed, vote info error");
        require(voter == address(this), "vote failed, voter error");
    }

    function testUcanDelegate(address oracle) external {
        powerVotingAddress.updateOracleContract(oracle);
        string memory ucanCid = "bafkreiaitdxatu7ylo47uybagqkvhuutqdplbt2d4oquhilikxhrk42gx4";
        powerVotingAddress.ucanDelegate(ucanCid);
    }

}