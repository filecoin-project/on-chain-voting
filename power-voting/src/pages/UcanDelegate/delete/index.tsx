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
import {useLocation, useNavigate, Link} from "react-router-dom";
import { ethers } from "ethers";
import { message } from 'antd';
import Table from '../../../components/Table';
import {useForm, Controller} from 'react-hook-form';
import classNames from 'classnames';
import {RadioGroup} from '@headlessui/react';
import {useNetwork, useAccount} from "wagmi";
import {useConnectModal} from "@rainbow-me/rainbowkit";
import {
  UCAN_GITHUB_STEP_1,
  UCAN_GITHUB_STEP_2,
  UCAN_TYPE_GITHUB_OPTIONS,
  UCAN_JWT_HEADER,
  UCAN_TYPE_FILECOIN_OPTIONS,
  STORING_DATA_MSG, OPERATION_CANCELED_MSG,
} from '../../../common/consts';
import './index.less';
import {stringToBase64Url} from "../../../utils";
import {getIpfsId, useDynamicContract} from "../../../hooks";
import LoadingButton from "../../../components/LoadingButton";

const UcanDelegate = () => {
  const {chain} = useNetwork();
  const {isConnected, address} = useAccount();
  const {openConnectModal} = useConnectModal();
  const navigate = useNavigate();
  const prevAddressRef = useRef(address);

  const location = useLocation();
  const params = location.state?.params;

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
    formState: {errors}
  } = useForm({
    defaultValues: {
      ...formValue,
    }
  });

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

  const validateValue = (value: string) => {
    return value?.trim() !== '';
  };

  const onSubmit = (values: any, githubStep?: number) => {
    if (params?.isGithubType) {
      switch (githubStep) {
        case UCAN_GITHUB_STEP_1:
          createSignature();
          break;
        case UCAN_GITHUB_STEP_2:
          deAuthorizeGithubUcan(values);
          break;
      }
    } else {
      deAuthorizeFilecoinUcan(values);
    }
  }

  const setUcan = async (ucan: string) => {
    const chainId = chain?.id || 0;
    const { ucanDelegate } = useDynamicContract(chainId);
    const cid = await getIpfsId(ucan) as any;
    const res = await ucanDelegate(cid);
    if (res.code === 200 && res.data?.hash) {
      message.success(STORING_DATA_MSG);
      navigate("/");
    } else {
      message.error(res.msg, 3);
    }
    setLoading(false);
  }

  const deAuthorizeFilecoinUcan = async (values:  any) => {
    setLoading(true);
    const { aud } = params;
    const { prf } = values;
    if (!aud || !prf) {
      return;
    }
    const ucanParams = {
      iss: address,
      aud,
      prf,
      act: 'del',
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

  const createSignature = async () => {
    setLoading(true);
    if (isConnected) {
      try {
        const { aud } = params;
        if (!aud) {
          return;
        }
        const signatureParams = {
          iss: address,
          aud,
          prf:'',
          act: 'del',
        }
        // @ts-ignore
        const provider = new ethers.providers.Web3Provider(window.ethereum);
        const signer = await provider.getSigner();
        const base64Header = stringToBase64Url(JSON.stringify(UCAN_JWT_HEADER));
        const base64Params = stringToBase64Url(JSON.stringify(signatureParams));
        let signature = '';
        try {
          signature = await signer.signMessage(`${base64Header}.${base64Params}`);
        } catch (e) {
          message.error(OPERATION_CANCELED_MSG);
          setLoading(false);
          return;
        }
        const base64Signature = stringToBase64Url(signature);
        const githubSignatureParams = `${base64Header}.${base64Params}.${base64Signature}`;
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

  const deAuthorizeGithubUcan = async (values: any) => {
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
        <RadioGroup className='flex'>
          {UCAN_TYPE_FILECOIN_OPTIONS.map(item => (
            <RadioGroup.Option
              key={item.label}
              value={item.value}
              className='relative flex items-center cursor-pointer p-4 focus:outline-none'
            >
              {() => (
                <>
                        <span
                          className='bg-[#45B753] border-transparent mt-0.5 h-4 w-4 shrink-0 cursor-pointer rounded-full border flex items-center justify-center'
                          aria-hidden='true'
                        >
                           <span className='rounded-full bg-white w-1.5 h-1.5'/>
                        </span>
                  <span className='ml-3'>
                          <RadioGroup.Label
                            as='span'
                            className='text-white'
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
        <input
          disabled
          className='form-input w-[520px] rounded bg-[#212B3C] border border-[#313D4F] cursor-not-allowed'
          value={params?.aud || ''}
        />
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
        <RadioGroup>
          {UCAN_TYPE_GITHUB_OPTIONS.map(item => (
            <RadioGroup.Option
              key={item.label}
              value={item.value}
              className='relative flex items-center cursor-pointer p-4 focus:outline-none'
            >
              {() => (
                <>
                        <span
                          className='bg-[#45B753] border-transparent mt-0.5 h-4 w-4 shrink-0 cursor-pointer rounded-full border flex items-center justify-center'
                          aria-hidden='true'
                        >
                           <span className='rounded-full bg-white w-1.5 h-1.5'/>
                        </span>
                  <span className='ml-3'>
                          <RadioGroup.Label
                            as='span'
                            className='text-white'
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
        <input
          disabled
          className='form-input w-[520px] rounded bg-[#212B3C] border border-[#313D4F] cursor-not-allowed'
          value={params?.aud || ''}
        />
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

  const renderFilecoinDeauthorize = () => {
    return (
      <form onSubmit={handleSubmit((value) => { onSubmit(value) })}>
        <div className='flow-root space-y-8'>
          <Table
            title='UCAN Delegates (Deauthorize)'
            link={{
              type: 'filecoin',
              action: 'deAuthorize',
              href: '/ucanDelegate/help'
            }}
            list={filecoinAuthorizeList}
          />

          <div className='text-center'>
            <LoadingButton className='!bg-red-500 !hover:bg-red-700' text='Deauthorize' loading={loading} />
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
            title='UCAN Delegates (Deauthorize)'
            link={{
              type: 'github',
              action: 'deAuthorize',
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

  const renderGithubDeauthorize = () => {
    return (
      <form onSubmit={handleSubmit((value) => { onSubmit(value, UCAN_GITHUB_STEP_2) })}>
        <div className='flow-root space-y-8'>
          <Table title='UCAN Delegates (Deauthorize)' list={githubAuthorizeList}/>

          <div className='text-center'>
            <button
              className={`h-[40px] bg-sky-500 hover:bg-sky-700 text-white py-2 px-6 rounded-xl disabled:opacity-50 mr-8 ${loading && 'cursor-not-allowed'}`}
              type='button' onClick={() => { setGithubStep(UCAN_GITHUB_STEP_1) }}>
              Previous
            </button>
            <LoadingButton className='!bg-red-500 !hover:bg-red-700' text='Deauthorize' loading={loading} />
          </div>
        </div>
      </form>
    )
  }

  const renderForm = () => {
    if (params?.isGithubType) {
      switch (githubStep) {
        case UCAN_GITHUB_STEP_1:
          return renderGithubSignature();
        case UCAN_GITHUB_STEP_2:
          return renderGithubDeauthorize()
      }
      return renderFilecoinDeauthorize();
    } else {
      return renderFilecoinDeauthorize();
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
                      d="m11 17l-5-5m0 0l5-5m-5 5h12" />
              </svg>
              Back
            </Link>
          </div>
        </button>
      </div>
      {
        renderForm()
      }
    </>
  )
}

export default UcanDelegate;
