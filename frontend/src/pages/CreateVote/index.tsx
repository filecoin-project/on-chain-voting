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
import { DatePicker, InputNumber, message } from "antd";
import axios from 'axios';
import classNames from 'classnames';
import dayjs from "dayjs";
import timezone from 'dayjs/plugin/timezone';
import utc from 'dayjs/plugin/utc';
import React, { useEffect, useRef, useState } from "react";
import { Controller, useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from "react-router-dom";
import { ProposalDraft } from "src/common/types";
import { UserRejectedRequestError } from "viem";
import type { BaseError } from "wagmi";
import { useAccount, useWaitForTransactionReceipt, useWriteContract } from "wagmi";
import fileCoinAbi from "../../common/abi/power-voting.json";

import CreateTable from "src/components/CreateTable";
import timezoneOption from '../../../public/json/timezons.json';
import {
  calibrationChainId,
  DEFAULT_TIMEZONE,
  githubApi,
  NOT_FIP_EDITOR_MSG,
  proposalDraftAddApi,
  proposalDraftGetApi,
  SAVE_DRAFT_FAIL,
  SAVE_DRAFT_SUCCESS,
  SAVE_DRAFT_TOO_LARGE,
  STORING_DATA_MSG,
  WRONG_EXPIRATION_TIME_MSG,
  WRONG_START_TIME_MSG
} from "../../common/consts";
import { useFipList, useSearchValue, useStoringCid, useVoterInfo } from "../../common/store";
import LoadingButton from "../../components/LoadingButton";
import Editor from '../../components/MDEditor';
import { getContractAddress, hexToString, validateValue } from "../../utils";
import './index.less';

dayjs.extend(utc);
dayjs.extend(timezone);

const { RangePicker } = DatePicker;

const CreateVote = () => {
  const { isConnected, address, chain } = useAccount();
  const chainId = chain?.id || calibrationChainId;
  const { t } = useTranslation();
  const { openConnectModal } = useConnectModal();
  const prevAddressRef = useRef(address);
  const [messageApi, contextHolder] = message.useMessage();
  const setSearchValue = useSearchValue((state: any) => state.setSearchValue);
  const voterInfo = useVoterInfo((state: any) => state.voterInfo);
  const addStoringCid = useStoringCid((state: any) => state.addStoringCid);
  const {
    register,
    handleSubmit,
    control,
    setValue,
    getValues,
    formState: { errors }
  } = useForm({
    defaultValues: {
      timezone: DEFAULT_TIMEZONE,
      time: [] as string[],
      name: '',
      descriptions: '',
      percent: {
        spPercentage: 25,
        clientPercentage: 25,
        developerPercentage: 25,
        tokenHolderPercentage: 25,
      },
      option: [
        { value: '' },
        { value: '' }
      ]
    }
  });

  const navigate = useNavigate();

  const { isFipEditorAddress } = useFipList((state: any) => state.data);

  const {
    data: hash,
    writeContract,
    isPending: writeContractPending,
    isSuccess: writeContractSuccess,
    error,
    reset
  } = useWriteContract();
  // const [cid, setCid] = useState('');
  const [loading, setLoading] = useState<boolean>(writeContractPending);
  const [isDraftSave, setDraftSave] = useState(false);
  const [hasDraft, setHasDraft] = useState(false);
  useEffect(() => {
    if (error) {
      const errorStr = JSON.stringify(error);
      // Get error cause
      if (error as UserRejectedRequestError) {
        messageApi.open({
          type: 'warning',
          content: t('content.rejectedSignature'),
        });

      } else {
        // Intercepts the first hexadecimal in the string
        const reg = /revert reason:\s*0x[0-9A-Fa-f]+/;
        const match = errorStr.match(reg) || [];
        messageApi.open({
          type: 'error',
          content: hexToString(match[0]) || (error as BaseError)?.shortMessage,
        });
      }
    }
    reset()
  }, [error])
  useEffect(() => {
    if (!isConnected) {
      navigate("/home");
      setSearchValue('')
      return;
    }
  }, []);

  useEffect(() => {
    const prevAddress = prevAddressRef.current;
    if (prevAddress !== address) {
      navigate("/home");
      setSearchValue('')
    }
  }, [address]);
  // const OPTION_SPLIT_TAG = "&%";
  const loadDraft = async () => {
    try {
      const resp = await axios.get(proposalDraftGetApi, {
        params: {
          chainId: chainId,
          address: address
        }
      });
      if (resp.data != null && resp.data.data) {
        const result = (resp.data.data as ProposalDraft)
        setValue("descriptions", result.content)
        setValue("name", result.title)
        if (result.startTime && result.endTime) {
          setValue("time", [dayjs(result.startTime * 1000).toString(), dayjs(result.endTime * 1000).toString()])
        }
        if (result.timezone) {
          setValue("timezone", result.timezone)
        }
        setValue("percent", {
          spPercentage: result?.spPercentage,
          clientPercentage: result?.clientPercentage,
          developerPercentage: result?.developerPercentage,
          tokenHolderPercentage: result?.tokenHolderPercentage,
        })

        setHasDraft(true)
      }
    } catch (e) {
      console.log(e)
    }
  }

  useEffect(() => {
    loadDraft();
  }, []);

  useEffect(() => {
    if (writeContractSuccess) {
      messageApi.open({
        type: 'success',
        content: t(STORING_DATA_MSG),
      });
      addStoringCid([{
        hash,
      }]);
      setTimeout(() => {
        navigate("/home")
      }, 1000);
    }
  }, [writeContractSuccess])

  /**
   * create proposal
   * @param values
   */
  const onSubmit = async (values: any) => {
    setLoading(true);

    // Calculate offset based on selected timezone
    const offset = dayjs().utcOffset() - dayjs().tz(values.timezone).utcOffset();
    const startTimestamp = dayjs(values.time[0]).add(offset, 'minute').unix();
    const expTimestamp = dayjs(values.time[1]).add(offset, 'minute').unix();
    const currentTime = Math.floor(Date.now() / 1000);
    // Check if the role proportion is equal to 100
    if (values.percent.clientPercentage + values.percent.developerPercentage + values.percent.spPercentage + values.percent.tokenHolderPercentage !== 100) {
      messageApi.open({
        type: "warning",
        content: t('content.infoPercentage'),
      });
      setLoading(false);
      return false;
    }
    // Check if current time is after start time
    if (currentTime > startTimestamp) {
      messageApi.open({
        type: 'warning',
        content: t(WRONG_START_TIME_MSG),
      });
      setLoading(false);
      return false;
    }

    // Check if current time is after expiration time
    if (currentTime > expTimestamp) {
      messageApi.open({
        type: 'warning',
        content: t(WRONG_EXPIRATION_TIME_MSG),
      });
      setLoading(false);
      return false;
    }
    if (startTimestamp >= expTimestamp) {
      messageApi.open({
        type: 'warning',
        content: t('content.endTimeLoss'),
      });
      setLoading(false);
      return false;
    }
    const githubObj = {
      githubName: '',
      githubAvatar: ''
    }

    if (voterInfo && voterInfo[0]) {
      const githubName = voterInfo[0];
      const { data } = await axios.get(`${githubApi}/${githubName}`);
      githubObj.githubName = githubName;
      githubObj.githubAvatar = data.avatar_url;
    }

    if (isConnected) {
      // Check if user is a FIP editor
      if (isFipEditorAddress) {
        // Create voting using dynamic contract API
        writeContract({
          abi: fileCoinAbi,
          address: getContractAddress(chain?.id || calibrationChainId, 'powerVoting'),
          functionName: 'createProposal',
          args: [
            startTimestamp,
            expTimestamp,
            // data.blockHeight,
            values.percent.tokenHolderPercentage * 100,
            values.percent.spPercentage * 100,
            values.percent.clientPercentage * 100,
            values.percent.developerPercentage * 100,
            values.descriptions,
            values.name,
          ],
        });

        //clear draft
        if (hasDraft) {
          clearDraft()
        }

        setSearchValue('')

      } else {
        messageApi.open({
          type: 'warning',
          content: t(NOT_FIP_EDITOR_MSG),
        });
      }
    } else {
      openConnectModal && openConnectModal();
    }
    setLoading(false);
  }
  const clearDraft = async () => {
    try {
      const data = {
        creator: address,
        title: '',
        content: '',
        startTime: 0,
        endTime: 0,
        chainId: chainId,
        spPercentage: 25,
        clientPercentage: 25,
        developerPercentage: 25,
        tokenHolderPercentage: 25,
        timezone: DEFAULT_TIMEZONE,
      }
      await axios.post(proposalDraftAddApi, data)
      setHasDraft(false);
    } catch (e) {
      console.log(e)
    }
  }

  const saveDraft = async () => {
    if (loading || writeContractPending || transactionLoading) {
      return
    }
    if (isDraftSave) {
      return
    }
    const values = getValues()
    if (!values.descriptions.length
      && !values.name
      && !values.time
    ) {
      return
    }

    if (values.descriptions.length >= 2048) {
      messageApi.open({
        type: "warning",
        content: t(SAVE_DRAFT_TOO_LARGE),
      });
      return
    }
    if (values.percent.clientPercentage + values.percent.developerPercentage + values.percent.spPercentage + values.percent.tokenHolderPercentage !== 100) {
      messageApi.open({
        type: "warning",
        content: t('content.infoPercentage'),
      });
      return
    }
    setDraftSave(true);
    const githubObj = {
      githubName: '',
      githubAvatar: ''
    }
    if (voterInfo && voterInfo[0]) {
      const githubName = voterInfo[0];
      const { data } = await axios.get(`${githubApi}/${githubName}`);
      githubObj.githubName = githubName;
      githubObj.githubAvatar = data.avatar_url;
    }
    const offset = dayjs().utcOffset() - dayjs().tz(values.timezone).utcOffset();
    const startTimestamp = dayjs(values.time[0]).add(offset, 'minute').unix();
    const expTimestamp = dayjs(values.time[1]).add(offset, 'minute').unix();
    const data = {
      creator: address,
      title: values.name,
      content: values.descriptions,
      // ...githubObj,
      // GMTOffset,
      startTime: startTimestamp,
      endTime: expTimestamp,
      chainId: chainId,
      // currentTime,
      // Time: (values.time ?? []).join(OPTION_SPLIT_TAG),
      spPercentage: values.percent.spPercentage,
      clientPercentage: values.percent.clientPercentage,
      developerPercentage: values.percent.developerPercentage,
      tokenHolderPercentage: values.percent.tokenHolderPercentage,
      timezone: values.timezone,
    }
    try {
      const res = await axios.post(proposalDraftAddApi, data)
      if (res.data != null && res.data.code === 0) {
        messageApi.open({
          type: "success",
          content: t(SAVE_DRAFT_SUCCESS),
        });
      } else {
        messageApi.open({
          type: "error",
          content: t(SAVE_DRAFT_FAIL),
        });
      }
    } catch (e) {
      console.log(e)
    }
    setTimeout(() => {
      setDraftSave(false)
    }, 3000)
  }

  const { isLoading: transactionLoading } =
    useWaitForTransactionReceipt({
      hash,
    })

  const inputRefs = useRef<(HTMLElement | null)[]>([]);

  useEffect(() => {
    let index = 0
    document.querySelectorAll('input').forEach(element => {
      inputRefs.current[index++] = element
    })
    //input => editor
    inputRefs.current.splice(1, 1, document.querySelector("textarea"))

    //timezone
    inputRefs.current[index++] = document.querySelector(".form-select")

    //creat btn
    inputRefs.current[index++] = document.querySelector(".create-submit")

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Tab') {
        const refList = inputRefs.current
        const length = refList.length ?? 0
        const currentIndex = refList.findIndex(ref => ref === document.activeElement);
        inputRefs.current[(currentIndex + 1) % length]?.focus();
        event.preventDefault();
      }
    };
    document.addEventListener('keydown', handleKeyDown);
    return () => {
      document.removeEventListener('keydown', handleKeyDown);
    };
  }, [])
  // const disabledDate = (current:any) => {
  //   return current && current < dayjs().startOf('day');
  // };
  //VOTE_OPTIONS
  const list = [
    {
      name: t('content.proposalTitle'),
      comp: (
        <Controller
          name="name"
          control={control}

          render={({ field: { value } }) => {
            return (
              <>
                <input
                  className={classNames(
                    'form-input w-full rounded !bg-[#ffffff] border-1 border-[#EEEEEE] text-[#4B535B]',
                    errors.name && 'border-red-500 focus:border-red-500'
                  )}
                  placeholder={t('content.proposalTitle')}
                  {...register('name', { required: true, validate: (v) => validateValue(v) && v.length <= 100 })}
                />
                {errors.name && (
                  value.length > 100 ? <p className='text-red-500 mt-1 text-sm'>{t('content.proposalTitleLengthValid')}</p> :
                    <p className='text-red-500 mt-1 text-sm'>{t('content.proposalTitleRequired')}</p>
                )}
              </>
            )
          }
          }
        />

      )
    },
    {
      name: t('content.description'),
      width: 280,
      desc: <div className="text-red">
        <span className="text-sm">
          {t('content.describeFIPObjectives')} <a target="_blank"
            rel="noopener" href="" className="text-sm" style={{ color: "blue" }}>{t('content.here')}↗</a>.
          <br /> {t('content.markdownFormattingInField')}.
        </span>

      </div>,
      comp:
        <Controller
          name='descriptions'
          control={control}
          rules={{
            required: true,
            validate: (v) => validateValue(v) && v.length < 10000
          }}
          render={({ field: { onChange, value } }) => {
            return (
              <>
                <Editor
                  style={{ height: 500 }} value={value} onChange={onChange} />
                {errors.descriptions && (
                  value.length >= 10000 ? <p className='text-red-500 mt-1 text-sm'>{t('content.proposalDescriptionLengthValid')}</p> :
                    <p className='text-red-500 mt-1 text-sm'>{t('content.proposalDescriptionRequired')}</p>
                )}
              </>
            )
          }}
        />
    },
    {
      name: t('content.powerPercentage'),
      comp: (
        <div className='flex items-center'>
          <div className='mr-2.5'>
            <Controller
              name='percent'
              control={control}
              rules={{
                required: true,
                validate: (v) => !!v.clientPercentage && !!v.developerPercentage && !!v.spPercentage && !!v.tokenHolderPercentage
              }}
              render={({ field: { onChange, value: data } }) => {
                return (
                  <div className="gap-[10px] flex">
                    <div>
                      <p className="text-[#4B535B] text-sm mb-[1px]">SP</p>
                      <InputNumber
                        className={classNames(
                          !data.spPercentage && '!border-red-500 focus:!border-red-500'
                        )}
                        value={data.spPercentage}
                        style={{ width: '148px', border: '1px solid #EEEEEE' }}
                        min={0}
                        max={100}
                        placeholder={'SP'}
                        suffix="%"
                        onChange={(v) => onChange({ ...data, spPercentage: v })}
                        precision={2}
                      />
                      {errors.percent && !data.spPercentage && (
                        <p className='text-red-500 mt-2 text-sm'>{t('content.percentageRequired')}</p>
                      )}
                    </div>
                    <div>
                      <p className="text-[#4B535B] text-sm mb-[1px]">Client</p>
                      <InputNumber
                        className={classNames(
                          !data.clientPercentage && '!border-red-500 focus:!border-red-500'
                        )}
                        value={data.clientPercentage}
                        type="number"
                        style={{ width: '148px' }}
                        min={0}
                        max={100}
                        placeholder={'Client'}
                        suffix="%"
                        onChange={(v) => onChange({ ...data, clientPercentage: v })}
                        precision={2}
                      />
                      {errors.percent && !data.clientPercentage && (
                        <p className='text-red-500 mt-2 text-sm'>{t('content.percentageRequired')}</p>
                      )}
                    </div>
                    <div>
                      <p className="text-[#4B535B] text-sm mb-[1px]">Developer</p>
                      <InputNumber
                        className={classNames(
                          !data.developerPercentage && '!border-red-500 focus:!border-red-500'
                        )}
                        value={data.developerPercentage}
                        type="number"
                        style={{ width: '148px' }}
                        min={0}
                        max={100}
                        placeholder={'Developer'}
                        suffix="%"
                        onChange={(v) => onChange({ ...data, developerPercentage: v })}
                        precision={2}
                      />
                      {errors.percent && !data.developerPercentage && (
                        <p className='text-red-500 mt-2 text-sm'>{t('content.percentageRequired')}</p>
                      )}
                    </div>
                    <div>
                      <p className="text-[#4B535B] text-sm mb-[1px]">TokenHolder</p>
                      <InputNumber
                        className={classNames(
                          !data.tokenHolderPercentage && '!border-red-500 focus:!border-red-500'
                        )}
                        type="number"
                        value={data.tokenHolderPercentage}
                        style={{ width: '148px' }}
                        min={0}
                        max={100}
                        placeholder={'TokenHolder'}
                        suffix="%"
                        onChange={(v) => onChange({ ...data, tokenHolderPercentage: v })}
                        precision={2}
                      />
                      {errors.percent && !data.tokenHolderPercentage && (
                        <p className='text-red-500 mt-2 text-sm'>{t('content.percentageRequired')}</p>
                      )}
                    </div>
                  </div>
                )
              }}
            />
          </div>
        </div>
      )
    },
    {
      name: t('content.votingTime'),
      comp: (
        <div className='flex items-center'>
          <div className='mr-2.5'>
            <Controller
              name='time'
              control={control}
              rules={{ required: true }}
              render={({ field: { onChange, value } }) => {
                const date = (value ?? []).filter(it => it !== "").map(it => dayjs(it))
                return (
                  <>
                    <RangePicker
                      showTime
                      // disabledDate={disabledDate}
                      format="YYYY-MM-DD HH:mm"
                      placeholder={[t('content.startTime'), t('content.endTime')]}
                      allowClear={true}
                      value={[date[0], date[1]]}
                      onChange={onChange}
                      className={classNames(
                        'form-input rounded w-[450px] !bg-[#ffffff] border border-[#eeeeee]',
                        errors.time && 'border-red-500 focus:border-red-500'
                      )}
                    />
                    {errors.time && (
                      <p className='text-red-500 mt-2 text-sm'>{t('content.proposalTimeRequired')}</p>
                    )}
                  </>
                )
              }}
            />
          </div>
        </div>
      )
    },
    {
      name: t('content.timezone'),
      comp: (
        <div className='flex items-center'>
          <div className='mr-2.5'>
            <Controller
              name='timezone'
              control={control}
              rules={{
                required: true,
              }}
              render={({ field: { onChange, value } }) => {
                return (
                  <>
                    <select
                      onChange={onChange}
                      value={value}
                      className={classNames(
                        'form-select w-[450px] rounded bg-[#ffffff] border border-[#eeeeee] text-[#4B535B]',
                        errors.timezone && 'border-red-500 focus:border-red-500'
                      )}
                    >
                      {timezoneOption.map((option: any) => (
                        <option value={option.value} key={option.value}>{option.text}</option>
                      ))}
                    </select>
                    {errors.timezone && (
                      <p className='text-red-500 mt-2'>{t('content.proposalTimeZoneRequired')}</p>
                    )}
                  </>
                )
              }}
            />
          </div>
        </div>
      )
    }
  ];

  return (
    <>
      {contextHolder}
      <div className="px-3 mb-6 md:px-0">
        <button>
          <div className="inline-flex items-center gap-1 text-skin-text hover:text-skin-link">
            <Link to="/home" className="flex items-center">
              <svg className="mr-1" viewBox="0 0 24 24" width="1.2em" height="1.2em">
                <path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m11 17l-5-5m0 0l5-5m-5 5h12" />
              </svg>
              {t('content.back')}
            </Link>
          </div>
        </button>
      </div>
      <form onSubmit={handleSubmit(onSubmit)} >
        <div className='flow-root space-y-8'>
          <CreateTable title={t('content.createProposal')} subTitle={<div className="text-base font-normal">
            {t('content.proposalsClear')} <a target="_blank"
              rel="noopener" href="" style={{ color: "blue" }}>{t('content.codePractices')}↗</a>.
          </div>} list={list} />

          <div className="flex justify-center items-center text-center ">
            <Link to="/home" >
              <div className="flex justify-center rounded items-center text-center  bg-[#EEEEEE] w-[101px] h-[40px] text-[#313D4F] mr-2 cursor-pointer" >{t('content.cancel')}</div>
            </Link>
            <div className='w-full items-center flex justify-end text-center'>
              <Link to={""}>
                <div className="text-[#313D4F] mr-[32px] font-semibold cursor-pointer" onClick={saveDraft} >
                  {t('content.saveDraft')}
                </div>
              </Link>
              <LoadingButton className="create-submit" text={t('content.createProposals')} loading={loading || writeContractPending || transactionLoading} />
            </div>
          </div>
        </div>
      </form>
    </>
  )
}

export default CreateVote;
