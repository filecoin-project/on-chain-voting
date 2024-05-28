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

import { PowerAPI } from "filecoin-solidity-api/contracts/v0.8/PowerAPI.sol";
import { PowerTypes } from "filecoin-solidity-api/contracts/v0.8/types/PowerTypes.sol";
import { CommonTypes } from "filecoin-solidity-api/contracts/v0.8/types/CommonTypes.sol";
import { FilAddresses } from "filecoin-solidity-api/contracts/v0.8/utils/FilAddresses.sol";
import { DataCapAPI } from "filecoin-solidity-api/contracts/v0.8/DataCapAPI.sol";
import { MinerAPI } from "filecoin-solidity-api/contracts/v0.8/MinerAPI.sol";
import { MinerTypes } from "filecoin-solidity-api/contracts/v0.8/types/MinerTypes.sol";
import {PrecompilesAPI} from "filecoin-solidity-api/contracts/v0.8/PrecompilesAPI.sol";


library Powers {
    /**
     * @notice Retrieves storage provider information for a specified miner.
     * @param minerID ID of the target miner.
     * @return The raw bytecode representation of the storage provider.
     */
    function getSp(uint64 minerID) external view returns(bytes memory) {
        (,PowerTypes.MinerRawPowerReturn memory sp )= PowerAPI.minerRawPower(minerID);
        return sp.raw_byte_power.val;
    }

    /**
     * @notice Retrieves client balance information for a specified actor.
     * @param actorID ID of the target actor.
     * @return The bytecode representation of the client balance.
     */
    function getClient(uint64 actorID) external view returns(bytes memory) {
        CommonTypes.FilAddress memory result = FilAddresses.fromActorID(actorID);
        (,CommonTypes.BigInt memory clientBalance) = DataCapAPI.balance(result);
        return clientBalance.val;
    }

    /**
     * @notice Retrieves owner information for a specified miner.
     * @param minerID ID of the target miner.
     * @return The ID of the miner's owner.
     */
    function getOwner(uint64 minerID) external view returns(uint64) {
        CommonTypes.FilActorId miner = CommonTypes.FilActorId.wrap(minerID);
        (, MinerTypes.GetOwnerReturn memory result )  = MinerAPI.getOwner(miner);
        uint64 ownerId = PrecompilesAPI.resolveAddress(result.owner);
        return ownerId;
    }

    /**
     * @notice Resolves an Ethereum address and returns the associated ID.
     * @param addr The Ethereum address to resolve.
     * @return The ID associated with the provided Ethereum address.
     */
    function resolveEthAddress(address addr) public view returns (uint64) {
        return PrecompilesAPI.resolveEthAddress(addr);
    }

}
