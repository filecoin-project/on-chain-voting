export interface ProposalFilter {
  value: number
  color: string,
  label: string
}
export interface ProposalOption {
  name: string
  count: number,
  disabled?: boolean
}

export interface ProposalResult {
  id: number
  proposalId: number
  optionId: number
  optionName?: string
  votes: number
}

export interface ProposalHistory extends ProposalResult {
  address: string
}

export interface ProposalData {
  id: number
  cid: number
  creator: string
  expTime: number
  proposalType: number
  proposalResults: ProposalResult[]
}

export interface ProposalList extends ProposalData {
  option: ProposalOption[]
  name: string
  address: string
  descriptions: string
  voteType: number
  proposalStatus: number
  showTime: string
  GMTOffset: string[]
  voteStatus?: number
}