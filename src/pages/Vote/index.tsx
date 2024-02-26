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

import React, { useEffect, useState } from "react";
import { useNavigate, Link, useParams } from "react-router-dom";
import { useForm } from "react-hook-form";
import { Input, InputNumber, message } from "antd";
import axios from 'axios';
import dayjs from 'dayjs';
import {getIpfsId, useDynamicContract, useStaticContract} from "../../hooks";
import { useAccount, useNetwork } from "wagmi";
import { useConnectModal, useChainModal } from "@rainbow-me/rainbowkit";
import MDEditor from "../../components/MDEditor";
import EllipsisMiddle from "../../components/EllipsisMiddle";
import {
  IN_PROGRESS_STATUS,
  SINGLE_VOTE,
  MULTI_VOTE,
  WRONG_NET_STATUS,
  web3AvatarUrl, VOTE_SUCCESS_MSG, CHOOSE_VOTE_MSG, WRONG_MINER_ID_MSG,
} from "../../common/consts";
import { timelockEncrypt, roundAt, mainnetClient, Buffer } from "../../../tlock-js/src";
import {ProposalList, ProposalOption} from "../../common/types";
import "./index.less";

const totalPercentValue = 100;

const Vote = () => {
  const { chain } = useNetwork();
  const chainId = chain?.id || 0;
  const { isConnected, address } = useAccount();

  const { id, cid } = useParams();
  const [votingData, setVotingData] = useState({} as ProposalList);
  const { openConnectModal } = useConnectModal();
  const { openChainModal } = useChainModal();

  const navigate = useNavigate();
  const [options, setOptions] = useState([] as ProposalOption[]);
  const [minerIds, setMinerIds] = useState(['']);
  const [minerError, setMinerError] = useState(false);

  const [loading, setLoading] = useState(false);

  const {
    formState: { errors }
  } = useForm({
    defaultValues: {
      option: votingData?.voteType === MULTI_VOTE ? [] : null
    }
  })

  useEffect(() => {
    initState();
  }, [chain]);

  const addMinerIdPrefix = (minerIds: number[]) => {
    return minerIds.map(minerId => `f0${minerId}`);
  }

  const removeMinerIdPrefix = (minerIds: string[]) => {
    return minerIds.map(minerId => {
      const str = minerId.replace(/f0/g, '');
      if (isNaN(Number(str))) {
        setMinerError(true);
        return 0;
      } else {
        setMinerError(false);
        return Number(str)
      }
    });
  }

  const initState = async () => {
    const { getMinerIds } = await useStaticContract(chainId);
    const { code, data: { minerIds } } = await getMinerIds(address);
    if (code === 200) {
      setMinerIds(addMinerIdPrefix(minerIds.map((id: any) => id.toNumber())));
    }
    const res = await axios.get(`https://${cid}.ipfs.nftstorage.link/`);
    const data = res.data;
    let voteStatus = null;
    if (data.chainId !== chainId) {
      voteStatus = WRONG_NET_STATUS;
      if (isConnected) {
        openChainModal && openChainModal();
      } else {
        openConnectModal && openConnectModal();
      }
    } else {
      voteStatus = IN_PROGRESS_STATUS;
    }
    const option = data.option?.map((item: string) => {
      return {
        name: item,
        count: 0,
      };
    });
    setVotingData({
      ...data,
      id,
      cid,
      option,
      voteStatus,
    });
    setOptions(option);
  }

  /**
   * timelock encrypt
   * @param value
   */
  const handleEncrypt = async (value: string[][]) => {
    const payload = Buffer.from(JSON.stringify(value));

    const chainInfo = await mainnetClient().chain().info();

    const time = votingData?.expTime ? new Date(votingData.expTime * 1000).valueOf() : 0;

    const roundNumber = roundAt(time, chainInfo);

    const ciphertext = await timelockEncrypt(
      roundNumber,
      payload,
      mainnetClient()
    )

    return ciphertext;
  }

  const startVoting = async () => {
    const countIndex = options.findIndex((item: ProposalOption) => item.count > 0);
    if (countIndex < 0) {
      message.warning(CHOOSE_VOTE_MSG);
    } else {
      if (minerError) {
        message.warning(WRONG_MINER_ID_MSG);
        return;
      }
      setLoading(true);
      // vote params
      let params: string[][] = [];
      if (votingData?.voteType === SINGLE_VOTE) {
        params.push([`${countIndex}`, `${options[countIndex].count}`])
      } else {
        options.map((item: ProposalOption, index: number) => {
          params.push([`${index}`, `${item.count}`])
        })
      }
      const encryptValue = await handleEncrypt(params);
      const optionId = await getIpfsId(encryptValue);
      const ids = removeMinerIdPrefix(minerIds);

      if (isConnected) {
        const { voteApi } = useDynamicContract(chainId);
        const res = await voteApi(Number(id), optionId, ids);
        if (res.code === 200) {
          message.success(VOTE_SUCCESS_MSG, 3);
          setTimeout(() => {
            navigate("/")
          }, 3000)
        } else {
          message.error(res.msg)
        }
        setLoading(false)
      } else {
        // @ts-ignore
        openConnectModal && openConnectModal()
      }
    }
  }

  const handleMinerChange = (value: string) => {
    const arr = value.split(',');
    setMinerIds(arr);
  }

  const handleOptionChange = (index: number, count: number) => {
    setOptions((prevState: ProposalOption[]) => {
      return prevState.map((item: ProposalOption, preIndex) => {
        let currentTotal = 0;
        currentTotal += count;
        if (preIndex === index) {
          return {
            ...item,
            count
          }
        } else {
          return {
            ...item,
            disabled: votingData?.voteType === SINGLE_VOTE && count > 0 || votingData?.voteType === MULTI_VOTE && currentTotal === 100
          }
        }
      })
    })
  }

  const handleCountChange = (type: string, index: number, item: ProposalOption) => {
    if (item.disabled) return false;

    let currentCount: number;
    const restTotal = options.reduce(((acc: number, current: ProposalOption) => acc + current.count), 0) - item.count;
    const max = totalPercentValue - restTotal;
    const min = 0
    if (type === "decrease") {
      currentCount = item.count - 1 < min ? min : item.count - 1;
    } else {
      currentCount = item.count + 1 > max ? max : item.count + 1;
    }
    handleOptionChange(index, currentCount);
  }

  const countMax = (options: ProposalOption[], count: number) => {
    const restTotal = options.reduce(((acc: number, current: ProposalOption) => acc + current.count), 0) - count;
    return totalPercentValue - restTotal;
  }

  const handleVoteStatusTag = (status: number) => {
    switch (status) {
      case WRONG_NET_STATUS:
        return {
          name: 'Wrong network',
          color: 'bg-red-700',
        };
      case IN_PROGRESS_STATUS:
        return {
          name: 'In Progress',
          color: 'bg-green-700',
        };
      default:
        return {
          name: '',
          color: '',
        };
    }
  }

  return (
    <div className="flex voting">
      <div className="relative w-full pr-5 lg:w-8/12">
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
        </div>
        <div className="px-3 md:px-0 ">
          <h1 className="mb-6 text-3xl text-white break-words break-all leading-12" style={{overflowWrap: 'break-word'}}>
            {votingData?.name}
          </h1>
          {
            (votingData?.voteStatus || votingData?.voteStatus === 0) &&
              <div className="flex justify-between mb-6">
                  <div className="flex items-center justify-between w-full mb-1 sm:mb-0">
                      <button
                          className={`${handleVoteStatusTag(votingData.voteStatus).color} bg-[#6D28D9] h-[26px] px-[12px] text-white rounded-xl mr-4`}>
                        {handleVoteStatusTag(votingData.voteStatus).name}
                      </button>
                      <div className="flex items-center justify-center">
                          <img className="w-[20px] h-[20px] rounded-full mr-2" src={`${web3AvatarUrl}:${votingData.address}`} alt="" />
                          <a
                              className="text-white"
                              target="_blank"
                              rel="noopener"
                              href={`${chain?.blockExplorers?.default.url}/address/${votingData?.address}`}
                          >
                            {EllipsisMiddle({ suffixCount: 4, children: votingData?.address })}
                          </a>
                      </div>
                  </div>
              </div>
          }
          <div className="MDEditor">
            <MDEditor
              className="border-none rounded-[16px] bg-transparent"
              style={{ height: 'auto' }}
              moreButton={true}
              value={votingData?.descriptions}
              readOnly={true}
              view={{ menu: false, md: false, html: true, both: false, fullScreen: true, hideMenu: false }}
              onChange={() => {
              }}
            />
          </div>
          {
            votingData?.voteStatus === IN_PROGRESS_STATUS &&
              <div className='mt-5'>
                  <Input
                      defaultValue={minerIds.toString()}
                      onChange={(e) => { handleMinerChange(e.target.value) }}
                      className='form-input w-full rounded bg-[#212B3C] border border-[#313D4F] text-white'
                      placeholder='Input minerId'
                  />
                  <div className="border-[#313D4F] mt-6 border-skin-border bg-skin-block-bg text-base md:rounded-xl md:border border-solid">
                      <div className="group flex h-[57px] !border-[#313D4F] justify-between items-center border-b px-4 pb-[12px] pt-3 border-solid">
                          <h4 className="text-xl">
                              Cast Your Vote
                          </h4>
                          <div className='text-base'>{totalPercentValue} %</div>
                      </div>
                      <div className="p-4 text-center">
                        {
                          options.map((item: ProposalOption, index: number) => {

                            return (
                              <div className="mb-4 space-y-3 leading-10" key={item.name + index}>
                                <div
                                  className="w-full h-[45px] !border-[#313D4F] flex justify-between items-center pl-4 md:border border-solid rounded-full">
                                  <div className="text-ellipsis h-[100%] overflow-hidden">{item.name}</div>
                                  <div className="w-[180px] h-[45px] flex items-center">
                                    <div onClick={() => {
                                      handleCountChange("decrease", index, item)
                                    }}
                                         className={`w-[35px] border-x border-solid !border-[#313D4F] text-white font-semibold ${item.disabled ? "cursor-not-allowed" : "cursor-pointer"}`}>-
                                    </div>
                                    <InputNumber
                                      disabled={item.disabled}
                                      className="text-white bg-transparent focus:outline-none"
                                      controls={false}
                                      min={0}
                                      max={countMax(options, item.count)}
                                      precision={0}
                                      value={item.count}
                                      onChange={(value) => {
                                        handleOptionChange(index, Number(value))
                                      }}
                                    />
                                    <div onClick={() => {
                                      handleCountChange("increase", index, item)
                                    }}
                                         className={`w-[35px] border-x border-solid !border-[#313D4F] text-white font-semibold ${item.disabled ? "cursor-not-allowed" : "cursor-pointer"}`}>+
                                    </div>
                                    <div className="w-[40px] text-center">%</div>
                                  </div>
                                </div>
                              </div>
                            )
                          })
                        }
                          <button onClick={startVoting} className="w-full h-[40px] bg-sky-500 hover:bg-sky-700 text-white py-2 px-6 rounded-full" type="submit" disabled={loading}>
                              Vote
                          </button>
                      </div>
                  </div>
              </div>
          }
        </div>
      </div>
      <div className="w-full lg:w-4/12 lg:min-w-[321px]">
        <div className="mt-4 space-y-4 lg:mt-0">
          <div
            className="text-base border-solid border-y border-skin-border bg-skin-block-bg md:rounded-xl md:border">
            <div
              className="group flex h-[57px] justify-between rounded-t-none border-b border-skin-border border-solid px-4 pb-[12px] pt-3 md:rounded-t-lg">
              <h4 className="flex items-center text-xl">
                <div>Message</div>
              </h4>
            </div>
            <div className="p-4 leading-6 sm:leading-8">
              <div className="space-y-1">
                <div>
                  <b>Vote Type</b>
                  <span
                    className="float-right text-white">{["Single", "Multiple"][votingData?.voteType - 1]} Choice Voting</span>
                </div>
                <div>
                  <b>Exp. Time</b>
                  <span className="float-right text-white">
                    {dayjs(votingData?.showTime).format('MMM.D, YYYY, h:mm A')}
                  </span>
                </div>
                <div>
                  <b>Exp. Timezone</b>
                  <span className="float-right text-white">
                    {votingData?.GMTOffset}
                  </span>
                </div>
                <div>
                  <b>Snapshot</b>
                  <span className="float-right text-white">Takes At Exp. Time</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Vote;
