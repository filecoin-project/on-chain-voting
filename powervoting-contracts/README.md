# PowerVoting Contract

PowerVoting is a smart contract written in Solidity for managing voting activities within the PowerVoting project. It provides functionalities for creating proposals, voting on proposals, adding miner IDs, and delegating UCAN CIDs to the Oracle.

## Features
- **Proposal Creation:** Users can create proposals by specifying the CID, start time, expiration time, and type of proposal.
- **Voting:** Users can vote on proposals. The contract ensures that  checks if the proposal is within the valid voting period.
- **Miner ID Management:** Users can add miner IDs to the Oracle contract.
- **F4 Task Addition:** The contract automatically adds an F4 task for the caller if necessary.
- **UCAN CID Delegation:** Users can delegate UCAN CIDs to the Oracle for processing.

## Deploy
For deployment instructions, please refer to this document: [Deployment Guide](Install.md)

## Implementation Details
- The contract is implemented using the OpenZeppelin library for upgradeability and access control.
- It utilizes Counters for managing proposal IDs.
- Function selectors are used to interact with the Oracle contract.
- Various modifiers are used to ensure the validity of inputs and permissions.

## Author
PowerVoting Contract is developed by StorSwift Inc.

## Version
This contract version is compatible with Solidity version ^0.8.19.
