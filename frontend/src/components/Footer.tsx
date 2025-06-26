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

import { useTranslation } from 'react-i18next';
const Footer = () => {
  const { t } = useTranslation();

  const resources = [
    {
      href: "",
      text: `${t('content.FAQs')} ↗`
    },
    {
      href: "",
      text: `${t('content.documentation')} ↗`
    },
    {
      href: "",
      text: `${t('content.resources')} ↗`
    }
  ]

  const contact = [
    {
      href: "https://github.com/black-domain/power-voting",
      text: "GitHub ↗"
    },
    {
      href: "https://discord.gg/S8NHC7fV26",
      text: `${t('content.discord')} ↗`
    },
    {
      href: "",
      text: `${t('content.slack')} ↗`
    }
  ]
  const legal = [
    {
      href: "",
      text: `${t('content.privacyTerms')} ↗`
    },
    {
      href: "",
      text: `${t('content.codeConduct')} ↗`
    },
  ]

  return (
    <footer className='p-8 bg-[#000000]'>
      <div className='max-w-[1032px] mx-auto md:grid-cols-2 items-center justify-center md:grid-rows-1 grid grid-rows-2 w-full gap-8'>
        <div className='order-2 md:order-1 grid grid-rows-3 gap-4 items-center max-w-[32rem]'>
        <p className='text-sm font-normal text-[#ffffff]'>{t('content.poweredBy')}</p>

        <div className="grid grid-cols-2 gap-2">
          <a target="_blank"
            rel="noopener" href="https://www.storswift.com"><img src="/images/logo_1.png" alt="StorSwift company logo" className='w-[144px] h-[31px] shrink-0' /></a>
          <a target="_blank"
            rel="noopener" href="https://fil.org/"><img src="/images/logo_2.png" alt="Filecoin Foundation logo" className='w-[120px] shrink-0' /></a>
        </div>
        <div style={{
          fontSize: "1.1rem",
          fontWeight: "bold",
          color: "#7F8FA3",
          maxWidth: "32rem",
        }}>
          <p className='text-sm font-normal'>{t('content.allRightReserved')}</p>
        </div>
        </div>
        <div className='order-1 md:order-2 grid grid-cols-3 gap-8'>
        <div >
          <h4 className='text-sm text-[#ffffff] mb-[12px]'>{t('content.partners')}</h4>
          <div className='justify-center text-xs'>
            {resources.map((partner, index) => (
              <a key={index} className='flex items-center hover:text-blue-300 mt-[16px] text-[#989898]' href={partner.href} target='_blank' rel="noreferrer" >
                {partner.text}
              </a>
            ))}
          </div>
        </div>
        <div >
          <h4 className='text-sm text-[#ffffff] mb-[12px]'>{t('content.contactSupport')}</h4>
          <div className='justify-center text-xs'>
            {contact.map((partner, index) => (
              <a key={index} className='flex items-center hover:text-blue-300 mt-[16px] text-[#989898]' href={partner.href} target='_blank' rel="noreferrer" >
                {partner.text}
              </a>
            ))}
          </div>
        </div>
        <div >
          <h4 className='text-sm text-[#ffffff] mb-[12px]'>{t('content.legal')}</h4>
          <div className='justify-center text-xs'>
            {legal.map((partner, index) => (
              <a key={index} className='flex items-center hover:text-blue-300 mt-[16px] text-[#989898]' href={partner.href} target='_blank' rel="noreferrer" >
                {partner.text}
              </a>
            ))}
          </div>
          </div>
        </div>
      </div>
    </footer>
  )
};

export default Footer;
