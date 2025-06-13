import { CORE_ORG_REPO, ECOSYSTEM_ORG, GITHUB_USER, POWER_VOTING_CONFIG } from "./constant";
import { getConstantJson } from "./utils";
import * as conf from "./power_voting_config.json";

async function main() {
    const constantJSON = getConstantJson();
    const POWER_VOTING_CONFIG_ADDRESS = constantJSON[POWER_VOTING_CONFIG];
    if (!ethers.isAddress(POWER_VOTING_CONFIG_ADDRESS)) {
        throw new Error(`Invalid contract address: ${POWER_VOTING_CONFIG_ADDRESS}`);
    }

    console.log(`POWER_VOTING_CONFIG_ADDRESS: ${POWER_VOTING_CONFIG_ADDRESS}`);

    const powerVotingConfigContract = await ethers.getContractAt("PowerVotingConf", POWER_VOTING_CONFIG_ADDRESS);

    await init_repo(powerVotingConfigContract);
    await init_algorithm(powerVotingConfigContract);
    await watchConfig(powerVotingConfigContract);
}

async function init_repo(contractAddr: any) {
    try {
        await contractAddr.batchAddGithubRepo(CORE_ORG_REPO, conf.github.coreOrg)
        console.log(`init CORE_ORG_REPO success`);

    } catch (error) {
        console.error(`init CORE_ORG_REPO failed`, error)
    }

    try {
        await contractAddr.batchAddGithubRepo(ECOSYSTEM_ORG, conf.github.ecosystemOrg)
        console.log(`init ECOSYSTEM_ORG success`);
    } catch (error) {
        console.error(`init ECOSYSTEM_ORG failed`, error)

    }

    try {
        await contractAddr.batchAddGithubRepo(GITHUB_USER, conf.github.githubUser)
        console.log(`init GITHUB_USER success`);
    } catch (error) {
        console.error(`init GITHUB_USER failed`, error)
    }
}

async function init_algorithm(contractAddr: any) {
    try {
        await contractAddr.setVotingCountingAlgorithm(conf.votingCountingAlgorithm)
        console.log(`init algorithm success`);
    } catch (error) {
        console.error(`init algorithm failed`, error)
    }
}

async function watchConfig(contractAddr: any) {
    try {
        const currentId = await contractAddr.githubRepoId();
        console.log(`core: ${(currentId)}`);

    } catch (error) {
        console.error(`watchConfig failed`, error)
    }
}

main().then(() => process.exit(0)).catch(error => {
    console.error(error);
    process.exit(1);
})