import { useReadContract, useReadContracts } from 'wagmi';
import {getContractAddress} from "../utils";
import fileCoinAbi from "./abi/power-voting.json";
import oracleAbi from "./abi/oracle.json";

export const useCheckFipAddress = (chainId: number, address: `0x${string}` | undefined) => {
  const { data: isFipAddress } = useReadContract({
    address: getContractAddress(chainId, 'powerVoting'),
    abi: fileCoinAbi,
    functionName: 'fipMap',
    args: [address]
  });
  return {
    isFipAddress
  };
}

export const useLatestId = (chainId: number) => {
  const { data: latestId, isLoading: getLatestIdLoading } = useReadContract({
    address: getContractAddress(chainId, 'powerVoting'),
    abi: fileCoinAbi,
    functionName: 'proposalId',
  });
  return {
    latestId,
    getLatestIdLoading
  };
}

export const useProposalDataSet = (params: any) => {
  const { chainId, total, page, pageSize } = params;
  const contracts: any[] = [];
  const offset = (page - 1) * pageSize;
  // Generate contract calls for fetching proposals based on pagination
  for (let i = total - offset; i > Math.max(total - offset - pageSize, 0); i--) {
    contracts.push({
      address: getContractAddress(chainId, 'powerVoting'),
      abi: fileCoinAbi,
      functionName: 'idToProposal',
      args: [i],
    });
  }
  const {
    data: proposalData,
    isLoading: getProposalIdLoading,
    isSuccess: getProposalIdSuccess,
    error,
  } = useReadContracts({
    contracts: contracts,
    query: { enabled: !!contracts.length }
  });
  return {
    proposalData: proposalData || [],
    getProposalIdLoading,
    getProposalIdSuccess,
    error,
  };
}

export const useMinerIdSet = (chainId: number, address: `0x${string}` | undefined) => {
  const { data: minerIdData, isLoading: getMinerIdsLoading, isSuccess: getMinerIdsSuccess } = useReadContract({
    address: getContractAddress(chainId || 0, 'oracle'),
    abi: oracleAbi,
    functionName: 'getVoterInfo',
    args: [address]
  });
  return {
    minerIdData: minerIdData as any,
    getMinerIdsLoading,
    getMinerIdsSuccess
  }
}

export const useOwnerDataSet = (contracts: any[]) => {
  const {
    data: ownerData,
    isLoading: getOwnerLoading,
    isSuccess: getOwnerSuccess
  } = useReadContracts({
    contracts: contracts,
    query: { enabled: !!contracts.length }
  });

  return {
    ownerData: ownerData || [],
    getOwnerLoading,
    getOwnerSuccess,
  };
}