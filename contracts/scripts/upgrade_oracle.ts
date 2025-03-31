import { POWER_VOTING_ORCAL, POWER_VOTING_POWER } from "./constant";
import { getConstantJson } from "./utils";
const { ethers, upgrades } = require("hardhat");

async function main() {
  const constantJSON = getConstantJson();
  const POWER_VOTING_ORACLE_ADDRESS = constantJSON[POWER_VOTING_ORCAL];
  console.log("Ugrade begin ,Current Address is ",POWER_VOTING_ORACLE_ADDRESS);
  const PowerVotingOracle = await ethers.getContractFactory("Oracle");
  await upgrades.upgradeProxy(POWER_VOTING_ORACLE_ADDRESS, PowerVotingOracle);
  console.log("Oracle ugrade success");
}

main().catch(console.error);
