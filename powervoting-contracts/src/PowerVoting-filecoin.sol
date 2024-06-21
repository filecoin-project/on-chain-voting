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

import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { UUPSUpgradeable } from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import { IPowerVoting } from "./interfaces/IPowerVoting-filecoin.sol";
import { Proposal, VoteInfo, VoterInfo, FipEditorProposal, HasVoted } from "./types.sol";


contract PowerVoting is IPowerVoting, Ownable2StepUpgradeable, UUPSUpgradeable {
    // Define constants for proposal types
    int8 public constant APPROVE_PROPOSAL_TYPE = 1;

    int8 public constant REVOKE_PROPOSAL_TYPE = 0;

    // proposal id
    uint256 public proposalId;

    // fip editor proposal id
    uint256 public fipEditorProposalId;

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

    // proposal id to vote, out key: proposal id, inner key: vote id, value: vote info
    mapping(uint256 => mapping(uint256 => VoteInfo)) public proposalToVote;

    // fip editor address mapping, key: address, value: boolean
    mapping(address => bool) public fipAddressMap;

    // fip editor proposal mapping, key: fip editor proposal id, value: fip editor proposal
    mapping(uint256 => FipEditorProposal) private idToFipEditorProposal;

    // proposal id to vote status, outer key: proposal id, inner key: address, value: HasVoted
    mapping(uint256 => HasVoted) private idToHasVoted;

    // active proposal mapping, key: address, value: boolean
    mapping(address => bool) public hasActiveProposal;

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

    /**
    * @dev Modifier to allow only FIP editors to call the function.
    * Reverts with a custom error message if the caller is not a FIP editor.
    */
    modifier onlyFIPEditors() {
        if(!fipAddressMap[msg.sender]){
            revert OnlyFipEditorsAllowed("Only FIP editors can call this function");
        }
        _;
    }

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

        address sender = msg.sender;
        fipAddressMap[sender] = true;
        fipAddressList.add(sender);

        __UUPSUpgradeable_init();
        __Ownable_init(sender);
    }

    /**
     * @notice Updates the address of the Oracle contract.
     * @param oracleAddress The new address of the Oracle contract.
     */
    function updateOracleContract(address oracleAddress) external onlyOwner nonZeroAddress (oracleAddress) {
        oracleContract = oracleAddress;
    }

    /**
    * @notice Creates a proposal to add or remove a FIP editor.
    * @param fipEditorAddress The address of the proposed FIP editor.
    * @param voterInfoCid The CID of the voter's information.
    * @param fipEditorProposalType The type of proposal: 1 for adding a FIP editor, 0 for removing.
    */
    function createFipEditorProposal(address fipEditorAddress, string calldata voterInfoCid, int8 fipEditorProposalType) override external onlyFIPEditors nonZeroAddress(fipEditorAddress) {
        // Check if the address is already a FIP editor and the proposal is not to revoke (0)
        if (fipAddressMap[fipEditorAddress] && fipEditorProposalType != REVOKE_PROPOSAL_TYPE) {
            revert AddressIsAlreadyFipEditor("Address is already a FIP editor");
        }

        // Check if the address already has an active proposal
        if (hasActiveProposal[fipEditorAddress]) {
            revert AddressHasActiveProposal("Address has an active proposal");
        }

        // Ensure the proposer is not proposing themselves
        if (fipEditorAddress == msg.sender) {
            revert CannotProposeToSelf("Cannot propose to self");
        }

        // Ensure the proposal type is valid (must be 1 for approval or 0 for revocation)
        if (fipEditorProposalType != APPROVE_PROPOSAL_TYPE && fipEditorProposalType != REVOKE_PROPOSAL_TYPE) {
            revert InvalidProposalType("Invalid proposal type: must be 0 or 1");
        }

        // There must be at minimum two votes to revoke a FIP Editor.
        if (fipEditorProposalType == REVOKE_PROPOSAL_TYPE && fipAddressList.length() <= 2) {
            revert InsufficientEditors("There must be more than two FIP editors to revoke an editor's FIP status");
        }

        // Increment the global proposal ID counter
        ++fipEditorProposalId;

        // Create a new proposal and store it in the mapping
        FipEditorProposal storage proposal = idToFipEditorProposal[fipEditorProposalId];
        proposal.fipEditorAddress = fipEditorAddress;
        proposal.voterInfoCid = voterInfoCid;
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
    function approveFipEditor(address fipEditorAddress, uint256 id) override public onlyFIPEditors {
        // Fetch the voting record for the proposal
        HasVoted storage hasVoted = idToHasVoted[id];

        // Check if the sender has already voted for this proposal
        if (hasVoted.hasVotedAddress[msg.sender]) {
            revert AddressHasActiveProposal("Address has already voted for this proposal");
        }

        // Ensure the proposal ID is valid for approval
        if (!approveProposalId.contains(id)) {
            revert InvalidProposalId("The proposal ID is not valid for approval");
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
    function revokeFipEditor(address fipEditorAddress, uint256 id) override public onlyFIPEditors {
        // Fetch the voting record for the proposal
        HasVoted storage hasVoted = idToHasVoted[id];

        // Check if the sender has already voted for this proposal
        if (hasVoted.hasVotedAddress[msg.sender]) {
            revert AddressHasActiveProposal("Address has already voted for this proposal");
        }

        // Ensure the proposal ID is valid for revocation
        if (!revokeProposalId.contains(id)) {
            revert InvalidProposalId("The proposal ID is not valid for revocation");
        }

        // Ensure the voter is not voting on their own proposal
        if (fipEditorAddress == msg.sender) {
            revert CannotVoteForOwnProposal("Cannot vote for own proposal");
        }

        // Fetch the proposal details
        FipEditorProposal storage proposal = idToFipEditorProposal[id];

        // Add the sender to the list of voters for this proposal
        proposal.voters.add(msg.sender);
        hasVoted.hasVotedAddress[msg.sender] = true;

        // Check if all FIP editors, except the one being revoked, have voted
        if (proposal.voters.length() == fipAddressList.length() - 1 && fipAddressList.length() > 2) {
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
    function _cleanupProposalsByType(address fipEditorAddress, EnumerableSet.UintSet storage proposalSet, bool isApproval) private {
        uint256[] memory proposalIds = proposalSet.values();

        for (uint256 i = 0; i < proposalIds.length; i++) {
            uint256 id = proposalIds[i];
            FipEditorProposal storage proposal = idToFipEditorProposal[id];

            if (proposal.voters.contains(fipEditorAddress)) {
                proposal.voters.remove(fipEditorAddress);
                HasVoted storage hasVoted = idToHasVoted[id];
                hasVoted.hasVotedAddress[fipEditorAddress] = false;
            }
            // Finalize the proposal if all FIP editors have voted (for approval) or all except one (for revocation)
            if ((isApproval && proposal.voters.length() == fipAddressList.length()) ||
                (!isApproval && proposal.voters.length() == fipAddressList.length() - 1) && fipAddressList.length() > 2) {
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
    function _finalizeProposal(address fipEditorAddress, uint256 id, bool isApproval) private {
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
    function getFipAddressList() override external view returns (address[] memory) {
        return fipAddressList.values();
    }

    /**
    * @notice Get the details of a proposal based on its ID.
    * @param id The ID of the proposal.
    * @return The details of the proposal.
    */
   function getFipEditorProposal(uint256 id) external view returns (uint256, address, string memory, address[] memory) {
        FipEditorProposal storage proposal = idToFipEditorProposal[id];
        return (
            proposal.proposalId,
            proposal.fipEditorAddress,
            proposal.voterInfoCid,
            proposal.voters.values()
        );
    }

    /**
    * @notice Get the list of IDs of all approved proposals.
    * @return An array containing the IDs of all approved proposals.
    */
    function getApproveProposalId() external override view returns (uint256[] memory) {
        return approveProposalId.values();
    }

    /**
    * @notice Get the list of IDs of all revoked proposals.
    * @return An array containing the IDs of all revoked proposals.
    */
    function getRevokeProposalId() external override view returns (uint256[] memory) {
        return revokeProposalId.values();
    }

    /**
     * @notice Creates a new proposal.
     * @param proposalCid The CID of the proposal.
     * @param startTime The start time of the proposal.
     * @param expTime The expiration time of the proposal.
     * @param proposalType The type of the proposal.
     */
    function createProposal(string calldata proposalCid, uint248 startTime, uint248 expTime, uint256 proposalType) onlyFIPEditors override external {
        ++proposalId;
        uint256 id = proposalId;
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
