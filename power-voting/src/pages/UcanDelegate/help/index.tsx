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

import type { FC } from "react";
import React from "react";
import { useTranslation } from 'react-i18next';
import { Link, useLocation, useNavigate } from "react-router-dom";
import MDEditor from "../../../components/MDEditor";
import './index.less';
const Index: FC = () => {
  const location = useLocation();
  const doc = location.state?.doc;
  const navigate = useNavigate();
  const { t } = useTranslation();
 
  return (
    <div>
      <div className="px-3 mb-6 md:px-0">
        <button>
          <div className="inline-flex items-center gap-1 text-skin-text hover:text-skin-link">
            <Link to="#" onClick={() => navigate(-1)} className="flex items-center">
              <svg className="mr-1" viewBox="0 0 24 24" width="1.2em" height="1.2em">
                <path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                      d="m11 17l-5-5m0 0l5-5m-5 5h12"></path>
              </svg>
              {t('content.back')}
            </Link>
          </div>
        </button>
      </div>
      <MDEditor
        className="doc"
        value={doc}
        moreButton={false}
        readOnly={true}
        view={{ menu: false, md: false, html: true, both: false, fullScreen: true, hideMenu: false }}
        onChange={() => { }}
      />
    </div>
  )
}

export default Index;