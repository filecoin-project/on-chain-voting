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

interface CommonError {
    // Zero address error
    error ZeroAddressError();

    //Only FIP editors can call this function
    error OnlyFipEditorsAllowedError();
}

interface IPowerVotingError is CommonError {
    //Proposal not start yet.
    error VotingTimeNotStartedError();

    //Proposal expiration time reached.
    error VotingAlreadyEndedError();

    //Proposal EndTime Invalid
    error InvalidProposalEndTimeError();

    //Proposal Time Invalid
    error InvalidProposalTimeError();

    // call other contract error
    error CallError(string);

    // Invalid Proposal Percentage error
    error InvalidProposalPercentageError();

    //Proposal Percentage out of range
    error PercentageOutOfRangeError();

    //The maximum length of the proposal title is ${len}
    error TitleLengthLimitError(uint256);

    //The maximum length of the proposal content is ${len}
    error ContentLengthLimitError(uint256);

    //invalid vote info
    error InvalidVoteInfoError();

    //snapshot time out of range
    error SnapshotTimeOutOfRangeError();
}

interface IPowerVotingFipError is CommonError {
    //The proposal ID is not valid for revocation
    error InvalidRevocationProposalId();

    //The proposal ID is not valid for approval
    error InvalidApprovalProposalId();
    //Address is already a FIP editor
    error AddressIsAlreadyFipEditorError();

    // Address has an active proposal
    error AddressHasActiveProposalError();

    // Cannot propose to self
    error CannotProposeToSelfError();

    //There must be more than two FIP editors to revoke an editor's FIP status
    error InsufficientEditorsError();

    //Invalid proposal type: must be 0 or 1
    error InvalidProposalTypeError();

    //Cannot vote for own proposal
    error CannotVoteForOwnProposalError();

    //The maximum length of the voter info is ${len}
    error VoterInfoLimitError(uint256);
}
