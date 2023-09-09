package model

import "math/big"

type GraphQLRequestBody struct {
	Query string `json:"query"`
}

type GraphQLProposalResponse struct {
	Data struct {
		Proposals []ProposalRes `json:"proposals"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type ProposalRes struct {
	ProposalID int64 `json:"proposalId"`
}

type GraphQLVoteResponse struct {
	Data struct {
		Votes []VoteInfo `json:"votes"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type VoteInfo struct {
	VoteInfo        string `json:"voteInfo"`
	TransactionHash string `json:"transactionHash"`
	Address         string `json:"id"`
}

// Vote VoteResult vote result
type Vote struct {
	OptionId        *big.Int
	Votes           *big.Int
	TransactionHash string
	Address         string
}
