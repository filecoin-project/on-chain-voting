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

import React, { useState, useEffect } from 'react';
import { Link, useLocation, useParams } from 'react-router-dom';
import { useAccount } from "wagmi";
import { useChainModal, useConnectModal } from "@rainbow-me/rainbowkit";
import axios from 'axios';
import dayjs from 'dayjs';
import EllipsisMiddle from "../../components/EllipsisMiddle";
import MDEditor from '../../components/MDEditor';
import {
  web3AvatarUrl,
  COMPLETED_STATUS,
  WRONG_NET_STATUS, VOTE_COUNTING_STATUS,
} from "../../common/consts";
import VoteList from "../../components/VoteList";
import {ProposalOption, ProposalResult, ProposalHistory} from "../../common/types";

const VotingResults = () => {
  const { chain, isConnected } = useAccount();
  const { openConnectModal } = useConnectModal();
  const { openChainModal } = useChainModal();
  const { id, cid } = useParams();
  const { state } = useLocation() || null;
  const [votingData, setVotingData] = useState(state);

  useEffect(() => {
    initState();
  }, [chain])

  const initState = async () => {
    const option: ProposalOption[] = [];
    let voteList: any[] = [];
    let voteStatus = null;
    const params = {
      proposalId: Number(id),
      network: chain?.id
    }

    const { data: proposalData } = await axios.get(`https://${cid}.ipfs.w3s.link/`);
    if (proposalData.chainId !== chain?.id) {
      voteStatus = WRONG_NET_STATUS;
      if (isConnected) {
        openChainModal && openChainModal();
      } else {
        openConnectModal && openConnectModal();
      }
    } else {
      const { data: { data: resultData } } = await axios.get('/api/proposal/result', {
        params,
      })
      voteStatus = resultData.length > 0 ? COMPLETED_STATUS : VOTE_COUNTING_STATUS;
      resultData.map((_: any, index: number) => {
        const voteItem = resultData.find((vote: ProposalResult) => vote.optionId === index);
        option.push({
          name: proposalData.option[voteItem.optionId],
          count: voteItem?.votes ? Number(voteItem.votes) : 0
        })
      })
      const { data: { data: historyData } } = await axios.get('/api/proposal/history', {
        params,
      });
      voteList = historyData?.votePowers?.map((item: ProposalHistory) => ({
        ...item,
        optionName: proposalData.option[item.optionId],
        address: item.address?.substring(0, 42),
        totalClientPower: historyData.totalClientPower,
        totalDeveloperPower: historyData.totalDeveloperPower,
        totalSpPower: historyData.totalSpPower,
        totalTokenHolderPower: historyData.totalTokenHolderPower,
        votePowers: historyData.votePowers
      }));
    }
    setVotingData({
      ...proposalData,
      id,
      cid,
      option,
      voteStatus,
      voteList: voteList?.sort((a: any, b: any) => b.votes - a.votes)
    })
  }

  const handleVoteStatusTag = (status: number) => {
    switch (status) {
      case WRONG_NET_STATUS:
        return {
          name: 'Wrong network',
          color: 'bg-red-700',
        };
      case VOTE_COUNTING_STATUS:
        return {
          name: 'Vote Counting',
          color: 'bg-yellow-700',
        };
      case COMPLETED_STATUS:
        return {
          name: 'Completed',
          color: 'bg-[#6D28D9]',
        };
      default:
        return {
          name: '',
          color: '',
        };
    }
  }

  return (
    <div className='flex voting-result'>
      <div className='relative w-full pr-4 lg:w-8/12'>
        <div className='px-3 mb-6 md:px-0'>
          <button>
            <div className='inline-flex items-center gap-1 text-skin-text hover:text-skin-link'>
              <Link to='/' className='flex items-center'>
                <svg className='mr-1' viewBox="0 0 24 24" width="1.2em" height="1.2em"><path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m11 17l-5-5m0 0l5-5m-5 5h12"></path></svg>
                Back
              </Link>
            </div>
          </button>
        </div>
        <div className='px-3 md:px-0'>
          <h1 className='mb-6 text-3xl font-semibold text-white break-words break-all leading-12'>
            {votingData?.name}
          </h1>
          {
            (votingData?.voteStatus || votingData?.voteStatus === 0) &&
            <div className="flex justify-between mb-6">
              <div className="flex items-center justify-between w-full mb-1 sm:mb-0">
                <button
                  className={`${handleVoteStatusTag(votingData?.voteStatus).color} bg-[#6D28D9] h-[26px] px-[12px] text-white rounded-xl mr-4`}>
                  {handleVoteStatusTag(votingData?.voteStatus).name}
                </button>
                <div className="flex items-center justify-center">
                  <img className="w-[20px] h-[20px] rounded-full mr-2" src={`${web3AvatarUrl}:${votingData?.address}`} alt="" />
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
          <div className='MDEditor mb-8'>
            <MDEditor
              className="border-none rounded-[16px] bg-transparent"
              style={{ height: 'auto' }}
              value={votingData?.descriptions}
              moreButton
              readOnly={true}
              view={{ menu: false, md: false, html: true, both: false, fullScreen: true, hideMenu: false }}
              onChange={() => {}}
            />
          </div>
          {
            votingData?.voteStatus === COMPLETED_STATUS && <VoteList voteList={votingData?.voteList} chain={chain} />
          }
        </div>
      </div>
      <div className='w-full lg:w-4/12 lg:min-w-[321px]'>
        <div className='mt-4 space-y-4 lg:mt-0'>
          <div className='text-base border-solid border-y border-skin-border bg-skin-block-bg md:rounded-xl md:border'>
            <div className='group flex h-[57px] justify-between rounded-t-none border-b border-skin-border border-solid px-4 pb-[12px] pt-3 md:rounded-t-lg'>
              <h4 className='flex items-center text-xl'>
                <div>Message</div>
              </h4>
              <div className='flex items-center'>

              </div>
            </div>
            <div className='p-4 leading-6 sm:leading-8'>
              <div className='space-y-1'>
                <div>
                  <b>Start Time</b>
                  <span className='float-right text-white'>{votingData?.showTime?.length && dayjs(votingData?.showTime[0]).format('MMM.D, YYYY, h:mm A')}</span>
                </div>
                <div>
                  <b>End Time</b>
                  <span className='float-right text-white'>{votingData?.showTime?.length && dayjs(votingData?.showTime[1]).format('MMM.D, YYYY, h:mm A')}</span>
                </div>
                <div>
                  <b>Timezone</b>
                  <span className='float-right text-white'>{votingData?.GMTOffset}</span>
                </div>
              </div>
            </div>
          </div>
          {
            votingData?.voteStatus === COMPLETED_STATUS &&
              <div className='text-base border-solid border-y border-skin-border bg-skin-block-bg md:rounded-xl md:border'>
                  <div className='group flex h-[57px] justify-between rounded-t-none border-b border-skin-border border-solid px-4 pb-[12px] pt-3 md:rounded-t-lg'>
                      <h4 className='flex items-center text-xl'>
                          <div>Results</div>
                      </h4>
                      <div className='flex items-center' />
                  </div>
                  <div className='p-4 leading-6 sm:leading-8'>
                      <div className='space-y-3'>
                        {
                          votingData?.option?.map((item: any, index: number) => {
                            return (
                              <div key={item.name + index}>
                                <div className='flex justify-between mb-1 text-skin-link'>
                                  <div className='w-[150px] flex items-center overflow-hidden'>
                                    <span className='mr-1 truncate'>{item.name}</span>
                                  </div>
                                  <div className='flex justify-end'>
                                    <div className='space-x-2'>
                                      <span>{item.count}%</span>
                                    </div>
                                  </div>
                                </div>
                                <div className='relative h-2 rounded bg-[#273141]'>
                                  {
                                    item.count ?
                                      <div
                                        className='absolute top-0 left-0 h-full rounded bg-[#384AFF]'
                                        style={{
                                          width: `${item.count}%`
                                        }}
                                      /> :
                                      <div
                                        className='absolute top-0 left-0 h-full rounded bg-[#273141]'
                                        style={{
                                          width: '100%'
                                        }}
                                      />
                                  }
                                </div>
                              </div>
                            )
                          })
                        }
                      </div>
                  </div>
              </div>
          }
        </div>
      </div>
    </div>
  )
}

export default VotingResults
