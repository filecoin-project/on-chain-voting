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

package config

import (
	"backend/models"
	"context"
	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

var Client models.Config
var GithubClient *github.Client

// InitConfig initializes the configuration by reading from a YAML file located at the specified path.
func InitConfig(path string) { // configuration file name
	viper.SetConfigName("configuration")

	viper.AddConfigPath(path)

	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		zap.L().Error("read config file error: ", zap.Error(err))
		return
	}

	err = viper.Unmarshal(&Client)
	if err != nil {
		zap.L().Error("unmarshal error: ", zap.Error(err))
		return
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: Client.Github.GithubToken})
	tc := oauth2.NewClient(context.Background(), ts)

	GithubClient = github.NewClient(tc)
}
