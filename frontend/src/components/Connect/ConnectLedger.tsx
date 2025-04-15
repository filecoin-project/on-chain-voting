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

import {
  useAccount,
  useDisconnect
} from "iso-filecoin-react"
import { useEffect, useState } from "react"
import { ConnectModal } from "./ConnectModal"

/**
 * Connect to the network.
 */
export default function ConnectLedger() {
  const [isOpen, setIsOpen] = useState(false)
  const { account, adapter } = useAccount()
  const disconnect = useDisconnect()

  console.log(account)
  console.log(adapter)
  useEffect(() => {
    if (account) {
      setIsOpen(false)
    }
  }, [account])

  return (
    <div>
      <div
        className="Cell50"
        style={{
          alignContent: "center",
          textAlign: "center"
        }}
      >
        {!account && (
          <button type="button" onClick={() => setIsOpen(true)}>
            Connect
          </button>
        )}
        {account && (
          <button
            type="button"
            onClick={() => disconnect.mutate()}
            disabled={disconnect.isPending}
          >
            Disconnect
          </button>
        )}
      </div>

      {account && account?.address.toString()}

      <ConnectModal isOpen={isOpen} setIsOpen={setIsOpen} />
    </div>
  )
}
