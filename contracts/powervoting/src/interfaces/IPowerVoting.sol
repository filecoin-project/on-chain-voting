// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {VoteResult} from "../types.sol";

import {IPowerVotingEvent} from "./IPowerVotingEvent.sol";
import {IPowerVotingError} from "./IPowerVotingError.sol";


interface IPowerVoting is IPowerVotingEvent, IPowerVotingError {
    function createProposal(string calldata proposalCid, uint248 expTime, uint256 chainId, uint256 proposalType) external;
    function updateProposalAllowList(address[] calldata addrList) external;
    function cancelProposal(uint256 id) external;
    function vote(uint256 id, string calldata voteInfo) external;
    function count(uint256 id, VoteResult[] calldata voteResults, string calldata voteListCid) external;
}