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
import dayjs from "dayjs";
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
import {useNavigate, Link} from "react-router-dom";
import Table from '../../components/Table';
import {useForm, Controller} from 'react-hook-form';
import classNames from 'classnames';
import {useAccount} from "wagmi";
import {useConnectModal} from "@rainbow-me/rainbowkit";
import Editor from '../../components/MDEditor';
import {
  DEFAULT_TIMEZONE,
  WRONG_EXPIRATION_TIME_MSG,
  NOT_FIP_EDITOR_MSG, VOTE_OPTIONS, WRONG_START_TIME_MSG,
} from '../../common/consts';
import {useStaticContract, useDynamicContract, getWeb3IpfsId} from "../../hooks";
import { useTimezoneSelect, allTimezones } from 'react-timezone-select';
import './index.less';
import LoadingButton from "../../components/LoadingButton";

dayjs.extend(utc);
dayjs.extend(timezone);

const { RangePicker } = DatePicker;

const CreateVote = () => {
  const {isConnected, address, chain} = useAccount();
  const {openConnectModal} = useConnectModal();
  const prevAddressRef = useRef(address);

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

  const [loading, setLoading] = useState<boolean>(false);

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

  /**
   * create proposal
   * @param values
   */
  const onSubmit = async (values: any) => {
    setLoading(true)
    const offset =  dayjs().tz(values.timezone).utcOffset() - dayjs().utcOffset();
    const startTimestamp = dayjs(values.time[0]).add(offset, 'minute').unix();
    const expTimestamp = dayjs(values.time[1]).add(offset, 'minute').unix();
    const currentTime = dayjs().unix();

    if (currentTime > startTimestamp) {
      message.warning(WRONG_START_TIME_MSG);
      setLoading(false);
      return false;
    }

    if (currentTime > expTimestamp) {
      message.warning(WRONG_EXPIRATION_TIME_MSG);
      setLoading(false);
      return false;
    }

    const chainId = chain?.id || 0;
    const label = options?.find(item => item.value === values.timezone)?.label || '';
    const regex = /(?<=\().*?(?=\))/g;
    const GMTOffset = label.match(regex);

    const _values = {
      ...values,
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
      const { isFipEditor } = await useStaticContract(chainId);
      const res = await isFipEditor(address || '');
      if (res.code === 200 && res.data) {
        const { createVotingApi } = useDynamicContract(chainId);
        const res1 = await createVotingApi(cid, startTimestamp, expTimestamp, 1);
        if (res1.code === 200 && res1.data?.hash) {
          message.success(res1.msg);
          navigate("/");
        } else {
          message.error(res1.msg);
        }
      } else {
        message.warning(NOT_FIP_EDITOR_MSG);
      }
      setLoading(false);
    } else {
      // @ts-ignore
      openConnectModal && openConnectModal();
    }
  }

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
                'form-input w-full rounded bg-[#212B3C] border border-[#313D4F]',
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
            <LoadingButton text='Create' loading={loading} />
          </div>
        </div>
      </form>
    </>
  )
}

export default CreateVote;
