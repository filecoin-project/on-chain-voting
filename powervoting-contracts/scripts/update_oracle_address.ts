import { getConstantJson } from "./utils";

const { ethers } = require("hardhat");
async function main() {
  const PowerVoting = await ethers.getContractFactory("PowerVoting");
  const constantJSON = getConstantJson();
  const POWER_VOTING_VOTE = constantJSON["POWER_VOTING_VOTE"];
  const POWER_VOTING_ORCAL = constantJSON["POWER_VOTING_ORCAL"];

  const powerVotingContract = await PowerVoting.attach(POWER_VOTING_VOTE);
  const updateTx = await powerVotingContract.updateOracleContract(
    POWER_VOTING_ORCAL
  );
  updateTx.await();
  
  const oracleAddress = await powerVotingContract.oracleContract();
  console.log("remote orcal=" + oracleAddress);
  console.log(" local orcal=" + POWER_VOTING_ORCAL);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
