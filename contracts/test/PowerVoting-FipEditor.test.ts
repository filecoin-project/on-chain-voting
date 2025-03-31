import { expect } from "chai";
import { ethers } from "hardhat";
import { upgrades } from "hardhat";

describe("PowerVotingFipEditor", function () {
  let powerVotingFipEditor: any;
  let accounts: any[];
  let owner: any;
  let otherAccounts: any;

  const STATUS = {
    REVOKED: 0,
    ADDING: 1,
    APPROVED: 2,
    REVOKING: 3,
  };

  const PROPOSAL_TYPE = {
    APPROVE: 1,
    REVOKE: 0,
  };

  beforeEach(async function () {
    accounts = await ethers.getSigners();
    [owner, ...otherAccounts] = accounts;

    const PowerVotingFipEditor = await ethers.getContractFactory(
      "PowerVotingFipEditor"
    );
    const contract = await upgrades.deployProxy(PowerVotingFipEditor, {
      initializer: "initialize",
    });
    powerVotingFipEditor = await contract.waitForDeployment();
  });

  describe("Initialization", function () {
    it("should initialize correctly", async function () {
      expect(await powerVotingFipEditor.fipEditorCount()).to.equal(1);
      expect(
        await powerVotingFipEditor.fipEditorStatusMap(owner.address)
      ).to.equal(STATUS.APPROVED);
    });
  });

  describe("createFipEditorProposal", function () {
    it("should auto-approve proposal when only one editor exists", async function () {
      const initialCount = await powerVotingFipEditor.fipEditorCount();

      await powerVotingFipEditor
        .connect(owner)
        .createFipEditorProposal(
          otherAccounts[0].address,
          "Auto-approve Test",
          PROPOSAL_TYPE.APPROVE
        );

      expect(await powerVotingFipEditor.fipEditorCount()).to.equal(
        parseInt(initialCount) + 1
      );
      expect(
        await powerVotingFipEditor.fipEditorStatusMap(otherAccounts[0].address)
      ).to.equal(STATUS.APPROVED);
    });

    it("should reject invalid proposal type", async function () {
      await expect(
        powerVotingFipEditor.createFipEditorProposal(
          otherAccounts[0].address,
          "Test",
          2
        )
      ).to.be.revertedWithCustomError(
        powerVotingFipEditor,
        "InvalidProposalTypeError"
      );
    });

    it("should prevent self-proposal", async function () {
      await expect(
        powerVotingFipEditor
          .connect(owner)
          .createFipEditorProposal(owner.address, "Test", PROPOSAL_TYPE.APPROVE)
      ).to.be.revertedWithCustomError(
        powerVotingFipEditor,
        "CannotProposeToSelfError"
      );
    });

    it("should validate candidate info length", async function () {
      await powerVotingFipEditor.setLengthLimits(5);
      await expect(
        powerVotingFipEditor.createFipEditorProposal(
          otherAccounts[0].address,
          "123456",
          PROPOSAL_TYPE.APPROVE
        )
      )
        .to.be.revertedWithCustomError(
          powerVotingFipEditor,
          "CandidateInfoLimitError"
        )
        .withArgs(5);
    });

    it("should reject duplicate approve proposal", async function () {
      await powerVotingFipEditor.createFipEditorProposal(
        otherAccounts[0].address,
        "Test",
        PROPOSAL_TYPE.APPROVE
      );
      await expect(
        powerVotingFipEditor.createFipEditorProposal(
          otherAccounts[0].address,
          "Test",
          PROPOSAL_TYPE.APPROVE
        )
      ).to.be.revertedWithCustomError(
        powerVotingFipEditor,
        "AddressIsAlreadyFipEditorError"
      );
    });
    it("should reject if candidate is not an fip editor", async function () {
      await expect(
        powerVotingFipEditor.createFipEditorProposal(
          otherAccounts[0].address,
          "Test",
          PROPOSAL_TYPE.REVOKE
        )
      ).to.be.revertedWithCustomError(
        powerVotingFipEditor,
        "AddressNotFipEditorError"
      );
    });
    it("should reject if candidate already in proposal", async function () {
      powerVotingFipEditor.createFipEditorProposal(
        otherAccounts[0].address,
        "Test",
        PROPOSAL_TYPE.APPROVE
      );
      powerVotingFipEditor.createFipEditorProposal(
        otherAccounts[1].address,
        "Test",
        PROPOSAL_TYPE.APPROVE
      );
      await expect(
        powerVotingFipEditor.createFipEditorProposal(
          otherAccounts[1].address,
          "Test",
          PROPOSAL_TYPE.APPROVE
        )
      ).to.be.revertedWithCustomError(
        powerVotingFipEditor,
        "AddressHasActiveProposalError"
      );
    });

    it("should reject the revoke proposal if there are fewer than 3 fip editor", async function () {
      powerVotingFipEditor.createFipEditorProposal(
        otherAccounts[0].address,
        "Test",
        PROPOSAL_TYPE.APPROVE
      );
      await expect(
        powerVotingFipEditor.createFipEditorProposal(
          otherAccounts[0].address,
          "Test",
          PROPOSAL_TYPE.REVOKE
        )
      ).to.be.revertedWithCustomError(
        powerVotingFipEditor,
        "InsufficientEditorsError"
      );
    });
  });

  describe("voteFipEditorProposal", function () {
    let proposalId: any;
    const fipEditors = new Set<number>();
    const createProposal = async (
      accountIndex: number,
      proposalType: number
    ): Promise<number> => {
      // create
      let candidateAddress = otherAccounts[accountIndex];
      await powerVotingFipEditor.createFipEditorProposal(
        candidateAddress,
        "111",
        proposalType
      );
      return powerVotingFipEditor.fipEditorProposalId();
    };
    const voteProposal = async (
      proposalId: number,
      votes: Set<number>,
      candidateIndex: number
    ) => {
      for (const index of votes) {
        if (candidateIndex != index) {
          await powerVotingFipEditor
            .connect(otherAccounts[index])
            .voteFipEditorProposal(proposalId);
        }
      }
    };
    const createAndvoteProposal = async (
      accountIndex: number,
      proposalType: number
    ) => {
      let proposalId = await createProposal(accountIndex, proposalType);
      let candidateAddress = otherAccounts[accountIndex];
      for (const index of fipEditors) {
        if (index != accountIndex) {
          await powerVotingFipEditor
            .connect(otherAccounts[index])
            .voteFipEditorProposal(proposalId);
        }
      }
      let stauts = await powerVotingFipEditor.fipEditorStatusMap(
        candidateAddress
      );
      if (proposalType == PROPOSAL_TYPE.APPROVE) {
        expect(stauts).to.equal(STATUS.APPROVED);
        fipEditors.add(accountIndex);
      } else {
        expect(stauts).to.equal(STATUS.REVOKED);
        fipEditors.delete(accountIndex);
      }

      const fipEditorCount = await powerVotingFipEditor.fipEditorCount();
      expect(fipEditorCount).to.equal(fipEditors.size + 1);
    };
    beforeEach(async function () {
      await powerVotingFipEditor
        .connect(owner)
        .createFipEditorProposal(
          otherAccounts[0].address,
          "Test Candidate",
          PROPOSAL_TYPE.APPROVE
        );
      fipEditors.add(0);
      await powerVotingFipEditor
        .connect(owner)
        .createFipEditorProposal(
          otherAccounts[1].address,
          "Test Candidate",
          PROPOSAL_TYPE.APPROVE
        );
      proposalId = await powerVotingFipEditor.fipEditorProposalId();
    });

    it("should auto-approve proposal when all editors vote", async function () {
      const editorCount = parseInt(await powerVotingFipEditor.fipEditorCount());
      await Promise.all(
        Array.from({ length: editorCount - 1 }, (_, i) =>
          powerVotingFipEditor
            .connect(otherAccounts[i])
            .voteFipEditorProposal(proposalId)
        )
      );
      expect(
        await powerVotingFipEditor.fipEditorStatusMap(otherAccounts[1].address)
      ).to.equal(STATUS.APPROVED);
    });

    it("should auto-approve proposal when all editors vote. (Special scene)", async function () {
      powerVotingFipEditor
        .connect(otherAccounts[0])
        .voteFipEditorProposal(proposalId);
      fipEditors.add(1);
      await createAndvoteProposal(2, PROPOSAL_TYPE.APPROVE);
      await createAndvoteProposal(3, PROPOSAL_TYPE.APPROVE);
      await createAndvoteProposal(4, PROPOSAL_TYPE.APPROVE);
      await createAndvoteProposal(5, PROPOSAL_TYPE.APPROVE);
      await createAndvoteProposal(6, PROPOSAL_TYPE.APPROVE);
      await createAndvoteProposal(7, PROPOSAL_TYPE.APPROVE);
      await createAndvoteProposal(8, PROPOSAL_TYPE.APPROVE);

      let proposalMap = new Map<number, number>();

      for (let accountIndex = 5; accountIndex <= 8; accountIndex++) {
        let proposalId = await createProposal(
          accountIndex,
          PROPOSAL_TYPE.REVOKE
        );
        proposalMap.set(proposalId, accountIndex);
      }

      //now have 0-8 + owener= 10 fip editor, revoke 4 ,remain 6
      let remainFipEditorCount = 6;
      let voters = new Set([0, 1, 2, 3, 4]);
      let beginVoterIndex = 5;
      //rovoke 5,6,7,8
      for (const [proposalId, accountIndex] of proposalMap.entries()) {
        voters.add(beginVoterIndex++);
        await voteProposal(proposalId, voters, accountIndex);
      }
      const fipEditorCount = await powerVotingFipEditor.fipEditorCount();
      expect(fipEditorCount).to.equal(remainFipEditorCount);
    });
    it("should reject invalid proposal ID", async function () {
      await expect(
        powerVotingFipEditor
          .connect(otherAccounts[0])
          .voteFipEditorProposal(999)
      ).to.be.revertedWithCustomError(
        powerVotingFipEditor,
        "InvalidApprovalProposalId"
      );
    });

    it("should prevent duplicate votes", async function () {
      await expect(
        powerVotingFipEditor.connect(owner).voteFipEditorProposal(proposalId)
      ).to.be.revertedWithCustomError(
        powerVotingFipEditor,
        "AddressHasActiveProposalError"
      );
    });

    it("should prevent self-vote", async function () {
      powerVotingFipEditor
        .connect(otherAccounts[0])
        .voteFipEditorProposal(proposalId);

      //
      await powerVotingFipEditor
        .connect(owner)
        .createFipEditorProposal(
          otherAccounts[1].address,
          "Test Candidate",
          PROPOSAL_TYPE.REVOKE
        );
      proposalId = await powerVotingFipEditor.fipEditorProposalId();
      await expect(
        powerVotingFipEditor
          .connect(otherAccounts[1])
          .voteFipEditorProposal(proposalId)
      ).to.be.revertedWithCustomError(
        powerVotingFipEditor,
        "CannotVoteForOwnProposalError"
      );
    });
  });
});
