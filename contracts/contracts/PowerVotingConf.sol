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

import {IPowerVotingConf} from "./interfaces/IPowerVotingConf.sol";
import {GithubRepoInfo} from "./types.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

contract PowerVotingConf is
    IPowerVotingConf,
    Ownable2StepUpgradeable,
    UUPSUpgradeable
{
    /// @dev Store GitHub repository name information
    mapping(uint256 => GithubRepoInfo) public githubInfos;

    /// @dev To determine whether a GitHub repository name already exists or has been deleted
    mapping(string => uint256) public githubRepoNameToId;

    /// @dev Storage the snapshot height by date. (e.g "20240101" => 123456)
    mapping(bytes3 => uint64) public shapshotDays;

    /// @dev Voting counted algorithm for result
    string public votingCountingAlgorithm;

    /// @dev  Maximum expiration days for snapshot data
    uint64 public expDays;

    /// @dev Repo info id
    uint256 public githubRepoId;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
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
    function initialize() public initializer {
        expDays = 60; // default 60 days
        __UUPSUpgradeable_init();
        __Ownable_init(msg.sender);
    }

    /// @notice Batch add github repos
    /// @dev  This function is only callable by the owner of the contract.
    /// @param orgType 0: "CoreFilecoinOrg", 1:"EcosystemOrg", 2: "GithubUser" ...
    /// @param values repos name
    function batchAddGithubRepo(
        uint8 orgType,
        string[] memory values
    ) external onlyOwner {
        require(values.length > 0, "repo is empty");
        _addGithubRepos(orgType, values);
    }

    /// @notice If repos is unavailable, you can use this function to remove it.
    /// @dev This function is only callable by the owner of the contract.
    /// @param ids ids needs to be removed repo ids
    function batchRemoveGithubRepos(uint256[] memory ids) external onlyOwner {
        require(ids.length > 0, "repo is empty");
        _removeGithubRepos(ids);
    }

    /// @notice Set snapshot height by date. (e.g "20240101" => 123456)
    /// @dev Explain to a developer any extra details
    /// @param dateStr Snapshot date string (e.g "20240101")
    /// @param height  Block height corresponding to the snapshot date
    function setSnapshotHeight(
        string memory dateStr,
        uint64 height
    ) external onlyOwner {
        require(bytes(dateStr).length == 8, "Invalid date format");
        shapshotDays[bytes3(bytes(dateStr))] = height;
        emit SnapshotDays(dateStr, height);
    }

    /// @notice Get snapshot height by date. (e.g "20240101" => 123456)
    /// @dev If you want to get the snapshot height by date, you can use this function.
    /// @param dateStr Get the height of which day, it must be a length of 8. (e.g "20240101")
    /// @return Snapshot height. e.g 123456
    function getSnapshotHeight(
        string memory dateStr
    ) external view returns (uint64) {
        require(bytes(dateStr).length == 8, "Invalid date format");
        return shapshotDays[bytes3(bytes(dateStr))];
    }

    /// @notice Set snapshot data expiration days
    /// @dev This function is only callable by the owner of the contract.
    /// @param newDays New expiration days
    function setExpDays(uint64 newDays) external onlyOwner {
        require(newDays > 0 && newDays <= 180, "Invalid expiration days");
        if (newDays == expDays) {
            return;
        }
        uint64 oldDay = expDays;
        expDays = newDays;
        emit SnapshotExpirationDay(oldDay, newDays);
    }

    /// @notice Set voting counted algorithm for result
    /// @dev This function is only callable by the owner of the contract.
    /// @param algorithm Voting counted algorithm for result
    function setVotingCountingAlgorithm(
        string memory algorithm
    ) external onlyOwner {
        require(bytes(algorithm).length > 0, "Invalid algorithm");
        string memory oldAlgorithm = votingCountingAlgorithm;
        votingCountingAlgorithm = algorithm;
        emit VotingCountingAlgorithm(oldAlgorithm, algorithm);
    }

    /// @notice The specific implementation of adding the Github repository to the contract
    /// @dev If repo if already exists, it will not be added again.
    /// @param orgType 0: "CoreFilecoinOrg", 1:"EcosystemOrg", 2: "GithubUser" ...
    /// @param values Repos name
    function _addGithubRepos(uint8 orgType, string[] memory values) internal {
        for (uint256 i = 0; i < values.length; i++) {
            string memory repoName = values[i];
            if (githubRepoNameToId[repoName] > 0) {
                continue;
            }

            ++githubRepoId;

            githubRepoNameToId[repoName] = githubRepoId;
            GithubRepoInfo memory repoInfo = GithubRepoInfo(repoName, orgType);
            githubInfos[githubRepoId] = repoInfo;

            emit GithubRepoAdded(githubRepoId, repoInfo);
        }
    }

    /// @notice If repos is unavailable, you can use this function to remove it.
    /// @dev _gitHubInfos is a simulated collection value. To delete it,
    /// you need to remove it from the array and also delete the key-value pair that determines whether the element exists.
    /// @param ids needs to be removed repo ids
    function _removeGithubRepos(uint256[] memory ids) internal {
        for (uint256 i = 0; i < ids.length; i++) {
            uint256 id = ids[i];
            GithubRepoInfo memory repoInfo = githubInfos[id];

            delete githubInfos[id];
            delete githubRepoNameToId[repoInfo.repoName];

            emit GithubRepoRemoved(id);
        }
    }
}
