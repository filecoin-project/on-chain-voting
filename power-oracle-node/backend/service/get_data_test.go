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

package service

import (
	"backend/config"
	"backend/contract"
	"backend/utils"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetWalletBalance(t *testing.T) {
	var err error
	config.InitConfig("../")
	contract.GoEthClient, err = contract.GetClient(314159)
	if err != nil {
		t.Error(err)
	}
	contract.LotusRpcClient = utils.NewClient(contract.GoEthClient.Rpc)

	tokenHolder, err := GetWalletBalance("t1wa4gvyeek4oh5zg375oo6lwhcmdwgxws5rgslyy", contract.LotusRpcClient)
	assert.Nil(t, err)

	res, err := decimal.NewFromString("49999999118307566058")
	assert.Nil(t, err)

	assert.Equal(t, res.BigInt(), tokenHolder)

	// test if addr not exist
	_, err = GetWalletBalance("", contract.LotusRpcClient)
	assert.NotNil(t, err)

	fmt.Println(tokenHolder)
}

func TestGetUcanIpfs(t *testing.T) {
	var err error
	config.InitConfig("../")
	contract.GoEthClient, err = contract.GetClient(314159)
	assert.Nil(t, err)

	rsp, err := GetUcanFromIpfs("bafkreidnrzkmshm36e4uzj5mxl6gmwv2b6uhkkwvvpxo3anjcxepbskn3q")
	assert.Nil(t, err)

	assert.NotEmpty(t, rsp)
}
