import { ethers, utils } from "ethers";
import * as zksync from "zksync-web3";
import { Provider } from "zksync-web3";
import { NFTStorage, Blob } from 'nft.storage';
import abi from "../../public/abi/power-voting.json";
import {
  filecoinMainnetContractAddress,
  contractAddressList,
  NFT_STORAGE_KEY,
} from "../common/consts";

export const useDynamicContract = (chainId: number) => {
  const contractAddress = contractAddressList.find((item: any) => item.id === chainId)?.address || filecoinMainnetContractAddress;

  // @ts-ignore
  const provider = new ethers.providers.Web3Provider(window.ethereum);
  const signer = provider.getSigner();

  const contract = new ethers.Contract(contractAddress, abi, signer);

  const decodeError = (data: string) => {
    const errorData = data.substring(0, 2) + data.substring(10);
    const defaultAbiCoder = new ethers.utils.AbiCoder();
    const decodedData =  defaultAbiCoder.decode(['string'], errorData)[0];
    return decodedData;
  }

  const handleReturn = (params: any) => {
    const { type, data } = params
    let code = 200;
    let msg = 'success';
    if (type === 'error') {
      console.log(data?.error?.data);
      const encodeData = data?.error?.data?.originalError?.data;
      if (encodeData) {
        code = 401
        msg = decodeError(encodeData);
      } else {
        code = 402
      }
    }
    return {
      code,
      msg,
      data
    }
  }

  const createVotingApi = async (proposalCid: string, timestamp: number, chainId:number, proposalType: number) => {
    try {
      const data = await contract.createProposal(proposalCid, timestamp, chainId, proposalType);
      return handleReturn({
        type: 'success',
        data
      })
    } catch (e: any) {
      return handleReturn({
        type: 'error',
        data: e
      })
    }
  }

  const voteApi = async (id: number, optionId: string) => {
    try {
      const data = await contract.vote(id, optionId);
      return handleReturn({
        type: 'success',
        data
      })
    } catch (e: any) {
      return handleReturn({
        type: 'error',
        data: e
      })
    }
  }

  const cancelVotingApi = async (id: number) => {
    try {
      const data = await contract.cancelProposal(id);
      return handleReturn({
        type: 'success',
        data
      })
    } catch (e: any) {
      return handleReturn({
        type: 'error',
        data: e
      })
    }
  }

  const zkSyncDepositApi = async () => {
    const zkSyncProvider = new Provider("https://testnet.era.zksync.dev");
    const ethProvider = ethers.getDefaultProvider("goerli");
    // const signingKey = await signer.getSigningKey();
    //console.log(signingKey);
    // const zkSyncWallet = new zksync.Wallet(signingKey, zkSyncProvider, ethProvider);


    // const deposit = await zkSyncWallet.deposit({
    //   token: zksync.utils.ETH_ADDRESS,
    //   amount: ethers.utils.parseEther("1.0"),
    // });
    //
    // return deposit;
  }

  return {
    createVotingApi,
    cancelVotingApi,
    voteApi,
    zkSyncDepositApi,
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