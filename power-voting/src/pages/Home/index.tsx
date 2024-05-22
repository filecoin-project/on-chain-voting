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
import { Row, Empty, Pagination, message } from "antd";
import {useAccount, useReadContract, useReadContracts, BaseError} from "wagmi";
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
  PENDING_STATUS,
  proposalResultApi, worldTimeApi
} from '../../common/consts';
import ListFilter from "../../components/ListFilter";
import EllipsisMiddle from "../../components/EllipsisMiddle";
import {ProposalData, ProposalFilter, ProposalList, ProposalOption, ProposalResult} from '../../common/types';
import Loading from "../../components/Loading";
import {markdownToText, getContractAddress} from "../../utils";
import fileCoinAbi from "../../common/abi/power-voting.json";
import {useCurrentTimezone} from "../../common/store";

dayjs.extend(utc);
dayjs.extend(timezone);

function useLatestId(chainId: number) {
  const { data: latestId, isLoading: getLatestIdLoading } = useReadContract({
    address: getContractAddress(chainId, 'powerVoting'),
    abi: fileCoinAbi,
    functionName: 'proposalId',
  });
  return {
    latestId,
    getLatestIdLoading
  };
}

function useProposalDataSet(params: any) {
  const { chainId, total, page, pageSize } = params;
  const contracts: any[] = [];
  const offset = (page - 1) * pageSize;
  // Generate contract calls for fetching proposals based on pagination
  for (let i = total - offset; i > Math.max(total - offset - pageSize, 0); i--) {
    contracts.push({
      address: getContractAddress(chainId, 'powerVoting'),
      abi: fileCoinAbi,
      functionName: 'idToProposal',
      args: [i],
    });
  }
  const {
    data: proposalData,
    isLoading: getProposalIdLoading,
    isSuccess: getProposalIdSuccess,
    error,
  } = useReadContracts({
    contracts: contracts,
    query: { enabled: !!contracts.length }
  });
  return {
    proposalData: proposalData || [],
    getProposalIdLoading,
    getProposalIdSuccess,
    error,
  };
}

const Home = () => {
  const navigate = useNavigate();

  const {chain, address, isConnected} = useAccount();
  const chainId = chain?.id || 0;
  const {openConnectModal} = useConnectModal();

  const [filterList, setFilterList] = useState([
    {
      label: "All",
      value: VOTE_ALL_STATUS
    }
  ])

  const [proposalStatus, setProposalStatus] = useState(VOTE_ALL_STATUS);
  const [proposalList, setProposalList] = useState<ProposalList[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(5);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const timezone = useCurrentTimezone((state: any) => state.timezone);
  const [messageApi, contextHolder] = message.useMessage();

  const { latestId, getLatestIdLoading } = useLatestId(chainId);
  const { proposalData, getProposalIdLoading, getProposalIdSuccess, error } = useProposalDataSet({
    chainId,
    total: Number(latestId),
    page,
    pageSize,
  });

  useEffect(() => {
    if (error) {
      messageApi.open({
        type: 'error',
        content: (error as BaseError)?.shortMessage || error?.message,
      });
    }
  }, [error]);

  useEffect(() => {
    if (getProposalIdSuccess) {
      getProposalList(page);
    }
  }, [getProposalIdSuccess]);

  useEffect(() => {
    if (isConnected && !loading && !getLatestIdLoading && !getProposalIdLoading) {
      getProposalList(page);
    }
  }, [chain, page, address, isConnected]);

  /**
   * get proposal list
   * @param page
   */
  const getProposalList = async (page: number) => {
    setLoading(true);
    // Convert latest ID to number
    const total = latestId ? Number(latestId) : 0;
    // Calculate the offset based on the current page number

    const offset = (page - 1) * pageSize;
    setTotal(total);
    try {
      const list = await Promise.all(proposalData.map(async(data, index) => {
        const { result } = data as any;
        const proposalId = total - offset - index;
        const params = {
          proposalId,
          network: chainId
        };
        // Fetch proposal results data from the API
        const { data: { data: resultData } } = await axios.get(proposalResultApi, { params });
        // Map proposal results data to a more structured format
        const proposalResults = resultData.map((item: ProposalResult) => ({
          optionId: item.optionId,
          votes: item.votes
        }));
        // Return formatted proposal object
        return {
          id: proposalId,
          cid: result[0],
          creator: result[2],
          startTime: Number(result[3]),
          expTime: Number(result[4]),
          proposalType: Number(result[1]),
          proposalResults
        };
      }));
      // Process and set the fetched proposal list
      const proposalsList = await getList(list);
      const originList = proposalsList || [];
      // Set filter list for proposal filtering
      setFilterList(VOTE_FILTER_LIST);
      // Set the proposal list state
      setProposalList(originList);
    } catch (e) {
      console.log(e);
    } finally {
      setLoading(false);
    }
  }

  /**
   * get proposal info
   * @param proposals
   */
  const getList = async (proposals: ProposalData[]) => {
    // IPFS URL list
    const ipfsUrls = proposals.map(
      (_item: ProposalData) => `https://${_item.cid}.ipfs.w3s.link/`
    );
    try {
      // IPFS data List
      const responses = await Promise.all(ipfsUrls.map((url: string) => axios.get(url)));
      const { data } = await axios.get(worldTimeApi);
      const results: ProposalList[] = responses.map((res, i: number) => {
        const  proposal = proposals[i];
        const now = data?.unixtime;
        let proposalStatus = 0;
        // Set proposal status
        if (now < proposal.startTime) {
          proposalStatus = PENDING_STATUS;
        } else {
          if (now >= proposal.expTime) {
            if (proposal.proposalResults.length === 0) {
              proposalStatus = VOTE_COUNTING_STATUS
            } else {
              proposalStatus = COMPLETED_STATUS
            }
          } else {
            proposalStatus = IN_PROGRESS_STATUS
          }
        }
        // Prepare option
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

  /**
   * page jump
   * @param item
   */
  const handleJump = (item: ProposalList) => {
    const router = `/${[PENDING_STATUS, IN_PROGRESS_STATUS].includes(item.proposalStatus) ? "vote" : "votingResults"}/${item.id}/${item.cid}`;
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
    // Reset vote status when page change
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
    }
    return list.map((item: ProposalList, index: number) => {
      const proposal = VOTE_LIST?.find((proposal: ProposalFilter) => proposal.value === item.proposalStatus);
      const maxOption = item.option.reduce((prev, current) => {
        return (prev.count > current.count) ? prev : current;
      });
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
          className="rounded-xl border border-[#313D4F] bg-[#273141] px-[30px] py-[12px] mb-8"
        >
          <div className="flex justify-between mb-3">
            <a
              target='_blank'
              rel="noopener"
              href={href}
              className="flex justify-center items-center"
            >
              <img className="w-[20px] h-[20px] rounded-full mr-2" src={img} alt="" />
              <div className="truncate text-white">
                {item.githubName || EllipsisMiddle({suffixCount: 4, children:  item.address})}
              </div>
            </a>
            <div
              className={`${proposal?.color} h-[26px] px-[12px] text-white rounded-xl`}>
              { proposal?.label }
            </div>
          </div>
          <div className="relative mb-4 line-clamp-2 break-words break-all text-lg pr-[80px] leading-7 cursor-pointer"
               onClick={() => {
                 handleJump(item);
               }}>
            <h3 className="inline pr-2 text-2xl font-semibold text-white">
              {item.name}
            </h3>
          </div>
          <div className="mb-2 line-clamp-2 break-words text-lg cursor-pointer" onClick={() => {
            handleJump(item)
          }}>
            {markdownToText(item.descriptions)}
          </div>
          {
            maxOption.count > 0 &&
              <div>
                {
                  item.option?.map((option: ProposalOption, index: number) => {
                    return (
                      <div className="relative mt-1 w-full" key={option.name + index}>
                        <div className="absolute ml-3 flex items-center leading-[43px] text-white">
                          {
                            option.count > 0 && option.count ===  maxOption.count &&
                              <svg viewBox="0 0 24 24" width="1.2em" height="1.2em" className="-ml-1 mr-2 text-sm">
                                  <path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round"
                                        strokeWidth="2" d="m5 13l4 4L19 7" />
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
            <span className="mr-2">End Time:</span>
            {dayjs(item.expTime * 1000).format('MMM.D, YYYY, h:mm A')} ({timezone})
          </div>
        </div>
      )
    })
  }

  const renderContent = () => {
    // Display loading when data is loading
    if (getProposalIdLoading || getLatestIdLoading || loading) {
      return (
        <Loading />
      );
    }

    // Display empty when data is empty
    if (proposalData.length === 0) {
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
      <div className='home-table overflow-auto'>
        {
          renderList(proposalList)
        }
        <Row justify='end'>
          <Pagination
            simple
            showSizeChanger={false}
            current={page}
            pageSize={pageSize}
            total={total}
            onChange={handlePageChange}
          />
        </Row>
      </div>
    );
  };

  return (
    <div className="home_container main">
      {contextHolder}
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
