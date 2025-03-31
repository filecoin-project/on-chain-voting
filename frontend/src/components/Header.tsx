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

import { SearchOutlined } from '@ant-design/icons';
import {
  ConnectButton,
  useConnectModal,
} from "@rainbow-me/rainbowkit";
import { Dropdown, Input, Modal } from 'antd';
import 'dayjs/locale/zh-cn';
import React, { useEffect, useRef, useState } from "react";
import Countdown from 'react-countdown';
import { useTranslation } from 'react-i18next';
import { Link, useLocation, useNavigate } from "react-router-dom";
import "tailwindcss/tailwind.css";
import { useAccount } from "wagmi";
import timezones from '../../public/json/timezons.json';
import {
  calibrationChainId,
  STORING_DATA_MSG
} from "../common/consts";
import { useVoterInfoSet } from "../common/hooks";
import {
  useCurrentTimezone,
  useFipList,
  useSearchValue,
  useVoterInfo
} from "../common/store";
import "../common/styles/reset.less";
import '../lang/config';

const Header = (props: any) => {
  const { changeLang } = props;
  // Destructure values from custom hooks
  const { chain, address, isConnected } = useAccount();
  const chainId = chain?.id || calibrationChainId;
  const prevAddressRef = useRef(address);
  const { openConnectModal } = useConnectModal();
  const navigate = useNavigate();
  const location = useLocation();

  // State variables
  const [expirationTime, setExpirationTime] = useState(0);
  const [modalOpen, setModalOpen] = useState(false);
  const [isFocus, setIsFocus] = useState<boolean>(false); // Determine whether the mouse has clicked on the search box
  // const [searchValue, setSearchValue] = useState<string>(); // Stores the value of the search box
  // Get the user's timezone
  const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  const text = timezones.find((item: any) => item.value === timezone)?.text;

  // Extract GMT offset from timezone
  const regex = /(?<=\().*?(?=\))/g;
  const GMTOffset = text?.match(regex);

  // Get voter information using custom hook
  const { voterInfo } = useVoterInfoSet(chainId, address);

  const { isFipEditorAddress } = useFipList((state: any) => state.data);

  // Update voter information in state
  const setVoterInfo = useVoterInfo((state: any) => state.setVoterInfo);
  const searchValue = useSearchValue((state: any) => state.searchValue);
  const setSearchValue = useSearchValue((state: any) => state.setSearchValue);


  // Update current timezone in state
  const setTimezone = useCurrentTimezone((state: any) => state.setTimezone);

  const { pathname } = useLocation();
  const { t } = useTranslation();



  useEffect(() => {
    window.scrollTo(0, 0);
  }, [pathname]);

  // Reload the page if address changes
  useEffect(() => {
    const prevAddress = prevAddressRef.current;
    if (address && prevAddress !== address) {
      window.location.reload();
    }
  }, [address]);

  // Update voter information when available
  useEffect(() => {
    if (voterInfo) {
      setVoterInfo(voterInfo);
    }
  }, [voterInfo]);

  // Set user's timezone based on GMT offset
  useEffect(() => {
    if (GMTOffset) {
      setTimezone(GMTOffset);
    }
  }, [GMTOffset])


  /**
   * Handle delegation action
   */
  const handleDelegate = async () => {
    // Prompt user to connect if not already connected
    if (!isConnected) {
      openConnectModal && openConnectModal();
      return;
    }

    // Retrieve gistStorageData from localStorage
    const gistStorageData = JSON.parse(localStorage.getItem('gistStorage') || '[]');
    const gistIndex = gistStorageData?.findIndex((item: any) => item.address === address);

    if (gistIndex > -1) {
      if (Date.now() < gistStorageData[gistIndex].timestamp) {
        setModalOpen(true);
        setExpirationTime(gistStorageData[gistIndex].timestamp);
        // Data has not expired
        return;
      } else {
        // Data has expired
        setExpirationTime(0);
        gistStorageData?.splice(gistIndex, 1);
        localStorage.setItem('gistStorage', JSON.stringify(gistStorageData));
      }
    }

    if (!voterInfo) {
      navigate('/gistDelegate/add');
      return
    }
    // Determine if the user has a GitHub account
    // const isGithubType = !!voterInfo[0];
    const isGithubType = true; //only display github
    if (voterInfo[2]) {
      if (isGithubType) {
        // Process GitHub data and navigate to appropriate page
        navigate('/gistDelegate/delete', {
          state: {
            params: {
              isGithubType,
              // aud,
              prf: ''
            }
          }
        });
      }
      else {
        // Process non-GitHub data and navigate to appropriate page
        // const decodeString = atob(data.split('.')[1]);
        // const payload = JSON.parse(decodeString);
        // const { aud, prf } = payload;
        navigate('/gistDelegate/delete', {
          state: {
            params: {
              isGithubType,
              // aud,
              // prf
            }
          }
        });
      }
    } else {
      // Navigate to add delegate page if no voter information is available
      navigate('/gistDelegate/add');
    }
  }

  const handleJump = (route: string) => {
    if (!isConnected) {
      openConnectModal && openConnectModal();
      return;
    }
    navigate(route);
  }

  const items: any = [
    {
      key: 'github',
      label: (
        <a
          onClick={handleDelegate}
        >
          {t('content.GithubDelegates')}
        </a>
      ),
    },
    {
      key: 'minerId',
      label: (
        <a
          onClick={() => { handleJump('/minerid') }}
        >
          {t('content.minerIDsManagement')}
        </a>
      ),
    },
    {
      key: 'fipEditorList',
      label: (
        <a
          onClick={() => { handleJump('/fip-editor/fipEditorList') }}
        >
          {t('content.fipEditorList')}
        </a>
      ),
    },
  ];

  if (isFipEditorAddress) {
    items.push({
      key: '3',
      label: t('content.fipEditorManagement'),
      children: [
        {
          key: '3-1',
          label: (
            <a
              onClick={() => { handleJump('/fip-editor/propose') }}
            >
              {t('content.createProposals')}
            </a>
          ),
        },
        {
          key: '3-2',
          label: (
            <a
              onClick={() => { handleJump('/fip-editor/approve') }}
            >
              {t('content.approveList')}
            </a>
          ),
        },
        {
          key: '3-3',
          label: (
            <a
              onClick={() => { handleJump('/fip-editor/revoke') }}
            >
              {t('content.revokeList')}
            </a>
          ),
        },
      ],
    })
  }

  const languageOptions = [
    { label: 'EN', value: 'en' },
    { label: '中文', value: 'zh' },
  ];

  const lang = localStorage.getItem("lang") || "en";
  const changeLanguage = (value: string) => {
    changeLang(value);
  };
  useEffect(() => {
    if (!chain || !isConnected) return
    setSearchValue('');
  }, [chain]);


  return (
    <>
      <header className='h-[96px] bg-[#ffffff] border-b border-solid border-[#DFDFDF]'>
        <div className='w-full h-[88px] flex items-center' style={{ justifyContent: "space-evenly" }}>
          <div className='flex items-center'>
            <div className='flex-shrink-0'>
              <Link to='/'>
                <img className="logo" src="/images/logo.png" alt="" />
              </Link>
            </div>
            <div className='ml-6 flex items-baseline space-x-20'>
              <Link
                to='/'
                className='text-black text-2xl font-semibold hover:opacity-80'
              >
                {t('content.powerVoting')}
              </Link>
            </div>
            {(location.pathname === '/home' || location.pathname === '/') &&
              <div className="ml-6">
                <Input
                  placeholder={t('content.searchProposals')}
                  size="large"
                  prefix={<SearchOutlined onClick={() => setSearchValue(searchValue)} className={`${isFocus ? "text-[#1677ff]" : "text-[#8b949e]"} text-xl hover:text-[#1677ff]`} />}
                  onClick={() => setIsFocus(true)}
                  onBlur={() => setIsFocus(false)}
                  onChange={(e) => setSearchValue(e.currentTarget.value)}
                  onPressEnter={() => setSearchValue(searchValue)}
                  value={searchValue}
                  className={`${isFocus ? 'w-[270px]' : "w-[180px]"} font-medium text-base item-center text-slate-800 bg-[#f7f7f7] rounded-lg`}
                />
              </div>
            }

          </div>
          <div className='flex items-center'>
            <Dropdown
              menu={{
                items,
              }}
              placement="bottomLeft"
              arrow
            >
              <button
                className="h-[40px] bg-sky-500 hover:bg-sky-700 text-white font-bold py-2 px-4 rounded-xl mr-4"
              >
                {t('content.tools')}
              </button>
            </Dropdown>
            <div className="connect flex items-center">
              <ConnectButton showBalance={false} label={t('content.connectWallet')} />
              <div className='px-4 py-2 h-full flex flex-nowrap text-sm'>
                {languageOptions.map((item) => {
                  return (
                    <div key={item.label} className={`h-full mr-1.5 cursor-pointer text-black font-semibold ${item.value === lang ? 'border-solid border-b-2 border-current' : 'border-none'}`} onClick={() => changeLanguage(item.value)}>
                      <div className='h-5 leading-6 text-center my-*'>{item.label}</div>
                    </div>
                  )
                })
                }
              </div>
            </div>
          </div>
          <Modal
            width={520}
            open={modalOpen}
            title={false}
            destroyOnClose={true}
            closeIcon={false}
            onCancel={() => { setModalOpen(false) }}
            footer={false}
            style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}
          >
            <p>{t(STORING_DATA_MSG)} {t('content.pleaseWait')}:&nbsp;
              <Countdown
                date={expirationTime}
                renderer={({ minutes, seconds, completed }) => {
                  if (completed) {
                    // Render a completed state
                    setModalOpen(false);
                  } else {
                    // Render a countdown
                    return <span>{minutes}:{seconds}</span>;
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
