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

import React, { useState, useEffect } from 'react';
import MdEditor from 'react-markdown-editor-lite';
// @ts-ignore
import markdownIt from 'markdown-it';
// @ts-ignore
import { full as emoji } from 'markdown-it-emoji';
// @ts-ignore
import footnote from 'markdown-it-footnote';
// @ts-ignore
import mdKatex from '@iktakahiro/markdown-it-katex';
// @ts-ignore
import subscript from 'markdown-it-sub';
// @ts-ignore
import superscript from 'markdown-it-sup';
// @ts-ignore
import deflist from 'markdown-it-deflist';
// @ts-ignore
import abbreviation from 'markdown-it-abbr';
// @ts-ignore
import insert from 'markdown-it-ins';
// @ts-ignore
import mark from 'markdown-it-mark';
// @ts-ignore
import tasklists from 'markdown-it-task-lists';
// @ts-ignore
import anchor from 'markdown-it-anchor';
// @ts-ignore
import linkAttributes  from 'markdown-it-link-attributes'; 

import 'katex/dist/katex.css';
import 'react-markdown-editor-lite/lib/index.css';
import './index.less';
const slugify = (text: string) => {
  return text
    .toString()
    .toLowerCase()
    .trim() //Trim leading and trailing whitespace
    .replace(/\s+/g, '-')     // Replace spaces with hyphens
    .replace(/[^\w-]+/g, '') //Remove any characters that are not word characters or hyphens
    .replace(/--+/g, '-') // Replace multiple consecutive hyphens with a single hyphen
    .replace(/^-+/, '') //Remove leading hyphens
    .replace(/-+$/, ''); // Remove trailing hyphens
};
const mdParser = markdownIt({
  html: true,
  linkify: true,
  typographer: true,
})
  .use(mdKatex)
  .use(emoji)
  .use(footnote)
  .use(subscript)
  .use(superscript)
  .use(deflist)
  .use(abbreviation)
  .use(insert)
  .use(mark)
  .use(linkAttributes, {
    pattern: /^https?:\/\//, 
    attrs: {
      target: '_blank',     
      rel: 'noopener' 
    }
  })
  .use(anchor, {
    slugify: slugify,
    permalink: false,
    permalinkClass: 'anchor',
    permalinkSymbol: '#',
    permalinkBefore: true,
    level: [1, 2, 3]
  })
  .use(tasklists);

interface Props {
  value: string;
  onChange: (value: string) => void;
  moreButton?: boolean;
  className?: string;
  style?: React.CSSProperties;
  readOnly?: boolean;
  view?: {
    menu: boolean;
    md: boolean;
    html: boolean;
    both: boolean;
    fullScreen: boolean;
    hideMenu: boolean;
  };
}

const Index: React.FC<Props> = ({ value = '', onChange, ...rest }) => {
  const { moreButton = true, className = '', style = {}, readOnly = false, view = { menu: true, md: true, html: true, both: true, fullScreen: true, hideMenu: true } } = rest;

  const mdEditor: any = React.useRef(null);
  const handleEditorChange = () => mdEditor.current?.getMdValue();
  const [currentValue, setCurrentValue] = useState(value)
  const [showMore, setShowMore] = useState(false);
  const handleClickShowMore = () => {
    setShowMore(prev => !prev);
  };

  const handleInternalLinks = () => {
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
      anchor.addEventListener('click', (e) => {
        e.preventDefault();
        const targetId = anchor.getAttribute('href')?.substring(1);
        if (targetId) {
          const targetElement = document.getElementById(targetId);
          if (targetElement) {
            window.scrollTo({
              top: targetElement.offsetTop,
              behavior: 'smooth',
            });
          }
        }
      });
    });
  }
  useEffect(() => {
    const timeId = setTimeout(() => {
      handleInternalLinks();
    }, 100);
    setCurrentValue(value)
    setShowMore(moreButton && value.length > 800)
    return () => clearTimeout(timeId);
  }, [moreButton, value]);
  const renderMoreButton = () => {
    if (value.length > 800) {
      return (
        <>
          <div className={`absolute bottom-0 h-[200px] w-full bg-gradient-to-t from-[#F6F6F6] ${showMore ? 'flex' : 'hidden'}`} />
          <div className={`flex w-full justify-center  ${showMore ? 'absolute -bottom-[50px]' : ''}`}>
            <button className="focus:outline-none border-[#DFDFDF] hover:border-[#DFDFDF] border-[1px] border-solid text-[#575757] mt-4 self-center rounded-xl py-2 px-4 font-semibold" onClick={handleClickShowMore}>
              {showMore ? "Show More" : "Show Less"}
            </button>
          </div>
        </>
      );
    }
  };
  return (
    readOnly ?
      <div className='relative'>
        <MdEditor
          className={`rcmd scrollD  ${className} ${moreButton && showMore ? 'mb-20' : ''}`}
          value={value}
          readOnly={readOnly}
          style={{ ...style, maxHeight: moreButton && showMore ? '70vh' : 'fit-content' }}
          renderHTML={text => {
            return mdParser.render(text)
          }}
          view={view}
        />
        {moreButton && renderMoreButton()}
        <button id="back-to-top-btn" className="fixed bottom-[11rem] z-40 right-[30%] p-2 bg-gray-600 text-white rounded-full hover:bg-gray-700 focus:outline-none ">
          <svg className="w-6 h-6" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
            <path d="M12 3l-8 8h5v10h6V11h5z" fill="currentColor" />
          </svg>
        </button>
      </div>
      :
      <MdEditor
        className={`MDEditor rcmd scrollD ${className}`}
        ref={mdEditor}
        value={currentValue}
        style={style}
        renderHTML={text => mdParser.render(text)}
        onChange={(v) => {
          setCurrentValue(v.text)
          onChange(handleEditorChange())
        }}
      />
  )
}

export default Index