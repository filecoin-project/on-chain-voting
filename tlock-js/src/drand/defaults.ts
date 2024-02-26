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

// config for the testnet chain info
import { ChainInfo } from "drand-client";

export const MAINNET_CHAIN_URL = "https://api.drand.secureweb3.com:6875/52db9ba70e0cc0f6eaf7803dd07447a1f5477735fd3f661792ba94600c84e971"
export const MAINNET_CHAIN_INFO: ChainInfo = {
    hash: "52db9ba70e0cc0f6eaf7803dd07447a1f5477735fd3f661792ba94600c84e971",
    public_key: "83cf0f2896adee7eb8b5f01fcad3912212c437e0073e911fb90022d3e760183c8c4b450b6a0a6c3ac6a5776a2d1064510d1fec758c921cc22b0e17e63aaf4bcb5ed66304de9cf809bd274ca73bab4af5a6e9c76a4bc09e76eae8991ef5ece45a",
    period: 3,
    genesis_time: 1692803367,
    groupHash: "f477d5c89f21a17c863a7f937c6a6d15859414d2be09cd448d4279af331c5d3e",
    schemeID: "bls-unchained-g1-rfc9380",
    metadata: {
        beaconID: "quicknet"

    }
};

export const defaultChainUrl = MAINNET_CHAIN_URL
export const defaultChainInfo = MAINNET_CHAIN_INFO

export const TESTNET_CHAIN_URL = "https://pl-us.testnet.drand.sh/7672797f548f3f4748ac4bf3352fc6c6b6468c9ad40ad456a397545c6e2df5bf"
export const TESTNET_CHAIN_INFO: ChainInfo = {
    hash: "7672797f548f3f4748ac4bf3352fc6c6b6468c9ad40ad456a397545c6e2df5bf",
    public_key: "8200fc249deb0148eb918d6e213980c5d01acd7fc251900d9260136da3b54836ce125172399ddc69c4e3e11429b62c11",
    genesis_time: 1651677099,
    period: 3,
    schemeID: "pedersen-bls-unchained",
    groupHash: "65083634d852ae169e21b6ce5f0410be9ed4cc679b9970236f7875cff667e13d",
    metadata: {
        beaconID: "testnet-unchained-3s"
    }
}
