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

import FipEditorList from "../../src/pages/Fip/fipEditorList";
import GistDelegateList from "../../src/pages/GistDelegate/list";
import Landing from "../../src/pages/Landing";
import CreateVote from "../pages/CreateVote";
import FipEditorApprove from "../pages/Fip/approve";
import FipEditorPropose from "../pages/Fip/propose";
import FipEditorRevoke from "../pages/Fip/revoke";
import GistDelegateAdd from "../pages/GistDelegate/add/index";
import Home from "../pages/Home";
import MinerId from "../pages/MinerId";
import Vote from "../pages/Vote";
import VotingResults from "../pages/VotingResults";

const routes = [
  {
    path: "/",
    element: <Landing />,
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
    path: "/fip-editor/fipEditorList",
    element: <FipEditorList />,
  },
  {
    path: "/gistDelegate/add",
    element: <GistDelegateAdd />,
  },
  {
    path: "/gistDelegate/list",
    element: <GistDelegateList />,
  },
  {
    path: "/vote/:id",
    element: <Vote />,
  },
  {
    path: "/votingResults/:id",
    element: <VotingResults />,
  },
  {
    path: "*",
    element: <Home />,
  }
]
export default routes;
