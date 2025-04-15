import { POWER_VOTING_VOTE, POWER_VOTING_FIP } from "./constant";
import { getConstantJson } from "./utils";

const { ethers } = require("hardhat");
async function main() {
  await check_fipeditor();
  await check_vote();
}

async function check_vote() {
  const PowerVoting = await ethers.getContractFactory("PowerVoting");
  const constantJSON = getConstantJson();
  const POWER_VOTING_VOTE_ADDRESS = constantJSON[POWER_VOTING_VOTE];
  const POWER_VOTING_FIP_ADDRESS = constantJSON[POWER_VOTING_FIP];

  const powerVotingContract = await PowerVoting.attach(
    POWER_VOTING_VOTE_ADDRESS
  );
 const count= await powerVotingContract.proposalId();

  const fipEditorAddress = await powerVotingContract.fipEditorContract();
  console.log("remote   fipEditor=" + fipEditorAddress);
  console.log("count=" + count);
}

async function check_fipeditor() {
  const PowerVotingFipEditor = await ethers.getContractFactory(
    "PowerVotingFipEditor"
  );
  const constantJSON = getConstantJson();
  const POWER_VOTING_FIP = constantJSON["POWER_VOTING_FIP"];
  const contract = await PowerVotingFipEditor.attach(POWER_VOTING_FIP);
  // const accounts = await ethers.getSigners();
  // const sender = accounts[0].address;
  const data = await contract.fipEditorCount();
  console.log(data)
  // console.log(sender,"status", data);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
