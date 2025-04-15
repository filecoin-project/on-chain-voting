import { POWER_VOTING_ORACLE, POWER_VOTING_POWER } from "./constant";
import { getConstantJson } from "./utils";
const { ethers, upgrades } = require("hardhat");

async function main() {
  const constantJSON = getConstantJson();
  const POWER_VOTING_ORACLE_ADDRESS = constantJSON[POWER_VOTING_ORACLE];
  if (!ethers.isAddress(POWER_VOTING_ORACLE_ADDRESS)) {
    throw new Error(`Invalid contract address: ${POWER_VOTING_ORACLE_ADDRESS}`);
  }
  console.log(`Target proxy contract address: ${POWER_VOTING_ORACLE_ADDRESS}`);
  const PowerVotingOracle = await ethers.getContractFactory("Oracle");
  // check forceImport
  try {
    const currentImplementationAddress =
      await upgrades.erc1967.getImplementationAddress(
        POWER_VOTING_ORACLE_ADDRESS
      );
    console.log("currentImplementationAddress ", currentImplementationAddress);
  } catch (error) {
    console.log("Force importing existing proxy...");
    await upgrades.forceImport(POWER_VOTING_ORACLE_ADDRESS, PowerVotingOracle);
    console.log("Proxy contract successfully registered");
  }
  await upgrades.upgradeProxy(POWER_VOTING_ORACLE_ADDRESS, PowerVotingOracle);
  console.log("Upgrade completed successfully");
}

main().catch(console.error);
