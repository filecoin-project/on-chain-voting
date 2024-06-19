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

import type { ReactNode } from 'react';
import React from 'react';
import { QuestionCircleOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { FILECOIN_AUTHORIZE_DOC, FILECOIN_DEAUTHORIZE_DOC, GITHUB_AUTHORIZE_DOC, GITHUB_DEAUTHORIZE_DOC } from "../common/consts";

export default function Table({ title = '', link = {} as { type: string, action: string, href: string }, list = [] as { name: string, hide?: boolean, comp: ReactNode, width?: number, desc?: ReactNode, }[], subTitle = <div/> }) {
  const navigate = useNavigate();
  const { type, action, href } = link;

  const handleJump = () => {
    let doc = '';
    if (type === 'filecoin') {
      doc = action === 'authorize' ? FILECOIN_AUTHORIZE_DOC : FILECOIN_DEAUTHORIZE_DOC;
    } else {
      doc = action === 'authorize' ? GITHUB_AUTHORIZE_DOC : GITHUB_DEAUTHORIZE_DOC;
    }
    navigate(href, {
      state: {
        doc
      }
    });
  }

  return (
    <table className='min-w-full bg-[#FFFFFF] rounded text-left'>
      <thead>
        <tr>
          <th scope='col' colSpan={2}>
            <div className='font-normal text-black px-8 py-7 text-2xl border-b border-[#313D4F] flex items-center'>
              <span>{title}</span>
             
              {
                href && (
                  <div className='flex items-start cursor-pointer' onClick={handleJump}>
                    <QuestionCircleOutlined className='text-[#8896AA] text-[16px] ml-2' />
                  </div>
                )
              }
            </div>
            <div className='px-8'>
            {subTitle && (
                <span className='text-[#4B535B]'>{subTitle}</span>
              )}
            </div>
          </th>
        </tr>
      </thead>
      <tbody className='divide-y divide-[#111111]'>
        {list.filter((item: { name: string, hide?: boolean, comp: ReactNode, width?: number, desc?: ReactNode }) => !item.hide).map((item: { name: string, hide?: boolean, comp: ReactNode, width?: number, desc?: ReactNode }) => (
          <tr key={item.name} className='divide-x divide-[#313D4F]  '>
            <td className={`${item.width ? `w-[${item.width}px]` : 'w-[280px]'} whitespace-nowrap py-9 px-8 text-base text-[#313D4F] align-top`}>
              {item.name}
              {item.desc && <div className={`${item.width ? `w-[${item.width}px]` : 'w-[280px]'} text-sm whitespace-normal text-[#4B535B] mt-8`}>
                {item.desc}
              </div>}

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
