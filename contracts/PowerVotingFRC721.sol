// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// https://docs.openzeppelin.com/contracts/4.x/erc721
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "hardhat/console.sol";

/// NFT 水龙头
contract PowerVotingFRC721 is ERC721URIStorage, Ownable {
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIds;

    struct  powerVotingFRC721NFT {
        address owner;
        uint256 tokenId;
        uint256 timestamp;
    }
    mapping(address => bool) public requestedAddress;   /// 记录领取过代币的地址
    powerVotingFRC721NFT[] public nftCollection;        /// NFT 记录
    mapping(address => powerVotingFRC721NFT[]) public nftCollectionByOwner;

    event NewPowerVotingFRC721NFTMinted(
        address indexed sender,
        uint256 indexed tokenId,
        uint256 indexed timestamp
    );
    // 加入投票项目验证后，再添加每个用户只能领一次
    modifier Received() {
        require(requestedAddress[msg.sender] == false, "Can't Request Multiple Times!");
        _;
    }

    constructor() ERC721("PowerVoting NFTs", "BAC") {
        console.log("Hello Fil-ders! Now creating PowerVoting FRC721 NFT contract!");
    }

    function mintPowerVotingNFT() public  returns(uint256) {

        uint256 newItemId = _tokenIds.current();

        powerVotingFRC721NFT memory newNFT = powerVotingFRC721NFT({
            owner: msg.sender,
            tokenId: newItemId,
            timestamp: block.timestamp
        });

        _mint(msg.sender, newItemId);
        nftCollectionByOwner[msg.sender].push(newNFT);

        _tokenIds.increment();

        nftCollection.push(newNFT);

        requestedAddress[msg.sender] = true;

        emit NewPowerVotingFRC721NFTMinted(msg.sender, newItemId, block.timestamp);

        return newItemId;
    }

    /**
     * @notice helper function to display NFTs for frontends
     */
    function getNFTCollection() public view returns (powerVotingFRC721NFT[] memory) {
        return nftCollection;
    }


}
