// Copyright (C) 2023-2024 StorSwift Inc.
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

package mock

import (
	"context"

	"powervoting-server/model"
	"powervoting-server/service"
)

type MockSyncService struct {
}

// AddProposal implements service.ISyncService.
func (m *MockSyncService) AddProposal(ctx context.Context, in *model.ProposalTbl) error {
	panic("unimplemented")
}

// AddVote implements service.ISyncService.
func (m *MockSyncService) AddVote(ctx context.Context, in *model.VoteTbl) error {
	panic("unimplemented")
}

// BatchUpdateVotes implements service.ISyncService.
func (m *MockSyncService) BatchUpdateVotes(ctx context.Context, votes []model.VoteTbl) error {
	return nil
}

// CreateSyncEventInfo implements service.ISyncService.
func (m *MockSyncService) CreateSyncEventInfo(ctx context.Context, in *model.SyncEventTbl) error {
	return nil
}

// GetSyncEventInfo implements service.ISyncService.
func (m *MockSyncService) GetSyncEventInfo(ctx context.Context, addr string) (*model.SyncEventTbl, error) {
	panic("unimplemented")
}

// GetUncountedVotedList implements service.ISyncService.
func (m *MockSyncService) GetUncountedVotedList(ctx context.Context, chainId int64, proposalId int64) ([]model.VoteTbl, error) {
	return []model.VoteTbl{
		{
			ProposalId: 1,
			ChainId:    1,
			Address:    "0x1234567890123456789012345678901234567890",
			VoteEncrypted: `-----BEGIN AGE ENCRYPTED FILE-----
YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IHRsb2NrIDE2MjU3MDEyIDUyZGI5YmE3
MGUwY2MwZjZlYWY3ODAzZGQwNzQ0N2ExZjU0Nzc3MzVmZDNmNjYxNzkyYmE5NDYw
MGM4NGU5NzEKbVEvRW8zR3gxZS9xTXhPb0kyT080MFhLdm5INUZiK0FQYW1qSSs3
WnlKeGVwbkMwYU1zV29LanhJNEljcFBNTApEL0pFNEdGK3hhSkdxS1FkbmhzdWpZ
WUtrcFM3WE5qUi9mSUFDRCttaEZiZzdCTUlYa2o4c0ZuSG9YZHUvWnBxCmVSWnpS
YnUweHlHN2x4bjJ5d1BuVmZvMGtQMC9OMlpMcS9hM0g2U3E3NlEKLS0tIFd5N1JJ
TnBkVlh3MXo0UWhLSEg2NnFHdDd4ZHRPZlJsRVQ1ZENvN01yRVkKZxIPZbUHXxCp
WJFtIRfXNUjhc/ClVYYKp7nrUcnWc/soVYTO3q7adl6i738=
-----END AGE ENCRYPTED FILE-----
`,
			SpPower:          "0",
			TokenHolderPower: "1000",
			ClientPower:      "0",
			DeveloperPower:   "0",
		},
		{
			ProposalId: 1,
			ChainId:    1,
			Address:    "0x1234567890123456789012345678901234567891",
			VoteEncrypted: `-----BEGIN AGE ENCRYPTED FILE-----
YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IHRsb2NrIDE2MzI5MjEyIDUyZGI5YmE3
MGUwY2MwZjZlYWY3ODAzZGQwNzQ0N2ExZjU0Nzc3MzVmZDNmNjYxNzkyYmE5NDYw
MGM4NGU5NzEKdVlwMzF6RzF0REdYbk5QcHpPam5mOFl1dE9CSFdHZGRuYVIzMllQ
UVZkbjBTWEVOWEZROUROQVpXdERuSFFQYQpBL1NsdXdjS25ZRXN6OU0rQ3JKY3dx
NnFxelY4Y2Nvc2JnMnZrdHMvdU1lVldvVXhENVdtRE9yTHRzNjFJVlY0CmM4Wjg0
ekxoazYyT2psYnNDUWNvdjBabjl3NHFBZHNZSlNJRmNyeFhDZzgKLS0tIG5qcDZX
VnBxT2VlUWphZTA5MWNueVZ5eXRERGpVUURTNEN6VG83QnZ5dEEKlsE1n8mpI4w4
ETGnT966XLd1P76JQujjdDMUb4yJdRgFGDHWkVH3SjPY1WD2
-----END AGE ENCRYPTED FILE-----
`,
			SpPower:          "0",
			TokenHolderPower: "9000",
			ClientPower:      "0",
			DeveloperPower:   "0",
		},
	}, nil
}

// UncountedProposalList implements service.ISyncService.
func (m *MockSyncService) UncountedProposalList(ctx context.Context, chainId int64, endTime int64) ([]model.ProposalTbl, error) {
	return []model.ProposalTbl{
		{
			ProposalId: 1,
			ChainId:    1,
		},
	}, nil
}

// UpdateProposal implements service.ISyncService.
func (m *MockSyncService) UpdateProposal(ctx context.Context, in *model.ProposalTbl) error {
	return nil
}

// UpdateSyncEventInfo implements service.ISyncService.
func (m *MockSyncService) UpdateSyncEventInfo(ctx context.Context, addr string, height int64) error {
	panic("unimplemented")
}

var _ service.ISyncService = (*MockSyncService)(nil)
