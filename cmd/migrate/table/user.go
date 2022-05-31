package table

import (
	"pluto/global"
)

type User struct {
	global.UUID
	UserName   string `json:"userName"`
	Mobile     string `json:"mobile" `
	Password   string `json:"-"  `
	RegisterIP string `gorm:"index"`
	HeaderImg  string `json:"headerImg" `
	Note       string `json:"note" `
	Status     string `json:"status" `
}

func (User) TableName() string {
	return "sys_user"
}
