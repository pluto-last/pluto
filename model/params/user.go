package params

import (
	"pluto/model/table"
)

// User register structure
type Register struct {
	Mobile   string `json:"mobile" form:"mobile"`
	Password string `json:"password" form:"password"`
	Sms      string `json:"sms" form:"sms"`
}

// User login structure
type Login struct {
	Mobile   string `json:"mobile" form:"mobile"`     // 用户名
	Password string `json:"password" form:"password"` // 密码
	Source   string `json:"source" form:"source"`     // 来源
}

type LoginResponse struct {
	User        table.User         `json:"user"`
	Token       string             `json:"token"`
	Permissions []table.Permission `json:"permissions"`
}

type GetUserList struct {
	PageInfo
	Mobile   string `json:"mobile" form:"mobile"`     // 手机号码
	UserName string `json:"userName" form:"userName"` // 用户名
}

type SetUserInfo struct {
	ID        string `json:"id" form:"id"`             // 用户id
	UserName  string `json:"userName" form:"userName"` // 用户名
	HeaderImg string `json:"headerImg" form:"headerImg"`
	Note      string `json:"note" form:"note"`
	Status    string `json:"status" form:"status"`
}
