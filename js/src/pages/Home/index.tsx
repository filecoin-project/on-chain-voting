import React, { useEffect, useState } from "react"
import { Button, Alert, Table } from "antd"
import type { ColumnsType, TablePaginationConfig } from "antd/es/table"
import Tabulation from "../../components/Tabulation"
import {
  useConnectModal,
  useAccountModal,
  useChainModal,
} from "@rainbow-me/rainbowkit"
import { useLocation, useNavigate } from "react-router-dom"
import { usePowerVotingContract } from "../../hooks"
import useGetWallet from "../../hooks/getWallet"
import axios from "axios"
import { mainnetClient, timelockDecrypt } from "tlock-js"
// @ts-ignore
import nftStorage from "../../utils/storeNFT.js"
import pagingConfig from "../../common/js/pagingConfig"

export default function Home() {
  const { openConnectModal } = useConnectModal()
  const navigate = useNavigate()
  const { state } = useLocation()
  const [ipfsCid, setIpfsCid] = useState<any>([])
  const [votingList, setVotingList] = useState<any>([])
  const [visibale, setVisibale] = useState(false)
  const [page, setPage] = useState(1)
  const [count, setCount] = useState(0)
  const pageSize = 10

  const { getVotingList, getVoteDataApi, updateVotingResultFun, isFinishVoteFun, updateVotingResultBatchFun } = usePowerVotingContract()
  // console.log(getVotingList(),'getVotingList()');
  useEffect(() => {
    // getList()
    getIpfsCid()
    if (state) {
      setVisibale(true)
      closeMessage()
    }
  }, [page])

  // 获取投票数据
  const getIpfsCid = async () => {
    if (getVotingList) {
      const res = await getVotingList()
      setIpfsCid(res)
      setCount(res.length)
      getList(res)
    }
  }
  const getList = async (prop: any) => {
    const date = new Date().getTime()
    console.log(page);
    const data = prop.slice((page - 1) * pageSize, page * pageSize);
    const arr = data.map(async (_item: any) => {
      // 所有的投票项目
      // 将cid进行拼接 https://bafkreihw6rmsatp7x43v4zf5d2d6xrsnilpzbhsduf6szowdggbx5ji4re.ipfs.nftstorage.link/
      const ipfs = `https://${_item.cid}.ipfs.nftstorage.link/`
      const res = await axios.get(ipfs)
      const dataa = { ...res.data.string, cid: _item.cid }
      return dataa
    })

    Promise.all(arr)
      .then((results) => {
        let arr = [] as any;
        // 每个投票项目对应的数据
        // 这里打印出每个Promise对象的结果 results
        setVotingList(results)
        results.map(async (item) => {
          // let myMap = new Map();
          if (item.Time <= date) {
            // 遍历每个投票项目有没有过期
            if (getVoteDataApi) {
              getVoteDataApi(item.cid).then((res) => {
                let myMap = new Map()
                const promises = res.map(async (i: any) => {
                  const ipfs = `https://${i}.ipfs.nftstorage.link/`
                  const r = await axios.get(ipfs)
                  const dataString = await timelockDecrypt(
                    r.data.string,
                    mainnetClient()
                  )
                  const data = JSON.parse(dataString)
                  if (myMap.get(data.index) === undefined) {
                    myMap.set(data.index, 1)
                  } else {
                    myMap.set(data.index, myMap.get(data.index) + 1)
                  }
                })

                Promise.all(promises)
                  .then(async () => {
                    const sortedArray = Array.from(myMap.entries())
                    const cid = await nftStorage(sortedArray)
                    arr = [...arr,[item.cid,cid]];
                    console.log(arr);
                    // const res = await updateVotingResultFun(item.cid, cid);


                  })
                  .catch((err) => {
                    console.error(err)
                  })
              })
            }
          }
        })
        contract(arr);
      })
      .catch((error) => {
        console.error(error)
      })
  }
  const contract = async (arr:any)=>{
    if(updateVotingResultBatchFun){
      const res = await updateVotingResultBatchFun(arr);
      console.log(res);
    }

  }

  // 判断是否登录了钱包
  const isLogin = (path: string, params?: any) => {
    const res = localStorage.getItem("isConnect")
    console.log(res)
    if (res !== "undefined") {
      console.log(res)
      console.log(params)
      params ? navigate(path, params) : navigate(path)
    } else if (openConnectModal) {
      openConnectModal()
    }
  }

  // 关闭提示通知
  const closeMessage = () => {
    setTimeout(() => {
      setVisibale(false)
    }, 10000)
  }

  // 分页
  const onchange = (pagination: any) => {
    console.log(pagination.current);
    pagination.current && setPage(pagination.current)
    getList(ipfsCid);
  }

  const cloumns = [
    {
      title: "Name",
      dataIndex: "Name",
    },
    {
      title: "Deadline",
      dataIndex: "Time",
      render: (text: number) => {
        return <>{new Date(text).toLocaleString()}</>
      },
    },
    {
      title: "Status",
      dataIndex: "status",
      render: (text: string, record: any) => {
        const date = new Date().getTime()
        return <div>{date >= record.Time ? "Completed" : "In Progress"}</div>
      },
    },
    {
      title: "Operations",
      dataIndex: "Operations",
      render: (text: string, record: any) => {
        const date = new Date().getTime()
        console.log(date, record.Time)
        return (
          <>
            {date <= record.Time ? (
              <div>
                <Button
                  onClick={() => {
                    isLogin("/acquireNFT", { state: record })
                  }}
                  className="menu_btn"
                  type="primary"
                >
                  Claim NFT
                </Button>
                <Button
                  onClick={() => {
                    isLogin("/vote", { state: record })
                  }}
                  className="menu_btn"
                  type="primary"
                >
                  Vote
                </Button>
              </div>
            ) : (
              <Button
                onClick={() => {
                  isLogin("/votingResults", { state: record })
                }}
                className="menu_btn"
                type="primary"
              >
                View
              </Button>
            )}
          </>
        )
      },
    },
  ]
  return (
    <div className="home_container main">
      {visibale ? (
        <Alert
          style={{ marginBottom: "10px", fontSize: "16px" }}
          banner={true}
          message="Need to waiting for the transaction to be chained!"
          type="warning"
        />
      ) : (
        ""
      )}
      {/* <Tabulation
        className="rowStyle"
        rowKey={(record) => {
          return record.Name
        }}
        dataConf={dataConf}
        onChange={onchange}
        tableDataTypeConfig={cloumns}
      /> */}
      <Table
        className="rowStyle"
        rowKey={(record) => record.Name}
        dataSource={votingList}
        columns={cloumns}
        pagination={pagingConfig({ count, page, pageSize })}
        onChange={onchange}
      />
    </div>
  )
}
