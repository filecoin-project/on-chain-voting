
import { POWER_VOTING_CONFIG } from "./constant";
import { updateConstant } from "./utils";
const { ethers, upgrades } = require("hardhat");
async function main() {
    console.log("PowerVotingConf deploy beging");
    const PowerVotingConfig = await ethers.getContractFactory("PowerVotingConf");
    const contract = await upgrades.deployProxy(PowerVotingConfig, {
        initializer: "initialize"
    });
    await contract.waitForDeployment();
    const address = await contract.getAddress();
    updateConstant(POWER_VOTING_CONFIG, address);
    console.log("PowerVotingConf deployed to:", address);
}

main().catch(console.error);
