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

import { useChainModal, useConnectModal } from "@rainbow-me/rainbowkit";
import axios from 'axios';
import dayjs from 'dayjs';
import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link, useLocation, useParams } from 'react-router-dom';
import Loading from 'src/components/Loading';
import VoteStatusBtn from 'src/components/VoteStatusBtn';
import { useAccount } from "wagmi";
import {
  COMPLETED_STATUS,
  PASSED_STATUS,
  REJECTED_STATUS,
  VOTE_COUNTING_STATUS,
  VOTE_OPTIONS,
  WRONG_NET_STATUS,
  proposalHistoryApi,
  proposalResultApi,
  web3AvatarUrl,
} from "../../common/consts";
import { useCurrentTimezone } from "../../common/store";
import type { ProposalHistory, ProposalOption, ProposalResult } from "../../common/types";
import EllipsisMiddle from "../../components/EllipsisMiddle";
import MDEditor from '../../components/MDEditor';
import VoteList from "../../components/VoteList";
const VotingResults = () => {
  const { chain, isConnected } = useAccount();
  const { openConnectModal } = useConnectModal();
  const { openChainModal } = useChainModal();
  const { id, cid } = useParams();
  const { state } = useLocation() || null;
  const { t } = useTranslation();
  const [votingData, setVotingData] = useState(state);
  const timezone = useCurrentTimezone((state: any) => state.timezone);
  const [loading, setLoading] = useState(true);

  const initState = async () => {
    const option: ProposalOption[] = [];
    let voteList: any[] = [];
    let voteStatus = null;
    let subStatus = 0;
    const params = {
      proposalId: Number(id),
      network: chain?.id
    }

    // Fetch proposal data from IPFS
    const { data: proposalData } = await axios.get(`https://${cid}.ipfs.w3s.link/`);
    // Check if the proposal chain ID matches the current chain ID
    if (proposalData.chainId !== chain?.id) {
      // If not, set vote status to wrong network status
      voteStatus = WRONG_NET_STATUS;
      if (isConnected) {
        openChainModal && openChainModal();
      } else {
        openConnectModal && openConnectModal();
      }
    } else {
      // If proposal chain ID matches, proceed with fetching voting data
      const { data: { data: resultData } } = await axios.get(proposalResultApi, {
        params,
      })
      // Determine vote status based on whether votes have been counted
      voteStatus = resultData.length > 0 ? COMPLETED_STATUS : VOTE_COUNTING_STATUS;

      // Map result data to populate option array
      resultData.map((_: any, index: number) => {
        const voteItem = resultData.find((vote: ProposalResult) => vote.optionId === index);
        option.push({
          name: proposalData.option[voteItem.optionId],
          count: voteItem?.votes ? Number(voteItem.votes) : 0
        })
      })
      if (voteStatus == COMPLETED_STATUS) {//
        const passedOption = option?.find((v: any) => { return v.name === VOTE_OPTIONS[0] })
        const rejectOption = option?.find((v: any) => { return v.name === VOTE_OPTIONS[1] })
        if ((passedOption?.count ?? 0) > (rejectOption?.count ?? 0)) {
          subStatus = PASSED_STATUS
        } else {
          subStatus = REJECTED_STATUS
        }
      }
      // Fetch voting history data
      const { data: { data: historyData } } = await axios.get(proposalHistoryApi, {
        params,
      });
      // Map history data to populate voteList array
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
    let powerBlockHeight = 0
    if (voteList?.length) {
      powerBlockHeight = voteList[0].powerBlockHeight
    }
    // Set voting data state
    setVotingData({
      ...proposalData,
      id,
      cid,
      option,
      voteStatus,
      subStatus,
      powerBlockHeight,
      // Sort voteList array by number of votes in descending order
      voteList: voteList?.sort((a: any, b: any) => b.votes - a.votes)
    })
    setLoading(false)
  }

  useEffect(() => {
    initState();
  }, [chain]);

  // const handleVoteStatusTag = (status: number) => {
  //   switch (status) {
  //     case WRONG_NET_STATUS:
  //       return {
  //         name: 'Wrong network',
  //         color: 'bg-red-700',
  //       };
  //     case VOTE_COUNTING_STATUS:
  //       return {
  //         name: 'Vote Counting',
  //         color: 'bg-yellow-700',
  //       };
  //     case COMPLETED_STATUS:
  //       return {
  //         name: 'Completed',
  //         color: 'bg-[#6D28D9]',
  //       };
  //     default:
  //       return {
  //         name: '',
  //         color: '',
  //       };
  //   }
  // }

  let href = '';
  let img = '';
  if (votingData?.githubName) {
    href = `https://github.com/${votingData.githubName}`;
    img = `${votingData?.githubAvatar}`;
  } else {
    href = `${chain?.blockExplorers?.default.url}/address/${votingData?.address}`;
    img = `${web3AvatarUrl}:${votingData?.address}`
  }
  if (loading) {
    return (
      <Loading />
    );
  }
  return (
    <div className='flex voting-result'>
      <div className='relative w-full pr-4 lg:w-8/12'>
        <div className='px-3 mb-6 md:px-0'>
          <button>
            <div className='inline-flex items-center gap-1 text-skin-text hover:text-skin-link'>
              <Link to='/home' className='flex items-center'>
                <svg className='mr-1' viewBox="0 0 24 24" width="1.2em" height="1.2em"><path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m11 17l-5-5m0 0l5-5m-5 5h12"></path></svg>
                {t('content.back')}
              </Link>
            </div>
          </button>
        </div>
        <div className='px-3 md:px-0'>
          <h1 className='mb-6 text-2xl font-semibold text-[#313D4F] break-words break-all leading-12'>
            {votingData?.name}
          </h1>
          {
            (votingData?.voteStatus || votingData?.voteStatus === 0) &&
            <div className="flex justify-between mb-6">
              <div className="flex items-center w-full mb-1 sm:mb-0">
                <VoteStatusBtn status={(votingData?.subStatus > 0) ? votingData?.subStatus : votingData?.voteStatus} />
                <div className="flex items-center justify-center ml-[12px]">
                  <div className='text-[#4B535B] text-[14px]'>{t('content.createdby')}</div>
                  <div className='p-[5px] ml-[8px] flex items-center justify-center bg-[#F5F5F5] rounded-full'>
                    <img className="w-[20px] h-[20px] rounded-full mr-[4px]" src={img} alt="" />
                    <a
                      className="text-[#313D4F]"
                      target="_blank"
                      rel="noreferrer"
                      href={href}
                    >
                      {votingData?.githubName || EllipsisMiddle({ suffixCount: 4, children: votingData?.address })}
                    </a>
                  </div>
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
              onChange={() => { }}
            />
          </div>
          {
            votingData?.voteStatus === COMPLETED_STATUS && <VoteList voteList={votingData?.voteList} chain={chain} />
          }
        </div>
      </div>
      <div className='w-full lg:w-4/12 lg:min-w-[321px]'>
        <div className='mt-4 space-y-4 lg:mt-0'>
          <div className='text-base border-solid border-[1px] border-[#DFDFDF] border-y  bg-skin-block-bg md:rounded-xl md:border'>
            <div className='group flex h-[57px] justify-between px-4 pb-[12px] pt-3 md:rounded-t-lg'>
              <h4 className='flex items-center font-medium text-[#313D4F]'>
                <div>{t('content.details')}</div>
              </h4>
              <div className='flex items-center'>
              </div>
            </div>
            <div className='h-[1px] bg-[#DFDFDF]' />
            <div className='p-4 leading-6 sm:leading-8'>
              <div className='space-y-1 text-sm font-medium'>
                <div className='flex justify-between'>
                  <div>{t('content.startTime')}</div>
                  <span className='text-[#313D4F] text-sm font-normal'>{votingData?.startTime && dayjs(votingData.startTime * 1000).format('MMM.D, YYYY, h:mm A')}</span>
                </div>
                <div className='flex justify-between'>
                  <div>{t('content.endTime')}</div>
                  <span className='text-[#313D4F] text-sm font-normal'>{votingData?.expTime && dayjs(votingData.expTime * 1000).format('MMM.D, YYYY, h:mm A')}</span>
                </div>
                <div className='flex justify-between'>
                  <div>{t('content.timezone')}</div>
                  <span className='text-[#313D4F] text-sm font-normal'>{timezone}</span>
                </div>
                {
                  votingData?.powerBlockHeight > 0 && (<div className='flex justify-between'>
                    <div className='text-sm font-medium'>{t('content.blockHeight')}</div>
                    <span className='text-[#313D4F] font-normal'>{votingData?.powerBlockHeight}</span>
                  </div>)
                }
              </div>

            </div>
          </div>
          {
            votingData?.voteStatus === COMPLETED_STATUS &&
            <div className='text-base border-solid border-[1px] border-[#DFDFDF] border-y  bg-skin-block-bg md:rounded-xl md:border'>
              <div className='group flex h-[57px] justify-between px-4 pb-[12px] pt-3 md:rounded-t-lg'>
                <h4 className='flex items-center font-medium'>
                  <div>{t('content.results')}</div>
                </h4>
                <div className='flex items-center' />
              </div>
              <div className='h-[1px] bg-[#DFDFDF]' />
              <div className='p-4 leading-6 sm:leading-8'>
                <div className='space-y-3'>
                  {
                    votingData?.option?.map((item: any, index: number) => {
                      return (
                        <div key={item.name + index}>
                          <div className='flex justify-between mb-1 text-skin-link'>
                            <div className='w-[150px] flex items-center overflow-hidden'>
                              <span className='mr-1 truncate text-sm'>{item.name}</span>
                            </div>
                            <div className='flex justify-end'>
                              <div className='space-x-2 text-sm font-medium'>
                                <span>{item.count}%</span>
                              </div>
                            </div>
                          </div>
                          <div className='relative h-2 rounded bg-[#EEEEEE]'>
                            {
                              item.count ?
                                <div
                                  className='absolute top-0 left-0 h-full rounded bg-[#0190FF]'
                                  style={{
                                    width: `${item.count}%`
                                  }}
                                /> :
                                <div
                                  className='absolute top-0 left-0 h-full rounded bg-[#EEEEEE]'
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
