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
import {FipEditorProposal, HasVoted} from "./types.sol";
contract PowerVotingFipEditor is
    Ownable2StepUpgradeable,
    UUPSUpgradeable,
    IPowerVotingFipEditor
{
    // Define constants for proposal types
    int8 public constant APPROVE_PROPOSAL_TYPE = 1;

    int8 public constant REVOKE_PROPOSAL_TYPE = 0;

    // fip editor proposal id
    uint256 public fipEditorProposalId;

    // fip editor address mapping, key: address, value: boolean
    mapping(address => bool) public fipAddressMap;

    // fip editor proposal mapping, key: fip editor proposal id, value: fip editor proposal
    mapping(uint256 => FipEditorProposal) private idToFipEditorProposal;
    // Use EnumerableSet library to handle  addresses
    using EnumerableSet for EnumerableSet.AddressSet;

    // Use EnumerableSet library to handle integers
    using EnumerableSet for EnumerableSet.UintSet;

    // Set to store addresses of FIP editors
    EnumerableSet.AddressSet fipAddressList;

    // Set to store IDs of approved proposals id
    EnumerableSet.UintSet approveProposalId;

    // Set to store IDs of revoked proposals id
    EnumerableSet.UintSet revokeProposalId;

    // proposal id to vote status, outer key: proposal id, inner key: address, value: HasVoted
    mapping(uint256 => HasVoted) private idToHasVoted;

    // active proposal mapping, key: address, value: boolean
    mapping(address => bool) public hasActiveProposal;

    uint256 public voterInfoMaxLength;

    /**
     * @notice Initializes the contract by setting up UUPS upgrade ability and ownership.
     */
    function initialize() public initializer {
        address sender = msg.sender;
        fipAddressMap[sender] = true;
        fipAddressList.add(sender);
        voterInfoMaxLength = 1000;
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
        if (!fipAddressMap[msg.sender]) {
            revert OnlyFipEditorsAllowedError();
        }
        _;
    }

    /**
    * @notice  update the length 
     * @param _voterInfoMaxLength  max length of voter info
     */
    function setLengthLimits(
        uint256 _voterInfoMaxLength
    ) external override onlyFIPEditors {
        voterInfoMaxLength = _voterInfoMaxLength;
    }

    /**
     * @notice Creates a proposal to add or remove a FIP editor.
     * @param fipEditorAddress The address of the proposed FIP editor.
     * @param voterInfo  the voter's information.
     * @param fipEditorProposalType The type of proposal: 1 for adding a FIP editor, 0 for removing.
     */
    function createFipEditorProposal(
        address fipEditorAddress,
        string calldata voterInfo,
        int8 fipEditorProposalType
    ) external override onlyFIPEditors nonZeroAddress(fipEditorAddress) {
        
        if (bytes(voterInfo).length > voterInfoMaxLength) {
            revert VoterInfoLimitError(voterInfoMaxLength);
        }

        // Check if the address is already a FIP editor and the proposal is not to revoke (0)
        if (
            fipAddressMap[fipEditorAddress] &&
            fipEditorProposalType != REVOKE_PROPOSAL_TYPE
        ) {
            revert AddressIsAlreadyFipEditorError();
        }

        // Check if the address already has an active proposal
        if (hasActiveProposal[fipEditorAddress]) {
            revert AddressHasActiveProposalError();
        }

        // Ensure the proposer is not proposing themselves
        if (fipEditorAddress == msg.sender) {
            revert CannotProposeToSelfError();
        }

        // Ensure the proposal type is valid (must be 1 for approval or 0 for revocation)
        if (
            fipEditorProposalType != APPROVE_PROPOSAL_TYPE &&
            fipEditorProposalType != REVOKE_PROPOSAL_TYPE
        ) {
            revert InvalidProposalTypeError();
        }

        // There must be at minimum two votes to revoke a FIP Editor.
        if (
            fipEditorProposalType == REVOKE_PROPOSAL_TYPE &&
            fipAddressList.length() <= 2
        ) {
            revert InsufficientEditorsError();
        }

        // Increment the global proposal ID counter
        ++fipEditorProposalId;

        // Create a new proposal and store it in the mapping
        FipEditorProposal storage proposal = idToFipEditorProposal[
            fipEditorProposalId
        ];
        proposal.fipEditorAddress = fipEditorAddress;
        proposal.voterInfo = voterInfo;
        proposal.proposalId = fipEditorProposalId;

        // Mark the address as having an active proposal
        hasActiveProposal[fipEditorAddress] = true;

        // If the proposal is to approve (1), add to the approval set and call approve function
        if (fipEditorProposalType == APPROVE_PROPOSAL_TYPE) {
            approveProposalId.add(fipEditorProposalId);
            approveFipEditor(fipEditorAddress, fipEditorProposalId);
        }

        // If the proposal is to revoke (0), add to the revocation set and call revoke function
        if (fipEditorProposalType == REVOKE_PROPOSAL_TYPE) {
            revokeProposalId.add(fipEditorProposalId);
            revokeFipEditor(fipEditorAddress, fipEditorProposalId);
        }
    }
    /**
     * @notice Approves the proposal to add a new FIP editor.
     * @param fipEditorAddress The address of the proposed FIP editor.
     * @param id The ID of the proposal.
     * Only FIP editors can call this function. The caller must not have already voted for this proposal.
     */
    function approveFipEditor(
        address fipEditorAddress,
        uint256 id
    ) public override onlyFIPEditors {
        // Fetch the voting record for the proposal
        HasVoted storage hasVoted = idToHasVoted[id];

        // Check if the sender has already voted for this proposal
        if (hasVoted.hasVotedAddress[msg.sender]) {
            revert AddressHasActiveProposalError();
        }

        // Ensure the proposal ID is valid for approval
        if (!approveProposalId.contains(id)) {
            //"The proposal ID is not valid for approval"
            revert InvalidApprovalProposalId();
        }

        // Fetch the proposal details
        FipEditorProposal storage proposal = idToFipEditorProposal[id];

        // Add the sender to the list of voters for this proposal
        proposal.voters.add(msg.sender);
        hasVoted.hasVotedAddress[msg.sender] = true;

        // Check if all FIP editors have voted
        if (proposal.voters.length() == fipAddressList.length()) {
            _finalizeProposal(fipEditorAddress, id, true);
        }
    }

    /**
     * @notice Revokes the proposal to remove a FIP editor.
     * @param fipEditorAddress The address of the FIP editor to be removed.
     * @param id The ID of the proposal.
     */
    function revokeFipEditor(
        address fipEditorAddress,
        uint256 id
    ) public override onlyFIPEditors {
        // Fetch the voting record for the proposal
        HasVoted storage hasVoted = idToHasVoted[id];

        // Check if the sender has already voted for this proposal
        if (hasVoted.hasVotedAddress[msg.sender]) {
            revert AddressHasActiveProposalError();
        }

        // Ensure the proposal ID is valid for revocation
        if (!revokeProposalId.contains(id)) {
            revert InvalidRevocationProposalId();
        }

        // Ensure the voter is not voting on their own proposal
        if (fipEditorAddress == msg.sender) {
            revert CannotVoteForOwnProposalError();
        }

        // Fetch the proposal details
        FipEditorProposal storage proposal = idToFipEditorProposal[id];

        // Add the sender to the list of voters for this proposal
        proposal.voters.add(msg.sender);
        hasVoted.hasVotedAddress[msg.sender] = true;

        // Check if all FIP editors, except the one being revoked, have voted
        if (
            proposal.voters.length() == fipAddressList.length() - 1 &&
            fipAddressList.length() > 2
        ) {
            _finalizeProposal(fipEditorAddress, id, false);

            // Clean up any remaining proposals for this address
            _cleanupProposals(fipEditorAddress);
        }
    }

    /**
     * @notice Cleans up any remaining proposals for a given FIP editor address.
     * @param fipEditorAddress The address of the FIP editor whose proposals need to be cleaned up.
     */
    function _cleanupProposals(address fipEditorAddress) internal {
        // Iterate through all approval proposals and remove those related to the given address
        _cleanupProposalsByType(fipEditorAddress, approveProposalId, true);

        // Iterate through all revocation proposals and remove those related to the given address
        _cleanupProposalsByType(fipEditorAddress, revokeProposalId, false);
    }

    /**
     * @notice Cleans up proposals of a specific type (approval or revocation) related to a given FIP editor address.
     * @param fipEditorAddress The address of the FIP editor whose proposals need to be cleaned up.
     * @param proposalSet The set of proposal IDs (either approval or revocation).
     * @param isApproval True if cleaning up approval proposals, false for revocation proposals.
     */
    function _cleanupProposalsByType(
        address fipEditorAddress,
        EnumerableSet.UintSet storage proposalSet,
        bool isApproval
    ) private {
        uint256[] memory proposalIds = proposalSet.values();

        for (uint256 i = 0; i < proposalIds.length; i++) {
            uint256 id = proposalIds[i];
            FipEditorProposal storage proposal = idToFipEditorProposal[id];

            if (proposal.voters.contains(fipEditorAddress)) {
                HasVoted storage hasVoted = idToHasVoted[id];
                hasVoted.hasVotedAddress[fipEditorAddress] = false;
                proposal.voters.remove(fipEditorAddress);
            }
            // Finalize the proposal if all FIP editors have voted (for approval) or all except one (for revocation)
            if (
                (isApproval &&
                    proposal.voters.length() == fipAddressList.length()) ||
                ((!isApproval &&
                    proposal.voters.length() == fipAddressList.length() - 1) &&
                    fipAddressList.length() > 2)
            ) {
                _finalizeProposal(proposal.fipEditorAddress, id, isApproval);
            }
        }
    }
    /**
     * @notice Finalizes the proposal by adding/removing the FIP editor and cleaning up storage.
     * @param fipEditorAddress The address of the FIP editor.
     * @param id The ID of the proposal.
     * @param isApproval True if it's an approval proposal, false for revocation.
     */
    function _finalizeProposal(
        address fipEditorAddress,
        uint256 id,
        bool isApproval
    ) private {
        if (isApproval) {
            fipAddressList.add(fipEditorAddress);
            approveProposalId.remove(id);
            fipAddressMap[fipEditorAddress] = true;
        } else {
            fipAddressList.remove(fipEditorAddress);
            revokeProposalId.remove(id);
            fipAddressMap[fipEditorAddress] = false;
        }

        delete idToHasVoted[id];
        delete idToFipEditorProposal[id];
        hasActiveProposal[fipEditorAddress] = false;
    }

    /**
     * @notice Get the list of addresses of all FIP editors.
     * @return An array containing the addresses of all FIP editors.
     */
    function getFipAddressList()
        external
        view
        override
        returns (address[] memory)
    {
        return fipAddressList.values();
    }

    /**
     * @notice Get the details of a proposal based on its ID.
     * @param id The ID of the proposal.
     * @return The details of the proposal.
     */
    function getFipEditorProposal(
        uint256 id
    )
        external
        view
        returns (uint256, address, string memory, address[] memory)
    {
        FipEditorProposal storage proposal = idToFipEditorProposal[id];
        return (
            proposal.proposalId,
            proposal.fipEditorAddress,
            proposal.voterInfo,
            proposal.voters.values()
        );
    }

    /**
     * @notice Get the list of IDs of all approved proposals.
     * @return An array containing the IDs of all approved proposals.
     */
    function getApproveProposalId()
        external
        view
        override
        returns (uint256[] memory)
    {
        return approveProposalId.values();
    }

    /**
     * @notice Get the list of IDs of all revoked proposals.
     * @return An array containing the IDs of all revoked proposals.
     */
    function getRevokeProposalId()
        external
        view
        override
        returns (uint256[] memory)
    {
        return revokeProposalId.values();
    }

    function isFipEditor(address sender) external view override returns (bool) {
        return fipAddressMap[sender];
    }
}
