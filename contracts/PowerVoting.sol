// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;


contract PowerVoting {

    // 项目
    struct Voting {
        string cid;
        string votingResult;
    }

    // 保存投票项目数据
    mapping(string => Voting) votingMap;
    // 保存投票项目列表
    string[] votingLists;
    // 保存NFT是否投过票
    mapping(string => mapping(string => bool)) nftVote;
    // 保存每一个项目对应的投票数据
    mapping(string => string[]) votingDataMap;


    // 获取投票列表
    function votingList() external view returns(Voting[] memory) {
        Voting[] memory res = new Voting[](votingLists.length);
        for(uint i = 0; i < votingLists.length; i++){
            Voting memory v = votingMap[votingLists[i]];
            res[i] = v;
        }
        return res;
    }

    // 创建投票
    function createVoting(string memory _cid) external {
        Voting memory voting;
        voting.cid = _cid;
        votingMap[_cid] = voting;
        votingLists.push(_cid);
    }

    // 校验NFT
    function nftVerify(string memory _cid, string memory _tokenId) external view returns(bool) {
        bool isUse = nftVote[_cid][_tokenId];
        return !isUse;
    }

    // 投票
    function vote(string memory _cid, string memory _votingData, string memory _tokenId) external {
        require(!nftVote[_cid][_tokenId], "The current NFT already in use");
        string[] storage data = votingDataMap[_cid];
        data.push(_votingData);
        votingDataMap[_cid] = data;
        nftVote[_cid][_tokenId] = true;
    }

    // 获取投票详情
    function getVote(string memory _cid) external view returns (Voting memory) {
        return votingMap[_cid];
    }

    // 更新投票结果
    function updateVotingResult(string memory _cid, string memory _votingResult) external {
        Voting memory voting = votingMap[_cid];
        voting.votingResult = _votingResult;
        votingMap[_cid] = voting;
    }

    // 批量更新投票数据
    function updateVotingResultBatch(string[][] memory _list) external {
        for (uint256 i = 0; i < _list.length; i++) {
            Voting memory voting = votingMap[_list[i][0]];
            voting.votingResult = _list[i][1];
            votingMap[_list[i][0]] = voting;
        }
    }

    // 返回投票数据
    function getVoteData(string memory _cid) external view returns(string[] memory){
        return votingDataMap[_cid];
    }

    // 判断是否需要计票
    function isFinishVote(string memory _cid) external view returns(bool) {
        if (bytes(votingMap[_cid].votingResult).length == 0) {
            return true;
        } else {
            return false;
        }
    }

}
