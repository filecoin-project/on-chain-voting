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

import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { UUPSUpgradeable } from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import { IPowerVoting } from "./interfaces/IPowerVoting-filecoin.sol";
import { Proposal, VoteInfo, VoterInfo } from "./types.sol";
import "@openzeppelin/contracts/utils/Counters.sol";


contract PowerVoting is IPowerVoting, Ownable2StepUpgradeable, UUPSUpgradeable {

    using Counters for Counters.Counter;

    // proposal id
    Counters.Counter public proposalId;

    // Power Oracle contract address
    address public oracleContract;

    // add task function selector
    bytes4 public immutable ADD_TASK_SELECTOR = bytes4(keccak256('addTask(string)'));

    // add f4 task function selector
    bytes4 public immutable ADD_F4_TASK_SELECTOR = bytes4(keccak256('addF4Task(address)'));

    // add miner id function selector
    bytes4 public immutable ADD_MINER_IDS_SELECTOR = bytes4(keccak256('addMinerIds(uint64[],address)'));

    // get voter info
    bytes4 public immutable GET_VOTER_INFO_SELECTOR = bytes4(keccak256('getVoterInfo(address)'));

    // proposal mapping, key: proposal id, value: Proposal
    mapping(uint256 => Proposal) public idToProposal;

    // fip map
    mapping(address => bool) public fipMap;

    // proposal id to vote, out key: proposal id, inner key: vote id, value: vote info
    mapping(uint256 => mapping(uint256 => VoteInfo)) public proposalToVote;

    /**
     * @dev Modifier that ensures the provided address is non-zero.
     * @param addr The address to check.
     */
    modifier nonZeroAddress(address addr) {
        if(addr == address(0)){
            revert ZeroAddressError("Zero address error.");
        }
        _;
    }

    /**
     * @notice Authorizes an upgrade to a new implementation contract.
     * @param newImplementation The address of the new implementation contract.
     */
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    /**
     * @notice Initializes the contract by setting up UUPS upgrade ability and ownership.
     */
    function initialize(address oracleAddress) public initializer nonZeroAddress (oracleAddress) {
        oracleContract = oracleAddress;
        __UUPSUpgradeable_init();
        __Ownable_init(msg.sender);
    }

    /**
     * @notice Updates the address of the Oracle contract.
     * @param oracleAddress The new address of the Oracle contract.
     */
    function updateOracleContract(address oracleAddress) external onlyOwner nonZeroAddress (oracleAddress) {
        oracleContract = oracleAddress;
    }

    /**
     * @notice Adds a new FIP address.
     * @param fipAddress The address of the new FIP.
     */
    function addFIP( address fipAddress ) external override onlyOwner nonZeroAddress (fipAddress) {
        bool exist = fipMap[fipAddress];
        // FIP Editor is not allowed to have other roles currently.
        if (exist) {
            revert AddFIPError("Add FIP editor error.");
        }
        fipMap[fipAddress] = true;
    }

    /**
     * @notice Removes the specified FIP address.
     * @param fipAddress The address of the FIP to be removed.
     */
    function removeFIP( address fipAddress ) external override onlyOwner nonZeroAddress(fipAddress) {
        fipMap[fipAddress] = false;
    }

    /**
     * @notice Creates a new proposal.
     * @param proposalCid The CID of the proposal.
     * @param startTime The start time of the proposal.
     * @param expTime The expiration time of the proposal.
     * @param proposalType The type of the proposal.
     */
    function createProposal(string calldata proposalCid, uint248 startTime, uint248 expTime, uint256 proposalType) override external {
        bool fip = fipMap[msg.sender];
        if(!fip){
            revert CallError("Not FIP.");
        }

        // increment proposal id
        proposalId.increment();
        uint256 id = proposalId.current();
        // create proposal
        Proposal storage proposal = idToProposal[id];
        proposal.cid = proposalCid;
        proposal.creator = msg.sender;
        proposal.startTime = startTime;
        proposal.expTime = expTime;
        proposal.proposalType = proposalType;
        emit ProposalCreate(id, proposal);
    }


    /**
     * @notice Voting rights for proposals.
     * @param id The ID of the proposal.
     * @param info Additional information related to the vote.
     */
    function vote(uint256 id, string calldata info) override external {
        Proposal storage proposal = idToProposal[id];
        // if proposal is not start, won't be allowed to vote
        if(proposal.startTime > block.timestamp){
            revert TimeError("Proposal not start yet.");
        }
        // if proposal is expired, won't be allowed to vote
        if(proposal.expTime <= block.timestamp){
            revert TimeError("Proposal expiration time reached.");
        }
        _addF4Task();
        // increment votesCount
        uint256 vid = ++proposal.votesCount;
        // use votesCount as vote id
        VoteInfo storage voteInfo = proposalToVote[id][vid];
        voteInfo.voteInfo = info;
        voteInfo.voter = msg.sender;
        emit Vote(id, msg.sender, info);
    }

    /**
     * @notice Adds miner IDs to the Oracle contract.
     * @param minerIds An array containing the miner IDs to be added.
     */
    function addMinerId(uint64[] memory minerIds) override external {
        _addF4Task();
        (bool addMinerSuccess, ) = oracleContract.call(abi.encodeWithSelector(ADD_MINER_IDS_SELECTOR, minerIds, msg.sender));
        if(!addMinerSuccess){
            revert CallError("Call oracle contract to add miner id failed.");
        }
    }

    /**
     * @notice Adds an F4 task for the caller if necessary.
     * @dev This function is called internally to check whether the caller needs to have an F4 task added.
     */
    function _addF4Task() private {
        (bool getVoterInfoSuccess, bytes memory data) = oracleContract.call(abi.encodeWithSelector(GET_VOTER_INFO_SELECTOR, msg.sender));
        if (!getVoterInfoSuccess) {
            revert CallError("Call oracle contract to get voter info failed.");
        }
        VoterInfo memory voterInfo = abi.decode(data, (VoterInfo));
        if (voterInfo.actorIds.length == 0) {
            (bool addF4TaskSuccess, ) = oracleContract.call(abi.encodeWithSelector(ADD_F4_TASK_SELECTOR, msg.sender));
            if (!addF4TaskSuccess) {
                revert CallError("Call oracle contract to add F4 task failed.");
            }
        }
    }

    /**
     * @notice Delegates the specified UCAN CID to the  Oracle for processing.
     * @param ucanCid The UCAN CID to be delegated.
     */
    function ucanDelegate(string calldata ucanCid) override external {
        // call kyc oracle to add task
        (bool success, ) = oracleContract.call(abi.encodeWithSelector(ADD_TASK_SELECTOR, ucanCid));
        if(!success){
            revert CallError("Call oracle contract to add task failed.");
        }
    }

}
