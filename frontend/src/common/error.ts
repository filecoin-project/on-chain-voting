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

import { ethers } from 'ethers';
import { BaseError } from "wagmi";
import powerVotingAbi from "./abi/power-voting.json"
import { message } from 'antd';
const extractRevertReason = (errorString: string) => {
    const match = errorString.match(/revert reason: (0x[0-9a-fA-F]+)/);
    return match ? match[1] : null;
};
export const parseError = (e: any): string => {
    if (e) {
        // console.log('errorerror===11111111111', e)
        const iface = new ethers.Interface(powerVotingAbi);
        // const decodedError = iface.decodeErrorResult('InvalidProposalPercentageError', e);
        // const decodedError2 = iface.decodeErrorResult('InvalidProposalId', e);

        const reason = extractRevertReason((e as BaseError)?.details)
        const error = iface.parseError(reason ?? "");
        // console.log('errorerror===222222222222222',decodedError,decodedError2, reason, error)
        message.open({
            type: 'warning',
            content: error?.args[0],
        });
        return error?.args[0]
    }
    return "unknow error"
}