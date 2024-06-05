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

import { useReadContract, useReadContracts } from 'wagmi';
import {getContractAddress} from "../utils";
import fileCoinAbi from "./abi/power-voting.json";
import oracleAbi from "./abi/oracle.json";

export const useVoterInfoSet = (chainId: number, address: `0x${string}` | undefined) => {
  const { data: voterInfo } = useReadContract({
    address: getContractAddress(chainId, 'oracle'),
    abi: oracleAbi,
    functionName: 'voterToInfo',
    args: [address]
  });
  return {
    voterInfo: voterInfo as any
  }
}

export const useCheckFipEditorAddress = (chainId: number, address: `0x${string}` | undefined) => {
  const { data: isFipEditorAddress, isSuccess: checkFipEditorAddressSuccess } = useReadContract({
    address: getContractAddress(chainId, 'powerVoting'),
    abi: fileCoinAbi,
    functionName: 'fipAddressMap',
    args: [address]
  });
  return {
    isFipEditorAddress,
    checkFipEditorAddressSuccess
  };
}

export const useLatestId = (chainId: number, enabled: boolean) => {
  const { data: latestId, isLoading: getLatestIdLoading, refetch } = useReadContract({
    address: getContractAddress(chainId, 'powerVoting'),
    abi: fileCoinAbi,
    functionName: 'proposalId',
    query: {
      enabled
    }
  });
  return {
    latestId,
    getLatestIdLoading,
    refetch
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
    query: { enabled: contracts.length > 0 }
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
    ownerData: ownerData as any[],
    getOwnerLoading,
    getOwnerSuccess,
  };
}

export const useFipEditors = (chainId: number) => {
  const { data: fipEditors } = useReadContract({
    address: getContractAddress(chainId, 'powerVoting'),
    abi: fileCoinAbi,
    functionName: 'getFipAddressList',
  });
  return {
    fipEditors: fipEditors as string[]
  };
}

export const useApproveProposalId = (chainId: number) => {
  const { data: approveProposalId, isLoading: getApproveProposalLoading } = useReadContract({
    address: getContractAddress(chainId, 'powerVoting'),
    abi: fileCoinAbi,
    functionName: 'getApproveProposalId',
  });
  return {
    approveProposalId,
    getApproveProposalLoading,
  } as {
    approveProposalId: string[],
    getApproveProposalLoading: boolean
  };
}

export const useRevokeProposalId = (chainId: number) => {
  const { data: revokeProposalId, isLoading: getRevokeProposalIdLoading } = useReadContract({
    address: getContractAddress(chainId, 'powerVoting'),
    abi: fileCoinAbi,
    functionName: 'getRevokeProposalId',
  });
  return {
    revokeProposalId,
    getRevokeProposalIdLoading,
  } as {
    revokeProposalId: string[],
    getRevokeProposalIdLoading: boolean
  };
}

export const useFipEditorProposalDataSet = (params: any) => {
  const { chainId, idList, page, pageSize } = params;
  const contracts: any[] = [];
  const offset = (page - 1) * pageSize;
  idList?.sort((a: bigint, b: bigint) => Number(b) - Number(a));

  // Generate contract calls for fetching proposals based on pagination
  const startIndex = offset;
  const endIndex = Math.min(startIndex + pageSize, idList?.length);
  for (let i = startIndex; i < endIndex; i++) {
    const id = Number(idList[i]);
    contracts.push({
      address: getContractAddress(chainId, 'powerVoting'),
      abi: fileCoinAbi,
      functionName: 'getFipEditorProposal',
      args: [
        id
      ],
    });
  }
  const {
    data: fipEditorProposalData,
    isLoading: getFipEditorProposalIdLoading,
    isSuccess: getFipEditorProposalIdSuccess,
    error,
  } = useReadContracts({
    contracts: contracts,
    query: { enabled: contracts.length > 0 }
  });
  return {
    fipEditorProposalData: fipEditorProposalData || [],
    getFipEditorProposalIdLoading,
    getFipEditorProposalIdSuccess,
    error,
  };
}