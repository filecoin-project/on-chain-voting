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

export interface ProposalFilter {
  value: number
  label: string
  bgColor: string
  textColor: string
  borderColor: string
  dotColor: string
}
export interface ProposalOption {
  name: string
  count: number,
  disabled?: boolean
}

export interface ProposalResult {
  id: number
  proposalId: number
  optionName?: string
  votes: number,
  vote_result: string,
}

export interface ProposalVotes extends ProposalResult {
  voterAddress: string,
  votedResult: string,
  blockHeight: string,
  percentage: string,
  clientPower: string,
  developerPower: string,
  tokenHolderPower: string,
  spPower: string,
  votedTime: string,
  totalSpPower: string,
  totalClientPower: string,
  totalDeveloperPower: string,
  totalTokenHolderPower: string,
  spPowerPercent: number,
  clientPowerPercent: number,
  developerPowerPercent: number,
  tokenHolderPowerPercent: number
}

export interface ProposalData {
  id: number
  cid: number
  creator: string
  startTime: number
  expTime: number
  proposalType: number
  proposalStatus?: number
  proposalResults: ProposalResult[]
}

export interface ProposalList extends ProposalData {
  option: ProposalOption[]
  title: string
  address: string
  githubName: string
  currentTime: number
  githubAvatar: string
  content: string
  proposalStatus: number
  showTime: string
  GMTOffset: string[]
  voteStatus?: number
  subStatus: number,
  endTime: number,
  startTime: number,
  snapshotInfo: {
    snapshotHeight: number,
    snapshotDay: string
  }
}

export interface ProposalDraft {
  timezone: string
  Time: string
  title: string
  content: string
  chainId: number,
  startTime: number,
  endTime: number,
  spPercentage: number,
  clientPercentage: number,
  developerPercentage: number,
  tokenHolderPercentage: number
}
export interface VotingList {
  proposalId: number,
  cid: string,
  address: string,
  startTime: number,
  endTime: number,
  chainId: number,
  title: string,
  timezone: string,
  content: string,
  githubName: string,
  githubAvatar: string,
  gmtOffset: string,
  height: number,
  createdAt: number,
  updatedAt: number,
  voteResult: Array<any>,
  time: Array<any>,
  option: Array<string>,
  showTime: Array<any>,
  status: number,
  voteCount: number,
  voteCountDay: string
  votePercentage: {
    approve: number,
    reject: number
  }
}