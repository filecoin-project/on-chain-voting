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
import { message, DatePicker } from "antd";
import axios from 'axios';
import dayjs from "dayjs";
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
import { useNavigate, Link } from "react-router-dom";
import Table from '../../components/CreateTable';
import { useForm, Controller } from 'react-hook-form';
import classNames from 'classnames';
import type { BaseError } from "wagmi";
import { useAccount, useWriteContract, useWaitForTransactionReceipt } from "wagmi";
import { useConnectModal } from "@rainbow-me/rainbowkit";
import Editor from '../../components/MDEditor';
import {
  DEFAULT_TIMEZONE,
  WRONG_EXPIRATION_TIME_MSG,
  NOT_FIP_EDITOR_MSG,
  VOTE_OPTIONS,
  WRONG_START_TIME_MSG,
  STORING_DATA_MSG, githubApi, worldTimeApi,
  proposalDraftGetApi,
  proposalDraftAddApi,
  SAVE_DRAFT_SUCCESS,
  SAVE_DRAFT_FAIL,
  UPLOAD_DATA_FAIL_MSG,
  SAVE_DRAFT_TOO_LARGE
} from '../../common/consts';
import { useStoringCid, useVoterInfo } from "../../common/store";
import timezoneOption from '../../../public/json/timezons.json';
import { validateValue, getContractAddress, getWeb3IpfsId } from '../../utils';
import './index.less';
import LoadingButton from "../../components/LoadingButton";
import fileCoinAbi from "../../common/abi/power-voting.json";
import { useCheckFipEditorAddress } from "../../common/hooks";
import { ProposalDraft } from "src/common/types";

dayjs.extend(utc);
dayjs.extend(timezone);

const { RangePicker } = DatePicker;

const CreateVote = () => {
  const { isConnected, address, chain } = useAccount();
  const chainId = chain?.id || 0;

  const { openConnectModal } = useConnectModal();
  const prevAddressRef = useRef(address);
  const [messageApi, contextHolder] = message.useMessage();

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
      option: [
        { value: '' },
        { value: '' }
      ]
    }
  });

  const navigate = useNavigate();

  const { isFipEditorAddress } = useCheckFipEditorAddress(chainId, address);

  const {
    data: hash,
    writeContract,
    error,
    isPending: writeContractPending,
    isSuccess: writeContractSuccess,
    reset
  } = useWriteContract();

  const [cid, setCid] = useState('');
  const [loading, setLoading] = useState<boolean>(writeContractPending);
  const [isDraftSave, setDraftSave] = useState(false)
  const [hasDraft, setHasDraft] = useState(false)
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
    if (error) {
      messageApi.open({
        type: 'error',
        content: (error as BaseError)?.shortMessage || error?.message,
      });
    }
    reset();
  }, [error]);

  const OPTION_SPLIT_TAG = "&%"
  const loadDraft = async () => {
    try {
      const resp = await axios.get(proposalDraftGetApi, {
        params: {
          chainId: chainId,
          address: address

        }
      })
      if (resp.data != null && resp.data.data?.length) {
        const result = (resp.data.data as ProposalDraft[])[0]
        setValue("descriptions", result.Descriptions)
        setValue("name", result.Name)
        if(result.Time.length){
          setValue("time", result.Time.split(OPTION_SPLIT_TAG) ?? [])
        }
        if (result.Timezone) {
          setValue("timezone", result.Timezone)
        }
        setHasDraft(true)
      }
    } catch (e) {
      console.log(e)
    }

  }
  useEffect(() => {
    loadDraft()
  }, [])

  useEffect(() => {
    if (writeContractSuccess) {
      messageApi.open({
        type: 'success',
        content: STORING_DATA_MSG,
      });
      addStoringCid([{
        hash,
        cid
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
    const { data } = await axios.get(worldTimeApi);
    const currentTime = data?.unixtime;

    // Check if current time is after start time
    if (currentTime > startTimestamp) {
      messageApi.open({
        type: 'warning',
        content: WRONG_START_TIME_MSG,
      });
      setLoading(false);
      return false;
    }

    // Check if current time is after expiration time
    if (currentTime > expTimestamp) {
      messageApi.open({
        type: 'warning',
        content: WRONG_EXPIRATION_TIME_MSG,
      });
      setLoading(false);
      return false;
    }

    // Get text for timezone array
    const text = timezoneOption?.find((item: any) => item.value === values.value)?.text || '';
    // Extract GMT offset from text using regex
    const regex = /(?<=\().*?(?=\))/g;
    const GMTOffset = text.match(regex);

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

    // Prepare values object with additional information
    const _values = {
      ...values,
      ...githubObj,
      GMTOffset,
      startTime: startTimestamp,
      expTime: expTimestamp,
      showTime: values.time,
      option: VOTE_OPTIONS,
      address: address,
      chainId: chainId,
      currentTime,
    };

    const cid = await getWeb3IpfsId(_values);

    if (!cid?.length) {
      messageApi.open({
        type: 'warning',
        content: UPLOAD_DATA_FAIL_MSG,
      });
      setLoading(false);
      return
    }

    setCid(cid);

    if (isConnected) {
      // Check if user is a FIP editor
      if (isFipEditorAddress) {
        // Create voting using dynamic contract API
        writeContract({
          abi: fileCoinAbi,
          address: getContractAddress(chain?.id || 0, 'powerVoting'),
          functionName: 'createProposal',
          args: [
            cid,
            startTimestamp,
            expTimestamp,
            1
          ],
        });
      } else {
        messageApi.open({
          type: 'warning',
          content: NOT_FIP_EDITOR_MSG,
        });
      }
    } else {
      openConnectModal && openConnectModal();
    }
    //clear draft
    if (hasDraft) {
      clearDraft()
    }
    setLoading(false);
  }
  const clearDraft = async () => {
    try {
      const data = {
        Timezone: "",
        Time: "",
        Name: "",
        Descriptions: "",
        Option: "",
        Address: address,
        ChainId: chainId,
      }
      await axios.post(proposalDraftAddApi, data)
      setHasDraft(false)
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

    if(values.descriptions.length>=2048){
      messageApi.open({
        type: "warning",
        content: SAVE_DRAFT_TOO_LARGE,
      });
      return
    }
    setDraftSave(true)
    const data = {
      Timezone: values.timezone,
      Time: (values.time ?? []).join(OPTION_SPLIT_TAG),
      Name: values.name,
      Descriptions: values.descriptions,
      Option: "",
      Address: address,
      ChainId: chainId,
    }
    try {
      const res = await axios.post(proposalDraftAddApi, data)
      if (res.data != null && res.data.data == true) {
        messageApi.open({
          type: "success",
          content: SAVE_DRAFT_SUCCESS,
        });
      } else {
        messageApi.open({
          type: "error",
          content: SAVE_DRAFT_FAIL,
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
      name: 'Proposal Title',
      comp: (
        <>
          <Controller
            name="name"
            control={control}

            render={() => <input
              className={classNames(
                'form-input w-full rounded !bg-[#ffffff] border-1 border-[#EEEEEE] text-[#4B535B]',
                errors.name && 'border-red-500 focus:border-red-500'
              )}
              placeholder='Proposal Title'

              {...register('name', { required: true, validate: validateValue })}
            />}
          />
          {errors.name && (
            <p className='text-red-500 mt-1'>Proposal Title is required</p>
          )}
        </>
      )
    },
    {
      name: 'Description',
      width: 280,
      desc: <div className="text-red">
        <span className="text-sm">
          Describe FIP objectives, implementation details, risks, and include GitHub links for transparency. See a template <a target="_blank"
            rel="noopener" href="" className="text-sm" style={{ color: "blue" }}>here↗</a>.
          <br /> You can use Markdown formatting in the text input field.
        </span>

      </div>,
      comp:
        <Controller
          name='descriptions'
          control={control}
          rules={{
            required: true,
            validate: validateValue
          }}
          render={({ field: { onChange, value } }) => {
            return (
              <>
                <Editor
                  style={{ height: 500 }} value={value} onChange={onChange} />
                {errors.descriptions && (
                  <p className='text-red-500 mt-2'>Proposal Description is required</p>
                )}
              </>
            )
          }}
        />
    },
    {
      name: 'Voting Time',
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
                      placeholder={['Start Time', 'End Time']}
                      allowClear={true}
                      value={[date[0], date[1]]}
                      onChange={onChange}
                      className={classNames(
                        'form-input rounded w-[450px] !bg-[#ffffff] border border-[#eeeeee]',
                        errors.time && 'border-red-500 focus:border-red-500'
                      )}
                    />
                    {errors.time && (
                      <p className='text-red-500 mt-2'>Proposal Time is required</p>
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
      name: 'Timezone',
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
                        'form-select rounded bg-[#ffffff] border border-[#eeeeee] text-[#4B535B]',
                        errors.timezone && 'border-red-500 focus:border-red-500'
                      )}
                    >
                      {timezoneOption.map((option: any) => (
                        <option value={option.value} key={option.value}>{option.text}</option>
                      ))}
                    </select>
                    {errors.timezone && (
                      <p className='text-red-500 mt-2'>Proposal Expiration TimeZone is required</p>
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
              Back
            </Link>
          </div>
        </button>
      </div>
      <form onSubmit={handleSubmit(onSubmit)} >
        <div className='flow-root space-y-8'>
          <Table title='Create A Proposal' subTitle={<div className="text-base font-normal">
            Proposals should be clear, concise, and focused on specific improvements or changes. FIPs must adhere to the Filecoin community's <a target="_blank"
              rel="noopener" href="" style={{ color: "blue" }}>code of conduct and best practices↗</a>.
          </div>} list={list} />

          <div className="flex justify-center items-center text-center ">
            <Link to="/home" >
              <div className="flex justify-center rounded items-center text-center  bg-[#EEEEEE] w-[101px] h-[40px] text-[#313D4F] mr-2 cursor-pointer" >Cancel</div>
            </Link>
            <div className='w-full items-center flex justify-end text-center'>
              <Link to={""}>
                <div className="text-[#313D4F] mr-[32px] font-semibold cursor-pointer" onClick={saveDraft} >
                  Save Draft
                </div>
              </Link>
              <LoadingButton className="create-submit" text='Create Proposal' loading={loading || writeContractPending || transactionLoading} />
            </div>
          </div>
        </div>
      </form>
    </>
  )
}

export default CreateVote;
