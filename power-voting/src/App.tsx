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

import React, { useState, useEffect, useRef } from "react";
import { useRoutes, useNavigate, Link, useLocation } from "react-router-dom";
import axios from "axios";
import {
  ConnectButton,
  useConnectModal
} from "@rainbow-me/rainbowkit";
import { ConfigProvider, theme, Modal, Dropdown, FloatButton } from 'antd';
import { useAccount } from "wagmi";
import Countdown from 'react-countdown';
import timezones from '../public/json/timezons.json';
import routes from "./router";
import Footer from './components/Footer';
import "./common/styles/reset.less";
import "tailwindcss/tailwind.css";
import { STORING_DATA_MSG } from "./common/consts";
import { useVoterInfo, useCurrentTimezone } from "./common/store";
import { useCheckFipEditorAddress, useVoterInfoSet } from "./common/hooks";

const App: React.FC = () => {
  // Destructure values from custom hooks
  const { chain, address, isConnected } = useAccount();
  const chainId = chain?.id || 0;
  const prevAddressRef = useRef(address);
  const { openConnectModal } = useConnectModal();
  const navigate = useNavigate();

  // Render routes based on URL
  const element = useRoutes(routes);

  const location = useLocation();
  const isLanding = location.pathname === "/" || element?.props?.match?.route?.path === "*"
  // State variables
  const [expirationTime, setExpirationTime] = useState(0);
  const [modalOpen, setModalOpen] = useState(false);

  // Get the user's timezone
  const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  const text = timezones.find((item: any) => item.value === timezone)?.text;

  // Extract GMT offset from timezone
  const regex = /(?<=\().*?(?=\))/g;
  const GMTOffset = text?.match(regex);

  // Get voter information using custom hook
  const { voterInfo } = useVoterInfoSet(chainId, address);

  const { isFipEditorAddress } = useCheckFipEditorAddress(chainId, address);

  // Update voter information in state
  const setVoterInfo = useVoterInfo((state: any) => state.setVoterInfo);

  // Update current timezone in state
  const setTimezone = useCurrentTimezone((state: any) => state.setTimezone);

  const { pathname } = useLocation();
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
    // Retrieve ucanStorageData from localStorage
    const ucanStorageData = JSON.parse(localStorage.getItem('ucanStorage') || '[]');
    const ucanIndex = ucanStorageData?.findIndex((item: any) => item.address === address);
    if (ucanIndex > -1) {
      if (Date.now() < ucanStorageData[ucanIndex].timestamp) {
        setModalOpen(true);
        setExpirationTime(ucanStorageData[ucanIndex].timestamp);
        // Data has not expired
        return;
      } else {
        // Data has expired
        setExpirationTime(0);
        ucanStorageData?.splice(ucanIndex, 1);
        localStorage.setItem('ucanStorage', JSON.stringify(ucanStorageData));
      }
    }

    // Prompt user to connect if not already connected
    if (!isConnected) {
      openConnectModal && openConnectModal();
      return;
    }


    if (!voterInfo) {
      navigate('/ucanDelegate/add');
      return
    }
    // Determine if the user has a GitHub account
    const isGithubType = !!voterInfo[0];
    if (voterInfo[2]) {
      // Fetch data from IPFS using voter's identifier
      const { data } = await axios.get(`https://${voterInfo[2]}.ipfs.w3s.link/`);
      if (isGithubType) {
        // Process GitHub data and navigate to appropriate page
        const regex = /\/([^/]+)\/([^/]+)\/git\/blobs\/(\w+)/;
        const result = data.match(regex);
        const aud = result[1];
        navigate('/ucanDelegate/delete', {
          state: {
            params: {
              isGithubType,
              aud,
              prf: ''
            }
          }
        });
      }
      else {
        navigate('/ucanDelegate/add');
        // Process non-GitHub data and navigate to appropriate page
        // const decodeString = atob(data.split('.')[1]);
        // const payload = JSON.parse(decodeString);
        // const { aud, prf } = payload;
        // navigate('/ucanDelegate/delete', { state: {
        //     params: {
        //       isGithubType,
        //       aud,
        //       prf
        //     }
        //   }
        // });
      }
    } else {
      // Navigate to add delegate page if no voter information is available
      navigate('/ucanDelegate/add');
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
      key: 'ucan',
      label: (
        <a
          onClick={handleDelegate}
        >
          Connect GitHub
        </a>
      ),
    },
    {
      key: 'minerId',
      label: (
        <a
          onClick={() => { handleJump('/minerid') }}
        >
          Miner IDs Management
        </a>
      ),
    },
  ];

  if (isFipEditorAddress) {
    items.push({
      key: '3',
      label: 'FIP Editor Management',
      children: [
        {
          key: '3-1',
          label: (
            <a
              onClick={() => { handleJump('/fip-editor/propose') }}
            >
              Propose
            </a>
          ),
        },
        {
          key: '3-2',
          label: (
            <a
              onClick={() => { handleJump('/fip-editor/approve') }}
            >
              Approve
            </a>
          ),
        },
        {
          key: '3-3',
          label: (
            <a
              onClick={() => { handleJump('/fip-editor/revoke') }}
            >
              Revoke
            </a>
          ),
        },
      ],
    })
  }

  return (
    <ConfigProvider theme={{ algorithm: theme.defaultAlgorithm }}>
      <div className="layout font-body">
        {!isLanding && <header className='h-[96px] bg-[#ffffff]'>
          <div className='w-[1000px] h-[88px] mx-auto flex items-center justify-between'>
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
                  Power Voting
                </Link>
              </div>
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
                  Tools
                </button>
              </Dropdown>
              <div className="connect flex items-center">
                <ConnectButton showBalance={false} />
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
              <p>{STORING_DATA_MSG} Please wait:&nbsp;
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
        </header>}
        <div className='content w-[1000px] mx-auto pt-10 pb-10'>
          {
            element
          }
        </div>
        <Footer />
        <FloatButton.BackTop style={{ bottom: 100 }} />
      </div>
    </ConfigProvider>
  )
}

export default App
