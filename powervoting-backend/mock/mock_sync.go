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

// AddVoterAddress implements service.ISyncService.
func (m *MockSyncService) AddVoterAddress(ctx context.Context, in *model.VoterAddressTbl) error {
	panic("unimplemented")
}

// GetVoterAddresss implements service.ISyncService.
func (m *MockSyncService) GetVoterAddresss(ctx context.Context, height int64) ([]model.VoterAddressTbl, int64, error) {
	panic("unimplemented")
}

// AddProposal implements service.ISyncService.
func (m MockSyncService) AddProposal(ctx context.Context, in *model.ProposalTbl) error {
	panic("unimplemented")
}

// AddVote implements service.ISyncService.
func (m MockSyncService) AddVote(ctx context.Context, in *model.VoteTbl) error {
	panic("unimplemented")
}

// BatchUpdateVotes implements service.ISyncService.
func (m MockSyncService) BatchUpdateVotes(ctx context.Context, votes []model.VoteTbl) error {
	return nil
}

// CreateSyncEventInfo implements service.ISyncService.
func (m MockSyncService) CreateSyncEventInfo(ctx context.Context, in *model.SyncEventTbl) error {
	return nil
}

// GetSyncEventInfo implements service.ISyncService.
func (m MockSyncService) GetSyncEventInfo(ctx context.Context, addr string) (*model.SyncEventTbl, error) {
	panic("unimplemented")
}

// GetUncountedVotedList implements service.ISyncService.
func (m MockSyncService) GetUncountedVotedList(ctx context.Context, chainId int64, proposalId int64) ([]model.VoteTbl, error) {
	return []model.VoteTbl{
		{
			ProposalId: 1,
			ChainId:    1,
			Address:    "0x1234567890123456789012345678901234567890",
			// approve
			VoteEncrypted: `-----BEGIN AGE ENCRYPTED FILE-----
YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IHRsb2NrIDE2MzI1ODEyIDUyZGI5YmE3
MGUwY2MwZjZlYWY3ODAzZGQwNzQ0N2ExZjU0Nzc3MzVmZDNmNjYxNzkyYmE5NDYw
MGM4NGU5NzEKcExwOUJkK2ZWbFpEMEx5OGZ3OVhCa09pWldRWnpHSGZtWGErcSsw
b3pzdEFON2VNTjU3T3RFbTRFQ2d2bElpMApBV2ZvU3A5cy9Ydjl0YkRrb3p1Tk5v
RjdaRUo0c2RtdWRWQ2hEOENQcm1FN1VBYUxtaXBLR29ScVFKdFFFbHUrCitDT2Vr
SENnYWE4eHNGNUxuVzdwWEhkd25sRzlsSmxXSy9wUVJIaEM2S2cKLS0tIHNUV2Yr
RVdEZlJQWm9yaS9ITnNqcGNuQmxNcUlob3VKYkErRUJDb243eUkK2rehvaQY2kad
04WmD3ZNvrbcQRZwWonF+Ww4UhpwUApwbCjpEx80W5jfIo+p
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
			// reject
			VoteEncrypted: `-----BEGIN AGE ENCRYPTED FILE-----
YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IHRsb2NrIDE2MzI1ODEyIDUyZGI5YmE3
MGUwY2MwZjZlYWY3ODAzZGQwNzQ0N2ExZjU0Nzc3MzVmZDNmNjYxNzkyYmE5NDYw
MGM4NGU5NzEKdFVxWng1VVpiSW1pRjk2ZmJSU2dXRVJrblRjeEFuWjdkblg5VHNx
djNaZkdzYlh1eVBUSUh4NHZabkJWWmFwVwpCa09MNGpZMUVtNkQ2cjdGK0Z6eEE3
Y0VBR3F3SCtabFMyenE5TjRtNEF4d0FoMkFBS1NhbkxXRVozRzBMSlViCnQ4ZEF1
MHhCbHhmTW1LQnJKRXN3NkFIV3h6VjNSNXlHdnlobk03SFY3b3cKLS0tIHFwbHJi
Qnhxd1JvZHVKaElYSVZOa3dMZkt0SHdIc2hMVzA0NHVmSmsvcTQKBVM4IBEuFQcP
l0YZKbPlAmnEhcp3EAwQ84BvVSibhTmIzq/MYdsHnTX/1O8=
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
