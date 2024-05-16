# Oracle Contract

The Oracle contract is a Solidity smart contract designed to manage various tasks related to the PowerVoting project. It facilitates task addition, management, voter information processing, power calculation, and interaction with PowerVoting contracts.

## Features
- **Task Management:** The contract allows for the addition of tasks associated with specific voters.
- **Power Calculation:** It calculates the power associated with voters based on their actor IDs and miner IDs.
- **Node Allow list:** Provides functionality to update the node allow list to control access to certain functions.
- **Voter Information:** Stores and retrieves information associated with voters, including actor IDs, GitHub accounts, and miner IDs.
- **Power History:** Maintains a history of power information for each voter.

## Deploy
For deployment instructions, please refer to this document: [Deployment Guide](Install.md)

## Implementation Details
- The contract uses OpenZeppelin libraries such as Counters and EnumerableSet for managing task IDs and sets of addresses.
- Function modifiers are used to ensure that only authorized nodes can perform certain actions.
- The contract is upgradeable using the UUPSUpgradeable pattern from OpenZeppelin.
- Various mappings and storage variables are used to efficiently store and retrieve information.

## Author
The Oracle Contract is developed by StorSwift Inc.

## Version
This contract version is compatible with Solidity version ^0.8.19.
