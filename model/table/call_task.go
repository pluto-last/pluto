package table

import (
	"pluto/global"
	"time"
)

// CallTask 外呼任务表
type CallTask struct {
	global.UUID
	UserID           string    `gorm:"index" json:"userID"`
	Name             string    `json:"name"`             // 任务名称
	SceneID          string    `json:"sceneID"  `        // 话术ID
	Status           string    `json:"status"`           // 状态
	ExecuteStartTime int       `json:"executeStartTime"` // 任务允许拨打的开始时间（存分钟数）
	ExecuteEndTime   int       `json:"executeEndTime"`   // 任务允许拨打的结束时间（存分钟数）
	Describe         string    `json:"describe"`         // 任务描述
	Sips             []TaskSip `json:"sips" gorm:"ForeignKey:ID;AssociationForeignKey:TaskID"`
}

func (CallTask) TableName() string {
	return "m_call_task"
}

// CallPlan 呼叫计划表
type CallPlan struct {
	global.UUID
	UserID       string    `gorm:"index" json:"userID"`
	TaskID       string    `json:"taskID"`    // 任务ID
	SceneID      string    `json:"sceneID"  ` // 话术ID
	SipID        string    `json:"sipID"`     // 线路ID
	Name         string    `json:"name"`
	Mobile       string    `json:"mobile"`
	Status       string    `json:"status"`       // 任务状态
	CallStatus   string    `json:"callStatus"`   // 呼叫状态
	DurationSec  int64     `json:"durationSec"`  // 时长(单位秒)
	CallAt       time.Time `json:"callAt"`       // 拨打时间
	IntentionTag string    `json:"intentionTag"` // 最终意向
	Note         string    `json:"note"`         // 通话备注
	Rounds       int       `json:"rounds"`       // 通话回合数
	Price        int       `json:"price"`        // 单价 每分钟
	Minute       int       `json:"minute"`       // 计费分钟数
	Cost         int       `json:"cost"`         // 费用
	IsRead       bool      `json:"isRead"`       // 是否已读
}

func (CallPlan) TableName() string {
	return "m_call_plan"
}

// SQLFields 实现xgorm批量插入的方法
func (plan CallPlan) SQLFields() (string, []string, []interface{}) {
	now := time.Now()
	if plan.ID == "" {
		plan.ID = global.RandUUID()
	}
	return "m_call_plan",
		[]string{
			"id", "user_id", "task_id", "name",
			"mobile", "status", "created_at", "updated_at",
		},
		[]interface{}{
			plan.ID, plan.UserID, plan.TaskID, plan.Name, plan.Mobile,
			plan.Status, now, now,
		}
}

// TaskSip 任务线路关联表
type TaskSip struct {
	global.UUID
	UserID      string  `gorm:"index" json:"userID"`
	TaskID      string  `json:"taskID"`     // 任务ID
	UserSipID   string  `json:"userSipID"`  // 用户线路ID
	Concurrent  int     `json:"concurrent"` // 并发数
	UserSipInfo UserSip `gorm:"ForeignKey:UserSipID;AssociationForeignKey:ID" json:"userSipInfo" `
}

func (TaskSip) TableName() string {
	return "m_task_sip"
}
