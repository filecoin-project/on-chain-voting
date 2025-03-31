import hre from "hardhat";
import { POWER_VOTING_FIP, POWER_VOTING_VOTE } from "./constant";
import { getConstantJson, updateConstant } from "./utils";
const { ethers, upgrades } = require("hardhat");
async function main() {
  const constantJSON = getConstantJson();
  const POWER_VOTING_FIP_ADDRESS = constantJSON[POWER_VOTING_FIP];
  console.log("PowerVoting deploy beging");
  console.log("FipEditor address is ", POWER_VOTING_FIP_ADDRESS);
  const PowerVoting = await ethers.getContractFactory("PowerVoting");
  const contract = await upgrades.deployProxy(
    PowerVoting,
    [POWER_VOTING_FIP_ADDRESS],
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
