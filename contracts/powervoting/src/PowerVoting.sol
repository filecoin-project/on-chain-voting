// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import { OwnableUpgradeable } from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import { UUPSUpgradeable } from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import { Initializable } from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {IPowerVoting} from "./interfaces/IPowerVoting.sol";
import {Proposal,VoteResult,ProposalStatus,RESULT_MAX_LENGTH,MAX_ADDRESS} from "./types.sol";

contract PowerVoting is IPowerVoting, Initializable, OwnableUpgradeable, UUPSUpgradeable {

    // override from UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    // proposal id, increment
    uint256 proposalId;


    // proposal mapping, key: proposal id, value: Proposal
    mapping(uint256 => Proposal) public proposalMap;
    // vote power mapping
    mapping(uint256 => mapping(uint256 => uint256)) votePowerMap;
    // voted mapping, out key: proposal id, inner key: address, value: whether voted
    // proposal voters addresses, key: proposal id, value: address array
    mapping(uint256 => address[]) proposalVotersAddrs;
    // proposal AllowList, key: address, value: whether can create proposal
    mapping(address => bool) proposalAllowList;

    function initialize() public initializer {
        proposalId = 1;
        __UUPSUpgradeable_init();
        __Ownable_init();
    }

    /**
     * create a proposal and store it into mapping
     *
     * @param proposalCid: proposal content is stored in ipfs, proposal cid is ipfs cid for proposal content
     * @param expTime: proposal expiration timestamp, second
     */
    function createProposal(string calldata proposalCid, uint248 expTime, uint256 chainId, uint256 proposalType) override external {
        if(!proposalAllowList[msg.sender]) {
            revert PermissionError("No permission to create a new proposal");
        }

        Proposal storage proposal = proposalMap[proposalId];

        proposal.cid = proposalCid;
        proposal.creator = msg.sender;
        proposal.expTime = expTime;
        proposal.chainId = chainId;
        proposal.proposalType = proposalType;
        proposal.status = ProposalStatus.InProgress;

        emit ProposalCreate(proposalId, ProposalStatus.InProgress, proposal);

        proposalId++;
    }

    /**
     * update proposal AllowList
     * @param addrList: address list
     */
    function updateProposalAllowList(address[] calldata addrList) onlyOwner override external{
        uint256 addrLen = addrList.length;
        if(addrLen > MAX_ADDRESS){
            revert AddressMaxError("address list over max address length");
        }
        for(uint256 i; i < addrLen; i++){
            proposalAllowList[addrList[i]] = true;
        }

    }


    /**
     * cancel proposal
     * @param id: proposal id
     */
    function cancelProposal(uint256 id) override external{
        Proposal storage proposal = proposalMap[id];
        if(msg.sender != proposal.creator){
            revert PermissionError("No permission to cancel");
        }
        // time
        if (block.timestamp > proposal.expTime) {
            revert TimeError("proposal is expired");
        }
        if (proposal.status != ProposalStatus.InProgress) {
            revert StatusError("proposal status not in progress");
        }
        proposal.status = ProposalStatus.Canceled;
        emit ProposalCancel(id, ProposalStatus.Canceled);
    }

    /**
     * vote
     * @param id: proposal id
     * @param voteInfo: vote info, IPFS cid
     */
    function vote(uint256 id, string calldata voteInfo) override external{
        Proposal storage proposal = proposalMap[id];
        if(proposal.status != ProposalStatus.InProgress){
            revert StatusError("proposal status not in progress");
        }
        if(!proposalAllowList[msg.sender]){
            revert PermissionError("No permission to vote");
        }
        // if proposal is expired, won't be allowed to vote
        if(proposal.expTime <= block.timestamp){
            revert TimeError("proposal is expired");
        }
        proposalVotersAddrs[id].push(msg.sender);
        emit Vote(id, msg.sender, voteInfo);
    }

    /**
     * count
     * @param id: proposal id
     * @param voteResults: vote result array
     */
    function count(uint256 id, VoteResult[] calldata voteResults, string calldata voteListCid) onlyOwner override external {
        // if proposal is now expired, cant't count
        Proposal storage proposal = proposalMap[id];
        if(proposal.expTime > block.timestamp){
            revert TimeError("It's not time for the vote count");
        }
        if(proposal.status != ProposalStatus.InProgress){
            revert StatusError("proposal status not in progress");
        }
        if(voteResults.length > RESULT_MAX_LENGTH){
            revert OptionLengthError("result length must <= 5");
        }
        proposal.status = ProposalStatus.Completed;
        uint256 resultLen = voteResults.length;
        for(uint256 i; i < resultLen; i++){
            VoteResult memory voteItem = voteResults[i];
            VoteResult memory result;
            result.optionId = voteItem.optionId;
            result.votes = voteItem.votes;
            proposal.voteResults.push(result);
        }
        emit ProposalCount(id, ProposalStatus.Completed, proposal.voteResults, voteListCid);
    }

}