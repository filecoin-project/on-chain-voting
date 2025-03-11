import hre from "hardhat";
import {
  POWER_VOTING_FIP,
  POWER_VOTING_ORCAL,
  POWER_VOTING_VOTE,
} from "./constant";
import { getConstantJson, updateConstant } from "./utils";
const { ethers, upgrades } = require("hardhat");
async function main() {
  const constantJSON = getConstantJson();
  const POWER_VOTING_FIP_ADDRESS = constantJSON[POWER_VOTING_FIP];
  const POWER_VOTING_ORCAL_ADDRESS = constantJSON[POWER_VOTING_ORCAL];
  console.log("FipEditor address = ", POWER_VOTING_FIP_ADDRESS);
  console.log("Oracle address = ", POWER_VOTING_ORCAL_ADDRESS);
  const PowerVoting = await ethers.getContractFactory("PowerVoting");
  console.log("wait deployProxy");
  const contract = await upgrades.deployProxy(
    PowerVoting,
    [POWER_VOTING_ORCAL_ADDRESS, POWER_VOTING_FIP_ADDRESS],
    {
      initializer: "initialize",
    }
  );
  console.log("wait deploy...");
  await contract.waitForDeployment();
  const address = await contract.getAddress();
  updateConstant(POWER_VOTING_VOTE, address);
  console.log("Vote deployed to:", address);
}

main().catch(console.error);
