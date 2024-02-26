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

import React, {useEffect, useState} from "react";
import { ConfigProvider, theme, Row, Empty, Pagination, Spin } from "antd";
import {useNetwork, useAccount} from "wagmi";
import {useConnectModal} from "@rainbow-me/rainbowkit";
import {useNavigate} from "react-router-dom";
import axios from "axios";
import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
import {
  VOTE_ALL_STATUS,
  VOTE_FILTER_LIST,
  VOTE_LIST,
  IN_PROGRESS_STATUS,
  VOTE_COUNTING_STATUS,
  COMPLETED_STATUS,
  web3AvatarUrl,
} from '../../common/consts';
import ListFilter from "../../components/ListFilter";
import EllipsisMiddle from "../../components/EllipsisMiddle";
import {useStaticContract} from "../../hooks";
import {ProposalData, ProposalFilter, ProposalList, ProposalOption, ProposalResult} from '../../common/types';
import Loading from "../../components/Loading";

dayjs.extend(utc);
dayjs.extend(timezone);

const Home = () => {
  const navigate = useNavigate();
  const {chain} = useNetwork();

  const {isConnected} = useAccount();
  const {openConnectModal} = useConnectModal();

  const [filterList, setFilterList] = useState([
    {
      label: "All",
      value: VOTE_ALL_STATUS
    }
  ])

  const [loading, setLoading] = useState(true);
  const [proposalStatus, setProposalStatus] = useState(VOTE_ALL_STATUS);
  const [proposalList, setProposalList] = useState<ProposalList[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(5);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    getProposalList(page);
  }, [chain, page]);

  /**
   * get proposal list
   * @param page
   */
  const getProposalList = async (page: number) => {
    setLoading(true);
    if (!isConnected) {
      openConnectModal && openConnectModal();
      setLoading(false);
      return false;
    }
    const chainId = chain?.id || 0;
    const { getLatestId, getProposal } = await useStaticContract(chainId);
    const res = await getLatestId();
    const total = res.data.toNumber();
    setTotal(total);
    const offset = (page - 1) * pageSize;

    const proposalRequests = [];

    for (let i = total - offset; i > Math.max(total - offset - pageSize, 0); i--) {
      proposalRequests.push(getProposal(i));
    }

    const proposals = await Promise.all(proposalRequests);

    const list = proposals.map(async ({ data }, index) => {
      const params = {
        proposalId: total - offset - index,
        network: chainId
      };

      const { data: { data: resultData } } = await axios.get('/api/proposal/result', { params });
      const proposalResults = resultData.map((item: ProposalResult) => ({
        optionId: item.optionId,
        votes: item.votes
      }));
      return {
        id: total - offset - index,
        cid: data.cid,
        creator: data.creator,
        expTime: data.expTime.toNumber(),
        proposalType: data.proposalType.toNumber(),
        proposalResults
      };
    });

    const proposalsList: ProposalData[] = await Promise.all(list);
    const originList: ProposalList[] = await getList(proposalsList) || [];

    setLoading(false);
    setFilterList(VOTE_FILTER_LIST);
    setProposalList(originList);
  }

  /**
   * get proposal info
   * @param proposals
   */
  const getList = async (proposals: ProposalData[]) => {
    const ipfsUrls = proposals.map(
      (_item: ProposalData) => `https://${_item.cid}.ipfs.nftstorage.link/`
    );
    try {
      const responses = await Promise.all(ipfsUrls.map((url: string) => axios.get(url)));
      const results: ProposalList[] = responses.map((res, i: number) => {
        const  proposal = proposals[i];
        const now = dayjs().unix();
        let proposalStatus = 0;
        if (now >= proposal.expTime) {
          if (proposal.proposalResults.length === 0) {
            proposalStatus = VOTE_COUNTING_STATUS
          } else {
            proposalStatus = COMPLETED_STATUS
          }
        } else {
          proposalStatus = IN_PROGRESS_STATUS
        }
        const option = res.data.option?.map((item: string, index: number) => {
          const proposalItem = proposal?.proposalResults?.find(
            (proposal: ProposalResult) => proposal.optionId === index
          );
          return {
            name: item,
            count: proposalItem?.votes ? Number(proposalItem.votes) : 0,
          };
        });
        return {
          ...res.data,
          id: proposal.id,
          cid: proposal.cid,
          option,
          proposalStatus,
        };
      });
      return results;
    } catch (error) {
      console.error(error);
    }
  };

  /**
   * filter proposal list
   * @param status
   */
  const handleFilter = async (status: number) => {
    setProposalStatus(status);
  }

  const handleJump = (item: ProposalList) => {
    const router = `/${item.proposalStatus === IN_PROGRESS_STATUS ? "vote" : "votingResults"}/${item.id}/${item.cid}`;
    navigate(router, {state: item});
  }

  const handleCreate = () => {
    if (!isConnected) {
      openConnectModal && openConnectModal();
      return false;
    }
    navigate("/createVote");
  }

  const handlePageChange = (page: number) => {
    setProposalStatus(VOTE_ALL_STATUS);
    setPage(page);
  }

  /**
   * render proposal list
   * @param list
   */
  const renderList = (list: ProposalList[]) => {
    if (proposalStatus !== VOTE_ALL_STATUS) {
      list = list.filter(item => item.proposalStatus === proposalStatus);
      if (list.length === 0) {
        return (
          <Empty
            className='empty'
            description={
              <span className='text-white'>No Data</span>
            }
          />
        );
      }
    }
    return list.map((item: ProposalList, index: number) => {
      const proposal = VOTE_LIST?.find((proposal: ProposalFilter) => proposal.value === item.proposalStatus);
      return (
        <div
          key={item.cid + index}
          className="rounded-xl border border-[#313D4F] bg-[#273141] px-[30px] py-[12px] mb-8"
        >
          <div className="flex justify-between mb-3">
            <a
              target='_blank'
              rel="noopener"
              href={`${chain?.blockExplorers?.default.url}/address/${item.address}`}
              className="flex justify-center items-center"
            >
              <img className="w-[20px] h-[20px] rounded-full mr-2" src={`${web3AvatarUrl}:${item.address}`} alt="" />
              <div className="truncate text-white">
                {EllipsisMiddle({suffixCount: 4, children: item.address})}
              </div>
            </a>
            <div
              className={`${proposal?.color} h-[26px] px-[12px] text-white rounded-xl`}>
              { proposal?.label }
            </div>
          </div>
          <div className="relative mb-4 line-clamp-2 break-words break-all text-lg pr-[80px] leading-7 cursor-pointer"
               onClick={() => {
                 handleJump(item)
               }}>
            <h3 className="inline pr-2 text-2xl font-semibold text-white">
              {item.name}
            </h3>
          </div>
          <div className="mb-2 line-clamp-2 break-words text-lg cursor-pointer" onClick={() => {
            handleJump(item)
          }}>
            {item.descriptions}
          </div>
          {
            item.proposalStatus === COMPLETED_STATUS &&
              <div>
                {
                  item.option?.map((option: ProposalOption, index: number) => {
                    return (
                      <div className="relative mt-1 w-full" key={option.name + index}>
                        <div className="absolute ml-3 flex items-center leading-[43px] text-white">
                          {
                            option.count > 0 &&
                              <svg viewBox="0 0 24 24" width="1.2em" height="1.2em" className="-ml-1 mr-2 text-sm">
                                  <path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round"
                                        strokeWidth="2" d="m5 13l4 4L19 7"></path>
                              </svg>
                          }
                          {option.name}</div>
                        <div className="absolute right-0 mr-3 leading-[40px] text-white">{option.count}%</div>
                        <div className="h-[40px] rounded-md bg-[#1b2331]" style={{width: `${option.count}%`}}></div>
                      </div>
                    )
                  })
                }
              </div>
          }
          <div className="text-[#8B949E] text-sm mt-4">

            <span className="mr-2">Expiration Time:</span>
            {dayjs(item.showTime).format('MMM.D, YYYY, h:mm A')} ({item.GMTOffset})
          </div>
        </div>
      )
    })
  }

  const renderContent = () => {
    if (loading) {
      return (
        <Loading />
      );
    }

    if (proposalList.length === 0) {
      return (
        <Empty
          className='empty'
          description={
            <span className='text-white'>No Data</span>
          }
        />
      );
    }

    return (
      <Spin spinning={loading}>
        <div className='home-table overflow-auto'>
          {
            renderList(proposalList)
          }
          <Row justify='end'>
            <ConfigProvider theme={{ algorithm: theme.darkAlgorithm }}>
              <Pagination
                simple
                current={page}
                pageSize={pageSize}
                total={total}
                onChange={handlePageChange}
              />
            </ConfigProvider>
          </Row>
        </div>
      </Spin>
    );
  };

  return (
    <div className="home_container main">
      <div className="flex justify-between items-center rounded-xl border border-[#313D4F] bg-[#273141] mb-8 px-[30px]">
        <div className="flex justify-between">
          <ListFilter
            name="Status"
            value={proposalStatus}
            list={filterList}
            onChange={handleFilter}
          />
        </div>
        <button
          className="h-[40px] bg-sky-500 hover:bg-sky-700 text-white py-2 px-4 rounded-xl"
          onClick={handleCreate}
        >
          Create A Proposal
        </button>
      </div>
      {
        renderContent()
      }
    </div>
  )
}

export default Home;
