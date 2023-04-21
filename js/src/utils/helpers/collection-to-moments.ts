import { MomentMetadata, MomentSwapFRC721NFT } from "../definitions/interfaces";

export const collectionToMoments = async (collection: MomentSwapFRC721NFT[]): Promise<MomentMetadata[]> => {
  let result: Array<MomentMetadata> = [];

  if (collection.length === 0) return result;
  const preMoments = [...collection].reverse().map<MomentMetadata>((item) => ({
    address: item[0],
    id: item[2].toString(),
    timestamp: item[3].toNumber(),
    metadataURL: `https://${item[1].split("/")[2]}.ipfs.dweb.link/metadata.json`,
  }));

  if (!preMoments) return result;
  for (let index = 0; index < preMoments.length; index++) {
    let metadata;
    try {
      metadata = await fetch(preMoments[index].metadataURL).then((res) => res.json());
    } catch (err) {
      console.error(`Failed get Metadata [${preMoments[index].id}] from [${preMoments[index].metadataURL}]`);
      continue;
    }
    preMoments[index].contentText = metadata.properties.content["text/markdown"];
    preMoments[index].media = `https://${metadata.properties.media.cid}.ipfs.dweb.link`;
    preMoments[index].mediaType = metadata.properties.media.type;
    result.push(preMoments[index]);
  }

  return result;
};
