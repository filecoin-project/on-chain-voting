/**
 * format 路由表
 */
import React from "react";
import Home from "../pages/Home";
import CreatePoll from "../pages/CreatePoll";
import Vote from "../pages/Vote";
import VotingResults from "../pages/VotingResults";
import PVDocument from "src/pages/Documents/PVDocument";

const routes = [
  {
    path: "/",
    element: <Home />,
  },
  {
    path: "/createPoll",
    element: <CreatePoll />,
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
    path: '/document',
    element: <PVDocument />
  },
  {
    path: "*",
    element: <Home />,
  }
]
export default routes;
