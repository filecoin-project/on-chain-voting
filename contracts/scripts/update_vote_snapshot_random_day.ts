import { POWER_VOTING_VOTE } from "./constant";
import { getConstantJson } from "./utils";

const { ethers } = require("hardhat");
const newSnapshotMaxRandomOffsetDays = 0; //Configure the maximum random number of days for snapshots, and the value should be greater than 0
async function main() {
  const PowerVoting = await ethers.getContractFactory("PowerVoting");
  const constantJSON = getConstantJson();
  const POWER_VOTING_VOTE_ADDRESS = constantJSON[POWER_VOTING_VOTE];

  const powerVotingContract = await PowerVoting.attach(
    POWER_VOTING_VOTE_ADDRESS
  );
  const currentRadomDay =
    await powerVotingContract.snapshotMaxRandomOffsetDays();
  console.log("Current Snapshot Random day ", currentRadomDay);

  if (newSnapshotMaxRandomOffsetDays <= 0) {
    throw Error("newSnapshotRandomDay should be greater than 0");
  }
  console.log("New Snapshot Random day ", newSnapshotMaxRandomOffsetDays);
  await powerVotingContract.setSnapshotMaxRandomOffsetDays(newSnapshotMaxRandomOffsetDays);
  console.log("Update Success");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
