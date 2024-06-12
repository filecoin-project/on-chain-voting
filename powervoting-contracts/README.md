# PowerVoting Contract

## Overview

PowerVoting is a smart contract written in Solidity for managing voting activities within the PowerVoting project. It provides functionalities for creating proposals, voting on proposals, adding miner IDs, and delegating UCAN CIDs to the Oracle.

## Features
- **Proposal Creation:** Users can create proposals by specifying the CID, start time, expiration time, and type of proposal.
- **Voting:** Users can vote on proposals. The contract ensures that  checks if the proposal is within the valid voting period.
- **Miner ID Management:** Users can add miner IDs to the Oracle contract.
- **F4 Task Addition:** The contract automatically adds an F4 task for the caller if necessary.
- **UCAN CID Delegation:** Users can delegate UCAN CIDs to the Oracle for processing.


## Implementation Details
- The contract is implemented using the OpenZeppelin library for upgradeability and access control.

- It utilizes Counters for managing proposal IDs.

- Function selectors are used to interact with the Oracle contract.

- Various modifiers are used to ensure the validity of inputs and permissions.

## Deploy
For deployment instructions, please refer to this document: [Deployment Guide](Install.md)

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
   cd powervoting-contracts
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
PowerVoting Contract is developed by StorSwift Inc.

## Version
This contract version is compatible with Solidity version ^0.8.19.
