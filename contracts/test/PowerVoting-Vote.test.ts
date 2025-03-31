import { upgrades } from "hardhat";

const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("PowerVotingVote", function () {
  let powerVoting: any;
  let powerVotingFipEditor: any;
  let fipEditor: any;
  let fipContractAddress: any;
  let accounts: any;
  let voter: any;

  beforeEach(async function () {
    //update fip editor first
    const PowerVotingFipEditor = await ethers.getContractFactory(
      "PowerVotingFipEditor"
    );
    const fipContract = await upgrades.deployProxy(PowerVotingFipEditor, {
      initializer: "initialize",
    });
    powerVotingFipEditor = await fipContract.waitForDeployment();
    fipContractAddress = await fipContract.getAddress();

    //update vote
    const PowerVotingVote = await ethers.getContractFactory("PowerVoting");
    const contract = await upgrades.deployProxy(
      PowerVotingVote,
      [fipContractAddress],
      {
        initializer: "initialize",
      }
    );
    powerVoting = await contract.waitForDeployment();
    accounts = await ethers.getSigners();
    fipEditor = accounts[0];
    voter = accounts[1];
  });

  it("Test Initialize success", async function () {
    const fipEditorCount = await powerVotingFipEditor.fipEditorCount();
    expect(fipEditorCount).to.equal(1);
    const stauts = await powerVotingFipEditor.fipEditorStatusMap(accounts[0]);
    expect(stauts).to.equal(2); //approved

    const remoteFipEditorAddress = await powerVoting.fipEditorContract();
    expect(remoteFipEditorAddress).to.equal(fipContractAddress);
  });

  describe("createProposal", function () {
    it("should create a proposal successfully", async function () {
      const startTime = Math.floor(Date.now() / 1000) + 60;
      const endTime = startTime + 3600;
      const tokenHolderPercentage = 2500;
      const spPercentage = 2500;
      const clientPercentage = 2500;
      const developerPercentage = 2500;
      const content = "Test proposal content";
      const title = "Test proposal title";

      await expect(
        powerVoting
          .connect(fipEditor)
          .createProposal(
            startTime,
            endTime,
            tokenHolderPercentage,
            spPercentage,
            clientPercentage,
            developerPercentage,
            content,
            title
          )
      ).to.emit(powerVoting, "ProposalCreate");
    });

    it("should revert if title length exceeds limit", async function () {
      const startTime = Math.floor(Date.now() / 1000) + 60;
      const endTime = startTime + 3600;
      const tokenHolderPercentage = 2500;
      const spPercentage = 2500;
      const clientPercentage = 2500;
      const developerPercentage = 2500;
      const content = "Test proposal content";
      const longTitle = "a".repeat(201);
      const titleMaxLength = await powerVoting.titleMaxLength();
      await expect(
        powerVoting
          .connect(fipEditor)
          .createProposal(
            startTime,
            endTime,
            tokenHolderPercentage,
            spPercentage,
            clientPercentage,
            developerPercentage,
            content,
            longTitle
          )
      )
        .to.be.revertedWithCustomError(powerVoting, "TitleLengthLimitError")
        .withArgs(titleMaxLength);
    });

    it("should revert if content length exceeds limit", async function () {
      const startTime = Math.floor(Date.now() / 1000) + 60;
      const endTime = startTime + 3600;
      const tokenHolderPercentage = 2500;
      const spPercentage = 2500;
      const clientPercentage = 2500;
      const developerPercentage = 2500;
      const longContent = "a".repeat(10001);
      const title = "Test proposal title";
      const contentMaxLength = await powerVoting.contentMaxLength();

      await expect(
        powerVoting
          .connect(fipEditor)
          .createProposal(
            startTime,
            endTime,
            tokenHolderPercentage,
            spPercentage,
            clientPercentage,
            developerPercentage,
            longContent,
            title
          )
      )
        .to.be.revertedWithCustomError(powerVoting, "ContentLengthLimitError")
        .withArgs(contentMaxLength);
    });

    it("should revert if percentage is out of range", async function () {
      const startTime = Math.floor(Date.now() / 1000) + 60;
      const endTime = startTime + 3600;
      const tokenHolderPercentage = 10001;
      const spPercentage = 2500;
      const clientPercentage = 2500;
      const developerPercentage = 2500;
      const content = "Test proposal content";
      const title = "Test proposal title";

      await expect(
        powerVoting
          .connect(fipEditor)
          .createProposal(
            startTime,
            endTime,
            tokenHolderPercentage,
            spPercentage,
            clientPercentage,
            developerPercentage,
            content,
            title
          )
      ).to.be.revertedWithCustomError(powerVoting, "PercentageOutOfRangeError");
    });

    it("should revert if total percentage is not 100%", async function () {
      const startTime = Math.floor(Date.now() / 1000) + 60;
      const endTime = startTime + 3600;
      const tokenHolderPercentage = 2500;
      const spPercentage = 2500;
      const clientPercentage = 2500;
      const developerPercentage = 2000;
      const content = "Test proposal content";
      const title = "Test proposal title";

      await expect(
        powerVoting
          .connect(fipEditor)
          .createProposal(
            startTime,
            endTime,
            tokenHolderPercentage,
            spPercentage,
            clientPercentage,
            developerPercentage,
            content,
            title
          )
      ).to.be.revertedWithCustomError(
        powerVoting,
        "InvalidProposalPercentageError"
      );
    });

    it("should revert if end time is in the past", async function () {
      const startTime = Math.floor(Date.now() / 1000) + 60;
      const endTime = Math.floor(Date.now() / 1000) - 60;
      const tokenHolderPercentage = 2500;
      const spPercentage = 2500;
      const clientPercentage = 2500;
      const developerPercentage = 2500;
      const content = "Test proposal content";
      const title = "Test proposal title";

      await expect(
        powerVoting
          .connect(fipEditor)
          .createProposal(
            startTime,
            endTime,
            tokenHolderPercentage,
            spPercentage,
            clientPercentage,
            developerPercentage,
            content,
            title
          )
      ).to.be.revertedWithCustomError(
        powerVoting,
        "InvalidProposalEndTimeError"
      );
    });

    it("should revert if start time is after end time", async function () {
      const startTime = Math.floor(Date.now() / 1000) + 3600;
      const endTime = Math.floor(Date.now() / 1000) + 1660;
      const tokenHolderPercentage = 2500;
      const spPercentage = 2500;
      const clientPercentage = 2500;
      const developerPercentage = 2500;
      const content = "Test proposal content";
      const title = "Test proposal title";

      await expect(
        powerVoting
          .connect(fipEditor)
          .createProposal(
            startTime,
            endTime,
            tokenHolderPercentage,
            spPercentage,
            clientPercentage,
            developerPercentage,
            content,
            title
          )
      ).to.be.revertedWithCustomError(powerVoting, "InvalidProposalTimeError");
    });
  });

  describe("vote", function () {
    let proposalId: any;
    let offsetStarTime = 3600;
    let offsetEndTime = 7200;
    let snapshotId: any;

    beforeEach(async function () {
      const startTime = Math.floor(Date.now() / 1000) + offsetStarTime;
      const endTime = startTime + offsetEndTime;
      const tokenHolderPercentage = 2500;
      const spPercentage = 2500;
      const clientPercentage = 2500;
      const developerPercentage = 2500;
      const content = "Test proposal content";
      const title = "Test proposal title";

      const tx = await powerVoting
        .connect(fipEditor)
        .createProposal(
          startTime,
          endTime,
          tokenHolderPercentage,
          spPercentage,
          clientPercentage,
          developerPercentage,
          content,
          title
        );
      await tx.wait();
      proposalId = await powerVoting.proposalId();
      snapshotId = await ethers.provider.send("evm_snapshot", []);
    });

    afterEach(async () => {
      await ethers.provider.send("evm_revert", [snapshotId]);
    });
    it("should vote successfully", async function () {
      await ethers.provider.send("evm_increaseTime", [offsetStarTime]);
      await ethers.provider.send("evm_mine");

      const voteInfo = "Yes";
      await expect(
        powerVoting.connect(voter).vote(proposalId, voteInfo)
      ).to.emit(powerVoting, "Vote");
    });

    it("should revert if vote info is empty", async function () {
      const voteInfo = "";
      await expect(
        powerVoting.connect(voter).vote(proposalId, voteInfo)
      ).to.be.revertedWithCustomError(powerVoting, "InvalidVoteInfoError");
    });

    it("should revert if voting time has not started", async function () {
      const voteInfo = "Yes";
      await expect(
        powerVoting.connect(voter).vote(proposalId, voteInfo)
      ).to.be.revertedWithCustomError(powerVoting, "VotingTimeNotStartedError");
    });

    it("should revert if voting time has ended", async function () {
      await ethers.provider.send("evm_increaseTime", [offsetEndTime+offsetStarTime]);
      await ethers.provider.send("evm_mine");

      const voteInfo = "Yes";
      await expect(
        powerVoting.connect(voter).vote(proposalId, voteInfo)
      ).to.be.revertedWithCustomError(powerVoting, "VotingAlreadyEndedError");
    });
  });
});
