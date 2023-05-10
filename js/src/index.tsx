import React from "react"
import ReactDOM from "react-dom/client"
import App from "./App"
import { BrowserRouter } from "react-router-dom"
import "@rainbow-me/rainbowkit/styles.css"
import {
  darkTheme,
  getDefaultWallets,
  RainbowKitProvider,
} from "@rainbow-me/rainbowkit"
import { Chain, configureChains, createClient, WagmiConfig } from "wagmi"
import { publicProvider } from "wagmi/providers/public"
import {getChain} from "./utils/helpers/chain"

const filecoinChain: Chain = getChain()
const { chains, provider } = configureChains(
  [filecoinChain],
  [publicProvider()]
)
const { connectors } = getDefaultWallets({
  appName: "My RainbowKit App",
  chains,
})

const wagmiClient = createClient({
  autoConnect: true,
  connectors,
  provider,
})

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <WagmiConfig client={wagmiClient}>
    <RainbowKitProvider
      theme={darkTheme({
        accentColor: "orange",
        accentColorForeground: "#fff",
        borderRadius: "small",
        fontStack: "system",
        overlayBlur: "small",
      })}
      chains={chains}
      modalSize="compact"
    >
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </RainbowKitProvider>
  </WagmiConfig>
)
