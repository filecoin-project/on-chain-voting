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

import React, { useState, useEffect, useRef } from "react";
import { Link, useNavigate } from "react-router-dom";
import Table from '../../../components/Table';
import { message } from 'antd';
import { useForm, Controller } from 'react-hook-form';
import classNames from 'classnames';
import { RadioGroup } from '@headlessui/react';
import type { BaseError } from "wagmi";
import { useAccount, useWriteContract, useWaitForTransactionReceipt, useSignMessage } from "wagmi";
import { useConnectModal } from "@rainbow-me/rainbowkit";
import {
  UCAN_JWT_HEADER,
  UCAN_TYPE_FILECOIN,
  UCAN_TYPE_FILECOIN_OPTIONS,
  UCAN_TYPE_GITHUB_OPTIONS,
  UCAN_GITHUB_STEP_1,
  UCAN_GITHUB_STEP_2,
  OPERATION_CANCELED_MSG,
  STORING_DATA_MSG,
  UPLOAD_DATA_FAIL_MSG,
} from '../../../common/consts';
import { stringToBase64Url, validateValue, getWeb3IpfsId, getContractAddress } from '../../../utils';
import './index.less';
import LoadingButton from "../../../components/LoadingButton";
import fileCoinAbi from "../../../common/abi/power-voting.json";
import { useTranslation } from 'react-i18next';
const UcanDelegate = () => {
  const { chain, isConnected, address } = useAccount();
  const { signMessageAsync } = useSignMessage();
  const { openConnectModal } = useConnectModal();
  const navigate = useNavigate();
  const prevAddressRef = useRef(address);
  const [messageApi, contextHolder] = message.useMessage();

  const [ucanType, setUcanType] = useState(UCAN_TYPE_FILECOIN);
  const [githubSignature, setGithubSignature] = useState('');
  const [githubStep, setGithubStep] = useState(UCAN_GITHUB_STEP_1);
  const [formValue] = useState({
    aud: '',
    prf: '',
    owner: '',
    repo: '',
    url: '',
    token: '',
  });

  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors }
  } = useForm({
    defaultValues: {
      ...formValue,
    }
  });

  const {
    data: hash,
    writeContract,
    error,
    isPending: writeContractPending,
    isSuccess: writeContractSuccess,
    reset: resetWriteContract
  } = useWriteContract();
  const [loading, setLoading] = useState<boolean>(writeContractPending);
  const { t } = useTranslation();
  useEffect(() => {
    renderForm();
  }, [githubStep]);

  useEffect(() => {
    const prevAddress = prevAddressRef.current;
    if (prevAddress !== address) {
      navigate("/home");
    }
  }, [address]);

  useEffect(() => {
    if (writeContractSuccess) {

      // save data to localStorage and set validity period to three minutes
      const ucanStorageData = JSON.parse(localStorage.getItem('ucanStorage') || '[]');
      // Calculate expiration time (three minutes from now)
      const expirationTime = Date.now() + 3 * 60 * 1000;
      // Push new data (timestamp and address) to the array
      ucanStorageData.push({ timestamp: expirationTime, address });
      // Save updated data to localStorage
      localStorage.setItem('ucanStorage', JSON.stringify(ucanStorageData));

      messageApi.open({
        type: 'success',
        content: t(STORING_DATA_MSG),
      });
      setTimeout(() => {
        navigate("/home");
      }, 3000);
    }
  }, [writeContractSuccess])

  useEffect(() => {
    if (error) {
      messageApi.open({
        type: 'error',
        content: (error as BaseError)?.shortMessage || error?.message,
      });
    }
    resetWriteContract();
  }, [error]);


  useEffect(() => {
    if (!isConnected) {
      navigate("/home");
      return;
    }
  }, []);

  const handleUcanTypeChange = (value: number) => {
    reset({ aud: '' });
    setUcanType(value);
  }

  /**
   * create ucan
    * @param values
   * @param githubStep
   */
  const onSubmit = (values: any, githubStep?: number) => {
    if (ucanType === UCAN_TYPE_FILECOIN) {
      authorizeFilecoinUcan(values);
    } else {
      switch (githubStep) {
        case UCAN_GITHUB_STEP_1:
          createSignature(values);
          break;
        case UCAN_GITHUB_STEP_2:
          authorizeGithubUcan(values);
          break;
      }
    }
  }

  /**
   * upload UCAN to IPFS
   * @param ucan
   */
  const setUcan = async (ucan: string) => {
    // Get the IPFS ID for the provided UCAN
    const cid = await getWeb3IpfsId(ucan);
    if (!cid?.length) {
      setLoading(false);
      messageApi.open({
        type: 'warning',
        content: t(UPLOAD_DATA_FAIL_MSG),
      });
      return;
    }

    writeContract({
      abi: fileCoinAbi,
      address: getContractAddress(chain?.id || 0, 'powerVoting'),
      functionName: 'ucanDelegate',
      args: [
        cid
      ],
    });
    setLoading(false);
  }

  const authorizeFilecoinUcan = async (values: any) => {
    setLoading(true);
    const { aud, prf } = values;
    if (!aud || !prf) {
      return;
    }
    const ucanParams = {
      iss: address,
      aud,
      prf,
      act: 'add',
    }
    const base64Header = stringToBase64Url(JSON.stringify(UCAN_JWT_HEADER));
    const base64Params = stringToBase64Url(JSON.stringify(ucanParams));
    let signature = '';
    try {
      signature = await signMessageAsync({ message: `${base64Header}.${base64Params}` })
    } catch (e) {
      messageApi.open({
        type: 'error',
        content: t(OPERATION_CANCELED_MSG),
      });
      setLoading(false);
      return;
    }
    const base64Signature = stringToBase64Url(signature);
    const ucan = `${base64Header}.${base64Params}.${base64Signature}`;
    setUcan(ucan);
  }

  const createSignature = async (values: any) => {
    setLoading(true);
    if (isConnected) {
      try {
        const { aud } = values;
        // If 'aud' is not provided, return without proceeding
        if (!aud) {
          return;
        }
        // Define signature parameters
        const signatureParams = {
          iss: address,
          aud,
          prf: '',
          act: 'add',
        }
        // Create a new Web3Provider using the current Ethereum provider

        // Convert header and params to base64 URL
        const base64Header = stringToBase64Url(JSON.stringify(UCAN_JWT_HEADER));
        const base64Params = stringToBase64Url(JSON.stringify(signatureParams));
        let signature = '';
        try {
          // Sign the concatenated header and params
          signature = await signMessageAsync({ message: `${base64Header}.${base64Params}` })
        } catch (e) {
          messageApi.open({
            type: 'error',
            content: t(OPERATION_CANCELED_MSG),
          });
          setLoading(false);
          return;
        }
        // Convert signature to base64 URL
        const base64Signature = stringToBase64Url(signature);

        // Concatenate header, params, and signature
        const githubSignatureParams = `${base64Header}.${base64Params}.${base64Signature}`;

        // Set GitHub signature and step
        setGithubSignature(githubSignatureParams);
        setGithubStep(UCAN_GITHUB_STEP_2);
      } catch (e) {
        console.log(e);
      }
    } else {
      openConnectModal && openConnectModal();
    }
    setLoading(false);
  }
  const authorizeGithubUcan = async (values: any) => {
    setLoading(true);
    const { url } = values;
    setUcan(url);
  }

  const filecoinAuthorizeList = [
    {
      name: t('content.ucatType'),
      width: 100,
      hide: false,
      comp: (
        <RadioGroup className='flex h-[30px] mt-[-5px]' value={ucanType} onChange={handleUcanTypeChange}>
          {[...(UCAN_TYPE_FILECOIN_OPTIONS.map((item) => {
            return {
              label: t(item.label), value: item.value
            }
          })), ...UCAN_TYPE_GITHUB_OPTIONS].map(item => (
            <RadioGroup.Option
              key={item.label}
              value={item.value}
              className='relative flex items-center cursor-pointer p-4 focus:outline-none'
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
                        checked ? 'text-[#4B535B]' : 'text-[#8896AA]'
                      }
                    >
                      {item.label}
                    </RadioGroup.Label>
                  </span>
                </>
              )}
            </RadioGroup.Option>
          ))}
        </RadioGroup>
      )
    },
    {
      name: t('content.issuer'),
      width: 100,
      comp: (
        <input
          disabled
          value={address}
          placeholder={t('content.filecoinAddress')}
          className='form-input w-[520px] rounded bg-[#ffffff] border border-[#eeeeee] text-[#4B535B] cursor-not-allowed'
        />
      )
    },
    {
      name: t('content.audience'),
      width: 100,
      comp: (
        <>
          <Controller
            name="aud"
            control={control}
            render={() => <input
              className={classNames(
                'form-input w-[520px] rounded bg-[#ffffff] border border-[#eeeeee] text-[#4B535B]',
                errors.aud && 'border-red-500 focus:border-red-500'
              )}
              {...register('aud', { required: true, validate: validateValue })}
            />}
          />
          {errors.aud && (
            <p className='text-red-500 mt-1'>{t('content.audRequired')}</p>
          )}
        </>
      )
    },
    {
      name: t('content.proof'),
      width: 100,
      comp: (
        <>
          <Controller
            name="prf"
            control={control}
            render={() => <textarea
              placeholder='The full UCAN content (include header, payload and signature) signed by your Filecoin private key.'
              className={classNames(
                'form-input h-[320px] w-full rounded bg-[#ffffff] border border-[#eeeeee] text-[#4B535B]',
                errors.prf && 'border-red-500 focus:border-red-500'
              )}
              {...register('prf', { required: true, validate: validateValue })}
            />}
          />
          {errors.prf && (
            <p className='text-red-500 mt-1'>{t('content.proofRequired')}</p>
          )}
        </>
      )
    },
  ];

  const githubSignatureList = [
    {
      name: t('content.ucatType'),
      width: 100,
      hide: false,
      comp: (
        <RadioGroup className='flex' value={ucanType} onChange={handleUcanTypeChange}>
          {[...(UCAN_TYPE_FILECOIN_OPTIONS.map((item) => {
            return {
              label: t(item.label), value: item.value
            }
          })), ...UCAN_TYPE_GITHUB_OPTIONS].map(item => (
            <RadioGroup.Option
              key={item.label}
              value={item.value}
              className='relative flex items-center cursor-pointer p-4 focus:outline-none'
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
                        checked ? 'text-[#4B535B]' : 'text-[#8896AA]'
                      }
                    >
                      {item.label}
                    </RadioGroup.Label>
                  </span>
                </>
              )}
            </RadioGroup.Option>
          ))}
        </RadioGroup>
      )
    },
    {
      name: t('content.issuer'),
      width: 100,
      comp: (
        <input
          disabled
          value={address}
          className='form-input w-[520px] rounded bg-[#FFFFFF] border border-[#EEEEEE] text-[#4B535B] cursor-not-allowed '
        />
      )
    },
    {
      name: t('content.audience'),
      width: 100,
      comp: (
        <>
          <Controller
            name="aud"
            control={control}
            render={() => <input
              placeholder={t('content.yourGithubAccount')}
              className={classNames(
                'form-input w-[520px] rounded bg-[#FFFFFF] border border-[#EEEEEE] text-[#4B535B]',
                errors.aud && 'border-red-500 focus:border-red-500'
              )}
              {...register('aud', { required: true, validate: validateValue })}
            />}
          />
          {errors.aud && (
            <p className='text-red-500 mt-1'>{t('content.audRequired')}</p>
          )}
        </>
      )
    },
  ];

  const githubAuthorizeList = [
    {
      name: t('content.signature'),
      width: 100,
      comp: (
        <textarea
          disabled
          value={githubSignature}
          className='form-input h-[320px] w-full rounded bg-[#ffffff] border border-[#eeeeee] text-black cursor-not-allowed'
        />
      )
    },
    {
      name: 'URL',
      width: 100,
      comp: (
        <>
          <Controller
            name="url"
            control={control}
            render={() => <input
              className={classNames(
                'form-input w-full rounded bg-[#ffffff] border border-[#eeeeee] text-black',
                errors.url && 'border-red-500 focus:border-red-500'
              )}
              {...register('url', { required: true, validate: validateValue })}
            />}
          />
          {errors.url && (
            <p className='text-red-500 mt-1'>{t('content.uRLRequired')}</p>
          )}
        </>
      )
    },
  ];

  const renderFilecoinAuthorize = () => {
    return (
      <form onSubmit={handleSubmit(value => { onSubmit(value) })}>
        <div className='flow-root space-y-8'>
          <Table
            title={t('content.ucanDelegatesAuthorize')}
            link={{
              type: 'filecoin',
              action: 'authorize',
              href: '/ucanDelegate/help'
            }}
            list={filecoinAuthorizeList}
          />

          <div className='text-center'>
            <LoadingButton text={t('content.authorize')} loading={loading || writeContractPending || transactionLoading} />
          </div>
        </div>
      </form>
    )
  }

  const renderGithubSignature = () => {
    return (
      <form onSubmit={handleSubmit(value => { onSubmit(value, UCAN_GITHUB_STEP_1) })}>
        <div className='flow-root space-y-8'>
          <Table
            title={t('content.ucanDelegatesAuthorize')}
            link={{
              type: 'github',
              action: 'authorize',
              href: '/ucanDelegate/help'
            }}
            list={githubSignatureList}
          />

          <div className='text-center'>
            <LoadingButton text={t('content.sign')} loading={loading || writeContractPending || transactionLoading} />
          </div>
        </div>
      </form>
    )
  }

  const renderGithubAuthorize = () => {
    return (
      <form onSubmit={handleSubmit(value => { onSubmit(value, UCAN_GITHUB_STEP_2) })}>
        <div className='flow-root space-y-8'>
          <Table
           title={t('content.ucanDelegatesAuthorize')}
            link={{
              type: 'github',
              action: 'authorize',
              href: '/ucanDelegate/help'
            }}
            list={githubAuthorizeList}
          />

          <div className='text-center'>
            <button
              className={`h-[40px] bg-sky-500 hover:bg-sky-700 text-white py-2 px-6 rounded-xl disabled:opacity-50 mr-8 ${loading && 'cursor-not-allowed'}`}
              type='button' onClick={() => { setGithubStep(UCAN_GITHUB_STEP_1) }}>
              {t('content.previous')}
            </button>
            <LoadingButton text={t('content.authorize')} loading={loading || writeContractPending || transactionLoading} />
          </div>
        </div>
      </form>
    )
  }

  const renderForm = () => {
    if (ucanType === UCAN_TYPE_FILECOIN) {
      return renderFilecoinAuthorize();
    } else {
      switch (githubStep) {
        case UCAN_GITHUB_STEP_1:
          return renderGithubSignature();
        case UCAN_GITHUB_STEP_2:
          return renderGithubAuthorize();

      }
    }
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
          <div className="inline-flex items-center gap-1 mb-8  text-skin-text hover:text-skin-link">
            <Link to="/home" className="flex items-center">
              <svg className="mr-1" viewBox="0 0 24 24" width="1.2em" height="1.2em">
                <path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                  d="m11 17l-5-5m0 0l5-5m-5 5h12"></path>
              </svg>
              {t('content.back')}
            </Link>
          </div>
        </button>
        {
          renderForm()
        }
      </div>
    </>
  )
}

export default UcanDelegate;
