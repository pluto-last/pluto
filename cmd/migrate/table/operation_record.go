package table

import (
	"pluto/global"
)

// OperationRecord 操作记录表
type OperationRecord struct {
	global.UUID
	UserID       string `json:"user_id" gorm:"index"` // 用户id
	Ip           string `json:"ip" `                  // 请求ip
	Method       string `json:"method" `              // 请求方法
	Path         string `json:"path" `                // 请求路径
	Status       int    `json:"status" `              // 请求状态
	Agent        string `json:"agent" `               // 代理
	ErrorMessage string `json:"error_message" `       // 错误信息
	Body         string `json:"body"`                 // 请求Body
	Resp         string `json:"resp" `                // 响应Body
}

func (OperationRecord) TableName() string {
	return "sys_operation_record"
}
