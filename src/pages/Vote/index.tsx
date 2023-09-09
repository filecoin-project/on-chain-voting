import React, { useEffect, useState } from "react";
import { useNavigate, Link, useParams } from "react-router-dom";
import { useForm } from "react-hook-form";
import { InputNumber, message } from "antd";
import axios from 'axios';
import dayjs from 'dayjs';
import {getIpfsId, useDynamicContract} from "../../hooks";
import { useAccount, useNetwork } from "wagmi";
import { useConnectModal, useChainModal } from "@rainbow-me/rainbowkit";
import MDEditor from "../../components/MDEditor";
import EllipsisMiddle from "../../components/EllipsisMiddle";
import {
  IN_PROGRESS_STATUS,
  SINGLE_VOTE,
  MULTI_VOTE,
  VOTE_COUNTING_STATUS, WRONG_NET_STATUS, web3AvatarUrl
} from "../../common/consts";
import { timelockEncrypt, roundAt, mainnetClient, Buffer } from "../../../tlock-js/src"
import "./index.less";

const totalPercentValue = 100;

const Vote = () => {
  const { chain } = useNetwork();
  const chainId = chain?.id || 0;
  const { isConnected } = useAccount();

  const { id, cid } = useParams();
  const [votingData, setVotingData] = useState({} as any);
  const { openConnectModal } = useConnectModal();
  const { openChainModal } = useChainModal();

  const navigate = useNavigate()
  const [options, setOptions] = useState([] as any)

  const [loading, setLoading] = useState(false);

  const {
    formState: { errors }
  } = useForm({
    defaultValues: {
      option: votingData?.VoteType === MULTI_VOTE ? [] : null
    }
  })

  useEffect(() => {
    initState();
  }, [chain])

  const initState = async () => {
    const res = await axios.get(`https://${cid}.ipfs.nftstorage.link/`);
    const data = res.data;
    const now = dayjs().unix();
    let voteStatus = null;
    if (data.chainId !== chainId) {
      voteStatus = WRONG_NET_STATUS;
      if (isConnected) {
        openChainModal && openChainModal();
      } else {
        openConnectModal && openConnectModal();
      }
    } else {
      if (now <= res.data.Time) {
        voteStatus = IN_PROGRESS_STATUS;
      } else {
        voteStatus = VOTE_COUNTING_STATUS;
      }
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

  const handleEncrypt = async (value: any) => {
    const payload = Buffer.from(JSON.stringify(value));

    const chainInfo = await mainnetClient().chain().info();

    const time = new Date().valueOf();

    const roundNumber = roundAt(time, chainInfo) // drand 随机数索引

    const ciphertext = await timelockEncrypt(
      roundNumber,
      payload,
      mainnetClient()
    )

    return ciphertext;
  }

  const startVoting = async () => {
    // 获取有投票的索引，判断是否填写了投票

    const countIndex = options.findIndex((item: any) => item.count > 0);
    if (countIndex < 0) {
      message.warning("Please choose a option to vote");
    } else {
      setLoading(true);
      // vote params
      let params = [];
      if (votingData?.VoteType === SINGLE_VOTE) {
        params.push([`${countIndex}`, `${options[countIndex].count}`])
      } else {
        options.map((item: any, index: number) => {
          params.push([`${index}`, `${item.count}`])
        })
      }
      const encryptValue = await handleEncrypt(params);
      const optionId = await getIpfsId(encryptValue);

      if (isConnected) {
        const { voteApi } = useDynamicContract(chainId);
        const res = await voteApi(Number(id), optionId);
        if (res.code === 200) {
          message.success("Vote successful!", 3);
          setTimeout(() => {
            navigate("/")
          }, 3000)
        } else if (res.code === 401) {
          message.error(res.msg)
        }
        setLoading(false)
      } else {
        // @ts-ignore
        openConnectModal && openConnectModal()
      }
    }
  }

  const cancelVoting= async () => {
    if (isConnected) {
      setLoading(true);
      const { cancelVotingApi } = useDynamicContract(chainId);
      const res = await cancelVotingApi(Number(id));
      if (res.code === 200) {
        message.success("Cancel successful!", 3);
        setTimeout(() => {
          navigate("/")
        }, 3000)
      } else if (res.code === 401) {
        message.error(res.msg)
      }
      setLoading(false)
    } else {
      // @ts-ignore
      openConnectModal && openConnectModal()
    }
  }

  const handleOptionChange = (index: number, count: number) => {
    setOptions((prevState: any[]) => {
      return prevState.map((item: any, preIndex) => {
        let currentTotal = 0;
        currentTotal += count;
        if (preIndex === index) {
          return {
            ...item,
            count
          }
        } else {
          return {
            ...item,
            disabled: votingData?.VoteType === SINGLE_VOTE && count > 0 || votingData?.VoteType === MULTI_VOTE && currentTotal === 100
          }
        }
      })
    })
  }

  const handleCountChange = (type: string, index: number, item: any) => {
    if (item.disabled) return false;

    let currentCount: number;
    const restTotal = options.reduce(((acc: number, current: any) => acc + current.count), 0) - item.count;
    const max = totalPercentValue - restTotal;
    const min = 0
    if (type === "decrease") {
      currentCount = item.count - 1 < min ? min : item.count - 1;
    } else {
      currentCount = item.count + 1 > max ? max : item.count + 1;
    }
    handleOptionChange(index, currentCount);
  }

  const countMax = (options: any, count: number) => {
    const restTotal = options.reduce(((acc: number, current: any) => acc + current.count), 0) - count;
    return totalPercentValue - restTotal;
  }

  const handleVoteStatusTag = (status: number) => {
    switch (status) {
      case WRONG_NET_STATUS:
        return {
          name: 'Wrong network',
          color: 'bg-red-700',
        };
      case IN_PROGRESS_STATUS:
        return {
          name: 'In Progress',
          color: 'bg-green-700',
        };
      case VOTE_COUNTING_STATUS:
        return {
          name: 'Vote Counting',
          color: 'bg-yellow-700',
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
            {votingData?.Name}
          </h1>
          {
            votingData.voteStatus || votingData.voteStatus === 0 &&
              <div className="flex justify-between mb-6">
                  <div className="flex items-center justify-between w-full mb-1 sm:mb-0">
                      <button
                          className={`${handleVoteStatusTag(votingData.voteStatus).color} bg-[#6D28D9] h-[26px] px-[12px] text-white rounded-xl mr-4`}>
                        {handleVoteStatusTag(votingData.voteStatus).name}
                      </button>
                      <div className="flex items-center justify-center">
                          <img className="w-[20px] h-[20px] rounded-full mr-2" src={`${web3AvatarUrl}:${votingData.Address}`} alt="" />
                          <a
                              className="text-white"
                              target="_blank"
                              rel="noopener"
                              href={`${chain?.blockExplorers?.default.url}/address/${votingData?.Address}`}
                          >
                            {EllipsisMiddle({ suffixCount: 4, children: votingData?.Address })}
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
              value={votingData?.Descriptions}
              readOnly={true}
              view={{ menu: false, md: false, html: true, both: false, fullScreen: true, hideMenu: false }}
              onChange={() => {
              }}
            />
          </div>
          {
            votingData?.voteStatus === IN_PROGRESS_STATUS &&
              <div className="border-[#313D4F] mt-6 border-skin-border bg-skin-block-bg text-base md:rounded-xl md:border border-solid">
                  <div className="group flex h-[57px] !border-[#313D4F] justify-between items-center border-b px-4 pb-[12px] pt-3 border-solid">
                      <h4 className="text-xl">
                          Cast Your Vote
                      </h4>
                      <div className='text-base'>{totalPercentValue} %</div>
                  </div>
                  <div className="p-4 text-center">
                    {
                      options.map((item: any, index: number) => {

                        return (
                          <div className="mb-4 space-y-3 leading-10" key={item.name + index}>
                            <div
                              className="w-full h-[45px] !border-[#313D4F] flex justify-between items-center pl-4 md:border border-solid rounded-full">
                              <div className="text-ellipsis h-[100%] overflow-hidden">{item.name}</div>
                              <div className="w-[180px] h-[45px] flex items-center">
                                <div onClick={() => {
                                  handleCountChange("decrease", index, item)
                                }}
                                     className={`w-[35px] border-x border-solid !border-[#313D4F] text-white font-semibold ${item.disabled ? "cursor-not-allowed" : "cursor-pointer"}`}>-
                                </div>
                                <InputNumber
                                  disabled={item.disabled}
                                  className="text-white bg-transparent focus:outline-none"
                                  controls={false}
                                  min={0}
                                  max={countMax(options, item.count)}
                                  precision={0}
                                  value={item.count}
                                  onChange={(value) => {
                                    handleOptionChange(index, Number(value))
                                  }}
                                />
                                <div onClick={() => {
                                  handleCountChange("increase", index, item)
                                }}
                                     className={`w-[35px] border-x border-solid !border-[#313D4F] text-white font-semibold ${item.disabled ? "cursor-not-allowed" : "cursor-pointer"}`}>+
                                </div>
                                <div className="w-[40px] text-center">%</div>
                              </div>
                            </div>
                          </div>
                        )
                      })
                    }
                    <button onClick={startVoting} className="w-full h-[40px] bg-sky-500 hover:bg-sky-700 text-white py-2 px-6 rounded-full" type="submit" disabled={loading}>
                        Vote
                    </button>
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
                  <b>Vote Type</b>
                  <span
                    className="float-right text-white">{["Single", "Multiple"][votingData?.VoteType - 1]} Choice Voting</span>
                </div>
                <div>
                  <b>Exp. Time</b>
                  <span className="float-right text-white">
                    {dayjs(votingData?.showTime).format('MMM.D, YYYY, h:mm A')}
                  </span>
                </div>
                <div>
                  <b>Exp. Timezone</b>
                  <span className="float-right text-white">
                    {votingData?.GMTOffset}
                  </span>
                </div>
                <div>
                  <b>Snapshot</b>
                  <span className="float-right text-white">Takes At Exp. Time</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        {
          votingData?.voteStatus === IN_PROGRESS_STATUS &&
          <button onClick={cancelVoting} className="w-full h-[40px] bg-red-500 hover:bg-red-700 text-white py-2 px-6 rounded-full mt-6" type="submit" disabled={loading}>
            Cancel Proposal
          </button>
        }
      </div>
    </div>
  )
}

export default Vote;
