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
import { ethers } from "ethers";
import {Link, useNavigate} from "react-router-dom";
import Table from '../../../components/Table';
import { message } from 'antd';
import {useForm, Controller} from 'react-hook-form';
import classNames from 'classnames';
import {RadioGroup} from '@headlessui/react';
import {useNetwork, useAccount} from "wagmi";
import {useConnectModal} from "@rainbow-me/rainbowkit";
import {
  UCAN_JWT_HEADER,
  UCAN_TYPE_FILECOIN,
  UCAN_TYPE_FILECOIN_OPTIONS,
  UCAN_TYPE_GITHUB_OPTIONS,
  UCAN_GITHUB_STEP_1,
  UCAN_GITHUB_STEP_2,
  OPERATION_CANCELED_MSG,
} from '../../../common/consts';
import { stringToBase64Url, validateValue } from '../../../utils';
import {getWeb3IpfsId, useDynamicContract} from "../../../hooks";
import './index.less';
import LoadingButton from "../../../components/LoadingButton";

const UcanDelegate = () => {
  const {chain} = useNetwork();
  const {isConnected, address} = useAccount();
  const {openConnectModal} = useConnectModal();
  const navigate = useNavigate();
  const prevAddressRef = useRef(address);

  const [ucanType, setUcanType] = useState(UCAN_TYPE_FILECOIN);
  const [loading, setLoading] = useState(false);
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
    formState: {errors}
  } = useForm({
    defaultValues: {
      ...formValue,
    }
  });

  useEffect(() => {
    renderForm();
  }, [githubStep]);

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
    const chainId = chain?.id || 0;
    const { ucanDelegate } = useDynamicContract(chainId);
    // Get the IPFS ID for the provided UCAN
    const cid = await getWeb3IpfsId(ucan);
    // Call the ucanDelegate function with the IPFS ID
    const res = await ucanDelegate(cid);
    if (res.code === 200 && res.data?.hash) {
      message.success(res.msg);
      navigate("/");

      // save data to localStorage and set validity period to three minutes
      const ucanStorageData = JSON.parse(localStorage.getItem('ucanStorage') || '[]');
      // Calculate expiration time (three minutes from now)
      const expirationTime = Date.now() + 3 * 60 * 1000;
      // Push new data (timestamp and address) to the array
      ucanStorageData.push({ timestamp: expirationTime, address });
      // Save updated data to localStorage
      localStorage.setItem('ucanStorage', JSON.stringify(ucanStorageData));
    } else {
      message.error(res.msg, 3);
    }
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
    // @ts-ignore
    const provider = new ethers.providers.Web3Provider(window.ethereum);
    const signer = await provider.getSigner();
    const base64Header = stringToBase64Url(JSON.stringify(UCAN_JWT_HEADER));
    const base64Params = stringToBase64Url(JSON.stringify(ucanParams));
    let signature = '';
    try {
      signature = await signer.signMessage(`${base64Header}.${base64Params}`);
    } catch (e) {
      message.error(OPERATION_CANCELED_MSG);
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
          prf:'',
          act: 'add',
        }
        // Create a new Web3Provider using the current Ethereum provider
        // @ts-ignore
        const provider = new ethers.providers.Web3Provider(window.ethereum);
        const signer = await provider.getSigner();

        // Convert header and params to base64 URL
        const base64Header = stringToBase64Url(JSON.stringify(UCAN_JWT_HEADER));
        const base64Params = stringToBase64Url(JSON.stringify(signatureParams));
        let signature = '';
        try {
          // Sign the concatenated header and params
          signature = await signer.signMessage(`${base64Header}.${base64Params}`);
        } catch (e) {
          message.error(OPERATION_CANCELED_MSG);
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
      // @ts-ignore
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
      name: 'UCAN Type',
      width: 100,
      hide: false,
      comp: (
        <RadioGroup className='flex' value={ucanType} onChange={handleUcanTypeChange}>
          {[...UCAN_TYPE_FILECOIN_OPTIONS, ...UCAN_TYPE_GITHUB_OPTIONS].map(item => (
            <RadioGroup.Option
              key={item.label}
              value={item.value}
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
      name: 'Issuer',
      width: 100,
      comp: (
        <input
          disabled
          value={address}
          placeholder='Your Filecoin address.'
          className='form-input w-[520px] rounded bg-[#212B3C] border border-[#313D4F] cursor-not-allowed'
        />
      )
    },
    {
      name: 'Audience',
      width: 100,
      comp: (
        <>
          <Controller
            name="aud"
            control={control}
            render={() => <input
              className={classNames(
                'form-input w-[520px] rounded bg-[#212B3C] border border-[#313D4F]',
                errors['aud'] && 'border-red-500 focus:border-red-500'
              )}
              {...register('aud', {required: true, validate: validateValue})}
            />}
          />
          {errors['aud'] && (
            <p className='text-red-500 mt-1'>Aud is required</p>
          )}
        </>
      )
    },
    {
      name: 'Proof',
      width: 100,
      comp: (
        <>
          <Controller
            name="prf"
            control={control}
            render={() => <textarea
              placeholder='The full UCAN content (include header, payload and signature) signed by your Filecoin private key.'
              className={classNames(
                'form-input h-[320px] w-full rounded bg-[#212B3C] border border-[#313D4F]',
                errors['prf'] && 'border-red-500 focus:border-red-500'
              )}
              {...register('prf', {required: true, validate: validateValue})}
            />}
          />
          {errors['prf'] && (
            <p className='text-red-500 mt-1'>Proof is required</p>
          )}
        </>
      )
    },
  ];

  const githubSignatureList = [
    {
      name: 'UCAN Type',
      width: 100,
      hide: false,
      comp: (
        <RadioGroup className='flex' value={ucanType} onChange={handleUcanTypeChange}>
          {[...UCAN_TYPE_FILECOIN_OPTIONS, ...UCAN_TYPE_GITHUB_OPTIONS].map(item => (
            <RadioGroup.Option
              key={item.label}
              value={item.value}
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
      name: 'Issuer',
      width: 100,
      comp: (
        <input
          disabled
          value={address}
          className='form-input w-[520px] rounded bg-[#212B3C] border border-[#313D4F] cursor-not-allowed'
        />
      )
    },
    {
      name: 'Audience',
      width: 100,
      comp: (
        <>
          <Controller
            name="aud"
            control={control}
            render={() => <input
              placeholder='Your github account.'
              className={classNames(
                'form-input w-[520px] rounded bg-[#212B3C] border border-[#313D4F]',
                errors['aud'] && 'border-red-500 focus:border-red-500'
              )}
              {...register('aud', {required: true, validate: validateValue})}
            />}
          />
          {errors['aud'] && (
            <p className='text-red-500 mt-1'>Aud is required</p>
          )}
        </>
      )
    },
  ];

  const githubAuthorizeList = [
    {
      name: 'Signature',
      width: 100,
      comp: (
        <textarea
          disabled
          value={githubSignature}
          className='form-input h-[320px] w-full rounded bg-[#212B3C] border border-[#313D4F] cursor-not-allowed'
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
                'form-input w-full rounded bg-[#212B3C] border border-[#313D4F]',
                errors['url'] && 'border-red-500 focus:border-red-500'
              )}
              {...register('url', {required: true, validate: validateValue})}
            />}
          />
          {errors['url'] && (
            <p className='text-red-500 mt-1'>URL is required</p>
          )}
        </>
      )
    },
  ];

  const renderFilecoinAuthorize = () => {
    return (
      <form onSubmit={handleSubmit((value) => { onSubmit(value) })}>
        <div className='flow-root space-y-8'>
          <Table
            title='UCAN Delegates (Authorize)'
            link={{
              type: 'filecoin',
              action: 'authorize',
              href: '/ucanDelegate/help'
            }}
            list={filecoinAuthorizeList}
          />

          <div className='text-center'>
            <LoadingButton text='Authorize' loading={loading} />
          </div>
        </div>
      </form>
    )
  }

  const renderGithubSignature = () => {
    return (
      <form onSubmit={handleSubmit((value) => { onSubmit(value, UCAN_GITHUB_STEP_1) })}>
        <div className='flow-root space-y-8'>
          <Table
            title='UCAN Delegates (Authorize)'
            link={{
              type: 'github',
              action: 'authorize',
              href: '/ucanDelegate/help'
            }}
            list={githubSignatureList}
          />

          <div className='text-center'>
            <LoadingButton text='Sign' loading={loading} />
          </div>
        </div>
      </form>
    )
  }

  const renderGithubAuthorize = () => {
    return (
      <form onSubmit={handleSubmit((value) => { onSubmit(value, UCAN_GITHUB_STEP_2) })}>
        <div className='flow-root space-y-8'>
          <Table
            title='UCAN Delegates (Authorize)'
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
              Previous
            </button>
            <LoadingButton text='Authorize' loading={loading} />
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

  return (
    <>
      <div className="px-3 mb-6 md:px-0">
        <button>
          <div className="inline-flex items-center gap-1 text-skin-text hover:text-skin-link">
            <Link to="/" className="flex items-center">
              <svg className="mr-1" viewBox="0 0 24 24" width="1.2em" height="1.2em">
                <path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                      d="m11 17l-5-5m0 0l5-5m-5 5h12"></path>
              </svg>
              Back
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
