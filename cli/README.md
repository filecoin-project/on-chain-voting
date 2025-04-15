### Fil Vote 

`fil-vote` is a command-line tool for interacting with the Power Voting system, supporting wallet management and proposal operations.

## Features Overview

- **Wallet Management**: Import wallets, list wallets, set the default wallet.
- **Proposal Management**: View proposals, vote approve/reject, view proposal voting results.

## Installation Steps

### 1. Install Golang 1.23.0

Ensure that Go version 1.23.0 is installed.

### 2. Clone the Repository

```bash
git clone https://github.com/yourusername/fil-vote.git
```

### 3. Navigate to the Project Directory

```bash
cd fil-vote
```

### 4. Install Dependencies

```bash
go mod tidy
```

### 5. Build the Project

```bash
go build -o fil-vote
```

This will generate an executable file named `fil-vote`.

## Configuration

Configure the network connection parameters for the Power Voting system in the `config/config.go` file. Modify these settings according to your environment.

```
network:
  chainID: {CHAIN_ID}
  rpc: "{RPC_URL}"
  token: "{JWT_TOKEN}"
  powerVotingContract: "{POWER_VOTING_CONTRACT}"
  powerBackendURL: "{POWER_BACKEND_URL}"

abiPath:
  powerVotingAbi: "{POWER_VOTING_ABI_PATH}"

drand:
  urls:
    - "{DRAND_URL_1}"
    - "{DRAND_URL_2}"
    - "{DRAND_URL_3}"
    - "{DRAND_URL_4}"
    - "{DRAND_URL_5}"
  chainHash: "{DRAND_CHAIN_HASH}"

```



## Usage

### Wallet Management

#### 1. Add Wallet

Import a wallet and specify the wallet type and private key:

```bash
fil-vote wallet add [walletType] [privateKey]
```

**Parameters**:

- `walletType`: The type of wallet.
- `privateKey`: The wallet's private key.

#### 2. List Wallets

List all wallets connected to the Lotus node:

```bash
fil-vote wallet ls
```

#### 3. Set Default Wallet

Set a wallet as the default wallet for the Lotus node:

```bash
fil-vote wallet use [walletAddress]
```

**Parameters**:

- `walletAddress`: The wallet address to set as the default.

### Proposal Management

#### 1. List Proposals

List all proposals, with pagination support. Press `n` to go to the next page, `p` to go back to the previous page, and `q` to quit:

```bash
fil-vote proposal ls
```

#### 2. View Proposal Results

View the detailed information and voting results of a specific proposal:

```bash
fil-vote proposal results [proposalID]
```

**Parameters**:

- `proposalID`: The ID of the proposal.

#### 3. Approve Proposal

Cast an approve vote for a proposal. If the `from` parameter is not specified, the default wallet will be used for voting:

```bash
fil-vote proposal approve --proposalId <proposalID> --from <walletAddress>
```

**Parameters**:

- `proposalId`: The ID of the proposal to vote on.
- `from`: The wallet address, defaulting to the default wallet.

#### 4. Reject Proposal

Cast a reject vote for a proposal:

```bash
fil-vote proposal reject --proposalId <proposalID> --from <walletAddress>
```

**Parameters**:

- `proposalId`: The ID of the proposal to vote on.
- `from`: The wallet address, defaulting to the default wallet.

### Command Help

You can view help information for each command by using the `-h` or `--help` flag. For example, to view help for the `proposal` command:

```bash
fil-vote proposal -h
```