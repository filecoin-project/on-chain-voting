import { message, Form, Input, Radio, Space, Button } from "antd"
import React, { useEffect, useState } from "react"
import { useLocation, useNavigate } from "react-router-dom"
import {
  timelockEncrypt,
  roundAt,
  HttpChainClient,
  mainnetClient,
} from "tlock-js"
// @ts-ignore
import nftStorage from "../../utils/storeNFT"
import { useDynamicContract } from "../../hooks/use-power-voting-contract"
import { ethers } from "ethers"
const Vote = () => {
  const [loading, setLoading] = useState<boolean>(false)
  const location = useLocation()
  const navigate = useNavigate()
  const { voteApi } = useDynamicContract()
  const [address, setAddress] = useState("")

  useEffect(() => {
    getConnectedAddress()
  }, [])

  const provider = new ethers.providers.Web3Provider(window.ethereum)
  const [value, setvalue] = useState("")
  async function getConnectedAddress() {
    // 获取当前连接的账户地址
    const signer = provider.getSigner()
    const address = await signer.getAddress()
    setAddress(address)
  }

  const onFinish = async (values: any) => {
    if (values === undefined) {
      message.error("Please confirm if you want to add a voting option")
    } else {
      // 调取接口发送数据
      setLoading(true)

      const payload = Buffer.from(
        JSON.stringify({
          cid: location.state.cid,
          index: values.option,
          address,
        })
      )

      const chainInfo = await mainnetClient().chain().info()

      const roundNumber = roundAt(location.state.Time, chainInfo) // drand 随机数索引

      const ciphertext = await timelockEncrypt(
        roundNumber,
        payload,
        mainnetClient()
      ) // 加密
      const res = await nftStorage(ciphertext) // 保存到 ipfs

      if (res) {
        if (voteApi) {
          voteApi(location.state.cid, res, address)
            .then((result) => {
              setLoading(false);
              message.success("Creation successful!", 3)
              setTimeout(() => {
                navigate("/")
              }, 3000)
            })
            .catch((error) => {
              setLoading(false);
              message.warning("A user cannot vote more than once!",3)
              setTimeout(() => {
                navigate("/")
              }, 3000)
            })
        }
      }
    }
  }

  const onFinishFailed = (errorInfo: any) => {
    console.log("Failed:", errorInfo)
  }

  return (
    <div className="main" style={{ width: "500px" }}>
      <h2 style={{ fontSize: "38px" }}>{location.state.Name}</h2>
      <br></br>
      <br></br>
      <p>Deadline: {new Date(location.state.Time).toLocaleString()}</p>
      <br></br>
      <p>{location.state.Description}</p>
      <br></br>
      <Form
        layout={"vertical"}
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
        labelCol={{ span: 12 }}
        wrapperCol={{ span: 24 }}
      >
        <Form.Item name="option" label="option:">
          <Radio.Group value={value} onChange={(e) => setvalue(e.target.value)}>
            <Space direction="vertical">
              {location.state.option.map((item: any, index: any) => {
                return (
                  <Radio key={index} value={index}>
                    {" "}
                    {item}{" "}
                  </Radio>
                )
              })}
            </Space>
          </Radio.Group>
        </Form.Item>
        <br></br>
        <br></br>
        <Form.Item>
          <Button
            loading={loading}
            htmlType="submit"
            className="menu_btn"
            style={{ color: "#fff" }}
          >
            Submit
          </Button>
        </Form.Item>
      </Form>
    </div>
  )
}

export default Vote
