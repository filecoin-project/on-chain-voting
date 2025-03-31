import hre from "hardhat";
import { POWER_VOTING_VOTE } from "./constant";
import { getConstantJson} from "./utils";
const { ethers, upgrades } = require("hardhat");

async function main() {
  const constantJSON = getConstantJson();
  const POWER_VOTING_VOTE_ADDRESS = constantJSON[POWER_VOTING_VOTE];
  console.log("Ugrade begin ,Current Address is ",POWER_VOTING_VOTE_ADDRESS);
  const PowerVoting = await ethers.getContractFactory("PowerVoting");
  await upgrades.upgradeProxy(POWER_VOTING_VOTE_ADDRESS, PowerVoting);
  console.log("Ugrade success");
}

main().catch(console.error);
