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

import React from "react";

const Footer = () => {
  const partners = [
    {
      href: 'https://protocol.ai',
      text: 'Protocol Labs',
      icon: '/images/protocol.svg',
    },
  ];

  return (
    <footer className='h-[135px] flex px-8 items-center justify-between bg-[#273141]'>
      <div className='flex items-center'>
        <img src="/images/logo.png" alt="" className='w-[100px] mr-8' />
        <div style={{
          fontSize: "1.1rem",
          fontWeight: "bold",
          color: "#7F8FA3",
          maxWidth: "32rem",
        }}> An infrastructure for DAO governance.
          <p className='text-[12px] font-normal'>© 2023 All rights reserved. StorSwift</p>
        </div>
      </div>
      <div className='flex items-center'>
        <div className='mr-6'>
          <h4 className='text-xl text-[#7F8FA3] mb-[12px]'>Partners</h4>
          <div className=' flex justify-center text-xs'>
            {partners.map((partner, index) => (
              <a key={index} className='flex items-center hover:text-blue-300' href={partner.href} target='_blank' >
                <img className='h-[14px] mr-2 ' src={partner.icon} alt="" />
                {partner.text}
              </a>
            ))}
          </div>
        </div>
        <div>
          <h4 className='text-xl text-[#7F8FA3] mb-[4px]'>Contact Us</h4>
          <div className='flex  m-auto'>
            <div className='mr-3'><a href="https://github.com/black-domain/power-voting" target='blank' ><img className='h-[24px]' src="/images/github.svg" alt="" /></a></div>
            <div className=''><a href="https://discord.gg/S8NHC7fV26" target='blank'><img className='h-[24px]' src="/images/discord.svg" alt="" /></a></div>
          </div>
        </div>
      </div>
    </footer>
  )
};

export default Footer;