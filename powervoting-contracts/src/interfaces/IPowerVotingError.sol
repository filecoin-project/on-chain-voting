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

interface IPowerVotingError {
    // time error
    error TimeError(string);

    // status error
    error StatusError(string);

    // zero address error
    error ZeroAddressError(string);

    // call other contract error
    error CallError(string);

    // add FIP editor role error
    error AddFIPError(string);

    // fip already exists error
    error AddressIsAlreadyFipEditor(string);

    // address has active proposal error
    error AddressHasActiveProposal(string);

    // cannot propose to self error
    error CannotProposeToSelf(string);

    // only fip editors allowed error
    error OnlyFipEditorsAllowed(string);

    // cannot vote for own proposal error
    error CannotVoteForOwnProposal(string);

    // Invalid Proposal id error
    error InvalidProposalId(string);

    // Invalid Proposal type error
    error InvalidProposalType(string);
}
