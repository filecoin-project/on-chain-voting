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

import { message, Table } from "antd";
import dayjs from "dayjs";
import { useEffect, useRef, useState } from "react";
import { useAddresses } from "iso-filecoin-react";
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from "react-router-dom";
import { useGistList, useTransactionHash } from "../../../common/store";
import { useAccount, useWaitForTransactionReceipt } from "wagmi";
import axios from "axios";
import { getGistListApi, STORING_FAILED_MSG } from "../../../common/consts";
import Loading from "../../../components/Loading";
import { isFilAddress } from "../../../utils"
const GistDelegateList = () => {
    const { isConnected, address } = useAccount();
    const { address0x} = useAddresses({ address: address as string });
    const { t } = useTranslation();
    const navigate = useNavigate();
    const prevAddressRef = useRef(address);
    const gistList = useGistList((state: any) => state.gistList);
    const setGistList = useGistList((state: any) => state.setGistList)
    const storingHash = useTransactionHash((state: any) => state.transactionHash)
    const [messageApi, contextHolder] = message.useMessage();
    const [loading, setLoading] = useState<boolean>(false);
    const setStoringHash = useTransactionHash((state: any) => state.setTransactionHash)

    const { isSuccess, isFetched, isError } =
        useWaitForTransactionReceipt({
            hash: storingHash?.gistAudHash,
        })
    const columns = [
        {
            title: t('content.githubName'),
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
                    return dayjs(Number(value) * 1000).format('MMM.D, YYYY, h:mm A')
                }
            }
        },
    ];
    const getGistList = async () => {
        setLoading(true)
        const params = {
            address: isFilAddress(address!) && address0x.data ? address0x.data.toString() : address
        }
        const { data: { data: gistList } } = await axios.get(getGistListApi, { params })
        if (gistList.gistSigObj) {
            setGistList([
                {
                    githubName: gistList.gistSigObj.githubName,
                    walletAddress: gistList.gistSigObj.walletAddress,
                    timestamp: gistList.gistSigObj.timestamp
                }
            ])
        } else {
            navigate("/gistDelegate/add");
        }
        setLoading(false)
    }
    useEffect(() => {
        const prevAddress = prevAddressRef.current;
        if (prevAddress !== address || !isConnected) {
            navigate("/home");
        } else {
            getGistList()
        }
    }, [address, isConnected]);
    useEffect(() => {
        if (isFetched) {
            if (isSuccess) {
                getGistList();
                setStoringHash({ 'gistAudHash': undefined })
            }
            // If the transaction fails, show an error message
            if (isError) {
                messageApi.open({
                    type: 'error',
                    content: t(STORING_FAILED_MSG)
                })
                setStoringHash({ 'gistAudHash': undefined })
            }
        }
    }, [isFetched]);
    return (
        <div className="px-3 mb-6 md:px-0">
            {contextHolder}
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
            {loading ? <Loading /> :
                <div>
                    <div className='min-w-full bg-[#ffffff] text-left rounded-xl border-[1px] border-solid border-[#DFDFDF] overflow-hidden'>
                        <div className='flow-root space-y-4'>
                            <div className='font-normal text-black px-8 py-7 text-2xl border-b border-[#eeeeee] flex items-center'>
                                <span>{t('content.gitHubDelegatesList')}</span>
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
            }

        </div>
    )
}

export default GistDelegateList;
