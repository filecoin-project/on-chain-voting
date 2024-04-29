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

package utils

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"powervoting-server/config"
	"testing"
)

func TestGetIpfsAndDecrypt(t *testing.T) {
	config.InitConfig("../")
	ipfs, err := GetIpfs("bafkreic2hs32eeortzls7utl5bu3yjxieb64k2q3afqn2l7enofeamvqjq")
	if err != nil {
		zap.L().Error("get ipfs error: ", zap.Error(err))
		return
	}
	decrypt, err := Decrypt(ipfs)
	if err != nil {
		zap.L().Error("decrypt error: ", zap.Error(err))
		return
	}
	fmt.Println("decrypt string: ", string(decrypt))
	var mapData [][]string
	err = json.Unmarshal(decrypt, &mapData)
	if err != nil {
		zap.L().Error("unmarshal error：", zap.Error(err))
		return
	}
	fmt.Println("Map data: ", mapData)
}

func TestGetOptions(t *testing.T) {
	config.InitConfig("../")
	options, err := GetOptions("bafkreihmdncwpk2kgos7ddzhu2tznirmjhevukhmvonvtspkeomsgxvsty")
	if err != nil {
		zap.L().Error("get option error: ", zap.Error(err))
		return
	}
	fmt.Println(options)
	fmt.Println(len(options))
}
