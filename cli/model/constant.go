package model

const (
	InvokeContract = 3844450837

	Page     = 1
	PageSize = 10

	Approved = "Approved"
	Rejected = "Rejected"

	Approve             = "approve"
	Reject              = "reject"
	BaseProposalAPIPath = "/api"

	GithubAPI = "https://api.github.com/users/"
)

const (
	ProposalStatusPending ProposalStatus = iota + 1
	ProposalStatusInProgress
	ProposalStatusCounting
	ProposalStatusCompleted
)

// ProposalStatus defines the various states a proposal can have.
type ProposalStatus int

// String returns the string representation of the ProposalStatus.
func (status ProposalStatus) String() string {
	switch status {
	case ProposalStatusPending:
		return "Upcoming"
	case ProposalStatusInProgress:
		return "In Progress"
	case ProposalStatusCounting:
		return "Vote Counting"
	case ProposalStatusCompleted:
		return "Complete"
	default:
		return "Unknown"
	}
}
