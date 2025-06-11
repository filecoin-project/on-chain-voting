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

import {IPowerVotingConfEvent} from "./IPowerVotingConfEvent.sol";
import {GithubRepoInfo} from "../types.sol";
interface IPowerVotingConf is IPowerVotingConfEvent {
    /// @notice Batch add github repos
    /// @dev  This function is only callable by the owner of the contract.
    /// @param orgType 0: "CoreFilecoinOrg", 1:"EcosystemOrg", 2: "GithubUser" ...
    /// @param values repos name
    function batchAddGithubRepo(uint8 orgType, string[] memory values) external;

    /// @notice If repos is unavailable, you can use this function to remove it.
    /// @dev This function is only callable by the owner of the contract.
    /// @param ids needs to be removed repo ids
    function batchRemoveGithubRepos(uint256[] memory ids) external;

    /// @notice Set snapshot height by date. (e.g "20240101" => 123456)
    /// @dev Explain to a developer any extra details
    /// @param dateStr Snapshot date string (e.g "20240101")
    /// @param height  Block height corresponding to the snapshot date
    function setSnapshotHeight(string memory dateStr, uint64 height) external;

    /// @notice Get snapshot height by date. (e.g "20240101" => 123456)
    /// @dev If you want to get the snapshot height by date, you can use this function.
    /// @param dateStr Get the height of which day, it must be a length of 8. (e.g "20240101")
    /// @return Snapshot height. e.g 123456
    function getSnapshotHeight(
        string memory dateStr
    ) external view returns (uint64);

    /// @notice Set snapshot data expiration days
    /// @dev This function is only callable by the owner of the contract.
    /// @param newDays New expiration days
    function setExpDays(uint64 newDays) external;
}
