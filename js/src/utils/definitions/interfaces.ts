import { BigNumber } from "ethers"

export type MomentSwapFRC721NFT = [string, string, BigNumber, BigNumber]

export interface CommentData {
  id?: string
  username?: string
  address?: string
  userImg?: string
  comment?: string
  timestamp?: string
}

export interface MomentMetadata {
  id: string
  address: string
  timestamp: number
  metadataURL: string
  contentText?: string
  username?: string
  userImg?: string
  media?: string
  mediaType?: "video" | "image"
}

export interface Media {
  cid: string
  url: string
  type: string
}
export interface votingListType {
  cid: string
  votingResult: string
}
