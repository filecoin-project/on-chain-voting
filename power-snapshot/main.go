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
	pb "power-snapshot/api/proto"
	"power-snapshot/config"
	"power-snapshot/handler"
	"power-snapshot/internal/data"
	"power-snapshot/internal/repo"
	"power-snapshot/internal/scheduler"
	"power-snapshot/internal/service"
	"power-snapshot/internal/task"
	"power-snapshot/utils"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	ctx := context.Background()
	// init logger
	config.InitLogger()
	// load config
	err := config.InitConfig("./")
	if err != nil {
		panic(err)
	}
	// init third-part util
	manager, err := utils.NewGoEthClientManager(config.Client.Network)
	if err != nil {
		panic(err)
	}
	// init datasource
	redisClient, err := data.NewRedisClient()
	if err != nil {
		panic(err)
	}
	defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {
			zap.S().Error("close redis client failed", zap.Error(err))
			return
		}
	}(redisClient)
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
	syncRepo, err := repo.NewSyncRepoImpl(manager.ListClientNetWork(), redisClient, jetstreamClient)
	if err != nil {
		panic(err)
	}
	queryRepo, err := repo.NewQueryRepoImpl(redisClient, manager)
	if err != nil {
		panic(err)
	}
	baseRepo, err := repo.NewBaseRepoImpl(manager, redisClient)
	if err != nil {
		panic(err)
	}
	// init service
	syncSrv := service.NewSyncService(baseRepo, syncRepo)
	for _, network := range config.Client.Network {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					zap.S().Error("recover from panic", zap.Any("err", r))
				}
			}()
			err := syncSrv.StartSyncWorker(ctx, network.Id)
			if err != nil {
				zap.S().Error("start sync worker failed", zap.Error(err))
				panic(err)
			}
		}()
	}

	querySrv := service.NewQueryService(baseRepo, queryRepo, syncSrv, redisClient)
	// init job
	jobs := []scheduler.JobInfo{
		{
			// every day 1 am
			Spec: "0 0 1 * * ?",
			Job:  task.SyncPower(syncSrv),
		},
		{
			// every hour
			Spec: "0 0 0/1 * * ?",
			Job:  task.SyncDevWeightStepDay(syncSrv),
		},
		{
			// every day 1 am
			Spec: "0 0 1 * * ?",
			Job:  task.SyncDelegateEvent(syncSrv),
		},
	}
	scheduler.TaskScheduler(jobs)
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
