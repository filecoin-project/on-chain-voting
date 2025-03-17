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

import { CheckCircleOutlined, CloseCircleOutlined, InfoCircleOutlined } from '@ant-design/icons';
import { Empty, Popover, Table } from 'antd';
import React from 'react';
import { useTranslation } from 'react-i18next';
import type { Chain } from "viem";
import { VOTE_OPTIONS, web3AvatarUrl } from "../../common/consts";
import type { ProposalVotes } from "../../common/types";
import { bigNumberToFloat, convertBytes } from "../../utils";
import EllipsisMiddle from "../EllipsisMiddle";
import './index.less';
interface Props {
  voteList: ProposalVotes[];
  chain?: Chain;
  powerPercent: {
    sp_percentage: number,
    client_percentage: number,
    developer_percentage: number,
    token_holder_percentage: number
  }
}

const VoteList: React.FC<Props> = ({ voteList, chain, powerPercent }) => {
  const { t } = useTranslation();
  const columns = [
    {
      title: t('content.role'),
      dataIndex: 'role',
      key: 'role',
    },
    {
      title: t('content.power'),
      dataIndex: 'power',
      key: 'power',
    },
    {
      title: t('content.totalPower'),
      dataIndex: 'total',
      key: 'total',
    },
    {
      title: t('content.percent'),
      dataIndex: 'percent',
      key: 'percent',
      render: (text: string, record: any) => {
        return text === '0%' && record.power === '0' ? 'NO VOTES' : text;
      }
    },
    {
      title: t('content.blockHeight'),
      dataIndex: 'powerBlockHeight',
      key: 'powerBlockHeight',
    },
  ];
  const getPowerData = (votePower: ProposalVotes) => {
    return [
      {
        key: 'sp',
        role: 'SP',
        powerBlockHeight: votePower.blockHeight,
        power: convertBytes(votePower.spPower),
        total: convertBytes(votePower.totalSpPower),
        percent: votePower.spPowerPercent === 0 ? '0%' : `${votePower.spPowerPercent.toFixed(2)}%`,
        powerPercent: powerPercent.sp_percentage
      },
      {
        key: 'client',
        role: 'Client',
        powerBlockHeight: votePower.blockHeight,
        power: convertBytes(Number(votePower.clientPower) / (10 ** 18)),
        total: convertBytes(Number(votePower.totalClientPower) / (10 ** 18)),
        percent: votePower.clientPowerPercent === 0 ? '0%' : `${votePower.clientPowerPercent.toFixed(2)}%`,
        powerPercent: powerPercent.client_percentage
      },
      {
        key: 'developer',
        role: 'Developer',
        powerBlockHeight: votePower.blockHeight,
        power: votePower.developerPower,
        total: votePower.totalDeveloperPower,
        percent: votePower.developerPowerPercent === 0 ? '0%' : `${votePower.developerPowerPercent.toFixed(2)}%`,
        powerPercent: powerPercent.developer_percentage
      },
      {
        key: 'tokenHolder',
        role: 'TokenHolder',
        powerBlockHeight: votePower.blockHeight,
        power: bigNumberToFloat(votePower.tokenHolderPower),
        total: bigNumberToFloat(votePower.totalTokenHolderPower),
        percent: votePower.tokenHolderPowerPercent === 0 ? '0%' : `${votePower.tokenHolderPowerPercent.toFixed(2)}%`,
        powerPercent: powerPercent.token_holder_percentage
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
    let totalPercent = `${t('content.totalPercent')} = `;
    // Initialize count for non-zero total values
    // Initialize an array to store non-zero percent values
    const arr: any[] = [];

    data.forEach(item => {
      const { total, percent, power } = item;
      // Check if total is not '0'
      if (total !== '0') {
        arr.push({ percent, total, power });
      }
    });
    let total = 0;
    let power = 0
    if (arr.length) {
      arr.forEach((item, index) => {
        // Check if it's not the last item in the array
        total += Number(item.total);
        power += Number(item.power);
        if (index < arr.length - 1) {
          // Append percent value and count with a plus sign

        } else {
          // Append percent value and count without a plus sign
          // totalPercent += `${item.total} * ${item.percent}`;
        }
      });
      // Append the final vote percentage
      totalPercent += `${power} / ${total} = ${votes.toFixed(2)}%`;
    } else {
      totalPercent += '0%';
    }

    return <div>{totalPercent}</div>;
  }
  return (
    <div className="border-y border-skin-border bg-skin-block-bg text-base md:rounded-xl md:border mt-[20px] mb-[20px]">
      <div className="group flex h-[57px] justify-between rounded-t-none border-b border-skin-border px-6 pb-[12px] pt-3 md:rounded-t-lg">
        <h4 className="flex items-center">
          <div className="font-medium">{t('content.votes')}</div>
        </h4>
        <div className="flex items-center" />
      </div>
      {
        voteList?.length > 0 ?
          <div className="voteList leading-5 sm:leading-6 overflow-auto">
            {
              voteList?.map((item: any, index: number) => {
                const powers = []
                if (item.tokenHolderPower > 0) {
                  powers.push("TokenHolder")
                }
                if (item.spPower > 0) {
                  powers.push("Sp")
                }
                if (item.developerPower > 0) {
                  powers.push("Developer")
                }
                if (item.clientPower > 0) {
                  powers.push("Client")
                }
                const isApprove = item.optionName === VOTE_OPTIONS[0]
                return (
                  <div className={`flex items-center gap-3 border-t px-4 py-[14px] ${index === 0 && '!border-0'}`} key={item.address + index}>
                    <div className="flex items-center">
                      <img className="w-[20px] h-[20px] rounded-full mr-2" src={`${web3AvatarUrl}:${item.address}`} alt="" />
                      <a
                        className="text-[#313D4F] text-sm font-medium"
                        target="_blank"
                        rel="noopener noreferrer"
                        href={`${chain?.blockExplorers?.default.url}/address/${item?.address}`}
                      >
                        {EllipsisMiddle({ suffixCount: 4, children: item?.address })}
                      </a>
                      <div className='ml-[5px]'>
                        {powers.length > 0 && <div className='flex'>
                          {
                            powers.slice(0, 4).map((v, index) => {
                              return <div
                                key={index}
                                style={{ marginLeft: "4px", borderColor: "#C3E5FF", backgroundColor: "#E7F4FF", color: "#005292" }}
                                className={`font-medium text-xs flex items-center justify-center border-solid h-[22px] px-[8px] rounded-full`}>
                                {v}
                              </div>
                            })
                          }
                        </div>}
                      </div>
                    </div>

                    <div className='absolute right-[10px] flex items-center'>
                      <div className="mr-[5px] flex min-w-[110px] items-center justify-end whitespace-nowrap text-center text-skin-link xs:w-[130px] xs:min-w-[130px] cursor-pointer">
                        <Popover content={
                          <Table
                            rowKey={(record: any) => record.key}
                            dataSource={getPowerData(item)}
                            columns={columns}
                            pagination={false}
                            footer={(currentData: any) => renderFooter(currentData, item.tokenHolderPowerPercent)}
                          />
                        }>
                          <span className='text-[14px] text-[#273141] text-sm'>{item.tokenHolderPowerPercent.toFixed(2)}% <InfoCircleOutlined style={{ fontSize: 14 }} /></span>
                        </Popover>
                      </div>

                      <div className="flex truncate px-2 text-skin-link ">
                        <div className="pr-[20px] text-right w-[100px] text-c truncate text-skin-link text-xs font-medium" style={{ color: isApprove ? "green" : "red" }}>
                          {isApprove ? <CheckCircleOutlined style={{ fontSize: 14, marginRight: "4px" }} /> : <CloseCircleOutlined style={{ fontSize: 14, marginRight: "4px" }} />}
                          {isApprove ? t('content.approved') : t('content.rejected')}</div>
                      </div>
                    </div>
                  </div>
                )
              })
            }
          </div> :
          <Empty
            className='my-12'
            description={
              <span className='text-[#313D4F]'>{t('content.noData')}</span>
            }
          />
      }
    </div>
  )
}

export default VoteList;