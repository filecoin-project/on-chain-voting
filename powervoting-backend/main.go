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
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

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

	proposalService := service.NewProposalService(proposalRepoImpl)
	voteService := service.NewVoteService(voteRepoImpl)
	syncService := service.NewSyncService(syncRepoImpl, voteRepoImpl, proposalRepoImpl)
	// initialization scheduled task
	go task.TaskScheduler(syncService)

	// default gin web
	r := gin.Default()
	r.Use(Cors())
	router.InitRouters(r, proposalService, voteService)
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
