import React, { useEffect, useState } from "react"
import { ethers } from "ethers"
export default function usegetWallet() {
  const [address, setAddress] = useState("")
  useEffect(() => {
    ;(async () => {
      const provider = new ethers.providers.Web3Provider(window.ethereum)
      const signer = provider.getSigner()
      const address = await signer.getAddress()
      setAddress(address);
    })()
  }, [address])

  return address;
}
