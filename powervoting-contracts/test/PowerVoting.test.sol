// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../PowerVoting-filecoin.sol";

contract TestPowerVoting {

    PowerVoting public powerVoting;

    address public owner;
    address public oracleAddress = YOUR_ORACLE_ADDRESS; // Replace with your Oracle contract address
    address public addr1 = YOUR_ADDRESS_1;
    address public addr2 = YOUR_ADDRESS_2;

    event Log(string message);

    constructor() {
        owner = msg.sender;

        // Deploy the PowerVoting contract
        powerVoting = new PowerVoting();

        // Initialize the PowerVoting contract
        powerVoting.initialize(oracleAddress);
    }

    function test_initialization() public {
        emit Log("Testing initialization...");
        require(powerVoting.oracleContract() == oracleAddress, "Oracle contract address mismatch");
    }

    function test_fip_management() public {
        emit Log("Testing FIP management...");

        // Add FIP address
        powerVoting.addFIP(addr1);
        require(powerVoting.fipMap(addr1) == true, "FIP address not added correctly");

        // Remove FIP address
        powerVoting.removeFIP(addr1);
        require(powerVoting.fipMap(addr1) == false, "FIP address not removed correctly");
    }

    function test_proposal_creation() public {
        emit Log("Testing proposal creation...");

        // Add FIP address
        powerVoting.addFIP(addr1);

        // Create proposal
        string memory proposalCid = "ProposalCID";
        uint248 startTime = uint248(block.timestamp + 60); // Start time 1 minute in the future
        uint248 expTime = uint248(startTime + 3600); // Expire 1 hour after start
        uint256 proposalType = 1;

        powerVoting.createProposal(proposalCid, startTime, expTime, proposalType);

        // Get proposal details
        (string memory cid, uint256 pType, address creator, uint248 sTime, uint248 eTime, uint256 votesCount) = powerVoting.idToProposal(1);

        assert(keccak256(abi.encodePacked(cid)) == keccak256(abi.encodePacked(proposalCid)));
        assert(creator == addr1);
        assert(sTime == startTime);
        assert(eTime == expTime);
        assert(pType == proposalType);
        assert(votesCount == 0);
    }


    function test_voting() public {
        emit Log("Testing voting...");

        // Add FIP address and create proposal
        powerVoting.addFIP(addr1);
        string memory proposalCid = "ProposalCID";
        uint248 startTime = uint248(block.timestamp - 60); // Start time 1 minute ago
        uint248 expTime = uint248(startTime + 3600); // Expire 1 hour after start
        uint256 proposalType = 1;

        powerVoting.createProposal(proposalCid, startTime, expTime, proposalType);

        // Cast vote
        string memory voteInfo = "I support this proposal!";
        powerVoting.vote(1, voteInfo);

        ( string memory info, address voter) = powerVoting.proposalToVote(1, 1);
        require(voter == address(this), "Voter address mismatch");
        require(keccak256(bytes(info)) == keccak256(bytes(voteInfo)), "Vote info mismatch");
    }

    function test_voting_time_constraints() public {
        emit Log("Testing voting time constraints...");

        // Add FIP address and create proposal with future start time
        powerVoting.addFIP(addr1);
        string memory proposalCid = "ProposalCID";
        uint248 startTime = uint248(block.timestamp + 60); // Start time 1 minute in the future
        uint248 expTime = uint248(startTime + 3600); // Expire 1 hour after start
        uint256 proposalType = 1;

        powerVoting.createProposal(proposalCid, startTime, expTime, proposalType);

        string memory voteInfo = "I support this proposal!";
        try powerVoting.vote(2, voteInfo) {
            revert("Voting before proposal start time should fail");
        } catch {}

        // Create proposal with expired time
        startTime = uint248(block.timestamp - 3600); // Start time 1 hour ago
        expTime = uint248(startTime + 1800); // Expired 30 minutes ago
        powerVoting.createProposal(proposalCid, startTime, expTime, proposalType);

        try powerVoting.vote(3, voteInfo) {
            revert("Voting after proposal expiration time should fail");
        } catch {}
    }
}
