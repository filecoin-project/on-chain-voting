# PowerVoting-Contract 

## Overview
The PowerVoting-Contract consists of two contracts:
1. **Vote Contract**: Responsible for creating proposals and conducting proposal voting.
2. **FipEditor Contract**: Handles permission management and supports adding and deleting FipEditors through proposals.

## Deployment and Upgrade Process

### Prerequisites
- Node.js version v18.13.0 or higher.

### Installation
Install the necessary libraries by running the following command:
```bash
npm install
```

### Environment Variable Configuration
1. Copy the `.env.example` file to `.env`:
```bash
cp .env.example .env
```
2. Configure the corresponding network keys in the `.env` file:
- `PRIVATE_KEY_TESTNET`: The private key of the account for testnet deployment.
- `PRIVATE_KEY_MAINNET`: The private key of the account for mainnet deployment.

### Contract Deployment
1. **Deploy the FipEditor Contract**
Run the following command to deploy the FipEditor contract:
```bash
npx hardhat run scripts/deploy_fip.ts --network [network_name]
```
2. **Deploy the Vote Contract**
Run the following command to deploy the Vote contract:
```bash
npx hardhat run scripts/deploy_vote.ts --network [network_name]
```
3. **Contract Address Storage**
After the contracts are deployed, the contract addresses will be saved in the `[network_name]_config.json` file in the `scripts` directory. The content format is as follows:
```json
{
  "POWER_VOTING_FIP": "",
  "POWER_VOTING_VOTE": ""
}
```

### Contract Upgrade
The contracts support upgrade via the UUPS (Universal Upgradeable Proxy Standard) pattern.
- **Upgrade the FipEditor Contract**
Run the following command to upgrade the FipEditor contract:
```bash
npx hardhat run scripts/upgrade_fip.ts --network [network_name]
```
- **Upgrade the Vote Contract**
Run the following command to upgrade the Vote contract:
```bash
npx hardhat run scripts/upgrade_vote.ts --network [network_name]
```

## Notes
- Replace `[network_name]` with the actual network name (e.g., `filecoin_testnet`, `filecoin_mainnet`) when running the deployment and upgrade commands.
- Ensure that the private keys in the `.env` file are kept secure and not exposed.