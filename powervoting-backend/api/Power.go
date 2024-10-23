package api

import (
	"powervoting-server/request"
	"powervoting-server/response"

	"powervoting-server/client"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetPower(c *gin.Context) {
	var req request.GetPower
	if err := c.ShouldBindQuery(&req); err != nil {
		zap.L().Error("get req error: ", zap.Error(err))
		response.SystemError(c)
		return
	}

	power, err := client.GetAddressPowerByDay(req.NetId, req.Address, req.Day)
	if err != nil {
		zap.L().Error("get req error: ", zap.Error(err))
		response.SystemError(c)
		return
	}

	res := response.Power{
		DeveloperPower:   power.DeveloperPower.String(),
		SpPower:          power.SpPower.String(),
		ClientPower:      power.ClientPower.String(),
		TokenHolderPower: power.TokenHolderPower.String(),
		BlockHeight:      power.BlockHeight.String(),
	}

	response.SuccessWithData(res, c)
}
