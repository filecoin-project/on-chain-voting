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

package service

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"powervoting-server/constant"
	"powervoting-server/model"
	"powervoting-server/model/api"
)

type FipRepo interface {
	CreateFipProposal(ctx context.Context, in *model.FipProposalTbl) (int64, error)
	GetFipProposalListWithPagination(ctx context.Context, req api.FipProposalListReq) ([]model.FipProposalVoted, int64, error)
	GetUnpassFipProposalList(ctx context.Context, chainId int64) ([]model.FipProposalTbl, error)
	UpdateFipProposalVoteByAddress(ctx context.Context, proposalId int64, address string) error
	UpdateFipEditorByAddress(ctx context.Context, address string) error
	UpdateStatusAndGetFipProposal(ctx context.Context, proposalId, chainId int64) (*model.FipProposalTbl, error)
	CreateFipProposalVote(ctx context.Context, in *model.FipProposalVoteTbl) (int64, error)
	CreateFipEditor(ctx context.Context, in *model.FipEditorTbl) (int64, error)

	GetValidFipEditorList(ctx context.Context, req api.FipEditorListReq) ([]model.FipEditorTbl, error)
	GetFipEditorCount(ctx context.Context, chainId int64) (int64, error)
	GetFipProposalVoteCount(ctx context.Context, chainId, proposalId int64) (int64, error)
}

type IFipService interface {
	GetFipProposalList(ctx context.Context, req api.FipProposalListReq) (*api.CountListRep, error)
	GetFipEditorList(ctx context.Context, req api.FipEditorListReq) ([]api.FipEditorRep, error)
}

type FipService struct {
	repo FipRepo
}

func NewFipService(repo FipRepo) *FipService {
	return &FipService{
		repo: repo,
	}
}

func (f *FipService) GetFipProposalList(ctx context.Context, req api.FipProposalListReq) (*api.CountListRep, error) {
	fipProposalList, total, err := f.repo.GetFipProposalListWithPagination(ctx, req)
	if err != nil {
		zap.L().Error("GetFipProposalList error", zap.Error(err))
		return nil, err
	}

	editorCount, err := f.repo.GetFipEditorCount(ctx, req.ChainId)
	if err != nil {
		zap.L().Error("GetFipEditorCount error", zap.Error(err))
		return nil, err
	}

	var res []api.FipProposalRep
	for _, fip := range fipProposalList {
		proposalVotedCount, err := f.repo.GetFipProposalVoteCount(ctx, req.ChainId, fip.ProposalId)
		if err != nil {
			zap.L().Error("GetFipProposalVoteCount error", zap.Error(err))
		}

		if fip.Status == constant.FipProposalPass {
			editorCount = proposalVotedCount
		}

		var votedAddress []string
		err = json.Unmarshal([]byte(fip.Voters), &votedAddress)
		if err != nil {
			zap.L().Error("Unmarshal voters error", zap.Error(err))
		}

		res = append(res, api.FipProposalRep{
			ProposalId:       fip.ProposalId,
			ChainId:          fip.ChainId,
			ProposalType:     fip.ProposalType,
			Creator:          fip.Creator,
			CandidateAddress: fip.CandidateAddress,
			CandidateInfo:    fip.CandidateInfo,
			Timestamp:        fip.Timestamp,
			VotedCount:       proposalVotedCount,
			EditorCount:      editorCount,
			Status:           fip.Status,
			VotedAddresss:    votedAddress,
		})
	}

	return &api.CountListRep{
		Total: total,
		List:  res,
	}, nil
}

func (f *FipService) GetFipEditorList(ctx context.Context, req api.FipEditorListReq) ([]api.FipEditorRep, error) {
	fipEditors, err := f.repo.GetValidFipEditorList(ctx, req)
	if err != nil {
		zap.L().Error("GetFipEditorList error", zap.Error(err))
		return nil, err
	}

	var res []api.FipEditorRep
	for _, fipEditor := range fipEditors {
		res = append(res, api.FipEditorRep{
			Editor:    fipEditor.Editor,
			ChainId:   fipEditor.ChainId,
			Timestamp: fipEditor.UpdatedAt.Unix(),
		})
	}

	return res, err
}

