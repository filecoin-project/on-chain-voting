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

import { InfoCircleOutlined } from '@ant-design/icons';
import { Button, Pagination, Popconfirm, Popover, Row, Table, Tooltip, message } from "antd";
import axios from "axios";
import React, { useEffect, useRef, useState } from "react";
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from "react-router-dom";
import type { BaseError } from "wagmi";
import { useAccount, useWaitForTransactionReceipt, useWriteContract } from "wagmi";
import fileCoinAbi from "../../../common/abi/power-voting.json";
import {
  HAVE_APPROVED_MSG,
  STORING_DATA_MSG,
  web3AvatarUrl,
} from "../../../common/consts";
import { useApproveProposalId, useCheckFipEditorAddress, useFipEditorProposalDataSet, useFipEditors } from "../../../common/hooks";
import EllipsisMiddle from "../../../components/EllipsisMiddle";
import Loading from "../../../components/Loading";
import { getContractAddress } from "../../../utils";
import "./index.less";
const FipEditorApprove = () => {
  const { isConnected, address, chain } = useAccount();
  const { t } = useTranslation();
  const chainId = chain?.id || 0;
  const navigate = useNavigate();
  const prevAddressRef = useRef(address);
  const [messageApi, contextHolder] = message.useMessage();

  const [fipProposalList, setFipProposalList] = useState<any[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(5);
  const [loading, setLoading] = useState(false);
  const [currentProposalId, setCurrentProposalId] = useState(null);

  const { isFipEditorAddress, checkFipEditorAddressSuccess } = useCheckFipEditorAddress(chainId, address);
  const { fipEditors } = useFipEditors(chainId);
  const { approveProposalId, getApproveProposalLoading } = useApproveProposalId(chainId);

  const { fipEditorProposalData, getFipEditorProposalIdLoading, getFipEditorProposalIdSuccess, error } = useFipEditorProposalDataSet({
    chainId,
    idList: approveProposalId,
    page,
    pageSize,
  });

  const {
    data: hash,
    writeContract,
    error: writeContractError,
    isPending: writeContractPending,
    isSuccess: writeContractSuccess,
    reset
  } = useWriteContract();

  const { isLoading: transactionLoading } =
    useWaitForTransactionReceipt({
      hash,
    })

  const isLoading = loading || writeContractPending || transactionLoading;

  const popoverColumns = [
    {
      title: t('content.FIPEditor'),
      dataIndex: 'address',
      key: 'address',
      width: 280,
      render: (value: string) => {
        return (
          <div className="w-[180px] flex items-center">
            <img className="w-[20px] h-[20px] rounded-full mr-2" src={`${web3AvatarUrl}:${value}`} alt="" />
            <a
              className="text-black hover:text-black"
              target="_blank"
              rel="noopener"
              href={`${chain?.blockExplorers?.default.url}/address/${value}`}
            >
              {EllipsisMiddle({ suffixCount: 8, children: value })}
            </a>
          </div>
        )
      }
    },
    {
      title: t('content.status'),
      dataIndex: 'status',
      key: 'status',
      render: (value: string) => value || '-'
    },
  ]

  const columns = [
    {
      title: t('content.address'),
      dataIndex: 'address',
      key: 'address',
      width: 280,
      render: (value: string) => {
        return (
          <div className="w-[180px] flex items-center">
            <img className="w-[20px] h-[20px] rounded-full mr-2" src={`${web3AvatarUrl}:${value}`} alt="" />
            <a
              className="text-black hover:text-black"
              target="_blank"
              rel="noopener"
              href={`${chain?.blockExplorers?.default.url}/address/${value}`}
            >
              {EllipsisMiddle({ suffixCount: 8, children: value })}
            </a>
          </div>
        )
      }
    },
    {
      title: <div><div>{t('content.info')}</div></div>,
      dataIndex: 'info',
      key: 'info',
      ellipsis: { showTitle: false },
      render: (value: string) => {
        return (
          value ? <Tooltip overlayClassName="custom-tooltip" color="#ffffff" placement="topLeft" title={value}>
            {value}
          </Tooltip> : '-'
        )
      }
    },
    {
      title: t('content.approveRatio'),
      dataIndex: 'ratio',
      key: 'ratio',
      render: (value: string, record: any) => {
        return (
          <div className='flex items-center gap-2'>
            <span>{value}</span>
            <Popover content={
              <Table
                rowKey={(record: any) => record.address}
                dataSource={record.voteList}
                columns={popoverColumns}
                pagination={false}
              />
            }>
              <InfoCircleOutlined style={{ fontSize: 14, cursor: 'pointer' }} />
            </Popover>
          </div>
        )
      }
    },
    {
      title: t('content.action'),
      key: 'total',
      align: 'center' as const,
      width: 120,
      render: (_: any, record: any) =>
        <Popconfirm
          title={t('content.approveFIPEditor')}
          description={t('content.isConfirmApprove')}
          onConfirm={() => { confirm(record) }}
          okText={t('content.yes')}
          cancelText={t('content.no')}
        >
          <Button type='primary' className='w-[80px] h-[24px] flex justify-center items-center' loading={record.proposalId === currentProposalId && isLoading} >{t('content.approve')}</Button>
        </Popconfirm>
    },
  ];

  useEffect(() => {
    if (!isConnected || (checkFipEditorAddressSuccess && !isFipEditorAddress)) {
      navigate("/home");
      return;
    }
  }, [isConnected, checkFipEditorAddressSuccess, isFipEditorAddress]);

  useEffect(() => {
    const prevAddress = prevAddressRef.current;
    if (prevAddress !== address) {
      navigate("/home");
    }
  }, [address]);

  useEffect(() => {
    if (error) {
      messageApi.open({
        type: 'error',
        content: (error as BaseError)?.shortMessage || error?.message,
      });
    }
  }, [error]);

  useEffect(() => {
    if (writeContractError) {
      messageApi.open({
        type: 'error',
        content: (writeContractError as BaseError)?.shortMessage || writeContractError?.message,
      });
    }
    reset();
  }, [writeContractError]);

  useEffect(() => {
    if (getFipEditorProposalIdSuccess) {
      initState();
    }
  }, [getFipEditorProposalIdSuccess]);

  useEffect(() => {
    if (isConnected && !loading && !getApproveProposalLoading && !getFipEditorProposalIdLoading) {
      initState();
    }
  }, [chain, page, address]);

  useEffect(() => {
    if (writeContractSuccess) {
      messageApi.open({
        type: 'success',
        content: t(STORING_DATA_MSG),
      });
      setTimeout(() => {
        navigate("/home")
      }, 1000);
    }
  }, [writeContractSuccess]);

  const initState = async () => {
    setLoading(true);
    const list: any = [];
    await Promise.all(fipEditorProposalData.map(async (item: any) => {
      try {
        const { result } = item;
        const obj = {
          proposalId: result[0],
          fipEditorAddress: result[1],
          voterInfoCid: result[2],
          voters: result[3],
        }
        const url = `https://${obj.voterInfoCid}.ipfs.w3s.link/`;
        const { data } = await axios.get(url);
        list.push({
          proposalId: obj.proposalId,
          address: obj.fipEditorAddress,
          info: data,
          voters: obj.voters,
          ratio: `${obj.voters?.length} / ${fipEditors?.length}`,
          voteList: fipEditors?.map((address: string) => {
            return { address, status: obj.voters?.includes(address) ? 'Approved' : '' }
          }).sort((a) => (a.status ? -1 : 1))
        });
      } catch (e) {
        console.log(e)
      }
    }));
    setFipProposalList(list);
    setLoading(false);
  }

  const handlePageChange = (page: number) => {
    setPage(page);
  }

  const confirm = (record: any) => {
    if (record.voters.includes(address)) {
      messageApi.open({
        type: 'warning',
        content: t(HAVE_APPROVED_MSG),
      });
      return;
    }
    writeContract({
      abi: fileCoinAbi,
      address: getContractAddress(chainId, 'powerVoting'),
      functionName: 'approveFipEditor',
      args: [
        record.address,
        record.proposalId,
      ],
    });
    setCurrentProposalId(record.proposalId);
  };

  return (
    loading ? <Loading /> : <div className="px-3 mb-6 md:px-0">
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
      <div className='min-w-full bg-[#ffffff] text-left rounded-xl border-[1px] border-solid border-[#DFDFDF] overflow-hidden'>
        <div className='flow-root space-y-4'>
          <div className='font-normal text-black px-8 py-7 text-2xl border-b border-[#eeeeee] flex items-center'>
            <span>{t('content.fipEditorApprove')}</span>
          </div>
          <div className='px-8 pb-4 !mt-0'>
            <Table
              className='mb-4'
              rowKey={(record: any) => record.proposalId}
              dataSource={fipProposalList}
              columns={columns}
              pagination={false}
            />
            {
              !!approveProposalId?.length && <Row justify='end'>
                <Pagination
                  simple
                  showSizeChanger={false}
                  current={page}
                  pageSize={pageSize}
                  total={approveProposalId.length}
                  onChange={handlePageChange}
                />
              </Row>
            }
          </div>
        </div>
      </div>
    </div>
  )
}

export default FipEditorApprove;
