import { Chain } from "wagmi"

export const filecoinHyperSpaceChain: Chain = {
  id: 3141,
  name: "Filecoin — HyperSpace testnet",
  network: "Filecoin — HyperSpace testnet",
  nativeCurrency: {
    decimals: 18,
    name: "Test Filecoin",
    symbol: "tFIL",
  },
  rpcUrls: {
    default: {
      http: ["https://api.hyperspace.node.glif.io/rpc/v1"],
    },
    chainstack: {
      http: ["https://filecoin-hyperspace.chainstacklabs.com/rpc/v1"],
    },
    public: {
      http: ["https://filecoin-hyperspace.chainstacklabs.com/rpc/v1"],
    },
  },
  blockExplorers: {
    default: {
      name: "ImFil Explorer",
      url: "https://imfil.io",
    },
    etherscan: {
      name: "ImFil Explorer",
      url: "https://imfil.io",
    },
  },
  testnet: true,
}

export const filecoinMainnetChain: Chain = {
  id: 314,
  name: "Filecoin — Mainnet",
  network: "Filecoin — Mainnet",
  nativeCurrency: {
    decimals: 18,
    name: "Filecoin",
    symbol: "FIL",
  },
  rpcUrls: {
    default: {
      http: ["https://api.node.glif.io/rpc/v1"],
    },
    chainstack: {
      http: ["https://filecoin-mainnet.chainstacklabs.com/rpc/v1"],
    },
    public: {
      http: ["https://filecoin-mainnet.chainstacklabs.com/rpc/v1"],
    },
  },
  blockExplorers: {
    default: {
      name: "ImFil Explorer",
      url: "https://imfil.io",
    },
    etherscan: {
      name: "ImFil Explorer",
      url: "https://imfil.io",
    },
  },
  testnet: false,
}
