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

import React from 'react';
import { Empty, Table, Popover } from 'antd';
import { CheckCircleOutlined, CloseCircleOutlined, InfoCircleOutlined } from '@ant-design/icons';
import EllipsisMiddle from "../EllipsisMiddle";
import { VOTE_OPTIONS, web3AvatarUrl } from "../../common/consts";
import type { Chain } from "viem";
import type { ProposalHistory } from "../../common/types";
import './index.less';
import { bigNumberToFloat, convertBytes } from "../../utils";

interface Props {
  voteList: ProposalHistory[];
  chain?: Chain;
}

const VoteList: React.FC<Props> = ({ voteList, chain }) => {

  const columns = [
    {
      title: 'Role',
      dataIndex: 'role',
      key: 'role',
    },
    {
      title: 'Power',
      dataIndex: 'power',
      key: 'power',
    },
    {
      title: 'Total Power',
      dataIndex: 'total',
      key: 'total',
    },
    {
      title: 'Percent',
      dataIndex: 'percent',
      key: 'percent',
      render: (text: string, record: any) => {
        return text === '0%' && record.power === '0' ? 'NO VOTES' : text;
      }
    },
    {
      title: 'Block Height',
      dataIndex: 'powerBlockHeight',
      key: 'powerBlockHeight',
    },
  ];

  const getPowerData = (votePower: any) => {
    return [
      {
        key: 'sp',
        role: 'SP',
        powerBlockHeight: votePower.powerBlockHeight,
        power: convertBytes(votePower.spPower),
        total: convertBytes(votePower.totalSpPower),
        percent: `${votePower.spPowerPercent}%`,
      },
      {
        key: 'client',
        role: 'Client',
        powerBlockHeight: votePower.powerBlockHeight,
        power: convertBytes(Number(votePower.clientPower) / (10 ** 18)),
        total: convertBytes(Number(votePower.totalClientPower) / (10 ** 18)),
        percent: `${votePower.clientPowerPercent}%`,
      },
      {
        key: 'developer',
        role: 'Developer',
        powerBlockHeight: votePower.powerBlockHeight,
        power: votePower.developerPower,
        total: votePower.totalDeveloperPower,
        percent: `${votePower.developerPowerPercent}%`,
      },
      {
        key: 'tokenHolder',
        role: 'TokenHolder',
        powerBlockHeight: votePower.powerBlockHeight,
        power: bigNumberToFloat(votePower.tokenHolderPower),
        total: bigNumberToFloat(votePower.totalTokenHolderPower),
        percent: `${votePower.tokenHolderPowerPercent}%`
      },
    ];
  }

  /**
   * Show the weight calculation process
   * @param data
   * @param votes
   */
  const renderFooter = (data: any[], votes: number) => {
    // Initialize the string for total percent calculation
    let totalPercent = "Total Percent = ";
    // Initialize count for non-zero total values
    let count = 0;
    // Initialize an array to store non-zero percent values
    const arr: string[] = [];


    data.forEach(item => {
      const { total, percent } = item;
      // Check if total is not '0'
      if (total !== '0') {
        arr.push(percent);
        count++;
      }
    });

    arr.forEach((item, index) => {
      // Check if it's not the last item in the array
      if (index < arr.length - 1) {
        // Append percent value and count with a plus sign
        totalPercent += `${item} / ${count} + `;
      } else {
        // Append percent value and count without a plus sign
        totalPercent += `${item} / ${count}`;
      }
    });
    // Append the final vote percentage
    totalPercent += `= ${votes}%`;

    return <div>{totalPercent}</div>;
  }

  return (
    <div className="border-y border-skin-border bg-skin-block-bg text-base md:rounded-xl md:border my-12">
      <div className="group flex h-[57px] justify-between rounded-t-none border-b border-skin-border px-6 pb-[12px] pt-3 md:rounded-t-lg">
        <h4 className="flex items-center">
          <div className="font-semibold">Votes</div>
        </h4>
        <div className="flex items-center" />
      </div>
      {
        voteList?.length > 0 ?
          <div className="voteList leading-5 sm:leading-6 max-h-[260px] overflow-auto">
            {
              voteList?.map((item: any, index: number) => {

                const powers = []

                if (item.totalTokenHolderPower > 0) {
                  powers.push("TokenHolder")
                }
                if (item.totalSpPower > 0) {
                  powers.push("Sp")
                }
                if (item.totalDeveloperPower > 0) {
                  powers.push("Developer")
                }
                if (item.totalClientPower > 0) {
                  powers.push("Client")
                }
                const isApprove = item.optionName === VOTE_OPTIONS[0]
                return (
                  <div className={`flex items-center gap-3 border-t px-4 py-[14px] ${index === 0 && '!border-0'}`} key={item.address + index}>
                    <div className="w-[300px] flex items-center">
                      <img className="w-[20px] h-[20px] rounded-full mr-2" src={`${web3AvatarUrl}:${item.address}`} alt="" />
                      <a
                        className="text-[#313D4F]"
                        target="_blank"
                        rel="noopener noreferrer"
                        href={`${chain?.blockExplorers?.default.url}/address/${item?.address}`}
                      >
                        {EllipsisMiddle({ suffixCount: 4, children: item?.address })}
                      </a>
                      <div >
                        {powers.length > 0 && <div className='flex'>
                          {
                            powers.slice(0, 2).map((v, index) => {
                              return <div
                                key={index}
                                style={{ marginLeft: "4px", borderColor: "#C3E5FF", backgroundColor: "#E7F4FF", color: "#005292" }}
                                className={`flex items-center justify-center border-solid h-[32px] px-[12px] rounded-full`}>
                                {v}
                              </div>
                            })
                          }
                        </div>}
                        {
                          powers.length > 2 && <div className='flex mt-[5px]'>
                            {
                              powers.slice(2).map((v, index) => {
                                return <div
                                  key={index}
                                  style={{ marginLeft: "4px", borderColor: "#C3E5FF", backgroundColor: "#E7F4FF", color: "#005292" }}
                                  className={`flex items-center justify-center border-solid h-[32px] px-[12px] rounded-full`}>
                                  {v}
                                </div>
                              })
                            }
                          </div>
                        }
                      </div>


                    </div>

                    <div className="flex min-w-[110px] items-center justify-end whitespace-nowrap text-center text-skin-link xs:w-[130px] xs:min-w-[130px] cursor-pointer">
                      <Popover content={
                        <Table
                          rowKey={(record: any) => record.key}
                          dataSource={getPowerData(item)}
                          columns={columns}
                          pagination={false}
                          footer={(currentData: any) => renderFooter(currentData, item.votes)}
                        />
                      }>
                        <span>{item.votes}% <InfoCircleOutlined style={{ fontSize: 14 }} /></span>
                      </Popover>
                    </div>

                    <div className="w-[180px] flex truncate px-2 justify-end text-skin-link">
                      <div className="w-[80px] text-c truncate text-skin-link" style={{ color: isApprove ? "green" : "red" }}>
                       {isApprove? <CheckCircleOutlined style={{ fontSize: 14, marginRight: "4px" }} />: <CloseCircleOutlined style={{ fontSize: 14, marginRight: "4px" }} />}
                        {item.optionName}</div>
                    </div>
                  </div>
                )
              })
            }
          </div> :
          <Empty
            className='my-12'
            description={
              <span className='text-[#313D4F]'>No Data</span>
            }
          />
      }
    </div>
  )
}

export default VoteList;