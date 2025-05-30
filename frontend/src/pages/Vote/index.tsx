import { InfoCircleOutlined } from "@ant-design/icons"
import { useChainModal, useConnectModal } from "@rainbow-me/rainbowkit"
import { message, Popover, Table } from "antd"
import axios from "axios"
import dayjs from "dayjs"
import { useEffect, useState } from "react"
import { useTranslation } from "react-i18next"
import { Link, useNavigate, useParams } from "react-router-dom"
import VoteStatusBtn from "../../../src/components/VoteStatusBtn"
import { Buffer } from "buffer"
import { mainnetClient, roundAt, timelockEncrypt } from "tlock-js"
import { UserRejectedRequestError } from "viem"
import type { BaseError } from "wagmi"
import { useAccount, useWriteContract } from "wagmi"
import { formatNumberWithCommas, getBlockExplorers, getContractAddress, isFilAddress } from "../../utils"
import { useSendMessage, useAddresses } from "iso-filecoin-react"
import fileCoinAbi from "../../common/abi/power-voting.json"
import {
  calibrationChainId,
  CHOOSE_VOTE_MSG,
  GETTING_POWER_MSG,
  getVoteDetail,
  githubApi,
  IN_PROGRESS_STATUS,
  PENDING_STATUS,
  VOTE_SUCCESS_MSG,
  votePowerGetApi,
  web3AvatarUrl,
  WRONG_NET_STATUS,
} from "../../common/consts"
import { useCurrentTimezone, useStoringHash } from "../../common/store"
import type { ProposalList, ProposalOption } from "../../common/types"
import EllipsisMiddle from "../../components/EllipsisMiddle"
import LoadingButton from "../../components/LoadingButton"
import MDEditor from "../../components/MDEditor"
import { bigNumberToFloat, convertBytes } from "../../utils"
import "./index.less"
import { useFilAddressMessage } from "../../common/hooks.ts"

const Vote = () => {

  const { mutateAsync: sendMessage } = useSendMessage();
  const { chain, isConnected, address } = useAccount();
  const { address0x } = useAddresses({ address: address as string })
  const chainId = chain?.id || calibrationChainId

  const { id } = useParams()
  const { t } = useTranslation()
  const [votingData, setVotingData] = useState({} as ProposalList)
  const [powerDetail, setPowerDetail] = useState<any>()
  const { openConnectModal } = useConnectModal()
  const { openChainModal } = useChainModal()
  const navigate = useNavigate();

  const [options, setOptions] = useState([] as ProposalOption[])
  const [selectedOptionIndex, setSelectedOptionIndex] = useState(-1)

  const [loading, setLoading] = useState(false)
  const timezone = useCurrentTimezone((state: any) => state.timezone)
  const [userInfo, setUserInfo] = useState<{ href: string; img: string }>({
    href: "",
    img: ""
  })

  const [messageApi, contextHolder] = message.useMessage()

  const columns = [
    {
      title: t("content.role"),
      dataIndex: "role",
      key: "role"
    },
    {
      title: t("content.power"),
      dataIndex: "power",
      key: "power"
    },
    {
      title: t("content.snapshotBlockHeight"),
      dataIndex: "blockHeight",
      key: "blockHeight"
    }
  ]

  const storingHash = useStoringHash((state: any) => state.storingHash)

  const getPowerData = (votePower: any) => {
    return [
      {
        key: "sp",
        role: "SP",
        blockHeight: formatNumberWithCommas({ value: votingData?.snapshotInfo?.snapshotHeight }),
        power: convertBytes(votePower?.spPower)
      },
      {
        key: "client",
        role: "Client",
        blockHeight: formatNumberWithCommas({ value: votingData?.snapshotInfo?.snapshotHeight }),
        power: convertBytes(Number(votePower?.clientPower) / 10 ** 18)
      },
      {
        key: "developer",
        role: "Developer",
        blockHeight: formatNumberWithCommas({ value: votingData?.snapshotInfo?.snapshotHeight }),
        power: votePower?.developerPower || 0,
      },
      {
        key: "tokenHolder",
        role: "TokenHolder",
        blockHeight: formatNumberWithCommas({ value: votingData?.snapshotInfo?.snapshotHeight }),
        power: bigNumberToFloat(votePower?.tokenHolderPower)
      }
    ]
  }
  const {
    writeContract,
    error,
    isPending: writeContractPending,
    isSuccess: writeContractSuccess,
    reset
  } = useWriteContract()

  useEffect(() => {
    initState()
  }, [chain, address])

  useEffect(() => {
    if (writeContractSuccess) {
      messageApi.open({
        type: "success",
        content: t(VOTE_SUCCESS_MSG)
      });
      setSelectedOptionIndex(-1);
    }
  }, [writeContractSuccess])

  useEffect(() => {
    if (error) {
      if (error as UserRejectedRequestError) {
        messageApi.open({
          type: "warning",
          content: t("content.rejectedSignature")
        })
      } else {
        messageApi.open({
          type: "error",
          content: (error as BaseError)?.shortMessage || error?.message
        })
      }
    }
    reset()
  }, [error])

  const initState = async () => {
    try {
      const {
        data: { data: voteDetail }
      } = await axios.get(`${getVoteDetail}`, {
        params: { chainId, proposalId: id }
      })

      let voteStatus = null
      // Check if the chain ID from the fetched data matches the current chain ID
      if (voteDetail?.chainId !== chainId) {
        // If not, set the vote status to indicate wrong network status
        voteStatus = WRONG_NET_STATUS
        // If user is connected, open the chain modal to prompt for network switch
        if (isConnected) {
          openChainModal && openChainModal()
        } else {
          openConnectModal && openConnectModal()
        }
      } else {
        // If chain ID matches, determine the vote status based on current time and start time
        voteStatus =
          Math.floor(Date.now() / 1000) < voteDetail?.startTime
            ? PENDING_STATUS
            : IN_PROGRESS_STATUS
      }
      // Map each option from the fetched data to include count initialized to 0
      const option = [
        {
          name: "Approve",
          count: 0
        },
        {
          name: "Reject",
          count: 0
        }
      ]
      // Set the voting data state with the fetched data and additional properties
      setVotingData({
        ...voteDetail,
        id,
        option,
        voteStatus
      })
      setOptions(option)
      const {
        data: { data: powerDetail }
      } = await axios.get(`${votePowerGetApi}`, {
        params: {
          chainId,
          address: isFilAddress(address!) && address0x.data ? address0x.data.toString() : address,
          powerDay: voteDetail?.snapshotInfo?.snapshotDay
        }
      })
      setPowerDetail(powerDetail)
    } catch (e) {
      console.log(e)
    }
  }

  /**
   * timelock encrypt
   * @param value
   */
  const handleEncrypt = async (value: string[][]) => {
    // Convert value to a JSON string and encode it as a Buffer
    const payload = Buffer.from(JSON.stringify(value))

    // Get chain information from the mainnet client
    const chainInfo = await mainnetClient().chain().info()

    // Calculate time for the voting expiration, or set to 0 if not available
    const time = votingData?.endTime
      ? new Date(votingData.endTime * 1000).valueOf()
      : 0

    // Determine the round number based on the time and chain information
    const roundNumber = roundAt(time, chainInfo)

    // Encrypt the payload using timelock encryption
    const ciphertext = await timelockEncrypt(
      roundNumber,
      payload,
      mainnetClient()
    )

    return ciphertext
  }

  const startVoting = async () => {
    const item = storingHash.find((item: any) => item.address === address)

    if (item) {
      messageApi.open({
        type: "warning",
        content: t(GETTING_POWER_MSG)
      })
      return
    }

    // Check if a valid option is selected
    if (selectedOptionIndex < 0) {
      // If not, display a warning message
      messageApi.open({
        type: "warning",
        content: t(CHOOSE_VOTE_MSG)
      })
    } else {
      // If a valid option is selected, proceed with voting
      setLoading(true)
      // Encrypt the selected option index and weight using handleEncrypt function
      const encryptValue = await handleEncrypt([
        [`${selectedOptionIndex === 0 ? "approve" : "reject"}`]
      ])

      // Check if user is connected to the network
      if (isConnected) {
        // f1 or f3 address
        if (address && isFilAddress(address)) {
          try {
            const { message } = await useFilAddressMessage({
              abi: fileCoinAbi,
              contractAddress: getContractAddress(chain?.id || calibrationChainId, 'powerVoting'),
              address: address as string,
              functionName: "vote",
              functionParams: [Number(id), encryptValue],
            })
            await sendMessage(message);
            messageApi.open({
              type: "success",
              content: t(VOTE_SUCCESS_MSG)
            })
            setSelectedOptionIndex(-1);
          } catch (error) {
            console.log(error)
            if (error as UserRejectedRequestError) {
              messageApi.open({
                type: "warning",
                content: t("content.rejectedSignature")
              })
            } else {
              messageApi.open({
                type: "error",
                content: (error as BaseError)?.shortMessage || JSON.stringify(error)
              })
            }
          } finally {
            setLoading(false);
          }
        } else {
          writeContract({
            abi: fileCoinAbi,
            address: getContractAddress(chain?.id || calibrationChainId, 'powerVoting'),
            functionName: 'vote',
            args: [
              Number(id),
              encryptValue,
            ],
          });
        }
        setLoading(false);
      } else {
        // If user is not connected, prompt to Connect
        openConnectModal && openConnectModal();
        setLoading(false);
      }
    }
  }

  const handleOptionClick = (index: number) => {
    setSelectedOptionIndex(index)
  }

  const getUserInfo = async () => {
    let href = ""
    let img = ""
    if (votingData?.githubName) {
      href = `https://github.com/${votingData.githubName}`;
      // const { data } = await axios.get(`${githubApi}/${votingData.githubName}`);
      // const githubAvatar = data.avatar_url;
      img = `${githubApi}/${votingData.githubName}`;
    } else {
      href = getBlockExplorers(chain, votingData?.address)
      img = `${web3AvatarUrl}:${votingData?.address}`
    }
    setUserInfo({ href, img })
  }

  useEffect(() => {
    getUserInfo()
  }, [votingData]);
  useEffect(() => {
    if (!isConnected) {
      navigate("/home");
      return;
    }
  }, [isConnected]);
  return (
    <div className="flex voting">
      {contextHolder}
      <div className="relative w-full pr-5 lg:w-8/12">
        <div className="px-3 mb-6 md:px-0">
          <button>
            <div className="inline-flex items-center gap-1 text-skin-text hover:text-skin-link">
              <Link to="/home" className="flex items-center">
                <svg
                  className="mr-1"
                  viewBox="0 0 24 24"
                  width="1.2em"
                  height="1.2em"
                >
                  <path
                    fill="none"
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="m11 17l-5-5m0 0l5-5m-5 5h12"
                  ></path>
                </svg>
                {t("content.back")}
              </Link>
            </div>
          </button>
        </div>
        <div className="px-3 md:px-0 ">
          <h1
            className="mb-6 text-2xl text-[#313D4F] break-words break-all leading-12"
            style={{ overflowWrap: "break-word" }}
          >
            {votingData?.title}
          </h1>
          {(votingData?.voteStatus || votingData?.voteStatus === 0) && (
            <div className="flex justify-between mb-6">
              <div className="flex items-center w-full mb-1 sm:mb-0">
                <VoteStatusBtn status={votingData.voteStatus} />

                <div className="flex items-center justify-center ml-[12px]">
                  <div className="text-[#4B535B] text-[14px]">
                    {t("content.createdby")}
                  </div>
                  <div className="ml-[8px] flex items-center justify-center bg-[#F5F5F5] rounded-full p-[5px]">
                    <img
                      className="w-[20px] h-[20px] rounded-full mr-[4px]"
                      src={userInfo?.img}
                      alt=""
                    />
                    <a
                      className="text-[#313D4F]"
                      target="_blank"
                      rel="noreferrer"
                      href={userInfo?.href}
                    >
                      {votingData?.githubName ||
                        EllipsisMiddle({
                          suffixCount: 4,
                          children: votingData?.address
                        })}
                    </a>
                  </div>
                </div>
              </div>
            </div>
          )}
          <div className="MDEditor">
            <MDEditor
              className="border-none rounded-[16px] bg-transparent"
              style={{ height: "auto" }}
              moreButton={true}
              value={votingData?.content}
              readOnly={true}
              view={{
                menu: false,
                md: false,
                html: true,
                both: false,
                fullScreen: true,
                hideMenu: false
              }}
              onChange={() => {
              }}
            />
          </div>
        </div>
      </div>
      <div className="w-full lg:w-4/12 lg:min-w-[321px]">
        <div className="mt-4 space-y-4 lg:mt-0">
          <div className="text-base border-solid border-y border-skin-border bg-skin-block-bg md:rounded-xl md:border">
            <div
              className="group flex h-[57px] justify-between rounded-t-none border-b border-skin-border border-solid px-4 pb-[12px] pt-3 md:rounded-t-lg">
              <h4 className="flex items-center font-medium">
                <div>{t("content.details")}</div>
              </h4>
            </div>
            <div className="p-4 leading-6 sm:leading-8">
              <div className="space-y-1 text-sm font-medium">
                <div className="flex justify-between">
                  <div>{t("content.startTime")}</div>
                  <span className="text-[#313D4F] text-sm font-normal">
                    {votingData?.startTime &&
                      dayjs(votingData.startTime * 1000).format(
                        "MMM. D, YYYY, h:mm A"
                      )}
                  </span>
                </div>
                <div className="flex justify-between">
                  <div>{t("content.endTime")}</div>
                  <span className="text-[#313D4F] text-sm font-normal">
                    {votingData?.endTime &&
                      dayjs(votingData.endTime * 1000).format(
                        "MMM. D, YYYY, h:mm A"
                      )}
                  </span>
                </div>
                <div className="flex justify-between">
                  <div>{t("content.timezone")}</div>
                  <span className="text-[#313D4F] text-sm font-normal">
                    {timezone}
                  </span>
                </div>
                <div className="flex justify-between">
                  <div className="text-sm font-medium">
                    {t("content.snapshotBlockHeight")}
                  </div>
                  <span className="text-[#313D4F] font-normal">
                    {formatNumberWithCommas({ value: votingData.snapshotInfo?.snapshotHeight }) || "-"}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
        {votingData?.voteStatus === IN_PROGRESS_STATUS && (
          <div className="mt-5">
            <div
              className="border-[#313D4F] mt-6 border-skin-border bg-skin-block-bg text-base md:rounded-xl md:border border-solid">
              <div
                className="group flex h-[57px] !border-[#eeeeee] justify-between items-center border-b px-4 pb-[12px] pt-3 border-solid">
                <h4 className="font-medium">{t("content.castVote")}</h4>
                <Popover
                  content={
                    <Table
                      rowKey={(record: any) => record.key}
                      dataSource={getPowerData(powerDetail)}
                      columns={columns}
                      pagination={false}
                    />
                  }
                >
                  <span className="text-[14px] text-[#273141] text-sm">
                    <InfoCircleOutlined style={{ fontSize: 14 }} />
                  </span>
                </Popover>
              </div>
              <div className="p-4 text-center">
                {options.map((item: ProposalOption, index: number) => {
                  return (
                    <div
                      className="mb-4 space-y-3 leading-10"
                      key={item.name + index}
                      onClick={() => {
                        handleOptionClick(index)
                      }}
                    >
                      <div
                        className={`w-full h-[45px] border-[#eeeeee] ${selectedOptionIndex === index
                          ? "border-[#0190FF] bg-[#F3FAFF]"
                          : ""
                          } hover:border-[#0190FF] flex justify-between items-center pl-8 pr-4 md:border border-solid rounded-full cursor-pointer`}
                      >
                        <div className="text-ellipsis h-[100%] overflow-hidden">
                          {item.name === "Approve"
                            ? t("content.approve")
                            : t("content.reject")}
                        </div>
                        {selectedOptionIndex === index && (
                          <svg
                            viewBox="0 0 24 24"
                            width="1.2em"
                            height="1.2em"
                            className="-ml-1 mr-2 text-md text-[#0190FF]"
                          >
                            <path
                              fill="none"
                              stroke="currentColor"
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth="2"
                              d="m5 13l4 4L19 7"
                            />
                          </svg>
                        )}
                      </div>
                    </div>
                  )
                })}
                <LoadingButton
                  text={t("content.vote")}
                  isFull={true}
                  loading={loading || writeContractPending}
                  handleClick={startVoting}
                />
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default Vote