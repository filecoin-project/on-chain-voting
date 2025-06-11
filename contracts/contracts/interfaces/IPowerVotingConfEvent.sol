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

import {GithubRepoInfo} from "../types.sol";

interface IPowerVotingConfEvent {
    /// @notice updates the snapshot expiration time
    /// @param oldExp old expiration days
    /// @param newExp new expiration days
    event SnapshotExpirationDay(uint256 oldExp, uint256 newExp);
    /// @notice updates the snapshot height
    /// @param dateStr old height
    /// @param height new height
    event SnapshotDays(string dateStr, uint64 height);
    /// @notice updates the github repo
    /// @param id autoIncrementId
    /// @param repoInfo github repo info
    event GithubRepoAdded(uint256 id, GithubRepoInfo repoInfo);
    /// @notice updates the github repo
    /// @param id autoIncrementId
    event GithubRepoRemoved(uint256 id);
    /// @notice updates the voting counting algorithm
    /// @param oldAlgorithm old algorithm
    /// @param algorithm new algorithm
    event VotingCountingAlgorithm(string oldAlgorithm, string algorithm);
}
