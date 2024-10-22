package response

import (
	"powervoting-server/model"
)

type Proposal struct {
	ProposalId   int64              `json:"proposalId"`   // Proposal ID
	Cid          string             `json:"cid"`          // CID
	Creator      string             `json:"address"`      // Creator address
	StartTime    int64              `json:"startTime"`    // Start time
	ExpTime      int64              `json:"expTime"`      // Expiry time
	Network      int64              `json:"chainId"`      // Network ID
	Name         string             `json:"name"`         // Name
	Timezone     string             `json:"timezone"`     // Timezone
	Descriptions string             `json:"descriptions"` // Descriptions
	GithubName   string             `json:"githubName"`   // Github name
	GithubAvatar string             `json:"githubAvatar"` // Github avatar
	GMTOffset    string             `json:"gmtOffset"`    // GMT offset
	CurrentTime  int64              `json:"currentTime"`  // Current time
	CreatedAt    int64              `json:"createdAt"`    // Created time
	UpdatedAt    int64              `json:"updatedAt"`    // Updated time
	VoteResult   []model.VoteResult `json:"voteResult"`   // Vote result
	Time         []string           `json:"time"`
	Option       []string           `json:"option"`
	ShowTime     []string           `json:"showTime"`
	Status       int                `json:"status"`
	VoteCount    int64              `json:"voteCount"`
	VoteCountDay string             `json:"voteCountDay"`
	Height       int64              `json:"height"`
}
