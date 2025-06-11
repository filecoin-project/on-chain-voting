import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
require("@nomicfoundation/hardhat-ethers");
require("@nomicfoundation/hardhat-ignition-ethers");
import '@openzeppelin/hardhat-upgrades';

import dotenv from "dotenv";
dotenv.config();

const PRIVATE_KEY_TESTNET = process.env.PRIVATE_KEY_TESTNET ?? "";

const PRIVATE_KEY_MAINNET = process.env.PRIVATE_KEY_MAINNET ?? "";

const config: HardhatUserConfig = {
  solidity: {
    version: "0.8.22",
    settings: {
      viaIR: true,
    },
  },
  networks: {
    filecoin_devnet: {
      url: "https://filecoin-calibration.chainup.net/rpc/v1",
      chainId: 314159,
      accounts: [PRIVATE_KEY_TESTNET],
    },
    filecoin_testnet: {
      url: "https://filecoin-calibration.chainup.net/rpc/v1",
      chainId: 314159,
      accounts: [PRIVATE_KEY_TESTNET],
    },
    filecoin_mainnet: {
      url: "https://api.node.glif.io/rpc/v1",
      chainId: 314,
      accounts: [PRIVATE_KEY_MAINNET],
    },
  },
};

export default config;
