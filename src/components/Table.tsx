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

import React, { ReactNode } from 'react';
import { QuestionCircleOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { FILECOIN_AUTHORIZE_DOC, FILECOIN_DEAUTHORIZE_DOC, GITHUB_AUTHORIZE_DOC, GITHUB_DEAUTHORIZE_DOC } from "../common/consts";

export default function Table ({ title = '', link= {} as { type: string, action: string, href: string }, list = [] as { name: string, comp: ReactNode, width?: number }[], subTitle = '' }) {
  const navigate = useNavigate();
  const { type, action, href } = link;

  const handleJump = () => {
    let doc = '';
    if (type === 'filecoin') {
      doc = action === 'authorize' ? FILECOIN_AUTHORIZE_DOC : FILECOIN_DEAUTHORIZE_DOC;
    } else {
      doc = action === 'authorize' ? GITHUB_AUTHORIZE_DOC : GITHUB_DEAUTHORIZE_DOC;
    }
    navigate(href, { state: {
      doc
      }
    });
  }

  return (
    <table className='min-w-full bg-[#273141] rounded text-left'>
      <thead>
        <tr>
          <th scope='col' colSpan={2}>
            <div className='font-normal text-white px-8 py-7 text-2xl border-b border-[#313D4F] flex items-center'>
              <span>{title}</span>
              {subTitle && (
                <span className='text-[#8896AA] text-xl px-1'> - {subTitle}</span>
              )}
              {
                href && (
                  <div className='flex items-start cursor-pointer' onClick={handleJump}>
                    <QuestionCircleOutlined className='text-[#8896AA] text-[16px] ml-2' />
                  </div>
                )
              }
            </div>
          </th>
        </tr>
      </thead>
      <tbody className='divide-y divide-[#313D4F]'>
        {list.map((item: { name: string, comp: ReactNode, width?: number }) => (
          <tr key={item.name} className='divide-x divide-[#313D4F]'>
            <td className={`${item.width ? `w-[${item.width}px]` : 'w-[280px]'} whitespace-nowrap py-9 px-8 text-xl text-[#8896AA]`}>
              {item.name}
            </td>
            <td className='py-5 px-4 text-xl text-white'>
              {item.comp}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
