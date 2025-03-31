import { POWER_VOTING_POWER } from "./constant";
import { updateConstant } from "./utils";
const { ethers } = require("hardhat");
async function main() {
  const PowerVotingPOWER = await ethers.getContractFactory("Powers");
  console.log("Powers deploy beging");
  const contract = await PowerVotingPOWER.deploy();
  await contract.waitForDeployment();
  const address = await contract.getAddress();
  updateConstant(POWER_VOTING_POWER, address);
  console.log("Powers deployed to:", address);
}

main().catch(console.error);
