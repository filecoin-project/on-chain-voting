import { HttpCachingChain, HttpChainClient, ChainInfo } from "drand-client"


export function mainnet(): HttpChainClient {
    const clientOpts = {
        disableBeaconVerification: false,
        noCache: false,
        chainVerificationParams: {
            chainHash: "8990e7a9aaed2ffed73dbd7092123d6f289930540d7651336225dc172e51b2ce",
            publicKey: "868f005eb8e6e4ca0a47c8a77ceaa5309a47978a7c71bc5cce96366b5d7a569937c529eeda66c7293784a9402801af31"
        }
    }
    // passing an empty httpOptions arg to strip the user agent header to stop CORS issues
    return new HttpChainClient(new HttpCachingChain(MAINNET_CHAIN_URL, clientOpts), clientOpts, {})
}



export const MAINNET_CHAIN_URL = "https://api.drand.sh/8990e7a9aaed2ffed73dbd7092123d6f289930540d7651336225dc172e51b2ce"
export const MAINNET_CHAIN_INFO: ChainInfo = {
    hash: "8990e7a9aaed2ffed73dbd7092123d6f289930540d7651336225dc172e51b2ce",
    public_key: "868f005eb8e6e4ca0a47c8a77ceaa5309a47978a7c71bc5cce96366b5d7a569937c529eeda66c7293784a9402801af31",
    period: 30,
    genesis_time: 1595431050,
    groupHash: "176f93498eac9ca337150b46d21dd58673ea4e3581185f869672e59fa4cb390a",
    schemeID: "pedersen-bls-chained",
    metadata: {
        beaconID: "default"
    }
};
