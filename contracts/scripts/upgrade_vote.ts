import { POWER_VOTING_VOTE } from "./constant";
import { getConstantJson } from "./utils";
const { ethers, upgrades } = require("hardhat");

async function main() {
  try {
    console.log("Starting upgrade process...");
    const constantJSON = getConstantJson();
    if (!constantJSON) {
      throw new Error("Failed to get constant JSON");
    }
    const POWER_VOTING_VOTE_ADDRESS = constantJSON[POWER_VOTING_VOTE];
    if (!ethers.isAddress(POWER_VOTING_VOTE_ADDRESS)) {
      throw new Error(`Invalid contract address: ${POWER_VOTING_VOTE_ADDRESS}`);
    }
    console.log(`Target proxy contract address: ${POWER_VOTING_VOTE_ADDRESS}`);

    const PowerVotingFactory = await ethers.getContractFactory("PowerVoting");

    try {
      const currentImplementationAddress =
        await upgrades.erc1967.getImplementationAddress(
          POWER_VOTING_VOTE_ADDRESS
        );
      console.log(
        "currentImplementationAddress ",
        currentImplementationAddress
      );
    } catch (error) {}
    console.log("Force importing existing proxy...");
    await upgrades.forceImport(POWER_VOTING_VOTE_ADDRESS, PowerVotingFactory);
    console.log("Proxy contract successfully registered");
    // upgrade
    console.log("Executing upgrade...");
    await upgrades.upgradeProxy(POWER_VOTING_VOTE_ADDRESS, PowerVotingFactory);
    console.log("Upgrade completed successfully");
  } catch (error: any) {
    console.error("Upgrade failed:", error);
    process.exit(1);
  }
}

main().catch(console.error);
