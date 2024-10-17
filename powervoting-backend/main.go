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
	"powervoting-server/client"
	"powervoting-server/config"
	"powervoting-server/db"
	"powervoting-server/routers"
	"powervoting-server/scheduler"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// initialization configuration
	config.InitConfig("./")
	// initialization logger
	config.InitLogger()
	// initialization mysql
	db.InitMysql()

	client.InitW3Client()
	// initialization scheduled task
	go scheduler.TaskScheduler()

	// default gin web
	r := gin.Default()
	r.Use(Cors())
	routers.InitRouters(r)
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
