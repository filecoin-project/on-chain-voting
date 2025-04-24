import { POWER_VOTING_FIP } from "./constant";
import { getConstantJson } from "./utils";
const { ethers, upgrades } = require("hardhat");

async function main() {
  const constantJSON = getConstantJson();
  const POWER_VOTING_FIP_ADDRESS = constantJSON[POWER_VOTING_FIP];
  if (!ethers.isAddress(POWER_VOTING_FIP_ADDRESS)) {
    throw new Error(`Invalid contract address: ${POWER_VOTING_FIP_ADDRESS}`);
  }
  console.log(`Target proxy contract address: ${POWER_VOTING_FIP_ADDRESS}`);
  const PowerVotingFipEditor = await ethers.getContractFactory(
    "PowerVotingFipEditor"
  );
  try {
    const currentImplementationAddress =
      await upgrades.erc1967.getImplementationAddress(POWER_VOTING_FIP_ADDRESS);
    console.log("currentImplementationAddress ", currentImplementationAddress);
  } catch (error) {}
  console.log("Force importing existing proxy...");
  await upgrades.forceImport(POWER_VOTING_FIP_ADDRESS, PowerVotingFipEditor);
  console.log("Proxy contract successfully registered");

  await upgrades.upgradeProxy(POWER_VOTING_FIP_ADDRESS, PowerVotingFipEditor);
  console.log("Upgrade completed successfully");
}

main().catch(console.error);
