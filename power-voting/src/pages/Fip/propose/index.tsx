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
import {RadioGroup} from '@headlessui/react';
import classNames from 'classnames';
import Table from '../../../components/Table';
import LoadingButton from '../../../components/LoadingButton';
import {useAccount} from "wagmi";

const FipPropose = () => {
  const {isConnected, address, chain} = useAccount();
  const chainId = chain?.id || 0;
  console.log(chainId);
  const navigate = useNavigate();
  const prevAddressRef = useRef(address);
  const [fipAddress, setFipAddress] = useState('');
  const [fipInfo, setFipInfo] = useState('');
  const [proposeType, setProposeType] = useState('approve');
  const [loading, setLoading] = useState(false);
  console.log(fipAddress);
  console.log(fipInfo);
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

  const handleProposeTypeChange = (value: string) => {
    setProposeType(value);
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
          title='FIP Editor Propose'
          list={[
            {
              name: 'Propose Type',
              width: 100,
              comp: (
                <RadioGroup className='flex' value={proposeType} onChange={handleProposeTypeChange}>
                  <RadioGroup.Option
                    key='approve'
                    value='approve'
                    className='relative flex items-center cursor-pointer p-4 focus:outline-none'
                  >
                    {({active, checked}) => (
                      <>
                        <span
                          className={classNames(
                            checked
                              ? 'bg-[#45B753] border-transparent'
                              : 'bg-[#212B3B] border-[#38485C]',
                            active ? 'ring-2 ring-offset-2 ring-[#45B753]' : '',
                            'mt-0.5 h-4 w-4 shrink-0 cursor-pointer rounded-full border flex items-center justify-center'
                          )}
                          aria-hidden='true'
                        >
                          {(active || checked) && (
                            <span className='rounded-full bg-white w-1.5 h-1.5'/>
                          )}
                        </span>
                        <span className='ml-3'>
                          <RadioGroup.Label
                            as='span'
                            className={
                              checked ? 'text-white' : 'text-[#8896AA]'
                            }
                          >
                            Approve
                          </RadioGroup.Label>
                        </span>
                      </>
                    )}
                  </RadioGroup.Option>
                  <RadioGroup.Option
                    key='revoke'
                    value='revoke'
                    className='relative flex items-center cursor-pointer p-4 focus:outline-none'
                  >
                    {({active, checked}) => (
                      <>
                        <span
                          className={classNames(
                            checked
                              ? 'bg-[#45B753] border-transparent'
                              : 'bg-[#212B3B] border-[#38485C]',
                            active ? 'ring-2 ring-offset-2 ring-[#45B753]' : '',
                            'mt-0.5 h-4 w-4 shrink-0 cursor-pointer rounded-full border flex items-center justify-center'
                          )}
                          aria-hidden='true'
                        >
                          {(active || checked) && (
                            <span className='rounded-full bg-white w-1.5 h-1.5'/>
                          )}
                        </span>
                        <span className='ml-3'>
                          <RadioGroup.Label
                            as='span'
                            className={
                              checked ? 'text-white' : 'text-[#8896AA]'
                            }
                          >
                            Revoke
                          </RadioGroup.Label>
                        </span>
                      </>
                    )}
                  </RadioGroup.Option>
                </RadioGroup>
              )
            },
            {
              name: 'Editor Address',
              width: 100,
              comp: (
                <input
                  placeholder='Input editor address'
                  className='form-input w-[520px] rounded bg-[#212B3C] border border-[#313D4F]'
                  onChange={(e) => { handleChange('fipAddress', e.target.value) }}
                />
              )
            },
            {
              name: 'Propose Info',
              width: 100,
              comp: (
                <textarea
                  placeholder='Input propose info'
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
