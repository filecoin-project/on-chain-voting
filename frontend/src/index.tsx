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

import ReactDOM from "react-dom/client"
import { type Wallet, connectorsForWallets } from "@rainbow-me/rainbowkit"
import { metaMaskWallet } from "@rainbow-me/rainbowkit/wallets";
import "@rainbow-me/rainbowkit/styles.css"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { FilecoinProvider } from 'iso-filecoin-react'
import { WalletAdapterFilsnap } from 'iso-filecoin-wallets'
import { WagmiProvider, http, createConfig, createConnector } from "wagmi"
import { getProvider } from 'filsnap-adapter'
import { getAddress } from "viem"
import { filecoin } from "wagmi/chains"
import { walletConnectProjectId } from "./common/consts"
import { BrowserRouter } from "react-router-dom"
import App from "./App"

const queryClient = new QueryClient();

const filSnapAdapter =  new WalletAdapterFilsnap();

const filSnap = (): Wallet => ({
  id: "filSnap",
  name: "FilSnap",
  iconUrl: "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzUiIGhlaWdodD0iMzQiIHZpZXdCb3g9IjAgMCAzNSAzNCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTMyLjcwNzcgMzIuNzUyMkwyNS4xNjg4IDMwLjUxNzRMMTkuNDgzMyAzMy45MDA4TDE1LjUxNjcgMzMuODk5MUw5LjgyNzkzIDMwLjUxNzRMMi4yOTIyNSAzMi43NTIyTDAgMjUuMDQ4OUwyLjI5MjI1IDE2LjQ5OTNMMCA5LjI3MDk0TDIuMjkyMjUgMC4zMTIyNTZMMTQuMDY3NCA3LjMxNTU0SDIwLjkzMjZMMzIuNzA3NyAwLjMxMjI1NkwzNSA5LjI3MDk0TDMyLjcwNzcgMTYuNDk5M0wzNSAyNS4wNDg5TDMyLjcwNzcgMzIuNzUyMloiIGZpbGw9IiNGRjVDMTYiLz4KPHBhdGggZD0iTTIuMjkzOTUgMC4zMTIyNTZMMTQuMDY5MSA3LjMyMDQ3TDEzLjYwMDggMTIuMTMwMUwyLjI5Mzk1IDAuMzEyMjU2WiIgZmlsbD0iI0ZGNUMxNiIvPgo8cGF0aCBkPSJNOS44Mjk1OSAyNS4wNTIyTDE1LjAxMDYgMjguOTgxMUw5LjgyOTU5IDMwLjUxNzVWMjUuMDUyMloiIGZpbGw9IiNGRjVDMTYiLz4KPHBhdGggZD0iTTE0LjU5NjYgMTguNTU2NUwxMy42MDA5IDEyLjEzMzNMNy4yMjY5MiAxNi41MDA5TDcuMjIzNjMgMTYuNDk5M1YxNi41MDI1TDcuMjQzMzUgMjAuOTk4M0w5LjgyODA5IDE4LjU1NjVIOS44Mjk3NEgxNC41OTY2WiIgZmlsbD0iI0ZGNUMxNiIvPgo8cGF0aCBkPSJNMzIuNzA3NyAwLjMxMjI1NkwyMC45MzI2IDcuMzIwNDdMMjEuMzk5MyAxMi4xMzAxTDMyLjcwNzcgMC4zMTIyNTZaIiBmaWxsPSIjRkY1QzE2Ii8+CjxwYXRoIGQ9Ik0yNS4xNzIyIDI1LjA1MjJMMTkuOTkxMiAyOC45ODExTDI1LjE3MjIgMzAuNTE3NVYyNS4wNTIyWiIgZmlsbD0iI0ZGNUMxNiIvPgo8cGF0aCBkPSJNMjcuNzc2NiAxNi41MDI1SDI3Ljc3ODNIMjcuNzc2NlYxNi40OTkzTDI3Ljc3NSAxNi41MDA5TDIxLjQwMSAxMi4xMzMzTDIwLjQwNTMgMTguNTU2NUgyNS4xNzIyTDI3Ljc1ODYgMjAuOTk4M0wyNy43NzY2IDE2LjUwMjVaIiBmaWxsPSIjRkY1QzE2Ii8+CjxwYXRoIGQ9Ik05LjgyNzkzIDMwLjUxNzVMMi4yOTIyNSAzMi43NTIyTDAgMjUuMDUyMkg5LjgyNzkzVjMwLjUxNzVaIiBmaWxsPSIjRTM0ODA3Ii8+CjxwYXRoIGQ9Ik0xNC41OTQ3IDE4LjU1NDlMMTYuMDM0MSAyNy44NDA2TDE0LjAzOTMgMjIuNjc3N0w3LjIzOTc1IDIwLjk5ODRMOS44MjYxMyAxOC41NTQ5SDE0LjU5M0gxNC41OTQ3WiIgZmlsbD0iI0UzNDgwNyIvPgo8cGF0aCBkPSJNMjUuMTcyMSAzMC41MTc1TDMyLjcwNzggMzIuNzUyMkwzNS4wMDAxIDI1LjA1MjJIMjUuMTcyMVYzMC41MTc1WiIgZmlsbD0iI0UzNDgwNyIvPgo8cGF0aCBkPSJNMjAuNDA1MyAxOC41NTQ5TDE4Ljk2NTggMjcuODQwNkwyMC45NjA3IDIyLjY3NzdMMjcuNzYwMiAyMC45OTg0TDI1LjE3MjIgMTguNTU0OUgyMC40MDUzWiIgZmlsbD0iI0UzNDgwNyIvPgo8cGF0aCBkPSJNMCAyNS4wNDg4TDIuMjkyMjUgMTYuNDk5M0g3LjIyMTgzTDcuMjM5OTEgMjAuOTk2N0wxNC4wMzk0IDIyLjY3NkwxNi4wMzQzIDI3LjgzODlMMTUuMDA4OSAyOC45NzZMOS44Mjc5MyAyNS4wNDcySDBWMjUuMDQ4OFoiIGZpbGw9IiNGRjhENUQiLz4KPHBhdGggZD0iTTM1LjAwMDEgMjUuMDQ4OEwzMi43MDc4IDE2LjQ5OTNIMjcuNzc4M0wyNy43NjAyIDIwLjk5NjdMMjAuOTYwNyAyMi42NzZMMTguOTY1OCAyNy44Mzg5TDE5Ljk5MTIgMjguOTc2TDI1LjE3MjIgMjUuMDQ3MkgzNS4wMDAxVjI1LjA0ODhaIiBmaWxsPSIjRkY4RDVEIi8+CjxwYXRoIGQ9Ik0yMC45MzI1IDcuMzE1NDNIMTcuNDk5OUgxNC4wNjczTDEzLjYwMDYgMTIuMTI1MUwxNi4wMzQyIDI3LjgzNEgxOC45NjU2TDIxLjQwMDggMTIuMTI1MUwyMC45MzI1IDcuMzE1NDNaIiBmaWxsPSIjRkY4RDVEIi8+CjxwYXRoIGQ9Ik0yLjI5MjI1IDAuMzEyMjU2TDAgOS4yNzA5NEwyLjI5MjI1IDE2LjQ5OTNINy4yMjE4M0wxMy41OTkxIDEyLjEzMDFMMi4yOTIyNSAwLjMxMjI1NloiIGZpbGw9IiM2NjE4MDAiLz4KPHBhdGggZD0iTTEzLjE3IDIwLjQxOTlIMTAuOTM2OUw5LjcyMDk1IDIxLjYwNjJMMTQuMDQwOSAyMi42NzI3TDEzLjE3IDIwLjQxODJWMjAuNDE5OVoiIGZpbGw9IiM2NjE4MDAiLz4KPHBhdGggZD0iTTMyLjcwNzcgMC4zMTIyNTZMMzQuOTk5OSA5LjI3MDk0TDMyLjcwNzcgMTYuNDk5M0gyNy43NzgxTDIxLjQwMDkgMTIuMTMwMUwzMi43MDc3IDAuMzEyMjU2WiIgZmlsbD0iIzY2MTgwMCIvPgo8cGF0aCBkPSJNMjEuODMzIDIwLjQxOTlIMjQuMDY5NEwyNS4yODUzIDIxLjYwNzlMMjAuOTYwNCAyMi42NzZMMjEuODMzIDIwLjQxODJWMjAuNDE5OVoiIGZpbGw9IiM2NjE4MDAiLz4KPHBhdGggZD0iTTE5LjQ4MTcgMzAuODM2MkwxOS45OTExIDI4Ljk3OTRMMTguOTY1OCAyNy44NDIzSDE2LjAzMjdMMTUuMDA3MyAyOC45Nzk0TDE1LjUxNjcgMzAuODM2MiIgZmlsbD0iIzY2MTgwMCIvPgo8cGF0aCBkPSJNMTkuNDgxNiAzMC44MzU5VjMzLjkwMjFIMTUuNTE2NlYzMC44MzU5SDE5LjQ4MTZaIiBmaWxsPSIjQzBDNENEIi8+CjxwYXRoIGQ9Ik05LjgyOTU5IDMwLjUxNDJMMTUuNTIgMzMuOTAwOFYzMC44MzQ2TDE1LjAxMDYgMjguOTc3OEw5LjgyOTU5IDMwLjUxNDJaIiBmaWxsPSIjRTdFQkY2Ii8+CjxwYXRoIGQ9Ik0yNS4xNzIxIDMwLjUxNDJMMTkuNDgxNyAzMy45MDA4VjMwLjgzNDZMMTkuOTkxMSAyOC45Nzc4TDI1LjE3MjEgMzAuNTE0MloiIGZpbGw9IiNFN0VCRjYiLz4KPC9zdmc+Cg==",
  iconBackground: "",
  downloadUrls: {
    browserExtension: "https://github.com/Chainsafe/filsnap",
    android: "https://play.google.com/store/apps/details?id=io.metamask",
    ios: "https://apps.apple.com/us/app/metamask-blockchain-wallet/id1438144202"
  },
  extension: {
    instructions: {
      learnMoreUrl: "https://github.com/Chainsafe/filsnap",
      steps: [
        {
          description: "FilSnap requires MetaMask. Install MetaMask first.",
          step: "install",
          title: "Install MetaMask"
        },
        {
          description: "After installing MetaMask, click below to install FilSnap.",
          step: "create",
          title: "Install FilSnap"
        },
        {
          description: "Once installed, click \"Connect\" to link your wallet.",
          step: "connect",
          title: "Connect Wallet"
        }
      ]
    }
  },
  createConnector: (walletDetails) => {
    return createConnector((config) => ({
      id: "filsnap",
      name: "FileSnap",
      type: "filsnap",
      network: "mainnet",
      isConnecting: false,
      connected: false,
      syncWithProvider: true,
      async connect( chain) {
        const snaps = await window.ethereum.request({
          method: "wallet_getSnaps"
        });
        const isFileSnapInstalled =
          snaps &&
          snaps["npm:filsnap"] &&
          snaps["npm:filsnap"].enabled

        if (!isFileSnapInstalled) {
          try {
            const installResult = await window.ethereum.request({
              method: "wallet_requestSnaps",
              params: {
                "npm:filsnap": {
                  version: "^1.6.0"
                }
              }
            })

            if (!installResult["npm:filsnap"]) {
              new Error("FileSnap installation failed")
            }
          } catch (error) {
            throw new Error("User rejected FileSnap installation")
          }
        }

        const chainId = chain?.chainId;
        let currentChainId = await this.getChainId();
        if (chainId && currentChainId !== chainId) {
          const chain = await this.switchChain!({ chainId: chainId })
          currentChainId = chain?.id ?? currentChainId;
        }

        const res = await filSnapAdapter.connect({ network: 'mainnet' });
        const accounts = [res?.account?.address.toString() as `0x${string}` || '0x'];
        config.emitter.emit('connect', { accounts, chainId: currentChainId })
        return {
          accounts,
          chainId: currentChainId
        }
      },
      async switchChain({ chainId }) {
        try {
          const chain = config.chains.find((x) => x.id === chainId);
          if (!chain) throw new Error("Switch chain failed");
          config.emitter.emit("change", {
            chainId,
          });
          await filSnapAdapter.changeNetwork('mainnet');
          config.emitter.emit("change", { chainId: Number(chainId), accounts: await this.getAccounts() })

          return chain;
        } catch (error: unknown) {
          throw new Error(JSON.stringify(error));
        }
      },
      async disconnect() {
        filSnapAdapter.disconnect();
        config.emitter.emit("disconnect");
      },
      async getAccounts() {
        const { result } = await window.ethereum.request({
          method: 'wallet_invokeSnap',
          params: {
            snapId: 'npm:filsnap',
            request: { method: 'fil_getAddress' }
          }
        })
        return [result];
      },
      async getChainId() {
        return Number(window.ethereum?.chainId) || 0
      },
      async getProvider() {
        return await getProvider();
      },
      async isAuthorized() {
        try {
          const accounts = await this.getAccounts()
          return !!accounts.length
        } catch {
          return false
        }
      },
      onAccountsChanged(accounts) {
        if (accounts.length === 0) config.emitter.emit("disconnect")
        else
          config.emitter.emit("change", { accounts: accounts.map((x) => getAddress(x)) })
      },
      onChainChanged(chainId) {
        config.emitter.emit("change", { chainId: Number(chainId) })
      },
      onDisconnect() {
        config.emitter.emit("disconnect")
      },
      ...walletDetails
    }))
  }
})

const connectors = connectorsForWallets(
  [
    {
      groupName: "Recommended",
      wallets: [
        metaMaskWallet,
        filSnap
      ]
    }
  ],
  {
    appName: "power-voting",
    projectId: walletConnectProjectId
  }
)

const config = createConfig({
  chains: [filecoin],
  transports: {
    [filecoin.id]: http()
  },
  multiInjectedProviderDiscovery: false,
  connectors: [
    ...connectors,
  ]
})

//dynamic add font
const style = document.createElement("style")
style.type = "text/css"
style.innerHTML = `
  @font-face {
    font-family: 'SuisseIntl';
    src: url('/fonts/SuisseIntl-Regular.ttf') format('truetype');
  }
`

document.head.appendChild(style)


ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <BrowserRouter>
    <WagmiProvider config={config}>
      <FilecoinProvider adapters={[filSnapAdapter]} network="mainnet" reconnectOnMount={true}>
        <QueryClientProvider client={queryClient}>
          <App />
        </QueryClientProvider>
      </FilecoinProvider>
    </WagmiProvider>
  </BrowserRouter>
)

