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
import { filecoinCalibration, filecoin } from 'wagmi/chains';
export const network = import.meta.env.VITE_CHAIN_NETWORK || 'mainnet';
export const powerVotingMainNetContractAddress = import.meta.env.VITE_POWER_VOTING_MAINNET_CONTRACT_ADDRESS || '';
export const oracleMainNetContractAddress = import.meta.env.VITE_ORACLE_MAINNET_CONTRACT_ADDRESS || '';
export const powerVotingCalibrationContractAddress = import.meta.env.VITE_POWER_VOTING_CALIBRATION_CONTRACT_ADDRESS || '';
export const oracleCalibrationContractAddress = import.meta.env.VITE_ORACLE_CALIBRATION_CONTRACT_ADDRESS || '';
export const powerVotingFipMainNetContractAddress = import.meta.env.VITE_POWER_VOTING_FIP_MAINNET_CONTRACT_ADDRESS || '';
export const powerVotingFipCalibrationContractAddress = import.meta.env.VITE_POWER_VOTING_FIP_CALIBRATION_CONTRACT_ADDRESS || '';
export const walletConnectProjectId = import.meta.env.VITE_WALLET_CONNECT_ID || '';
export const filecoinId = filecoin.id;
export const calibrationChainId = filecoinCalibration.id;
export const filecoinChainRpc = filecoin.rpcUrls.default.http[0];
export const calibrationChainRpc = filecoinCalibration.rpcUrls.default.http[0];
export const githubApi = 'https://avatars.githubusercontent.com';
export const baseUrl = import.meta.env.VITE_BASE_API_URL || ''
export const proposalListApi = `${baseUrl}/proposal/list`;
export const proposalVoteDataApi = `${baseUrl}/proposal/votes`;
export const proposalDraftAddApi = `${baseUrl}/proposal/draft/add`;
export const proposalDraftGetApi = `${baseUrl}/proposal/draft/get`;
export const votePowerGetApi = `${baseUrl}/power/getPower`;
export const getVoteDetail = `${baseUrl}/proposal/details`;
export const getFipListApi = `${baseUrl}/fipEditor/list`;
export const getFipProposalApi = `${baseUrl}/fipProposal/list`;
export const getGistListApi = `${baseUrl}/voter/info`
export const checkGistApi = `${baseUrl}/fipEditor/checkGist`
export const IN_PROGRESS_STATUS = 2;
export const COMPLETED_STATUS = 4;
export const PENDING_STATUS = 1;
export const VOTE_COUNTING_STATUS = 3;
export const VOTE_ALL_STATUS = 7;
export const WRONG_NET_STATUS = 5;
export const STORING_STATUS = 0;
export const PASSED_STATUS = 6;
export const REJECTED_STATUS = 5;
export const VOTE_OPTIONS = ['Approve', 'Reject'];
export const VOTE_LIST = [
  // {
  //   value: WRONG_NET_STATUS,
  //   label: 'content.wrongNetwork',
  //   bgColor: "#FFF3F3",
  //   textColor: "#AA0101",
  //   borderColor: "#FFDBDB",
  //   dotColor: "#FF0000"

  // },
  {
    value: PENDING_STATUS,
    label: "content.pending",
    bgColor: "#FFF1CE",
    textColor: "#7C4300",
    borderColor: "#FFDD87",
    dotColor: "#FFC327"

  },
  {
    value: IN_PROGRESS_STATUS,
    label: 'content.progressing',
    bgColor: "#FFF1CE",
    textColor: "#7C4300",
    borderColor: "#FFDD87",
    dotColor: "#FFC327"
  },
  {
    value: VOTE_COUNTING_STATUS,
    label: 'content.voteCounting',
    bgColor: "#FFF1CE",
    textColor: "#7C4300",
    borderColor: "#FFDD87",
    dotColor: "#FFC327"
  },
  {
    value: COMPLETED_STATUS,
    label: 'content.complete',
    bgColor: "#E7F4FF",
    textColor: "#005292",
    borderColor: "#C3E5FF",
    dotColor: "#0190FF",

  },
  {
    value: STORING_STATUS,
    label: 'content.storing',
    bgColor: "#FFF1CE",
    textColor: "#7C4300",
    borderColor: "#FFDD87",
    dotColor: "#FFC327"
  },
  {
    value: PASSED_STATUS,
    label: 'content.passed',
    bgColor: "#E3FFEE",
    textColor: "#006227",
    borderColor: "#87FFBE",
    dotColor: "#00C951"
  },
  {
    value: REJECTED_STATUS,
    label: 'content.rejected',
    bgColor: "#FFF3F3",
    textColor: "#AA0101",
    borderColor: "#FFDBDB",
    dotColor: "#FF0000"
  },
]
export const VOTE_FILTER_LIST = [
  {
    label: 'content.all',
    value: VOTE_ALL_STATUS
  },
  {
    label: 'content.pending',
    value: PENDING_STATUS
  },
  {
    label: 'content.progressing',
    value: IN_PROGRESS_STATUS
  },
  {
    label: 'content.voteCounting',
    value: VOTE_COUNTING_STATUS
  },
  {
    label: 'content.complete',
    value: COMPLETED_STATUS
  }
];


export const GITHUB_OPTIONS = [
  {
    label: 'Github',
    value: 2
  }
];
export const VoteOptionItem: { [key: string]: string } = {
  'approve': 'Approve',
  'reject': 'Reject',
}
export const GITHUB_STEP_1 = 1;
export const GITHUB_STEP_2 = 2;
export const FIP_EDITOR_REVOKE_TYPE = 0;
export const FIP_EDITOR_APPROVE_TYPE = 1;
export const DEFAULT_TIMEZONE = Intl.DateTimeFormat().resolvedOptions().timeZone;
export const web3AvatarUrl = 'https://cdn.stamp.fyi/avatar/eth';

export const JWT_HEADER = {
  alg: 'ecdsa',
  type: 'JWT',
  version: '0.0.1'
};
export const OPERATION_CANCELED_MSG = 'content.operationCanceled';
export const STORING_DATA_MSG = 'content.storingChain';
export const STORING_SUCCESS_MSG = 'content.storedSuccessfully';
export const STORING_FAILED_MSG = 'content.dataStoredFailed';
export const VOTE_SUCCESS_MSG = 'content.voteSuccessful';
export const GETTING_POWER_MSG = 'content.gettingPower';
export const CHOOSE_VOTE_MSG = 'content.chooseVote';
export const WRONG_BLOCK_HEIGHT = 'content.blockHeightLoss';
export const WRONG_START_TIME_MSG = 'content.startTimeLoss';
export const WRONG_EXPIRATION_TIME_MSG = 'content.expirationTime';
export const WRONG_MINER_ID_MSG = 'content.checkID';
export const DUPLICATED_MINER_ID_MSG = 'content.iDDuplicated';
export const NOT_FIP_EDITOR_MSG = 'content.fipCreateProposals';
export const NO_FIP_EDITOR_APPROVE_ADDRESS_MSG = 'content.inputAddress';
export const NO_FIP_EDITOR_REVOKE_ADDRESS_MSG = 'content.selectAddress';
export const NO_ENOUGH_FIP_EDITOR_REVOKE_ADDRESS_MSG = 'content.twoFIPRevoke';
export const FIP_ALREADY_EXECUTE_MSG = "content.activeProposal"
export const FIP_APPROVE_SELF_MSG = "content.noPropose"
export const FIP_APPROVE_ALREADY_MSG = "content.addressEditor"
export const HAVE_APPROVED_MSG = 'content.alreadyApproved';
export const HAVE_REVOKED_MSG = 'content.alreadyRevoked';
export const CAN_NOT_REVOKE_YOURSELF_MSG = 'content.revokeYourself';
export const SAVE_DRAFT_SUCCESS = "content.saveSuccess";
export const SAVE_DRAFT_TOO_LARGE = "content.savedDescriptionCharacters";
export const SAVE_DRAFT_FAIL = "content.saveFail";
export const UPLOAD_DATA_FAIL_MSG = "content.saveDataFail";
export const NO_FIP_INfO_MSG = "content.inputFipInfo"
