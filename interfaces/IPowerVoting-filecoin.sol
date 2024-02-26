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


import {IPowerVotingEvent} from "./IPowerVotingEvent.sol";
import {IPowerVotingError} from "./IPowerVotingError.sol";


interface IPowerVoting is IPowerVotingEvent, IPowerVotingError {
    function createProposal(string calldata proposalCid, uint248 expTime, uint256 proposalType) external;
    function vote(uint256 id, string calldata info, uint64[] memory minerIds) external;
    function ucanDelegate(string calldata ucanCid) external;
    function updateOracleContract(address oracleAddress) external;
}