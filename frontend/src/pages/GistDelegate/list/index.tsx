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
import dayjs from "dayjs";
import { useEffect, useRef } from "react";

import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from "react-router-dom";
import { useGistList } from "../../../common/store";
import { useAccount } from "wagmi";
const GistDelegateList = () => {
    const { isConnected, address } = useAccount();
    const { t } = useTranslation();
    const navigate = useNavigate();
    const prevAddressRef = useRef(address);
    const gistList = useGistList((state: any) => state.gistList);
    const columns = [
        {
            title: 'GithubName',
            dataIndex: 'githubName',
            key: 'githubName',
            width: 280,
        },
        {
            title: t('content.address'),
            dataIndex: 'walletAddress',
            key: 'walletAddress',
            width: 280,
        },
        {
            title: t('content.time'),
            dataIndex: 'timestamp',
            key: 'timestamp',
            width: 280,
            render: (value: string) => {
                if (value) {
                    return dayjs(Number(value) * 1000).toString()
                }
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
                        <span>GitHub Delegates List</span>
                    </div>
                    <div className='px-8 pb-4 !mt-0'>
                        <Table
                            className='mb-4'
                            rowKey={(record: any) => record.editor}
                            dataSource={gistList}
                            columns={columns}
                            pagination={false}
                        />
                    </div>
                </div>
            </div>
            <div className="text-center mt-10">
                <button
                    className={`h-[40px] bg-sky-500 hover:bg-sky-700 text-white py-2 px-6 rounded-xl`}
                    type='button' onClick={() => { navigate('/gistDelegate/add') }}>
                    {t('content.authorizeAgain')}
                </button>
            </div>
        </div>
    )
}

export default GistDelegateList;
