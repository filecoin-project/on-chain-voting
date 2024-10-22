package api

import (
	"math/rand"
	"powervoting-server/client"
	"powervoting-server/constant"
	"powervoting-server/request"
	"powervoting-server/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"go.uber.org/zap"
)

func GetHeight(c *gin.Context) {
	var req request.GetHeightReq
	if err := c.ShouldBindQuery(&req); err != nil {
		zap.L().Error("get height : ", zap.Error(err))
		response.SystemError(c)
		return
	}

	num := rand.Intn(int(constant.Period))
	num += 1

	dayStr := carbon.Now().SubDays(num).EndOfDay().ToShortDateString()
	height, err := client.GetDataHeight(req.NetId, dayStr)
	if err != nil {
		zap.L().Error("get dataheight error : ", zap.Error(err))
		response.SystemError(c)
		return
	}

	if height == 0 {
		zap.L().Error("fail to get dataheight")
		response.SystemError(c)
		return
	}

	res := response.DataHeightRep{
		Day:    dayStr,
		Height: height,
		NetId:  req.NetId,
	}

	response.SuccessWithData(res, c)
}
