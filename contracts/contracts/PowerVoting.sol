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

pragma solidity ^0.8.20;

import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IPowerVoting} from "./interfaces/IPowerVoting-Vote.sol";
import {IPowerVotingFipEditor} from "./interfaces/IPowerVoting-Fip.sol";
import {Proposal, VoteInfo, VoterInfo,ProposalEventInfo} from "./types.sol";
contract PowerVoting is IPowerVoting, Ownable2StepUpgradeable, UUPSUpgradeable {

    //fipediotr contract
    IPowerVotingFipEditor public fipEditorContract;

    // proposal id
    uint256 public proposalId;

    // title max length
    uint256 public titleMaxLength;
    // content max length
    uint256 public contentMaxLength;

    // proposal mapping, key: proposal id, value: Proposal
    mapping(uint256 => Proposal) public idToProposal;

    // proposal id to vote, out key: proposal id, inner key: vote id, value: vote info
    mapping(uint256 => mapping(address => string)) public proposalToVote;

    /**
     * @dev A constant value serving as a multiplier to establish the precision of percentage values
     * within the process of proposal creation. This multiplier dictates the number of decimal places
     * that the percentage values can possess.
     * - When `PERCENTAGE_100` is set to 10000, percentage values can have up to two decimal places. 
     *   For example, a value of 5000 stands for 50.00%, and 5025 represents 50.25%.
     * Once the contract is deployed, this value cannot be changed.
     */
    uint16 public constant PERCENTAGE_100 = 10000; 

    uint16 public snapshotMaxRandomOffsetDays;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }
    
    /**
     * @dev Modifier that ensures the provided address is non-zero.
     * @param addr The address to check.
     */
    modifier nonZeroAddress(address addr) {
        if (addr == address(0)) {
            revert ZeroAddressError();
        }
        _;
    }

    /**
     * @notice Authorizes an upgrade to a new implementation contract.
     * @param newImplementation The address of the new implementation contract.
     */
    function _authorizeUpgrade(
        address newImplementation
    ) internal override onlyOwner {}

    /**
     * @notice Initializes the contract by setting up UUPS upgrade ability and ownership.
     */
    function initialize(
        address fipEditorAddress
    ) public initializer nonZeroAddress(fipEditorAddress) {
        //uups must init here 
        proposalId = 0;
        snapshotMaxRandomOffsetDays = 45;
        titleMaxLength = 200;
        contentMaxLength = 10000;
        fipEditorContract = IPowerVotingFipEditor(fipEditorAddress);
        address sender = msg.sender;
        __UUPSUpgradeable_init();
        __Ownable_init(sender);
    }

     modifier onlyFIPEditors() {
        if (!fipEditorContract.isFipEditor(msg.sender)) {
            revert OnlyFipEditorsAllowedError();
        }
        _;
    }
       /**
     * @notice Set the maximum number of random offset days for a snapshot.
     *         The random offset days will be in the range from 1 to the specified value.
     * @dev This function allows the contract administrator to configure the upper limit
     *      of the random offset days used to determine snapshot timestamps.
     * @param _snapshotMaxRandomOffsetDays The maximum number of random offset days for a snapshot.
     *                                     It must be a positive integer.
     */
    function setSnapshotMaxRandomOffsetDays(
        uint16 _snapshotMaxRandomOffsetDays
    ) external  override onlyFIPEditors {
        snapshotMaxRandomOffsetDays = _snapshotMaxRandomOffsetDays;
    }
     /**
    * @notice set proposal content and title max length
     * @param _titleMaxLength  max length of title
     * @param _contentMaxLength max length of content
     */
    function setLengthLimits(uint256 _titleMaxLength, uint256 _contentMaxLength) external override onlyFIPEditors  {
        titleMaxLength = _titleMaxLength;
        contentMaxLength = _contentMaxLength;
    }


  /**
     * @notice Updates the address of the FipEditor contract.
     * @param fipEditorAddress The new address of the FipEditor contract.
     */
    function updateFipEditorContract(address fipEditorAddress) external onlyOwner nonZeroAddress(fipEditorAddress) {
        fipEditorContract = IPowerVotingFipEditor(fipEditorAddress);
    }
    
   /**
     * @notice Creates a new proposal.
     * @param startTime The start time of the proposal. It is expected to be a valid Unix timestamp representing when the proposal officially begins to accept votes or be considered active.
     * @param endTime The expiration time of the proposal. This Unix timestamp indicates when the proposal will no longer be open for voting or further actions.
    *  @param tokenHolderPercentage The token's weight when the proposal is counted. This value should be within 0 - 100 * DECIMAL_MULTIPLIER (representing 0 - 100.00 with precision determined by DECIMAL_MULTIPLIER).
     * @param spPercentage  The sp's weight when the proposal is counted. This value should be within 0 - 100 * DECIMAL_MULTIPLIER (representing 0 - 100.00 with precision determined by DECIMAL_MULTIPLIER).
     * @param clientPercentage The client's weight when the proposal is counted. This value should be within 0 - 100 * DECIMAL_MULTIPLIER (representing 0 - 100.00 with precision determined by DECIMAL_MULTIPLIER).
     * @param developerPercentage The developer's weight when the proposal is counted. This value should be within 0 - 100 * DECIMAL_MULTIPLIER (representing 0 - 100.00 with precision determined by DECIMAL_MULTIPLIER).
     * @param content A detailed description of the proposal. 
     * @param title A short and descriptive title for the proposal. 
     */
    function createProposal(
        uint256 startTime,
        uint256 endTime,
        uint16 tokenHolderPercentage,
        uint16 spPercentage,
        uint16 clientPercentage,
        uint16 developerPercentage,
        string memory content,
        string memory title
    ) external override  onlyFIPEditors{
        ++proposalId;

        //check title and content
        if(bytes(title).length > titleMaxLength){
             revert TitleLengthLimitError(titleMaxLength);
        }
        if(bytes(content).length > contentMaxLength){
             revert ContentLengthLimitError(contentMaxLength);
        }

        //check percentage
        if (
            tokenHolderPercentage > PERCENTAGE_100 ||
            spPercentage > PERCENTAGE_100 ||
            clientPercentage > PERCENTAGE_100 ||
            developerPercentage > PERCENTAGE_100
        ) {
            revert PercentageOutOfRangeError();
        }
        if (
            tokenHolderPercentage +
            spPercentage +
            clientPercentage +
            developerPercentage != PERCENTAGE_100
        ) {
            revert InvalidProposalPercentageError();
        }

        //check time
        if (block.timestamp >= endTime){
            revert InvalidProposalEndTimeError();
        }
        if (startTime >= endTime) {
            revert InvalidProposalTimeError();
        }

        // create proposal
        Proposal memory newProposal = Proposal({
            creator: msg.sender,
            startTime: startTime,
            endTime: endTime
        });

        //store proposal
        idToProposal[proposalId] = newProposal;

        // Define the number of seconds in a day (24 * 60 * 60 = 86400 seconds).
        uint256 oneDayInSeconds = 86400;
        // Get the timestamp of the current block.
        uint256 blockTimestamp = block.timestamp;
        // Generate a pseudo - random number using block data, sender, and proposal ID.
        uint256 randomValue = uint256(keccak256(abi.encodePacked(
            blockhash(block.number - 1),
            blockTimestamp,
            msg.sender,
            proposalId
        )));
        // Calculate a random offset in days (1 to snapshotMaxRandomOffsetDays).
        uint256 randomOffset = randomValue % snapshotMaxRandomOffsetDays + 1;
        // Calculate the timestamp of the previous day.
        uint256 previousDayTimestamp = block.timestamp - oneDayInSeconds;
        // Calculate the passed time since midnight of the previous day.
        uint256 previousDayPassedTime = previousDayTimestamp % oneDayInSeconds;
        // Get the timestamp of midnight of the previous day.
        uint256 snapshotEndDayTimestamp = previousDayTimestamp - previousDayPassedTime;
        //It is possible that snapshotMaxRandomOffsetDays is configured with an excessively large number
        if(randomOffset * oneDayInSeconds >= snapshotEndDayTimestamp ){
            revert SnapshotTimeOutOfRangeError();
        }
        // Calculate the final random snapshot timestamp.
        uint256 snapshotTimestamp = snapshotEndDayTimestamp - randomOffset * oneDayInSeconds;

        // proposal create event
        ProposalEventInfo memory eventInfo = ProposalEventInfo({
            creator: msg.sender,
            startTime: startTime,
            endTime: endTime,
            timestamp: blockTimestamp,
            content: content,
            title: title,
            snapshotTimestamp: snapshotTimestamp,
            tokenHolderPercentage:tokenHolderPercentage,
            spPercentage:spPercentage,
            clientPercentage:clientPercentage,
            developerPercentage:developerPercentage
        });
        emit ProposalCreate(proposalId, eventInfo);
    }

    /**
     * @notice Voting rights for proposals.
     * @param id The ID of the proposal.
     * @param info Additional information related to the vote.
     */
    function vote(uint256 id, string calldata info) external override {

        // info length must more than zero
        if(bytes(info).length == 0){
             revert InvalidVoteInfoError();
        }
        
        Proposal storage proposal = idToProposal[id];
        // if proposal is not start, won't be allowed to vote
        if (proposal.startTime > block.timestamp) {
            revert VotingTimeNotStartedError();
        }
        // if proposal is expired, won't be allowed to vote
        if (proposal.endTime <= block.timestamp) {
            revert VotingAlreadyEndedError();
        }
        //Record the voting results of users. Only the latest voting results are kept
        proposalToVote[id][msg.sender]=info;
        //vote event
        emit Vote(id, msg.sender, info);
    }

  
    
}
