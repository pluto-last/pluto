package params

type CreateRecharge struct {
	UserID string `json:"userID" form:"userID"` // 充值的用户 id
	Type   string `json:"type" form:"type"`     // 充值类型
	Amount int    `json:"amount" form:"amount"` // 充值金额
}

type GetRechargeList struct {
	PageInfo
	Mobile string `json:"mobile" form:"mobile"` // 用户电话
	Type   string `json:"type" form:"type"`     // 充值类型
}
