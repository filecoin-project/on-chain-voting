import { getConstantJson } from "./utils";

const { ethers } = require("hardhat");
async function main() {
  await check_fipeditor();
  await check_vote();
}

async function check_vote() {
  const PowerVoting = await ethers.getContractFactory("PowerVoting");
  const constantJSON = getConstantJson();
  const POWER_VOTING_VOTE = constantJSON["POWER_VOTING_VOTE"];
  const POWER_VOTING_ORCAL = constantJSON["POWER_VOTING_ORCAL"];
  const POWER_VOTING_FIP = constantJSON["POWER_VOTING_FIP"];

  const powerVotingContract = await PowerVoting.attach(POWER_VOTING_VOTE);
  const oracleAddress = await powerVotingContract.oracleContract();
  const fipEditorAddress = await powerVotingContract.fipEditorContract();
  console.log(
    "remote orcal=" + oracleAddress + "  fipEditor=" + fipEditorAddress
  );
  console.log(
    " local orcal=" + POWER_VOTING_ORCAL + "  fipEditor=" + POWER_VOTING_FIP
  );
}

async function check_fipeditor() {
  const PowerVotingFipEditor = await ethers.getContractFactory(
    "PowerVotingFipEditor"
  );
  const constantJSON = getConstantJson();
  const POWER_VOTING_FIP = constantJSON["POWER_VOTING_FIP"];
  const contract = await PowerVotingFipEditor.attach(POWER_VOTING_FIP);
  const accounts = await ethers.getSigners();
  const sender = accounts[0].address;
  const data = await contract.fipAddressMap(sender);
  console.log(accounts[0].address, data);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
