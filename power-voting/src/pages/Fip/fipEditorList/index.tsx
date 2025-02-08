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

import { Table } from "antd";
import React, { useEffect, useRef } from "react";
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from "react-router-dom";
import { useAccount } from "wagmi";
import {
  calibrationChainId,
  web3AvatarUrl
} from "../../../common/consts";
import { useFipEditors } from "../../../common/hooks";
const FipEditorList = () => {
  const { isConnected, address, chain } = useAccount();
  const { t } = useTranslation();
  const chainId = chain?.id || calibrationChainId;
  const navigate = useNavigate();
  const prevAddressRef = useRef(address);

  const { fipEditors } = useFipEditors(chainId);

  const columns = [
    {
      title: t('content.FIPEditor'),
      dataIndex: 'address',
      key: 'address',
      width: 280,
      render: (record: any, value: string) => {
        return (
          <div className="w-[180px] flex items-center">
            <img className="w-[20px] h-[20px] rounded-full mr-2" src={`${web3AvatarUrl}:${value}`} alt="" />
            <a
              className="text-black hover:text-blue"
              target="_blank"
              rel="noopener"
              href={`${chain?.blockExplorers?.default.url}/address/${value}`}
            >
              {value}
            </a>
          </div>
        )
      }
    },
  ];

  useEffect(() => {
    const prevAddress = prevAddressRef.current;
    if (prevAddress !== address || !isConnected) {
      navigate("/home");
    }
  }, [address, isConnected]);



  return (
    <div className="px-3 mb-6 md:px-0">
      <button>
        <div className="inline-flex items-center mb-8 gap-1 text-skin-text hover:text-skin-link">
          <Link to="/home" className="flex items-center">
            <svg className="mr-1" viewBox="0 0 24 24" width="1.2em" height="1.2em">
              <path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                d="m11 17l-5-5m0 0l5-5m-5 5h12" />
            </svg>
            {t('content.back')}
          </Link>
        </div>
      </button>
      <div className='min-w-full bg-[#ffffff] text-left rounded-xl border-[1px] border-solid border-[#DFDFDF] overflow-hidden'>
        <div className='flow-root space-y-4'>
          <div className='font-normal text-black px-8 py-7 text-2xl border-b border-[#eeeeee] flex items-center'>
            <span>{t('content.fipEditorList')}</span>
          </div>
          <div className='px-8 pb-4 !mt-0'>
            <Table
              className='mb-4'
              rowKey={(record: any) => record.proposalId}
              dataSource={fipEditors}
              columns={columns}
              pagination={false}
            />
          </div>
        </div>
      </div>
    </div>
  )
}

export default FipEditorList;
