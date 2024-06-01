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
import {message} from "antd";
import {Link, useNavigate} from "react-router-dom";
import {RadioGroup} from '@headlessui/react';
import classNames from 'classnames';
import Table from '../../../components/Table';
import LoadingButton from '../../../components/LoadingButton';
import {useAccount, useWriteContract, useWaitForTransactionReceipt} from "wagmi";
import type { BaseError } from "wagmi";
import { useFipEditors } from "../../../common/hooks"
import fileCoinAbi from "../../../common/abi/power-voting.json";
import {getContractAddress, getWeb3IpfsId} from "../../../utils";
import {
  CAN_NOT_REVOKE_YOURSELF_MSG,
  FIP_APPROVE_TYPE,
  FIP_REVOKE_TYPE,
  NO_FIP_EDITOR_PROPOSAL_ADDRESS_MSG,
  STORING_DATA_MSG,
} from "../../../common/consts";

const FipPropose = () => {
  const {isConnected, address, chain} = useAccount();
  const chainId = chain?.id || 0;

  const navigate = useNavigate();
  const prevAddressRef = useRef(address);
  const [messageApi, contextHolder] = message.useMessage();


  const [fipAddress, setFipAddress] = useState('');
  const [selectedAddress, setSelectedAddress] = useState('');
  const [fipInfo, setFipInfo] = useState('');
  const [fipProposalType, setFipProposeType] = useState(FIP_APPROVE_TYPE);
  const { fipEditors } = useFipEditors(chainId);

  const {
    data: hash,
    writeContract,
    error,
    isPending: writeContractPending,
    isSuccess: writeContractSuccess,
    reset
  } = useWriteContract();

  const [loading, setLoading] = useState(writeContractPending);

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
    if (writeContractSuccess) {
      messageApi.open({
        type: 'success',
        content: STORING_DATA_MSG,
      });
      setTimeout(() => {
        navigate("/")
      }, 1000);
    }
  }, [writeContractSuccess]);

  useEffect(() => {
    if (error) {
      messageApi.open({
        type: 'error',
        content: (error as BaseError)?.shortMessage || error?.message,
      });
    }
    reset();
  }, [error]);

  const handleChange = (type: string, value: string) => {
    type === 'fipAddress' ? setFipAddress(value) : setFipInfo(value);
  }

  const handleProposeTypeChange = (value: number) => {
    setFipProposeType(value);
    value === FIP_APPROVE_TYPE ? setSelectedAddress('') : setFipAddress('');
    handleChange('fipInfo', '');
  }

  /**
   * Set miner ID
   */
  const onSubmit = async () => {
    // Check if required fields are filled based on proposal type
    if ((fipProposalType === FIP_APPROVE_TYPE && !fipAddress) || (fipProposalType === FIP_REVOKE_TYPE && !selectedAddress)) {
      messageApi.open({
        type: 'warning',
        // Prompt user to fill required fields
        content: NO_FIP_EDITOR_PROPOSAL_ADDRESS_MSG,
      });
      return;
    }

    // Check if the user revoke himself
    if (fipProposalType === FIP_REVOKE_TYPE && selectedAddress === address) {
      messageApi.open({
        type: 'warning',
        // Prompt user to fill required fields
        content: CAN_NOT_REVOKE_YOURSELF_MSG,
      });
      return;
    }

    // Set loading state to true while submitting proposal
    setLoading(true);

    // Get the IPFS CID for the proposal information
    const cid = await getWeb3IpfsId(fipInfo);

    // Construct the arguments and call the writeContract function to create the proposal
    const proposalArgs = [
      // Use appropriate address based on proposal type
      fipProposalType === FIP_APPROVE_TYPE ? fipAddress : selectedAddress,
      cid,
      fipProposalType, // Proposal type (approve or revoke)
    ];

    // Write the contract based on the proposal type
    writeContract({
      abi: fileCoinAbi,
      address: getContractAddress(chainId, 'powerVoting'),
      functionName: 'createFipEditorProposal',
      args: proposalArgs,
    });

    setLoading(false);
  }

  const { isLoading: transactionLoading } =
    useWaitForTransactionReceipt({
      hash,
    })

  return (
    <>
      {contextHolder}
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
                comp: (
                  <RadioGroup className='flex' value={fipProposalType} onChange={handleProposeTypeChange}>
                    <RadioGroup.Option
                      key='approve'
                      value={FIP_APPROVE_TYPE}
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
                      value={FIP_REVOKE_TYPE}
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
                hide: fipProposalType === FIP_REVOKE_TYPE,
                comp: (
                  <input
                    placeholder='Input editor address'
                    className='form-input w-[520px] rounded bg-[#212B3C] border border-[#313D4F]'
                    onChange={(e) => { handleChange('fipAddress', e.target.value) }}
                  />
                )
              },
              {
                name: 'FIP Editor Address',
                hide: fipProposalType === FIP_APPROVE_TYPE,
                comp: (
                  <select
                    onChange={(e: any) => { setSelectedAddress(e.target.value) }}
                    value={selectedAddress}
                    className={classNames(
                      'form-select w-[520px] rounded bg-[#212B3C] border border-[#313D4F]'
                    )}
                  >
                    <option style={{ display: 'none' }}></option>
                    {fipEditors?.map((address: string) => (
                      <option value={address} key={address}>{address}</option>
                    ))}
                  </select>
                )
              },
              {
                name: 'Propose Info',
                width: 100,
                comp: (
                  <textarea
                    value={fipInfo}
                    placeholder='Input propose info'
                    className='form-input h-[320px] w-full rounded bg-[#212B3C] border border-[#313D4F]'
                    onChange={(e) => { handleChange('fipInfo', e.target.value) }}
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
    </>
  )
}

export default FipPropose;
