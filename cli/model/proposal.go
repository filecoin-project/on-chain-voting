package model

// ProposalsResponse represents the response format for a list of proposals.
type ProposalsResponse struct {
	Code    int    `json:"code"`    // Response code
	Message string `json:"message"` // Response message
	Data    struct {
		Total int        `json:"total"` // Total number of proposals
		List  []Proposal `json:"list"`  // List of proposals
	} `json:"data"`
}

// Proposal represents the structure of a proposal within the system.
type Proposal struct {
	ProposalId     int64          `json:"proposalId"` // Unique identifier for the proposal
	Address        string         `json:"address"`    // Address of the proposer
	GithubName     string         `json:"githubName"` // Github username associated with the proposal
	StartTime      int64          `json:"startTime"`  // Start time of the proposal in Unix timestamp
	EndTime        int64          `json:"endTime"`    // End time of the proposal in Unix timestamp
	ChainId        int            `json:"chainId"`    // Chain ID associated with the proposal
	Title          string         `json:"title"`      // Title of the proposal
	Content        string         `json:"content"`    // Content of the proposal
	CreatedAt      int64          `json:"createdAt"`  // Proposal creation time in Unix timestamp
	UpdatedAt      int64          `json:"updatedAt"`  // Last update time of the proposal in Unix timestamp
	Voted          bool           `json:"voted"`      // Whether the proposal has been voted on
	Status         ProposalStatus `json:"status"`     // Current status of the proposal (e.g., pending, active, completed)
	VotePercentage VotePercentage `json:"votePercentage"`
	SnapshotInfo   SnapshotInfo   `json:"snapshotInfo"`
	Percentage     Percentage     `json:"percentage"`
	TotalPower     TotalPower     `json:"totalPower"`
}

// VotePercentage represents the percentage of votes for approve and reject actions.
type VotePercentage struct {
	Approve float64 `json:"approve"` // Percentage of approve votes
	Reject  float64 `json:"reject"`  // Percentage of reject votes
}

// SnapshotInfo contains details of the snapshot when the proposal was created.
type SnapshotInfo struct {
	SnapshotDay    string `json:"snapshotDay"`    // Day of the snapshot (e.g., "2025-04-01")
	SnapshotHeight int    `json:"snapshotHeight"` // Block height of the snapshot
}

// Percentage represents the percentage breakdown of different stakeholder categories.
type Percentage struct {
	TokenHolderPercentage int `json:"tokenHolderPercentage"` // Percentage of token holders
	SpPercentage          int `json:"spPercentage"`          // Percentage of service providers
	ClientPercentage      int `json:"clientPercentage"`      // Percentage of clients
	DeveloperPercentage   int `json:"developerPercentage"`   // Percentage of developers
}

// TotalPower represents the total power of various stakeholders (e.g., service providers, token holders).
type TotalPower struct {
	SpPower          string `json:"spPower"`          // Total power of service providers
	TokenHolderPower string `json:"tokenHolderPower"` // Total power of token holders
	DeveloperPower   string `json:"developerPower"`   // Total power of developers
	ClientPower      string `json:"clientPower"`      // Total power of clients
}

// VoteResponse represents the response format for a list of votes on proposals.
type VoteResponse struct {
	Code    int    `json:"code"`    // Response code
	Message string `json:"message"` // Response message
	Data    []Vote `json:"data"`    // List of votes
}

// Vote represents the details of a vote on a proposal.
type Vote struct {
	ProposalId       int    `json:"proposalId"`       // ID of the proposal being voted on
	ChainId          int    `json:"chainId"`          // Chain ID
	VoterAddress     string `json:"voterAddress"`     // Address of the voter
	VotedResult      string `json:"votedResult"`      // Result of the vote (e.g., approve, reject)
	Percentage       string `json:"percentage"`       // Voter's voting percentage
	VotedTime        int    `json:"votedTime"`        // Time the vote was cast (Unix timestamp)
	DeveloperPower   string `json:"developerPower"`   // Power of the developer stakeholder
	SpPower          string `json:"spPower"`          // Power of the service provider stakeholder
	ClientPower      string `json:"clientPower"`      // Power of the client stakeholder
	TokenHolderPower string `json:"tokenHolderPower"` // Power of the token holder stakeholder
}
