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

package api

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "power-snapshot/api/proto"
	"power-snapshot/config"
	"power-snapshot/constant"
	models "power-snapshot/internal/model"
)

var (
	backendClient pb.BackendClient
	clientOnce    sync.Once
	logger        = zap.L().With(zap.String("gRPC", "backend"))
)

// getClient returns a singleton gRPC client instance.
func getClient() pb.BackendClient {
	clientOnce.Do(func() {
		conn, err := grpc.NewClient(
			config.Client.Server.RpcUri,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			zap.L().Error("failed to connect to gRPC server", zap.Error(err))
		}
		backendClient = pb.NewBackendClient(conn)
	})
	return backendClient
}

func GetAllVoterAddresss(chainId int64) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constant.TimeoutSecond)
	defer cancel()
	grpcReq := &pb.GetAllVoterAddressRequest{
		ChainId: chainId,
	}

	grpcResp, err := getClient().GetAllVoterAddresss(ctx, grpcReq)
	if err != nil {
		logger.Error("GetAllVoterAddresss failed", zap.Error(err))
		return nil, err
	}

	return grpcResp.Address, nil
}

func GetVoterInfo(address string) (models.VoterInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constant.TimeoutSecond)
	defer cancel()
	grpcReq := &pb.GetVoterInfoRequest{
		Address: address,
	}

	grpcResp, err := getClient().GetVoterInfo(ctx, grpcReq)
	if err != nil {
		logger.Error("GetVoterInfo failed", zap.Error(err))
		return models.VoterInfo{}, err
	}

	return models.VoterInfo{
		EthAddress:    common.HexToAddress(address),
		ActorIds:      []uint64{grpcResp.ActorId},
		MinerIds:      grpcResp.MinerIds,
		GithubAccount: grpcResp.GithubAccount,
	}, nil
}
