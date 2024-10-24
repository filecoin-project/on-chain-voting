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
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"powervoting-server/client"
	"powervoting-server/constant"
	"powervoting-server/db"
	"powervoting-server/model"
	"powervoting-server/request"
	"powervoting-server/response"
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
		zap.L().Info("get upload file error: ", zap.Error(err))
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
	cid, err := client.W3.Upload(absolutePath)
	if err != nil {
		os.Remove(absolutePath)
		zap.L().Info("get upload file error: ", zap.Error(err))
		response.SystemError(c)
	}

	os.Remove(absolutePath)
	response.SuccessWithData(cid, c)
}

func AddDraft(c *gin.Context) {
	var draft model.ProposalDraft
	if err := c.ShouldBindJSON(&draft); err != nil {
		zap.L().Error("add draft error: ", zap.Error(err))
		response.SystemError(c)
		return
	}
	var count int64
	db.Engine.Model(model.ProposalDraft{}).Where("chain_id", draft.ChainId).Where("address", draft.Address).Count(&count)
	if count == 0 {
		result := db.Engine.Model(model.ProposalDraft{}).Create(&draft)
		if result.Error != nil {
			zap.L().Error("insert draft error: ", zap.Error(result.Error))
			response.SystemError(c)
			return
		}
	} else {
		result := db.Engine.Model(model.ProposalDraft{}).Where("chain_id", draft.ChainId).Where("address", draft.Address).Select("timezone", "time", "name", "descriptions", "option", "start_time", "exp_time").Updates(&draft)
		if result.Error != nil {
			zap.L().Error("update draft error: ", zap.Error(result.Error))
			response.SystemError(c)
			return
		}
	}

	response.SuccessWithData(true, c)
}

func AddProposal(c *gin.Context) {
	var proposalReq request.Proposal
	if err := c.ShouldBindJSON(&proposalReq); err != nil {
		zap.L().Error("add draft error: ", zap.Error(err))
		response.SystemError(c)
		return
	}

	var proposal model.Proposal
	db.Engine.Model(model.Proposal{}).Where("cid", proposalReq.Cid).Take(&proposal)
	if proposal.Id != 0 {
		zap.L().Error("proposal already exists")
		response.Error(errors.New("proposal already exists"), c)
		return
	}

	proposal = model.Proposal{
		Name:         proposalReq.Name,
		Descriptions: proposalReq.Descriptions,
		Network:      proposalReq.Network,
		Timezone:     proposalReq.Timezone,
		GithubName:   proposalReq.GithubName,
		GithubAvatar: proposalReq.GithubAvatar,
		GMTOffset:    proposalReq.GMTOffset,
		Cid:          proposalReq.Cid,
		Creator:      proposalReq.Creator,
		StartTime:    proposalReq.StartTime,
		ExpTime:      proposalReq.ExpTime,
		CurrentTime:  proposalReq.CurrentTime,
		VoteCountDay: proposalReq.VoteCountDay,
		Height:       proposalReq.Height,
	}

	result := db.Engine.Model(model.Proposal{}).Create(&proposal)
	if result.Error != nil {
		zap.L().Error("insert proposal error: ", zap.Error(result.Error))
		response.SystemError(c)
		return
	}

	db.Engine.Model(model.SnapshotByDay{}).Where("day", proposal.VoteCountDay).Delete(&model.SnapshotByDay{})

	snapshotTask := model.SnapshotByDay{
		Day:    proposal.VoteCountDay,
		NetId:  proposal.Network,
		Height: proposal.Height,
	}
	res := db.Engine.Model(model.SnapshotByDay{}).Create(&snapshotTask)
	if res.Error != nil {
		zap.L().Error("add snapshotTask error: ", zap.Error(result.Error))
		response.SystemError(c)
		return
	}

	response.SuccessWithData(true, c)
}

func GetDraft(c *gin.Context) {
	var req request.GetDraft
	if err := c.ShouldBindQuery(&req); err != nil {
		zap.L().Error("add draft error: ", zap.Error(err))
		response.SystemError(c)
		return
	}

	var result []model.ProposalDraft

	tx := db.Engine.Model(model.ProposalDraft{}).Where("chain_id", req.ChainId).Where("Address", req.Address).Find(&result)

	if tx.Error != nil {
		zap.L().Error("Get draft result error: ", zap.Error(tx.Error))
		response.SystemError(c)
		return
	}
	response.SuccessWithData(result, c)
}

func ProposalList(c *gin.Context) {
	var req request.ProposalList

	if err := c.ShouldBindQuery(&req); err != nil {
		zap.L().Error("Param error: ", zap.Error(err))
		return
	}

	queryCount := db.Engine.Model(model.Proposal{})
	queryList := db.Engine.Model(model.Proposal{}).Order("created_at desc").Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)

	if req.SearchKey != "" {
		queryCount.Where("name like ?", "%"+req.SearchKey+"%")
		queryList.Where("name like ?", "%"+req.SearchKey+"%")
	}

	switch req.Status {
	case constant.ProposalStatusPending:
		queryCount.Where("status = ?", constant.ProposalStatusPending).Where("start_time > ?", time.Now().Unix())
		queryList.Where("status = ?", constant.ProposalStatusPending).Where("start_time > ?", time.Now().Unix())
	case constant.ProposalStatusInProgress:
		queryCount.Where("status = ?", constant.ProposalStatusPending).Where("start_time < ?", time.Now().Unix()).Where("exp_time > ?", time.Now().Unix()).Where("exp_time > ?", time.Now().Unix())
		queryList.Where("status = ?", constant.ProposalStatusPending).Where("start_time < ?", time.Now().Unix()).Where("exp_time > ?", time.Now().Unix()).Where("exp_time > ?", time.Now().Unix())
	case constant.ProposalStatusCounting:
		queryCount.Where("status = ?", constant.ProposalStatusPending).Where("exp_time < ?", time.Now().Unix())
		queryList.Where("status = ?", constant.ProposalStatusPending).Where("exp_time < ?", time.Now().Unix())
	case constant.ProposalStatusCompleted:
		queryCount.Where("status = ?", constant.ProposalStatusCompleted)
		queryList.Where("status = ?", constant.ProposalStatusCompleted)
	}

	var count int64
	tx := queryCount.Count(&count)
	if tx.Error != nil {
		zap.L().Error("Get proposal result error: ", zap.Error(tx.Error))
		response.SystemError(c)
		return
	}

	var proposals []model.Proposal
	result := make(map[string]any)
	result["total"] = count

	if count == 0 {
		result["list"] = []response.Proposal{}
		response.SuccessWithData(result, c)
		return
	}

	tx = queryList.Find(&proposals)
	if tx.Error != nil {
		zap.L().Error("Get proposal result error: ", zap.Error(tx.Error))
		response.SystemError(c)
		return
	}

	proposalIds := []int64{}
	for _, v := range proposals {
		if v.Status == constant.ProposalStatusCompleted ||
			(v.Status == constant.ProposalStatusPending && v.ExpTime < time.Now().Unix() && v.VoteCount != 0) {
			proposalIds = append(proposalIds, v.ProposalId)
		}
	}

	var voteResult []model.VoteResult
	db.Engine.Model(model.VoteResult{}).Where("proposal_id in ?", proposalIds).Find(&voteResult)

	var voteMap = make(map[int64][]model.VoteResult)
	for _, v := range voteResult {
		voteMap[v.ProposalId] = append(voteMap[v.ProposalId], v)
	}

	var list []response.Proposal
	for _, v := range proposals {
		temp := response.Proposal{
			ProposalId:   v.ProposalId,
			Cid:          v.Cid,
			Creator:      v.Creator,
			StartTime:    v.StartTime,
			ExpTime:      v.ExpTime,
			Network:      v.Network,
			Name:         v.Name,
			Timezone:     v.Timezone,
			Descriptions: v.Descriptions,
			GithubName:   v.GithubName,
			GithubAvatar: v.GithubAvatar,
			GMTOffset:    v.GMTOffset,
			CurrentTime:  v.CurrentTime,
			CreatedAt:    v.CreatedAt.Unix(),
			UpdatedAt:    v.UpdatedAt.Unix(),
			VoteResult:   []model.VoteResult{},
			Time:         []string{},
			Option:       []string{},
			ShowTime:     []string{},
			Status:       v.Status,
			VoteCount:    v.VoteCount,
		}

		startTimeFormat := time.Unix(v.StartTime, 0).In(time.UTC).Format(time.RFC3339)
		expTimeFormant := time.Unix(v.ExpTime, 0).In(time.UTC).Format(time.RFC3339)
		temp.Time = []string{
			startTimeFormat,
			expTimeFormant,
		}
		temp.ShowTime = []string{
			startTimeFormat,
			expTimeFormant,
		}

		if temp.Status == constant.ProposalStatusPending {
			if temp.StartTime < time.Now().Unix() && temp.ExpTime > time.Now().Unix() {
				temp.Status = constant.ProposalStatusInProgress
			} else if temp.ExpTime < time.Now().Unix() {
				temp.Status = constant.ProposalStatusCounting
			}
		}

		temp.Option = []string{
			"Approve",
			"Reject",
		}

		var approve float64
		var reject float64
		if voteMap[v.ProposalId] != nil {
			temp.VoteResult = voteMap[v.ProposalId]
			for _, v := range temp.VoteResult {
				if v.OptionId == constant.VoteApprove {
					approve = v.Votes
				}
				if v.OptionId == constant.VoteReject {
					reject = v.Votes
				}
			}
		}

		if temp.Status == constant.ProposalStatusCompleted {
			temp.Status = constant.ProposalStatusPassed
			if approve == 0 || approve <= reject {
				temp.Status = constant.ProposalStatusRejected
			}
		}

		list = append(list, temp)
	}

	result["list"] = list
	response.SuccessWithData(result, c)
}
