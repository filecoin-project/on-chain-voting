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

export default function ListFilter (props: { name: string, value: number, list: { label: string, value: number }[], onChange: (status: number) => Promise<void> }) {

  const { value, list, onChange } = props;

  return (
    <div className='flex text-base py-4'>
      <div className='flex'>
        {list.map((item: any, index: number) => {
          return (
            <button
              onClick={() => onChange(item.value)}
              type='button'
              key={index}
              className={`ml-[20px]  hover:text-blue-300 cursor-pointer relative ${value === item.value
                  ? 'text-#005292 before:absolute before:inset-x-0 before:-bottom-4 before:h-1 before:bg-[#2DA1F7]'
              : 'text-[#4B535B]'}`}
            >
              {item.label}
            </button>
          )
        })}
      </div>
    </div>
  )
}
