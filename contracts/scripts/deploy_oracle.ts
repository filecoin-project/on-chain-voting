import { POWER_VOTING_ORACLE } from "./constant";
import { updateConstant } from "./utils";
const { ethers, upgrades } = require("hardhat");
async function main() {
  console.log("Oracle deploy beging");
  const PowerVotingORACLE= await ethers.getContractFactory("Oracle");
  const contract = await upgrades.deployProxy(PowerVotingORACLE, {
    initializer: "initialize"
  });
  await contract.waitForDeployment();
  const address = await contract.getAddress();
  updateConstant(POWER_VOTING_ORACLE, address);
  console.log("Oracle deployed to:", address);
}

main().catch(console.error);
