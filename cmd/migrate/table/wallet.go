package table

import "pluto/global"

type Wallet struct {
	global.UUID
	UserID  string `json:"userID"`
	Balance int    `json:"balance"`
	Point   int    `json:"point"`
}

func (Wallet) TableName() string {
	return "sys_wallet"
}
