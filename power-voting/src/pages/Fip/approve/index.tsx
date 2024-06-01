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

import React, {useState, useEffect, useRef} from "react";
import {Link, useNavigate} from "react-router-dom";
import { message, Popover, Table, Tooltip, Popconfirm, Button, Row, Pagination } from 'antd';
import { InfoCircleOutlined } from '@ant-design/icons';
import axios from "axios";
import {useAccount, useWriteContract, useWaitForTransactionReceipt} from "wagmi";
import type { BaseError} from "wagmi";
import {
  HAVE_APPROVED_MSG,
  STORING_DATA_MSG,
  web3AvatarUrl,
} from "../../../common/consts";
import Loading from "../../../components/Loading";
import EllipsisMiddle from "../../../components/EllipsisMiddle";
import {useFipEditors, useApproveFipId, useFipProposalDataSet, useCheckFipAddress} from "../../../common/hooks";
import fileCoinAbi from "../../../common/abi/power-voting.json";
import {getContractAddress} from "../../../utils";

const FipApprove = () => {
  const {isConnected, address, chain} = useAccount();
  const chainId = chain?.id || 0;
  const navigate = useNavigate();
  const prevAddressRef = useRef(address);
  const [messageApi, contextHolder] = message.useMessage();

  const [fipProposalList, setFipProposalList] = useState<any[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(5);
  const [loading, setLoading] = useState(false);
  const [currentProposalId, setCurrentProposalId] = useState(null);

  const { isFipAddress } = useCheckFipAddress(chainId, address);
  const { fipEditors } = useFipEditors(chainId);

  const { approveFipId, getApproveFipIdLoading } = useApproveFipId(chainId);
  const { fipProposalData, getFipProposalIdLoading, getFipProposalIdSuccess, error } = useFipProposalDataSet({
    chainId,
    idList: approveFipId,
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
      title: 'FIP Editor',
      dataIndex: 'address',
      key: 'address',
      width: 280,
      render: (value: string) => {
        return (
          <div className="w-[180px] flex items-center">
            <img className="w-[20px] h-[20px] rounded-full mr-2" src={`${web3AvatarUrl}:${value}`} alt="" />
            <a
              className="text-white hover:text-white"
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
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (value: string) => value || '-'
    },
  ]

  const columns = [
    {
      title: 'Address',
      dataIndex: 'address',
      key: 'address',
      width: 280,
      render: (value: string) => {
        return (
          <div className="w-[180px] flex items-center">
            <img className="w-[20px] h-[20px] rounded-full mr-2" src={`${web3AvatarUrl}:${value}`} alt="" />
            <a
              className="text-white hover:text-white"
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
      title: 'Info',
      dataIndex: 'info',
      key: 'info',
      ellipsis: true,
      render: (value: string) => {
        return (
          <Tooltip placement="topLeft" title={value}>
            {value}
          </Tooltip>
        )
      }
    },
    {
      title: 'Approve Ratio',
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
      title: 'Action',
      key: 'total',
      align: 'center' as const,
      width: 120,
      render: (_: any, record: any) =>
        <Popconfirm
          title="Approve FIP editor"
          description="Are you sure to approve?"
          onConfirm={() => { confirm(record) }}
          okText="Yes"
          cancelText="No"
        >
          <Button type='primary' className='w-[80px] h-[24px] flex justify-center items-center' loading={record.proposalId === currentProposalId && isLoading} >Approve</Button>
        </Popconfirm>
    },
  ];

  useEffect(() => {
    if (!isConnected || !isFipAddress) {
      navigate("/home");
      return;
    }
  }, []);

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
    if (getFipProposalIdSuccess) {
      initState();
    }
  }, [getFipProposalIdSuccess]);

  useEffect(() => {
    if (isConnected && !loading && !getApproveFipIdLoading && !getFipProposalIdLoading) {
      initState();
    }
  }, [chain,  page, address]);

  useEffect(() => {
    if (writeContractSuccess) {
      messageApi.open({
        type: 'success',
        content: STORING_DATA_MSG,
      });
      setTimeout(() => {
        navigate("/")
      }, 1000);
    }
  }, [writeContractSuccess])

  const initState = async () => {
    setLoading(true);
    setFipProposalList([]);
    await Promise.all(fipProposalData.map(async (item: any) => {
      const { result } = item;
      const url = `https://${result.voterInfoCid}.ipfs.w3s.link/`;
      const { data } = await axios.get(url);
      fipProposalList.push({
        proposalId: result.proposalId,
        address: result.fipEditorAddress,
        info: data,
        voters: result.voters,
        ratio: `${result.voters?.length} / ${fipEditors?.length}`,
        voteList: fipEditors?.map((address: string) => {
          return { address, status: result.voters?.includes(address) ? 'Approved' : '' }
        }).sort((a) => (a.status ? -1 : 1))
      });
    }))
    setFipProposalList(fipProposalList);
    setLoading(false);
  }

  const handlePageChange = (page: number) => {
    setPage(page);
  }

  const confirm = (record: any) => {
    if (record.voters.includes(address)) {
      messageApi.open({
        type: 'warning',
        content: HAVE_APPROVED_MSG,
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
        <div className="inline-flex items-center gap-1 text-skin-text hover:text-skin-link">
          <Link to="/" className="flex items-center">
            <svg className="mr-1" viewBox="0 0 24 24" width="1.2em" height="1.2em">
              <path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                    d="m11 17l-5-5m0 0l5-5m-5 5h12" />
            </svg>
            Back
          </Link>
        </div>
      </button>
      <div className='min-w-full bg-[#273141] rounded text-left'>
        <div className='flow-root space-y-4'>
          <div className='font-normal text-white px-8 py-7 text-2xl border-b border-[#313D4F] flex items-center'>
            <span>FIP Editor Approve</span>
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
              !!fipProposalData?.length && <Row justify='end'>
                    <Pagination
                        simple
                        showSizeChanger={false}
                        current={page}
                        pageSize={pageSize}
                        total={fipProposalData.length}
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

export default FipApprove;
