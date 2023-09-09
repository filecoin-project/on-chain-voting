import React from "react";
import { Link } from 'react-router-dom';

const Footer = () => {

  return (
    <footer className='flex h-[135px]  px-8 items-center justify-between bg-[#273141]'>
      <div className='flex items-center'>
        <img src="/images/logo1.png" alt="" className='w-[100px] mr-8' />
        <div style={{
          fontSize: "1.1rem",
          fontWeight: "bold",
          color: "#7F8FA3",
          maxWidth: "32rem",
        }}> An infrasturcture for DAO governace.
        <p className='text-[12px] font-normal'>© 2023 All rights reserved. StorSwift</p>
        </div>
      </div>
      <div className='flex items-center'>
        <div className='mr-6'>
          <h4 className='text-xl text-[#7F8FA3]  mb-[12px]'>Links</h4>
          <div className=' flex justify-center text-xs'>
            <Link to={'/document'} className='flex items-center hover:text-blue-300'>
              <img className='h-[16px] mr-[6px]' src="/images/document.svg" alt="" />
              Document
            </Link>
          </div>
        </div>
        <div className='mr-6'>
          <h4 className='text-xl text-[#7F8FA3] mb-[12px]'>Partners</h4>
          <div className=' flex justify-center text-xs'>
            <a className='flex items-center   hover:text-blue-300' href="https://protocol.ai" target='_blank' >
              <img className='h-[14px] mr-2 ' src="/images/protocol.svg" alt="" />
              Protocol Labs
            </a>
          </div>
        </div>
        <div>
          <h4 className='text-xl text-[#7F8FA3] mb-[4px]'>Contact Us</h4>
          <div className='flex  m-auto'>
            {/*<div className='mr-3'><a href="https://twitter.com/SwiftNFTMarket" target='blank' ><img className='h-[24px]' src="/images/twitter.svg" alt="" /></a></div>*/}
            <div className='mr-3'><a href="https://github.com/black-domain/power-voting" target='blank' ><img className='h-[24px]' src="/images/github.svg" alt="" /></a></div>
            <div className=''><a href="https://discord.gg/S8NHC7fV26" target='blank'><img className='h-[24px]' src="/images/discord.svg" alt="" /></a></div>
          </div>
        </div>
      </div>
    </footer>
  )

};
export default Footer;
