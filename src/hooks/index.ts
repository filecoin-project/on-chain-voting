// Copyright (C) 2023-2024 StorSwift Inc.
// This file is part of the PowerVoting library.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { ethers } from "ethers";
import { NFTStorage, Blob } from 'nft.storage';
// import { filesFromPaths } from 'files-from-path';
import { create } from '@web3-storage/w3up-client';
import fileCoinAbi from "../common/abi/power-voting.json";
import oracleAbi from "../common/abi/oracle.json";
import oraclePowerAbi from "../common/abi/oracle-powers.json";
import {
  powerVotingMainNetContractAddress,
  oracleCalibrationContractAddress,
  contractAddressList,
  oracleMainNetContractAddress,
  SUCCESS_INFO,
  ERROR_INFO,
  NFT_STORAGE_KEY,
  OPERATION_FAILED_MSG, STORING_DATA_MSG,
  oraclePowerCalibrationContractAddress,
  oraclePowerMainNetContractAddress,
  web3StorageEmail,
} from "../common/consts";
import { filecoin, filecoinCalibration } from 'wagmi/chains';
import {extractRevertReason} from "../utils";

const decodeError = (data: string) => {
  const errorData = data.substring(0, 2) + data.substring(10);
  const abiCoder = new ethers.utils.AbiCoder();
  let decodedData =  '';
  try {
    decodedData = abiCoder.decode(['string'], errorData)[0];
  } catch (e) {
    console.log(e);
  }
  return decodedData;
}

// @ts-ignore
const handleReturn = ({ type, data }) => {
  let code = 200;
  let msg = STORING_DATA_MSG;
  if (type === ERROR_INFO) {
    const encodeData = data?.error?.data?.message;
    const reason: any = extractRevertReason(encodeData);
    if (encodeData) {
      code = 401
      msg = decodeError(reason) || OPERATION_FAILED_MSG;
    } else {
      code = 402;
      msg = OPERATION_FAILED_MSG;
    }
  }

  return {
    code,
    msg,
    data
  }
}

export const useStaticContract = async (chainId: number) => {
  const rpcUrl = chainId === filecoin.id ? filecoin.rpcUrls.default.http[0] : filecoinCalibration.rpcUrls.default.http[0];
  const provider = new ethers.providers.JsonRpcProvider(rpcUrl);

  const powerVotingContractAddress = contractAddressList.find(item => item.id === chainId)?.address || powerVotingMainNetContractAddress;
  const powerVotingContract = new ethers.Contract(powerVotingContractAddress, fileCoinAbi, provider);

  const oracleContractAddress = chainId === filecoin.id ? oracleMainNetContractAddress : oracleCalibrationContractAddress;
  const oracleContract = new ethers.Contract(oracleContractAddress, oracleAbi, provider);

  const oraclePowerContractAddress = chainId === filecoin.id ? oraclePowerMainNetContractAddress : oraclePowerCalibrationContractAddress;
  const oraclePowerContract = new ethers.Contract(oraclePowerContractAddress, oraclePowerAbi, provider);

  /**
   * get latest proposal id
   */
  const getLatestId = async () => {
    try {
      const data = await powerVotingContract.proposalId();
      return handleReturn({
        type: SUCCESS_INFO,
        data
      })
    } catch (e) {
      return handleReturn({
        type: ERROR_INFO,
        data: e
      })
    }
  }

  /**
   * check FIP editor
   */
  const isFipEditor = async (address: string) => {
    try {
      const data = await powerVotingContract.fipMap(address);
      return handleReturn({
        type: SUCCESS_INFO,
        data
      })
    } catch (e) {
      return handleReturn({
        type: ERROR_INFO,
        data: e
      })
    }
  }

  /**
   * get proposal detail
   * @param id
   *
   */
  const getProposal = async (id: number) => {
    try {
      const data = await powerVotingContract.idToProposal(id);
      return handleReturn({
        type: SUCCESS_INFO,
        data
      })
    } catch (e) {
      console.log(e);
      return handleReturn({
        type: ERROR_INFO,
        data: e
      })
    }
  }

  /**
   * get UCAN data
   * @param address
   */
  const getOracleAuthorize = async (address: any) => {
    try {
      const data = await oracleContract.voterToInfo(address);
      return handleReturn({
        type: SUCCESS_INFO,
        data
      })
    } catch (e) {
      console.log(e);
      return handleReturn({
        type: ERROR_INFO,
        data: e
      })
    }
  }

  /**
   * get miner ID
   * @param address
   */
  const getMinerIds = async (address: any) => {
    try {
      const data = await oracleContract.getVoterInfo(address);
      return handleReturn({
        type: SUCCESS_INFO,
        data
      })
    } catch (e) {
      console.log(e);
      return handleReturn({
        type: ERROR_INFO,
        data: e
      })
    }
  }

  /**
   * get miner ID owner
   * @param minerId
   */
  const getMinerIdOwner = async (minerId: number) => {
    try {
      const data = await oraclePowerContract.getOwner(minerId);
      return handleReturn({
        type: SUCCESS_INFO,
        data
      })
    } catch (e) {
      console.log(e);
      return handleReturn({
        type: ERROR_INFO,
        data: e
      })
    }
  }

  return {
    getLatestId,
    isFipEditor,
    getProposal,
    getOracleAuthorize,
    getMinerIds,
    getMinerIdOwner,
  }
}

export const useDynamicContract = (chainId: number) => {
  const contractAddress = contractAddressList.find(item => item.id === chainId)?.address || powerVotingMainNetContractAddress;
  // @ts-ignore
  const provider = new ethers.providers.Web3Provider(window.ethereum);
  const signer = provider.getSigner();
  const contract = new ethers.Contract(contractAddress, fileCoinAbi, signer);

  /**
   * create proposal
   * @param proposalCid
   * @param startTime
   * @param expTimestamp
   * @param proposalType
   */
  const createVotingApi = async (proposalCid: string, startTime: number, expTimestamp: number, proposalType: number) => {
    try {
      const data = await contract.createProposal(proposalCid, startTime, expTimestamp, proposalType);
      return handleReturn({
        type: SUCCESS_INFO,
        data
      })
    } catch (e) {
      console.log(e);
      return handleReturn({
        type: ERROR_INFO,
        data: e
      })
    }
  }

  /**
   * proposal vote
   * @param proposalId
   * @param optionId
   */
  const voteApi = async (proposalId: number, optionId: string) => {
    try {
      const data = await contract.vote(proposalId, optionId);
      return handleReturn({
        type: SUCCESS_INFO,
        data
      })
    } catch (e) {
      return handleReturn({
        type: ERROR_INFO,
        data: e
      })
    }
  }

  /**
   * UCAN authorize or deAuthorize
   * @param ucanCid
   */
  const ucanDelegate = async (ucanCid: string) => {
    try {
      const data = await contract.ucanDelegate(ucanCid);
      return handleReturn({
        type: SUCCESS_INFO,
        data
      })
    } catch (e) {
      return handleReturn({
        type: ERROR_INFO,
        data: e
      })
    }
  }

  /**
   * Add miner Ids
   * @param minerIds
   */
  const addMinerId = async (minerIds: number[]) => {
    try {
      const data = await contract.addMinerId(minerIds);
      return handleReturn({
        type: SUCCESS_INFO,
        data
      })
    } catch (e) {
      return handleReturn({
        type: ERROR_INFO,
        data: e
      })
    }
  }

  return {
    addMinerId,
    ucanDelegate,
    createVotingApi,
    voteApi,
  }
}

/**
 * Reads an image file from `imagePath` and stores an NFT with the given name and description.
 * @param {object} params for the NFT
 */
const storeIpfs = (params: object) => {
  const json = JSON.stringify(params);
  const data = new Blob([json]);

  const nftStorage = new NFTStorage({ token: NFT_STORAGE_KEY });

  return nftStorage.storeBlob(data);
}

/**
 * The main entry point for the script that checks the command line arguments and
 * calls storeNFT.
 *
 * To simplify the example, we don't do any fancy command line parsing. Just three
 * positional arguments for imagePath, name, and description
 */
export const getIpfsId = async (props: any) => {
  const result = await storeIpfs(props);
  return result;
}


/**
 * Store data into web3.storage
 * @param params for the data
 */
export const getWeb3IpfsId = async (params: object | string) => {
  const client = await create();
  console.log(client.uploadFile);
  // first time setup!
  if (!Object.keys(client.accounts()).length) {
    // waits for you to click the link in your email to verify your identity
    const account = await client.login(web3StorageEmail)
    // create a space for your uploads
    const space = await client.createSpace('power-voting')
    // save the space to the store, and set as "current"
    await space.save()
    // associate this space with your account
    await account.provision(space.did())
  }

  const json = JSON.stringify(params);
  const data = new Blob([json]);
  const cid = await client.uploadFile(data);
  console.log(cid);
  return cid;
}