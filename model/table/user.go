package table

import (
	"pluto/global"
)

type User struct {
	global.UUID
	UserName   string `json:"userName"` // 用户姓名
	Mobile     string `json:"mobile" `  // 手机号码
	Password   string `json:"password"  `
	RegisterIP string `gorm:"index"`      // 注册IP
	HeaderImg  string `json:"headerImg" ` // 头像
	Note       string `json:"note" `      // 备注
	Status     string `json:"status" `    // 状态
	WalletInfo Wallet `json:"wallet" gorm:"ForeignKey:UserID" `
}

func (User) TableName() string {
	return "sys_user"
}
