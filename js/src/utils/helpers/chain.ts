import {filecoinMainnetChain, filecoinHyperSpaceChain} from '../definitions/consts'

export const getChain = () => {
    const network = process.env.FILECOIN_NETWORK
    if (network == "mainnet") {
      return filecoinMainnetChain
    }
    return filecoinHyperSpaceChain
  }
  
  