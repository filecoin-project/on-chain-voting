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

import { Pagination, Row, Table } from "antd";
import { useEffect, useRef, useState } from "react";
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from "react-router-dom";
import { useFipList } from "../../../common/store";
import Loading from "../../../../src/components/Loading";
import { useAccount } from "wagmi";
import {
  calibrationChainId,
  getFipListApi,
  web3AvatarUrl
} from "../../../common/consts";
import { useAddresses } from "iso-filecoin-react"
import axios from "axios";
import { getBlockExplorers, isFilAddress } from "../../../utils"
const FipEditorList = () => {
  const { isConnected, address, chain } = useAccount();
  const { address0x } = useAddresses({ address: address as string })
  const { t } = useTranslation();
  const chainId = chain?.id || calibrationChainId;
  const navigate = useNavigate();
  const prevAddressRef = useRef(address);
  const [loading, setLoading] = useState<boolean>(false)
  const { fipList, totalSize } = useFipList((state: any) => state.data);
  const setFipList = useFipList((state: any) => state.setFipList)
  const [page, setPage] = useState(1);
  const [pageSize] = useState(5);
  const [showList, setShowList] = useState([])
  const getFipList = async () => {
    const params = {
      chainId,
    }
    setLoading(true)
    const { data: { data: fipList } } = await axios.get(getFipListApi, { params });
    if (isFilAddress(address!) && address0x.data) {
      setFipList(fipList, address0x.data.toString())
    } else {
      setFipList(fipList, address)
    }
    if (fipList.length > 5) {
      setShowList(fipList.slice(0, pageSize));
    } else {
      setShowList(fipList)
    }
    setLoading(false)
  }
  useEffect(() => {
    if (!address || !chainId) return
    getFipList()
  }, [address, chainId])
  const columns = [
    {
      title: t('content.FIPEditor'),
      dataIndex: 'editor',
      key: 'editor',
      width: 280,
      render: (value: string) => {
        return (
          <div className="w-[180px] flex items-center">
            <img className="w-[20px] h-[20px] rounded-full mr-2" src={`${web3AvatarUrl}:${value}`} alt="" />
            <a
              className="text-black hover:text-blue"
              target="_blank"
              rel="noopener"
              href={getBlockExplorers(chain, value)}
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
  const handlePageChange = (page: number) => {
    setPage(page);
    const startIndex = (page - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    const newData = fipList.slice(startIndex, endIndex)
    setShowList(newData);
  }
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
            {loading ? <Loading /> : (
              <Table
                className='mb-4'
                rowKey={(record: any) => record.editor}
                dataSource={showList}
                columns={columns}
                pagination={false}
              />
            )
            }
            {
              totalSize > 5 && <Row justify='end'>
                <Pagination
                  simple
                  showSizeChanger={false}
                  current={page}
                  pageSize={pageSize}
                  total={totalSize}
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

export default FipEditorList;
