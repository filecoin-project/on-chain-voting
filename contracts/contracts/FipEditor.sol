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
import {IPowerVotingFipEditor} from "./interfaces/IPowerVoting-Fip.sol";
import {FipEditorProposal, HasVoted, FipEditorProposalCreateInfo, FipEditorProposalVoteInfo} from "./types.sol";

contract PowerVotingFipEditor is
    Ownable2StepUpgradeable,
    UUPSUpgradeable,
    IPowerVotingFipEditor
{
    // Define constants for proposal types
    int8 public constant PROPOSAL_TYPE_APPROVE = 1;
    int8 public constant PROPOSAL_TYPE_REVOKE = 0;
    //fipeditor status
    int8 public constant FITEITOR_STATUS_REVOKED = 0; //default status
    int8 public constant FITEITOR_STATUS_ADDING = 1;
    int8 public constant FITEITOR_STATUS_APPROVED = 2;
    int8 public constant FITEITOR_STATUS_REVOKING = 3;
    // Use EnumerableSet library to handle integers
    using EnumerableSet for EnumerableSet.UintSet;
    using EnumerableSet for EnumerableSet.AddressSet;

    uint256 public candidateInfoMaxLength;
    mapping(address => int8) public fipEditorStatusMap;
    // fip editor proposal id
    uint256 public fipEditorProposalId;
    uint256 public fipEditorCount;
    // fip editor proposal mapping, key: fip editor proposal id, value: fip editor proposal
    mapping(uint256 => FipEditorProposal) idToFipEditorProposal;
    // Set to store IDs of  proposals id
    EnumerableSet.UintSet proposalIdSet;
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }
    /**
     * @notice Initializes the contract by setting up UUPS upgrade ability and ownership.
     */
    function initialize() public initializer {
        address sender = msg.sender;
        fipEditorCount = 1;
        fipEditorStatusMap[sender] = FITEITOR_STATUS_APPROVED;
        candidateInfoMaxLength = 1000;
        __UUPSUpgradeable_init();
        __Ownable_init(sender);
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
     * @dev Modifier to allow only FIP editors to call the function.
     * Reverts with a custom error message if the caller is not a FIP editor.
     */
    modifier onlyFIPEditors() {
        if (
            fipEditorStatusMap[msg.sender] != FITEITOR_STATUS_APPROVED &&
            fipEditorStatusMap[msg.sender] != FITEITOR_STATUS_REVOKING
        ) {
            revert OnlyFipEditorsAllowedError();
        }
        _;
    }

    /**
     * @notice  update the length
     * @param _candidateInfoMaxLength  max length of candidate info
     */
    function setLengthLimits(
        uint256 _candidateInfoMaxLength
    ) external override onlyFIPEditors {
        candidateInfoMaxLength = _candidateInfoMaxLength;
    }

    /**
     * create a proposal to add or remove fipeditor
     * @param candidateAddress candidate address
     * @param candidateInfo  candidate info
     * @param fipEditorProposalType fipeditor proposal type:PROPOSAL_TYPE_APPROVE or PROPOSAL_TYPE_REVOKE
     */
    function createFipEditorProposal(
        address candidateAddress,
        string calldata candidateInfo,
        int8 fipEditorProposalType
    ) external override onlyFIPEditors nonZeroAddress(candidateAddress) {
        // Ensure the proposal type is valid (must be 1 for approval or 0 for revocation)
        if (
            fipEditorProposalType != PROPOSAL_TYPE_APPROVE &&
            fipEditorProposalType != PROPOSAL_TYPE_REVOKE
        ) {
            revert InvalidProposalTypeError();
        }
        // Ensure the proposer is not proposing themselves
        if (candidateAddress == msg.sender) {
            revert CannotProposeToSelfError();
        }
        if (bytes(candidateInfo).length > candidateInfoMaxLength) {
            revert CandidateInfoLimitError(candidateInfoMaxLength);
        }

        int8 candidateFipEditorStatus = fipEditorStatusMap[candidateAddress];
        // Check if the address is already a FIP editor and the proposal is not to revoke (0)
        if (
            candidateFipEditorStatus == FITEITOR_STATUS_APPROVED &&
            fipEditorProposalType != PROPOSAL_TYPE_REVOKE
        ) {
            revert AddressIsAlreadyFipEditorError();
        }

        // Check if the address already has an active proposal
        if (
            candidateFipEditorStatus == FITEITOR_STATUS_ADDING ||
            candidateFipEditorStatus == FITEITOR_STATUS_REVOKING
        ) {
            revert AddressHasActiveProposalError();
        }

         // If the address is not an FIP editor, it cannot be revoked
        if (
            candidateFipEditorStatus != FITEITOR_STATUS_APPROVED &&
            fipEditorProposalType == PROPOSAL_TYPE_REVOKE
        ) {
            revert AddressNotFipEditorError();
        }

        // There must be at minimum two votes to revoke a FIP Editor.
        if (
            fipEditorProposalType == PROPOSAL_TYPE_REVOKE && fipEditorCount <= 2
        ) {
            revert InsufficientEditorsError();
        }

        ++fipEditorProposalId;

        //update FipEditor status
        if (fipEditorProposalType == PROPOSAL_TYPE_APPROVE) {
            fipEditorStatusMap[candidateAddress] = FITEITOR_STATUS_ADDING;
        } else {
            fipEditorStatusMap[candidateAddress] = FITEITOR_STATUS_REVOKING;
        }
        //keep active proposal
        proposalIdSet.add(fipEditorProposalId);
        // Create a new proposal and store it in the mapping
        FipEditorProposal storage proposal = idToFipEditorProposal[
            fipEditorProposalId
        ];
        proposal.candidateAddress = candidateAddress;
        proposal.proposalId = fipEditorProposalId;
        proposal.proposalType = fipEditorProposalType;

        //emit create event
        FipEditorProposalCreateInfo
            memory eventInfo = FipEditorProposalCreateInfo({
                proposalId: fipEditorProposalId,
                proposalType: fipEditorProposalType,
                creator: msg.sender,
                candidateInfo: candidateInfo,
                candidateAddress: candidateAddress
            });

        emit FipEditorProposalCreateEvent(eventInfo);

        // creating a proposal defaults to voting on the proposal
        voteFipEditorProposal(fipEditorProposalId);
    }

    /**
     * vote on the proposal
     * @param proposalId proposal
     */
    function voteFipEditorProposal(
        uint256 proposalId
    ) public override onlyFIPEditors {
        if (!proposalIdSet.contains(proposalId)) {
            revert InvalidApprovalProposalId();
        }
        FipEditorProposal storage proposal = idToFipEditorProposal[proposalId];
        // Check if the sender has already voted for this proposal
        if (proposal.votedAddress[msg.sender]) {
            revert AddressHasActiveProposalError();
        }
        // Ensure the voter is not voting on their own proposal
        if (proposal.candidateAddress == msg.sender) {
            revert CannotVoteForOwnProposalError();
        }

        proposal.voters.add(msg.sender);
        proposal.votedAddress[msg.sender] = true;

        FipEditorProposalVoteInfo memory voteInfo = FipEditorProposalVoteInfo({
            voter: msg.sender,
            proposalId: proposalId
        });
        emit FipEditorProposalVoteEvent(voteInfo);

        //check proposal result
        _checkProposalResult(proposalId);
    }

    /**
     * Check the status of the proposal and calculate the result
     * @param proposalId proposal id
     */
    function _checkProposalResult(uint256 proposalId) private {
        FipEditorProposal storage proposal = idToFipEditorProposal[proposalId];
        uint256 voterCount = proposal.voters.length();
        if (
            proposal.proposalType == PROPOSAL_TYPE_REVOKE &&
            voterCount == fipEditorCount - 1 &&
            fipEditorCount > 2
        ) {
            address candidateAddress = proposal.candidateAddress;
            //update fipeditor count
            fipEditorCount--;
            //will delete proposal
            _processPassedProposal(
                proposalId,
                proposal.candidateAddress,
                FITEITOR_STATUS_REVOKED
            );
            //to pass a revoke proposal, need to review all proposals
            _cleanupProposals(candidateAddress);
        } else if (
            //everyone agreed
            proposal.proposalType == PROPOSAL_TYPE_APPROVE &&
            voterCount == fipEditorCount
        ) {
            //update fipeditor count
            fipEditorCount++;

            _processPassedProposal(
                proposalId,
                proposal.candidateAddress,
                FITEITOR_STATUS_APPROVED
            );
        }
    }

    /**
     * After the proposal is passed, the corresponding status needs to be updated
     * @param proposalId proposal id
     * @param candidateAddress candidate address
     * @param status update status
     */
    function _processPassedProposal(
        uint256 proposalId,
        address candidateAddress,
        int8 status
    ) private {
        fipEditorStatusMap[candidateAddress] = status;
        proposalIdSet.remove(proposalId);
        delete idToFipEditorProposal[proposalId];
        emit FipEditorProposalPassedEvent(proposalId);
    }

    /**
     * After the fipeditor is revoked,  needs to clean up all his votes
     * @param candidateAddress address of the candidate
     */
    function _cleanupProposals(address candidateAddress) private {
        uint256[] memory proposalIds = proposalIdSet.values();
        for (uint256 i = 0; i < proposalIds.length; i++) {
            uint256 id = proposalIds[i];
            FipEditorProposal storage proposal = idToFipEditorProposal[id];
            if (proposal.voters.contains(candidateAddress)) {
                proposal.votedAddress[msg.sender] = false;
                proposal.voters.remove(candidateAddress);
            }
            //After the voter or fipeditor changes, the proposal results should be rechecked
            _checkProposalResult(id);
        }
    }

    function isFipEditor(address sender) external view override returns (bool) {
        return
            fipEditorStatusMap[sender] == FITEITOR_STATUS_APPROVED ||
            fipEditorStatusMap[sender] == FITEITOR_STATUS_REVOKING;
    }
}
