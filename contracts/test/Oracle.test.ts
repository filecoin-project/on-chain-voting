const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("Oracle", function () {
  let oracle: any;
  let owner: any;
  let voter: any;

  beforeEach(async function () {
    [owner, voter] = await ethers.getSigners();
    const PowerVotingOracle = await ethers.getContractFactory("Oracle");
    const contract = await upgrades.deployProxy(PowerVotingOracle, {
      initializer: "initialize",
    });
    oracle=await contract.waitForDeployment();
  });

  it("Should initialize correctly", async function () {
    expect(await oracle.owner()).to.equal(owner.address);
  });

  it("Should emit UpdateMinerIdsEvent when updating miner IDs", async function () {
    const minerIds = [1, 2, 3];
    await expect(oracle.connect(voter).updateMinerIds(minerIds))
      .to.emit(oracle, "UpdateMinerIdsEvent")
      .withArgs(voter.address, minerIds);
  });

  it("Should emit UpdateGistIdsEvent when updating gist ID", async function () {
    const gistId = "testGistId";
    await expect(oracle.connect(voter).updateGistId(gistId))
      .to.emit(oracle, "UpdateGistIdsEvent")
      .withArgs(voter.address, gistId);
  });
});
