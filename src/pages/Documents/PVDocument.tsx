import React from "react";
import { FC } from "react";
import { Link } from "react-router-dom";
import { PVText } from "./PVText/PVText";
import MDEditor from "../../components/MDEditor";
import './pv.less';

const PVDocument: FC = () => {
  const pvText = new PVText();
 
  return (
    <div>
      <div className="px-3 mb-6 md:px-0">
        <button>
          <div className="inline-flex items-center gap-1 text-skin-text hover:text-skin-link">
            <Link to="/" className="flex items-center">
              <svg className="mr-1" viewBox="0 0 24 24" width="1.2em" height="1.2em">
                <path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                      d="m11 17l-5-5m0 0l5-5m-5 5h12"></path>
              </svg>
              Back
            </Link>
          </div>
        </button>
      </div>
      <MDEditor
        className="pvdoc"
        value={pvText.pvText()}
        readOnly={true}
        view={{ menu: false, md: false, html: true, both: false, fullScreen: true, hideMenu: false }}
        onChange={() => { }}
      />
    </div>
  )
}

export default PVDocument;