package params

type GetSipList struct {
	PageInfo
	Mobile string `json:"mobile" form:"mobile"` // 线路显示号码
	Name   string `json:"name" form:"name"`     // 线路名称
}

type SetSipInfo struct {
	ID          string `json:"id" form:"id"`                   // sipID
	Name        string `json:"name" form:"name"`               // 线路名称
	Mobile      string `json:"mobile" form:"mobile"`           // 线路显示号码
	IntervalSec int    `json:"intervalSec" form:"intervalSec"` // 呼叫间隔
	SipIP       string `json:"sipIP" form:"sipIP"`             // 线路IP
	SipPort     string `json:"sipPort" form:"sipPort"`         // 线路端口
	Note        string `json:"note" form:"note"`               //备注
}

type GetUserSipList struct {
	PageInfo
	UserID string `json:"userID" form:"userID"` // 用户ID
	Name   string `json:"name" form:"name"`     // 线路名称
}

type SetUserSip struct {
	ID         string `json:"id" form:"id"` // id
	UserID     string `json:"userID" form:"userID"`
	SipID      string `json:"sipID" form:"sipID"`
	ExpireAt   string `json:"expireAt" form:"expireAt"`     // 过期时间
	Concurrent int    `json:"concurrent" form:"concurrent"` // 并发数
	Price      int    `json:"price" form:"price"`           // 单价 每分钟
}
