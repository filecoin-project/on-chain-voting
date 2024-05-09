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
import { message } from "antd";
import axios from 'axios';
import dayjs from 'dayjs';
import {getWeb3IpfsId, useDynamicContract} from "../../hooks";
import { useAccount } from "wagmi";
import { useConnectModal, useChainModal } from "@rainbow-me/rainbowkit";
import MDEditor from "../../components/MDEditor";
import EllipsisMiddle from "../../components/EllipsisMiddle";
import LoadingButton from "../../components/LoadingButton";
import {
  IN_PROGRESS_STATUS,
  WRONG_NET_STATUS,
  web3AvatarUrl,
  VOTE_SUCCESS_MSG,
  CHOOSE_VOTE_MSG, PENDING_STATUS,
} from "../../common/consts";
import { timelockEncrypt, roundAt, mainnetClient, Buffer } from "tlock-js";
import {ProposalList, ProposalOption} from "../../common/types";
import "./index.less";

const Vote = () => {
  const { chain, isConnected } = useAccount();
  const chainId = chain?.id || 0;
  const { id, cid } = useParams();
  const [votingData, setVotingData] = useState({} as ProposalList);
  const { openConnectModal } = useConnectModal();
  const { openChainModal } = useChainModal();

  const navigate = useNavigate();
  const [options, setOptions] = useState([] as ProposalOption[]);
  const [selectedOptionIndex, setSelectedOptionIndex] = useState(-1);

  const [loading, setLoading] = useState(false);

  useEffect(() => {
    initState();
  }, [chain]);

  const initState = async () => {
    const res = await axios.get(`https://${cid}.ipfs.w3s.link/`);
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
      voteStatus = dayjs().unix() < data?.startTime ? PENDING_STATUS : IN_PROGRESS_STATUS;
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
    if (selectedOptionIndex < 0) {
      message.warning(CHOOSE_VOTE_MSG);
    } else {
      setLoading(true);

      // vote params
      const encryptValue = await handleEncrypt([[`${selectedOptionIndex}`, `100`]]);
      const optionId = await getWeb3IpfsId(encryptValue);
      if (isConnected) {
        const { voteApi } = useDynamicContract(chainId);
        const res = await voteApi(Number(id), optionId);
        if (res.code === 200 && res.data?.hash) {
          message.success(VOTE_SUCCESS_MSG, 3);
          setTimeout(() => {
            navigate("/");
          }, 3000);
        } else {
          message.error(res.msg, 3);
        }
        setLoading(false);
      } else {
        // @ts-ignore
        openConnectModal && openConnectModal();
      }
    }
  }

  const handleOptionClick = (index: number) => {
    setSelectedOptionIndex(index);
  }

  const handleVoteStatusTag = (status: number) => {
    switch (status) {
      case WRONG_NET_STATUS:
        return {
          name: 'Wrong network',
          color: 'bg-red-700',
        };
      case PENDING_STATUS:
        return {
          name: 'Pending',
          color: 'bg-cyan-700',
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
                  <div className="border-[#313D4F] mt-6 border-skin-border bg-skin-block-bg text-base md:rounded-xl md:border border-solid">
                      <div className="group flex h-[57px] !border-[#313D4F] justify-between items-center border-b px-4 pb-[12px] pt-3 border-solid">
                          <h4 className="text-xl">
                              Cast Your Vote
                          </h4>
                      </div>
                      <div className="p-4 text-center">
                        {
                          options.map((item: ProposalOption, index: number) => {
                            return (
                              <div className="mb-4 space-y-3 leading-10" key={item.name + index} onClick={() => { handleOptionClick(index) }}>
                                <div
                                  className={`w-full h-[45px] border-[#313D4F] ${selectedOptionIndex === index ? 'border-blue-500' : ''} hover:border-blue-500 flex justify-between items-center pl-8 pr-4 md:border border-solid rounded-full cursor-pointer`}>
                                  <div className="text-ellipsis h-[100%] overflow-hidden">{item.name}</div>
                                  {
                                    selectedOptionIndex === index &&
                                    <svg viewBox="0 0 24 24" width="1.2em" height="1.2em" className="-ml-1 mr-2 text-md">
                                      <path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round"
                                            strokeWidth="2" d="m5 13l4 4L19 7" />
                                    </svg>
                                  }
                                </div>
                              </div>
                            )
                          })
                        }
                          <LoadingButton text='Vote' isFull={true} loading={loading} handleClick={startVoting} />
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
                  <b>Start Time</b>
                  <span className="float-right text-white">
                    {votingData?.showTime?.length && dayjs(votingData.showTime[0]).format('MMM.D, YYYY, h:mm A')}
                  </span>
                </div>
                <div>
                  <b>End Time</b>
                  <span className="float-right text-white">
                    {votingData?.showTime?.length && dayjs(votingData.showTime[1]).format('MMM.D, YYYY, h:mm A')}
                  </span>
                </div>
                <div>
                  <b>Timezone</b>
                  <span className="float-right text-white">
                    {votingData?.GMTOffset}
                  </span>
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
