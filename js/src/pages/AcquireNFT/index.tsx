import { message, Form, Input, Radio, Space, Button } from "antd"
import React, { useEffect, useState } from "react"
import { useLocation, useNavigate } from "react-router-dom"
import { usePowerVotingContract } from "../../hooks"

const contractAddress = process.env.VOTING_CONTRACT_ADDRESS
if (!contractAddress) {
  throw new Error("Please set VOTING_CONTRACT_ADDRESS in a .env file")
}

const AcquireNFT = () => {
  const navigate = useNavigate()
  const { VotingNFT } = usePowerVotingContract()

  const { state } = useLocation()
  const [loading, setLoading] = useState(false)

  const onClaim = async (values: any) => {
    setLoading(true)
    if (VotingNFT) {
      await VotingNFT()
        .then(() => {
          message.success("Waiting for the transaction to be chained!")
          setTimeout(() => {
            setLoading(false)
            message.success("Claim NFT successful!")
            navigate("/")
          }, 3000)
        })
        .catch((e) => {
          console.log(e)
          message.warning("Please link to the wallet!")
        })
    }
  }
  return (
    <div className="acquireNFT_container main">
      <h2 style={{ fontSize: "38px" }}>{state.Name}</h2>
      <br></br>
      <br></br>
      <p>{new Date(state.Time).toLocaleString()}</p>
      <br></br>
      <p style={{ wordWrap: "break-word" }}>{state.Description}</p>
      <br></br>
      <Button
        loading={loading}
        className="menu_btn"
        type="primary"
        onClick={onClaim}
      >
        Claim NFT
      </Button>
    </div>
  )
}

export default AcquireNFT
