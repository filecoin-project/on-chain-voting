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

	"go.uber.org/zap"

	"powervoting-server/model"
)

type RpcService struct {
	voteRepo VoteRepo
	logger   *zap.Logger
}

func NewRpcService(voteRepo VoteRepo) *RpcService {
	return &RpcService{
		voteRepo: voteRepo,
		logger:   zap.L().With(zap.String("service", "RPC")),
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
	return r.voteRepo.GetVoterInfoByAddress(ctx, address)
}
