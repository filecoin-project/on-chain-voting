/**
 * format 路由表
 */
import { Navigate, RouteObject } from "react-router-dom";
import React from "react";
import Home from "../pages/Home";
import AcquireNFT from "../pages/AcquireNFT";
import CreatePoll from "../pages/CreatePoll";
import Vote from "../pages/Vote";
import VotingResults from "../pages/VotingResults";


const routes: RouteObject[] = [
  {
    path: "/",
    element: <Home />,
  },
  {
    path: "/acquirenft",
    element: <AcquireNFT />
  },
  {
    path: "/createPoll",
    element: <CreatePoll />,
  },
  {
    path: "/vote",
    element: <Vote />,
  },
  {
    path: "/votingResults",
    element: <VotingResults />,
  },
  {
    path: "*",
    element: <Home />,
  }
]


export default routes;
