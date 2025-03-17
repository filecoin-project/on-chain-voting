import { POWER_VOTING_FIP } from "./constant";
import { updateConstant } from "./utils";
const { ethers, upgrades } = require("hardhat");

async function deployContract() {
 
}

async function main() {
  const PowerVotingFipEditor = await ethers.getContractFactory(
    "PowerVotingFipEditor"
  );
  const contract = await upgrades.deployProxy(PowerVotingFipEditor, {
    initializer: "initialize",
  });
  await contract.waitForDeployment();
  const fipAddress = await contract.getAddress();
  updateConstant(POWER_VOTING_FIP, fipAddress);
  console.log("FipEditor deployed to:", fipAddress);
}

main().catch(console.error);
