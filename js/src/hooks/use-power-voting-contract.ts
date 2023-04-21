import { ethers } from "ethers";

import PowerVoting from "@Contracts/PowerVoting.sol/PowerVoting.json"
import PowerVotingFRC721 from "@Contracts/PowerVotingFRC721.sol/PowerVotingFRC721.json"

const contractAddress = process.env.VOTING_CONTRACT_ADDRESS
const contractNFTAddress = process.env.NFT_CONTRACT_ADDRESS


if (!contractAddress) {
  throw new Error(
    "Please set VOTING_CONTRACT_ADDRESS in a .env file"
  )
}

if (!contractNFTAddress) {
  throw new Error(
    "Please set VOTING_CONTRACT_ADDRESS in a .env file"
  )
}

export const usePowerVotingContract = () => {

  if (!window.ethereum) {
    return {}
  } 
  const provider = new ethers.providers.Web3Provider(window.ethereum);
  const signer = provider.getSigner()
  // 创建合约实例
  const contract = new ethers.Contract(contractAddress, PowerVoting.abi, signer)
  const contractNFT = new ethers.Contract(contractNFTAddress, PowerVotingFRC721.abi, signer)

  // 封装获取投票列表方法
  const getVotingList = async () => {
    const data = await contract.votingList()
    return data
  }

  // 封装创建投票方法
  const createVotingApi = async (cid: string) => {
    const data = await contract.createVoting(cid)
    return data
  }

  // 投票详情
  const getVoteApi = async (cid: string) => {
    const data = await contract.getVote(cid)
    return data
  }
  // 投票
  const voteApi = async (id: string, cid: string, address: string) => {
    const data = await contract.vote(id, cid, address)
    return data
  }

  const VotingNFT = async () => {
    const data = await contractNFT.mintPowerVotingNFT()
    return data
  }
  // 获取投票项目数据
  const getVoteDataApi = async (cid: string) => {
    const data = await contract.getVoteData(cid)
    return data
  }

  // 更新投票结果
  const updateVotingResultFun = async (_cid: string, cid: string) => {
    // _cid 为投票的cid cid为ipfs返回的cid
    const data = await contract.updateVotingResult(_cid, cid)
    return data
  }

  // 判断是否需要计票
  const isFinishVoteFun = async (_cid: string) => {
    const data = await contract.isFinishVote(_cid)
    return data
  }
  // 批量更新投票数据
  const updateVotingResultBatchFun = async (arr: string[][]) => {
    const data = await contract.updateVotingResultBatch(arr)
    return data
  }


  // 导出模块方法
  return {
    getVotingList,
    createVotingApi,
    getVoteApi,
    voteApi,
    VotingNFT,
    getVoteDataApi,
    updateVotingResultFun,
    isFinishVoteFun,
    updateVotingResultBatchFun
  }

}
