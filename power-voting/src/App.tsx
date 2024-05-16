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

import React,{ useState } from "react";
import {useRoutes, useNavigate, Link} from "react-router-dom";
import axios from "axios";
import {
  ConnectButton,
  useConnectModal
} from "@rainbow-me/rainbowkit";
import { ConfigProvider, theme, Modal, Dropdown, FloatButton } from 'antd';
import { useAccount, useReadContract } from "wagmi";
import Countdown from 'react-countdown';
import routes from "./router";
import Footer from './components/Footer';
import "./common/styles/reset.less";
import "tailwindcss/tailwind.css";
import Loading from "./components/Loading";
import {STORING_DATA_MSG} from "./common/consts";
import oracleAbi from "./common/abi/oracle.json";
import {getContractAddress} from "./utils";

const App: React.FC = () => {
  const { chain, address, isConnected} = useAccount();

  const {openConnectModal} = useConnectModal();
  const navigate = useNavigate();
  const element = useRoutes(routes);
  const [spinning, setSpinning] = useState(false);
  const [expirationTime, setExpirationTime] = useState(0);
  const [modalOpen, setModalOpen] = useState(false);

  const { data: voterInfo, isSuccess } = useReadContract({
    // @ts-ignore
    address: getContractAddress(chain?.id || 0, 'oracle'),
    abi: oracleAbi,
    functionName: 'voterToInfo',
    args: [address]
  }) as any;

  const handleDelegate = async () => {
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

    setSpinning(true);
    if (!isConnected) {
      openConnectModal && openConnectModal();
      setSpinning(false);
      return;
    }
    setSpinning(false);
    const isGithubType = !!voterInfo[0];
    if (voterInfo[2]) {
      const { data } = await axios.get(`https://${voterInfo[2]}.ipfs.w3s.link/`);
      if (isGithubType) {
        const regex = /\/([^\/]+)\/([^\/]+)\/git\/blobs\/(\w+)/;
        const result = data.match(regex);
        const aud = result[1];
        navigate('/ucanDelegate/delete', { state: {
            params: {
              isGithubType,
              aud,
              prf: ''
            }
          }
        });
      } else {
        const decodeString = atob(data.split('.')[1]);
        const payload = JSON.parse(decodeString);
        const { aud, prf } = payload;
        navigate('/ucanDelegate/delete', { state: {
            params: {
              isGithubType,
              aud,
              prf
            }
          }
        });
      }
    } else {
      navigate('/ucanDelegate/add');
    }
  }

  const handleMinerId = () => {
    if (!isConnected) {
      openConnectModal && openConnectModal();
      setSpinning(false);
      return;
    }
    navigate('/minerid');
  }

  const items = [
    {
      key: 'ucan',
      label: (
        <a
          onClick={handleDelegate}
        >
          UCAN Delegates
        </a>
      ),
    },
    {
      key: 'minerId',
      label: (
        <a
          onClick={handleMinerId}
        >
          Miner IDs Management
        </a>
      ),
    },
  ];

  return (
    <ConfigProvider theme={{ algorithm: theme.darkAlgorithm }}>
      <div className="layout font-body">
        <header className='h-[96px] bg-[#273141]'>
          <div className='w-[1000px] h-[96px] mx-auto flex items-center justify-between'>
            <div className='flex items-center'>
              <div className='flex-shrink-0'>
                <Link to='/'>
                  <img className="logo" src="/images/logo.png" alt=""/>
                </Link>
              </div>
              <div className='ml-6 flex items-baseline space-x-20'>
                <Link
                  to='/'
                  className='text-white text-2xl font-semibold hover:opacity-80'
                >
                  Power Voting
                </Link>
              </div>
            </div>
            <div className='flex items-center'>
              {
                isSuccess &&
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
              }
              <div className="connect flex items-center">
                <ConnectButton />
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
        </header>
        <div className='content w-[1000px] mx-auto pt-10 pb-10'>
          {
            element
          }
        </div>
        <Footer/>
        <FloatButton.BackTop style={{ bottom: 100 }} />
      </div>
    </ConfigProvider>
  )
}

export default App
