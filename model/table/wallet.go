package table

import (
	"pluto/global"
)

// Wallet 用户钱包表
type Wallet struct {
	global.UUID
	UserID  string `json:"userID"`
	Balance int    `json:"balance"` // 余额
}

func (Wallet) TableName() string {
	return "sys_wallet"
}

// Recharge 充值记录表
type Recharge struct {
	global.UUID
	UserID   string `json:"userID"`  // 充值的用户 id
	Type     string `json:"type"`    // 充值类型
	Amount   int    `json:"amount"`  // 充值金额
	Balance  int    `json:"balance"` // 余额
	Status   string `json:"status"`  // 充值状态
	UserInfo User   `gorm:"ForeignKey:UserID" json:"userInfo"`
}

func (Recharge) TableName() string {
	return "m_recharge"
}
