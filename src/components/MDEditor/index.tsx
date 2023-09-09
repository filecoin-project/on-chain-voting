import React, { useState } from 'react'
import MdEditor from 'react-markdown-editor-lite'
import MarkdownIt from 'markdown-it'
import emoji from 'markdown-it-emoji';
import footnote from 'markdown-it-footnote'
// @ts-ignore
import mdKatex from 'markdown-it-katex'
// @ts-ignore
import subscript from 'markdown-it-sub'
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

// import hljs from 'highlight.js'
// import 'highlight.js/styles/atom-one-light.css'
// import 'highlight.js/styles/github.css'
import 'katex/dist/katex.css';
import 'react-markdown-editor-lite/lib/index.css';
import './index.less';

const mdParser = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
})
  .use(emoji)
  .use(mdKatex)
  .use(footnote)
  .use(subscript)
  .use(superscript)
  .use(deflist)
  .use(abbreviation)
  .use(insert)
  .use(mark)
  .use(tasklists);

// @ts-ignore
const Index = ({ value = '', onChange, moreButton = undefined as unknown as boolean, className = '', style = {}, readOnly = false, view = { menu: true, md: true, html: true, both: true, fullScreen: true, hideMenu: true } }) => {
  const mdEditor: any = React.useRef(null);

  const handleEditorChange = () => {
    if (mdEditor.current) {
      return mdEditor.current.getMdValue();
    }
  }

  const [showMore, setShowMore] = useState(moreButton !== false);

  const backToTopButton = document.querySelector("#back-to-top-btn");

  window.addEventListener("scroll", () => {
    if (window.scrollY > 400) {
      backToTopButton?.classList.add("!flex", "scale-100");
    } else {
      backToTopButton?.classList.remove("!flex", "scale-100");
    }
  });

  backToTopButton?.addEventListener("click", () => {
    window.scrollTo({ top: 0, behavior: "smooth" });
  });

  const handleClickShowMore = () => {
    setShowMore((prev) => !prev);
  }

  return (
    readOnly ?
      <div className='relative'>
        <MdEditor
          className={`rcmd scrollD  ${className} ${moreButton && showMore ? 'mb-10' : ''}`}
          value={value}
          readOnly={readOnly}
          style={{ ...style, maxHeight: moreButton && showMore ? '70vh' : 'fit-content' }}
          renderHTML={(text) => {
            return mdParser.render(text)
          }}
          view={view}
        />
        {
          moreButton !== undefined && value.length > 800 ? (
            <>
              <div className={`absolute bottom-0 h-[80px] w-full bg-gradient-to-t from-[#1b2331] ${showMore ? 'flex' : 'hidden'}`} />
              <div className={`flex w-full justify-center  ${showMore ? 'absolute -bottom-5' : ''}`}>
                <button className="border-[#313D4F] hover:border-[#8896AA] border-[1px] border-solid text-white mt-4 self-center rounded-xl py-2 px-4" onClick={handleClickShowMore}>
                  {showMore ? "Show More" : "Show Less"}
                </button>
              </div>
            </>
          ) : null
        }
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
        style={style}
        renderHTML={(text) => mdParser.render(text)}
        onChange={() => { onChange(handleEditorChange()) }}
      />
  )
}

export default Index