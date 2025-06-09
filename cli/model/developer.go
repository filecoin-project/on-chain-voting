package model

type SignatureData struct {
	WalletAddress string `json:"walletAddress"`
	GithubName    string `json:"githubName"`
	Timestamp     int64  `json:"timestamp"`
}
