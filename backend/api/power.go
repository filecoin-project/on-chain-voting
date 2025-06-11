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
	"go.uber.org/zap"

	snapshot "powervoting-server/api/rpc"
	"powervoting-server/constant"
	"powervoting-server/model/api"
)

// GetPower handles the request to get power information for a specific address on a specific day.
func GetAddressPower(c *constant.Context) {
	// Declare a variable of type request.GetPower to hold the request parameters.
	var req api.GetPowerReq
	if err := c.BindAndValidate(&req); err != nil {
		zap.L().Error("GetAddressPower bind parmas error: ", zap.Errors("errors", err.Errors()))
		ParamError(c.Context)
		return
	}

	ethAddr, err := req.AddressReq.ToEthAddr()
	if err != nil {
		zap.L().Error("GetAddressPower invalid address: ", zap.String("address", req.AddressReq.Address), zap.Error(err))
		Error(c.Context, err)
		return
	}

	// Call the client's GetAddressPowerByDay method to retrieve power information.
	power, err := snapshot.GetAddressPowerByDay(req.ChainId, ethAddr, req.PowerDay)
	if err != nil {
		zap.L().Error(
			"get snapshot power error ",
			zap.String("address", ethAddr),
			zap.String("power day", req.PowerDay),
			zap.Int64("chain id", req.ChainId),
			zap.Error(err),
		)
		SystemError(c.Context)
		return
	}

	res := api.PowerRep{
		DeveloperPower:   power.DeveloperPower.String(),
		SpPower:          power.SpPower.String(),
		ClientPower:      power.ClientPower.String(),
		TokenHolderPower: power.TokenHolderPower.String(),
	}

	SuccessWithData(c.Context, res)
}
