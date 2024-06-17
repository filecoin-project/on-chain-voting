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
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"powervoting-server/db"
	"powervoting-server/model"
	"powervoting-server/response"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// VoteResult handles the HTTP request to retrieve the result of a vote for a specific proposal on a given network.
func VoteResult(c *gin.Context) {
	proposalId := c.Query("proposalId")
	network := c.Query("network")
	if proposalId == "" || network == "" {
		zap.L().Error("Param error, proposalId: ", zap.String("proposalId", proposalId))
		response.ParamError(c)
		return
	}
	var result []model.VoteResult
	tx := db.Engine.Model(model.VoteResult{}).Where("proposal_id", proposalId).Where("network", network).Find(&result)
	if tx.Error != nil {
		zap.L().Error("Get vote result error: ", zap.Error(tx.Error))
		response.SystemError(c)
		return
	}
	response.SuccessWithData(result, c)
}

// VoteHistory function handles an HTTP request to retrieve the voting history of a specific proposal on a given network.
func VoteHistory(c *gin.Context) {
	proposalId := c.Query("proposalId")
	network := c.Query("network")
	if proposalId == "" || network == "" {
		zap.L().Error("Param error, proposalId: ", zap.String("proposalId", proposalId))
		response.ParamError(c)
		return
	}
	var history model.VoteCompleteHistory
	tx := db.Engine.Preload("VotePowers").Where("proposal_id", proposalId).Where("network", network).Find(&history)
	if tx.Error != nil {
		zap.L().Error("Get vote result error: ", zap.Error(tx.Error))
		response.SystemError(c)
		return
	}
	response.SuccessWithData(history, c)
}

func W3Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		zap.L().Error("upload file error: ", zap.Error(err))
		response.SystemError(c)
		return
	}

	//rand file name
	randSource := rand.NewSource(time.Now().UnixNano())
	r := rand.New(randSource)
	randomNumber := r.Intn(1000000)
	timeStamp := time.Now().Unix()

	filePath := fmt.Sprintf("./uploads/%d_%d", timeStamp, randomNumber)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		zap.L().Error("save file error: ", zap.Error(err))
		response.SystemError(c)
		return
	}

	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		zap.L().Error("get file path error: ", zap.Error(err))
		response.SystemError(c)
		return
	}
	zap.L().Info("upload with w3")
	cmd := exec.Command("w3", "upload", absolutePath, "--json", "--no-wrap")

	//execut w3 upload xxxx
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	if err != nil {
		os.Remove(absolutePath)
		zap.L().Error("upload file error: ", zap.Error(err))
		response.SystemError(c)
		return
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(outBuf.Bytes()), &jsonData)
	if err != nil {
		os.Remove(absolutePath)
		zap.L().Error("upload file error: ", zap.Error(err))
		response.SystemError(c)
		return
	}
	os.Remove(absolutePath)
	response.SuccessWithData(jsonData, c)
}

func AddDraft(c *gin.Context) {
	var draft model.ProposalDraft
	if err := c.ShouldBindJSON(&draft); err != nil {
		zap.L().Error("add draft error: ", zap.Error(err))
		response.SystemError(c)
		return
	}
	var count int64
	db.Engine.Model(model.ProposalDraft{}).Where("chain_id", draft.ChainId).Where("Address", draft.Address).Count(&count)
	if count == 0 {
		result := db.Engine.Model(model.ProposalDraft{}).Create(&draft)
		if result.Error != nil {
			zap.L().Error("insert draft error: ", zap.Error(result.Error))
			response.SystemError(c)
			return
		}
	} else {
		printStruct(draft)
		result := db.Engine.Model(model.ProposalDraft{}).Where("chain_id", draft.ChainId).Where("address", draft.Address).Select("Timezone", "Time", "Name", "Descriptions", "Option").Updates(&draft)
		if result.Error != nil {
			zap.L().Error("update draft error: ", zap.Error(result.Error))
			response.SystemError(c)
			return
		}
	}

	response.SuccessWithData(true, c)
}

func GetDraft(c *gin.Context) {
	chainId := c.Query("chainId")
	address := c.Query("address")
	var result []model.ProposalDraft

	tx := db.Engine.Model(model.ProposalDraft{}).Where("chain_id", chainId).Where("Address", address).Find(&result)

	if tx.Error != nil {
		zap.L().Error("Get draft result error: ", zap.Error(tx.Error))
		response.SystemError(c)
		return
	}
	response.SuccessWithData(result, c)
}
func printStruct(s interface{}) {
	// 获取结构体类型信息
	val := reflect.ValueOf(s)
	typ := val.Type()

	// 遍历结构体的字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		fmt.Printf("%s: %v\n", fieldName, field.Interface())
	}
}
func deleteDraft(c *gin.Context) {
	// chainId := c.Query("chainId")
	// address := c.Query("address")
	// var result []model.ProposalDraft

	// tx := db.Engine.Model(model.ProposalDraft{}).Where("chainId", chainId).Where("address", address).Delete()

	// if tx.Error != nil {
	// 	zap.L().Error("Get draft result error: ", zap.Error(tx.Error))
	// 	response.SystemError(c)
	// 	return
	// }
	// response.SuccessWithData(result, c)
}
