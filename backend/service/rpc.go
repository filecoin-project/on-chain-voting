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
	"errors"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"powervoting-server/config"
	"powervoting-server/constant"
	"powervoting-server/model"
	"powervoting-server/utils"
)

type RpcService struct {
	voteRepo  VoteRepo
	syncRepo  SyncRepo
	lotusRepo LotusRepo
	logger    *zap.Logger
}

func NewRpcService(voteRepo VoteRepo, syncRepo SyncRepo, lotusRepo LotusRepo) *RpcService {
	return &RpcService{
		voteRepo:  voteRepo,
		syncRepo:  syncRepo,
		lotusRepo: lotusRepo,
		logger:    zap.L().With(zap.String("service", "RPC")),
	}
}

// GetAllVoterAddresss retrieves a list of voter addresses that were created after a specified block height.
// It delegates the operation to the `GetAllVoterAddresss` method of the `voteRepo` and handles any errors that may occur.
// If no voter addresses are found, it returns nil for the list and 0 for the height.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - chainId: The chain ID for which to retrieve the voter addresses.
//
// Returns:
//   - []string: A list of voter addresses that were created after the specified block height.
//   - error: An error object if the operation fails. Returns nil on success.
func (r *RpcService) GetAllVoterAddresss(ctx context.Context, chainId int64) ([]string, error) {
	data, err := r.voteRepo.GetAllVoterAddresss(ctx, chainId)
	if err != nil {
		r.logger.Error("GetAllVoterAddresss error", zap.Error(err))
		return nil, err
	}

	if len(data) == 0 {
		r.logger.Warn("GetAllVoterAddresss empty")
		return nil, nil
	}

	var addresss []string
	for _, d := range data {
		addresss = append(addresss, d.Address)
	}

	return addresss, nil
}

func (r *RpcService) GetVoterInfoByAddress(ctx context.Context, address string) (*model.VoterInfoTbl, error) {
	voterInfo, err := r.voteRepo.GetVoterInfoByAddress(ctx, address)
	if err != nil {
		r.logger.Error("GetVoterInfoByAddress error", zap.Error(err))
		return nil, err
	}

	if len(voterInfo.MinerIds) != 0 {
		minerIds := make([]uint64, 0, len(voterInfo.MinerIds))
		for _, minerId := range voterInfo.MinerIds {
			cutId := strings.ReplaceAll(minerId, config.Client.Network.MinerIdPrefix, "")
			id, err := strconv.ParseInt(cutId, 10, 64)
			if err != nil {
				r.logger.Error("ParseInt error", zap.Error(err), zap.Any("minerId", minerId))
				continue
			}
			minerIds = append(minerIds, uint64(id))
		}

		validMinerIds, err := r.lotusRepo.GetValidMinerIds(ctx, voterInfo.OwnerId, minerIds)
		if err != nil {
			r.logger.Error("verfy minerIds error", zap.Error(err))
		} else {
			voterInfo.MinerIds = validMinerIds

			if err := r.voteRepo.UpdateVoterByMinerInfo(ctx, voterInfo); err != nil {
				r.logger.Warn("UpdateVoterByMinerInfo error for GetVoterInfoByAddress", zap.Any("voterInfo", voterInfo), zap.Error(err))
			}
		}
	}

	clearGist := func() {
		voterInfo.GistId = ""
		voterInfo.GithubName = ""
		voterInfo.GistInfo = ""
	}

	if voterInfo.GistId != "" {
		gist, err := utils.FetchGistInfoByGistId(voterInfo.GistId)
		if err != nil {
			r.logger.Error("GetGistInfoByGistId error", zap.Error(err))

			clearGist()
			if errors.Is(err, constant.ErrGistNotFound) {
				if err := r.voteRepo.UpdateVoterByGistInfo(ctx, voterInfo); err != nil {
					r.logger.Warn("UpdateVoterByGistInfo error for GetVoterInfoByAddress", zap.Any("voterInfo", voterInfo), zap.Error(err))
				}
			}

			return voterInfo, nil
		}

		isValid := utils.VerifyAuthorizeAllow(voterInfo.GithubName, gist, func(gistAddr string) bool {
			if gistAddr == voterInfo.Address {
				return true
			}

			gistAddrActorId, err := r.lotusRepo.GetActorIdByAddress(ctx, gistAddr)
			if err != nil {
				zap.L().Error("GetActorIdByAddress failed by VerifyGist", zap.String("address", gistAddr), zap.Error(err))
				return false
			}

			return voterInfo.OwnerId == gistAddrActorId
		})
		if !isValid {
			r.logger.Warn("VerifyAuthorizeAllow error", zap.Any("gist", gist), zap.Any("voterInfo", voterInfo))
			clearGist()
			if err := r.voteRepo.UpdateVoterByGistInfo(ctx, voterInfo); err != nil {
				r.logger.Warn("UpdateVoterByGistInfo error for GetVoterInfoByAddress", zap.Any("voterInfo", voterInfo), zap.Error(err))
			}
		}
	}

	return voterInfo, nil
}

func (r *RpcService) GetGithubRepoName(ctx context.Context, orgType int) ([]model.GithubRepos, error) {
	res, err := r.syncRepo.GetGithubRepoName(ctx, orgType)
	if err != nil {
		zap.L().Error("GetGithubRepoName failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}
