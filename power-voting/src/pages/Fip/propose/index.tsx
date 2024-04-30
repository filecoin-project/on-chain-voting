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
import Table from '../../../components/Table';
import LoadingButton from '../../../components/LoadingButton';
import {useNetwork, useAccount} from "wagmi";
import {
  DUPLICATED_MINER_ID_MSG,
  filecoinCalibrationChain,
  STORING_DATA_MSG,
  WRONG_MINER_ID_MSG
} from "../../../common/consts";
import {useDynamicContract, useStaticContract} from "../../../hooks";
import Loading from "../../../components/Loading";
import {hasDuplicates} from "../../../utils";

const FipPropose = () => {
  const {chain} = useNetwork();
  const chainId = chain?.id || 0;
  const {isConnected, address} = useAccount();
  const navigate = useNavigate();
  const prevAddressRef = useRef(address);
  const [fipAddress, setFipAddress] = useState('');
  const [fipInfo, setFipInfo] = useState('');
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

  const handleChange = (type: string, value: string) => {
    type === 'fipAddress' ? setFipAddress(value) : setFipInfo(value);
  }

  /**
   * Set miner ID
   */
  const onSubmit = async () => {
    setLoading(true);

    setLoading(false);
  }

  return (
    <div className="px-3 mb-6 md:px-0">
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
          title='FIP Propose'
          list={[
            {
              name: 'FIP Address',
              width: 100,
              comp: (
                <input
                  placeholder='Input FIP Address.'
                  className='form-input w-[520px] rounded bg-[#212B3C] border border-[#313D4F] cursor-not-allowed'
                  onChange={(e) => { handleChange('fipAddress', e.target.value) }}
                />
              )
            },
            {
              name: 'Address Info.',
              width: 100,
              comp: (
                <textarea
                  placeholder='Input Address Info.'
                  className='form-input h-[320px] w-full rounded bg-[#212B3C] border border-[#313D4F]'
                  onChange={(e) => { handleChange('fipInfo', e.target.value) }}
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

export default FipPropose;
