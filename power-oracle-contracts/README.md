# Oracle Contract

## Overview

The Oracle contract is a Solidity smart contract designed to manage various tasks related to the PowerVoting project. It facilitates task addition, management, voter information processing, power calculation, and interaction with PowerVoting contracts.

## Features

- **Task Management:** Allows for the addition and management of tasks associated with specific voters.
- **Power Calculation:** Calculates the power associated with voters based on their actor IDs and miner IDs.
- **Node Allow List:** Provides functionality to update the node allow list, controlling access to certain functions.
- **Voter Information:** Stores and retrieves information associated with voters, including actor IDs, GitHub accounts, and miner IDs.
- **Power History:** Maintains a history of power information for each voter.


## Implementation Details

- **Libraries:** Utilizes OpenZeppelin libraries such as `Counters` and `EnumerableSet` for managing task IDs and sets of addresses.
- **Modifiers:** Uses function modifiers to ensure that only authorized nodes can perform certain actions.
- **Upgradeability:** The contract is upgradeable using the UUPSUpgradeable pattern from OpenZeppelin.
- **Storage:** Employs various mappings and storage variables to efficiently store and retrieve information.

## Deployment

For deployment instructions, please refer to the [Deployment Guide](Install.md).

## Foundry Testing

Foundry is a fast, portable, and modular toolkit for Ethereum application development. The Oracle contract can be tested using Foundry to ensure its functionality and reliability.

### Prerequisites

Ensure you have Foundry installed. If not, you can install it using:

```
curl -L https://foundry.paradigm.xyz | bash
foundryup
```

### Setting Up the Test Environment

1. **Clone the Repository:**

   ```
   git clone https://github.com/filecoin-project/on-chain-voting.git
   cd power-oracle-contracts
   ```

2. **Install Dependencies:**

   Foundry uses a dependency management tool called `forge`. Install the dependencies with:

   ```
   forge install
   ```

3. **Compile the Contract:**

   Before running the tests, compile the contract to ensure there are no syntax errors:

   ```
   forge build
   ```

### Running Tests

Execute the tests using the following command:

```
forge test
```

This will run all the test cases present in the `test` directory and provide you with a detailed output of the test results.

## Author

The Oracle Contract is developed by StorSwift Inc.

## Version

This contract version is compatible with Solidity version ^0.8.19.
