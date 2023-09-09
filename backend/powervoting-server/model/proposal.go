package model

import "time"

type Proposal struct {
	Id             int64     `json:"id"`
	ExpirationTime int64     `json:"expirationTime"`
	Net            int       `json:"net"`
	ProposalCid    string    `json:"proposalCid"`
	Status         int       `json:"status"`
	CreateTime     time.Time `json:"createTime"`
	UpdateTime     time.Time `json:"updateTime"`
}
