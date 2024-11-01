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
import { Empty, Pagination, Row, message } from "antd";
import axios from "axios";
import dayjs from 'dayjs';
import timezone from 'dayjs/plugin/timezone';
import utc from 'dayjs/plugin/utc';
import React, { useEffect, useState } from "react";
import { useTranslation } from 'react-i18next';
import { useNavigate } from "react-router-dom";
import VoteStatusBtn from "src/components/VoteStatusBtn";
import { useAccount, useWaitForTransactionReceipt } from "wagmi";
import {
  IN_PROGRESS_STATUS, calibrationChainId,
  PENDING_STATUS,
  STORING_DATA_MSG,
  STORING_FAILED_MSG,
  STORING_STATUS,
  STORING_SUCCESS_MSG,
  VOTE_ALL_STATUS,
  VOTE_FILTER_LIST,
  web3AvatarUrl
} from "../../common/consts";
import { useCheckFipEditorAddress } from "../../common/hooks"
import { useCurrentTimezone, useProposalStatus, useStoringCid, useVotingList } from "../../common/store";
import type { ProposalResult, VotingList } from '../../common/types';
import EllipsisMiddle from "../../components/EllipsisMiddle";
import ListFilter from "../../components/ListFilter";
import Loading from "../../components/Loading";
import { markdownToText } from "../../utils";
dayjs.extend(utc);
dayjs.extend(timezone);

const Home = () => {
  const navigate = useNavigate();
  const { t, i18n } = useTranslation();
  const { chain, address, isConnected } = useAccount();
  const chainId = chain?.id || calibrationChainId;

  const { openConnectModal } = useConnectModal();

  const items = VOTE_FILTER_LIST.map((item) => {
    return {
      label: t(item.label),
      value: item.value
    }
  });

  const [proposalStatus, setProposalStatus] = useState(VOTE_ALL_STATUS);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [messageApi, contextHolder] = message.useMessage();
  const timezone = useCurrentTimezone((state: any) => state.timezone);
  const storingCid = useStoringCid((state: any) => state.storingCid);
  const setStoringCid = useStoringCid((state: any) => state.setStoringCid);
  const { votingList, totalPage, searchKey } = useVotingList((state: any) => state.votingData);
  const setVotingList = useVotingList((state: any) => state.setVotingList);
  const setStatusList = useProposalStatus((state: any) => state.setStatusList);
  const { isFipEditorAddress } = useCheckFipEditorAddress(chainId, address);

  const { isFetched, isSuccess, isError } = useWaitForTransactionReceipt({
    hash: storingCid[0]?.hash
  });

  useEffect(() => {
    if (isConnected && !loading) {
      queryVotingList(page, proposalStatus);
    }
  }, [chain, page, address, isConnected, i18n.language]);

  useEffect(() => {
    if (isFetched) {
      // If data is fetched, remove the last item from storingCid array
      storingCid.splice(storingCid.length - 1, 1);
      setStoringCid(storingCid);
      // If the transaction is successful, show a success message
      if (isSuccess) {
        messageApi.open({
          type: 'success',
          content: t(STORING_SUCCESS_MSG)
        });
        setTimeout(() => {
          queryVotingList(page, proposalStatus);
        }, 3000)
      }
      // If the transaction fails, show an error message
      if (isError) {
        messageApi.open({
          type: 'error',
          content: t(STORING_FAILED_MSG)
        })
      }
    }
  }, [isFetched]);

  const queryVotingList = async (page: number, proposalStatus: number) => {
    setLoading(true);
    const params = {
      chainId,
      page,
      pageSize: 5,
      searchKey: searchKey,
      status: proposalStatus === VOTE_ALL_STATUS ? 0 : proposalStatus,
    }
    const { data: { data: votingData } } = await axios.get('/api/proposal/list', { params });
    setVotingList({ votingList: votingData?.list || [], totalPage: votingData?.total, searchKey: searchKey });
    setLoading(false);
  }
  /**
   * filter proposal list
   * @param status
   */
  const handleFilter = async (status: number) => {
    setProposalStatus(status);
    queryVotingList(1, status);
    setStatusList(status)
    setPage(1);
  }
  useEffect(() => {
    //When the search value changes, display all by default
    setPage(1);
  }, [searchKey])
  /**
   * page jump
   * @param item
   */
  const handleJump = (item: VotingList) => {
    if (item.status === STORING_STATUS) {
      messageApi.open({
        type: 'warning',
        content: t(STORING_DATA_MSG),
      });
      return;
    }
    const router = `/${[PENDING_STATUS, IN_PROGRESS_STATUS].includes(item.status) ? "vote" : "votingResults"}/${item.proposalId}/${item.cid}`;
    navigate(router, { state: item });
  }

  const handleCreate = () => {
    if (!isConnected) {
      openConnectModal && openConnectModal();
      return false;
    }
    navigate("/createVote");
  }

  const handlePageChange = async (page: number) => {
    // Reset vote status when page change
    // setProposalStatus(VOTE_ALL_STATUS);
    setPage(page);
    queryVotingList(page, proposalStatus);
  }

  /**
   * render proposal list
   * @param list
   */
  const renderList = (list: VotingList[]) => {
    // if (proposalStatus !== VOTE_ALL_STATUS) {
    //   list = list.filter(item => item.status === proposalStatus);
    // }
    if (!list.length) {
      return (
        <div className='empty mt-20'>
          <Empty
            description={
              <span className='text-black'>{t('content.noData')}</span>
            }
          />
        </div>
      );
    }
    return list.map((item: VotingList, index: number) => {
      const maxOption = (item?.voteResult || [])?.reduce((prev, current) => {
        return (prev.votes > current.votes) ? prev : current;
      }, 0);

      let href = '';
      let img = '';
      if (item?.githubName) {
        href = `https://github.com/${item.githubName}`;
        img = `${item.githubAvatar}`;
      } else {
        href = `${chain?.blockExplorers?.default.url}/address/${item.address}`;
        img = `${web3AvatarUrl}:${item.address}`
      }
      return (
        <div
          key={item.cid + index}
          className="rounded-xl border-[1px] border-solid border-[#DFDFDF] bg-[#FFFFFF] px-[30px] py-[12px] mb-[16px]"
        >
          <div className="flex justify-between mb-3">
            <div
              className="flex justify-center items-center"
            >
              <a
                target='_blank'
                rel="noopener noreferrer"
                href={href}
              >
                <div className="bg-[#F5F5F5] rounded-full  flex p-[5px] justify-center items-center">
                  <img className="w-[20px] h-[20px] rounded-full mr-2" src={img} alt="" />
                  <div className="truncate text-[#313D4F] mr-[5px]">
                    {item.githubName || EllipsisMiddle({ suffixCount: 4, children: item.address })}
                  </div>
                </div>
              </a>
              <div className="truncate text-[#4B535B] text-sm ml-5">
                {t('content.created')} {dayjs(item.currentTime * 1000).format('YYYY-MM-D')}
              </div>
            </div>
            <VoteStatusBtn status={item.status} />

          </div>
          <div className="relative mb-4 line-clamp-2 break-words break-all text-lg pr-[80px] leading-7 cursor-pointer"
            onClick={() => {
              handleJump(item);
            }}>
            <h3 className="inline pr-2 text-2xl font-semibold text-[#313D4F]">
              {item.name}
            </h3>
          </div>
          <div className="mb-2 line-clamp-2 break-words text-normal text-lg cursor-pointer" onClick={() => {
            handleJump(item)
          }}>
            {markdownToText(item.descriptions)}
          </div>
          {
            maxOption?.votes > 0 &&
            <div>
              {
                item.voteResult?.map((option: ProposalResult, index: number) => {
                  const isapprove = option.optionId == 0; //0 approve 1 reject
                  const passed = maxOption.optionId == 0;
                  let bgColor = "#F7F7F7";
                  let txColor = "#273141";
                  let borderColor = "#F7F7F7";
                  if (isapprove && passed) {
                    bgColor = "#E3FFEE";
                    txColor = "#006227";
                    borderColor = "#87FFBE";
                  } else if (!isapprove && !passed) {
                    bgColor = "#FFF3F3";
                    txColor = "#AA0101";
                    borderColor = "#FFDBDB";
                  }
                  return (
                    <div className="h-[35px] relative mt-1 w-full" key={index}>
                      <div
                        style={{ color: txColor }}
                        className='absolute ml-3 flex items-center leading-[35px] font-semibold'>
                        {
                          ((maxOption.votes === 50 && option.optionId === 1) || (maxOption.votes > 50 && option.votes > 0 && option.votes === maxOption.votes)) &&
                          <svg viewBox="0 0 24 24" width="1.2em" height="1.2em" className="-ml-1 mr-2 text-sm">
                            <path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round"
                              strokeWidth="2" d="m5 13l4 4L19 7" />
                          </svg>
                        }
                        {option.optionId === 0 ? t('content.approve') : t('content.rejected')}
                      </div>
                      <div className="font-semibold absolute right-0 mr-3 leading-[35px]" style={{ color: txColor }}>{option.votes}%</div>
                      {option.votes > 0 && <div className="h-[35px] border-[1px] border-solid rounded-md bg-[#E3FFEE]" style={{ width: `${option.votes}%`, backgroundColor: bgColor, borderColor: borderColor }} />
                      }
                    </div>
                  )
                })
              }
            </div>
          }
          <div className="text-[#4B535B] text-sm mt-4">
            <span className="mr-2">{t('content.endTime')}:</span>
            {dayjs(item.expTime * 1000).format('MMM.D, YYYY, h:mm A')} ({timezone})
          </div>
        </div >
      )
    })
  }

  const renderContent = () => {
    // Display loading when data is loading
    if (loading) {
      return (
        <Loading />
      );
    }

    // Display empty when data is empty
    if (!votingList.length) {
      return (
        <div className='empty mt-20'>
          <Empty
            description={
              <span className='text-black'>{t('content.noData')}</span>
            }
          />
        </div>
      );
    }

    return (
      <div className='home-table overflow-auto'>
        {
          renderList(votingList)
        }
        <Row justify='end'>
          <Pagination
            simple
            showSizeChanger={false}
            current={page}
            pageSize={5}
            total={totalPage}
            onChange={handlePageChange}
          />
        </Row>
      </div>
    );
  };
  return (
    <div className="home_container main">
      {contextHolder}
      <div className="flex justify-between items-center rounded-xl border-[1px] border-solid border-[#DFDFDF] bg-[#ffffff] mb-[32px] px-[12px]">
        <ListFilter
          name="Status"
          value={proposalStatus}
          list={items}
          onChange={handleFilter}
        />
        {
          !!isFipEditorAddress &&
          <button
            className="h-[40px] bg-sky-500 hover:bg-sky-700 text-white py-2 px-4 rounded-xl"
            onClick={handleCreate}
          >
            {t('content.createProposal')}
          </button>
        }
      </div>
      {
        renderContent()
      }
    </div>
  )
}

export default Home;
