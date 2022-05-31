package table

import "pluto/global"

type Robot struct {
	global.UUID
	Name         string
	SoftID       string
	URL          string
	WorkStatus   string // free\busy
	OnlineStatus string // online\offline
	SystemStatus string // 机器人在内部的状态
}

func (Robot) TableName() string {
	return "t_robot"
}
