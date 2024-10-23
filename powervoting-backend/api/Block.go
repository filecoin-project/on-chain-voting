package api

import (
	"crypto/rand"
	"math/big"
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

	num, err := rand.Int(rand.Reader, big.NewInt(constant.RandomDay+1))
	if err != nil {
		zap.L().Error("get random day error : ", zap.Error(err))
		response.SystemError(c)
		return
	}

	dayStr := carbon.Now().SubDays(int(num.Int64())).EndOfDay().ToShortDateString()
	height, err := client.GetDataHeight(req.NetId, dayStr)
	if err != nil {
		zap.L().Error("get dataheight error : ", zap.Error(err))
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
