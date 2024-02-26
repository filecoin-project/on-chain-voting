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

import React, {useState} from "react";
import {message, DatePicker} from "antd";
import dayjs from "dayjs";
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
import {useNavigate, Link} from "react-router-dom";
import Table from '../../components/Table';
import {useForm, Controller, useFieldArray} from 'react-hook-form';
import classNames from 'classnames';
import {RadioGroup} from '@headlessui/react';
import {useAccount, useNetwork} from "wagmi";
import {useConnectModal} from "@rainbow-me/rainbowkit";
import Editor from '../../components/MDEditor';
import {
  VOTE_TYPE_OPTIONS,
  SINGLE_VOTE,
  DEFAULT_TIMEZONE,
  STORING_DATA_MSG,
  WRONG_EXPIRATION_TIME_MSG,
} from '../../common/consts';
import {useDynamicContract, getIpfsId} from "../../hooks";
import { useTimezoneSelect, allTimezones } from 'react-timezone-select';
import './index.less';

dayjs.extend(utc);
dayjs.extend(timezone);

const CreateVote = () => {
  const {chain} = useNetwork();
  const {isConnected, address} = useAccount();
  const {openConnectModal} = useConnectModal();

  const {
    register,
    handleSubmit,
    control,
    formState: {errors}
  } = useForm({
    defaultValues: {
      proposalType: 1,
      voteType: SINGLE_VOTE,
      timezone: DEFAULT_TIMEZONE,
      expTime: '',
      name: '',
      descriptions: '',
      option: [
        {value: ''},
        {value: ''}
      ]
    }
  })

  const labelStyle = 'original'

  const { options } = useTimezoneSelect({ labelStyle, timezones: allTimezones });

  // ts-ignore
  const {fields, append, remove} = useFieldArray({
    name: 'option',
    control,
    rules: {
      required: true
    }
  });

  const navigate = useNavigate();

  const [loading, setLoading] = useState<boolean>(false);

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
    const expTimestamp = dayjs(values.expTime).add(offset, 'minute').unix();
    const currentTime = dayjs().unix();

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
      expTime: expTimestamp,
      showTime: values.expTime,
      option: values.option.map((item: { value: string }) => item.value),
      address: address,
      chainId: chainId,
      currentTime,
    };

    const cid = await getIpfsId(_values);

    if (isConnected) {
      const { createVotingApi } = useDynamicContract(chainId);
      const res = await createVotingApi(cid, expTimestamp, values.proposalType);
      if (res.code === 200) {
        message.success(STORING_DATA_MSG);
        navigate("/");
      } else {
        message.error(res.msg);
      }
    } else {
      // @ts-ignore
      openConnectModal && openConnectModal();
    }
  }

  const list = [
    {
      name: 'Proposal Type',
      comp: (
        <div className='flex items-center'>
          <div className='w-full'>
            <Controller
              name='proposalType'
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
                        'form-select w-full rounded bg-[#212B3C] border border-[#313D4F]',
                        errors['proposalType'] && 'border-red-500 focus:border-red-500'
                      )}
                    >
                      <option value={1} key='proposal'>Proposal</option>
                    </select>
                    {errors['proposalType'] && (
                      <p className='text-red-500 mt-2'>Proposal Type is required</p>
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
      name: 'Proposal Title',
      comp: (
        <>
          <Controller
            name="name"
            control={control}
            render={({ field }) => <input
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
      name: 'Proposal Expiration Time',
      comp: (
        <div className='flex items-center'>
          <div className='mr-2.5'>
            <Controller
              name='expTime'
              control={control}
              rules={{ required: true }}
              render={({field: {onChange, value}}) => {
                const dateValue = value ? dayjs(value) : null;
                return (
                  <>
                    <DatePicker
                      showTime
                      format="YYYY-MM-DD HH:mm"
                      allowClear={false}
                      value={dateValue}
                      onChange={onChange}
                      // disabledDate={disabledDate}
                      // disabledTime={disabledTime}
                      className={classNames(
                        'form-input rounded bg-[#212B3C] border border-[#313D4F]',
                        errors['expTime'] && 'border-red-500 focus:border-red-500'
                      )}
                      style={{color: 'red'}}
                      placeholder='Pick Date'
                    />
                    {errors['expTime'] && (
                      <p className='text-red-500 mt-2'>Proposal Expiration Time is required</p>
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
      name: 'Proposal Expiration Timezone',
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
    {
      name: 'Voting Type',
      comp: (
        <Controller
          name='voteType'
          control={control}
          render={({field: {onChange, value}}) => {
            return (
              <RadioGroup className='flex' value={value} onChange={onChange}>
                {VOTE_TYPE_OPTIONS.map(poll => (
                  <RadioGroup.Option
                    key={poll.label}
                    value={poll.value}
                    className={
                      'relative flex items-center cursor-pointer p-4 focus:outline-none'
                    }
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
                            {poll.label}
                          </RadioGroup.Label>
                        </span>
                      </>
                    )}
                  </RadioGroup.Option>
                ))}
              </RadioGroup>
            )
          }}
        />
      )
    },
    {
      name: 'Proposal Options',
      comp: (
        <>
          <div className='rounded border border-[#313D4F] divide-y divide-[#212B3C]'>
            <div className='flex justify-between bg-[#293545] text-base text-[#8896AA] px-5 py-4'>
              <span>Options</span>
              <span>Operations</span>
            </div>
            {fields.map((field: any, index: number) => (
              <div key={field.id}>
                <div className='flex items-center pl-2.5 py-2.5 pr-5'>
                  <input
                    type='text'
                    maxLength={40}
                    className={classNames(
                      'form-input flex-auto rounded bg-[#212B3C] border border-[#313D4F]',
                      errors.option && errors.option[index]?.value &&
                      'border-red-500 focus:border-red-500'
                    )}
                    placeholder='Edit Option'
                    {...register(`option.${index}.value`, {required: true, validate: validateValue})}
                  />
                  {
                    fields.length > 2 &&
                      <button
                          type='button'
                          onClick={() => remove(index)}
                          className='ml-3 w-[50px] h-[50px] flex justify-center items-center bg-[#212B3C] rounded-full'
                      >
                          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                               stroke="currentColor" aria-hidden="true"
                               className="h-5 w-5 text-[#8896AA] hover:opacity-80">
                              <path stroke-linecap="round" stroke-linejoin="round"
                                    d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"></path>
                          </svg>
                      </button>
                  }
                </div>
                {errors.option && errors.option[index]?.value && (
                  <div className='px-2.5 pb-3'>
                    <p className='text-red-500 text-base'>
                      Option Name is required
                    </p>
                  </div>
                )}
              </div>
            ))}
          </div>
          {
            fields.length < 5 &&
              <div className='pl-2.5 py-4'>
                  <button
                      type='button'
                      onClick={() => append({ value: '' })}
                      className='px-8 py-3 rounded border border-[#313D4F] bg-[#3B495B] text-base text-white hover:opacity-80'
                  >
                      Add Option
                  </button>
              </div>
          }
        </>
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
            <button
              className={`h-[40px] bg-sky-500 hover:bg-sky-700 text-white py-2 px-6 rounded-xl disabled:opacity-50 ${loading && 'cursor-not-allowed'}`}
              type='submit' disabled={loading}>
              Create
            </button>
          </div>
        </div>
      </form>
    </>
  )
}

export default CreateVote;
