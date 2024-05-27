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

import React, {useState, useEffect, useRef} from "react";
import {Link, useNavigate} from "react-router-dom";
import { message } from 'antd';
import Table from '../../components/Table';
import LoadingButton from '../../components/LoadingButton';
import type { BaseError} from "wagmi";
import {useAccount, useReadContract, useReadContracts, useWriteContract, useWaitForTransactionReceipt} from "wagmi";
import {filecoinCalibration} from "wagmi/chains";
import {
  DUPLICATED_MINER_ID_MSG,
  STORING_DATA_MSG,
  WRONG_MINER_ID_MSG
} from "../../common/consts";
import Loading from "../../components/Loading";
import {getContractAddress, hasDuplicates} from "../../utils";
import fileCoinAbi from "../../common/abi/power-voting.json";
import oracleAbi from "../../common/abi/oracle.json";
import oraclePowerAbi from "../../common/abi/oracle-powers.json";

function useMinerIdSet(chainId: number, address: `0x${string}` | undefined) {
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

function useOwnerDataSet(contracts: any[]) {
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

const MinerId = () => {
  const {chain, isConnected, address} = useAccount();
  const chainId = chain?.id || 0;
  const navigate = useNavigate();
  const prevAddressRef = useRef(address);
  const [minerIds, setMinerIds] = useState(['']);
  const [contracts, setContracts] = useState([] as any);

  const { minerIdData, getMinerIdsLoading, getMinerIdsSuccess } = useMinerIdSet(chainId, address);
  const { ownerData } = useOwnerDataSet(contracts);
  const [messageApi, contextHolder] = message.useMessage()

  const {
    data: hash,
    writeContract,
    error,
    isPending: writeContractPending,
    isSuccess: writeContractSuccess,
    reset
  } = useWriteContract();

  const [loading, setLoading] = useState<boolean>(writeContractPending);

  useEffect(() => {
    if (error) {
      messageApi.open({
        type: 'error',
        content: (error as BaseError)?.shortMessage || error?.message,
      });
    }
    reset();
  }, [error]);

  useEffect(() => {
    if (!isConnected) {
      navigate("/home");
      return;
    }
  }, [isConnected]);

  useEffect(() => {
    const prevAddress = prevAddressRef.current;
    if (prevAddress !== address) {
      navigate("/home");
    }
  }, [address]);

  useEffect(() => {
    initState();
  }, [chain, getMinerIdsSuccess]);

  useEffect(() => {
    if (writeContractSuccess) {
      messageApi.open({
        type: 'success',
        content: STORING_DATA_MSG,
      });
      setTimeout(() => {
        navigate("/");
      }, 3000);
    }
  }, [writeContractSuccess])

  const initState = async () => {
    setMinerIds(addMinerIdPrefix(minerIdData?.minerIds?.map((id: any) => Number(id))));
  }

  /**
   * Handle changes in the miner IDs input field
   * @param ids
   */
  const handleMinerChange = (ids: string) => {
    const arr = ids ? ids.split(',') : [];
    setMinerIds(arr);

    const { value } = removeMinerIdPrefix(arr);
    const contracts = value.map((item: number) => {
      return {
        address: getContractAddress(chain?.id || 0, 'oraclePower'),
        abi: oraclePowerAbi,
        functionName: 'getOwner',
        args: [item],
      } as const;
    });
    setContracts(contracts);
  }

  /**
   * Get the prefix for a given chain ID
   * @param chainId
   */
  const getMinerIdPrefix = (chainId: number) => {
    return chainId === filecoinCalibration.id ? 't0' : 'f0';
  }

  /**
   * Add prefix to each miner ID based on the chain ID
   * @param minerIds
   */
  const addMinerIdPrefix = (minerIds: number[]) => {
    return minerIds?.length ? minerIds.map(minerId => `${getMinerIdPrefix(chainId)}${minerId}`) : [];
  }

  /**
   * Remove prefix from each miner ID and validate the format
   * @param minerIds
   */
  const removeMinerIdPrefix = (minerIds: string[]) => {
    const prefix = getMinerIdPrefix(chainId);
    const prefixRegex = new RegExp('^' + prefix);
    let hasError = false;

    const arr = minerIds?.length > 0 ? minerIds.map(minerId => {
      const str = minerId.replace(prefixRegex, '');
      if (isNaN(Number(str)) || str?.length > 7) {
        hasError = true;
      }
      return Number(str);
    }) : [];

    return {
      value: arr,
      hasError
    };
  }

  /**
   * Set miner ID
   */
  const onSubmit = async () => {
    // Check for duplicate miner IDs
    if (minerIds.length && hasDuplicates(minerIds)) {
      messageApi.open({
        type: 'error',
        content: DUPLICATED_MINER_ID_MSG,
      });
      return;
    }

    // Set loading state to true to indicate loading
    setLoading(true);

    const { value, hasError } = removeMinerIdPrefix(minerIds);
    // Remove prefix from miner IDs and check for errors
    if (hasError) {
      messageApi.open({
        type: 'warning',
        content: WRONG_MINER_ID_MSG,
      });
      setLoading(false);
      return;
    }
    try {
      // Check if all requests were successful
      const allSuccessful = ownerData?.every((res: any) => {
        return res.status === 'success';
      });

      if (!allSuccessful) {
        messageApi.open({
          type: 'error',
          content: WRONG_MINER_ID_MSG,
        });
        setLoading(false);
        return;
      } else {
        writeContract({
          abi: fileCoinAbi,
          address: getContractAddress(chain?.id || 0, 'powerVoting'),
          functionName: 'addMinerId',
          args: [
            value
          ],
        });
      }
    } catch (error) {
      console.log(error);
    }
    setLoading(false);
  }

  const { isLoading: transactionLoading } =
    useWaitForTransactionReceipt({
      hash,
    })

  return (
    getMinerIdsLoading ? <Loading /> : <div className="px-3 mb-6 md:px-0">
      {contextHolder}
      <button>
        <div className="inline-flex items-center gap-1 text-skin-text hover:text-skin-link">
          <Link to="/" className="flex items-center">
            <svg className="mr-1" viewBox="0 0 24 24" width="1.2em" height="1.2em">
              <path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                    d="m11 17l-5-5m0 0l5-5m-5 5h12" />
            </svg>
            Back
          </Link>
        </div>
      </button>
      <div className='flow-root space-y-8'>
        <Table
          title='Miner IDs Management'
          list={[
            {
              name: 'Miner IDs',
              width: 100,
              comp: (
                <textarea
                  defaultValue={minerIds}
                  placeholder='Input miner ID (For multiple miner IDs, use commas to separate them.)'
                  className='form-input h-[320px] w-full rounded bg-[#212B3C] border border-[#313D4F]'
                  onBlur={e => { handleMinerChange(e.target.value) }}
                />
              )
            }
          ]}
        />

        <div className='text-center'>
          <LoadingButton text='Submit' loading={loading || writeContractPending || transactionLoading} handleClick={onSubmit} />
        </div>
      </div>
    </div>
  )
}

export default MinerId;
