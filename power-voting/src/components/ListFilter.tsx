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

// @ts-ignore
export default function ListFilter ({ name, value, list, onChange, }) {

  return (
    <div className='flex text-base pt-6 pb-5'>
      <div className='text-[#7F8FA3]'>{name}:</div>
      <div className='flex'>
        {list.map((item: any, index: number) => {
          return (
            <button
              onClick={() => onChange(item.value)}
              type='button'
              key={index}
              className={`ml-[20px]  hover:text-blue-300 cursor-pointer relative ${value === item.value
                  ? 'text-white before:absolute before:inset-x-0 before:-top-6 before:h-1 before:bg-[#2DA1F7]'
              : 'text-[#7F8FA3]'}`}
            >
              {item.label}
            </button>
          )
        })}
      </div>
    </div>
  )
}
