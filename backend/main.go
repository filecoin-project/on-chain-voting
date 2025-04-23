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
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"powervoting-server/api/rpc"
	pb "powervoting-server/api/rpc/proto"
	"powervoting-server/config"
	"powervoting-server/data"
	"powervoting-server/repo"
	"powervoting-server/router"
	"powervoting-server/service"
	"powervoting-server/task"
)

func main() {
	// initialization configuration
	config.InitConfig("./")
	// initialization logger
	config.InitLogger()
	// initialization mysql
	mydb := data.NewMysql()

	proposalRepoImpl := repo.NewProposalRepo(mydb)
	voteRepoImpl := repo.NewVoteRepo(mydb)
	syncRepoImpl := repo.NewSyncRepo(mydb)
	fipRepoImpl := repo.NewFipRepo(mydb)
	lotusRepoImpl := repo.NewLotusRPCRepo()
	proposalService := service.NewProposalService(proposalRepoImpl)
	voteService := service.NewVoteService(voteRepoImpl, lotusRepoImpl)
	syncService := service.NewSyncService(
		syncRepoImpl,
		voteRepoImpl,
		proposalRepoImpl,
		fipRepoImpl,
		lotusRepoImpl,
	)
	fipService := service.NewFipService(fipRepoImpl)
	// initialization scheduled task
	go task.TaskScheduler(syncService)

	// initialization grpc server
	go RpcServer(rpc.NewBackendRpc(service.NewRpcService(voteRepoImpl, lotusRepoImpl)))
	// default gin web
	r := gin.Default()
	r.Use(Cors())
	router.InitRouters(r, proposalService, voteService, fipService)
	err := r.Run(config.Client.Server.Port)
	if err != nil {
		zap.L().Error("start web server failed: ", zap.Error(err))
	}
}

// Cors is a middleware function that sets CORS headers.
func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method

		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
	}
}

func RpcServer(rpc *rpc.BackendRpc) {
	_ = grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(
			func(p any) (err error) {
				zap.S().Error("recovery form panic", zap.Error(err))
				return status.Errorf(codes.Internal, "internal error")
			},
		)),
	)

	server := grpc.NewServer()
	pb.RegisterBackendServer(server, rpc)

	lis, err := net.Listen("tcp", config.Client.Server.RpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("rpc server start on port: %v\n", config.Client.Server.RpcPort)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("rpc server shutdown")
}
