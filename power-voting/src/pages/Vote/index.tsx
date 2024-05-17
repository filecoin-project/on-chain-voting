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
import { useAccount, useWriteContract, useWaitForTransactionReceipt, BaseError } from "wagmi";
import { useConnectModal, useChainModal } from "@rainbow-me/rainbowkit";
import {getContractAddress, getWeb3IpfsId} from '../../utils';
import MDEditor from "../../components/MDEditor";
import EllipsisMiddle from "../../components/EllipsisMiddle";
import LoadingButton from "../../components/LoadingButton";
import {
  IN_PROGRESS_STATUS,
  WRONG_NET_STATUS,
  web3AvatarUrl,
  CHOOSE_VOTE_MSG,
  PENDING_STATUS, STORING_DATA_MSG,
} from "../../common/consts";
import { timelockEncrypt, roundAt, mainnetClient, Buffer } from "tlock-js";
import {ProposalList, ProposalOption} from "../../common/types";
import "./index.less";
import fileCoinAbi from "../../common/abi/power-voting.json";

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

  const {
    data: hash,
    writeContract,
    error,
    isPending: writeContractPending,
    isSuccess: writeContractSuccess,
    reset
  } = useWriteContract();

  useEffect(() => {
    initState();
  }, [chain]);

  useEffect(() => {
    if (writeContractSuccess) {
      message.success(STORING_DATA_MSG);
      navigate("/");
    }
  }, [writeContractSuccess])

  useEffect(() => {
    if (error) {
      message.error((error as BaseError)?.shortMessage || error?.message);
    }
    reset();
  }, [error]);

  const initState = async () => {
    // Fetch data from the IPFS link using the provided CID
    const res = await axios.get(`https://${cid}.ipfs.w3s.link/`);
    const data = res.data;
    let voteStatus = null;
    // Check if the chain ID from the fetched data matches the current chain ID
    if (data.chainId !== chainId) {
      // If not, set the vote status to indicate wrong network status
      voteStatus = WRONG_NET_STATUS;
      // If user is connected, open the chain modal to prompt for network switch
      if (isConnected) {
        openChainModal && openChainModal();
      } else {
        openConnectModal && openConnectModal();
      }
    } else {
      // If chain ID matches, determine the vote status based on current time and start time
      voteStatus = dayjs().unix() < data?.startTime ? PENDING_STATUS : IN_PROGRESS_STATUS;
    }
    // Map each option from the fetched data to include count initialized to 0
    const option = data.option?.map((item: string) => {
      return {
        name: item,
        count: 0,
      };
    });
    // Set the voting data state with the fetched data and additional properties
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
    // Convert value to a JSON string and encode it as a Buffer
    const payload = Buffer.from(JSON.stringify(value));

    // Get chain information from the mainnet client
    const chainInfo = await mainnetClient().chain().info();

    // Calculate time for the voting expiration, or set to 0 if not available
    const time = votingData?.expTime ? new Date(votingData.expTime * 1000).valueOf() : 0;

    // Determine the round number based on the time and chain information
    const roundNumber = roundAt(time, chainInfo);

    // Encrypt the payload using timelock encryption
    const ciphertext = await timelockEncrypt(
      roundNumber,
      payload,
      mainnetClient()
    )

    return ciphertext;
  }

  const startVoting = async () => {
    // Check if a valid option is selected
    if (selectedOptionIndex < 0) {
      // If not, display a warning message
      message.warning(CHOOSE_VOTE_MSG);
    } else {
      // If a valid option is selected, proceed with voting
      setLoading(true);
      // Encrypt the selected option index and weight using handleEncrypt function
      const encryptValue = await handleEncrypt([[`${selectedOptionIndex}`, `100`]]);
      // Get the IPFS ID for the encrypted value
      const optionId = await getWeb3IpfsId(encryptValue);
      // Check if user is connected to the network
      if (isConnected) {
        writeContract({
          abi: fileCoinAbi,
          address: getContractAddress(chain?.id || 0, 'powerVoting'),
          functionName: 'vote',
          args: [
            Number(id),
            optionId,
          ],
        });
        setLoading(false);
      } else {
        // If user is not connected, prompt to connect
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

  const { isLoading: transactionLoading } =
    useWaitForTransactionReceipt({
      hash,
    })

  let href = '';
  let img = '';
  if (votingData?.githubName) {
    href = `https://github.com/${votingData.githubName}`;
    img = `${votingData?.githubAvatar}`;
  } else {
    href = `${chain?.blockExplorers?.default.url}/address/${votingData?.address}`;
    img = `${web3AvatarUrl}:${votingData?.address}`
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
                          <img className="w-[20px] h-[20px] rounded-full mr-2" src={img} alt="" />
                          <a
                              className="text-white"
                              target="_blank"
                              rel="noopener"
                              href={href}
                          >
                            {votingData?.githubName || EllipsisMiddle({suffixCount: 4, children: votingData?.address})}
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
                          <LoadingButton text='Vote' isFull={true} loading={loading || writeContractPending || transactionLoading} handleClick={startVoting} />
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
