package request

type GetHeightReq struct {
	NetId int64 `form:"chainId" binding:"required"`
}
