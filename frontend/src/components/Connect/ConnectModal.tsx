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

import { Button, Modal } from 'antd';
import { useConnect } from 'iso-filecoin-react'

/**
 * @typedef {Object} Inputs
 * @property {string} mnemonic
 * @property {string} password
 * @property {number} index
 */

/**
 * Modal component for connecting wallets
 *
 * @param {object} props
 * @param {boolean} props.isOpen - Whether the modal is open
 * @param {(isOpen: boolean) => void} props.setIsOpen - Function to set modal open state
 */
export function ConnectModal({ isOpen, setIsOpen }: { isOpen: boolean, setIsOpen: any }) {
  const {
    adapters,
    mutate: connect,
    isPending,
    loading,
    reset,
  } = useConnect()

  const handleClose = () => {
    setIsOpen(false)
    reset()
  }

  return (
    <Modal open={isOpen} onCancel={handleClose}>
      <h4>Connect a wallet</h4>
      {adapters.map((adapter) => (
        <Button
          key={adapter.name}
          title={
            adapter.support === 'NotDetected' && adapter.name === 'Filsnap'
              ? 'Install Metamask'
              : adapter.name
          }
          onClick={() => {
            connect({ adapter })
          }}
          disabled={isPending || adapter.support === 'NotDetected' || loading}
        >
          <span>{adapter.name}</span>
        </Button>
      ))}
    </Modal>
  )
}