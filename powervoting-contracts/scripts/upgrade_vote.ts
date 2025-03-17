import hre from "hardhat";
import { POWER_VOTING_VOTE } from "./constant";
import { getConstantJson} from "./utils";
const { ethers, upgrades } = require("hardhat");

async function main() {
  const constantJSON = getConstantJson();
  const POWER_VOTING_VOTE_ADDRESS = constantJSON[POWER_VOTING_VOTE];
  const PowerVoting = await ethers.getContractFactory("PowerVoting");
  console.log("wait upgrde...")
  await upgrades.upgradeProxy(POWER_VOTING_VOTE_ADDRESS, PowerVoting);
  console.log("Vote ugrade success");
}

main().catch(console.error);
