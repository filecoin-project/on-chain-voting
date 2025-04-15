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
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// export client
var Client Config

func initEnv(envFile string) {
	if err := godotenv.Load(envFile); err != nil {
		log.Default().Printf("Unable to load .env file: %v", zap.Error(err))
	}
}

// InitConfig initializes the configuration by reading from a YAML file located at the specified path.
func InitConfig(path string) {
	initEnv(path + "/.env")
	// configuration file name
	viper.SetConfigName("configuration-backend")

	viper.AddConfigPath(path)

	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		zap.L().Error("read config file error:", zap.Error(err))
		return
	}

	for _, key := range viper.AllKeys() {
		value := viper.GetString(key)
		if value == "" {
			continue
		}

		replacedValue := replaceEnvVariables(value)
		replacedValue = strings.ReplaceAll(replacedValue, "'", "")
		replacedValue = strings.ReplaceAll(replacedValue, "\"", "")

		if key == "xxl-job.serveraddrs" {
			viper.Set(key, strings.Split(replacedValue, ","))
		} else {
			viper.Set(key, replacedValue)
		}

	}

	err = viper.Unmarshal(&Client)
	if err != nil {
		zap.L().Error("unmarshal error:", zap.Error(err))
		return
	}

}

func replaceEnvVariables(value string) string {
	return os.Expand(value, func(key string) string {
		return os.Getenv(key)
	})
}
