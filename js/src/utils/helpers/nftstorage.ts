import { NFTStorage } from "nft.storage";

import { Media } from "../definitions/interfaces";

const key = process.env.NEXT_PUBLIC_NFT_STORAGE_API_KEY;
const NFTStorageClient = new NFTStorage({
  token: key || "undefined",
});

// NFT.Stroage
// We can fetch storage deal IDs and pinning info from NFT.Storage
// Would be super handy for future when on mainnet for expanded use case
export const fetchNFTStoreStatus = async (ipfsCID: string) => {
  const nftStatus = await NFTStorageClient.status(ipfsCID);
  return nftStatus;
};

export const storeMediaToIPFS = async (mediaFile: File) => {
  const mediaType = mediaFile.type;
  const mediaBlob = new Blob([mediaFile], { type: mediaType });
  console.info("Storing media to IPFS...");
  const mediaCID = await NFTStorageClient.storeBlob(mediaBlob);
  console.info("Success! Media CID:", mediaCID);

  return { mediaCID, mediaType };
};

// Create NFT Metadata
export const createMomentSwapMetadata = (owner: string, contentText: string, media?: Media) => {
  return {
    name: "MomentSwap Hyperspace NFTs 2023",
    description: contentText,
    image: new Blob(),
    properties: {
      authors: [{ addres: owner }],
      content: {
        "text/markdown": contentText,
      },
      media: {
        cid: media?.cid,
        type: media?.type,
      },
    },
  };
};

// Store NFT Metadata to NFT.Storage
export const storeMetadataToIPFS = async (metadata: any) => {
  console.info("Storing Metadata to IPFS...");
  const token = await NFTStorageClient.store(metadata);
  console.info("Success! Metadata IPFS URL:", token.url);
  return token.url;
};

export const ipfsCidToHttpUrl = (ipfsCid: string) => {
  return `https://${ipfsCid}.ipfs.dweb.link`;
};
