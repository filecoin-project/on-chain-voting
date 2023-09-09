// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ProposalStatus, Proposal, VoteResult} from "../types.sol";

interface IPowerVotingEvent {
    /**
     * vote event
     * @param id: proposal id
     * @param voteInfo: vote info, IPFS cid
     */
    event Vote(uint256 id, address voter, string voteInfo);

    /**
     * create proposal event
     * @param id: proposal id
     * @param status: proposal status
     * @param proposal: proposal detail
     */
    event ProposalCreate(uint256 id, ProposalStatus status, Proposal proposal);

    /**
     * cancel proposal event
     * @param id: proposal id
     * @param status: proposal status
     */
    event ProposalCancel(uint256 id, ProposalStatus status);

    /**
     * count event
     * @param id: proposal id
     * @param status: proposal status
     * @param voteResult: vote result
     */
    event ProposalCount(uint256 id, ProposalStatus status, VoteResult[] voteResult, string voteListCid);
}