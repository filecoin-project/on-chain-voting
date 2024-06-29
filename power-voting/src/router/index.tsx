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
import Home from "../pages/Home";
import CreateVote from "../pages/CreateVote";
import UcanDelegateAdd from "../pages/UcanDelegate/add/index";
import UcanDelegateDelete from "../pages/UcanDelegate/delete/index";
import UcanDelegateHelp from "../pages/UcanDelegate/help/index";
import Vote from "../pages/Vote";
import VotingResults from "../pages/VotingResults";
import MinerId from "../pages/MinerId";
import FipEditorPropose from "../pages/Fip/propose";
import FipEditorApprove from "../pages/Fip/approve";
import FipEditorRevoke from "../pages/Fip/revoke";
// import Landing from "src/pages/Landing";

const routes = [
  {
    path: "/",
    element: <Home />,
  },
  {
    path: "/home",
    element: <Home />,
  },
  {
    path: "/createVote",
    element: <CreateVote />,
  },
  {
    path: "/minerid",
    element: <MinerId />,
  },
  {
    path: "/fip-editor/propose",
    element: <FipEditorPropose />,
  },
  {
    path: "/fip-editor/approve",
    element: <FipEditorApprove />,
  },
  {
    path: "/fip-editor/revoke",
    element: <FipEditorRevoke />,
  },
  {
    path: "/ucanDelegate/add",
    element: <UcanDelegateAdd />,
  },
  {
    path: "/ucanDelegate/delete",
    element: <UcanDelegateDelete />,
  },
  {
    path: "/ucanDelegate/help",
    element: <UcanDelegateHelp />,
  },
  {
    path: "/vote/:id/:cid",
    element: <Vote />,
  },
  {
    path: "/votingResults/:id/:cid",
    element: <VotingResults />,
  },
  {
    path: "*",
    element: <Home />,
  }
]
export default routes;
