import { Progress } from 'antd'
import React, { useEffect, useState } from 'react'
import { useLocation } from 'react-router-dom'
import { usePowerVotingContract } from '../../hooks'
import axios from "axios"
const VotingResults = () => {

  const { state } = useLocation()
  const { getVoteApi } = usePowerVotingContract()
  const [data1, setdata] = useState({
    Total: 0,
    data: {} as any,
    option: [[]]
  })


  useEffect(() => {

    getVoteApi && getVoteApi(state.cid).then(async (res) => {
      const ipfs = `https://${res.votingResult}.ipfs.nftstorage.link/`  //每一项投票的cid 取调取ipfs获取具体数据
      const resa = await axios.get(ipfs)
      let num = 0
      let asd = {} as any

      if (resa.data.string) {
        resa.data.string.map((item: any, index: number) => {
          num += item[1] //获得一个多少投票
          asd[index] = item[1]
          setdata({ ...data1, data: asd, Total: num, option: state.option })
        })
      }

    }).catch((err) => {

    })

    setdata({ ...data1, option: state.option })
  }, [state.cid])

  return (
    <div style={{ display: "flex", justifyContent: "center" }}>
      <div style={{ width: "500px" }}>
        <h2 style={{ fontSize: "38px" }}>{state.Name}</h2>
        <br></br>
        <br></br>
        <p>Deadline: {new Date(state.Time).toLocaleString()}</p>
        <br></br>
        <p >{state.Description}</p>
        <br></br>
        <br></br>
        <h3 style={{fontSize:"22px"}}>{'Voting Result'}</h3>
        <br></br>
        <h3 >Total Votes:{"  "+data1.Total}</h3>
        <br></br>
        <>
          {
            data1.option && data1.option.map((item: any, index: number) => {
              // console.log((data1.data[index] ? data1.data[index] : 0) / data1.Total * 100);

              return <div>
                {item} <Progress percent={(data1.data[index] ? data1.data[index] : 0) / data1.Total * 100} format={(percent) => `${data1.data[index] ? data1.data[index] : 0}`} />
              </div>
            })
          }

        </>
      </div >
    </div>
  )
}

export default VotingResults
