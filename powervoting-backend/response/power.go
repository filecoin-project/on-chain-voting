package response

type Power struct {
	DeveloperPower   string `json:"developerPower"`   // Developer power
	SpPower          string `json:"spPower"`          // SP power
	ClientPower      string `json:"clientPower"`      // Client power
	TokenHolderPower string `json:"tokenHolderPower"` // Token holder power
	BlockHeight      string `json:"blockHeight"`      // Block height
}

type DataHeightRep struct {
	Day    string `json:"day"`
	Height int64  `json:"blockHeight"`
	NetId  int64  `json:"chainId"`
}
