// Import the NFTStorage class and File constructor from the 'nft.storage' package
import { NFTStorage, Blob, File } from 'nft.storage'

// The 'mime' npm package helps us set the correct file type on our File objects


// The 'fs' builtin module on Node.js provides access to the file system
// The 'path' module provides helpers for manipulating filesystem paths

// Paste your NFT.Storage API key into the quotes:
const NFT_STORAGE_KEY = process.env.NFT_STORAGE_KEY;

/**
  * Reads an image file from `imagePath` and stores an NFT with the given name and description.
  * @param {string} string the path to an image file
  * @param {string} name a name for the NFT
  * @param {string} description a text description for the NFT
  */
async function storeNFT(string, name, description) {
  // load the file from disk
  // const image = await fileFromPath(imagePath)
  const json = JSON.stringify({ string, name, description });
  const someData = new Blob([json]);

  // create a new NFTStorage client using our API key
  const nftstorage = new NFTStorage({ token: NFT_STORAGE_KEY })

  // call client.store, passing in the image & metadata
  return nftstorage.storeBlob(someData);
}


/**
 * The main entry point for the script that checks the command line arguments and
 * calls storeNFT.
 * 
 * To simplify the example, we don't do any fancy command line parsing. Just three
 * positional arguments for imagePath, name, and description
 */
async function nftStorage(props) {
  // const args = process.argv.slice(2)
  // if (args.length !== 3) {
  //   console.error(`usage: ${process.argv[0]} ${process.argv[1]} <image-path> <name> <description>`)
  //   process.exit(1)
  // }
  const result = await storeNFT(props);
  return result;
}
// Don't forget to actually call the main function!
// We can't `await` things at the top level, so this adds
// a .catch() to grab any errors and print them to the console.

export default nftStorage;