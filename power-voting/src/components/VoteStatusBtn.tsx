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

import React from 'react';
import { VOTE_LIST } from 'src/common/consts';
import { ProposalFilter } from 'src/common/types';
import { useTranslation } from 'react-i18next';

export default function VoteStatusBtn({ status = 0 }) {
    const { t } = useTranslation();
    const proposal = VOTE_LIST?.find((proposal: ProposalFilter) => proposal.value === status);
    return (
        <div
            style={{ borderColor: proposal?.bgColor, backgroundColor: proposal?.bgColor, color: proposal?.textColor }}
            className={`flex items-center justify-center border-solid h-[32px] px-[12px] rounded-full font-medium text-base`}>
            <div className='rounded-full w-[10px] h-[10px] mr-[5px]' style={{ backgroundColor: proposal?.dotColor }} />
            {proposal?.label && t(proposal?.label)}
        </div>
    )
}
