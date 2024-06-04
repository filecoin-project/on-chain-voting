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
import "forge-std/Test.sol";

contract TestPowerVoting is Test {
    PowerVoting public powerVoting;
    address fipEditorAddressOne = address(0x123);
    address fipEditorAddressTwo = address(0x123456);
    constructor() {
        powerVoting = new PowerVoting();
        powerVoting.initialize(address(this));
    }

    function test_initialization() public view {
        require(powerVoting.oracleContract() == address(this), "Oracle contract address mismatch");
    }


    function test_create_fip_editor_proposal() public {
        string memory voterInfoCid = "test";
        int8 fipEditorProposalType = 1;

        powerVoting.createFipEditorProposal(fipEditorAddressOne, voterInfoCid, fipEditorProposalType);
        require(powerVoting.fipAddressMap(fipEditorAddressOne), "Not a fip address");
    }

    function test_approve_fip_editor() public {
        string memory voterInfoCid = "test";
        int8 fipEditorProposalType = 1;

        powerVoting.createFipEditorProposal(fipEditorAddressOne, voterInfoCid, fipEditorProposalType);

        powerVoting.createFipEditorProposal(fipEditorAddressTwo, voterInfoCid, fipEditorProposalType);
        vm.prank(fipEditorAddressOne);
        powerVoting.approveFipEditor(fipEditorAddressTwo, 2);

        require(powerVoting.fipAddressMap(fipEditorAddressOne), "Not a fip address");
        require(powerVoting.fipAddressMap(fipEditorAddressTwo), "Not a fip address");
    }

    function test_revoke_Fip_editor() public {
        string memory voterInfoCid = "test";

        powerVoting.createFipEditorProposal(fipEditorAddressOne, voterInfoCid, 1);

        powerVoting.createFipEditorProposal(fipEditorAddressTwo, voterInfoCid, 1);
        vm.prank(fipEditorAddressOne);
        powerVoting.approveFipEditor(fipEditorAddressTwo, 2);

        powerVoting.createFipEditorProposal(fipEditorAddressOne, voterInfoCid, 0);
        vm.prank(fipEditorAddressTwo);
        powerVoting.revokeFipEditor(fipEditorAddressOne, 3);
        require(!powerVoting.fipAddressMap(fipEditorAddressOne), "Not a fip address");
    }

    function test_proposal_creation() public {
        string memory proposalCid = "ProposalCID";
        uint248 startTime = uint248(block.timestamp + 60);
        uint248 expTime = uint248(startTime + 3600);
        uint256 proposalType = 1;

        powerVoting.createProposal(proposalCid, startTime, expTime, proposalType);

        (string memory cid, uint256 pType, address creator, uint248 sTime, uint248 eTime, uint256 votesCount) = powerVoting.idToProposal(1);
        assert(keccak256(abi.encodePacked(cid)) == keccak256(abi.encodePacked(proposalCid)));
        assert(creator == address(this));
        assert(sTime == startTime);
        assert(eTime == expTime);
        assert(pType == proposalType);
        assert(votesCount == 0);
    }
}
