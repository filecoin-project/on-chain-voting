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

import { message } from 'antd';
import React, { useEffect, useRef, useState } from "react";
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from "react-router-dom";
import type { BaseError } from "wagmi";
import { useAccount, useWaitForTransactionReceipt, useWriteContract } from "wagmi";
import { filecoinCalibration } from "wagmi/chains";
import oraclePowerAbi from "../../common/abi/oracle-powers.json";
import fileCoinAbi from "../../common/abi/power-voting.json";
import {
  DUPLICATED_MINER_ID_MSG,
  STORING_DATA_MSG,
  WRONG_MINER_ID_MSG
} from "../../common/consts";
import { useMinerIdSet, useOwnerDataSet } from "../../common/hooks";
import Loading from "../../components/Loading";
import LoadingButton from '../../components/LoadingButton';
import Table from '../../components/Table';
import { getContractAddress, hasDuplicates } from "../../utils";
const MinerId = () => {
  const { chain, isConnected, address } = useAccount();
  const chainId = chain?.id || 0;
  const navigate = useNavigate();
  const { t } = useTranslation();
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
        content: t(STORING_DATA_MSG),
      });
      setTimeout(() => {
        navigate("/home");
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
        content: t(DUPLICATED_MINER_ID_MSG),
      });
      return;
    }

    // Set loading state to true to indicate loading
    setLoading(true);

    const { value, hasError } = removeMinerIdPrefix(minerIds);
    // Remove prefix from miner IDs and check for errors
    if (hasError) {
      console.log(hasError)
      messageApi.open({
        type: 'warning',
        content: t(WRONG_MINER_ID_MSG),
      });
      setLoading(false);
      return;
    }
    try {
      // Check if all requests were successful
      let allSuccessful = ownerData?.every((res: any) => {
        return res.status === 'success';
      });
      //if empty ,ignore check
      if (!minerIds.length) {
        allSuccessful = true
      }
      if (!allSuccessful) {
        messageApi.open({
          type: 'error',
          content: t(WRONG_MINER_ID_MSG),
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
        <div className="inline-flex items-center mb-8 gap-1 text-skin-text hover:text-skin-link">
          <Link to="/home" className="flex items-center">
            <svg className="mr-1" viewBox="0 0 24 24" width="1.2em" height="1.2em">
              <path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                d="m11 17l-5-5m0 0l5-5m-5 5h12" />
            </svg>
            {t('content.back')}
          </Link>
        </div>
      </button>
      <div className='flow-root space-y-8'>
        <Table
          title={t('content.minerIDsManagement')}
          list={[
            {
              name: 'Miner IDs',
              width: 100,
              comp: (
                <textarea
                  defaultValue={minerIds}
                  placeholder={t('content.inputMinerIDmultiplEseparate')}
                  className='form-input h-[320px] w-full rounded bg-[#ffffff] border border-[#EEEEEE] text-black'
                  onBlur={e => { handleMinerChange(e.target.value) }}
                />
              )
            }
          ]}
        />

        <div className='text-center'>
          <LoadingButton text={t('content.submit')} loading={loading || writeContractPending || transactionLoading} handleClick={onSubmit} />
        </div>
      </div>
    </div>
  )
}

export default MinerId;
