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
import {message, DatePicker} from "antd";
import axios from 'axios';
import dayjs from "dayjs";
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
import {useNavigate, Link} from "react-router-dom";
import Table from '../../components/Table';
import {useForm, Controller} from 'react-hook-form';
import classNames from 'classnames';
import {useAccount, useReadContract, useWriteContract, useWaitForTransactionReceipt, BaseError} from "wagmi";
import {useConnectModal} from "@rainbow-me/rainbowkit";
import Editor from '../../components/MDEditor';
import {
  DEFAULT_TIMEZONE,
  WRONG_EXPIRATION_TIME_MSG,
  NOT_FIP_EDITOR_MSG,
  VOTE_OPTIONS,
  WRONG_START_TIME_MSG,
  STORING_DATA_MSG, githubApi
} from '../../common/consts';
import { useVoterInfo } from "../../common/store";
import { useTimezoneSelect, allTimezones } from 'react-timezone-select';
import { validateValue, getContractAddress, getWeb3IpfsId } from '../../utils';
import './index.less';
import LoadingButton from "../../components/LoadingButton";
import fileCoinAbi from "../../common/abi/power-voting.json";

dayjs.extend(utc);
dayjs.extend(timezone);

const { RangePicker } = DatePicker;

function useCheckFipAddress(chainId: number, address: `0x${string}` | undefined) {
  const { data: isFipAddress } = useReadContract({
    // @ts-ignore
    address: getContractAddress(chainId, 'powerVoting'),
    abi: fileCoinAbi,
    functionName: 'fipMap',
    args: [address]
  });
  return {
    isFipAddress
  };
}

const CreateVote = () => {
  const {isConnected, address, chain} = useAccount();
  const chainId = chain?.id || 0;

  const {openConnectModal} = useConnectModal();
  const prevAddressRef = useRef(address);
  const voterInfo = useVoterInfo((state: any) => state.voterInfo);
  const {
    register,
    handleSubmit,
    control,
    formState: {errors}
  } = useForm({
    defaultValues: {
      timezone: DEFAULT_TIMEZONE,
      time: '',
      name: '',
      descriptions: '',
      option: [
        {value: ''},
        {value: ''}
      ]
    }
  });

  const labelStyle = 'original';

  const { options } = useTimezoneSelect({ labelStyle, timezones: allTimezones });

  const navigate = useNavigate();

  const { isFipAddress } = useCheckFipAddress(chainId, address);

  const {
    data: hash,
    writeContract,
    error,
    isPending: writeContractPending,
    isSuccess: writeContractSuccess,
    reset
  } = useWriteContract();

  const [loading, setLoading] = useState<boolean>(writeContractPending);

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
      message.error((error as BaseError)?.shortMessage || error?.message);
    }
    reset();
  }, [error]);

  useEffect(() => {
    if (writeContractSuccess) {
      message.success(STORING_DATA_MSG);
      navigate("/");
    }
  }, [writeContractSuccess])

  /**
   * create proposal
   * @param values
   */
  const onSubmit = async (values: any) => {
    setLoading(true);
    // Calculate offset based on selected timezone
    const offset =  dayjs().tz(values.timezone).utcOffset() - dayjs().utcOffset();
    const startTimestamp = dayjs(values.time[0]).add(offset, 'minute').unix();
    const expTimestamp = dayjs(values.time[1]).add(offset, 'minute').unix();
    const currentTime = dayjs().unix();

    // Check if current time is after start time
    if (currentTime > startTimestamp) {
      message.warning(WRONG_START_TIME_MSG);
      setLoading(false);
      return false;
    }

    // Check if current time is after expiration time
    if (currentTime > expTimestamp) {
      message.warning(WRONG_EXPIRATION_TIME_MSG);
      setLoading(false);
      return false;
    }

    // Get chain ID
    const chainId = chain?.id || 0;
    // Get label for selected timezone
    const label = options?.find(item => item.value === values.timezone)?.label || '';
    // Extract GMT offset from label using regex
    const regex = /(?<=\().*?(?=\))/g;
    const GMTOffset = label.match(regex);

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

    if (isConnected) {
      // Check if user is a FIP editor
      if (isFipAddress) {
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
        message.warning(NOT_FIP_EDITOR_MSG);
      }
      setLoading(false);
    } else {
      // @ts-ignore
      openConnectModal && openConnectModal();
    }
  }

  const { isLoading: transactionLoading } =
    useWaitForTransactionReceipt({
      hash,
    })

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
                'form-input w-full rounded !bg-[#212B3C] border border-[#313D4F]',
                errors['name'] && 'border-red-500 focus:border-red-500'
              )}
              placeholder='Proposal Title'
              {...register('name', {required: true, validate: validateValue})}
            />}
          />
          {errors['name'] && (
            <p className='text-red-500 mt-1'>Proposal Title is required</p>
          )}
        </>
      )
    },
    {
      name: 'Proposal Description',
      comp:
        <Controller
          name='descriptions'
          control={control}
          rules={{
            required: true,
            validate: validateValue
          }}
          render={({field: {onChange, value}}) => {
            return (
              <>
                <Editor style={{height: 500}} value={value} onChange={onChange}/>
                {errors['descriptions'] && (
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
              render={({field: {onChange}}) => {
                return (
                  <>
                    <RangePicker
                      showTime
                      format="YYYY-MM-DD HH:mm"
                      placeholder={['Start Time', 'End Time']}
                      allowClear={false}
                      onChange={onChange}
                      className={classNames(
                        'form-input rounded !bg-[#212B3C] border border-[#313D4F]',
                        errors['time'] && 'border-red-500 focus:border-red-500'
                      )}
                      style={{color: 'red'}}
                    />
                    {errors['time'] && (
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
              render={({field: {onChange, value}}) => {
                return (
                  <>
                    <select
                      onChange={onChange}
                      value={value}
                      className={classNames(
                        'form-select rounded bg-[#212B3C] border border-[#313D4F]',
                        errors['timezone'] && 'border-red-500 focus:border-red-500'
                      )}
                    >
                      {options.map(option => (
                        <option value={option.value} key={option.value}>{option.label}</option>
                      ))}
                    </select>
                    {errors['timezone'] && (
                      <p className='text-red-500 mt-2'>Proposal Expiration TimeZone is required</p>
                    )}
                  </>
                )
              }}
            />
          </div>
        </div>
      )
    },
  ];

  return (
    <>
      <div className="px-3 mb-6 md:px-0">
        <button>
          <div className="inline-flex items-center gap-1 text-skin-text hover:text-skin-link">
            <Link to="/" className="flex items-center">
              <svg className="mr-1" viewBox="0 0 24 24" width="1.2em" height="1.2em">
                <path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m11 17l-5-5m0 0l5-5m-5 5h12" />
              </svg>
              Back
            </Link>
          </div>
        </button>
      </div>
      <form onSubmit={handleSubmit(onSubmit)}>
        <div className='flow-root space-y-8'>
          <Table title='Create A Proposal' list={list}/>

          <div className='text-center'>
            <LoadingButton text='Create' loading={loading || writeContractPending || transactionLoading} />
          </div>
        </div>
      </form>
    </>
  )
}

export default CreateVote;
