import { expect } from "chai";
import { ethers } from "hardhat";
import { upgrades } from "hardhat";


// const 
describe("PowerVotingConf", function () {
    let powerVotingConf: any, owner, user: any;
    beforeEach(async function () {
        [owner, user] = await ethers.getSigners();
        const PowerVotingConf = await ethers.getContractFactory("PowerVotingConf");
        const confProxy = await upgrades.deployProxy(PowerVotingConf, { initializer: 'initialize' });
        powerVotingConf = await confProxy.waitForDeployment();
    })

    describe('batchAddGithubRepo', () => {
        it('should add github repo', async () => {
            await expect(powerVotingConf.batchAddGithubRepo(0, ["repo1"]))
                .to.emit(powerVotingConf, 'GithubRepoAdded')
                .withArgs(1, ['repo1', 0]);
        })
        it("should revert if not owner", async () => {
            await expect(powerVotingConf.connect(user).batchAddGithubRepo(1, ["repo1", "repo2"])).to.be.rejected;
        })
        it("should revert if repo is empty", async () => {
            await expect(powerVotingConf.batchAddGithubRepo(1, [])).to.be.revertedWith("repo is empty");
        })
        it("should countiue if repo is already added", async () => {
            await powerVotingConf.batchAddGithubRepo(0, ["repo1", "repo2"]);
            await powerVotingConf.batchAddGithubRepo(0, ["repo2"]);
            expect(await powerVotingConf.githubRepoId()).to.be.equal(2);
        })
    })

    describe("batchRemoveGithubRepos", () => {
        beforeEach(async () => {
            await powerVotingConf.batchAddGithubRepo(0, ["repo1", "repo2", "repo3"]);
        })
        it("should remove github repo", async () => {
            await expect(powerVotingConf.batchRemoveGithubRepos([1])).
                to.emit(powerVotingConf, "GithubRepoRemoved").withArgs(1);
        })
        it("should revert if not owner", async () => {
            await expect(powerVotingConf.connect(user).
                batchRemoveGithubRepos([1, 2, 3])).to.be.rejected;
        })
        it("should revert if repo is empty", async () => {
            await expect(powerVotingConf.batchRemoveGithubRepos([])).
                to.be.revertedWith("repo is empty");
        })

        it("should continue if repo is not exist", async () => {
            expect(await powerVotingConf.batchRemoveGithubRepos([1, 4])).
                to.emit(powerVotingConf, "GithubRepoRemoved").withArgs(1);;
        })
    })

    describe("setSnapshotHeight", async function () {
        it("should set snapshot height", async function () {
            await expect(powerVotingConf.setSnapshotHeight("20250101", 123456))
                .to.be.emit(powerVotingConf, "SnapshotDays").withArgs(
                    "20250101", 123456
                );
            expect(await powerVotingConf.getSnapshotHeight("20250101")).to.be.equal(123456)
        })
        it("should revert if not owner", async function () {
            await expect(powerVotingConf.connect(user).setSnapshotHeight("20250101", 123456)).to.be.reverted;
        })
        it("should revert if dateStr length is not 8 bits", async function () {
            await expect(powerVotingConf.setSnapshotHeight("202501", 123456)).to.be.rejectedWith("Invalid date format");
        })
    })

    describe("setExpDays", async () => {
        it("should set exp days", async () => {
            await expect(powerVotingConf.setExpDays(120))
                .to.be.emit(powerVotingConf, "SnapshotExpirationDay").withArgs(
                    60, 120
                );
            expect(await powerVotingConf.expDays()).to.be.equal(120)
        })
        it("should revert if not owner", async () => {
            await expect(powerVotingConf.connect(user).setExpDays(120)).to.be.reverted;
        })
        it("should revert if exp days is less than 1 or more than 180", async () => {
            await expect(powerVotingConf.setExpDays(0)).to.be.revertedWith("Invalid expiration days");
            await expect(powerVotingConf.setExpDays(181)).to.be.revertedWith("Invalid expiration days");
        })
    })

    describe("setVotingCountingAlgorithm", () => {
        it("should set voting counting algorithm", async () => {
            await powerVotingConf.setVotingCountingAlgorithm("y = x^2");
            expect(await powerVotingConf.votingCountingAlgorithm()).to.be.equal("y = x^2")
        })
    })
})