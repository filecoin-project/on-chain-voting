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
import {useNetwork, useAccount} from "wagmi";
import {
  DUPLICATED_MINER_ID_MSG,
  filecoinCalibrationChain,
  STORING_DATA_MSG,
  WRONG_MINER_ID_MSG
} from "../../common/consts";
import {useDynamicContract, useStaticContract} from "../../hooks";
import Loading from "../../components/Loading";
import {hasDuplicates} from "../../utils";

const MinerId = () => {
  const {chain} = useNetwork();
  const chainId = chain?.id || 0;
  const {isConnected, address} = useAccount();
  const navigate = useNavigate();
  const prevAddressRef = useRef(address);
  const [minerIds, setMinerIds] = useState(['']);
  const [spinning, setSpinning] = useState(true);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!isConnected) {
      navigate("/home");
      return;
    }
  }, []);

  useEffect(() => {
    const prevAddress = prevAddressRef.current;
    if (prevAddress !== address) {
      navigate("/home");
    }
  }, [address]);

  useEffect(() => {
    initState();
  }, [chain]);


  const initState = async () => {
    // Fetch the miner IDs from the static contract based on the chain ID
    const { getMinerIds } = await useStaticContract(chainId);
    const { code, data: { minerIds } } = await getMinerIds(address);

    // If the response code is 200, set the miner IDs after adding the prefix
    if (code === 200) {
      setMinerIds(addMinerIdPrefix(minerIds?.map((id: any) => id.toNumber())));
    }

    // Stop the spinner indicating loading
    setSpinning(false);
  }

  /**
   * Handle changes in the miner IDs input field
   * @param value
   */
  const handleMinerChange = (value: string) => {
    const arr = value ? value.split(',') : [];
    setMinerIds(arr);
  }

  /**
   * Get the prefix for a given chain ID
   * @param chainId
   */
  const getMinerIdPrefix = (chainId: number) => {
    return chainId === filecoinCalibrationChain.id ? 't0' : 'f0';
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
    // Set loading state to true to indicate loading
    setLoading(true);

    // Fetch static and dynamic contract methods
    const { getMinerIdOwner } = await useStaticContract(chainId);
    const { addMinerId } = useDynamicContract(chainId);

    // Check for duplicate miner IDs
    if (minerIds.length && hasDuplicates(minerIds)) {
      message.error(DUPLICATED_MINER_ID_MSG, 3);
      setLoading(false);
      return;
    }

    // Remove prefix from miner IDs and check for errors
    const { value, hasError } = removeMinerIdPrefix(minerIds);
    if (minerIds.length > 0) {
      if (hasError) {
        message.warning(WRONG_MINER_ID_MSG);
        setLoading(false);
        return;
      }

      try {
        // Fetch owner information for each miner ID
        const promises = value.map(async (item) => {
          const res = await getMinerIdOwner(item);
          return res;
        });
        const results = await Promise.all(promises);

        // Check if all requests were successful
        const allSuccessful = results.every((res) => {
          return res.code === 200;
        });

        if (!allSuccessful) {
          message.error(WRONG_MINER_ID_MSG, 3);
          setLoading(false);
          return;
        }
      } catch (error) {
        console.log(error);
      }
    }

    // Add miner IDs to the dynamic contract
    const res = await addMinerId(value);

    if (res.code === 200 && res.data?.hash) {
      // Show success message and navigate to home page after a delay
      message.success(STORING_DATA_MSG, 3);
      setTimeout(() => {
        navigate("/");
      }, 3000);
    } else {
      message.error(res.msg, 3);
    }

    setLoading(false);
  }

  return (
    spinning ? <Loading /> : <div className="px-3 mb-6 md:px-0">
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
                  onChange={(e) => { handleMinerChange(e.target.value) }}
                />
              )
            }
          ]}
        />

        <div className='text-center'>
          <LoadingButton text='Submit' loading={loading} handleClick={onSubmit} />
        </div>
      </div>
    </div>
  )
}

export default MinerId;
