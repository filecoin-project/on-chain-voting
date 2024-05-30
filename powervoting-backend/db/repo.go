package db

import "powervoting-server/model"

type DataRepo interface {
	GetProposalList(network int64, timestamp int64) ([]model.Proposal, error)
	GetVoteList(network, proposalId int64) ([]model.Vote, error)
	VoteResult(proposalId int64, history model.VoteCompleteHistory, result []model.VoteResult) (int64, error)

	GetDict(key string) (*model.Dict, error)
	CreateDict(in *model.Dict) (int64, error)
	UpdateDict(name string, value string) error

	CountVotes(filter map[string]any) (int64, error)
	UpdateVoteInfo(filter map[string]any, in string) error
	CreateVote(in *model.Vote) (int64, error)

	CountProposal(filter map[string]any) (int64, error)
	CreateProposal(in *model.Proposal) (int64, error)
}
