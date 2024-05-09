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
import { ConfigProvider, theme, Modal, Dropdown } from 'antd';
import { useAccount } from "wagmi";
import Countdown from 'react-countdown';
import routes from "./router";
import Footer from './components/Footer';
import "./common/styles/reset.less";
import "tailwindcss/tailwind.css";
import {useStaticContract} from "./hooks";
import Loading from "./components/Loading";
import {STORING_DATA_MSG} from "./common/consts";


const App: React.FC = () => {
  const { chain, address, isConnected} = useAccount();
  const {openConnectModal} = useConnectModal();
  const navigate = useNavigate();
  const element = useRoutes(routes);
  const [showButton, setShowButton] = useState(false);
  const [spinning, setSpinning] = useState(false);
  const [expirationTime, setExpirationTime] = useState(0);
  const [modalOpen, setModalOpen] = useState(false);

  const scrollRef = useRef(null);

  const scrollToTop = () => {
    const element = scrollRef.current;
    // @ts-ignore
    element.scrollTo({
      top: 0,
      behavior: 'smooth'
    });
  };

  useEffect(() => {
    const handleScroll = () => {
      const element = scrollRef.current;
      // @ts-ignore
      setShowButton(element.scrollTop > 300)
    };

    // @ts-ignore
    scrollRef.current.addEventListener('scroll', handleScroll);

    return () => {
      // @ts-ignore
      scrollRef.current.removeEventListener('scroll', handleScroll);
    };
  }, []);


  useEffect(() => {
    scrollToTop();
  }, [location]);

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
    const chainId = chain?.id || 0;
    const { getOracleAuthorize }  = await useStaticContract(chainId);
    const { data: { githubAccount, ucanCid } } = await getOracleAuthorize(address);
    setSpinning(false);
    const isGithubType = !!githubAccount;
    if (ucanCid) {
      const { data } = await axios.get(`https://${ucanCid}.ipfs.w3s.link/`);
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
      key: '1',
      label: (
        <a
          onClick={handleDelegate}
        >
          UCAN Delegates
        </a>
      ),
    },
    {
      key: '2',
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
      <div className="layout font-body" id='scrollBox' ref={scrollRef}>
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
        <button onClick={scrollToTop} className={`${showButton ? '' : 'hidden'} fixed bottom-[6rem] right-[6rem] z-40  p-2 bg-gray-600 text-white rounded-full hover:bg-gray-700 focus:outline-none`}>
          <svg className="w-6 h-6" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
            <path d="M12 3l-8 8h5v10h6V11h5z" fill="currentColor" />
          </svg>
        </button>

      </div>
    </ConfigProvider>
  )
}

export default App
