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

import { SearchOutlined } from "@ant-design/icons"
import {
  ConnectButton,
  useConnectModal
} from "@rainbow-me/rainbowkit"
import { useAddresses } from 'iso-filecoin-react'
import { Dropdown, Input, Modal, Typography } from "antd"
import "dayjs/locale/zh-cn"
import { useEffect, useState } from "react"
import Countdown from "react-countdown"
import { useTranslation } from "react-i18next"
import { Link, useLocation, useNavigate } from "react-router-dom"
import "tailwindcss/tailwind.css"
import { useAccount } from "wagmi"
import timezones from "../json/timezons.json"
import {
  calibrationChainId,
  getGistListApi, network,
  STORING_DATA_MSG
} from "../common/consts"
import { useVoterInfoSet } from "../common/hooks";
import {
  useCurrentTimezone,
  useFipList,
  useGistList,
  useSearchValue,
  useVoterInfo
} from "../common/store"
import "../common/styles/reset.less"
import "../lang/config"
import axios from "axios"
import { getBlockExplorers, isFilAddress } from "../utils"

const Header = (props: any) => {
  const { changeLang } = props
  // Destructure values from custom hooks
  const { chain, address, isConnected } = useAccount();
  const { address0x } = useAddresses({ address: address as string })
  const chainId = chain?.id || calibrationChainId
  const { openConnectModal } = useConnectModal()
  const navigate = useNavigate()
  const location = useLocation()
  // State variables
  // const [expirationTime, setExpirationTime] = useState(0);
  const [modalOpen, setModalOpen] = useState(false)
  const [isFocus, setIsFocus] = useState<boolean>(false) // Determine whether the mouse has clicked on the search box
  // const [searchValue, setSearchValue] = useState<string>(); // Stores the value of the search box
  // Get the user's timezone
  const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone
  const text = timezones.find((item: any) => item.value === timezone)?.text

  // Extract GMT offset from timezone
  const regex = /(?<=\().*?(?=\))/g
  const GMTOffset = text?.match(regex)

  // Get voter information using custom hook
  const { voterInfo } = useVoterInfoSet(chainId, address)

  const { isFipEditorAddress } = useFipList((state: any) => state.data)

  // Update voter information in state
  const setVoterInfo = useVoterInfo((state: any) => state.setVoterInfo)
  const searchValue = useSearchValue((state: any) => state.searchValue)
  const setSearchValue = useSearchValue((state: any) => state.setSearchValue)
  const setGistList = useGistList((state: any) => state.setGistList)
  // Update current timezone in state
  const setTimezone = useCurrentTimezone((state: any) => state.setTimezone)

  const { pathname } = useLocation()
  const { t } = useTranslation()

  useEffect(() => {
    window.scrollTo(0, 0)
  }, [pathname])

  // Update voter information when available
  useEffect(() => {
    if (voterInfo) {
      setVoterInfo(voterInfo)
    }
  }, [voterInfo])

  // Set user's timezone based on GMT offset
  useEffect(() => {
    if (GMTOffset) {
      setTimezone(GMTOffset)
    }
  }, [GMTOffset])


  /**
   * Handle delegation action
   */
  const handleDelegate = async () => {
    // Prompt user to Connect if not already connected
    if (!isConnected) {
      openConnectModal && openConnectModal()
      return
    }
    const params = {
      address: isFilAddress(address!) && address0x.data ? address0x.data.toString() : address
    }
    try {
      const { data: { data: gistList } } = await axios.get(getGistListApi, { params })
      if (gistList?.gistSigObj) {
        setGistList([
          {
            githubName: gistList?.gistSigObj?.githubName,
            walletAddress: gistList?.gistSigObj?.walletAddress,
            timestamp: gistList?.gistSigObj?.timestamp
          }
        ])
        navigate("/gistDelegate/list")
      } else {
        setGistList([])
        navigate("/gistDelegate/add")
      }
    } catch (e) {
      navigate("/gistDelegate/add")
    }
    return
  }

  const handleJump = (route: string) => {
    if (!isConnected) {
      openConnectModal && openConnectModal()
      return
    }
    navigate(route)
  }

  const items: any = [
    {
      key: "github",
      label: (
        <a
          onClick={handleDelegate}
        >
          {t("content.GithubDelegates")}
        </a>
      )
    },
    {
      key: "minerId",
      label: (
        <a
          onClick={() => handleJump("/minerid")}
        >
          {t("content.minerIDsManagement")}
        </a>
      )
    },
    {
      key: "fipEditorList",
      label: (
        <a
          onClick={() => {
            handleJump("/fip-editor/fipEditorList")
          }}
        >
          {t("content.fipEditorList")}
        </a>
      )
    }
  ]

  if (isFipEditorAddress) {
    items.push({
      key: "3",
      label: t("content.fipEditorManagement"),
      children: [
        {
          key: "3-1",
          label: (
            <a
              onClick={() => {
                handleJump("/fip-editor/propose")
              }}
            >
              {t("content.createProposals")}
            </a>
          )
        },
        {
          key: "3-2",
          label: (
            <a
              onClick={() => {
                handleJump("/fip-editor/approve")
              }}
            >
              {t("content.approveList")}
            </a>
          )
        },
        {
          key: "3-3",
          label: (
            <a
              onClick={() => {
                handleJump("/fip-editor/revoke")
              }}
            >
              {t("content.revokeList")}
            </a>
          )
        }
      ]
    })
  }

  const languageOptions = [
    { label: "EN", value: "en" },
    { label: "中文", value: "zh" }
  ]

  const lang = localStorage.getItem("lang") || "en"
  const changeLanguage = (value: string) => {
    changeLang(value)
  }
  useEffect(() => {
    if (!chain || !isConnected) return
    setSearchValue("")
  }, [chain])


  return (
    <>
      <header className="h-[96px] bg-[#ffffff] border-b border-solid border-[#DFDFDF]">
        <div className="w-full h-[88px] flex items-center" style={{ justifyContent: "space-evenly" }}>
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <Link to="/">
                <img className="logo" src="/images/logo.png" alt="" />
              </Link>
            </div>
            <div className="ml-6 flex items-baseline space-x-20">
              <Link
                to="/"
                className="text-black text-2xl font-semibold hover:opacity-80"
              >
                {t("content.powerVoting")}
              </Link>
            </div>
            {(location.pathname === "/home" || location.pathname === "/") &&
              <div className="ml-6">
                <Input
                  placeholder={t("content.searchProposals")}
                  size="large"
                  prefix={<SearchOutlined onClick={() => setSearchValue(searchValue)}
                    className={`${isFocus ? "text-[#1677ff]" : "text-[#8b949e]"} text-xl hover:text-[#1677ff]`} />}
                  onClick={() => setIsFocus(true)}
                  onBlur={() => setIsFocus(false)}
                  onChange={(e) => setSearchValue(e.currentTarget.value)}
                  onPressEnter={() => setSearchValue(searchValue)}
                  value={searchValue}
                  className={`${isFocus ? "w-[270px]" : "w-[180px]"} font-medium text-base item-center text-slate-800 bg-[#f7f7f7] rounded-lg`}
                />
              </div>
            }

          </div>
          <div className="flex items-center">
            <Dropdown
              menu={{
                items
              }}
              placement="bottomLeft"
              arrow
            >
              <button
                className="h-[40px] bg-primary hover:bg-primary/80 text-black font-bold py-2 px-4 rounded-xl mr-4"
              >
                {t("content.tools")}
              </button>
            </Dropdown>
            <div className="connect flex items-center justify-center">
              <ConnectButton showBalance={false} label={t("content.connectWallet")} />
              {
                address0x.data && isFilAddress(address!) && <a
                  target='_blank'
                  rel="noopener noreferrer"
                  href={getBlockExplorers(chain, address!)}
                  className="ml-4 py-2 text-[#25292E] text-[14px] font-[700] flex items-center">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"

                    width={props.width ?? 24}
                    height={props.height ?? 24}
                    viewBox="0 0 24 24"
                  >
                    <title>Ethereum Token</title>
                    <path
                      fill="currentColor"
                      d="M12 3v6.652l5.625 2.516zm0 0l-5.625 9.166L12 9.652zm0 13.478V21l5.625-7.785zM12 21v-4.522l-5.625-3.263z"
                    />
                    <path
                      fill="currentColor"
                      d="m12 15.43l5.625-3.263L12 9.652zm-5.625-3.263L12 15.43V9.652z"
                    />
                    <path
                      fill="currentColor"
                      fillRule="evenodd"
                      d="m12 15.43l-5.625-3.262L12 3l5.625 9.166zm-5.25-3.528l5.162-8.41v6.115zm-.077.229l5.239-2.327v5.364zm5.418-2.327v5.364l5.233-3.037zm0-.197l5.162 2.295l-5.162-8.41z"
                      clipRule="evenodd"
                    />
                    <path
                      fill="currentColor"
                      fillRule="evenodd"
                      d="m12 16.407l-5.625-3.195L12 21l5.625-7.789zm-4.995-2.633l4.906 2.79v4.005zm5.085 2.79v4.005l4.904-6.795z"
                      clipRule="evenodd"
                    />
                  </svg>
                  <Typography.Paragraph style={{ margin: 0 }} copyable={{ text: address0x.data.toString() }}>({address0x.data.toString().substring(0, 4)}...{address0x.data.toString().substring(38)})</Typography.Paragraph>
                </a>
              }
              <div className="px-4 py-2 h-full flex flex-nowrap text-sm">
                {languageOptions.map((item) => {
                  return (
                    <div key={item.label}
                      className={`h-full mr-1.5 cursor-pointer text-black font-semibold ${item.value === lang ? "border-solid border-b-2 border-current" : "border-none"}`}
                      onClick={() => changeLanguage(item.value)}>
                      <div className="h-5 leading-6 text-center my-*">{item.label}</div>
                    </div>
                  )
                })
                }
              </div>

              <a
                target='_blank'
                rel="noopener noreferrer"
                href={network === 'testnet' ? "https://vote.fil.org/" : "https://vote.storswift.io/"}
                className="py-2 text-[#0000FF] text-[14px] flex items-center">
                {network === 'testnet' ? 'Mainnet↗' : 'Calibration↗'}

              </a>
            </div>
          </div>
          <Modal
            width={520}
            open={modalOpen}
            title={false}
            destroyOnClose={true}
            closeIcon={false}
            onCancel={() => {
              setModalOpen(false)
            }}
            footer={false}
            style={{ display: "flex", alignItems: "center", justifyContent: "center" }}
          >
            <p>{t(STORING_DATA_MSG)} {t("content.pleaseWait")}:&nbsp;
              <Countdown
                date={0}
                renderer={({ minutes, seconds, completed }) => {
                  if (completed) {
                    // Render a completed state
                    setModalOpen(false)
                  } else {
                    // Render a countdown
                    return <span>{minutes}:{seconds}</span>
                  }
                }}
              />
            </p>
          </Modal>
        </div>
      </header>
    </>
  )
}

export default Header