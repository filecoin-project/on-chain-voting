package request

type ProposalList struct {
	PageReq
	Status    int    `form:"status" binding:"oneof=0 1 2 3 4"`
	SearchKey string `form:"searchKey"`
}

type PageReq struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type Proposal struct {
	Cid          string `json:"cid" binding:"required"`          // CID
	Creator      string `json:"address" binding:"required"`      // Creator address
	StartTime    int64  `json:"startTime" binding:"required"`    // Start time
	ExpTime      int64  `json:"expTime" binding:"required"`      // Expiry time
	Network      int64  `json:"chainId" binding:"required"`      // Network ID
	Name         string `json:"name" binding:"required"`         // Name
	Timezone     string `json:"timezone"`                        // Timezone
	Descriptions string `json:"descriptions" binding:"required"` // Descriptions
	GithubName   string `json:"githubName"`                      // Github name
	GithubAvatar string `json:"githubAvatar"`                    // Github avatar
	GMTOffset    string `json:"gmtOffset"`                       // GMT offset
	CurrentTime  int64  `json:"currentTime"`                     // Current time
}

type GetDraft struct {
	ChainId string `form:"chainId" binding:"required"`
	Address string `form:"address" binding:"required"`
}
