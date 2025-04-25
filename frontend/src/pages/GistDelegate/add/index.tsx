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
import { useEffect, useRef, useState } from "react";
import { Controller, useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from "react-router-dom";
import type { BaseError } from "wagmi";
import { useAccount, useSignMessage, useWaitForTransactionReceipt, useWriteContract } from "wagmi";
import type { UserRejectedRequestError } from 'viem';
import { useSendMessage, useSign } from "iso-filecoin-react"
import {
  OPERATION_CANCELED_MSG,
  STORING_DATA_MSG,
  GITHUB_STEP_1,
  GITHUB_STEP_2,
  calibrationChainId,
  checkGistApi
} from "../../../common/consts";
import LoadingButton from "../../../components/LoadingButton";
import Table from '../../../components/Table';
import { getContractAddress, isFilAddress, validateValue } from "../../../utils"
import './index.less';
import oracleAbi from "../../../common/abi/oracle.json";
import MDEditor from "../../../../src/components/MDEditor";
import { useGistList, useTransactionHash } from "../../../common/store.ts";
import axios from "axios";
import { CopyButton } from "../../../components/CopyButton.tsx";
import { useFilAddressMessage } from "../../../common/hooks.ts"

const GistDelegate = () => {
  const { isConnected, address, chain } = useAccount();
  const { signMessageAsync } = useSignMessage();
  const { mutateAsync: sendMessage } = useSendMessage();
  const { mutateAsync: sign } = useSign();
  const { openConnectModal } = useConnectModal();
  const navigate = useNavigate();
  const prevAddressRef = useRef(address);
  const [messageApi, contextHolder] = message.useMessage();

  const [githubSignature, setGithubSignature] = useState('');
  const [githubStep, setGithubStep] = useState(GITHUB_STEP_1);
  const [filOperationSuccess, setFilOperationSuccess] = useState(false);
  const setGistList = useGistList((state: any) => state.setGistList);
  const gistList = useGistList((state: any) => state.gistList);
  const setStoringHash = useTransactionHash((state: any) => state.setTransactionHash)

  const [formValue] = useState({
    aud: '',
    prf: '',
    owner: '',
    repo: '',
    gistId: '',
    token: '',
  });

  const {
    register,
    handleSubmit,
    control,
    setValue: setFormValue,
    formState: { errors }
  } = useForm({
    defaultValues: {
      ...formValue,
    }
  });

  const {
    data: hash,
    error,
    writeContract,
    isPending: writeContractPending,
    isSuccess: writeContractSuccess,
    reset: resetWriteContract
  } = useWriteContract();
  const [loading, setLoading] = useState<boolean>(false);
  const [authorizeLoading, setAuthorizeLoading] = useState<boolean>(writeContractPending);
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
    if (writeContractSuccess || filOperationSuccess) {
      messageApi.open({
        type: 'success',
        content: t(STORING_DATA_MSG),
      });
      setStoringHash({ 'gistAudHash': hash })
      setTimeout(() => {
        navigate("/home");
      }, 3000);
    }
  }, [writeContractSuccess, filOperationSuccess])

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
    const params = {
      gistId: value.trim(),
      address
    }

    const { data } = await axios.get(checkGistApi, { params });
    if (!data.data) {
      messageApi.open({
        type: 'error',
        content: t('content.validGistIdInfo'),
      });
      setAuthorizeLoading(false);
      return
    }
    if (address && isFilAddress(address)) {
      try {
        const { message } = await useFilAddressMessage({
          abi: oracleAbi,
          contractAddress: getContractAddress(chain?.id || calibrationChainId, 'oracle'),
          address: address as string,
          functionName: "updateGistId",
          functionParams: [value],
        })
        await sendMessage(message);
        setFilOperationSuccess(true);
      } catch (error) {
        console.log(error)
        if (error as UserRejectedRequestError) {
          messageApi.open({
            type: "warning",
            content: t("content.rejectedSignature")
          })
        } else {
          messageApi.open({
            type: "error",
            content: (error as BaseError)?.shortMessage || JSON.stringify(error)
          })
        }
      } finally {
        setLoading(false);
      }
    } else {
      writeContract({
        abi: oracleAbi,
        address: getContractAddress(chain?.id || calibrationChainId, 'oracle'),
        functionName: 'updateGistId',
        args: [
          value
        ],
      });
    }
    setGistList([])
    setAuthorizeLoading(false);
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
          walletAddress: address,
          githubName: aud,
          timestamp: Math.floor(Date.now() / 1000)
        }
        let signature = "";
        try {
          // Sign the concatenated header and params
          if (address && isFilAddress(address)) {
            const jsonData = await sign(`${JSON.stringify(signatureParams)}`);
            signature = jsonData.toLotusHex();
          } else {
            signature = await signMessageAsync({ message: `${JSON.stringify(signatureParams)}` });
          }
        } catch (e) {
          messageApi.open({
            type: 'error',
            content: t(OPERATION_CANCELED_MSG),
          });
          setLoading(false);
          return;
        }

        const signatureStr = `
      I hereby claim:

        * I am ${aud} on Github.
        * I control ${address} (Filecoin wallet address).

      To claim this, I am signing this object

      ${JSON.stringify(signatureParams)}

      with my Filecoin wallet's private key, yielding the signature:

      ${signature}

      And finally, I am proving ownership of the github account by posting this as a gist.`
        // Set GitHub signature and step
        setGithubSignature(signatureStr);
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
    setAuthorizeLoading(true);
    const { gistId } = values;
    setValue(gistId);
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
            <p className='text-red-500 mt-1 text-sm'>{t('content.audRequired')}</p>
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
        <div style={{ position: 'relative', width: '830px' }}>
          <div style={{ position: 'absolute', right: 10, top: 10, zIndex: 100 }}><CopyButton text={githubSignature} /></div>
          <MDEditor
            value={githubSignature}
            readOnly={true}
            onChange={() => undefined}
            className="border-none rounded-[16px] bg-transparent"
            style={{ height: '100%', width: '100%', whiteSpace: 'pre-wrap', wordWrap: 'break-word' }}
            view={{ menu: false, md: false, html: true, both: false, fullScreen: true, hideMenu: false }}
          >
          </MDEditor>
          <a target="_blank" href="https://gist.github.com/" style={{ color: "blue", fontSize: '16px' }}>{t('content.getGistID')}↗</a>
        </div>
      )
    },
    {
      name: 'GistId',
      width: 100,
      comp: (
        <>
          <Controller
            name="gistId"
            control={control}
            render={() => <input
              className={classNames(
                'form-input rounded bg-[#ffffff] border border-[#eeeeee] text-black w-[830px]',
                errors.gistId && 'border-red-500 focus:border-red-500'
              )}
              {...register('gistId', { required: true, validate: validateValue })}
            />}
          />
          {errors.gistId && (
            <p className='text-red-500 mt-1 text-sm'>{t('content.gistIdRequired')}</p>
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
            {
              gistList && gistList.length > 0 && (
                <button
                  className={`h-[40px] bg-sky-500 hover:bg-sky-700 text-white py-2 px-6 rounded-xl disabled:opacity-50 mr-8 ${loading && 'cursor-not-allowed'}`}
                  type='button' onClick={() => { navigate('/gistDelegate/list'); }}>
                  {t('content.previous')}
                </button>
              )
            }
            <LoadingButton text={t('content.sign')} loading={loading} />
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
              className={`h-[40px] bg-sky-500 hover:bg-sky-700 text-white py-2 px-6 rounded-xl disabled:opacity-50 mr-8 ${authorizeLoading && 'cursor-not-allowed'}`}
              type='button' onClick={() => { setGithubStep(GITHUB_STEP_1); setFormValue('gistId', '') }}>
              {t('content.previous')}
            </button>
            <LoadingButton text={t('content.authorize')} loading={authorizeLoading || writeContractPending || transactionLoading} />
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