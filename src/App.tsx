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

import React,{ useState, useEffect, useRef } from "react";
import {useRoutes, useNavigate, Link} from "react-router-dom";
import axios from "axios";
import {
  ConnectButton,
  useConnectModal
} from "@rainbow-me/rainbowkit";
import "@rainbow-me/rainbowkit/styles.css";
import { ConfigProvider, theme, Modal } from 'antd';
import { useNetwork, useAccount } from "wagmi";
import Countdown from 'react-countdown';
import routes from "./router";
import Footer from './components/Footer';
import "@rainbow-me/rainbowkit/styles.css";
import "./common/styles/reset.less";
import "tailwindcss/tailwind.css";
import {useStaticContract} from "./hooks";
import Loading from "./components/Loading";
import {STORING_DATA_MSG} from "./common/consts";


const App: React.FC = () => {
  const { chain } = useNetwork();
  const { address, isConnected} = useAccount();
  const {openConnectModal} = useConnectModal();
  const navigate = useNavigate();
  const element = useRoutes(routes);
  const [spinning, setSpinning] = useState(false);
  const [expirationTime, setExpirationTime] = useState(0);
  const [modalOpen, setModalOpen] = useState(false);

  const prevAddressRef = useRef(address);

  useEffect(() => {
    const prevAddress = prevAddressRef.current;
    if (prevAddress !== address) {
      window.location.reload();
    }
  }, [address]);

  const handleDelegate = async () => {
    const ucanStorageData = JSON.parse(localStorage.getItem('ucanStorage') || '[]');
    const ucanIndex = ucanStorageData.findIndex((item: any) => item.address === address);
    if (ucanIndex > -1) {
      if (Date.now() < ucanStorageData[ucanIndex].timestamp) {
        setModalOpen(true);
        setExpirationTime(ucanStorageData[ucanIndex].timestamp);
        // Data has not expired
        return;
      } else {
        // Data has expired
        ucanStorageData.splice(ucanIndex, 1);
        localStorage.setItem('ucanStorage', JSON.stringify(ucanStorageData));
        setExpirationTime(0);
      }
    }

    setSpinning(true);
    if (!isConnected) {
      openConnectModal && openConnectModal();
      setSpinning(false);
      return;
    }
    const chainId = chain?.id || 0;
    const { getOracleAuthorize }  = await useStaticContract(chainId);
    const { data: { githubAccount, ucanCid } } = await getOracleAuthorize(address);
    setSpinning(false);
    const isGithubType = !!githubAccount;
    if (ucanCid) {
      const { data } = await axios.get(`https://${ucanCid}.ipfs.nftstorage.link/`);
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

  setTimeout(() => {
    const elementToRemove = document.getElementById('okx-inject');
    elementToRemove?.remove();
  }, 3000);

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
            <button
              className="h-[40px] bg-sky-500 hover:bg-sky-700 text-white font-bold py-2 px-4 rounded-xl mr-4"
              onClick={handleDelegate}
            >
              UCAN Delegates
            </button>
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
          spinning ? <Loading /> : element
        }
      </div>
      <Footer/>
    </div>
    </ConfigProvider>
  )
}

export default App
