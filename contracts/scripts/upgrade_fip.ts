import { POWER_VOTING_FIP } from "./constant";
import { getConstantJson} from "./utils";
const { ethers, upgrades } = require("hardhat");

async function main() {
  const constantJSON = getConstantJson();
  const POWER_VOTING_FIP_ADDRESS = constantJSON[POWER_VOTING_FIP];
  console.log("Ugrade begin ,Current Address is ",POWER_VOTING_FIP_ADDRESS);
  const PowerVotingFipEditor = await ethers.getContractFactory(
    "PowerVotingFipEditor"
  );
  await upgrades.upgradeProxy(POWER_VOTING_FIP_ADDRESS, PowerVotingFipEditor);
  console.log("FipEditor ugrade success");
}

main().catch(console.error);
