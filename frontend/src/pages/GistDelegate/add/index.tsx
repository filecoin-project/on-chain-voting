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

import { useConnectModal } from "@rainbow-me/rainbowkit";
import { message } from 'antd';
import classNames from 'classnames';
import React, { useEffect, useRef, useState } from "react";
import { Controller, useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from "react-router-dom";
import type { BaseError } from "wagmi";
import { useAccount, useSignMessage, useWaitForTransactionReceipt, useWriteContract } from "wagmi";
import {
  OPERATION_CANCELED_MSG,
  STORING_DATA_MSG,
  GITHUB_STEP_1,
  GITHUB_STEP_2,
  JWT_HEADER
} from "../../../common/consts";
import LoadingButton from "../../../components/LoadingButton";
import Table from '../../../components/Table';
import { stringToBase64Url, validateValue } from '../../../utils';
import './index.less';
const GistDelegate = () => {
  const { isConnected, address } = useAccount();
  const { signMessageAsync } = useSignMessage();
  const { openConnectModal } = useConnectModal();
  const navigate = useNavigate();
  const prevAddressRef = useRef(address);
  const [messageApi, contextHolder] = message.useMessage();

  const [githubSignature, setGithubSignature] = useState('');
  const [githubStep, setGithubStep] = useState(GITHUB_STEP_1);
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
    formState: { errors }
  } = useForm({
    defaultValues: {
      ...formValue,
    }
  });

  const {
    data: hash,
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
      const gistStorageData = JSON.parse(localStorage.getItem('gistStorage') || '[]');
      // Calculate expiration time (three minutes from now)
      const expirationTime = Date.now() + 3 * 60 * 1000;
      // Push new data (timestamp and address) to the array
      gistStorageData.push({ timestamp: expirationTime, address });
      // Save updated data to localStorage
      localStorage.setItem('gistStorage', JSON.stringify(gistStorageData));

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

  /**
   * create gist
    * @param values
   * @param githubStep
   */
  const onSubmit = (values: any, githubStep?: number) => {
    switch (githubStep) {
      case GITHUB_STEP_1:
        createSignature(values);
        break;
      case GITHUB_STEP_2:
        authorizeGithubValue(values);
        break;
    }

  }

  const setValue = async (value: string) => {
    console.log('value', value)
    //TODO

    setLoading(false);
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
        const base64Header = stringToBase64Url(JSON.stringify(JWT_HEADER));
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
        setGithubStep(GITHUB_STEP_2);
      } catch (e) {
        console.log(e);
      }
    } else {
      openConnectModal && openConnectModal();
    }
    setLoading(false);
  }

  const authorizeGithubValue = async (values: any) => {
    setLoading(true);
    const { url } = values;
    setValue(url);
  }

  const githubSignatureList = [
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


  const renderGithubSignature = () => {
    return (
      <form onSubmit={handleSubmit(value => { onSubmit(value, GITHUB_STEP_1) })}>
        <div className='flow-root space-y-8'>
          <Table
            title={t('content.githubDelegatesAuthorize')}
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
      <form onSubmit={handleSubmit(value => { onSubmit(value, GITHUB_STEP_2) })}>
        <div className='flow-root space-y-8'>
          <Table
            title={t('content.githubDelegatesAuthorize')}
            list={githubAuthorizeList}
          />

          <div className='text-center'>
            <button
              className={`h-[40px] bg-sky-500 hover:bg-sky-700 text-white py-2 px-6 rounded-xl disabled:opacity-50 mr-8 ${loading && 'cursor-not-allowed'}`}
              type='button' onClick={() => { setGithubStep(GITHUB_STEP_1) }}>
              {t('content.previous')}
            </button>
            <LoadingButton text={t('content.authorize')} loading={loading || writeContractPending || transactionLoading} />
          </div>
        </div>
      </form>
    )
  }

  const renderForm = () => {
    switch (githubStep) {
      case GITHUB_STEP_1:
        return renderGithubSignature();
      case GITHUB_STEP_2:
        return renderGithubAuthorize();

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

export default GistDelegate;
