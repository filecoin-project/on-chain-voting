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

import { RadioGroup } from '@headlessui/react';
import { message } from "antd";
import classNames from 'classnames';
import { useEffect, useRef, useState } from "react";
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from "react-router-dom";
import { useFipList } from '../../../common/store';
import { isAddress } from 'viem';
import type { BaseError } from "wagmi";
import { useAccount, useWaitForTransactionReceipt, useWriteContract } from "wagmi";
import votingFipeditorAbi from "../../../common/abi/power-voting-fipeditor.json";
import {
  calibrationChainId,
  FIP_ALREADY_EXECUTE_MSG,
  FIP_APPROVE_ALREADY_MSG,
  FIP_APPROVE_SELF_MSG,
  FIP_EDITOR_APPROVE_TYPE,
  FIP_EDITOR_REVOKE_TYPE,
  getFipProposalApi,
  NO_ENOUGH_FIP_EDITOR_REVOKE_ADDRESS_MSG,
  NO_FIP_EDITOR_APPROVE_ADDRESS_MSG,
  NO_FIP_EDITOR_REVOKE_ADDRESS_MSG,
  NO_FIP_INfO_MSG,
  STORING_DATA_MSG
} from "../../../common/consts";
import LoadingButton from '../../../components/LoadingButton';
import Table from '../../../components/Table';
import { getContractAddress, hexToString } from "../../../utils";
import axios from "axios";

const FipEditorPropose = () => {
  const { isConnected, address, chain } = useAccount();
  const { t } = useTranslation();
  const chainId = chain?.id || calibrationChainId;

  const navigate = useNavigate();
  const prevAddressRef = useRef(address);
  const [messageApi, contextHolder] = message.useMessage();


  const [fipAddress, setFipAddress] = useState('');
  const [selectedAddress, setSelectedAddress] = useState('');
  const [fipInfo, setFipInfo] = useState('');
  const [fipProposalType, setFipEditorProposeType] = useState(FIP_EDITOR_APPROVE_TYPE);

  const { fipList, isFipEditorAddress } = useFipList((state: any) => state.data);

  const {
    data: hash,
    writeContract,
    error,
    isPending: writeContractPending,
    isSuccess: writeContractSuccess,
    reset,
  } = useWriteContract();

  const [loading, setLoading] = useState(writeContractPending);
  const [approveResult, setApproveResult] = useState([])
  const [revokeResult, setRevokeResult] = useState([])
  const getApproveList = async () => {
    const params = {
      chainId,
      page: 1,
      pageSize: 20,
      proposalType: 1
    }
    const { data: { data: { list } } } = await axios.get(getFipProposalApi, { params });
    setApproveResult(list)
  }
  const getRevokeList = async () => {
    const params = {
      chainId,
      page: 1,
      pageSize: 20,
      proposalType: 0
    }
    const { data: { data: { list } } } = await axios.get(getFipProposalApi, { params });

    setRevokeResult(list)
  }
  useEffect(() => {
    if (!chainId || !address) return
    getApproveList()
    getRevokeList()
  }, [address, chainId])


  useEffect(() => {
    if (!isConnected || !isFipEditorAddress) {
      navigate("/home");
      return;
    }
  }, [isConnected, isFipEditorAddress]);

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
        content: t(STORING_DATA_MSG),
      });
      setTimeout(() => {
        navigate("/home")
      }, 1000);
    }
  }, [writeContractSuccess]);

  useEffect(() => {
    if (error) {
      // Get error cause
      const errorStr = JSON.stringify(error);
      // Intercepts the first hexadecimal in the string
      const reg = /revert reason:\s*0x[0-9A-Fa-f]+/;
      const match = errorStr.match(reg) || [];
      messageApi.open({
        type: 'error',
        content: hexToString(match[0]) || (error as BaseError)?.shortMessage,
      });
    }
    reset();
  }, [error]);

  const handleChange = (type: string, value: string) => {
    type === 'fipAddress' ? setFipAddress(value) : setFipInfo(value);
  }

  const handleProposeTypeChange = (value: number) => {
    setFipEditorProposeType(value);
    value === FIP_EDITOR_APPROVE_TYPE ? setSelectedAddress('') : setFipAddress('');
    handleChange('fipInfo', '');
  }

  /**
   * Set miner ID
   */
  const onSubmit = async () => {
    // Check if required fields are filled based on proposal type
    if (fipProposalType === FIP_EDITOR_APPROVE_TYPE && !fipAddress) {
      messageApi.open({
        type: 'warning',
        // Prompt user to fill required fields
        content: t(NO_FIP_EDITOR_APPROVE_ADDRESS_MSG),
      });
      return;
    }
    if (fipProposalType === FIP_EDITOR_APPROVE_TYPE && !isAddress(fipAddress)) {
      messageApi.open({
        type: 'warning',
        // Prompt user to fill required fields
        content: t('content.inputValidAddress'),
      });
      return
    }
    if (fipProposalType === FIP_EDITOR_REVOKE_TYPE && !selectedAddress) {
      messageApi.open({
        type: 'warning',
        // Prompt user to fill required fields
        content: t(NO_FIP_EDITOR_REVOKE_ADDRESS_MSG),
      });
      return;
    }
    if (!fipInfo || !fipInfo.trim()) {
      messageApi.open({
        type: 'warning',
        // Prompt user to fill required fields
        content: t(NO_FIP_INfO_MSG),
      });
      return;
    }
    if (fipProposalType === FIP_EDITOR_REVOKE_TYPE && fipList.length <= 2) {
      messageApi.open({
        type: 'warning',
        // must more than 2
        content: t(NO_ENOUGH_FIP_EDITOR_REVOKE_ADDRESS_MSG),
      });
      return;
    }



    if (fipProposalType === FIP_EDITOR_REVOKE_TYPE && revokeResult && revokeResult.length > 0) {
      const find = revokeResult?.find((v: any) => v.candidateAddress === selectedAddress)
      if (find) {
        messageApi.open({
          type: 'warning',
          content: t(FIP_ALREADY_EXECUTE_MSG),
        });
        return;
      }
    }
    if (fipProposalType === FIP_EDITOR_APPROVE_TYPE && approveResult && approveResult.length > 0) {
      const find = approveResult?.find((v: any) => v.candidateAddress === fipAddress)
      if (find) {
        messageApi.open({
          type: 'warning',
          content: t(FIP_ALREADY_EXECUTE_MSG),
        });
        return;
      }
    }

    if (fipProposalType === FIP_EDITOR_APPROVE_TYPE && fipAddress === address) {
      messageApi.open({
        type: 'warning',
        content: t(FIP_APPROVE_SELF_MSG),
      });
      return;
    }

    //fipEditors
    if (fipProposalType === FIP_EDITOR_APPROVE_TYPE && !!fipList.filter((item: any) => item.editor === fipAddress)[0]) {
      messageApi.open({
        type: 'warning',
        content: t(FIP_APPROVE_ALREADY_MSG),
      });
      return;
    }

    // Set loading state to true while submitting proposal
    setLoading(true);

    // Construct the arguments and call the writeContract function to create the proposal
    const proposalArgs = [
      // Use appropriate address based on proposal type
      fipProposalType === FIP_EDITOR_APPROVE_TYPE ? fipAddress : selectedAddress,
      fipInfo,
      fipProposalType, // Proposal type (approve or revoke)
    ];

    // Write the contract based on the proposal type
    //TODO
    writeContract({
      abi: votingFipeditorAbi,
      address: getContractAddress(chainId, 'powerVotingFip'),
      functionName: 'createFipEditorProposal',
      args: proposalArgs,
    });
    setLoading(false);
  }

  const { isLoading: transactionLoading } =
    useWaitForTransactionReceipt({
      hash,
    })
  const isLoading = loading || writeContractPending || transactionLoading;

  return (
    <>
      {contextHolder}
      <div className="px-3 mb-6 md:px-0">
        <button>
          <div className="inline-flex items-center gap-1 mb-8 text-skin-text hover:text-skin-link">
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
            title={t('content.fipEditorPropose')}
            list={[
              {
                name: t('content.proposeType'),
                width: 100,
                comp: (
                  <RadioGroup className='flex h-[30px] mt-[-5px]' value={fipProposalType} onChange={handleProposeTypeChange}>
                    <RadioGroup.Option
                      key='approve'
                      disabled={isLoading}
                      value={FIP_EDITOR_APPROVE_TYPE}
                      className='relative flex items-center cursor-pointer p-4 focus:outline-none data-[disabled]:cursor-not-allowed'
                    >
                      {({ active, checked }) => (
                        <>
                          <span
                            className={classNames(
                              checked
                                ? 'bg-[#45B753] border-transparent'
                                : 'bg-[#eeeeee] border-transparent]',
                              active ? 'ring-2 ring-offset-2 ring-[#ffffff]' : '',
                              'mt-0.5 h-4 w-4 shrink-0 cursor-pointer rounded-full border flex items-center justify-center'
                            )}
                            aria-hidden='true'
                          >
                            {(active || checked) && (
                              <span className='rounded-full bg-white w-1.5 h-1.5' />
                            )}
                          </span>
                          <span className='ml-3'>
                            <RadioGroup.Label
                              as='span'
                              className={
                                checked ? 'text-black' : 'text-[#8896AA]'
                              }
                            >
                              {t('content.approve')}
                            </RadioGroup.Label>
                          </span>
                        </>
                      )}
                    </RadioGroup.Option>
                    <RadioGroup.Option
                      key='revoke'
                      disabled={isLoading}
                      value={FIP_EDITOR_REVOKE_TYPE}
                      className='relative flex items-center cursor-pointer p-4 focus:outline-none data-[disabled]:cursor-not-allowed'
                    >
                      {({ active, checked }) => (
                        <>
                          <span
                            className={classNames(
                              checked
                                ? 'bg-[#45B753] border-transparent'
                                : 'bg-[#eeeeee] border-transparent]',
                              active ? 'ring-2 ring-offset-2 ring-[#ffffff]' : '',
                              'mt-0.5 h-4 w-4 shrink-0 cursor-pointer rounded-full border flex items-center justify-center'
                            )}
                            aria-hidden='true'
                          >
                            {(active || checked) && (
                              <span className='rounded-full bg-white w-1.5 h-1.5' />
                            )}
                          </span>
                          <span className='ml-3'>
                            <RadioGroup.Label
                              as='span'
                              className={
                                checked ? 'text-black' : 'text-[#8896AA]'
                              }
                            >
                              {t('content.revoke')}
                            </RadioGroup.Label>
                          </span>
                        </>
                      )}
                    </RadioGroup.Option>
                  </RadioGroup>
                )
              },
              {
                name: t('content.editorAddress'),
                hide: fipProposalType === FIP_EDITOR_REVOKE_TYPE,
                comp: (
                  <input
                    placeholder={t('content.inputEditorAddress')}
                    className='form-input w-[520px] rounded bg-[#ffffff] border border-[#eeeeee] text-black'
                    onChange={(e) => { handleChange('fipAddress', e.target.value) }}
                  />
                )
              },
              {
                name: t('content.fipEditorAddress'),
                hide: fipProposalType === FIP_EDITOR_APPROVE_TYPE,
                comp: (
                  <select
                    onChange={(e: any) => { setSelectedAddress(e.target.value) }}
                    value={selectedAddress}
                    className={classNames(
                      'form-select w-[520px] rounded bg-[#ffffff] border border-[#eeeeee] text-black'
                    )}
                  >
                    <option style={{ display: 'none' }}></option>
                    {fipList?.map((item: any) => (
                      <option
                        disabled={address === item.editor}
                        value={item.editor}
                        key={item.editor}
                      >
                        {item.editor}
                      </option>
                    ))}
                  </select>
                )
              },
              {
                name: t('content.proposeInfo'),
                width: 100,
                comp: (
                  <textarea
                    value={fipInfo}
                    maxLength={300}
                    placeholder={t('content.inputProposeInfo')}
                    className='form-input h-[320px] w-full rounded bg-[#ffffff] border border-[#eeeeee] text-black'
                    onChange={(e) => { handleChange('fipInfo', e.target.value) }}
                  />
                )
              }
            ]}
          />

          <div className='text-center'>
            <LoadingButton text={t('content.submit')} loading={isLoading} handleClick={onSubmit} />
          </div>
        </div>
      </div>
    </>
  )
}

export default FipEditorPropose;
