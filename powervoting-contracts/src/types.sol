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


struct Proposal {
    // proposal cid
    string cid;
    // proposal type
    uint256 proposalType;
    // proposal creator
    address creator;
    // proposal start timestamp
    uint248 startTime;
    // proposal expiration timestamp, second
    uint248 expTime;
    // votes count
    uint256 votesCount;
}

struct VoteInfo {
    // vote info
    string voteInfo;
    // vote address
    address voter;
}

// voter info
struct VoterInfo {
    uint64[] actorIds;
    uint64[] minerIds;
    string githubAccount;
    address ethAddress;
    string ucanCid;
}

struct FipEditorProposal {
     // Unique identifier for the proposal
    uint256 proposalId;        
    // Address of the FIP editor
    address fipEditorAddress;   
    // CID (Content Identifier) of the voter's information
    string voterInfoCid;        
    // Array containing addresses of voters
    address[] voters;           
}

struct HasVoted {
    mapping(address => bool) hasVotedAddress;
}