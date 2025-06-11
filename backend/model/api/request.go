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
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"

	"powervoting-server/config"
	"powervoting-server/utils"
)

// ProposalListReq represents a request for listing proposals with pagination, status filter, and search functionality.
type ProposalListReq struct {
	Status    int    `form:"status" validate:"oneof=0 1 2 3 4"` // Status filter (0: all, 1: pending, 2: in progress, 3: counting, 4: completed)
	SearchKey string `form:"searchKey"`                         // Keyword for fuzzy search in proposal titles
	AddressReq
	PageReq      // Embedded pagination request
	ChainIdParam // Embedded chain ID parameter
}

// PageReq represents a pagination request with page number and page size.
type PageReq struct {
	Page     int `form:"page"`     // Current page number
	PageSize int `form:"pageSize"` // Number of items per page
}

// AddressReq represents a request for retrieving a proposal draft by creator address.
type AddressReq struct {
	Address string `form:"address"` // Creator's address
}

// VerifyGistReq verify that Gist is valid
type VerifyGistReq struct {
	GistId string `form:"gistId" validate:"required"` // Gist ID
	AddressReq
}

// GetPowerReq represents a request for retrieving power information for a specific address and day.
type GetPowerReq struct {
	PowerDay string `form:"powerDay" validate:"required"` // Day for which power information is requested
	ChainId  int64  `form:"chainId" validate:"required"`  // Chain ID to filter power information
	AddressReq
}

// ChainIdParam represents a request parameter for chain ID.
type ChainIdParam struct {
	ChainId int64 `form:"chainId" validate:"required,gt=0"` // Chain ID
}

// ProposalReq represents a request for retrieving a specific proposal by its ID and chain ID.
type ProposalReq struct {
	ProposalId   int64 `form:"proposalId" validate:"required,gt=0"` // Proposal ID
	ChainIdParam       // Embedded chain ID parameter
}

// AddProposalDraftReq represents a request for creating or updating a proposal draft.
type AddProposalDraftReq struct {
	Creator   string `json:"creator" validate:"required"`                // Creator address
	StartTime int64  `json:"startTime" validate:"required"`              // Start time of the proposal
	EndTime   int64  `json:"endTime" validate:"required"`                // End time of the proposal
	Timezone  string `json:"timezone" validate:"required"`               // Timezone of the proposal
	Title     string `json:"title" validate:"required,min=1,max=254"`    // Title of the proposal
	Content   string `json:"content" validate:"required,min=1,max=2000"` // Description of the proposal
	ChainIdParam
	ProposalPercentage
}

type DelProposalDraftReq struct {
    AddressReq
	ChainIdParam
}

type ProposalPercentage struct {
	TokenHolderPercentage uint16 `json:"tokenHolderPercentage" validate:"number,is-integer"` // Voting power percentage for token holders
	SpPercentage          uint16 `json:"spPercentage" validate:"number,is-integer"`          // Voting power percentage for SPs
	ClientPercentage      uint16 `json:"clientPercentage" validate:"number,is-integer"`      // Voting power percentage for clients
	DeveloperPercentage   uint16 `json:"developerPercentage" validate:"number,is-integer"`   // Voting power percentage for developers
}

type FipProposalListReq struct {
	PageReq
	ProposalType int   `form:"proposalType" validate:"oneof=0 1"` // FIP type filter (0: revoke, 1: approve)
	ChainId      int64 `form:"chainId" validate:"required"`
}

type FipEditorListReq struct {
	ChainId int64 `form:"chainId" validate:"required"`
}

// Offset calculates the offset for pagination based on the current page and page size.
// It ensures the page and page size are within valid ranges.
//
// Returns:
//   - int: The calculated offset for pagination.
func (p *PageReq) Offset() int {
	// Ensure page is at least 1
	if p.Page <= 0 {
		p.Page = 1
	}

	// Ensure page size is at least 10 and at most 50
	if p.PageSize <= 0 {
		p.PageSize = 10
	} else if p.PageSize > 50 {
		p.PageSize = 50
	}

	// Calculate the offset
	return (p.Page - 1) * p.PageSize
}

func (a *AddressReq) ToEthAddr() (string, error) {
	if a == nil {
		return "", nil
	}

	if a.Address == "" {
		return "", nil
	}
	
	lotusClient := jsonrpc.NewClientWithOpts(config.Client.Network.Rpc, &jsonrpc.RPCClientOpts{})

	if strings.HasPrefix(a.Address, "0x") {
		return utils.EthStandardAddressToHex(a.Address), nil
	}

	resp, err := lotusClient.Call(context.Background(), "Filecoin.FilecoinAddressToEthAddress", a.Address)
	if err != nil {
		zap.L().Error("FilcoinAddressToEthAddress: lotus rpc error", zap.String("address", a.Address), zap.Error(err))
		return "", errors.New("lotus rpc error")
	}

	if resp.Error != nil {
		zap.L().Error("FilcoinAddressToEthAddress error", zap.String("address", a.Address), zap.Error(resp.Error))
		return "", fmt.Errorf("get eth address error: %s", resp.Error.Message)
	}

	return utils.EthStandardAddressToHex(resp.Result.(string)), nil
}
