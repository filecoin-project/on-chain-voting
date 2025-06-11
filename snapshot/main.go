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

package main

import (
	"context"
	"log"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"power-snapshot/api"
	pb "power-snapshot/api/proto"
	"power-snapshot/config"
	"power-snapshot/constant"
	"power-snapshot/handler"
	"power-snapshot/internal/data"
	"power-snapshot/internal/repo"
	"power-snapshot/internal/service"
	"power-snapshot/internal/task"
)

func main() {
	// init logger
	initEvn()

	// init third-part util
	manager, err := data.NewGoEthClientManager()
	if err != nil {
		panic(err)
	}

	if len(manager.GetClient().QueryClient) == 0 {
		panic("no geth client")
	}

	contractRepo := repo.NewContractRepoImpl(manager.GetClient().QueryClient[0])
	expDates, err := contractRepo.GetExpirationData()
	if err != nil {
		zap.L().Error("get expiration data failed", zap.Error(err))
		expDates = constant.DataExpiredDuration
	}

	// init datasource
	redisClient, err := data.NewRedisClient(expDates)
	if err != nil {
		panic(err)
	}
	defer func(redisClient *data.RedisClient) {
		err := redisClient.Close()
		if err != nil {
			zap.S().Error("close redis client failed", zap.Error(err))
			return
		}
	}(redisClient)

	// init jetstream ant nats
	jetstreamClient, err := data.NewJetstreamClient()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := jetstreamClient.Drain()
		if err != nil {
			zap.S().Error("drain jetstream client failed", zap.Error(err))
			return
		}
	}()

	// init repo
	syncRepo, err := repo.NewSyncRepoImpl(manager.GetChainId(), redisClient, jetstreamClient)
	if err != nil {
		panic(err)
	}

	queryRepo, err := repo.NewQueryRepoImpl(redisClient, manager)
	if err != nil {
		panic(err)
	}

	baseRepo := repo.NewBaseRepoImpl(manager, redisClient)

	mysqlRepo := repo.NewMysqlRepoImpl(data.NewMysql())

	lotusRepo := repo.NewLotusRPCRepo(redisClient)

	// init service
	syncSrv := service.NewSyncService(baseRepo, syncRepo, mysqlRepo, lotusRepo, contractRepo, api.NewBackendGRPCClient())

	go func() {
		defer func() {
			if r := recover(); r != nil {
				zap.S().Error("recover from panic", zap.Any("err", r))
			}
		}()
		err := syncSrv.StartSyncWorker(context.Background(), config.Client.Network.ChainId)
		if err != nil {
			zap.S().Error("start sync worker failed", zap.Error(err))
			panic(err)
		}
	}()

	// init job
	go task.TaskScheduler(syncSrv)

	// init query service
	querySrv := service.NewQueryService(baseRepo, queryRepo, syncSrv, lotusRepo)

	// init handler
	snapshotHandler := handler.NewSnapshot(querySrv, syncSrv)

	// init grpc
	opt := grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
			zap.S().Error("recovery from panic", zap.Any("p", p))
			return status.Errorf(codes.Internal, "internal error")
		})),
	)
	server := grpc.NewServer(opt)
	pb.RegisterSnapshotServer(server, snapshotHandler)

	// start
	lis, err := net.Listen("tcp", config.Client.Server.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Starting gRPC server...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initEvn() {
	config.InitLogger()
	// load config
	err := config.InitConfig("./")
	if err != nil {
		panic(err)
	}
}
