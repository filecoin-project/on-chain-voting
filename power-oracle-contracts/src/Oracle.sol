// SPDX-License-Identifier: Apache-2.0
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

pragma solidity ^0.8.19;

import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "./interfaces/IOracle.sol";
import "./Powers.sol";
import { UUPSUpgradeable } from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";

contract Oracle is IOracle, Ownable2StepUpgradeable, UUPSUpgradeable {

    using EnumerableSet for EnumerableSet.UintSet;
    using EnumerableSet for EnumerableSet.AddressSet;
    using Powers for uint64;
    using Powers for address;

    // max day, 60, Save once every day, with a maximum storage period of 60 days.
    uint64 constant public MAX_HISTORY = 60;

    // task id
    uint256 private  _taskId;
    // history power id
    int8 private  _historyPowerId;
    // PowerVoting contract address
    address public powerVotingContract;

    // oracle node allowlist
    mapping(address => bool) public nodeAllowList;
    // address status map, key: voter address value: block height
    mapping(address => uint256) public voterAddressToBlockHeight;
    // task list
    EnumerableSet.UintSet taskIdList;
    // f4 task list
    EnumerableSet.UintSet f4TaskIdList;
    // task id to ucan cid
    mapping(uint256 => string) public taskIdToUcanCid;
    // task id to address
    mapping(uint256 => address) public f4TaskIdToAddress;
    // voter list
    EnumerableSet.AddressSet voterList;
    // voter info map
    mapping(address => VoterInfo) public voterToInfo;
    // voter to miner id list
    mapping(address => uint64[]) public voterToMinerIds;
    // voter to history power
    mapping(address => mapping(uint256 => Power)) public voterToHistoryPower;
    // voter to id
    mapping(address => PowerStatus) public voterToPowerStatus;
    // actor id list
    mapping(uint64 => bool) public actorIdList;
    // github account list
    mapping(string => bool) public githubAccountList;

    /**
     * @dev Modifier that allows a function to be called only by addresses in the node allow list.
     */
    modifier onlyInAllowList(){
        if (!nodeAllowList[msg.sender]) {
            revert PermissionError("Not in allow list error.");
        }
        _;
    }

    /**
     * @dev Modifier that ensures the provided address is non-zero.
     * @param addr The address to check.
     */
    modifier nonZeroAddress(address addr){
        if(addr == address(0)){
            revert ZeroAddressError("Zero address error.");
        }
        _;
    }

    /**
     * @notice Authorizes an upgrade to a new implementation contract.
     * @param newImplementation The address of the new implementation contract.
     */
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    /**
     * @notice Initializes the contract by setting up UUPS upgrade ability and ownership.
     */
    function initialize() public initializer {
        __UUPSUpgradeable_init();
        __Ownable_init(msg.sender);
    }

    /**
     * @notice Updates the address of the PowerVoting contract.
     * @param powerVotingAddress The new address of the PowerVoting contract.
     */
    function updatePowerVotingContract(address powerVotingAddress) external override onlyOwner nonZeroAddress(powerVotingAddress) {
        powerVotingContract = powerVotingAddress;
    }

    /**
     * @notice Adds a list of miner IDs for a specific voter.
     * @param minerIds List of miner IDs to be added.
     * @param voter Address of the voter.
     */
    function addMinerIds(uint64[] memory minerIds, address voter) external override {
        if(msg.sender != powerVotingContract) {
            revert PermissionError("Permission error.");
        }
        voterToMinerIds[voter] = minerIds;
        _updateMinerId(voter);
    }

    /**
     * @notice Adds a task with the specified UCAN CID.
     * @param ucanCid The UCAN CID associated with the task.
     */
    function addTask(string calldata ucanCid) external override {
        if (msg.sender != powerVotingContract) {
            revert PermissionError("Permission error.");
        }
        ++_taskId;
        uint256 taskId = _taskId;
        taskIdList.add(taskId);
        taskIdToUcanCid[taskId] = ucanCid;
    }

    /**
     * @notice Retrieves the list of task IDs.
     * @return taskIds List of task IDs.
     */
    function getTasks() external override view returns (uint256[] memory) {
        return taskIdList.values();
    }

    /**
     * @notice Adds an F4 task for the specified voter.
     * @param voter The address of the voter.
     */
    function addF4Task(address voter) external override {
        if(msg.sender != powerVotingContract) {
            revert PermissionError("Permission error.");
        }
        ++_taskId;
        uint256 f4TaskId = _taskId;
        f4TaskIdList.add(f4TaskId);
        f4TaskIdToAddress[f4TaskId] = voter;
    }

    /**
     * @notice Retrieves the list of F4 task IDs.
     * @return An array containing the F4 task IDs.
     */
    function getF4Tasks() external view override returns(uint256[] memory) {
        return f4TaskIdList.values();
    }

    /**
     * @notice Callback function for updating task information.
     * @param voterInfoParam Voter information containing the Ethereum address and other details.
     * @param taskId The ID of the task being updated.
     * @param powerParam Power information associated with the task.
     */
    function taskCallback(VoterInfo calldata voterInfoParam, uint256 taskId, Power calldata powerParam) external onlyInAllowList override {

        address voterAddress = voterInfoParam.ethAddress;
        if (voterAddressToBlockHeight[voterAddress] == block.number) {
            revert StatusError("Has already been updated by other nodes.");
        }

        //Check if it is a task of f4
        if (bytes(voterInfoParam.ucanCid).length != 0) {
            // check if account exist
            bool exist = false;
            uint256  actorIdsLength = voterInfoParam.actorIds.length;
            for (uint256 i = 0; i < actorIdsLength; i++) {
                if (actorIdList[voterInfoParam.actorIds[i]]) {
                    exist = true;
                    break;
                }
            }
            if (!exist && githubAccountList[voterInfoParam.githubAccount]) {
                exist = true;
            }
            if (exist) {
                delete taskIdToUcanCid[taskId];
                delete f4TaskIdToAddress[taskId];
                taskIdList.remove(taskId);
                f4TaskIdList.remove(taskId);
                return;
            }

            // save account to list
            for (uint256 i = 0; i <actorIdsLength; i++) {
                actorIdList[voterInfoParam.actorIds[i]] = true;
            }
            if (bytes(voterInfoParam.githubAccount).length != 0) {
                githubAccountList[voterInfoParam.githubAccount] = true;
            }
        }

        // update voter info
        voterToInfo[voterAddress] = voterInfoParam;

        // update miner id
        _updateMinerId(voterAddress);

        Power memory power = _calcPower(voterAddress, powerParam);
        uint256 id = _getDayId(voterAddress);
        voterToHistoryPower[voterAddress][id] = power;

        // add to voter list for schedule task
        voterList.add(voterAddress);
        voterAddressToBlockHeight[voterAddress] = block.number;

        // delete task id
        delete taskIdToUcanCid[taskId];
        delete f4TaskIdToAddress[taskId];
        taskIdList.remove(taskId);
        f4TaskIdList.remove(taskId);
        emit CreateDelegate(voterAddress,voterInfoParam.actorIds,voterInfoParam.githubAccount);
    }

    /**
     * @notice Removes a voter along with associated task information.
     * @param voterAddress Address of the voter to be removed.
     * @param taskId ID of the task associated with the voter.
     */
    function removeVoter(address voterAddress, uint256 taskId) external override onlyInAllowList nonZeroAddress(voterAddress) {
        VoterInfo storage voterInfo = voterToInfo[voterAddress];
        uint64[] memory actorIds = voterInfo.actorIds;
        uint256 actorIdsLength = actorIds.length;
        for (uint256 i = 0; i < actorIdsLength; i++) {
            actorIdList[actorIds[i]] = false;
        }
        githubAccountList[voterInfo.githubAccount]=false;
        VoterInfo memory newVoterInfo = VoterInfo(new uint64[](0),new uint64[](0),"",address(0),"");
        voterToInfo[voterAddress] = newVoterInfo;
        PowerStatus storage powerStatus = voterToPowerStatus[voterAddress];
        powerStatus.dayId = 0;
        powerStatus.hasFullRound = 0;
        voterList.remove(voterAddress);
        delete taskIdToUcanCid[taskId];
        taskIdList.remove(taskId);
        emit DeleteDelegate(voterAddress,voterInfo.actorIds,voterInfo.minerIds,voterInfo.githubAccount);
    }

    /**
     * @notice Retrieves the power information for a specific voter and day.
     * @param voterAddress Address of the voter.
     * @param id ID of the day for which power information is requested.
     * @return Power structure containing the power information.
     */
    function getPower(address voterAddress, uint256 id) external view override returns(Power memory){
        PowerStatus storage powerStatus = voterToPowerStatus[voterAddress];
        if (powerStatus.hasFullRound == 0 && id > powerStatus.dayId) {
            return Power(0,new bytes[](0),new bytes[](0),0,0);
        }
        return voterToHistoryPower[voterAddress][id];
    }

    /**
     * @notice Updates the node allow list by adding or removing a node.
     * @param nodeAddress Address of the node to be added or removed.
     * @param allow Boolean indicating whether to allow or disallow the node.
     */
    function updateNodeAllowList(address nodeAddress, bool allow) external override onlyOwner nonZeroAddress(nodeAddress) {
        nodeAllowList[nodeAddress] = allow;
    }

    /**
     * @notice Retrieves the list of voter addresses.
     * @return An array containing the addresses of all voters.
     */
    function getVoterAddresses() external override view returns(address[] memory){
        return voterList.values();
    }

    /**
     * @notice Retrieves the information associated with a specific voter.
     * @param voter The address of the voter.
     * @return VoterInfo The information associated with the specified voter.
     */
    function getVoterInfo(address voter) external override view returns(VoterInfo memory){
        return voterToInfo[voter];
    }

    /**
     * @notice Saves the power information associated with a voter.
     * @param voterAddress The address of the voter.
     * @param powerParam The power information to be saved.
     */
    function savePower(address voterAddress, Power calldata powerParam) external onlyInAllowList nonZeroAddress(voterAddress) override {
        if (voterAddressToBlockHeight[voterAddress] == block.number) {
            revert StatusError("Has already been updated by other nodes.");
        }
        Power memory power = _calcPower(voterAddress, powerParam);
        uint256 id = _getDayId(voterAddress);
        voterToHistoryPower[voterAddress][id] = power;
        voterAddressToBlockHeight[voterAddress] = block.number;
    }

    /**
     * @notice Increments the day ID for the specified voter and returns the updated day ID.
     * @param voter The address of the voter.
     * @return The updated day ID.
     */
    function _getDayId(address voter) private returns(uint256){
        PowerStatus storage powerStatus = voterToPowerStatus[voter];
        ++powerStatus.dayId;
        uint256 id = powerStatus.dayId % MAX_HISTORY;
        if (id == 0) {
            id = MAX_HISTORY;
            powerStatus.dayId = 0;
            powerStatus.hasFullRound = 0;
        }
        return id;
    }

    /**
     * @notice Calculates the power for the specified voter address.
     * @param voterAddress The address of the voter.
     * @param power The Power struct containing the power information.
     * @return The updated Power struct.
     */
    function _calcPower(address voterAddress, Power memory power) private returns(Power memory){
        VoterInfo storage voterInfo = voterToInfo[voterAddress];
        uint64[] memory actorList = voterInfo.actorIds;
        uint256 actorIdsLength = actorList.length;
        _filterAndSetMinerIds(voterAddress, voterInfo, voterInfo.minerIds, voterInfo.minerIds.length, actorIdsLength);
        bytes[] memory clientPowerList = new bytes[](actorIdsLength);
        for (uint256 i = 0; i < actorIdsLength; i++) {
            bytes memory clientPower = actorList[i].getClient();
            clientPowerList[i] = clientPower;
        }
        uint64[] memory minerList = voterInfo.minerIds;
        uint256 minerListLength = minerList.length;
        bytes[] memory spPowerList = new bytes[](minerListLength);
        for (uint256 i = 0; i < minerListLength; i++) {
            bytes memory spPower = minerList[i].getSp();
            spPowerList[i] = spPower;
        }
        power.clientPower = clientPowerList;
        power.spPower = spPowerList;
        power.blockHeight = block.number;
        return power;
    }

    /**
     * @notice Updates the miner IDs associated with a voter based on their actor IDs.
     * @param voterAddress The address of the voter.
     */
    function _updateMinerId(address voterAddress) private {
        VoterInfo storage voterInfo = voterToInfo[voterAddress];
        uint64[] storage actorIds = voterInfo.actorIds;
        uint64[] storage minerIds = voterToMinerIds[voterAddress];
        uint256 minerIdsLength = minerIds.length;
        uint256 actorIdsLength = actorIds.length;

        if (actorIdsLength == 0){
            return;
        }

        if (minerIdsLength == 0)  {
            delete voterInfo.minerIds;
            emit UpdateMinerId(voterAddress,minerIds);
            return;
        }

        _filterAndSetMinerIds(voterAddress, voterInfo, minerIds, minerIdsLength, actorIdsLength);
    }

    /**
     * @notice Filters and sets the miner IDs for a voter based on their associated actor IDs.
     * @param voterAddress The address of the voter.
     * @param voterInfo The storage reference to the voter's information.
     * @param minerIds The storage reference to the original miner IDs.
     * @param minerIdsLength The length of the original miner IDs array.
     * @param actorIdsLength The length of the actor IDs array associated with the voter.
     */
    function _filterAndSetMinerIds(address voterAddress, VoterInfo storage voterInfo, uint64[] storage minerIds, uint256 minerIdsLength, uint256 actorIdsLength) private {
        uint64[] memory minerIdsRes = new uint64[](minerIdsLength);
        uint256 index = 0;
        for (uint256 i = 0; i < minerIdsLength; i++) {
            uint64 actorId = minerIds[i].getOwner();
            for (uint256 j = 0; j < actorIdsLength; j++) {
                if (actorId == voterInfo.actorIds[j]) {
                    minerIdsRes[index++] = minerIds[i];
                    break;
                }
            }
        }

        uint64[] memory matchedMinerIds = new uint64[](index);
        for (uint256 i = 0; i < index; i++) {
            matchedMinerIds[i] = minerIdsRes[i];
        }

        voterInfo.minerIds = matchedMinerIds;
        emit UpdateMinerId(voterAddress,voterInfo.minerIds);
        delete voterToMinerIds[voterAddress];
    }

    /**
     * @notice Resolves the Ethereum address to an actor ID.
     * @param addr The Ethereum address to resolve.
     * @return The resolved actor ID.
     */
    function resolveEthAddress(address addr) external view returns (uint64) {
        uint64 actorId = addr.resolveEthAddress();
        return actorId;
    }

}
