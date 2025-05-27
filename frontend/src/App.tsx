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

import {
  lightTheme,
  RainbowKitProvider,
} from "@rainbow-me/rainbowkit";
import { useConnect as FilUseConnect, useAddresses } from "iso-filecoin-react"
import { ConfigProvider, FloatButton, theme } from 'antd';
import {
  useAccount as useFilAccount
} from 'iso-filecoin-react'
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';
import enUS from 'antd/locale/en_US';
import zhCN from 'antd/locale/zh_CN';
import { useTranslation } from 'react-i18next';
import { useEffect, useRef } from "react";
import { useLocation, useRoutes } from "react-router-dom";
import "tailwindcss/tailwind.css";
import { useAccount, useConnect } from "wagmi";
import timezones from './json/timezons.json';
import { calibrationChainId, getFipListApi } from "./common/consts"
import { useVoterInfoSet } from "./common/hooks"
import { useCurrentTimezone, useFipList, useVoterInfo } from "./common/store";
import "./common/styles/reset.less";
import Header from "./components/Header";
import Footer from './components/Footer';
import './lang/config';
import routes from "./router";
import axios from "axios";
import { isFilAddress } from "./utils"

const lang = localStorage.getItem("lang") || "en"
dayjs.locale(lang === 'en' ? lang : "zh-cn")

const App: React.FC = () => {
  // Destructure values from custom hooks
  const { chain, address, isConnected } = useAccount();
  const { adapter, state} = useFilAccount();
  const { connect, connectors } = useConnect()
  const { address0x} = useAddresses({ address: address as string });
  const { mutate: FilConnect, adapters } = FilUseConnect();
  const chainId = chain?.id || calibrationChainId;
  const prevAddressRef = useRef(address);
  const setFipList = useFipList((state: any) => state.setFipList)
  const { i18n } = useTranslation();
  const adapterId = window.localStorage.getItem('adapter');

  useEffect(()=>{
    if  (state === "connected" && isConnected){
      connect({connector:connectors.find(item => item.id === adapterId) || connectors[0]})
    }
  },[isConnected, state])

  useEffect(() => {
    if (isConnected && isFilAddress(address!)) {
      if (adapter) {
        FilConnect({ adapter: adapter });
      } else {
        FilConnect({ adapter: adapters.find(item => item.id === adapterId) || adapters[0] });
      }
    }
  }, [isConnected, adapter, address])

  // Render routes based on URL
  const element = useRoutes(routes);

  const isLanding = location.pathname === "/" || element?.props?.match?.route?.path === "*";

  // Get the user's timezone
  const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  const text = timezones.find((item: any) => item.value === timezone)?.text;

  // Extract GMT offset from timezone
  const regex = /(?<=\().*?(?=\))/g;
  const GMTOffset = text?.match(regex);

  // Get voter information using custom hook
  const { voterInfo } = useVoterInfoSet(chainId, address);

  // Update voter information in state
  const setVoterInfo = useVoterInfo((state: any) => state.setVoterInfo);

  // Update current timezone in state
  const setTimezone = useCurrentTimezone((state: any) => state.setTimezone);

  const { pathname } = useLocation();

  const handleChange = (value: string) => {
    i18n.changeLanguage(value);
    localStorage.setItem("lang", value);
    if (value === 'en') {
      dayjs.locale('en');
    } else if (value === 'zh') {
      dayjs.locale('zh-cn');
    }
  }

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [pathname]);

  // Reload the page if address changes
  useEffect(() => {
    const prevAddress = prevAddressRef.current;
    if (address && address.indexOf('0x') > 0 && prevAddress !== address) {
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
  const getFipList = async () => {
    const params = {
      chainId,
    }
    const { data: { data: fipList } } = await axios.get(getFipListApi, { params });
    if (isFilAddress(address!) && address0x.data) {
      setFipList(fipList, address0x.data.toString())
    } else {
      setFipList(fipList, address)
    }
  }
  useEffect(() => {
    if (!address || !chainId || !isConnected) {
      setFipList([], address);
      return
    }
    getFipList()
  }, [address, chainId, isConnected])

  const lang = localStorage.getItem("lang") || "en";

  return (
    <RainbowKitProvider
      locale={lang === "en" ? "en-US" : "zh-CN"}
      theme={lightTheme({
        accentColor: "#7b3fe4",
        accentColorForeground: "white",
      })}
      modalSize="compact"
    >
      <ConfigProvider theme={{
        algorithm: theme.defaultAlgorithm,
        components: {
          Radio: {
            buttonSolidCheckedBg: ''
          }
        }
      }} locale={lang === "en" ? enUS : zhCN}>
        <div className="layout font-body">
          {!isLanding && <Header changeLang={handleChange} />}
          <div className='content w-full mx-auto max-w-[1032px] pt-10 px-8 pb-24'>
            {
              element
            }
          </div>
          <Footer />
          <FloatButton.BackTop style={{ bottom: 100 }} />
        </div>
      </ConfigProvider>
    </RainbowKitProvider>
  )
}

export default App
