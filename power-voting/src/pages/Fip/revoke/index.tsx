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
import { Popover, Table, Popconfirm, Tooltip } from 'antd';
import { InfoCircleOutlined } from '@ant-design/icons';
import {useAccount} from "wagmi";
import {
  web3AvatarUrl,
} from "../../../common/consts";
import Loading from "../../../components/Loading";
import EllipsisMiddle from "../../../components/EllipsisMiddle";

const FipRevoke = () => {
  const {isConnected, address, chain} = useAccount();
  const chainId = chain?.id || 0;

  const navigate = useNavigate();
  const prevAddressRef = useRef(address);
  const [minerIds, setMinerIds] = useState(['']);
  const [spinning, setSpinning] = useState(false);
  const [loading, setLoading] = useState(false);

  const confirm = (e: any) => {
    console.log(e);
  };
  const cancel = (e: any) => {
    console.log(e);
  };

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
    },
  ]

  const popoverData = [
    {
      address: '0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307',
      status: 'Revoked',
    },
    {
      address: '0x2E4f5898ec86A71d4D0681B33DAeD47845357BaC',
      status: 'Revoked',
    },
    {
      address: '0xe4c7b2bb1d600bCD0A9af60dda3874e369C37bc4',
      status: 'Revoked',
    },
    {
      address: '0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307',
      status: '-',
    },
    {
      address: '0xe4c7b2bb1d600bCD0A9af60dda3874e369C37bc4',
      status: '-',
    },
    {
      address: '0x2E4f5898ec86A71d4D0681B33DAeD47845357BaC',
      status: '-',
    },
    {
      address: '0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307',
      status: '-',
    },
  ]

  const columns = [
    {
      title: 'FIP Editor Address',
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
      title: 'Revoke Ratio',
      dataIndex: 'ratio',
      key: 'ratio',
      render: (value: string) => {
        return (
          <div className='flex items-center gap-2'>
            <span>{value} </span>
            <Popover content={
              <Table
                dataSource={popoverData}
                columns={popoverColumns}
                pagination={false}
              />
            }>
            </Popover>
          </div>
        )
      }
    },
    {
      title: 'Action',
      key: 'total',
      width: 100,
      render: (_: any, record: any) => <a className='hover:text-white' onClick={() => handleRevoke(record.address)}>
        <Popconfirm
          title="Revoke FIP editor"
          description="Are you sure to revoke?"
          onConfirm={confirm}
          onCancel={cancel}
          okText="Yes"
          cancelText="No"
        >
          <button className='w-[80px] h-[24px] bg-sky-500 hover:bg-sky-700 text-white py-2 px-6 rounded-xl flex justify-center items-center'>Revoke</button>
        </Popconfirm>
      </a>
    },
  ];

  const dataSource = [
    {
      address: '0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307',
      info: 'test1',
      ratio: '3 / 7'
    },
    {
      address: '0x2E4f5898ec86A71d4D0681B33DAeD47845357BaC',
      info: 'test2test2test2test2test2test2test2test2test2test2test2test2test2test2',

      ratio: '5 / 7'
    },
    {
      address: '0xe4c7b2bb1d600bCD0A9af60dda3874e369C37bc4',
      info: 'test3',
      ratio: '6 / 7'
    },
  ]

  useEffect(() => {
    if (!isConnected) {
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
    initState();
  }, [chain]);


  const initState = async () => {
    // const { getMinerIds } = await useStaticContract(chainId);
    // const { code, data: { minerIds } } = await getMinerIds(address);
    // setSpinning(false);
  }

  const handleRevoke = (address: string) => {

  }

  return (
    spinning ? <Loading /> : <div className="px-3 mb-6 md:px-0">
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
      <div className='flow-root space-y-8'>
        <div className='min-w-full bg-[#273141] rounded text-left'>
          <div className='flow-root space-y-8'>
            <div className='font-normal text-white px-8 py-7 text-2xl border-b border-[#313D4F] flex items-center'>
              <span>FIP Editor Revoke</span>
            </div>
            <div className='px-8 pb-10 !mt-0'>
              <Table
                dataSource={dataSource}
                columns={columns}
                pagination={false}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default FipRevoke;
