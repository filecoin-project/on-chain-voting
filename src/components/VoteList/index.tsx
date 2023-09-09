import React from 'react';
import { Empty } from 'antd';
import EllipsisMiddle from "../EllipsisMiddle";
import {web3AvatarUrl} from "../../common/consts";
import './index.less';

const VoteList = (props: any) => {
  const { voteList, chain } = props;

  const totalVotes = voteList?.reduce(((acc: number, current: any) => acc + current.value), 0) || 0;
  return (
    <div className="border-y border-skin-border bg-skin-block-bg text-base md:rounded-xl md:border my-12">
      <div className="group flex h-[57px] justify-between rounded-t-none border-b border-skin-border px-6 pb-[12px] pt-3 md:rounded-t-lg">
        <h4 className="flex items-center">
          <div className="font-semibold">Votes</div>
          <div className="h-[20px] min-w-[20px] rounded-full bg-[#8b949e] px-1 text-center text-xs leading-5 text-white ml-2 inline-block">{totalVotes}</div>
        </h4>
        <div className="flex items-center" />
      </div>
      {
        voteList?.length > 0 ?
          <div className="voteList leading-5 sm:leading-6 max-h-[260px] overflow-auto">
            {
              voteList?.map((item: any, index: number) => {
                return (
                  <div className={`flex items-center gap-3 border-t px-4 py-[14px] ${index === 0 && '!border-0'}`} key={item.TransactionHash + index}>
                    <div className="flex items-center">
                      <img className="w-[20px] h-[20px] rounded-full mr-2" src={`${web3AvatarUrl}:${item.Address}`} alt="" />
                      <a
                        className="text-white"
                        target="_blank"
                        rel="noopener"
                        href={`${chain?.blockExplorers?.default.url}/address/${item?.Address}`}
                      >
                        {EllipsisMiddle({ suffixCount: 4, children: item?.Address })}
                      </a>
                    </div>
                    <div className="flex-auto truncate px-2 text-center text-skin-link">
                      <div className="truncate text-center text-skin-link">{item.label}</div>
                    </div>
                    <div className="flex w-[110px] min-w-[110px] items-center justify-end whitespace-nowrap text-center text-skin-link xs:w-[130px] xs:min-w-[130px]">
                      <span>{item.value} Vote(s)</span>
                    </div>
                    <div className="flex items-center w-[100px] items-center justify-end whitespace-nowrap text-right text-skin-link">
                      <a
                        className="text-white"
                        target="_blank"
                        rel="noopener"
                        href={`${chain?.blockExplorers?.default.url}/tx/${item?.TransactionHash}`}
                      >
                        <svg viewBox="0 0 24 24" width="1.2em" height="1.2em" className="mb-[2px] ml-1 inline-block text-xs"><path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 6H6a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2v-4M14 4h6m0 0v6m0-6L10 14"></path></svg>
                      </a>
                    </div>
                  </div>
                )
              })
            }
          </div> :
          <Empty
            className='my-12'
            description={
              <span className='text-white'>No Data</span>
            }
          />
      }
    </div>
  )
}

export default VoteList;