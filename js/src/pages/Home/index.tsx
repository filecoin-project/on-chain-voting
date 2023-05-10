import React, { useEffect, useState } from "react"
import { Button, Alert, Table, message } from "antd"
import type { ColumnsType, TablePaginationConfig } from "antd/es/table"
import Tabulation from "../../components/Tabulation"
import MyButton from "../../components/MyButton"
import { useConnectModal } from "@rainbow-me/rainbowkit"
import { useLocation, useNavigate } from "react-router-dom"
import { usePowerVotingContract } from "../../hooks"
import useGetWallet from "../../hooks/getWallet"
import axios from "axios"
import { mainnetClient, timelockDecrypt } from "tlock-js"
// @ts-ignore
import nftStorage from "../../utils/storeNFT.js"
import pagingConfig from "../../common/js/pagingConfig"
import { getChain } from "../../utils/helpers/chain"
import { Chain } from "wagmi"

export default function Home() {
  const { openConnectModal } = useConnectModal()
  const navigate = useNavigate()
  const [addr, setAddr] = useState(false)
  const { state } = useLocation()
  const [ipfsCid, setIpfsCid] = useState<any>([])
  const [votingList, setVotingList] = useState<any>([])
  const [visibale, setVisibale] = useState(false)
  const [page, setPage] = useState(1)
  const [count, setCount] = useState(0)
  const [loading, setLoading] = useState(true)
  const [change, setChange] = useState(true)
  const pageSize = 10
  const network = {
    chainId: "0xc45", // 此处为链ID
    chainName: "Filecoin — HyperSpace testnet", // 此处为网络名称
    rpcUrls: [
      "https://api.hyperspace.node.glif.io/rpc/v1",
      "https://filecoin-hyperspace.chainstacklabs.com/rpc/v1",
      "https://filecoin-hyperspace.chainstacklabs.com/rpc/v1",
    ], // 此处为RPC URL
    nativeCurrency: {
      name: "Test Filecoin", // 此处为货币名称
      symbol: "tFIL", // 此处为货币符号
      decimals: 18,
    },
    blockExplorerUrls: ["https://imfil.io"], // 此处为区块浏览器URL
  }

  const {
    getVotingList,
    getVoteDataApi,
    updateVotingResultFun,
    isFinishVoteFun,
  } = usePowerVotingContract()
  // console.log(getVotingList(),'getVotingList()');
  useEffect(() => {
    getIpfsCid()
    if (state) {
      setVisibale(true)
      closeMessage()
    }
  }, [page])

  useEffect(() => {
    isMetaMask()
  }, [])

  // 判断是否安装小狐狸插件

  const isMetaMask = async () => {
    const provider = await window.ethereum
    if (typeof window.ethereum == "undefined") {
      console.log("1")
      // 小狐狸钱包未安装
      isLogin()
    } else {
      // 小狐狸钱包已经安装
      console.log("2")
      if (!provider.selectedAddress) {
        // 钱包未链接
        console.log("3")
        // window.ethereum.enable()
        await provider.request({
          method: "wallet_addEthereumChain",
          params: [network],
        })
        await provider.request({
          method: "eth_requestAccounts",
        })
      }
    }
  }

  // 获取投票数据
  const getIpfsCid = async () => {
    if (getVotingList) {
      const res = await getVotingList()
      setIpfsCid(res)
      setCount(res.length)
      const list = await getList(res)
      setLoading(false)
      setVotingList(list)
    }
  }
  // 获取投票项目列表
  const getList = async (prop: any) => {
    setLoading(true)
    const data = prop.slice((page - 1) * pageSize, page * pageSize)
    const ipfsUrls = data.map(
      (_item: any) => `https://${_item.cid}.ipfs.nftstorage.link/`
    )
    try {
      const responses = await Promise.all(
        ipfsUrls.map((url: string) => axios.get(url))
      )
      const results = []
      if (isFinishVoteFun) {
        for (let i = 0; i < responses.length; i++) {
          const bool = await isFinishVoteFun(data[i].cid)
          results.push({
            ...responses[i].data.string,
            cid: data[i].cid,
            bool,
          })
        }
      }
      return results
    } catch (error) {
      console.error(error)
    }
  }

  // 点击计票按钮,开始计票
  const startCounting = async (record: any) => {
    let myMap = new Map()
    if (isLogin()) {
      if (getVoteDataApi) {
        setLoading(true)
        message.success("Waiting for confirmation of transactions", 3)
        // 获取投票数据
        const res = await getVoteDataApi(record.cid)
        res.map(async (_item: any) => {
          // 生成ipfs 请求得到原始数据
          const ipfs = `https://${_item}.ipfs.nftstorage.link/`
          const r = await axios.get(ipfs)
          // 进行解密
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
        // 将计票结果上传nftStorage
        const sortedArray = Array.from(myMap.entries())
        const cid = await nftStorage(sortedArray)
        const result = await updateVotingResultFun(record.cid, cid)
        if (result) {
          setLoading(false)
          setChange(false)
          message.success("Successful vote counting", 3)
          setVisibale(true)
          closeMessage()
        }
        console.log(result)
      }
    }
  }

  // 判断是否登录了钱包
  const isLogin = () => {
    const res = localStorage.getItem("isConnect")
    console.log(res)
    if ((res == "undefined" || res == "false") && openConnectModal) {
      openConnectModal()
    } else {
      return true
    }
  }

  // 处理函数
  const handlerNavigate = (path: string, params?: any) => {
    if (isLogin()) {
      params ? navigate(path, params) : navigate(path)
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
    console.log(pagination.current)
    pagination.current && setPage(pagination.current)
    getList(ipfsCid)
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
        console.log(record.bool)
        return (
          <>
            {date <= record.Time ? (
              <div>
                <Button
                  onClick={() => {
                    handlerNavigate("/acquireNFT", { state: record })
                  }}
                  className="menu_btn"
                  type="primary"
                >
                  Claim NFT
                </Button>
                <Button
                  onClick={() => {
                    handlerNavigate("/vote", { state: record })
                  }}
                  className="menu_btn"
                  type="primary"
                >
                  Vote
                </Button>
              </div>
            ) : // <MyButton
            //   startCounting={() => {
            //     startCounting(record)
            //   }}
            //   handlerNavigate={() => {
            //     handlerNavigate("/votingResults", { state: record })
            //   }}
            //   change={change}
            // />
            record.bool ? (
              <MyButton
                startCounting={() => {
                  startCounting(record)
                }}
                handlerNavigate={() => {
                  handlerNavigate("/votingResults", { state: record })
                }}
                change={change}
              />
            ) : (
              <Button
                className="menu_btn"
                type="primary"
                onClick={() => {
                  handlerNavigate("/votingResults", { state: record })
                }}
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
      <Table
        className="rowStyle"
        rowKey={(record) => record.Name}
        dataSource={votingList}
        columns={cloumns}
        pagination={pagingConfig({ count, page, pageSize })}
        onChange={onchange}
        loading={loading}
      />
    </div>
  )
}
